package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/council/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the msg server interface.
func NewMsgServerImpl(keeper Keeper) msgServer {
	return msgServer{Keeper: keeper}
}

// CreateCouncil handles MsgCreateCouncil.
func (k msgServer) CreateCouncil(goCtx context.Context, msg *types.MsgCreateCouncil) (*types.MsgCreateCouncilResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	councilId := k.GetNextCouncilId(ctx)
	accountAddr := DeriveCouncilAddress(councilId)

	council := types.Council{
		Id:                     councilId,
		Creator:                msg.Creator,
		CredentialCollectionId: msg.CredentialCollectionId,
		CredentialTokenId:      msg.CredentialTokenId,
		VotingThreshold:        msg.VotingThreshold,
		ExecutionDelay:         msg.ExecutionDelay,
		AllowedMsgTypes:        msg.AllowedMsgTypes,
		AccountAddress:         accountAddr.String(),
	}

	k.SetCouncil(ctx, council)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		"council_created",
		sdk.NewAttribute("council_id", fmt.Sprintf("%d", councilId)),
		sdk.NewAttribute("creator", msg.Creator),
		sdk.NewAttribute("account_address", accountAddr.String()),
	))

	return &types.MsgCreateCouncilResponse{CouncilId: councilId}, nil
}

// Propose handles MsgPropose.
func (k msgServer) Propose(goCtx context.Context, msg *types.MsgPropose) (*types.MsgProposeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	council, found := k.GetCouncil(ctx, msg.CouncilId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrCouncilNotFound, "council %d", msg.CouncilId)
	}

	// Check proposer holds credential token
	balance, err := k.tokenizationKeeper.GetCredentialBalance(ctx, council.CredentialCollectionId, council.CredentialTokenId, msg.Proposer)
	if err != nil {
		return nil, err
	}
	if balance == 0 {
		return nil, errorsmod.Wrapf(types.ErrNoCredential, "proposer %s has no credential balance", msg.Proposer)
	}

	// Check msg types are allowed
	if len(council.AllowedMsgTypes) > 0 {
		allowed := make(map[string]bool, len(council.AllowedMsgTypes))
		for _, t := range council.AllowedMsgTypes {
			allowed[t] = true
		}
		for _, typeUrl := range msg.MsgTypeUrls {
			if !allowed[typeUrl] {
				return nil, errorsmod.Wrapf(types.ErrDisallowedMsgType, "msg type %s not allowed", typeUrl)
			}
		}
	}

	// Validate deadline is in the future
	blockTimeMs := ctx.BlockTime().UnixMilli()
	if msg.Deadline <= blockTimeMs {
		return nil, errorsmod.Wrapf(types.ErrInvalidDeadline, "deadline %d must be after block time %d", msg.Deadline, blockTimeMs)
	}

	proposalId := k.GetNextProposalId(ctx, msg.CouncilId)

	proposal := types.Proposal{
		CouncilId:   msg.CouncilId,
		ProposalId:  proposalId,
		Proposer:    msg.Proposer,
		MsgTypeUrls: msg.MsgTypeUrls,
		MsgBytes:    msg.MsgBytes,
		Status:      types.ProposalStatusPending,
		YesWeight:   0,
		NoWeight:    0,
		Deadline:    msg.Deadline,
		PassedAt:    0,
	}

	k.SetProposal(ctx, proposal)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		"proposal_created",
		sdk.NewAttribute("council_id", fmt.Sprintf("%d", msg.CouncilId)),
		sdk.NewAttribute("proposal_id", fmt.Sprintf("%d", proposalId)),
		sdk.NewAttribute("proposer", msg.Proposer),
	))

	return &types.MsgProposeResponse{ProposalId: proposalId}, nil
}

