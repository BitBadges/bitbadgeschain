package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//TODO: This file should probably be refactored a lot, but it currently works.
//			It is also not user-facing or dev-facing, so I am okay with how it is now

//For permissions, we have many types of permissions that are all very similar to each other
//Here, we abstract all those permissions to a UniversalPermission struct in order to reuse code.
//When casting to a UniversalPermission, we use fake dummy values for the unused values to avoid messing up the logic
//
//This file implements certain logic using UniversalPermissions such as overlaps and getting first match only
//This is used in many places around the codebase

type UniversalPermission struct {
	TokenIds        []*UintRange
	TimelineTimes   []*UintRange
	TransferTimes   []*UintRange
	OwnershipTimes  []*UintRange
	ToList          *AddressList
	FromList        *AddressList
	InitiatedByList *AddressList

	ApprovalIdList *AddressList

	PermanentlyPermittedTimes []*UintRange
	PermanentlyForbiddenTimes []*UintRange

	UsesTokenIds        bool
	UsesTimelineTimes   bool
	UsesTransferTimes   bool
	UsesToList          bool
	UsesFromList        bool
	UsesInitiatedByList bool
	UsesOwnershipTimes  bool

	UsesApprovalId bool

	ArbitraryValue interface{}
}

type UniversalPermissionDetails struct {
	TokenId         *UintRange
	TimelineTime    *UintRange
	TransferTime    *UintRange
	OwnershipTime   *UintRange
	ToList          *AddressList
	FromList        *AddressList
	InitiatedByList *AddressList

	ApprovalIdList *AddressList

	//These fields are not actually used in the overlapping logic.
	//They are just along for the ride and used later for lookups
	PermanentlyPermittedTimes []*UintRange
	PermanentlyForbiddenTimes []*UintRange

	ArbitraryValue interface{}
}

type Overlap struct {
	Overlap       *UniversalPermissionDetails
	FirstDetails  *UniversalPermissionDetails
	SecondDetails *UniversalPermissionDetails
}

// Get (overlaps, inOldButNotNew, inNewButNotOld)
func GetOverlapsAndNonOverlaps(ctx sdk.Context, firstDetails, secondDetails []*UniversalPermissionDetails) ([]*Overlap, []*UniversalPermissionDetails, []*UniversalPermissionDetails) {
	inOldButNotNew := make([]*UniversalPermissionDetails, len(firstDetails))
	inNewButNotOld := make([]*UniversalPermissionDetails, len(secondDetails))
	copy(inOldButNotNew, firstDetails)
	copy(inNewButNotOld, secondDetails)

	allOverlaps := []*Overlap{}

	// Find all overlaps between old and new permissions
	for _, oldPerm := range firstDetails {
		for _, newPerm := range secondDetails {
			_, overlaps := UniversalRemoveOverlaps(ctx, newPerm, oldPerm)
			for _, overlap := range overlaps {
				allOverlaps = append(allOverlaps, &Overlap{
					Overlap:       overlap,
					FirstDetails:  oldPerm,
					SecondDetails: newPerm,
				})
			}
		}
	}

	// Remove overlapping sections from both sets
	for _, overlap := range allOverlaps {
		inOldButNotNew, _ = UniversalRemoveOverlapFromValues(ctx, overlap.Overlap, inOldButNotNew)
		inNewButNotOld, _ = UniversalRemoveOverlapFromValues(ctx, overlap.Overlap, inNewButNotOld)
	}

	return allOverlaps, inOldButNotNew, inNewButNotOld
}

// Remove overlap from value -> remaining, removed
func UniversalRemoveOverlapFromValues(ctx sdk.Context, overlap *UniversalPermissionDetails, valuesToCheck []*UniversalPermissionDetails) ([]*UniversalPermissionDetails, []*UniversalPermissionDetails) {
	newValuesToCheck := []*UniversalPermissionDetails{}
	removed := []*UniversalPermissionDetails{}
	for _, base := range valuesToCheck {
		remaining, overlaps := UniversalRemoveOverlaps(ctx, overlap, base)
		newValuesToCheck = append(newValuesToCheck, remaining...)
		removed = append(removed, overlaps...)
	}

	return newValuesToCheck, removed
}

