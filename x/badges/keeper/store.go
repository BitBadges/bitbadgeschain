package keeper

import (
	"strconv"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkerrors "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkmath "cosmossdk.io/math"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
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

	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	store.Set(collectionStoreKey(collection.CollectionId), marshaled_badge)
	return nil
}

// Gets a badge from the store according to the collectionId.
func (k Keeper) GetCollectionFromStore(ctx sdk.Context, collectionId sdkmath.Uint) (*types.BadgeCollection, bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
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
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	iterator := storetypes.KVStorePrefixIterator(store, CollectionKey)
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
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	return store.Has(collectionStoreKey(collectionId))
}

// DeleteCollectionFromStore deletes a badge from the store.
func (k Keeper) DeleteCollectionFromStore(ctx sdk.Context, collectionId sdkmath.Uint) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	store.Delete(collectionStoreKey(collectionId))
}

/****************************************USER BALANCES****************************************/

// Sets a user balance in the store using UserBalanceKey ([]byte{0x02}) as the prefix. No check if store has key already.
func (k Keeper) SetUserBalanceInStore(ctx sdk.Context, balanceKey string, UserBalance *types.UserBalanceStore) error {
	//HACK: We always store a non-nil permissions object to avoid the case where everything is nil -> marshaled len = 0 -> default balances get populated again
	if UserBalance.UserPermissions == nil {
		UserBalance.UserPermissions = &types.UserPermissions{
			CanUpdateOutgoingApprovals:                         []*types.UserOutgoingApprovalPermission{},
			CanUpdateIncomingApprovals:                         []*types.UserIncomingApprovalPermission{},
			CanUpdateAutoApproveSelfInitiatedOutgoingTransfers: []*types.ActionPermission{},
			CanUpdateAutoApproveSelfInitiatedIncomingTransfers: []*types.ActionPermission{},
			CanUpdateAutoApproveAllIncomingTransfers:           []*types.ActionPermission{},
		}
	}

	marshaled_badge_balance_info, err := k.cdc.Marshal(UserBalance)
	if err != nil {
		return sdkerrors.Wrap(err, "Marshal types.UserBalanceStore failed")
	}

	//Prevent accidental non-BitBadges addresses from being stored
	if GetDetailsFromBalanceKey(balanceKey).address != "Mint" && GetDetailsFromBalanceKey(balanceKey).address != "Total" {
		if err = types.ValidateAddress(GetDetailsFromBalanceKey(balanceKey).address, false); err != nil {
			return sdkerrors.Wrap(err, "Invalid address")
		}
	}

	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	store.Set(userBalanceStoreKey(balanceKey), marshaled_badge_balance_info)
	return nil
}

// Gets a user balance from the store according to the balanceID.
func (k Keeper) GetUserBalanceFromStore(ctx sdk.Context, balanceKey string) (*types.UserBalanceStore, bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
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
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	iterator := storetypes.KVStorePrefixIterator(store, UserBalanceKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var UserBalance types.UserBalanceStore
		k.cdc.MustUnmarshal(iterator.Value(), &UserBalance)
		balances = append(balances, &UserBalance)

		balanceKeyDetails := GetDetailsFromBalanceKey(string(iterator.Key()[1:]))
		ids = append(ids, balanceKeyDetails.collectionId)
		addresses = append(addresses, balanceKeyDetails.address)
	}
	return
}

// GetUserBalanceIdsFromStore defines a method for returning all keys of all user balances.
func (k Keeper) GetUserBalanceIdsFromStore(ctx sdk.Context) (ids []string) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})

	iterator := storetypes.KVStorePrefixIterator(store, UserBalanceKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		ids = append(ids, string(iterator.Key()[1:]))
	}
	return
}

// StoreHasUserBalanceID determines whether the specified user balanceID exists in the store
func (k Keeper) StoreHasUserBalance(ctx sdk.Context, balanceKey string) bool {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	return store.Has(userBalanceStoreKey(balanceKey))
}

// DeleteUserBalanceFromStore deletes a user balance from the store.
func (k Keeper) DeleteUserBalanceFromStore(ctx sdk.Context, balanceKey string) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	store.Delete(userBalanceStoreKey(balanceKey))
}

/****************************************NEXT COLLECTION ID****************************************/

