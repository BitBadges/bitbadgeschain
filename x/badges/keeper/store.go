package keeper

import (
	"math"
	"strconv"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkerrors "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkmath "cosmossdk.io/math"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
)

// The following methods are used for the store and everything associated with tokens.
// All preconditions and checks must be handled before these functions are called.
// This file handles storing collections, balances, approvals, used challenges, next collection ID, etc.

// All the following CRUD operations must obey the key prefixes defined in keys.go.

/****************************************COLLECTIONS****************************************/

// validateCollectionBeforeStore validates a collection before storing it
func (k Keeper) validateCollectionBeforeStore(ctx sdk.Context, collection *types.TokenCollection) error {
	// Validate collection approvals with invariants
	if collection.Invariants != nil && collection.Invariants.NoCustomOwnershipTimes {
		if err := types.ValidateCollectionApprovalsWithInvariants(ctx, collection.CollectionApprovals, false, collection); err != nil {
			return sdkerrors.Wrap(err, "collection approval validation failed")
		}
	}
	return nil
}

// Sets a token in the store using BadgeKey ([]byte{0x01}) as the prefix. No check if store has key already.
func (k Keeper) SetCollectionInStore(ctx sdk.Context, collection *types.TokenCollection) error {
	// Validate collection before storing
	if err := k.validateCollectionBeforeStore(ctx, collection); err != nil {
		return err
	}
	marshaled_token, err := k.cdc.Marshal(collection)
	if err != nil {
		return sdkerrors.Wrap(err, "Marshal types.TokenCollection failed")
	}

	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	store.Set(collectionStoreKey(collection.CollectionId), marshaled_token)
	return nil
}

// Gets a token from the store according to the collectionId.
func (k Keeper) GetCollectionFromStore(ctx sdk.Context, collectionId sdkmath.Uint) (*types.TokenCollection, bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	marshaled_collection := store.Get(collectionStoreKey(collectionId))

	var collection types.TokenCollection
	if len(marshaled_collection) == 0 {
		return &collection, false
	}
	k.cdc.MustUnmarshal(marshaled_collection, &collection)
	return &collection, true
}

// GetCollectionsFromStore defines a method for returning all tokens information by key.
func (k Keeper) GetCollectionsFromStore(ctx sdk.Context) (collections []*types.TokenCollection) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	iterator := storetypes.KVStorePrefixIterator(store, CollectionKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var collection types.TokenCollection
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

// DeleteCollectionFromStore deletes a token from the store.
func (k Keeper) DeleteCollectionFromStore(ctx sdk.Context, collectionId sdkmath.Uint) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	store.Delete(collectionStoreKey(collectionId))
}

/****************************************USER BALANCES****************************************/

// validateUserBalanceBeforeStore validates a user balance before storing it
func (k Keeper) validateUserBalanceBeforeStore(ctx sdk.Context, balanceKey string, userBalance *types.UserBalanceStore, collection *types.TokenCollection) error {
	// Get collection if not provided
	if collection == nil {
		balanceKeyDetails, err := GetDetailsFromBalanceKey(balanceKey)
		if err != nil {
			return sdkerrors.Wrapf(err, "invalid balance key format")
		}
		collectionId := balanceKeyDetails.collectionId
		var found bool
		collection, found = k.GetCollectionFromStore(ctx, collectionId)
		if !found {
			return sdkerrors.Wrapf(types.ErrInvalidRequest, "collection not found for balance validation")
		}
	}

	// Check invariants if enabled
	if collection.Invariants != nil {
		// Check noCustomOwnershipTimes invariant
		if collection.Invariants.NoCustomOwnershipTimes {
			// Validate incoming approvals
			for _, incomingApproval := range userBalance.IncomingApprovals {
				if err := types.ValidateNoCustomOwnershipTimesInvariant(incomingApproval.OwnershipTimes, true); err != nil {
					return sdkerrors.Wrap(err, "incoming approval ownership times validation failed")
				}
			}

			// Validate outgoing approvals
			for _, outgoingApproval := range userBalance.OutgoingApprovals {
				if err := types.ValidateNoCustomOwnershipTimesInvariant(outgoingApproval.OwnershipTimes, true); err != nil {
					return sdkerrors.Wrap(err, "outgoing approval ownership times validation failed")
				}
			}

			// Validate balances ownership times
			for _, balance := range userBalance.Balances {
				if err := types.ValidateNoCustomOwnershipTimesInvariant(balance.OwnershipTimes, true); err != nil {
					return sdkerrors.Wrap(err, "balance ownership times validation failed")
				}
			}
		}

		// Check maxSupplyPerId invariant if we're setting Total address balances
		if !collection.Invariants.MaxSupplyPerId.IsNil() && !collection.Invariants.MaxSupplyPerId.IsZero() {
			balanceKeyDetails, err := GetDetailsFromBalanceKey(balanceKey)
			if err != nil {
				return sdkerrors.Wrapf(err, "invalid balance key format")
			}
			if types.IsTotalAddress(balanceKeyDetails.address) {
				// Validate that no balance amount exceeds maxSupplyPerId
				for _, balance := range userBalance.Balances {
					if balance.Amount.GT(collection.Invariants.MaxSupplyPerId) {
						return sdkerrors.Wrapf(types.ErrInvalidRequest, "maxSupplyPerId invariant violation: balance amount %s exceeds maximum supply per ID %s", balance.Amount.String(), collection.Invariants.MaxSupplyPerId.String())
					}
				}
			}
		}
	}
	return nil
}

