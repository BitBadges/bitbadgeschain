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

	currProtocol, found := k.GetProtocolFromStore(ctx, protocolName)
	if !found {
		return nil, sdkerrors.Wrap(ErrProtocolDoesNotExist, "Protocol does not exist")
	}

	//Check if user is creator of protocol
	if currProtocol.CreatedBy != msg.Creator {
		return nil, sdkerrors.Wrap(ErrNotProtocolCreator, "Not protocol creator")
	}

	if currProtocol.IsFrozen {
		return nil, sdkerrors.Wrap(ErrProtocolIsFrozen, "Protocol is frozen")
	}

	newProtocol := types.Protocol{
		Name:       msg.Name,
		Uri:        msg.Uri,
		CustomData: msg.CustomData,
		CreatedBy:  msg.Creator,
		IsFrozen: 	msg.IsFrozen,
	}

	//Update protocol in store
	err := k.SetProtocolInStore(ctx, &newProtocol)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "Failed to update protocol in store")
	}

	return &types.MsgUpdateProtocolResponse{}, nil
}
