package customhooks

import (
	"testing"

	transfertypes "github.com/cosmos/ibc-go/v11/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v11/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v11/modules/core/04-channel/types"
	"github.com/stretchr/testify/require"
)

func makePacket(denom string, sourcePort, sourceChannel, destPort, destChannel string) channeltypes.Packet {
	packetData := transfertypes.FungibleTokenPacketData{
		Denom:    denom,
		Amount:   "1000",
		Sender:   "cosmos1sender",
		Receiver: "bb1receiver",
	}
	data, _ := transfertypes.ModuleCdc.MarshalJSON(&packetData)
	return channeltypes.Packet{
		Sequence:           1,
		SourcePort:         sourcePort,
		SourceChannel:      sourceChannel,
		DestinationPort:    destPort,
		DestinationChannel: destChannel,
		Data:               data,
		TimeoutHeight:      clienttypes.Height{},
	}
}

// Case 1: Native token returning home
// ubadge sent to Osmosis via transfer/channel-0, coming back
// Packet denom: "transfer/channel-0/ubadge" (Osmosis prefixed it with its source port/channel)
// Expected local denom: "ubadge"
func TestExtractDenom_NativeReturning(t *testing.T) {
	packet := makePacket(
		"transfer/channel-0/ubadge", // packet denom (prefixed by sender chain)
		"transfer", "channel-0",     // source (sender's port/channel)
		"transfer", "channel-1",     // dest (our port/channel)
	)

	denom, err := extractDenomFromPacketOnRecv(packet)
	require.NoError(t, err)
	require.Equal(t, "ubadge", denom, "native token returning home should resolve to base denom")
}

// Case 2: Multi-hop native token returning
// ubadge went BB -> Osmosis(channel-0) -> Hub -> back to BB
// Hub sends it with denom "transfer/channel-0/transfer/channel-99/ubadge"
// ReceiverChainIsSource = true (first prefix matches source port/channel)
// After stripping first hop: remaining trace has transfer/channel-99 prefix = IBC denom
func TestExtractDenom_MultiHopNativeReturning(t *testing.T) {
	packet := makePacket(
		"transfer/channel-0/transfer/channel-99/ubadge",
		"transfer", "channel-0", // source (Hub's side — matches first trace hop)
		"transfer", "channel-1", // dest (our side)
	)

	denom, err := extractDenomFromPacketOnRecv(packet)
	require.NoError(t, err)

	// After stripping the source hop, remaining trace is transfer/channel-99/ubadge
	// This should be an IBC denom (has a trace), not the bare "ubadge"
	require.Contains(t, denom, "ibc/", "multi-hop native should be an IBC hash denom")
	require.NotEqual(t, "ubadge", denom, "should not resolve to bare base denom")
	// Verify it matches what ibc-go would produce: NewDenom with the remaining hop
	expected := transfertypes.NewDenom("ubadge", transfertypes.NewHop("transfer", "channel-99")).IBCDenom()
	require.Equal(t, expected, denom)
}

// Case 3: Foreign token arriving for the first time
// Cosmos Hub sends "uatom" to BitBadges
// Packet denom: "uatom" (native on sender, no prefix)
// Expected local denom: ibc/HASH(transfer/channel-M/uatom)
func TestExtractDenom_ForeignFirstTime(t *testing.T) {
	destPort := "transfer"
	destChannel := "channel-5"

	packet := makePacket(
		"uatom",                     // raw denom — no prefix
		"transfer", "channel-0",     // source (Hub's side)
		destPort, destChannel,       // dest (our side)
	)

	denom, err := extractDenomFromPacketOnRecv(packet)
	require.NoError(t, err)

	// Should be ibc/HASH(transfer/channel-5/uatom) — our dest port/channel prepended
	expected := transfertypes.NewDenom("uatom", transfertypes.NewHop(destPort, destChannel)).IBCDenom()
	require.Equal(t, expected, denom, "foreign token should resolve to ibc/HASH(destPort/destChannel/denom)")
	require.Contains(t, denom, "ibc/", "should be an IBC hash denom")
	require.NotEqual(t, "uatom", denom, "must NOT return raw packet denom")
}

// Case 4: Multi-hop foreign token
// uatom went Hub -> Osmosis(channel-X) -> BitBadges
// Packet denom: "transfer/channel-X/uatom" (Osmosis's prefix for its IBC atom)
// Expected: ibc/HASH(transfer/channel-M/transfer/channel-X/uatom)
func TestExtractDenom_MultiHopForeign(t *testing.T) {
	destPort := "transfer"
	destChannel := "channel-7"

	packet := makePacket(
		"transfer/channel-X/uatom",  // Osmosis's IBC atom
		"transfer", "channel-3",     // source (Osmosis's side)
		destPort, destChannel,       // dest (our side)
	)

	denom, err := extractDenomFromPacketOnRecv(packet)
	require.NoError(t, err)

	// Should be ibc/HASH(transfer/channel-7/transfer/channel-X/uatom)
	expected := transfertypes.NewDenom("uatom",
		transfertypes.NewHop(destPort, destChannel),
		transfertypes.NewHop("transfer", "channel-X"),
	).IBCDenom()
	require.Equal(t, expected, denom, "multi-hop foreign should prepend dest port/channel to full trace")
	require.Contains(t, denom, "ibc/", "should be an IBC hash denom")
}

// Case 5: Invalid packet data
func TestExtractDenom_InvalidPacketData(t *testing.T) {
	packet := channeltypes.Packet{
		Data:               []byte("not valid json"),
		SourcePort:         "transfer",
		SourceChannel:      "channel-0",
		DestinationPort:    "transfer",
		DestinationChannel: "channel-1",
	}

	_, err := extractDenomFromPacketOnRecv(packet)
	require.Error(t, err, "should error on invalid packet data")
}

// Case 6: Verify consistency between source and sink branches
// Same base denom through two different paths should produce different local denoms
func TestExtractDenom_DifferentPathsDifferentDenoms(t *testing.T) {
	// uatom arriving via channel-1
	packet1 := makePacket("uatom", "transfer", "channel-0", "transfer", "channel-1")
	denom1, err := extractDenomFromPacketOnRecv(packet1)
	require.NoError(t, err)

	// uatom arriving via channel-2
	packet2 := makePacket("uatom", "transfer", "channel-0", "transfer", "channel-2")
	denom2, err := extractDenomFromPacketOnRecv(packet2)
	require.NoError(t, err)

	require.NotEqual(t, denom1, denom2, "same base denom via different channels should produce different local denoms")
}
