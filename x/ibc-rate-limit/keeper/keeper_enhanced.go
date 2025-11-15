package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/types"
)

const (
	// DefaultBlockTimeSeconds is the default block time in seconds (6 seconds for Cosmos chains)
	DefaultBlockTimeSeconds int64 = 6
)

// GetChannelFlowWithTimeframe gets the current flow state for a channel, denom, and timeframe
func (k Keeper) GetChannelFlowWithTimeframe(ctx sdk.Context, channelID, denom string, timeframeType types.TimeframeType, timeframeDuration int64) (types.ChannelFlow, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.ChannelFlowKey(channelID, denom, int32(timeframeType), timeframeDuration)
	bz := store.Get(key)
	if bz == nil {
		return types.ChannelFlow{NetFlow: sdkmath.ZeroInt()}, false
	}

	var flow types.ChannelFlow
	k.cdc.MustUnmarshal(bz, &flow)
	return flow, true
}

// SetChannelFlowWithTimeframe sets the flow state for a channel, denom, and timeframe
func (k Keeper) SetChannelFlowWithTimeframe(ctx sdk.Context, channelID, denom string, timeframeType types.TimeframeType, timeframeDuration int64, flow types.ChannelFlow) {
	store := ctx.KVStore(k.storeKey)
	key := types.ChannelFlowKey(channelID, denom, int32(timeframeType), timeframeDuration)
	bz := k.cdc.MustMarshal(&flow)
	store.Set(key, bz)
}

// GetChannelFlowWindowWithTimeframe gets the time window for a channel, denom, and timeframe
func (k Keeper) GetChannelFlowWindowWithTimeframe(ctx sdk.Context, channelID, denom string, timeframeType types.TimeframeType, timeframeDuration int64) (types.ChannelFlowWindow, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.ChannelFlowWindowKey(channelID, denom, int32(timeframeType), timeframeDuration)
	bz := store.Get(key)
	if bz == nil {
		return types.ChannelFlowWindow{}, false
	}

	var window types.ChannelFlowWindow
	k.cdc.MustUnmarshal(bz, &window)
	return window, true
}

// SetChannelFlowWindowWithTimeframe sets the time window for a channel, denom, and timeframe
func (k Keeper) SetChannelFlowWindowWithTimeframe(ctx sdk.Context, channelID, denom string, timeframeType types.TimeframeType, timeframeDuration int64, window types.ChannelFlowWindow) {
	store := ctx.KVStore(k.storeKey)
	key := types.ChannelFlowWindowKey(channelID, denom, int32(timeframeType), timeframeDuration)
	bz := k.cdc.MustMarshal(&window)
	store.Set(key, bz)
}

// ResetChannelFlowWindowWithTimeframe resets the flow window for a channel, denom, and timeframe if it has expired
func (k Keeper) ResetChannelFlowWindowWithTimeframe(ctx sdk.Context, channelID, denom string, timeframeType types.TimeframeType, timeframeDuration int64) {
	blockTimeSeconds := DefaultBlockTimeSeconds
	windowDurationBlocks := types.TimeframeDurationInBlocks(timeframeType, timeframeDuration, blockTimeSeconds)

	currentHeight := ctx.BlockHeight()
	window, found := k.GetChannelFlowWindowWithTimeframe(ctx, channelID, denom, timeframeType, timeframeDuration)

	// If no window exists or window has expired, create/reset it
	if !found || currentHeight >= window.WindowStart+windowDurationBlocks {
		newWindow := types.ChannelFlowWindow{
			WindowStart:    currentHeight,
			WindowDuration: windowDurationBlocks,
		}
		k.SetChannelFlowWindowWithTimeframe(ctx, channelID, denom, timeframeType, timeframeDuration, newWindow)
		// Reset flow to zero
		k.SetChannelFlowWithTimeframe(ctx, channelID, denom, timeframeType, timeframeDuration, types.ChannelFlow{NetFlow: sdkmath.ZeroInt()})
	}
}

// GetUniqueSenders gets the unique senders for a channel and timeframe
func (k Keeper) GetUniqueSenders(ctx sdk.Context, channelID string, timeframeType types.TimeframeType, timeframeDuration int64) (types.UniqueSenders, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.UniqueSendersKey(channelID, int32(timeframeType), timeframeDuration)
	bz := store.Get(key)
	if bz == nil {
		return types.UniqueSenders{Senders: []string{}}, false
	}

	var senders types.UniqueSenders
	k.cdc.MustUnmarshal(bz, &senders)
	return senders, true
}

// SetUniqueSenders sets the unique senders for a channel and timeframe
func (k Keeper) SetUniqueSenders(ctx sdk.Context, channelID string, timeframeType types.TimeframeType, timeframeDuration int64, senders types.UniqueSenders) {
	store := ctx.KVStore(k.storeKey)
	key := types.UniqueSendersKey(channelID, int32(timeframeType), timeframeDuration)
	bz := k.cdc.MustMarshal(&senders)
	store.Set(key, bz)
}

// AddUniqueSender adds a sender to the unique senders list if not already present
func (k Keeper) AddUniqueSender(ctx sdk.Context, channelID, senderAddr string, timeframeType types.TimeframeType, timeframeDuration int64) {
	senders, _ := k.GetUniqueSenders(ctx, channelID, timeframeType, timeframeDuration)

	// Check if sender already exists
	for _, addr := range senders.Senders {
		if addr == senderAddr {
			return // Already exists
		}
	}

	// Add new sender
	senders.Senders = append(senders.Senders, senderAddr)
	k.SetUniqueSenders(ctx, channelID, timeframeType, timeframeDuration, senders)
}

