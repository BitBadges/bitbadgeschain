package keeper

import (
	"context"
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateDisallowedTransfers(goCtx context.Context, msg *types.MsgUpdateDisallowedTransfers) (*types.MsgUpdateDisallowedTransfersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, badge, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:             msg.Creator,
		CollectionId:        msg.CollectionId,
		MustBeManager:       true,
		CanUpdateDisallowed: true,
	})
	if err != nil {
		return nil, err
	}

	badge.DisallowedTransfers = msg.DisallowedTransfers

	err = k.SetCollectionInStore(ctx, badge)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeyAction, "UpdateDisallowedTransfers"),
			sdk.NewAttribute("BadgeID", fmt.Sprint(msg.CollectionId)),
			sdk.NewAttribute("DisallowedTransfers", fmt.Sprint(msg.DisallowedTransfers)),
		),
	)

	return &types.MsgUpdateDisallowedTransfersResponse{}, nil
}
