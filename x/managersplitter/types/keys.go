package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ManagerSplitterKey is the prefix for manager splitter storage
	ManagerSplitterKey = "ManagerSplitter-value-"
	// ManagerSplitterCountKey is the prefix for the next manager splitter ID
	ManagerSplitterCountKey = "ManagerSplitter-count-"
)

// GetManagerSplitterKey returns the key for a manager splitter by address
func GetManagerSplitterKey(addr string) []byte {
	return append([]byte(ManagerSplitterKey), []byte(addr)...)
}

// GetManagerSplitterCountKey returns the key for the next manager splitter ID
func GetManagerSplitterCountKey() []byte {
	return []byte(ManagerSplitterCountKey)
}

// DeriveManagerSplitterAddress derives a module address for a manager splitter from its ID
func DeriveManagerSplitterAddress(id sdkmath.Uint) string {
	// Convert ID to bytes (big endian)
	idBytes := []byte(id.String())
	
	// Use the module address derivation function
	addr := address.Module(ModuleName, idBytes)
	return sdk.AccAddress(addr).String()
}

