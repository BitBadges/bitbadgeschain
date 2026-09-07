package keeper

import (
	"context"
	"fmt"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CastVote(goCtx context.Context, msg *types.MsgCastVote) (*types.MsgCastVoteResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate the message
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// Get the collection
	collection, found := k.GetCollectionFromStore(ctx, msg.CollectionId)
	if !found {
		return nil, sdkerrors.Wrap(ErrCollectionNotExists, fmt.Sprintf("Collection %s not found", msg.CollectionId.String()))
	}

	// Find the approval and verify the voting challenge exists
	var votingChallenge *types.VotingChallenge
	foundChallenge := false

	// Check collection-level approvals
	if msg.ApprovalLevel == "collection" {
		for _, app := range collection.CollectionApprovals {
			if app.ApprovalId == msg.ApprovalId {
				if app.ApprovalCriteria != nil {
					for _, challenge := range app.ApprovalCriteria.VotingChallenges {
						if challenge != nil && challenge.ProposalId == msg.ProposalId {
							votingChallenge = challenge
							foundChallenge = true
							break
						}
					}
				}
				break
			}
		}
	} else {
		// Check user-level approvals (incoming/outgoing)
		balanceStore, _, err := k.GetBalanceOrApplyDefault(ctx, collection, msg.ApproverAddress)
		if err != nil {
			return nil, err
		}

		if msg.ApprovalLevel == "incoming" {
			for _, app := range balanceStore.IncomingApprovals {
				if app.ApprovalId == msg.ApprovalId {
					if app.ApprovalCriteria != nil {
						for _, challenge := range app.ApprovalCriteria.VotingChallenges {
							if challenge != nil && challenge.ProposalId == msg.ProposalId {
								votingChallenge = challenge
								foundChallenge = true
								break
							}
						}
					}
					break
				}
			}
		} else if msg.ApprovalLevel == "outgoing" {
			for _, app := range balanceStore.OutgoingApprovals {
				if app.ApprovalId == msg.ApprovalId {
					if app.ApprovalCriteria != nil {
						for _, challenge := range app.ApprovalCriteria.VotingChallenges {
							if challenge != nil && challenge.ProposalId == msg.ProposalId {
								votingChallenge = challenge
								foundChallenge = true
								break
							}
						}
					}
					break
				}
			}
		} else {
			return nil, sdkerrors.Wrap(types.ErrInvalidRequest, fmt.Sprintf("Invalid approval level: %s", msg.ApprovalLevel))
		}
	}

	if !foundChallenge || votingChallenge == nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidRequest, "Voting challenge not found")
	}

	// Verify the voter is in the voters list
	var voterAddresses []string
	voterWeight := sdkmath.ZeroUint()
	for _, voter := range votingChallenge.Voters {
		if voter != nil {
			address, _ := canonicalBech32(voter.Address)
			if address == msg.Creator {
				voterAddresses = append(voterAddresses, voter.Address)
				voterWeight = voterWeight.Add(voter.Weight)
			}
		}
	}

	if len(voterAddresses) == 0 {
		return nil, sdkerrors.Wrap(types.ErrInvalidRequest, "Voter not found in voting challenge voters list")
	}

	// Validate yesWeight is 0-100
	if msg.YesWeight.GT(sdkmath.NewUint(100)) {
		return nil, sdkerrors.Wrap(types.ErrInvalidRequest, fmt.Sprintf("yesWeight must be between 0 and 100, got %s", msg.YesWeight.String()))
	}

	now := sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli()))

	// Keep legacy voter slots separate so quorum retains its per-slot rounding.
	for _, voterAddress := range voterAddresses {
		voteKey := ConstructVotingTrackerKey(msg.CollectionId, msg.ApproverAddress, msg.ApprovalLevel, msg.ApprovalId, msg.ProposalId, voterAddress)
		vote := &types.VoteProof{
			ProposalId: msg.ProposalId,
			Voter:      voterAddress,
			YesWeight:  msg.YesWeight,
			VotedAt:    now,
		}
		if err := k.SetVoteInStore(ctx, voteKey, vote); err != nil {
			return nil, sdkerrors.Wrap(err, "Failed to store vote")
		}
	}

	// Update quorum tracking for challenges with delay or reset
	if votingChallenge.DelayAfterQuorum.GT(sdkmath.NewUint(0)) || votingChallenge.ResetAfterExecution {
		trackerKey := ConstructVotingChallengeTrackerKey(
			msg.CollectionId,
			msg.ApproverAddress,
			msg.ApprovalLevel,
			msg.ApprovalId,
			msg.ProposalId,
		)

		// Calculate current quorum to determine if threshold is met
		quorumMet := k.calculateQuorumMet(ctx, votingChallenge, collection.CollectionId, msg.ApproverAddress, msg.ApprovalLevel, msg.ApprovalId)

		tracker, found := k.GetVotingChallengeTrackerFromStore(ctx, trackerKey)
		if !found {
			tracker = &types.VotingChallengeTracker{
				QuorumReachedTimestamp: sdkmath.NewUint(0),
			}
		}

		if quorumMet {
			// Only set quorumReachedTimestamp if it wasn't already set
			if tracker.QuorumReachedTimestamp.IsZero() {
				tracker.QuorumReachedTimestamp = now
			}
		} else {
			// Quorum dropped — clear the timestamp (delay resets)
			tracker.QuorumReachedTimestamp = sdkmath.NewUint(0)
		}

		if err := k.SetVotingChallengeTrackerInStore(ctx, trackerKey, tracker); err != nil {
			return nil, sdkerrors.Wrap(err, "Failed to store voting challenge tracker")
		}
	}

	// Emit event
	msgStr, err := MarshalMessageForEvent(msg)
	if err != nil {
		return nil, err
	}

	EmitMessageAndIndexerEvents(ctx,
		sdk.NewAttribute(sdk.AttributeKeyModule, "tokenization"),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		sdk.NewAttribute("msg_type", "cast_vote"),
		sdk.NewAttribute("msg", msgStr),
		sdk.NewAttribute("collection_id", msg.CollectionId.String()),
		sdk.NewAttribute("approval_level", msg.ApprovalLevel),
		sdk.NewAttribute("approver_address", msg.ApproverAddress),
		sdk.NewAttribute("approval_id", msg.ApprovalId),
		sdk.NewAttribute("proposal_id", msg.ProposalId),
		sdk.NewAttribute("voter", msg.Creator),
		sdk.NewAttribute("yes_weight", msg.YesWeight.String()),
		sdk.NewAttribute("voter_weight", voterWeight.String()),
	)

	return &types.MsgCastVoteResponse{}, nil
}

