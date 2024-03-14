package types

import (
	fmt "fmt"
	"math"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
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
	BadgeIds           []*UintRange
	TimelineTimes      []*UintRange
	TransferTimes      []*UintRange
	OwnershipTimes     []*UintRange
	ToList          *AddressList
	FromList        *AddressList
	InitiatedByList *AddressList

	ApprovalIdList    *AddressList
	ChallengeTrackerIdList *AddressList
	AmountTrackerIdList *AddressList

	PermanentlyPermittedTimes []*UintRange
	PermanentlyForbiddenTimes []*UintRange

	UsesBadgeIds           bool
	UsesTimelineTimes      bool
	UsesTransferTimes      bool
	UsesToList          bool
	UsesFromList        bool
	UsesInitiatedByList 	 bool
	UsesOwnershipTimes     bool

	UsesApprovalId    bool
	UsesAmountTrackerId bool
	UsesChallengeTrackerId   bool

	ArbitraryValue interface{}
}

type UniversalPermissionDetails struct {
	BadgeId            *UintRange
	TimelineTime       *UintRange
	TransferTime       *UintRange
	OwnershipTime      *UintRange
	ToList          *AddressList
	FromList        *AddressList
	InitiatedByList *AddressList

	ApprovalIdList    *AddressList
	ChallengeTrackerIdList *AddressList
	AmountTrackerIdList *AddressList

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

func GetOverlapsAndNonOverlaps(ctx sdk.Context,	firstDetails []*UniversalPermissionDetails, secondDetails []*UniversalPermissionDetails) ([]*Overlap, []*UniversalPermissionDetails, []*UniversalPermissionDetails) {
	//TODO: deep copy???
	inOldButNotNew := make([]*UniversalPermissionDetails, len(firstDetails))
	copy(inOldButNotNew, firstDetails)
	inNewButNotOld := make([]*UniversalPermissionDetails, len(secondDetails))
	copy(inNewButNotOld, secondDetails)

	allOverlaps := []*Overlap{}
	for _, oldPermission := range firstDetails {
		for _, newPermission := range secondDetails {
			_, overlaps := UniversalRemoveOverlaps(ctx, newPermission, oldPermission)
			// fmt.Println("UniversalRemoveOverlaps", time.Since(startTime))
			for _, overlap := range overlaps {
				allOverlaps = append(allOverlaps, &Overlap{
					Overlap:       overlap,
					FirstDetails:  oldPermission,
					SecondDetails: newPermission,
				})
			}
		}
	}

	for _, overlap := range allOverlaps {
		inOldButNotNew, _ = UniversalRemoveOverlapFromValues(ctx, overlap.Overlap, inOldButNotNew)
		inNewButNotOld, _ = UniversalRemoveOverlapFromValues(ctx, overlap.Overlap, inNewButNotOld)
	}

	return allOverlaps, inOldButNotNew, inNewButNotOld
}

func UniversalRemoveOverlapFromValues(ctx sdk.Context, handled *UniversalPermissionDetails, valuesToCheck []*UniversalPermissionDetails) ([]*UniversalPermissionDetails, []*UniversalPermissionDetails) {
	newValuesToCheck := []*UniversalPermissionDetails{}
	removed := []*UniversalPermissionDetails{}
	for _, valueToCheck := range valuesToCheck {
		remaining, overlaps := UniversalRemoveOverlaps(ctx, handled, valueToCheck)
		newValuesToCheck = append(newValuesToCheck, remaining...)
		removed = append(removed, overlaps...)
	}

	return newValuesToCheck, removed
}

func IsAddressListEmpty(list *AddressList) bool {
	return len(list.Addresses) == 0 && list.Whitelist
}

func UniversalRemoveOverlaps(ctx sdk.Context, handled *UniversalPermissionDetails, valueToCheck *UniversalPermissionDetails) ([]*UniversalPermissionDetails, []*UniversalPermissionDetails) {
	if !ctx.IsZero() {
		ctx.GasMeter().ConsumeGas(500, "UniversalRemoveOverlaps")
	}

	remaining := []*UniversalPermissionDetails{}

	timelineTimesAfterRemoval, removedTimelineTimes := RemoveUintRangeFromUintRange(handled.TimelineTime, valueToCheck.TimelineTime)
	if len(removedTimelineTimes) == 0 {
		remaining = append(remaining, valueToCheck)
		return remaining, []*UniversalPermissionDetails{}
	}


	badgesAfterRemoval, removedBadges := RemoveUintRangeFromUintRange(handled.BadgeId, valueToCheck.BadgeId)
	if len(removedBadges) == 0 {
		remaining = append(remaining, valueToCheck)
		return remaining, []*UniversalPermissionDetails{}
	}
	
	transferTimesAfterRemoval, removedTransferTimes := RemoveUintRangeFromUintRange(handled.TransferTime, valueToCheck.TransferTime)
	if len(removedTransferTimes) == 0 {
		remaining = append(remaining, valueToCheck)
		return remaining, []*UniversalPermissionDetails{}
	}
	
	ownershipTimesAfterRemoval, removedOwnershipTimes := RemoveUintRangeFromUintRange(handled.OwnershipTime, valueToCheck.OwnershipTime)
	if len(removedOwnershipTimes) == 0 {
		remaining = append(remaining, valueToCheck)
		return remaining, []*UniversalPermissionDetails{}
	}
	
	toListAfterRemoval, removedToList := RemoveAddressListFromAddressList(handled.ToList, valueToCheck.ToList)
	fromListAfterRemoval, removedFromList := RemoveAddressListFromAddressList(handled.FromList, valueToCheck.FromList)
	initiatedByListAfterRemoval, removedInitiatedByList := RemoveAddressListFromAddressList(handled.InitiatedByList, valueToCheck.InitiatedByList)
	
	approvalIdListAfterRemoval, removedApprovalIdList := RemoveAddressListFromAddressList(handled.ApprovalIdList, valueToCheck.ApprovalIdList)
	amountTrackerIdListAfterRemoval, removedAmountTrackerIdList := RemoveAddressListFromAddressList(handled.AmountTrackerIdList, valueToCheck.AmountTrackerIdList)
	challengeTrackerIdListAfterRemoval, removedChallengeTrackerIdList := RemoveAddressListFromAddressList(handled.ChallengeTrackerIdList, valueToCheck.ChallengeTrackerIdList)
	
	toListRemoved := !IsAddressListEmpty(removedToList)
	fromListRemoved := !IsAddressListEmpty(removedFromList)
	initiatedByListRemoved := !IsAddressListEmpty(removedInitiatedByList)
	approvalIdListRemoved := !IsAddressListEmpty(removedApprovalIdList)
	amountTrackerIdListRemoved := !IsAddressListEmpty(removedAmountTrackerIdList)
	challengeTrackerIdListRemoved := !IsAddressListEmpty(removedChallengeTrackerIdList)

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

	
	//If some field does not overlap, we simply end up with the original values because it is only considered an overlap if all fields overlap.
	//The function would work fine without this but it makes it more efficient and less complicated because it will not get broken down further
	if len(removedTimelineTimes) == 0 || len(removedBadges) == 0 || len(removedTransferTimes) == 0 || len(removedOwnershipTimes) == 0 || !toListRemoved || !fromListRemoved || !initiatedByListRemoved || !approvalIdListRemoved || !amountTrackerIdListRemoved || !challengeTrackerIdListRemoved {
		remaining = append(remaining, valueToCheck)
		return remaining, []*UniversalPermissionDetails{}
	}

	for _, timelineTimeAfterRemoval := range timelineTimesAfterRemoval {
		remaining = append(remaining, &UniversalPermissionDetails{
			TimelineTime:              timelineTimeAfterRemoval,
			BadgeId:                   valueToCheck.BadgeId,
			TransferTime:              valueToCheck.TransferTime,
			OwnershipTime:             valueToCheck.OwnershipTime,
			ToList:                 valueToCheck.ToList,
			FromList:               valueToCheck.FromList,
			InitiatedByList:        valueToCheck.InitiatedByList,
			ApprovalIdList: 			 valueToCheck.ApprovalIdList,
			AmountTrackerIdList: 	 valueToCheck.AmountTrackerIdList,
			ChallengeTrackerIdList:  valueToCheck.ChallengeTrackerIdList,
			ArbitraryValue: valueToCheck.ArbitraryValue,
		})
	}

	for _, badgeAfterRemoval := range badgesAfterRemoval {
		remaining = append(remaining, &UniversalPermissionDetails{
			TimelineTime:              removedTimelineTimes[0], //We know there is only one because there can only be one interesection between two ranges
			BadgeId:                   badgeAfterRemoval,
			TransferTime:              valueToCheck.TransferTime,
			OwnershipTime:             valueToCheck.OwnershipTime,
			ToList:                 valueToCheck.ToList,
			FromList:               valueToCheck.FromList,
			InitiatedByList:        valueToCheck.InitiatedByList,
			ApprovalIdList: valueToCheck.ApprovalIdList,
			AmountTrackerIdList: 	 valueToCheck.AmountTrackerIdList,
			ChallengeTrackerIdList:  valueToCheck.ChallengeTrackerIdList,
			ArbitraryValue: valueToCheck.ArbitraryValue,
		})
	}

	for _, transferTimeAfterRemoval := range transferTimesAfterRemoval {
		remaining = append(remaining, &UniversalPermissionDetails{
			TimelineTime:              removedTimelineTimes[0], //We know there is only one because there can only be one interesection between two ranges
			BadgeId:                   removedBadges[0],        //We know there is only one because there can only be one interesection between two ranges
			TransferTime:              transferTimeAfterRemoval,
			OwnershipTime:             valueToCheck.OwnershipTime,
			ToList:                 valueToCheck.ToList,
			FromList:               valueToCheck.FromList,
			InitiatedByList:        valueToCheck.InitiatedByList,
			ApprovalIdList: valueToCheck.ApprovalIdList,
			AmountTrackerIdList: 	 valueToCheck.AmountTrackerIdList,
			ChallengeTrackerIdList:  valueToCheck.ChallengeTrackerIdList,
			ArbitraryValue: valueToCheck.ArbitraryValue,
		})
	}

	for _, ownershipTimeAfterRemoval := range ownershipTimesAfterRemoval {
		remaining = append(remaining, &UniversalPermissionDetails{
			TimelineTime:              removedTimelineTimes[0], //We know there is only one because there can only be one interesection between two ranges
			BadgeId:                   removedBadges[0],        //We know there is only one because there can only be one interesection between two ranges
			TransferTime:              removedTransferTimes[0], //We know there is only one because there can only be one interesection between two ranges
			OwnershipTime:             ownershipTimeAfterRemoval,
			ToList:                 valueToCheck.ToList,
			FromList:               valueToCheck.FromList,
			InitiatedByList:        valueToCheck.InitiatedByList,
			ApprovalIdList: valueToCheck.ApprovalIdList,
			AmountTrackerIdList: 	 valueToCheck.AmountTrackerIdList,
			ChallengeTrackerIdList:  valueToCheck.ChallengeTrackerIdList,
			ArbitraryValue: valueToCheck.ArbitraryValue,
		})
	}

	if !IsAddressListEmpty(toListAfterRemoval) {
		remaining = append(remaining, &UniversalPermissionDetails{
			TimelineTime:              removedTimelineTimes[0], //We know there is only one because there can only be one interesection between two ranges
			BadgeId:                   removedBadges[0],        //We know there is only one because there can only be one interesection between two ranges
			TransferTime:              removedTransferTimes[0], //We know there is only one because there can only be one interesection between two ranges
			OwnershipTime:             removedOwnershipTimes[0],
			ToList:                 toListAfterRemoval,
			FromList:               valueToCheck.FromList,
			InitiatedByList:        valueToCheck.InitiatedByList,
			ApprovalIdList: valueToCheck.ApprovalIdList,
			AmountTrackerIdList: 	 valueToCheck.AmountTrackerIdList,
			ChallengeTrackerIdList:  valueToCheck.ChallengeTrackerIdList,
			ArbitraryValue: valueToCheck.ArbitraryValue,
		})
	}

	if !IsAddressListEmpty(fromListAfterRemoval) {
		remaining = append(remaining, &UniversalPermissionDetails{
			TimelineTime:              removedTimelineTimes[0], //We know there is only one because there can only be one interesection between two ranges
			BadgeId:                   removedBadges[0],        //We know there is only one because there can only be one interesection between two ranges
			TransferTime:              removedTransferTimes[0], //We know there is only one because there can only be one interesection between two ranges
			OwnershipTime:             removedOwnershipTimes[0],
			ToList:                 removedToList,
			FromList:               fromListAfterRemoval,
			InitiatedByList:        valueToCheck.InitiatedByList,
			ApprovalIdList: valueToCheck.ApprovalIdList,
			AmountTrackerIdList: 	 valueToCheck.AmountTrackerIdList,
			ChallengeTrackerIdList:  valueToCheck.ChallengeTrackerIdList,
			ArbitraryValue: valueToCheck.ArbitraryValue,
		})
	}

	if !IsAddressListEmpty(initiatedByListAfterRemoval) {
		remaining = append(remaining, &UniversalPermissionDetails{
			TimelineTime:              removedTimelineTimes[0], //We know there is only one because there can only be one interesection between two ranges
			BadgeId:                   removedBadges[0],        //We know there is only one because there can only be one interesection between two ranges
			TransferTime:              removedTransferTimes[0], //We know there is only one because there can only be one interesection between two ranges
			OwnershipTime:             removedOwnershipTimes[0],
			ToList:                 removedToList,
			FromList:               removedFromList,
			InitiatedByList:        initiatedByListAfterRemoval,
			ApprovalIdList: valueToCheck.ApprovalIdList,
			AmountTrackerIdList: 	 valueToCheck.AmountTrackerIdList,
			ChallengeTrackerIdList:  valueToCheck.ChallengeTrackerIdList,
			ArbitraryValue: valueToCheck.ArbitraryValue,
		})
	}

	if !IsAddressListEmpty(approvalIdListAfterRemoval) {
		remaining = append(remaining, &UniversalPermissionDetails{
			TimelineTime:              removedTimelineTimes[0],
			BadgeId:                   removedBadges[0],
			TransferTime:              removedTransferTimes[0],
			OwnershipTime:             removedOwnershipTimes[0],
			ToList:                 removedToList,
			FromList:               removedFromList,
			InitiatedByList:        removedInitiatedByList,
			ApprovalIdList: approvalIdListAfterRemoval,
			AmountTrackerIdList: 	 valueToCheck.AmountTrackerIdList,
			ChallengeTrackerIdList:  valueToCheck.ChallengeTrackerIdList,

			ArbitraryValue: valueToCheck.ArbitraryValue,
		})
	}

	if !IsAddressListEmpty(amountTrackerIdListAfterRemoval) {
		remaining = append(remaining, &UniversalPermissionDetails{
			TimelineTime:              removedTimelineTimes[0],
			BadgeId:                   removedBadges[0],
			TransferTime:              removedTransferTimes[0],
			OwnershipTime:             removedOwnershipTimes[0],
			ToList:                 removedToList,
			FromList:               removedFromList,
			InitiatedByList:        removedInitiatedByList,
			ApprovalIdList: removedApprovalIdList,
			AmountTrackerIdList: amountTrackerIdListAfterRemoval,
			ChallengeTrackerIdList:  valueToCheck.ChallengeTrackerIdList,

			ArbitraryValue: valueToCheck.ArbitraryValue,
		})
	}

	if !IsAddressListEmpty(challengeTrackerIdListAfterRemoval) {
		remaining = append(remaining, &UniversalPermissionDetails{
			TimelineTime:              removedTimelineTimes[0],
			BadgeId:                   removedBadges[0],
			TransferTime:              removedTransferTimes[0],
			OwnershipTime:             removedOwnershipTimes[0],
			ToList:                 removedToList,
			FromList:               removedFromList,
			InitiatedByList:        removedInitiatedByList,
			ApprovalIdList: removedApprovalIdList,
			AmountTrackerIdList: removedAmountTrackerIdList,
			ChallengeTrackerIdList:  challengeTrackerIdListAfterRemoval,

			ArbitraryValue: valueToCheck.ArbitraryValue,
		})
	}

	removedDetails := []*UniversalPermissionDetails{}
	for _, removedTimelineTime := range removedTimelineTimes {
		for _, removedBadge := range removedBadges {
			for _, removedTransferTime := range removedTransferTimes {
				for _, removedOwnershipTime := range removedOwnershipTimes {
					removedDetails = append(removedDetails, &UniversalPermissionDetails{
						TimelineTime:              removedTimelineTime,
						BadgeId:                   removedBadge,
						TransferTime:              removedTransferTime,
						OwnershipTime:             removedOwnershipTime,
						ToList:                 removedToList,
						FromList:               removedFromList,
						InitiatedByList:        removedInitiatedByList,
						ApprovalIdList: removedApprovalIdList,
						AmountTrackerIdList: removedAmountTrackerIdList,
						ChallengeTrackerIdList: removedChallengeTrackerIdList,

						ArbitraryValue: valueToCheck.ArbitraryValue,
					})
				}
			}
		}
	}

	return remaining, removedDetails
}

func GetUintRangesWithOptions(ranges []*UintRange, uses bool) []*UintRange {
	if !uses {
		ranges = []*UintRange{{Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64)}} //dummy range
		return ranges
	} else {
		return ranges
	}
}

