package keeper

import (
	"context"

	"bitbadgeschain/x/maps/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Queries a balance for the given address and badgeId and returns its contents.
func (k Keeper) Map(goCtx context.Context, req *types.QueryGetMapRequest) (*types.QueryGetMapResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	fetchedMap, found := k.GetMapFromStore(ctx, req.MapId)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	return &types.QueryGetMapResponse{
		Map: fetchedMap,
	}, nil
}
