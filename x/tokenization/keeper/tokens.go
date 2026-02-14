package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

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

	// Only check permissions for NEW token IDs (not already in collection.ValidTokenIds)
	existingTokenIds := types.DeepCopyRanges(collection.ValidTokenIds)
	newTokenIdsCopy := types.DeepCopyRanges(newValidTokenIds)
	newTokenIdsOnly, _ := types.RemoveUintRangesFromUintRanges(existingTokenIds, newTokenIdsCopy)

	// Only check permissions if there are actually new token IDs
	if len(newTokenIdsOnly) > 0 {
		detailsToCheck := []*types.UniversalPermissionDetails{}
		for _, tokenIdRange := range newTokenIdsOnly {
			detailsToCheck = append(detailsToCheck, &types.UniversalPermissionDetails{
				TokenId: tokenIdRange,
			})
		}

		err = k.CheckIfTokenIdsActionPermissionPermits(ctx, detailsToCheck, collection.CollectionPermissions.CanUpdateValidTokenIds, "can create more tokens")
		if err != nil {
			return &types.TokenCollection{}, err
		}
	}

	collection.ValidTokenIds = allTokenIds

	return collection, nil
}
