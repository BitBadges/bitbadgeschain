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
func (k Keeper) GetCollectionFromStore(ctx sdk.Context, collectionId sdk.Uint) (types.BadgeCollection, bool) {
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
func (k Keeper) StoreHasCollectionID(ctx sdk.Context, collectionId sdk.Uint) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(collectionStoreKey(collectionId))
}

// DeleteCollectionFromStore deletes a badge from the store.
func (k Keeper) DeleteCollectionFromStore(ctx sdk.Context, collectionId sdk.Uint) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(collectionStoreKey(collectionId))
}

/****************************************CLAIMS****************************************/

// Sets a user balance in the store using ClaimKey ([]byte{0x02}) as the prefix. No check if store has key already.
func (k Keeper) SetClaimInStore(ctx sdk.Context, collectionId sdk.Uint, claimId sdk.Uint, Claim types.Claim) error {
	marshaled_claim_info, err := k.cdc.Marshal(&Claim)
	if err != nil {
		return sdkerrors.Wrap(err, "Marshal types.Claim failed")
	}

	store := ctx.KVStore(k.storeKey)
	claimKey := ConstructClaimKey(collectionId, claimId)
	store.Set(claimStoreKey(claimKey), marshaled_claim_info)
	return nil
}

// Sets a user balance in the store using ClaimKey ([]byte{0x02}) as the prefix. No check if store has key already.
func (k Keeper) SetClaimInStoreWithKey(ctx sdk.Context, claimKey string, Claim types.Claim) error {
	marshaled_claim_info, err := k.cdc.Marshal(&Claim)
	if err != nil {
		return sdkerrors.Wrap(err, "Marshal types.Claim failed")
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(claimStoreKey(claimKey), marshaled_claim_info)
	return nil
}

// Gets a user balance from the store according to the balanceID.
func (k Keeper) GetClaimFromStore(ctx sdk.Context, collectionId sdk.Uint, claimId sdk.Uint) (types.Claim, bool) {
	store := ctx.KVStore(k.storeKey)
	claimKey := ConstructClaimKey(collectionId, claimId)
	marshaled_claim_info := store.Get(claimStoreKey(claimKey))

	var Claim types.Claim
	if len(marshaled_claim_info) == 0 {
		return Claim, false
	}
	k.cdc.MustUnmarshal(marshaled_claim_info, &Claim)
	return Claim, true
}

// GetClaimsFromStore defines a method for returning all user balances information by key.
func (k Keeper) GetClaimsFromStore(ctx sdk.Context) (claims []*types.Claim, collectionIds []sdk.Uint, claimIds []sdk.Uint) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, ClaimKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var Claim types.Claim
		k.cdc.MustUnmarshal(iterator.Value(), &Claim)
		claims = append(claims, &Claim)

		claimKeyDetails := GetDetailsFromClaimKey(string(iterator.Key()))
		collectionIds = append(collectionIds, claimKeyDetails.collectionId)
		claimIds = append(claimIds, claimKeyDetails.claimId)
	}
	return
}

// GetClaimIdsFromStore defines a method for returning all keys of all user balances.
func (k Keeper) GetClaimIdsFromStore(ctx sdk.Context) (ids []string) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, ClaimKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		ids = append(ids, string(iterator.Key()))
	}
	return
}

// StoreHasClaimID determines whether the specified user balanceID exists in the store
func (k Keeper) StoreHasClaim(ctx sdk.Context, collectionId sdk.Uint, claimId sdk.Uint) bool {
	store := ctx.KVStore(k.storeKey)
	claimKey := ConstructClaimKey(collectionId, claimId)
	return store.Has(claimStoreKey(claimKey))
}

// DeleteClaimFromStore deletes a user balance from the store.
func (k Keeper) DeleteClaimFromStore(ctx sdk.Context, collectionId sdk.Uint, claimId sdk.Uint) {
	store := ctx.KVStore(k.storeKey)
	claimKey := ConstructClaimKey(collectionId, claimId)
	store.Delete(claimStoreKey(claimKey))
}

/****************************************USER BALANCES****************************************/

