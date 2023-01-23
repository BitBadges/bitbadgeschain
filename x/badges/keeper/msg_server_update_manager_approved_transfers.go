package keeper

import (
	"context"
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateManagerApprovedTransfers(goCtx context.Context, msg *types.MsgUpdateManagerApprovedTransfers) (*types.MsgUpdateManagerApprovedTransfersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, badge, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:       msg.Creator,
		CollectionId:       msg.CollectionId,
		MustBeManager: true,
	})
	if err != nil {
		return nil, err
	}

	//TODO: check for only removing addresses

	badge.ManagerApprovedTransfers = msg.ManagerApprovedTransfers

	err = k.SetCollectionInStore(ctx, badge)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeyAction, "UpdateManagerApprovedTransfers"),
			sdk.NewAttribute("BadgeID", fmt.Sprint(msg.CollectionId)),
			sdk.NewAttribute("ManagerApprovedTransfers", fmt.Sprint(msg.ManagerApprovedTransfers)),
		),
	)

	return &types.MsgUpdateManagerApprovedTransfersResponse{}, nil
}
