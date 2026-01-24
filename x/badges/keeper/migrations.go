package keeper

import (
	"context"
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	newtypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
	oldtypes "github.com/bitbadges/bitbadgeschain/x/badges/types/v22"
)

// MigrateBadgesKeeper migrates the badges keeper from v21 to current version
func (k Keeper) MigrateBadgesKeeper(ctx sdk.Context) error {
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

	if err := MigrateUbadgeCoins(ctx, k); err != nil {
		return err
	}

	return nil
}

func MigrateIncomingApprovals(incomingApprovals []*newtypes.UserIncomingApproval) []*newtypes.UserIncomingApproval {
	for _, approval := range incomingApprovals {
		if approval.ApprovalCriteria == nil {
			continue
		}

		// Migrate DynamicStoreChallenge ownershipCheckParty to empty string (defaults to initiator)
		for _, challenge := range approval.ApprovalCriteria.DynamicStoreChallenges {
			if challenge != nil {
				challenge.OwnershipCheckParty = ""
			}
		}
	}

	return incomingApprovals
}

func MigrateOutgoingApprovals(outgoingApprovals []*newtypes.UserOutgoingApproval) []*newtypes.UserOutgoingApproval {
	for _, approval := range outgoingApprovals {
		if approval.ApprovalCriteria == nil {
			continue
		}

		// Migrate DynamicStoreChallenge ownershipCheckParty to empty string (defaults to initiator)
		for _, challenge := range approval.ApprovalCriteria.DynamicStoreChallenges {
			if challenge != nil {
				challenge.OwnershipCheckParty = ""
			}
		}
	}

	return outgoingApprovals
}

func MigrateApprovals(collectionApprovals []*newtypes.CollectionApproval) []*newtypes.CollectionApproval {
	for _, approval := range collectionApprovals {
		if approval.ApprovalCriteria == nil {
			// For backwards compatibility, create ApprovalCriteria with allowBackedMinting and allowSpecialWrapping set to true
			approval.ApprovalCriteria = &newtypes.ApprovalCriteria{
				AllowBackedMinting:   true,
				AllowSpecialWrapping: true,
			}
			continue
		}

		// Migrate DynamicStoreChallenge ownershipCheckParty to empty string (defaults to initiator)
		for _, challenge := range approval.ApprovalCriteria.DynamicStoreChallenges {
			if challenge != nil {
				challenge.OwnershipCheckParty = ""
			}
		}

		// For backwards compatibility, set allowBackedMinting and allowSpecialWrapping to true if not already set
		// This ensures existing approvals continue to work as they did before
		approval.ApprovalCriteria.AllowBackedMinting = true
		approval.ApprovalCriteria.AllowSpecialWrapping = true
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

// MigrateUbadgeCoins transfers all "ubadge" coins from one address to another
func MigrateUbadgeCoins(ctx sdk.Context, k Keeper) error {
	fromAddress := "bb1pkqancsm6lzmkz24x43ymjc8t8nykwye2vn2sx"
	toAddress := "bb1esaf3yswwkusfx68jusacfkfgs37j07rlrchj5"

	// Parse addresses
	fromAddr, err := sdk.AccAddressFromBech32(fromAddress)
	if err != nil {
		return err
	}

	toAddr, err := sdk.AccAddressFromBech32(toAddress)
	if err != nil {
		return err
	}

	// Get all balances for the source address
	allBalances := k.bankKeeper.GetAllBalances(ctx, fromAddr)

	// Get the "ubadge" balance amount
	ubadgeAmount := allBalances.AmountOf("ubadge")

	// If no ubadge balance found, nothing to migrate
	if ubadgeAmount.IsZero() {
		ctx.Logger().Info("No ubadge coins to migrate", "from", fromAddress)
		return nil
	}

	// Transfer all ubadge coins
	ubadgeCoin := sdk.NewCoin("ubadge", ubadgeAmount)
	coinsToTransfer := sdk.NewCoins(ubadgeCoin)
	if err := k.bankKeeper.SendCoins(ctx, fromAddr, toAddr, coinsToTransfer); err != nil {
		return err
	}

	ctx.Logger().Info("Migrated ubadge coins", "from", fromAddress, "to", toAddress, "amount", ubadgeAmount.String())
	return nil
}
