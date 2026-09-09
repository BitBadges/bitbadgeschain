package keeper

import (
	"context"
	"strings"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

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
	collectionId, err := parseQueryUint(req.CollectionId, "CollectionId")
	if err != nil {
		return nil, err
	}

	legacyProposalPrefix := strings.Join([]string{
		collectionId.String(),
		req.ApproverAddress,
		req.ApprovalLevel,
		req.ApprovalId,
		req.ProposalId,
	}, BalanceKeyDelimiter) + BalanceKeyDelimiter
	v2ProposalPrefix := constructV2TrackerKey(
		collectionId,
		req.ApproverAddress,
		req.ApprovalLevel,
		req.ApprovalId,
		req.ProposalId,
	)

	// Get all votes from store
	allVotes, allKeys := k.GetVotesFromStore(ctx)

	var matchingVotes []*types.VoteProof
	v2Voters := map[string]bool{}
	for i, key := range allKeys {
		if !strings.HasPrefix(key, v2ProposalPrefix) {
			continue
		}
		_, fields, ok := decodeV2TrackerKey(key, 5)
		if !ok || fields[0] != req.ApproverAddress || fields[1] != req.ApprovalLevel || fields[2] != req.ApprovalId || fields[3] != req.ProposalId {
			continue
		}
		voterAddress := fields[4]
		if allVotes[i] != nil && allVotes[i].ProposalId == req.ProposalId && allVotes[i].Voter == voterAddress {
			matchingVotes = append(matchingVotes, allVotes[i])
			v2Voters[voterAddress] = true
		}
	}

	// Legacy entries remain readable until a voter next writes this proposal.
	for i, key := range allKeys {
		if strings.HasPrefix(key, legacyProposalPrefix) {
			voterAddress := strings.TrimPrefix(key, legacyProposalPrefix)
			// Verify the vote matches the proposal
			if !v2Voters[voterAddress] && allVotes[i] != nil && allVotes[i].ProposalId == req.ProposalId && allVotes[i].Voter == voterAddress {
				matchingVotes = append(matchingVotes, allVotes[i])
			}
		}
	}

	return &types.QueryGetVotesResponse{
		Votes: matchingVotes,
	}, nil
}
