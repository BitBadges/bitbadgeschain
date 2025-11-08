package badges

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"

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
		if err := k.SetCollectionInStore(ctx, collection); err != nil {
			panic(err)
		}
	}

	for idx, balance := range genState.Balances {
		if err := k.SetUserBalanceInStore(ctx, genState.BalanceStoreKeys[idx], balance); err != nil {
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
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.NextCollectionId = k.GetNextCollectionId(ctx)
	genesis.NextDynamicStoreId = k.GetNextDynamicStoreId(ctx)
	genesis.PortId = k.GetPort(ctx)

	genesis.Collections = k.GetCollectionsFromStore(ctx)
	addresses := []string{}
	balanceIds := []sdkmath.Uint{}
	genesis.Balances, addresses, balanceIds = k.GetUserBalancesFromStore(ctx)

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

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
