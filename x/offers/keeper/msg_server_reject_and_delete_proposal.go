package keeper

import (
	"context"

	"bitbadgeschain/x/offers/types"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) RejectAndDeleteProposal(goCtx context.Context, msg *types.MsgRejectAndDeleteProposal) (*types.MsgRejectAndDeleteProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Convert the ID string to sdkmath.Uint
	proposalId := msg.Id

	// Get the proposal from the store
	proposal, found := k.GetProposalFromStore(ctx, proposalId)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrProposalNotFound, "proposal ID %s not found", msg.Id)
	}

	// Check if the message creator is a party to the proposal
	isParty := false
	for _, party := range proposal.Parties {
		if party.Creator == msg.Creator {
			isParty = true
			break
		}
	}

	if !isParty {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "creator %s is not a party to this proposal", msg.Creator)
	}

	// Delete the proposal from the store
	k.DeleteProposalFromStore(ctx, proposalId)

	return &types.MsgRejectAndDeleteProposalResponse{}, nil
}
