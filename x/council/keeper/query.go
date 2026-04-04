package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/council/types"
)

// QueryCouncilRequest is the request type for querying a single council.
type QueryCouncilRequest struct {
	CouncilId uint64 `json:"councilId"`
}

// QueryCouncilResponse is the response type.
type QueryCouncilResponse struct {
	Council types.Council `json:"council"`
}

// QueryProposalRequest is the request type for querying a single proposal.
type QueryProposalRequest struct {
	CouncilId  uint64 `json:"councilId"`
	ProposalId uint64 `json:"proposalId"`
}

// QueryProposalResponse is the response type.
type QueryProposalResponse struct {
	Proposal types.Proposal `json:"proposal"`
}

// QueryCouncil returns a council by ID.
func (k Keeper) QueryCouncil(goCtx context.Context, req *QueryCouncilRequest) (*QueryCouncilResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	council, found := k.GetCouncil(ctx, req.CouncilId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrCouncilNotFound, "council %d", req.CouncilId)
	}

	return &QueryCouncilResponse{Council: council}, nil
}

// QueryProposal returns a proposal by council ID and proposal ID.
func (k Keeper) QueryProposal(goCtx context.Context, req *QueryProposalRequest) (*QueryProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	proposal, found := k.GetProposal(ctx, req.CouncilId, req.ProposalId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrProposalNotFound, "proposal %d/%d", req.CouncilId, req.ProposalId)
	}

	return &QueryProposalResponse{Proposal: proposal}, nil
}
