package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkmath "cosmossdk.io/math"
)

// Determines what to validate for each Msg
type UniversalValidationParams struct {
	Creator                              string
	CollectionId                         sdkmath.Uint
	MustBeManager                        bool
	OverrideArchive 										 bool
}

// Validates everything about the Msg is valid and returns (creatorNum, collection, permissions, error).
func (k Keeper) UniversalValidate(ctx sdk.Context, params UniversalValidationParams) (*types.BadgeCollection, error) {
	// Assert collection and badgeId ranges exist and are well-formed
	collection, found := k.GetCollectionFromStore(ctx, params.CollectionId)
	if !found {
		return &types.BadgeCollection{}, ErrCollectionNotExists
	}

	if !params.OverrideArchive {
		isArchived := types.GetIsArchived(ctx, collection)
		if isArchived {
			return &types.BadgeCollection{}, ErrCollectionIsArchived
		}
	}

	// Assert all permissions
	if params.MustBeManager {
		currManager := types.GetCurrentManager(ctx, collection)
		if currManager != params.Creator {
			return &types.BadgeCollection{}, ErrSenderIsNotManager
		}
	}

	return collection, nil
}