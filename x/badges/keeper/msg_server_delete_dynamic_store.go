package keeper

import (
	"context"
	"encoding/json"

	sdkerrors "cosmossdk.io/errors"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	oldtypes "github.com/bitbadges/bitbadgeschain/x/badges/types/v13"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CastOldDeleteDynamicStoreToNewType(oldMsg *oldtypes.MsgDeleteDynamicStore) (*types.MsgDeleteDynamicStore, error) {
	// Convert to JSON
	jsonBytes, err := json.Marshal(oldMsg)
	if err != nil {
		return nil, err
	}

	// Unmarshal into new type
	var newMsg types.MsgDeleteDynamicStore
	if err := json.Unmarshal(jsonBytes, &newMsg); err != nil {
		return nil, err
	}

	return &newMsg, nil
}

func (k msgServer) DeleteDynamicStoreV13(goCtx context.Context, msg *oldtypes.MsgDeleteDynamicStore) (*types.MsgDeleteDynamicStoreResponse, error) {
	newMsg, err := CastOldDeleteDynamicStoreToNewType(msg)
	if err != nil {
		return nil, err
	}
	return k.DeleteDynamicStore(goCtx, newMsg)
}

func (k msgServer) DeleteDynamicStoreV14(goCtx context.Context, msg *types.MsgDeleteDynamicStore) (*types.MsgDeleteDynamicStoreResponse, error) {
	return k.DeleteDynamicStore(goCtx, msg)
}

func (k msgServer) DeleteDynamicStore(goCtx context.Context, msg *types.MsgDeleteDynamicStore) (*types.MsgDeleteDynamicStoreResponse, error) {
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
		return nil, sdkerrors.Wrap(types.ErrInvalidDynamicStoreID, "Only the creator can delete the dynamic store")
	}

	// Delete the dynamic store
	k.DeleteDynamicStoreFromStore(ctx, msg.StoreId)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
			sdk.NewAttribute("msg_type", "delete_dynamic_store"),
			sdk.NewAttribute("store_id", msg.StoreId.String()),
		),
	)

	return &types.MsgDeleteDynamicStoreResponse{}, nil
}
