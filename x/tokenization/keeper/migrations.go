package keeper

import (
	"context"
	"encoding/json"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	newtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	oldtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types/v23"
)

// MigrateTokenizationKeeper migrates the tokenization keeper from v21 to current version
func (k Keeper) MigrateTokenizationKeeper(ctx sdk.Context) error {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})

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

// MigrateCollectionStats computes and stores holder count and circulating supply for all existing collections.
// Used by the v24->v25 module migration.
func (k Keeper) MigrateCollectionStats(ctx sdk.Context) error {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})

	collectionIterator := storetypes.KVStorePrefixIterator(store, CollectionKey)
	defer func() {
		if err := collectionIterator.Close(); err != nil {
			k.Logger().Error("failed to close collection stats migration iterator", "error", err)
		}
	}()

	for ; collectionIterator.Valid(); collectionIterator.Next() {
		var collection newtypes.TokenCollection
		k.cdc.MustUnmarshal(collectionIterator.Value(), &collection)

		holderCount := sdkmath.ZeroUint()
		circulatingSupply := sdkmath.ZeroUint()

		// Count holders by iterating balances for this collection
		balanceIterator := storetypes.KVStorePrefixIterator(store, UserBalanceKey)
		for ; balanceIterator.Valid(); balanceIterator.Next() {
			balanceKey := string(balanceIterator.Key()[len(UserBalanceKey):])
			details, err := GetDetailsFromBalanceKey(balanceKey)
			if err != nil || !details.collectionId.Equal(collection.CollectionId) {
				continue
			}

			// Skip excluded addresses
			if newtypes.IsMintOrTotalAddress(details.address) {
				continue
			}
			if k.IsBackingPathAddress(ctx, &collection, details.address) {
				continue
			}
			if k.IsWrappingPathAddress(ctx, &collection, details.address) {
				continue
			}

			var balance newtypes.UserBalanceStore
			k.cdc.MustUnmarshal(balanceIterator.Value(), &balance)

			if hasNonZeroBalance(&balance) {
				holderCount = holderCount.Add(sdkmath.OneUint())
			}
		}
		balanceIterator.Close()

		// Get Total balance for circulating supply - use Balance[] for proper range handling
		var circulatingBalances []*newtypes.Balance
		totalKey := ConstructBalanceKey(newtypes.TotalAddress, collection.CollectionId)
		totalBalance, found := k.GetUserBalanceFromStore(ctx, totalKey)
		if found {
			circulatingBalances = newtypes.DeepCopyBalances(totalBalance.Balances)
		}

		// Subtract backed address balance using proper balance subtraction
		if collection.Invariants != nil && collection.Invariants.CosmosCoinBackedPath != nil {
			backedKey := ConstructBalanceKey(collection.Invariants.CosmosCoinBackedPath.Address, collection.CollectionId)
			backedBalance, found := k.GetUserBalanceFromStore(ctx, backedKey)
			if found && len(backedBalance.Balances) > 0 {
				var err error
				circulatingBalances, err = newtypes.SubtractBalancesWithZeroForUnderflows(ctx, backedBalance.Balances, circulatingBalances)
				if err != nil {
					return err
				}
			}
		}

		// Store stats with Balance[]
		stats := &newtypes.CollectionStats{
			HolderCount: holderCount,
			Balances:    circulatingBalances,
		}
		if err := k.SetCollectionStatsInStore(ctx, collection.CollectionId, stats); err != nil {
			return err
		}

		// Calculate total for logging
		for _, bal := range circulatingBalances {
			circulatingSupply = circulatingSupply.Add(bal.Amount)
		}

		ctx.Logger().Info("Migrated collection stats",
			"collectionId", collection.CollectionId,
			"holderCount", holderCount,
			"circulatingSupply", circulatingSupply)
	}

	return nil
}

