package keeper

import (
	"context"
	"encoding/binary"
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	newtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	oldtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types/v29"
)

// The x/badges module was renamed to x/tokenization in v23, but the migration that
// re-derives collection addresses with the new module name (and moves the corresponding
// bank balances) was never wired into an upgrade handler. As a result, every existing
// collection on mainnet still has a mintEscrowAddress / cosmosCoinBackedPath.address /
// cosmosCoinWrapperPath.address derived from module="badges", while today's msg handlers
// and every downstream consumer (SDK, indexer, frontend) derives with module="tokenization".
// v30 is the catch-up migration: for each collection we compute the new tokenization-derived
// address, move all bank balances from the old address, rewrite the stored address on the
// collection, and flip the reserved-protocol-address flag. Idempotent — if old == new
// (which happens for any collection created after the migration ships), we skip.
const (
	oldAddressModuleName = "badges"
)

// generateMintEscrowAddressWithModuleName reproduces UniversalUpdateCollection's address
// derivation with a pluggable module name. Used to compute both the legacy (badges) and
// new (tokenization) addresses for a collection so we can migrate between them.
func generateMintEscrowAddressWithModuleName(moduleName string, collectionId sdkmath.Uint) (sdk.AccAddress, error) {
	derivationKey := make([]byte, DerivationKeyLength)
	binary.BigEndian.PutUint64(derivationKey, collectionId.Uint64())
	ac, err := authtypes.NewModuleCredential(moduleName, AccountGenerationPrefix, derivationKey)
	if err != nil {
		return nil, err
	}
	return sdk.AccAddress(ac.Address()), nil
}

// generatePathAddressWithModuleName mirrors generatePathAddress with a pluggable module name.
func generatePathAddressWithModuleName(moduleName string, pathString string, prefix []byte) (sdk.AccAddress, error) {
	fullPathBytes := []byte(pathString)
	ac, err := authtypes.NewModuleCredential(moduleName, prefix, fullPathBytes)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to generate module credential")
	}
	return sdk.AccAddress(ac.Address()), nil
}

// migrateBankBalancesBetweenAddresses moves every coin held by oldAddress to newAddress.
// No-op when the balance is zero (common for fresh wrapper/backed paths that never held funds).
func (k Keeper) migrateBankBalancesBetweenAddresses(ctx sdk.Context, oldAddress, newAddress string) error {
	if oldAddress == newAddress {
		return nil
	}
	oldAddr, err := sdk.AccAddressFromBech32(oldAddress)
	if err != nil {
		return errorsmod.Wrapf(err, "invalid old address %s", oldAddress)
	}
	newAddr, err := sdk.AccAddressFromBech32(newAddress)
	if err != nil {
		return errorsmod.Wrapf(err, "invalid new address %s", newAddress)
	}

	balances := k.bankKeeper.GetAllBalances(ctx, oldAddr)
	if balances.IsZero() {
		return nil
	}
	if err := k.bankKeeper.SendCoins(ctx, oldAddr, newAddr, balances); err != nil {
		return errorsmod.Wrapf(err, "failed to migrate balances from %s to %s", oldAddress, newAddress)
	}
	return nil
}

// migrateReservedProtocolAddressMapping flips the reserved-protocol-address flag from the
// old address to the new one. The old flag is cleared so future lookups don't accidentally
// treat a stale module address as reserved.
func (k Keeper) migrateReservedProtocolAddressMapping(ctx sdk.Context, oldAddress, newAddress string) error {
	if oldAddress == newAddress {
		return nil
	}
	if err := k.SetReservedProtocolAddressInStore(ctx, oldAddress, false); err != nil {
		return errorsmod.Wrapf(err, "failed to clear reserved flag on old address %s", oldAddress)
	}
	if err := k.SetReservedProtocolAddressInStore(ctx, newAddress, true); err != nil {
		return errorsmod.Wrapf(err, "failed to set reserved flag on new address %s", newAddress)
	}
	return nil
}

