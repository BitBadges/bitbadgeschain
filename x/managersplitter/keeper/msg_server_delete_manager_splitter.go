package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	"github.com/bitbadges/bitbadgeschain/x/managersplitter/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) DeleteManagerSplitter(goCtx context.Context, msg *types.MsgDeleteManagerSplitter) (*types.MsgDeleteManagerSplitterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate admin address
	_, err := sdk.AccAddressFromBech32(msg.Admin)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidAdmin, err.Error())
	}
	if err := types.ValidateCanonicalAddresses(nil, msg.Admin, msg.Address); err != nil {
		return nil, err
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

	managed, err := k.tokenizationKeeper.HasCollectionsManagedBy(ctx, msg.Address)
	if err != nil {
		return nil, err
	}
	if managed {
		return nil, sdkerrors.Wrap(types.ErrUnauthorized, "transfer collection management before deleting the manager splitter")
	}

	// Delete manager splitter
	k.DeleteManagerSplitterFromStore(ctx, msg.Address)

	return &types.MsgDeleteManagerSplitterResponse{}, nil
}
