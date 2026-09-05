package hooks

import (
	"errors"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v11/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v11/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v11/modules/core/exported"
	"github.com/stretchr/testify/require"

	ibchooks "github.com/bitbadges/bitbadgeschain/x/ibc-hooks"
	ratelimittypes "github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/types"
)

func TestReceiveDenomMatchesLocalTrace(t *testing.T) {
	packet := buildICS20Packet(t, "1")
	for _, tc := range []struct{ wire, local string }{
		{"uatom", "transfer/channel-99/uatom"},
		{"transfer/channel-8/uatom", "transfer/channel-99/transfer/channel-8/uatom"},
		{"transfer/" + packet.SourceChannel + "/uatom", "uatom"},
		{"transfer/" + packet.SourceChannel + "/transfer/channel-8/uatom", "transfer/channel-8/uatom"},
	} {
		require.Equal(t, transfertypes.ExtractDenomFromPath(tc.local).IBCDenom(), extractDenomFromPacketOnRecv(packet, tc.wire))
	}
}

type receiveModule struct {
	noopIBCModule
	ack     ibcexported.Acknowledgement
	version string
	calls   int
}

func (m *receiveModule) OnRecvPacket(_ sdk.Context, version string, _ channeltypes.Packet, _ sdk.AccAddress) ibcexported.Acknowledgement {
	m.version = version
	m.calls++
	return m.ack
}

func TestReceiveUsesDestinationChannelAndPreservesVersion(t *testing.T) {
	for _, amount := range []string{"500", "1000001"} {
		t.Run(amount, func(t *testing.T) {
			h, k, ctx := newTestHooks(t)
			packet := buildICS20Packet(t, amount)
			params := k.GetParams(ctx)
			params.RateLimits[0].ChannelId = packet.DestinationChannel
			params.RateLimits[0].Denom = extractDenomFromPacketOnRecv(packet, testDenom)
			k.SetParams(ctx, params)
			app := &receiveModule{ack: channeltypes.NewResultAcknowledgement([]byte{1})}
			ics4 := ibchooks.NewICS4Middleware(stubChannel{}, h)
			im := ibchooks.NewIBCMiddleware(app, &ics4)
			ack := im.OnRecvPacket(ctx, "ics20-1", packet, nil)
			if amount == "1000001" {
				require.False(t, ack.Success())
				require.Zero(t, app.calls)
				return
			}
			require.True(t, ack.Success())
			require.Equal(t, "ics20-1", app.version)
			flow, _ := k.GetChannelFlowWithTimeframe(ctx, packet.DestinationChannel, params.RateLimits[0].Denom, ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1000)
			require.True(t, flow.NetFlow.Equal(sdkmath.NewInt(500)))
		})
	}
}

func TestReceivePendingAcknowledgement(t *testing.T) {
	h, k, ctx := newTestHooks(t)
	packet := buildICS20Packet(t, "500")
	params := k.GetParams(ctx)
	params.RateLimits[0].ChannelId = packet.DestinationChannel
	params.RateLimits[0].Denom = extractDenomFromPacketOnRecv(packet, testDenom)
	k.SetParams(ctx, params)
	ics4 := ibchooks.NewICS4Middleware(stubChannel{}, h)
	im := ibchooks.NewIBCMiddleware(&receiveModule{}, &ics4)
	require.NotPanics(t, func() { require.Nil(t, im.OnRecvPacket(ctx, "ics20-1", packet, nil)) })
	flow, _ := k.GetChannelFlowWithTimeframe(ctx, packet.DestinationChannel, params.RateLimits[0].Denom, ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1000)
	require.True(t, flow.NetFlow.Equal(sdkmath.NewInt(500)))
	errAck := channeltypes.NewErrorAcknowledgementWithCodespace(errors.New("transfer rejected"))
	require.NoError(t, ics4.WriteAcknowledgement(ctx, packet, errAck))
	flow, _ = k.GetChannelFlowWithTimeframe(ctx, packet.DestinationChannel, params.RateLimits[0].Denom, ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1000)
	require.True(t, flow.NetFlow.IsZero())
	require.NoError(t, ics4.WriteAcknowledgement(ctx, packet, errAck))
	flow, _ = k.GetChannelFlowWithTimeframe(ctx, packet.DestinationChannel, params.RateLimits[0].Denom, ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1000)
	require.True(t, flow.NetFlow.IsZero())
}

type failingAckChannel struct{ stubChannel }

func (failingAckChannel) WriteAcknowledgement(sdk.Context, ibcexported.PacketI, ibcexported.Acknowledgement) error {
	return errors.New("acknowledgement write failed")
}

func TestReceiveAcknowledgementWindowLifecycle(t *testing.T) {
	for _, outcome := range []string{"success", "error", "expired", "write_error"} {
		t.Run(outcome, func(t *testing.T) {
			h, k, ctx := newTestHooks(t)
			packet := buildICS20Packet(t, "500")
			params := k.GetParams(ctx)
			denom := extractDenomFromPacketOnRecv(packet, testDenom)
			params.RateLimits[0].ChannelId = packet.DestinationChannel
			params.RateLimits[0].Denom = denom
			k.SetParams(ctx, params)
			ics4 := ibchooks.NewICS4Middleware(stubChannel{}, h)
			im := ibchooks.NewIBCMiddleware(&receiveModule{}, &ics4)
			require.Nil(t, im.OnRecvPacket(ctx, "ics20-1", packet, nil))
			if outcome == "expired" {
				ctx = ctx.WithBlockHeight(2001)
				other := buildICS20Packet(t, "200")
				other.Sequence = 2
				require.Nil(t, im.OnRecvPacket(ctx, "ics20-1", other, nil))
			}
			var ack ibcexported.Acknowledgement = channeltypes.NewErrorAcknowledgementWithCodespace(errors.New("transfer rejected"))
			if outcome == "success" {
				ack = channeltypes.NewResultAcknowledgement([]byte{1})
			}
			if outcome == "write_error" {
				ics4 = ibchooks.NewICS4Middleware(failingAckChannel{}, h)
				require.Error(t, ics4.WriteAcknowledgement(ctx, packet, ack))
			} else {
				require.NoError(t, ics4.WriteAcknowledgement(ctx, packet, ack))
			}
			flow, _ := k.GetChannelFlowWithTimeframe(ctx, packet.DestinationChannel, denom, ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1000)
			address, _ := k.GetAddressTransferData(ctx, testSender, packet.DestinationChannel, denom, ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1000)
			expected := int64(500)
			if outcome == "error" {
				expected = 0
			}
			if outcome == "expired" {
				expected = 200
			}
			require.Equal(t, sdkmath.NewInt(expected), flow.NetFlow)
			require.Equal(t, sdkmath.NewInt(expected), address.TotalAmount)
			_, pending := k.GetPendingReceiveWindow(ctx, packet.DestinationPort, packet.DestinationChannel, packet.Sequence, ratelimittypes.PendingSendScopeSupplyShift, ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1000)
			require.Equal(t, outcome == "write_error", pending)
		})
	}
}
