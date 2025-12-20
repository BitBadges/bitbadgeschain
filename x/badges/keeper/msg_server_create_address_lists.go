package keeper

import (
	"context"
	"encoding/json"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateAddressLists(goCtx context.Context, msg *types.MsgCreateAddressLists) (*types.MsgCreateAddressListsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	for _, addressListInput := range msg.AddressLists {
		// Convert AddressListInput to AddressList by adding createdBy field
		addressList := &types.AddressList{
			ListId:     addressListInput.ListId,
			Addresses:  addressListInput.Addresses,
			Whitelist:  addressListInput.Whitelist,
			Uri:        addressListInput.Uri,
			CustomData: addressListInput.CustomData,
			CreatedBy:  msg.Creator,
		}
		if err := k.CreateAddressList(ctx, addressList); err != nil {
			return nil, err
		}
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	EmitMessageAndIndexerEvents(ctx,
		sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		sdk.NewAttribute("msg_type", "create_address_lists"),
		sdk.NewAttribute("msg", string(msgBytes)),
	)

	return &types.MsgCreateAddressListsResponse{}, nil
}