// migrateInvariants handles migration of CollectionInvariants fields.
// MaxSupplyPerId and EvmQueryChallenges are now at the base CollectionInvariants level.
// The JSON marshal/unmarshal handles field mapping automatically, so no manual migration is needed.
// This function is kept for potential future migrations.
func migrateInvariants(newCollection *newtypes.TokenCollection, oldCollection *oldtypes.TokenCollection) error {
	// v23 already has MaxSupplyPerId at base level, and the new version keeps it there.
	// EvmQueryChallenges is a new field that will default to nil/empty if not present.
	// JSON marshal/unmarshal handles all of this automatically.
	return nil
}

// migrateIncomingApprovalCriteria handles WASM contract check field removal and EVM contract check field defaults
// Note: The JSON marshal/unmarshal process automatically drops fields that don't exist in the target struct,
// so the removed mustBeWasmContract and mustNotBeWasmContract fields will be automatically ignored.
// New EVM contract check fields (mustBeEvmContract, mustNotBeEvmContract) default to false if not present.
func migrateIncomingApprovalCriteria(approvalCriteria *newtypes.IncomingApprovalCriteria) {
	if approvalCriteria == nil {
		return
	}
}

// migrateOutgoingApprovalCriteria handles WASM contract check field removal and EVM contract check field defaults
// Note: The JSON marshal/unmarshal process automatically drops fields that don't exist in the target struct,
// so the removed mustBeWasmContract and mustNotBeWasmContract fields will be automatically ignored.
// New EVM contract check fields (mustBeEvmContract, mustNotBeEvmContract) default to false if not present.
func migrateOutgoingApprovalCriteria(approvalCriteria *newtypes.OutgoingApprovalCriteria) {
	if approvalCriteria == nil {
		return
	}
}

// migrateApprovalCriteria handles WASM contract check field removal and EVM contract check field defaults
// Note: The JSON marshal/unmarshal process automatically drops fields that don't exist in the target struct,
// so the removed mustBeWasmContract and mustNotBeWasmContract fields will be automatically ignored.
// New EVM contract check fields (mustBeEvmContract, mustNotBeEvmContract) default to false if not present.
func migrateApprovalCriteria(approvalCriteria *newtypes.ApprovalCriteria) {
	if approvalCriteria == nil {
		return
	}
}

func MigrateIncomingApprovals(incomingApprovals []*newtypes.UserIncomingApproval) []*newtypes.UserIncomingApproval {
	for _, approval := range incomingApprovals {
		if approval.ApprovalCriteria == nil {
			continue
		}
		migrateIncomingApprovalCriteria(approval.ApprovalCriteria)
	}

	return incomingApprovals
}

func MigrateOutgoingApprovals(outgoingApprovals []*newtypes.UserOutgoingApproval) []*newtypes.UserOutgoingApproval {
	for _, approval := range outgoingApprovals {
		if approval.ApprovalCriteria == nil {
			continue
		}
		migrateOutgoingApprovalCriteria(approval.ApprovalCriteria)
	}

	return outgoingApprovals
}

func MigrateApprovals(collectionApprovals []*newtypes.CollectionApproval) []*newtypes.CollectionApproval {
	for _, approval := range collectionApprovals {
		if approval.ApprovalCriteria == nil {
			continue
		}
		migrateApprovalCriteria(approval.ApprovalCriteria)
	}

	return collectionApprovals
}

func MigrateCollections(ctx sdk.Context, store storetypes.KVStore, k Keeper) error {
	iterator := storetypes.KVStorePrefixIterator(store, CollectionKey)
	defer func() {
		if err := iterator.Close(); err != nil {
			// Log error but don't fail migration
			k.Logger().Error("failed to close collection migration iterator", "error", err)
		}
	}()

	const oldModuleName = "badges"
	const newModuleName = newtypes.ModuleName // "tokenization"

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

		// Migrate invariants (now a no-op since maxSupplyPerId stays at base level)
		if err := migrateInvariants(&newCollection, &oldCollection); err != nil {
			return err
		}

		// Save the updated collection (with migrated addresses)
		if err := k.SetCollectionInStore(ctx, &newCollection, true); err != nil {
			return err
		}
	}

	return nil
}

