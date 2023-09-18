package keeper

import (
	"math"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	proto "github.com/gogo/protobuf/proto"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//This file is responsible for verifying that if we go from Value A to Value B for a timeline, that the update is valid
//So here, we check the following:
//-Assert the collection has the correct balances type, if we are updating a balance type-specific field
//-For all updates, check that we are able to update according to the permissions.
// This means we have to verify that the current permissions do not forbid the update.
//
//To do this, we have to do the following:
//-1) Get the combination of values which are "updated". Note this depends on the field, so this is kept generic through a function (GetUpdateCombinationsToCheck)
//		We are also dealing with timelines, so "updated" depends on the respective times, as well as the field value.
//		We do this in a multi step-process:
//		-First, we cast the timeline to UniversalPermission using only the TimelineTimes. We store the timeline VALUE in the ArbitraryValue field.
//		-Second, we get the overlaps and non-overlaps between the (old, new) timeline times.
//		 Note we also have to handle edge cases (in one but not the other). We add empty values where needed.
//		 This then leaves us with a list of all the (timeA-timeB, valueX) - (timeA-timeB, valueY) pairs we need to check.
//		-Third, we compare all valueX and valueY values to see if the actual value was updated.
//		 If the value was not updated, then for timeA-timeB, we do not need to check the permissions.
//		 If it was updated, we need to check the permissions for timeA-timeB.
//		-Lastly, if it was updated, in addition to just simply checking timeA-timeB, we may also have to be more specific with what that we need to check.
//		 Ex: If we go from [badgeIDs 1 to 10 -> www.example.com] to [badgeIDs 1 to 2 -> www.example2.com, badgeIDs 3 to 10 -> www.example.com],
//				 we only need to check badgeIDs 1 to 2 from timeA-timeB
//		 We eventually end with a (timeA-timeB, badgeIds, transferTimes, toMapping, fromMapping, initiatedByMapping) tuples array[] that we need to check, adding dummy values where needed.
//		 This step and the third step is field-specific, so that is why we do it via a generic custom function (GetUpdatedStringCombinations, GetUpdatedBoolCombinations, etc...)
//-2) For all the values that are considered "updated", we check if we are allowed to update them, according to the permissions.
//		This is done by fetching wherever the returned tuples from above overlaps any of the permission's (timelineTime, badgeIds, transferTimes, toMapping, fromMapping, initiatedByMapping) tuples, again adding dummy values where needed.
//		For all overlaps, we then assert that the current block time is NOT forbidden (permitted or undefined both correspond to allowed)
//		If all are not forbidden, it is a valid update.

// To make it easier, we first
func GetPotentialUpdatesForTimelineValues(times [][]*types.UintRange, values []interface{}) []*types.UniversalPermissionDetails {
	castedPermissions := []*types.UniversalPermission{}
	for idx, time := range times {
		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			DefaultValues: &types.UniversalDefaultValues{
				TimelineTimes:     time,
				ArbitraryValue:    values[idx],
				UsesTimelineTimes: true,
			},
			Combinations: []*types.UniversalCombination{{}},
		})
	}

	firstMatches := types.GetFirstMatchOnly(castedPermissions) //I think this is unnecessary because we already disallow duplicate timeline times in ValidateBasic but if we allow duplicates, this may be needed

	return firstMatches
}

