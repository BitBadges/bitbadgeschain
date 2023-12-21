package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/protocols/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Queries a balance for the given address and badgeId and returns its contents.
func (k Keeper) GetProtocol(goCtx context.Context, req *types.QueryGetProtocolRequest) (*types.QueryGetProtocolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	protocol, found := k.GetProtocolFromStore(ctx, req.Name)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	return &types.QueryGetProtocolResponse{
		Protocol: *protocol,
	}, nil
}