package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Queries how many times a leaf has been used for a challenge
func (k Keeper) GetNumUsedForMerkleChallenge(goCtx context.Context, req *types.QueryGetNumUsedForMerkleChallengeRequest) (*types.QueryGetNumUsedForMerkleChallengeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	numUsed, err := k.GetNumUsedForMerkleChallengeFromStore(ctx, req.CollectionId, req.ApproverAddress, req.ApprovalLevel, req.ChallengeTrackerId, req.LeafIndex)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	return &types.QueryGetNumUsedForMerkleChallengeResponse{
		NumUsed: numUsed,
	}, nil
}
