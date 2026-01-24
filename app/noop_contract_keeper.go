package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibccallbackstypes "github.com/cosmos/ibc-go/v10/modules/apps/callbacks/types"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"
)

// noopContractKeeper is a no-op implementation of ContractKeeper for IBC callbacks
// This is used when callbacks are not needed (e.g., when WasmKeeper doesn't implement ContractKeeper)
type noopContractKeeper struct{}

var _ ibccallbackstypes.ContractKeeper = (*noopContractKeeper)(nil)

// NewNoopContractKeeper creates a new no-op ContractKeeper
func NewNoopContractKeeper() ibccallbackstypes.ContractKeeper {
	return &noopContractKeeper{}
}

// IBCSendPacketCallback implements ibccallbackstypes.ContractKeeper
func (n *noopContractKeeper) IBCSendPacketCallback(
	cachedCtx sdk.Context,
	sourcePort string,
	sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
	packetData []byte,
	contractAddress string,
	packetSenderAddress string,
	version string,
) error {
	return nil
}

// IBCOnAcknowledgementPacketCallback implements ibccallbackstypes.ContractKeeper
func (n *noopContractKeeper) IBCOnAcknowledgementPacketCallback(
	cachedCtx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
	contractAddress string,
	packetSenderAddress string,
	version string,
) error {
	return nil
}

// IBCOnTimeoutPacketCallback implements ibccallbackstypes.ContractKeeper
func (n *noopContractKeeper) IBCOnTimeoutPacketCallback(
	cachedCtx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
	contractAddress string,
	packetSenderAddress string,
	version string,
) error {
	return nil
}

// IBCReceivePacketCallback implements ibccallbackstypes.ContractKeeper
func (n *noopContractKeeper) IBCReceivePacketCallback(
	cachedCtx sdk.Context,
	packet ibcexported.PacketI,
	ack ibcexported.Acknowledgement,
	contractAddress string,
	version string,
) error {
	return nil
}
