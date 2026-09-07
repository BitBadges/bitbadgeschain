package v35

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	precisebanktypes "github.com/cosmos/evm/contrib/x/precisebank/types"
)

// MigrateV35PreciseBankModulePermissions grants the precisebank and EVM modules
// account the minter and burner permissions its reserve needs.
//
// precisebank keeps every fractional balance backed by whole integer coins
// held in its own module account, and tops that reserve up (or trims it) by
// minting or burning one integer coin at a time. x/bank checks those
// permissions on the stored module account object, not on the app's
// configuration, so an account created before the permissions were configured
// keeps its empty list until it is rewritten. Idempotent.
func MigrateV35PreciseBankModulePermissions(ctx sdk.Context, ak authkeeper.AccountKeeper) error {
	for _, moduleName := range []string{precisebanktypes.ModuleName, "evm"} {
		acc := ak.GetModuleAccount(ctx, moduleName)
		if acc == nil {
			return fmt.Errorf("module account %s does not exist", moduleName)
		}
		if acc.HasPermission(authtypes.Minter) && acc.HasPermission(authtypes.Burner) {
			continue
		}

		base := authtypes.NewBaseAccount(acc.GetAddress(), nil, acc.GetAccountNumber(), acc.GetSequence())
		permissions := append([]string(nil), acc.GetPermissions()...)
		for _, permission := range []string{authtypes.Minter, authtypes.Burner} {
			if !acc.HasPermission(permission) {
				permissions = append(permissions, permission)
			}
		}
		ak.SetModuleAccount(ctx, authtypes.NewModuleAccount(base, moduleName, permissions...))
	}
	return nil
}
