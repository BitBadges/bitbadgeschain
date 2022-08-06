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

func GetBadgeBalanceFromIDsAndBalancesForSubbadgeId(subbadgeId uint64, ids []*types.NumberRange, balances []uint64) uint64 {
	//TODO: binary search
	for i, idObject := range ids {
		if idObject.End >= subbadgeId && idObject.Start <= subbadgeId {
			return balances[i]
		} else if idObject.End == 0 && idObject.Start == subbadgeId {
			return balances[i]
		}
	}
	return 0
}

func RemoveBadgeBalanceBySubbadgeId(subbadgeId uint64, ids []*types.NumberRange, balances []uint64) ([]*types.NumberRange, []uint64) {
	for _, idObject := range ids {
		if idObject.End == 0 {
			idObject.End = idObject.Start;
		}
	}

	newIds := []*types.NumberRange{}
	newBalances := []uint64{}
	for i, idObject := range ids {
		if idObject.End >= subbadgeId && idObject.Start <= subbadgeId {
			//Found current subbadge

			//If we still have an existing before range, keep that up until subbadge - 1
			if subbadgeId >= 1 && subbadgeId - 1 >= idObject.Start {
				newIds = append(newIds, &types.NumberRange{
					Start: idObject.Start,
					End:   subbadgeId - 1,
				})
				newBalances = append(newBalances, balances[i])
			}

			//If we still have an existing after range, start that at subbadge + 1
			if subbadgeId <= idObject.End - 1 {
				newIds = append(newIds, &types.NumberRange{
					Start: subbadgeId + 1,
					End:   idObject.End,
				})
				newBalances = append(newBalances, balances[i])
			}
		} else {
			newIds = append(newIds, idObject)
			newBalances = append(newBalances, balances[i])
		}
	}
	return newIds, newBalances
}

//Precondition: Must be removed already (balance == 0)
func SetBadgeBalanceBySubbadgeId(subbadgeId uint64, amount uint64, ids []*types.NumberRange, balances []uint64) ([]*types.NumberRange, []uint64) {
	newIds := []*types.NumberRange{}
	newBalances := []uint64{}

	if len(ids) == 0 {
		newIds = append(newIds, &types.NumberRange{
			Start: subbadgeId,
			End:   0,
		})
		newBalances = append(newBalances, amount)

		return newIds, newBalances
	}
	
	if len(ids) > 0 && ids[0].Start > subbadgeId {
		newIds = append(newIds, &types.NumberRange{
			Start: subbadgeId,
			End:   subbadgeId,
		})
		newBalances = append(newBalances, amount)
	}

	for i := 0; i < len(ids); i++ {
		if i >= 1 && subbadgeId > ids[i - 1].End && subbadgeId < ids[i].Start {
			newIds = append(newIds, &types.NumberRange{
				Start: subbadgeId,
				End:   subbadgeId,
			})
			newBalances = append(newBalances, amount)
		}

		newIds = append(newIds, &types.NumberRange{
			Start: ids[i].Start,
			End:   ids[i].End,
		})
		newBalances = append(newBalances, balances[i])
	}

	if len(ids) > 0 && ids[len(ids)-1].End < subbadgeId {
		newIds = append(newIds, &types.NumberRange{
			Start: subbadgeId,
			End:   subbadgeId,
		})
		newBalances = append(newBalances, amount)
	}

	mergedIds := []*types.NumberRange{
		newIds[0],
	}
	mergedBalances := []uint64{
		newBalances[0],
	}
	for idx := 1; idx < len(newIds); idx++ {
		if newIds[idx].Start == mergedIds[len(mergedIds)-1].End + 1 && newBalances[idx] == mergedBalances[len(mergedBalances)-1] {
			mergedIds[len(mergedIds)-1].End = newIds[idx].End
		} else {
			mergedIds = append(mergedIds, newIds[idx])
			mergedBalances = append(mergedBalances, newBalances[idx])
		}
	}

	for idx := 0; idx < len(mergedIds); idx++ {
		if mergedIds[idx].End == mergedIds[idx].Start {
			mergedIds[idx].End = 0;
		}
	}

	return mergedIds, mergedBalances
}


