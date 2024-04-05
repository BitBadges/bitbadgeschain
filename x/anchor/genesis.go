package anchor

import (
	"github.com/bitbadges/bitbadgeschain/x/anchor/keeper"
	"github.com/bitbadges/bitbadgeschain/x/anchor/types"
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

	genState.NextLocationId = sdk.NewUint(uint64(len(genState.AnchorData) + 1))

	for i, anchor := range genState.AnchorData {
		k.SetAnchorLocation(ctx, sdk.NewUint(uint64(i+1)), anchor.Data, anchor.Creator)
	}
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.PortId = k.GetPort(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	if k.GetNextAnchorId(ctx).IsZero() {
		genesis.NextLocationId = sdk.NewUint(1)
	} else {
		genesis.NextLocationId = k.GetNextAnchorId(ctx)
	}

	if genesis.NextLocationId.IsZero() {
		genesis.NextLocationId = sdk.NewUint(1)
	}

	return genesis
}