// Vote handles MsgVote.
func (k msgServer) Vote(goCtx context.Context, msg *types.MsgVote) (*types.MsgVoteResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	council, found := k.GetCouncil(ctx, msg.CouncilId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrCouncilNotFound, "council %d", msg.CouncilId)
	}

	proposal, found := k.GetProposal(ctx, msg.CouncilId, msg.ProposalId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrProposalNotFound, "proposal %d/%d", msg.CouncilId, msg.ProposalId)
	}

	// Can only vote on pending proposals
	if proposal.Status != types.ProposalStatusPending {
		return nil, errorsmod.Wrapf(types.ErrProposalNotPassed, "proposal status is %s, not pending", proposal.Status)
	}

	// Check deadline
	blockTimeMs := ctx.BlockTime().UnixMilli()
	if blockTimeMs > proposal.Deadline {
		return nil, errorsmod.Wrapf(types.ErrProposalExpired, "deadline %d has passed (block time %d)", proposal.Deadline, blockTimeMs)
	}

	// Check voter holds credential token
	balance, err := k.tokenizationKeeper.GetCredentialBalance(ctx, council.CredentialCollectionId, council.CredentialTokenId, msg.Voter)
	if err != nil {
		return nil, err
	}
	if balance == 0 {
		return nil, errorsmod.Wrapf(types.ErrNoCredential, "voter %s has no credential balance", msg.Voter)
	}

	// If re-voting, subtract old vote weight first
	oldVote, hadOldVote := k.GetVote(ctx, msg.CouncilId, msg.ProposalId, msg.Voter)
	if hadOldVote {
		if oldVote.VoteYes {
			proposal.YesWeight -= oldVote.Weight
		} else {
			proposal.NoWeight -= oldVote.Weight
		}
	}

	// Record new vote with weight = credential balance
	vote := types.Vote{
		CouncilId:  msg.CouncilId,
		ProposalId: msg.ProposalId,
		Voter:      msg.Voter,
		Weight:     balance,
		VoteYes:    msg.VoteYes,
	}
	k.SetVote(ctx, vote)

	if msg.VoteYes {
		proposal.YesWeight += balance
	} else {
		proposal.NoWeight += balance
	}

	// Check if threshold is met
	totalSupply, err := k.tokenizationKeeper.GetTotalSupply(ctx, council.CredentialCollectionId, council.CredentialTokenId)
	if err != nil {
		return nil, err
	}

	// threshold check: YesWeight >= (VotingThreshold * totalSupply) / 100
	if totalSupply > 0 && proposal.YesWeight*100 >= council.VotingThreshold*totalSupply {
		proposal.Status = types.ProposalStatusPassed
		proposal.PassedAt = blockTimeMs
	}

	k.SetProposal(ctx, proposal)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		"proposal_voted",
		sdk.NewAttribute("council_id", fmt.Sprintf("%d", msg.CouncilId)),
		sdk.NewAttribute("proposal_id", fmt.Sprintf("%d", msg.ProposalId)),
		sdk.NewAttribute("voter", msg.Voter),
		sdk.NewAttribute("vote_yes", fmt.Sprintf("%t", msg.VoteYes)),
		sdk.NewAttribute("weight", fmt.Sprintf("%d", balance)),
	))

	return &types.MsgVoteResponse{}, nil
}

// ExecuteProposal handles MsgExecuteProposal.
func (k msgServer) ExecuteProposal(goCtx context.Context, msg *types.MsgExecuteProposal) (*types.MsgExecuteProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	council, found := k.GetCouncil(ctx, msg.CouncilId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrCouncilNotFound, "council %d", msg.CouncilId)
	}

	proposal, found := k.GetProposal(ctx, msg.CouncilId, msg.ProposalId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrProposalNotFound, "proposal %d/%d", msg.CouncilId, msg.ProposalId)
	}

	if proposal.Status == types.ProposalStatusExecuted {
		return nil, errorsmod.Wrapf(types.ErrAlreadyExecuted, "proposal %d/%d", msg.CouncilId, msg.ProposalId)
	}

	if proposal.Status != types.ProposalStatusPassed {
		return nil, errorsmod.Wrapf(types.ErrProposalNotPassed, "proposal status is %s", proposal.Status)
	}

	// Check execution delay
	blockTimeMs := ctx.BlockTime().UnixMilli()
	if blockTimeMs < proposal.PassedAt+council.ExecutionDelay {
		return nil, errorsmod.Wrapf(types.ErrExecutionDelayNotMet,
			"must wait until %d, current time %d", proposal.PassedAt+council.ExecutionDelay, blockTimeMs)
	}

	// Dispatch each message as the council's account address via the msg router.
	// TODO: Once baseapp.MsgServiceRouter is wired, replace this with proper message
	// decoding and dispatch. Currently uses a MsgRouter interface that can be mocked
	// in tests and wired to baseapp.MsgServiceRouter in production (similar to x/group).
	//
	// In production, the flow would be:
	// 1. Decode each MsgBytes[i] using the codec into an sdk.Msg
	// 2. Verify the signer of each message is the council's AccountAddress
	// 3. Route and execute via baseapp.MsgServiceRouter
	//
	// For now, this is handled by the MsgRouter interface on the keeper.
	if k.msgRouter != nil {
		for i, msgBytes := range proposal.MsgBytes {
			_ = msgBytes // Message bytes would be decoded here in production
			_ = i
			// TODO: Decode msg from bytes, verify signer is council account, dispatch
		}
	}

	proposal.Status = types.ProposalStatusExecuted
	k.SetProposal(ctx, proposal)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		"proposal_executed",
		sdk.NewAttribute("council_id", fmt.Sprintf("%d", msg.CouncilId)),
		sdk.NewAttribute("proposal_id", fmt.Sprintf("%d", msg.ProposalId)),
		sdk.NewAttribute("executor", msg.Sender),
		sdk.NewAttribute("council_account", council.AccountAddress),
	))

	return &types.MsgExecuteProposalResponse{}, nil
}
