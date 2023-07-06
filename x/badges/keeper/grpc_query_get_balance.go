package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
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

	userBalanceKey := ConstructBalanceKey(req.Address, req.CollectionId)
	userBalance, found := k.GetUserBalanceFromStore(ctx, userBalanceKey)
	if found {
		return &types.QueryGetBalanceResponse{
			Balance: userBalance,
		}, nil
	} else {
		collection, found := k.GetCollectionFromStore(ctx, req.CollectionId)
		if !found {
			return nil, status.Error(codes.NotFound, "collection and balances not found")
		}

		//TODO: Recursive get inherited balance for user w/ inherited balances
		if !IsStandardBalances(collection) {
			return nil, status.Error(codes.NotFound, "collection and balances not found: unsupported balances type")
		}

		blankUserBalance := &types.UserBalanceStore{
			Balances: []*types.Balance{},
			ApprovedOutgoingTransfersTimeline: collection.DefaultUserApprovedOutgoingTransfersTimeline,
			ApprovedIncomingTransfersTimeline: collection.DefaultUserApprovedIncomingTransfersTimeline,
			Permissions: &types.UserPermissions{},
		}
		return &types.QueryGetBalanceResponse{
			Balance: blankUserBalance,
		}, nil
	}
}
