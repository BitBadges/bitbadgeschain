package hooks

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"

	ibchooks "github.com/bitbadges/bitbadgeschain/x/ibc-hooks"
	"github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/keeper"
	ratelimittypes "github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/types"
)

var _ ibchooks.OnRecvPacketBeforeHooks = &RateLimitHooks{}
var _ ibchooks.SendPacketBeforeHooks = &RateLimitHooks{}

type RateLimitHooks struct {
	keeper keeper.Keeper
}

func NewRateLimitHooks(keeper keeper.Keeper) *RateLimitHooks {
	return &RateLimitHooks{
		keeper: keeper,
	}
}

// OnRecvPacketBeforeHook implements OnRecvPacketBeforeHooks
// This is called before processing an incoming IBC packet
func (h *RateLimitHooks) OnRecvPacketBeforeHook(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) {
	// Parse ICS20 packet
	var data transfertypes.FungibleTokenPacketData
	if err := json.Unmarshal(packet.GetData(), &data); err != nil {
		// Not an ICS20 packet, skip
		return
	}

	// Extract amount
	amount, ok := sdkmath.NewIntFromString(data.Amount)
	if !ok {
		// Invalid amount, will fail later
		return
	}

	// Get channel ID
	channelID := packet.GetDestChannel()

	// Extract denom (convert to local denom)
	denom := extractDenomFromPacketOnRecv(packet, data.Denom)

	// Extract sender address
	senderAddr := data.Sender

	// Check rate limit for inflow
	if err := h.keeper.CheckRateLimit(ctx, channelID, denom, amount, true, senderAddr); err != nil {
		// Rate limit exceeded - this will be caught by the override hook
		// We log here but the actual rejection happens in the override hook
		h.keeper.Logger(ctx).Error("rate limit check failed on receive", "error", err, "channel", channelID, "denom", denom, "amount", amount)
	}
}

// SendPacketBeforeHook implements SendPacketBeforeHooks
// This is called before sending an IBC packet
func (h *RateLimitHooks) SendPacketBeforeHook(ctx sdk.Context, chanCap *capabilitytypes.Capability, sourcePort string, sourceChannel string, timeoutHeight clienttypes.Height, timeoutTimestamp uint64, data []byte) {
	// Parse ICS20 packet
	var packetData transfertypes.FungibleTokenPacketData
	if err := json.Unmarshal(data, &packetData); err != nil {
		// Not an ICS20 packet, skip
		return
	}

	// Extract amount
	amount, ok := sdkmath.NewIntFromString(packetData.Amount)
	if !ok {
		// Invalid amount, will fail later
		return
	}

	// Extract denom (for sends, denom is already in local format or needs parsing)
	denom := packetData.Denom
	// If it's an IBC denom, we need to parse it
	if transfertypes.ReceiverChainIsSource(sourcePort, sourceChannel, denom) {
		// This is a native token being sent out
		// Remove the IBC prefix if present
		voucherPrefix := transfertypes.GetDenomPrefix(sourcePort, sourceChannel)
		if len(denom) > len(voucherPrefix) && denom[:len(voucherPrefix)] == voucherPrefix {
			denom = denom[len(voucherPrefix):]
		}
		denomTrace := transfertypes.ParseDenomTrace(denom)
		if denomTrace.Path != "" {
			denom = denomTrace.IBCDenom()
		}
	}

	// Extract sender address
	senderAddr := packetData.Sender

	// Check rate limit for outflow
	if err := h.keeper.CheckRateLimit(ctx, sourceChannel, denom, amount, false, senderAddr); err != nil {
		// Rate limit exceeded - log but don't fail here as we need to fail in override hook
		h.keeper.Logger(ctx).Error("rate limit check failed on send", "error", err, "channel", sourceChannel, "denom", denom, "amount", amount)
	}
}

// OnRecvPacketOverrideHook implements OnRecvPacketOverrideHooks to actually reject packets
type RateLimitOverrideHooks struct {
	keeper keeper.Keeper
}

func NewRateLimitOverrideHooks(keeper keeper.Keeper) *RateLimitOverrideHooks {
	return &RateLimitOverrideHooks{
		keeper: keeper,
	}
}

// OnRecvPacketOverride implements OnRecvPacketOverrideHooks
func (h *RateLimitOverrideHooks) OnRecvPacketOverride(im ibchooks.IBCMiddleware, ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) ibcexported.Acknowledgement {
	// Parse ICS20 packet
	var data transfertypes.FungibleTokenPacketData
	if err := json.Unmarshal(packet.GetData(), &data); err != nil {
		// Not an ICS20 packet, pass through
		return im.App.OnRecvPacket(ctx, packet, relayer)
	}

	// Extract amount
	amount, ok := sdkmath.NewIntFromString(data.Amount)
	if !ok {
		// Invalid amount, pass through (will fail later)
		return im.App.OnRecvPacket(ctx, packet, relayer)
	}

	// Get channel ID
	channelID := packet.GetDestChannel()

	// Extract denom
	denom := extractDenomFromPacketOnRecv(packet, data.Denom)

	// Extract sender address (from packet data)
	senderAddr := data.Sender

	// Check rate limit for inflow
	if err := h.keeper.CheckRateLimit(ctx, channelID, denom, amount, true, senderAddr); err != nil {
		// Rate limit exceeded - reject packet
		return channeltypes.NewErrorAcknowledgement(errorsmod.Wrap(ratelimittypes.ErrRateLimitExceeded, err.Error()))
	}

	// Rate limit check passed, process packet
	ack := im.App.OnRecvPacket(ctx, packet, relayer)

	// If packet was successful, update all tracking
	if ack.Success() {
		h.updateTrackingAfterTransfer(ctx, channelID, denom, amount, true, senderAddr)
	}

	return ack
}

