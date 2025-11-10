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
	ibchookstypes "github.com/bitbadges/bitbadgeschain/x/ibc-hooks/types"
)

// CustomHooks implements OnRecvPacketAfterHooks to execute custom hooks after IBC transfer
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

// OnRecvPacketAfterHook implements OnRecvPacketAfterHooks interface
// This is called after the IBC transfer is successfully received
func (h *CustomHooks) OnRecvPacketAfterHook(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress, ack ibcexported.Acknowledgement) {
	// Only process if the packet was successful
	if !ack.Success() {
		return
	}

	// Check if this is an ICS20 packet
	isIcs20, data := isIcs20Packet(packet.GetData())
	if !isIcs20 {
		return
	}

	// Parse hook data from memo
	hookData, err := customhookstypes.ParseHookDataFromMemo(data.GetMemo())
	if err != nil || hookData == nil {
		// No custom hooks in memo, or error parsing (non-fatal)
		return
	}

	// Extract amount and denom from packet
	amount, ok := osmomath.NewIntFromString(data.Amount)
	if !ok {
		h.keeper.Logger(ctx).Error("failed to parse amount from packet", "amount", data.Amount)
		return
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
		return
	}

	// Convert sender string to AccAddress
	senderAddr, err := sdk.AccAddressFromBech32(senderBech32)
	if err != nil {
		h.keeper.Logger(ctx).Error("failed to convert sender to AccAddress", "error", err)
		return
	}

	// Execute custom hooks after transfer
	// Convert osmomath.Int to sdkmath.Int for sdk.Coin
	amountSDK := sdkmath.NewIntFromBigInt(amount.BigInt())
	if err := h.keeper.ExecuteHook(ctx, senderAddr, hookData, sdk.NewCoin(denom, amountSDK)); err != nil {
		// Log error but don't fail the packet - custom hooks are optional
		h.keeper.Logger(ctx).Error("failed to execute custom hook", "error", err)
	}
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
