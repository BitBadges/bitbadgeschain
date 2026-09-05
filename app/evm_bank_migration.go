package app

import (
	"fmt"

	precisebanktypes "github.com/cosmos/evm/contrib/x/precisebank/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// MigrateV35PreciseBankModulePermissions grants the precisebank module
// account the minter and burner permissions its reserve needs.
//
// precisebank keeps every fractional balance backed by whole integer coins
// held in its own module account, and tops that reserve up (or trims it) by
// minting or burning one integer coin at a time. x/bank checks those
// permissions on the stored module account object, not on the app's
// configuration, so an account created before the permissions were configured
// keeps its empty list until it is rewritten. Idempotent.
func MigrateV35PreciseBankModulePermissions(ctx sdk.Context, ak authkeeper.AccountKeeper) error {
	acc := ak.GetModuleAccount(ctx, precisebanktypes.ModuleName)
	if acc == nil {
		return fmt.Errorf("module account %s does not exist", precisebanktypes.ModuleName)
	}
	if acc.HasPermission(authtypes.Minter) && acc.HasPermission(authtypes.Burner) {
		return nil
	}

	base := authtypes.NewBaseAccount(acc.GetAddress(), nil, acc.GetAccountNumber(), acc.GetSequence())
	ak.SetModuleAccount(ctx, authtypes.NewModuleAccount(base, precisebanktypes.ModuleName, authtypes.Minter, authtypes.Burner))
	return nil
}
