package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetAllReservedProtocolAddresses queries all reserved protocol addresses.
func (k Keeper) GetAllReservedProtocolAddresses(goCtx context.Context, req *types.QueryGetAllReservedProtocolAddressesRequest) (*types.QueryGetAllReservedProtocolAddressesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	addresses := k.GetAllReservedProtocolAddressesFromStore(ctx)

	return &types.QueryGetAllReservedProtocolAddressesResponse{
		Addresses: addresses,
	}, nil
}
