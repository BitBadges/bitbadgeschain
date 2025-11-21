package ibc_hooks

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
)

// MustExtractDenomFromPacketOnRecv extracts the denom from an IBC packet on receive
func MustExtractDenomFromPacketOnRecv(packet channeltypes.Packet) string {
	var data transfertypes.FungibleTokenPacketData
	if err := json.Unmarshal(packet.GetData(), &data); err != nil {
		panic(err)
	}

	// Extract denom from the packet
	denom := data.Denom
	if transfertypes.ReceiverChainIsSource(packet.GetSourcePort(), packet.GetSourceChannel(), denom) {
		// remove prefix added by sender chain
		voucherPrefix := transfertypes.GetDenomPrefix(packet.GetSourcePort(), packet.GetSourceChannel())
		unprefixedDenom := denom[len(voucherPrefix):]

		// coin denomination used in sending from the source chain is not prefixed
		denom = unprefixedDenom
	} else {
		// since SendPacket did not prefix denomination, we must prefix denomination here
		denom = transfertypes.GetDenomPrefix(packet.GetDestPort(), packet.GetDestChannel()) + denom
	}

	return denom
}

// IsAckError checks if an acknowledgement is an error
func IsAckError(acknowledgement []byte) bool {
	var ack channeltypes.Acknowledgement
	if err := json.Unmarshal(acknowledgement, &ack); err != nil {
		return true
	}
	return !ack.Success()
}

// NewSuccessAckRepresentingAnError creates a success ack that represents an error
func NewSuccessAckRepresentingAnError(ctx sdk.Context, err errorsmod.Error, errorResponse []byte, errorDescription string) channeltypes.Acknowledgement {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"ibc-hooks-error",
			sdk.NewAttribute("error", err.Error()),
			sdk.NewAttribute("description", errorDescription),
		),
	)
	return channeltypes.NewResultAcknowledgement(errorResponse)
}
