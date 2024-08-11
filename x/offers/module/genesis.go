package offers

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbadgeschain/x/offers/keeper"
	"bitbadgeschain/x/offers/types"

	sdkmath "cosmossdk.io/math"
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

	// Set if defined; default 0
	if genState.NextProposalId.Equal(sdkmath.NewUint(0)) {
		genState.NextProposalId = sdkmath.NewUint(1)
	}
	k.SetNextProposalId(ctx, genState.NextProposalId)

	for _, proposal := range genState.Proposals {
		if err := k.SetProposalInStore(ctx, proposal); err != nil {
			panic(err)
		}
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.PortId = k.GetPort(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	genesis.NextProposalId = k.GetNextProposalId(ctx)

	genesis.Proposals = k.GetProposalsFromStore(ctx)
	return genesis
}
