package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

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
	collectionId, err := parseQueryUint(req.CollectionId, "CollectionId")
	if err != nil {
		return nil, err
	}
	leafIndex, err := parseQueryUint(req.LeafIndex, "LeafIndex")
	if err != nil {
		return nil, err
	}
	numUsed, err := k.GetChallengeTrackerFromStore(ctx, collectionId, req.ApproverAddress, req.ApprovalLevel, req.ApprovalId, req.ChallengeTrackerId, leafIndex)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	return &types.QueryGetChallengeTrackerResponse{
		NumUsed: numUsed.String(),
	}, nil
}
