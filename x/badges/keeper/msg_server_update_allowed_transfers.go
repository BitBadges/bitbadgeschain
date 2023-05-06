package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateAllowedTransfers(goCtx context.Context, msg *types.MsgUpdateAllowedTransfers) (*types.MsgUpdateAllowedTransfersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	collection, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:             msg.Creator,
		CollectionId:        msg.CollectionId,
		MustBeManager:       true,
		CanUpdateAllowed:    true,
	})
	if err != nil {
		return nil, err
	}

	collection.AllowedTransfers = msg.AllowedTransfers

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

	return &types.MsgUpdateAllowedTransfersResponse{}, nil
}
