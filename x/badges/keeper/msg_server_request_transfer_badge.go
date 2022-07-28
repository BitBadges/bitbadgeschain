package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) RequestTransferBadge(goCtx context.Context, msg *types.MsgRequestTransferBadge) (*types.MsgRequestTransferBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	to, err := sdk.AccAddressFromBech32(msg.To)
	if err != nil {
		return nil, err
	}

	//TODO: add msg.Creator === from here?
	err = k.Keeper.RequestTransferBadge(ctx, from, to, msg.Amount, msg.BadgeId, msg.SubbadgeId)
	if err != nil {
		return nil, err
	}

	_ = ctx
	return &types.MsgRequestTransferBadgeResponse{
		Message: "Success!",
	}, nil
}