func (k Keeper) ValidateCollectionApprovedTransfersUpdate(ctx sdk.Context, collection *types.BadgeCollection, oldApprovedTransfers []*types.CollectionApprovedTransfer, newApprovedTransfers []*types.CollectionApprovedTransfer, CanUpdateCollectionApprovedTransfers []*types.CollectionApprovedTransferPermission, managerAddress string) error {
	if !IsStandardBalances(collection) && newApprovedTransfers != nil && len(newApprovedTransfers) > 0 {
		return sdkerrors.Wrapf(ErrWrongBalancesType, "collection %s does not have standard balances", collection.CollectionId)
	}

	x := [][]*types.UintRange{}
	x = append(x, []*types.UintRange{
		//Dummmy range
		{
			Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64),
		},
	})

	y := [][]*types.UintRange{}
	y = append(y, []*types.UintRange{
		//Dummmy range
		{
			Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64),
		},
	})

	//This is just to maintain consistency with the legacy features when we used to have timeline times
	oldTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(x, []interface{}{oldApprovedTransfers})

	newTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(y, []interface{}{newApprovedTransfers})

	detailsToCheck, err := GetUpdateCombinationsToCheck(ctx, oldTimelineFirstMatches, newTimelineFirstMatches, []*types.CollectionApprovedTransfer{}, managerAddress, func(ctx sdk.Context, oldValue interface{}, newValue interface{}, managerAddress string) ([]*types.UniversalPermissionDetails, error) {
		//Cast to UniversalPermissionDetails for comaptibility with these overlap functions and get first matches only (i.e. first match for each badge ID)
		oldApprovedTransfers, err := k.CastCollectionApprovedTransferToUniversalPermission(ctx, oldValue.([]*types.CollectionApprovedTransfer), managerAddress)
		if err != nil {
			return nil, err
		}
		firstMatchesForOld := types.GetFirstMatchOnly(oldApprovedTransfers)

		newApprovedTransfers, err := k.CastCollectionApprovedTransferToUniversalPermission(ctx, newValue.([]*types.CollectionApprovedTransfer), managerAddress)
		if err != nil {
			return nil, err
		}
		firstMatchesForNew := types.GetFirstMatchOnly(newApprovedTransfers)

		//For every badge, we need to check if the new provided value is different in any way from the old value for each badge ID
		//The overlapObjects from GetOverlapsAndNonOverlaps will return which badge IDs overlap
		detailsToReturn := []*types.UniversalPermissionDetails{}
		overlapObjects, inOldButNotNew, inNewButNotOld := types.GetOverlapsAndNonOverlaps(firstMatchesForOld, firstMatchesForNew)
		for _, overlapObject := range overlapObjects {
			overlap := overlapObject.Overlap
			oldDetails := overlapObject.FirstDetails
			newDetails := overlapObject.SecondDetails
			different := false
			if (oldDetails.ArbitraryValue == nil && newDetails.ArbitraryValue != nil) || (oldDetails.ArbitraryValue != nil && newDetails.ArbitraryValue == nil) {
				different = true
			} else {
				oldVal := oldDetails.ArbitraryValue.(*types.CollectionApprovedTransfer)
				newVal := newDetails.ArbitraryValue.(*types.CollectionApprovedTransfer)

				if len(oldVal.ApprovalDetails) != len(newVal.ApprovalDetails) {
					different = true
				} else {
					for i := 0; i < len(oldVal.ApprovalDetails); i++ {
						if proto.MarshalTextString(oldVal.ApprovalDetails[i]) != proto.MarshalTextString(newVal.ApprovalDetails[i]) {
							different = true
						}
					}
				}
			}

			if different {
				detailsToReturn = append(detailsToReturn, overlap)
			}
		}

		//If there are combinations in old but not new, then it is considered updated. If it is in new but not old, then it is considered updated.
		detailsToReturn = append(detailsToReturn, inOldButNotNew...)
		detailsToReturn = append(detailsToReturn, inNewButNotOld...)

		return detailsToReturn, nil
	})
	if err != nil {
		return err
	}

	err = k.CheckCollectionApprovedTransferPermission(ctx, detailsToCheck, CanUpdateCollectionApprovedTransfers, managerAddress, "update collection approved transfers")
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) ValidateUserApprovedOutgoingTransfersUpdate(ctx sdk.Context, _oldApprovedTransfers []*types.UserApprovedOutgoingTransfer, _newApprovedTransfers []*types.UserApprovedOutgoingTransfer, CanUpdateCollectionApprovedTransfers []*types.UserApprovedOutgoingTransferPermission, managerAddress string) error {
	x := [][]*types.UintRange{}
	x = append(x, []*types.UintRange{
		//Dummmy range
		{ Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64), },
	})

	y := [][]*types.UintRange{}
	y = append(y, []*types.UintRange{
		//Dummmy range
		{ Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64), },
	})

	//This is just to maintain consistency with the legacy features when we used to have timeline times
	oldTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(x, []interface{}{_oldApprovedTransfers})

	newTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(y, []interface{}{_newApprovedTransfers})


	detailsToCheck, err := GetUpdateCombinationsToCheck(ctx, oldTimelineFirstMatches, newTimelineFirstMatches, []*types.UserApprovedOutgoingTransfer{}, managerAddress, func(ctx sdk.Context, oldValue interface{}, newValue interface{}, managerAddress string) ([]*types.UniversalPermissionDetails, error) {
		//Cast to UniversalPermissionDetails for comaptibility with these overlap functions and get first matches only (i.e. first match for each badge ID)
		oldApprovedTransfers, err := k.CastUserApprovedOutgoingTransferToUniversalPermission(ctx, oldValue.([]*types.UserApprovedOutgoingTransfer), managerAddress)
		if err != nil {
			return nil, err
		}
		firstMatchesForOld := types.GetFirstMatchOnly(oldApprovedTransfers)

		newApprovedTransfers, err := k.CastUserApprovedOutgoingTransferToUniversalPermission(ctx, newValue.([]*types.UserApprovedOutgoingTransfer), managerAddress)
		if err != nil {
			return nil, err
		}
		firstMatchesForNew := types.GetFirstMatchOnly(newApprovedTransfers)

		//For every badge, we need to check if the new provided value is different in any way from the old value for each badge ID
		//The overlapObjects from GetOverlapsAndNonOverlaps will return which badge IDs overlap
		detailsToReturn := []*types.UniversalPermissionDetails{}
		overlapObjects, inOldButNotNew, inNewButNotOld := types.GetOverlapsAndNonOverlaps(firstMatchesForOld, firstMatchesForNew)
		for _, overlapObject := range overlapObjects {
			overlap := overlapObject.Overlap
			oldDetails := overlapObject.FirstDetails
			newDetails := overlapObject.SecondDetails
			different := false
			if (oldDetails.ArbitraryValue == nil && newDetails.ArbitraryValue != nil) || (oldDetails.ArbitraryValue != nil && newDetails.ArbitraryValue == nil) {
				different = true
			} else {
				oldVal := oldDetails.ArbitraryValue.(*types.UserApprovedOutgoingTransfer)
				newVal := newDetails.ArbitraryValue.(*types.UserApprovedOutgoingTransfer)

				if len(oldVal.ApprovalDetails) != len(newVal.ApprovalDetails) {
					different = true
				} else {
					for i := 0; i < len(oldVal.ApprovalDetails); i++ {
						if proto.MarshalTextString(oldVal.ApprovalDetails[i]) != proto.MarshalTextString(newVal.ApprovalDetails[i]) {
							different = true
						}
					}
				}
			}

			if different {
				detailsToReturn = append(detailsToReturn, overlap)
			}
		}

		//If there are combinations in old but not new, then it is considered updated. If it is in new but not old, then it is considered updated.
		detailsToReturn = append(detailsToReturn, inOldButNotNew...)
		detailsToReturn = append(detailsToReturn, inNewButNotOld...)

		return detailsToReturn, nil
	})
	if err != nil {
		return err
	}

	err = k.CheckUserApprovedOutgoingTransferPermission(ctx, detailsToCheck, CanUpdateCollectionApprovedTransfers, managerAddress, "update user approved outgoing transfers")
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) ValidateUserApprovedIncomingTransfersUpdate(ctx sdk.Context, _oldApprovedTransfers []*types.UserApprovedIncomingTransfer, _newApprovedTransfers []*types.UserApprovedIncomingTransfer, CanUpdateCollectionApprovedTransfers []*types.UserApprovedIncomingTransferPermission, managerAddress string) error {
	x := [][]*types.UintRange{}
	x = append(x, []*types.UintRange{
		//Dummmy range
		{ Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64), },
	})

	y := [][]*types.UintRange{}
	y = append(y, []*types.UintRange{
		//Dummmy range
		{ Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64), },
	})

	//This is just to maintain consistency with the legacy features when we used to have timeline times
	oldTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(x, []interface{}{_oldApprovedTransfers})

	newTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(y, []interface{}{_newApprovedTransfers})

	detailsToCheck, err := GetUpdateCombinationsToCheck(ctx, oldTimelineFirstMatches, newTimelineFirstMatches, []*types.UserApprovedIncomingTransfer{}, managerAddress, func(ctx sdk.Context, oldValue interface{}, newValue interface{}, managerAddress string) ([]*types.UniversalPermissionDetails, error) {
		//Cast to UniversalPermissionDetails for comaptibility with these overlap functions and get first matches only (i.e. first match for each badge ID)
		oldApprovedTransfers, err := k.CastUserApprovedIncomingTransferToUniversalPermission(ctx, oldValue.([]*types.UserApprovedIncomingTransfer), managerAddress)
		if err != nil {
			return nil, err
		}
		firstMatchesForOld := types.GetFirstMatchOnly(oldApprovedTransfers)

		newApprovedTransfers, err := k.CastUserApprovedIncomingTransferToUniversalPermission(ctx, newValue.([]*types.UserApprovedIncomingTransfer), managerAddress)
		if err != nil {
			return nil, err
		}
		firstMatchesForNew := types.GetFirstMatchOnly(newApprovedTransfers)

		//For every badge, we need to check if the new provided value is different in any way from the old value for each badge ID
		//The overlapObjects from GetOverlapsAndNonOverlaps will return which badge IDs overlap
		detailsToReturn := []*types.UniversalPermissionDetails{}
		overlapObjects, inOldButNotNew, inNewButNotOld := types.GetOverlapsAndNonOverlaps(firstMatchesForOld, firstMatchesForNew)
		for _, overlapObject := range overlapObjects {
			overlap := overlapObject.Overlap
			oldDetails := overlapObject.FirstDetails
			newDetails := overlapObject.SecondDetails
			different := false
			if (oldDetails.ArbitraryValue == nil && newDetails.ArbitraryValue != nil) || (oldDetails.ArbitraryValue != nil && newDetails.ArbitraryValue == nil) {
				different = true
			} else {

				oldVal := oldDetails.ArbitraryValue.(*types.UserApprovedIncomingTransfer)
				newVal := newDetails.ArbitraryValue.(*types.UserApprovedIncomingTransfer)

				if len(oldVal.ApprovalDetails) != len(newVal.ApprovalDetails) {
					different = true
				} else {
					for i := 0; i < len(oldVal.ApprovalDetails); i++ {
						if proto.MarshalTextString(oldVal.ApprovalDetails[i]) != proto.MarshalTextString(newVal.ApprovalDetails[i]) {
							different = true
						}
					}
				}
			}

			if different {
				detailsToReturn = append(detailsToReturn, overlap)
			}
		}

		//If there are combinations in old but not new, then it is considered updated. If it is in new but not old, then it is considered updated.
		detailsToReturn = append(detailsToReturn, inOldButNotNew...)
		detailsToReturn = append(detailsToReturn, inNewButNotOld...)

		return detailsToReturn, nil
	})
	if err != nil {
		return err
	}

	err = k.CheckUserApprovedIncomingTransferPermission(ctx, detailsToCheck, CanUpdateCollectionApprovedTransfers, managerAddress, "update user approved incoming transfers")
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) ValidateBadgeMetadataUpdate(ctx sdk.Context, oldBadgeMetadata []*types.BadgeMetadataTimeline, newBadgeMetadata []*types.BadgeMetadataTimeline, canUpdateBadgeMetadata []*types.TimedUpdateWithBadgeIdsPermission) error {
	oldTimes, oldValues := types.GetBadgeMetadataTimesAndValues(oldBadgeMetadata)
	oldTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(oldTimes, oldValues)

	newTimes, newValues := types.GetBadgeMetadataTimesAndValues(newBadgeMetadata)
	newTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(newTimes, newValues)

	detailsToCheck, err := GetUpdateCombinationsToCheck(ctx, oldTimelineFirstMatches, newTimelineFirstMatches, []*types.BadgeMetadata{}, "", func(ctx sdk.Context, oldValue interface{}, newValue interface{}, managerAddress string) ([]*types.UniversalPermissionDetails, error) {
		//Cast to UniversalPermissionDetails for comaptibility with these overlap functions and get first matches only (i.e. first match for each badge ID)
		oldBadgeMetadata := oldValue.([]*types.BadgeMetadata)
		firstMatchesForOld := types.GetFirstMatchOnly(k.CastBadgeMetadataToUniversalPermission(oldBadgeMetadata))

		newBadgeMetadata := newValue.([]*types.BadgeMetadata)
		firstMatchesForNew := types.GetFirstMatchOnly(k.CastBadgeMetadataToUniversalPermission(newBadgeMetadata))

		//For every badge, we need to check if the new provided value is different in any way from the old value for each badge ID
		//The overlapObjects from GetOverlapsAndNonOverlaps will return which badge IDs overlap
		detailsToReturn := []*types.UniversalPermissionDetails{}
		overlapObjects, inOldButNotNew, inNewButNotOld := types.GetOverlapsAndNonOverlaps(firstMatchesForOld, firstMatchesForNew)
		for _, overlapObject := range overlapObjects {
			overlap := overlapObject.Overlap
			oldDetails := overlapObject.FirstDetails
			newDetails := overlapObject.SecondDetails
			//HACK: We set it to a string beforehand
			if (oldDetails.ArbitraryValue == nil && newDetails.ArbitraryValue != nil) || (oldDetails.ArbitraryValue != nil && newDetails.ArbitraryValue == nil) {
				detailsToReturn = append(detailsToReturn, overlap)
			} else {
				oldVal := oldDetails.ArbitraryValue.(string)
				newVal := newDetails.ArbitraryValue.(string)

				if newVal != oldVal {
					detailsToReturn = append(detailsToReturn, overlap)
				}
			}
		}

		//If metadata is in old but not new, then it is considered updated. If it is in new but not old, then it is considered updated.
		detailsToReturn = append(detailsToReturn, inOldButNotNew...)
		detailsToReturn = append(detailsToReturn, inNewButNotOld...)

		return detailsToReturn, nil
	})
	if err != nil {
		return err
	}

	err = k.CheckTimedUpdateWithBadgeIdsPermission(ctx, detailsToCheck, canUpdateBadgeMetadata, "update badge metadata")
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) ValidateCollectionMetadataUpdate(ctx sdk.Context, oldCollectionMetadata []*types.CollectionMetadataTimeline, newCollectionMetadata []*types.CollectionMetadataTimeline, canUpdateCollectionMetadata []*types.TimedUpdatePermission) error {
	oldTimes, oldValues := types.GetCollectionMetadataTimesAndValues(oldCollectionMetadata)
	oldTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(oldTimes, oldValues)

	newTimes, newValues := types.GetCollectionMetadataTimesAndValues(newCollectionMetadata)
	newTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(newTimes, newValues)

	detailsToCheck, err := GetUpdateCombinationsToCheck(ctx, oldTimelineFirstMatches, newTimelineFirstMatches, &types.CollectionMetadata{}, "", func(ctx sdk.Context, oldValue interface{}, newValue interface{}, managerAddress string) ([]*types.UniversalPermissionDetails, error) {
		detailsToCheck := []*types.UniversalPermissionDetails{}
		if oldValue == nil && newValue != nil {
			detailsToCheck = append(detailsToCheck, &types.UniversalPermissionDetails{})
		} else {
			oldVal := oldValue.(*types.CollectionMetadata)
			newVal := newValue.(*types.CollectionMetadata)

			if oldVal.Uri != newVal.Uri || oldVal.CustomData != newVal.CustomData {
				detailsToCheck = append(detailsToCheck, &types.UniversalPermissionDetails{})
			}
		}
		return detailsToCheck, nil
	})
	if err != nil {
		return err
	}

	err = k.CheckTimedUpdatePermission(ctx, detailsToCheck, canUpdateCollectionMetadata, "update collection metadata")
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) ValidateOffChainBalancesMetadataUpdate(ctx sdk.Context, collection *types.BadgeCollection, oldOffChainBalancesMetadata []*types.OffChainBalancesMetadataTimeline, newOffChainBalancesMetadata []*types.OffChainBalancesMetadataTimeline, canUpdateOffChainBalancesMetadata []*types.TimedUpdatePermission) error {
	if !IsOffChainBalances(collection) {
		if len(oldOffChainBalancesMetadata) > 0 || len(newOffChainBalancesMetadata) > 0 {
			return sdkerrors.Wrapf(ErrWrongBalancesType, "off chain balances are being set but collection %s does not have off chain balances", collection.CollectionId)
		}
		return nil
	}

	oldTimes, oldValues := types.GetOffChainBalancesMetadataTimesAndValues(oldOffChainBalancesMetadata)
	oldTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(oldTimes, oldValues)

	newTimes, newValues := types.GetOffChainBalancesMetadataTimesAndValues(newOffChainBalancesMetadata)
	newTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(newTimes, newValues)

	detailsToCheck, err := GetUpdateCombinationsToCheck(ctx, oldTimelineFirstMatches, newTimelineFirstMatches, &types.OffChainBalancesMetadata{}, "", func(ctx sdk.Context, oldValue interface{}, newValue interface{}, managerAddress string) ([]*types.UniversalPermissionDetails, error) {
		detailsToCheck := []*types.UniversalPermissionDetails{}
		if oldValue == nil && newValue != nil {
			detailsToCheck = append(detailsToCheck, &types.UniversalPermissionDetails{})
		} else {

			oldVal := oldValue.(*types.OffChainBalancesMetadata)
			newVal := newValue.(*types.OffChainBalancesMetadata)

			if oldVal.Uri != newVal.Uri || oldVal.CustomData != newVal.CustomData {
				detailsToCheck = append(detailsToCheck, &types.UniversalPermissionDetails{})
			}
		}
		return detailsToCheck, nil
	})
	if err != nil {
		return err
	}

	err = k.CheckTimedUpdatePermission(ctx, detailsToCheck, canUpdateOffChainBalancesMetadata, "update off chain balances metadata")
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) ValidateInheritedBalancesUpdate(ctx sdk.Context, collection *types.BadgeCollection, oldInheritedBalances []*types.InheritedBalancesTimeline, newInheritedBalances []*types.InheritedBalancesTimeline, canUpdateInheritedBalances []*types.TimedUpdateWithBadgeIdsPermission) error {
	if !IsInheritedBalances(collection) {
		if len(oldInheritedBalances) > 0 || len(newInheritedBalances) > 0 {
			return sdkerrors.Wrapf(ErrWrongBalancesType, "inherited balances are being set but collection %s does not have inherited balances", collection.CollectionId)
		}
		return nil
	}

	//Enforce that badge IDs are sequential starting from 1
	for _, timelineVal := range newInheritedBalances {
		allBadgeIds := []*types.UintRange{}
		for _, inheritedBalance := range timelineVal.InheritedBalances {
			allBadgeIds = append(allBadgeIds, inheritedBalance.BadgeIds...)
		}

		allBadgeIds = types.SortAndMergeOverlapping(allBadgeIds)
		if len(allBadgeIds) > 1 || (len(allBadgeIds) == 1 && !allBadgeIds[0].Start.Equal(sdkmath.NewUint(1))) {
			return sdkerrors.Wrapf(types.ErrNotSupported, "BadgeIds must be sequential starting from 1")
		}
	}

	oldTimes, oldValues := types.GetInheritedBalancesTimesAndValues(oldInheritedBalances)
	oldTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(oldTimes, oldValues)

	newTimes, newValues := types.GetInheritedBalancesTimesAndValues(newInheritedBalances)
	newTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(newTimes, newValues)

	detailsToCheck, err := GetUpdateCombinationsToCheck(ctx, oldTimelineFirstMatches, newTimelineFirstMatches, []*types.InheritedBalance{}, "", func(ctx sdk.Context, oldValue interface{}, newValue interface{}, managerAddress string) ([]*types.UniversalPermissionDetails, error) {
		//Cast to UniversalPermissionDetails for comaptibility with these overlap functions and get first matches only (i.e. first match for each badge ID)
		oldInheritedBalances := oldValue.([]*types.InheritedBalance)
		firstMatchesForOld := types.GetFirstMatchOnly(k.CastInheritedBalancesToUniversalPermission(oldInheritedBalances))

		newInheritedBalances := newValue.([]*types.InheritedBalance)
		firstMatchesForNew := types.GetFirstMatchOnly(k.CastInheritedBalancesToUniversalPermission(newInheritedBalances))

		//For every badge, we need to check if the new provided value is different in any way from the old value for each badge ID
		//The overlapObjects from GetOverlapsAndNonOverlaps will return which badge IDs overlap
		detailsToReturn := []*types.UniversalPermissionDetails{}
		overlapObjects, inOldButNotNew, inNewButNotOld := types.GetOverlapsAndNonOverlaps(firstMatchesForOld, firstMatchesForNew)
		for _, overlapObject := range overlapObjects {
			overlap := overlapObject.Overlap
			oldDetails := overlapObject.FirstDetails
			newDetails := overlapObject.SecondDetails
			different := false

			if (oldDetails.ArbitraryValue == nil && newDetails.ArbitraryValue != nil) || (oldDetails.ArbitraryValue != nil && newDetails.ArbitraryValue == nil) {
				different = true
			} else {
				oldVal := oldDetails.ArbitraryValue.(*types.InheritedBalance)
				newVal := newDetails.ArbitraryValue.(*types.InheritedBalance)

				if proto.MarshalTextString(newVal) != proto.MarshalTextString(oldVal) {
					different = true

				}
			}

			if different {
				detailsToReturn = append(detailsToReturn, overlap)
			}
		}

		//If metadata is in old but not new, then it is considered updated. If it is in new but not old, then it is considered updated.
		detailsToReturn = append(detailsToReturn, inOldButNotNew...)
		detailsToReturn = append(detailsToReturn, inNewButNotOld...)

		return detailsToReturn, nil
	})
	if err != nil {
		return err
	}

	err = k.CheckTimedUpdateWithBadgeIdsPermission(ctx, detailsToCheck, canUpdateInheritedBalances, "update inherited balances")
	if err != nil {
		return err
	}

	return nil
}

