package keeper

import (
	"context"
	"encoding/json"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	oldtypes "github.com/bitbadges/bitbadgeschain/x/badges/types/v13"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CastOldCreateDynamicStoreToNewType(oldMsg *oldtypes.MsgCreateDynamicStore) (*types.MsgCreateDynamicStore, error) {
	// Convert to JSON
	jsonBytes, err := json.Marshal(oldMsg)
	if err != nil {
		return nil, err
	}

	// Unmarshal into new type
	var newMsg types.MsgCreateDynamicStore
	if err := json.Unmarshal(jsonBytes, &newMsg); err != nil {
		return nil, err
	}

	return &newMsg, nil
}

func (k msgServer) CreateDynamicStoreV13(goCtx context.Context, msg *oldtypes.MsgCreateDynamicStore) (*types.MsgCreateDynamicStoreResponse, error) {
	newMsg, err := CastOldCreateDynamicStoreToNewType(msg)
	if err != nil {
		return nil, err
	}
	return k.CreateDynamicStore(goCtx, newMsg)
}

func (k msgServer) CreateDynamicStoreV14(goCtx context.Context, msg *types.MsgCreateDynamicStore) (*types.MsgCreateDynamicStoreResponse, error) {
	return k.CreateDynamicStore(goCtx, msg)
}

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
		StoreId:      nextStoreId,
		CreatedBy:    msg.Creator,
		DefaultValue: msg.DefaultValue,
	}

	// Store the dynamic store
	if err := k.SetDynamicStoreInStore(ctx, dynamicStore); err != nil {
		return nil, sdkerrors.Wrap(err, "Failed to store dynamic store")
	}

	// Increment the next dynamic store ID
	k.IncrementNextDynamicStoreId(ctx)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
			sdk.NewAttribute("msg_type", "create_dynamic_store"),
			sdk.NewAttribute("store_id", nextStoreId.String()),
		),
	)

	return &types.MsgCreateDynamicStoreResponse{
		StoreId: nextStoreId,
	}, nil
}
