package app

import (
	v35 "github.com/bitbadges/bitbadgeschain/app/upgrades/v35"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	precisebanktypes "github.com/cosmos/evm/contrib/x/precisebank/types"
	"github.com/stretchr/testify/require"
)

// TestMigrateV35PreciseBankModulePermissions starts from the shape the module
// account has on chain today — created with an empty permission list — and
// checks the migration lets precisebank grow its reserve.
func TestMigrateV35PreciseBankModulePermissions(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)

	existing := app.AccountKeeper.GetModuleAccount(ctx, precisebanktypes.ModuleName)
	require.NotNil(t, existing)
	base := authtypes.NewBaseAccount(existing.GetAddress(), nil, existing.GetAccountNumber(), existing.GetSequence())
	app.AccountKeeper.SetModuleAccount(ctx, authtypes.NewModuleAccount(base, precisebanktypes.ModuleName))
	require.False(t, app.AccountKeeper.GetModuleAccount(ctx, precisebanktypes.ModuleName).HasPermission(authtypes.Minter),
		"precondition: stored module account has no permissions")

	// Minting a fractional amount forces a one-coin reserve top-up.
	addr := sdk.AccAddress(newEVMAddress(t).Bytes())
	half := sdk.NewCoins(sdk.NewCoin(precisebanktypes.ExtendedCoinDenom(), sdkmath.NewIntFromBigInt(ubadgeToWeiFactor).QuoRaw(2)))
	require.Panics(t, func() {
		_ = app.PreciseBankKeeper.MintCoins(ctx, "evm", half)
	}, "precondition: reserve mint is refused before the migration")

	require.NoError(t, v35.MigrateV35PreciseBankModulePermissions(ctx, app.AccountKeeper))

	migrated := app.AccountKeeper.GetModuleAccount(ctx, precisebanktypes.ModuleName)
	require.Equal(t, existing.GetAccountNumber(), migrated.GetAccountNumber(), "account number must be preserved")
	require.Equal(t, existing.GetAddress(), migrated.GetAddress())
	require.True(t, migrated.HasPermission(authtypes.Minter))
	require.True(t, migrated.HasPermission(authtypes.Burner))

	require.NoError(t, app.PreciseBankKeeper.MintCoins(ctx, "evm", half))
	require.NoError(t, app.PreciseBankKeeper.SendCoinsFromModuleToAccount(ctx, "evm", addr, half))
	require.Equal(t, half.AmountOf(precisebanktypes.ExtendedCoinDenom()).String(),
		app.PreciseBankKeeper.GetBalance(ctx, addr, precisebanktypes.ExtendedCoinDenom()).Amount.String())

	// Running it again is a no-op.
	require.NoError(t, v35.MigrateV35PreciseBankModulePermissions(ctx, app.AccountKeeper))
}
