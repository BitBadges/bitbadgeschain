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
// - Claims

// All the following CRUD operations must obey the key prefixes defined in keys.go.

/****************************************BADGES****************************************/

// Sets a badge in the store using BadgeKey ([]byte{0x01}) as the prefix. No check if store has key already.
func (k Keeper) SetCollectionInStore(ctx sdk.Context, collection types.BadgeCollection) error {
	marshaled_badge, err := k.cdc.Marshal(&collection)
	if err != nil {
		return sdkerrors.Wrap(err, "Marshal types.BadgeCollection failed")
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(collectionStoreKey(collection.CollectionId), marshaled_badge)
	return nil
}

// Gets a badge from the store according to the collectionId.
func (k Keeper) GetCollectionFromStore(ctx sdk.Context, collectionId uint64) (types.BadgeCollection, bool) {
	store := ctx.KVStore(k.storeKey)
	marshaled_collection := store.Get(collectionStoreKey(collectionId))

	var collection types.BadgeCollection
	if len(marshaled_collection) == 0 {
		return collection, false
	}
	k.cdc.MustUnmarshal(marshaled_collection, &collection)
	return collection, true
}

// GetCollectionsFromStore defines a method for returning all badges information by key.
func (k Keeper) GetCollectionsFromStore(ctx sdk.Context) (collections []*types.BadgeCollection) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, CollectionKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var collection types.BadgeCollection
		k.cdc.MustUnmarshal(iterator.Value(), &collection)
		collections = append(collections, &collection)
	}
	return
}

// StoreHasCollectionID determines whether the specified collectionId exists
func (k Keeper) StoreHasCollectionID(ctx sdk.Context, collectionId uint64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(collectionStoreKey(collectionId))
}

// DeleteCollectionFromStore deletes a badge from the store.
func (k Keeper) DeleteCollectionFromStore(ctx sdk.Context, collectionId uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(collectionStoreKey(collectionId))
}

/****************************************USER BALANCES****************************************/

// Sets a user balance in the store using UserBalanceKey ([]byte{0x02}) as the prefix. No check if store has key already.
func (k Keeper) SetUserBalanceInStore(ctx sdk.Context, balanceKey string, UserBalance types.UserBalance) error {
	marshaled_badge_balance_info, err := k.cdc.Marshal(&UserBalance)
	if err != nil {
		return sdkerrors.Wrap(err, "Marshal types.UserBalance failed")
	}

	

	store := ctx.KVStore(k.storeKey)
	store.Set(userBalanceStoreKey(balanceKey), marshaled_badge_balance_info)
	return nil
}

// Gets a user balance from the store according to the balanceID.
func (k Keeper) GetUserBalanceFromStore(ctx sdk.Context, balanceKey string) (types.UserBalance, bool) {
	store := ctx.KVStore(k.storeKey)
	marshaled_badge_balance_info := store.Get(userBalanceStoreKey(balanceKey))

	var UserBalance types.UserBalance
	if len(marshaled_badge_balance_info) == 0 {
		return UserBalance, false
	}
	k.cdc.MustUnmarshal(marshaled_badge_balance_info, &UserBalance)
	return UserBalance, true
}

// GetUserBalancesFromStore defines a method for returning all user balances information by key.
func (k Keeper) GetUserBalancesFromStore(ctx sdk.Context) (balances []*types.UserBalance, accNums []uint64, ids []uint64) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, UserBalanceKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var UserBalance types.UserBalance
		k.cdc.MustUnmarshal(iterator.Value(), &UserBalance)
		balances = append(balances, &UserBalance)

		
		balanceKeyDetails := GetDetailsFromBalanceKey(string(iterator.Key()))
		ids = append(ids, balanceKeyDetails.collectionId)
		accNums = append(accNums, balanceKeyDetails.accountNum)
	}
	return
}

// GetUserBalanceIdsFromStore defines a method for returning all keys of all user balances.
func (k Keeper) GetUserBalanceIdsFromStore(ctx sdk.Context) (ids []string) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, UserBalanceKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		ids = append(ids, string(iterator.Key()))
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
func (k Keeper) HasAddressRequestedManagerTransfer(ctx sdk.Context, collectionId uint64, address uint64) bool {
	store := ctx.KVStore(k.storeKey)
	key := ConstructTransferManagerRequestKey(collectionId, address)
	return store.Has(managerTransferRequestKey(key))
}

// Creates a transfer manager request for the given address and collectionId.
func (k Keeper) CreateTransferManagerRequest(ctx sdk.Context, collectionId uint64, address uint64) error {
	request := []byte{}
	store := ctx.KVStore(k.storeKey)
	key := ConstructTransferManagerRequestKey(collectionId, address)
	store.Set(managerTransferRequestKey(key), request)
	return nil
}

// Deletes a transfer manager request for the given address and collectionId.
func (k Keeper) RemoveTransferManagerRequest(ctx sdk.Context, collectionId uint64, address uint64) error {
	key := ConstructTransferManagerRequestKey(collectionId, address)
	store := ctx.KVStore(k.storeKey)

	store.Delete(managerTransferRequestKey(key))
	return nil
}

/****************************************NEXT ASSET ID****************************************/

