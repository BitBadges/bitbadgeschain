package keeper

import (
	"context"

	"bitbadgeschain/x/badges/types"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Queries a balance for the given address and badgeId and returns its contents.
func (k Keeper) GetBalance(goCtx context.Context, req *types.QueryGetBalanceRequest) (*types.QueryGetBalanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	//Assert that initiatedBy owns the required badges
	balances := &types.UserBalanceStore{}
	currCollectionId := req.CollectionId
	collection, found := k.GetCollectionFromStore(ctx, currCollectionId)
	if !found {
		return nil, sdkerrors.Wrapf(ErrCollectionNotExists, "collection %s not found", currCollectionId)
	} else {
		isStandardBalances := collection.BalancesType == "Standard"
		if isStandardBalances || req.Address == "Mint" || req.Address == "Total" {
			balances = k.GetBalanceOrApplyDefault(ctx, collection, req.Address)
		} else {
			return nil, sdkerrors.Wrapf(ErrWrongBalancesType, "unsupported balances type %s %s", collection.BalancesType, collection.CollectionId)
		}
	}

	return &types.QueryGetBalanceResponse{
		Balance: balances,
	}, nil
}
