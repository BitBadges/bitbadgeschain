package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"bitbadgeschain/x/offers/types"
)

func (k Keeper) GetProposal(goCtx context.Context, req *types.QueryGetProposalRequest) (*types.QueryGetProposalResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the proposal from the store
	proposal, found := k.GetProposalFromStore(ctx, req.Id)
	if !found {
		return nil, status.Error(codes.NotFound, "proposal not found")
	}

	return &types.QueryGetProposalResponse{Proposal: proposal}, nil
}