// Gets the next collection ID.
func (k Keeper) GetNextCollectionId(ctx sdk.Context) sdkmath.Uint {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	nextCollectionId := store.Get(nextCollectionIdKey())
	nextCollectionIdStr := string((nextCollectionId))
	nextID := types.NewUintFromString(nextCollectionIdStr)
	return nextID
}

// Sets the next asset ID. Should only be used in InitGenesis. Everything else should call IncrementNextAssetID()
func (k Keeper) SetNextCollectionId(ctx sdk.Context, nextID sdkmath.Uint) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	store.Set(nextCollectionIdKey(), []byte(nextID.String()))
}

// Increments the next collection ID by 1.
func (k Keeper) IncrementNextCollectionId(ctx sdk.Context) {
	nextID := k.GetNextCollectionId(ctx)
	k.SetNextCollectionId(ctx, nextID.AddUint64(1)) //susceptible to overflow but by that time we will have 2^64 badges which isn't totally feasible
}

/****************************************NEXT LIST ID****************************************/

// Gets the next collection ID.
func (k Keeper) GetNextAddressListCounter(ctx sdk.Context) sdkmath.Uint {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	nextID := types.NewUintFromString(string((store.Get(nextAddressListCounterKey()))))
	return nextID
}

// Sets the next asset ID. Should only be used in InitGenesis. Everything else should call IncrementNextAssetID()
func (k Keeper) SetNextAddressListCounter(ctx sdk.Context, nextID sdkmath.Uint) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	store.Set(nextAddressListCounterKey(), []byte(nextID.String()))
}

// Increments the next collection ID by 1.
func (k Keeper) IncrementNextAddressListCounter(ctx sdk.Context) {
	nextID := k.GetNextAddressListCounter(ctx)
	k.SetNextAddressListCounter(ctx, nextID.AddUint64(1)) //susceptible to overflow but by that time we will have 2^64 badges which isn't totally feasible
}

/********************************************************************************/
// Sets a usedClaimData in the store using UsedClaimDataKey ([]byte{0x07}) as the prefix. No check if store has key already.
func (k Keeper) IncrementChallengeTrackerInStore(ctx sdk.Context, collectionId sdkmath.Uint, addressForChallenge string, approvalLevel string, approvalId, challengeId string, leafIndex sdkmath.Uint) (sdkmath.Uint, error) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	currBytes := store.Get(usedClaimChallengeStoreKey(ConstructUsedClaimChallengeKey(collectionId, addressForChallenge, approvalLevel, approvalId, challengeId, leafIndex)))
	curr := sdkmath.NewUint(0)
	if currBytes != nil {
		currUint, err := strconv.ParseUint(string((currBytes)), 10, 64)
		if err != nil {
			panic("Failed to parse num used")
		}

		curr = sdkmath.NewUint(currUint)
	}
	incrementedNum := curr.AddUint64(1)
	store.Set(usedClaimChallengeStoreKey(ConstructUsedClaimChallengeKey(collectionId, addressForChallenge, approvalLevel, approvalId, challengeId, leafIndex)), []byte(curr.Incr().String()))
	return incrementedNum, nil
}

func (k Keeper) GetChallengeTrackerFromStore(ctx sdk.Context, collectionId sdkmath.Uint, addressForChallenge string, approvalLevel string, approvalId, challengeId string, leafIndex sdkmath.Uint) (sdkmath.Uint, error) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	currBytes := store.Get(usedClaimChallengeStoreKey(ConstructUsedClaimChallengeKey(collectionId, addressForChallenge, approvalLevel, approvalId, challengeId, leafIndex)))
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

func (k Keeper) GetChallengeTrackersFromStore(ctx sdk.Context) (numUsed []sdkmath.Uint, ids []string) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	iterator := storetypes.KVStorePrefixIterator(store, UsedClaimChallengeKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		curr, err := strconv.ParseUint(string((iterator.Value())), 10, 64)
		if err != nil {
			panic("Failed to parse num used")
		}
		numUsed = append(numUsed, sdkmath.NewUint(curr))
		ids = append(ids, string(iterator.Key()[1:]))
	}
	return
}

