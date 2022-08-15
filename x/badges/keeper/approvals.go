package keeper

import (
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

// Sets an approval amount for an address.
func SetApproval(ctx sdk.Context, userBalanceInfo types.UserBalanceInfo, amount uint64, addressNum uint64, subbadgeRange *types.IdRange) (types.UserBalanceInfo, error) {
	idx, found := SearchApprovals(addressNum, userBalanceInfo.Approvals)

	if found {
		//Update existing approval object for address
		approval := userBalanceInfo.Approvals[idx]
		newAmounts := approval.ApprovalAmounts
		newAmounts = UpdateBalancesForIdRanges([]*types.IdRange{subbadgeRange}, amount, newAmounts)
		userBalanceInfo.Approvals[idx].ApprovalAmounts = newAmounts
	} else {
		//Add new approval object for address at idx, if amount != 0
		newApprovals := []*types.Approval{}
		newApprovals = append(newApprovals, userBalanceInfo.Approvals[:idx]...)
		if amount != 0 {
			newApprovals = append(newApprovals, &types.Approval{
				Address: addressNum,
				ApprovalAmounts: []*types.BalanceObject{
					{
						Balance:  amount,
						IdRanges: []*types.IdRange{subbadgeRange},
					},
				},
			})
		}
		newApprovals = append(newApprovals, userBalanceInfo.Approvals[idx:]...)

		userBalanceInfo.Approvals = newApprovals
	}

	return userBalanceInfo, nil
}

//Remove a balance from the approval amount for address
func RemoveBalanceFromApproval(ctx sdk.Context, userBalanceInfo types.UserBalanceInfo, amountToRemove uint64, addressNum uint64, subbadgeRange *types.IdRange) (types.UserBalanceInfo, error) {
	if amountToRemove == 0 {
		return userBalanceInfo, nil
	}
	
	idx, found := SearchApprovals(addressNum, userBalanceInfo.Approvals)
	if !found {
		return userBalanceInfo, ErrApprovalForAddressDoesntExist
	}

	approval := userBalanceInfo.Approvals[idx]	
	
	//This may be a bit confusing because we have the following structure:
	//	userBalanceInfo.Approvals is of type []Approval
	//	Approval is defined as { Address: uint64; ApprovalAmounts: []*types.BalanceObject }
	
	//Basic flow is we get the current approval amounts and ranges in currApprovalAmounts for all IDs in our specified subbadgeRange, 
	//and for each unique balance found (which also has its own corresponding []IdRange), we update the balances to balance - amountToRemove.
	currApprovalAmounts := GetBalancesForIdRanges([]*types.IdRange{subbadgeRange}, approval.ApprovalAmounts)
	for _, currApprovalAmountObj := range currApprovalAmounts {
		newBalance, err := SafeSubtract(currApprovalAmountObj.Balance, amountToRemove)
		if err != nil {
			return userBalanceInfo, err
		}

		approval.ApprovalAmounts = UpdateBalancesForIdRanges(currApprovalAmountObj.IdRanges, newBalance, approval.ApprovalAmounts)
	}

	userBalanceInfo.Approvals[idx].ApprovalAmounts = approval.ApprovalAmounts
	if len(approval.ApprovalAmounts) == 0 {
		//If we end up in the event where this address does not have any more approvals after being removed, we don't have to store it anymore
		userBalanceInfo.Approvals = append(userBalanceInfo.Approvals[:idx], userBalanceInfo.Approvals[idx+1:]...)
	}

	return userBalanceInfo, nil
}

//Add a balance to the approval amount
func AddBalanceToApproval(ctx sdk.Context, userBalanceInfo types.UserBalanceInfo, amountToAdd uint64, addressNum uint64, subbadgeRange *types.IdRange) (types.UserBalanceInfo, error) {
	if amountToAdd == 0 {
		return userBalanceInfo, nil
	}
	
	idx, found := SearchApprovals(addressNum, userBalanceInfo.Approvals)
	if !found {
		//We just need to add a new approval for this address with only this approval amount
		newApprovals := []*types.Approval{}
		newApprovals = append(newApprovals, userBalanceInfo.Approvals[:idx]...)
		newApprovals = append(newApprovals, &types.Approval{
			Address: addressNum,
			ApprovalAmounts: []*types.BalanceObject{
				{
					Balance:  amountToAdd,
					IdRanges: []*types.IdRange{subbadgeRange},
				},
			},
		})
		newApprovals = append(newApprovals, userBalanceInfo.Approvals[idx:]...)
		userBalanceInfo.Approvals = newApprovals
		return userBalanceInfo, ErrInsufficientApproval
	}

	//This may be a bit confusing because we have the following structure:
	//	userBalanceInfo.Approvals is of type []Approval
	//	Approval is defined as { Address: uint64; ApprovalAmounts: []*types.BalanceObject }
	
	//Basic flow is we get the current approval amounts and ranges in currApprovalAmounts for all IDs in our specified subbadgeRange, 
	//and for each unique balance found (which also has its own corresponding []IdRange), we update the balances to balance + amountToAdd
	approval := userBalanceInfo.Approvals[idx]
	newAmounts := approval.ApprovalAmounts
	currApprovalAmounts := GetBalancesForIdRanges([]*types.IdRange{subbadgeRange}, approval.ApprovalAmounts)
	for _, currApprovalAmountObj := range currApprovalAmounts {
		newBalance, err := SafeAdd(currApprovalAmountObj.Balance, amountToAdd)
		if err != nil {
			newBalance = math.MaxUint64
		}

		newAmounts = UpdateBalancesForIdRanges(currApprovalAmountObj.IdRanges, newBalance, approval.ApprovalAmounts)
	}

	userBalanceInfo.Approvals[idx].ApprovalAmounts = newAmounts

	return userBalanceInfo, nil
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
