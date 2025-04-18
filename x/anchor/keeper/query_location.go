package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/anchor/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdkmath "cosmossdk.io/math"
)

func (k Keeper) GetValueAtLocation(goCtx context.Context, req *types.QueryGetValueAtLocationRequest) (*types.QueryGetValueAtLocationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	locationId := sdkmath.NewUintFromString(req.LocationId)
	val := k.GetAnchorLocation(ctx, locationId)
	return &types.QueryGetValueAtLocationResponse{AnchorData: val}, nil
}
