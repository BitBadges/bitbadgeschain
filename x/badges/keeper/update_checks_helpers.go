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
//		 We eventually end with a (timeA-timeB, badgeIds, transferTimes, toList, fromList, initiatedByList) tuples array[] that we need to check, adding dummy values where needed.
//		 This step and the third step is field-specific, so that is why we do it via a generic custom function (GetUpdatedStringCombinations, GetUpdatedBoolCombinations, etc...)
//-2) For all the values that are considered "updated", we check if we are allowed to update them, according to the permissions.
//		This is done by fetching wherever the returned tuples from above overlaps any of the permission's (timelineTime, badgeIds, transferTimes, toList, fromList, initiatedByList) tuples, again adding dummy values where needed.
//		For all overlaps, we then assert that the current block time is NOT forbidden (permitted or undefined both correspond to allowed)
//		If all are not forbidden, it is a valid update.

// To make it easier, we first
func GetPotentialUpdatesForTimelineValues(ctx sdk.Context,
	times [][]*types.UintRange, values []interface{}) []*types.UniversalPermissionDetails {
	castedPermissions := []*types.UniversalPermission{}
	for idx, time := range times {
		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			TimelineTimes:     time,
			ArbitraryValue:    values[idx],
			UsesTimelineTimes: true,
		})
	}

	firstMatches := types.GetFirstMatchOnly(ctx, castedPermissions) //I think this is unnecessary because we already disallow duplicate timeline times in ValidateBasic but if we allow duplicates, this may be needed

	return firstMatches
}

// Make a struct witha  bool flag isApproved and an approval details arr
type ApprovalCriteriaWithIsApproved struct {
	ApprovalCriteria *types.ApprovalCriteria
}

