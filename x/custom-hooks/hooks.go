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
// Both IBC transfer and hook execution are wrapped in a cached context to ensure atomicity
func (h *CustomHooks) OnRecvPacketOverride(im ibchooks.IBCMiddleware, ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) ibcexported.Acknowledgement {
	h.keeper.Logger(ctx).Info("custom-hooks: OnRecvPacketOverride", "packet", packet)

	// Check if this is an ICS20 packet and parse hook data before executing IBC transfer
	// This allows us to determine if we need atomicity before committing IBC transfer state
	isIcs20, data := isIcs20Packet(packet.GetData())
	var hookData *customhookstypes.HookData
	var parseErr error
	if isIcs20 {
		hookData, parseErr = customhookstypes.ParseHookDataFromMemo(data.GetMemo())
		if parseErr != nil {
			h.keeper.Logger(ctx).Error("custom-hooks: failed to parse memo", "error", parseErr)
			// Return custom error acknowledgement with deterministic error string
			return customhookstypes.NewCustomErrorAcknowledgement("failed to parse hook data from memo")
		}
	}

	// If no hook data, execute IBC transfer normally (no atomicity needed)
	if !isIcs20 || hookData == nil {
		return im.App.OnRecvPacket(ctx, packet, relayer)
	}

	// We have hook data - need atomicity between IBC transfer and hook execution
	// Wrap both operations in a cached context to ensure atomicity
	cacheCtx, writeCache := ctx.CacheContext()

	// Execute IBC transfer in cached context
	ack := im.App.OnRecvPacket(cacheCtx, packet, relayer)

	// If the IBC transfer itself failed, discard cache and return error
	if !ack.Success() {
		// Cache context is automatically discarded - no state changes committed
		return ack
	}

	// IBC transfer succeeded in cached context - now execute hook in same cached context
	// Extract amount and denom from packet
	amount, ok := osmomath.NewIntFromString(data.Amount)
	if !ok {
		h.keeper.Logger(ctx).Error("custom-hooks: failed to parse amount from packet")
		// Cache context is automatically discarded - no state changes committed
		return customhookstypes.NewCustomErrorAcknowledgement(fmt.Sprintf("failed to parse amount from packet: %s", data.Amount))
	}

	// The packet's denom is the denom in the sender chain. This needs to be converted to the local denom.
	denom, err := extractDenomFromPacketOnRecv(packet)
	if err != nil {
		h.keeper.Logger(ctx).Error("custom-hooks: failed to extract denom from packet", "error", err)
		// Cache context is automatically discarded - no state changes committed
		return customhookstypes.NewCustomErrorAcknowledgement(fmt.Sprintf("failed to extract denom from packet: %s", err.Error()))
	}

	// Get sender address from packet
	sender := data.GetSender()
	channel := packet.GetDestChannel()

	// log sender and channel
	h.keeper.Logger(ctx).Info("custom-hooks: sender", "sender", sender)
	h.keeper.Logger(ctx).Info("custom-hooks: channel", "channel", channel)

	// Store original sender address as string (from source chain, can't convert to AccAddress)
	originalSender := sender

	// Derive intermediate sender address
	// TODO: I think we can explore removing this
	// This was from the Osmosis IBC hooks implementation
	// Their reasons was that they were doing any WASM contract call, but we don't use WASM contracts here.
	senderBech32, err := ibchookstypes.DeriveIntermediateSender(channel, sender, h.bech32PrefixAccAddr)
	if err != nil {
		h.keeper.Logger(ctx).Error("custom-hooks: failed to derive intermediate sender", "error", err)
		// Cache context is automatically discarded - no state changes committed
		return customhookstypes.NewCustomErrorAcknowledgement("failed to derive intermediate sender")
	}

	// Convert sender string to AccAddress
	senderAddr, err := sdk.AccAddressFromBech32(senderBech32)
	if err != nil {
		h.keeper.Logger(ctx).Error("custom-hooks: failed to convert sender to AccAddress", "error", err)
		// Cache context is automatically discarded - no state changes committed
		return customhookstypes.NewCustomErrorAcknowledgement(fmt.Sprintf("invalid sender address: %s", senderBech32))
	}

	// Convert osmomath.Int to sdkmath.Int for sdk.Coin
	amountSDK := sdkmath.NewIntFromBigInt(amount.BigInt())
	tokenCoin := sdk.NewCoin(denom, amountSDK)

	// Execute hook in the same cached context as IBC transfer
	hookAck := h.keeper.ExecuteHook(cacheCtx, senderAddr, hookData, tokenCoin, originalSender)
	if !hookAck.Success() {
		h.keeper.Logger(ctx).Error("custom-hooks: failed to execute hook")
		// Cache context is automatically discarded - no state changes committed
		// Both IBC transfer and hook state changes are rolled back atomically
		return hookAck
	}

	// Both IBC transfer and hook succeeded in cached context
	// Commit all state changes atomically
	writeCache()

	// Return the IBC transfer ack (which we know is successful)
	// Both IBC transfer and hook state are now committed atomically
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

// extractDenomFromPacketOnRecv extracts the denom from the packet on receive
// This is similar to the wasm hooks implementation
// Returns the denom and an error if the packet data cannot be unmarshalled
func extractDenomFromPacketOnRecv(packet channeltypes.Packet) (string, error) {
	var data transfertypes.FungibleTokenPacketData
	if err := json.Unmarshal(packet.GetData(), &data); err != nil {
		return "", fmt.Errorf("cannot unmarshal ICS-20 transfer packet data: %w", err)
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
	return denom, nil
}
