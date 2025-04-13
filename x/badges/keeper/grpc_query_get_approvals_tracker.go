package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdkmath "cosmossdk.io/math"
)

// Queries a balance for the given address and badgeId and returns its contents.
func (k Keeper) GetApprovalTracker(goCtx context.Context, req *types.QueryGetApprovalTrackerRequest) (*types.QueryGetApprovalTrackerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	collectionId := sdkmath.NewUintFromString(req.CollectionId)
	address, found := k.GetApprovalTrackerFromStore(ctx, collectionId, req.ApproverAddress, req.ApprovalId, req.AmountTrackerId, req.ApprovalLevel, req.TrackerType, req.ApprovedAddress)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	return &types.QueryGetApprovalTrackerResponse{
		Tracker: &address,
	}, nil
}
