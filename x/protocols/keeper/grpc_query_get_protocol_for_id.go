package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/protocols/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Queries a balance for the given address and badgeId and returns its contents.
func (k Keeper) GetCollectionIdForProtocol(goCtx context.Context, req *types.QueryGetCollectionIdForProtocolRequest) (*types.QueryGetCollectionIdForProtocolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	collectionId := k.GetProtocolCollectionFromStore(ctx, req.Name, req.Address)
	return &types.QueryGetCollectionIdForProtocolResponse{
		CollectionId: collectionId,
	}, nil
}