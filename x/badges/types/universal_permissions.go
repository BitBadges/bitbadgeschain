package types

import (
	fmt "fmt"
	"math"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
)

//TODO: This file should probably be refactored a lot, but it currently works.
//			It is also not user-facing or dev-facing, so I am okay with how it is now

//For permissions, we have many types of permissions that are all very similar to each other
//Here, we abstract all those permissions to a UniversalPermission struct in order to reuse code.
//When casting to a UniversalPermission, we use fake dummy values for the unused values to avoid messing up the logic
//
//This file implements certain logic using UniversalPermissions such as overlaps and getting first match only
//This is used in many places around the codebase

type UniversalCombination struct {
	TimelineTimesOptions *ValueOptions

	FromMappingOptions        *ValueOptions
	ToMappingOptions          *ValueOptions
	InitiatedByMappingOptions *ValueOptions
	TransferTimesOptions      *ValueOptions
	BadgeIdsOptions           *ValueOptions
	OwnedTimesOptions     *ValueOptions

	PermittedTimesOptions *ValueOptions
	ForbiddenTimesOptions *ValueOptions
}

type UniversalPermission struct {
	DefaultValues *UniversalDefaultValues
	Combinations  []*UniversalCombination
}

type UniversalDefaultValues struct {
	BadgeIds           []*UintRange
	TimelineTimes      []*UintRange
	TransferTimes      []*UintRange
	OwnedTimes     []*UintRange
	ToMapping          *AddressMapping
	FromMapping        *AddressMapping
	InitiatedByMapping *AddressMapping

	PermittedTimes []*UintRange
	ForbiddenTimes []*UintRange

	UsesBadgeIds           bool
	UsesTimelineTimes      bool
	UsesTransferTimes      bool
	UsesToMapping          bool
	UsesFromMapping        bool
	UsesInitiatedByMapping bool
	UsesOwnedTimes     bool

	ArbitraryValue interface{}
}

type UniversalPermissionDetails struct {
	BadgeId            *UintRange
	TimelineTime       *UintRange
	TransferTime       *UintRange
	OwnershipTime      *UintRange
	ToMapping          *AddressMapping
	FromMapping        *AddressMapping
	InitiatedByMapping *AddressMapping

	//These fields are not actually used in the overlapping logic.
	//They are just along for the ride and used later for lookups
	PermittedTimes []*UintRange
	ForbiddenTimes []*UintRange

	ArbitraryValue interface{}
}

type Overlap struct {
	Overlap       *UniversalPermissionDetails
	FirstDetails  *UniversalPermissionDetails
	SecondDetails *UniversalPermissionDetails
}

func GetOverlapsAndNonOverlaps(firstDetails []*UniversalPermissionDetails, secondDetails []*UniversalPermissionDetails) ([]*Overlap, []*UniversalPermissionDetails, []*UniversalPermissionDetails) {
	//TODO: deep copy???
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
					Overlap:       overlap,
					FirstDetails:  oldPermission,
					SecondDetails: newPermission,
				})
				inOldButNotNew, _ = UniversalRemoveOverlapFromValues(overlap, inOldButNotNew)
				inNewButNotOld, _ = UniversalRemoveOverlapFromValues(overlap, inNewButNotOld)
			}
		}
	}

	return allOverlaps, inOldButNotNew, inNewButNotOld
}

func UniversalRemoveOverlapFromValues(handled *UniversalPermissionDetails, valuesToCheck []*UniversalPermissionDetails) ([]*UniversalPermissionDetails, []*UniversalPermissionDetails) {
	newValuesToCheck := []*UniversalPermissionDetails{}
	removed := []*UniversalPermissionDetails{}
	for _, valueToCheck := range valuesToCheck {
		remaining, overlaps := UniversalRemoveOverlaps(handled, valueToCheck)
		newValuesToCheck = append(newValuesToCheck, remaining...)
		removed = append(removed, overlaps...)
	}

	return newValuesToCheck, removed
}

func IsAddressMappingEmpty(mapping *AddressMapping) bool {
	return len(mapping.Addresses) == 0 && mapping.IncludeAddresses
}

