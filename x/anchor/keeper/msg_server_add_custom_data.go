package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/anchor/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) AddCustomData(goCtx context.Context, msg *types.MsgAddCustomData) (*types.MsgAddCustomDataResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	nextLocationId := k.GetNextAnchorId(ctx)
	k.SetAnchorLocation(ctx, nextLocationId, msg.Data, msg.Creator)

	return &types.MsgAddCustomDataResponse{
		LocationId: sdk.NewUint(nextLocationId.Uint64()),
	}, nil
}
