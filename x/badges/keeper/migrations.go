package keeper

import (
	"context"
	"encoding/json"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	storetypes "cosmossdk.io/store/types"
	v8types "github.com/bitbadges/bitbadgeschain/x/badges/types"
	v7types "github.com/bitbadges/bitbadgeschain/x/badges/types/v7"
)

// MigrateBadgesKeeper migrates the badges keeper to set all approval versions to 0
func (k Keeper) MigrateBadgesKeeper(ctx sdk.Context) error {

	// Get all collections
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

	return nil
}

func MigrateIncomingApprovals(incomingApprovals []*v8types.UserIncomingApproval) []*v8types.UserIncomingApproval {
	for _, approval := range incomingApprovals {
		if approval.ApprovalCriteria == nil {
			continue
		}

		if approval.ApprovalCriteria.AutoDeletionOptions == nil {
			continue
		}

		approval.ApprovalCriteria.AutoDeletionOptions = &v8types.AutoDeletionOptions{
			AfterOneUse:                 approval.ApprovalCriteria.AutoDeletionOptions.AfterOneUse,
			AfterOverallMaxNumTransfers: false,
		}
	}

	return incomingApprovals
}

func MigrateOutgoingApprovals(outgoingApprovals []*v8types.UserOutgoingApproval) []*v8types.UserOutgoingApproval {
	for _, approval := range outgoingApprovals {
		if approval.ApprovalCriteria == nil {
			continue
		}

		if approval.ApprovalCriteria.AutoDeletionOptions == nil {
			continue
		}

		approval.ApprovalCriteria.AutoDeletionOptions = &v8types.AutoDeletionOptions{
			AfterOneUse:                 approval.ApprovalCriteria.AutoDeletionOptions.AfterOneUse,
			AfterOverallMaxNumTransfers: false,
		}
	}

	return outgoingApprovals
}

func MigrateApprovals(collectionApprovals []*v8types.CollectionApproval) []*v8types.CollectionApproval {
	for _, approval := range collectionApprovals {
		if approval.ApprovalCriteria == nil {
			continue
		}

		if approval.ApprovalCriteria.AutoDeletionOptions == nil {
			continue
		}

		approval.ApprovalCriteria.AutoDeletionOptions = &v8types.AutoDeletionOptions{
			AfterOneUse:                 approval.ApprovalCriteria.AutoDeletionOptions.AfterOneUse,
			AfterOverallMaxNumTransfers: false,
		}
	}

	return collectionApprovals
}

func MigrateCollections(ctx sdk.Context, store storetypes.KVStore, k Keeper) error {
	iterator := storetypes.KVStorePrefixIterator(store, CollectionKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		// First unmarshal into v7 type
		var v7Collection v7types.BadgeCollection
		k.cdc.MustUnmarshal(iterator.Value(), &v7Collection)

		// Convert to JSON
		jsonBytes, err := json.Marshal(v7Collection)
		if err != nil {
			return err
		}

		// Unmarshal into v7 type
		var v8Collection v8types.BadgeCollection
		if err := json.Unmarshal(jsonBytes, &v8Collection); err != nil {
			return err
		}

		// Set all approval versions to 0
		v8Collection.CollectionApprovals = MigrateApprovals(v8Collection.CollectionApprovals)
		v8Collection.DefaultBalances.IncomingApprovals = MigrateIncomingApprovals(v8Collection.DefaultBalances.IncomingApprovals)
		v8Collection.DefaultBalances.OutgoingApprovals = MigrateOutgoingApprovals(v8Collection.DefaultBalances.OutgoingApprovals)

		for _, cosmosCoinWrapperPath := range v8Collection.CosmosCoinWrapperPaths {
			cosmosCoinWrapperPath.Symbol = cosmosCoinWrapperPath.Denom
			cosmosCoinWrapperPath.DenomUnits = []*v8types.DenomUnit{}
		}

		// Save the updated collection
		if err := k.SetCollectionInStore(ctx, &v8Collection); err != nil {
			return err
		}
	}

	return nil
}

func MigrateBalances(ctx context.Context, store storetypes.KVStore, k Keeper) error {
	// iterator := storetypes.KVStorePrefixIterator(store, UserBalanceKey)
	// defer iterator.Close()

	// for ; iterator.Valid(); iterator.Next() {
	// 	var UserBalance v7types.UserBalanceStore
	// 	k.cdc.MustUnmarshal(iterator.Value(), &UserBalance)

	// 	// Convert to JSON
	// 	jsonBytes, err := json.Marshal(UserBalance)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// Unmarshal into v7 type
	// 	var v7Balance v8types.UserBalanceStore
	// 	if err := json.Unmarshal(jsonBytes, &v7Balance); err != nil {
	// 		return err
	// 	}

	// 	v7Balance.IncomingApprovals = MigrateIncomingApprovals(v7Balance.IncomingApprovals)
	// 	v7Balance.OutgoingApprovals = MigrateOutgoingApprovals(v7Balance.OutgoingApprovals)

	// 	store.Set(iterator.Key(), k.cdc.MustMarshal(&v7Balance))
	// }

	return nil
}

func MigrateAddressLists(ctx context.Context, store storetypes.KVStore, k Keeper) error {
	// iterator := storetypes.KVStorePrefixIterator(store, AddressListKey)
	// defer iterator.Close()

	// for ; iterator.Valid(); iterator.Next() {
	// 	var AddressList v7types.AddressList
	// 	k.cdc.MustUnmarshal(iterator.Value(), &AddressList)

	// 	// Convert to JSON
	// 	jsonBytes, err := json.Marshal(AddressList)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// Unmarshal into v7 type
	// 	var v7AddressList v8types.AddressList
	// 	if err := json.Unmarshal(jsonBytes, &v7AddressList); err != nil {
	// 		return err
	// 	}

	// 	store.Set(iterator.Key(), k.cdc.MustMarshal(&v7AddressList))
	// }
	return nil
}

func MigrateApprovalTrackers(ctx context.Context, store storetypes.KVStore, k Keeper) error {
	// iterator := storetypes.KVStorePrefixIterator(store, ApprovalTrackerKey)
	// defer iterator.Close()

	// for ; iterator.Valid(); iterator.Next() {
	// 	var ApprovalTracker v7types.ApprovalTracker
	// 	k.cdc.MustUnmarshal(iterator.Value(), &ApprovalTracker)

	// 	// Convert to JSON
	// 	jsonBytes, err := json.Marshal(ApprovalTracker)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// Unmarshal into v7 type
	// 	var v7ApprovalTracker v8types.ApprovalTracker
	// 	if err := json.Unmarshal(jsonBytes, &v7ApprovalTracker); err != nil {
	// 		return err
	// 	}

	// 	wctx := sdk.UnwrapSDKContext(ctx)
	// 	nowUnixMilli := wctx.BlockTime().UnixMilli()
	// 	v7ApprovalTracker.LastUpdatedAt = sdkmath.NewUint(uint64(nowUnixMilli))

	// 	store.Set(iterator.Key(), k.cdc.MustMarshal(&v7ApprovalTracker))
	// }
	return nil
}
