package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"
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

	if collection.BalancesType != sdkmath.NewUint(0) {
		return nil, ErrOffChainBalances
	}

	for _, addressMapping := range msg.AddressMappings {
		if err := k.CreateAddressMapping(ctx, addressMapping); err != nil {
			return nil, err
		}
	}

	if err := ValidateCollectionApprovedTransfersUpdate(ctx, collection.ApprovedTransfersTimeline, msg.ApprovedTransfersTimeline, collection.Permissions.CanUpdateCollectionApprovedTransfers); err != nil {
		return nil, err
	}
	collection.ApprovedTransfersTimeline = msg.ApprovedTransfersTimeline

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