// Sets a user balance in the store using UserBalanceKey ([]byte{0x02}) as the prefix. No check if store has key already.
func (k Keeper) SetUserBalanceInStore(ctx sdk.Context, balanceKey string, UserBalance *types.UserBalanceStore) error {
	// Validate user balance before storing
	if err := k.validateUserBalanceBeforeStore(ctx, balanceKey, UserBalance, nil); err != nil {
		return err
	}

	// NOTE: We always store a non-nil permissions object to prevent issues where
	// nil permissions would marshal to zero length, causing default balances to be
	// incorrectly populated again during deserialization.
	if UserBalance.UserPermissions == nil {
		UserBalance.UserPermissions = &types.UserPermissions{
			CanUpdateOutgoingApprovals:                         []*types.UserOutgoingApprovalPermission{},
			CanUpdateIncomingApprovals:                         []*types.UserIncomingApprovalPermission{},
			CanUpdateAutoApproveSelfInitiatedOutgoingTransfers: []*types.ActionPermission{},
			CanUpdateAutoApproveSelfInitiatedIncomingTransfers: []*types.ActionPermission{},
			CanUpdateAutoApproveAllIncomingTransfers:           []*types.ActionPermission{},
		}
	}

	marshaled_token_balance_info, err := k.cdc.Marshal(UserBalance)
	if err != nil {
		return sdkerrors.Wrap(err, "Marshal types.UserBalanceStore failed")
	}

	//Prevent accidental non-BitBadges addresses from being stored
	balanceKeyDetails, err := GetDetailsFromBalanceKey(balanceKey)
	if err != nil {
		return sdkerrors.Wrapf(err, "invalid balance key format")
	}
	if !types.IsSpecialAddress(balanceKeyDetails.address) {
		if err = types.ValidateAddress(balanceKeyDetails.address, false); err != nil {
			return sdkerrors.Wrap(err, "Invalid address")
		}
	}

	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	store.Set(userBalanceStoreKey(balanceKey), marshaled_token_balance_info)
	return nil
}

// Gets a user balance from the store according to the balanceID.
func (k Keeper) GetUserBalanceFromStore(ctx sdk.Context, balanceKey string) (*types.UserBalanceStore, bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	marshaled_token_balance_info := store.Get(userBalanceStoreKey(balanceKey))

	var UserBalance types.UserBalanceStore
	if len(marshaled_token_balance_info) == 0 {
		return &UserBalance, false
	}
	k.cdc.MustUnmarshal(marshaled_token_balance_info, &UserBalance)
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

		balanceKeyDetails, err := GetDetailsFromBalanceKey(string(iterator.Key()[1:]))
		if err != nil {
			// Log error and continue processing other keys
			// This could indicate data corruption, but we continue to avoid breaking the entire operation
			k.logger.Error("failed to parse balance key", "key", string(iterator.Key()), "error", err)
			continue
		}
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
func (k Keeper) IncrementNextCollectionId(ctx sdk.Context) error {
	nextID := k.GetNextCollectionId(ctx)

	// Check for overflow before incrementing
	if nextID.Equal(sdkmath.NewUint(math.MaxUint64)) {
		return sdkerrors.Wrapf(types.ErrOverflow, "collection ID overflow: maximum number of collections reached")
	}

	newID := nextID.AddUint64(1)
	k.SetNextCollectionId(ctx, newID)
	return nil
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

// Increments the next address list counter by 1.
func (k Keeper) IncrementNextAddressListCounter(ctx sdk.Context) error {
	nextID := k.GetNextAddressListCounter(ctx)
	newID := nextID.AddUint64(1)

	// Check for overflow: if the new ID is less than the original, we've overflowed
	if newID.LT(nextID) {
		return sdkerrors.Wrapf(types.ErrOverflow, "address list counter overflow: maximum number of address lists reached")
	}

	k.SetNextAddressListCounter(ctx, newID)
	return nil
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
			return sdkmath.NewUint(0), sdkerrors.Wrapf(err, "failed to parse num used")
		}

		curr = sdkmath.NewUint(currUint)
	}
	incrementedNum := curr.AddUint64(1)
	store.Set(usedClaimChallengeStoreKey(ConstructUsedClaimChallengeKey(collectionId, addressForChallenge, approvalLevel, approvalId, challengeId, leafIndex)), []byte(incrementedNum.String()))
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
			return sdkmath.NewUint(0), sdkerrors.Wrapf(err, "failed to parse num used")
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
			// Log error and continue processing - this is a data corruption issue
			// We continue to avoid breaking the entire operation, but log for monitoring
			k.logger.Error("failed to parse challenge tracker value", "value", string(iterator.Value()), "error", err)
			continue
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
			// Return 0 on parse error
			return sdkmath.NewUint(0)
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
			// Log error but continue processing - this is a data corruption issue
			continue
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
		// Return false on parse error - this indicates data corruption
		return sdkmath.NewUint(0), false
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

	// Check for overflow before incrementing
	if nextID.Equal(sdkmath.NewUint(math.MaxUint64)) {
		panic("dynamic store ID overflow: cannot increment beyond MaxUint64")
	}

	k.SetNextDynamicStoreId(ctx, nextID.AddUint64(1))
}

