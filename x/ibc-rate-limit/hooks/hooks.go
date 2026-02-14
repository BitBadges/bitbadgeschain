package hooks

import (
	"encoding/json"
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"

	ibchooks "github.com/bitbadges/bitbadgeschain/x/ibc-hooks"
	"github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/keeper"
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
func (h *RateLimitHooks) OnRecvPacketBeforeHook(ctx sdk.Context, channelID string, packet channeltypes.Packet, relayer sdk.AccAddress) {
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

	// Extract denom (convert to local denom)
	// channelID is already provided as a parameter
	denom := extractDenomFromPacketOnRecv(packet, data.Denom)

	// Extract sender address
	senderAddr := data.Sender

	// Check rate limit for inflow
	ack := h.keeper.CheckRateLimit(ctx, channelID, denom, amount, true, senderAddr)
	if !ack.Success() {
		// Rate limit exceeded - this will be caught by the override hook
		h.keeper.Logger(ctx).Error("rate limit check failed on receive", "channel", channelID, "denom", denom, "amount", amount)
	}
}

// SendPacketBeforeHook implements SendPacketBeforeHooks
// This is called before sending an IBC packet
func (h *RateLimitHooks) SendPacketBeforeHook(ctx sdk.Context, sourcePort string, sourceChannel string, timeoutHeight clienttypes.Height, timeoutTimestamp uint64, data []byte) {
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
		if denomTrace.Path() != "" {
			denom = denomTrace.IBCDenom()
		}
	}

	// Extract sender address
	senderAddr := packetData.Sender

	// Check rate limit for outflow
	ack := h.keeper.CheckRateLimit(ctx, sourceChannel, denom, amount, false, senderAddr)
	if !ack.Success() {
		// Rate limit exceeded - log but don't fail here as we need to fail in override hook
		h.keeper.Logger(ctx).Error("rate limit check failed on send", "channel", sourceChannel, "denom", denom, "amount", amount)
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
func (h *RateLimitOverrideHooks) OnRecvPacketOverride(im ibchooks.IBCMiddleware, ctx sdk.Context, channelID string, packet channeltypes.Packet, relayer sdk.AccAddress) ibcexported.Acknowledgement {
	// Parse ICS20 packet
	var data transfertypes.FungibleTokenPacketData
	if err := json.Unmarshal(packet.GetData(), &data); err != nil {
		// Not an ICS20 packet, pass through
		return im.App.OnRecvPacket(ctx, channelID, packet, relayer)
	}

	// Extract amount
	amount, ok := sdkmath.NewIntFromString(data.Amount)
	if !ok {
		// Invalid amount, pass through (will fail later)
		return im.App.OnRecvPacket(ctx, channelID, packet, relayer)
	}

	// Extract denom
	denom := extractDenomFromPacketOnRecv(packet, data.Denom)

	// Extract sender address (from packet data)
	senderAddr := data.Sender

	// Check rate limit for inflow
	rateLimitAck := h.keeper.CheckRateLimit(ctx, channelID, denom, amount, true, senderAddr)
	if !rateLimitAck.Success() {
		// Rate limit exceeded - reject packet
		h.keeper.Logger(ctx).Error("rate limit exceeded, rejecting packet",
			"channel", channelID,
			"denom", denom,
			"amount", amount.String(),
			"sender", senderAddr,
		)
		// Return the error acknowledgement from CheckRateLimit (already has deterministic error message)
		return rateLimitAck
	}

	// Rate limit check passed, process packet
	ack := im.App.OnRecvPacket(ctx, channelID, packet, relayer)

	// If packet was successful, update all tracking
	if ack.Success() {
		h.updateTrackingAfterTransfer(ctx, channelID, denom, amount, true, senderAddr)
	}

	return ack
}

// SendPacketOverride implements SendPacketOverrideHooks
func (h *RateLimitOverrideHooks) SendPacketOverride(i ibchooks.ICS4Middleware, ctx sdk.Context, sourcePort string, sourceChannel string, timeoutHeight clienttypes.Height, timeoutTimestamp uint64, data []byte) (uint64, error) {
	// Parse ICS20 packet
	var packetData transfertypes.FungibleTokenPacketData
	if err := json.Unmarshal(data, &packetData); err != nil {
		// Not an ICS20 packet, pass through
		return i.SendPacket(ctx, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
	}

	// Extract amount
	amount, ok := sdkmath.NewIntFromString(packetData.Amount)
	if !ok {
		// Invalid amount, pass through (will fail later)
		return i.SendPacket(ctx, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
	}

	// Extract denom
	denom := packetData.Denom
	if transfertypes.ReceiverChainIsSource(sourcePort, sourceChannel, denom) {
		voucherPrefix := transfertypes.GetDenomPrefix(sourcePort, sourceChannel)
		if len(denom) > len(voucherPrefix) && denom[:len(voucherPrefix)] == voucherPrefix {
			denom = denom[len(voucherPrefix):]
		}
		denomTrace := transfertypes.ParseDenomTrace(denom)
		if denomTrace.Path() != "" {
			denom = denomTrace.IBCDenom()
		}
	}

	// Extract sender address (from packet data)
	senderAddr := packetData.Sender

	// Check rate limit for outflow
	ack := h.keeper.CheckRateLimit(ctx, sourceChannel, denom, amount, false, senderAddr)
	if !ack.Success() {
		// Rate limit exceeded - reject send
		h.keeper.Logger(ctx).Error("rate limit exceeded, rejecting send",
			"channel", sourceChannel,
			"denom", denom,
			"amount", amount.String(),
			"sender", senderAddr,
		)
		// Extract error message from acknowledgement for return
		channelAck, ok := ack.(channeltypes.Acknowledgement)
		if ok {
			if errResp, ok := channelAck.Response.(*channeltypes.Acknowledgement_Error); ok {
				return 0, fmt.Errorf("%s: %s", "rate limit exceeded", errResp.Error)
			}
		}
		// Fallback if acknowledgement format is unexpected
		return 0, fmt.Errorf("rate limit exceeded: channel=%s, denom=%s, amount=%s, sender=%s", sourceChannel, denom, amount.String(), senderAddr)
	}

	// Rate limit check passed, send packet
	seq, err := i.SendPacket(ctx, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
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
		if denomTrace.Path() != "" {
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
