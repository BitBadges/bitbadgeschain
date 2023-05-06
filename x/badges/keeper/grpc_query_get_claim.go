package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetClaim(goCtx context.Context, req *types.QueryGetClaimRequest) (*types.QueryGetClaimResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	claim, found := k.GetClaimFromStore(ctx, req.CollectionId, req.ClaimId);
	if !found {
		return nil, ErrClaimNotExists
	}

	return &types.QueryGetClaimResponse{
		Claim: &claim,
	}, nil
}
