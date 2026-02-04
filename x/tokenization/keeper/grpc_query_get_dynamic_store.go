package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Queries a dynamic store by its ID and returns its contents.
func (k Keeper) GetDynamicStore(goCtx context.Context, req *types.QueryGetDynamicStoreRequest) (*types.QueryGetDynamicStoreResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	storeId := sdkmath.NewUintFromString(req.StoreId)
	dynamicStore, found := k.GetDynamicStoreFromStore(ctx, storeId)
	if !found {
		return nil, status.Error(codes.NotFound, "dynamic store not found")
	}

	return &types.QueryGetDynamicStoreResponse{
		Store: &dynamicStore,
	}, nil
}
