package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Gets the balances for a specific ID. Assumes balances are sorted
func GetBalanceForId(id sdk.Uint, balances []*types.Balance) sdk.Uint {
	for _, balance := range balances {
		_, found := SearchIdRangesForId(id, balance.BadgeIds)
		if found {
			return balance.Amount
		}
	}
	return sdk.NewUint(0)
}

// Updates the balance for a specific ids from what it currently is to newAmount.
func UpdateBalancesForIdRanges(ranges []*types.IdRange, newAmount sdk.Uint, balances []*types.Balance) (newBalances []*types.Balance, err error) {
	//Can maybe optimize this in the future by doing this all in one loop instead of deleting then setting
	ranges = SortAndMergeOverlapping(ranges)
	balances = DeleteBalanceForIdRanges(ranges, balances)
	balances, err = SetBalanceForIdRanges(ranges, newAmount, balances)

	return balances, err
}

// Gets the balances for specified ID ranges. Returns a new []*types.Balance where only the specified ID ranges and their balances are included. Appends balance == 0 objects so all IDs are accounted for, even if not found.
func GetBalancesForIdRanges(idRanges []*types.IdRange, balances []*types.Balance) (newBalances []*types.Balance, err error) {
	fetchedBalances := []*types.Balance{}
	idRanges = SortAndMergeOverlapping(idRanges)
	idsWithZeroBalance := idRanges //We use this to keep track of which IDs we have already found a balance for. Will deduct from this as we find balances

	for _, balanceObj := range balances {
		//For each specified range, search the current balanceObj's IdRanges to see if there is any overlap.
		//If so, we add the overlapping range and current balance to the new []*types.Balances to be returned.

		for _, idRange := range idRanges {
			idxSpan, found := GetIdxSpanForRange(idRange, balanceObj.BadgeIds)
			if found {
				//Set newIdRanges to all ranges where there is overlap
				newIdRanges := balanceObj.BadgeIds[idxSpan.Start.Uint64():idxSpan.End.AddUint64(1).Uint64()] // + 1 since the slice is non-inclusive but idxSpan is

				//Since GetIdxSpanForRange only returns the indexes of the overlapping ranges,
				//we need to remove the non-overlapping portions of the first and last ranges.

				//Remove everything before the start of the range. Only need to remove from idx 0 since it is sorted.
				if !idRange.Start.IsZero() && len(newIdRanges) > 0 {
					everythingBefore := &types.IdRange{
						Start: sdk.NewUint(0),
						End:   idRange.Start.SubUint64(1),
					}
					removedRanges, _ := RemoveIdsFromIdRange(everythingBefore, newIdRanges[0])
					newIdRanges = append(removedRanges, newIdRanges[1:]...)
				}

				//Remove everything after the end of the range. Only need to remove from last idx since it is sorted.
				if len(newIdRanges) > 0 {
					rangeToTrim := newIdRanges[len(newIdRanges)-1]
					if idRange.End.LT(rangeToTrim.End) {

						everythingAfter := &types.IdRange{
							Start: idRange.End.AddUint64(1),
							End:   rangeToTrim.End,
						}

						removedRanges, _ := RemoveIdsFromIdRange(everythingAfter, rangeToTrim)
						newIdRanges = append(newIdRanges[0:len(newIdRanges)-1], removedRanges...)
					}
				}

				//If we found any overlapping ranges, remove the IDs from the list of IDs with zero balance
				for _, idRange := range newIdRanges {
					newIdsWithZeroBalance := []*types.IdRange{}
					for _, currRangeWithZeroBalance := range idsWithZeroBalance {
						removedRanges, _ := RemoveIdsFromIdRange(idRange, currRangeWithZeroBalance)
						newIdsWithZeroBalance = append(newIdsWithZeroBalance, removedRanges...)
					}
					idsWithZeroBalance = newIdsWithZeroBalance
				}

				//Update the fetchedBalances with the IDs which we found
				fetchedBalances, err = UpdateBalancesForIdRanges(newIdRanges, balanceObj.Amount, fetchedBalances)
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
				Amount:   sdk.NewUint(0),
				BadgeIds: idsWithZeroBalance,
			},
		}, fetchedBalances...)
	}

	return fetchedBalances, nil
}