func (k Keeper) SetChallengeTrackerInStore(ctx sdk.Context, key string, numUsed sdkmath.Uint) error {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	store.Set(usedClaimChallengeStoreKey(key), []byte(numUsed.String()))
	return nil
}

/****************************************ADDRESS LISTS****************************************/

func (k Keeper) SetAddressListInStore(ctx sdk.Context, addressList types.AddressList) error {
	marshaled_address_list, err := k.cdc.Marshal(&addressList)
	if err != nil {
		return sdkerrors.Wrap(err, "Marshal types.AddressList failed")
	}

	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	store.Set(addressListStoreKey(addressList.ListId), marshaled_address_list)
	return nil
}

func (k Keeper) GetAddressListFromStore(ctx sdk.Context, addressListId string) (types.AddressList, bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	marshaled_address_list := store.Get(addressListStoreKey(addressListId))

	var addressList types.AddressList
	if len(marshaled_address_list) == 0 {
		return addressList, false
	}
	k.cdc.MustUnmarshal(marshaled_address_list, &addressList)
	return addressList, true
}

func (k Keeper) GetAddressListsFromStore(ctx sdk.Context) (addressLists []*types.AddressList) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	iterator := storetypes.KVStorePrefixIterator(store, AddressListKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var addressList types.AddressList
		k.cdc.MustUnmarshal(iterator.Value(), &addressList)
		addressLists = append(addressLists, &addressList)
	}
	return
}

func (k Keeper) StoreHasAddressList(ctx sdk.Context, addressListId string) bool {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	return store.Has(addressListStoreKey(addressListId))
}

func (k Keeper) DeleteAddressListFromStore(ctx sdk.Context, addressListId string) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	store.Delete(addressListStoreKey(addressListId))
}

/****************************************TRANSFER TRACKERS****************************************/

func (k Keeper) SetApprovalTrackerInStore(ctx sdk.Context, collectionId sdkmath.Uint, addressForApproval string, approvalId, amountTrackerId string, approvalTracker types.ApprovalTracker, level string, trackerType string, address string) error {
	marshaled_transfer_tracker, err := k.cdc.Marshal(&approvalTracker)
	if err != nil {
		return sdkerrors.Wrap(err, "Marshal types.ApprovalTracker failed")
	}

	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	store.Set(approvalTrackerStoreKey(ConstructApprovalTrackerKey(collectionId, addressForApproval, approvalId, amountTrackerId, level, trackerType, address)), marshaled_transfer_tracker)
	return nil
}

func (k Keeper) SetApprovalTrackerInStoreViaKey(ctx sdk.Context, key string, approvalTracker types.ApprovalTracker) error {
	marshaled_transfer_tracker, err := k.cdc.Marshal(&approvalTracker)
	if err != nil {
		return sdkerrors.Wrap(err, "Marshal types.ApprovalTracker failed")
	}

	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	store.Set(approvalTrackerStoreKey(key), marshaled_transfer_tracker)
	return nil
}

func (k Keeper) GetApprovalTrackerFromStore(ctx sdk.Context, collectionId sdkmath.Uint, addressForApproval string, approvalId string, amountTrackerId string, level string, trackerType string, address string) (types.ApprovalTracker, bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	marshaled_transfer_tracker := store.Get(approvalTrackerStoreKey(ConstructApprovalTrackerKey(collectionId, addressForApproval, approvalId, amountTrackerId, level, trackerType, address)))

	var approvalTracker types.ApprovalTracker
	if len(marshaled_transfer_tracker) == 0 {
		return approvalTracker, false
	}
	k.cdc.MustUnmarshal(marshaled_transfer_tracker, &approvalTracker)
	return approvalTracker, true
}

func (k Keeper) GetApprovalTrackersFromStore(ctx sdk.Context) (approvalTrackers []*types.ApprovalTracker, ids []string) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	iterator := storetypes.KVStorePrefixIterator(store, ApprovalTrackerKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var approvalTracker types.ApprovalTracker
		k.cdc.MustUnmarshal(iterator.Value(), &approvalTracker)
		approvalTrackers = append(approvalTrackers, &approvalTracker)

		ids = append(ids, string(iterator.Key()[1:]))
	}
	return
}

