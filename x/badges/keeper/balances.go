package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"math"
)

//TODO: should really combine times and ranges and probably amount into one struct - Balance

// Gets the balances for a specific ID. Assumes balances are sorted
func GetBalancesForId(id sdk.Uint, balances []*types.Balance) []*types.Balance {
	matchingBalances := []*types.Balance{}
	for _, balance := range balances {
		_, found := types.SearchIdRangesForId(id, balance.BadgeIds)
		if found {
			matchingBalances = append(matchingBalances, &types.Balance{
				Amount:   balance.Amount,
				BadgeIds: []*types.IdRange{&types.IdRange{Start: id, End: id}},
				Times: balance.Times,
			})
		}
	}
	return matchingBalances
}

// Updates the balance for a specific ids from what it currently is to newAmount.
func UpdateBalancesForIdRanges(ranges []*types.IdRange, newTimes []*types.IdRange, newAmount sdk.Uint,  balances []*types.Balance) (newBalances []*types.Balance, err error) {
	//Can maybe optimize this in the future by doing this all in one loop instead of deleting then setting
	ranges = types.SortAndMergeOverlapping(ranges)
	balances = DeleteBalanceForIdRanges(ranges, newTimes, balances)
	balances, err = SetBalanceForIdRanges(ranges, newTimes,  newAmount, balances)

	return balances, err
}

// Gets the balances for specified ID ranges. Returns a new []*types.Balance where only the specified ID ranges and their balances are included. Appends balance == 0 objects so all IDs are accounted for, even if not found.
func GetBalancesForIdRanges(idRanges []*types.IdRange, times []*types.IdRange, balances []*types.Balance) (newBalances []*types.Balance, err error) {
	fetchedBalances := []*types.Balance{}

	for _, balanceObj := range balances {
		currPermissionDetails := []*types.UniversalPermissionDetails{}
		for _, currRange := range balanceObj.BadgeIds {
			for _, currTime := range balanceObj.Times {
				currPermissionDetails = append(currPermissionDetails, &types.UniversalPermissionDetails{
					BadgeId: currRange,
					TimelineTime: currTime,
					TransferTime: &types.IdRange{ Start: sdk.NewUint(math.MaxUint64), End: sdk.NewUint(math.MaxUint64) }, //dummy range
				})
			}
		}

		toFetchPermissionDetails := []*types.UniversalPermissionDetails{}
		for _, rangeToFetch := range idRanges {
			for _, timeToFetch := range times {
				toFetchPermissionDetails = append(toFetchPermissionDetails, &types.UniversalPermissionDetails{
						BadgeId: rangeToFetch,
						TimelineTime: timeToFetch,
						TransferTime: &types.IdRange{ Start: sdk.NewUint(math.MaxUint64), End: sdk.NewUint(math.MaxUint64) }, //dummy range
					},
				)
			}
		}

		overlaps, _, inNewButNotOld := types.GetOverlapsAndNonOverlaps(currPermissionDetails, toFetchPermissionDetails)
		for _, overlapObject := range overlaps {
			overlap := overlapObject.Overlap

			fetchedBalances = append(fetchedBalances, &types.Balance{
				Amount:   balanceObj.Amount,
				BadgeIds: []*types.IdRange{overlap.BadgeId},
				Times: []*types.IdRange{overlap.TimelineTime},
			})
		}

		for _, detail := range inNewButNotOld {
			fetchedBalances = append(fetchedBalances, &types.Balance{
				Amount:   balanceObj.Amount,
				BadgeIds: []*types.IdRange{detail.BadgeId},
				Times: []*types.IdRange{detail.TimelineTime},
			})
		}
	}

	return fetchedBalances, nil
}

// Adds a balance to all ids specified in []ranges
func AddBalancesForIdRanges(balances []*types.Balance, ranges []*types.IdRange, times []*types.IdRange, balanceToAdd sdk.Uint) ([]*types.Balance, error) {
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
func SubtractBalancesForIdRanges(balances []*types.Balance, ranges []*types.IdRange, times []*types.IdRange, balanceToRemove sdk.Uint) ([]*types.Balance, error) {
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
func DeleteBalanceForIdRanges(rangesToDelete []*types.IdRange, timesToDelete []*types.IdRange, balances []*types.Balance) []*types.Balance {
	newBalances := []*types.Balance{}

	for _, balanceObj := range balances {
		currPermissionDetails := []*types.UniversalPermissionDetails{}
		for _, currRange := range balanceObj.BadgeIds {
			for _, currTime := range balanceObj.Times {
				currPermissionDetails = append(currPermissionDetails, &types.UniversalPermissionDetails{
					BadgeId: currRange,
					TimelineTime: currTime,
					TransferTime: &types.IdRange{ Start: sdk.NewUint(math.MaxUint64), End: sdk.NewUint(math.MaxUint64) }, //dummy range
				})
			}
		}

		toDeletePermissionDetails := []*types.UniversalPermissionDetails{}
		for _, rangeToDelete := range rangesToDelete {
			for _, timeToDelete := range timesToDelete {
				toDeletePermissionDetails = append(toDeletePermissionDetails, &types.UniversalPermissionDetails{
						BadgeId: rangeToDelete,
						TimelineTime: timeToDelete,
						TransferTime: &types.IdRange{ Start: sdk.NewUint(math.MaxUint64), End: sdk.NewUint(math.MaxUint64) }, //dummy range
					},
				)
			}
		}

		_, inOldButNotNew, _ := types.GetOverlapsAndNonOverlaps(currPermissionDetails, toDeletePermissionDetails)
		for _, remainingBalance := range inOldButNotNew {
			newBalances = append(newBalances, &types.Balance{
				Amount:   balanceObj.Amount,
				BadgeIds: []*types.IdRange{remainingBalance.BadgeId},
				Times: []*types.IdRange{remainingBalance.TimelineTime},
			})
		}
	}

	return newBalances
}

// Sets the balance for a specific id. Assumes balance does not exist.
func SetBalanceForIdRanges(ranges []*types.IdRange, times []*types.IdRange, amount sdk.Uint, balances []*types.Balance) ([]*types.Balance, error) {
	//TODO: ?
	// err := *new(error)
	
	// for _, balance := range balances {
	// 	if balance.Amount.Equal(amount) {
	// 		if types.IdRangeEquals(times, balance.Times) {
	// 			for _, rangeToAdd := range ranges {
	// 				balance.BadgeIds, err = types.InsertRangeToIdRanges(rangeToAdd, balance.BadgeIds)
	// 				if err != nil {
	// 					return nil, err
	// 				}
	// 			}

	// 			return balances, nil
	// 		}
	// 	}
	// }

	balances = append(balances, &types.Balance{
		Amount:   amount,
		BadgeIds: ranges,
		Times: times,
	})

	return balances, nil
}
