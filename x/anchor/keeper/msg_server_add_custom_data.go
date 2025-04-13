package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/anchor/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) AddCustomData(goCtx context.Context, msg *types.MsgAddCustomData) (*types.MsgAddCustomDataResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	nextLocationId, err := k.GetNextAnchorId(ctx)
	if err != nil {
		return nil, err
	}

	k.SetAnchorLocation(ctx, nextLocationId, msg.Data, msg.Creator)

	return &types.MsgAddCustomDataResponse{
		LocationId: sdkmath.NewUint(nextLocationId.Uint64()),
	}, nil
}