func (k Keeper) StoreHasApprovalTracker(ctx sdk.Context, collectionId sdkmath.Uint, addressForApproval string, approvalId, amountTrackerId string, level string, trackerType string, address string) bool {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	return store.Has(approvalTrackerStoreKey(ConstructApprovalTrackerKey(collectionId, addressForApproval, approvalId, amountTrackerId, level, trackerType, address)))
}

func (k Keeper) DeleteApprovalTrackerFromStore(ctx sdk.Context, collectionId sdkmath.Uint, addressForApproval string, approvalId, amountTrackerId string, level string, trackerType string, address string) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	store.Delete(approvalTrackerStoreKey(ConstructApprovalTrackerKey(collectionId, addressForApproval, approvalId, amountTrackerId, level, trackerType, address)))
}

/** -------------------------------------- VERSION TRACKERS FOR APPROVAL IDS -------------------------------------- */

func (k Keeper) IncrementApprovalVersion(ctx sdk.Context, collectionId sdkmath.Uint, approvalLevel string, approverAddress string, approvalId string) sdkmath.Uint {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	version := store.Get(approvalVersionStoreKey(ConstructApprovalVersionKey(collectionId, approvalLevel, approverAddress, approvalId)))
	if version == nil {
		store.Set(approvalVersionStoreKey(ConstructApprovalVersionKey(collectionId, approvalLevel, approverAddress, approvalId)), []byte("0"))
		return sdkmath.NewUint(0)
	} else {
		versionUint, err := strconv.ParseUint(string(version), 10, 64)
		if err != nil {
			panic("Failed to parse version")
		}

		versionUint++
		newVersion := sdkmath.NewUint(versionUint)
		store.Set(approvalVersionStoreKey(ConstructApprovalVersionKey(collectionId, approvalLevel, approverAddress, approvalId)), []byte(newVersion.String()))
		return newVersion
	}
}

func (k Keeper) ResetApprovalVersion(ctx sdk.Context, collectionId sdkmath.Uint, approvalLevel string, approverAddress string, approvalId string) sdkmath.Uint {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	store.Set(approvalVersionStoreKey(ConstructApprovalVersionKey(collectionId, approvalLevel, approverAddress, approvalId)), []byte("0"))
	return sdkmath.NewUint(0)
}

func (k Keeper) GetApprovalTrackerVersionsFromStore(ctx sdk.Context) (versions []sdkmath.Uint, ids []string) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	iterator := storetypes.KVStorePrefixIterator(store, ApprovalVersionKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		version, err := strconv.ParseUint(string(iterator.Value()), 10, 64)
		if err != nil {
			panic("Failed to parse version")
		}
		versions = append(versions, sdkmath.NewUint(version))
		ids = append(ids, string(iterator.Key()[1:]))
	}
	return
}

func (k Keeper) SetApprovalTrackerVersionInStore(ctx sdk.Context, key string, version sdkmath.Uint) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	store.Set(approvalVersionStoreKey(key), []byte(version.String()))
}

func (k Keeper) GetApprovalTrackerVersionFromStore(ctx sdk.Context, key string) (sdkmath.Uint, bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	version := store.Get(approvalVersionStoreKey(key))
	if version == nil {
		return sdkmath.NewUint(0), false
	}
	versionUint, err := strconv.ParseUint(string(version), 10, 64)
	if err != nil {
		panic("Failed to parse version")
	}
	return sdkmath.NewUint(versionUint), true
}

/****************************************DYNAMIC STORES****************************************/

func (k Keeper) SetDynamicStoreInStore(ctx sdk.Context, dynamicStore types.DynamicStore) error {
	marshaled_dynamic_store, err := k.cdc.Marshal(&dynamicStore)
	if err != nil {
		return sdkerrors.Wrap(err, "Marshal types.DynamicStore failed")
	}

	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	store.Set(dynamicStoreStoreKey(dynamicStore.StoreId), marshaled_dynamic_store)
	return nil
}

func (k Keeper) GetDynamicStoreFromStore(ctx sdk.Context, storeId sdkmath.Uint) (types.DynamicStore, bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	marshaled_dynamic_store := store.Get(dynamicStoreStoreKey(storeId))

	var dynamicStore types.DynamicStore
	if len(marshaled_dynamic_store) == 0 {
		return dynamicStore, false
	}
	k.cdc.MustUnmarshal(marshaled_dynamic_store, &dynamicStore)
	return dynamicStore, true
}

