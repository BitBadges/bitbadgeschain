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

	address, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}

	err = k.AssertBadgeAndSubBadgeExists(ctx, req.BadgeId, req.SubbadgeId)
	if err != nil {
		return nil, err
	}

	account := k.accountKeeper.GetAccount(ctx, address)

	if account != nil {
		full_id := GetBalanceKey(
			account.GetAccountNumber(),
			req.BadgeId,
			req.SubbadgeId,
		)

		badgeBalanceInfo, found := k.GetBadgeBalanceFromStore(ctx, full_id)
		if found {
			return &types.QueryGetBalanceResponse{
				BalanceInfo: &badgeBalanceInfo,
				Message:     "Successfully queried badge balance info.",
			}, nil
		}
	}

	_ = ctx

	return &types.QueryGetBalanceResponse{
		BalanceInfo: &types.BadgeBalanceInfo{
			Balance:      0,
			PendingNonce: 0,
			Pending:      []*types.PendingTransfer{},
			Approvals:    []*types.Approval{},
		},
		Message: "Badge balance was not found so returning default balance == 0.",
	}, nil
}
