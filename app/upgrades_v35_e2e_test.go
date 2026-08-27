package app

import (
	"testing"

	"cosmossdk.io/collections"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/stretchr/testify/require"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	balancer "github.com/bitbadges/bitbadgeschain/x/gamm/poolmodels/balancer"
	ratelimittypes "github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/types"
	tokenizationkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	appparams "github.com/bitbadges/bitbadgeschain/app/params"
	v35 "github.com/bitbadges/bitbadgeschain/app/upgrades/v35"
)

// v35Keepers assembles the keeper set the upgrade handler takes, so the test
// exercises the same entry point the chain does rather than calling the
// individual migrations.
// v35Keepers delegates to the production wiring on purpose. Building an
// equivalent struct here is what previously let the tests pass while the chain
// ran the migration with nil keepers.
func v35Keepers(app *App) v35.Keepers {
	return app.V35Keepers()
}

// seedLegacyChainState makes the test app look like the real chain before the
// upgrade: bond denom set to the legacy denom, and a validator with tokens and
// shares in the staking store. Setup() leaves the SDK default "stake" and no
// bonded validator, so without this the migration correctly finds nothing to do
// and the test would assert on a no-op.
// seedLegacyBondDenom puts the staking module on the denom this upgrade
// redenominates.
//
// Every test that runs the handler needs this, not just the ones that assert on
// staking. Setup() leaves the SDK default "stake", and the handler now refuses
// outright to run on a chain whose bond denom it does not redenominate, because
// PowerReduction is compiled in at 10^15 and consensus power is
// tokens/PowerReduction. A test app bonding "stake" at 10^15 is exactly the
// configuration the guard exists to reject, so seeding it here is not a
// workaround — it is the test finally describing a chain this upgrade can run
// on.
func seedLegacyBondDenom(t *testing.T, app *App, ctx sdk.Context) {
	t.Helper()

	params, err := app.StakingKeeper.GetParams(ctx)
	require.NoError(t, err)
	params.BondDenom = legacyDenom
	require.NoError(t, app.StakingKeeper.SetParams(ctx, params))
}

