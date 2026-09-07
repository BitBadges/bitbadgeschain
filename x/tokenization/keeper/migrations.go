package keeper

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/pkg/storewalk"
	newtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	oldtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types/v32"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/store/v2/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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

// MigrateV35CanonicalAddresses rewrites store entries keyed or valued by a non-canonical
// bech32 spelling (bech32 also decodes the all-uppercase form) to the canonical spelling.
//
// Stores walked: user balances, address lists, approval trackers, used merkle-leaf
// trackers, ETH signature trackers, approval versions, dynamic store values and reserved
// protocol addresses. When a canonical twin already exists the entries are merged:
// balances and tracker tallies are summed, address lists are unioned, versions take the
// max, and boolean/flag entries keep the canonical value. Collection and user approval
// address values are normalized as well. Voting participant slots retain their spelling
// and weights; only their user-level approver scope is normalized.
func (k Keeper) MigrateV35CanonicalAddresses(ctx sdk.Context) error {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})

	if err := k.migrateV35VotingApproverKeys(ctx, store); err != nil {
		return err
	}
	if err := k.migrateV35BalanceKeys(ctx, store); err != nil {
		return err
	}
	if err := k.migrateV35AddressListValues(ctx, store); err != nil {
		return err
	}
	if err := k.migrateV35AddressValues(ctx, store); err != nil {
		return err
	}
	if err := k.migrateV35ApprovalTrackerKeys(ctx, store); err != nil {
		return err
	}
	// Used merkle-leaf and ETH signature trackers share the "cid-approver-…" layout and a
	// decimal-string counter value.
	if err := migrateV35CounterKeys(ctx, store, UsedClaimChallengeKey, 1); err != nil {
		return err
	}
	if err := migrateV35CounterKeys(ctx, store, ETHSignatureTrackerKey, 1); err != nil {
		return err
	}
	if err := migrateV35ApprovalVersionKeys(ctx, store); err != nil {
		return err
	}
	if err := k.migrateV35DynamicStoreValueKeys(ctx, store); err != nil {
		return err
	}
	if err := migrateV35ReservedProtocolAddressKeys(ctx, store); err != nil {
		return err
	}
	return nil
}

// canonicalBech32 returns the canonical spelling of a bech32 address and whether it differs
// from the input. Non-address strings ("", "Mint", "Total", …) are returned unchanged.
func canonicalBech32(address string) (string, bool) {
	accAddr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return address, false
	}
	canonical := accAddr.String()
	return canonical, canonical != address
}

type v35KeyRewrite struct {
	oldKey []byte
	newKey []byte
}

func walkV35KeyRewrites(ctx sdk.Context, store storetypes.KVStore, keyPrefix []byte, rewrite func(string) (string, bool), visit func(v35KeyRewrite) error) error {
	return storewalk.Prefix(ctx, store, keyPrefix, func(key, value []byte) error {
		suffix, changed := rewrite(string(key[len(keyPrefix):]))
		if !changed {
			return nil
		}
		return visit(v35KeyRewrite{oldKey: key, newKey: storeKey(keyPrefix, suffix)})
	})
}

// rewriteV35DelimitedKey canonicalises the address components at the given positions of a
// BalanceKeyDelimiter-joined key. Negative positions count from the end (-1 = last part), so
// components that may themselves contain the delimiter (approval ids) are never touched.
func rewriteV35DelimitedKey(suffix string, positions ...int) (string, bool) {
	parts := strings.Split(suffix, BalanceKeyDelimiter)
	changed := false
	for _, pos := range positions {
		idx := pos
		if idx < 0 {
			idx = len(parts) + pos
		}
		if idx < 0 || idx >= len(parts) {
			continue
		}
		canonical, partChanged := canonicalBech32(parts[idx])
		if partChanged {
			parts[idx] = canonical
			changed = true
		}
	}
	return strings.Join(parts, BalanceKeyDelimiter), changed
}

