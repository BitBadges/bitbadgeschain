package wasmx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	wasmxkeeper "github.com/bitbadges/bitbadgeschain/x/wasmx/keeper"
	"github.com/bitbadges/bitbadgeschain/x/wasmx/types"
)

func InitGenesis(ctx sdk.Context, keeper wasmxkeeper.Keeper, data types.GenesisState) {
	keeper.SetParams(ctx, data.Params)

}

func ExportGenesis(ctx sdk.Context, k wasmxkeeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		Params: k.GetParams(ctx),
	}
}
