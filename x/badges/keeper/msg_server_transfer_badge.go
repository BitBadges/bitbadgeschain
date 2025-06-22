package keeper

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) TransferBadges(goCtx context.Context, msg *types.MsgTransferBadges) (*types.MsgTransferBadgesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := msg.CheckAndCleanMsg(ctx, true)
	if err != nil {
		return nil, err
	}

	/*
		We support specifying collectionId === 0 for multi-msg transactions where you do not know the collection ID yet upon creation
	*/
	collectionId := msg.CollectionId
	if collectionId.Equal(sdkmath.NewUint(0)) {
		nextCollectionId := k.GetNextCollectionId(ctx)
		collectionId = nextCollectionId.Sub(sdkmath.NewUint(1))
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

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
			sdk.NewAttribute("msg_type", "transfer_badges"),
			sdk.NewAttribute("msg", string(msgBytes)),
			sdk.NewAttribute("collectionId", fmt.Sprint(collectionId)),
		),
	)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent("indexer",
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
			sdk.NewAttribute("msg_type", "transfer_badges"),
			sdk.NewAttribute("msg", string(msgBytes)),
			sdk.NewAttribute("collectionId", fmt.Sprint(collectionId)),
		),
	)

	return &types.MsgTransferBadgesResponse{}, nil
}