func UniversalRemoveOverlaps(handled *UniversalPermissionDetails, valueToCheck *UniversalPermissionDetails) ([]*UniversalPermissionDetails, []*UniversalPermissionDetails) {
	timelineTimesAfterRemoval, removedTimelineTimes := RemoveUintsFromUintRange(handled.TimelineTime, valueToCheck.TimelineTime)
	badgesAfterRemoval, removedBadges := RemoveUintsFromUintRange(handled.BadgeId, valueToCheck.BadgeId)
	transferTimesAfterRemoval, removedTransferTimes := RemoveUintsFromUintRange(handled.TransferTime, valueToCheck.TransferTime)
	ownedTimesAfterRemoval, removedOwnedTimes := RemoveUintsFromUintRange(handled.OwnershipTime, valueToCheck.OwnershipTime)

	toMappingAfterRemoval, removedToMapping := RemoveAddressMappingFromAddressMapping(handled.ToMapping, valueToCheck.ToMapping)
	fromMappingAfterRemoval, removedFromMapping := RemoveAddressMappingFromAddressMapping(handled.FromMapping, valueToCheck.FromMapping)
	initiatedByMappingAfterRemoval, removedInitiatedByMapping := RemoveAddressMappingFromAddressMapping(handled.InitiatedByMapping, valueToCheck.InitiatedByMapping)

	toMappingRemoved := !IsAddressMappingEmpty(removedToMapping)
	fromMappingRemoved := !IsAddressMappingEmpty(removedFromMapping)
	initiatedByMappingRemoved := !IsAddressMappingEmpty(removedInitiatedByMapping)

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
	if len(removedTimelineTimes) == 0 || len(removedBadges) == 0 || len(removedTransferTimes) == 0 || len(removedOwnedTimes) == 0 || !toMappingRemoved || !fromMappingRemoved || !initiatedByMappingRemoved {
		remaining = append(remaining, valueToCheck)
		return remaining, []*UniversalPermissionDetails{}
	}



	for _, timelineTimeAfterRemoval := range timelineTimesAfterRemoval {
		remaining = append(remaining, &UniversalPermissionDetails{
			TimelineTime:       timelineTimeAfterRemoval,
			BadgeId:            valueToCheck.BadgeId,
			TransferTime:       valueToCheck.TransferTime,
			OwnershipTime:      valueToCheck.OwnershipTime,
			ToMapping:          valueToCheck.ToMapping,
			FromMapping:        valueToCheck.FromMapping,
			InitiatedByMapping: valueToCheck.InitiatedByMapping,

			ArbitraryValue: valueToCheck.ArbitraryValue,
		})
	}

	for _, badgeAfterRemoval := range badgesAfterRemoval {
		remaining = append(remaining, &UniversalPermissionDetails{
			TimelineTime:       removedTimelineTimes[0], //We know there is only one because there can only be one interesection between two ranges
			BadgeId:            badgeAfterRemoval,
			TransferTime:       valueToCheck.TransferTime,
			OwnershipTime:      valueToCheck.OwnershipTime,
			ToMapping:          valueToCheck.ToMapping,
			FromMapping:        valueToCheck.FromMapping,
			InitiatedByMapping: valueToCheck.InitiatedByMapping,

			ArbitraryValue: valueToCheck.ArbitraryValue,
		})
	}

	for _, transferTimeAfterRemoval := range transferTimesAfterRemoval {
		remaining = append(remaining, &UniversalPermissionDetails{
			TimelineTime:       removedTimelineTimes[0], //We know there is only one because there can only be one interesection between two ranges
			BadgeId:            removedBadges[0],        //We know there is only one because there can only be one interesection between two ranges
			TransferTime:       transferTimeAfterRemoval,
			OwnershipTime:      valueToCheck.OwnershipTime,
			ToMapping:          valueToCheck.ToMapping,
			FromMapping:        valueToCheck.FromMapping,
			InitiatedByMapping: valueToCheck.InitiatedByMapping,

			ArbitraryValue: valueToCheck.ArbitraryValue,
		})
	}

	for _, ownershipTimeAfterRemoval := range ownedTimesAfterRemoval {
		remaining = append(remaining, &UniversalPermissionDetails{
			TimelineTime:       removedTimelineTimes[0], //We know there is only one because there can only be one interesection between two ranges
			BadgeId:            removedBadges[0],        //We know there is only one because there can only be one interesection between two ranges
			TransferTime:       removedTransferTimes[0], //We know there is only one because there can only be one interesection between two ranges
			OwnershipTime:      ownershipTimeAfterRemoval,
			ToMapping:          valueToCheck.ToMapping,
			FromMapping:        valueToCheck.FromMapping,
			InitiatedByMapping: valueToCheck.InitiatedByMapping,

			ArbitraryValue: valueToCheck.ArbitraryValue,
		})
	}

	if !IsAddressMappingEmpty(toMappingAfterRemoval) {
		remaining = append(remaining, &UniversalPermissionDetails{
			TimelineTime:       removedTimelineTimes[0], //We know there is only one because there can only be one interesection between two ranges
			BadgeId:            removedBadges[0],        //We know there is only one because there can only be one interesection between two ranges
			TransferTime:       removedTransferTimes[0], //We know there is only one because there can only be one interesection between two ranges
			OwnershipTime:      removedOwnedTimes[0],
			ToMapping:          toMappingAfterRemoval,
			FromMapping:        valueToCheck.FromMapping,
			InitiatedByMapping: valueToCheck.InitiatedByMapping,

			ArbitraryValue: valueToCheck.ArbitraryValue,
		})
	}

	if !IsAddressMappingEmpty(fromMappingAfterRemoval) {
		remaining = append(remaining, &UniversalPermissionDetails{
			TimelineTime:       removedTimelineTimes[0], //We know there is only one because there can only be one interesection between two ranges
			BadgeId:            removedBadges[0],        //We know there is only one because there can only be one interesection between two ranges
			TransferTime:       removedTransferTimes[0], //We know there is only one because there can only be one interesection between two ranges
			OwnershipTime:      removedOwnedTimes[0],
			ToMapping:          toMappingAfterRemoval,
			FromMapping:        fromMappingAfterRemoval,
			InitiatedByMapping: valueToCheck.InitiatedByMapping,

			ArbitraryValue: valueToCheck.ArbitraryValue,
		})
	}

	if !IsAddressMappingEmpty(initiatedByMappingAfterRemoval) {
		remaining = append(remaining, &UniversalPermissionDetails{
			TimelineTime:       removedTimelineTimes[0], //We know there is only one because there can only be one interesection between two ranges
			BadgeId:            removedBadges[0],        //We know there is only one because there can only be one interesection between two ranges
			TransferTime:       removedTransferTimes[0], //We know there is only one because there can only be one interesection between two ranges
			OwnershipTime:      removedOwnedTimes[0],
			ToMapping:          toMappingAfterRemoval,
			FromMapping:        fromMappingAfterRemoval,
			InitiatedByMapping: initiatedByMappingAfterRemoval,

			ArbitraryValue: valueToCheck.ArbitraryValue,
		})
	}

	removedDetails := []*UniversalPermissionDetails{}
	for _, removedTimelineTime := range removedTimelineTimes {
		for _, removedBadge := range removedBadges {
			for _, removedTransferTime := range removedTransferTimes {
				for _, removedOwnershipTime := range removedOwnedTimes {
					removedDetails = append(removedDetails, &UniversalPermissionDetails{
						TimelineTime:       removedTimelineTime,
						BadgeId:            removedBadge,
						TransferTime:       removedTransferTime,
						OwnershipTime:      removedOwnershipTime,
						ToMapping:          removedToMapping,
						FromMapping:        removedFromMapping,
						InitiatedByMapping: removedInitiatedByMapping,

						ArbitraryValue: valueToCheck.ArbitraryValue,
					})
				}
			}
		}
	}

	return remaining, removedDetails
}

