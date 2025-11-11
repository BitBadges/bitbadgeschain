package customhooks

import (
	"encoding/json"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
	// Log packet receipt for debugging duplicate detection
	h.keeper.Logger(ctx).Info("custom-hooks: OnRecvPacketOverride called",
		"packet_sequence", packet.GetSequence(),
		"packet_source_port", packet.GetSourcePort(),
		"packet_source_channel", packet.GetSourceChannel(),
		"block_height", ctx.BlockHeight())

	// First, process the IBC transfer by calling the underlying app
	ack := im.App.OnRecvPacket(ctx, packet, relayer)

	// If the IBC transfer itself failed, return the error acknowledgement
	if !ack.Success() {
		h.keeper.Logger(ctx).Info("custom-hooks: IBC transfer failed, skipping hook execution")
		return ack
	}

	// Check if this is an ICS20 packet
	isIcs20, data := isIcs20Packet(packet.GetData())
	if !isIcs20 {
		h.keeper.Logger(ctx).Info("custom-hooks: not an ICS20 packet, returning success ack")
		return ack
	}

	// Parse hook data from memo
	hookData, err := customhookstypes.ParseHookDataFromMemo(data.GetMemo())
	if err != nil {
		h.keeper.Logger(ctx).Error("custom-hooks: failed to parse memo", "error", err, "memo", data.GetMemo())
		// Return error acknowledgement - hook execution is required if memo is malformed
		return channeltypes.NewErrorAcknowledgement(errorsmod.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("failed to parse hook data from memo: %v", err)))
	}

	// If no hook data, just return success (normal IBC transfer)
	if hookData == nil {
		h.keeper.Logger(ctx).Info("custom-hooks: no hook data in memo, returning success ack")
		return ack
	}

	// Extract amount and denom from packet
	amount, ok := osmomath.NewIntFromString(data.Amount)
	if !ok {
		h.keeper.Logger(ctx).Error("failed to parse amount from packet", "amount", data.Amount)
		return channeltypes.NewErrorAcknowledgement(errorsmod.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("failed to parse amount from packet: %s", data.Amount)))
	}

	// The packet's denom is the denom in the sender chain. This needs to be converted to the local denom.
	denom := mustExtractDenomFromPacketOnRecv(packet)

	// Get sender address from packet
	sender := data.GetSender()
	channel := packet.GetDestChannel()

	// Derive intermediate sender address
	senderBech32, err := ibchookstypes.DeriveIntermediateSender(channel, sender, h.bech32PrefixAccAddr)
	if err != nil {
		h.keeper.Logger(ctx).Error("failed to derive intermediate sender", "error", err)
		return channeltypes.NewErrorAcknowledgement(errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("failed to derive intermediate sender: %v", err)))
	}

	// Convert sender string to AccAddress
	senderAddr, err := sdk.AccAddressFromBech32(senderBech32)
	if err != nil {
		h.keeper.Logger(ctx).Error("failed to convert sender to AccAddress", "error", err)
		return channeltypes.NewErrorAcknowledgement(errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid sender address: %v", err)))
	}

	// Convert osmomath.Int to sdkmath.Int for sdk.Coin
	amountSDK := sdkmath.NewIntFromBigInt(amount.BigInt())
	tokenCoin := sdk.NewCoin(denom, amountSDK)

	// Execute custom hooks from intermediate sender address
	// If this fails, the entire packet (including IBC transfer) will be rolled back
	if err := h.keeper.ExecuteHook(ctx, senderAddr, hookData, tokenCoin); err != nil {
		h.keeper.Logger(ctx).Error("custom-hooks: failed to execute hook, failing packet", "error", err)
		// Return error acknowledgement - this will roll back the IBC transfer
		errorAck := channeltypes.NewErrorAcknowledgement(errorsmod.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("hook execution failed: %v", err)))
		// Log the error acknowledgement details for debugging
		h.keeper.Logger(ctx).Error("custom-hooks: returning error acknowledgement",
			"ack_success", errorAck.Success(),
			"ack_type", fmt.Sprintf("%T", errorAck.Response))
		return errorAck
	}

	// Hook succeeded, return success acknowledgement
	h.keeper.Logger(ctx).Info("custom-hooks: hook executed successfully")
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
