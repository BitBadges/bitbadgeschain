package wasmx

import (
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/wasmx/keeper"
)

type BlockHandler struct {
	k keeper.Keeper
}

func NewBlockHandler(k keeper.Keeper) *BlockHandler {
	return &BlockHandler{
		k: k,
	}
}

func (h *BlockHandler) BeginBlocker(ctx sdk.Context, block abci.RequestBeginBlock) {
}

func (h *BlockHandler) EndBlocker(ctx sdk.Context, block abci.RequestEndBlock) {
}