func GetUintRangesWithOptions(ranges []*UintRange, options *ValueOptions, uses bool) []*UintRange {
	if !uses {
		ranges = []*UintRange{{Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64)}} //dummy range
		return ranges
	}

	if options == nil {
		return ranges
	}

	if options.AllValues {
		ranges = []*UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}}
	}

	if options.InvertDefault {
		ranges = InvertUintRanges(ranges, sdkmath.NewUint(1), sdkmath.NewUint(math.MaxUint64))
	}

	if options.NoValues {
		ranges = []*UintRange{}
	}

	return ranges
}

func GetMappingIdWithOptions(mappingId string, options *ValueOptions, uses bool) string {
	if !uses {
		mappingId = "All"
	}

	if options == nil {
		return mappingId
	}

	if options.AllValues {
		mappingId = "All"
	}

	if options.InvertDefault {
		mappingId = "!" + mappingId
	}

	if options.NoValues {
		mappingId = "None"
	}

	return mappingId
}

func GetMappingWithOptions(mapping *AddressMapping, options *ValueOptions, uses bool) *AddressMapping {
	if !uses {
		mapping = &AddressMapping{Addresses: []string{}, IncludeAddresses: false} //All addresses
	}

	if options == nil {
		return mapping
	}

	if options.AllValues {
		mapping = &AddressMapping{Addresses: []string{}, IncludeAddresses: false} //All addresses
	}

	if options.InvertDefault {
		mapping = &AddressMapping{Addresses: mapping.Addresses, IncludeAddresses: !mapping.IncludeAddresses} //Invert
	}

	if options.NoValues {
		mapping = &AddressMapping{Addresses: []string{}, IncludeAddresses: true} //No addresses
	}

	return mapping
}

