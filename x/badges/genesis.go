package badges

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.

//NOTE: We assume that all badges are validly formed here
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set if defined; default 0
	k.SetNextBadgeId(ctx, genState.NextBadgeId)
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

	for _, badge := range genState.Badges {
		if err := k.SetBadgeInStore(ctx, *badge); err != nil {
			panic(err)
		}
	}

	for idx, balance := range genState.Balances {
		if err := k.SetUserBalanceInStore(ctx, genState.BalanceIds[idx], *balance); err != nil {
			panic(err)
		}
	}

	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.PortId = k.GetPort(ctx)
	genesis.NextBadgeId = k.GetNextBadgeId(ctx)

	genesis.Badges = k.GetBadgesFromStore(ctx)
	genesis.Balances = k.GetUserBalancesFromStore(ctx)
	genesis.BalanceIds = k.GetUserBalanceIdsFromStore(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
