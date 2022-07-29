package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

//TODO: overflow / underflow handlers

func GetEmptyBadgeBalanceTemplate() types.BadgeBalanceInfo {
	badgeBalanceInfo := types.BadgeBalanceInfo{
		Balance:      0,
		PendingNonce: 0,
		Pending:      []*types.PendingTransfer{},
		Approvals:    []*types.Approval{},
	}
	return badgeBalanceInfo
}

func (k Keeper) AddToBadgeBalance(ctx sdk.Context, balance_key string, balance_to_add uint64) error {
	if balance_to_add == 0 {
		return ErrBalanceIsZero
	}

	badgeBalanceInfo, found := k.GetBadgeBalanceFromStore(ctx, balance_key)
	if !found {
		badgeBalanceInfo = GetEmptyBadgeBalanceTemplate()
		badgeBalanceInfo.Balance += balance_to_add
		err := k.CreateBadgeBalanceInStore(ctx, balance_key, badgeBalanceInfo)
		if err != nil {
			return err
		}
	} else {
		badgeBalanceInfo.Balance += balance_to_add
		err := k.UpdateBadgeBalanceInStore(ctx, balance_key, badgeBalanceInfo)
		if err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) RemoveFromBadgeBalance(ctx sdk.Context, balance_key string, balance_to_remove uint64) error {
	if balance_to_remove == 0 {
		return ErrBalanceIsZero
	}

	badgeBalanceInfo, found := k.GetBadgeBalanceFromStore(ctx, balance_key)
	if !found {
		return ErrBadgeBalanceNotExists
	} else {
		if badgeBalanceInfo.Balance < balance_to_remove {
			return ErrBadgeBalanceTooLow
		}

		badgeBalanceInfo.Balance -= balance_to_remove
		err := k.UpdateBadgeBalanceInStore(ctx, balance_key, badgeBalanceInfo)
		if err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) AddToBothPendingBadgeBalances(ctx sdk.Context, badgeId uint64, subbadgeId uint64, to uint64, from uint64, amount uint64, approvedBy uint64, sentByFrom bool) error {
	if amount == 0 {
		return ErrBalanceIsZero
	}

	//Get "From" balance info
	from_balance_key := GetBalanceKey(from, badgeId, subbadgeId)
	fromBadgeBalanceInfo, fromFound := k.GetBadgeBalanceFromStore(ctx, from_balance_key)
	if !fromFound {
		fromBadgeBalanceInfo = GetEmptyBadgeBalanceTemplate()
	}

	//Get "To" balance info
	to_balance_key := GetBalanceKey(to, badgeId, subbadgeId)
	toBadgeBalanceInfo, toFound := k.GetBadgeBalanceFromStore(ctx, to_balance_key)
	if !toFound {
		toBadgeBalanceInfo = GetEmptyBadgeBalanceTemplate()
	}

	//Append pending transfers and update nonces
	fromBadgeBalanceInfo.Pending = append(fromBadgeBalanceInfo.Pending, &types.PendingTransfer{
		Amount:           	amount,
		ApprovedBy:       	approvedBy,
		SendRequest: 	  	sentByFrom,
		To: 		 	  	to,
		From: 		  	  	from,
		ThisPendingNonce: 	fromBadgeBalanceInfo.PendingNonce,
		OtherPendingNonce:  toBadgeBalanceInfo.PendingNonce,
	})

	toBadgeBalanceInfo.Pending = append(toBadgeBalanceInfo.Pending, &types.PendingTransfer{
		Amount:           	amount,
		ApprovedBy:       	approvedBy,
		SendRequest: 	  	!sentByFrom,
		To: 		 	  	to,
		From: 		  	  	from,
		ThisPendingNonce: 	toBadgeBalanceInfo.PendingNonce,
		OtherPendingNonce:  fromBadgeBalanceInfo.PendingNonce,
	})

	fromBadgeBalanceInfo.PendingNonce += 1
	toBadgeBalanceInfo.PendingNonce += 1

	//Finally, update the stores
	if !fromFound {
		err := k.CreateBadgeBalanceInStore(ctx, from_balance_key, fromBadgeBalanceInfo)
		if err != nil {
			return err
		}
	} else {
		err := k.UpdateBadgeBalanceInStore(ctx, from_balance_key, fromBadgeBalanceInfo)
		if err != nil {
			return err
		}
	}

	if !toFound {
		err := k.CreateBadgeBalanceInStore(ctx, to_balance_key, toBadgeBalanceInfo)
		if err != nil {
			return err
		}
	} else {
		err := k.UpdateBadgeBalanceInStore(ctx, to_balance_key, toBadgeBalanceInfo)
		if err != nil {
			return err
		}
	}

	return nil
}




func (k Keeper) RemovePending(ctx sdk.Context, balance_key string, this_nonce uint64, other_nonce uint64) error {
	badgeBalanceInfo, found := k.GetBadgeBalanceFromStore(ctx, balance_key)
	if !found {
		return ErrBadgeBalanceNotExists
	} else {
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
			return ErrPendingNotFound
		}

		badgeBalanceInfo.Pending = new_pending
		err := k.UpdateBadgeBalanceInStore(ctx, balance_key, badgeBalanceInfo)
		if err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) SetApproval(ctx sdk.Context, balance_key string, amount uint64, address_num uint64) error {
	badgeBalanceInfo, found := k.GetBadgeBalanceFromStore(ctx, balance_key)
	if !found {
		return ErrBadgeBalanceNotExists
	} else {
		new_approvals := []*types.Approval{}
		//check for approval with same address / amount
		for _, approval := range badgeBalanceInfo.Approvals {
			if approval.AddressNum != address_num {
				new_approvals = append(new_approvals, approval)
			}
		}

		new_approvals = append(new_approvals, &types.Approval{
			Amount:     amount,
			AddressNum: address_num,
		})

		badgeBalanceInfo.Approvals = new_approvals
		k.UpdateBadgeBalanceInStore(ctx, balance_key, badgeBalanceInfo)
		return nil
	}
}

//Will return an error if isn't approved for amounts
func (k Keeper) RemoveBalanceFromApproval(ctx sdk.Context, balance_key string, amount_to_remove uint64, address_num uint64) error {
	badgeBalanceInfo, found := k.GetBadgeBalanceFromStore(ctx, balance_key)
	if !found {
		return ErrBadgeBalanceNotExists
	} else {
		new_approvals := []*types.Approval{}
		//check for approval with same address / amount
		for _, approval := range badgeBalanceInfo.Approvals {
			if approval.AddressNum == address_num {
				if approval.Amount < amount_to_remove {
					return ErrInsufficientApproval
				}

				if approval.Amount-amount_to_remove > 0 {
					new_approvals = append(new_approvals, &types.Approval{
						Amount:     approval.Amount - amount_to_remove,
						AddressNum: address_num,
					})
				}
			} else {
				new_approvals = append(new_approvals, approval)
			}
		}

		badgeBalanceInfo.Approvals = new_approvals
		k.UpdateBadgeBalanceInStore(ctx, balance_key, badgeBalanceInfo)
		return nil
	}
}

func (k Keeper) AddBalanceToApproval(ctx sdk.Context, balance_key string, amount_to_add uint64, address_num uint64) error {
	badgeBalanceInfo, found := k.GetBadgeBalanceFromStore(ctx, balance_key)
	if !found {
		return ErrBadgeBalanceNotExists
	} else {
		new_approvals := []*types.Approval{}
		found := false
		//check for approval with same address / amount
		for _, approval := range badgeBalanceInfo.Approvals {
			if approval.AddressNum == address_num {
				new_approvals = append(new_approvals, &types.Approval{
					Amount:     approval.Amount + amount_to_add,
					AddressNum: address_num,
				})
			} else {
				new_approvals = append(new_approvals, approval)
			}
		}
		if !found {
			new_approvals = append(new_approvals, &types.Approval{
				Amount:     amount_to_add,
				AddressNum: address_num,
			})
		}

		badgeBalanceInfo.Approvals = new_approvals
		k.UpdateBadgeBalanceInStore(ctx, balance_key, badgeBalanceInfo)
		return nil
	}
}
