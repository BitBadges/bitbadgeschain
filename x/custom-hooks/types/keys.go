package types

import (
	storetypes "cosmossdk.io/store/types"
)

const (
	ModuleName = "customhooks"
	RouterKey  = ModuleName
	StoreKey   = ModuleName
)

var (
	// TransientStoreKey is the transient store key for storing deterministic error messages
	// Transient stores are automatically cleared at the end of each transaction
	TransientStoreKey = storetypes.NewTransientStoreKey("customhooks_transient")

	// DeterministicErrorKey is the key used to store deterministic error messages in the transient store
	DeterministicErrorKey = []byte("deterministic_error")
)
