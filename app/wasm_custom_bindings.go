package app

import (
	anchorKeeper "bitbadgeschain/x/anchor/keeper"
	badgeKeeper "bitbadgeschain/x/badges/keeper"
	mapsKeeper "bitbadgeschain/x/maps/keeper"

	wasmKeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"

	customBindings "bitbadgeschain/custom-bindings"
)

func GetCustomMsgEncodersOptions() []wasmKeeper.Option {
	badgeEncodingOptions := wasmKeeper.WithMessageEncoders(badgeMsgEncoders())
	return []wasmKeeper.Option{badgeEncodingOptions}
}

func GetCustomMsgQueryOptions(keeper badgeKeeper.Keeper, anchorKeeper anchorKeeper.Keeper, mapsKeeper mapsKeeper.Keeper) []wasmKeeper.Option {
	badgeQueryOptions := wasmKeeper.WithQueryPlugins(badgeQueryPlugins(keeper, anchorKeeper, mapsKeeper))
	return []wasmKeeper.Option{badgeQueryOptions}
}

func badgeMsgEncoders() *wasmKeeper.MessageEncoders {
	return &wasmKeeper.MessageEncoders{
		Custom: customBindings.EncodeBitBadgesModuleMessage(),
	}
}

// badgeQueryPlugins needs to be registered in test setup to handle custom query callbacks
func badgeQueryPlugins(bk badgeKeeper.Keeper, anchorKeeper anchorKeeper.Keeper, mapsKeeper mapsKeeper.Keeper) *wasmKeeper.QueryPlugins {
	return &wasmKeeper.QueryPlugins{
		Custom: customBindings.PerformCustomBitBadgesModuleQuery(bk, anchorKeeper, mapsKeeper),
	}
}
