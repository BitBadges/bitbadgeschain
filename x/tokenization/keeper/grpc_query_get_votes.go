package keeper

import (
	"context"
	"strings"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetVotes queries all votes for a proposal
func (k Keeper) GetVotes(goCtx context.Context, req *types.QueryGetVotesRequest) (*types.QueryGetVotesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	collectionId := sdkmath.NewUintFromString(req.CollectionId)

	// Construct the prefix for all votes for this proposal
	// The key format is: collectionId-approverAddress-approvalLevel-approvalId-proposalId-voterAddress
	// So we need to find all keys that start with: collectionId-approverAddress-approvalLevel-approvalId-proposalId-
	proposalPrefix := strings.Join([]string{
		collectionId.String(),
		req.ApproverAddress,
		req.ApprovalLevel,
		req.ApprovalId,
		req.ProposalId,
	}, BalanceKeyDelimiter) + BalanceKeyDelimiter

	// Get all votes from store
	allVotes, allKeys := k.GetVotesFromStore(ctx)

	// Filter votes that match the proposal prefix
	var matchingVotes []*types.VoteProof
	for i, key := range allKeys {
		if strings.HasPrefix(key, proposalPrefix) {
			// Extract voter address from key (it's the last part after the prefix)
			voterAddress := strings.TrimPrefix(key, proposalPrefix)
			// Verify the vote matches the proposal
			if allVotes[i] != nil && allVotes[i].ProposalId == req.ProposalId && allVotes[i].Voter == voterAddress {
				matchingVotes = append(matchingVotes, allVotes[i])
			}
		}
	}

	return &types.QueryGetVotesResponse{
		Votes: matchingVotes,
	}, nil
}
