package keeper

import (
	"context"
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) RequestTransferManager(goCtx context.Context, msg *types.MsgRequestTransferManager) (*types.MsgRequestTransferManagerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	CreatorAccountNum, badge, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator: msg.Creator,
		BadgeId: msg.BadgeId,
	})
	if err != nil {
		return nil, err
	}

	if msg.Add {
		permissions := types.GetPermissions(badge.Permissions)
		if !permissions.CanManagerTransfer {
			return nil, ErrInvalidPermissions //Manager can never transfer, so we don't unnecessarily store stuff
		}

		if err := k.CreateTransferManagerRequest(ctx, msg.BadgeId, CreatorAccountNum); err != nil {
			return nil, err
		}
	} else {
		if err := k.RemoveTransferManagerRequest(ctx, msg.BadgeId, CreatorAccountNum); err != nil {
			return nil, err
		}
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeyAction, "RequestTransferManager"),
			sdk.NewAttribute("Creator", fmt.Sprint(CreatorAccountNum)),
			sdk.NewAttribute("BadgeId", fmt.Sprint(msg.BadgeId)),
			sdk.NewAttribute("AddRequest", fmt.Sprint(msg.Add)),
		),
	)

	return &types.MsgRequestTransferManagerResponse{}, nil
}
