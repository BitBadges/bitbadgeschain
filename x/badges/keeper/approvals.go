package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

// Sets an approval amount for an address, expirationTime pair.
func (k Keeper) SetApproval(ctx sdk.Context, userBalanceInfo types.UserBalanceInfo, amount uint64, address_num uint64, subbadgeRange types.IdRange, expirationTime uint64) (types.UserBalanceInfo, error) {
	newApprovals := []*types.Approval{}
	found := false
	//check for approval with same address / amount

	//TODO: binary search
	for _, approval := range userBalanceInfo.Approvals {
		if approval.Address != address_num || approval.ExpirationTime != expirationTime {
			newApprovals = append(newApprovals, approval)
		} else {
			found = true
			//Remove completely if setting to zero
			if amount != 0 {
				newAmounts := approval.ApprovalAmounts
				for i := subbadgeRange.Start; i <= subbadgeRange.End; i++ {
					newAmounts = UpdateBalanceForId(i, amount, newAmounts)
				}

				approval.ApprovalAmounts = newAmounts

				newApprovals = append(newApprovals, approval)
			}
		}
	}

	if !found {
		//Add new approval
		newApprovals = append(newApprovals, &types.Approval{
			Address: address_num,
			ApprovalAmounts: []*types.BalanceObject{
				{
					Balance: amount,
					IdRanges: []*types.IdRange{&subbadgeRange},
				},
			},
			ExpirationTime: expirationTime,
		})
	}

	//TODO: sort by address_num

	userBalanceInfo.Approvals = newApprovals

	return userBalanceInfo, nil
}

//Will return an error if isn't approved for amounts
func (k Keeper) RemoveBalanceFromApproval(ctx sdk.Context, userBalanceInfo types.UserBalanceInfo, amount_to_remove uint64, address_num uint64, subbadgeRange types.IdRange) (types.UserBalanceInfo, error) {
	newApprovals := []*types.Approval{}
	removed := false

	//check for approval with same address / amount
	for _, approval := range userBalanceInfo.Approvals {
		if approval.Address == address_num {
			newAmounts := approval.ApprovalAmounts
			for i := subbadgeRange.Start; i <= subbadgeRange.End; i++ {
				currAmount := GetBalanceForId(i, approval.ApprovalAmounts)
				if currAmount < amount_to_remove {
					return userBalanceInfo, ErrInsufficientApproval
				}

				newAmount, err := SafeSubtract(currAmount, amount_to_remove)
				if err != nil {
					return userBalanceInfo, err
				}

				newAmounts = UpdateBalanceForId(i, newAmount, newAmounts)
			}

			approval.ApprovalAmounts = newAmounts

			newApprovals = append(newApprovals, approval)

			removed = true
		} else {
			newApprovals = append(newApprovals, approval)
		}
	}

	if !removed {
		return userBalanceInfo, ErrInsufficientApproval
	}

	if len(newApprovals) == 0 {
		userBalanceInfo.Approvals = nil
	} else {
		userBalanceInfo.Approvals = newApprovals
	}

	return userBalanceInfo, nil
}

func (k Keeper) AddBalanceToApproval(ctx sdk.Context, userBalanceInfo types.UserBalanceInfo, amount_to_add uint64, address_num uint64, subbadgeRange types.IdRange) (types.UserBalanceInfo, error) {
	ctx.GasMeter().ConsumeGas(SimpleAdjustBalanceOrApproval, "adjust approval")

	newApprovals := []*types.Approval{}
	//check for approval with same address / amount
	for _, approval := range userBalanceInfo.Approvals {
		if approval.Address == address_num {
			newAmounts := approval.ApprovalAmounts
			for i := subbadgeRange.Start; i <= subbadgeRange.End; i++ {
				currAmount := GetBalanceForId(i, newAmounts)
				newAmount, err := SafeAdd(currAmount, amount_to_add)
				if err != nil {
					return userBalanceInfo, err
				}

				newAmounts = UpdateBalanceForId(i, newAmount, newAmounts)
			}

			approval.ApprovalAmounts = newAmounts

			newApprovals = append(newApprovals, approval)
		} else {
			newApprovals = append(newApprovals, approval)
		}
	}

	userBalanceInfo.Approvals = newApprovals

	return userBalanceInfo, nil

}

//Prune expired approvals to save space
func PruneExpiredApprovals(currTime uint64, approvals []*types.Approval) []*types.Approval {
	prunedApprovals := make([]*types.Approval, 0)
	for _, approval := range approvals {
		if approval.ExpirationTime != 0 && approval.ExpirationTime < currTime {
			continue
		} else {
			prunedApprovals = append(prunedApprovals, approval)
		}
	}
	return prunedApprovals
}