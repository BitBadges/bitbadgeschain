package keeper

import (
	"strconv"

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

func GetEmptyBadgeBalanceTemplate() types.BadgeBalanceInfo {
	badgeBalanceInfo := types.BadgeBalanceInfo{
		Balance:      0,
		PendingNonce: 0,
		Pending:      []*types.PendingTransfer{},
		Approvals:    []*types.Approval{},
	}
	return badgeBalanceInfo
}

//assumes already in store
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

//assumes already in store
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


//TODO: many of these functions should not be exposed to the keeper aka outside world
//Permissions are not handled here. This is only called to handle the balance transfers. Assumed to be approved to transfer
func (k Keeper) TransferBadge(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, amount uint64, badge_id uint64, subasset_id uint64, forceful_transfer bool) error {
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
	to_balance_id := GetFullSubassetID(to_account.GetAccountNumber(), badge_id, subasset_id)

	from_account := k.accountKeeper.GetAccount(ctx, from)
	if from_account == nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "account %s does not exist", from)
	}
	from_balance_id := GetFullSubassetID(from_account.GetAccountNumber(), badge_id, subasset_id)

	if (forceful_transfer) {
		k.AddToBadgeBalance(ctx, to_balance_id, amount)
		k.RemoveFromBadgeBalance(ctx, from_balance_id, amount)
	} else {
		//TODO: handle pending transfers
	}

	return nil;
}

//Transfer

//Accept / Reject / Cancel

//Approve / Update Approval

//TODO: GarbageCollect (may do this natively if everything is null)
func (k Keeper) GarbageCollectAddressIfPossible(ctx sdk.Context, address string) error {
	
	return nil
}
