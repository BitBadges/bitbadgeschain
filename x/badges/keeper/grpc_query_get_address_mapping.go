package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Queries a balance for the given address and badgeId and returns its contents.
func (k Keeper) GetAddressMapping(goCtx context.Context, req *types.QueryGetAddressMappingRequest) (*types.QueryGetAddressMappingResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	if req.MappingId == "Manager" {
		return nil, status.Error(codes.InvalidArgument, "invalid request. this query does not support the manager mapping")
	}

	addressMapping, err := k.GetAddressMappingById(ctx, req.MappingId, "")
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	return &types.QueryGetAddressMappingResponse{
		Mapping: addressMapping,
	}, nil
}
