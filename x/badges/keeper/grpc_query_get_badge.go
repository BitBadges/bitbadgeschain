package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Queries a badge by its ID and returns its contents.
func (k Keeper) GetCollection(goCtx context.Context, req *types.QueryGetCollectionRequest) (*types.QueryGetCollectionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	badge, found := k.GetCollectionFromStore(ctx, req.Id)
	if !found {
		return nil, ErrCollectionNotExists
	}

	return &types.QueryGetCollectionResponse{
		Collection: &badge,
	}, nil
}
