package keeper

import (
	"context"

	"bitbadgeschain/x/badges/types"

	sdkmath "cosmossdk.io/math"
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
	collectionId := sdkmath.NewUintFromString(req.CollectionId)
	leafIndex := sdkmath.NewUintFromString(req.LeafIndex)
	numUsed, err := k.GetChallengeTrackerFromStore(ctx, collectionId, req.ApproverAddress, req.ApprovalLevel, req.ApprovalId, req.ChallengeTrackerId, leafIndex)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	return &types.QueryGetChallengeTrackerResponse{
		NumUsed: numUsed.String(),
	}, nil
}
