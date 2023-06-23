package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"math"

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
}

func UniversalRemoveOverlaps(handled *UniversalPermissionDetails, valueToCheck *UniversalPermissionDetails) []*UniversalPermissionDetails {
	timelineTimesAfterRemoval, removed := RemoveIdsFromIdRange(handled.TimelineTime, valueToCheck.TimelineTime)
	badgesAfterRemoval, badgeRemoved := RemoveIdsFromIdRange(handled.BadgeId, valueToCheck.BadgeId)
	transferTimesAfterRemoval, transferTimeRemoved := RemoveIdsFromIdRange(handled.TransferTime, valueToCheck.TransferTime)

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
	if !removed || !badgeRemoved || !transferTimeRemoved || !toMappingRemoved || !fromMappingRemoved || !initiatedByMappingRemoved {
		remaining = append(remaining, valueToCheck)
		return remaining
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
			TimelineTime: handled.TimelineTime,
			BadgeId: badgeAfterRemoval,
			TransferTime: valueToCheck.TransferTime,
			ToMappingId: valueToCheck.ToMappingId,
			FromMappingId: valueToCheck.FromMappingId,
			InitiatedByMappingId: valueToCheck.InitiatedByMappingId,
		})
	}

	for _, transferTimeAfterRemoval := range transferTimesAfterRemoval {
		remaining = append(remaining, &UniversalPermissionDetails{
			TimelineTime: handled.TimelineTime,
			BadgeId: handled.BadgeId,
			TransferTime: transferTimeAfterRemoval,
			ToMappingId: valueToCheck.ToMappingId,
			FromMappingId: valueToCheck.FromMappingId,
			InitiatedByMappingId: valueToCheck.InitiatedByMappingId,
		})
	}

			
	//For the mapping IDs, it is an all or nothing system. Either we remove the whole mapping ID or we don't
	//If we reach here, we know that toMappingRemoved, fromMappingRemoved, and initiatedByMappingRemoved are all true
	//Thus, there is nothing left to add here
	
	return remaining
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
							},
						}

						for _, handledPermission := range handled {
							newBrokenDown := []*UniversalPermissionDetails{}
							for _, broken := range brokenDown {

								for _, remaining := range UniversalRemoveOverlaps(&UniversalPermissionDetails{
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
								}) {
									newBrokenDown = append(newBrokenDown, &UniversalPermissionDetails{
										TimelineTime: remaining.TimelineTime,
										BadgeId: remaining.BadgeId,
										TransferTime: remaining.TransferTime,
										ToMappingId: remaining.ToMappingId,
										FromMappingId: remaining.FromMappingId,
										InitiatedByMappingId: remaining.InitiatedByMappingId,
										PermittedTimes: permittedTimes,
										ForbiddenTimes: forbiddenTimes,
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

func ValidateUniversalPermissionUpdate(oldPermissions []*UniversalPermissionDetails, newPermissions []*UniversalPermissionDetails) error {

	for len(oldPermissions) > 0 {
		badgeIdToCheck := oldPermissions[0].BadgeId.Start
		timelineTimeToCheck := oldPermissions[0].TimelineTime.Start
		transferTimeToCheck := oldPermissions[0].TransferTime.Start
		toMappingIdToCheck := oldPermissions[0].ToMappingId
		fromMappingIdToCheck := oldPermissions[0].FromMappingId
		initiatedByMappingIdToCheck := oldPermissions[0].InitiatedByMappingId

		oldPermittedTimes := oldPermissions[0].PermittedTimes
		oldForbiddenTimes := oldPermissions[0].ForbiddenTimes

		newPermittedTimes := []*IdRange{}
		newForbiddenTimes := []*IdRange{}
		newBadgeIds := &IdRange{ Start: sdk.NewUint(math.MaxUint64), End: sdk.NewUint(math.MaxUint64) }
		newTimelineTimes := &IdRange{ Start: sdk.NewUint(math.MaxUint64), End: sdk.NewUint(math.MaxUint64) }
		newTransferTimes := &IdRange{ Start: sdk.NewUint(math.MaxUint64), End: sdk.NewUint(math.MaxUint64) }
		newToMappingId := ""
		newFromMappingId := ""
		newInitiatedByMappingId := ""
		newFound := false
		for _, newPermission := range newPermissions {
			_, found := SearchIdRangesForId(badgeIdToCheck, []*IdRange{newPermission.BadgeId})
			_, timelineTimeFound := SearchIdRangesForId(timelineTimeToCheck, []*IdRange{newPermission.TimelineTime})
			_, transferTimeFound := SearchIdRangesForId(transferTimeToCheck, []*IdRange{newPermission.TransferTime})
			toMappingIdFound := toMappingIdToCheck == newPermission.ToMappingId
			fromMappingIdFound := fromMappingIdToCheck == newPermission.FromMappingId
			initiatedByMappingIdFound := initiatedByMappingIdToCheck == newPermission.InitiatedByMappingId

			if found && timelineTimeFound && transferTimeFound && toMappingIdFound && fromMappingIdFound && initiatedByMappingIdFound {
				newFound = true
				newBadgeIds = newPermission.BadgeId
				newTimelineTimes = newPermission.TimelineTime
				newTransferTimes = newPermission.TransferTime
				newToMappingId = newPermission.ToMappingId
				newFromMappingId = newPermission.FromMappingId
				newInitiatedByMappingId = newPermission.InitiatedByMappingId
				newPermittedTimes = append(newPermittedTimes, newPermission.PermittedTimes...)
				newForbiddenTimes = append(newForbiddenTimes, newPermission.ForbiddenTimes...)
			}
		}

		if !newFound {
			errMsg := "permission ( "
			//If we are using the dummy range, we know that the field is not used
			if oldPermissions[0].BadgeId.Start != sdk.NewUint(math.MaxUint64) || oldPermissions[0].BadgeId.End != sdk.NewUint(math.MaxUint64) {
				errMsg += "badgeId: " + oldPermissions[0].BadgeId.Start.String() + " "
			}

			if oldPermissions[0].TimelineTime.Start != sdk.NewUint(math.MaxUint64) || oldPermissions[0].TimelineTime.End != sdk.NewUint(math.MaxUint64) {
				errMsg += "timelineTime: " + oldPermissions[0].TimelineTime.Start.String() + " "
			}

			if oldPermissions[0].TransferTime.Start != sdk.NewUint(math.MaxUint64) || oldPermissions[0].TransferTime.End != sdk.NewUint(math.MaxUint64) {
				errMsg += "transferTime: " + oldPermissions[0].TransferTime.Start.String() + " "
			}

			if oldPermissions[0].ToMappingId != "" {
				errMsg += "toMappingId: " + oldPermissions[0].ToMappingId + " "
			}

			if oldPermissions[0].FromMappingId != "" {
				errMsg += "fromMappingId: " + oldPermissions[0].FromMappingId + " "
			}

			if oldPermissions[0].InitiatedByMappingId != "" {
				errMsg += "initiatedByMappingId: " + oldPermissions[0].InitiatedByMappingId + " "
			}

			errMsg += ") found in old permissions but not in new permissions"
			
			return sdkerrors.Wrapf(ErrInvalidPermissions, errMsg)
		}

		//Check if all old permitted times are in new permitted times
		//Check if all old forbidden times are in new forbidden times
		leftoverPermittedTimes := RemoveIdRangeFromIdRange(newPermittedTimes, oldPermittedTimes)
		leftoverForbiddenTimes := RemoveIdRangeFromIdRange(newForbiddenTimes, oldForbiddenTimes)

		if len(leftoverPermittedTimes) > 0 || len(leftoverForbiddenTimes) > 0 {
			errMsg := ""
			if len(leftoverPermittedTimes) > 0 {
				errMsg += "the times ( "
				for _, oldPermittedTime := range leftoverPermittedTimes {
					errMsg += oldPermittedTime.Start.String() + "-" + oldPermittedTime.End.String() + " "
				}
				errMsg += ") were found in old permitted times but not in new permitted times."
			}
			if len(leftoverForbiddenTimes) > 0 {
				errMsg += "the times ( "
				for _, oldForbiddenTime := range leftoverForbiddenTimes {
					errMsg += oldForbiddenTime.Start.String() + "-" + oldForbiddenTime.End.String() + " "	
				}
				errMsg += ") were found in old forbidden times but not in new forbidden times."
			}

			return sdkerrors.Wrapf(ErrInvalidPermissions, "all old permission times must be found in new permissions %s %s", oldPermittedTimes, oldForbiddenTimes)
		}

		//Even though we searched for just handled the start, we also handled any overlap between the two badgeIds
		//So we remove the overlap from the unhandled badgeIds

		remaining := UniversalRemoveOverlaps(&UniversalPermissionDetails{
			TimelineTime: newTimelineTimes,
			BadgeId: newBadgeIds,
			TransferTime: newTransferTimes,
			ToMappingId: newToMappingId,
			FromMappingId: newFromMappingId,
			InitiatedByMappingId: newInitiatedByMappingId,
		}, &UniversalPermissionDetails{
			TimelineTime: oldPermissions[0].TimelineTime,
			BadgeId: oldPermissions[0].BadgeId,
			TransferTime: oldPermissions[0].TransferTime,
			ToMappingId: oldPermissions[0].ToMappingId,
			FromMappingId: oldPermissions[0].FromMappingId,
			InitiatedByMappingId: oldPermissions[0].InitiatedByMappingId,
		})

		

		newUnhandledPermissions := []*UniversalPermissionDetails{}
		for _, remainingPermission := range remaining {
			newUnhandledPermissions = append(newUnhandledPermissions, &UniversalPermissionDetails{
				TimelineTime: remainingPermission.TimelineTime,
				BadgeId: remainingPermission.BadgeId,
				TransferTime: remainingPermission.TransferTime,
				ToMappingId: remainingPermission.ToMappingId,
				FromMappingId: remainingPermission.FromMappingId,
				InitiatedByMappingId: remainingPermission.InitiatedByMappingId,
				PermittedTimes: oldPermissions[0].PermittedTimes,
				ForbiddenTimes: oldPermissions[0].ForbiddenTimes,
			})
		}
		newUnhandledPermissions = append(newUnhandledPermissions, oldPermissions[1:]...)
		oldPermissions = newUnhandledPermissions
	}

	return nil
}

func ValidateActionWithBadgeIdsPermissionUpdate(oldPermissions []*ActionWithBadgeIdsPermission, newPermissions []*ActionWithBadgeIdsPermission) error {
	if err := ValidateActionWithBadgeIdsPermission(oldPermissions); err != nil {
		return err
	}

	if err := ValidateActionWithBadgeIdsPermission(newPermissions); err != nil {
		return err
	}

	castedOldPermissions := []*UniversalPermission{}
	for _, oldPermission := range oldPermissions {
		castedCombinations := []*UniversalCombination{}
		for _, oldCombination := range oldPermission.Combinations {
			castedCombinations = append(castedCombinations, &UniversalCombination{
				BadgeIdsOptions: oldCombination.BadgeIdsOptions,
				PermittedTimesOptions: oldCombination.PermittedTimesOptions,
				ForbiddenTimesOptions: oldCombination.ForbiddenTimesOptions,
			})
		}

		castedOldPermissions = append(castedOldPermissions, &UniversalPermission{
			DefaultValues: &UniversalDefaultValues{
				BadgeIds: oldPermission.DefaultValues.BadgeIds,
				UsesBadgeIds: true,
				PermittedTimes: oldPermission.DefaultValues.PermittedTimes,
				ForbiddenTimes: oldPermission.DefaultValues.ForbiddenTimes,
			},
			Combinations: castedCombinations,
		})
	}

	castedNewPermissions := []*UniversalPermission{}
	for _, newPermission := range newPermissions {
		castedCombinations := []*UniversalCombination{}
		for _, newCombination := range newPermission.Combinations {
			castedCombinations = append(castedCombinations, &UniversalCombination{
				BadgeIdsOptions: newCombination.BadgeIdsOptions,
				PermittedTimesOptions: newCombination.PermittedTimesOptions,
				ForbiddenTimesOptions: newCombination.ForbiddenTimesOptions,
			})
		}

		castedNewPermissions = append(castedNewPermissions, &UniversalPermission{
			DefaultValues: &UniversalDefaultValues{
				BadgeIds: newPermission.DefaultValues.BadgeIds,
				UsesBadgeIds: true,
				PermittedTimes: newPermission.DefaultValues.PermittedTimes,
				ForbiddenTimes: newPermission.DefaultValues.ForbiddenTimes,
			},
			Combinations: castedCombinations,
		})
	}

	
	err := ValidateUniversalPermissionUpdate(GetFirstMatchOnly(castedOldPermissions), GetFirstMatchOnly(castedNewPermissions))
	if err != nil {
		return err
	}

	return nil
}


func ValidateTimedUpdatePermissionUpdate(oldPermissions []*TimedUpdatePermission, newPermissions []*TimedUpdatePermission) error {
	if err := ValidateTimedUpdatePermission(oldPermissions); err != nil {
		return err
	}

	if err := ValidateTimedUpdatePermission(newPermissions); err != nil {
		return err
	}

	castedOldPermissions := []*UniversalPermission{}
	for _, oldPermission := range oldPermissions {
		castedCombinations := []*UniversalCombination{}
		for _, oldCombination := range oldPermission.Combinations {
			castedCombinations = append(castedCombinations, &UniversalCombination{
				PermittedTimesOptions: oldCombination.PermittedTimesOptions,
				ForbiddenTimesOptions: oldCombination.ForbiddenTimesOptions,
				TimelineTimesOptions: oldCombination.TimelineTimesOptions,
			})
		}

		castedOldPermissions = append(castedOldPermissions, &UniversalPermission{
			DefaultValues: &UniversalDefaultValues{
				TimelineTimes: oldPermission.DefaultValues.TimelineTimes,
				UsesTimelineTimes: true,
				PermittedTimes: oldPermission.DefaultValues.PermittedTimes,
				ForbiddenTimes: oldPermission.DefaultValues.ForbiddenTimes,
			},
			Combinations: castedCombinations,
		})
	}

	castedNewPermissions := []*UniversalPermission{}
	for _, newPermission := range newPermissions {
		castedCombinations := []*UniversalCombination{}
		for _, newCombination := range newPermission.Combinations {
			castedCombinations = append(castedCombinations, &UniversalCombination{
				PermittedTimesOptions: newCombination.PermittedTimesOptions,
				ForbiddenTimesOptions: newCombination.ForbiddenTimesOptions,
				TimelineTimesOptions: newCombination.TimelineTimesOptions,
			})
		}

		castedNewPermissions = append(castedNewPermissions, &UniversalPermission{
			DefaultValues: &UniversalDefaultValues{
				TimelineTimes: newPermission.DefaultValues.TimelineTimes,
				UsesTimelineTimes: true,
				PermittedTimes: newPermission.DefaultValues.PermittedTimes,
				ForbiddenTimes: newPermission.DefaultValues.ForbiddenTimes,
			},
			Combinations: castedCombinations,
		})
	}

	err := ValidateUniversalPermissionUpdate(GetFirstMatchOnly(castedOldPermissions), GetFirstMatchOnly(castedNewPermissions))
	if err != nil {
		return err
	}

	return nil
}

func ValidateTimedUpdateWithBadgeIdsPermissionUpdate(oldPermissions []*TimedUpdateWithBadgeIdsPermission, newPermissions []*TimedUpdateWithBadgeIdsPermission) error {
	if err := ValidateTimedUpdateWithBadgeIdsPermission(oldPermissions); err != nil {
		return err
	}

	if err := ValidateTimedUpdateWithBadgeIdsPermission(newPermissions); err != nil {
		return err
	}

	castedOldPermissions := []*UniversalPermission{}
	for _, oldPermission := range oldPermissions {
		castedCombinations := []*UniversalCombination{}
		for _, oldCombination := range oldPermission.Combinations {
			castedCombinations = append(castedCombinations, &UniversalCombination{
				BadgeIdsOptions: oldCombination.BadgeIdsOptions,
				PermittedTimesOptions: oldCombination.PermittedTimesOptions,
				ForbiddenTimesOptions: oldCombination.ForbiddenTimesOptions,
				TimelineTimesOptions: oldCombination.TimelineTimesOptions,
			})
		}

		castedOldPermissions = append(castedOldPermissions, &UniversalPermission{
			DefaultValues: &UniversalDefaultValues{
				TimelineTimes: oldPermission.DefaultValues.TimelineTimes,
				BadgeIds: oldPermission.DefaultValues.BadgeIds,
				UsesTimelineTimes: true,
				UsesBadgeIds: true,
				PermittedTimes: oldPermission.DefaultValues.PermittedTimes,
				ForbiddenTimes: oldPermission.DefaultValues.ForbiddenTimes,
			},
			Combinations: castedCombinations,
		})
	}

	castedNewPermissions := []*UniversalPermission{}
	for _, newPermission := range newPermissions {
		castedCombinations := []*UniversalCombination{}
		for _, newCombination := range newPermission.Combinations {
			castedCombinations = append(castedCombinations, &UniversalCombination{
				BadgeIdsOptions: newCombination.BadgeIdsOptions,
				PermittedTimesOptions: newCombination.PermittedTimesOptions,
				ForbiddenTimesOptions: newCombination.ForbiddenTimesOptions,
				TimelineTimesOptions: newCombination.TimelineTimesOptions,
			})
		}

		castedNewPermissions = append(castedNewPermissions, &UniversalPermission{
			DefaultValues: &UniversalDefaultValues{
				TimelineTimes: newPermission.DefaultValues.TimelineTimes,
				BadgeIds: newPermission.DefaultValues.BadgeIds,
				UsesTimelineTimes: true,
				UsesBadgeIds: true,
				PermittedTimes: newPermission.DefaultValues.PermittedTimes,
				ForbiddenTimes: newPermission.DefaultValues.ForbiddenTimes,
			},
			Combinations: castedCombinations,
		})
	}

	err := ValidateUniversalPermissionUpdate(GetFirstMatchOnly(castedOldPermissions), GetFirstMatchOnly(castedNewPermissions))
	if err != nil {
		return err
	}

	return nil
}

func ValidateCollectionApprovedTransferPermissionsUpdate(oldPermissions []*CollectionApprovedTransferPermission, newPermissions []*CollectionApprovedTransferPermission) error {
	if err := ValidateCollectionApprovedTransferPermissions(oldPermissions); err != nil {
		return err
	}

	if err := ValidateCollectionApprovedTransferPermissions(newPermissions); err != nil {
		return err
	}

	castedOldPermissions := []*UniversalPermission{}
	for _, oldPermission := range oldPermissions {
		castedCombinations := []*UniversalCombination{}
		for _, oldCombination := range oldPermission.Combinations {
			castedCombinations = append(castedCombinations, &UniversalCombination{
				BadgeIdsOptions: oldCombination.BadgeIdsOptions,
				PermittedTimesOptions: oldCombination.PermittedTimesOptions,
				ForbiddenTimesOptions: oldCombination.ForbiddenTimesOptions,
				TimelineTimesOptions: oldCombination.TimelineTimesOptions,
				TransferTimesOptions: oldCombination.TransferTimesOptions,
				ToMappingIdOptions: oldCombination.ToMappingIdOptions,
				FromMappingIdOptions: oldCombination.FromMappingIdOptions,
				InitiatedByMappingIdOptions: oldCombination.InitiatedByMappingIdOptions,
			})
		}

		castedOldPermissions = append(castedOldPermissions, &UniversalPermission{
			DefaultValues: &UniversalDefaultValues{
				BadgeIds: oldPermission.DefaultValues.BadgeIds,
				TimelineTimes: oldPermission.DefaultValues.TimelineTimes,
				TransferTimes: oldPermission.DefaultValues.TransferTimes,
				ToMappingId: oldPermission.DefaultValues.ToMappingId,
				FromMappingId: oldPermission.DefaultValues.FromMappingId,
				InitiatedByMappingId: oldPermission.DefaultValues.InitiatedByMappingId,
				UsesBadgeIds: true,
				UsesTimelineTimes: true,
				UsesTransferTimes: true,
				UsesToMappingId: true,
				UsesFromMappingId: true,
				UsesInitiatedByMappingId: true,
				PermittedTimes: oldPermission.DefaultValues.PermittedTimes,
				ForbiddenTimes: oldPermission.DefaultValues.ForbiddenTimes,
			},
			Combinations: castedCombinations,
		})
	}

	castedNewPermissions := []*UniversalPermission{}
	for _, newPermission := range newPermissions {
		castedCombinations := []*UniversalCombination{}
		for _, newCombination := range newPermission.Combinations {
			castedCombinations = append(castedCombinations, &UniversalCombination{
				BadgeIdsOptions: newCombination.BadgeIdsOptions,
				PermittedTimesOptions: newCombination.PermittedTimesOptions,
				ForbiddenTimesOptions: newCombination.ForbiddenTimesOptions,
				TimelineTimesOptions: newCombination.TimelineTimesOptions,
				TransferTimesOptions: newCombination.TransferTimesOptions,
				ToMappingIdOptions: newCombination.ToMappingIdOptions,
				FromMappingIdOptions: newCombination.FromMappingIdOptions,
				InitiatedByMappingIdOptions: newCombination.InitiatedByMappingIdOptions,
			})
		}

		castedNewPermissions = append(castedNewPermissions, &UniversalPermission{
			DefaultValues: &UniversalDefaultValues{
				BadgeIds: newPermission.DefaultValues.BadgeIds,
				TimelineTimes: newPermission.DefaultValues.TimelineTimes,
				TransferTimes: newPermission.DefaultValues.TransferTimes,
				ToMappingId: newPermission.DefaultValues.ToMappingId,
				FromMappingId: newPermission.DefaultValues.FromMappingId,
				InitiatedByMappingId: newPermission.DefaultValues.InitiatedByMappingId,
				UsesBadgeIds: true,
				UsesTimelineTimes: true,
				UsesTransferTimes: true,
				UsesToMappingId: true,
				UsesFromMappingId: true,
				UsesInitiatedByMappingId: true,
				PermittedTimes: newPermission.DefaultValues.PermittedTimes,
				ForbiddenTimes: newPermission.DefaultValues.ForbiddenTimes,
			},
			Combinations: castedCombinations,
		})
	}

	err := ValidateUniversalPermissionUpdate(GetFirstMatchOnly(castedOldPermissions), GetFirstMatchOnly(castedNewPermissions))
	if err != nil {
		return err
	}

	return nil
}

func ValidateActionPermissionUpdate(oldPermissions []*ActionPermission, newPermissions []*ActionPermission) error {
	if err := ValidateActionPermission(oldPermissions); err != nil {
		return err
	}

	if err := ValidateActionPermission(newPermissions); err != nil {
		return err
	}

	castedOldPermissions := []*UniversalPermission{}
	for _, oldPermission := range oldPermissions {
		castedCombinations := []*UniversalCombination{}
		for _, oldCombination := range oldPermission.Combinations {
			castedCombinations = append(castedCombinations, &UniversalCombination{
				PermittedTimesOptions: oldCombination.PermittedTimesOptions,
				ForbiddenTimesOptions: oldCombination.ForbiddenTimesOptions,
			})
		}

		castedOldPermissions = append(castedOldPermissions, &UniversalPermission{
			DefaultValues: &UniversalDefaultValues{
				PermittedTimes: oldPermission.DefaultValues.PermittedTimes,
				ForbiddenTimes: oldPermission.DefaultValues.ForbiddenTimes,
			},
			Combinations: castedCombinations,
		})

	}

	castedNewPermissions := []*UniversalPermission{}
	for _, newPermission := range newPermissions {
		castedCombinations := []*UniversalCombination{}
		for _, newCombination := range newPermission.Combinations {
			castedCombinations = append(castedCombinations, &UniversalCombination{
				PermittedTimesOptions: newCombination.PermittedTimesOptions,
				ForbiddenTimesOptions: newCombination.ForbiddenTimesOptions,
			})
		}

		castedNewPermissions = append(castedNewPermissions, &UniversalPermission{
			DefaultValues: &UniversalDefaultValues{
				PermittedTimes: newPermission.DefaultValues.PermittedTimes,
				ForbiddenTimes: newPermission.DefaultValues.ForbiddenTimes,
			},
			Combinations: castedCombinations,
		})
	}

	err := ValidateUniversalPermissionUpdate(GetFirstMatchOnly(castedOldPermissions), GetFirstMatchOnly(castedNewPermissions))
	if err != nil {
		return err
	}

	return nil
}

func ValidateUserApprovedTransferPermissionsUpdate(oldPermissions []*UserApprovedTransferPermission, newPermissions []*UserApprovedTransferPermission) error {
	if err := ValidateUserApprovedTransferPermissions(oldPermissions); err != nil {
		return err
	}

	if err := ValidateUserApprovedTransferPermissions(newPermissions); err != nil {
		return err
	}

	castedOldPermissions := []*UniversalPermission{}
	for _, oldPermission := range oldPermissions {
		castedCombinations := []*UniversalCombination{}
		for _, oldCombination := range oldPermission.Combinations {
			castedCombinations = append(castedCombinations, &UniversalCombination{
				BadgeIdsOptions: oldCombination.BadgeIdsOptions,
				PermittedTimesOptions: oldCombination.PermittedTimesOptions,
				ForbiddenTimesOptions: oldCombination.ForbiddenTimesOptions,
				TimelineTimesOptions: oldCombination.TimelineTimesOptions,
				TransferTimesOptions: oldCombination.TransferTimesOptions,
				ToMappingIdOptions: oldCombination.ToMappingIdOptions,
				InitiatedByMappingIdOptions: oldCombination.InitiatedByMappingIdOptions,
			})
		}

		castedOldPermissions = append(castedOldPermissions, &UniversalPermission{
			DefaultValues: &UniversalDefaultValues{
				BadgeIds: oldPermission.DefaultValues.BadgeIds,
				TimelineTimes: oldPermission.DefaultValues.TimelineTimes,
				TransferTimes: oldPermission.DefaultValues.TransferTimes,
				ToMappingId: oldPermission.DefaultValues.ToMappingId,
				InitiatedByMappingId: oldPermission.DefaultValues.InitiatedByMappingId,
				UsesBadgeIds: true,
				UsesTimelineTimes: true,
				UsesTransferTimes: true,
				UsesToMappingId: true,
				UsesInitiatedByMappingId: true,
				PermittedTimes: oldPermission.DefaultValues.PermittedTimes,
				ForbiddenTimes: oldPermission.DefaultValues.ForbiddenTimes,
			},
			Combinations: castedCombinations,
		})
	}

	castedNewPermissions := []*UniversalPermission{}
	for _, newPermission := range newPermissions {
		castedCombinations := []*UniversalCombination{}
		for _, newCombination := range newPermission.Combinations {
			castedCombinations = append(castedCombinations, &UniversalCombination{
				BadgeIdsOptions: newCombination.BadgeIdsOptions,
				PermittedTimesOptions: newCombination.PermittedTimesOptions,
				ForbiddenTimesOptions: newCombination.ForbiddenTimesOptions,
				TimelineTimesOptions: newCombination.TimelineTimesOptions,
				TransferTimesOptions: newCombination.TransferTimesOptions,
				ToMappingIdOptions: newCombination.ToMappingIdOptions,
				InitiatedByMappingIdOptions: newCombination.InitiatedByMappingIdOptions,
			})
		}

		castedNewPermissions = append(castedNewPermissions, &UniversalPermission{
			DefaultValues: &UniversalDefaultValues{
				BadgeIds: newPermission.DefaultValues.BadgeIds,
				TimelineTimes: newPermission.DefaultValues.TimelineTimes,
				TransferTimes: newPermission.DefaultValues.TransferTimes,
				ToMappingId: newPermission.DefaultValues.ToMappingId,
				InitiatedByMappingId: newPermission.DefaultValues.InitiatedByMappingId,
				UsesBadgeIds: true,
				UsesTimelineTimes: true,
				UsesTransferTimes: true,
				UsesToMappingId: true,
				UsesInitiatedByMappingId: true,
				PermittedTimes: newPermission.DefaultValues.PermittedTimes,
				ForbiddenTimes: newPermission.DefaultValues.ForbiddenTimes,
			},
			Combinations: castedCombinations,
		})
	}

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

	if oldPermissions.CanUpdateApprovedTransfers != nil && newPermissions.CanUpdateApprovedTransfers != nil {
		if err := ValidateUserApprovedTransferPermissionsUpdate(oldPermissions.CanUpdateApprovedTransfers, newPermissions.CanUpdateApprovedTransfers); err != nil {
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

	if oldPermissions.CanUpdateCollectionMetadata != nil && newPermissions.CanUpdateCollectionMetadata != nil {
		if err := ValidateTimedUpdatePermissionUpdate(oldPermissions.CanUpdateCustomData, newPermissions.CanUpdateCustomData); err != nil {
			return err
		}
	}

	if oldPermissions.CanUpdateOffChainBalancesMetadata != nil && newPermissions.CanUpdateOffChainBalancesMetadata != nil {
		if err := ValidateTimedUpdatePermissionUpdate(oldPermissions.CanUpdateManager, newPermissions.CanUpdateManager); err != nil {
			return err
		}
	}

	if oldPermissions.CanUpdateContractAddress != nil && newPermissions.CanUpdateContractAddress != nil {
		if err := ValidateTimedUpdatePermissionUpdate(oldPermissions.CanUpdateCollectionMetadata, newPermissions.CanUpdateCollectionMetadata); err != nil {
			return err
		}
	}

	if oldPermissions.CanArchive != nil && newPermissions.CanArchive != nil {
		if err := ValidateTimedUpdatePermissionUpdate(oldPermissions.CanUpdateOffChainBalancesMetadata, newPermissions.CanUpdateOffChainBalancesMetadata); err != nil {
			return err
		}
	}

	if oldPermissions.CanUpdateBadgeMetadata != nil && newPermissions.CanUpdateBadgeMetadata != nil {
		if err := ValidateTimedUpdatePermissionUpdate(oldPermissions.CanUpdateContractAddress, newPermissions.CanUpdateContractAddress); err != nil {
			return err
		}
	}

	if oldPermissions.CanCreateMoreBadges != nil && newPermissions.CanCreateMoreBadges != nil {
		if err := ValidateTimedUpdatePermissionUpdate(oldPermissions.CanArchive, newPermissions.CanArchive); err != nil {
			return err
		}
	}

	if oldPermissions.CanUpdateInheritedBalances != nil && newPermissions.CanUpdateInheritedBalances != nil {
		if err := ValidateTimedUpdateWithBadgeIdsPermissionUpdate(oldPermissions.CanUpdateBadgeMetadata, newPermissions.CanUpdateBadgeMetadata); err != nil {
			return err
		}
	}

	if oldPermissions.CanUpdateApprovedTransfers != nil && newPermissions.CanUpdateApprovedTransfers != nil {
		if err := ValidateActionWithBadgeIdsPermissionUpdate(oldPermissions.CanCreateMoreBadges, newPermissions.CanCreateMoreBadges); err != nil {
			return err
		}
	}

	if oldPermissions.CanUpdateApprovedTransfers != nil && newPermissions.CanUpdateApprovedTransfers != nil {
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