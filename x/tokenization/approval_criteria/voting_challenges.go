package approval_criteria

import (
	"fmt"
	"math/big"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// VotingService defines the interface for accessing vote data and constructing keys
type VotingService interface {
	GetVoteFromStore(ctx sdk.Context, key string) (*types.VoteProof, bool)
	GetVotingChallengeTrackerFromStore(ctx sdk.Context, key string) (*types.VotingChallengeTracker, bool)
	ConstructVoteKey(collectionId sdkmath.Uint, approverAddress, approvalLevel, approvalId, proposalId, voterAddress string) string
	ConstructChallengeTrackerKey(collectionId sdkmath.Uint, approverAddress, approvalLevel, approvalId, proposalId string) string
}

// VotingChallengesChecker implements ApprovalCriteriaChecker for VotingChallenges
type VotingChallengesChecker struct {
	votingService VotingService
}

// NewVotingChallengesChecker creates a new VotingChallengesChecker
func NewVotingChallengesChecker(votingService VotingService) *VotingChallengesChecker {
	return &VotingChallengesChecker{
		votingService: votingService,
	}
}

// Name returns the name of this checker
func (c *VotingChallengesChecker) Name() string {
	return "VotingChallenges"
}

// Check validates voting challenges by looking up votes from the store
func (c *VotingChallengesChecker) Check(ctx sdk.Context, approval *types.CollectionApproval, collection *types.TokenCollection, to string, from string, initiator string, approvalLevel string, approverAddress string, merkleProofs []*types.MerkleProof, ethSignatureProofs []*types.ETHSignatureProof, memo string, isPrioritized bool) (string, error) {
	if approval.ApprovalCriteria == nil || len(approval.ApprovalCriteria.VotingChallenges) == 0 {
		return "", nil
	}

	challenges := approval.ApprovalCriteria.VotingChallenges
	for _, challenge := range challenges {
		if challenge == nil {
			detErrMsg := "voting challenge is nil"
			return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
		}

		// Construct the scoped proposal ID
		proposalId := challenge.ProposalId

		// Calculate total possible weight from all voters using arbitrary-precision big.Int
		totalPossibleWeight := new(big.Int)
		for _, voter := range challenge.Voters {
			if voter == nil {
				detErrMsg := "voter is nil in voting challenge"
				return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
			}
			if voter.Address == "" {
				detErrMsg := "voter address is empty"
				return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
			}
			if voter.Weight.IsNil() || voter.Weight.IsZero() {
				detErrMsg := fmt.Sprintf("voter %s has zero weight", voter.Address)
				return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
			}
			totalPossibleWeight.Add(totalPossibleWeight, voter.Weight.BigInt())
		}

		if totalPossibleWeight.Sign() <= 0 {
			detErrMsg := "total possible weight is zero - no voters defined"
			return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
		}

		// Retrieve all votes for this proposal and sum voterWeight * yesWeight without intermediate truncation.
		// Iterate over the original slice for deterministic ordering.
		weightedYesSum := new(big.Int)
		collectionId := collection.CollectionId
		for _, voter := range challenge.Voters {
			voterAddress := voter.Address
			voteKey := c.votingService.ConstructVoteKey(collectionId, approverAddress, approvalLevel, approval.ApprovalId, proposalId, voterAddress)

			vote, found := c.votingService.GetVoteFromStore(ctx, voteKey)
			if found && vote != nil {
				// Validate the vote
				if vote.ProposalId != proposalId {
					detErrMsg := fmt.Sprintf("vote proposal ID mismatch: expected %s, got %s", proposalId, vote.ProposalId)
					return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
				}
				if vote.Voter != voterAddress {
					detErrMsg := fmt.Sprintf("vote voter mismatch: expected %s, got %s", voterAddress, vote.Voter)
					return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
				}
				if vote.YesWeight.IsNil() {
					detErrMsg := "vote yesWeight is nil"
					return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
				}
				if vote.YesWeight.GT(sdkmath.NewUint(100)) {
					detErrMsg := fmt.Sprintf("vote yesWeight exceeds 100: %s", vote.YesWeight.String())
					return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
				}

				weightedYesSum.Add(weightedYesSum, new(big.Int).Mul(voter.Weight.BigInt(), vote.YesWeight.BigInt()))
			}
		}

		quorumThreshold := big.NewInt(0)
		if !challenge.QuorumThreshold.IsNil() {
			quorumThreshold = challenge.QuorumThreshold.BigInt()
		}
		requiredWeightedYes := new(big.Int).Mul(totalPossibleWeight, quorumThreshold)

		// Compare sum(voterWeight * yesWeight) >= totalPossibleWeight * quorumThreshold directly
		if weightedYesSum.Cmp(requiredWeightedYes) < 0 {
			percentage := new(big.Int).Quo(weightedYesSum, totalPossibleWeight)
			detErrMsg := fmt.Sprintf("voting challenge threshold not met: got %s%%, need %s%%", percentage.String(), quorumThreshold.String())
			return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
		}

		// Check delay-after-quorum if configured
		if !challenge.DelayAfterQuorum.IsNil() && challenge.DelayAfterQuorum.GT(sdkmath.NewUint(0)) {
			trackerKey := c.votingService.ConstructChallengeTrackerKey(collectionId, approverAddress, approvalLevel, approval.ApprovalId, proposalId)

			tracker, found := c.votingService.GetVotingChallengeTrackerFromStore(ctx, trackerKey)
			if !found || tracker == nil || tracker.QuorumReachedTimestamp.IsNil() || tracker.QuorumReachedTimestamp.IsZero() {
				detErrMsg := "voting challenge delay: quorum was reached but tracker not found (vote again to initialize)"
				return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
			}

			now := new(big.Int).SetUint64(uint64(ctx.BlockTime().UnixMilli()))
			executeAfter := new(big.Int).Add(tracker.QuorumReachedTimestamp.BigInt(), challenge.DelayAfterQuorum.BigInt())

			if now.Cmp(executeAfter) < 0 {
				remainingMs := new(big.Int).Sub(executeAfter, now)
				detErrMsg := fmt.Sprintf("voting challenge delay not elapsed: %s ms remaining (quorum reached at %s, delay is %s ms)",
					remainingMs.String(), tracker.QuorumReachedTimestamp.String(), challenge.DelayAfterQuorum.String())
				return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
			}
		}
	}

	return "", nil
}
