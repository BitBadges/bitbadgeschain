package keeper

import (
	"strconv"

	sdkerrors "cosmossdk.io/errors"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkmath "cosmossdk.io/math"
)

// The following methods are used for the badge store and everything associated with badges.
// All preconditions and checks must be handled before these functions are called.
// This file handles storing collections, balances, approvals, used challenges, next collection ID, etc.

// All the following CRUD operations must obey the key prefixes defined in keys.go.

/****************************************COLLECTIONS****************************************/

// Sets a badge in the store using BadgeKey ([]byte{0x01}) as the prefix. No check if store has key already.
func (k Keeper) SetCollectionInStore(ctx sdk.Context, collection *types.BadgeCollection) error {
	marshaled_badge, err := k.cdc.Marshal(collection)
	if err != nil {
		return sdkerrors.Wrap(err, "Marshal types.BadgeCollection failed")
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(collectionStoreKey(collection.CollectionId), marshaled_badge)
	return nil
}

// Gets a badge from the store according to the collectionId.
func (k Keeper) GetCollectionFromStore(ctx sdk.Context, collectionId sdkmath.Uint) (*types.BadgeCollection, bool) {
	store := ctx.KVStore(k.storeKey)
	marshaled_collection := store.Get(collectionStoreKey(collectionId))

	var collection types.BadgeCollection
	if len(marshaled_collection) == 0 {
		return &collection, false
	}
	k.cdc.MustUnmarshal(marshaled_collection, &collection)
	return &collection, true
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
func (k Keeper) StoreHasCollectionID(ctx sdk.Context, collectionId sdkmath.Uint) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(collectionStoreKey(collectionId))
}

// DeleteCollectionFromStore deletes a badge from the store.
func (k Keeper) DeleteCollectionFromStore(ctx sdk.Context, collectionId sdkmath.Uint) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(collectionStoreKey(collectionId))
}

/****************************************USER BALANCES****************************************/

// Sets a user balance in the store using UserBalanceKey ([]byte{0x02}) as the prefix. No check if store has key already.
func (k Keeper) SetUserBalanceInStore(ctx sdk.Context, balanceKey string, UserBalance *types.UserBalanceStore) error {
	marshaled_badge_balance_info, err := k.cdc.Marshal(UserBalance)
	if err != nil {
		return sdkerrors.Wrap(err, "Marshal types.UserBalanceStore failed")
	}

	//Prevent accidental non-cosmos addresses from being stored
	if err = types.ValidateAddress(GetDetailsFromBalanceKey(balanceKey).address, false); err != nil {
		return sdkerrors.Wrap(err, "Invalid address")
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(userBalanceStoreKey(balanceKey), marshaled_badge_balance_info)
	return nil
}

// Gets a user balance from the store according to the balanceID.
func (k Keeper) GetUserBalanceFromStore(ctx sdk.Context, balanceKey string) (*types.UserBalanceStore, bool) {
	store := ctx.KVStore(k.storeKey)
	marshaled_badge_balance_info := store.Get(userBalanceStoreKey(balanceKey))

	var UserBalance types.UserBalanceStore
	if len(marshaled_badge_balance_info) == 0 {
		return &UserBalance, false
	}
	k.cdc.MustUnmarshal(marshaled_badge_balance_info, &UserBalance)
	return &UserBalance, true
}

// GetUserBalancesFromStore defines a method for returning all user balances information by key.
func (k Keeper) GetUserBalancesFromStore(ctx sdk.Context) (balances []*types.UserBalanceStore, addresses []string, ids []sdkmath.Uint) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, UserBalanceKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var UserBalance types.UserBalanceStore
		k.cdc.MustUnmarshal(iterator.Value(), &UserBalance)
		balances = append(balances, &UserBalance)

		balanceKeyDetails := GetDetailsFromBalanceKey(string(iterator.Key()))
		ids = append(ids, balanceKeyDetails.collectionId)
		addresses = append(addresses, balanceKeyDetails.address)
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

/****************************************NEXT COLLECTION ID****************************************/

// Gets the next collection ID.
func (k Keeper) GetNextCollectionId(ctx sdk.Context) sdkmath.Uint {
	store := ctx.KVStore(k.storeKey)
	nextID := types.NewUintFromString(string((store.Get(nextCollectionIdKey()))))
	return nextID
}

// Sets the next asset ID. Should only be used in InitGenesis. Everything else should call IncrementNextAssetID()
func (k Keeper) SetNextCollectionId(ctx sdk.Context, nextID sdkmath.Uint) {
	store := ctx.KVStore(k.storeKey)
	store.Set(nextCollectionIdKey(), []byte(nextID.String()))
}

// Increments the next collection ID by 1.
func (k Keeper) IncrementNextCollectionId(ctx sdk.Context) {
	nextID := k.GetNextCollectionId(ctx)
	k.SetNextCollectionId(ctx, nextID.AddUint64(1)) //susceptible to overflow but by that time we will have 2^64 badges which isn't totally feasible
}

/********************************************************************************/
// Sets a usedClaimData in the store using UsedClaimDataKey ([]byte{0x07}) as the prefix. No check if store has key already.
func (k Keeper) IncrementNumUsedForChallengeInStore(ctx sdk.Context, collectionId sdkmath.Uint, challengeId string, leafIndex sdkmath.Uint, level string) (sdkmath.Uint, error) {
	store := ctx.KVStore(k.storeKey)
	currBytes := store.Get(usedClaimChallengeStoreKey(ConstructUsedClaimChallengeKey(collectionId, challengeId, leafIndex, level)))
	curr := sdkmath.NewUint(0)
	if currBytes != nil {
		currUint, err := strconv.ParseUint(string((currBytes)), 10, 64)
		if err != nil {
			panic("Failed to parse num used")
		}

		curr = sdkmath.NewUint(currUint)
	}
	incrementedNum := curr.AddUint64(1)
	store.Set(usedClaimChallengeStoreKey(ConstructUsedClaimChallengeKey(collectionId, challengeId, leafIndex, level)), []byte(curr.Incr().String()))
	return incrementedNum, nil
}

func (k Keeper) GetNumUsedForChallengeFromStore(ctx sdk.Context, collectionId sdkmath.Uint, challengeId string, leafIndex sdkmath.Uint, level string) (sdkmath.Uint, error) {
	store := ctx.KVStore(k.storeKey)
	currBytes := store.Get(usedClaimChallengeStoreKey(ConstructUsedClaimChallengeKey(collectionId, challengeId, leafIndex, level)))
	curr := sdkmath.NewUint(0)
	if currBytes != nil {
		currUint, err := strconv.ParseUint(string((currBytes)), 10, 64)
		if err != nil {
			panic("Failed to parse num used")
		}

		curr = sdkmath.NewUint(currUint)
	}
	return curr, nil
}

func (k Keeper) GetNumUsedForChallengesFromStore(ctx sdk.Context) (numUsed []sdkmath.Uint, ids []string) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, UsedClaimChallengeKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		curr, err := strconv.ParseUint(string((iterator.Value())), 10, 64)
		if err != nil {
			panic("Failed to parse num used")
		}
		numUsed = append(numUsed, sdkmath.NewUint(curr))
		ids = append(ids, string(iterator.Key()))
	}
	return
}

