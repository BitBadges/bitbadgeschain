package app

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/stretchr/testify/require"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	appparams "github.com/bitbadges/bitbadgeschain/app/params"
	v35 "github.com/bitbadges/bitbadgeschain/app/upgrades/v35"
)

// v35Keepers assembles the keeper set the upgrade handler takes, so the test
// exercises the same entry point the chain does rather than calling the
// individual migrations.
func v35Keepers(app *App) v35.Keepers {
	return v35.Keepers{
		Bank:         app.BankKeeper.(bankkeeper.BaseKeeper),
		Staking:      app.StakingKeeper,
		Mint:         app.MintKeeper,
		Gov:          app.GovKeeper,
		Distribution: app.DistrKeeper,
		EVM:          app.EVMKeeper,
		Gamm:         app.GammKeeper,
		PoolManager:  app.PoolManagerKeeper,
		Tokenization: *app.TokenizationKeeper,
	}
}

// seedLegacyChainState makes the test app look like the real chain before the
// upgrade: bond denom set to the legacy denom, and a validator with tokens and
// shares in the staking store. Setup() leaves the SDK default "stake" and no
// bonded validator, so without this the migration correctly finds nothing to do
// and the test would assert on a no-op.
func seedLegacyChainState(t *testing.T, app *App, ctx sdk.Context) stakingtypes.Validator {
	t.Helper()

	params, err := app.StakingKeeper.GetParams(ctx)
	require.NoError(t, err)
	params.BondDenom = legacyDenom
	require.NoError(t, app.StakingKeeper.SetParams(ctx, params))

	mintParams, err := app.MintKeeper.Params.Get(ctx)
	require.NoError(t, err)
	mintParams.MintDenom = legacyDenom
	require.NoError(t, app.MintKeeper.Params.Set(ctx, mintParams))

	pk := ed25519.GenPrivKey().PubKey()
	pkAny, err := codectypes.NewAnyWithValue(pk)
	require.NoError(t, err)

	val := stakingtypes.Validator{
		OperatorAddress:   sdk.ValAddress(pk.Address()).String(),
		ConsensusPubkey:   pkAny,
		Status:            stakingtypes.Bonded,
		Tokens:            sdkmath.NewInt(1_000_000_000),
		DelegatorShares:   sdkmath.LegacyNewDec(1_000_000_000),
		MinSelfDelegation: sdkmath.NewInt(1),
		Commission:        stakingtypes.NewCommission(sdkmath.LegacyZeroDec(), sdkmath.LegacyOneDec(), sdkmath.LegacyZeroDec()),
	}
	require.NoError(t, app.StakingKeeper.SetValidator(ctx, val))
	return val
}

// TestV35UpgradeConservesValue is the assertion the whole migration exists to
// satisfy: after redenominating, every holder owns exactly 10^9 times what they
// owned before, and the supply agrees with the sum of balances.
//
// Running the real handler rather than the individual migrations is deliberate.
// The ordering between steps is load-bearing — params are moved last so nothing
// upstream reads a half-updated denom — and calling the pieces directly would
// not exercise it.
func TestV35UpgradeConservesValue(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	bk := app.BankKeeper.(bankkeeper.BaseKeeper)

	alice, bob := randAddr(), randAddr()
	fundLegacy(t, ctx, bk, alice, 1_000_000_000) // 1 BADGE at 9 decimals
	fundLegacy(t, ctx, bk, bob, 1)               // smallest possible unit

	before := map[string]sdkmath.Int{
		alice.String(): bk.GetBalance(ctx, alice, legacyDenom).Amount,
		bob.String():   bk.GetBalance(ctx, bob, legacyDenom).Amount,
	}
	legacySupplyBefore := bk.GetSupply(ctx, legacyDenom).Amount

	require.NoError(t, v35.CustomUpgradeHandlerLogic(ctx, v35Keepers(app)))

	factor := v35.ConversionFactor

	for addrStr, amount := range before {
		addr, err := sdk.AccAddressFromBech32(addrStr)
		require.NoError(t, err)
		require.Equal(t, amount.Mul(factor), bk.GetBalance(ctx, addr, appparams.BaseCoinUnit).Amount,
			"holder %s must own exactly 10^9x what it owned before", addrStr)
		require.True(t, bk.GetBalance(ctx, addr, legacyDenom).IsZero(),
			"holder %s must hold none of the retired denom", addrStr)
	}

	require.True(t, bk.GetSupply(ctx, legacyDenom).Amount.IsZero(), "legacy supply must be retired")
	newSupply := bk.GetSupply(ctx, appparams.BaseCoinUnit).Amount
	require.Equal(t, legacySupplyBefore.Mul(factor), newSupply, "supply must scale by exactly 10^9")

	// The invariant that catches a migration which quietly invents or destroys
	// value: every balance in the new denom must add up to the supply.
	summed := sdkmath.ZeroInt()
	bk.IterateAllBalances(ctx, func(_ sdk.AccAddress, coin sdk.Coin) bool {
		if coin.Denom == appparams.BaseCoinUnit {
			summed = summed.Add(coin.Amount)
		}
		return false
	})
	require.Equal(t, newSupply, summed, "balances must sum to supply after the upgrade")
}

