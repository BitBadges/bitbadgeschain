package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateManager(goCtx context.Context, msg *types.MsgUpdateManager) (*types.MsgUpdateManagerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	collection, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:                 msg.Creator,
		CollectionId:            msg.CollectionId,
		MustBeManager:           true,
	})
	if err != nil {
		return nil, err
	}

	if err := k.ValidateManagerUpdate(ctx, collection.ManagerTimeline, msg.ManagerTimeline, collection.Permissions.CanUpdateManager); err != nil {
		return nil, err
	}
	collection.ManagerTimeline = msg.ManagerTimeline


	if err := k.SetCollectionInStore(ctx, collection); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		),
	)

	return &types.MsgUpdateManagerResponse{}, nil
}