func (k Keeper) GetDynamicStoresFromStore(ctx sdk.Context) (dynamicStores []*types.DynamicStore) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	iterator := storetypes.KVStorePrefixIterator(store, DynamicStoreKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var dynamicStore types.DynamicStore
		k.cdc.MustUnmarshal(iterator.Value(), &dynamicStore)
		dynamicStores = append(dynamicStores, &dynamicStore)
	}
	return
}

func (k Keeper) StoreHasDynamicStore(ctx sdk.Context, storeId sdkmath.Uint) bool {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	return store.Has(dynamicStoreStoreKey(storeId))
}

func (k Keeper) DeleteDynamicStoreFromStore(ctx sdk.Context, storeId sdkmath.Uint) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	store.Delete(dynamicStoreStoreKey(storeId))
}

/****************************************NEXT DYNAMIC STORE ID****************************************/

func (k Keeper) GetNextDynamicStoreId(ctx sdk.Context) sdkmath.Uint {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	nextID := types.NewUintFromString(string((store.Get(nextDynamicStoreIdKey()))))
	return nextID
}

func (k Keeper) SetNextDynamicStoreId(ctx sdk.Context, nextID sdkmath.Uint) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	store.Set(nextDynamicStoreIdKey(), []byte(nextID.String()))
}

func (k Keeper) IncrementNextDynamicStoreId(ctx sdk.Context) {
	nextID := k.GetNextDynamicStoreId(ctx)
	k.SetNextDynamicStoreId(ctx, nextID.AddUint64(1))
}

/****************************************DYNAMIC STORE VALUES****************************************/

// Sets a dynamic store value in the store using DynamicStoreValueKey ([]byte{0x0F}) as the prefix.
func (k Keeper) SetDynamicStoreValueInStore(ctx sdk.Context, storeId sdkmath.Uint, address string, value bool) error {
	dynamicStoreValue := types.DynamicStoreValue{
		StoreId: storeId,
		Address: address,
		Value:   value,
	}

	marshaled_value, err := k.cdc.Marshal(&dynamicStoreValue)
	if err != nil {
		return sdkerrors.Wrap(err, "Marshal types.DynamicStoreValue failed")
	}

	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	store.Set(dynamicStoreValueStoreKey(storeId, address), marshaled_value)
	return nil
}

// Gets a dynamic store value from the store according to the storeId and address.
func (k Keeper) GetDynamicStoreValueFromStore(ctx sdk.Context, storeId sdkmath.Uint, address string) (types.DynamicStoreValue, bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	marshaled_value := store.Get(dynamicStoreValueStoreKey(storeId, address))

	var dynamicStoreValue types.DynamicStoreValue
	if len(marshaled_value) == 0 {
		return dynamicStoreValue, false
	}
	k.cdc.MustUnmarshal(marshaled_value, &dynamicStoreValue)
	return dynamicStoreValue, true
}

// GetDynamicStoreValuesFromStore defines a method for returning all dynamic store values for a given store.
func (k Keeper) GetDynamicStoreValuesFromStore(ctx sdk.Context, storeId sdkmath.Uint) (dynamicStoreValues []*types.DynamicStoreValue) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	iterator := storetypes.KVStorePrefixIterator(store, append(DynamicStoreValueKey, []byte(storeId.String())...))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var dynamicStoreValue types.DynamicStoreValue
		k.cdc.MustUnmarshal(iterator.Value(), &dynamicStoreValue)
		dynamicStoreValues = append(dynamicStoreValues, &dynamicStoreValue)
	}
	return
}

// StoreHasDynamicStoreValue determines whether the specified dynamic store value exists in the store
func (k Keeper) StoreHasDynamicStoreValue(ctx sdk.Context, storeId sdkmath.Uint, address string) bool {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	return store.Has(dynamicStoreValueStoreKey(storeId, address))
}

// DeleteDynamicStoreValueFromStore deletes a dynamic store value from the store.
func (k Keeper) DeleteDynamicStoreValueFromStore(ctx sdk.Context, storeId sdkmath.Uint, address string) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	store.Delete(dynamicStoreValueStoreKey(storeId, address))
}
