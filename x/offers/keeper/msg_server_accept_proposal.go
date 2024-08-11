package keeper

import (
	"context"

	"bitbadgeschain/x/offers/types"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) AcceptProposal(goCtx context.Context, msg *types.MsgAcceptProposal) (*types.MsgAcceptProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Convert the ID string to sdkmath.Uint
	proposalId := msg.Id

	// Get the proposal from the store
	proposal, found := k.GetProposalFromStore(ctx, proposalId)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrProposalNotFound, "proposal ID %s not found", msg.Id)
	}

	// Find the party of the message creator and mark it as accepted
	partyFound := false
	for i, party := range proposal.Parties {
		if party.Creator == msg.Creator {
			proposal.Parties[i].Accepted = true
			partyFound = true
			break
		}
	}

	if !partyFound {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "creator %s is not a party to this proposal", msg.Creator)
	}

	// Update the proposal in the store
	err := k.SetProposalInStore(ctx, proposal)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to update proposal in store")
	}

	return &types.MsgAcceptProposalResponse{}, nil
}