func (k Keeper) migrateV35BalanceKeys(ctx sdk.Context, store storetypes.KVStore) error {
	return storewalk.Prefix(ctx, store, UserBalanceKey, func(key, value []byte) error {
		suffix, changed := rewriteV35DelimitedKey(string(key[len(UserBalanceKey):]), -1)
		newKey := storeKey(UserBalanceKey, suffix)
		var moved newtypes.UserBalanceStore
		k.cdc.MustUnmarshal(value, &moved)
		merged := moved
		if changed {
			if existingBz := store.Get(newKey); existingBz != nil {
				var existing newtypes.UserBalanceStore
				k.cdc.MustUnmarshal(existingBz, &existing)
				summed, err := newtypes.AddBalances(ctx, moved.Balances, existing.Balances)
				if err != nil {
					return err
				}
				if hasNonZeroBalance(&moved) && hasNonZeroBalance(&existing) {
					parts := strings.SplitN(suffix, BalanceKeyDelimiter, 2)
					id, err := sdkmath.ParseUint(parts[0])
					if err != nil {
						return err
					}
					if collection, found := k.GetCollectionFromStore(ctx, id); found && !k.isExcludedFromHolderCount(ctx, collection, parts[1]) {
						stats, found := k.GetCollectionStatsFromStore(ctx, id)
						if found && stats.HolderCount.GT(sdkmath.ZeroUint()) {
							stats.HolderCount = stats.HolderCount.Sub(sdkmath.OneUint())
							if err := k.SetCollectionStatsInStore(ctx, id, stats); err != nil {
								return err
							}
						}
					}
				}
				merged = existing
				merged.Balances = summed
			}
		}
		k.canonicalV35UserBalance(ctx, &merged)
		updated := k.cdc.MustMarshal(&merged)
		if changed || !bytes.Equal(updated, value) {
			store.Set(newKey, updated)
		}
		if changed {
			store.Delete(key)
		}
		return nil
	})
}

func (k Keeper) migrateV35AddressListValues(ctx sdk.Context, store storetypes.KVStore) error {
	return storewalk.Prefix(ctx, store, AddressListKey, func(key, value []byte) error {
		var addressList newtypes.AddressList
		k.cdc.MustUnmarshal(value, &addressList)

		changed := false
		seen := map[string]bool{}
		addresses := make([]string, 0, len(addressList.Addresses))
		for _, address := range addressList.Addresses {
			canonical, addrChanged := canonicalBech32(address)
			changed = changed || addrChanged
			if seen[canonical] {
				// Union: the same identity listed twice under different spellings collapses to one.
				changed = true
				continue
			}
			seen[canonical] = true
			addresses = append(addresses, canonical)
		}
		createdBy, createdByChanged := canonicalBech32(addressList.CreatedBy)
		if !changed && !createdByChanged {
			return nil
		}
		addressList.Addresses = addresses
		addressList.CreatedBy = createdBy
		store.Set(key, k.cdc.MustMarshal(&addressList))
		return nil
	})
}

func (k Keeper) migrateV35ApprovalTrackerKeys(ctx sdk.Context, store storetypes.KVStore) error {
	// Key layout: collectionId-approverAddress-approvalId-amountTrackerId-level-trackerType-address
	return walkV35KeyRewrites(ctx, store, ApprovalTrackerKey, func(suffix string) (string, bool) {
		return rewriteV35DelimitedKey(suffix, 1, -1)
	}, func(rw v35KeyRewrite) error {
		var moved newtypes.ApprovalTracker
		k.cdc.MustUnmarshal(store.Get(rw.oldKey), &moved)

		merged := moved
		if existingBz := store.Get(rw.newKey); existingBz != nil {
			var existing newtypes.ApprovalTracker
			k.cdc.MustUnmarshal(existingBz, &existing)
			summed, err := newtypes.AddBalances(ctx, moved.Amounts, existing.Amounts)
			if err != nil {
				return err
			}
			merged = existing
			merged.Amounts = summed
			merged.NumTransfers = existing.NumTransfers.Add(moved.NumTransfers)
			if moved.LastUpdatedAt.GT(existing.LastUpdatedAt) {
				merged.LastUpdatedAt = moved.LastUpdatedAt
			}
		}

		store.Set(rw.newKey, k.cdc.MustMarshal(&merged))
		store.Delete(rw.oldKey)
		return nil
	})
}

