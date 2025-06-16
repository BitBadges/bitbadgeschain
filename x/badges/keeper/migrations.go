package keeper

import (
	"context"
	"encoding/json"

	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	storetypes "cosmossdk.io/store/types"

	v6types "github.com/bitbadges/bitbadgeschain/x/badges/types"
	v5types "github.com/bitbadges/bitbadgeschain/x/badges/types/v5"
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

func MigrateIncomingApprovals(incomingApprovals []*v6types.UserIncomingApproval) []*v6types.UserIncomingApproval {
	for _, approval := range incomingApprovals {
		if approval.ApprovalCriteria == nil {
			continue
		}

		approval.ApprovalCriteria.MustOwnBadges = []*v6types.MustOwnBadges{}
	}

	return incomingApprovals
}

func MigrateOutgoingApprovals(outgoingApprovals []*v6types.UserOutgoingApproval) []*v6types.UserOutgoingApproval {
	for _, approval := range outgoingApprovals {
		if approval.ApprovalCriteria == nil {
			continue
		}

		approval.ApprovalCriteria.MustOwnBadges = []*v6types.MustOwnBadges{}
	}

	return outgoingApprovals
}

func MigrateApprovals(collectionApprovals []*v6types.CollectionApproval) []*v6types.CollectionApproval {
	for _, approval := range collectionApprovals {
		if approval.ApprovalCriteria == nil {
			continue
		}

		approval.ApprovalCriteria.MustOwnBadges = []*v6types.MustOwnBadges{}

		approval.ApprovalCriteria.UserRoyalties = &v6types.UserRoyalties{
			PayoutAddress: "",
			Percentage:    sdkmath.NewUint(0),
		}
	}

	return collectionApprovals
}

func MigrateCollections(ctx sdk.Context, store storetypes.KVStore, k Keeper) error {
	iterator := storetypes.KVStorePrefixIterator(store, CollectionKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		// First unmarshal into v5 type
		var v5Collection v5types.BadgeCollection
		k.cdc.MustUnmarshal(iterator.Value(), &v5Collection)

		// Convert to JSON
		jsonBytes, err := json.Marshal(v5Collection)
		if err != nil {
			return err
		}

		// Unmarshal into v6 type
		var v6Collection v6types.BadgeCollection
		if err := json.Unmarshal(jsonBytes, &v6Collection); err != nil {
			return err
		}

		// Set all approval versions to 0
		v6Collection.CollectionApprovals = MigrateApprovals(v6Collection.CollectionApprovals)
		v6Collection.DefaultBalances.IncomingApprovals = MigrateIncomingApprovals(v6Collection.DefaultBalances.IncomingApprovals)
		v6Collection.DefaultBalances.OutgoingApprovals = MigrateOutgoingApprovals(v6Collection.DefaultBalances.OutgoingApprovals)

		// Save the updated collection
		if err := k.SetCollectionInStore(ctx, &v6Collection); err != nil {
			return err
		}
	}

	return nil
}

func MigrateBalances(ctx context.Context, store storetypes.KVStore, k Keeper) error {
	iterator := storetypes.KVStorePrefixIterator(store, UserBalanceKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var UserBalance v5types.UserBalanceStore
		k.cdc.MustUnmarshal(iterator.Value(), &UserBalance)

		// Convert to JSON
		jsonBytes, err := json.Marshal(UserBalance)
		if err != nil {
			return err
		}

		// Unmarshal into v6 type
		var v6Balance v6types.UserBalanceStore
		if err := json.Unmarshal(jsonBytes, &v6Balance); err != nil {
			return err
		}

		v6Balance.IncomingApprovals = MigrateIncomingApprovals(v6Balance.IncomingApprovals)
		v6Balance.OutgoingApprovals = MigrateOutgoingApprovals(v6Balance.OutgoingApprovals)

		store.Set(iterator.Key(), k.cdc.MustMarshal(&v6Balance))
	}
	return nil
}

func MigrateAddressLists(ctx context.Context, store storetypes.KVStore, k Keeper) error {
	// iterator := storetypes.KVStorePrefixIterator(store, AddressListKey)
	// defer iterator.Close()

	// for ; iterator.Valid(); iterator.Next() {
	// 	var AddressList v5types.AddressList
	// 	k.cdc.MustUnmarshal(iterator.Value(), &AddressList)

	// 	// Convert to JSON
	// 	jsonBytes, err := json.Marshal(AddressList)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// Unmarshal into v6 type
	// 	var v6AddressList v6types.AddressList
	// 	if err := json.Unmarshal(jsonBytes, &v6AddressList); err != nil {
	// 		return err
	// 	}

	// 	store.Set(iterator.Key(), k.cdc.MustMarshal(&v6AddressList))
	// }
	return nil
}

func MigrateApprovalTrackers(ctx context.Context, store storetypes.KVStore, k Keeper) error {
	// iterator := storetypes.KVStorePrefixIterator(store, ApprovalTrackerKey)
	// defer iterator.Close()

	// for ; iterator.Valid(); iterator.Next() {
	// 	var ApprovalTracker v5types.ApprovalTracker
	// 	k.cdc.MustUnmarshal(iterator.Value(), &ApprovalTracker)

	// 	// Convert to JSON
	// 	jsonBytes, err := json.Marshal(ApprovalTracker)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// Unmarshal into v6 type
	// 	var v6ApprovalTracker v6types.ApprovalTracker
	// 	if err := json.Unmarshal(jsonBytes, &v6ApprovalTracker); err != nil {
	// 		return err
	// 	}

	// 	wctx := sdk.UnwrapSDKContext(ctx)
	// 	nowUnixMilli := wctx.BlockTime().UnixMilli()
	// 	v6ApprovalTracker.LastUpdatedAt = sdkmath.NewUint(uint64(nowUnixMilli))

	// 	store.Set(iterator.Key(), k.cdc.MustMarshal(&v6ApprovalTracker))
	// }
	return nil
}
