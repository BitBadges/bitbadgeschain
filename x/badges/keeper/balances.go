package keeper

import (
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

//TODO: overflow / underflow handlers

func GetFullSubassetID(accountNumber uint64, id uint64, subasset_id uint64) string {
	badge_id_str := strconv.FormatUint(id, 10)
	subasset_id_str := strconv.FormatUint(subasset_id, 10)
	account_num_str := strconv.FormatUint(accountNumber, 10)
	return account_num_str + "-" + badge_id_str + "-" + subasset_id_str
}

func GetPendingID(first uint64, second uint64) string {
	first_nonce := strconv.FormatUint(first, 10)
	second_nonce := strconv.FormatUint(second, 10)
	return first_nonce + "-" + second_nonce
}

type PendingIdDetails struct {
	FirstNonce  uint64
	SecondNonce uint64
}

func GetDetailsFromPendingID(id string) PendingIdDetails {
	result := strings.Split(id, "-")
	first, _ := strconv.ParseUint(result[0], 10, 64)
	second, _ := strconv.ParseUint(result[1], 10, 64)

	return PendingIdDetails{
		FirstNonce:  first,
		SecondNonce: second,
	}
}

type FullIdDetails struct {
	badge_id    uint64
	subasset_id uint64
	account_num uint64
}

func GetDetailsFromFullSubassetID(id string) FullIdDetails {
	result := strings.Split(id, "-")
	account_num, _ := strconv.ParseUint(result[0], 10, 64)
	badge_id, _ := strconv.ParseUint(result[1], 10, 64)
	subasset_id, _ := strconv.ParseUint(result[2], 10, 64)

	return FullIdDetails{
		account_num: account_num,
		badge_id:    badge_id,
		subasset_id: subasset_id,
	}
}

func GetEmptyBadgeBalanceTemplate() types.BadgeBalanceInfo {
	badgeBalanceInfo := types.BadgeBalanceInfo{
		Balance:      0,
		PendingNonce: 0,
		Pending:      []*types.PendingTransfer{},
		Approvals:    []*types.Approval{},
	}
	return badgeBalanceInfo
}

func (k Keeper) AddToBadgeBalance(ctx sdk.Context, balance_id string, balance_to_add uint64) error {
	if balance_to_add == 0 {
		return ErrBalanceIsZero
	}

	badgeBalanceInfo, found := k.GetBadgeBalanceFromStore(ctx, balance_id)
	if !found {
		badgeBalanceInfo = GetEmptyBadgeBalanceTemplate()
		badgeBalanceInfo.Balance += balance_to_add
		err := k.CreateBadgeBalanceInStore(ctx, balance_id, badgeBalanceInfo)
		if err != nil {
			return err
		}
	} else {
		badgeBalanceInfo.Balance += balance_to_add
		err := k.UpdateBadgeBalanceInStore(ctx, balance_id, badgeBalanceInfo)
		if err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) RemoveFromBadgeBalance(ctx sdk.Context, balance_id string, balance_to_remove uint64) error {
	if balance_to_remove == 0 {
		return ErrBalanceIsZero
	}

	badgeBalanceInfo, found := k.GetBadgeBalanceFromStore(ctx, balance_id)
	if !found {
		return ErrBadgeBalanceNotExists
	} else {
		if badgeBalanceInfo.Balance < balance_to_remove {
			return ErrBadgeBalanceTooLow
		}

		badgeBalanceInfo.Balance -= balance_to_remove
		err := k.UpdateBadgeBalanceInStore(ctx, balance_id, badgeBalanceInfo)
		if err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) AddToPendingBadgeBalance(ctx sdk.Context, balance_id string, pending_info types.PendingTransfer) error {
	if pending_info.Amount == 0 {
		return ErrBalanceIsZero
	}

	badgeBalanceInfo, found := k.GetBadgeBalanceFromStore(ctx, balance_id)
	if !found {
		badgeBalanceInfo = GetEmptyBadgeBalanceTemplate()
		badgeBalanceInfo.Pending = append(badgeBalanceInfo.Pending, &pending_info)
		badgeBalanceInfo.PendingNonce += 1

		err := k.CreateBadgeBalanceInStore(ctx, balance_id, badgeBalanceInfo)
		if err != nil {
			return err
		}
	} else {
		badgeBalanceInfo.Pending = append(badgeBalanceInfo.Pending, &pending_info)
		badgeBalanceInfo.PendingNonce += 1
		err := k.UpdateBadgeBalanceInStore(ctx, balance_id, badgeBalanceInfo)
		if err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) RemoveFromPendingBadgeBalance(ctx sdk.Context, balance_id string, pending_id string) error {
	badgeBalanceInfo, found := k.GetBadgeBalanceFromStore(ctx, balance_id)
	if !found {
		return ErrBadgeBalanceNotExists
	} else {
		new_pending := []*types.PendingTransfer{}
		for _, pending_info := range badgeBalanceInfo.Pending {
			if pending_info.Id != pending_id {
				new_pending = append(new_pending, pending_info)
			}
		}
		badgeBalanceInfo.Pending = new_pending
		err := k.UpdateBadgeBalanceInStore(ctx, balance_id, badgeBalanceInfo)
		if err != nil {
			return err
		}
	}

	return nil
}

//TODO: many of these functions should not be exposed to the keeper aka outside world
//Permissions are not handled here. This is only called to handle the balance transfers.
//Only handles from => to (pending and forceful) (not other way around)
func (k Keeper) TransferBadge(ctx sdk.Context, tx_signer sdk.AccAddress, from sdk.AccAddress, to sdk.AccAddress, amount uint64, badge_id uint64, subasset_id uint64) error {
	err := k.AssertBadgeAndSubBadgeExists(ctx, badge_id, subasset_id)
	if err != nil {
		return err
	}

	//TODO: assert can transfer? or check revoke status? or freeze? etc?

	badge, _ := k.GetBadgeFromStore(ctx, badge_id)

	//TODO: In some instances, you may want to transfer to an unregistered account
	// 	In this case, we should register a new account and not throw
	to_account := k.accountKeeper.GetAccount(ctx, to)
	if to_account == nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "account %s does not exist", to)
	}
	to_account_num := to_account.GetAccountNumber()
	to_balance_id := GetFullSubassetID(to_account_num, badge_id, subasset_id)

	from_account := k.accountKeeper.GetAccount(ctx, from)
	if from_account == nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "account %s does not exist", from)
	}
	from_account_num := from_account.GetAccountNumber()
	from_balance_id := GetFullSubassetID(from_account_num, badge_id, subasset_id)

	//TODO: check approvals
	if !tx_signer.Equals(from) {
		//throw if not approved
	}

	//TODO: check if the account is frozen

	manager_address, err := sdk.AccAddressFromBech32(badge.Manager)
	if err != nil {
		return err
	}

	permissions := GetPermissions(badge.PermissionFlags)

	//Forceful transfers only when permitted to or "burning" (aka sending back to manager)
	if permissions.forceful_transfers || manager_address.Equals(to) {
		k.AddToBadgeBalance(ctx, to_balance_id, amount)
		k.RemoveFromBadgeBalance(ctx, from_balance_id, amount)
	} else {
		//TODO: validate memo and other field
		//Both nonces from balance info
		toPendingNonce := uint64(0)
		fromPendingNonce := uint64(0)

		toBadgeBalanceInfo, found := k.GetBadgeBalanceFromStore(ctx, to_balance_id)
		if found {
			toPendingNonce = toBadgeBalanceInfo.PendingNonce
		}

		fromBadgeBalanceInfo, found := k.GetBadgeBalanceFromStore(ctx, from_balance_id)
		if found {
			fromPendingNonce = fromBadgeBalanceInfo.PendingNonce
		}

		to_transfer := types.PendingTransfer{
			Id:          GetPendingID(toPendingNonce, fromPendingNonce),
			Amount:      amount,
			SendRequest: false,
			To:          to_account_num,
			From:        from_account_num,
			Memo:        "",
		}

		from_transfer := types.PendingTransfer{
			Id:          GetPendingID(fromPendingNonce, toPendingNonce),
			Amount:      amount,
			SendRequest: true,
			To:          to_account_num,
			From:        from_account_num,
			Memo:        "",
		}

		if to_account_num == from_account_num {
			return ErrSenderAndReceiverSame
		}

		//Remove from from's balance
		k.RemoveFromBadgeBalance(ctx, from_balance_id, amount)
		k.AddToPendingBadgeBalance(ctx, to_balance_id, to_transfer)
		k.AddToPendingBadgeBalance(ctx, from_balance_id, from_transfer)
	}

	return nil
}

