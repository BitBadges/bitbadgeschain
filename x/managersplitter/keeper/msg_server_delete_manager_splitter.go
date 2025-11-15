package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/managersplitter/types"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) DeleteManagerSplitter(goCtx context.Context, msg *types.MsgDeleteManagerSplitter) (*types.MsgDeleteManagerSplitterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate admin address
	_, err := sdk.AccAddressFromBech32(msg.Admin)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidAdmin, err.Error())
	}

	// Get existing manager splitter
	managerSplitter, found := k.GetManagerSplitterFromStore(ctx, msg.Address)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrManagerSplitterNotFound, msg.Address)
	}

	// Check authorization - only admin can delete
	if managerSplitter.Admin != msg.Admin {
		return nil, sdkerrors.Wrap(types.ErrUnauthorized, "only admin can delete manager splitter")
	}

	// Delete manager splitter
	k.DeleteManagerSplitterFromStore(ctx, msg.Address)

	return &types.MsgDeleteManagerSplitterResponse{}, nil
}
