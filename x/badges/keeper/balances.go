package keeper

import (
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
func UpdateBalanceForId(id uint64, newAmount uint64, balanceObjects []*types.BalanceObject) []*types.BalanceObject {
	balanceObjects = DeleteBalanceForId(id, balanceObjects)
	balanceObjects = SetBalanceForId(id, newAmount, balanceObjects)
	return balanceObjects
}

// Gets a balance for a specific id
func GetBalanceForId(id uint64, balanceObjects []*types.BalanceObject) uint64 {
	for _, balanceObj := range balanceObjects {
		balanceObj.IdRanges = GetIdRangesWithOmitEmptyCaseHandled(balanceObj.IdRanges)

		_, found := SearchIdRangesForId(id, balanceObj.IdRanges)
		if found {
			return balanceObj.Balance
		}
	}

	return 0 //Not found; return 0
}

// Adds a balance for the id
func AddBalanceForId(ctx sdk.Context, userBalanceInfo types.UserBalanceInfo, id uint64, balanceToAdd uint64) (types.UserBalanceInfo, error) {
	currBalance := GetBalanceForId(id, userBalanceInfo.BalanceAmounts)
	newBalance, err := SafeAdd(currBalance, balanceToAdd)
	if err != nil {
		return userBalanceInfo, err
	}

	userBalanceInfo.BalanceAmounts = UpdateBalanceForId(id, newBalance, userBalanceInfo.BalanceAmounts)
	return userBalanceInfo, nil
}

// Subtracts a balance for the id
func SubtractBalanceForId(ctx sdk.Context, userBalanceInfo types.UserBalanceInfo, id uint64, balanceToRemove uint64) (types.UserBalanceInfo, error) {
	currBalance := GetBalanceForId(id, userBalanceInfo.BalanceAmounts)
	newBalance, err := SafeSubtract(currBalance, balanceToRemove)
	if err != nil {
		return userBalanceInfo, err
	}

	userBalanceInfo.BalanceAmounts = UpdateBalanceForId(id, newBalance, userBalanceInfo.BalanceAmounts)
	return userBalanceInfo, nil
}

// Deletes the balance for a specific id.
func DeleteBalanceForId(id uint64, balanceObjects []*types.BalanceObject) []*types.BalanceObject {
	newBalanceObjects := []*types.BalanceObject{}
	for _, balanceObj := range balanceObjects {
		balanceObj.IdRanges = GetIdRangesWithOmitEmptyCaseHandled(balanceObj.IdRanges)

		idx, found := SearchIdRangesForId(id, balanceObj.IdRanges)
		if found {
			newIdRanges := append([]*types.IdRange{}, balanceObj.IdRanges[:idx]...)
			newIdRanges = append(newIdRanges, RemoveIdFromIdRange(id, *balanceObj.IdRanges[idx])...)
			newIdRanges = append(newIdRanges, balanceObj.IdRanges[idx+1:]...)
			balanceObj.IdRanges = newIdRanges
		}

		if len(balanceObj.IdRanges) > 0 {
			newBalanceObjects = append(newBalanceObjects, balanceObj)
		}
	}
	return newBalanceObjects
}

// Sets the balance for a specific id.
func SetBalanceForId(id uint64, amount uint64, balanceObjects []*types.BalanceObject) []*types.BalanceObject {
	if amount == 0 {
		return balanceObjects
	}

	idx, found := SearchBalanceObjectsForBalanceAndGetIdxToInsertIfNotFound(amount, balanceObjects)
	newBalanceObjects := []*types.BalanceObject{}
	if !found {
		newBalanceObjects = append(newBalanceObjects, balanceObjects[:idx]...)
		newBalanceObjects = append(newBalanceObjects, &types.BalanceObject{
			Balance:  amount,
			IdRanges: []*types.IdRange{GetIdRangeToInsert(id, id)},
		})
		newBalanceObjects = append(newBalanceObjects, balanceObjects[idx:]...)
	} else {
		newBalanceObjects = balanceObjects
		oldIdRanges := GetIdRangesWithOmitEmptyCaseHandled(newBalanceObjects[idx].IdRanges)
		newBalanceObjects[idx].IdRanges = InsertIdRange(id, oldIdRanges)
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