// Sets a user balance in the store using UserBalanceKey ([]byte{0x02}) as the prefix. No check if store has key already.
func (k Keeper) SetUserBalanceInStore(ctx sdk.Context, balanceKey string, UserBalance types.UserBalanceStore) error {
	marshaled_badge_balance_info, err := k.cdc.Marshal(&UserBalance)
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
func (k Keeper) GetUserBalanceFromStore(ctx sdk.Context, balanceKey string) (types.UserBalanceStore, bool) {
	store := ctx.KVStore(k.storeKey)
	marshaled_badge_balance_info := store.Get(userBalanceStoreKey(balanceKey))

	var UserBalance types.UserBalanceStore
	if len(marshaled_badge_balance_info) == 0 {
		return UserBalance, false
	}
	k.cdc.MustUnmarshal(marshaled_badge_balance_info, &UserBalance)
	return UserBalance, true
}

// GetUserBalancesFromStore defines a method for returning all user balances information by key.
func (k Keeper) GetUserBalancesFromStore(ctx sdk.Context) (balances []*types.UserBalanceStore, addresses []string, ids []sdk.Uint) {
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

/****************************************TRANSFER MANAGER REQUESTS****************************************/

// Checks if a certain address has requested a managerial transfer
func (k Keeper) HasAddressRequestedManagerTransfer(ctx sdk.Context, collectionId sdk.Uint, address string) bool {
	store := ctx.KVStore(k.storeKey)
	key := ConstructUpdateManagerRequestKey(collectionId, address)
	return store.Has(managerTransferRequestKey(key))
}

// Creates a transfer manager request for the given address and collectionId.
func (k Keeper) CreateUpdateManagerRequest(ctx sdk.Context, collectionId sdk.Uint, address string) error {
	request := []byte{}
	store := ctx.KVStore(k.storeKey)
	key := ConstructUpdateManagerRequestKey(collectionId, address)
	store.Set(managerTransferRequestKey(key), request)
	return nil
}

// Deletes a transfer manager request for the given address and collectionId.
func (k Keeper) RemoveUpdateManagerRequest(ctx sdk.Context, collectionId sdk.Uint, address string) error {
	key := ConstructUpdateManagerRequestKey(collectionId, address)
	store := ctx.KVStore(k.storeKey)

	store.Delete(managerTransferRequestKey(key))
	return nil
}

/****************************************NEXT ASSET ID****************************************/

// Gets the next badge ID.
func (k Keeper) GetNextCollectionId(ctx sdk.Context) sdk.Uint {
	store := ctx.KVStore(k.storeKey)
	nextID := types.NewUintFromString(string((store.Get(nextCollectionIdKey()))))
	return nextID
}

// Sets the next asset ID. Should only be used in InitGenesis. Everything else should call IncrementNextAssetID()
func (k Keeper) SetNextCollectionId(ctx sdk.Context, nextID sdk.Uint) {
	store := ctx.KVStore(k.storeKey)
	store.Set(nextCollectionIdKey(), []byte(nextID.String()))
}

// Increments the next badge ID by 1.
func (k Keeper) IncrementNextCollectionId(ctx sdk.Context) {
	nextID := k.GetNextCollectionId(ctx)
	k.SetNextCollectionId(ctx, nextID.AddUint64(1)) //susceptible to overflow but by that time we will have 2^64 badges which isn't totally feasible
}

/****************************************NEXT CLAIM ID****************************************/
//Note these claim IDs are different than the ones within each collection.
//These are the IDs for the claim stores themselves, not the claim IDs for within a collection.

// Gets the next badge ID.
func (k Keeper) GetNextClaimId(ctx sdk.Context) sdk.Uint {
	store := ctx.KVStore(k.storeKey)
	nextID := types.NewUintFromString(string((store.Get(nextClaimIdKey()))))
	return nextID
}

// Sets the next asset ID. Should only be used in InitGenesis. Everything else should call IncrementNextAssetID()
func (k Keeper) SetNextClaimId(ctx sdk.Context, nextID sdk.Uint) {
	store := ctx.KVStore(k.storeKey)
	store.Set(nextClaimIdKey(), []byte(nextID.String()))
}

// Increments the next badge ID by 1.
func (k Keeper) IncrementNextClaimId(ctx sdk.Context) {
	nextID := k.GetNextClaimId(ctx)
	k.SetNextClaimId(ctx, nextID.Incr()) //susceptible to overflow but by that time we will have 2^64 badges which isn't totally feasible
}

/****************************************Claims****************************************/
// Sets a usedClaimData in the store using UsedClaimDataKey ([]byte{0x07}) as the prefix. No check if store has key already.
func (k Keeper) IncrementNumUsedForClaimInStore(ctx sdk.Context, collectionId sdk.Uint, claimId sdk.Uint) (sdk.Uint, error) {
	store := ctx.KVStore(k.storeKey)
	currBytes := store.Get(usedClaimDataStoreKey(ConstructUsedClaimDataKey(collectionId, claimId)))
	curr := sdk.NewUint(0)
	if currBytes != nil {
		currUint, err := strconv.ParseUint(string((currBytes)), 10, 64)
		if err != nil {
			panic("Failed to parse num used")
		}

		curr = sdk.NewUint(currUint)
	}
	incrementedNum := curr.AddUint64(1)
	store.Set(usedClaimDataStoreKey(ConstructUsedClaimDataKey(collectionId, claimId)), []byte(curr.Incr().String()))
	return incrementedNum, nil
}

func (k Keeper) IncrementNumUsedForChallengeInStore(ctx sdk.Context, collectionId sdk.Uint, claimId sdk.Uint, challengeId sdk.Uint, leafIndex sdk.Uint) (sdk.Uint, error) {
	store := ctx.KVStore(k.storeKey)
	currBytes := store.Get(usedClaimChallengeStoreKey(ConstructUsedClaimChallengeKey(collectionId, claimId, challengeId, leafIndex)))
	curr := sdk.NewUint(0)
	if currBytes != nil {
		currUint, err := strconv.ParseUint(string((currBytes)), 10, 64)
		if err != nil {
			panic("Failed to parse num used")
		}

		curr = sdk.NewUint(currUint)
	}
	incrementedNum := curr.AddUint64(1)
	store.Set(usedClaimChallengeStoreKey(ConstructUsedClaimChallengeKey(collectionId, claimId, challengeId, leafIndex)), []byte(curr.Incr().String()))
	return incrementedNum, nil
}

func (k Keeper) IncrementNumUsedForAddressInStore(ctx sdk.Context, collectionId sdk.Uint, claimId sdk.Uint, address string) (sdk.Uint, error) {
	store := ctx.KVStore(k.storeKey)
	currBytes := store.Get(usedClaimAddressStoreKey(ConstructUsedClaimAddressKey(collectionId, claimId, address)))
	curr := sdk.NewUint(0)
	if currBytes != nil {
		currUint, err := strconv.ParseUint(string((currBytes)), 10, 64)
		if err != nil {
			panic("Failed to parse num used")
		}

		curr = sdk.NewUint(currUint)
	}
	incrementedNum := curr.AddUint64(1)
	store.Set(usedClaimAddressStoreKey(ConstructUsedClaimAddressKey(collectionId, claimId, address)), []byte(curr.Incr().String()))
	return incrementedNum, nil
}

func (k Keeper) IncrementNumUsedForWhitelistIndexInStore(ctx sdk.Context, collectionId sdk.Uint, claimId sdk.Uint, whitelistLeafIndex sdk.Uint) (sdk.Uint, error) {
	store := ctx.KVStore(k.storeKey)
	currBytes := store.Get(usedWhitelistIndexStoreKey(ConstructUsedWhitelistIndexKey(collectionId, claimId, whitelistLeafIndex)))
	curr := sdk.NewUint(0)
	if currBytes != nil {
		currUint, err := strconv.ParseUint(string((currBytes)), 10, 64)
		if err != nil {
			panic("Failed to parse num used")
		}

		curr = sdk.NewUint(currUint)
	}
	incrementedNum := curr.AddUint64(1)
	store.Set(usedWhitelistIndexStoreKey(ConstructUsedWhitelistIndexKey(collectionId, claimId, whitelistLeafIndex)), []byte(curr.Incr().String()))
	return incrementedNum, nil
}
