package keeper

import (
	"context"

	"bitbadgeschain/x/badges/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateAddressLists(goCtx context.Context, msg *types.MsgCreateAddressLists) (*types.MsgCreateAddressListsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := k.UniversalValidateNotHalted(ctx)
	if err != nil {
		return nil, err
	}

	for _, addressList := range msg.AddressLists {
		addressList.CreatedBy = msg.Creator
		if err := k.CreateAddressList(ctx, addressList); err != nil {
			return nil, err
		}
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		),
	)

	return &types.MsgCreateAddressListsResponse{}, nil
}
