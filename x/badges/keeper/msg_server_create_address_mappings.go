package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateAddressMappings(goCtx context.Context, msg *types.MsgCreateAddressMappings) (*types.MsgCreateAddressMappingsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	for _, addressMapping := range msg.AddressMappings {
		addressMapping.CreatedBy = msg.Creator
		if err := k.CreateAddressMapping(ctx, addressMapping); err != nil {
			return nil, err
		}
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		),
	)

	return &types.MsgCreateAddressMappingsResponse{}, nil
}
