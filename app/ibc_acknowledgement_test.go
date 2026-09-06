package app

import (
	sdkmath "cosmossdk.io/math"
	"errors"
	ibchooks "github.com/bitbadges/bitbadgeschain/x/ibc-hooks"
	types "github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v11/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v11/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v11/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v11/modules/core/exported"
	"github.com/stretchr/testify/require"
	"testing"
)

type acknowledgementSink struct{ porttypes.ICS4Wrapper }

func (acknowledgementSink) WriteAcknowledgement(sdk.Context, ibcexported.PacketI, ibcexported.Acknowledgement) error {
	return nil
}

func TestCombinedHooksFinalizePendingReceive(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false).WithBlockHeight(1)
	packet := channeltypes.Packet{Sequence: 7, SourcePort: "transfer", SourceChannel: "channel-0", DestinationPort: "transfer", DestinationChannel: "channel-1"}
	data := transfertypes.FungibleTokenPacketData{Denom: "transfer/channel-0/uatom", Amount: "500", Sender: "sender", Receiver: "receiver"}
	var err error
	packet.Data, err = transfertypes.ModuleCdc.MarshalJSON(&data)
	require.NoError(t, err)
	params := types.DefaultParams()
	params.RateLimits = []types.RateLimitConfig{{ChannelId: packet.DestinationChannel, Denom: "uatom", SupplyShiftLimits: []types.TimeframeLimit{{MaxAmount: sdkmath.NewInt(1000), TimeframeType: types.TimeframeType_TIMEFRAME_TYPE_BLOCK, TimeframeDuration: 1000}}}}
	k := app.IBCRateLimitKeeper
	k.SetParams(ctx, params)
	k.ResetChannelFlowWindowWithTimeframe(ctx, packet.DestinationChannel, "uatom", types.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1000)
	flow, _ := k.GetChannelFlowWithTimeframe(ctx, packet.DestinationChannel, "uatom", types.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1000)
	flow.NetFlow = sdkmath.NewInt(500)
	k.SetChannelFlowWithTimeframe(ctx, packet.DestinationChannel, "uatom", types.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1000, flow)
	window, _ := k.GetChannelFlowWindowWithTimeframe(ctx, packet.DestinationChannel, "uatom", types.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1000)
	k.SetPendingReceiveWindow(ctx, packet.DestinationPort, packet.DestinationChannel, packet.Sequence, types.PendingSendScopeSupplyShift, types.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1000, window.WindowStart)
	wrapper := ibchooks.NewICS4Middleware(acknowledgementSink{}, app.HooksICS4Wrapper.Hooks)
	require.NoError(t, wrapper.WriteAcknowledgement(ctx, packet, channeltypes.NewErrorAcknowledgementWithCodespace(errors.New("transfer rejected"))))
	flow, _ = k.GetChannelFlowWithTimeframe(ctx, packet.DestinationChannel, "uatom", types.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1000)
	require.True(t, flow.NetFlow.IsZero())
}
