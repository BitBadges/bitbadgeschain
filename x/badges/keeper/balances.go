package keeper

import (
	"math"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

// Gets the balances for a specific ID.
func GetBalanceForId(id uint64, currentUserBalances []*types.Balance) (uint64, error) {
	for _, balance := range currentUserBalances {
		_, found := SearchIdRangesForId(id, balance.BadgeIds)
		if found {
			return balance.Amount, nil
		}
	}
	return 0, nil
}

// Updates the balance for a specific ids from what it currently is to newAmount.
func UpdateBalancesForIdRanges(ranges []*types.IdRange, newAmount uint64, balances []*types.Balance) ([]*types.Balance, error) {
	err := *new(error)
	ranges = SortAndMergeOverlapping(ranges)
	//Can maybe optimize this in the future by doing this all in one loop instead of deleting then setting
	balances = DeleteBalanceForIdRanges(ranges, balances)
	balances, err = SetBalanceForIdRanges(ranges, newAmount, balances)

	return balances, err
}

// Gets the balances for specified ID ranges. Returns a new []*types.Balance where only the specified ID ranges and their balances are included. Appends balance == 0 objects so all IDs are accounted for, even if not found.
func GetBalancesForIdRanges(idRanges []*types.IdRange, currentUserBalances []*types.Balance) ([]*types.Balance, error) {

	// 				//Remove everything before the start of the range. Only need to remove from idx 0 since it is sorted.
	//

	// 				//Remove everything after the end of the range. Only need to remove from last idx since it is sorted.
	//

	err := *new(error)
	fetchedBalances := []*types.Balance{}
	idRanges = SortAndMergeOverlapping(idRanges)
	idsWithZeroBalance := idRanges
	for _, userBalanceObj := range currentUserBalances {
		//For each specified range, search the current userBalanceObj's IdRanges to see if there is any overlap.
		//If so, we add the overlapping range and current balance to the new []*types.Balances to be returned.

		for _, idRange := range idRanges {
			idxSpan, found := GetIdxSpanForRange(idRange, userBalanceObj.BadgeIds)
			if found {
				//Set newIdRanges to all ranges where there is overlap
				newIdRanges := userBalanceObj.BadgeIds[idxSpan.Start : idxSpan.End+1] //+ 1 since the slice is non-inclusive

				//Since GetIdxSpanForRange only returns the indexes of the overlapping ranges,
				//we need to remove the non-overlapping portions of the first and last ranges.

				//Remove everything before the start of the range. Only need to remove from idx 0 since it is sorted.
				if idRange.Start > 0 && len(newIdRanges) > 0 {
					everythingBefore := &types.IdRange{
						Start: 0,
						End:   idRange.Start - 1,
					}
					removedRanges, _ := RemoveIdsFromIdRange(everythingBefore, newIdRanges[0])
					newIdRanges = append([]*types.IdRange{}, removedRanges...)
					newIdRanges = append(newIdRanges, newIdRanges[1:]...)
				}

				//Remove everything after the end of the range. Only need to remove from last idx since it is sorted.
				if idRange.End < math.MaxUint64 && len(newIdRanges) > 0 {
					everythingAfter := &types.IdRange{
						Start: idRange.End + 1,
						End:   math.MaxUint64,
					}

					removedRanges, _ := RemoveIdsFromIdRange(everythingAfter, newIdRanges[len(newIdRanges)-1])
					newIdRanges = append([]*types.IdRange{}, newIdRanges[0:len(newIdRanges)-1]...)
					newIdRanges = append(newIdRanges, removedRanges...)
				}

				//If we found any overlapping ranges, remove the IDs from the list of IDs with zero balance
				for _, foundIdRange := range newIdRanges {
					newIdsWithZeroBalance := []*types.IdRange{}
					for _, currRangeWithZeroBalance := range idsWithZeroBalance {
						removedRanges, _ := RemoveIdsFromIdRange(foundIdRange, currRangeWithZeroBalance)
						newIdsWithZeroBalance = append(newIdsWithZeroBalance, removedRanges...)
					}
					idsWithZeroBalance = newIdsWithZeroBalance
				}

				//Update the fetchedBalances with the IDs which we found
				fetchedBalances, err = UpdateBalancesForIdRanges(newIdRanges, userBalanceObj.Amount, fetchedBalances)
				if err != nil {
					return fetchedBalances, err
				}
			}
		}
	}

	//Update balance objects with IDs where balance == 0
	if len(idsWithZeroBalance) > 0 {
		fetchedBalances = append([]*types.Balance{
			{
				Amount:  0,
				BadgeIds: idsWithZeroBalance,
			},
		}, fetchedBalances...)
	}

	return fetchedBalances, nil
}

// Adds a balance to all ids specified in []ranges
func AddBalancesForIdRanges(userBalance types.UserBalanceStore, ranges []*types.IdRange, balanceToAdd uint64) (types.UserBalanceStore, error) {
	currBalances, err := GetBalancesForIdRanges(ranges, userBalance.Balances)
	if err != nil {
		return userBalance, err
	}

	for _, balance := range currBalances {
		newBalance, err := SafeAdd(balance.Amount, balanceToAdd)
		if err != nil {
			return userBalance, err
		}

		userBalance.Balances, err = UpdateBalancesForIdRanges(balance.BadgeIds, newBalance, userBalance.Balances)
		if err != nil {
			return userBalance, err
		}
	}
	return userBalance, nil
}

// Subtracts a balance to all ids specified in []ranges
func SubtractBalancesForIdRanges(UserBalance types.UserBalanceStore, ranges []*types.IdRange, balanceToRemove uint64) (types.UserBalanceStore, error) {
	currBalances, err := GetBalancesForIdRanges(ranges, UserBalance.Balances)
	if err != nil {
		return UserBalance, err
	}

	for _, currBalanceObj := range currBalances {
		newBalance, err := SafeSubtract(currBalanceObj.Amount, balanceToRemove)
		if err != nil {
			return UserBalance, err
		}

		UserBalance.Balances, err = UpdateBalancesForIdRanges(currBalanceObj.BadgeIds, newBalance, UserBalance.Balances)
		if err != nil {
			return UserBalance, err
		}
	}
	return UserBalance, nil
}

// Deletes the balance for a specific id.
func DeleteBalanceForIdRanges(ranges []*types.IdRange, balances []*types.Balance) []*types.Balance {
	newBalances := []*types.Balance{}
	for _, balanceObj := range balances {
		for _, rangeToDelete := range ranges {
			currRanges := balanceObj.BadgeIds
			idxSpan, found := GetIdxSpanForRange(rangeToDelete, currRanges)
			if found {
				if idxSpan.End == 0 {
					idxSpan.End = idxSpan.Start
				}

				//Remove the ids within the rangeToDelete from existing ranges
				newIdRanges := append([]*types.IdRange{}, currRanges[:idxSpan.Start]...)
				for i := idxSpan.Start; i <= idxSpan.End; i++ {
					removedRanges, _ := RemoveIdsFromIdRange(rangeToDelete, currRanges[i])
					newIdRanges = append(newIdRanges, removedRanges...)
				}
				newIdRanges = append(newIdRanges, currRanges[idxSpan.End+1:]...)
				balanceObj.BadgeIds = newIdRanges
			}
		}

		// If we don't have any corresponding IDs, don't store this anymore
		if len(balanceObj.BadgeIds) > 0 {
			newBalances = append(newBalances, balanceObj)
		}
	}

	return newBalances
}

// Sets the balance for a specific id. Assumes balance does not exist.
func SetBalanceForIdRanges(ranges []*types.IdRange, amount uint64, balances []*types.Balance) ([]*types.Balance, error) {
	if amount == 0 {
		return balances, nil
	}
	err := *new(error)

	idx, found := SearchBalances(amount, balances)
	newBalances := []*types.Balance{}
	if !found {
		//We don't have an existing object with such a balance
		newBalances = append(newBalances, balances[:idx]...)
		rangesToInsert := []*types.IdRange{}
		for _, rangeToAdd := range ranges {
			rangesToInsert = append(rangesToInsert, CreateIdRange(rangeToAdd.Start, rangeToAdd.End))
		}
		newBalances = append(newBalances, &types.Balance{
			Amount:  amount,
			BadgeIds: rangesToInsert,
		})
		newBalances = append(newBalances, balances[idx:]...)
	} else {
		//Update existing balance object
		newBalances = balances
		for _, rangeToAdd := range ranges {
			newBalances[idx].BadgeIds, err = InsertRangeToIdRanges(rangeToAdd, newBalances[idx].BadgeIds)
			if err != nil {
				return nil, err
			}
		}
	}

	return newBalances, nil
}

// Balances will be sorted, so we can binary search to get the targetIdx.
// If found, returns (the index it was found at, true). Else, returns (index to insert at, false).
func SearchBalances(targetAmount uint64, balances []*types.Balance) (int, bool) {
	balanceLow := 0
	balanceHigh := len(balances) - 1
	median := 0
	hasEntryWithSameBalance := false
	idx := 0
	for balanceLow <= balanceHigh {
		median = int(uint(balanceLow+balanceHigh) >> 1)
		if balances[median].Amount == targetAmount {
			hasEntryWithSameBalance = true
			break
		} else if balances[median].Amount > targetAmount {
			balanceHigh = median - 1
		} else {
			balanceLow = median + 1
		}
	}

	if len(balances) != 0 {
		idx = median + 1
		if targetAmount <= balances[median].Amount {
			idx = median
		}
	}

	return idx, hasEntryWithSameBalance
}