/** Everything below here is pretty standard because all we need to compare is primitive types **/

func GetUpdatedStringCombinations(ctx sdk.Context, oldValue interface{}, newValue interface{}, managerAddress string) ([]*types.UniversalPermissionDetails, error) {
	x := []*types.UniversalPermissionDetails{}
	if (oldValue == nil && newValue != nil) || (oldValue != nil && newValue == nil) {
		x = append(x, &types.UniversalPermissionDetails{})
	} else if oldValue.(string) != newValue.(string) {
		x = append(x, &types.UniversalPermissionDetails{})
	}
	return x, nil
}

func GetUpdatedBoolCombinations(ctx sdk.Context, oldValue interface{}, newValue interface{}, managerAddress string) ([]*types.UniversalPermissionDetails, error) {
	if (oldValue == nil && newValue != nil) || (oldValue != nil && newValue == nil) {
		return []*types.UniversalPermissionDetails{{}}, nil
	}

	oldVal := oldValue.(bool)
	newVal := newValue.(bool)
	if oldVal != newVal {
		return []*types.UniversalPermissionDetails{
			{},
		}, nil
	}
	return []*types.UniversalPermissionDetails{}, nil
}

func (k Keeper) ValidateManagerUpdate(ctx sdk.Context, oldManager []*types.ManagerTimeline, newManager []*types.ManagerTimeline, canUpdateManager []*types.TimedUpdatePermission) error {
	oldTimes, oldValues := types.GetManagerTimesAndValues(oldManager)
	oldTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(oldTimes, oldValues)

	newTimes, newValues := types.GetManagerTimesAndValues(newManager)
	newTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(newTimes, newValues)

	updatedTimelineTimes, err := GetUpdateCombinationsToCheck(ctx, oldTimelineFirstMatches, newTimelineFirstMatches, "", "", GetUpdatedStringCombinations)
	if err != nil {
		return err
	}

	if err = k.CheckTimedUpdatePermission(ctx, updatedTimelineTimes, canUpdateManager, "update manager"); err != nil {
		return err
	}

	return nil
}

