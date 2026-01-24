package app

import (
	errorsmod "cosmossdk.io/errors"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	channelkeeper "github.com/cosmos/ibc-go/v10/modules/core/04-channel/keeper"
	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

// channelKeeperV2Adapter implements wasmd's ChannelKeeperV2 interface for IBC v2 (async packets).
// This is required by wasmd's NewKeeper, but IBC v2 is not yet supported in this chain.
// The adapter returns an error if IBC v2 is attempted, making the limitation explicit.
type channelKeeperV2Adapter struct{}

var _ wasmtypes.ChannelKeeperV2 = (*channelKeeperV2Adapter)(nil)

// NewChannelKeeperV2Adapter creates a new adapter for ChannelKeeperV2.
// Required by wasmd's NewKeeper, but IBC v2 (async packets) is not yet supported.
func NewChannelKeeperV2Adapter(_ *channelkeeper.Keeper) wasmtypes.ChannelKeeperV2 {
	return &channelKeeperV2Adapter{}
}

// WriteAcknowledgement implements wasmtypes.ChannelKeeperV2.
// IBC v2 (async packets) is not supported - returns an error to prevent silent failures.
func (*channelKeeperV2Adapter) WriteAcknowledgement(
	ctx sdk.Context,
	clientID string,
	sequence uint64,
	_ channeltypesv2.Acknowledgement,
) error {
	return errorsmod.Wrapf(
		channeltypes.ErrInvalidChannelState,
		"IBC v2 (async packets) not supported: cannot write acknowledgement for clientID=%s sequence=%d",
		clientID, sequence,
	)
}
