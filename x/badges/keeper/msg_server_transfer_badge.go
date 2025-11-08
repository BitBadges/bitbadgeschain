package keeper

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// DefaultCollectionId represents the default collection ID used for multi-msg transactions
	DefaultCollectionId = 0
)

func (k msgServer) TransferTokens(goCtx context.Context, msg *types.MsgTransferTokens) (*types.MsgTransferTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := msg.CheckAndCleanMsg(ctx, true)
	if err != nil {
		return nil, err
	}

	/*
		We support specifying collectionId === 0 for multi-msg transactions where you do not know the collection ID yet upon creation
	*/
	collectionId := msg.CollectionId
	if collectionId.Equal(sdkmath.NewUint(DefaultCollectionId)) {
		nextCollectionId := k.GetNextCollectionId(ctx)
		// Check for potential underflow before subtracting
		if nextCollectionId.IsZero() {
			return nil, fmt.Errorf("no collections available: next collection ID is zero")
		}
		collectionId = nextCollectionId.Sub(sdkmath.NewUint(1))
	}

	collection, found := k.GetCollectionFromStore(ctx, collectionId)
	if !found {
		return nil, ErrCollectionNotExists
	}

	if err := k.Keeper.HandleTransfers(ctx, collection, msg.Transfers, msg.Creator); err != nil {
		return nil, err
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	EmitMessageAndIndexerEvents(ctx,
		sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		sdk.NewAttribute("msg_type", "transfer_tokens"),
		sdk.NewAttribute("msg", string(msgBytes)),
		sdk.NewAttribute("collectionId", fmt.Sprint(collectionId)),
	)

	return &types.MsgTransferTokensResponse{}, nil
}