// The staking module keeps its own token accounting that no bank balance
// reflects. If it is not rescaled alongside the pools, the bonded pool holds
// 10^9x what the validators believe they have.
func TestV35UpgradeRescalesStakingWithBankPools(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)

	seeded := seedLegacyChainState(t, app, ctx)
	tokensBefore := seeded.Tokens
	sharesBefore := seeded.DelegatorShares
	require.True(t, tokensBefore.IsPositive())

	require.NoError(t, v35.CustomUpgradeHandlerLogic(ctx, v35Keepers(app)))

	valsAfter, err := app.StakingKeeper.GetAllValidators(ctx)
	require.NoError(t, err)
	require.Len(t, valsAfter, 1)

	factor := v35.ConversionFactor
	require.Equal(t, tokensBefore.Mul(factor), valsAfter[0].Tokens,
		"validator tokens must scale with the bonded pool")
	require.Equal(t, sharesBefore.Mul(sdkmath.LegacyNewDecFromInt(factor)), valsAfter[0].DelegatorShares,
		"delegator shares must scale with tokens so the exchange rate is unchanged")

	// The property that actually matters: scaling both leaves each share worth
	// the same, so no delegation is repriced.
	rateBefore := sharesBefore.Quo(sdkmath.LegacyNewDecFromInt(tokensBefore))
	rateAfter := valsAfter[0].DelegatorShares.Quo(sdkmath.LegacyNewDecFromInt(valsAfter[0].Tokens))
	require.Equal(t, rateBefore, rateAfter, "the shares-per-token exchange rate must not move")
}

// Params are migrated last, but they must actually end up migrated — a staking
// bond denom left pointing at the retired denom would make every delegation
// fail against a denom with zero supply.
func TestV35UpgradeRepointsParams(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	seedLegacyChainState(t, app, ctx)

	require.NoError(t, v35.CustomUpgradeHandlerLogic(ctx, v35Keepers(app)))

	stakingParams, err := app.StakingKeeper.GetParams(ctx)
	require.NoError(t, err)
	require.Equal(t, appparams.BaseCoinUnit, stakingParams.BondDenom)

	mintParams, err := app.MintKeeper.Params.Get(ctx)
	require.NoError(t, err)
	require.Equal(t, appparams.BaseCoinUnit, mintParams.MintDenom)

	evmParams := app.EVMKeeper.GetParams(ctx)
	require.Equal(t, appparams.BaseCoinUnit, evmParams.EvmDenom)
	require.NotNil(t, evmParams.ExtendedDenomOptions)
	require.Equal(t, appparams.BaseCoinUnit, evmParams.ExtendedDenomOptions.ExtendedDenom,
		"at 18 decimals x/vm requires base == extended denom")

	tokParams := app.TokenizationKeeper.GetParams(ctx)
	require.NotContains(t, tokParams.AllowedDenoms, legacyDenom,
		"the retired denom must not remain allowed")
	require.Contains(t, tokParams.AllowedDenoms, appparams.BaseCoinUnit)
}

// Running the handler twice must not double-scale. Upgrade handlers get re-run
// during replay and recovery, and a migration that is not idempotent turns a
// routine restart into a 10^18 inflation event.
func TestV35UpgradeIsIdempotent(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	bk := app.BankKeeper.(bankkeeper.BaseKeeper)

	alice := randAddr()
	fundLegacy(t, ctx, bk, alice, 12_345)

	require.NoError(t, v35.CustomUpgradeHandlerLogic(ctx, v35Keepers(app)))
	afterFirst := bk.GetBalance(ctx, alice, appparams.BaseCoinUnit).Amount
	supplyAfterFirst := bk.GetSupply(ctx, appparams.BaseCoinUnit).Amount

	require.NoError(t, v35.CustomUpgradeHandlerLogic(ctx, v35Keepers(app)))

	require.Equal(t, afterFirst, bk.GetBalance(ctx, alice, appparams.BaseCoinUnit).Amount,
		"a second run must not scale balances again")
	require.Equal(t, supplyAfterFirst, bk.GetSupply(ctx, appparams.BaseCoinUnit).Amount,
		"a second run must not scale supply again")
}
