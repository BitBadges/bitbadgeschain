package hooks

import (
	"context"
	"encoding/json"
	"testing"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/store/v2"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdkmath "cosmossdk.io/math"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"
	"github.com/stretchr/testify/require"

	ibchooks "github.com/bitbadges/bitbadgeschain/x/ibc-hooks"
	"github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/keeper"
	ratelimittypes "github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/types"
)

// --- test scaffolding ------------------------------------------------------

const (
	testChannelID = "channel-0"
	testSourcePort = "transfer"
	testDenom      = "uatom"
	testSender     = "cosmos1sender000000000000000000000000000000"
)

// mockBank satisfies ratelimittypes.BankKeeper with no-op behavior.
type mockBank struct{}

func (mockBank) GetSupply(ctx context.Context, denom string) sdk.Coin {
	return sdk.Coin{Denom: denom, Amount: sdkmath.ZeroInt()}
}
func (mockBank) GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins {
	return sdk.Coins{}
}
func (mockBank) MintCoins(ctx sdk.Context, moduleName string, coins sdk.Coins) error {
	return nil
}

func newTestKeeper(t *testing.T) (keeper.Keeper, sdk.Context) {
	t.Helper()

	storeKey := storetypes.NewKVStoreKey(ratelimittypes.StoreKey)
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)
	authority := "cosmos1w6t0l7z0yerj49ehnqwqaayxqpe3u7e23edgma"

	k := keeper.NewKeeper(cdc, storeKey, mockBank{}, authority)
	ctx := sdk.NewContext(stateStore, cmtproto.Header{Height: 1}, false, log.NewNopLogger())

	params := ratelimittypes.DefaultParams()
	params.RateLimits = []ratelimittypes.RateLimitConfig{{
		ChannelId: testChannelID,
		Denom:     testDenom,
		SupplyShiftLimits: []ratelimittypes.TimeframeLimit{{
			MaxAmount:         sdkmath.NewInt(1_000_000),
			TimeframeType:     ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK,
			TimeframeDuration: 1000,
		}},
		AddressLimits: []ratelimittypes.AddressLimit{{
			MaxTransfers:      10,
			MaxAmount:         sdkmath.NewInt(1_000_000),
			TimeframeType:     ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK,
			TimeframeDuration: 1000,
		}},
	}}
	k.SetParams(ctx, params)
	return k, ctx
}

func newTestHooks(t *testing.T) (*RateLimitOverrideHooks, keeper.Keeper, sdk.Context) {
	t.Helper()
	k, ctx := newTestKeeper(t)
	return NewRateLimitOverrideHooks(k), k, ctx
}

func buildICS20Packet(t *testing.T, amount string) channeltypes.Packet {
	t.Helper()
	data := transfertypes.FungibleTokenPacketData{
		Denom:    testDenom,
		Amount:   amount,
		Sender:   testSender,
		Receiver: "cosmos1receiver00000000000000000000000000",
	}
	bz, err := json.Marshal(data)
	require.NoError(t, err)
	return channeltypes.Packet{
		Sequence:           1,
		SourcePort:         testSourcePort,
		SourceChannel:      testChannelID,
		DestinationPort:    "transfer",
		DestinationChannel: "channel-99",
		Data:               bz,
		TimeoutHeight:      clienttypes.Height{RevisionNumber: 0, RevisionHeight: 1000},
		TimeoutTimestamp:   0,
	}
}

// seedOutboundTracking simulates the state that updateTrackingAfterTransfer
// would have written after a successful SendPacket: supply-shift NetFlow
// decreased by `amount` and per-address TransferCount/TotalAmount incremented.
func seedOutboundTracking(t *testing.T, k keeper.Keeper, ctx sdk.Context, amount sdkmath.Int) {
	t.Helper()
	k.SetChannelFlowWithTimeframe(ctx, testChannelID, testDenom,
		ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1000,
		ratelimittypes.ChannelFlow{NetFlow: amount.Neg()})
	k.SetAddressTransferData(ctx, testSender, testChannelID, testDenom,
		ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1000,
		ratelimittypes.AddressTransferData{
			TransferCount: 1,
			TotalAmount:   amount,
		})
}

