package keeper

import (
	"context"
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	newtypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
	oldtypes "github.com/bitbadges/bitbadgeschain/x/badges/types/v19"
	ibcratelimittypes "github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/types"
	poolmanagertypes "github.com/bitbadges/bitbadgeschain/x/poolmanager/types"
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

// PoolManagerKeeperI defines the interface needed for poolmanager migrations
type PoolManagerKeeperI interface {
	GetParams(ctx sdk.Context) poolmanagertypes.Params
	SetParams(ctx sdk.Context, params poolmanagertypes.Params)
}

// MigratePoolManagerTakerFee updates the poolmanager default taker fee to 0.1% (0.001)
func MigratePoolManagerTakerFee(ctx sdk.Context, poolManagerKeeper PoolManagerKeeperI) error {
	poolManagerParams := poolManagerKeeper.GetParams(ctx)
	poolManagerParams.TakerFeeParams.DefaultTakerFee = osmomath.MustNewDecFromStr("0.001")
	poolManagerKeeper.SetParams(ctx, poolManagerParams)
	return nil
}

// IBCRateLimitKeeperI defines the interface needed for IBC rate limit migrations
type IBCRateLimitKeeperI interface {
	GetParams(ctx sdk.Context) ibcratelimittypes.Params
	SetParams(ctx sdk.Context, params ibcratelimittypes.Params) error
}

// MigrateIBCRateLimit adds or updates an IBC rate limit configuration
// If a rate limit with the same channel_id and denom exists, it will be updated
func MigrateIBCRateLimit(ctx sdk.Context, rateLimitKeeper IBCRateLimitKeeperI, rateLimitConfig ibcratelimittypes.RateLimitConfig) error {
	params := rateLimitKeeper.GetParams(ctx)

	// Find if a rate limit with the same channel_id and denom exists
	foundIndex := -1
	for i, config := range params.RateLimits {
		if config.ChannelId == rateLimitConfig.ChannelId && config.Denom == rateLimitConfig.Denom {
			foundIndex = i
			break
		}
	}

	// Update or append the rate limit
	if foundIndex >= 0 {
		// Update existing rate limit
		params.RateLimits[foundIndex] = rateLimitConfig
	} else {
		// Append new rate limit
		params.RateLimits = append(params.RateLimits, rateLimitConfig)
	}

	// Validate and set the updated params
	return rateLimitKeeper.SetParams(ctx, params)
}

func MigrateIncomingApprovals(incomingApprovals []*newtypes.UserIncomingApproval) []*newtypes.UserIncomingApproval {
	for _, approval := range incomingApprovals {
		if approval.ApprovalCriteria == nil {
			continue
		}

		// Ensure altTimeChecks is properly initialized (nil for migrated data)
		// This field is new and won't exist in old data, so it will be nil by default
		// We explicitly set it to nil to ensure proper migration
		if approval.ApprovalCriteria.AltTimeChecks == nil {
			approval.ApprovalCriteria.AltTimeChecks = nil
		}
	}

	return incomingApprovals
}

func MigrateOutgoingApprovals(outgoingApprovals []*newtypes.UserOutgoingApproval) []*newtypes.UserOutgoingApproval {
	for _, approval := range outgoingApprovals {
		if approval.ApprovalCriteria == nil {
			continue
		}

		// Ensure altTimeChecks is properly initialized (nil for migrated data)
		// This field is new and won't exist in old data, so it will be nil by default
		// We explicitly set it to nil to ensure proper migration
		if approval.ApprovalCriteria.AltTimeChecks == nil {
			approval.ApprovalCriteria.AltTimeChecks = nil
		}
	}

	return outgoingApprovals
}

func MigrateApprovals(collectionApprovals []*newtypes.CollectionApproval) []*newtypes.CollectionApproval {
	for _, approval := range collectionApprovals {
		if approval.ApprovalCriteria == nil {
			continue
		}

		// Ensure altTimeChecks is properly initialized (nil for migrated data)
		// This field is new and won't exist in old data, so it will be nil by default
		// We explicitly set it to nil to ensure proper migration
		if approval.ApprovalCriteria.AltTimeChecks == nil {
			approval.ApprovalCriteria.AltTimeChecks = nil
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
	// Iterate through pool IDs from 1 to a reasonable upper bound
	// We check up to 10000 pools - if there are more, they will be handled when created
	// maxPoolId := uint64(10000)
	// for poolId := uint64(1); poolId < maxPoolId; poolId++ {
	// 	pool, err := k.gammKeeper.GetPool(ctx, poolId)
	// 	if err != nil {
	// 		// Pool doesn't exist, continue to next ID
	// 		continue
	// 	}

	// 	// Get pool address
	// 	poolAddress := pool.GetAddress().String()

	// 	// Set pool address as reserved protocol address
	// 	if err := k.SetReservedProtocolAddressInStore(ctx, poolAddress, true); err != nil {
	// 		// Log error but continue - don't fail migration for individual pools
	// 		ctx.Logger().Error(fmt.Sprintf("Failed to set pool %d address as reserved protocol: %v", poolId, err))
	// 		continue
	// 	}

	// 	// Cache the pool address -> pool ID mapping
	// 	k.SetPoolAddressInCache(ctx, poolAddress, poolId)
	// }

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
