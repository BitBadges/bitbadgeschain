package types

import (
	"fmt"

	types "cosmossdk.io/math"
	host "github.com/cosmos/ibc-go/v11/modules/core/24-host"
	// this line is used by starport scaffolding # genesis/types/import
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		PortId: PortID,
		// this line is used by starport scaffolding # genesis/types/default
		Params:                 DefaultParams(),
		NextCollectionId:       types.NewUint(1),
		NextDynamicStoreId:     types.NewUint(1),
		NextAddressListCounter: types.NewUint(0),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.

// IMPORTANT: We assume tokens are well-formed and validated here
func (gs GenesisState) Validate() error {
	if err := host.PortIdentifierValidator(gs.PortId); err != nil {
		return err
	}

	// Parallel arrays are indexed together in InitGenesis
	parallel := []struct {
		name string
		a, b int
	}{
		{"balances / balanceStoreKeys", len(gs.Balances), len(gs.BalanceStoreKeys)},
		{"challengeTrackers / challengeTrackerStoreKeys", len(gs.ChallengeTrackers), len(gs.ChallengeTrackerStoreKeys)},
		{"approvalTrackers / approvalTrackerStoreKeys", len(gs.ApprovalTrackers), len(gs.ApprovalTrackerStoreKeys)},
		{"approvalTrackerVersions / approvalTrackerVersionsStoreKeys", len(gs.ApprovalTrackerVersions), len(gs.ApprovalTrackerVersionsStoreKeys)},
		{"ethSignatureTrackers / ethSignatureTrackerStoreKeys", len(gs.EthSignatureTrackers), len(gs.EthSignatureTrackerStoreKeys)},
		{"votingTrackers / votingTrackerStoreKeys", len(gs.VotingTrackers), len(gs.VotingTrackerStoreKeys)},
		{"collectionStats / collectionStatsIds", len(gs.CollectionStats), len(gs.CollectionStatsIds)},
		{"votingChallengeTrackers / votingChallengeTrackerStoreKeys", len(gs.VotingChallengeTrackers), len(gs.VotingChallengeTrackerStoreKeys)},
	}
	for _, p := range parallel {
		if p.a != p.b {
			return fmt.Errorf("genesis %s length mismatch: %d vs %d", p.name, p.a, p.b)
		}
	}

	maxCollectionId := types.ZeroUint()
	seenCollections := map[string]bool{}
	for i, collection := range gs.Collections {
		if collection == nil {
			return fmt.Errorf("genesis collection at index %d is nil", i)
		}
		if collection.CollectionId.IsNil() || collection.CollectionId.IsZero() {
			return fmt.Errorf("genesis collection at index %d has an invalid id", i)
		}
		id := collection.CollectionId.String()
		if seenCollections[id] {
			return fmt.Errorf("duplicate genesis collection id %s", id)
		}
		seenCollections[id] = true
		if collection.CollectionId.GT(maxCollectionId) {
			maxCollectionId = collection.CollectionId
		}
	}
	if len(gs.Collections) > 0 && (gs.NextCollectionId.IsNil() || !gs.NextCollectionId.GT(maxCollectionId)) {
		return fmt.Errorf("nextCollectionId must be greater than the largest collection id %s", maxCollectionId.String())
	}

	maxStoreId := types.ZeroUint()
	seenStores := map[string]bool{}
	for i, store := range gs.DynamicStores {
		if store == nil {
			return fmt.Errorf("genesis dynamic store at index %d is nil", i)
		}
		if store.StoreId.IsNil() || store.StoreId.IsZero() {
			return fmt.Errorf("genesis dynamic store at index %d has an invalid id", i)
		}
		id := store.StoreId.String()
		if seenStores[id] {
			return fmt.Errorf("duplicate genesis dynamic store id %s", id)
		}
		seenStores[id] = true
		if store.StoreId.GT(maxStoreId) {
			maxStoreId = store.StoreId
		}
	}
	if len(gs.DynamicStores) > 0 && (gs.NextDynamicStoreId.IsNil() || !gs.NextDynamicStoreId.GT(maxStoreId)) {
		return fmt.Errorf("nextDynamicStoreId must be greater than the largest dynamic store id %s", maxStoreId.String())
	}

	seenLists := map[string]bool{}
	for i, list := range gs.AddressLists {
		if list == nil {
			return fmt.Errorf("genesis address list at index %d is nil", i)
		}
		if seenLists[list.ListId] {
			return fmt.Errorf("duplicate genesis address list id %s", list.ListId)
		}
		seenLists[list.ListId] = true
	}

	seenBalanceKeys := map[string]bool{}
	for _, key := range gs.BalanceStoreKeys {
		if seenBalanceKeys[key] {
			return fmt.Errorf("duplicate genesis balance store key %s", key)
		}
		seenBalanceKeys[key] = true
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
