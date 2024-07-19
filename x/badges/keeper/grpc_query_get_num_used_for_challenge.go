package keeper

import (
	"context"

	"bitbadgeschain/x/badges/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Queries how many times a leaf has been used for a challenge
func (k Keeper) GetChallengeTracker(goCtx context.Context, req *types.QueryGetChallengeTrackerRequest) (*types.QueryGetChallengeTrackerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	numUsed, err := k.GetChallengeTrackerFromStore(ctx, req.CollectionId, req.ApproverAddress, req.ApprovalLevel, req.ApprovalId, req.ChallengeTrackerId, req.LeafIndex)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	return &types.QueryGetChallengeTrackerResponse{
		NumUsed: numUsed,
	}, nil
}
