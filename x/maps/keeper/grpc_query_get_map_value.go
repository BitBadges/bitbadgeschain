package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/maps/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Queries a balance for the given address and tokenId and returns its contents.
func (k Keeper) MapValue(goCtx context.Context, req *types.QueryGetMapValueRequest) (*types.QueryGetMapValueResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	val := k.GetMapValueFromStore(ctx, req.MapId, req.Key)

	if val.Value == "" {
		currMap, found := k.GetMapFromStore(ctx, req.MapId)
		if !found {
			return nil, status.Error(codes.InvalidArgument, "invalid request")
		}

		val.Value = currMap.DefaultValue
	}

	return &types.QueryGetMapValueResponse{
		Value: &val,
	}, nil
}