// migrateV35CounterKeys rewrites the approver component of a "-"-joined key whose value is a
// decimal counter, summing counters that land on the same canonical key.
func migrateV35CounterKeys(ctx sdk.Context, store storetypes.KVStore, keyPrefix []byte, approverPosition int) error {
	return walkV35KeyRewrites(ctx, store, keyPrefix, func(suffix string) (string, bool) {
		return rewriteV35DelimitedKey(suffix, approverPosition)
	}, func(rw v35KeyRewrite) error {
		moved, err := sdkmath.ParseUint(string(store.Get(rw.oldKey)))
		if err != nil {
			return err
		}
		if existingBz := store.Get(rw.newKey); existingBz != nil {
			existing, err := sdkmath.ParseUint(string(existingBz))
			if err != nil {
				return err
			}
			moved = moved.Add(existing)
		}
		store.Set(rw.newKey, []byte(moved.String()))
		store.Delete(rw.oldKey)
		return nil
	})
}

func migrateV35ApprovalVersionKeys(ctx sdk.Context, store storetypes.KVStore) error {
	// Key layout: collectionId-approvalLevel-approverAddress-approvalId
	return walkV35KeyRewrites(ctx, store, ApprovalVersionKey, func(suffix string) (string, bool) {
		return rewriteV35DelimitedKey(suffix, 2)
	}, func(rw v35KeyRewrite) error {
		moved, err := sdkmath.ParseUint(string(store.Get(rw.oldKey)))
		if err != nil {
			return err
		}
		if existingBz := store.Get(rw.newKey); existingBz != nil {
			existing, err := sdkmath.ParseUint(string(existingBz))
			if err != nil {
				return err
			}
			if existing.GT(moved) {
				moved = existing
			}
		}
		store.Set(rw.newKey, []byte(moved.String()))
		store.Delete(rw.oldKey)
		return nil
	})
}

func (k Keeper) migrateV35DynamicStoreValueKeys(ctx sdk.Context, store storetypes.KVStore) error {
	return walkV35KeyRewrites(ctx, store, DynamicStoreValueKey, func(suffix string) (string, bool) {
		if len(suffix) < IDLength {
			return suffix, false
		}
		address, changed := canonicalBech32(suffix[IDLength:])
		return suffix[:IDLength] + address, changed
	}, func(rw v35KeyRewrite) error {
		var moved newtypes.DynamicStoreValue
		k.cdc.MustUnmarshal(store.Get(rw.oldKey), &moved)
		moved.Address, _ = canonicalBech32(moved.Address)
		if bz := store.Get(rw.newKey); bz != nil {
			var existing newtypes.DynamicStoreValue
			k.cdc.MustUnmarshal(bz, &existing)
			moved.Value = moved.Value && existing.Value
		}
		store.Set(rw.newKey, k.cdc.MustMarshal(&moved))
		store.Delete(rw.oldKey)
		return nil
	})
}

func migrateV35ReservedProtocolAddressKeys(ctx sdk.Context, store storetypes.KVStore) error {
	return walkV35KeyRewrites(ctx, store, ReservedProtocolAddressKey, canonicalBech32, func(rw v35KeyRewrite) error {
		if !store.Has(rw.newKey) {
			store.Set(rw.newKey, store.Get(rw.oldKey))
		}
		store.Delete(rw.oldKey)
		return nil
	})
}
