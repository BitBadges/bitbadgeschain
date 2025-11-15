package module

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/keeper"
	ibcratelimittypes "github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/types"
)

// InitGenesis initializes the module's state from a genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState *ibcratelimittypes.GenesisState) {
	if err := genState.Validate(); err != nil {
		panic(err)
	}
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis exports the module's state to a genesis state.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) ibcratelimittypes.GenesisState {
	return ibcratelimittypes.GenesisState{
		Params: k.GetParams(ctx),
	}
}
