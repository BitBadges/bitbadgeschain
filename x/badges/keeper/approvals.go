package keeper

import (
	"math"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

// Sets an approval amount for an address.
func SetApproval(userBalance types.UserBalanceStore, amount uint64, address string, badgeIds []*types.IdRange) (types.UserBalanceStore, error) {
	idx, found := SearchApprovals(address, userBalance.Approvals)
	err := *new(error)

	approval := &types.Approval{
		Address: address,
		Balances: []*types.Balance{},
	}

	if found {
		approval = userBalance.Approvals[idx]
	}

	approval.Balances, err = UpdateBalancesForIdRanges(badgeIds, amount, approval.Balances)
	if err != nil {
		return userBalance, err
	}

	//If this is a new entry, we simply appends, otherwise we update the existing entry
	if !found {
		//Insert at idx
		newApprovals := []*types.Approval{}
		newApprovals = append(newApprovals, userBalance.Approvals[:idx]...)
		newApprovals = append(newApprovals, approval)
		newApprovals = append(newApprovals, userBalance.Approvals[idx:]...)
		userBalance.Approvals = newApprovals
	} else {
		userBalance.Approvals[idx] = approval

		if len(userBalance.Approvals[idx].Balances) == 0 {
			//If we end up in the event where this address does not have any more approvals after being removed, we don't have to store it anymore
			userBalance.Approvals = append(userBalance.Approvals[:idx], userBalance.Approvals[idx+1:]...)
		}
	}

	return userBalance, nil
}

// Remove a balance from the approval amount for address
func RemoveBalanceFromApproval(userBalance types.UserBalanceStore, amountToRemove uint64, address string, badgeIdRanges []*types.IdRange) (types.UserBalanceStore, error) {
	err := *new(error)
	if amountToRemove == 0 {
		return userBalance, nil
	}

	idx, found := SearchApprovals(address, userBalance.Approvals)
	if !found {
		return userBalance, ErrApprovalForAddressDoesntExist
	}

	approval := userBalance.Approvals[idx]

	//This may be a bit confusing because we have the following structure:
	//	userBalance.Approvals is of type []Approval
	//	Approval is defined as { Address: string; Balances: []*types.Balance }

	//Basic flow is we get the current approval amounts and ranges in currApprovalAmounts for all IDs in our specified badgeIdRange,
	//and for each unique balance found (which also has its own corresponding []IdRange), we update the balances to balance - amountToRemove.
	currApprovalAmounts, err := GetBalancesForIdRanges(badgeIdRanges, approval.Balances)
	if err != nil {
		return userBalance, err
	}

	for _, currApprovalAmountObj := range currApprovalAmounts {
		newBalance, err := SafeSubtract(currApprovalAmountObj.Amount, amountToRemove)
		if err != nil {
			return userBalance, err
		}

		approval.Balances, err = UpdateBalancesForIdRanges(currApprovalAmountObj.BadgeIds, newBalance, approval.Balances)
		if err != nil {
			return userBalance, err
		}
	}

	userBalance.Approvals[idx].Balances = approval.Balances
	
	//If we end up in the event where this address does not have any more approvals after being removed, we don't have to store it anymore
	if len(approval.Balances) == 0 {
		userBalance.Approvals = append(userBalance.Approvals[:idx], userBalance.Approvals[idx+1:]...)
	}

	return userBalance, nil
}

// Add a balance to the approval amount
func AddBalanceToApproval(userBalance types.UserBalanceStore, amountToAdd uint64, address string, badgeIdRanges []*types.IdRange) (types.UserBalanceStore, error) {
	err := *new(error)
	if amountToAdd == 0 {
		return userBalance, nil
	}

	idx, found := SearchApprovals(address, userBalance.Approvals)
	if !found {
		newApprovals := []*types.Approval{}
		newApprovals = append(newApprovals, userBalance.Approvals[:idx]...)
		newApprovals = append(newApprovals, &types.Approval{
			Address: address,
			Balances: []*types.Balance{
				{
					Amount:  amountToAdd,
					BadgeIds: badgeIdRanges,
				},
			},
		})
		newApprovals = append(newApprovals, userBalance.Approvals[idx:]...)
		userBalance.Approvals = newApprovals
		return userBalance, nil
	}

	//This may be a bit confusing because we have the following structure:
	//	userBalance.Approvals is of type []Approval
	//	Approval is defined as { Address: uint64; Balances: []*types.Balance }

	//Basic flow is we get the current approval amounts and ranges in currApprovalAmounts for all IDs in our specified badgeIdRange,
	//and for each unique balance found (which also has its own corresponding []IdRange), we update the balances to balance + amountToAdd
	approval := userBalance.Approvals[idx]
	currApprovalAmounts, err := GetBalancesForIdRanges(badgeIdRanges, approval.Balances)
	if err != nil {
		return userBalance, err
	}

	for _, currApprovalAmountObj := range currApprovalAmounts {
		newBalance, err := SafeAdd(currApprovalAmountObj.Amount, amountToAdd)
		if err != nil {
			newBalance = math.MaxUint64
		}

		approval.Balances, err = UpdateBalancesForIdRanges(currApprovalAmountObj.BadgeIds, newBalance, approval.Balances)
		if err != nil {
			return userBalance, err
		}
	}

	userBalance.Approvals[idx].Balances = approval.Balances

	return userBalance, nil
}

// Approvals will be sorted alphabetically, so we can binary search to get the targetIdx.
// If found, returns (the index it was found at, true). Else, returns (index to insert at (i.e. the end), false).
func SearchApprovals(targetAddress string, approvals []*types.Approval) (int, bool) {
	low := 0
	high := len(approvals) - 1
	median := 0
	matchingEntry := false
	idx := 0
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

	//Adjust idx to be the insertionIdx, if not found
	if len(approvals) != 0 {
		idx = median + 1
		if targetAddress <= approvals[median].Address {
			idx = median
		}
	}

	return idx, matchingEntry
}
