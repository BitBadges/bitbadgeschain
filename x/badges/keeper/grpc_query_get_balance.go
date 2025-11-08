package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Queries a balance for the given address and tokenId and returns its contents.
func (k Keeper) GetBalance(goCtx context.Context, req *types.QueryGetBalanceRequest) (*types.QueryGetBalanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	//Assert that initiatedBy owns the required tokens
	balances := &types.UserBalanceStore{}
	collectionId := sdkmath.NewUintFromString(req.CollectionId)
	collection, found := k.GetCollectionFromStore(ctx, collectionId)
	if !found {
		return nil, sdkerrors.Wrapf(ErrCollectionNotExists, "collection %s not found", req.CollectionId)
	} else {
		// All collections now use Standard balances
		// We technically want to apply versions here, but it doesn't actually set them in store (due to queries)
		balances, _ = k.GetBalanceOrApplyDefault(ctx, collection, req.Address)
	}

	return &types.QueryGetBalanceResponse{
		Balance: balances,
	}, nil
}
