package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) ArchiveCollection(goCtx context.Context, msg *types.MsgArchiveCollection) (*types.MsgArchiveCollectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgArchiveCollectionResponse{}, nil
}