//Gets the next badge ID.
func (k Keeper) GetNextCollectionId(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	nextID, err := strconv.ParseUint(string((store.Get(nextCollectionIdKey()))), 10, 64)
	if err != nil {
		panic("Failed to parse next asset ID")
	}
	return nextID
}

// Sets the next asset ID. Should only be used in InitGenesis. Everything else should call IncrementNextAssetID()
func (k Keeper) SetNextCollectionId(ctx sdk.Context, nextID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(nextCollectionIdKey(), []byte(strconv.FormatInt(int64(nextID), 10)))
}

// Increments the next badge ID by 1.
func (k Keeper) IncrementNextCollectionId(ctx sdk.Context) {
	nextID := k.GetNextCollectionId(ctx)
	k.SetNextCollectionId(ctx, nextID+1) //susceptible to overflow but by that time we will have 2^64 badges which isn't totally feasible
}

/****************************************NEXT CLAIM ID****************************************/

//Gets the next badge ID.
func (k Keeper) GetNextClaimId(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	nextID, err := strconv.ParseUint(string((store.Get(nextClaimIdKey()))), 10, 64)
	if err != nil {
		panic("Failed to parse next asset ID")
	}
	return nextID
}

// Sets the next asset ID. Should only be used in InitGenesis. Everything else should call IncrementNextAssetID()
func (k Keeper) SetNextClaimId(ctx sdk.Context, nextID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(nextClaimIdKey(), []byte(strconv.FormatInt(int64(nextID), 10)))
}

// Increments the next badge ID by 1.
func (k Keeper) IncrementNextClaimId(ctx sdk.Context) {
	nextID := k.GetNextClaimId(ctx)
	k.SetNextClaimId(ctx, nextID+1) //susceptible to overflow but by that time we will have 2^64 badges which isn't totally feasible
}

/****************************************Claims****************************************/
// Sets a usedClaimData in the store using UsedClaimDataKey ([]byte{0x07}) as the prefix. No check if store has key already.
func (k Keeper) IncrementNumUsedForClaimInStore(ctx sdk.Context, collectionId uint64, claimId uint64) (uint64, error) {
	store := ctx.KVStore(k.storeKey)
	currBytes := store.Get(usedClaimDataStoreKey(ConstructUsedClaimDataKey(collectionId, claimId)))
	curr := uint64(0)
	err := error(nil)
	if currBytes != nil {
		curr, err = strconv.ParseUint(string((currBytes)), 10, 64)
		if err != nil {
			panic("Failed to parse num used")
		}
	}
	incrementedNum := curr + 1
	store.Set(usedClaimDataStoreKey(ConstructUsedClaimDataKey(collectionId, claimId)), []byte(strconv.FormatInt(int64(curr+1), 10)))
	return incrementedNum, nil
}

func (k Keeper) IncrementNumUsedForCodeInStore(ctx sdk.Context, collectionId uint64, claimId uint64, codeLeafIndex uint64) (uint64, error) {
	store := ctx.KVStore(k.storeKey)
	currBytes := store.Get(usedClaimCodeStoreKey(ConstructUsedClaimCodeKey(collectionId, claimId, codeLeafIndex)))
	curr := uint64(0)
	err := error(nil)
	if currBytes != nil {
		curr, err = strconv.ParseUint(string((currBytes)), 10, 64)
		if err != nil {
			panic("Failed to parse num used")
		}
	}
	incrementedNum := curr + 1
	store.Set(usedClaimCodeStoreKey(ConstructUsedClaimCodeKey(collectionId, claimId, codeLeafIndex)), []byte(strconv.FormatInt(int64(curr+1), 10)))
	return incrementedNum, nil
}

func (k Keeper) IncrementNumUsedForAddressInStore(ctx sdk.Context, collectionId uint64, claimId uint64, address string) (uint64, error) {
	store := ctx.KVStore(k.storeKey)
	currBytes := store.Get(usedClaimAddressStoreKey(ConstructUsedClaimAddressKey(collectionId, claimId, address)))
	curr := uint64(0)
	err := error(nil)
	if currBytes != nil {
		curr, err = strconv.ParseUint(string((currBytes)), 10, 64)
		if err != nil {
			panic("Failed to parse num used")
		}
	}
	incrementedNum := curr + 1
	store.Set(usedClaimAddressStoreKey(ConstructUsedClaimAddressKey(collectionId, claimId, address)), []byte(strconv.FormatInt(int64(curr+1), 10)))
	return incrementedNum, nil
}



func (k Keeper) IncrementNumUsedForWhitelistIndexInStore(ctx sdk.Context, collectionId uint64, claimId uint64, whitelistLeafIndex uint64) (uint64, error) {
	store := ctx.KVStore(k.storeKey)
	currBytes := store.Get(usedWhitelistIndexStoreKey(ConstructUsedWhitelistIndexKey(collectionId, claimId, whitelistLeafIndex)))
	curr := uint64(0)
	err := error(nil)
	if currBytes != nil {
		curr, err = strconv.ParseUint(string((currBytes)), 10, 64)
		if err != nil {
			panic("Failed to parse num used")
		}
	}
	incrementedNum := curr + 1
	store.Set(usedWhitelistIndexStoreKey(ConstructUsedWhitelistIndexKey(collectionId, claimId, whitelistLeafIndex)), []byte(strconv.FormatInt(int64(curr+1), 10)))
	return incrementedNum, nil
}