func GetListIdWithOptions(listId string, uses bool) string {
	if !uses {
		listId = "All"
		return listId
	} else {
		return listId
	}
}

func GetListWithOptions(list *AddressList, uses bool) *AddressList {
	if !uses {
		list = &AddressList{Addresses: []string{}, Whitelist: false} //All addresses
	}

	return list
}

func ApplyManipulations(permissions []*UniversalPermission) []*UniversalPermissionDetails {
	handled := []*UniversalPermissionDetails{}
	for _, permission := range permissions {
		badgeIds := GetUintRangesWithOptions(permission.BadgeIds, permission.UsesBadgeIds)
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
		amountTrackerIdList := GetListWithOptions(permission.AmountTrackerIdList, permission.UsesAmountTrackerId)
		challengeTrackerIdList := GetListWithOptions(permission.ChallengeTrackerIdList, permission.UsesChallengeTrackerId)

		for _, badgeId := range badgeIds {
			for _, timelineTime := range timelineTimes {
				for _, transferTime := range transferTimes {
					for _, ownershipTime := range ownershipTimes {
						brokenDown := []*UniversalPermissionDetails{
							{
								BadgeId:                   badgeId,
								TimelineTime:              timelineTime,
								TransferTime:              transferTime,
								OwnershipTime:             ownershipTime,
								ToList:                 toList,
								FromList:               fromList,
								InitiatedByList:        initiatedByList,
								ApprovalIdList: approvalIdList,
								AmountTrackerIdList: amountTrackerIdList,
								ChallengeTrackerIdList: challengeTrackerIdList,

								//Appended for future lookups (not involved in overlap logic)
								PermanentlyPermittedTimes: permanentlyPermittedTimes,
								PermanentlyForbiddenTimes: permanentlyForbiddenTimes,
								ArbitraryValue: arbitraryValue,
							},
						}

						handled = append(handled, brokenDown...)
					}
				}
			}
		}

	}

	return handled
}

