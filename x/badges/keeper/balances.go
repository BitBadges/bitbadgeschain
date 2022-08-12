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
		return 0, ErrOverflow
	}
	return left - right, nil
}

// Updates the balance for a specific id from what it currently is to newAmount.
func UpdateBalancesForIdRanges(ranges []*types.IdRange, newAmount uint64, balanceObjects []*types.BalanceObject) []*types.BalanceObject {
	ranges = SortIdRangesAndMergeIfNecessary(ranges)
	balanceObjects = DeleteBalanceForIdRanges(ranges, balanceObjects)
	balanceObjects = SetBalanceForIdRanges(ranges, newAmount, balanceObjects)
	return balanceObjects
}

// Gets a balance for a specific id
func GetBalancesForIdRanges(ranges []*types.IdRange, balanceObjects []*types.BalanceObject) []*types.BalanceObject {
	rangeBalanceObjects := []*types.BalanceObject{}

	for _, balanceObj := range balanceObjects {
		balanceObj.IdRanges = GetIdRangesWithOmitEmptyCaseHandled(balanceObj.IdRanges)
		
		
		for _, rangeToGet := range ranges {
			idxRange, found := GetIdxSpanForRange(*rangeToGet, balanceObj.IdRanges)
			if found {
				everythingBefore := types.IdRange{
					Start: 0, 
					End: rangeToGet.Start - 1, 
				}

				everythingAfter := types.IdRange{
					Start: rangeToGet.End + 1, 
					End: math.MaxUint64,
				}
				newIdRanges := []*types.IdRange{}
				for i := idxRange.Start; i <= idxRange.End; i++ {
					currRange := balanceObj.IdRanges[i]
					if rangeToGet.Start > 0 {
						newIdRanges = append(newIdRanges, RemoveIdsFromIdRange(everythingBefore, *currRange)...)
					}

					
					if rangeToGet.End < math.MaxUint64 {
						if len(newIdRanges) == 0 {
							newIdRanges = append(newIdRanges, RemoveIdsFromIdRange(everythingAfter, *currRange)...)
						} else {
							finalIdRanges := []*types.IdRange{}
							for _, newIdRange := range newIdRanges {
								finalIdRanges = append(finalIdRanges, RemoveIdsFromIdRange(everythingAfter, *newIdRange)...)
							}
							newIdRanges = finalIdRanges
						}
					}
				}

				rangeBalanceObjects = UpdateBalancesForIdRanges(newIdRanges, balanceObj.Balance, rangeBalanceObjects)
				// newIdRanges = append(newIdRanges, RemoveIdsFromIdRange(*rangeToGet, *balanceObj.IdRanges[idxRange.End])...)
			}
		}
	}

	if len(rangeBalanceObjects) == 0 {
		rangeBalanceObjects = append(rangeBalanceObjects, &types.BalanceObject{
			Balance:  0,
			IdRanges: ranges,
		})
	}

	return rangeBalanceObjects
}

// Adds a balance for the id
func AddBalancesForIdRanges(ctx sdk.Context, userBalanceInfo types.UserBalanceInfo, ranges []*types.IdRange, balanceToAdd uint64) (types.UserBalanceInfo, error) {
	currBalances := GetBalancesForIdRanges(ranges, userBalanceInfo.BalanceAmounts)
	for _, balanceObj := range currBalances {
		newBalance, err := SafeAdd(balanceObj.Balance, balanceToAdd)
		if err != nil {
			return userBalanceInfo, err
		}

		userBalanceInfo.BalanceAmounts = UpdateBalancesForIdRanges(balanceObj.IdRanges, newBalance, userBalanceInfo.BalanceAmounts)
	}
	return userBalanceInfo, nil
}

// Subtracts a balance for the id
func SubtractBalancesForIdRanges(ctx sdk.Context, userBalanceInfo types.UserBalanceInfo, ranges []*types.IdRange, balanceToRemove uint64) (types.UserBalanceInfo, error) {
	currBalances := GetBalancesForIdRanges(ranges, userBalanceInfo.BalanceAmounts)
	
	for _, balanceObj := range currBalances {
		newBalance, err := SafeSubtract(balanceObj.Balance, balanceToRemove)
		if err != nil {
			return userBalanceInfo, err
		}

		userBalanceInfo.BalanceAmounts = UpdateBalancesForIdRanges(balanceObj.IdRanges, newBalance, userBalanceInfo.BalanceAmounts)
	}
	return userBalanceInfo, nil
}

// Deletes the balance for a specific id.
func DeleteBalanceForIdRanges(ranges []*types.IdRange, balanceObjects []*types.BalanceObject) []*types.BalanceObject {
	newBalanceObjects := []*types.BalanceObject{}
	for _, balanceObj := range balanceObjects {
		balanceObj.IdRanges = GetIdRangesWithOmitEmptyCaseHandled(balanceObj.IdRanges)

		for _, rangeToDelete := range ranges {
			idxRange, found := GetIdxSpanForRange(*rangeToDelete, balanceObj.IdRanges)
			if found {
				if idxRange.End == 0 {
					idxRange.End = idxRange.Start
				}

				newIdRanges := append([]*types.IdRange{}, balanceObj.IdRanges[:idxRange.Start]...)
				for i := idxRange.Start; i <= idxRange.End; i++ {
					newIdRanges = append(newIdRanges, RemoveIdsFromIdRange(*rangeToDelete, *balanceObj.IdRanges[i])...)
				}
				newIdRanges = append(newIdRanges, balanceObj.IdRanges[idxRange.End + 1:]...)
				balanceObj.IdRanges = newIdRanges
			}
		}

		if len(balanceObj.IdRanges) > 0 {
			newBalanceObjects = append(newBalanceObjects, balanceObj)
		}
	}
	return newBalanceObjects
}

// Sets the balance for a specific id.
func SetBalanceForIdRanges(ranges []*types.IdRange, amount uint64, balanceObjects []*types.BalanceObject) []*types.BalanceObject {
	if amount == 0 {
		return balanceObjects
	}

	ranges = SortIdRangesAndMergeIfNecessary(ranges)
	
	idx, found := SearchBalanceObjectsForBalanceAndGetIdxToInsertIfNotFound(amount, balanceObjects)
	newBalanceObjects := []*types.BalanceObject{}
	if !found {
		newBalanceObjects = append(newBalanceObjects, balanceObjects[:idx]...)
		newBalanceObjects = append(newBalanceObjects, &types.BalanceObject{
			Balance:  amount,
			IdRanges: ranges,
		})
		newBalanceObjects = append(newBalanceObjects, balanceObjects[idx:]...)
	} else {
		newBalanceObjects = balanceObjects
		newBalanceObjects[idx].IdRanges = GetIdRangesWithOmitEmptyCaseHandled(newBalanceObjects[idx].IdRanges)
		for _, rangeToAdd := range ranges {
			newBalanceObjects[idx].IdRanges = InsertRangeToIdRanges(*rangeToAdd, newBalanceObjects[idx].IdRanges)
		}
	}

	return newBalanceObjects
}

// Balances will be sorted, so we can binary search to get the targetIdx. Returns the index to insert at if not found
func SearchBalanceObjectsForBalanceAndGetIdxToInsertIfNotFound(targetAmount uint64, balanceObjects []*types.BalanceObject) (int, bool) {
	balanceLow := 0
	balanceHigh := len(balanceObjects) - 1
	median := 0
	hasEntryWithSameBalance := false
	setIdx := 0
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
		setIdx = median + 1
		if targetAmount <= balanceObjects[median].Balance {
			setIdx = median
		}
	}

	return setIdx, hasEntryWithSameBalance
}
