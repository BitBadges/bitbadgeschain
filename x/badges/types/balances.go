package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"math"
)

//TODO: should really combine times and ranges and probably amount into one struct - Balance


func HandleDuplicateBadgeIds(balances []*Balance) ([]*Balance, error) {
	newBalances := []*Balance{}
	err := *new(error)
	for _, balance := range balances {
		for _, badgeId := range balance.BadgeIds {
			newBalances, err = AddBalancesForIdRanges(newBalances, []*IdRange{badgeId}, balance.Times, balance.Amount)
			if err != nil {
				return newBalances, err
			}
		}
	}

	return newBalances, nil
}

// Gets the balances for a specific ID. Assumes balances are sorted
func GetBalancesForId(id sdk.Uint, balances []*Balance) []*Balance {
	matchingBalances := []*Balance{}
	for _, balance := range balances {
		found := SearchIdRangesForId(id, balance.BadgeIds)
		if found {
			matchingBalances = append(matchingBalances, &Balance{
				Amount:   balance.Amount,
				BadgeIds: []*IdRange{&IdRange{Start: id, End: id}},
				Times: balance.Times,
			})
		}
	}
	return matchingBalances
}

// Updates the balance for a specific ids from what it currently is to newAmount.
func UpdateBalancesForIdRanges(ranges []*IdRange, newTimes []*IdRange, newAmount sdk.Uint,  balances []*Balance) (newBalances []*Balance, err error) {
	//Can maybe optimize this in the future by doing this all in one loop instead of deleting then setting
	// ranges = SortAndMergeOverlapping(ranges)
	balances = DeleteBalanceForIdRanges(ranges, newTimes, balances)
	balances, err = SetBalanceForIdRanges(ranges, newTimes,  newAmount, balances)

	return balances, err
}

// Gets the balances for specified ID ranges. Returns a new []*Balance where only the specified ID ranges and their balances are included. Appends balance == 0 objects so all IDs are accounted for, even if not found.
func GetBalancesForIdRanges(idRanges []*IdRange, times []*IdRange, balances []*Balance) (newBalances []*Balance, err error) {
	fetchedBalances := []*Balance{}

	for _, balanceObj := range balances {
		currPermissionDetails := []*UniversalPermissionDetails{}
		for _, currRange := range balanceObj.BadgeIds {
			for _, currTime := range balanceObj.Times {
				currPermissionDetails = append(currPermissionDetails, &UniversalPermissionDetails{
					BadgeId: currRange,
					TimelineTime: currTime,
					TransferTime: &IdRange{ Start: sdk.NewUint(math.MaxUint64), End: sdk.NewUint(math.MaxUint64) }, //dummy range
				})
			}
		}

		toFetchPermissionDetails := []*UniversalPermissionDetails{}
		for _, rangeToFetch := range idRanges {
			for _, timeToFetch := range times {
				toFetchPermissionDetails = append(toFetchPermissionDetails, &UniversalPermissionDetails{
						BadgeId: rangeToFetch,
						TimelineTime: timeToFetch,
						TransferTime: &IdRange{ Start: sdk.NewUint(math.MaxUint64), End: sdk.NewUint(math.MaxUint64) }, //dummy range
					},
				)
			}
		}

		overlaps, _, inNewButNotOld := GetOverlapsAndNonOverlaps(currPermissionDetails, toFetchPermissionDetails)
		for _, overlapObject := range overlaps {
			overlap := overlapObject.Overlap

			fetchedBalances = append(fetchedBalances, &Balance{
				Amount:   balanceObj.Amount,
				BadgeIds: []*IdRange{overlap.BadgeId},
				Times: []*IdRange{overlap.TimelineTime},
			})
		}

		for _, detail := range inNewButNotOld {
			fetchedBalances = append(fetchedBalances, &Balance{
				Amount:   balanceObj.Amount,
				BadgeIds: []*IdRange{detail.BadgeId},
				Times: []*IdRange{detail.TimelineTime},
			})
		}
	}

	return fetchedBalances, nil
}