func UpdateBadgeBalanceBySubbadgeId(subbadgeId uint64, newAmount uint64, ids []*types.NumberRange, balances []uint64) ([]*types.NumberRange, []uint64) {
	ids, balances = RemoveBadgeBalanceBySubbadgeId(subbadgeId, ids, balances)
	if newAmount != 0 {
		ids, balances = SetBadgeBalanceBySubbadgeId(subbadgeId, newAmount, ids, balances)
	}
	return ids, balances
}

func (k Keeper) AddToBadgeBalance(ctx sdk.Context, badgeBalanceInfo types.BadgeBalanceInfo, subbadgeId uint64, balance_to_add uint64) (types.BadgeBalanceInfo, error) {
	ctx.GasMeter().ConsumeGas(SimpleAdjustBalanceOrApproval, "simple add balance")
	if balance_to_add == 0 {
		return badgeBalanceInfo, ErrBalanceIsZero
	}

	currBalance := GetBadgeBalanceFromIDsAndBalancesForSubbadgeId(subbadgeId, badgeBalanceInfo.IdsForBalances, badgeBalanceInfo.Balances)
	newBalance, err := SafeAdd(currBalance, balance_to_add)
	if err != nil {
		return badgeBalanceInfo, err
	}

	newIds, newBalances := UpdateBadgeBalanceBySubbadgeId(subbadgeId, newBalance, badgeBalanceInfo.IdsForBalances, badgeBalanceInfo.Balances)

	badgeBalanceInfo.Balances = newBalances
	badgeBalanceInfo.IdsForBalances = newIds
	
	return badgeBalanceInfo, nil
}

func (k Keeper) RemoveFromBadgeBalance(ctx sdk.Context, badgeBalanceInfo types.BadgeBalanceInfo, subbadgeId uint64, balance_to_remove uint64) (types.BadgeBalanceInfo, error) {
	ctx.GasMeter().ConsumeGas(SimpleAdjustBalanceOrApproval, "simple remove balance")
	if balance_to_remove == 0 {
		return badgeBalanceInfo, ErrBalanceIsZero
	}

	
	currBalance := GetBadgeBalanceFromIDsAndBalancesForSubbadgeId(subbadgeId, badgeBalanceInfo.IdsForBalances, badgeBalanceInfo.Balances)
	if currBalance < balance_to_remove {
		return badgeBalanceInfo, ErrBadgeBalanceTooLow
	}

	newBalance, err := SafeSubtract(currBalance, balance_to_remove)
	newIds, newBalances := UpdateBadgeBalanceBySubbadgeId(subbadgeId, newBalance, badgeBalanceInfo.IdsForBalances, badgeBalanceInfo.Balances)
	
	badgeBalanceInfo.Balances = newBalances
	badgeBalanceInfo.IdsForBalances = newIds
	if err != nil {
		return badgeBalanceInfo, err
	}

	return badgeBalanceInfo, nil
}

