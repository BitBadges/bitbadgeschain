package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

// Handles a full transfer from one user to another. If it is to be a forceful transfer, it will transfer the balances and approvals. If it is a pending transfer, it will add it to the pending transfers.
func (k Keeper) HandleTransfer(ctx sdk.Context, badge types.BitBadge, subbadgeRange types.IdRange, fromUserBalanceInfo types.UserBalanceInfo, toUserBalanceInfo types.UserBalanceInfo, amount uint64, from uint64, to uint64, approvedBy uint64, expirationTime uint64) (types.UserBalanceInfo, types.UserBalanceInfo, error) {
	permissions := types.GetPermissions(badge.Permissions)
	err := *new(error)
	sendingToReservedAddress := false //TODO: implement this; Check if to Address is reserved; if so, we automatically forceful transfer
	
	//If we can forceful transfer, do it. Else, add it to pending.
	if sendingToReservedAddress || permissions.ForcefulTransfers() || badge.Manager == to {
		fromUserBalanceInfo, toUserBalanceInfo, err = k.ForcefulTransfer(ctx, badge, subbadgeRange, fromUserBalanceInfo, toUserBalanceInfo, amount, from, to, approvedBy, expirationTime)
	} else {
		fromUserBalanceInfo, toUserBalanceInfo, err = k.PendingTransfer(ctx, badge, subbadgeRange, fromUserBalanceInfo, toUserBalanceInfo, amount, from, to, approvedBy, expirationTime)
	}

	if err != nil {
		return types.UserBalanceInfo{}, types.UserBalanceInfo{}, err
	}

	return fromUserBalanceInfo, toUserBalanceInfo, nil
}

//Forceful transfers will transfer the balances and deduct from approvals directly without adding it to pending.
func (k Keeper) ForcefulTransfer(ctx sdk.Context, badge types.BitBadge, subbadgeRange types.IdRange, fromUserBalanceInfo types.UserBalanceInfo, toUserBalanceInfo types.UserBalanceInfo, amount uint64, from uint64, to uint64, approvedBy uint64, expirationTime uint64) (types.UserBalanceInfo, types.UserBalanceInfo, error) {
	// 1. Check if the from address is frozen
	// 2. Remove approvals if approvedBy != from
	// 3. Deduct from "From" balance
	// 4. Add to "To" balance
	for currSubbadgeId := subbadgeRange.Start; currSubbadgeId <= subbadgeRange.End; currSubbadgeId++ {
		err := k.AssertAccountNotFrozen(ctx, badge, from)
		if err != nil {
			return types.UserBalanceInfo{}, types.UserBalanceInfo{}, err
		}

		fromUserBalanceInfo, err = k.DeductApprovals(ctx, fromUserBalanceInfo, badge, badge.Id, currSubbadgeId, from, to, approvedBy, amount)
		if err != nil {
			return types.UserBalanceInfo{}, types.UserBalanceInfo{}, err
		}

		fromUserBalanceInfo, err = SubtractBalanceForId(ctx, fromUserBalanceInfo, currSubbadgeId, amount)
		if err != nil {
			return types.UserBalanceInfo{}, types.UserBalanceInfo{}, err
		}

		toUserBalanceInfo, err = AddBalanceForId(ctx, toUserBalanceInfo, currSubbadgeId, amount)
		if err != nil {
			return types.UserBalanceInfo{}, types.UserBalanceInfo{}, err
		}
	}

	return fromUserBalanceInfo, toUserBalanceInfo, nil
}