func (k Keeper) ValidateCustomDataUpdate(ctx sdk.Context, oldCustomData []*types.CustomDataTimeline, newCustomData []*types.CustomDataTimeline, canUpdateCustomData []*types.TimedUpdatePermission) error {
	oldTimes, oldValues := types.GetCustomDataTimesAndValues(oldCustomData)
	oldTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(oldTimes, oldValues)

	newTimes, newValues := types.GetCustomDataTimesAndValues(newCustomData)
	newTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(newTimes, newValues)

	updatedTimelineTimes, err := GetUpdateCombinationsToCheck(ctx, oldTimelineFirstMatches, newTimelineFirstMatches, "", "", GetUpdatedStringCombinations)
	if err != nil {
		return err
	}

	if err = k.CheckTimedUpdatePermission(ctx, updatedTimelineTimes, canUpdateCustomData, "update custom data"); err != nil {
		return err
	}

	return nil
}

func (k Keeper) ValidateStandardsUpdate(ctx sdk.Context, oldStandards []*types.StandardsTimeline, newStandards []*types.StandardsTimeline, canUpdateStandards []*types.TimedUpdatePermission) error {
	oldTimes, oldValues := types.GetStandardsTimesAndValues(oldStandards)
	oldTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(oldTimes, oldValues)

	newTimes, newValues := types.GetStandardsTimesAndValues(newStandards)
	newTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(newTimes, newValues)

	updatedTimelineTimes, err := GetUpdateCombinationsToCheck(ctx, oldTimelineFirstMatches, newTimelineFirstMatches, "", "", func(ctx sdk.Context, oldValue interface{}, newValue interface{}, managerAddress string) ([]*types.UniversalPermissionDetails, error) {
		if (oldValue == nil && newValue != nil) || (oldValue != nil && newValue == nil) {
			return []*types.UniversalPermissionDetails{{}}, nil
		} else {
			oldVal := oldValue.([]string)
			newVal := newValue.([]string)

			if len(oldVal) != len(newVal) {
				return []*types.UniversalPermissionDetails{{}}, nil
			} else {
				for i := 0; i < len(oldVal); i++ {
					if oldVal[i] != newVal[i] {
						return []*types.UniversalPermissionDetails{{}}, nil
					}
				}
			}
		}

		return []*types.UniversalPermissionDetails{}, nil
	})
	if err != nil {
		return err
	}

	if err = k.CheckTimedUpdatePermission(ctx, updatedTimelineTimes, canUpdateStandards, "update standards"); err != nil {
		return err
	}

	return nil
}