func GetFirstMatchOnlyWithApprovalCriteria(ctx sdk.Context, permissions []*types.UniversalPermission) []*types.UniversalPermissionDetails {
	handled := []*types.UniversalPermissionDetails{}
	for _, permission := range permissions {
		badgeIds := types.GetUintRangesWithOptions(permission.BadgeIds, permission.UsesBadgeIds)
		timelineTimes := types.GetUintRangesWithOptions(permission.TimelineTimes, permission.UsesTimelineTimes)
		transferTimes := types.GetUintRangesWithOptions(permission.TransferTimes, permission.UsesTransferTimes)
		ownershipTimes := types.GetUintRangesWithOptions(permission.OwnershipTimes, permission.UsesOwnershipTimes)
		permanentlyPermittedTimes := types.GetUintRangesWithOptions(permission.PermanentlyPermittedTimes, true)
		permanentlyForbiddenTimes := types.GetUintRangesWithOptions(permission.PermanentlyForbiddenTimes, true)

		toList := types.GetListWithOptions(permission.ToList, permission.UsesToList)
		fromList := types.GetListWithOptions(permission.FromList, permission.UsesFromList)
		initiatedByList := types.GetListWithOptions(permission.InitiatedByList, permission.UsesInitiatedByList)

		approvalIdList := types.GetListWithOptions(permission.ApprovalIdList, permission.UsesApprovalId)
		amountTrackerIdList := types.GetListWithOptions(permission.AmountTrackerIdList, permission.UsesAmountTrackerId)
		challengeTrackerIdList := types.GetListWithOptions(permission.ChallengeTrackerIdList, permission.UsesChallengeTrackerId)

		for _, badgeId := range badgeIds {
			for _, timelineTime := range timelineTimes {
				for _, transferTime := range transferTimes {
					for _, ownershipTime := range ownershipTimes {
						arbValue := []*ApprovalCriteriaWithIsApproved{
							{
								ApprovalCriteria: permission.ArbitraryValue.(*types.CollectionApproval).ApprovalCriteria,
							},
						}

						brokenDown := []*types.UniversalPermissionDetails{
							{
								BadgeId:                badgeId,
								TimelineTime:           timelineTime,
								TransferTime:           transferTime,
								OwnershipTime:          ownershipTime,
								ToList:                 toList,
								FromList:               fromList,
								InitiatedByList:        initiatedByList,
								ApprovalIdList:         approvalIdList,
								AmountTrackerIdList:    amountTrackerIdList,
								ChallengeTrackerIdList: challengeTrackerIdList,

								ArbitraryValue: arbValue,
							},
						}

						overlaps, inBrokenDownButNotHandled, inHandledButNotBrokenDown := types.GetOverlapsAndNonOverlaps(ctx, brokenDown, handled)
						handled = []*types.UniversalPermissionDetails{}
						//if no overlaps, we can just append all of them
						handled = append(handled, inHandledButNotBrokenDown...)
						handled = append(handled, inBrokenDownButNotHandled...)

						//for overlaps, we append approval details
						for _, overlap := range overlaps {
							mergedApprovalCriteria := overlap.SecondDetails.ArbitraryValue.([]*ApprovalCriteriaWithIsApproved)

							for _, approvalDetail := range overlap.FirstDetails.ArbitraryValue.([]*ApprovalCriteriaWithIsApproved) {
								mergedApprovalCriteria = append(mergedApprovalCriteria, approvalDetail)
							}

							newArbValue := mergedApprovalCriteria

							handled = append(handled, &types.UniversalPermissionDetails{
								TimelineTime:    overlap.Overlap.TimelineTime,
								BadgeId:         overlap.Overlap.BadgeId,
								TransferTime:    overlap.Overlap.TransferTime,
								OwnershipTime:   overlap.Overlap.OwnershipTime,
								ToList:          overlap.Overlap.ToList,
								FromList:        overlap.Overlap.FromList,
								InitiatedByList: overlap.Overlap.InitiatedByList,

								ApprovalIdList:         overlap.Overlap.ApprovalIdList,
								AmountTrackerIdList:    overlap.Overlap.AmountTrackerIdList,
								ChallengeTrackerIdList: overlap.Overlap.ChallengeTrackerIdList,

								//Appended for future lookups (not involved in overlap logic)
								PermanentlyPermittedTimes: permanentlyPermittedTimes,
								PermanentlyForbiddenTimes: permanentlyForbiddenTimes,
								ArbitraryValue:            newArbValue,
							})
						}
					}
				}
			}
		}

	}

	//It is first match only, so we can do this
	//To help with determinism in comparing later, we sort by badge ID
	//Thanks ChatGPT
	returnArr := []*types.UniversalPermissionDetails{}
	for _, handledItem := range handled {
		idxToInsert := 0
		for idxToInsert < len(returnArr) && handledItem.BadgeId.Start.GT(returnArr[idxToInsert].BadgeId.Start) {
			idxToInsert++
		}

		returnArr = append(returnArr, nil)
		copy(returnArr[idxToInsert+1:], returnArr[idxToInsert:])
		returnArr[idxToInsert] = handledItem
	}

	return handled
}
func (k Keeper) GetDetailsToCheck(ctx sdk.Context, collection *types.BadgeCollection, oldApprovals []*types.CollectionApproval, newApprovals []*types.CollectionApproval) ([]*types.UniversalPermissionDetails, error) {
	if !IsStandardBalances(collection) && newApprovals != nil && len(newApprovals) > 0 {
		return nil, sdkerrors.Wrapf(ErrWrongBalancesType, "collection %s does not have standard balances", collection.CollectionId)
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
	oldTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(ctx, x, []interface{}{oldApprovals})

	newTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(ctx, y, []interface{}{newApprovals})

	detailsToCheck, err := GetUpdateCombinationsToCheck(ctx, oldTimelineFirstMatches, newTimelineFirstMatches, []*types.CollectionApproval{}, func(ctx sdk.Context, oldValue interface{}, newValue interface{}) ([]*types.UniversalPermissionDetails, error) {
		//This is a little different from the other functions because it is not first match only

		//Expand all collection approved transfers so that they are manipulated according to options and approvalCriteria / allowedCombinations are len 1
		oldApprovals := oldValue.([]*types.CollectionApproval)
		newApprovals := newValue.([]*types.CollectionApproval)

		//Step 1: Merge so we get approvalCriteria arrays of proper length such that it is first match and each (to, from, init, time, ids, ownershipTimes) is only seen once
		//Step 2: Compare as we had previously

		//Step 1:
		oldApprovalsCasted, err := k.CastCollectionApprovalToUniversalPermission(ctx, oldApprovals)
		if err != nil {
			return nil, err
		}
		firstMatchesForOld := GetFirstMatchOnlyWithApprovalCriteria(ctx, oldApprovalsCasted)

		newApprovalsCasted, err := k.CastCollectionApprovalToUniversalPermission(ctx, newApprovals)
		if err != nil {
			return nil, err
		}
		firstMatchesForNew := GetFirstMatchOnlyWithApprovalCriteria(ctx, newApprovalsCasted)

		//Step 2:
		//For every badge, we need to check if the new provided value is different in any way from the old value for each badge ID
		//The overlapObjects from GetOverlapsAndNonOverlaps will return which badge IDs overlap
		//Note this okay since we already converted everything to first match only in the previous step
		detailsToReturn := []*types.UniversalPermissionDetails{}
		overlapObjects, inOldButNotNew, inNewButNotOld := types.GetOverlapsAndNonOverlaps(ctx, firstMatchesForOld, firstMatchesForNew)
		for _, overlapObject := range overlapObjects {
			overlap := overlapObject.Overlap
			oldDetails := overlapObject.FirstDetails
			newDetails := overlapObject.SecondDetails
			different := false
			if (oldDetails.ArbitraryValue == nil && newDetails.ArbitraryValue != nil) || (oldDetails.ArbitraryValue != nil && newDetails.ArbitraryValue == nil) {
				different = true
			} else {
				oldArbVal := oldDetails.ArbitraryValue.([]*ApprovalCriteriaWithIsApproved)
				newArbVal := newDetails.ArbitraryValue.([]*ApprovalCriteriaWithIsApproved)

				oldVal := oldArbVal
				newVal := newArbVal

				//TODO: Eventually we should make this more flexible instead of simply stringifying
				//For example, does it really matter what order they are in if approved? What about simply changing details that have no impact like customdata?
				//Or, if we have two empty approval details (no restrictions) and update to just one. That really does not matter.

				//Go one by one comparing old to new as flat array (if 2d array is empty we still treat it as an empty element
				if len(oldVal) != len(newVal) {
					different = true
				} else {

					//Decided against allowing flexible order here because if we use a linear match approahc, chanigng order might cause unexpected behavior
					//Even though, the user can choose which approval to select, it is still better to be consistent. Can change in the future though.
					//The only thing is I am not too sure how deterministic the GetFirstMatchOnlyWithApprovalCriteria function is.
					//TODO: Determine best path forward
					for i := 0; i < len(oldVal); i++ {
						oldApprovalCriteria := oldVal[i].ApprovalCriteria
						newApprovalCriteria := newVal[i].ApprovalCriteria
						if proto.MarshalTextString(oldApprovalCriteria) != proto.MarshalTextString(newApprovalCriteria) {
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
		return nil, err
	}

	return detailsToCheck, nil
}

func (k Keeper) ValidateCollectionApprovalsUpdate(ctx sdk.Context, collection *types.BadgeCollection, oldApprovals []*types.CollectionApproval, newApprovals []*types.CollectionApproval, CanUpdateCollectionApprovals []*types.CollectionApprovalPermission) error {
	detailsToCheck, err := k.GetDetailsToCheck(ctx, collection, oldApprovals, newApprovals)
	if err != nil {
		return err
	}

	err = k.CheckIfCollectionApprovalPermissionPermits(ctx, detailsToCheck, CanUpdateCollectionApprovals, "update collection approved transfers")
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) ValidateUserOutgoingApprovalsUpdate(ctx sdk.Context, collection *types.BadgeCollection, oldApprovals []*types.UserOutgoingApproval, newApprovals []*types.UserOutgoingApproval, CanUpdateCollectionApprovals []*types.UserOutgoingApprovalPermission, fromAddress string) error {
	old := types.CastOutgoingTransfersToCollectionTransfers(oldApprovals, fromAddress)
	new := types.CastOutgoingTransfersToCollectionTransfers(newApprovals, fromAddress)

	detailsToCheck, err := k.GetDetailsToCheck(ctx, collection, old, new)
	if err != nil {
		return err
	}

	err = k.CheckIfUserOutgoingApprovalPermissionPermits(ctx, detailsToCheck, CanUpdateCollectionApprovals, "update collection approved transfers")
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) ValidateUserIncomingApprovalsUpdate(ctx sdk.Context, collection *types.BadgeCollection, oldApprovals []*types.UserIncomingApproval, newApprovals []*types.UserIncomingApproval, CanUpdateCollectionApprovals []*types.UserIncomingApprovalPermission, toAddress string) error {
	old := types.CastIncomingTransfersToCollectionTransfers(oldApprovals, toAddress)
	new := types.CastIncomingTransfersToCollectionTransfers(newApprovals, toAddress)

	detailsToCheck, err := k.GetDetailsToCheck(ctx, collection, old, new)
	if err != nil {
		return err
	}

	err = k.CheckIfUserIncomingApprovalPermissionPermits(ctx, detailsToCheck, CanUpdateCollectionApprovals, "update collection approved transfers")
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) ValidateBadgeMetadataUpdate(ctx sdk.Context, oldBadgeMetadata []*types.BadgeMetadataTimeline, newBadgeMetadata []*types.BadgeMetadataTimeline, canUpdateBadgeMetadata []*types.TimedUpdateWithBadgeIdsPermission) error {
	oldTimes, oldValues := types.GetBadgeMetadataTimesAndValues(oldBadgeMetadata)
	oldTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(ctx, oldTimes, oldValues)

	newTimes, newValues := types.GetBadgeMetadataTimesAndValues(newBadgeMetadata)
	newTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(ctx, newTimes, newValues)

	detailsToCheck, err := GetUpdateCombinationsToCheck(ctx, oldTimelineFirstMatches, newTimelineFirstMatches, []*types.BadgeMetadata{}, func(ctx sdk.Context, oldValue interface{}, newValue interface{}) ([]*types.UniversalPermissionDetails, error) {
		//Cast to UniversalPermissionDetails for comaptibility with these overlap functions and get first matches only (i.e. first match for each badge ID)
		oldBadgeMetadata := oldValue.([]*types.BadgeMetadata)
		firstMatchesForOld := types.GetFirstMatchOnly(ctx, k.CastBadgeMetadataToUniversalPermission(oldBadgeMetadata))

		newBadgeMetadata := newValue.([]*types.BadgeMetadata)
		firstMatchesForNew := types.GetFirstMatchOnly(ctx, k.CastBadgeMetadataToUniversalPermission(newBadgeMetadata))

		//For every badge, we need to check if the new provided value is different in any way from the old value for each badge ID
		//The overlapObjects from GetOverlapsAndNonOverlaps will return which badge IDs overlap
		detailsToReturn := []*types.UniversalPermissionDetails{}
		overlapObjects, inOldButNotNew, inNewButNotOld := types.GetOverlapsAndNonOverlaps(ctx, firstMatchesForOld, firstMatchesForNew)
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

	err = k.CheckIfTimedUpdateWithBadgeIdsPermissionPermits(ctx, detailsToCheck, canUpdateBadgeMetadata, "update badge metadata")
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) ValidateCollectionMetadataUpdate(ctx sdk.Context, oldCollectionMetadata []*types.CollectionMetadataTimeline, newCollectionMetadata []*types.CollectionMetadataTimeline, canUpdateCollectionMetadata []*types.TimedUpdatePermission) error {
	oldTimes, oldValues := types.GetCollectionMetadataTimesAndValues(oldCollectionMetadata)
	oldTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(ctx, oldTimes, oldValues)

	newTimes, newValues := types.GetCollectionMetadataTimesAndValues(newCollectionMetadata)
	newTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(ctx, newTimes, newValues)

	detailsToCheck, err := GetUpdateCombinationsToCheck(ctx, oldTimelineFirstMatches, newTimelineFirstMatches, &types.CollectionMetadata{}, func(ctx sdk.Context, oldValue interface{}, newValue interface{}) ([]*types.UniversalPermissionDetails, error) {
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

	err = k.CheckIfTimedUpdatePermissionPermits(ctx, detailsToCheck, canUpdateCollectionMetadata, "update collection metadata")
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) ValidateOffChainBalancesMetadataUpdate(ctx sdk.Context, collection *types.BadgeCollection, oldOffChainBalancesMetadata []*types.OffChainBalancesMetadataTimeline, newOffChainBalancesMetadata []*types.OffChainBalancesMetadataTimeline, canUpdateOffChainBalancesMetadata []*types.TimedUpdatePermission) error {
	if !IsOffChainBalances(collection) && !IsNonIndexedBalances(collection) {
		if len(oldOffChainBalancesMetadata) > 0 || len(newOffChainBalancesMetadata) > 0 {
			return sdkerrors.Wrapf(ErrWrongBalancesType, "off chain balances are being set but collection %s does not have off chain balances", collection.CollectionId)
		}
		return nil
	}

	oldTimes, oldValues := types.GetOffChainBalancesMetadataTimesAndValues(oldOffChainBalancesMetadata)
	oldTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(ctx, oldTimes, oldValues)

	newTimes, newValues := types.GetOffChainBalancesMetadataTimesAndValues(newOffChainBalancesMetadata)
	newTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(ctx, newTimes, newValues)

	detailsToCheck, err := GetUpdateCombinationsToCheck(ctx, oldTimelineFirstMatches, newTimelineFirstMatches, &types.OffChainBalancesMetadata{}, func(ctx sdk.Context, oldValue interface{}, newValue interface{}) ([]*types.UniversalPermissionDetails, error) {
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

	err = k.CheckIfTimedUpdatePermissionPermits(ctx, detailsToCheck, canUpdateOffChainBalancesMetadata, "update off chain balances metadata")
	if err != nil {
		return err
	}

	return nil
}

/** Everything below here is pretty standard because all we need to compare is primitive types **/

func GetUpdatedStringCombinations(ctx sdk.Context, oldValue interface{}, newValue interface{}) ([]*types.UniversalPermissionDetails, error) {
	x := []*types.UniversalPermissionDetails{}
	if (oldValue == nil && newValue != nil) || (oldValue != nil && newValue == nil) {
		x = append(x, &types.UniversalPermissionDetails{})
	} else if oldValue.(string) != newValue.(string) {
		x = append(x, &types.UniversalPermissionDetails{})
	}
	return x, nil
}

func GetUpdatedBoolCombinations(ctx sdk.Context, oldValue interface{}, newValue interface{}) ([]*types.UniversalPermissionDetails, error) {
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
	oldTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(ctx, oldTimes, oldValues)

	newTimes, newValues := types.GetManagerTimesAndValues(newManager)
	newTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(ctx, newTimes, newValues)

	updatedTimelineTimes, err := GetUpdateCombinationsToCheck(ctx, oldTimelineFirstMatches, newTimelineFirstMatches, "", GetUpdatedStringCombinations)
	if err != nil {
		return err
	}

	if err = k.CheckIfTimedUpdatePermissionPermits(ctx, updatedTimelineTimes, canUpdateManager, "update manager"); err != nil {
		return err
	}

	return nil
}

func (k Keeper) ValidateCustomDataUpdate(ctx sdk.Context, oldCustomData []*types.CustomDataTimeline, newCustomData []*types.CustomDataTimeline, canUpdateCustomData []*types.TimedUpdatePermission) error {
	oldTimes, oldValues := types.GetCustomDataTimesAndValues(oldCustomData)
	oldTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(ctx, oldTimes, oldValues)

	newTimes, newValues := types.GetCustomDataTimesAndValues(newCustomData)
	newTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(ctx, newTimes, newValues)

	updatedTimelineTimes, err := GetUpdateCombinationsToCheck(ctx, oldTimelineFirstMatches, newTimelineFirstMatches, "", GetUpdatedStringCombinations)
	if err != nil {
		return err
	}

	if err = k.CheckIfTimedUpdatePermissionPermits(ctx, updatedTimelineTimes, canUpdateCustomData, "update custom data"); err != nil {
		return err
	}

	return nil
}

func (k Keeper) ValidateStandardsUpdate(ctx sdk.Context, oldStandards []*types.StandardsTimeline, newStandards []*types.StandardsTimeline, canUpdateStandards []*types.TimedUpdatePermission) error {
	oldTimes, oldValues := types.GetStandardsTimesAndValues(oldStandards)
	oldTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(ctx, oldTimes, oldValues)

	newTimes, newValues := types.GetStandardsTimesAndValues(newStandards)
	newTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(ctx, newTimes, newValues)

	updatedTimelineTimes, err := GetUpdateCombinationsToCheck(ctx, oldTimelineFirstMatches, newTimelineFirstMatches, []string{}, func(ctx sdk.Context, oldValue interface{}, newValue interface{}) ([]*types.UniversalPermissionDetails, error) {
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

	if err = k.CheckIfTimedUpdatePermissionPermits(ctx, updatedTimelineTimes, canUpdateStandards, "update standards"); err != nil {
		return err
	}

	return nil
}

func (k Keeper) ValidateIsArchivedUpdate(ctx sdk.Context, oldIsArchived []*types.IsArchivedTimeline, newIsArchived []*types.IsArchivedTimeline, canUpdateIsArchived []*types.TimedUpdatePermission) error {
	oldTimes, oldValues := types.GetIsArchivedTimesAndValues(oldIsArchived)
	oldTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(ctx, oldTimes, oldValues)

	newTimes, newValues := types.GetIsArchivedTimesAndValues(newIsArchived)
	newTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(ctx, newTimes, newValues)

	updatedTimelineTimes, err := GetUpdateCombinationsToCheck(ctx, oldTimelineFirstMatches, newTimelineFirstMatches, false, GetUpdatedBoolCombinations)
	if err != nil {
		return err
	}

	if err = k.CheckIfTimedUpdatePermissionPermits(ctx, updatedTimelineTimes, canUpdateIsArchived, "update is archived"); err != nil {
		return err
	}

	return nil
}
