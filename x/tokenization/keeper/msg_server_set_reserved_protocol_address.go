package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetReservedProtocolAddress sets or unsets a reserved protocol address (governance-only).
func (k msgServer) SetReservedProtocolAddress(goCtx context.Context, msg *types.MsgSetReservedProtocolAddress) (*types.MsgSetReservedProtocolAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check that the sender is the authority
	if k.GetAuthority() != msg.Authority {
		return nil, errorsmod.Wrapf(types.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.GetAuthority(), msg.Authority)
	}

	// Validate address
	if err := types.ValidateAddress(msg.Address, false); err != nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidRequest, "invalid address: %s", err)
	}

	// Set or unset the reserved protocol address
	err := k.SetReservedProtocolAddressInStore(ctx, msg.Address, msg.IsReservedProtocol)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "failed to set reserved protocol address")
	}

	return &types.MsgSetReservedProtocolAddressResponse{}, nil
}
