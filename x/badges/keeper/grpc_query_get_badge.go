package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Queries a badge by its ID and returns its contents.
func (k Keeper) GetBadge(goCtx context.Context, req *types.QueryGetBadgeRequest) (*types.QueryGetBadgeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	badge, found := k.GetBadgeFromStore(ctx, req.Id)
	if !found {
		return nil, ErrBadgeNotExists
	}

	return &types.QueryGetBadgeResponse{
		Badge: &badge,
	}, nil
}
