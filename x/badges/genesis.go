package badges

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.

// NOTE: We assume that all badges are validly formed here
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set if defined; default 0
	if genState.NextCollectionId == 0 {
		genState.NextCollectionId = 1
	}
	if genState.NextClaimId == 0 {
		genState.NextClaimId = 1
	}

	k.SetNextCollectionId(ctx, genState.NextCollectionId)
	// Set if defined; default 0
	k.SetNextClaimId(ctx, genState.NextCollectionId)
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

	for _, badge := range genState.Collections {
		if err := k.SetCollectionInStore(ctx, *badge); err != nil {
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
	genesis.NextCollectionId = k.GetNextCollectionId(ctx)
	genesis.NextClaimId = k.GetNextClaimId(ctx)

	genesis.Collections = k.GetCollectionsFromStore(ctx)
	accNums := []uint64{}
	ids := []uint64{}
	genesis.Balances, accNums, ids = k.GetUserBalancesFromStore(ctx)

	for i, accNum := range accNums {
		genesis.BalanceIds = append(genesis.BalanceIds, keeper.ConstructBalanceKey(accNum, ids[i]))
	}
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
