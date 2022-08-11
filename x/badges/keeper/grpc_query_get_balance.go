package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Queries a balance for the given address and badgeId and returns its contents.
func (k Keeper) GetBalance(goCtx context.Context, req *types.QueryGetBalanceRequest) (*types.QueryGetBalanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	userBalanceKey := ConstructBalanceKey(req.Address, req.BadgeId)
	userBalanceInfo, found := k.GetUserBalanceFromStore(ctx, userBalanceKey)
	if found {
		return &types.QueryGetBalanceResponse{
			BalanceInfo: &userBalanceInfo,
		}, nil
	} else {
		blankUserBalanceInfo := &types.UserBalanceInfo{}
		return &types.QueryGetBalanceResponse{
			BalanceInfo: blankUserBalanceInfo,
		}, nil
	}
}
