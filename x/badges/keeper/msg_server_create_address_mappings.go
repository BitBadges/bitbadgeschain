package keeper

import (
	"context"
	"encoding/json"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	oldtypes "github.com/bitbadges/bitbadgeschain/x/badges/types/v13"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CastOldTypeToNewType(oldMsg *oldtypes.MsgCreateAddressLists) (*types.MsgCreateAddressLists, error) {
	// Convert to JSON
	jsonBytes, err := json.Marshal(oldMsg)
	if err != nil {
		return nil, err
	}

	// Unmarshal into new type
	var newCollection types.MsgCreateAddressLists
	if err := json.Unmarshal(jsonBytes, &newCollection); err != nil {
		return nil, err
	}

	return &newCollection, nil
}

func (k msgServer) CreateAddressListsV13(goCtx context.Context, msg *oldtypes.MsgCreateAddressLists) (*types.MsgCreateAddressListsResponse, error) {
	newMsg, err := CastOldTypeToNewType(msg)
	if err != nil {
		return nil, err
	}
	return k.CreateAddressLists(goCtx, newMsg)
}

func (k msgServer) CreateAddressListsV14(goCtx context.Context, msg *types.MsgCreateAddressLists) (*types.MsgCreateAddressListsResponse, error) {
	return k.CreateAddressLists(goCtx, msg)
}

func (k msgServer) CreateAddressLists(goCtx context.Context, msg *types.MsgCreateAddressLists) (*types.MsgCreateAddressListsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	for _, addressList := range msg.AddressLists {
		addressList.CreatedBy = msg.Creator
		if err := k.CreateAddressList(ctx, addressList); err != nil {
			return nil, err
		}
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
			sdk.NewAttribute("msg_type", "create_address_lists"),
			sdk.NewAttribute("msg", string(msgBytes)),
		),
	)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent("indexer",
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
			sdk.NewAttribute("msg_type", "create_address_lists"),
			sdk.NewAttribute("msg", string(msgBytes)),
		),
	)

	return &types.MsgCreateAddressListsResponse{}, nil
}