func MigrateBalances(ctx context.Context, store storetypes.KVStore, k Keeper) error {
	iterator := storetypes.KVStorePrefixIterator(store, UserBalanceKey)
	defer func() {
		if err := iterator.Close(); err != nil {
			// Log error but don't fail migration
			k.Logger().Error("failed to close balance migration iterator", "error", err)
		}
	}()

	for ; iterator.Valid(); iterator.Next() {
		var oldBalance oldtypes.UserBalanceStore
		k.cdc.MustUnmarshal(iterator.Value(), &oldBalance)

		// Convert to JSON
		jsonBytes, err := json.Marshal(oldBalance)
		if err != nil {
			return err
		}

		// Unmarshal into new type
		var newBalance newtypes.UserBalanceStore
		if err := json.Unmarshal(jsonBytes, &newBalance); err != nil {
			return err
		}

		// Migrate approvals
		newBalance.IncomingApprovals = MigrateIncomingApprovals(newBalance.IncomingApprovals)
		newBalance.OutgoingApprovals = MigrateOutgoingApprovals(newBalance.OutgoingApprovals)

		store.Set(iterator.Key(), k.cdc.MustMarshal(&newBalance))
	}

	return nil
}

func MigrateAddressLists(ctx context.Context, store storetypes.KVStore, k Keeper) error {
	iterator := storetypes.KVStorePrefixIterator(store, AddressListKey)
	defer func() {
		if err := iterator.Close(); err != nil {
			// Log error but don't fail migration
			k.Logger().Error("failed to close address list migration iterator", "error", err)
		}
	}()

	for ; iterator.Valid(); iterator.Next() {
		var oldAddressList oldtypes.AddressList
		k.cdc.MustUnmarshal(iterator.Value(), &oldAddressList)

		// Convert to JSON
		jsonBytes, err := json.Marshal(oldAddressList)
		if err != nil {
			return err
		}

		// Unmarshal into new type
		var newAddressList newtypes.AddressList
		if err := json.Unmarshal(jsonBytes, &newAddressList); err != nil {
			return err
		}

		store.Set(iterator.Key(), k.cdc.MustMarshal(&newAddressList))
	}

	return nil
}

func MigrateApprovalTrackers(ctx context.Context, store storetypes.KVStore, k Keeper) error {
	iterator := storetypes.KVStorePrefixIterator(store, ApprovalTrackerKey)
	defer func() {
		if err := iterator.Close(); err != nil {
			k.Logger().Error("failed to close approval tracker migration iterator", "error", err)
		}
	}()

	for ; iterator.Valid(); iterator.Next() {
		var oldApprovalTracker oldtypes.ApprovalTracker
		k.cdc.MustUnmarshal(iterator.Value(), &oldApprovalTracker)

		// Convert to JSON
		jsonBytes, err := json.Marshal(oldApprovalTracker)
		if err != nil {
			return err
		}

		// Unmarshal into new type
		var newApprovalTracker newtypes.ApprovalTracker
		if err := json.Unmarshal(jsonBytes, &newApprovalTracker); err != nil {
			return err
		}

		store.Set(iterator.Key(), k.cdc.MustMarshal(&newApprovalTracker))
	}

	return nil
}

func MigrateDynamicStores(ctx context.Context, store storetypes.KVStore, k Keeper) error {
	// Migrate base dynamic stores
	iterator := storetypes.KVStorePrefixIterator(store, DynamicStoreKey)
	defer func() {
		if err := iterator.Close(); err != nil {
			k.Logger().Error("failed to close dynamic store migration iterator", "error", err)
		}
	}()
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
	defer func() {
		if err := valueIterator.Close(); err != nil {
			k.Logger().Error("failed to close dynamic store value migration iterator", "error", err)
		}
	}()
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
