package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) RequestTransferManager(goCtx context.Context, msg *types.MsgRequestTransferManager) (*types.MsgRequestTransferManagerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	CreatorAccountNum, badge, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:      msg.Creator,
		CollectionId: msg.CollectionId,
	})
	if err != nil {
		return nil, err
	}

	if msg.AddRequest {
		permissions := types.GetPermissions(badge.Permissions)
		if !permissions.CanManagerBeTransferred {
			return nil, ErrInvalidPermissions //Manager can never transfer, so we don't unnecessarily store stuff
		}

		if err := k.CreateTransferManagerRequest(ctx, msg.CollectionId, CreatorAccountNum); err != nil {
			return nil, err
		}
	} else {
		if err := k.RemoveTransferManagerRequest(ctx, msg.CollectionId, CreatorAccountNum); err != nil {
			return nil, err
		}
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
		),
	)

	return &types.MsgRequestTransferManagerResponse{}, nil
}
