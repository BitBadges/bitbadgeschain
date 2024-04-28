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

/** ****************************** GLOBAL ARCHIVE ****************************** **/
func (k Keeper) SetGlobalArchiveInStore(ctx sdk.Context, archive bool) error {
	store := ctx.KVStore(k.storeKey)
	store.Set(GlobalArchiveKey, []byte(strconv.FormatBool(archive)))
	return nil
}
func (k Keeper) GetGlobalArchiveFromStore(ctx sdk.Context) bool {
	store := ctx.KVStore(k.storeKey)
	archive := store.Get(GlobalArchiveKey)
	if archive == nil {
		return false
	}
	return archive[0] == 't'
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
		}
	}

	marshaled_badge_balance_info, err := k.cdc.Marshal(UserBalance)
	if err != nil {
		return sdkerrors.Wrap(err, "Marshal types.UserBalanceStore failed")
	}

	//Prevent accidental non-cosmos addresses from being stored
	if GetDetailsFromBalanceKey(balanceKey).address != "Mint" && GetDetailsFromBalanceKey(balanceKey).address != "Total" {
		if err = types.ValidateAddress(GetDetailsFromBalanceKey(balanceKey).address, false); err != nil {
			return sdkerrors.Wrap(err, "Invalid address")
		}
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

		balanceKeyDetails := GetDetailsFromBalanceKey(string(iterator.Key()[1:]))
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
		ids = append(ids, string(iterator.Key()[1:]))
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

/****************************************NEXT LIST ID****************************************/

// Gets the next collection ID.
func (k Keeper) GetNextAddressListCounter(ctx sdk.Context) sdkmath.Uint {
	store := ctx.KVStore(k.storeKey)
	nextID := types.NewUintFromString(string((store.Get(nextAddressListCounterKey()))))
	return nextID
}

// Sets the next asset ID. Should only be used in InitGenesis. Everything else should call IncrementNextAssetID()
func (k Keeper) SetNextAddressListCounter(ctx sdk.Context, nextID sdkmath.Uint) {
	store := ctx.KVStore(k.storeKey)
	store.Set(nextAddressListCounterKey(), []byte(nextID.String()))
}

// Increments the next collection ID by 1.
func (k Keeper) IncrementNextAddressListCounter(ctx sdk.Context) {
	nextID := k.GetNextAddressListCounter(ctx)
	k.SetNextAddressListCounter(ctx, nextID.AddUint64(1)) //susceptible to overflow but by that time we will have 2^64 badges which isn't totally feasible
}

/*********************************USED ZKPS*********************************/
func (k Keeper) SetZKPAsUsedInStore(ctx sdk.Context, collectionId sdkmath.Uint, addressForZKP string, approvalLevel string, approvalId, zkpId string, proofHash string) error {
	store := ctx.KVStore(k.storeKey)
	store.Set(usedZKPTrackerStoreKey(ConstructZKPTreeTrackerKey(collectionId, addressForZKP, approvalLevel, approvalId, zkpId, proofHash)), []byte("1"))
	return nil
}

func (k Keeper) GetZKPFromStore(ctx sdk.Context, collectionId sdkmath.Uint, addressForZKP string, approvalLevel string, approvalId, zkpId string, proofHash string) (bool, error) {
	store := ctx.KVStore(k.storeKey)
	return store.Has(usedZKPTrackerStoreKey(ConstructZKPTreeTrackerKey(collectionId, addressForZKP, approvalLevel, approvalId, zkpId, proofHash))), nil
}

/********************************************************************************/
// Sets a usedClaimData in the store using UsedClaimDataKey ([]byte{0x07}) as the prefix. No check if store has key already.
func (k Keeper) IncrementChallengeTrackerInStore(ctx sdk.Context, collectionId sdkmath.Uint, addressForChallenge string, approvalLevel string, approvalId, challengeId string, leafIndex sdkmath.Uint) (sdkmath.Uint, error) {
	store := ctx.KVStore(k.storeKey)
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
	store := ctx.KVStore(k.storeKey)
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
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, UsedClaimChallengeKey)
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
	store := ctx.KVStore(k.storeKey)
	store.Set(usedClaimChallengeStoreKey(key), []byte(numUsed.String()))
	return nil
}

