package app

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/stretchr/testify/require"

	appparams "github.com/bitbadges/bitbadgeschain/app/params"
)

// farFutureUnix keeps the vesting account locked for the whole test without
// reading the clock: 2100-01-01T00:00:00Z.
const farFutureUnix = int64(4102444800)

// TestEVMBankKeeperRespectsLockedCoins checks the x/vm bank keeper on an
// account with vesting-locked funds: the locked amount is reported in
// extended units, a balance write only moves spendable funds, and a write
// that would reach into the locked portion is refused.
func TestEVMBankKeeperRespectsLockedCoins(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	bk := newEVMBankKeeper(app.PreciseBankKeeper)

	integerDenom := appparams.BaseCoinUnit
	extDenom := evmtypes.GetEVMCoinExtendedDenom()
	factor := sdkmath.NewIntFromBigInt(ubadgeToWeiFactor)

	addr := sdk.AccAddress(newEVMAddress(t).Bytes())
	const total, locked = int64(1_000), int64(400)
	vesting := sdk.NewCoins(sdk.NewCoin(integerDenom, sdkmath.NewInt(locked)))
	base, ok := app.AccountKeeper.NewAccountWithAddress(ctx, addr).(*authtypes.BaseAccount)
	require.True(t, ok)
	acc, err := vestingtypes.NewDelayedVestingAccount(base, vesting, farFutureUnix)
	require.NoError(t, err)
	app.AccountKeeper.SetAccount(ctx, acc)

	funds := sdk.NewCoins(sdk.NewCoin(integerDenom, sdkmath.NewInt(total)))
	require.NoError(t, app.BankKeeper.MintCoins(ctx, "mint", funds))
	require.NoError(t, app.BankKeeper.SendCoinsFromModuleToAccount(ctx, "mint", addr, funds))
	require.Equal(t, sdkmath.NewInt(locked).String(),
		app.BankKeeper.LockedCoins(ctx, addr).AmountOf(integerDenom).String(), "precondition: funds are locked")

	// Locked coins come back scaled, under the EVM denom key x/vm reads.
	require.Equal(t, sdkmath.NewInt(locked).Mul(factor).String(),
		bk.LockedCoins(ctx, addr).AmountOf(integerDenom).String())

	// Spend 5 units: total 1000 -> 995, locked untouched.
	target := sdkmath.NewInt(total - 5).Mul(factor)
	require.NoError(t, bk.UncheckedSetBalance(ctx, addr, sdk.NewCoin(extDenom, target)))
	require.Equal(t, sdkmath.NewInt(total-5).String(), app.BankKeeper.GetBalance(ctx, addr, integerDenom).Amount.String())
	require.Equal(t, sdkmath.NewInt(locked).String(), app.BankKeeper.LockedCoins(ctx, addr).AmountOf(integerDenom).String())
	require.True(t, app.BankKeeper.GetAllBalances(ctx, addr).AmountOf(extDenom).IsZero())

	// A target below the locked amount would have to spend locked funds.
	require.Error(t, bk.UncheckedSetBalance(ctx, addr, sdk.NewCoin(extDenom, sdkmath.NewInt(locked-1).Mul(factor))))
	require.Equal(t, sdkmath.NewInt(total-5).String(), app.BankKeeper.GetBalance(ctx, addr, integerDenom).Amount.String(),
		"a refused write must not move anything")
}