// Adds a balance to all ids specified in []ranges
func AddBalancesForIdRanges(balances []*types.Balance, ranges []*types.IdRange, balanceToAdd sdk.Uint) ([]*types.Balance, error) {
	currBalances, err := GetBalancesForIdRanges(ranges, balances)
	if err != nil {
		return balances, err
	}

	for _, balance := range currBalances {
		newBalanceAmount, err := SafeAdd(balance.Amount, balanceToAdd)
		if err != nil {
			return balances, err
		}

		balances, err = UpdateBalancesForIdRanges(balance.BadgeIds, newBalanceAmount, balances)
		if err != nil {
			return balances, err
		}
	}
	return balances, nil
}

// Subtracts a balance to all ids specified in []ranges
func SubtractBalancesForIdRanges(balances []*types.Balance, ranges []*types.IdRange, balanceToRemove sdk.Uint) ([]*types.Balance, error) {
	currBalances, err := GetBalancesForIdRanges(ranges, balances)
	if err != nil {
		return balances, err
	}

	for _, currBalanceObj := range currBalances {
		newBalance, err := SafeSubtract(currBalanceObj.Amount, balanceToRemove)
		if err != nil {
			return balances, err
		}

		balances, err = UpdateBalancesForIdRanges(currBalanceObj.BadgeIds, newBalance, balances)
		if err != nil {
			return balances, err
		}
	}
	return balances, nil
}

// Deletes the balance for a specific id.
func DeleteBalanceForIdRanges(rangesToDelete []*types.IdRange, balances []*types.Balance) []*types.Balance {
	newBalances := []*types.Balance{}
	for _, balanceObj := range balances {
		for _, rangeToDelete := range rangesToDelete {
			currRanges := balanceObj.BadgeIds
			//Remove the ids within the rangeToDelete from existing ranges
			idxSpan, found := GetIdxSpanForRange(rangeToDelete, currRanges)
			if found {

				newIdRanges := append([]*types.IdRange{}, currRanges[:idxSpan.Start.Uint64()]...)

				//For the overlapping ranges, remove the ids within the rangeToDelete
				for i := idxSpan.Start; i.LTE(idxSpan.End); i = i.Incr() {
					removedRanges, _ := RemoveIdsFromIdRange(rangeToDelete, currRanges[i.Uint64()])
					newIdRanges = append(newIdRanges, removedRanges...)
				}

				newIdRanges = append(newIdRanges, currRanges[idxSpan.End.Incr().Uint64():]...)

				balanceObj.BadgeIds = newIdRanges
			}
		}

		// If we don't have any corresponding badge IDs, remove balance from balance array
		if len(balanceObj.BadgeIds) > 0 {
			newBalances = append(newBalances, balanceObj)
		}
	}

	return newBalances
}

// Sets the balance for a specific id. Assumes balance does not exist.
func SetBalanceForIdRanges(ranges []*types.IdRange, amount sdk.Uint, balances []*types.Balance) ([]*types.Balance, error) {
	if amount.IsZero() {
		return balances, nil
	}

	err := *new(error)
	idx, found := SearchBalances(amount, balances)

	newBalances := []*types.Balance{}
	if !found {
		//We don't have an existing object with such a balance, so we add it at idx
		newBalances = append(newBalances, balances[:idx]...)
		newBalances = append(newBalances, &types.Balance{
			Amount:   amount,
			BadgeIds: ranges,
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
func SearchBalances(targetAmount sdk.Uint, balances []*types.Balance) (int, bool) {
	balanceLow := 0
	balanceHigh := len(balances) - 1
	median := 0
	hasEntryWithSameBalance := false
	idx := 0
	for balanceLow <= balanceHigh {
		median = int(uint(balanceLow+balanceHigh) >> 1)
		if balances[median].Amount.Equal(targetAmount) {
			hasEntryWithSameBalance = true
			break
		} else if balances[median].Amount.GT(targetAmount) {
			balanceHigh = median - 1
		} else {
			balanceLow = median + 1
		}
	}

	if len(balances) != 0 {
		idx = median + 1
		if targetAmount.LTE(balances[median].Amount) {
			idx = median
		}
	}

	return idx, hasEntryWithSameBalance
}
