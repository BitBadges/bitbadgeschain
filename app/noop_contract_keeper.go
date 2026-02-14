package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibccallbackstypes "github.com/cosmos/ibc-go/v10/modules/apps/callbacks/types"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"
)

// NoopContractKeeper is a no-op implementation of ibccallbackstypes.ContractKeeper
// This is used for IBC callbacks middleware (no contract callbacks are executed)
type NoopContractKeeper struct{}

// NewNoopContractKeeper creates a new no-op contract keeper
func NewNoopContractKeeper() *NoopContractKeeper {
	return &NoopContractKeeper{}
}

// IBCSendPacketCallback implements ibccallbackstypes.ContractKeeper
func (k *NoopContractKeeper) IBCSendPacketCallback(ctx sdk.Context, sourcePort string, sourceChannel string, timeoutHeight clienttypes.Height, timeoutTimestamp uint64, packetData []byte, contractAddress string, callbackGasLimit string, sourceChannelID string) error {
	return nil
}

// IBCOnAcknowledgementPacketCallback implements ibccallbackstypes.ContractKeeper
func (k *NoopContractKeeper) IBCOnAcknowledgementPacketCallback(ctx sdk.Context, packet channeltypes.Packet, acknowledgement []byte, relayer sdk.AccAddress, contractAddress string, callbackGasLimit string, sourcePort string) error {
	return nil
}

// IBCOnTimeoutPacketCallback implements ibccallbackstypes.ContractKeeper
func (k *NoopContractKeeper) IBCOnTimeoutPacketCallback(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress, contractAddress string, callbackGasLimit string, sourcePort string) error {
	return nil
}

// IBCReceivePacketCallback implements ibccallbackstypes.ContractKeeper
func (k *NoopContractKeeper) IBCReceivePacketCallback(ctx sdk.Context, packet ibcexported.PacketI, acknowledgement ibcexported.Acknowledgement, contractAddress string, callbackGasLimit string) error {
	return nil
}

// Ensure NoopContractKeeper implements ibccallbackstypes.ContractKeeper
var _ ibccallbackstypes.ContractKeeper = (*NoopContractKeeper)(nil)

