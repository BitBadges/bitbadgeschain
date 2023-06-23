package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) RequestUpdateManager(goCtx context.Context, msg *types.MsgRequestUpdateManager) (*types.MsgRequestUpdateManagerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:                 msg.Creator,
		CollectionId:            msg.CollectionId,
		CanUpdateManager: msg.AddRequest,
	})
	if err != nil {
		return nil, err
	}

	if msg.AddRequest {
		if err := k.CreateUpdateManagerRequest(ctx, msg.CollectionId, msg.Creator); err != nil {
			return nil, err
		}
	} else {
		if err := k.RemoveUpdateManagerRequest(ctx, msg.CollectionId, msg.Creator); err != nil {
			return nil, err
		}
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
		),
	)

	return &types.MsgRequestUpdateManagerResponse{}, nil
}
