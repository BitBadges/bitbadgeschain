package tokenization

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkmath "cosmossdk.io/math"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set if defined; default 0
	if genState.NextCollectionId.Equal(sdkmath.NewUint(0)) {
		genState.NextCollectionId = sdkmath.NewUint(1)
	}
	k.SetNextCollectionId(ctx, genState.NextCollectionId)

	// Set next dynamic store ID if defined; default 0
	if genState.NextDynamicStoreId.Equal(sdkmath.NewUint(0)) {
		genState.NextDynamicStoreId = sdkmath.NewUint(1)
	}
	k.SetNextDynamicStoreId(ctx, genState.NextDynamicStoreId)

	// Set params
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}

	// Set port ID for IBC
	k.SetPort(ctx, genState.PortId)
	// Only try to bind to port if it is not already bound, since we may already own
	// port capability from capability InitGenesis
	if k.ShouldBound(ctx, genState.PortId) {
		// module binds to the port on InitChain
		// and claims the returned capability
		err := k.BindPort(ctx, genState.PortId)
		if err != nil {
			panic("could not claim port capability: " + err.Error())
		}
	}

	for _, collection := range genState.Collections {
		if err := k.SetCollectionInStore(ctx, collection, true); err != nil {
			panic(err)
		}
	}

	for idx, balance := range genState.Balances {
		if err := k.SetUserBalanceInStore(ctx, genState.BalanceStoreKeys[idx], balance, true); err != nil {
			panic(err)
		}
	}

	for idx, numUsed := range genState.ChallengeTrackers {
		if err := k.SetChallengeTrackerInStore(ctx, genState.ChallengeTrackerStoreKeys[idx], numUsed); err != nil {
			panic(err)
		}
	}

	for _, addressList := range genState.AddressLists {
		if err := k.SetAddressListInStore(ctx, *addressList); err != nil {
			panic(err)
		}
	}

	for idx, approvalTracker := range genState.ApprovalTrackers {
		if err := k.SetApprovalTrackerInStoreViaKey(ctx, genState.ApprovalTrackerStoreKeys[idx], *approvalTracker); err != nil {
			panic(err)
		}
	}

	for idx, version := range genState.ApprovalTrackerVersions {
		k.SetApprovalTrackerVersionInStore(ctx, genState.ApprovalTrackerVersionsStoreKeys[idx], version)
	}

	// Initialize dynamic stores
	for _, dynamicStore := range genState.DynamicStores {
		if err := k.SetDynamicStoreInStore(ctx, *dynamicStore); err != nil {
			panic(err)
		}
	}

	// Initialize dynamic store values
	for _, dynamicStoreValue := range genState.DynamicStoreValues {
		if err := k.SetDynamicStoreValueInStore(ctx, dynamicStoreValue.StoreId, dynamicStoreValue.Address, dynamicStoreValue.Value); err != nil {
			panic(err)
		}
	}

	// Initialize ETH signature trackers
	for idx, numUsed := range genState.EthSignatureTrackers {
		if err := k.SetETHSignatureTrackerInStore(ctx, genState.EthSignatureTrackerStoreKeys[idx], numUsed); err != nil {
			panic(err)
		}
	}

	// Initialize voting trackers
	for idx, vote := range genState.VotingTrackers {
		if err := k.SetVoteInStore(ctx, genState.VotingTrackerStoreKeys[idx], vote); err != nil {
			panic(err)
		}
	}

	// Initialize collection stats
	if len(genState.CollectionStats) != len(genState.CollectionStatsIds) {
		panic(fmt.Errorf("genesis collection stats and collection stats ids length mismatch: %d vs %d",
			len(genState.CollectionStats), len(genState.CollectionStatsIds)))
	}
	for idx, stats := range genState.CollectionStats {
		if err := k.SetCollectionStatsInStore(ctx, genState.CollectionStatsIds[idx], stats); err != nil {
			panic(err)
		}
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.NextCollectionId = k.GetNextCollectionId(ctx)
	genesis.NextDynamicStoreId = k.GetNextDynamicStoreId(ctx)
	genesis.PortId = k.GetPort(ctx)

	genesis.Collections = k.GetCollectionsFromStore(ctx)
	var addresses []string
	var balanceIds []sdkmath.Uint
	var err error
	genesis.Balances, addresses, balanceIds, err = k.GetUserBalancesFromStore(ctx)
	if err != nil {
		// Panic on genesis export errors to prevent corrupted state export
		// This ensures balances, ids, and addresses arrays remain synchronized
		panic(fmt.Errorf("failed to export user balances: %w", err))
	}

	for i, address := range addresses {
		genesis.BalanceStoreKeys = append(genesis.BalanceStoreKeys, keeper.ConstructBalanceKey(address, balanceIds[i]))
	}

	challengeTrackers, challengeTrackerStoreKeys := k.GetChallengeTrackersFromStore(ctx)
	genesis.ChallengeTrackers = challengeTrackers
	genesis.ChallengeTrackerStoreKeys = challengeTrackerStoreKeys
	genesis.AddressLists = k.GetAddressListsFromStore(ctx)
	genesis.ApprovalTrackers, genesis.ApprovalTrackerStoreKeys = k.GetApprovalTrackersFromStore(ctx)
	genesis.ApprovalTrackerVersions, genesis.ApprovalTrackerVersionsStoreKeys = k.GetApprovalTrackerVersionsFromStore(ctx)

	// Export dynamic stores
	genesis.DynamicStores = k.GetDynamicStoresFromStore(ctx)

	// Export dynamic store values
	genesis.DynamicStoreValues = k.GetAllDynamicStoreValuesFromStore(ctx)

	// Export ETH signature trackers
	genesis.EthSignatureTrackers, genesis.EthSignatureTrackerStoreKeys = k.GetETHSignatureTrackersFromStore(ctx)

	// Export voting trackers
	genesis.VotingTrackers, genesis.VotingTrackerStoreKeys = k.GetVotesFromStore(ctx)

	// Export collection stats
	stats, ids := k.GetAllCollectionStatsFromStore(ctx)
	genesis.CollectionStats = stats
	genesis.CollectionStatsIds = make([]types.Uint, len(ids))
	copy(genesis.CollectionStatsIds, ids)

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
