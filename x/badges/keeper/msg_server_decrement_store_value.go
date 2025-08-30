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

func CastOldDecrementStoreValueToNewType(oldMsg *oldtypes.MsgDecrementStoreValue) (*types.MsgDecrementStoreValue, error) {
	// Convert to JSON
	jsonBytes, err := json.Marshal(oldMsg)
	if err != nil {
		return nil, err
	}

	// Unmarshal into new type
	var newMsg types.MsgDecrementStoreValue
	if err := json.Unmarshal(jsonBytes, &newMsg); err != nil {
		return nil, err
	}

	return &newMsg, nil
}

func (k msgServer) DecrementStoreValueV13(goCtx context.Context, msg *oldtypes.MsgDecrementStoreValue) (*types.MsgDecrementStoreValueResponse, error) {
	newMsg, err := CastOldDecrementStoreValueToNewType(msg)
	if err != nil {
		return nil, err
	}
	return k.DecrementStoreValue(goCtx, newMsg)
}

func (k msgServer) DecrementStoreValueV14(goCtx context.Context, msg *types.MsgDecrementStoreValue) (*types.MsgDecrementStoreValueResponse, error) {
	return k.DecrementStoreValue(goCtx, msg)
}

func (k msgServer) DecrementStoreValue(goCtx context.Context, msg *types.MsgDecrementStoreValue) (*types.MsgDecrementStoreValueResponse, error) {
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
		return nil, sdkerrors.Wrap(types.ErrInvalidDynamicStoreID, "Only the creator can decrement values in the dynamic store")
	}

	// Validate the address
	if err := types.ValidateAddress(msg.Address, false); err != nil {
		return nil, sdkerrors.Wrap(err, "Invalid address")
	}

	// Get current value
	dynamicStoreValue, found := k.GetDynamicStoreValueFromStore(ctx, msg.StoreId, msg.Address)
	var currentValue sdkmath.Uint
	if found {
		currentValue = dynamicStoreValue.Value
	} else {
		currentValue = dynamicStore.DefaultValue
	}

	// Decrement the value
	var newValue sdkmath.Uint
	if currentValue.LT(msg.Amount) {
		// Underflow case
		if msg.SetToZeroOnUnderflow {
			newValue = sdkmath.NewUint(0)
		} else {
			return nil, sdkerrors.Wrapf(types.ErrInvalidRequest, "decrement would cause underflow: current value %s, decrement amount %s", currentValue.String(), msg.Amount.String())
		}
	} else {
		newValue = currentValue.Sub(msg.Amount)
	}

	// Set the new value
	if err := k.SetDynamicStoreValueInStore(ctx, msg.StoreId, msg.Address, newValue); err != nil {
		return nil, sdkerrors.Wrap(err, "Failed to store decremented dynamic store value")
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
			sdk.NewAttribute("msg_type", "decrement_store_value"),
			sdk.NewAttribute("store_id", msg.StoreId.String()),
			sdk.NewAttribute("address", msg.Address),
			sdk.NewAttribute("amount", msg.Amount.String()),
			sdk.NewAttribute("new_value", newValue.String()),
			sdk.NewAttribute("set_to_zero_on_underflow", func() string {
				if msg.SetToZeroOnUnderflow {
					return "true"
				}
				return "false"
			}()),
		),
	)

	return &types.MsgDecrementStoreValueResponse{}, nil
}