/****************************************DYNAMIC STORE VALUES****************************************/

// Sets a dynamic store value in the store using DynamicStoreValueKey ([]byte{0x0F}) as the prefix.
func (k Keeper) SetDynamicStoreValueInStore(ctx sdk.Context, storeId sdkmath.Uint, address string, value sdkmath.Uint) error {
	// Validate inputs
	if storeId.IsZero() {
		return sdkerrors.Wrapf(types.ErrInvalidRequest, "store ID cannot be zero")
	}
	if address == "" {
		return sdkerrors.Wrapf(types.ErrInvalidRequest, "address cannot be empty")
	}
	if value.IsNil() {
		return sdkerrors.Wrapf(types.ErrInvalidRequest, "value cannot be nil")
	}

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

// GetAllDynamicStoreValuesFromStore defines a method for returning all dynamic store values across all stores.
func (k Keeper) GetAllDynamicStoreValuesFromStore(ctx sdk.Context) (dynamicStoreValues []*types.DynamicStoreValue) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	iterator := storetypes.KVStorePrefixIterator(store, DynamicStoreValueKey)
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

/****************************************ETH SIGNATURE TRACKERS****************************************/

// SetETHSignatureTrackerInStore sets an ETH signature tracker in the store using ETHSignatureTrackerKey ([]byte{0x10}) as the prefix.
func (k Keeper) SetETHSignatureTrackerInStore(ctx sdk.Context, key string, numUsed sdkmath.Uint) error {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	store.Set(ethSignatureTrackerStoreKey(key), []byte(numUsed.String()))
	return nil
}

// GetETHSignatureTrackerFromStore gets an ETH signature tracker from the store.
func (k Keeper) GetETHSignatureTrackerFromStore(ctx sdk.Context, key string) (sdkmath.Uint, bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	marshaled_data := store.Get(ethSignatureTrackerStoreKey(key))

	if len(marshaled_data) == 0 {
		return sdkmath.NewUint(0), false
	}

	numUsed, err := sdkmath.ParseUint(string(marshaled_data))
	if err != nil {
		return sdkmath.NewUint(0), false
	}
	return numUsed, true
}

// GetETHSignatureTrackersFromStore gets all ETH signature trackers from the store.
func (k Keeper) GetETHSignatureTrackersFromStore(ctx sdk.Context) (numUsed []sdkmath.Uint, ids []string) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	iterator := storetypes.KVStorePrefixIterator(store, ETHSignatureTrackerKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		tracker, err := sdkmath.ParseUint(string(iterator.Value()))
		if err != nil {
			continue
		}
		numUsed = append(numUsed, tracker)
		ids = append(ids, string(iterator.Key()[len(ETHSignatureTrackerKey):]))
	}
	return
}

// IncrementETHSignatureTrackerInStore increments the usage count for an ETH signature tracker.
func (k Keeper) IncrementETHSignatureTrackerInStore(ctx sdk.Context, key string) (sdkmath.Uint, error) {
	currentNumUsed, exists := k.GetETHSignatureTrackerFromStore(ctx, key)
	if !exists {
		currentNumUsed = sdkmath.NewUint(0)
	}

	newNumUsed := currentNumUsed.Add(sdkmath.NewUint(1))
	err := k.SetETHSignatureTrackerInStore(ctx, key, newNumUsed)
	if err != nil {
		return sdkmath.NewUint(0), err
	}

	return newNumUsed, nil
}
