package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/anchor/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetValueAtLocation(goCtx context.Context, req *types.QueryGetValueAtLocationRequest) (*types.QueryGetValueAtLocationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val := k.GetAnchorLocation(ctx, req.LocationId)
	return &types.QueryGetValueAtLocationResponse{AnchorData: val}, nil
}
