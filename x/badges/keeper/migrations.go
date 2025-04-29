package keeper

import (
	"context"
	"encoding/json"

	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	storetypes "cosmossdk.io/store/types"

	v1types "github.com/bitbadges/bitbadgeschain/x/badges/types"
	v0types "github.com/bitbadges/bitbadgeschain/x/badges/types/v0"
)

// MigrateBadgesKeeper migrates the badges keeper to set all approval versions to 0
func (k Keeper) MigrateBadgesKeeper(ctx sdk.Context) error {

	// Get all collections
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	// Unchanged: params, nextCollectionId, challengeTrackers, approvalTrackers

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

	return nil
}

func MigrateCollections(ctx sdk.Context, store storetypes.KVStore, k Keeper) error {
	iterator := storetypes.KVStorePrefixIterator(store, CollectionKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		// First unmarshal into v0 type
		var v0Collection v0types.BadgeCollection
		k.cdc.MustUnmarshal(iterator.Value(), &v0Collection)

		// Convert to JSON
		jsonBytes, err := json.Marshal(v0Collection)
		if err != nil {
			return err
		}

		// Unmarshal into v1 type
		var v1Collection v1types.BadgeCollection
		if err := json.Unmarshal(jsonBytes, &v1Collection); err != nil {
			return err
		}

		// Set all approval versions to 0
		for _, approval := range v1Collection.CollectionApprovals {
			k.IncrementApprovalVersion(ctx, v1Collection.CollectionId, "collection", "", approval.ApprovalId)
			approval.Version = sdkmath.NewUint(0)
			approval.ApprovalCriteria.PredeterminedBalances.IncrementedBalances.ApprovalDurationFromNow = sdkmath.NewUint(0)
		}

		for _, approval := range v1Collection.DefaultBalances.IncomingApprovals {
			k.IncrementApprovalVersion(ctx, v1Collection.CollectionId, "incoming", "", approval.ApprovalId)
			approval.ApprovalCriteria.PredeterminedBalances.IncrementedBalances.ApprovalDurationFromNow = sdkmath.NewUint(0)
			approval.Version = sdkmath.NewUint(0)
		}

		for _, approval := range v1Collection.DefaultBalances.OutgoingApprovals {
			k.IncrementApprovalVersion(ctx, v1Collection.CollectionId, "outgoing", "", approval.ApprovalId)
			approval.ApprovalCriteria.PredeterminedBalances.IncrementedBalances.ApprovalDurationFromNow = sdkmath.NewUint(0)
			approval.Version = sdkmath.NewUint(0)
		}

		// Save the updated collection
		if err := k.SetCollectionInStore(ctx, &v1Collection); err != nil {
			return err
		}
	}

	return nil
}

func MigrateBalances(ctx context.Context, store storetypes.KVStore, k Keeper) error {
	iterator := storetypes.KVStorePrefixIterator(store, UserBalanceKey)
	defer iterator.Close()

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	for ; iterator.Valid(); iterator.Next() {
		var UserBalance v0types.UserBalanceStore
		k.cdc.MustUnmarshal(iterator.Value(), &UserBalance)

		// Convert to JSON
		jsonBytes, err := json.Marshal(UserBalance)
		if err != nil {
			return err
		}

		// Unmarshal into v1 type
		var v1Balance v1types.UserBalanceStore
		if err := json.Unmarshal(jsonBytes, &v1Balance); err != nil {
			return err
		}

		balanceKeyDetails := GetDetailsFromBalanceKey(string(iterator.Key()[1:]))
		collectionId := balanceKeyDetails.collectionId
		approverAddress := balanceKeyDetails.address

		for _, approval := range v1Balance.IncomingApprovals {
			k.IncrementApprovalVersion(sdkCtx, collectionId, "incoming", approverAddress, approval.ApprovalId)
			approval.ApprovalCriteria.PredeterminedBalances.IncrementedBalances.ApprovalDurationFromNow = sdkmath.NewUint(0)
			approval.Version = sdkmath.NewUint(0)
		}

		for _, approval := range v1Balance.OutgoingApprovals {
			k.IncrementApprovalVersion(sdkCtx, collectionId, "outgoing", approverAddress, approval.ApprovalId)
			approval.ApprovalCriteria.PredeterminedBalances.IncrementedBalances.ApprovalDurationFromNow = sdkmath.NewUint(0)
			approval.Version = sdkmath.NewUint(0)
		}

		store.Set(iterator.Key(), k.cdc.MustMarshal(&v1Balance))
	}
	return nil
}

func MigrateAddressLists(ctx context.Context, store storetypes.KVStore, k Keeper) error {
	iterator := storetypes.KVStorePrefixIterator(store, AddressListKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var AddressList v0types.AddressList
		k.cdc.MustUnmarshal(iterator.Value(), &AddressList)

		// Convert to JSON
		jsonBytes, err := json.Marshal(AddressList)
		if err != nil {
			return err
		}

		// Unmarshal into v1 type
		var v1AddressList v1types.AddressList
		if err := json.Unmarshal(jsonBytes, &v1AddressList); err != nil {
			return err
		}

		store.Set(iterator.Key(), k.cdc.MustMarshal(&v1AddressList))
	}
	return nil
}

func MigrateApprovalTrackers(ctx context.Context, store storetypes.KVStore, k Keeper) error {
	iterator := storetypes.KVStorePrefixIterator(store, ApprovalTrackerKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var ApprovalTracker v0types.ApprovalTracker
		k.cdc.MustUnmarshal(iterator.Value(), &ApprovalTracker)

		// Convert to JSON
		jsonBytes, err := json.Marshal(ApprovalTracker)
		if err != nil {
			return err
		}

		// Unmarshal into v1 type
		var v1ApprovalTracker v1types.ApprovalTracker
		if err := json.Unmarshal(jsonBytes, &v1ApprovalTracker); err != nil {
			return err
		}

		store.Set(iterator.Key(), k.cdc.MustMarshal(&v1ApprovalTracker))
	}
	return nil
}