/****************************************ADDRESS LISTS****************************************/

func (k Keeper) SetAddressListInStore(ctx sdk.Context, addressList types.AddressList) error {
	marshaled_address_list, err := k.cdc.Marshal(&addressList)
	if err != nil {
		return sdkerrors.Wrap(err, "Marshal types.AddressList failed")
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(addressListStoreKey(addressList.ListId), marshaled_address_list)
	return nil
}

func (k Keeper) GetAddressListFromStore(ctx sdk.Context, addressListId string) (types.AddressList, bool) {
	store := ctx.KVStore(k.storeKey)
	marshaled_address_list := store.Get(addressListStoreKey(addressListId))

	var addressList types.AddressList
	if len(marshaled_address_list) == 0 {
		return addressList, false
	}
	k.cdc.MustUnmarshal(marshaled_address_list, &addressList)
	return addressList, true
}

func (k Keeper) GetAddressListsFromStore(ctx sdk.Context) (addressLists []*types.AddressList) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, AddressListKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var addressList types.AddressList
		k.cdc.MustUnmarshal(iterator.Value(), &addressList)
		addressLists = append(addressLists, &addressList)
	}
	return
}

func (k Keeper) StoreHasAddressList(ctx sdk.Context, addressListId string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(addressListStoreKey(addressListId))
}

func (k Keeper) DeleteAddressListFromStore(ctx sdk.Context, addressListId string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(addressListStoreKey(addressListId))
}

/****************************************TRANSFER TRACKERS****************************************/

func (k Keeper) SetApprovalTrackerInStore(ctx sdk.Context, collectionId sdkmath.Uint, addressForApproval string, approvalId, amountTrackerId string, approvalTracker types.ApprovalTracker, level string, trackerType string, address string) error {
	marshaled_transfer_tracker, err := k.cdc.Marshal(&approvalTracker)
	if err != nil {
		return sdkerrors.Wrap(err, "Marshal types.ApprovalTracker failed")
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(approvalTrackerStoreKey(ConstructApprovalTrackerKey(collectionId, addressForApproval, approvalId, amountTrackerId, level, trackerType, address)), marshaled_transfer_tracker)
	return nil
}

func (k Keeper) SetApprovalTrackerInStoreViaKey(ctx sdk.Context, key string, approvalTracker types.ApprovalTracker) error {
	marshaled_transfer_tracker, err := k.cdc.Marshal(&approvalTracker)
	if err != nil {
		return sdkerrors.Wrap(err, "Marshal types.ApprovalTracker failed")
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(approvalTrackerStoreKey(key), marshaled_transfer_tracker)
	return nil
}

func (k Keeper) GetApprovalTrackerFromStore(ctx sdk.Context, collectionId sdkmath.Uint, addressForApproval string, approvalId string, amountTrackerId string, level string, trackerType string, address string) (types.ApprovalTracker, bool) {
	store := ctx.KVStore(k.storeKey)
	marshaled_transfer_tracker := store.Get(approvalTrackerStoreKey(ConstructApprovalTrackerKey(collectionId, addressForApproval, approvalId, amountTrackerId, level, trackerType, address)))

	var approvalTracker types.ApprovalTracker
	if len(marshaled_transfer_tracker) == 0 {
		return approvalTracker, false
	}
	k.cdc.MustUnmarshal(marshaled_transfer_tracker, &approvalTracker)
	return approvalTracker, true
}

func (k Keeper) GetApprovalTrackersFromStore(ctx sdk.Context) (approvalTrackers []*types.ApprovalTracker, ids []string) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, ApprovalTrackerKey)
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
	store := ctx.KVStore(k.storeKey)
	return store.Has(approvalTrackerStoreKey(ConstructApprovalTrackerKey(collectionId, addressForApproval, approvalId, amountTrackerId, level, trackerType, address)))
}

func (k Keeper) DeleteApprovalTrackerFromStore(ctx sdk.Context, collectionId sdkmath.Uint, addressForApproval string, approvalId, amountTrackerId string, level string, trackerType string, address string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(approvalTrackerStoreKey(ConstructApprovalTrackerKey(collectionId, addressForApproval, approvalId, amountTrackerId, level, trackerType, address)))
}
