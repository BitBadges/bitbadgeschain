package app

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmKeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	badgeCustomBindings "github.com/bitbadges/bitbadgeschain/x/badges/custom-bindings"
	badgeKeeper "github.com/bitbadges/bitbadgeschain/x/badges/keeper"
)

func GetCustomMsgEncodersOptions() []wasmKeeper.Option {
	badgeEncodingOptions := wasmKeeper.WithMessageEncoders(badgeMsgEncoders())
	return []wasm.Option{badgeEncodingOptions}
}

func GetCustomMsgQueryOptions(keeper badgeKeeper.Keeper) []wasmKeeper.Option {
	badgeQueryOptions := wasmKeeper.WithQueryPlugins(badgeQueryPlugins(keeper))
	return []wasm.Option{badgeQueryOptions}
}

func badgeMsgEncoders() *wasmKeeper.MessageEncoders {
	return &wasmKeeper.MessageEncoders{
		Custom: badgeCustomBindings.EncodeBadgeMessage(),
	}
}

// badgeQueryPlugins needs to be registered in test setup to handle custom query callbacks
func badgeQueryPlugins(keeper badgeKeeper.Keeper) *wasmKeeper.QueryPlugins {
	return &wasmKeeper.QueryPlugins{
		Custom: badgeCustomBindings.PerformCustomBadgeQuery(keeper),
	}
}
