package app

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// Import legacy types packages - their init() functions register proto types
	badgestypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
	wasmtypes "github.com/bitbadges/bitbadgeschain/x/wasm/types"
)

// Legacy types for backward compatibility during genesis export.
// These types are registered so that old governance proposals referencing
// removed/renamed module messages can be properly marshaled/unmarshaled.
//
// These stubs are needed because:
// - WASM module has been removed from the chain
// - badges module has been renamed to tokenization
// Old governance proposals may still reference these old message types.

// RegisterLegacyWasmInterfaces registers the legacy WASM types for backward
// compatibility when exporting/importing genesis with old governance proposals.
func RegisterLegacyWasmInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&wasmtypes.MsgUpdateParams{},
	)
}

// RegisterLegacyBadgesInterfaces registers the legacy badges types for backward
// compatibility when exporting/importing genesis with old governance proposals.
func RegisterLegacyBadgesInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&badgestypes.MsgUpdateParams{},
	)
}
