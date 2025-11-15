package types

import (
	"fmt"
)

const (
	ModuleName = "ibcratelimit"
	RouterKey  = ModuleName
	StoreKey   = ModuleName
)

var (
	// ParamsKey stores the module parameters
	ParamsKey = []byte("p_ibcratelimit")

	// KeyPrefixChannelFlow stores the net flow for each channel+denom+timeframe combination
	// Key: channelFlowKey(channelID, denom, timeframeType, timeframeDuration) -> ChannelFlow
	KeyPrefixChannelFlow = []byte{0x01}

	// KeyPrefixChannelFlowWindow stores the time window for each channel+denom+timeframe combination
	// Key: channelFlowWindowKey(channelID, denom, timeframeType, timeframeDuration) -> ChannelFlowWindow
	KeyPrefixChannelFlowWindow = []byte{0x02}

	// KeyPrefixUniqueSenders stores unique sender addresses for each channel+timeframe combination
	// Key: uniqueSendersKey(channelID, timeframeType, timeframeDuration) -> UniqueSenders
	KeyPrefixUniqueSenders = []byte{0x03}

	// KeyPrefixUniqueSendersWindow stores the time window for unique sender tracking
	// Key: uniqueSendersWindowKey(channelID, timeframeType, timeframeDuration) -> ChannelFlowWindow
	KeyPrefixUniqueSendersWindow = []byte{0x04}

	// KeyPrefixAddressTransferData stores transfer data for each address+channel+denom+timeframe combination
	// Key: addressTransferDataKey(address, channelID, denom, timeframeType, timeframeDuration) -> AddressTransferData
	KeyPrefixAddressTransferData = []byte{0x05}

	// KeyPrefixAddressTransferWindow stores the time window for address transfer tracking
	// Key: addressTransferWindowKey(address, channelID, denom, timeframeType, timeframeDuration) -> ChannelFlowWindow
	KeyPrefixAddressTransferWindow = []byte{0x06}
)

// ChannelFlowKeyLegacy returns the key for storing channel flow state (backward compatibility)
func ChannelFlowKeyLegacy(channelID, denom string) []byte {
	key := append(KeyPrefixChannelFlow, []byte(channelID)...)
	key = append(key, []byte("|")...)
	key = append(key, []byte(denom)...)
	return key
}

// ChannelFlowWindowKeyLegacy returns the key for storing channel flow window (backward compatibility)
func ChannelFlowWindowKeyLegacy(channelID, denom string) []byte {
	key := append(KeyPrefixChannelFlowWindow, []byte(channelID)...)
	key = append(key, []byte("|")...)
	key = append(key, []byte(denom)...)
	return key
}

// channelFlowKey returns the key for storing channel flow state for a specific channel, denom, and timeframe
func ChannelFlowKey(channelID, denom string, timeframeType int32, timeframeDuration int64) []byte {
	key := append(KeyPrefixChannelFlow, []byte(channelID)...)
	key = append(key, []byte("|")...)
	key = append(key, []byte(denom)...)
	key = append(key, []byte("|")...)
	key = append(key, []byte(fmt.Sprintf("%d|%d", timeframeType, timeframeDuration))...)
	return key
}

// channelFlowWindowKey returns the key for storing channel flow window for a specific channel, denom, and timeframe
func ChannelFlowWindowKey(channelID, denom string, timeframeType int32, timeframeDuration int64) []byte {
	key := append(KeyPrefixChannelFlowWindow, []byte(channelID)...)
	key = append(key, []byte("|")...)
	key = append(key, []byte(denom)...)
	key = append(key, []byte("|")...)
	key = append(key, []byte(fmt.Sprintf("%d|%d", timeframeType, timeframeDuration))...)
	return key
}

// uniqueSendersKey returns the key for storing unique senders for a specific channel and timeframe
func UniqueSendersKey(channelID string, timeframeType int32, timeframeDuration int64) []byte {
	key := append(KeyPrefixUniqueSenders, []byte(channelID)...)
	key = append(key, []byte("|")...)
	key = append(key, []byte(fmt.Sprintf("%d|%d", timeframeType, timeframeDuration))...)
	return key
}

// uniqueSendersWindowKey returns the key for storing unique senders window for a specific channel and timeframe
func UniqueSendersWindowKey(channelID string, timeframeType int32, timeframeDuration int64) []byte {
	key := append(KeyPrefixUniqueSendersWindow, []byte(channelID)...)
	key = append(key, []byte("|")...)
	key = append(key, []byte(fmt.Sprintf("%d|%d", timeframeType, timeframeDuration))...)
	return key
}

// addressTransferDataKey returns the key for storing address transfer data for a specific address, channel, denom, and timeframe
func AddressTransferDataKey(address, channelID, denom string, timeframeType int32, timeframeDuration int64) []byte {
	key := append(KeyPrefixAddressTransferData, []byte(address)...)
	key = append(key, []byte("|")...)
	key = append(key, []byte(channelID)...)
	key = append(key, []byte("|")...)
	key = append(key, []byte(denom)...)
	key = append(key, []byte("|")...)
	key = append(key, []byte(fmt.Sprintf("%d|%d", timeframeType, timeframeDuration))...)
	return key
}

// addressTransferWindowKey returns the key for storing address transfer window for a specific address, channel, denom, and timeframe
func AddressTransferWindowKey(address, channelID, denom string, timeframeType int32, timeframeDuration int64) []byte {
	key := append(KeyPrefixAddressTransferWindow, []byte(address)...)
	key = append(key, []byte("|")...)
	key = append(key, []byte(channelID)...)
	key = append(key, []byte("|")...)
	key = append(key, []byte(denom)...)
	key = append(key, []byte("|")...)
	key = append(key, []byte(fmt.Sprintf("%d|%d", timeframeType, timeframeDuration))...)
	return key
}
