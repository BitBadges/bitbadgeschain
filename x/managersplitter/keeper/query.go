package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/managersplitter/types"

	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}

func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	// For now, return default params
	// In the future, if params are stored, retrieve them here
	return types.DefaultParams()
}

func (k Keeper) ManagerSplitter(c context.Context, req *types.QueryGetManagerSplitterRequest) (*types.QueryGetManagerSplitterResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	managerSplitter, found := k.GetManagerSplitterFromStore(ctx, req.Address)
	if !found {
		return nil, status.Error(codes.NotFound, "manager splitter not found")
	}

	return &types.QueryGetManagerSplitterResponse{ManagerSplitter: managerSplitter}, nil
}

func (k Keeper) AllManagerSplitters(c context.Context, req *types.QueryAllManagerSplittersRequest) (*types.QueryAllManagerSplittersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	managerSplitters := k.GetAllManagerSplittersFromStore(ctx)

	var pagination *query.PageResponse
	if req.Pagination != nil {
		// Simple pagination - return all for now
		// In a production system, implement proper pagination
		pagination = &query.PageResponse{
			Total: uint64(len(managerSplitters)),
		}
	}

	return &types.QueryAllManagerSplittersResponse{
		ManagerSplitters: managerSplitters,
		Pagination:        pagination,
	}, nil
}

