package testutil

import (
	"fmt"

	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"
)

// GenerateTestPacket generates a test packet for testing
func GenerateTestPacket(portID, channelID string, data []byte) channeltypes.Packet {
	return channeltypes.Packet{
		Sequence:           1,
		SourcePort:         portID,
		SourceChannel:      channelID,
		DestinationPort:    "transfer",
		DestinationChannel: "channel-1",
		Data:               data,
		TimeoutHeight:      clienttypes.Height{},
		TimeoutTimestamp:   0,
	}
}

// GenerateTestCapability generates a test capability for testing
func GenerateTestCapability() *capabilitytypes.Capability {
	return &capabilitytypes.Capability{
		Index: 1,
	}
}

// GenerateTestAcknowledgement generates a test acknowledgement
func GenerateTestAcknowledgement(success bool) ibcexported.Acknowledgement {
	if success {
		return channeltypes.NewResultAcknowledgement([]byte("success"))
	}
	return channeltypes.NewErrorAcknowledgement(fmt.Errorf("error"))
}

