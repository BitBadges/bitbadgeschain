package keeper

import (
	"context"
	"encoding/json"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	storetypes "cosmossdk.io/store/types"

	v4types "github.com/bitbadges/bitbadgeschain/x/badges/types"
	v3types "github.com/bitbadges/bitbadgeschain/x/badges/types/v3"
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
		// First unmarshal into v3 type
		var v3Collection v3types.BadgeCollection
		k.cdc.MustUnmarshal(iterator.Value(), &v3Collection)

		// Convert to JSON
		jsonBytes, err := json.Marshal(v3Collection)
		if err != nil {
			return err
		}

		// Unmarshal into v4 type
		var v4Collection v4types.BadgeCollection
		if err := json.Unmarshal(jsonBytes, &v4Collection); err != nil {
			return err
		}

		// Set all approval versions to 0
		for _, approval := range v4Collection.CollectionApprovals {
			// correspondingv3Approval := &v3types.CollectionApproval{}
			// for _, v3Approval := range v3Collection.CollectionApprovals {
			// 	if v3Approval.ApprovalId == approval.ApprovalId {
			// 		correspondingv3Approval = v3Approval
			// 		break
			// 	}
			// }

			approval.ApprovalCriteria.AutoDeletionOptions = &v4types.AutoDeletionOptions{AfterOneUse: false}
		}

		for _, approval := range v4Collection.DefaultBalances.IncomingApprovals {
			// correspondingv3Approval := &v3types.UserIncomingApproval{}
			// for _, v3Approval := range v3Collection.DefaultBalances.IncomingApprovals {
			// 	if v3Approval.ApprovalId == approval.ApprovalId {
			// 		correspondingv3Approval = v3Approval
			// 		break
			// 	}
			// }

			approval.ApprovalCriteria.AutoDeletionOptions = &v4types.AutoDeletionOptions{AfterOneUse: false}
		}

		for _, approval := range v4Collection.DefaultBalances.OutgoingApprovals {
			// correspondingv3Approval := &v3types.UserOutgoingApproval{}
			// for _, v3Approval := range v3Collection.DefaultBalances.OutgoingApprovals {
			// 	if v3Approval.ApprovalId == approval.ApprovalId {
			// 		correspondingv3Approval = v3Approval
			// 		break
			// 	}
			// }

			approval.ApprovalCriteria.AutoDeletionOptions = &v4types.AutoDeletionOptions{AfterOneUse: false}
		}

		// Save the updated collection
		if err := k.SetCollectionInStore(ctx, &v4Collection); err != nil {
			return err
		}
	}

	return nil
}

func MigrateBalances(ctx context.Context, store storetypes.KVStore, k Keeper) error {
	iterator := storetypes.KVStorePrefixIterator(store, UserBalanceKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var UserBalance v3types.UserBalanceStore
		k.cdc.MustUnmarshal(iterator.Value(), &UserBalance)

		// Convert to JSON
		jsonBytes, err := json.Marshal(UserBalance)
		if err != nil {
			return err
		}

		// Unmarshal into v4 type
		var v4Balance v4types.UserBalanceStore
		if err := json.Unmarshal(jsonBytes, &v4Balance); err != nil {
			return err
		}

		for _, approval := range v4Balance.IncomingApprovals {
			// correspondingv3Approval := &v3types.UserIncomingApproval{}
			// for _, v3Approval := range UserBalance.IncomingApprovals {
			// 	if v3Approval.ApprovalId == approval.ApprovalId {
			// 		correspondingv3Approval = v3Approval
			// 		break
			// 	}
			// }
			approval.ApprovalCriteria.AutoDeletionOptions = &v4types.AutoDeletionOptions{AfterOneUse: false}
		}

		for _, approval := range v4Balance.OutgoingApprovals {
			// correspondingv3Approval := &v3types.UserOutgoingApproval{}
			// for _, v3Approval := range UserBalance.OutgoingApprovals {
			// 	if v3Approval.ApprovalId == approval.ApprovalId {
			// 		correspondingv3Approval = v3Approval
			// 		break
			// 	}
			// }
			approval.ApprovalCriteria.AutoDeletionOptions = &v4types.AutoDeletionOptions{AfterOneUse: false}
		}

		store.Set(iterator.Key(), k.cdc.MustMarshal(&v4Balance))
	}
	return nil
}

func MigrateAddressLists(ctx context.Context, store storetypes.KVStore, k Keeper) error {
	// iterator := storetypes.KVStorePrefixIterator(store, AddressListKey)
	// defer iterator.Close()

	// for ; iterator.Valid(); iterator.Next() {
	// 	var AddressList v3types.AddressList
	// 	k.cdc.MustUnmarshal(iterator.Value(), &AddressList)

	// 	// Convert to JSON
	// 	jsonBytes, err := json.Marshal(AddressList)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// Unmarshal into v4 type
	// 	var v4AddressList v4types.AddressList
	// 	if err := json.Unmarshal(jsonBytes, &v4AddressList); err != nil {
	// 		return err
	// 	}

	// 	store.Set(iterator.Key(), k.cdc.MustMarshal(&v4AddressList))
	// }
	return nil
}

func MigrateApprovalTrackers(ctx context.Context, store storetypes.KVStore, k Keeper) error {
	// iterator := storetypes.KVStorePrefixIterator(store, ApprovalTrackerKey)
	// defer iterator.Close()

	// for ; iterator.Valid(); iterator.Next() {
	// 	var ApprovalTracker v3types.ApprovalTracker
	// 	k.cdc.MustUnmarshal(iterator.Value(), &ApprovalTracker)

	// 	// Convert to JSON
	// 	jsonBytes, err := json.Marshal(ApprovalTracker)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// Unmarshal into v4 type
	// 	var v4ApprovalTracker v4types.ApprovalTracker
	// 	if err := json.Unmarshal(jsonBytes, &v4ApprovalTracker); err != nil {
	// 		return err
	// 	}

	// 	wctx := sdk.UnwrapSDKContext(ctx)
	// 	nowUnixMilli := wctx.BlockTime().UnixMilli()
	// 	v4ApprovalTracker.LastUpdatedAt = sdkmath.NewUint(uint64(nowUnixMilli))

	// 	store.Set(iterator.Key(), k.cdc.MustMarshal(&v4ApprovalTracker))
	// }
	return nil
}
