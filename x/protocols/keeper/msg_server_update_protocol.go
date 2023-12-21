package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/protocols/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkerrors "cosmossdk.io/errors"
)

func (k msgServer) UpdateProtocol(goCtx context.Context, msg *types.MsgUpdateProtocol) (*types.MsgUpdateProtocolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	protocolName := msg.Name

	//Check if protocol exists
	if !k.StoreHasProtocolID(ctx, protocolName) {
		return nil, sdkerrors.Wrap(ErrProtocolDoesNotExist, "Protocol does not exist")
	}

	newProtocol := types.Protocol{
		Name:    		msg.Name,
		Uri: 	 			msg.Uri,
		CustomData: msg.CustomData,
	}

	//Update protocol in store
	err := k.SetProtocolInStore(ctx, &newProtocol)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "Failed to update protocol in store")
	}

	return &types.MsgUpdateProtocolResponse{}, nil
}
