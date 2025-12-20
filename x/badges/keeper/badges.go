package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Create tokens and update the unminted / total supplys for the collection
func (k Keeper) CreateTokens(ctx sdk.Context, collection *types.TokenCollection, newValidTokenIds []*types.UintRange) (*types.TokenCollection, error) {
	var err error
	allTokenIds := []*types.UintRange{}
	allTokenIds = append(allTokenIds, newValidTokenIds...)
	allTokenIds, err = types.SortUintRangesAndMerge(allTokenIds, true)
	if err != nil {
		return &types.TokenCollection{}, err
	}

	// Ensure the token ids are sequential starting from 1
	if len(allTokenIds) > 1 || (len(allTokenIds) == 1 && !allTokenIds[0].Start.Equal(sdkmath.NewUint(1))) {
		return &types.TokenCollection{}, sdkerrors.Wrapf(types.ErrNotSupported, "Ids must be sequential starting from 1")
	}

	if len(newValidTokenIds) == 0 {
		return collection, nil
	}

	//Check if we are allowed to create these tokens
	detailsToCheck := []*types.UniversalPermissionDetails{}
	for _, tokenIdRange := range newValidTokenIds {
		detailsToCheck = append(detailsToCheck, &types.UniversalPermissionDetails{
			TokenId: tokenIdRange,
		})
	}

	err = k.CheckIfTokenIdsActionPermissionPermits(ctx, detailsToCheck, collection.CollectionPermissions.CanUpdateValidTokenIds, "can create more tokens")
	if err != nil {
		return &types.TokenCollection{}, err
	}

	collection.ValidTokenIds = allTokenIds

	return collection, nil
}