func (k Keeper) ValidateContractAddressUpdate(ctx sdk.Context, oldContractAddress []*types.ContractAddressTimeline, newContractAddress []*types.ContractAddressTimeline, canUpdateContractAddress []*types.TimedUpdatePermission) error {
	oldTimes, oldValues := types.GetContractAddressTimesAndValues(oldContractAddress)
	oldTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(oldTimes, oldValues)

	newTimes, newValues := types.GetContractAddressTimesAndValues(newContractAddress)
	newTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(newTimes, newValues)

	updatedTimelineTimes, err := GetUpdateCombinationsToCheck(ctx, oldTimelineFirstMatches, newTimelineFirstMatches, "", "", GetUpdatedStringCombinations)
	if err != nil {
		return err
	}

	if err = k.CheckTimedUpdatePermission(ctx, updatedTimelineTimes, canUpdateContractAddress, "update contract address"); err != nil {
		return err
	}

	return nil
}

func (k Keeper) ValidateIsArchivedUpdate(ctx sdk.Context, oldIsArchived []*types.IsArchivedTimeline, newIsArchived []*types.IsArchivedTimeline, canUpdateIsArchived []*types.TimedUpdatePermission) error {
	oldTimes, oldValues := types.GetIsArchivedTimesAndValues(oldIsArchived)
	oldTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(oldTimes, oldValues)

	newTimes, newValues := types.GetIsArchivedTimesAndValues(newIsArchived)
	newTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(newTimes, newValues)

	updatedTimelineTimes, err := GetUpdateCombinationsToCheck(ctx, oldTimelineFirstMatches, newTimelineFirstMatches, false, "", GetUpdatedBoolCombinations)
	if err != nil {
		return err
	}

	if err = k.CheckTimedUpdatePermission(ctx, updatedTimelineTimes, canUpdateIsArchived, "update is archived"); err != nil {
		return err
	}

	return nil
}
