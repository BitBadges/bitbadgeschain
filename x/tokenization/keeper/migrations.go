package keeper

import (
	"context"
	"encoding/binary"
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	newtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	oldtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types/v22"
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

func MigrateIncomingApprovals(incomingApprovals []*newtypes.UserIncomingApproval) []*newtypes.UserIncomingApproval {
	for _, approval := range incomingApprovals {
		if approval.ApprovalCriteria == nil {
			continue
		}
	}

	return incomingApprovals
}

func MigrateOutgoingApprovals(outgoingApprovals []*newtypes.UserOutgoingApproval) []*newtypes.UserOutgoingApproval {
	for _, approval := range outgoingApprovals {
		if approval.ApprovalCriteria == nil {
			continue
		}
	}

	return outgoingApprovals
}

func MigrateApprovals(collectionApprovals []*newtypes.CollectionApproval) []*newtypes.CollectionApproval {
	for _, approval := range collectionApprovals {
		if approval.ApprovalCriteria == nil {
			continue
		}
	}

	return collectionApprovals
}

// generateMintEscrowAddressWithModuleName generates mint escrow address with specified module name
func generateMintEscrowAddressWithModuleName(moduleName string, collectionId sdkmath.Uint) (sdk.AccAddress, error) {
	derivationKey := make([]byte, 8) // DerivationKeyLength = 8
	binary.BigEndian.PutUint64(derivationKey, collectionId.Uint64())
	ac, err := authtypes.NewModuleCredential(moduleName, AccountGenerationPrefix, derivationKey)
	if err != nil {
		return nil, err
	}
	return sdk.AccAddress(ac.Address()), nil
}

// generatePathAddressWithModuleName generates path address with specified module name
func generatePathAddressWithModuleName(moduleName string, pathString string, prefix []byte) (sdk.AccAddress, error) {
	fullPathBytes := []byte(pathString)
	ac, err := authtypes.NewModuleCredential(moduleName, prefix, fullPathBytes)
	if err != nil {
		return nil, err
	}
	return sdk.AccAddress(ac.Address()), nil
}

// migrateBankBalances migrates all bank balances from old address to new address
func (k Keeper) migrateBankBalances(ctx sdk.Context, oldAddress, newAddress string) error {
	oldAddr := sdk.MustAccAddressFromBech32(oldAddress)
	newAddr := sdk.MustAccAddressFromBech32(newAddress)

	// Get all balances from old address
	balances := k.bankKeeper.GetAllBalances(ctx, oldAddr)

	if !balances.IsZero() {
		// Transfer all balances to new address
		if err := k.bankKeeper.SendCoins(ctx, oldAddr, newAddr, balances); err != nil {
			return errorsmod.Wrapf(err, "failed to migrate balances from %s to %s", oldAddress, newAddress)
		}
	}

	return nil
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

		// Track address mappings for reserved protocol address updates
		addressMappings := make(map[string]string) // old -> new

		// 1. Migrate MintEscrowAddress
		if newCollection.MintEscrowAddress != "" {
			oldMintAddr, err := generateMintEscrowAddressWithModuleName(oldModuleName, newCollection.CollectionId)
			if err != nil {
				return errorsmod.Wrapf(err, "failed to generate old mint address for collection %s", newCollection.CollectionId)
			}

			newMintAddr, err := generateMintEscrowAddressWithModuleName(newModuleName, newCollection.CollectionId)
			if err != nil {
				return errorsmod.Wrapf(err, "failed to generate new mint address for collection %s", newCollection.CollectionId)
			}

			oldMintAddrStr := oldMintAddr.String()
			newMintAddrStr := newMintAddr.String()

			// Only migrate if addresses are different
			if oldMintAddrStr != newMintAddrStr {
				// Migrate bank balances
				if err := k.migrateBankBalances(ctx, oldMintAddrStr, newMintAddrStr); err != nil {
					return errorsmod.Wrapf(err, "failed to migrate mint escrow balances for collection %s", newCollection.CollectionId)
				}

				// Update collection
				newCollection.MintEscrowAddress = newMintAddrStr
				addressMappings[oldMintAddrStr] = newMintAddrStr
			}
		}

		// 2. Migrate CosmosCoinBackedPath address
		if newCollection.Invariants != nil && newCollection.Invariants.CosmosCoinBackedPath != nil {
			backedPath := newCollection.Invariants.CosmosCoinBackedPath
			if backedPath.Conversion != nil && backedPath.Conversion.SideA != nil {
				denom := backedPath.Conversion.SideA.Denom

				oldBackedAddr, err := generatePathAddressWithModuleName(oldModuleName, denom, BackedPathGenerationPrefix)
				if err != nil {
					return errorsmod.Wrapf(err, "failed to generate old backed path address for denom %s", denom)
				}

				newBackedAddr, err := generatePathAddressWithModuleName(newModuleName, denom, BackedPathGenerationPrefix)
				if err != nil {
					return errorsmod.Wrapf(err, "failed to generate new backed path address for denom %s", denom)
				}

				oldBackedAddrStr := oldBackedAddr.String()
				newBackedAddrStr := newBackedAddr.String()

				if oldBackedAddrStr != newBackedAddrStr {
					// Migrate bank balances
					if err := k.migrateBankBalances(ctx, oldBackedAddrStr, newBackedAddrStr); err != nil {
						return errorsmod.Wrapf(err, "failed to migrate backed path balances for denom %s", denom)
					}

					// Update collection
					backedPath.Address = newBackedAddrStr
					addressMappings[oldBackedAddrStr] = newBackedAddrStr
				}
			}
		}

		// 3. Migrate CosmosCoinWrapperPath addresses
		for i, wrapperPath := range newCollection.CosmosCoinWrapperPaths {
			if wrapperPath.Denom != "" {
				oldWrapperAddr, err := generatePathAddressWithModuleName(oldModuleName, wrapperPath.Denom, WrapperPathGenerationPrefix)
				if err != nil {
					return errorsmod.Wrapf(err, "failed to generate old wrapper path address for denom %s", wrapperPath.Denom)
				}

				newWrapperAddr, err := generatePathAddressWithModuleName(newModuleName, wrapperPath.Denom, WrapperPathGenerationPrefix)
				if err != nil {
					return errorsmod.Wrapf(err, "failed to generate new wrapper path address for denom %s", wrapperPath.Denom)
				}

				oldWrapperAddrStr := oldWrapperAddr.String()
				newWrapperAddrStr := newWrapperAddr.String()

				if oldWrapperAddrStr != newWrapperAddrStr {
					// Migrate bank balances
					if err := k.migrateBankBalances(ctx, oldWrapperAddrStr, newWrapperAddrStr); err != nil {
						return errorsmod.Wrapf(err, "failed to migrate wrapper path balances for denom %s", wrapperPath.Denom)
					}

					// Update collection
					newCollection.CosmosCoinWrapperPaths[i].Address = newWrapperAddrStr
					addressMappings[oldWrapperAddrStr] = newWrapperAddrStr
				}
			}
		}

		// 4. Update reserved protocol addresses
		for oldAddr, newAddr := range addressMappings {
			// Unset old address
			if err := k.SetReservedProtocolAddressInStore(ctx, oldAddr, false); err != nil {
				return errorsmod.Wrapf(err, "failed to unset old reserved protocol address %s", oldAddr)
			}

			// Set new address
			if err := k.SetReservedProtocolAddressInStore(ctx, newAddr, true); err != nil {
				return errorsmod.Wrapf(err, "failed to set new reserved protocol address %s", newAddr)
			}
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
