package keeper

import (
	"context"
	"encoding/json"

	sdkerrors "cosmossdk.io/errors"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	oldtypes "github.com/bitbadges/bitbadgeschain/x/badges/types/v13"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CastOldUpdateDynamicStoreToNewType(oldMsg *oldtypes.MsgUpdateDynamicStore) (*types.MsgUpdateDynamicStore, error) {
	// Convert to JSON
	jsonBytes, err := json.Marshal(oldMsg)
	if err != nil {
		return nil, err
	}

	// Unmarshal into new type
	var newMsg types.MsgUpdateDynamicStore
	if err := json.Unmarshal(jsonBytes, &newMsg); err != nil {
		return nil, err
	}

	return &newMsg, nil
}

func (k msgServer) UpdateDynamicStoreV13(goCtx context.Context, msg *oldtypes.MsgUpdateDynamicStore) (*types.MsgUpdateDynamicStoreResponse, error) {
	newMsg, err := CastOldUpdateDynamicStoreToNewType(msg)
	if err != nil {
		return nil, err
	}
	return k.UpdateDynamicStore(goCtx, newMsg)
}

func (k msgServer) UpdateDynamicStoreV14(goCtx context.Context, msg *types.MsgUpdateDynamicStore) (*types.MsgUpdateDynamicStoreResponse, error) {
	return k.UpdateDynamicStore(goCtx, msg)
}

func (k msgServer) UpdateDynamicStore(goCtx context.Context, msg *types.MsgUpdateDynamicStore) (*types.MsgUpdateDynamicStoreResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate the message
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// Get the existing dynamic store
	dynamicStore, found := k.GetDynamicStoreFromStore(ctx, msg.StoreId)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrInvalidDynamicStoreID, "Dynamic store not found")
	}

	// Check if the creator is the owner
	if dynamicStore.CreatedBy != msg.Creator {
		return nil, sdkerrors.Wrap(types.ErrInvalidDynamicStoreID, "Only the creator can update the dynamic store")
	}

	// Update the default value if set (always set, since proto3 bools default to false)
	dynamicStore.DefaultValue = msg.DefaultValue

	// Store the updated dynamic store
	if err := k.SetDynamicStoreInStore(ctx, dynamicStore); err != nil {
		return nil, sdkerrors.Wrap(err, "Failed to store dynamic store")
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
			sdk.NewAttribute("msg_type", "update_dynamic_store"),
			sdk.NewAttribute("store_id", msg.StoreId.String()),
		),
	)

	return &types.MsgUpdateDynamicStoreResponse{}, nil
}
