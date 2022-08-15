package keeper

import (
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

// Safe adds two uint64s and returns an error if the result overflows uint64.
func SafeAdd(left uint64, right uint64) (uint64, error) {
	sum := left + right
	if sum < left {
		return 0, ErrOverflow
	}
	return sum, nil
}

// Safe subtracts two uint64s and returns an error if the result underflows uint64.
func SafeSubtract(left uint64, right uint64) (uint64, error) {
	if right > left {
		return 0, ErrUnderflow
	}
	return left - right, nil
}

// Updates the balance for a specific id from what it currently is to newAmount.
func UpdateBalancesForIdRanges(ranges []*types.IdRange, newAmount uint64, balanceObjects []*types.BalanceObject) []*types.BalanceObject {
	ranges = SortIdRangesAndMergeIfNecessary(ranges)
	//Can maybe optimize this in the future by doing this all in one loop instead of deleting then setting
	balanceObjects = DeleteBalanceForIdRanges(ranges, balanceObjects)
	balanceObjects = SetBalanceForIdRanges(ranges, newAmount, balanceObjects)
	return balanceObjects
}

// Gets the balances for specified ID ranges. Returns a new []*types.BalanceObject where only the specified ID ranges and their balances are included. Appends balance == 0 objects so all IDs are accounted for, even if not found. 
func GetBalancesForIdRanges(idRanges []*types.IdRange, currentUserBalanceObjects []*types.BalanceObject) []*types.BalanceObject {
	balanceObjectsForSpecifiedRanges := []*types.BalanceObject{}
	idRangesNotFound := idRanges
	for _, userBalanceObj := range currentUserBalanceObjects {
		userBalanceObj.IdRanges = GetIdRangesWithOmitEmptyCaseHandled(userBalanceObj.IdRanges)

		//For each specified range, search the current userBalanceObj's IdRanges to see if there is any overlap. 
		//If so, we add the overlapping range and current balance to the new []*types.BalanceObjects to be returned.
		
		for _, idRange := range idRanges {
			idxSpan, found := GetIdxSpanForRange(idRange, userBalanceObj.IdRanges)
			if found {
				if idxSpan.End == 0 {
					idxSpan.End = idxSpan.Start
				}

				//Set newIdRanges to the ranges where there is overlap
				newIdRanges := userBalanceObj.IdRanges[idxSpan.Start:idxSpan.End + 1]

				//Remove everything before the start of the range. Only need to remove from idx 0 since it is sorted.
				if idRange.Start > 0 && len(newIdRanges) > 0 {
					everythingBefore := &types.IdRange{
						Start: 0,
						End: idRange.Start - 1, 
					}
					newIdRanges = append(newIdRanges, RemoveIdsFromIdRange(everythingBefore, newIdRanges[0])...)
				}
				
				//Remove everything after the end of the range. Only need to remove from last idx since it is sorted.
				if idRange.End < math.MaxUint64 && len(newIdRanges) > 0 {
					everythingAfter := &types.IdRange{
						Start: idRange.End + 1, 
						End: math.MaxUint64,
					}
					newIdRanges = append(newIdRanges, RemoveIdsFromIdRange(everythingAfter, newIdRanges[len(newIdRanges)-1])...)
				}

				for _, newIdRange := range newIdRanges {
					newNotFoundRanges := []*types.IdRange{}
					for _, idRangeNotFound := range idRangesNotFound {
						newNotFoundRanges = append(newNotFoundRanges, RemoveIdsFromIdRange(newIdRange, idRangeNotFound)...)
					}
					idRangesNotFound = newNotFoundRanges
				}
				balanceObjectsForSpecifiedRanges = UpdateBalancesForIdRanges(newIdRanges, userBalanceObj.Balance, balanceObjectsForSpecifiedRanges)
			}
		}
	}

	//Update balance objects with IDs where balance == 0
	if len(idRangesNotFound) > 0 {
		updatedBalanceObjects := []*types.BalanceObject{}
		updatedBalanceObjects = append(updatedBalanceObjects, &types.BalanceObject{
			Balance: 0,
			IdRanges: idRangesNotFound,
		})
		updatedBalanceObjects = append(updatedBalanceObjects, balanceObjectsForSpecifiedRanges...)
		return updatedBalanceObjects
	} else {
		return balanceObjectsForSpecifiedRanges
	}
}

// Adds a balance to all ids specified in []ranges
func AddBalancesForIdRanges(ctx sdk.Context, userBalanceInfo types.UserBalanceInfo, ranges []*types.IdRange, balanceToAdd uint64) (types.UserBalanceInfo, error) {
	currBalances := GetBalancesForIdRanges(ranges, userBalanceInfo.BalanceAmounts)
	for _, currBalanceObj := range currBalances {
		newBalance, err := SafeAdd(currBalanceObj.Balance, balanceToAdd)
		if err != nil {
			return userBalanceInfo, err
		}

		userBalanceInfo.BalanceAmounts = UpdateBalancesForIdRanges(currBalanceObj.IdRanges, newBalance, userBalanceInfo.BalanceAmounts)
	}
	return userBalanceInfo, nil
}

// Subtracts a balance to all ids specified in []ranges
func SubtractBalancesForIdRanges(ctx sdk.Context, userBalanceInfo types.UserBalanceInfo, ranges []*types.IdRange, balanceToRemove uint64) (types.UserBalanceInfo, error) {
	currBalances := GetBalancesForIdRanges(ranges, userBalanceInfo.BalanceAmounts)
	for _, currBalanceObj := range currBalances {
		newBalance, err := SafeSubtract(currBalanceObj.Balance, balanceToRemove)
		if err != nil {
			return userBalanceInfo, err
		}

		userBalanceInfo.BalanceAmounts = UpdateBalancesForIdRanges(currBalanceObj.IdRanges, newBalance, userBalanceInfo.BalanceAmounts)
	}
	return userBalanceInfo, nil
}

// Deletes the balance for a specific id.
func DeleteBalanceForIdRanges(ranges []*types.IdRange, balanceObjects []*types.BalanceObject) []*types.BalanceObject {
	newBalanceObjects := []*types.BalanceObject{}
	for _, balanceObj := range balanceObjects {
		balanceObj.IdRanges = GetIdRangesWithOmitEmptyCaseHandled(balanceObj.IdRanges)

		for _, rangeToDelete := range ranges {
			idxSpan, found := GetIdxSpanForRange(rangeToDelete, balanceObj.IdRanges)
			if found {
				if idxSpan.End == 0 {
					idxSpan.End = idxSpan.Start
				}

				//Remove the ids within the rangeToDelete from existing ranges
				newIdRanges := append([]*types.IdRange{}, balanceObj.IdRanges[:idxSpan.Start]...)
				for i := idxSpan.Start; i <= idxSpan.End; i++ {
					newIdRanges = append(newIdRanges, RemoveIdsFromIdRange(rangeToDelete, balanceObj.IdRanges[i])...)
				}
				newIdRanges = append(newIdRanges, balanceObj.IdRanges[idxSpan.End + 1:]...)
				balanceObj.IdRanges = newIdRanges
			}
		}

		// If we don't have any corresponding IDs, don't store this anymore
		if len(balanceObj.IdRanges) > 0 {
			newBalanceObjects = append(newBalanceObjects, balanceObj)
		}
	}
	return newBalanceObjects
}

// Sets the balance for a specific id. Assumes balance does not exist.
func SetBalanceForIdRanges(ranges []*types.IdRange, amount uint64, balanceObjects []*types.BalanceObject) []*types.BalanceObject {
	if amount == 0 {
		return balanceObjects
	}

	ranges = SortIdRangesAndMergeIfNecessary(ranges)
	
	idx, found := SearchBalanceObjectsForBalanceAndGetIdxToInsertIfNotFound(amount, balanceObjects)
	newBalanceObjects := []*types.BalanceObject{}
	if !found {
		//We don't have an existing object with such a balance
		newBalanceObjects = append(newBalanceObjects, balanceObjects[:idx]...)
		newBalanceObjects = append(newBalanceObjects, &types.BalanceObject{
			Balance:  amount,
			IdRanges: ranges,
		})
		newBalanceObjects = append(newBalanceObjects, balanceObjects[idx:]...)
	} else {
		//Update existing balance object
		newBalanceObjects = balanceObjects
		newBalanceObjects[idx].IdRanges = GetIdRangesWithOmitEmptyCaseHandled(newBalanceObjects[idx].IdRanges)
		for _, rangeToAdd := range ranges {
			newBalanceObjects[idx].IdRanges = InsertRangeToIdRanges(rangeToAdd, newBalanceObjects[idx].IdRanges)
		}
	}

	return newBalanceObjects
}

// Balances will be sorted, so we can binary search to get the targetIdx.
// If found, returns (the index it was found at, true). Else, returns (index to insert at, false).
func SearchBalanceObjectsForBalanceAndGetIdxToInsertIfNotFound(targetAmount uint64, balanceObjects []*types.BalanceObject) (int, bool) {
	balanceLow := 0
	balanceHigh := len(balanceObjects) - 1
	median := 0
	hasEntryWithSameBalance := false
	idx := 0
	for balanceLow <= balanceHigh {
		median = int(uint(balanceLow+balanceHigh) >> 1)
		if balanceObjects[median].Balance == targetAmount {
			hasEntryWithSameBalance = true
			break
		} else if balanceObjects[median].Balance > targetAmount {
			balanceHigh = median - 1
		} else {
			balanceLow = median + 1
		}
	}

	if len(balanceObjects) != 0 {
		idx = median + 1
		if targetAmount <= balanceObjects[median].Balance {
			idx = median
		}
	}

	return idx, hasEntryWithSameBalance
}