func UniversalRemoveOverlaps(ctx sdk.Context, toRemove *UniversalPermissionDetails, base *UniversalPermissionDetails) ([]*UniversalPermissionDetails, []*UniversalPermissionDetails) {
	if !ctx.IsZero() {
		ctx.GasMeter().ConsumeGas(500, "UniversalRemoveOverlaps")
	}

	toRemove = AddDefaultsIfNil(toRemove)
	base = AddDefaultsIfNil(base)
	remaining := []*UniversalPermissionDetails{}

	timelineTimesAfterRemoval, removedTimelineTimes := RemoveUintRangeFromUintRange(toRemove.TimelineTime, base.TimelineTime)
	if len(removedTimelineTimes) == 0 {
		remaining = append(remaining, base)
		return remaining, []*UniversalPermissionDetails{}
	}

	badgesAfterRemoval, removedBadges := RemoveUintRangeFromUintRange(toRemove.TokenId, base.TokenId)
	if len(removedBadges) == 0 {
		remaining = append(remaining, base)
		return remaining, []*UniversalPermissionDetails{}
	}

	transferTimesAfterRemoval, removedTransferTimes := RemoveUintRangeFromUintRange(toRemove.TransferTime, base.TransferTime)
	if len(removedTransferTimes) == 0 {
		remaining = append(remaining, base)
		return remaining, []*UniversalPermissionDetails{}
	}

	ownershipTimesAfterRemoval, removedOwnershipTimes := RemoveUintRangeFromUintRange(toRemove.OwnershipTime, base.OwnershipTime)
	if len(removedOwnershipTimes) == 0 {
		remaining = append(remaining, base)
		return remaining, []*UniversalPermissionDetails{}
	}

	toListAfterRemoval, removedToList := RemoveAddressListFromAddressList(toRemove.ToList, base.ToList)
	fromListAfterRemoval, removedFromList := RemoveAddressListFromAddressList(toRemove.FromList, base.FromList)
	initiatedByListAfterRemoval, removedInitiatedByList := RemoveAddressListFromAddressList(toRemove.InitiatedByList, base.InitiatedByList)

	approvalIdListAfterRemoval, removedApprovalIdList := RemoveAddressListFromAddressList(toRemove.ApprovalIdList, base.ApprovalIdList)

	toListRemoved := !IsAddressListEmpty(removedToList)
	fromListRemoved := !IsAddressListEmpty(removedFromList)
	initiatedByListRemoved := !IsAddressListEmpty(removedInitiatedByList)
	approvalIdListRemoved := !IsAddressListEmpty(removedApprovalIdList)

	//Approach: Iterate through each field one by one. Attempt to remove the overlap. We'll call each field by an ID number corresponding to its order
	//					Order doesn't matter as long as all fields are toRemove
	//          For each field N, we have the following the cases:
	//            1. For anything remaining after removal for field N is attempted (i.e. the stuff that does not overlap), we need to add
	//               it to the returned array with (0 to N-1) fields filled with removed values and (N+1 to end) fields filled with the original values
	//
	// 						   We only use the removed values of fields 0 to N - 1 because we already toRemove the other fields (via this step in previous iterations)
	//               and we don't want to double count.
	//							 Ex: [0: {1 to 10}, 1: {1 to 10}, 2: {1 to 10}] and we are removing [0: {1 to 5}, 1: {1 to 5}, 2: {1 to 5}]
	//							  	 Let's say we are on field 1. We would add [0: {1 to 5}, 1: {6 to 10}, 2: {1 to 10}] to the returned array
	//						2. If we have removed anything at all, we need to continue to test field N + 1 (i.e. the next field) for overlap
	//							 This is because we have not yet toRemove the cases for values which overlap with field N and field N + 1
	//					  3. If we have not removed anything, we add the original value as outlined in 1) but we do not need to continue to test field N + 1
	//							 because there are no cases untoRemove now where values overlap with field N and field N + 1 becuase nothing overlaps with N.
	//							 If we do end up with this case, it means we end up with the original values because to overlap, it needs to overlap with all fields
	//
	//							 We optimize step 3) by checking right away if something does not overlap with some field. If it does not overlap with some field,
	//							 we can just add the original values and be done with it. If it does overlap with all fields, we need to execute the algorithm

	//If some field does not overlap, we simply end up with the original values because it is only considered an overlap if all fields overlap.
	//The function would work fine without this but it makes it more efficient and less complicated because it will not get broken down further
	if len(removedTimelineTimes) == 0 || len(removedBadges) == 0 || len(removedTransferTimes) == 0 || len(removedOwnershipTimes) == 0 || !toListRemoved || !fromListRemoved || !initiatedByListRemoved || !approvalIdListRemoved {
		remaining = append(remaining, base)
		return remaining, []*UniversalPermissionDetails{}
	}

	for _, timelineTimeAfterRemoval := range timelineTimesAfterRemoval {
		remaining = append(remaining, &UniversalPermissionDetails{
			TimelineTime:    timelineTimeAfterRemoval,
			TokenId:         base.TokenId,
			TransferTime:    base.TransferTime,
			OwnershipTime:   base.OwnershipTime,
			ToList:          base.ToList,
			FromList:        base.FromList,
			InitiatedByList: base.InitiatedByList,
			ApprovalIdList:  base.ApprovalIdList,
			ArbitraryValue:  base.ArbitraryValue,
		})
	}

	for _, badgeAfterRemoval := range badgesAfterRemoval {
		remaining = append(remaining, &UniversalPermissionDetails{
			TimelineTime:    removedTimelineTimes[0], //We know there is only one because there can only be one interesection between two ranges
			TokenId:         badgeAfterRemoval,
			TransferTime:    base.TransferTime,
			OwnershipTime:   base.OwnershipTime,
			ToList:          base.ToList,
			FromList:        base.FromList,
			InitiatedByList: base.InitiatedByList,
			ApprovalIdList:  base.ApprovalIdList,
			ArbitraryValue:  base.ArbitraryValue,
		})
	}

	for _, transferTimeAfterRemoval := range transferTimesAfterRemoval {
		remaining = append(remaining, &UniversalPermissionDetails{
			TimelineTime:    removedTimelineTimes[0], //We know there is only one because there can only be one interesection between two ranges
			TokenId:         removedBadges[0],        //We know there is only one because there can only be one interesection between two ranges
			TransferTime:    transferTimeAfterRemoval,
			OwnershipTime:   base.OwnershipTime,
			ToList:          base.ToList,
			FromList:        base.FromList,
			InitiatedByList: base.InitiatedByList,
			ApprovalIdList:  base.ApprovalIdList,
			ArbitraryValue:  base.ArbitraryValue,
		})
	}

	for _, ownershipTimeAfterRemoval := range ownershipTimesAfterRemoval {
		remaining = append(remaining, &UniversalPermissionDetails{
			TimelineTime:    removedTimelineTimes[0], //We know there is only one because there can only be one interesection between two ranges
			TokenId:         removedBadges[0],        //We know there is only one because there can only be one interesection between two ranges
			TransferTime:    removedTransferTimes[0], //We know there is only one because there can only be one interesection between two ranges
			OwnershipTime:   ownershipTimeAfterRemoval,
			ToList:          base.ToList,
			FromList:        base.FromList,
			InitiatedByList: base.InitiatedByList,
			ApprovalIdList:  base.ApprovalIdList,
			ArbitraryValue:  base.ArbitraryValue,
		})
	}

	if !IsAddressListEmpty(toListAfterRemoval) {
		remaining = append(remaining, &UniversalPermissionDetails{
			TimelineTime:    removedTimelineTimes[0], //We know there is only one because there can only be one interesection between two ranges
			TokenId:         removedBadges[0],        //We know there is only one because there can only be one interesection between two ranges
			TransferTime:    removedTransferTimes[0], //We know there is only one because there can only be one interesection between two ranges
			OwnershipTime:   removedOwnershipTimes[0],
			ToList:          toListAfterRemoval,
			FromList:        base.FromList,
			InitiatedByList: base.InitiatedByList,
			ApprovalIdList:  base.ApprovalIdList,
			ArbitraryValue:  base.ArbitraryValue,
		})
	}

	if !IsAddressListEmpty(fromListAfterRemoval) {
		remaining = append(remaining, &UniversalPermissionDetails{
			TimelineTime:    removedTimelineTimes[0], //We know there is only one because there can only be one interesection between two ranges
			TokenId:         removedBadges[0],        //We know there is only one because there can only be one interesection between two ranges
			TransferTime:    removedTransferTimes[0], //We know there is only one because there can only be one interesection between two ranges
			OwnershipTime:   removedOwnershipTimes[0],
			ToList:          removedToList,
			FromList:        fromListAfterRemoval,
			InitiatedByList: base.InitiatedByList,
			ApprovalIdList:  base.ApprovalIdList,
			ArbitraryValue:  base.ArbitraryValue,
		})
	}

	if !IsAddressListEmpty(initiatedByListAfterRemoval) {
		remaining = append(remaining, &UniversalPermissionDetails{
			TimelineTime:    removedTimelineTimes[0], //We know there is only one because there can only be one interesection between two ranges
			TokenId:         removedBadges[0],        //We know there is only one because there can only be one interesection between two ranges
			TransferTime:    removedTransferTimes[0], //We know there is only one because there can only be one interesection between two ranges
			OwnershipTime:   removedOwnershipTimes[0],
			ToList:          removedToList,
			FromList:        removedFromList,
			InitiatedByList: initiatedByListAfterRemoval,
			ApprovalIdList:  base.ApprovalIdList,
			ArbitraryValue:  base.ArbitraryValue,
		})
	}

	if !IsAddressListEmpty(approvalIdListAfterRemoval) {
		remaining = append(remaining, &UniversalPermissionDetails{
			TimelineTime:    removedTimelineTimes[0],
			TokenId:         removedBadges[0],
			TransferTime:    removedTransferTimes[0],
			OwnershipTime:   removedOwnershipTimes[0],
			ToList:          removedToList,
			FromList:        removedFromList,
			InitiatedByList: removedInitiatedByList,
			ApprovalIdList:  approvalIdListAfterRemoval,

			ArbitraryValue: base.ArbitraryValue,
		})
	}

	removedDetails := []*UniversalPermissionDetails{}
	for _, removedTimelineTime := range removedTimelineTimes {
		for _, removedBadge := range removedBadges {
			for _, removedTransferTime := range removedTransferTimes {
				for _, removedOwnershipTime := range removedOwnershipTimes {
					removedDetails = append(removedDetails, &UniversalPermissionDetails{
						TimelineTime:    removedTimelineTime,
						TokenId:         removedBadge,
						TransferTime:    removedTransferTime,
						OwnershipTime:   removedOwnershipTime,
						ToList:          removedToList,
						FromList:        removedFromList,
						InitiatedByList: removedInitiatedByList,
						ApprovalIdList:  removedApprovalIdList,

						ArbitraryValue: base.ArbitraryValue,
					})
				}
			}
		}
	}

	return remaining, removedDetails
}

