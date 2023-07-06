package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Queries a balance for the given address and badgeId and returns its contents.
func (k Keeper) GetApprovalsTracker(goCtx context.Context, req *types.QueryGetApprovalsTrackerRequest) (*types.QueryGetApprovalsTrackerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	address, found := k.GetApprovalsTrackerFromStore(ctx, req.CollectionId, req.TrackerId, req.Level, req.Depth, req.Address)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	return &types.QueryGetApprovalsTrackerResponse{
		Tracker: &address,
	}, nil
}
