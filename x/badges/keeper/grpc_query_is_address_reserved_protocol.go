package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// IsAddressReservedProtocol queries if an address is a reserved protocol address.
func (k Keeper) IsAddressReservedProtocol(goCtx context.Context, req *types.QueryIsAddressReservedProtocolRequest) (*types.QueryIsAddressReservedProtocolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	isReservedProtocol := k.IsAddressReservedProtocolInStore(ctx, req.Address)

	return &types.QueryIsAddressReservedProtocolResponse{
		IsReservedProtocol: isReservedProtocol,
	}, nil
}
