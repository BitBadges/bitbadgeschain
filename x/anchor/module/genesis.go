package anchor

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/anchor/keeper"
	"github.com/bitbadges/bitbadgeschain/x/anchor/types"
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

	genState.NextLocationId = sdkmath.NewUint(uint64(len(genState.AnchorData) + 1))

	for i, anchor := range genState.AnchorData {
		k.SetAnchorLocation(ctx, sdkmath.NewUint(uint64(i+1)), anchor.Data, anchor.Creator)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.PortId = k.GetPort(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	nextAnchorId, err := k.GetNextAnchorId(ctx)
	if err != nil {
		panic(err)
	}

	if nextAnchorId.IsZero() {
		genesis.NextLocationId = sdkmath.NewUint(1)
	} else {
		genesis.NextLocationId = nextAnchorId
	}

	if genesis.NextLocationId.IsZero() {
		genesis.NextLocationId = sdkmath.NewUint(1)
	}

	return genesis
}
