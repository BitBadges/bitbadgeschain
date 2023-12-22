package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/protocols/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkerrors "cosmossdk.io/errors"
)

func (k msgServer) CreateProtocol(goCtx context.Context, msg *types.MsgCreateProtocol) (*types.MsgCreateProtocolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_ = ctx

	protocolToAdd := types.Protocol{
		Name:       msg.Name,
		Uri:        msg.Uri,
		CustomData: msg.CustomData,
		CreatedBy:  msg.Creator,
	}

	//Check if protocol already exists
	if k.StoreHasProtocolID(ctx, msg.Name) {
		return nil, sdkerrors.Wrap(ErrProtocolExists, "Protocol already exists")
	}

	//Add protocol to store
	err := k.SetProtocolInStore(ctx, &protocolToAdd)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "Failed to add protocol to store")
	}

	return &types.MsgCreateProtocolResponse{}, nil
}
