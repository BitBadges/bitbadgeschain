package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/managersplitter/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateManagerSplitter(goCtx context.Context, msg *types.MsgCreateManagerSplitter) (*types.MsgCreateManagerSplitterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate admin address
	_, err := sdk.AccAddressFromBech32(msg.Admin)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidAdmin, err.Error())
	}
	if err := types.ValidateCanonicalAddresses(msg.Permissions, msg.Admin); err != nil {
		return nil, err
	}

	// Get next ID and derive address
	nextId := k.GetNextManagerSplitterId(ctx)
	address := types.DeriveManagerSplitterAddress(nextId)

	// Check if address already exists (shouldn't happen, but be safe)
	_, exists := k.GetManagerSplitterFromStore(ctx, address)
	if exists {
		return nil, sdkerrors.Wrap(types.ErrManagerSplitterExists, "manager splitter already exists")
	}

	// Create manager splitter
	managerSplitter := &types.ManagerSplitter{
		Address:     address,
		Admin:       msg.Admin,
		Permissions: msg.Permissions,
	}

	// Set permissions to empty if nil
	if managerSplitter.Permissions == nil {
		managerSplitter.Permissions = &types.ManagerSplitterPermissions{}
	}

	// Store manager splitter
	err = k.SetManagerSplitterInStore(ctx, managerSplitter)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to store manager splitter")
	}

	// Increment next ID
	k.SetNextManagerSplitterId(ctx, nextId.Add(sdkmath.NewUint(1)))

	return &types.MsgCreateManagerSplitterResponse{
		Address: address,
	}, nil
}
