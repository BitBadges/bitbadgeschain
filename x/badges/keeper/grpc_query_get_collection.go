package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Queries a collection by its ID and returns its contents.
func (k Keeper) GetCollection(goCtx context.Context, req *types.QueryGetCollectionRequest) (*types.QueryGetCollectionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	collectionId := sdkmath.NewUintFromString(req.CollectionId)
	collection, found := k.GetCollectionFromStore(ctx, collectionId)
	if !found {
		return nil, types.ErrInvalidCollectionID
	}

	return &types.QueryGetCollectionResponse{
		Collection: collection,
	}, nil
}
