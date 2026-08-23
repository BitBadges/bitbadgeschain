package app

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/stretchr/testify/require"

	appparams "github.com/bitbadges/bitbadgeschain/app/params"
	v35 "github.com/bitbadges/bitbadgeschain/app/upgrades/v35"
)

const legacyDenom = "ubadge"

func randAddr() sdk.AccAddress {
	return sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
}

// fundLegacy mints legacy coins and places them on addr, keeping supply and
// balances consistent so the migration's own accounting check is meaningful.
func fundLegacy(t *testing.T, ctx sdk.Context, bk bankkeeper.BaseKeeper, addr sdk.AccAddress, amount int64) {
	t.Helper()
	coins := sdk.NewCoins(sdk.NewCoin(legacyDenom, sdkmath.NewInt(amount)))
	require.NoError(t, bk.MintCoins(ctx, "mint", coins))
	require.NoError(t, bk.SendCoinsFromModuleToAccount(ctx, "mint", addr, coins))
}

func TestV35RedenominateBank(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	bk := app.BankKeeper.(bankkeeper.BaseKeeper)

	// Wipe whatever genesis put in place so the arithmetic below is exact.
	startingSupply := bk.GetSupply(ctx, legacyDenom).Amount
	require.True(t, startingSupply.IsZero(), "expected a clean legacy supply at genesis, got %s", startingSupply)

	alice, bob, carol := randAddr(), randAddr(), randAddr()
	fundLegacy(t, ctx, bk, alice, 1_000_000_000) // 1 BADGE at 9 decimals
	fundLegacy(t, ctx, bk, bob, 1)               // smallest possible unit
	fundLegacy(t, ctx, bk, carol, 123_456_789)

	legacyTotal := sdkmath.NewInt(1_000_000_000 + 1 + 123_456_789)
	require.Equal(t, legacyTotal, bk.GetSupply(ctx, legacyDenom).Amount)

	res, err := v35.RedenominateBank(ctx, bk)
	require.NoError(t, err)
	require.Equal(t, 3, res.Holders)

	factor := v35.ConversionFactor
	require.Equal(t, sdkmath.NewInt(1_000_000_000), factor, "conversion factor must be 10^9")

	// Every holder converted at exactly 10^9, including the 1-unit dust case
	// that a naive truncating conversion would round away.
	require.Equal(t, sdkmath.NewInt(1_000_000_000).Mul(factor),
		bk.GetBalance(ctx, alice, appparams.BaseCoinUnit).Amount)
	require.Equal(t, factor, bk.GetBalance(ctx, bob, appparams.BaseCoinUnit).Amount,
		"1 ubadge must become exactly 10^9 abadge")
	require.Equal(t, sdkmath.NewInt(123_456_789).Mul(factor),
		bk.GetBalance(ctx, carol, appparams.BaseCoinUnit).Amount)

	// The legacy denom is fully retired — no balances, no supply.
	require.True(t, bk.GetBalance(ctx, alice, legacyDenom).IsZero())
	require.True(t, bk.GetBalance(ctx, bob, legacyDenom).IsZero())
	require.True(t, bk.GetBalance(ctx, carol, legacyDenom).IsZero())
	require.True(t, bk.GetSupply(ctx, legacyDenom).Amount.IsZero(),
		"legacy supply must be zero after migration")

	// Supply moved with the balances, and equals the sum of them. This is the
	// assertion that catches a migration that mints or burns value.
	newSupply := bk.GetSupply(ctx, appparams.BaseCoinUnit).Amount
	require.Equal(t, legacyTotal.Mul(factor), newSupply)

	summed := sdkmath.ZeroInt()
	bk.IterateAllBalances(ctx, func(_ sdk.AccAddress, coin sdk.Coin) bool {
		if coin.Denom == appparams.BaseCoinUnit {
			summed = summed.Add(coin.Amount)
		}
		return false
	})
	require.Equal(t, newSupply, summed, "balances must sum to supply — no value created or destroyed")
}

