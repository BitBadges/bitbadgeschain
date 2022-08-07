package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k Keeper) SetBadgeBalanceInStore(ctx sdk.Context, balance_id string, badgeBalanceInfo types.BadgeBalanceInfo) error {
	currTime := uint64(ctx.BlockTime().Unix())

	//If you have received a pending request that is expired, first write to this store no matter who it is from can just delete it
	//Or, if you sent a transfer request (which has no funds in escrow), you can just delete it
	prunedPending := make([]*types.PendingTransfer, 0)
	prunedApprovals := make([]*types.Approval, 0)

	details := GetDetailsFromBalanceKey(balance_id)
	thisAccountNum := details.account_num
	for _, pendingTransfer := range badgeBalanceInfo.Pending {
		if pendingTransfer.ExpirationTime != 0 && pendingTransfer.ExpirationTime < currTime && !pendingTransfer.SendRequest {
			continue
		} else if pendingTransfer.ExpirationTime != 0 && pendingTransfer.ExpirationTime < currTime && pendingTransfer.SendRequest && pendingTransfer.From == thisAccountNum {
			continue
		} else {
			prunedPending = append(prunedPending, pendingTransfer)
		}
	}
	badgeBalanceInfo.Pending = prunedPending

	//Remove any approvals that are expired
	for _, approval := range badgeBalanceInfo.Approvals {
		if approval.ExpirationTime != 0 && approval.ExpirationTime < currTime {
			continue
		} else {
			prunedApprovals = append(prunedApprovals, approval)
		}
	}
	badgeBalanceInfo.Approvals = prunedApprovals

	marshaled_badge_balance_info, err := k.cdc.Marshal(&badgeBalanceInfo)
	if err != nil {
		return sdkerrors.Wrap(err, "Marshal types.BadgeBalanceInfo failed")
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(badgeBalanceStoreKey(balance_id), marshaled_badge_balance_info)
	return nil
}

func (k Keeper) GetBadgeBalanceFromStore(ctx sdk.Context, balance_id string) (types.BadgeBalanceInfo, bool) {
	store := ctx.KVStore(k.storeKey)
	marshaled_badge_balance_info := store.Get(badgeBalanceStoreKey(balance_id))

	var badgeBalanceInfo types.BadgeBalanceInfo
	if len(marshaled_badge_balance_info) == 0 {
		return badgeBalanceInfo, false
	}
	k.cdc.MustUnmarshal(marshaled_badge_balance_info, &badgeBalanceInfo)
	return badgeBalanceInfo, true
}

func (k Keeper) GetBadgeBalancesFromStore(ctx sdk.Context) (addresses []*types.BadgeBalanceInfo) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, BadgeBalanceKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var badgeBalanceInfo types.BadgeBalanceInfo
		k.cdc.MustUnmarshal(iterator.Value(), &badgeBalanceInfo)
		addresses = append(addresses, &badgeBalanceInfo)
	}
	return
}

func (k Keeper) GetBadgeBalanceIdsFromStore(ctx sdk.Context) (ids []string) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, BadgeBalanceKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		ids = append(ids, string(iterator.Value()))
	}
	return
}

// func (k Keeper) StoreHasBadgeBalance(ctx sdk.Context, balance_id string) bool {
// 	store := ctx.KVStore(k.storeKey)
// 	return store.Has(badgeBalanceStoreKey(balance_id))
// }

// HasBadge determines whether the specified badgeID exists
func (k Keeper) DeleteBadgeBalanceFromStore(ctx sdk.Context, balance_id string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(badgeBalanceStoreKey(balance_id))
}