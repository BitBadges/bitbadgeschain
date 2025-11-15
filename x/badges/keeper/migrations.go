package keeper

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	newtypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
	oldtypes "github.com/bitbadges/bitbadgeschain/x/badges/types/v18"
)

// MigrateBadgesKeeper migrates the tokens keeper to set all approval versions to 0
func (k Keeper) MigrateBadgesKeeper(ctx sdk.Context) error {

	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})

	if err := MigratePools(ctx, k); err != nil {
		return err
	}

	if err := MigrateCollections(ctx, store, k); err != nil {
		return err
	}

	if err := MigrateBalances(ctx, store, k); err != nil {
		return err
	}

	if err := MigrateAddressLists(ctx, store, k); err != nil {
		return err
	}

	if err := MigrateApprovalTrackers(ctx, store, k); err != nil {
		return err
	}

	if err := MigrateDynamicStores(ctx, store, k); err != nil {
		return err
	}

	return nil
}

func MigrateIncomingApprovals(incomingApprovals []*newtypes.UserIncomingApproval) []*newtypes.UserIncomingApproval {
	for _, approval := range incomingApprovals {
		if approval.ApprovalCriteria == nil {
			continue
		}

		// Ensure address checks are properly initialized (nil is fine, but ensure they're not uninitialized)
		// Address checks will be nil for old data, which is the correct default
		if approval.ApprovalCriteria.SenderChecks == nil {
			// Keep as nil - this is the default for new fields
		}
		if approval.ApprovalCriteria.InitiatorChecks == nil {
			// Keep as nil - this is the default for new fields
		}
	}

	return incomingApprovals
}

func MigrateOutgoingApprovals(outgoingApprovals []*newtypes.UserOutgoingApproval) []*newtypes.UserOutgoingApproval {
	for _, approval := range outgoingApprovals {
		if approval.ApprovalCriteria == nil {
			continue
		}

		// Ensure address checks are properly initialized (nil is fine, but ensure they're not uninitialized)
		// Address checks will be nil for old data, which is the correct default
		if approval.ApprovalCriteria.RecipientChecks == nil {
			// Keep as nil - this is the default for new fields
		}
		if approval.ApprovalCriteria.InitiatorChecks == nil {
			// Keep as nil - this is the default for new fields
		}
	}

	return outgoingApprovals
}

func MigrateApprovals(collectionApprovals []*newtypes.CollectionApproval) []*newtypes.CollectionApproval {
	for _, approval := range collectionApprovals {
		if approval.ApprovalCriteria == nil {
			continue
		}

		// Ensure address checks are properly initialized (nil is fine, but ensure they're not uninitialized)
		// Address checks will be nil for old data, which is the correct default
		if approval.ApprovalCriteria.SenderChecks == nil {
			// Keep as nil - this is the default for new fields
		}
		if approval.ApprovalCriteria.RecipientChecks == nil {
			// Keep as nil - this is the default for new fields
		}
		if approval.ApprovalCriteria.InitiatorChecks == nil {
			// Keep as nil - this is the default for new fields
		}
	}

	return collectionApprovals
}

// convertUintRange converts old v9 UintRange to new UintRange
func convertUintRange(oldRange *oldtypes.UintRange) *newtypes.UintRange {
	return &newtypes.UintRange{
		Start: newtypes.Uint(oldRange.Start),
		End:   newtypes.Uint(oldRange.End),
	}
}

// convertUintRanges converts a slice of old v9 UintRange to new UintRange
func convertUintRanges(oldRanges []*oldtypes.UintRange) []*newtypes.UintRange {
	newRanges := make([]*newtypes.UintRange, len(oldRanges))
	for i, oldRange := range oldRanges {
		newRanges[i] = convertUintRange(oldRange)
	}
	return newRanges
}

