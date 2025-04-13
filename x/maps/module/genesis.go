package maps

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/maps/keeper"
	"github.com/bitbadges/bitbadgeschain/x/maps/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
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
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}

	for _, mapToStore := range genState.Maps {
		if err := k.SetMapInStore(ctx, mapToStore); err != nil {
			panic(err)
		}
	}

	for idx, fullKey := range genState.FullKeys {
		mapId, key := keeper.GetDetailsFromKey(fullKey)
		if err := k.SetMapValueInStore(ctx, mapId, key, genState.Values[idx].Value, genState.Values[idx].LastSetBy); err != nil {
			panic(err)
		}
	}

	for _, mapId := range genState.DuplicatesFullKeys {
		mapId, value := keeper.GetDetailsFromKey(mapId)
		if err := k.SetMapDuplicateValueInStore(ctx, mapId, value); err != nil {
			panic(err)
		}
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.PortId = k.GetPort(ctx)

	genesis.Maps = k.GetMapsFromStore(ctx)

	mapIds, keys, values := k.GetMapKeysAndValuesFromStore(ctx)

	genesis.FullKeys = []string{}
	for idx, mapId := range mapIds {
		genesis.FullKeys = append(genesis.FullKeys, keeper.ConstructMapValueKey(mapId, keys[idx]))
	}

	genesis.Values = values

	mapIds, duplicateValues := k.GetMapDuplicateKeysAndValuesFromStore(ctx)

	genesis.DuplicatesFullKeys = []string{}
	for idx, mapId := range mapIds {
		genesis.DuplicatesFullKeys = append(genesis.DuplicatesFullKeys, keeper.ConstructMapValueDuplicatesKey(mapId, duplicateValues[idx]))
	}
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
