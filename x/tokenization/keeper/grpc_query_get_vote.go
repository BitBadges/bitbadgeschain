package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetVote queries a vote by collection ID, approval level, approver address, approval ID, proposal ID, and voter address
func (k Keeper) GetVote(goCtx context.Context, req *types.QueryGetVoteRequest) (*types.QueryGetVoteResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	collectionId, err := parseQueryUint(req.CollectionId, "CollectionId")
	if err != nil {
		return nil, err
	}

	// Construct the vote key
	voteKey := ConstructVotingTrackerKey(
		collectionId,
		req.ApproverAddress,
		req.ApprovalLevel,
		req.ApprovalId,
		req.ProposalId,
		req.VoterAddress,
	)

	vote, found := k.GetVoteFromStore(ctx, voteKey)
	if !found {
		return nil, status.Error(codes.NotFound, "vote not found")
	}

	return &types.QueryGetVoteResponse{
		Vote: vote,
	}, nil
}
