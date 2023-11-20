package keeper

import (
	"context"
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) TransferBadges(goCtx context.Context, msg *types.MsgTransferBadges) (*types.MsgTransferBadgesResponse, error) {
	err := msg.CheckAndCleanMsg(true)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	collection, found := k.GetCollectionFromStore(ctx, msg.CollectionId)
	if !found {
		return nil, ErrCollectionNotExists
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
			sdk.NewAttribute("collection_id", fmt.Sprint(msg.CollectionId)),
		),
	)

	return &types.MsgTransferBadgesResponse{}, nil
}
