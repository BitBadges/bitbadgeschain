package keeper

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

// The following methods are used for the badge store and everything associated with badges.
// All permissions and checks must be handled before these functions are called.
// This file handles:
// - Storing badges in the store
// - Storing balances in the store
// - Storing transfer manager requests in the store
// - Storing the next asset ID in the store

// All the following CRUD operations must obey the key prefixes defined in keys.go.

/****************************************BADGES****************************************/

// Sets a badge in the store using BadgeKey ([]byte{0x01}) as the prefix. No check if store has key already.
func (k Keeper) SetBadgeInStore(ctx sdk.Context, badge types.BitBadge) error {
	marshaled_badge, err := k.cdc.Marshal(&badge)
	if err != nil {
		return sdkerrors.Wrap(err, "Marshal types.BitBadge failed")
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(badgeStoreKey(badge.Id), marshaled_badge)
	return nil
}

// Gets a badge from the store according to the badgeID.
func (k Keeper) GetBadgeFromStore(ctx sdk.Context, badgeID uint64) (types.BitBadge, bool) {
	store := ctx.KVStore(k.storeKey)
	marshaled_badge := store.Get(badgeStoreKey(badgeID))

	var badge types.BitBadge
	if len(marshaled_badge) == 0 {
		return badge, false
	}
	k.cdc.MustUnmarshal(marshaled_badge, &badge)
	return badge, true
}

// GetBadgesFromStore defines a method for returning all badges information by key. 
func (k Keeper) GetBadgesFromStore(ctx sdk.Context) (badges []*types.BitBadge) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, BadgeKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var badge types.BitBadge
		k.cdc.MustUnmarshal(iterator.Value(), &badge)
		badges = append(badges, &badge)
	}
	return
}

// StoreHasBadgeID determines whether the specified badgeID exists
func (k Keeper) StoreHasBadgeID(ctx sdk.Context, badgeID uint64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(badgeStoreKey(badgeID))
}

// DeleteBadgeFromStore deletes a badge from the store.
func (k Keeper) DeleteBadgeFromStore(ctx sdk.Context, badgeID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(badgeStoreKey(badgeID))
}

/****************************************BADGE BALANCES****************************************/
// Sets a badge balance in the store using BadgeBalanceKey ([]byte{0x02}) as the prefix. No check if store has key already.
func (k Keeper) SetBadgeBalanceInStore(ctx sdk.Context, balanceKey string, badgeBalanceInfo types.BadgeBalanceInfo) error {
	currTime := uint64(ctx.BlockTime().Unix())
	//TODO: neaten up and comment
	//If you have received a pending request that is expired, first write to this store no matter who it is from can just delete it
	//Or, if you sent a transfer request (which has no funds in escrow), you can just delete it
	prunedPending := make([]*types.PendingTransfer, 0)
	prunedApprovals := make([]*types.Approval, 0)

	details := GetDetailsFromBalanceKey(balanceKey)
	thisAccountNum := details.accountNum
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
	store.Set(badgeBalanceStoreKey(balanceKey), marshaled_badge_balance_info)
	return nil
}

// Gets a badge balance from the store according to the balanceID.
func (k Keeper) GetBadgeBalanceFromStore(ctx sdk.Context, balanceKey string) (types.BadgeBalanceInfo, bool) {
	store := ctx.KVStore(k.storeKey)
	marshaled_badge_balance_info := store.Get(badgeBalanceStoreKey(balanceKey))

	var badgeBalanceInfo types.BadgeBalanceInfo
	if len(marshaled_badge_balance_info) == 0 {
		return badgeBalanceInfo, false
	}
	k.cdc.MustUnmarshal(marshaled_badge_balance_info, &badgeBalanceInfo)
	return badgeBalanceInfo, true
}

// GetBadgeBalancesFromStore defines a method for returning all badge balances information by key.
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

// GetBadgeBalanceIdsFromStore defines a method for returning all keys of all badge balances.
func (k Keeper) GetBadgeBalanceIdsFromStore(ctx sdk.Context) (ids []string) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, BadgeBalanceKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		ids = append(ids, string(iterator.Value()))
	}
	return
}

// StoreHasBadgeBalanceID determines whether the specified badge balanceID exists in the store
func (k Keeper) StoreHasBadgeBalance(ctx sdk.Context, balanceKey string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(badgeBalanceStoreKey(balanceKey))
}

// DeleteBadgeBalanceFromStore deletes a badge balance from the store.
func (k Keeper) DeleteBadgeBalanceFromStore(ctx sdk.Context, balanceKey string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(badgeBalanceStoreKey(balanceKey))
}

/****************************************TRANSFER MANAGER REQUESTS****************************************/

// Checks if a certain address has requested a managerial transfer
func (k Keeper) HasAddressRequestedManagerTransfer(ctx sdk.Context, badgeId uint64, address uint64) bool {
	store := ctx.KVStore(k.storeKey)
	key := ConstructTransferManagerRequestKey(badgeId, address)
	return store.Has(managerTransferRequestKey(key))
}

// Creates a transfer manager request for the given address and badgeID.
func (k Keeper) CreateTransferManagerRequest(ctx sdk.Context, badgeId uint64, address uint64) error {
	request := []byte{}
	store := ctx.KVStore(k.storeKey)
	key := ConstructTransferManagerRequestKey(badgeId, address)
	store.Set(managerTransferRequestKey(key), request)
	return nil
}

// Deletes a transfer manager request for the given address and badgeID.
func (k Keeper) RemoveTransferManagerRequest(ctx sdk.Context, badgeId uint64, address uint64) error {
	key := ConstructTransferManagerRequestKey(badgeId, address)
	store := ctx.KVStore(k.storeKey)

	store.Delete(managerTransferRequestKey(key))
	return nil
}

/****************************************NEXT ASSET ID****************************************/

//Gets the next badge ID.
func (k Keeper) GetNextBadgeId(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	nextID, err := strconv.ParseUint(string((store.Get(nextAssetIDKey()))), 10, 64)
	if err != nil {
		panic("Failed to parse next asset ID")
	}
	return nextID
}

// Sets the next asset ID. Should only be used in InitGenesis. Everything else should call IncrementNextAssetID()
func (k Keeper) SetNextBadgeId(ctx sdk.Context, nextID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(nextAssetIDKey(), []byte(strconv.FormatInt(int64(nextID), 10)))
}

// Increments the next badge ID by 1. 
func (k Keeper) IncrementNextBadgeId(ctx sdk.Context) {
	nextID := k.GetNextBadgeId(ctx)
	k.SetNextBadgeId(ctx, nextID + 1)
}
