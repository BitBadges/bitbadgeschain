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
		CanUpdateCollectionApprovedTransfers: true,
	})
	if err != nil {
		return nil, err
	}

	if collection.BalancesType.LTE(sdk.NewUint(0)) {
		return nil, ErrOffChainBalances
	}

	// newApprovedTransfers, needToValidateUpdateCollectionApprovedTransfers := GetApprovedTransfersToStore(collection, msg.ApprovedTransfers)

	// _, err = k.UniversalValidate(ctx, UniversalValidationParams{
	// 	Creator:                              msg.Creator,
	// 	CollectionId:                         msg.CollectionId,
	// 	MustBeManager:                        true,
	// 	CanUpdateCollectionApprovedTransfers: needToValidateUpdateCollectionApprovedTransfers,
	// })

	// err = AssertIsFrozenLogicForApprovedTransfers(collection.ApprovedTransfers, newApprovedTransfers)
	// if err != nil {
	// 	return nil, err
	// }

	// collection.ApprovedTransfers = newApprovedTransfers

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