// --- noop IBC module -------------------------------------------------------

// noopIBCModule implements porttypes.IBCModule with no-op behavior.
// Used to satisfy the IBCMiddleware.App interface in override-hook tests.
type noopIBCModule struct{}

func (noopIBCModule) OnChanOpenInit(sdk.Context, channeltypes.Order, []string, string, string, channeltypes.Counterparty, string) (string, error) {
	return "", nil
}
func (noopIBCModule) OnChanOpenTry(sdk.Context, channeltypes.Order, []string, string, string, channeltypes.Counterparty, string) (string, error) {
	return "", nil
}
func (noopIBCModule) OnChanOpenAck(sdk.Context, string, string, string, string) error { return nil }
func (noopIBCModule) OnChanOpenConfirm(sdk.Context, string, string) error              { return nil }
func (noopIBCModule) OnChanCloseInit(sdk.Context, string, string) error                { return nil }
func (noopIBCModule) OnChanCloseConfirm(sdk.Context, string, string) error             { return nil }
func (noopIBCModule) OnRecvPacket(sdk.Context, string, channeltypes.Packet, sdk.AccAddress) ibcexported.Acknowledgement {
	return nil
}
func (noopIBCModule) OnAcknowledgementPacket(sdk.Context, string, channeltypes.Packet, []byte, sdk.AccAddress) error {
	return nil
}
func (noopIBCModule) OnTimeoutPacket(sdk.Context, string, channeltypes.Packet, sdk.AccAddress) error {
	return nil
}

func newTestMiddleware() ibchooks.IBCMiddleware {
	return ibchooks.IBCMiddleware{
		App:            noopIBCModule{},
		ICS4Middleware: &ibchooks.ICS4Middleware{},
	}
}

// --- tests -----------------------------------------------------------------

func TestRefundSendTracking_ReversesSupplyShiftAndAddress(t *testing.T) {
	h, k, ctx := newTestHooks(t)
	amount := sdkmath.NewInt(500)
	seedOutboundTracking(t, k, ctx, amount)

	h.refundSendTracking(ctx, buildICS20Packet(t, "500"))

	flow, _ := k.GetChannelFlowWithTimeframe(ctx, testChannelID, testDenom,
		ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1000)
	require.True(t, flow.NetFlow.IsZero(), "supply-shift NetFlow must be zero after refund, got %s", flow.NetFlow.String())

	data, _ := k.GetAddressTransferData(ctx, testSender, testChannelID, testDenom,
		ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1000)
	require.Equal(t, int64(0), data.TransferCount)
	require.True(t, data.TotalAmount.IsZero(), "address TotalAmount must be zero after refund, got %s", data.TotalAmount.String())
}

func TestRefundSendTracking_NoConfigIsNoop(t *testing.T) {
	h, k, ctx := newTestHooks(t)

	// Use a channel/denom with no configured limits — the hook must early-return
	// without panicking or touching unrelated state.
	data := transfertypes.FungibleTokenPacketData{
		Denom:    "uunknown",
		Amount:   "100",
		Sender:   testSender,
		Receiver: "cosmos1receiver00000000000000000000000000",
	}
	bz, err := json.Marshal(data)
	require.NoError(t, err)
	packet := channeltypes.Packet{
		Sequence:      1,
		SourcePort:    testSourcePort,
		SourceChannel: "channel-99",
		Data:          bz,
	}

	require.NotPanics(t, func() {
		h.refundSendTracking(ctx, packet)
	})

	// State for the real channel must be untouched.
	flow, found := k.GetChannelFlowWithTimeframe(ctx, testChannelID, testDenom,
		ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1000)
	require.False(t, found, "unrelated flow must not be created, got %v", flow)
}