// MigratePools iterates through all existing pools and sets their addresses as reserved protocol addresses
// and caches them in the pool address cache
func MigratePools(ctx sdk.Context, k Keeper) error {
	if k.gammKeeper == nil {
		// If GammKeeper is not set, skip pool migration
		// This allows the migration to work even if gamm module is not available
		return nil
	}

	// Iterate through pool IDs from 1 to a reasonable upper bound
	// We check up to 10000 pools - if there are more, they will be handled when created
	maxPoolId := uint64(10000)
	for poolId := uint64(1); poolId < maxPoolId; poolId++ {
		pool, err := k.gammKeeper.GetPool(ctx, poolId)
		if err != nil {
			// Pool doesn't exist, continue to next ID
			continue
		}

		// Get pool address
		poolAddress := pool.GetAddress().String()

		// Set pool address as reserved protocol address
		if err := k.SetReservedProtocolAddressInStore(ctx, poolAddress, true); err != nil {
			// Log error but continue - don't fail migration for individual pools
			ctx.Logger().Error(fmt.Sprintf("Failed to set pool %d address as reserved protocol: %v", poolId, err))
			continue
		}

		// Cache the pool address -> pool ID mapping
		k.SetPoolAddressInCache(ctx, poolAddress, poolId)
	}

	return nil
}

func MigrateCollections(ctx sdk.Context, store storetypes.KVStore, k Keeper) error {
	iterator := storetypes.KVStorePrefixIterator(store, CollectionKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var oldCollection oldtypes.TokenCollection
		k.cdc.MustUnmarshal(iterator.Value(), &oldCollection)

		// Convert to JSON
		jsonBytes, err := json.Marshal(oldCollection)
		if err != nil {
			return err
		}

		// Unmarshal into new type
		var newCollection newtypes.TokenCollection
		if err := json.Unmarshal(jsonBytes, &newCollection); err != nil {
			return err
		}

		newCollection.CollectionApprovals = MigrateApprovals(newCollection.CollectionApprovals)
		newCollection.DefaultBalances.IncomingApprovals = MigrateIncomingApprovals(newCollection.DefaultBalances.IncomingApprovals)
		newCollection.DefaultBalances.OutgoingApprovals = MigrateOutgoingApprovals(newCollection.DefaultBalances.OutgoingApprovals)

		// Migrate invariants: set new fields to default values
		if newCollection.Invariants != nil {
			// Set cosmosCoinBackedPath to nil (default for message type)
			newCollection.Invariants.CosmosCoinBackedPath = nil
			// Set noForcefulPostMintTransfers to false (default for bool)
			newCollection.Invariants.NoForcefulPostMintTransfers = false
			// Set disablePoolCreation to false (default for bool - allows pool creation by default)
			newCollection.Invariants.DisablePoolCreation = false
		}

		// Set reserved protocol addresses for cosmos coin backed path and wrapper paths
		// This mirrors the logic in msg_server_universal_update_collection.go
		if newCollection.Invariants != nil && newCollection.Invariants.CosmosCoinBackedPath != nil {
			backedPathAddress := newCollection.Invariants.CosmosCoinBackedPath.Address
			if backedPathAddress != "" {
				if err := k.SetReservedProtocolAddressInStore(ctx, backedPathAddress, true); err != nil {
					// Log error but continue - don't fail migration for individual collections
					ctx.Logger().Error(fmt.Sprintf("Failed to set cosmos coin backed path address as reserved protocol for collection %s: %v", newCollection.CollectionId.String(), err))
				}
			}
		}

		// Set reserved protocol addresses for cosmos coin wrapper paths
		for _, wrapperPath := range newCollection.CosmosCoinWrapperPaths {
			if wrapperPath.Address != "" {
				if err := k.SetReservedProtocolAddressInStore(ctx, wrapperPath.Address, true); err != nil {
					// Log error but continue - don't fail migration for individual collections
					ctx.Logger().Error(fmt.Sprintf("Failed to set cosmos coin wrapper path address as reserved protocol for collection %s: %v", newCollection.CollectionId.String(), err))
				}
			}
		}

		// Save the updated collection
		if err := k.SetCollectionInStore(ctx, &newCollection); err != nil {
			return err
		}
	}

	return nil
}