// Removes balances and approvals, and puts them in escrow. Adds a pending transfer to both parties' pending.
func (k Keeper) PendingTransfer(ctx sdk.Context, badge types.BitBadge, subbadgeRange types.IdRange, fromUserBalanceInfo types.UserBalanceInfo, toUserBalanceInfo types.UserBalanceInfo, amount uint64, from uint64, to uint64, approvedBy uint64, expirationTime uint64) (types.UserBalanceInfo, types.UserBalanceInfo, error) {
	err := *new(error)
	// 1. Check if the from address is frozen
	// 2. Remove approvals if approvedBy != from
	// 3. Deduct from "From" balance
	// 4. Append pending transfers to both parties
	// 5. If the pending tranfer is eventually accepted, we simply add the balance. If it is removed, we revert the balance and approvals.
	for currSubbadgeId := subbadgeRange.Start; currSubbadgeId <= subbadgeRange.End; currSubbadgeId++ {
		err := k.AssertAccountNotFrozen(ctx, badge, from)
		if err != nil {
			return types.UserBalanceInfo{}, types.UserBalanceInfo{}, err
		}
		
		fromUserBalanceInfo, err = k.DeductApprovals(ctx, fromUserBalanceInfo, badge, badge.Id, currSubbadgeId, from, to, approvedBy, amount)
		if err != nil {
			return types.UserBalanceInfo{}, types.UserBalanceInfo{}, err
		}

		fromUserBalanceInfo, err = SubtractBalanceForId(ctx, fromUserBalanceInfo, currSubbadgeId, amount)
		if err != nil {
			return types.UserBalanceInfo{}, types.UserBalanceInfo{}, err
		}
	}

	fromUserBalanceInfo, toUserBalanceInfo, err = k.AppendPendingTransferForBothParties(ctx, fromUserBalanceInfo, toUserBalanceInfo, subbadgeRange, to, from, amount, approvedBy, true, expirationTime)
	if err != nil {
		return types.UserBalanceInfo{}, types.UserBalanceInfo{}, err
	}

	return fromUserBalanceInfo, toUserBalanceInfo, nil
}

// Deduct approvals fromrequester if requester != from
func (k Keeper) DeductApprovals(ctx sdk.Context, userBalanceInfo types.UserBalanceInfo, badge types.BitBadge, badgeId uint64, subbadgeId uint64, from uint64, to uint64, requester uint64, amount uint64) (types.UserBalanceInfo, error) {
	newUserBalanceInfo := userBalanceInfo

	if from != requester {
		postApprovalUserBalanceInfo, err := k.RemoveBalanceFromApproval(ctx, newUserBalanceInfo, amount, requester, types.IdRange{Start: subbadgeId, End: subbadgeId})
		newUserBalanceInfo = postApprovalUserBalanceInfo
		if err != nil {
			return userBalanceInfo, err
		}
	}

	return newUserBalanceInfo, nil
}

// Deduct approvals fromrequester if requester != from
func (k Keeper) RevertEscrowedBalancesAndApprovals(ctx sdk.Context, userBalanceInfo types.UserBalanceInfo, id uint64, from uint64, approvedBy uint64, amount uint64) (types.UserBalanceInfo, error) {
	err := *new(error)
	userBalanceInfo, err = AddBalanceForId(ctx, userBalanceInfo, id, amount)
	if err != nil {
		return types.UserBalanceInfo{}, err
	}

	//If it was sent via an approval, we need to add the approval back
	if approvedBy != from {
		userBalanceInfo, err = k.AddBalanceToApproval(ctx, userBalanceInfo, amount, approvedBy, types.IdRange{Start: id, End: id})
		if err != nil {
			return types.UserBalanceInfo{}, err
		}
	}
	
	return userBalanceInfo, nil
}

// Checks if account is frozen or not. 
func IsAccountFrozen(badge types.BitBadge, permissions types.Permissions, address uint64) bool {
	frozenByDefault := permissions.FrozenByDefault()

	accountIsFrozen := false
	if frozenByDefault {
		_, found := SearchIdRangesForId(address, badge.FreezeRanges)	
		if !found {
			accountIsFrozen = true
		}
	} else {
		_, found := SearchIdRangesForId(address, badge.FreezeRanges)	
		if found {
			accountIsFrozen = true
		}
	}
	return accountIsFrozen
}

// Returns an error if account is Frozen
func (k Keeper) AssertAccountNotFrozen(ctx sdk.Context, badge types.BitBadge, from uint64) (error) {
	permissions := types.GetPermissions(badge.Permissions)

	accountIsFrozen := IsAccountFrozen(badge, permissions, from)
	if accountIsFrozen {
		return ErrAddressFrozen
	}

	return nil
}
