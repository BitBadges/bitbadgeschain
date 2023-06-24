package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetClaimNumProcessed(goCtx context.Context, req *types.QueryGetClaimNumProcessedRequest) (*types.QueryGetClaimNumProcessedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	//TODO:
	_ = ctx

	return nil, ErrNotImplemented

	return &types.QueryGetClaimNumProcessedResponse{
		NumProcessed: sdk.NewUint(0),
	}, nil
}