func seedLegacyChainState(t *testing.T, app *App, ctx sdk.Context) stakingtypes.Validator {
	t.Helper()

	seedLegacyBondDenom(t, app, ctx)

	mintParams, err := app.MintKeeper.Params.Get(ctx)
	require.NoError(t, err)
	mintParams.MintDenom = legacyDenom
	mintParams.MaxSupply = sdkmath.NewInt(10_000_000_000_000)
	require.NoError(t, app.MintKeeper.Params.Set(ctx, mintParams))

	minter, err := app.MintKeeper.Minter.Get(ctx)
	require.NoError(t, err)
	minter.AnnualProvisions = sdkmath.LegacyNewDec(500_000_000)
	require.NoError(t, app.MintKeeper.Minter.Set(ctx, minter))

	// The EVM's gas token is the legacy denom before the upgrade, and the
	// feemarket prices are denominated in it. Both are what the real chain
	// looks like at the upgrade height.
	evmParams := app.EVMKeeper.GetParams(ctx)
	evmParams.EvmDenom = legacyDenom
	require.NoError(t, app.EVMKeeper.SetParams(ctx, evmParams))

	feeParams := app.FeeMarketKeeper.GetParams(ctx)
	feeParams.BaseFee = sdkmath.LegacyNewDec(1_000)
	feeParams.MinGasPrice = sdkmath.LegacyNewDec(10)
	require.NoError(t, app.FeeMarketKeeper.SetParams(ctx, feeParams))

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

	seedLegacyBondDenom(t, app, ctx)

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
//
// Seeding the full legacy chain state rather than a bank balance alone is the
// point of this test. RedenominateBank early-returns once the legacy supply is
// zero, so a bank-only seed exercises the one path that is idempotent by
// construction and proves nothing about the unconditional Mul sites in
// staking, distribution, mint and feemarket.
func TestV35UpgradeIsIdempotent(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	bk := app.BankKeeper.(bankkeeper.BaseKeeper)

	// A validator, delegations with pending rewards, and the mint/feemarket
	// figures. seedDelegationsWithPendingRewards returns a context one block
	// later, which is what makes the reward calculation non-zero.
	ctx, seededVal, delegators := seedDelegationsWithPendingRewards(t, app, ctx, 1_000, 2)
	valAddr, err := sdk.ValAddressFromBech32(seededVal.OperatorAddress)
	require.NoError(t, err)

	alice := randAddr()
	fundLegacy(t, ctx, bk, alice, 12_345)

	// A gamm pool, so the pool record's own copy of its reserves is in the
	// snapshot rather than only the bank balances behind it.
	poolID := seedIdempotencyGammPool(t, app, ctx)

	// A gov proposal with a paid deposit: the module account holds the coins but
	// gov keeps its own copy in TotalDeposit and the Deposit record.
	proposalID, depositor := seedIdempotencyGovDeposit(t, app, ctx, bk, 5_000)

	// A backed collection, whose escrow balance moves to a re-derived address.
	escrowTotal := int64(4_000_000_000)
	oldEscrow, newEscrow := seedIdempotencyBackedCollection(t, app, ctx, bk, escrowTotal)

	// x/precisebank fractional balances, paid out and then gone.
	seedPreciseBankFractionalBalances(t, app, ctx, bk, map[string]int64{
		alice.String(): 1_000_000_000,
	}, 0)

	// Rate limit caps and an accrued taker fee: both are amounts that live
	// outside the bank entirely.
	seedIdempotencyRateLimit(t, app, ctx)
	require.NoError(t, app.PoolManagerKeeper.UpdateTakerFeeTrackerForStakersByDenom(
		ctx, legacyDenom, osmomath.NewInt(7_000)))

	snapshot := func() map[string]string {
		vals, err := app.StakingKeeper.GetAllValidators(ctx)
		require.NoError(t, err)
		require.Len(t, vals, 1)

		minter, err := app.MintKeeper.Minter.Get(ctx)
		require.NoError(t, err)
		mintParams, err := app.MintKeeper.Params.Get(ctx)
		require.NoError(t, err)
		feeParams := app.FeeMarketKeeper.GetParams(ctx)

		// The stake recorded against the delegator, which is what the reward
		// calculation multiplies the ratio delta by. Scaling it twice is the
		// exact shape of the bug this test is looking for.
		startingInfo, err := app.DistrKeeper.GetDelegatorStartingInfo(ctx, valAddr, delegators[0])
		require.NoError(t, err)

		pool, err := app.GammKeeper.GetPoolAndPoke(ctx, poolID)
		require.NoError(t, err)

		proposal, err := app.GovKeeper.Proposals.Get(ctx, proposalID)
		require.NoError(t, err)
		deposit, err := app.GovKeeper.Deposits.Get(ctx, collections.Join(proposalID, depositor))
		require.NoError(t, err)

		takerFee, err := app.PoolManagerKeeper.GetTakerFeeTrackerForStakersByDenom(ctx, appparams.BaseCoinUnit)
		require.NoError(t, err)

		rateLimits := app.IBCRateLimitKeeper.GetParams(ctx).RateLimits
		require.Len(t, rateLimits, 1)

		return map[string]string{
			"balance":            bk.GetBalance(ctx, alice, appparams.BaseCoinUnit).Amount.String(),
			"supply":             bk.GetSupply(ctx, appparams.BaseCoinUnit).Amount.String(),
			"validator_tokens":   vals[0].Tokens.String(),
			"validator_shares":   vals[0].DelegatorShares.String(),
			"min_self_delegate":  vals[0].MinSelfDelegation.String(),
			"annual_provisions":  minter.AnnualProvisions.String(),
			"max_supply":         mintParams.MaxSupply.String(),
			"base_fee":           feeParams.BaseFee.String(),
			"min_gas_price":      feeParams.MinGasPrice.String(),
			"starting_stake":     startingInfo.Stake.String(),
			"pending_reward":     pendingReward(t, app, ctx, valAddr, delegators[0]).String(),
			"pool_liquidity":     pool.GetTotalPoolLiquidity(ctx).String(),
			"pool_shares":        pool.GetTotalShares().String(),
			"gov_total_deposit":  sdk.Coins(proposal.TotalDeposit).String(),
			"gov_deposit":        sdk.Coins(deposit.Amount).String(),
			"escrow_old":         bk.GetBalance(ctx, oldEscrow, appparams.BaseCoinUnit).Amount.String(),
			"escrow_new":         bk.GetBalance(ctx, newEscrow, appparams.BaseCoinUnit).Amount.String(),
			"precisebank_paid":   bk.GetBalance(ctx, alice, appparams.BaseCoinUnit).Amount.String(),
			"taker_fee_stakers":  takerFee.Amount.String(),
			"rate_limit_cap":     rateLimits[0].SupplyShiftLimits[0].MaxAmount.String(),
			"rate_limit_addr":    rateLimits[0].AddressLimits[0].MaxAmount.String(),
			"gamm_liquidity_idx": app.GammKeeper.GetDenomLiquidity(ctx, appparams.BaseCoinUnit).String(),
		}
	}

	require.NoError(t, v35.CustomUpgradeHandlerLogic(ctx, v35Keepers(app)))
	afterFirst := snapshot()

	// Every seeded quantity must actually have moved on the first run, or the
	// idempotency assertion below would hold vacuously. This is the half of the
	// test that decays silently as state is added, so each new key gets a
	// liveness assertion alongside it.
	factor := v35.ConversionFactor
	require.Equal(t, seededVal.Tokens.Mul(factor).String(), afterFirst["validator_tokens"])
	require.NotEqual(t, "0.000000000000000000", afterFirst["annual_provisions"])
	require.NotEqual(t, "0", afterFirst["max_supply"])
	require.NotEqual(t, "0.000000000000000000", afterFirst["base_fee"])
	require.Equal(t,
		sdkmath.LegacyNewDecFromInt(seededVal.Tokens.Quo(sdkmath.NewInt(2)).Mul(factor)).String(),
		afterFirst["starting_stake"],
		"the delegator's recorded stake must have scaled exactly once")
	require.Contains(t, afterFirst["pool_liquidity"], appparams.BaseCoinUnit)
	require.Contains(t, afterFirst["gov_total_deposit"], appparams.BaseCoinUnit)
	require.Equal(t, sdkmath.NewInt(5_000).Mul(factor).String()+appparams.BaseCoinUnit, afterFirst["gov_deposit"])
	require.Equal(t, "0", afterFirst["escrow_old"], "the retired escrow must be empty")
	require.Equal(t, sdkmath.NewInt(escrowTotal).Mul(factor).String(), afterFirst["escrow_new"],
		"and the re-derived one must hold everything")
	require.Equal(t, sdkmath.NewInt(7_000).Mul(factor).String(), afterFirst["taker_fee_stakers"])
	require.Equal(t, sdkmath.NewInt(1_000_000).Mul(factor).String(), afterFirst["rate_limit_cap"])
	require.Equal(t, sdkmath.NewInt(500_000).Mul(factor).String(), afterFirst["rate_limit_addr"])
	require.Equal(t, sdkmath.NewInt(2_500).Mul(factor).String(), afterFirst["gamm_liquidity_idx"])

	require.NoError(t, v35.CustomUpgradeHandlerLogic(ctx, v35Keepers(app)))
	afterSecond := snapshot()

	require.Equal(t, afterFirst, afterSecond,
		"a second run must not scale anything again")
}

// seedIdempotencyGammPool creates a balancer pool holding the legacy denom, so
// the pool record's own copy of its reserves is exercised.
func seedIdempotencyGammPool(t *testing.T, app *App, ctx sdk.Context) uint64 {
	t.Helper()

	const foreignDenom = "ibc/ABCDEF0123456789"
	creator := randAddr()
	assets := sdk.NewCoins(
		sdk.NewCoin(legacyDenom, sdkmath.NewInt(2_500)),
		sdk.NewCoin(foreignDenom, sdkmath.NewInt(2_500)),
	)
	require.NoError(t, app.BankKeeper.(bankkeeper.BaseKeeper).MintCoins(ctx, "mint", assets))
	require.NoError(t, app.BankKeeper.(bankkeeper.BaseKeeper).SendCoinsFromModuleToAccount(ctx, "mint", creator, assets))

	msg := balancer.NewMsgCreateBalancerPool(creator,
		balancer.PoolParams{
			SwapFee: osmomath.MustNewDecFromStr("0.003"),
			ExitFee: osmomath.ZeroDec(),
		},
		[]balancer.PoolAsset{
			{Weight: sdkmath.NewInt(1), Token: sdk.NewCoin(legacyDenom, sdkmath.NewInt(2_500))},
			{Weight: sdkmath.NewInt(1), Token: sdk.NewCoin(foreignDenom, sdkmath.NewInt(2_500))},
		})

	poolID, err := app.PoolManagerKeeper.CreatePool(ctx, msg)
	require.NoError(t, err)
	return poolID
}

// seedIdempotencyGovDeposit files a proposal and pays a deposit on it, so gov's
// own copy of the amount is exercised alongside the coins on the module account.
func seedIdempotencyGovDeposit(
	t *testing.T, app *App, ctx sdk.Context, bk bankkeeper.BaseKeeper, amount int64,
) (uint64, sdk.AccAddress) {
	t.Helper()

	// Gov only accepts deposits in the denoms its MinDeposit names, which is
	// "stake" out of Setup(). The real chain's is the legacy denom.
	govParams, err := app.GovKeeper.Params.Get(ctx)
	require.NoError(t, err)
	govParams.MinDeposit = sdk.NewCoins(sdk.NewCoin(legacyDenom, sdkmath.NewInt(1)))
	require.NoError(t, app.GovKeeper.Params.Set(ctx, govParams))

	depositor := randAddr()
	deposit := sdk.NewCoins(sdk.NewCoin(legacyDenom, sdkmath.NewInt(amount)))
	require.NoError(t, bk.MintCoins(ctx, "mint", deposit))
	require.NoError(t, bk.SendCoinsFromModuleToAccount(ctx, "mint", depositor, deposit))

	proposal, err := app.GovKeeper.SubmitProposal(ctx, nil, "", "idempotency", "idempotency", depositor, false)
	require.NoError(t, err)

	_, err = app.GovKeeper.AddDeposit(ctx, proposal.Id, depositor, deposit)
	require.NoError(t, err)

	return proposal.Id, depositor
}

// seedIdempotencyBackedCollection puts a backed collection and its escrowed
// coins in place, and returns the addresses the escrow moves between.
func seedIdempotencyBackedCollection(
	t *testing.T, app *App, ctx sdk.Context, bk bankkeeper.BaseKeeper, escrowed int64,
) (sdk.AccAddress, sdk.AccAddress) {
	t.Helper()

	oldEscrow, err := tokenizationkeeper.DerivePathAddress(legacyDenom, tokenizationkeeper.BackedPathGenerationPrefix)
	require.NoError(t, err)
	newEscrow, err := tokenizationkeeper.DerivePathAddress(appparams.BaseCoinUnit, tokenizationkeeper.BackedPathGenerationPrefix)
	require.NoError(t, err)

	require.NoError(t, app.TokenizationKeeper.SetCollectionInStore(ctx, &tokenizationtypes.TokenCollection{
		CollectionId: sdkmath.NewUint(1),
		Invariants: &tokenizationtypes.CollectionInvariants{
			CosmosCoinBackedPath: &tokenizationtypes.CosmosCoinBackedPath{
				Address: oldEscrow.String(),
				Conversion: &tokenizationtypes.Conversion{
					SideA: &tokenizationtypes.ConversionSideAWithDenom{
						Amount: sdkmath.NewUint(1_000_000_000),
						Denom:  legacyDenom,
					},
				},
			},
		},
	}, true))

	fundLegacy(t, ctx, bk, oldEscrow, escrowed)
	return oldEscrow, newEscrow
}

// seedIdempotencyRateLimit installs a config on the legacy denom whose caps are
// amounts and must scale exactly once.
func seedIdempotencyRateLimit(t *testing.T, app *App, ctx sdk.Context) {
	t.Helper()
	require.NoError(t, app.IBCRateLimitKeeper.SetParams(ctx, ratelimittypes.Params{
		RateLimits: []ratelimittypes.RateLimitConfig{{
			ChannelId: "channel-0",
			Denom:     legacyDenom,
			SupplyShiftLimits: []ratelimittypes.TimeframeLimit{{
				MaxAmount:         sdkmath.NewInt(1_000_000),
				TimeframeType:     ratelimittypes.TimeframeType_TIMEFRAME_TYPE_DAY,
				TimeframeDuration: 1,
			}},
			AddressLimits: []ratelimittypes.AddressLimit{{
				MaxTransfers:      10,
				MaxAmount:         sdkmath.NewInt(500_000),
				TimeframeType:     ratelimittypes.TimeframeType_TIMEFRAME_TYPE_HOUR,
				TimeframeDuration: 1,
			}},
		}},
	}))
}