func GetFirstMatchOnly(ctx sdk.Context, permissions []*UniversalPermission) []*UniversalPermissionDetails {
	handled := []*UniversalPermissionDetails{}
	for _, permission := range permissions {

		badgeIds := GetUintRangesWithOptions(permission.BadgeIds, permission.UsesBadgeIds)
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
		amountTrackerIdList := GetListWithOptions(permission.AmountTrackerIdList, permission.UsesAmountTrackerId)
		challengeTrackerIdList := GetListWithOptions(permission.ChallengeTrackerIdList, permission.UsesChallengeTrackerId)

		for _, badgeId := range badgeIds {
			for _, timelineTime := range timelineTimes {
				for _, transferTime := range transferTimes {
					for _, ownershipTime := range ownershipTimes {
						brokenDown := []*UniversalPermissionDetails{
							{
								BadgeId:                   badgeId,
								TimelineTime:              timelineTime,
								TransferTime:              transferTime,
								OwnershipTime:             ownershipTime,
								ToList:                 toList,
								FromList:               fromList,
								InitiatedByList:        initiatedByList,
								ApprovalIdList: approvalIdList,
								AmountTrackerIdList: amountTrackerIdList,
								ChallengeTrackerIdList: challengeTrackerIdList,
							},
						}

						_, inBrokenDownButNotHandled, _ := GetOverlapsAndNonOverlaps(ctx, brokenDown, handled)
						for _, remaining := range inBrokenDownButNotHandled {
							handled = append(handled, &UniversalPermissionDetails{
								TimelineTime:              remaining.TimelineTime,
								BadgeId:                   remaining.BadgeId,
								TransferTime:              remaining.TransferTime,
								OwnershipTime:             remaining.OwnershipTime,
								ToList:                 remaining.ToList,
								FromList:               remaining.FromList,
								InitiatedByList:        remaining.InitiatedByList,
								ApprovalIdList: remaining.ApprovalIdList,
								AmountTrackerIdList: remaining.AmountTrackerIdList,
								ChallengeTrackerIdList: remaining.ChallengeTrackerIdList,

								//Appended for future lookups (not involved in overlap logic)
								PermanentlyPermittedTimes: permanentlyPermittedTimes,
								PermanentlyForbiddenTimes: permanentlyForbiddenTimes,
								ArbitraryValue: arbitraryValue,
							})
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

	if permission.ToList != nil {
		str += "toList: "
		if !permission.ToList.Whitelist {
			str += fmt.Sprint(len(permission.ToList.Addresses)) + " addresses "
		} else {
			str += "all except " + fmt.Sprint(len(permission.ToList.Addresses)) + " addresses "
		}

		if len(permission.ToList.Addresses) > 0 && len(permission.ToList.Addresses) <= 5 {
			str += "("
			for _, address := range permission.ToList.Addresses {
				str += address + " "
			}
			str += ")"
		}
	}

	if permission.FromList != nil {
		str += "fromList: "
		if !permission.FromList.Whitelist {
			str += fmt.Sprint(len(permission.FromList.Addresses)) + " addresses "
		} else {
			str += "all except " + fmt.Sprint(len(permission.FromList.Addresses)) + " addresses "
		}

		if len(permission.FromList.Addresses) > 0 && len(permission.FromList.Addresses) <= 5 {
			str += "("
			for _, address := range permission.FromList.Addresses {
				str += address + " "
			}
			str += ")"
		}
	}

	if permission.InitiatedByList != nil {
		str += "initiatedByList: "
		if !permission.InitiatedByList.Whitelist {
			str += fmt.Sprint(len(permission.InitiatedByList.Addresses)) + " addresses "
		} else {
			str += "all except " + fmt.Sprint(len(permission.InitiatedByList.Addresses)) + " addresses "
		}

		if len(permission.InitiatedByList.Addresses) > 0 && len(permission.InitiatedByList.Addresses) <= 5 {
			str += "("
			for _, address := range permission.InitiatedByList.Addresses {
				str += address + " "
			}
			str += ")"
		}
	}

	str += ") "

	return str
}

// IMPORTANT PRECONDITION: Must be first match only
func ValidateUniversalPermissionUpdate(ctx sdk.Context, oldPermissions []*UniversalPermissionDetails, newPermissions []*UniversalPermissionDetails) error {
	allOverlaps, inOldButNotNew, _ := GetOverlapsAndNonOverlaps(ctx,  oldPermissions, newPermissions) //we don't care about new not in old

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

		leftoverPermanentlyPermittedTimes, _ := RemoveUintRangesFromUintRanges(newPermission.PermanentlyPermittedTimes, oldPermission.PermanentlyPermittedTimes)
		leftoverPermanentlyForbiddenTimes, _ := RemoveUintRangesFromUintRanges(newPermission.PermanentlyForbiddenTimes, oldPermission.PermanentlyForbiddenTimes)

		if len(leftoverPermanentlyPermittedTimes) > 0 || len(leftoverPermanentlyForbiddenTimes) > 0 {
			errMsg := "permission "
			errMsg += GetPermissionString(oldPermission)
			errMsg += "found in both new and old permissions but "
			if len(leftoverPermanentlyPermittedTimes) > 0 {
				errMsg += "previously explicitly allowed the times ( "
				for _, oldPermittedTime := range leftoverPermanentlyPermittedTimes {
					errMsg += oldPermittedTime.Start.String() + "-" + oldPermittedTime.End.String() + " "
				}
				errMsg += ") which are now set to disallowed"
			}
			if len(leftoverPermanentlyForbiddenTimes) > 0 && len(leftoverPermanentlyPermittedTimes) > 0 {
				errMsg += " and"
			}
			if len(leftoverPermanentlyForbiddenTimes) > 0 {
				errMsg += " previously explicitly disallowed the times ( "
				for _, oldForbiddenTime := range leftoverPermanentlyForbiddenTimes {
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
