package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Queries a balance for the given address and badgeId and returns its contents.
func (k Keeper) IsClaimDataUsed(goCtx context.Context, req *types.QueryIsClaimDataUsedRequest) (*types.QueryIsClaimDataUsedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	used := k.StoreHasUsedClaimData(ctx, req.CollectionId, req.ClaimId, req.ClaimData)
	return &types.QueryIsClaimDataUsedResponse{
		Used: used,
	}, nil
}
