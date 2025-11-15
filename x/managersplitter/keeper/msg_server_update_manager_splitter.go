package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/managersplitter/types"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateManagerSplitter(goCtx context.Context, msg *types.MsgUpdateManagerSplitter) (*types.MsgUpdateManagerSplitterResponse, error) {
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

	// Check authorization - only admin can update
	if managerSplitter.Admin != msg.Admin {
		return nil, sdkerrors.Wrap(types.ErrUnauthorized, "only admin can update manager splitter")
	}

	// Update permissions
	if msg.Permissions != nil {
		managerSplitter.Permissions = msg.Permissions
	}

	// Store updated manager splitter
	err = k.SetManagerSplitterInStore(ctx, managerSplitter)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to update manager splitter")
	}

	return &types.MsgUpdateManagerSplitterResponse{}, nil
}

