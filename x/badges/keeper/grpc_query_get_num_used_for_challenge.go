package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Queries how many times a leaf has been used for a challenge
func (k Keeper) GetNumUsedForChallenge(goCtx context.Context, req *types.QueryGetNumUsedForChallengeRequest) (*types.QueryGetNumUsedForChallengeResponse, error) {
	// if req == nil {
	// 	return nil, status.Error(codes.InvalidArgument, "invalid request")
	// }

	// ctx := sdk.UnwrapSDKContext(goCtx)

	// numUsed, err := k.GetNumUsedForChallengeFromStore(ctx, req.CollectionId, req.ChallengeId, req.LeafIndex, req.Level)
	// if err != nil {
	// 	return nil, status.Error(codes.InvalidArgument, "invalid request")
	// }

	// return &types.QueryGetNumUsedForChallengeResponse{
	// 	NumUsed: numUsed,
	// }, nil

	return nil, status.Error(codes.InvalidArgument, "invalid request")
}
