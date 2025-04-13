package keeper

import (
	"bitbadgeschain/x/badges/types"

	sdkerrors "cosmossdk.io/errors"
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
