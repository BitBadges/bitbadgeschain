package keeper

import (
	"context"

	"bitbadgeschain/x/offers/types"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateProposal(goCtx context.Context, msg *types.MsgCreateProposal) (*types.MsgCreateProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the next proposal ID
	nextProposalId := k.GetNextProposalId(ctx)

	// Create a new Proposal
	proposal := &types.Proposal{
		Id:         nextProposalId,
		Parties:    msg.Parties,
		ValidTimes: msg.ValidTimes,
		CreatedBy:  msg.Creator,
	}

	// Set the proposal in the store
	err := k.SetProposalInStore(ctx, proposal)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to set proposal in store")
	}

	// Increment the next proposal ID
	k.IncrementNextProposalId(ctx)

	return &types.MsgCreateProposalResponse{
		Id: nextProposalId,
	}, nil
}
