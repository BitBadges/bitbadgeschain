package keeper

import (
	"context"
	"encoding/json"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	storetypes "cosmossdk.io/store/types"
	newtypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
	oldtypes "github.com/bitbadges/bitbadgeschain/x/badges/types/v10"
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

func MigrateIncomingApprovals(incomingApprovals []*newtypes.UserIncomingApproval) []*newtypes.UserIncomingApproval {
	// for _, approval := range incomingApprovals {
	// 	if approval.ApprovalCriteria == nil {
	// 		continue
	// 	}

	// 	if approval.ApprovalCriteria.AutoDeletionOptions == nil {
	// 		continue
	// 	}

	// 	approval.ApprovalCriteria.AutoDeletionOptions = &newtypes.AutoDeletionOptions{
	// 		AfterOneUse:                 approval.ApprovalCriteria.AutoDeletionOptions.AfterOneUse,
	// 		AfterOverallMaxNumTransfers: false,
	// 	}
	// }

	return incomingApprovals
}

func MigrateOutgoingApprovals(outgoingApprovals []*newtypes.UserOutgoingApproval) []*newtypes.UserOutgoingApproval {
	// for _, approval := range outgoingApprovals {
	// 	if approval.ApprovalCriteria == nil {
	// 		continue
	// 	}

	// 	if approval.ApprovalCriteria.AutoDeletionOptions == nil {
	// 		continue
	// 	}

	// 	approval.ApprovalCriteria.AutoDeletionOptions = &newtypes.AutoDeletionOptions{
	// 		AfterOneUse:                 approval.ApprovalCriteria.AutoDeletionOptions.AfterOneUse,
	// 		AfterOverallMaxNumTransfers: false,
	// 	}
	// }

	return outgoingApprovals
}

func MigrateApprovals(collectionApprovals []*newtypes.CollectionApproval) []*newtypes.CollectionApproval {
	// for _, approval := range collectionApprovals {
	// 	if approval.ApprovalCriteria == nil {
	// 		continue
	// 	}

	// 	if approval.ApprovalCriteria.AutoDeletionOptions == nil {
	// 		continue
	// 	}

	// 	approval.ApprovalCriteria.AutoDeletionOptions = &newtypes.AutoDeletionOptions{
	// 		AfterOneUse:                 approval.ApprovalCriteria.AutoDeletionOptions.AfterOneUse,
	// 		AfterOverallMaxNumTransfers: false,
	// 	}
	// }

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

func MigrateCollections(ctx sdk.Context, store storetypes.KVStore, k Keeper) error {
	iterator := storetypes.KVStorePrefixIterator(store, CollectionKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var oldCollection oldtypes.BadgeCollection
		k.cdc.MustUnmarshal(iterator.Value(), &oldCollection)

		// Convert to JSON
		jsonBytes, err := json.Marshal(oldCollection)
		if err != nil {
			return err
		}

		// Unmarshal into new type
		var newCollection newtypes.BadgeCollection
		if err := json.Unmarshal(jsonBytes, &newCollection); err != nil {
			return err
		}

		// Set all approval versions to 0
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
	// iterator := storetypes.KVStorePrefixIterator(store, UserBalanceKey)
	// defer iterator.Close()

	// for ; iterator.Valid(); iterator.Next() {
	// 	var UserBalance oldtypes.UserBalanceStore
	// 	k.cdc.MustUnmarshal(iterator.Value(), &UserBalance)

	// 	// Convert to JSON
	// 	jsonBytes, err := json.Marshal(UserBalance)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// Unmarshal into old type
	// 	var oldBalance newtypes.UserBalanceStore
	// 	if err := json.Unmarshal(jsonBytes, &oldBalance); err != nil {
	// 		return err
	// 	}

	// 	oldBalance.IncomingApprovals = MigrateIncomingApprovals(oldBalance.IncomingApprovals)
	// 	oldBalance.OutgoingApprovals = MigrateOutgoingApprovals(oldBalance.OutgoingApprovals)

	// 	store.Set(iterator.Key(), k.cdc.MustMarshal(&oldBalance))
	// }

	return nil
}

func MigrateAddressLists(ctx context.Context, store storetypes.KVStore, k Keeper) error {
	// iterator := storetypes.KVStorePrefixIterator(store, AddressListKey)
	// defer iterator.Close()

	// for ; iterator.Valid(); iterator.Next() {
	// 	var AddressList oldtypes.AddressList
	// 	k.cdc.MustUnmarshal(iterator.Value(), &AddressList)

	// 	// Convert to JSON
	// 	jsonBytes, err := json.Marshal(AddressList)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// Unmarshal into old type
	// 	var oldAddressList newtypes.AddressList
	// 	if err := json.Unmarshal(jsonBytes, &oldAddressList); err != nil {
	// 		return err
	// 	}

	// 	store.Set(iterator.Key(), k.cdc.MustMarshal(&oldAddressList))
	// }
	return nil
}

func MigrateApprovalTrackers(ctx context.Context, store storetypes.KVStore, k Keeper) error {
	// iterator := storetypes.KVStorePrefixIterator(store, ApprovalTrackerKey)
	// defer iterator.Close()

	// for ; iterator.Valid(); iterator.Next() {
	// 	var ApprovalTracker oldtypes.ApprovalTracker
	// 	k.cdc.MustUnmarshal(iterator.Value(), &ApprovalTracker)

	// 	// Convert to JSON
	// 	jsonBytes, err := json.Marshal(ApprovalTracker)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// Unmarshal into old type
	// 	var oldApprovalTracker newtypes.ApprovalTracker
	// 	if err := json.Unmarshal(jsonBytes, &oldApprovalTracker); err != nil {
	// 		return err
	// 	}

	// 	wctx := sdk.UnwrapSDKContext(ctx)
	// 	nowUnixMilli := wctx.BlockTime().UnixMilli()
	// 	oldApprovalTracker.LastUpdatedAt = sdkmath.NewUint(uint64(nowUnixMilli))

	// 	store.Set(iterator.Key(), k.cdc.MustMarshal(&oldApprovalTracker))
	// }
	return nil
}
