package protocols

import (
	"github.com/bitbadges/bitbadgeschain/x/protocols/keeper"
	"github.com/bitbadges/bitbadgeschain/x/protocols/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	k.SetPort(ctx, genState.PortId)
	// Only try to bind to port if it is not already bound, since we may already own
	// port capability from capability InitGenesis
	if !k.IsBound(ctx, genState.PortId) {
		// module binds to the port on InitChain
		// and claims the returned capability
		err := k.BindPort(ctx, genState.PortId)
		if err != nil {
			panic("could not claim port capability: " + err.Error())
		}
	}
	k.SetParams(ctx, genState.Params)

	for _, protocol := range genState.Protocols {
		if err := k.SetProtocolInStore(ctx, protocol); err != nil {
			panic(err)
		}
	}

	for idx, collectionId := range genState.CollectionIdsForProtocols {
		name, address := keeper.GetDetailsFromKey(genState.CollectionIdsForProtocolsKeys[idx])
		if err := k.SetProtocolCollectionInStore(ctx, name, address, collectionId); err != nil {
			panic(err)
		}
	}
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.PortId = k.GetPort(ctx)

	genesis.Protocols = k.GetProtocolsFromStore(ctx)

	names, addresses, collectionIds := k.GetProtocolCollectionsFromStore(ctx)
	genesis.CollectionIdsForProtocolsKeys = []string{}
	for idx, name := range names {
		genesis.CollectionIdsForProtocolsKeys = append(genesis.CollectionIdsForProtocolsKeys, keeper.ConstructCollectionIdForProtocolKey(name, addresses[idx]))
	}

	genesis.CollectionIdsForProtocols = collectionIds
	
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