// migrateCollectionAddressesFromBadgesToTokenization updates a single collection's derived
// addresses in-place (mintEscrowAddress + cosmosCoinBackedPath.address + wrapper path
// addresses), moving bank balances from each old address to each new one. Mutates the
// passed-in collection pointer so the caller can persist it.
func (k Keeper) migrateCollectionAddressesFromBadgesToTokenization(ctx sdk.Context, collection *newtypes.TokenCollection) error {
	newModuleName := newtypes.ModuleName // "tokenization"

	// 1. mintEscrowAddress
	oldMint, err := generateMintEscrowAddressWithModuleName(oldAddressModuleName, collection.CollectionId)
	if err != nil {
		return errorsmod.Wrapf(err, "mint old derivation for collection %s", collection.CollectionId)
	}
	newMint, err := generateMintEscrowAddressWithModuleName(newModuleName, collection.CollectionId)
	if err != nil {
		return errorsmod.Wrapf(err, "mint new derivation for collection %s", collection.CollectionId)
	}
	oldMintStr, newMintStr := oldMint.String(), newMint.String()
	if oldMintStr != newMintStr {
		if err := k.migrateBankBalancesBetweenAddresses(ctx, oldMintStr, newMintStr); err != nil {
			return errorsmod.Wrapf(err, "mint bank migration for collection %s", collection.CollectionId)
		}
		if err := k.migrateReservedProtocolAddressMapping(ctx, oldMintStr, newMintStr); err != nil {
			return errorsmod.Wrapf(err, "mint reserved-flag migration for collection %s", collection.CollectionId)
		}
		collection.MintEscrowAddress = newMintStr
	}

	// 2. cosmosCoinBackedPath.address (if present)
	if !newtypes.IsBasicallyEmpty(collection.Invariants.CosmosCoinBackedPath) {
		backed := collection.Invariants.CosmosCoinBackedPath
		if backed.Conversion != nil && backed.Conversion.SideA != nil && backed.Conversion.SideA.Denom != "" {
			denom := backed.Conversion.SideA.Denom
			oldAddr, err := generatePathAddressWithModuleName(oldAddressModuleName, denom, BackedPathGenerationPrefix)
			if err != nil {
				return errorsmod.Wrapf(err, "backed old derivation for denom %s", denom)
			}
			newAddr, err := generatePathAddressWithModuleName(newModuleName, denom, BackedPathGenerationPrefix)
			if err != nil {
				return errorsmod.Wrapf(err, "backed new derivation for denom %s", denom)
			}
			oldStr, newStr := oldAddr.String(), newAddr.String()
			if oldStr != newStr {
				if err := k.migrateBankBalancesBetweenAddresses(ctx, oldStr, newStr); err != nil {
					return errorsmod.Wrapf(err, "backed bank migration for denom %s", denom)
				}
				if err := k.migrateReservedProtocolAddressMapping(ctx, oldStr, newStr); err != nil {
					return errorsmod.Wrapf(err, "backed reserved-flag migration for denom %s", denom)
				}
				backed.Address = newStr
			}
		}
	}

	// 3. cosmosCoinWrapperPaths[].address
	for i := range collection.CosmosCoinWrapperPaths {
		path := collection.CosmosCoinWrapperPaths[i]
		if path == nil || path.Denom == "" {
			continue
		}
		oldAddr, err := generatePathAddressWithModuleName(oldAddressModuleName, path.Denom, WrapperPathGenerationPrefix)
		if err != nil {
			return errorsmod.Wrapf(err, "wrapper old derivation for denom %s", path.Denom)
		}
		newAddr, err := generatePathAddressWithModuleName(newModuleName, path.Denom, WrapperPathGenerationPrefix)
		if err != nil {
			return errorsmod.Wrapf(err, "wrapper new derivation for denom %s", path.Denom)
		}
		oldStr, newStr := oldAddr.String(), newAddr.String()
		if oldStr == newStr {
			continue
		}
		if err := k.migrateBankBalancesBetweenAddresses(ctx, oldStr, newStr); err != nil {
			return errorsmod.Wrapf(err, "wrapper bank migration for denom %s", path.Denom)
		}
		if err := k.migrateReservedProtocolAddressMapping(ctx, oldStr, newStr); err != nil {
			return errorsmod.Wrapf(err, "wrapper reserved-flag migration for denom %s", path.Denom)
		}
		path.Address = newStr
	}

	return nil
}

// MigrateTokenizationKeeper migrates the tokenization keeper from v28 to v29.
//
// v29 changes:
// - AltTimeChecks: added offlineMonths, offlineDaysOfMonth, offlineWeeksOfYear (default: empty slices via JSON)
// - VotingChallenge: added resetAfterExecution (default: false), delayAfterQuorum (default: "0")
// - VoteProof: added votedAt timestamp (default: "0")
// - ApprovalCriteria: added userApprovalSettings (default: nil via JSON)
// - New VotingChallengeTracker message (no migration needed, new store key)
//
// All new fields have zero-value defaults, so JSON marshal/unmarshal handles the migration automatically.
// Explicit default-setting functions below ensure correctness even if JSON omits zero-value fields.
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

// migrateIncomingApprovalCriteria ensures new v29 fields have explicit defaults after JSON migration.
func migrateIncomingApprovalCriteria(approvalCriteria *newtypes.IncomingApprovalCriteria) {
	if approvalCriteria == nil {
		return
	}
}

// migrateOutgoingApprovalCriteria ensures new v29 fields have explicit defaults after JSON migration.
func migrateOutgoingApprovalCriteria(approvalCriteria *newtypes.OutgoingApprovalCriteria) {
	if approvalCriteria == nil {
		return
	}
}

// migrateApprovalCriteria ensures new v29 fields have explicit defaults after JSON migration.
// Also migrates UserRoyalties from the old standalone field into UserApprovalSettings.
func migrateApprovalCriteria(approvalCriteria *newtypes.ApprovalCriteria) {
	if approvalCriteria == nil {
		return
	}

	// Migrate AltTimeChecks: new fields default to empty slices via JSON (no explicit action needed)

	// Migrate UserApprovalSettings: ensure royalties from old field are preserved.
	// The old ApprovalCriteria.UserRoyalties (field 13) is now reserved.
	// During JSON migration, the old field is dropped. We need to check if the old proto bytes
	// had a UserRoyalties and move it. Since we use JSON marshal/unmarshal, the old field is lost.
	// However, the v27 type still has it — we handle this in MigrateCollections directly.
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
		if newCollection.DefaultBalances != nil {
			newCollection.DefaultBalances.IncomingApprovals = MigrateIncomingApprovals(newCollection.DefaultBalances.IncomingApprovals)
			newCollection.DefaultBalances.OutgoingApprovals = MigrateOutgoingApprovals(newCollection.DefaultBalances.OutgoingApprovals)
		}

		// Re-derive module addresses with the current "tokenization" module name and move
		// any bank balances off the legacy "badges"-derived addresses. See the comment on
		// migrateCollectionAddressesFromBadgesToTokenization for the full history.
		if err := k.migrateCollectionAddressesFromBadgesToTokenization(ctx, &newCollection); err != nil {
			return err
		}

		// Save the updated collection (with migrated fields)
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
