package keeper

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/council/types"
)

// proposalKey returns the store key for a proposal.
func proposalKey(councilId, proposalId uint64) []byte {
	key := make([]byte, 0, len(types.ProposalKeyPrefix)+8+1+8)
	key = append(key, types.ProposalKeyPrefix...)
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, councilId)
	key = append(key, buf...)
	key = append(key, '/')
	binary.BigEndian.PutUint64(buf, proposalId)
	key = append(key, buf...)
	return key
}

// voteKey returns the store key for a vote.
func voteKey(councilId, proposalId uint64, voter string) []byte {
	key := make([]byte, 0, len(types.VoteKeyPrefix)+8+1+8+1+len(voter))
	key = append(key, types.VoteKeyPrefix...)
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, councilId)
	key = append(key, buf...)
	key = append(key, '/')
	binary.BigEndian.PutUint64(buf, proposalId)
	key = append(key, buf...)
	key = append(key, '/')
	key = append(key, []byte(voter)...)
	return key
}

// GetProposal retrieves a proposal. Returns (proposal, found).
func (k Keeper) GetProposal(ctx sdk.Context, councilId, proposalId uint64) (types.Proposal, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(proposalKey(councilId, proposalId))
	if bz == nil {
		return types.Proposal{}, false
	}
	var proposal types.Proposal
	if err := json.Unmarshal(bz, &proposal); err != nil {
		panic(fmt.Sprintf("failed to unmarshal proposal %d/%d: %v", councilId, proposalId, err))
	}
	return proposal, true
}

// SetProposal stores a proposal.
func (k Keeper) SetProposal(ctx sdk.Context, proposal types.Proposal) {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(proposal)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal proposal %d/%d: %v", proposal.CouncilId, proposal.ProposalId, err))
	}
	store.Set(proposalKey(proposal.CouncilId, proposal.ProposalId), bz)
}

// GetVote retrieves a vote. Returns (vote, found).
func (k Keeper) GetVote(ctx sdk.Context, councilId, proposalId uint64, voter string) (types.Vote, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(voteKey(councilId, proposalId, voter))
	if bz == nil {
		return types.Vote{}, false
	}
	var vote types.Vote
	if err := json.Unmarshal(bz, &vote); err != nil {
		panic(fmt.Sprintf("failed to unmarshal vote %d/%d/%s: %v", councilId, proposalId, voter, err))
	}
	return vote, true
}

// SetVote stores a vote.
func (k Keeper) SetVote(ctx sdk.Context, vote types.Vote) {
	store := ctx.KVStore(k.storeKey)
	bz, err := json.Marshal(vote)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal vote: %v", err))
	}
	store.Set(voteKey(vote.CouncilId, vote.ProposalId, vote.Voter), bz)
}
