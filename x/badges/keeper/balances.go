package keeper

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

//TODO: overflow / underflow handlers

func GetFullSubassetID(accountNumber uint64, id uint64, subasset_id uint64) string {
	badge_id_str := strconv.FormatUint(id, 10)
	subasset_id_str := strconv.FormatUint(subasset_id, 10)
	account_num_str := strconv.FormatUint(accountNumber, 10)
	return account_num_str + "-" + badge_id_str + "-" + subasset_id_str
}


func CreateNewBadgeBalanceIfEmpty(k Keeper, ctx sdk.Context, balance_id string, balance_to_add uint64) error {
	if balance_to_add == 0 {
		return ErrBalanceIsZero
	}

	if !k.StoreHasBadgeBalance(ctx, balance_id) {
		badgeBalanceInfo := types.BadgeBalanceInfo{
			Balance: balance_to_add,
			PendingNonce: 0,
			Pending: []*types.PendingTransfer{},
			Approvals: []*types.Approval{},
		}
		if err := k.CreateBadgeBalanceInStore(ctx, balance_id, badgeBalanceInfo); err != nil {
			return err
		}
	}

	return nil
}

//assumes already in store
func (k Keeper) AddToBadgeBalance(ctx sdk.Context, balance_id string, balance_to_add uint64) error {
	if balance_to_add == 0 {
		return ErrBalanceIsZero
	}

	badgeBalanceInfo, found := k.GetBadgeBalanceFromStore(ctx, balance_id)
	if !found {
		err := CreateNewBadgeBalanceIfEmpty(k, ctx, balance_id, balance_to_add); if err != nil {
			return err
		}
	} else {
		badgeBalanceInfo.Balance += balance_to_add
		err := k.UpdateBadgeBalanceInStore(ctx, balance_id, badgeBalanceInfo); if err != nil {
			return err
		}
	}

	return nil
}

//MintBadgeToManager

//Transfer

//Accept / Reject

//Approve / Update Approval

//TODO: GarbageCollect
func (k Keeper) GarbageCollectAddressIfPossible(ctx sdk.Context, address string) error {
	
	return nil
}