// SendPacketOverride implements SendPacketOverrideHooks
func (h *RateLimitOverrideHooks) SendPacketOverride(i ibchooks.ICS4Middleware, ctx sdk.Context, chanCap *capabilitytypes.Capability, sourcePort string, sourceChannel string, timeoutHeight clienttypes.Height, timeoutTimestamp uint64, data []byte) (uint64, error) {
	// Parse ICS20 packet
	var packetData transfertypes.FungibleTokenPacketData
	if err := json.Unmarshal(data, &packetData); err != nil {
		// Not an ICS20 packet, pass through
		return i.SendPacket(ctx, chanCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
	}

	// Extract amount
	amount, ok := sdkmath.NewIntFromString(packetData.Amount)
	if !ok {
		// Invalid amount, pass through (will fail later)
		return i.SendPacket(ctx, chanCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
	}

	// Extract denom
	denom := packetData.Denom
	if transfertypes.ReceiverChainIsSource(sourcePort, sourceChannel, denom) {
		voucherPrefix := transfertypes.GetDenomPrefix(sourcePort, sourceChannel)
		if len(denom) > len(voucherPrefix) && denom[:len(voucherPrefix)] == voucherPrefix {
			denom = denom[len(voucherPrefix):]
		}
		denomTrace := transfertypes.ParseDenomTrace(denom)
		if denomTrace.Path != "" {
			denom = denomTrace.IBCDenom()
		}
	}

	// Extract sender address (from packet data)
	senderAddr := packetData.Sender

	// Check rate limit for outflow
	if err := h.keeper.CheckRateLimit(ctx, sourceChannel, denom, amount, false, senderAddr); err != nil {
		// Rate limit exceeded - reject send
		return 0, errorsmod.Wrap(ratelimittypes.ErrRateLimitExceeded, err.Error())
	}

	// Rate limit check passed, send packet
	seq, err := i.SendPacket(ctx, chanCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
	if err == nil {
		// Update all tracking (negative for outflow)
		h.updateTrackingAfterTransfer(ctx, sourceChannel, denom, amount.Neg(), false, senderAddr)
	}

	return seq, err
}

// updateTrackingAfterTransfer updates all tracking after a successful transfer
func (h *RateLimitOverrideHooks) updateTrackingAfterTransfer(ctx sdk.Context, channelID, denom string, amount sdkmath.Int, isInflow bool, senderAddr string) {
	params := h.keeper.GetParams(ctx)
	config := params.FindMatchingConfig(channelID, denom)

	if config == nil {
		return // No config, no tracking
	}

	// Update multiple timeframe supply shift limits
	for _, limit := range config.SupplyShiftLimits {
		if limit.MaxAmount.IsZero() {
			continue
		}
		h.keeper.ResetChannelFlowWindowWithTimeframe(ctx, channelID, denom, limit.TimeframeType, limit.TimeframeDuration)
		flow, _ := h.keeper.GetChannelFlowWithTimeframe(ctx, channelID, denom, limit.TimeframeType, limit.TimeframeDuration)
		flow.NetFlow = flow.NetFlow.Add(amount)
		h.keeper.SetChannelFlowWithTimeframe(ctx, channelID, denom, limit.TimeframeType, limit.TimeframeDuration, flow)
	}

	// Update unique sender tracking (only for inflows)
	if isInflow && senderAddr != "" {
		for _, limit := range config.UniqueSenderLimits {
			if limit.MaxUniqueSenders == 0 {
				continue
			}
			h.keeper.ResetUniqueSendersWindow(ctx, channelID, limit.TimeframeType, limit.TimeframeDuration)
			h.keeper.AddUniqueSender(ctx, channelID, senderAddr, limit.TimeframeType, limit.TimeframeDuration)
		}
	}

	// Update per-address tracking
	if senderAddr != "" {
		for _, limit := range config.AddressLimits {
			h.keeper.ResetAddressTransferWindow(ctx, senderAddr, channelID, denom, limit.TimeframeType, limit.TimeframeDuration)
			data, _ := h.keeper.GetAddressTransferData(ctx, senderAddr, channelID, denom, limit.TimeframeType, limit.TimeframeDuration)
			data.TransferCount++
			data.TotalAmount = data.TotalAmount.Add(amount.Abs()) // Use absolute value for tracking
			h.keeper.SetAddressTransferData(ctx, senderAddr, channelID, denom, limit.TimeframeType, limit.TimeframeDuration, data)
		}
	}
}

// extractDenomFromPacketOnRecv extracts the local denom from a received packet
func extractDenomFromPacketOnRecv(packet channeltypes.Packet, packetDenom string) string {
	denom := packetDenom
	if transfertypes.ReceiverChainIsSource(packet.GetSourcePort(), packet.GetSourceChannel(), denom) {
		// Remove prefix added by sender chain
		voucherPrefix := transfertypes.GetDenomPrefix(packet.GetSourcePort(), packet.GetSourceChannel())
		unprefixedDenom := denom[len(voucherPrefix):]

		// coin denomination used in sending from the source chain
		denom = unprefixedDenom

		// The denomination used to send the coins is either the native denom or the hash of the path
		// if the denomination is not native.
		denomTrace := transfertypes.ParseDenomTrace(unprefixedDenom)
		if denomTrace.Path != "" {
			denom = denomTrace.IBCDenom()
		} else {
			denom = unprefixedDenom
		}
	} else {
		// Chain is sink for the denom, build local denom
		denomTrace := transfertypes.ParseDenomTrace(denom)
		denom = denomTrace.IBCDenom()
	}
	return denom
}
