package keeper

import (
	sdkerrors "cosmossdk.io/errors"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Determines what to validate for each Msg
type UniversalValidationParams struct {
	Creator       string
	MustBeManager bool
}

// Validates everything about the Msg is valid and returns (creatorNum, collection, permissions, error).
func (k Keeper) UniversalValidate(ctx sdk.Context, collection *types.BadgeCollection, params UniversalValidationParams) error {
	// Assert all permissions
	if params.MustBeManager {
		currManager := types.GetCurrentManager(ctx, collection)
		if currManager != params.Creator {
			return sdkerrors.Wrapf(ErrSenderIsNotManager, "current manager is %s but got %s", currManager, params.Creator)
		}
	}

	return nil
}


func (k Keeper) UniversalValidateNotHalted(ctx sdk.Context) error {
	halted := k.GetGlobalArchiveFromStore(ctx)

	if halted {
		return sdkerrors.Wrap(ErrGlobalArchive, "this action is not executable while the chain is halted")
	}

	return nil
}