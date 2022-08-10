package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

// Sets an approval amount for an address, expirationTime pair.
func (k Keeper) SetApproval(ctx sdk.Context, badgeBalanceInfo types.BadgeBalanceInfo, amount uint64, address_num uint64, subbadgeRange types.NumberRange, expirationTime uint64) (types.BadgeBalanceInfo, error) {
	new_approvals := []*types.Approval{}
	found := false
	//check for approval with same address / amount

	//TODO: binary search
	for _, approval := range badgeBalanceInfo.Approvals {
		if approval.Address != address_num || approval.ExpirationTime != expirationTime {
			new_approvals = append(new_approvals, approval)
		} else {
			found = true
			//Remove completely if setting to zero
			if amount != 0 {
				newAmounts := approval.ApprovalAmounts
				for i := subbadgeRange.Start; i <= subbadgeRange.End; i++ {
					newAmounts = UpdateBalanceForSubbadgeId(i, amount, newAmounts)
				}

				approval.ApprovalAmounts = newAmounts

				new_approvals = append(new_approvals, approval)
			}
		}
	}

	if !found {
		//Add new approval
		new_approvals = append(new_approvals, &types.Approval{
			Address: address_num,
			ApprovalAmounts: []*types.BalanceToIds{
				{
					Balance: amount,
					Ids: []*types.NumberRange{&subbadgeRange},
				},
			},
			ExpirationTime: expirationTime,
		})
	}

	//TODO: sort by address_num

	badgeBalanceInfo.Approvals = new_approvals

	return badgeBalanceInfo, nil
}

//Will return an error if isn't approved for amounts
func (k Keeper) RemoveBalanceFromApproval(ctx sdk.Context, badgeBalanceInfo types.BadgeBalanceInfo, amount_to_remove uint64, address_num uint64, subbadgeRange types.NumberRange) (types.BadgeBalanceInfo, error) {
	ctx.GasMeter().ConsumeGas(SimpleAdjustBalanceOrApproval, "adjust approval")

	new_approvals := []*types.Approval{}
	removed := false

	//check for approval with same address / amount
	for _, approval := range badgeBalanceInfo.Approvals {
		if approval.Address == address_num {
			newAmounts := approval.ApprovalAmounts
			for i := subbadgeRange.Start; i <= subbadgeRange.End; i++ {
				currAmount := GetBalanceForSubbadgeId(i, approval.ApprovalAmounts)
				if currAmount < amount_to_remove {
					return badgeBalanceInfo, ErrInsufficientApproval
				}

				newAmount, err := SafeSubtract(currAmount, amount_to_remove)
				if err != nil {
					return badgeBalanceInfo, err
				}

				newAmounts = UpdateBalanceForSubbadgeId(i, newAmount, newAmounts)
			}

			approval.ApprovalAmounts = newAmounts

			new_approvals = append(new_approvals, approval)

			removed = true
		} else {
			new_approvals = append(new_approvals, approval)
		}
	}

	if !removed {
		return badgeBalanceInfo, ErrInsufficientApproval
	}

	if len(new_approvals) == 0 {
		badgeBalanceInfo.Approvals = nil
	} else {
		badgeBalanceInfo.Approvals = new_approvals
	}

	return badgeBalanceInfo, nil
}

func (k Keeper) AddBalanceToApproval(ctx sdk.Context, badgeBalanceInfo types.BadgeBalanceInfo, amount_to_add uint64, address_num uint64, subbadgeRange types.NumberRange) (types.BadgeBalanceInfo, error) {
	ctx.GasMeter().ConsumeGas(SimpleAdjustBalanceOrApproval, "adjust approval")

	new_approvals := []*types.Approval{}
	//check for approval with same address / amount
	for _, approval := range badgeBalanceInfo.Approvals {
		if approval.Address == address_num {
			newAmounts := approval.ApprovalAmounts
			for i := subbadgeRange.Start; i <= subbadgeRange.End; i++ {
				currAmount := GetBalanceForSubbadgeId(i, newAmounts)
				newAmount, err := SafeAdd(currAmount, amount_to_add)
				if err != nil {
					return badgeBalanceInfo, err
				}

				newAmounts = UpdateBalanceForSubbadgeId(i, newAmount, newAmounts)
			}

			approval.ApprovalAmounts = newAmounts

			new_approvals = append(new_approvals, approval)
		} else {
			new_approvals = append(new_approvals, approval)
		}
	}

	badgeBalanceInfo.Approvals = new_approvals

	return badgeBalanceInfo, nil

}
