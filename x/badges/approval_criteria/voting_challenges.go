package approval_criteria

import (
	"fmt"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// VotingService defines the interface for accessing vote data
type VotingService interface {
	GetVoteFromStore(ctx sdk.Context, key string) (*types.VoteProof, bool)
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

		// Construct the scoped proposal ID: collectionId-approverAddress-approvalLevel-approvalId-challengeId
		// The challenge.ProposalId should already be scoped, but we use it as-is
		proposalId := challenge.ProposalId

		// Calculate total possible weight from all voters
		totalPossibleWeight := sdkmath.NewUint(0)
		voterMap := make(map[string]sdkmath.Uint) // voter address -> weight
		for _, voter := range challenge.Voters {
			if voter == nil {
				detErrMsg := "voter is nil in voting challenge"
				return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
			}
			if voter.Address == "" {
				detErrMsg := "voter address is empty"
				return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
			}
			if voter.Weight.IsZero() {
				detErrMsg := fmt.Sprintf("voter %s has zero weight", voter.Address)
				return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
			}
			totalPossibleWeight = totalPossibleWeight.Add(voter.Weight)
			voterMap[voter.Address] = voter.Weight
		}

		if totalPossibleWeight.IsZero() {
			detErrMsg := "total possible weight is zero - no voters defined"
			return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
		}

		// Retrieve all votes for this proposal and calculate total yes weight
		totalYesWeight := sdkmath.NewUint(0)
		collectionId := collection.CollectionId
		for voterAddress := range voterMap {
			// Construct the vote key using the same pattern as ConstructVotingTrackerKey
			// Format: collectionId-approverAddress-approvalLevel-approvalId-proposalId-voterAddress
			voteKey := fmt.Sprintf("%s-%s-%s-%s-%s-%s",
				collectionId.String(),
				approverAddress,
				approvalLevel,
				approval.ApprovalId,
				proposalId,
				voterAddress)

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
				if vote.YesWeight.GT(sdkmath.NewUint(100)) {
					detErrMsg := fmt.Sprintf("vote yesWeight exceeds 100: %s", vote.YesWeight.String())
					return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
				}

				// Calculate the yes weight contribution from this voter
				// yesWeight is a percentage (0-100), so we multiply voter weight by (yesWeight/100)
				voterWeight := voterMap[voterAddress]
				yesWeightPercent := vote.YesWeight
				// Calculate: (voterWeight * yesWeightPercent) / 100
				yesContribution := voterWeight.Mul(yesWeightPercent).Quo(sdkmath.NewUint(100))
				totalYesWeight = totalYesWeight.Add(yesContribution)
			}
		}

		// Calculate the percentage of total possible weight that voted yes
		// totalYesWeight * 100 / totalPossibleWeight
		percentage := totalYesWeight.Mul(sdkmath.NewUint(100)).Quo(totalPossibleWeight)

		// Check if percentage meets the quorum threshold
		if percentage.LT(challenge.QuorumThreshold) {
			detErrMsg := fmt.Sprintf("voting challenge threshold not met: got %s%%, need %s%%", percentage.String(), challenge.QuorumThreshold.String())
			return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
		}
	}

	return "", nil
}

