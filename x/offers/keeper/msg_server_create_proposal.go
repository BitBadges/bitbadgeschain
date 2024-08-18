package keeper

import (
	"context"
	"encoding/binary"

	"bitbadgeschain/x/offers/types"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func (k msgServer) CreateProposal(goCtx context.Context, msg *types.MsgCreateProposal) (*types.MsgCreateProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the next proposal ID
	nextProposalId := k.GetNextProposalId(ctx)

	// From cosmos SDK x/group module
	// Generate account address for list
	var accountAddr sdk.AccAddress
	// loop here in the rare case where a ADR-028-derived address creates a
	// collision with an existing address.
	for {
		derivationKey := make([]byte, 8)
		nextId := nextProposalId
		binary.BigEndian.PutUint64(derivationKey, nextId.Uint64())

		ac, err := authtypes.NewModuleCredential(types.ModuleName, GenerationPrefix, derivationKey)
		if err != nil {
			return nil, sdkerrors.Wrap(err, "failed to create proposal")
		}
		//generate the address from the credential
		accountAddr = sdk.AccAddress(ac.Address())

		break
	}

	contractAddress := accountAddr.String()

	// Create a new Proposal
	proposal := &types.Proposal{
		Id:                  nextProposalId,
		Parties:             msg.Parties,
		ValidTimes:          msg.ValidTimes,
		CreatedBy:           msg.Creator,
		CreatorMustFinalize: msg.CreatorMustFinalize,
		AnyoneCanFinalize:   msg.AnyoneCanFinalize,
		ContractAddress:     contractAddress,
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
