package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Queries a balance for the given address and badgeId and returns its contents.
func (k Keeper) GetAddressById(goCtx context.Context, req *types.QueryGetAddressByIdRequest) (*types.QueryGetAddressByIdResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	address := k.accountKeeper.GetAccountAddressByID(ctx, req.Id)
	return &types.QueryGetAddressByIdResponse{
		Address: address,
	}, nil
}
