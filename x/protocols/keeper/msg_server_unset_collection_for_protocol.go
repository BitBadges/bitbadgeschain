package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/protocols/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UnsetCollectionForProtocol(goCtx context.Context, msg *types.MsgUnsetCollectionForProtocol) (*types.MsgUnsetCollectionForProtocolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	protocolName := msg.Name
	k.DeleteProtocolCollectionFromStore(ctx, protocolName, msg.Creator)

	return &types.MsgUnsetCollectionForProtocolResponse{}, nil
}
