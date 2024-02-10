package app

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmKeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	badgeKeeper "github.com/bitbadges/bitbadgeschain/x/badges/keeper"

	protocolKeeper "github.com/bitbadges/bitbadgeschain/x/protocols/keeper"

	customBindings "github.com/bitbadges/bitbadgeschain/custom-bindings"
)

func GetCustomMsgEncodersOptions() []wasmKeeper.Option {
	badgeEncodingOptions := wasmKeeper.WithMessageEncoders(badgeMsgEncoders())
	return []wasm.Option{badgeEncodingOptions}
}

func GetCustomMsgQueryOptions(keeper badgeKeeper.Keeper, keeper2 protocolKeeper.Keeper) []wasmKeeper.Option {
	badgeQueryOptions := wasmKeeper.WithQueryPlugins(badgeQueryPlugins(keeper ,keeper2))
	return []wasm.Option{badgeQueryOptions}
}

func badgeMsgEncoders() *wasmKeeper.MessageEncoders {
	return &wasmKeeper.MessageEncoders{
		Custom: customBindings.EncodeBitBadgesModuleMessage(),
	}
}

// badgeQueryPlugins needs to be registered in test setup to handle custom query callbacks
func badgeQueryPlugins(bk badgeKeeper.Keeper, pk protocolKeeper.Keeper) *wasmKeeper.QueryPlugins {
	return &wasmKeeper.QueryPlugins{
		Custom: customBindings.PerformCustomBitBadgesModuleQuery(bk, pk),
	}
}