// calculateQuorumMet checks if the current votes meet the quorum threshold
func (k Keeper) calculateQuorumMet(ctx sdk.Context, challenge *types.VotingChallenge, collectionId sdkmath.Uint, approverAddress string, approvalLevel string, approvalId string) bool {
	totalPossibleWeight := sdkmath.NewUint(0)
	totalYesWeight := sdkmath.NewUint(0)

	for _, voter := range challenge.Voters {
		if voter == nil {
			continue
		}
		totalPossibleWeight = totalPossibleWeight.Add(voter.Weight)

		voteKey := ConstructVotingTrackerKey(collectionId, approverAddress, approvalLevel, approvalId, challenge.ProposalId, voter.Address)
		vote, found := k.GetVoteFromStore(ctx, voteKey)
		if found && vote != nil && vote.YesWeight.GT(sdkmath.NewUint(0)) {
			yesContribution := voter.Weight.Mul(vote.YesWeight).Quo(sdkmath.NewUint(100))
			totalYesWeight = totalYesWeight.Add(yesContribution)
		}
	}

	if totalPossibleWeight.IsZero() {
		return false
	}

	percentage := totalYesWeight.Mul(sdkmath.NewUint(100)).Quo(totalPossibleWeight)
	return percentage.GTE(challenge.QuorumThreshold)
}
