package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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
	EmitMessageAndIndexerEvents(ctx,
		sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		sdk.NewAttribute("msg_type", "update_dynamic_store"),
		sdk.NewAttribute("store_id", msg.StoreId.String()),
	)

	return &types.MsgUpdateDynamicStoreResponse{}, nil
}
