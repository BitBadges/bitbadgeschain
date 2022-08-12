package keeper

import (
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

// Sets an approval amount for an address, expirationTime pair.
func SetApproval(ctx sdk.Context, userBalanceInfo types.UserBalanceInfo, amount uint64, address_num uint64, subbadgeRange types.IdRange, expirationTime uint64) (types.UserBalanceInfo, error) {
	idx, found := SearchApprovalsForMatchingeAndGetIdxToInsertIfNotFound(address_num, userBalanceInfo.Approvals)

	if found {
		approval := userBalanceInfo.Approvals[idx]
		if amount != 0 {
			newAmounts := approval.ApprovalAmounts
			newAmounts = UpdateBalancesForIdRanges([]*types.IdRange{&subbadgeRange}, amount, newAmounts)
			userBalanceInfo.Approvals[idx].ApprovalAmounts = newAmounts
		}
	} else {
		newApprovals := []*types.Approval{}
		newApprovals = append(newApprovals, userBalanceInfo.Approvals[:idx]...)
		if amount != 0 {
			newApprovals = append(newApprovals, &types.Approval{
				Address: address_num,
				ApprovalAmounts: []*types.BalanceObject{
					{
						Balance:  amount,
						IdRanges: []*types.IdRange{&subbadgeRange},
					},
				},
			})
		}
		newApprovals = append(newApprovals, userBalanceInfo.Approvals[idx:]...)

		userBalanceInfo.Approvals = newApprovals
	}

	return userBalanceInfo, nil
}

//Remove a balance from the approval amount
func RemoveBalanceFromApproval(ctx sdk.Context, userBalanceInfo types.UserBalanceInfo, amount_to_remove uint64, address_num uint64, subbadgeRange types.IdRange) (types.UserBalanceInfo, error) {
	idx, found := SearchApprovalsForMatchingeAndGetIdxToInsertIfNotFound(address_num, userBalanceInfo.Approvals)
	if !found {
		return userBalanceInfo, ErrInsufficientApproval
	}

	approval := userBalanceInfo.Approvals[idx]
	newAmounts := approval.ApprovalAmounts
	currBalances := GetBalancesForIdRanges([]*types.IdRange{&subbadgeRange}, approval.ApprovalAmounts)
	for _, balanceObj := range currBalances {
		newBalance, err := SafeSubtract(balanceObj.Balance, amount_to_remove)
		if err != nil {
			return userBalanceInfo, err
		}

		newAmounts = UpdateBalancesForIdRanges(balanceObj.IdRanges, newBalance, approval.ApprovalAmounts)
	}

	userBalanceInfo.Approvals[idx].ApprovalAmounts = newAmounts

	if len(newAmounts) == 0 {
		userBalanceInfo.Approvals = append(userBalanceInfo.Approvals[:idx], userBalanceInfo.Approvals[idx+1:]...)
	}

	return userBalanceInfo, nil
}

//Add a balance to the approval amount
func AddBalanceToApproval(ctx sdk.Context, userBalanceInfo types.UserBalanceInfo, amount_to_add uint64, address_num uint64, subbadgeRange types.IdRange) (types.UserBalanceInfo, error) {
	idx, found := SearchApprovalsForMatchingeAndGetIdxToInsertIfNotFound(address_num, userBalanceInfo.Approvals)
	if !found {
		return userBalanceInfo, ErrInsufficientApproval
	}

	approval := userBalanceInfo.Approvals[idx]
	newAmounts := approval.ApprovalAmounts
	currBalances := GetBalancesForIdRanges([]*types.IdRange{&subbadgeRange}, approval.ApprovalAmounts)
	for _, balanceObj := range currBalances {
		newBalance, err := SafeAdd(balanceObj.Balance, amount_to_add)
		if err != nil {
			newBalance = math.MaxUint64
		}

		newAmounts = UpdateBalancesForIdRanges(balanceObj.IdRanges, newBalance, approval.ApprovalAmounts)
	}

	userBalanceInfo.Approvals[idx].ApprovalAmounts = newAmounts

	return userBalanceInfo, nil
}

// Approvals will be sorted, so we can binary search to get the targetIdx and expirationTime. Returns the index to insert at if not found
func SearchApprovalsForMatchingeAndGetIdxToInsertIfNotFound(targetAddress uint64, approvals []*types.Approval) (int, bool) {
	low := 0
	high := len(approvals) - 1
	median := 0
	matchingEntry := false
	setIdx := 0
	for low <= high {
		median = int(uint(low+high) >> 1)

		if approvals[median].Address == targetAddress {
			matchingEntry = true
			break
		} else if approvals[median].Address > targetAddress {
			high = median - 1
		} else {
			low = median + 1
		}
	}

	if len(approvals) != 0 {
		setIdx = median + 1
		if targetAddress <= approvals[median].Address {
			setIdx = median
		}
	}

	return setIdx, matchingEntry
}
