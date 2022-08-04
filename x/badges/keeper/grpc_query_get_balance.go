package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetBalance(goCtx context.Context, req *types.QueryGetBalanceRequest) (*types.QueryGetBalanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Verify that the from and to addresses are registered;
	account_nums := []uint64{}
	account_nums = append(account_nums, req.Address)
	err := k.AssertAccountNumbersAreRegistered(ctx, account_nums)
	if err != nil {
		return nil, err
	}

	_, found := k.GetBadgeFromStore(ctx, req.BadgeId)
	if !found {
		return nil, ErrBadgeNotExists
	}

	full_id := GetBalanceKey(
		req.Address,
		req.BadgeId,
	)
	badgeBalanceInfo, found := k.GetBadgeBalanceFromStore(ctx, full_id)
	if found {
		return &types.QueryGetBalanceResponse{
			BalanceInfo: &badgeBalanceInfo,
		}, nil
	}

	return &types.QueryGetBalanceResponse{
		BalanceInfo: &types.BadgeBalanceInfo{},
	}, nil
}