// Adds a balance to all ids specified in []ranges
func AddBalancesForIdRanges(balances []*Balance, ranges []*IdRange, times []*IdRange, balanceToAdd sdk.Uint) ([]*Balance, error) {
	currBalances, err := GetBalancesForIdRanges(ranges, times, balances)
	if err != nil {
		return balances, err
	}

	for _, balance := range currBalances {
		newBalanceAmount, err := SafeAdd(balance.Amount, balanceToAdd)
		if err != nil {
			return balances, err
		}

		balances, err = UpdateBalancesForIdRanges(balance.BadgeIds, times, newBalanceAmount, balances)
		if err != nil {
			return balances, err
		}
	}
	return balances, nil
}

// Subtracts a balance to all ids specified in []ranges
func SubtractBalancesForIdRanges(balances []*Balance, ranges []*IdRange, times []*IdRange, balanceToRemove sdk.Uint) ([]*Balance, error) {
	currBalances, err := GetBalancesForIdRanges(ranges, times, balances)
	if err != nil {
		return balances, err
	}

	for _, currBalanceObj := range currBalances {
		newBalance, err := SafeSubtract(currBalanceObj.Amount, balanceToRemove)
		if err != nil {
			return balances, err
		}

		balances, err = UpdateBalancesForIdRanges(currBalanceObj.BadgeIds, times, newBalance, balances)
		if err != nil {
			return balances, err
		}
	}
	return balances, nil
}

// Deletes the balance for a specific id.
func DeleteBalanceForIdRanges(rangesToDelete []*IdRange, timesToDelete []*IdRange, balances []*Balance) []*Balance {
	newBalances := []*Balance{}

	for _, balanceObj := range balances {
		currPermissionDetails := []*UniversalPermissionDetails{}
		for _, currRange := range balanceObj.BadgeIds {
			for _, currTime := range balanceObj.Times {
				currPermissionDetails = append(currPermissionDetails, &UniversalPermissionDetails{
					BadgeId: currRange,
					TimelineTime: currTime,
					TransferTime: &IdRange{ Start: sdk.NewUint(math.MaxUint64), End: sdk.NewUint(math.MaxUint64) }, //dummy range
				})
			}
		}

		toDeletePermissionDetails := []*UniversalPermissionDetails{}
		for _, rangeToDelete := range rangesToDelete {
			for _, timeToDelete := range timesToDelete {
				toDeletePermissionDetails = append(toDeletePermissionDetails, &UniversalPermissionDetails{
						BadgeId: rangeToDelete,
						TimelineTime: timeToDelete,
						TransferTime: &IdRange{ Start: sdk.NewUint(math.MaxUint64), End: sdk.NewUint(math.MaxUint64) }, //dummy range
					},
				)
			}
		}

		_, inOldButNotNew, _ := GetOverlapsAndNonOverlaps(currPermissionDetails, toDeletePermissionDetails)
		for _, remainingBalance := range inOldButNotNew {
			newBalances = append(newBalances, &Balance{
				Amount:   balanceObj.Amount,
				BadgeIds: []*IdRange{remainingBalance.BadgeId},
				Times: []*IdRange{remainingBalance.TimelineTime},
			})
		}
	}

	return newBalances
}

// Sets the balance for a specific id. Assumes balance does not exist.
func SetBalanceForIdRanges(ranges []*IdRange, times []*IdRange, amount sdk.Uint, balances []*Balance) ([]*Balance, error) {
	//TODO: ?
	// err := *new(error)
	
	// for _, balance := range balances {
	// 	if balance.Amount.Equal(amount) {
	// 		if IdRangeEquals(times, balance.Times) {
	// 			for _, rangeToAdd := range ranges {
	// 				balance.BadgeIds, err = InsertRangeToIdRanges(rangeToAdd, balance.BadgeIds)
	// 				if err != nil {
	// 					return nil, err
	// 				}
	// 			}

	// 			return balances, nil
	// 		}
	// 	}
	// }

	balances = append(balances, &Balance{
		Amount:   amount,
		BadgeIds: ranges,
		Times: times,
	})

	return balances, nil
}
