package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Queries a dynamic store value by its store ID and address.
func (k Keeper) GetDynamicStoreValue(goCtx context.Context, req *types.QueryGetDynamicStoreValueRequest) (*types.QueryGetDynamicStoreValueResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	storeId := sdkmath.NewUintFromString(req.StoreId)

	dynamicStoreValue, found := k.GetDynamicStoreValueFromStore(ctx, storeId, req.Address)
	if !found {
		dynamicStore, found := k.GetDynamicStoreFromStore(ctx, storeId)
		if !found {
			return nil, status.Error(codes.NotFound, "dynamic store not found")
		}

		// Return the default value from the store
		return &types.QueryGetDynamicStoreValueResponse{
			Value: &types.DynamicStoreValue{
				StoreId: storeId,
				Address: req.Address,
				Value:   dynamicStore.DefaultValue,
			},
		}, nil
	}

	return &types.QueryGetDynamicStoreValueResponse{
		Value: &dynamicStoreValue,
	}, nil
}
