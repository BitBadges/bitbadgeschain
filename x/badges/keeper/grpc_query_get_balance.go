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
	UserBalance, found := k.GetUserBalanceFromStore(ctx, userBalanceKey)
	if found {
		return &types.QueryGetBalanceResponse{
			Balance: &UserBalance,
		}, nil
	} else {
		blankUserBalance := &types.UserBalanceStore{}
		return &types.QueryGetBalanceResponse{
			Balance: blankUserBalance,
		}, nil
	}
}
