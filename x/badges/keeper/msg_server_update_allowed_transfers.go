package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)



func (k msgServer) UpdateCollectionApprovedTransfers(goCtx context.Context, msg *types.MsgUpdateCollectionApprovedTransfers) (*types.MsgUpdateCollectionApprovedTransfersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	collection, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:                              msg.Creator,
		CollectionId:                         msg.CollectionId,
		MustBeManager:                        true,
	})
	if err != nil {
		return nil, err
	}

	for _, addressMapping := range msg.AddressMappings {
		if err := k.CreateAddressMapping(ctx, addressMapping); err != nil {
			return nil, err
		}
	}

	if msg.ApprovedTransfersTimeline != nil && len(msg.ApprovedTransfersTimeline) > 0 {
		if err := k.ValidateCollectionApprovedTransfersUpdate(ctx, collection, collection.ApprovedTransfersTimeline, msg.ApprovedTransfersTimeline, collection.Permissions.CanUpdateCollectionApprovedTransfers, msg.Creator); err != nil {
			return nil, err
		}
		collection.ApprovedTransfersTimeline = msg.ApprovedTransfersTimeline
	}

	err = k.SetCollectionInStore(ctx, collection)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		),
	)

	return &types.MsgUpdateCollectionApprovedTransfersResponse{}, nil
}
