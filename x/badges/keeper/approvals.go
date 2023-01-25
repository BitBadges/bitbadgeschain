package keeper

import (
	"math"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

// Sets an approval amount for an address.
func SetApproval(userBalance types.UserBalance, amount uint64, addressNum uint64, badgeIdRange *types.IdRange) (types.UserBalance, error) {
	idx, found := SearchApprovals(addressNum, userBalance.Approvals)

	if found {
		//Update existing approval object for address
		approval := userBalance.Approvals[idx]
		newAmounts := approval.Balances
		newAmounts = UpdateBalancesForIdRanges([]*types.IdRange{badgeIdRange}, amount, newAmounts)
		userBalance.Approvals[idx].Balances = newAmounts

		if len(userBalance.Approvals[idx].Balances) == 0 {
			//If we end up in the event where this address does not have any more approvals after being removed, we don't have to store it anymore
			userBalance.Approvals = append(userBalance.Approvals[:idx], userBalance.Approvals[idx+1:]...)
		}
	} else {
		//Add new approval object for address at idx, if amount != 0
		newApprovals := []*types.Approval{}
		newApprovals = append(newApprovals, userBalance.Approvals[:idx]...)
		if amount != 0 {
			newApprovals = append(newApprovals, &types.Approval{
				Address: addressNum,
				Balances: []*types.Balance{
					{
						Balance:  amount,
						BadgeIds: []*types.IdRange{CreateIdRange(badgeIdRange.Start, badgeIdRange.End)},
					},
				},
			})
		}
		newApprovals = append(newApprovals, userBalance.Approvals[idx:]...)

		userBalance.Approvals = newApprovals
	}

	return GetBalanceToInsertToStorage(userBalance), nil
}

//Remove a balance from the approval amount for address
func RemoveBalanceFromApproval(userBalance types.UserBalance, amountToRemove uint64, addressNum uint64, badgeIdRanges []*types.IdRange) (types.UserBalance, error) {
	if amountToRemove == 0 {
		return userBalance, nil
	}

	idx, found := SearchApprovals(addressNum, userBalance.Approvals)
	if !found {
		return userBalance, ErrApprovalForAddressDoesntExist
	}

	approval := userBalance.Approvals[idx]

	//This may be a bit confusing because we have the following structure:
	//	userBalance.Approvals is of type []Approval
	//	Approval is defined as { Address: uint64; Balances: []*types.Balance }

	//Basic flow is we get the current approval amounts and ranges in currApprovalAmounts for all IDs in our specified badgeIdRange,
	//and for each unique balance found (which also has its own corresponding []IdRange), we update the balances to balance - amountToRemove.
	currApprovalAmounts := GetBalancesForIdRanges(badgeIdRanges, approval.Balances)
	for _, currApprovalAmountObj := range currApprovalAmounts {
		newBalance, err := SafeSubtract(currApprovalAmountObj.Balance, amountToRemove)
		if err != nil {
			return userBalance, err
		}

		approval.Balances = UpdateBalancesForIdRanges(currApprovalAmountObj.BadgeIds, newBalance, approval.Balances)
	}

	userBalance.Approvals[idx].Balances = approval.Balances
	if len(approval.Balances) == 0 {
		//If we end up in the event where this address does not have any more approvals after being removed, we don't have to store it anymore
		userBalance.Approvals = append(userBalance.Approvals[:idx], userBalance.Approvals[idx+1:]...)
	}

	return GetBalanceToInsertToStorage(userBalance), nil
}

//Add a balance to the approval amount
func AddBalanceToApproval(userBalance types.UserBalance, amountToAdd uint64, addressNum uint64, badgeIdRanges []*types.IdRange) (types.UserBalance, error) {
	if amountToAdd == 0 {
		return userBalance, nil
	}

	idx, found := SearchApprovals(addressNum, userBalance.Approvals)
	if !found {
		//We just need to add a new approval for this address with only this approval amount
		newApprovals := []*types.Approval{}
		newApprovals = append(newApprovals, userBalance.Approvals[:idx]...)
		idRangesToInsert := []*types.IdRange{}
		for _, badgeIdRange := range badgeIdRanges {
			idRangesToInsert = append(idRangesToInsert, CreateIdRange(badgeIdRange.Start, badgeIdRange.End))
		}

		newApprovals = append(newApprovals, &types.Approval{
			Address: addressNum,
			Balances: []*types.Balance{
				{
					Balance:  amountToAdd,
					BadgeIds: idRangesToInsert,
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
	currApprovalAmounts := GetBalancesForIdRanges(badgeIdRanges, approval.Balances)
	for _, currApprovalAmountObj := range currApprovalAmounts {
		newBalance, err := SafeAdd(currApprovalAmountObj.Balance, amountToAdd)
		if err != nil {
			newBalance = math.MaxUint64
		}

		approval.Balances = UpdateBalancesForIdRanges(currApprovalAmountObj.BadgeIds, newBalance, approval.Balances)
	}

	userBalance.Approvals[idx].Balances = approval.Balances

	return GetBalanceToInsertToStorage(userBalance), nil
}

// Approvals will be sorted, so we can binary search to get the targetIdx.
// If found, returns (the index it was found at, true). Else, returns (index to insert at, false).
func SearchApprovals(targetAddress uint64, approvals []*types.Approval) (int, bool) {
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
