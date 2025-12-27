package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// SendWithAliasRouting handles MsgSendWithAliasRouting by routing through sendmanager
// to support both standard coins and alias denoms (e.g., badgeslp:).
func (k msgServer) SendWithAliasRouting(goCtx context.Context, msg *types.MsgSendWithAliasRouting) (*types.MsgSendWithAliasRoutingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	fromAddress, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil, sdkerrors.Wrapf(errortypes.ErrInvalidAddress, "invalid from address: %s", err)
	}

	toAddress, err := sdk.AccAddressFromBech32(msg.ToAddress)
	if err != nil {
		return nil, sdkerrors.Wrapf(errortypes.ErrInvalidAddress, "invalid to address: %s", err)
	}

	coins := msg.Amount
	if err := coins.Validate(); err != nil {
		return nil, sdkerrors.Wrapf(errortypes.ErrInvalidCoins, "invalid coins: %s", err)
	}

	// Validate that coins is not empty
	if coins.IsZero() {
		return nil, sdkerrors.Wrapf(errortypes.ErrInvalidCoins, "coins cannot be empty")
	}

	// Use sendmanager's SendCoinsWithAliasRouting which handles both standard coins and alias denoms
	if err := k.SendCoinsWithAliasRouting(ctx, fromAddress, toAddress, coins); err != nil {
		return nil, err
	}

	return &types.MsgSendWithAliasRoutingResponse{}, nil
}