// GetUniqueSendersWindow gets the time window for unique sender tracking
func (k Keeper) GetUniqueSendersWindow(ctx sdk.Context, channelID string, timeframeType types.TimeframeType, timeframeDuration int64) (types.ChannelFlowWindow, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.UniqueSendersWindowKey(channelID, int32(timeframeType), timeframeDuration)
	bz := store.Get(key)
	if bz == nil {
		return types.ChannelFlowWindow{}, false
	}

	var window types.ChannelFlowWindow
	k.cdc.MustUnmarshal(bz, &window)
	return window, true
}

// SetUniqueSendersWindow sets the time window for unique sender tracking
func (k Keeper) SetUniqueSendersWindow(ctx sdk.Context, channelID string, timeframeType types.TimeframeType, timeframeDuration int64, window types.ChannelFlowWindow) {
	store := ctx.KVStore(k.storeKey)
	key := types.UniqueSendersWindowKey(channelID, int32(timeframeType), timeframeDuration)
	bz := k.cdc.MustMarshal(&window)
	store.Set(key, bz)
}

// ResetUniqueSendersWindow resets the unique senders window if it has expired
func (k Keeper) ResetUniqueSendersWindow(ctx sdk.Context, channelID string, timeframeType types.TimeframeType, timeframeDuration int64) {
	blockTimeSeconds := DefaultBlockTimeSeconds
	windowDurationBlocks := types.TimeframeDurationInBlocks(timeframeType, timeframeDuration, blockTimeSeconds)

	currentHeight := ctx.BlockHeight()
	window, found := k.GetUniqueSendersWindow(ctx, channelID, timeframeType, timeframeDuration)

	// If no window exists or window has expired, create/reset it
	if !found || currentHeight >= window.WindowStart+windowDurationBlocks {
		newWindow := types.ChannelFlowWindow{
			WindowStart:    currentHeight,
			WindowDuration: windowDurationBlocks,
		}
		k.SetUniqueSendersWindow(ctx, channelID, timeframeType, timeframeDuration, newWindow)
		// Reset unique senders
		k.SetUniqueSenders(ctx, channelID, timeframeType, timeframeDuration, types.UniqueSenders{Senders: []string{}})
	}
}

// GetAddressTransferData gets the transfer data for an address, channel, denom, and timeframe
func (k Keeper) GetAddressTransferData(ctx sdk.Context, address, channelID, denom string, timeframeType types.TimeframeType, timeframeDuration int64) (types.AddressTransferData, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.AddressTransferDataKey(address, channelID, denom, int32(timeframeType), timeframeDuration)
	bz := store.Get(key)
	if bz == nil {
		return types.AddressTransferData{
			TransferCount: 0,
			TotalAmount:   sdkmath.ZeroInt(),
		}, false
	}

	var data types.AddressTransferData
	k.cdc.MustUnmarshal(bz, &data)
	return data, true
}

// SetAddressTransferData sets the transfer data for an address, channel, denom, and timeframe
func (k Keeper) SetAddressTransferData(ctx sdk.Context, address, channelID, denom string, timeframeType types.TimeframeType, timeframeDuration int64, data types.AddressTransferData) {
	store := ctx.KVStore(k.storeKey)
	key := types.AddressTransferDataKey(address, channelID, denom, int32(timeframeType), timeframeDuration)
	bz := k.cdc.MustMarshal(&data)
	store.Set(key, bz)
}

// GetAddressTransferWindow gets the time window for address transfer tracking
func (k Keeper) GetAddressTransferWindow(ctx sdk.Context, address, channelID, denom string, timeframeType types.TimeframeType, timeframeDuration int64) (types.ChannelFlowWindow, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.AddressTransferWindowKey(address, channelID, denom, int32(timeframeType), timeframeDuration)
	bz := store.Get(key)
	if bz == nil {
		return types.ChannelFlowWindow{}, false
	}

	var window types.ChannelFlowWindow
	k.cdc.MustUnmarshal(bz, &window)
	return window, true
}

// SetAddressTransferWindow sets the time window for address transfer tracking
func (k Keeper) SetAddressTransferWindow(ctx sdk.Context, address, channelID, denom string, timeframeType types.TimeframeType, timeframeDuration int64, window types.ChannelFlowWindow) {
	store := ctx.KVStore(k.storeKey)
	key := types.AddressTransferWindowKey(address, channelID, denom, int32(timeframeType), timeframeDuration)
	bz := k.cdc.MustMarshal(&window)
	store.Set(key, bz)
}

// ResetAddressTransferWindow resets the address transfer window if it has expired
func (k Keeper) ResetAddressTransferWindow(ctx sdk.Context, address, channelID, denom string, timeframeType types.TimeframeType, timeframeDuration int64) {
	blockTimeSeconds := DefaultBlockTimeSeconds
	windowDurationBlocks := types.TimeframeDurationInBlocks(timeframeType, timeframeDuration, blockTimeSeconds)

	currentHeight := ctx.BlockHeight()
	window, found := k.GetAddressTransferWindow(ctx, address, channelID, denom, timeframeType, timeframeDuration)

	// If no window exists or window has expired, create/reset it
	if !found || currentHeight >= window.WindowStart+windowDurationBlocks {
		newWindow := types.ChannelFlowWindow{
			WindowStart:    currentHeight,
			WindowDuration: windowDurationBlocks,
		}
		k.SetAddressTransferWindow(ctx, address, channelID, denom, timeframeType, timeframeDuration, newWindow)
		// Reset transfer data
		k.SetAddressTransferData(ctx, address, channelID, denom, timeframeType, timeframeDuration, types.AddressTransferData{
			TransferCount: 0,
			TotalAmount:   sdkmath.ZeroInt(),
		})
	}
}
