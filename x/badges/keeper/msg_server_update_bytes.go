package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateBytes(goCtx context.Context, msg *types.MsgUpdateBytes) (*types.MsgUpdateBytesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, badge, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:        msg.Creator,
		CollectionId:   msg.CollectionId,
		MustBeManager:  true,
		CanUpdateBytes: true,
	})
	if err != nil {
		return nil, err
	}

	err = types.ValidateBytes(msg.NewBytes)
	if err != nil {
		return nil, err
	}

	badge.Bytes = msg.NewBytes

	if err := k.SetCollectionInStore(ctx, badge); err != nil {
		return nil, err
	}


	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		),
	)

	return &types.MsgUpdateBytesResponse{}, nil
}
