package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)


func GetBadgeMetadataTimesAndValues(timeline []*types.BadgeMetadataTimeline) ([][]*types.IdRange, []interface{}) {
	times := [][]*types.IdRange{}
	values := []interface{}{}
	for _, timelineVal := range timeline {
		times = append(times, timelineVal.Times)
		values = append(values, timelineVal.BadgeMetadata)
	}
	return times, values
}


func (k msgServer) MintAndDistributeBadges(goCtx context.Context, msg *types.MsgMintAndDistributeBadges) (*types.MsgMintAndDistributeBadgesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	checkIfCanCreateMoreBadges := msg.BadgesToCreate != nil && len(msg.BadgesToCreate) > 0
	collection, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:             msg.Creator,
		CollectionId:        msg.CollectionId,
		MustBeManager:       true,
	})
	if err != nil {
		return nil, err
	}

	if collection.BalancesType.LTE(sdk.NewUint(0)) {
		return nil, ErrOffChainBalances
	}

	if checkIfCanCreateMoreBadges {
		badgeIdsToCheck := []*types.IdRange{}
		for _, badge := range msg.BadgesToCreate {
			badgeIdsToCheck = append(badgeIdsToCheck, badge.BadgeIds...)
		}
		badgeIdsToCheck = types.SortAndMergeOverlapping(badgeIdsToCheck)

		oldTimes, oldValues := GetBadgeMetadataTimesAndValues(collection.BadgeMetadataTimeline)
		oldTimelineFirstMatches := GetFirstMatchOnlyForTimeline(oldTimes, oldValues)

		newTimes, newValues := GetBadgeMetadataTimesAndValues(msg.BadgeMetadataTimeline)
		newTimelineFirstMatches := GetFirstMatchOnlyForTimeline(newTimes, newValues)

		detailsToCheck := GetDetailsToCheck(oldTimelineFirstMatches, newTimelineFirstMatches, func(oldValue interface{}, newValue interface{}) []*types.UniversalPermissionDetails {
			detailsToReturn := []*types.UniversalPermissionDetails{}
			oldBadgeMetadata := oldValue.([]*types.BadgeMetadata)
			castedPermissions := []*types.UniversalPermission{}
			for _, badgeMetadata := range oldBadgeMetadata {
				castedPermissions = append(castedPermissions, &types.UniversalPermission{
					DefaultValues: &types.UniversalDefaultValues{
						BadgeIds: badgeMetadata.BadgeIds,
						UsesBadgeIds: true,
						ArbitraryValue: badgeMetadata.Uri + "<><><>" + badgeMetadata.CustomData,
					},
					Combinations: []*types.UniversalCombination{{}},
				})
			}
			firstMatchesForOld := types.GetFirstMatchOnly(castedPermissions)

			newBadgeMetadata := newValue.([]*types.BadgeMetadata)
			castedPermissions = []*types.UniversalPermission{}
			for _, badgeMetadata := range newBadgeMetadata {
				castedPermissions = append(castedPermissions, &types.UniversalPermission{
					DefaultValues: &types.UniversalDefaultValues{
						BadgeIds: badgeMetadata.BadgeIds,
						UsesBadgeIds: true,
						ArbitraryValue: badgeMetadata.Uri + "<><><>" + badgeMetadata.CustomData,
					},
					Combinations: []*types.UniversalCombination{{}},
				})
			}
			firstMatchesForNew := types.GetFirstMatchOnly(castedPermissions)

			inOldButNotNew := []*types.UniversalPermissionDetails{}
			inNewButNotOld := []*types.UniversalPermissionDetails{}
			for _, oldDetails := range firstMatchesForOld {
				inOldButNotNew = append(inOldButNotNew, oldDetails)
			}

			for _, newDetails := range firstMatchesForNew {
				inNewButNotOld = append(inNewButNotOld, newDetails)
			}

			overlapping := []*types.UniversalPermissionDetails{}
			for _, oldDetails := range firstMatchesForOld {
				for _, newDetails := range firstMatchesForNew {
					//We use the fact that no two details will have duplicate overlaps because we use GetFirstMatchOnly
					_, overlaps := types.UniversalRemoveOverlaps(oldDetails, newDetails)
					oldVal := oldDetails.ArbitraryValue.(string)
					newVal := newDetails.ArbitraryValue.(string)
					if newVal != oldVal {
						for _, overlap := range overlaps {
							detailsToReturn = append(detailsToReturn, overlap)
						}
					}

					for _, overlap := range overlaps {
						overlapping = append(overlapping, overlap)

						newInOld := []*types.UniversalPermissionDetails{}
						for _, inOld := range inOldButNotNew {
							_, removed := types.UniversalRemoveOverlaps(inOld, overlap)
							for _, removedVal := range removed {
								newInOld = append(newInOld, removedVal)
							}
						}
						inOldButNotNew = newInOld

						newInNew := []*types.UniversalPermissionDetails{}
						for _, inNew := range inNewButNotOld {
							_, removed := types.UniversalRemoveOverlaps(inNew, overlap)
							for _, removedVal := range removed {
								newInNew = append(newInNew, removedVal)
							}
						}
						inNewButNotOld = newInNew
					}
				}
			}

			return detailsToReturn
		})

		err = CheckActionWithBadgeIdsPermission(ctx, detailsToCheck, collection.Permissions.CanCreateMoreBadges)
		if err != nil {
			return nil, err
		}

		//TODO: Add create badges here?
	}

	//TODO:
	// newCollectionMetadata, newBadgeMetadata, newOffChainBalancesMetadata, needToValidateUpdateCollectionMetadata, needToValidateUpdateBadgeMetadata, needToValidateUpdateBalanceUri := GetUrisToStoreAndPermissionsToCheck(collection, msg.CollectionMetadata, msg.BadgeMetadata, msg.OffChainBalancesMetadata)
	// newApprovedTransfers, needToValidateUpdateCollectionApprovedTransfers := GetApprovedTransfersToStore(collection, msg.ApprovedTransfers)

	// _, err = k.UniversalValidate(ctx, UniversalValidationParams{
	// 	Creator:                              msg.Creator,
	// 	CollectionId:                         msg.CollectionId,
	// 	MustBeManager:                        true,
	// 	CanUpdateOffChainBalancesMetadata:            needToValidateUpdateBalanceUri,
	// 	CanUpdateBadgeMetadata:               needToValidateUpdateBadgeMetadata,
	// 	CanUpdateCollectionMetadata:          needToValidateUpdateCollectionMetadata,
	// 	CanUpdateCollectionApprovedTransfers: needToValidateUpdateCollectionApprovedTransfers,
	// })
	// if err != nil {
	// 	return nil, err
	// }

	// //Check badge metadata for isFrozen logic
	// err = AssertIsFrozenLogicIsMaintained(collection.BadgeMetadata, newBadgeMetadata)
	// if err != nil {
	// 	return nil, err
	// }

	// err = AssertIsFrozenLogicForApprovedTransfers(collection.ApprovedTransfers, newApprovedTransfers)
	// if err != nil {
	// 	return nil, err
	// }

	// collection.BadgeMetadata = newBadgeMetadata
	// collection.CollectionMetadata = newCollectionMetadata
	// collection.OffChainBalancesMetadata = newOffChainBalancesMetadata
	// collection.ApprovedTransfers = newApprovedTransfers

	collection, err = k.CreateBadges(ctx, collection, msg.BadgesToCreate, msg.Transfers, msg.Creator)
	if err != nil {
		return nil, err
	}

	if err := k.SetCollectionInStore(ctx, collection); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		),
	)

	return &types.MsgMintAndDistributeBadgesResponse{
		NextBadgeId: collection.NextBadgeId,
	}, nil
}