func GetFirstMatchOnly(ctx sdk.Context, permissions []*UniversalPermission) []*UniversalPermissionDetails {
	toRemove := []*UniversalPermissionDetails{}
	for _, permission := range permissions {

		tokenIds := GetUintRangesWithOptions(permission.TokenIds, permission.UsesTokenIds)
		timelineTimes := GetUintRangesWithOptions(permission.TimelineTimes, permission.UsesTimelineTimes)
		transferTimes := GetUintRangesWithOptions(permission.TransferTimes, permission.UsesTransferTimes)
		ownershipTimes := GetUintRangesWithOptions(permission.OwnershipTimes, permission.UsesOwnershipTimes)
		permanentlyPermittedTimes := GetUintRangesWithOptions(permission.PermanentlyPermittedTimes, true)
		permanentlyForbiddenTimes := GetUintRangesWithOptions(permission.PermanentlyForbiddenTimes, true)
		arbitraryValue := permission.ArbitraryValue

		toList := GetListWithOptions(permission.ToList, permission.UsesToList)
		fromList := GetListWithOptions(permission.FromList, permission.UsesFromList)
		initiatedByList := GetListWithOptions(permission.InitiatedByList, permission.UsesInitiatedByList)

		approvalIdList := GetListWithOptions(permission.ApprovalIdList, permission.UsesApprovalId)

		for _, tokenId := range tokenIds {
			for _, timelineTime := range timelineTimes {
				for _, transferTime := range transferTimes {
					for _, ownershipTime := range ownershipTimes {
						brokenDown := []*UniversalPermissionDetails{
							{
								TokenId:         tokenId,
								TimelineTime:    timelineTime,
								TransferTime:    transferTime,
								OwnershipTime:   ownershipTime,
								ToList:          toList,
								FromList:        fromList,
								InitiatedByList: initiatedByList,
								ApprovalIdList:  approvalIdList,
							},
						}

						_, inBrokenDownButNottoRemove, _ := GetOverlapsAndNonOverlaps(ctx, brokenDown, toRemove)
						for _, remaining := range inBrokenDownButNottoRemove {
							toRemove = append(toRemove, &UniversalPermissionDetails{
								TimelineTime:    remaining.TimelineTime,
								TokenId:         remaining.TokenId,
								TransferTime:    remaining.TransferTime,
								OwnershipTime:   remaining.OwnershipTime,
								ToList:          remaining.ToList,
								FromList:        remaining.FromList,
								InitiatedByList: remaining.InitiatedByList,
								ApprovalIdList:  remaining.ApprovalIdList,

								//Appended for future lookups (not involved in overlap logic)
								PermanentlyPermittedTimes: permanentlyPermittedTimes,
								PermanentlyForbiddenTimes: permanentlyForbiddenTimes,
								ArbitraryValue:            arbitraryValue,
							})
						}
					}
				}
			}
		}

	}

	return toRemove
}

// IMPORTANT PRECONDITION: Must be first match only
func ValidateUniversalPermissionUpdate(ctx sdk.Context, oldPermissions, newPermissions []*UniversalPermissionDetails) error {
	allOverlaps, inOldButNotNew, _ := GetOverlapsAndNonOverlaps(ctx, oldPermissions, newPermissions)

	// Check none in old but not new
	if err := ValidateNoMissingPermissions(inOldButNotNew); err != nil {
		return err
	}

	// Check no overlapping permissions
	if err := ValidateOverlappingPermissions(allOverlaps); err != nil {
		return err
	}

	return nil
}