func (k Keeper) AddToBothPendingBadgeBalances(ctx sdk.Context, fromBadgeBalanceInfo types.BadgeBalanceInfo, toBadgeBalanceInfo types.BadgeBalanceInfo, subbadgeRange types.NumberRange, to uint64, from uint64, amount uint64, approvedBy uint64, sentByFrom bool) (types.BadgeBalanceInfo, types.BadgeBalanceInfo, error) {
	ctx.GasMeter().ConsumeGas(AddOrRemovePending * 2, "add to both pending balances")
	if amount == 0 {
		return fromBadgeBalanceInfo, toBadgeBalanceInfo, ErrBalanceIsZero
	}

	//Append pending transfers and update nonces
	fromBadgeBalanceInfo.Pending = append(fromBadgeBalanceInfo.Pending, &types.PendingTransfer{
		NumberRange: &subbadgeRange,
		Amount:            amount,
		ApprovedBy:        approvedBy,
		SendRequest:       sentByFrom,
		To:                to,
		From:              from,
		ThisPendingNonce:  fromBadgeBalanceInfo.PendingNonce,
		OtherPendingNonce: toBadgeBalanceInfo.PendingNonce,
	})

	toBadgeBalanceInfo.Pending = append(toBadgeBalanceInfo.Pending, &types.PendingTransfer{
		NumberRange: &subbadgeRange,
		Amount:            amount,
		ApprovedBy:        approvedBy,
		SendRequest:       !sentByFrom,
		To:                to,
		From:              from,
		ThisPendingNonce:  toBadgeBalanceInfo.PendingNonce,
		OtherPendingNonce: fromBadgeBalanceInfo.PendingNonce,
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

func (k Keeper) SetApproval(ctx sdk.Context, badgeBalanceInfo types.BadgeBalanceInfo, amount uint64, address_num uint64, subbadgeRange types.NumberRange) (types.BadgeBalanceInfo, error) {
	ctx.GasMeter().ConsumeGas(SimpleAdjustBalanceOrApproval, "adjust approval")
	

	new_approvals := []*types.Approval{}
	found := false
	//check for approval with same address / amount
	for _, approval := range badgeBalanceInfo.Approvals {
		if approval.Address != address_num {
			new_approvals = append(new_approvals, approval)
		} else {
			found = true
			//Remove completely if setting to zero
			if amount != 0 {
				newSubbadgeRanges := approval.SubbadgeRanges
				newAmounts := approval.Amounts
				for i := subbadgeRange.Start; i <= subbadgeRange.End; i++ {
					newSubbadgeRanges, newAmounts = UpdateBadgeBalanceBySubbadgeId(i, amount, newSubbadgeRanges, newAmounts)
				}

				approval.Amounts = newAmounts;
				approval.SubbadgeRanges = newSubbadgeRanges;

				new_approvals = append(new_approvals, approval)
			}

		}
	}

	if !found {
		//Add new approval
		new_approvals = append(new_approvals, &types.Approval{
			Address: address_num,
			Amounts: []uint64{amount},
			SubbadgeRanges: []*types.NumberRange{&subbadgeRange},
		})
	}


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
			newSubbadgeRanges := approval.SubbadgeRanges
			newAmounts := approval.Amounts
			for i := subbadgeRange.Start; i <= subbadgeRange.End; i++ {
				currAmount := GetBadgeBalanceFromIDsAndBalancesForSubbadgeId(i, approval.SubbadgeRanges, approval.Amounts)
				if currAmount < amount_to_remove {
					return badgeBalanceInfo, ErrInsufficientApproval
				}
				
				
				newAmount, err := SafeSubtract(currAmount, amount_to_remove)
				if err != nil {
					return badgeBalanceInfo, err
				}

				newSubbadgeRanges, newAmounts  = UpdateBadgeBalanceBySubbadgeId(i, newAmount, approval.SubbadgeRanges, approval.Amounts)
			}

			approval.Amounts = newAmounts;
			approval.SubbadgeRanges = newSubbadgeRanges;

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
			newSubbadgeRanges := approval.SubbadgeRanges
			newAmounts := approval.Amounts
			for i := subbadgeRange.Start; i <= subbadgeRange.End; i++ {
				currAmount := GetBadgeBalanceFromIDsAndBalancesForSubbadgeId(i, approval.SubbadgeRanges, approval.Amounts)				
				newAmount, err := SafeAdd(currAmount, amount_to_add)
				if err != nil {
					return badgeBalanceInfo, err
				}

				newSubbadgeRanges, newAmounts  = UpdateBadgeBalanceBySubbadgeId(i, newAmount, approval.SubbadgeRanges, approval.Amounts)
			}

			approval.Amounts = newAmounts;
			approval.SubbadgeRanges = newSubbadgeRanges;

			new_approvals = append(new_approvals, approval)
		} else {
			new_approvals = append(new_approvals, approval)
		}
	}

	badgeBalanceInfo.Approvals = new_approvals
	

	return badgeBalanceInfo, nil

}
