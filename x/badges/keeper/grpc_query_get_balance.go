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

	BalanceKey := ConstructBalanceKey(req.Address, req.BadgeId)
	badgeBalanceInfo, found := k.GetBadgeBalanceFromStore(ctx, BalanceKey)
	if found {
		return &types.QueryGetBalanceResponse{
			BalanceInfo: &badgeBalanceInfo,
		}, nil
	} else {
		blankBadgeBalanceInfo := &types.BadgeBalanceInfo{}
		return &types.QueryGetBalanceResponse{
			BalanceInfo: blankBadgeBalanceInfo,
		}, nil
	}
}