func TestRefundSendTracking_TotalAmountDoesNotUnderflow(t *testing.T) {
	h, k, ctx := newTestHooks(t)

	// Seed with a smaller TotalAmount than the refund amount. The refund path
	// must clamp to zero instead of producing a negative TotalAmount (which
	// would never happen in practice, but we defend against it).
	k.SetAddressTransferData(ctx, testSender, testChannelID, testDenom,
		ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1000,
		ratelimittypes.AddressTransferData{
			TransferCount: 0,
			TotalAmount:   sdkmath.NewInt(100),
		})

	h.refundSendTracking(ctx, buildICS20Packet(t, "500"))

	data, _ := k.GetAddressTransferData(ctx, testSender, testChannelID, testDenom,
		ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1000)
	require.True(t, data.TotalAmount.IsZero(), "TotalAmount must clamp to zero, got %s", data.TotalAmount.String())
	require.Equal(t, int64(0), data.TransferCount, "TransferCount must not go below zero")
}

func TestOnAcknowledgementPacketOverride_SuccessAckDoesNotRefund(t *testing.T) {
	h, k, ctx := newTestHooks(t)
	amount := sdkmath.NewInt(500)
	seedOutboundTracking(t, k, ctx, amount)

	// Build a success ack and send it through the override.
	successAck := channeltypes.NewResultAcknowledgement([]byte{0x01})
	ackBz := transfertypes.ModuleCdc.MustMarshalJSON(&successAck)

	err := h.OnAcknowledgementPacketOverride(newTestMiddleware(), ctx, testChannelID, buildICS20Packet(t, "500"), ackBz, sdk.AccAddress{})
	require.NoError(t, err)

	// Tracking must still reflect the outbound transfer — nothing refunded.
	flow, _ := k.GetChannelFlowWithTimeframe(ctx, testChannelID, testDenom,
		ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1000)
	require.True(t, flow.NetFlow.Equal(amount.Neg()), "success ack must not refund supply-shift, got %s", flow.NetFlow.String())

	data, _ := k.GetAddressTransferData(ctx, testSender, testChannelID, testDenom,
		ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1000)
	require.Equal(t, int64(1), data.TransferCount)
	require.True(t, data.TotalAmount.Equal(amount))
}

func TestOnAcknowledgementPacketOverride_ErrorAckRefunds(t *testing.T) {
	h, k, ctx := newTestHooks(t)
	amount := sdkmath.NewInt(500)
	seedOutboundTracking(t, k, ctx, amount)

	errorAck := channeltypes.NewErrorAcknowledgement(transfertypes.ErrReceiveDisabled)
	ackBz := transfertypes.ModuleCdc.MustMarshalJSON(&errorAck)

	err := h.OnAcknowledgementPacketOverride(newTestMiddleware(), ctx, testChannelID, buildICS20Packet(t, "500"), ackBz, sdk.AccAddress{})
	require.NoError(t, err)

	flow, _ := k.GetChannelFlowWithTimeframe(ctx, testChannelID, testDenom,
		ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1000)
	require.True(t, flow.NetFlow.IsZero(), "error ack must refund supply-shift to zero, got %s", flow.NetFlow.String())

	data, _ := k.GetAddressTransferData(ctx, testSender, testChannelID, testDenom,
		ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1000)
	require.Equal(t, int64(0), data.TransferCount)
	require.True(t, data.TotalAmount.IsZero())
}

func TestOnTimeoutPacketOverride_Refunds(t *testing.T) {
	h, k, ctx := newTestHooks(t)
	amount := sdkmath.NewInt(500)
	seedOutboundTracking(t, k, ctx, amount)

	err := h.OnTimeoutPacketOverride(newTestMiddleware(), ctx, testChannelID, buildICS20Packet(t, "500"), sdk.AccAddress{})
	require.NoError(t, err)

	flow, _ := k.GetChannelFlowWithTimeframe(ctx, testChannelID, testDenom,
		ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1000)
	require.True(t, flow.NetFlow.IsZero(), "timeout must refund supply-shift to zero, got %s", flow.NetFlow.String())

	data, _ := k.GetAddressTransferData(ctx, testSender, testChannelID, testDenom,
		ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1000)
	require.Equal(t, int64(0), data.TransferCount)
	require.True(t, data.TotalAmount.IsZero())
}