// The staging module account is borrowed to move supply. If it happens to hold
// a balance of its own, that balance must survive at the converted amount and
// must not absorb the minted total.
func TestV35RedenominateBank_StagingModuleAccountBalanceIsPreserved(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	bk := app.BankKeeper.(bankkeeper.BaseKeeper)

	// Instantiate the module account properly; funding it via
	// SendCoinsFromModuleToAccount would write a plain BaseAccount over the
	// module account record and break BurnCoins.
	moduleAcc := app.AccountKeeper.GetModuleAccount(ctx, "evm")
	require.NotNil(t, moduleAcc)
	moduleAddr := moduleAcc.GetAddress()

	alice := randAddr()
	fundLegacy(t, ctx, bk, alice, 500)

	staged := sdk.NewCoins(sdk.NewCoin(legacyDenom, sdkmath.NewInt(250)))
	require.NoError(t, bk.MintCoins(ctx, "mint", staged))
	require.NoError(t, bk.SendCoinsFromModuleToModule(ctx, "mint", "evm", staged))

	res, err := v35.RedenominateBank(ctx, bk)
	require.NoError(t, err)
	require.Equal(t, 2, res.Holders)

	factor := v35.ConversionFactor
	require.Equal(t, sdkmath.NewInt(500).Mul(factor), bk.GetBalance(ctx, alice, appparams.BaseCoinUnit).Amount)
	require.Equal(t, sdkmath.NewInt(250).Mul(factor), bk.GetBalance(ctx, moduleAddr, appparams.BaseCoinUnit).Amount,
		"the staging account's own balance must convert, not be zeroed or inflated")

	require.Equal(t, sdkmath.NewInt(750).Mul(factor), bk.GetSupply(ctx, appparams.BaseCoinUnit).Amount)
}

// A supply that does not match the sum of balances means something is
// unaccounted for. Converting anyway would invent or destroy value, so the
// migration must refuse rather than guess.
func TestV35RedenominateBank_RefusesWhenSupplyDisagrees(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	bk := app.BankKeeper.(bankkeeper.BaseKeeper)

	alice := randAddr()
	fundLegacy(t, ctx, bk, alice, 100)

	// Inflate one balance without touching supply.
	require.NoError(t, bk.UncheckedSetBalance(ctx, alice, sdk.NewCoin(legacyDenom, sdkmath.NewInt(999))))

	_, err := v35.RedenominateBank(ctx, bk)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unaccounted supply")

	// And it must not have half-migrated anything.
	require.True(t, bk.GetSupply(ctx, appparams.BaseCoinUnit).Amount.IsZero(),
		"no new denom should exist after a refused migration")
}

// Nothing to do on a chain with no legacy balances; must not error or mint.
func TestV35RedenominateBank_NoLegacySupplyIsNoop(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	bk := app.BankKeeper.(bankkeeper.BaseKeeper)

	res, err := v35.RedenominateBank(ctx, bk)
	require.NoError(t, err)
	require.Equal(t, 0, res.Holders)
	require.True(t, bk.GetSupply(ctx, appparams.BaseCoinUnit).Amount.IsZero())
}

// Non-BADGE denominations must not move. An IBC voucher is worth what it is
// worth regardless of what the native token's decimals are.
func TestV35RedenominateBank_ForeignDenomsUntouched(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	bk := app.BankKeeper.(bankkeeper.BaseKeeper)

	alice := randAddr()
	fundLegacy(t, ctx, bk, alice, 10)

	const ibcDenom = "ibc/1234567890ABCDEF"
	voucher := sdk.NewCoins(sdk.NewCoin(ibcDenom, sdkmath.NewInt(777)))
	require.NoError(t, bk.MintCoins(ctx, "mint", voucher))
	require.NoError(t, bk.SendCoinsFromModuleToAccount(ctx, "mint", alice, voucher))

	_, err := v35.RedenominateBank(ctx, bk)
	require.NoError(t, err)

	require.Equal(t, sdkmath.NewInt(777), bk.GetBalance(ctx, alice, ibcDenom).Amount,
		"IBC voucher balance must be untouched")
	require.Equal(t, sdkmath.NewInt(777), bk.GetSupply(ctx, ibcDenom).Amount,
		"IBC voucher supply must be untouched")
}
