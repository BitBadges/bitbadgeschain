package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateDynamicStore(goCtx context.Context, msg *types.MsgCreateDynamicStore) (*types.MsgCreateDynamicStoreResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate the message
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// Get the next dynamic store ID
	nextStoreId := k.GetNextDynamicStoreId(ctx)
	if nextStoreId.Equal(sdkmath.NewUint(0)) {
		nextStoreId = sdkmath.NewUint(1)
	}

	// Create the dynamic store
	dynamicStore := types.DynamicStore{
		StoreId:       nextStoreId,
		CreatedBy:     msg.Creator,
		DefaultValue:  msg.DefaultValue,
		GlobalEnabled: true, // Default to enabled for new stores
		Uri:           msg.Uri,
		CustomData:    msg.CustomData,
	}

	// Store the dynamic store
	if err := k.SetDynamicStoreInStore(ctx, dynamicStore); err != nil {
		return nil, sdkerrors.Wrap(err, "Failed to store dynamic store")
	}

	// Increment the next dynamic store ID
	if err := k.IncrementNextDynamicStoreId(ctx); err != nil {
		return nil, sdkerrors.Wrap(err, "Failed to increment dynamic store ID")
	}

	// Emit event
	EmitMessageAndIndexerEvents(ctx,
		sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		sdk.NewAttribute("msg_type", "create_dynamic_store"),
		sdk.NewAttribute("store_id", nextStoreId.String()),
	)

	return &types.MsgCreateDynamicStoreResponse{
		StoreId: nextStoreId,
	}, nil
}
