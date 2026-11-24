package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	"github.com/bitbadges/bitbadgeschain/x/managersplitter/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateManagerSplitter(goCtx context.Context, msg *types.MsgUpdateManagerSplitter) (*types.MsgUpdateManagerSplitterResponse, error) {
	if msg == nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidRequest, "message cannot be nil")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate admin address
	adminAddr, err := sdk.AccAddressFromBech32(msg.Admin)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidAdmin, err.Error())
	}
	if adminAddr.Empty() {
		return nil, sdkerrors.Wrap(types.ErrInvalidAdmin, "admin address cannot be empty")
	}
	addr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil || addr.Empty() {
		return nil, sdkerrors.Wrap(types.ErrInvalidAddress, "invalid manager splitter address")
	}
	if msg.Admin == msg.Address {
		return nil, sdkerrors.Wrap(types.ErrInvalidRequest, "admin and manager splitter address cannot be the same")
	}
	if err := types.ValidateOptionalPermissions(msg.Permissions); err != nil {
		return nil, err
	}
	if err := types.ValidateCanonicalAddresses(msg.Permissions, msg.Admin, msg.Address); err != nil {
		return nil, err
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
