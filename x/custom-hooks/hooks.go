package customhooks

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	"github.com/bitbadges/bitbadgeschain/x/custom-hooks/keeper"
	customhookstypes "github.com/bitbadges/bitbadgeschain/x/custom-hooks/types"
	ibchooks "github.com/bitbadges/bitbadgeschain/x/ibc-hooks"
	ibchookstypes "github.com/bitbadges/bitbadgeschain/x/ibc-hooks/types"
)

// CustomHooks implements OnRecvPacketOverrideHooks to execute custom hooks and control acknowledgement
type CustomHooks struct {
	keeper              keeper.Keeper
	bech32PrefixAccAddr string
}

// NewCustomHooks creates a new CustomHooks instance
func NewCustomHooks(keeper keeper.Keeper, bech32PrefixAccAddr string) *CustomHooks {
	return &CustomHooks{
		keeper:              keeper,
		bech32PrefixAccAddr: bech32PrefixAccAddr,
	}
}

// OnRecvPacketOverride implements OnRecvPacketOverrideHooks interface
// This allows us to control the acknowledgement and fail the packet if the hook fails
func (h *CustomHooks) OnRecvPacketOverride(im ibchooks.IBCMiddleware, ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) ibcexported.Acknowledgement {
	h.keeper.Logger(ctx).Info("custom-hooks: OnRecvPacketOverride", "packet", packet)

	// First, process the IBC transfer by calling the underlying app
	ack := im.App.OnRecvPacket(ctx, packet, relayer)

	// If the IBC transfer itself failed, return the error acknowledgement
	if !ack.Success() {
		return ack
	}

	// Check if this is an ICS20 packet
	isIcs20, data := isIcs20Packet(packet.GetData())
	if !isIcs20 {
		return ack
	}

	// Parse hook data from memo
	hookData, err := customhookstypes.ParseHookDataFromMemo(data.GetMemo())
	if err != nil {
		h.keeper.Logger(ctx).Error("custom-hooks: failed to parse memo", "error", err)
		// Return custom error acknowledgement with deterministic error string
		return customhookstypes.NewCustomErrorAcknowledgement(fmt.Sprintf("failed to parse hook data from memo"))
	}

	// If no hook data, just return success (normal IBC transfer)
	if hookData == nil {
		return ack
	}

	// Extract amount and denom from packet
	amount, ok := osmomath.NewIntFromString(data.Amount)
	if !ok {
		h.keeper.Logger(ctx).Error("custom-hooks: failed to parse amount from packet")
		return customhookstypes.NewCustomErrorAcknowledgement(fmt.Sprintf("failed to parse amount from packet: %s", data.Amount))
	}

	// The packet's denom is the denom in the sender chain. This needs to be converted to the local denom.
	denom := mustExtractDenomFromPacketOnRecv(packet)

	// Get sender address from packet
	sender := data.GetSender()
	channel := packet.GetDestChannel()

	//log sender and channel
	h.keeper.Logger(ctx).Info("custom-hooks: sender", "sender", sender)
	h.keeper.Logger(ctx).Info("custom-hooks: channel", "channel", channel)

	// Store original sender address as string (from source chain, can't convert to AccAddress)
	originalSender := sender

	// Derive intermediate sender address
	// TODO: I think we can explore removing this
	// This was from the Osmosis IBC hooks implementation
	// Their reasons was that they were doing any WASM contract call, so they obviously can't trust any other chain
	// arbitrarily specifying a sender address.
	//
	// However, we are not doing any WASM contract calls. We are only ever doing swaps and transfers
	// with the received funds.
	senderBech32, err := ibchookstypes.DeriveIntermediateSender(channel, sender, h.bech32PrefixAccAddr)
	if err != nil {
		h.keeper.Logger(ctx).Error("custom-hooks: failed to derive intermediate sender", "error", err)
		return customhookstypes.NewCustomErrorAcknowledgement("failed to derive intermediate sender")
	}

	// Convert sender string to AccAddress
	senderAddr, err := sdk.AccAddressFromBech32(senderBech32)
	if err != nil {
		h.keeper.Logger(ctx).Error("custom-hooks: failed to convert sender to AccAddress", "error", err)
		return customhookstypes.NewCustomErrorAcknowledgement(fmt.Sprintf("invalid sender address: %s", senderBech32))
	}

	// Convert osmomath.Int to sdkmath.Int for sdk.Coin
	amountSDK := sdkmath.NewIntFromBigInt(amount.BigInt())
	tokenCoin := sdk.NewCoin(denom, amountSDK)

	// Execute custom hooks from intermediate sender address in a cache context
	// This ensures state changes are only committed if the hook succeeds
	cacheCtx, writeCache := ctx.CacheContext()
	hookAck := h.keeper.ExecuteHook(cacheCtx, senderAddr, hookData, tokenCoin, originalSender)
	if !hookAck.Success() {
		h.keeper.Logger(ctx).Error("custom-hooks: failed to execute hook")
		// Cache context is automatically discarded on error - no state changes committed
		// Return the error acknowledgement - this will roll back the IBC transfer
		return hookAck
	}

	// Hook succeeded - write cache to commit state changes
	writeCache()

	return ack
}

// isIcs20Packet checks if the packet data is an ICS20 transfer packet
func isIcs20Packet(data []byte) (bool, transfertypes.FungibleTokenPacketData) {
	var packetData transfertypes.FungibleTokenPacketData
	if err := json.Unmarshal(data, &packetData); err != nil {
		return false, packetData
	}
	return true, packetData
}

// mustExtractDenomFromPacketOnRecv extracts the denom from the packet on receive
// This is similar to the wasm hooks implementation
func mustExtractDenomFromPacketOnRecv(packet channeltypes.Packet) string {
	var data transfertypes.FungibleTokenPacketData
	if err := json.Unmarshal(packet.GetData(), &data); err != nil {
		panic(fmt.Errorf("cannot unmarshal ICS-20 transfer packet data: %w", err))
	}

	// The denom in the packet is the denom in the sender chain.
	// We need to convert it to the local denom.
	denom := data.Denom
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
	}
	return denom
}