//Requests to receive a badge. Precondition: from must == msg.Creator
func (k Keeper) RequestTransferBadge(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, amount uint64, badge_id uint64, subasset_id uint64) error {
	err := k.AssertBadgeAndSubBadgeExists(ctx, badge_id, subasset_id)
	if err != nil {
		return err
	}

	//TODO: In some instances, you may want to transfer to an unregistered account
	// 	In this case, we should register a new account and not throw
	to_account := k.accountKeeper.GetAccount(ctx, to)
	if to_account == nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "account %s does not exist", to)
	}
	to_account_num := to_account.GetAccountNumber()
	to_balance_id := GetFullSubassetID(to_account_num, badge_id, subasset_id)

	from_account := k.accountKeeper.GetAccount(ctx, from)
	if from_account == nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "account %s does not exist", from)
	}
	from_account_num := from_account.GetAccountNumber()
	from_balance_id := GetFullSubassetID(from_account_num, badge_id, subasset_id)

	//TODO: validate memo and other field
	//Both nonces from balance info
	toPendingNonce := uint64(0)
	fromPendingNonce := uint64(0)

	toBadgeBalanceInfo, found := k.GetBadgeBalanceFromStore(ctx, to_balance_id)
	if found {
		toPendingNonce = toBadgeBalanceInfo.PendingNonce
	}

	fromBadgeBalanceInfo, found := k.GetBadgeBalanceFromStore(ctx, from_balance_id)
	if found {
		fromPendingNonce = fromBadgeBalanceInfo.PendingNonce
	}

	if to_account_num == from_account_num {
		return ErrSenderAndReceiverSame
	}

	to_transfer := types.PendingTransfer{
		Id:          GetPendingID(toPendingNonce, fromPendingNonce),
		Amount:      amount,
		SendRequest: true,
		To:          to_account_num,
		From:        from_account_num,
		Memo:        "",
	}

	from_transfer := types.PendingTransfer{
		Id:          GetPendingID(fromPendingNonce, toPendingNonce),
		Amount:      amount,
		SendRequest: false,
		To:          to_account_num,
		From:        from_account_num,
		Memo:        "",
	}

	//Remove from from's balance
	k.AddToPendingBadgeBalance(ctx, to_balance_id, to_transfer)
	k.AddToPendingBadgeBalance(ctx, from_balance_id, from_transfer)

	return nil
}

