package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// GetCollectionStats returns holder count and circulating supply for a collection.
func (k Keeper) GetCollectionStats(goCtx context.Context, req *types.QueryGetCollectionStatsRequest) (*types.QueryGetCollectionStatsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	collectionId, err := sdkmath.ParseUint(req.CollectionId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid collection ID")
	}

	if !k.StoreHasCollectionID(ctx, collectionId) {
		return nil, types.ErrInvalidCollectionID
	}

	stats, _ := k.GetCollectionStatsFromStore(ctx, collectionId)

	return &types.QueryGetCollectionStatsResponse{
		Stats: stats,
	}, nil
}
