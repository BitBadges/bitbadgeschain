package keeper

import (
	"context"
	"encoding/json"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateAddressLists(goCtx context.Context, msg *types.MsgCreateAddressLists) (*types.MsgCreateAddressListsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	creator, err := k.GetCreator(ctx, msg.Creator, msg.CreatorOverride)
	if err != nil {
		return nil, err
	}
	msg.Creator = creator

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
