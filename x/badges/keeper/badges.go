package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Create tokens and update the unminted / total supplys for the collection
func (k Keeper) CreateBadges(ctx sdk.Context, collection *types.BadgeCollection, newValidBadgeIds []*types.UintRange) (*types.BadgeCollection, error) {
	//For readability, we do not allow transfers to happen on-chain, if not defined in the collection
	if !IsStandardBalances(collection) {
		if len(collection.CollectionApprovals) > 0 {
			return &types.BadgeCollection{}, ErrWrongBalancesType
		}
	}

	var err error
	allBadgeIds := []*types.UintRange{}
	allBadgeIds = append(allBadgeIds, newValidBadgeIds...)
	allBadgeIds, err = types.SortUintRangesAndMerge(allBadgeIds, true)
	if err != nil {
		return &types.BadgeCollection{}, err
	}

	if len(allBadgeIds) > 1 || (len(allBadgeIds) == 1 && !allBadgeIds[0].Start.Equal(sdkmath.NewUint(1))) {
		return &types.BadgeCollection{}, sdkerrors.Wrapf(types.ErrNotSupported, "Ids must be sequential starting from 1")
	}

	if len(newValidBadgeIds) == 0 {
		return collection, nil
	}

	//Check if we are allowed to create these tokens
	detailsToCheck := []*types.UniversalPermissionDetails{}
	for _, badgeIdRange := range newValidBadgeIds {
		detailsToCheck = append(detailsToCheck, &types.UniversalPermissionDetails{
			BadgeId: badgeIdRange,
		})
	}

	err = k.CheckIfBadgeIdsActionPermissionPermits(ctx, detailsToCheck, collection.CollectionPermissions.CanUpdateValidBadgeIds, "can create more tokens")
	if err != nil {
		return &types.BadgeCollection{}, err
	}

	collection.ValidBadgeIds = allBadgeIds

	return collection, nil
}
