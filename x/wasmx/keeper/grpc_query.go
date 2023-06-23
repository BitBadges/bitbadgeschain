package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/wasmx/types"
)

var _ types.QueryServer = &Keeper{}

func (k *Keeper) WasmxParams(c context.Context, _ *types.QueryWasmxParamsRequest) (*types.QueryWasmxParamsResponse, error) {

	ctx := sdk.UnwrapSDKContext(c)

	params := k.GetParams(ctx)

	res := &types.QueryWasmxParamsResponse{
		Params: params,
	}
	return res, nil
}

func (k *Keeper) WasmxModuleState(c context.Context, _ *types.QueryModuleStateRequest) (*types.QueryModuleStateResponse, error) {

	ctx := sdk.UnwrapSDKContext(c)

	res := &types.QueryModuleStateResponse{
		State: k.ExportGenesis(ctx),
	}
	return res, nil
}