/*
	Sender cancels their transfer
	Sender accepts their transfer (doesn't make sensse; throw)
	Receiver accepts / rejects their transfer

	TODO: think about adding approvals for pending / rejecting
*/
func (k Keeper) HandlePendingTransfer(ctx sdk.Context, accept bool, balance_id string, pending_id string) error {
	//TODO: ensure balance_id matches msg_creator or whoever has valid permission's balance ID
	//For now, we will always assume balance_id is the one who wants to accept / reject
	badgeBalanceInfo, found := k.GetBadgeBalanceFromStore(ctx, balance_id)
	if !found {
		return ErrBadgeBalanceNotExists
	} else {
		for _, pending_info := range badgeBalanceInfo.Pending {
			if pending_info.Id == pending_id {
				target_pending := pending_info

				if target_pending.SendRequest && accept {
					return ErrCantAcceptOwnTransferRequest
				}

				pending_id_details := GetDetailsFromPendingID(pending_id)
				balance_id_details := GetDetailsFromFullSubassetID(balance_id)

				other_account_num := uint64(0)
				if target_pending.From == balance_id_details.account_num {
					other_account_num = target_pending.To
				} else if target_pending.To == balance_id_details.account_num {
					other_account_num = target_pending.From
				} else {
					return ErrInvalidPermissions
				}
				
				other_balance_id := GetFullSubassetID(other_account_num, balance_id_details.badge_id, balance_id_details.subasset_id)
				other_pending_id := GetPendingID(pending_id_details.SecondNonce, pending_id_details.FirstNonce)

				//Remove from pending
				k.RemoveFromPendingBadgeBalance(ctx, balance_id, pending_id)
				k.RemoveFromPendingBadgeBalance(ctx, other_balance_id, other_pending_id)

				//Sent request and want to cancel
				if target_pending.SendRequest && !accept {
					//Cancel a transfer request
					if balance_id_details.account_num == target_pending.To {
						//Do nothing since balance wasn't removed yet
					} else if balance_id_details.account_num == target_pending.From {
						// Cancel an outgoing transfer
						k.AddToBadgeBalance(ctx, balance_id, target_pending.Amount) //Add back funds to this account
					}
				} else if !target_pending.SendRequest {
					to_balance_id := ""
					from_balance_id := ""
					if balance_id_details.account_num == target_pending.To {
						to_balance_id = balance_id
						from_balance_id = other_balance_id
					} else if balance_id_details.account_num == target_pending.From {
						to_balance_id = other_balance_id
						from_balance_id = balance_id
					} else {
						return ErrInvalidPermissions
					}

					if accept {
						//transfer funds to to_account
						k.AddToBadgeBalance(ctx, to_balance_id, target_pending.Amount)
						k.RemoveFromBadgeBalance(ctx, from_balance_id, target_pending.Amount)
					} else {
						if balance_id_details.account_num == target_pending.To {
							k.AddToBadgeBalance(ctx, from_balance_id, target_pending.Amount) //Add back funds to from account
						} else if balance_id_details.account_num == target_pending.From {
							// Do nothing; ignore request for transfer
						}
					}
				}
				return nil
			}
		}
		return sdkerrors.ErrNotFound
	}
}

//Approve / Update Approval