func GetFirstMatchOnly(permissions []*UniversalPermission) []*UniversalPermissionDetails {
	handled := []*UniversalPermissionDetails{}
	for _, permission := range permissions {
		for _, combination := range permission.Combinations {
			badgeIds := GetUintRangesWithOptions(permission.DefaultValues.BadgeIds, combination.BadgeIdsOptions, permission.DefaultValues.UsesBadgeIds)
			timelineTimes := GetUintRangesWithOptions(permission.DefaultValues.TimelineTimes, combination.TimelineTimesOptions, permission.DefaultValues.UsesTimelineTimes)
			transferTimes := GetUintRangesWithOptions(permission.DefaultValues.TransferTimes, combination.TransferTimesOptions, permission.DefaultValues.UsesTransferTimes)
			ownedTimes := GetUintRangesWithOptions(permission.DefaultValues.OwnedTimes, combination.OwnedTimesOptions, permission.DefaultValues.UsesOwnedTimes)
			permittedTimes := GetUintRangesWithOptions(permission.DefaultValues.PermittedTimes, combination.PermittedTimesOptions, true)
			forbiddenTimes := GetUintRangesWithOptions(permission.DefaultValues.ForbiddenTimes, combination.ForbiddenTimesOptions, true)
			arbitraryValue := permission.DefaultValues.ArbitraryValue

			toMapping := GetMappingWithOptions(permission.DefaultValues.ToMapping, combination.ToMappingOptions, permission.DefaultValues.UsesToMapping)
			fromMapping := GetMappingWithOptions(permission.DefaultValues.FromMapping, combination.FromMappingOptions, permission.DefaultValues.UsesFromMapping)
			initiatedByMapping := GetMappingWithOptions(permission.DefaultValues.InitiatedByMapping, combination.InitiatedByMappingOptions, permission.DefaultValues.UsesInitiatedByMapping)

			for _, badgeId := range badgeIds {
				for _, timelineTime := range timelineTimes {
					for _, transferTime := range transferTimes {
						for _, ownershipTime := range ownedTimes {
							brokenDown := []*UniversalPermissionDetails{
								{
									BadgeId:            badgeId,
									TimelineTime:       timelineTime,
									TransferTime:       transferTime,
									OwnershipTime:      ownershipTime,
									ToMapping:          toMapping,
									FromMapping:        fromMapping,
									InitiatedByMapping: initiatedByMapping,
								},
							}

							_, remainingAfterHandledIsRemoed, _ := GetOverlapsAndNonOverlaps(brokenDown, handled)
							for _, remaining := range remainingAfterHandledIsRemoed {
								handled = append(handled, &UniversalPermissionDetails{
									TimelineTime:       remaining.TimelineTime,
									BadgeId:            remaining.BadgeId,
									TransferTime:       remaining.TransferTime,
									OwnershipTime:      remaining.OwnershipTime,
									ToMapping:          remaining.ToMapping,
									FromMapping:        remaining.FromMapping,
									InitiatedByMapping: remaining.InitiatedByMapping,

									//Appended for future lookups (not involved in overlap logic)
									PermittedTimes: permittedTimes,
									ForbiddenTimes: forbiddenTimes,
									ArbitraryValue: arbitraryValue,
								})
							}
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
	if permission.BadgeId.Start.Equal(sdkmath.NewUint(math.MaxUint64)) || permission.BadgeId.End.Equal(sdkmath.NewUint(math.MaxUint64)) {
		str += "badgeId: " + permission.BadgeId.Start.String() + " "
	}

	if permission.TimelineTime.Start.Equal(sdkmath.NewUint(math.MaxUint64)) || permission.TimelineTime.End.Equal(sdkmath.NewUint(math.MaxUint64)) {
		str += "timelineTime: " + permission.TimelineTime.Start.String() + " "
	}

	if permission.TransferTime.Start.Equal(sdkmath.NewUint(math.MaxUint64)) || permission.TransferTime.End.Equal(sdkmath.NewUint(math.MaxUint64)) {
		str += "transferTime: " + permission.TransferTime.Start.String() + " "
	}

	if permission.OwnershipTime.Start.Equal(sdkmath.NewUint(math.MaxUint64)) || permission.OwnershipTime.End.Equal(sdkmath.NewUint(math.MaxUint64)) {
		str += "ownershipTime: " + permission.OwnershipTime.Start.String() + " "
	}

	if permission.ToMapping != nil {
		str += "toMapping: "
		if !permission.ToMapping.IncludeAddresses {
			str += fmt.Sprint(len(permission.ToMapping.Addresses)) + " addresses "
		} else {
			str += "all except " + fmt.Sprint(len(permission.ToMapping.Addresses)) + " addresses "
		}

		if len(permission.ToMapping.Addresses) > 0 && len(permission.ToMapping.Addresses) <= 5 {
			str += "("
			for _, address := range permission.ToMapping.Addresses {
				str += address + " "
			}
			str += ")"
		}
	}

	if permission.FromMapping != nil {
		str += "fromMapping: "
		if !permission.FromMapping.IncludeAddresses {
			str += fmt.Sprint(len(permission.FromMapping.Addresses)) + " addresses "
		} else {
			str += "all except " + fmt.Sprint(len(permission.FromMapping.Addresses)) + " addresses "
		}

		if len(permission.FromMapping.Addresses) > 0 && len(permission.FromMapping.Addresses) <= 5 {
			str += "("
			for _, address := range permission.FromMapping.Addresses {
				str += address + " "
			}
			str += ")"
		}
	}

	if permission.InitiatedByMapping != nil {
		str += "initiatedByMapping: "
		if !permission.InitiatedByMapping.IncludeAddresses {
			str += fmt.Sprint(len(permission.InitiatedByMapping.Addresses)) + " addresses "
		} else {
			str += "all except " + fmt.Sprint(len(permission.InitiatedByMapping.Addresses)) + " addresses "
		}

		if len(permission.InitiatedByMapping.Addresses) > 0 && len(permission.InitiatedByMapping.Addresses) <= 5 {
			str += "("
			for _, address := range permission.InitiatedByMapping.Addresses {
				str += address + " "
			}
			str += ")"
		}
	}

	str += ") "

	return str
}

// IMPORTANT PRECONDITION: Must be first match only
func ValidateUniversalPermissionUpdate(oldPermissions []*UniversalPermissionDetails, newPermissions []*UniversalPermissionDetails) error {
	allOverlaps, inOldButNotNew, _ := GetOverlapsAndNonOverlaps(oldPermissions, newPermissions) //we don't care about new not in old

	if len(inOldButNotNew) > 0 {
		errMsg := "permission "
		errMsg += GetPermissionString(inOldButNotNew[0])
		errMsg += "found in old permissions but not in new permissions"
		if len(inOldButNotNew) > 1 {
			errMsg += " (along with " + sdkmath.NewUint(uint64(len(inOldButNotNew)-1)).String() + " more)"
		}

		return sdkerrors.Wrapf(ErrInvalidPermissions, errMsg)
	}

	//For everywhere we detected an overlap, we need to check if the new permissions are valid
	//(i.e. they only explicitly define more permitted or forbidden times and do not remove any)
	for _, overlapObj := range allOverlaps {
		oldPermission := overlapObj.FirstDetails
		newPermission := overlapObj.SecondDetails

		leftoverPermittedTimes, _ := RemoveUintRangeFromUintRange(newPermission.PermittedTimes, oldPermission.PermittedTimes)
		leftoverForbiddenTimes, _ := RemoveUintRangeFromUintRange(newPermission.ForbiddenTimes, oldPermission.ForbiddenTimes)

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


