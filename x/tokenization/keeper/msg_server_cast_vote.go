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
	voterFound := false
	var voterWeight sdkmath.Uint
	for _, voter := range votingChallenge.Voters {
		if voter != nil && voter.Address == msg.Creator {
			voterFound = true
			voterWeight = voter.Weight
			break
		}
	}

	if !voterFound {
		return nil, sdkerrors.Wrap(types.ErrInvalidRequest, "Voter not found in voting challenge voters list")
	}

	// Validate yesWeight is 0-100
	if msg.YesWeight.GT(sdkmath.NewUint(100)) {
		return nil, sdkerrors.Wrap(types.ErrInvalidRequest, fmt.Sprintf("yesWeight must be between 0 and 100, got %s", msg.YesWeight.String()))
	}

	// Construct the vote key
	voteKey := ConstructVotingTrackerKey(
		msg.CollectionId,
		msg.ApproverAddress,
		msg.ApprovalLevel,
		msg.ApprovalId,
		msg.ProposalId,
		msg.Creator,
	)

	// Create the vote
	vote := &types.VoteProof{
		ProposalId: msg.ProposalId,
		Voter:      msg.Creator,
		YesWeight:  msg.YesWeight,
	}

	// Store the vote (this will overwrite any existing vote)
	if err := k.SetVoteInStore(ctx, voteKey, vote); err != nil {
		return nil, sdkerrors.Wrap(err, "Failed to store vote")
	}

	// Emit event
	EmitMessageAndIndexerEvents(ctx,
		sdk.NewAttribute(sdk.AttributeKeyModule, "tokenization"),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		sdk.NewAttribute("msg_type", "cast_vote"),
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
