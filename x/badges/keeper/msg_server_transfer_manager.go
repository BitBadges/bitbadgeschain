package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) TransferManager(goCtx context.Context, msg *types.MsgTransferManager) (*types.MsgTransferManagerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, badge, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator: msg.Creator,
		BadgeId: msg.BadgeId,
		MustBeManager: true,
		CanManagerTransfer: true,
	})
	if err != nil {
		return nil, err
	}

	requested := k.HasAddressRequestedManagerTransfer(ctx, msg.BadgeId, msg.Address)
	if !requested {
		return nil, ErrAddressNeedsToOptInAndRequestManagerTransfer
	}

	badge.Manager = msg.Address

	if err := k.RemoveTransferManagerRequest(ctx, msg.BadgeId, msg.Address); err != nil {
		return nil, err
	}

	if err := k.SetBadgeInStore(ctx, badge); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeyAction, "TransferManager"),
			sdk.NewAttribute("BadgeId", fmt.Sprint(msg.BadgeId)),
		),
	)

	return &types.MsgTransferManagerResponse{}, nil
}
