package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) TransferManager(goCtx context.Context, msg *types.MsgTransferManager) (*types.MsgTransferManagerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	collection, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:                     msg.Creator,
		CollectionId:                msg.CollectionId,
		MustBeManager:               true,
		CanManagerBeTransferred:     true,
	})
	if err != nil {
		return nil, err
	}

	requested := k.HasAddressRequestedManagerTransfer(ctx, msg.CollectionId, msg.Address)
	if !requested {
		return nil, ErrAddressNeedsToOptInAndRequestManagerTransfer
	}

	collection.Manager = msg.Address

	if err := k.RemoveTransferManagerRequest(ctx, msg.CollectionId, msg.Address); err != nil {
		return nil, err
	}

	if err := k.SetCollectionInStore(ctx, collection); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		),
	)

	return &types.MsgTransferManagerResponse{}, nil
}
