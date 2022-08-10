package keeper

import (
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

// Balances take the form of BalanceToIds[] in which we have a map of each balance to its corresuponding SubbadgeRange[]

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

// Updates the balance for a specific subbadgeId from what it currently is to newAmount. 
func UpdateBalanceForSubbadgeId(subbadgeId uint64, newAmount uint64, amounts []*types.BalanceToIds) []*types.BalanceToIds {
	amounts = RemoveBalanceForSubbadgeId(subbadgeId, amounts)
	amounts = SetBalanceForSubbadgeId(subbadgeId, newAmount, amounts)
	return amounts
}

// Gets a balance for a specific subbadgeId.
func GetBalanceForSubbadgeId(subbadgeId uint64, balancesArr []*types.BalanceToIds) uint64 {
	for _, balanceObj := range balancesArr {
		currBalance := balanceObj.Balance
		currIds := balanceObj.Ids

		// Handle the case where it omits an empty NumberRange because Start && End == 0
		if len(currIds) == 0 {
			currIds = append(currIds, &types.NumberRange{})
		}

		// SubbadgeRanges should be sorted by start IDs (and endIds but endIds can == 0), so we can binary search.
		low := 0
		high := len(currIds) - 1

		for low <= high {
			median := int(uint(low + high) >> 1)
			currRange := currIds[median]
			if currRange.End == 0 {
				currRange.End = currRange.Start // If end == 0, set it to start by convention (done in order to save space)
			}

			if currIds[median].Start <= subbadgeId && currIds[median].End >= subbadgeId {
				return currBalance
			} else if currIds[median].Start > subbadgeId {
				high = median - 1
			} else {
				low = median + 1
			}
		}
	
		if low == len(currIds) || !(currIds[low].Start <= subbadgeId && currIds[low].End >= subbadgeId) {
			continue
		} else {
			return currBalance
		}
	}

	return 0 //Not found; return 0
}

// Remove the balance for a specific subbadgeId.
func RemoveBalanceForSubbadgeId(subbadgeId uint64, balanceObjs []*types.BalanceToIds) []*types.BalanceToIds {
	newAmounts := []*types.BalanceToIds{}
	for _, balanceObj := range balanceObjs {
		// Handle the case where it omits an empty NumberRange because Start && End == 0
		if len(balanceObj.Ids) == 0 {
			balanceObj.Ids = append(balanceObj.Ids, &types.NumberRange{})
		}

		// SubbadgeRanges should be sorted by start IDs (and endIds but endIds can == 0 with our added convention), so we can binary search.
		low := 0
		high := len(balanceObj.Ids) - 1
		

		for low <= high {
			median := int(uint(low + high) >> 1)
			currRange := balanceObj.Ids[median]
			if currRange.End == 0 {
				currRange.End = currRange.Start // If end == 0, set it to start by convention (done in order to save space)
			}

			if balanceObj.Ids[median].Start <= subbadgeId && balanceObj.Ids[median].End >= subbadgeId {
				newIds := balanceObj.Ids[:median]

				//If we still have an existing before range, keep that up until subbadge - 1
				if subbadgeId >= 1 && subbadgeId - 1 >= currRange.Start {
					//Don't store endIdx if end == start
					endIdx := subbadgeId - 1
					if endIdx == currRange.Start {
						endIdx = 0
					}

					newIds = append(newIds, &types.NumberRange{
						Start: currRange.Start,
						End:   endIdx,
					})
				}

				//If we still have an existing after range, start that at subbadge + 1
				if subbadgeId <= math.MaxUint64-1 && subbadgeId + 1 <= currRange.End {
					//Don't store endIdx if end == start
					endIdx := currRange.End
					if endIdx == subbadgeId + 1 {
						endIdx = 0
					}

					newIds = append(newIds, &types.NumberRange{
						Start: subbadgeId + 1,
						End:   endIdx,
					})
				}

				newIds = append(newIds, balanceObj.Ids[median+1:]...) 
				balanceObj.Ids = newIds
				break
			} else if balanceObj.Ids[median].Start > subbadgeId {
				high = median - 1
			} else {
				low = median + 1
			}
		}
		
		if len(balanceObj.Ids) > 0 {
			newAmounts = append(newAmounts, balanceObj)
		}
	}
	return newAmounts
}

// Sets the balance for a specific subbadgeId. Assumes RemoveBalanceForSubbadgeId() was previously called.
func SetBalanceForSubbadgeId(subbadgeId uint64, amount uint64, amounts []*types.BalanceToIds) []*types.BalanceToIds {
	// Don't store if balance == 0
	if amount == 0 {
		return amounts
	}

	// Balances will be sorted, so we can binary search to get the targetIdx
	balanceLow := 0
	balanceHigh := len(amounts) - 1
	median := 0
	hasEntryWithSameBalance := false
	setIdx := 0
	for balanceLow <= balanceHigh {
		median = int(uint(balanceLow + balanceHigh) >> 1)
		if amounts[median].Balance == amount {
			hasEntryWithSameBalance = true
			break;
		} else if amounts[median].Balance > amount {
			balanceHigh = median - 1
		} else {
			balanceLow = median + 1
		}
	}
	
	if len(amounts) != 0 {
		setIdx = median + 1
		if (amount <= amounts[median].Balance) {
			setIdx = median
		}
	}


	newAmounts := []*types.BalanceToIds{}
	if !hasEntryWithSameBalance {
		newAmounts = append(newAmounts, amounts[:setIdx]...)
		newAmounts = append(newAmounts, &types.BalanceToIds{
			Balance: amount,
			Ids: []*types.NumberRange{{Start: subbadgeId}},
		})
		newAmounts = append(newAmounts, amounts[setIdx:]...)
	} else {
		newAmounts = amounts

		//Handle case where we have it omits the empty case where Start && End == 0
		if len(newAmounts[setIdx].Ids) == 0 {
			newAmounts[setIdx].Ids = append(newAmounts[setIdx].Ids, &types.NumberRange{})
		}

		newIds := []*types.NumberRange{}

		insertIdx := 0
		if newAmounts[setIdx].Ids[0].Start > subbadgeId {
			//Handle case where subbadgeId is less than the starting subbadgeId (append to front)
			insertIdx = 0
			newIds = append(newIds, &types.NumberRange{Start: subbadgeId})
			newIds = append(newIds, newAmounts[setIdx].Ids...)
		} else if newAmounts[setIdx].Ids[len(newAmounts[setIdx].Ids)-1].End < subbadgeId {
			//Handle case where subbadgeId is less than the starting subbadgeId (append to end)
			insertIdx = len(newAmounts[setIdx].Ids)
			newIds = append(newIds, newAmounts[setIdx].Ids...)
			newIds = append(newIds, &types.NumberRange{Start: subbadgeId})
		} else {
			//Handle case where subbadgeId is somewhere within existing ranges (append and merge if necessary)
			//Binary search to find the insertion index
			low := 0
			high := len(newAmounts[setIdx].Ids) - 2
			// We assume it is removed so it won't be in the middle of an existing range
			for low <= high {
				median = int(uint(low + high) >> 1)
				if newAmounts[setIdx].Ids[median].Start < subbadgeId && newAmounts[setIdx].Ids[median + 1].Start > subbadgeId {
					break;
				} else if newAmounts[setIdx].Ids[median].Start > subbadgeId {
					high = median - 1
				} else {
					low = median + 1
				}
			}
			
			//insertIdx is the index where we need to insert the new subbadgeId
			//Ex: [10, 20, 30] and we need to insert 25, insertIdx == 2
			insertIdx = median + 1
			if (newAmounts[setIdx].Ids[median].Start <= subbadgeId) {
				insertIdx = median
			}
			newIds = append(newIds, newAmounts[setIdx].Ids[:insertIdx]...)
			newIds = append(newIds, &types.NumberRange{
				Start: subbadgeId,
			})
			newIds = append(newIds, newAmounts[setIdx].Ids[insertIdx:]...)
		}

		//Handle cases where we need to merge with the previous or next range
		needToMergeWithPrev := false
		needToMergeWithNext := false
		prevStartIdx := uint64(0)
		nextEndIdx := uint64(0)

		if insertIdx > 0 {
			prevStartIdx = newIds[insertIdx - 1].Start
			prevEndIdx := newIds[insertIdx - 1].End
			if prevEndIdx == 0 {
				prevEndIdx = newIds[insertIdx - 1].Start
			}

			if prevEndIdx + 1 == subbadgeId {
				needToMergeWithPrev = true
			}
		}

		if insertIdx < len(newAmounts[setIdx].Ids) - 1 {
			nextStartIdx := newIds[insertIdx + 1].Start
			nextEndIdx = newIds[insertIdx + 1].End
			if nextEndIdx == 0 {
				nextEndIdx = newIds[insertIdx + 1].Start
			}

			if nextStartIdx - 1 == subbadgeId {
				needToMergeWithNext = true
			}
		}


		mergedIds := []*types.NumberRange{}
		// 4 Cases: Need to merge with both, just next, just prev, or none
		if needToMergeWithPrev && needToMergeWithNext {
			mergedIds = append(mergedIds, newIds[:insertIdx - 1]...)
			mergedIds = append(mergedIds, &types.NumberRange{
				Start: prevStartIdx,
				End:   nextEndIdx,
			})
			mergedIds = append(mergedIds, newIds[insertIdx + 2:]...)
		} else if needToMergeWithPrev {
			mergedIds = append(mergedIds, newIds[:insertIdx - 1]...)
			mergedIds = append(mergedIds, &types.NumberRange{
				Start: prevStartIdx,
				End:   subbadgeId,
			})
			mergedIds = append(mergedIds, newIds[insertIdx + 1:]...)
		} else if needToMergeWithNext {
			mergedIds = append(mergedIds, newIds[:insertIdx]...)
			mergedIds = append(mergedIds, &types.NumberRange{
				Start: subbadgeId,
				End:   nextEndIdx,
			})
			mergedIds = append(mergedIds, newIds[insertIdx + 2:]...)
		} else {
			mergedIds = newIds
		}

		newAmounts[setIdx].Ids = mergedIds
	}
	
	return newAmounts
}

// Adds a balance for the subbadgeID
func (k Keeper) AddBalanceForSubbadgeId(ctx sdk.Context, badgeBalanceInfo types.BadgeBalanceInfo, subbadgeId uint64, balanceToAdd uint64) (types.BadgeBalanceInfo, error) {
	if balanceToAdd == 0 {
		return badgeBalanceInfo, ErrBalanceIsZero
	}

	currBalance := GetBalanceForSubbadgeId(subbadgeId, badgeBalanceInfo.BalanceAmounts)
	newBalance, err := SafeAdd(currBalance, balanceToAdd)
	if err != nil {
		return badgeBalanceInfo, err
	}

	newAmounts := UpdateBalanceForSubbadgeId(subbadgeId, newBalance, badgeBalanceInfo.BalanceAmounts)
	badgeBalanceInfo.BalanceAmounts = newAmounts

	return badgeBalanceInfo, nil
}

// Removes a balance for the subbadgeID
func (k Keeper) RemoveBalanceForSubbadgeId(ctx sdk.Context, badgeBalanceInfo types.BadgeBalanceInfo, subbadgeId uint64, balanceToRemove uint64) (types.BadgeBalanceInfo, error) {
	if balanceToRemove == 0 {
		return badgeBalanceInfo, ErrBalanceIsZero
	}

	currBalance := GetBalanceForSubbadgeId(subbadgeId, badgeBalanceInfo.BalanceAmounts)
	if currBalance < balanceToRemove {
		return badgeBalanceInfo, ErrBadgeBalanceTooLow
	}

	newBalance, err := SafeSubtract(currBalance, balanceToRemove)
	if err != nil {
		return badgeBalanceInfo, err
	}
	
	newAmounts := UpdateBalanceForSubbadgeId(subbadgeId, newBalance, badgeBalanceInfo.BalanceAmounts)
	badgeBalanceInfo.BalanceAmounts = newAmounts

	return badgeBalanceInfo, nil
}