func MigrateBalances(ctx context.Context, store storetypes.KVStore, k Keeper) error {
	iterator := storetypes.KVStorePrefixIterator(store, UserBalanceKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var UserBalance oldtypes.UserBalanceStore
		k.cdc.MustUnmarshal(iterator.Value(), &UserBalance)

		// Convert to JSON
		jsonBytes, err := json.Marshal(UserBalance)
		if err != nil {
			return err
		}

		// Unmarshal into old type
		var oldBalance newtypes.UserBalanceStore
		if err := json.Unmarshal(jsonBytes, &oldBalance); err != nil {
			return err
		}

		oldBalance.IncomingApprovals = MigrateIncomingApprovals(oldBalance.IncomingApprovals)
		oldBalance.OutgoingApprovals = MigrateOutgoingApprovals(oldBalance.OutgoingApprovals)

		store.Set(iterator.Key(), k.cdc.MustMarshal(&oldBalance))
	}

	return nil
}

func MigrateAddressLists(ctx context.Context, store storetypes.KVStore, k Keeper) error {
	iterator := storetypes.KVStorePrefixIterator(store, AddressListKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var AddressList oldtypes.AddressList
		k.cdc.MustUnmarshal(iterator.Value(), &AddressList)

		// Convert to JSON
		jsonBytes, err := json.Marshal(AddressList)
		if err != nil {
			return err
		}

		// Unmarshal into old type
		var oldAddressList newtypes.AddressList
		if err := json.Unmarshal(jsonBytes, &oldAddressList); err != nil {
			return err
		}

		store.Set(iterator.Key(), k.cdc.MustMarshal(&oldAddressList))
	}

	return nil
}

func MigrateApprovalTrackers(ctx context.Context, store storetypes.KVStore, k Keeper) error {
	iterator := storetypes.KVStorePrefixIterator(store, ApprovalTrackerKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var ApprovalTracker oldtypes.ApprovalTracker
		k.cdc.MustUnmarshal(iterator.Value(), &ApprovalTracker)

		// Convert to JSON
		jsonBytes, err := json.Marshal(ApprovalTracker)
		if err != nil {
			return err
		}

		// Unmarshal into old type
		var oldApprovalTracker newtypes.ApprovalTracker
		if err := json.Unmarshal(jsonBytes, &oldApprovalTracker); err != nil {
			return err
		}

		store.Set(iterator.Key(), k.cdc.MustMarshal(&oldApprovalTracker))
	}

	return nil
}

func MigrateDynamicStores(ctx context.Context, store storetypes.KVStore, k Keeper) error {
	// Migrate base dynamic stores
	iterator := storetypes.KVStorePrefixIterator(store, DynamicStoreKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var oldDynamicStore oldtypes.DynamicStore
		k.cdc.MustUnmarshal(iterator.Value(), &oldDynamicStore)

		// Convert to JSON
		jsonBytes, err := json.Marshal(oldDynamicStore)
		if err != nil {
			return err
		}

		// Unmarshal into new type
		var newDynamicStore newtypes.DynamicStore
		if err := json.Unmarshal(jsonBytes, &newDynamicStore); err != nil {
			return err
		}

		// Save the updated dynamic store
		if err := k.SetDynamicStoreInStore(sdk.UnwrapSDKContext(ctx), newDynamicStore); err != nil {
			return err
		}
	}

	// Migrate dynamic store values
	valueIterator := storetypes.KVStorePrefixIterator(store, DynamicStoreValueKey)
	defer valueIterator.Close()
	for ; valueIterator.Valid(); valueIterator.Next() {
		var oldDynamicStoreValue oldtypes.DynamicStoreValue
		k.cdc.MustUnmarshal(valueIterator.Value(), &oldDynamicStoreValue)

		// Convert to JSON
		jsonBytes, err := json.Marshal(oldDynamicStoreValue)
		if err != nil {
			return err
		}

		// Unmarshal into new type
		var newDynamicStoreValue newtypes.DynamicStoreValue
		if err := json.Unmarshal(jsonBytes, &newDynamicStoreValue); err != nil {
			return err
		}

		// Save the updated dynamic store value
		if err := k.SetDynamicStoreValueInStore(sdk.UnwrapSDKContext(ctx), newDynamicStoreValue.StoreId, newDynamicStoreValue.Address, newDynamicStoreValue.Value); err != nil {
			return err
		}
	}

	return nil
}
