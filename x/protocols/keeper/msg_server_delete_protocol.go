package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/protocols/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkerrors "cosmossdk.io/errors"
)

func (k msgServer) DeleteProtocol(goCtx context.Context, msg *types.MsgDeleteProtocol) (*types.MsgDeleteProtocolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	protocolName := msg.Name

	//Check if protocol exists
	if !k.StoreHasProtocolID(ctx, protocolName) {
		return nil, sdkerrors.Wrap(ErrProtocolDoesNotExist, "Protocol does not exist")
	}

	//Delete protocol from store
	k.DeleteProtocolFromStore(ctx, protocolName)

	return &types.MsgDeleteProtocolResponse{}, nil
}
