package keeper

import (
	"context"
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) TransferBadges(goCtx context.Context, msg *types.MsgTransferBadges) (*types.MsgTransferBadgesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := msg.CheckAndCleanMsg(ctx, true)
	if err != nil {
		return nil, err
	}
	collectionId := msg.CollectionId
	if collectionId.Equal(sdk.NewUint(0)) {
		//Get next collection id - 1 from badges keeper
		nextCollectionId := k.GetNextCollectionId(ctx)

		collectionId = nextCollectionId.Sub(sdk.NewUint(1))
	}

	collection, found := k.GetCollectionFromStore(ctx, collectionId)
	if !found {
		return nil, ErrCollectionNotExists
	}

	isArchived := types.GetIsArchived(ctx, collection)
	if isArchived {
		return nil, ErrCollectionIsArchived
	}

	if !IsStandardBalances(collection) {
		return nil, ErrWrongBalancesType
	}

	if err := k.Keeper.HandleTransfers(ctx, collection, msg.Transfers, msg.Creator); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute("collection_id", fmt.Sprint(collectionId)),
		),
	)

	return &types.MsgTransferBadgesResponse{}, nil
}
