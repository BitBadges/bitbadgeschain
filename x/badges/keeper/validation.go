package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Determines what to validate for each Msg
type UniversalValidationParams struct {
	Creator       string
	MustBeManager bool
}

// Validates everything about the Msg is valid and returns (creatorNum, collection, permissions, error).
func (k Keeper) UniversalValidate(ctx sdk.Context, collection *types.TokenCollection, params UniversalValidationParams) error {
	// Assert all permissions
	if params.MustBeManager {
		// Check if gov authority
		govAuthorized := k.GetAuthority() == params.Creator
		if govAuthorized {
			// Gov authority is allowed to do anything
			// Technically, the gov authority is allowed to do anything already through manual upgrades, etc
			// This just streamlines it to be able to use execute Msgs directly from the proposal
			return nil
		}

		currManager := types.GetCurrentManager(ctx, collection)
		if currManager != params.Creator {
			return sdkerrors.Wrapf(ErrSenderIsNotManager, "current manager is %s but got %s", currManager, params.Creator)
		}
	}

	return nil
}
