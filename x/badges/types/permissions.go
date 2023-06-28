package types

import (
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type UniversalCombination struct {
	TimelineTimesOptions *ValueOptions
	FromMappingIdOptions *ValueOptions
	ToMappingIdOptions *ValueOptions
	InitiatedByMappingIdOptions *ValueOptions
	TransferTimesOptions *ValueOptions
	BadgeIdsOptions *ValueOptions

	PermittedTimesOptions *ValueOptions
	ForbiddenTimesOptions *ValueOptions
}

type UniversalPermission struct {
	DefaultValues *UniversalDefaultValues
	Combinations []*UniversalCombination
}

type UniversalDefaultValues struct {
	BadgeIds []*IdRange
	TimelineTimes []*IdRange
	TransferTimes []*IdRange
	ToMappingId string
	FromMappingId string
	InitiatedByMappingId string

	PermittedTimes []*IdRange
	ForbiddenTimes []*IdRange

	UsesBadgeIds bool
	UsesTimelineTimes bool
	UsesTransferTimes bool
	UsesToMappingId bool
	UsesFromMappingId bool
	UsesInitiatedByMappingId bool

	ArbitraryValue interface{}
}

type UniversalPermissionDetails struct {
	BadgeId *IdRange
	TimelineTime *IdRange
	TransferTime *IdRange
	ToMappingId string
	FromMappingId string
	InitiatedByMappingId string

	PermittedTimes []*IdRange
	ForbiddenTimes []*IdRange

	ArbitraryValue interface{}
}

type Overlap struct {
	Overlap *UniversalPermissionDetails
	FirstDetails *UniversalPermissionDetails
	SecondDetails *UniversalPermissionDetails
}

func GetOverlapsAndNonOverlaps(firstDetails []*UniversalPermissionDetails, secondDetails []*UniversalPermissionDetails) ([]*Overlap, []*UniversalPermissionDetails, []*UniversalPermissionDetails) {
	inOldButNotNew := make([]*UniversalPermissionDetails, len(firstDetails))
	copy(inOldButNotNew, firstDetails)
	inNewButNotOld := make([]*UniversalPermissionDetails, len(secondDetails))
	copy(inNewButNotOld, secondDetails)

	allOverlaps := []*Overlap{}

	for _, oldPermission := range firstDetails {
		for _, newPermission := range secondDetails {
			_, overlaps := UniversalRemoveOverlaps(newPermission, oldPermission)
			
			for _, overlap := range overlaps {
				allOverlaps = append(allOverlaps, &Overlap{
					Overlap: overlap,
					FirstDetails: oldPermission,
					SecondDetails: newPermission,
				})
				inOldButNotNew = UniversalRemoveOverlapFromValues(overlap, inOldButNotNew)
				inNewButNotOld = UniversalRemoveOverlapFromValues(overlap, inNewButNotOld)
			}
		}
	}

	return allOverlaps, inOldButNotNew, inNewButNotOld
}

func UniversalRemoveOverlapFromValues(handled *UniversalPermissionDetails, valuesToCheck []*UniversalPermissionDetails) ([]*UniversalPermissionDetails) {
	newValuesToCheck := []*UniversalPermissionDetails{}
	for _, valueToCheck := range valuesToCheck {
		remaining, _ := UniversalRemoveOverlaps(handled, valueToCheck)
		for _, val := range remaining {
			newValuesToCheck = append(newValuesToCheck, val)
		}
	}

	return newValuesToCheck
}

func UniversalRemoveOverlaps(handled *UniversalPermissionDetails, valueToCheck *UniversalPermissionDetails) ([]*UniversalPermissionDetails, []*UniversalPermissionDetails) {
	timelineTimesAfterRemoval, removedTimelineTimes := RemoveIdsFromIdRange(handled.TimelineTime, valueToCheck.TimelineTime)
	badgesAfterRemoval, removedBadges := RemoveIdsFromIdRange(handled.BadgeId, valueToCheck.BadgeId)
	transferTimesAfterRemoval, removedTransferTimes := RemoveIdsFromIdRange(handled.TransferTime, valueToCheck.TransferTime)

	//Due to lack of determinism over who owns what badge in each address mapping, we only check if handled.MappingId === valueToCheck.MappingId
	//If they are the same, we consider it completely overlapping
	//If they are different, we consider it not overlapping at all (even though it might be somehwat in reality)
	//As a result, we don't have mapping contents after removal like we do with the ID ranges. It is just all or nothing
	toMappingRemoved := handled.ToMappingId == valueToCheck.ToMappingId
	fromMappingRemoved := handled.FromMappingId == valueToCheck.FromMappingId
	initiatedByMappingRemoved := handled.InitiatedByMappingId == valueToCheck.InitiatedByMappingId

	//Approach: Iterate through each field one by one. Attempt to remove the overlap. We'll call each field by an ID number corresponding to its order
	//					Order doesn't matter as long as all fields are handled
	//          For each field N, we have the following the cases:
	//            1. For anything remaining after removal for field N is attempted (i.e. the stuff that does not overlap), we need to add 
	//               it to the returned array with (0 to N-1) fields filled with removed values and (N+1 to end) fields filled with the original values
	//               
	// 						   We only use the removed values of fields 0 to N - 1 because we already handled the other fields (via this step in previous iterations) 
	//               and we don't want to double count.
	//							 Ex: [0: {1 to 10}, 1: {1 to 10}, 2: {1 to 10}] and we are removing [0: {1 to 5}, 1: {1 to 5}, 2: {1 to 5}]
	//							  	 Let's say we are on field 1. We would add [0: {1 to 5}, 1: {6 to 10}, 2: {1 to 10}] to the returned array
	//						2. If we have removed anything at all, we need to continue to test field N + 1 (i.e. the next field) for overlap
	//							 This is because we have not yet handled the cases for values which overlap with field N and field N + 1
	//					  3. If we have not removed anything, we add the original value as outlined in 1) but we do not need to continue to test field N + 1
	//							 because there are no cases unhandled now where values overlap with field N and field N + 1 becuase nothing overlaps with N.
	//							 If we do end up with this case, it means we end up with the original values because to overlap, it needs to overlap with all fields
	//
	//							 We optimize step 3) by checking right away if something does not overlap with some field. If it does not overlap with some field,
	//							 we can just add the original values and be done with it. If it does overlap with all fields, we need to execute the algorithm

	remaining := []*UniversalPermissionDetails{}

	//If some field does not overlap, we simply end up with the original values because it is only considered an overlap if all fields overlap.
	//The function would work fine without this but it makes it more efficient and less complicated because it will not get broken down further
	if len(removedTimelineTimes) == 0 || len(removedBadges) == 0 || len(removedTransferTimes) == 0 || !toMappingRemoved || !fromMappingRemoved || !initiatedByMappingRemoved {
		remaining = append(remaining, valueToCheck)
		return remaining, []*UniversalPermissionDetails{}
	}

	for _, timelineTimeAfterRemoval := range timelineTimesAfterRemoval {
		remaining = append(remaining, &UniversalPermissionDetails{
			TimelineTime: timelineTimeAfterRemoval,
			BadgeId: valueToCheck.BadgeId,
			TransferTime: valueToCheck.TransferTime,
			ToMappingId: valueToCheck.ToMappingId,
			FromMappingId: valueToCheck.FromMappingId,
			InitiatedByMappingId: valueToCheck.InitiatedByMappingId,
		})
	}

	for _, badgeAfterRemoval := range badgesAfterRemoval {
		remaining = append(remaining, &UniversalPermissionDetails{
			TimelineTime: removedTimelineTimes[0], //We know there is only one because there can only be one interesection between two ranges
			BadgeId: badgeAfterRemoval,
			TransferTime: valueToCheck.TransferTime,
			ToMappingId: valueToCheck.ToMappingId,
			FromMappingId: valueToCheck.FromMappingId,
			InitiatedByMappingId: valueToCheck.InitiatedByMappingId,
		})
	}

	for _, transferTimeAfterRemoval := range transferTimesAfterRemoval {
		remaining = append(remaining, &UniversalPermissionDetails{
			TimelineTime: removedTimelineTimes[0], //We know there is only one because there can only be one interesection between two ranges
			BadgeId: removedBadges[0], //We know there is only one because there can only be one interesection between two ranges
			TransferTime: transferTimeAfterRemoval,
			ToMappingId: valueToCheck.ToMappingId,
			FromMappingId: valueToCheck.FromMappingId,
			InitiatedByMappingId: valueToCheck.InitiatedByMappingId,
		})
	}

			
	//For the mapping IDs, it is an all or nothing system. Either we remove the whole mapping ID or we don't
	//If we reach here, we know that toMappingRemoved, fromMappingRemoved, and initiatedByMappingRemoved are all true
	//Thus, there is nothing left to add here

	removedDetails := []*UniversalPermissionDetails{}
	for _, removedTimelineTime := range removedTimelineTimes {
		for _, removedBadge := range removedBadges {
			for _, removedTransferTime := range removedTransferTimes {
				removedDetails = append(removedDetails, &UniversalPermissionDetails{
					TimelineTime: removedTimelineTime,
					BadgeId: removedBadge,
					TransferTime: removedTransferTime,
					ToMappingId: valueToCheck.ToMappingId,
					FromMappingId: valueToCheck.FromMappingId,
					InitiatedByMappingId: valueToCheck.InitiatedByMappingId,
				})
			}
		}
	}
	
	return remaining, removedDetails
}

func GetIdRangesWithOptions(ranges []*IdRange, options *ValueOptions, uses bool) []*IdRange {
	if !uses {
		ranges = []*IdRange{&IdRange{Start: sdk.NewUint(math.MaxUint64), End: sdk.NewUint(math.MaxUint64)}} //dummy range
		return ranges
	}
	
	if options == nil {
		return ranges
	}

	if options.AllValues {
		ranges = []*IdRange{&IdRange{Start: sdk.NewUint(0), End: sdk.NewUint(math.MaxUint64)}}
	}

	if options.InvertDefault {
		ranges = InvertIdRanges(ranges, sdk.NewUint(math.MaxUint64))
	}

	if options.NoValues {
		ranges = []*IdRange{}
	}

	return ranges
}

func GetMappingWithOptions(mappingId string, options *ValueOptions, uses bool) string {
	if !uses {
		mappingId = "" //dummy mappingId
	}
	
	if options == nil {
		return mappingId
	}

	if options.AllValues {
		mappingId = "All"
	}

	if options.InvertDefault {
		if mappingId == "All" {
			mappingId = "None"
		} else if mappingId == "None" {
			mappingId = "All"
		} else {
			mappingId = "!" + mappingId
		}
	}

	if options.NoValues {
		mappingId = "None"
	}

	return mappingId
}

func GetFirstMatchOnly(permissions []*UniversalPermission) ([]*UniversalPermissionDetails) {
	handled := []*UniversalPermissionDetails{}
	for _, permission := range permissions {
		for _, combination := range permission.Combinations {
			badgeIds := GetIdRangesWithOptions(permission.DefaultValues.BadgeIds, combination.BadgeIdsOptions, permission.DefaultValues.UsesBadgeIds)
			timelineTimes := GetIdRangesWithOptions(permission.DefaultValues.TimelineTimes, combination.TimelineTimesOptions, permission.DefaultValues.UsesTimelineTimes)
			transferTimes := GetIdRangesWithOptions(permission.DefaultValues.TransferTimes, combination.TransferTimesOptions, permission.DefaultValues.UsesTransferTimes)
			permittedTimes := GetIdRangesWithOptions(permission.DefaultValues.PermittedTimes, combination.PermittedTimesOptions, true)
			forbiddenTimes := GetIdRangesWithOptions(permission.DefaultValues.ForbiddenTimes, combination.ForbiddenTimesOptions, true)
			arbitraryValue := permission.DefaultValues.ArbitraryValue

			toMappingId := GetMappingWithOptions(permission.DefaultValues.ToMappingId, combination.ToMappingIdOptions, permission.DefaultValues.UsesToMappingId)
			fromMappingId := GetMappingWithOptions(permission.DefaultValues.FromMappingId, combination.FromMappingIdOptions, permission.DefaultValues.UsesFromMappingId)
			initiatedByMappingId := GetMappingWithOptions(permission.DefaultValues.InitiatedByMappingId, combination.InitiatedByMappingIdOptions, permission.DefaultValues.UsesInitiatedByMappingId)

			for _, badgeId := range badgeIds {
				for _, timelineTime := range timelineTimes {
					for _, transferTime := range transferTimes {
						brokenDown := []*UniversalPermissionDetails{
							&UniversalPermissionDetails{
								BadgeId: badgeId,
								TimelineTime: timelineTime,
								TransferTime: transferTime,
								ToMappingId: toMappingId,
								FromMappingId: fromMappingId,
								InitiatedByMappingId: initiatedByMappingId,
								PermittedTimes: permittedTimes,
								ForbiddenTimes: forbiddenTimes,
								ArbitraryValue: arbitraryValue,
							},
						}

						for _, handledPermission := range handled {
							newBrokenDown := []*UniversalPermissionDetails{}
							for _, broken := range brokenDown {
								remainingVals, _ := UniversalRemoveOverlaps(&UniversalPermissionDetails{
									TimelineTime: handledPermission.TimelineTime,
									BadgeId: handledPermission.BadgeId,
									TransferTime: handledPermission.TransferTime,
									ToMappingId: handledPermission.ToMappingId,
									FromMappingId: handledPermission.FromMappingId,
									InitiatedByMappingId: handledPermission.InitiatedByMappingId,
								}, &UniversalPermissionDetails{
									TimelineTime: broken.TimelineTime,
									BadgeId: broken.BadgeId,
									TransferTime: broken.TransferTime,
									ToMappingId: broken.ToMappingId,
									FromMappingId: broken.FromMappingId,
									InitiatedByMappingId: broken.InitiatedByMappingId,
								})
								for _, remaining := range remainingVals {
									newBrokenDown = append(newBrokenDown, &UniversalPermissionDetails{
										TimelineTime: remaining.TimelineTime,
										BadgeId: remaining.BadgeId,
										TransferTime: remaining.TransferTime,
										ToMappingId: remaining.ToMappingId,
										FromMappingId: remaining.FromMappingId,
										InitiatedByMappingId: remaining.InitiatedByMappingId,
										PermittedTimes: permittedTimes,
										ForbiddenTimes: forbiddenTimes,
										ArbitraryValue: arbitraryValue,
									})
								}
							}

							brokenDown = newBrokenDown
						}

						if len(brokenDown) > 0 {
							handled = append(handled, brokenDown...)
						}
					}
				}
			}
		}
	}

	return handled
}

func GetPermissionString(permission *UniversalPermissionDetails) string {
	str := "("
	if permission.BadgeId.Start != sdk.NewUint(math.MaxUint64) || permission.BadgeId.End != sdk.NewUint(math.MaxUint64) {
		str += "badgeId: " + permission.BadgeId.Start.String() + " "
	}

	if permission.TimelineTime.Start != sdk.NewUint(math.MaxUint64) || permission.TimelineTime.End != sdk.NewUint(math.MaxUint64) {
		str += "timelineTime: " + permission.TimelineTime.Start.String() + " "
	}

	if permission.TransferTime.Start != sdk.NewUint(math.MaxUint64) || permission.TransferTime.End != sdk.NewUint(math.MaxUint64) {
		str += "transferTime: " + permission.TransferTime.Start.String() + " "
	}

	if permission.ToMappingId != "" {
		str += "toMappingId: " + permission.ToMappingId + " "
	}

	if permission.FromMappingId != "" {
		str += "fromMappingId: " + permission.FromMappingId + " "
	}

	if permission.InitiatedByMappingId != "" {
		str += "initiatedByMappingId: " + permission.InitiatedByMappingId + " "
	}

	str += ") "

	return str
}


//IMPORTANT PRECONDITION: Must be first match only
func ValidateUniversalPermissionUpdate(oldPermissions []*UniversalPermissionDetails, newPermissions []*UniversalPermissionDetails) error {
	allOverlaps, inOldButNotNew, _ := GetOverlapsAndNonOverlaps(oldPermissions, newPermissions)  //we don't care about new not in old
	
	if len(inOldButNotNew) > 0 {
		errMsg := "permission "
		errMsg += GetPermissionString(inOldButNotNew[0])
		errMsg += "found in old permissions but not in new permissions"

		return sdkerrors.Wrapf(ErrInvalidPermissions, errMsg)
	}
	
	//For everywhere we detected an overlap, we need to check if the new permissions are valid (i.e. they only explicitly define more times and do not remove any)
	for _, overlapObj := range allOverlaps {
		oldPermission := overlapObj.FirstDetails
		newPermission := overlapObj.SecondDetails
		oldPermittedTimes := oldPermission.PermittedTimes
		oldForbiddenTimes := oldPermission.ForbiddenTimes
		newPermittedTimes := newPermission.PermittedTimes
		newForbiddenTimes := newPermission.ForbiddenTimes

		leftoverPermittedTimes, _ := RemoveIdRangeFromIdRange(newPermittedTimes, oldPermittedTimes)
		leftoverForbiddenTimes, _ := RemoveIdRangeFromIdRange(newForbiddenTimes, oldForbiddenTimes)

		if len(leftoverPermittedTimes) > 0 || len(leftoverForbiddenTimes) > 0 {
			errMsg := "permission "
			errMsg += GetPermissionString(oldPermission)
			errMsg += "found in both new and old permissions but "
			if len(leftoverPermittedTimes) > 0 {
				errMsg += "previously explicitly allowed the times ( "
				for _, oldPermittedTime := range leftoverPermittedTimes {
					errMsg += oldPermittedTime.Start.String() + "-" + oldPermittedTime.End.String() + " "
				}
				errMsg += ") which are now set to disallowed"
			}
			if len(leftoverForbiddenTimes) > 0 && len(leftoverPermittedTimes) > 0 {
				errMsg += " and"
			}
			if len(leftoverForbiddenTimes) > 0 {
				errMsg += " previously explicitly disallowed the times ( "
				for _, oldForbiddenTime := range leftoverForbiddenTimes {
					errMsg += oldForbiddenTime.Start.String() + "-" + oldForbiddenTime.End.String() + " "
				}
				errMsg += ") which are now set to allowed."
			}
		
			return sdkerrors.Wrapf(ErrInvalidPermissions, errMsg)
		}
	}

	//Note we do not care about inNewButNotOld because it is fine to add new permissions that were not explicitly allowed/disallowed before

	return nil
}


func CastActionWithBadgeIdsPermissionToUniversalPermission(actionWithBadgeIdsPermission []*ActionWithBadgeIdsPermission) []*UniversalPermission {
	castedPermissions := []*UniversalPermission{}
	for _, actionWithBadgeIdsPermission := range actionWithBadgeIdsPermission {
		castedCombinations := []*UniversalCombination{}
		for _, actionWithBadgeIdsCombination := range actionWithBadgeIdsPermission.Combinations {
			castedCombinations = append(castedCombinations, &UniversalCombination{
				BadgeIdsOptions: actionWithBadgeIdsCombination.BadgeIdsOptions,
				PermittedTimesOptions: actionWithBadgeIdsCombination.PermittedTimesOptions,
				ForbiddenTimesOptions: actionWithBadgeIdsCombination.ForbiddenTimesOptions,
			})
		}

		castedPermissions = append(castedPermissions, &UniversalPermission{
			DefaultValues: &UniversalDefaultValues{
				BadgeIds: actionWithBadgeIdsPermission.DefaultValues.BadgeIds,
				UsesBadgeIds: true,
				PermittedTimes: actionWithBadgeIdsPermission.DefaultValues.PermittedTimes,
				ForbiddenTimes: actionWithBadgeIdsPermission.DefaultValues.ForbiddenTimes,
			},
			Combinations: castedCombinations,
		})
	}
	return castedPermissions
}

func ValidateActionWithBadgeIdsPermissionUpdate(oldPermissions []*ActionWithBadgeIdsPermission, newPermissions []*ActionWithBadgeIdsPermission) error {
	if err := ValidateActionWithBadgeIdsPermission(oldPermissions); err != nil {
		return err
	}

	if err := ValidateActionWithBadgeIdsPermission(newPermissions); err != nil {
		return err
	}

	castedOldPermissions := CastActionWithBadgeIdsPermissionToUniversalPermission(oldPermissions)
	castedNewPermissions := CastActionWithBadgeIdsPermissionToUniversalPermission(newPermissions)
	
	err := ValidateUniversalPermissionUpdate(GetFirstMatchOnly(castedOldPermissions), GetFirstMatchOnly(castedNewPermissions))
	if err != nil {
		return err
	}

	return nil
}

func CastTimedUpdatePermissionToUniversalPermission(timedUpdatePermission []*TimedUpdatePermission) []*UniversalPermission {
	castedPermissions := []*UniversalPermission{}
	for _, timedUpdatePermission := range timedUpdatePermission {
		castedCombinations := []*UniversalCombination{}
		for _, timedUpdateCombination := range timedUpdatePermission.Combinations {
			castedCombinations = append(castedCombinations, &UniversalCombination{
				PermittedTimesOptions: timedUpdateCombination.PermittedTimesOptions,
				ForbiddenTimesOptions: timedUpdateCombination.ForbiddenTimesOptions,
				TimelineTimesOptions: timedUpdateCombination.TimelineTimesOptions,
			})
		}

		castedPermissions = append(castedPermissions, &UniversalPermission{
			DefaultValues: &UniversalDefaultValues{
				TimelineTimes: timedUpdatePermission.DefaultValues.TimelineTimes,
				UsesTimelineTimes: true,
				PermittedTimes: timedUpdatePermission.DefaultValues.PermittedTimes,
				ForbiddenTimes: timedUpdatePermission.DefaultValues.ForbiddenTimes,
			},
			Combinations: castedCombinations,
		})
	}
	return castedPermissions
}

func ValidateTimedUpdatePermissionUpdate(oldPermissions []*TimedUpdatePermission, newPermissions []*TimedUpdatePermission) error {
	if err := ValidateTimedUpdatePermission(oldPermissions); err != nil {
		return err
	}

	if err := ValidateTimedUpdatePermission(newPermissions); err != nil {
		return err
	}

	castedOldPermissions := CastTimedUpdatePermissionToUniversalPermission(oldPermissions)
	castedNewPermissions := CastTimedUpdatePermissionToUniversalPermission(newPermissions)

	err := ValidateUniversalPermissionUpdate(GetFirstMatchOnly(castedOldPermissions), GetFirstMatchOnly(castedNewPermissions))
	if err != nil {
		return err
	}

	return nil
}

func CastTimedUpdateWithBadgeIdsPermissionToUniversalPermission(timedUpdateWithBadgeIdsPermission []*TimedUpdateWithBadgeIdsPermission) []*UniversalPermission {
	castedPermissions := []*UniversalPermission{}
	for _, timedUpdateWithBadgeIdsPermission := range timedUpdateWithBadgeIdsPermission {
		castedCombinations := []*UniversalCombination{}
		for _, timedUpdateWithBadgeIdsCombination := range timedUpdateWithBadgeIdsPermission.Combinations {
			castedCombinations = append(castedCombinations, &UniversalCombination{
				BadgeIdsOptions: timedUpdateWithBadgeIdsCombination.BadgeIdsOptions,
				PermittedTimesOptions: timedUpdateWithBadgeIdsCombination.PermittedTimesOptions,
				ForbiddenTimesOptions: timedUpdateWithBadgeIdsCombination.ForbiddenTimesOptions,
				TimelineTimesOptions: timedUpdateWithBadgeIdsCombination.TimelineTimesOptions,
			})
		}

		castedPermissions = append(castedPermissions, &UniversalPermission{
			DefaultValues: &UniversalDefaultValues{
				TimelineTimes: timedUpdateWithBadgeIdsPermission.DefaultValues.TimelineTimes,
				BadgeIds: timedUpdateWithBadgeIdsPermission.DefaultValues.BadgeIds,
				UsesTimelineTimes: true,
				UsesBadgeIds: true,
				PermittedTimes: timedUpdateWithBadgeIdsPermission.DefaultValues.PermittedTimes,
				ForbiddenTimes: timedUpdateWithBadgeIdsPermission.DefaultValues.ForbiddenTimes,
			},
			Combinations: castedCombinations,
		})
	}
	return castedPermissions
}

func ValidateTimedUpdateWithBadgeIdsPermissionUpdate(oldPermissions []*TimedUpdateWithBadgeIdsPermission, newPermissions []*TimedUpdateWithBadgeIdsPermission) error {
	if err := ValidateTimedUpdateWithBadgeIdsPermission(oldPermissions); err != nil {
		return err
	}

	if err := ValidateTimedUpdateWithBadgeIdsPermission(newPermissions); err != nil {
		return err
	}

	castedOldPermissions := CastTimedUpdateWithBadgeIdsPermissionToUniversalPermission(oldPermissions)
	castedNewPermissions := CastTimedUpdateWithBadgeIdsPermissionToUniversalPermission(newPermissions)

	err := ValidateUniversalPermissionUpdate(GetFirstMatchOnly(castedOldPermissions), GetFirstMatchOnly(castedNewPermissions))
	if err != nil {
		return err
	}

	return nil
}

func CastCollectionApprovedTransferPermissionToUniversalPermission(collectionUpdatePermission []*CollectionApprovedTransferPermission) []*UniversalPermission {
	castedPermissions := []*UniversalPermission{}
	for _, collectionUpdatePermission := range collectionUpdatePermission {
		castedCombinations := []*UniversalCombination{}
		for _, collectionUpdateCombination := range collectionUpdatePermission.Combinations {
			castedCombinations = append(castedCombinations, &UniversalCombination{
				PermittedTimesOptions: collectionUpdateCombination.PermittedTimesOptions,
				ForbiddenTimesOptions: collectionUpdateCombination.ForbiddenTimesOptions,
				TimelineTimesOptions: collectionUpdateCombination.TimelineTimesOptions,
				TransferTimesOptions: collectionUpdateCombination.TransferTimesOptions,
				ToMappingIdOptions: collectionUpdateCombination.ToMappingIdOptions,
				FromMappingIdOptions: collectionUpdateCombination.FromMappingIdOptions,
				InitiatedByMappingIdOptions: collectionUpdateCombination.InitiatedByMappingIdOptions,
				BadgeIdsOptions: collectionUpdateCombination.BadgeIdsOptions,
			})
		}

		castedPermissions = append(castedPermissions, &UniversalPermission{
			DefaultValues: &UniversalDefaultValues{
				TimelineTimes: collectionUpdatePermission.DefaultValues.TimelineTimes,
				TransferTimes: collectionUpdatePermission.DefaultValues.TransferTimes,
				ToMappingId: collectionUpdatePermission.DefaultValues.ToMappingId,
				FromMappingId: collectionUpdatePermission.DefaultValues.FromMappingId,
				InitiatedByMappingId: collectionUpdatePermission.DefaultValues.InitiatedByMappingId,
				BadgeIds: collectionUpdatePermission.DefaultValues.BadgeIds,
				UsesBadgeIds: true,
				UsesTimelineTimes: true,
				UsesTransferTimes: true,
				UsesToMappingId: true,
				UsesFromMappingId: true,
				UsesInitiatedByMappingId: true,
				PermittedTimes: collectionUpdatePermission.DefaultValues.PermittedTimes,
				ForbiddenTimes: collectionUpdatePermission.DefaultValues.ForbiddenTimes,
			},
			Combinations: castedCombinations,
		})
	}
	return castedPermissions
}

func ValidateCollectionApprovedTransferPermissionsUpdate(oldPermissions []*CollectionApprovedTransferPermission, newPermissions []*CollectionApprovedTransferPermission) error {
	if err := ValidateCollectionApprovedTransferPermissions(oldPermissions); err != nil {
		return err
	}

	if err := ValidateCollectionApprovedTransferPermissions(newPermissions); err != nil {
		return err
	}

	castedOldPermissions := CastCollectionApprovedTransferPermissionToUniversalPermission(oldPermissions)
	castedNewPermissions := CastCollectionApprovedTransferPermissionToUniversalPermission(newPermissions)

	err := ValidateUniversalPermissionUpdate(GetFirstMatchOnly(castedOldPermissions), GetFirstMatchOnly(castedNewPermissions))
	if err != nil {
		return err
	}

	return nil
}

func CastActionPermissionToUniversalPermission(actionPermission []*ActionPermission) []*UniversalPermission {
	castedPermissions := []*UniversalPermission{}
	for _, actionPermission := range actionPermission {
		castedCombinations := []*UniversalCombination{}
		for _, actionCombination := range actionPermission.Combinations {
			castedCombinations = append(castedCombinations, &UniversalCombination{
				PermittedTimesOptions: actionCombination.PermittedTimesOptions,
				ForbiddenTimesOptions: actionCombination.ForbiddenTimesOptions,
			})
		}

		castedPermissions = append(castedPermissions, &UniversalPermission{
			DefaultValues: &UniversalDefaultValues{
				PermittedTimes: actionPermission.DefaultValues.PermittedTimes,
				ForbiddenTimes: actionPermission.DefaultValues.ForbiddenTimes,
			},
			Combinations: castedCombinations,
		})
	}
	return castedPermissions
}

func ValidateActionPermissionUpdate(oldPermissions []*ActionPermission, newPermissions []*ActionPermission) error {
	if err := ValidateActionPermission(oldPermissions); err != nil {
		return err
	}

	if err := ValidateActionPermission(newPermissions); err != nil {
		return err
	}

	castedOldPermissions := CastActionPermissionToUniversalPermission(oldPermissions)
	castedNewPermissions := CastActionPermissionToUniversalPermission(newPermissions)

	err := ValidateUniversalPermissionUpdate(GetFirstMatchOnly(castedOldPermissions), GetFirstMatchOnly(castedNewPermissions))
	if err != nil {
		return err
	}

	return nil
}

func CastUserApprovedTransferPermissionToUniversalPermission(UserApprovedTransferPermission []*UserApprovedTransferPermission) []*UniversalPermission {
	castedPermissions := []*UniversalPermission{}
	for _, UserApprovedTransferPermission := range UserApprovedTransferPermission {
		castedCombinations := []*UniversalCombination{}
		for _, UserApprovedOutgoingTransferCombination := range UserApprovedTransferPermission.Combinations {
			castedCombinations = append(castedCombinations, &UniversalCombination{
				BadgeIdsOptions: UserApprovedOutgoingTransferCombination.BadgeIdsOptions,
				PermittedTimesOptions: UserApprovedOutgoingTransferCombination.PermittedTimesOptions,
				ForbiddenTimesOptions: UserApprovedOutgoingTransferCombination.ForbiddenTimesOptions,
				TimelineTimesOptions: UserApprovedOutgoingTransferCombination.TimelineTimesOptions,
				TransferTimesOptions: UserApprovedOutgoingTransferCombination.TransferTimesOptions,
				ToMappingIdOptions: UserApprovedOutgoingTransferCombination.ToMappingIdOptions,
				InitiatedByMappingIdOptions: UserApprovedOutgoingTransferCombination.InitiatedByMappingIdOptions,
			})
		}

		castedPermissions = append(castedPermissions, &UniversalPermission{
			DefaultValues: &UniversalDefaultValues{
				BadgeIds: UserApprovedTransferPermission.DefaultValues.BadgeIds,
				TimelineTimes: UserApprovedTransferPermission.DefaultValues.TimelineTimes,
				TransferTimes: UserApprovedTransferPermission.DefaultValues.TransferTimes,
				ToMappingId: UserApprovedTransferPermission.DefaultValues.ToMappingId,
				InitiatedByMappingId: UserApprovedTransferPermission.DefaultValues.InitiatedByMappingId,
				UsesBadgeIds: true,
				UsesTimelineTimes: true,
				UsesTransferTimes: true,
				UsesToMappingId: true,
				UsesInitiatedByMappingId: true,
				PermittedTimes: UserApprovedTransferPermission.DefaultValues.PermittedTimes,
				ForbiddenTimes: UserApprovedTransferPermission.DefaultValues.ForbiddenTimes,
			},
			Combinations: castedCombinations,
		})
	}
	return castedPermissions
}

func ValidateUserApprovedTransferPermissionsUpdate(oldPermissions []*UserApprovedTransferPermission, newPermissions []*UserApprovedTransferPermission) error {
	if err := ValidateUserApprovedTransferPermissions(oldPermissions); err != nil {
		return err
	}

	if err := ValidateUserApprovedTransferPermissions(newPermissions); err != nil {
		return err
	}

	castedOldPermissions := CastUserApprovedTransferPermissionToUniversalPermission(oldPermissions)
	castedNewPermissions := CastUserApprovedTransferPermissionToUniversalPermission(newPermissions)

	err := ValidateUniversalPermissionUpdate(GetFirstMatchOnly(castedOldPermissions), GetFirstMatchOnly(castedNewPermissions))
	if err != nil {
		return err
	}

	return nil
}

func ValidateUserPermissionsUpdate(oldPermissions *UserPermissions, newPermissions *UserPermissions, canBeNil bool) error {
	if err := ValidateUserPermissions(oldPermissions, canBeNil); err != nil {
		return err
	}

	if err := ValidateUserPermissions(newPermissions, canBeNil); err != nil {
		return err
	}

	if oldPermissions.CanUpdateApprovedIncomingTransfers != nil && newPermissions.CanUpdateApprovedIncomingTransfers != nil {
		if err := ValidateUserApprovedTransferPermissionsUpdate(oldPermissions.CanUpdateApprovedIncomingTransfers, newPermissions.CanUpdateApprovedIncomingTransfers); err != nil {
			return err
		}
	}

	if oldPermissions.CanUpdateApprovedOutgoingTransfers != nil && newPermissions.CanUpdateApprovedOutgoingTransfers != nil {
		if err := ValidateUserApprovedTransferPermissionsUpdate(oldPermissions.CanUpdateApprovedOutgoingTransfers, newPermissions.CanUpdateApprovedOutgoingTransfers); err != nil {
			return err
		}
	}

	return nil
}


// Validate that the new permissions are valid and is not changing anything that they can't.
func ValidatePermissionsUpdate(oldPermissions *CollectionPermissions, newPermissions *CollectionPermissions, canBeNil bool) error {
	if err := ValidatePermissions(newPermissions, canBeNil); err != nil {
		return err
	}

	if err := ValidatePermissions(oldPermissions, canBeNil); err != nil {
		return err
	}

	if oldPermissions.CanDeleteCollection != nil && newPermissions.CanDeleteCollection != nil {
		if err := ValidateActionPermissionUpdate(oldPermissions.CanDeleteCollection, newPermissions.CanDeleteCollection); err != nil {
			return err
		}
	}

	if oldPermissions.CanUpdateManager != nil && newPermissions.CanUpdateManager != nil {
		if err := ValidateTimedUpdatePermissionUpdate(oldPermissions.CanUpdateManager, newPermissions.CanUpdateManager); err != nil {
			return err
		}
	}

	if oldPermissions.CanUpdateCustomData != nil && newPermissions.CanUpdateCustomData != nil {
		if err := ValidateTimedUpdatePermissionUpdate(oldPermissions.CanUpdateCustomData, newPermissions.CanUpdateCustomData); err != nil {
			return err
		}
	}
	
	if oldPermissions.CanUpdateStandard != nil && newPermissions.CanUpdateStandard != nil {
		if err := ValidateTimedUpdatePermissionUpdate(oldPermissions.CanUpdateStandard, newPermissions.CanUpdateStandard); err != nil {
			return err
		}
	}

	if oldPermissions.CanArchive != nil && newPermissions.CanArchive != nil {
		if err := ValidateTimedUpdatePermissionUpdate(oldPermissions.CanArchive, newPermissions.CanArchive); err != nil {
			return err
		}
	}

	if oldPermissions.CanUpdateOffChainBalancesMetadata != nil && newPermissions.CanUpdateOffChainBalancesMetadata != nil {
		if err := ValidateTimedUpdatePermissionUpdate(oldPermissions.CanUpdateOffChainBalancesMetadata, newPermissions.CanUpdateOffChainBalancesMetadata); err != nil {
			return err
		}
	}

	if oldPermissions.CanUpdateCollectionMetadata != nil && newPermissions.CanUpdateCollectionMetadata != nil {
		if err := ValidateTimedUpdatePermissionUpdate(oldPermissions.CanUpdateCollectionMetadata, newPermissions.CanUpdateCollectionMetadata); err != nil {
			return err
		}
	}

	if oldPermissions.CanUpdateContractAddress != nil && newPermissions.CanUpdateContractAddress != nil {
		if err := ValidateTimedUpdatePermissionUpdate(oldPermissions.CanUpdateContractAddress, newPermissions.CanUpdateContractAddress); err != nil {
			return err
		}
	}

	if oldPermissions.CanCreateMoreBadges != nil && newPermissions.CanCreateMoreBadges != nil {
		if err := ValidateActionWithBadgeIdsPermissionUpdate(oldPermissions.CanCreateMoreBadges, newPermissions.CanCreateMoreBadges); err != nil {
			return err
		}
	}

	if oldPermissions.CanUpdateBadgeMetadata != nil && newPermissions.CanUpdateBadgeMetadata != nil {
		if err := ValidateTimedUpdateWithBadgeIdsPermissionUpdate(oldPermissions.CanUpdateBadgeMetadata, newPermissions.CanUpdateBadgeMetadata); err != nil {
			return err
		}
	}

	if oldPermissions.CanUpdateInheritedBalances != nil && newPermissions.CanUpdateInheritedBalances != nil {
		if err := ValidateTimedUpdateWithBadgeIdsPermissionUpdate(oldPermissions.CanUpdateInheritedBalances, newPermissions.CanUpdateInheritedBalances); err != nil {
			return err
		}
	}

	if oldPermissions.CanUpdateApprovedTransfers != nil && newPermissions.CanUpdateApprovedTransfers != nil {
		if err := ValidateCollectionApprovedTransferPermissionsUpdate(oldPermissions.CanUpdateApprovedTransfers, newPermissions.CanUpdateApprovedTransfers); err != nil {
			return err
		}
	}

	return nil
}