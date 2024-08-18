package keeper

import (
	"context"

	"bitbadgeschain/x/offers/types"

	badgetypes "bitbadgeschain/x/badges/types"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) ExecuteProposal(goCtx context.Context, msg *types.MsgExecuteProposal) (*types.MsgExecuteProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Convert the ID string to sdkmath.Uint
	proposalId := msg.Id

	// Get the proposal from the store
	proposal, found := k.GetProposalFromStore(ctx, proposalId)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrProposalNotFound, "proposal ID %s not found", msg.Id)
	}

	// Check if all parties have accepted
	for _, party := range proposal.Parties {
		if !party.Accepted {
			return nil, sdkerrors.Wrapf(types.ErrProposalNotAccepted, "not all parties have accepted the proposal")
		}
	}

	if !proposal.AnyoneCanFinalize {
		// Ensure that the creator of the proposal is the same as one of the parties
		creatorIsParty := false
		for _, party := range proposal.Parties {
			if party.Creator == msg.Creator {
				creatorIsParty = true
				break
			}
		}

		if !creatorIsParty {
			return nil, sdkerrors.Wrapf(types.ErrProposalNotAccepted, "creator of the proposal is not one of the parties")
		}

		if proposal.CreatorMustFinalize && msg.Creator != proposal.CreatedBy {
			return nil, sdkerrors.Wrapf(types.ErrProposalNotAccepted, "creator of the proposal must finalize the proposal")
		}
	}

	castedTimes := types.CastUintRanges(proposal.ValidTimes)
	found, err := badgetypes.SearchUintRangesForUint(sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli())), castedTimes)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to search for valid times")
	}

	if !found {
		return nil, sdkerrors.Wrapf(types.ErrProposalNotValidAtThisTime, "proposal is not valid at this time")
	}

	// Execute all messages for all parties w/ correct creator address
	for _, party := range proposal.Parties {
		for _, msg := range party.MsgsToExecute {
			msgsToExecute := []*codectypes.Any{msg.Msg}
			creatorAddr := party.Creator
			if msg.UseContractAddress {
				creatorAddr = proposal.ContractAddress
			}

			err := k.ExecuteGenericMsgs(ctx, msgsToExecute, creatorAddr)
			if err != nil {
				return nil, sdkerrors.Wrapf(err, "failed to execute messages for party %s", party.Creator)
			}
		}
	}

	// Delete the proposal from the store as it has been executed
	k.DeleteProposalFromStore(ctx, proposalId)

	return &types.MsgExecuteProposalResponse{}, nil
}
