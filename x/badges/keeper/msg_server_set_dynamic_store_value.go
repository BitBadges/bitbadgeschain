package keeper

import (
	"context"
	"encoding/json"

	sdkerrors "cosmossdk.io/errors"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	oldtypes "github.com/bitbadges/bitbadgeschain/x/badges/types/v13"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CastOldSetDynamicStoreValueToNewType(oldMsg *oldtypes.MsgSetDynamicStoreValue) (*types.MsgSetDynamicStoreValue, error) {
	// Convert to JSON
	jsonBytes, err := json.Marshal(oldMsg)
	if err != nil {
		return nil, err
	}

	// Unmarshal into new type
	var newMsg types.MsgSetDynamicStoreValue
	if err := json.Unmarshal(jsonBytes, &newMsg); err != nil {
		return nil, err
	}

	return &newMsg, nil
}

func (k msgServer) SetDynamicStoreValueV13(goCtx context.Context, msg *oldtypes.MsgSetDynamicStoreValue) (*types.MsgSetDynamicStoreValueResponse, error) {
	newMsg, err := CastOldSetDynamicStoreValueToNewType(msg)
	if err != nil {
		return nil, err
	}
	return k.SetDynamicStoreValue(goCtx, newMsg)
}

func (k msgServer) SetDynamicStoreValueV14(goCtx context.Context, msg *types.MsgSetDynamicStoreValue) (*types.MsgSetDynamicStoreValueResponse, error) {
	return k.SetDynamicStoreValue(goCtx, msg)
}

func (k msgServer) SetDynamicStoreValue(goCtx context.Context, msg *types.MsgSetDynamicStoreValue) (*types.MsgSetDynamicStoreValueResponse, error) {
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
		return nil, sdkerrors.Wrap(types.ErrInvalidDynamicStoreID, "Only the creator can set values in the dynamic store")
	}

	// Validate the address
	if err := types.ValidateAddress(msg.Address, false); err != nil {
		return nil, sdkerrors.Wrap(err, "Invalid address")
	}

	// Set the dynamic store value
	if err := k.SetDynamicStoreValueInStore(ctx, msg.StoreId, msg.Address, msg.Value); err != nil {
		return nil, sdkerrors.Wrap(err, "Failed to store dynamic store value")
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
			sdk.NewAttribute("msg_type", "set_dynamic_store_value"),
			sdk.NewAttribute("store_id", msg.StoreId.String()),
			sdk.NewAttribute("address", msg.Address),
			sdk.NewAttribute("value", msg.Value.String()),
		),
	)

	return &types.MsgSetDynamicStoreValueResponse{}, nil
}