/****************************************ADDRESS MAPPINGS****************************************/

func (k Keeper) SetAddressMappingInStore(ctx sdk.Context, addressMapping types.AddressMapping) error {
	marshaled_address_mapping, err := k.cdc.Marshal(&addressMapping)
	if err != nil {
		return sdkerrors.Wrap(err, "Marshal types.AddressMapping failed")
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(addressMappingStoreKey(addressMapping.MappingId), marshaled_address_mapping)
	return nil
}

func (k Keeper) GetAddressMappingFromStoreFromStore(ctx sdk.Context, addressMappingId string) (types.AddressMapping, bool) {
	store := ctx.KVStore(k.storeKey)
	marshaled_address_mapping := store.Get(addressMappingStoreKey(addressMappingId))

	var addressMapping types.AddressMapping
	if len(marshaled_address_mapping) == 0 {
		return addressMapping, false
	}
	k.cdc.MustUnmarshal(marshaled_address_mapping, &addressMapping)
	return addressMapping, true
}

func (k Keeper) GetAddressMappingFromStoresFromStore(ctx sdk.Context) (addressMappings []*types.AddressMapping) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, AddressMappingKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var addressMapping types.AddressMapping
		k.cdc.MustUnmarshal(iterator.Value(), &addressMapping)
		addressMappings = append(addressMappings, &addressMapping)
	}
	return
}

func (k Keeper) StoreHasAddressMapping(ctx sdk.Context, addressMappingId string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(addressMappingStoreKey(addressMappingId))
}

func (k Keeper) DeleteAddressMappingFromStore(ctx sdk.Context, addressMappingId string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(addressMappingStoreKey(addressMappingId))
}

/****************************************TRANSFER TRACKERS****************************************/

func (k Keeper) SetTransferTrackerInStore(ctx sdk.Context, collectionId sdkmath.Uint, trackerId string, transferTracker types.ApprovalsTracker, level string, depth string, address string) error {
	marshaled_transfer_tracker, err := k.cdc.Marshal(&transferTracker)
	if err != nil {
		return sdkerrors.Wrap(err, "Marshal types.ApprovalsTracker failed")
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(transferTrackerStoreKey(ConstructTransferTrackerKey(collectionId, trackerId, level, depth, address)), marshaled_transfer_tracker)
	return nil
}

func (k Keeper) GetTransferTrackerFromStore(ctx sdk.Context, collectionId sdkmath.Uint, trackerId string, level string, depth string, address string) (types.ApprovalsTracker, bool) {
	store := ctx.KVStore(k.storeKey)
	marshaled_transfer_tracker := store.Get(transferTrackerStoreKey(ConstructTransferTrackerKey(collectionId, trackerId, level, depth, address)))

	var transferTracker types.ApprovalsTracker
	if len(marshaled_transfer_tracker) == 0 {
		return transferTracker, false
	}
	k.cdc.MustUnmarshal(marshaled_transfer_tracker, &transferTracker)
	return transferTracker, true
}

func (k Keeper) GetTransferTrackersFromStore(ctx sdk.Context) (transferTrackers []*types.ApprovalsTracker, ids []string) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, TransferTrackerKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var transferTracker types.ApprovalsTracker
		k.cdc.MustUnmarshal(iterator.Value(), &transferTracker)
		transferTrackers = append(transferTrackers, &transferTracker)

		ids = append(ids, string(iterator.Key()))
	}
	return
}

func (k Keeper) StoreHasTransferTracker(ctx sdk.Context, collectionId sdkmath.Uint, trackerId string, level string, depth string, address string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(transferTrackerStoreKey(ConstructTransferTrackerKey(collectionId, trackerId, level, depth, address)))
}

func (k Keeper) DeleteTransferTrackerFromStore(ctx sdk.Context, collectionId sdkmath.Uint, trackerId string, level string, depth string, address string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(transferTrackerStoreKey(ConstructTransferTrackerKey(collectionId, trackerId, level, depth, address)))
}
