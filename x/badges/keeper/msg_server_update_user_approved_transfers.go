package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateUserApprovedTransfers(goCtx context.Context, msg *types.MsgUpdateUserApprovedTransfers) (*types.MsgUpdateUserApprovedTransfersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgUpdateUserApprovedTransfersResponse{}, nil
}
