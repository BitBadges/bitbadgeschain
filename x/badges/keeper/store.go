package keeper

import (
	"strconv"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

/****************************************USER BALANCES****************************************/
// Sets a user balance in the store using UserBalanceKey ([]byte{0x02}) as the prefix. No check if store has key already.
func (k Keeper) SetUserBalanceInStore(ctx sdk.Context, balanceKey string, userBalanceInfo types.UserBalanceInfo) error {
	currTime := uint64(ctx.BlockTime().Unix())
	userBalanceInfo.Pending = PruneExpiredPending(currTime, GetDetailsFromBalanceKey(balanceKey).accountNum, userBalanceInfo.Pending)

	marshaled_badge_balance_info, err := k.cdc.Marshal(&userBalanceInfo)
	if err != nil {
		return sdkerrors.Wrap(err, "Marshal types.UserBalanceInfo failed")
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(userBalanceStoreKey(balanceKey), marshaled_badge_balance_info)
	return nil
}

// Gets a user balance from the store according to the balanceID.
func (k Keeper) GetUserBalanceFromStore(ctx sdk.Context, balanceKey string) (types.UserBalanceInfo, bool) {
	store := ctx.KVStore(k.storeKey)
	marshaled_badge_balance_info := store.Get(userBalanceStoreKey(balanceKey))

	var userBalanceInfo types.UserBalanceInfo
	if len(marshaled_badge_balance_info) == 0 {
		return userBalanceInfo, false
	}
	k.cdc.MustUnmarshal(marshaled_badge_balance_info, &userBalanceInfo)
	return userBalanceInfo, true
}

// GetUserBalancesFromStore defines a method for returning all user balances information by key.
func (k Keeper) GetUserBalancesFromStore(ctx sdk.Context) (addresses []*types.UserBalanceInfo) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, UserBalanceKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var userBalanceInfo types.UserBalanceInfo
		k.cdc.MustUnmarshal(iterator.Value(), &userBalanceInfo)
		addresses = append(addresses, &userBalanceInfo)
	}
	return
}

// GetUserBalanceIdsFromStore defines a method for returning all keys of all user balances.
func (k Keeper) GetUserBalanceIdsFromStore(ctx sdk.Context) (ids []string) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, UserBalanceKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		ids = append(ids, string(iterator.Value()))
	}
	return
}

// StoreHasUserBalanceID determines whether the specified user balanceID exists in the store
func (k Keeper) StoreHasUserBalance(ctx sdk.Context, balanceKey string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(userBalanceStoreKey(balanceKey))
}

// DeleteUserBalanceFromStore deletes a user balance from the store.
func (k Keeper) DeleteUserBalanceFromStore(ctx sdk.Context, balanceKey string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(userBalanceStoreKey(balanceKey))
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
	k.SetNextBadgeId(ctx, nextID+1) //susceptible to overflow but by that time we will have 2^64 badges which isn't totally feasible
}
