package managersplitter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkmath "cosmossdk.io/math"

	"github.com/bitbadges/bitbadgeschain/x/managersplitter/keeper"
	"github.com/bitbadges/bitbadgeschain/x/managersplitter/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set params
	// For now, params are empty, so we skip setting them

	// Set next manager splitter ID if defined; default to 1
	if genState.NextManagerSplitterId.Equal(sdkmath.NewUint(0)) {
		genState.NextManagerSplitterId = sdkmath.NewUint(1)
	}
	k.SetNextManagerSplitterId(ctx, genState.NextManagerSplitterId)

	// Initialize manager splitters
	for _, ms := range genState.ManagerSplitters {
		if err := k.SetManagerSplitterInStore(ctx, ms); err != nil {
			panic(err)
		}
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)
	genesis.NextManagerSplitterId = k.GetNextManagerSplitterId(ctx)
	genesis.ManagerSplitters = k.GetAllManagerSplittersFromStore(ctx)
	return genesis
}

