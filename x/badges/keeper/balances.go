package keeper

import (
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func GetEmptyBadgeBalanceTemplate() types.BadgeBalanceInfo {
	return types.BadgeBalanceInfo{}
}

//Don't think it'll ever reach an overflow but this is here just in case
func SafeAdd(left uint64, right uint64) (uint64, error) {
	if left > math.MaxUint64-right {
		return 0, ErrOverflow
	}
	return left + right, nil
}

func SafeSubtract(left uint64, right uint64) (uint64, error) {
	if right > left {
		return 0, ErrOverflow
	}
	return left - right, nil
}

func GetBadgeBalanceFromBalanceAmountsForSubbadgeId(subbadgeId uint64, amounts []*types.RangesToAmounts) uint64 {
	//TODO: binary search
	for _, amountObj := range amounts {
		if len(amountObj.Ranges) == 0 {
			amountObj.Ranges = append(amountObj.Ranges, &types.NumberRange{Start: 0, End: 0})
		}

		for _, idObject := range amountObj.Ranges {
			if idObject.End >= subbadgeId && idObject.Start <= subbadgeId {
				return amountObj.Amount
			} else if idObject.End == 0 && idObject.Start == subbadgeId {
				return amountObj.Amount
			}
		}
	}
	return 0
}

func RemoveBadgeBalanceBySubbadgeId(subbadgeId uint64, amounts []*types.RangesToAmounts) []*types.RangesToAmounts {
	new_amounts := []*types.RangesToAmounts{}
	for _, amountObj := range amounts {
		for _, idObject := range amountObj.Ranges {
			if idObject.End == 0 {
				idObject.End = idObject.Start
			}
		}

		if len(amountObj.Ranges) == 0 {
			amountObj.Ranges = append(amountObj.Ranges, &types.NumberRange{Start: 0, End: 0})
		}

		newIds := []*types.NumberRange{}
		for _, idObject := range amountObj.Ranges {
			if idObject.End >= subbadgeId && idObject.Start <= subbadgeId {
				//Found current subbadge

				//If we still have an existing before range, keep that up until subbadge - 1
				if subbadgeId >= 1 && subbadgeId-1 >= idObject.Start {
					newIds = append(newIds, &types.NumberRange{
						Start: idObject.Start,
						End:   subbadgeId - 1,
					})
				}

				//If we still have an existing after range, start that at subbadge + 1
				if subbadgeId <= math.MaxUint64-1 && subbadgeId+1 <= idObject.End {
					newIds = append(newIds, &types.NumberRange{
						Start: subbadgeId + 1,
						End:   idObject.End,
					})
				}
			} else {
				newIds = append(newIds, idObject)
			}
		}
		if len(newIds) > 0 {
			new_amounts = append(new_amounts, &types.RangesToAmounts{
				Ranges: newIds,
				Amount: amountObj.Amount,
			})
		}
	}
	return new_amounts
}

//Precondition: Must be removed already (balance == 0)
func SetBadgeBalanceBySubbadgeId(subbadgeId uint64, amount uint64, amounts []*types.RangesToAmounts) []*types.RangesToAmounts {
	new_amounts := []*types.RangesToAmounts{}
	balanceFound := false

	for _, amountObj := range amounts {
		if amountObj.Amount != amount {
			new_amounts = append(new_amounts, amountObj)
			continue
		}
		balanceFound = true

		newIds := []*types.NumberRange{}

		if len(amountObj.Ranges) == 0 {
			newIds = append(newIds, &types.NumberRange{
				Start: subbadgeId,
				End:   0,
			})

			new_amounts = append(new_amounts, &types.RangesToAmounts{
				Ranges: newIds,
				Amount: amountObj.Amount,
			})
		} else {
			if len(amountObj.Ranges) > 0 && amountObj.Ranges[0].Start > subbadgeId {
				newIds = append(newIds, &types.NumberRange{
					Start: subbadgeId,
					End:   subbadgeId,
				})
			}

			for i := 0; i < len(amountObj.Ranges); i++ {
				if i >= 1 && subbadgeId > amountObj.Ranges[i-1].End && subbadgeId < amountObj.Ranges[i].Start {
					newIds = append(newIds, &types.NumberRange{
						Start: subbadgeId,
						End:   subbadgeId,
					})
				}

				newIds = append(newIds, &types.NumberRange{
					Start: amountObj.Ranges[i].Start,
					End:   amountObj.Ranges[i].End,
				})
			}

			if len(amountObj.Ranges) > 0 && amountObj.Ranges[len(amountObj.Ranges)-1].End < subbadgeId {
				newIds = append(newIds, &types.NumberRange{
					Start: subbadgeId,
					End:   subbadgeId,
				})
			}

			mergedIds := []*types.NumberRange{
				newIds[0],
			}
			for idx := 1; idx < len(newIds); idx++ {
				if newIds[idx].Start == mergedIds[len(mergedIds)-1].End+1 {
					mergedIds[len(mergedIds)-1].End = newIds[idx].End
				} else {
					mergedIds = append(mergedIds, newIds[idx])
				}
			}

			for idx := 0; idx < len(mergedIds); idx++ {
				if mergedIds[idx].End == mergedIds[idx].Start {
					mergedIds[idx].End = 0
				}
			}

			new_amounts = append(new_amounts, &types.RangesToAmounts{
				Ranges: mergedIds,
				Amount: amountObj.Amount,
			})
		}
	}

	if !balanceFound {
		new_amounts = append(new_amounts, &types.RangesToAmounts{
			Amount: amount,
			Ranges: []*types.NumberRange{{Start: subbadgeId}},
		})
	}

	return new_amounts
}

func UpdateBadgeBalanceBySubbadgeId(subbadgeId uint64, newAmount uint64, amounts []*types.RangesToAmounts) []*types.RangesToAmounts {
	amounts = RemoveBadgeBalanceBySubbadgeId(subbadgeId, amounts)
	if newAmount != 0 {
		amounts = SetBadgeBalanceBySubbadgeId(subbadgeId, newAmount, amounts)
	}

	//TODO: make sure this is all sorted by subbadgeId
	return amounts
}

func (k Keeper) AddToBadgeBalance(ctx sdk.Context, badgeBalanceInfo types.BadgeBalanceInfo, subbadgeId uint64, balance_to_add uint64) (types.BadgeBalanceInfo, error) {
	ctx.GasMeter().ConsumeGas(SimpleAdjustBalanceOrApproval, "simple add balance")
	if balance_to_add == 0 {
		return badgeBalanceInfo, ErrBalanceIsZero
	}

	currBalance := GetBadgeBalanceFromBalanceAmountsForSubbadgeId(subbadgeId, badgeBalanceInfo.BalanceAmounts)
	newBalance, err := SafeAdd(currBalance, balance_to_add)
	if err != nil {
		return badgeBalanceInfo, err
	}

	newAmounts := UpdateBadgeBalanceBySubbadgeId(subbadgeId, newBalance, badgeBalanceInfo.BalanceAmounts)

	badgeBalanceInfo.BalanceAmounts = newAmounts

	return badgeBalanceInfo, nil
}

func (k Keeper) RemoveFromBadgeBalance(ctx sdk.Context, badgeBalanceInfo types.BadgeBalanceInfo, subbadgeId uint64, balance_to_remove uint64) (types.BadgeBalanceInfo, error) {
	ctx.GasMeter().ConsumeGas(SimpleAdjustBalanceOrApproval, "simple remove balance")
	if balance_to_remove == 0 {
		return badgeBalanceInfo, ErrBalanceIsZero
	}

	currBalance := GetBadgeBalanceFromBalanceAmountsForSubbadgeId(subbadgeId, badgeBalanceInfo.BalanceAmounts)
	if currBalance < balance_to_remove {
		return badgeBalanceInfo, ErrBadgeBalanceTooLow
	}

	newBalance, err := SafeSubtract(currBalance, balance_to_remove)
	newAmounts := UpdateBadgeBalanceBySubbadgeId(subbadgeId, newBalance, badgeBalanceInfo.BalanceAmounts)

	badgeBalanceInfo.BalanceAmounts = newAmounts

	if err != nil {
		return badgeBalanceInfo, err
	}

	return badgeBalanceInfo, nil
}

func (k Keeper) AddToBothPendingBadgeBalances(ctx sdk.Context, fromBadgeBalanceInfo types.BadgeBalanceInfo, toBadgeBalanceInfo types.BadgeBalanceInfo, subbadgeRange types.NumberRange, to uint64, from uint64, amount uint64, approvedBy uint64, sentByFrom bool, expirationTime uint64) (types.BadgeBalanceInfo, types.BadgeBalanceInfo, error) {
	ctx.GasMeter().ConsumeGas(AddOrRemovePending*2, "add to both pending balances")
	if amount == 0 {
		return fromBadgeBalanceInfo, toBadgeBalanceInfo, ErrBalanceIsZero
	}

	//Append pending transfers and update nonces
	fromBadgeBalanceInfo.Pending = append(fromBadgeBalanceInfo.Pending, &types.PendingTransfer{
		SubbadgeRange:     &subbadgeRange,
		Amount:            amount,
		ApprovedBy:        approvedBy,
		SendRequest:       sentByFrom,
		To:                to,
		From:              from,
		ThisPendingNonce:  fromBadgeBalanceInfo.PendingNonce,
		OtherPendingNonce: toBadgeBalanceInfo.PendingNonce,
		ExpirationTime:   expirationTime,
	})

	toBadgeBalanceInfo.Pending = append(toBadgeBalanceInfo.Pending, &types.PendingTransfer{
		SubbadgeRange:     &subbadgeRange,
		Amount:            amount,
		ApprovedBy:        approvedBy,
		SendRequest:       !sentByFrom,
		To:                to,
		From:              from,
		ThisPendingNonce:  toBadgeBalanceInfo.PendingNonce,
		OtherPendingNonce: fromBadgeBalanceInfo.PendingNonce,
		ExpirationTime:   expirationTime,
	})

	fromBadgeBalanceInfo.PendingNonce += 1
	toBadgeBalanceInfo.PendingNonce += 1

	return fromBadgeBalanceInfo, toBadgeBalanceInfo, nil
}

func (k Keeper) RemovePending(ctx sdk.Context, badgeBalanceInfo types.BadgeBalanceInfo, this_nonce uint64, other_nonce uint64) (types.BadgeBalanceInfo, error) {
	ctx.GasMeter().ConsumeGas(AddOrRemovePending, "remove from pending")

	new_pending := []*types.PendingTransfer{}
	found := false
	for _, pending_info := range badgeBalanceInfo.Pending {
		if pending_info.ThisPendingNonce != this_nonce || pending_info.OtherPendingNonce != other_nonce {
			new_pending = append(new_pending, pending_info)
		} else {
			found = true
		}
	}

	if !found {
		return badgeBalanceInfo, ErrPendingNotFound
	}

	if len(new_pending) == 0 {
		badgeBalanceInfo.Pending = nil
		badgeBalanceInfo.PendingNonce = 0
	} else {
		badgeBalanceInfo.Pending = new_pending
	}

	return badgeBalanceInfo, nil
}

func (k Keeper) SetApproval(ctx sdk.Context, badgeBalanceInfo types.BadgeBalanceInfo, amount uint64, address_num uint64, subbadgeRange types.NumberRange, expirationTime uint64) (types.BadgeBalanceInfo, error) {
	ctx.GasMeter().ConsumeGas(SimpleAdjustBalanceOrApproval, "adjust approval")

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
					newAmounts = UpdateBadgeBalanceBySubbadgeId(i, amount, newAmounts)
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
			ApprovalAmounts: []*types.RangesToAmounts{
				{
					Amount: amount,
					Ranges: []*types.NumberRange{&subbadgeRange},
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
				currAmount := GetBadgeBalanceFromBalanceAmountsForSubbadgeId(i, approval.ApprovalAmounts)
				if currAmount < amount_to_remove {
					return badgeBalanceInfo, ErrInsufficientApproval
				}

				newAmount, err := SafeSubtract(currAmount, amount_to_remove)
				if err != nil {
					return badgeBalanceInfo, err
				}

				newAmounts = UpdateBadgeBalanceBySubbadgeId(i, newAmount, newAmounts)
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
				currAmount := GetBadgeBalanceFromBalanceAmountsForSubbadgeId(i, newAmounts)
				newAmount, err := SafeAdd(currAmount, amount_to_add)
				if err != nil {
					return badgeBalanceInfo, err
				}

				newAmounts = UpdateBadgeBalanceBySubbadgeId(i, newAmount, newAmounts)
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
