package app

import (
	wasmKeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	anchorKeeper "github.com/bitbadges/bitbadgeschain/x/anchor/keeper"
	badgeKeeper "github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	gammKeeper "github.com/bitbadges/bitbadgeschain/x/gamm/keeper"
	managersplitterKeeper "github.com/bitbadges/bitbadgeschain/x/managersplitter/keeper"
	mapsKeeper "github.com/bitbadges/bitbadgeschain/x/maps/keeper"

	customBindings "github.com/bitbadges/bitbadgeschain/custom-bindings"
)

func GetCustomMsgEncodersOptions() []wasmKeeper.Option {
	badgeEncodingOptions := wasmKeeper.WithMessageEncoders(badgeMsgEncoders())
	return []wasmKeeper.Option{badgeEncodingOptions}
}

func GetCustomMsgQueryOptions(keeper badgeKeeper.Keeper, anchorKeeper anchorKeeper.Keeper, mapsKeeper mapsKeeper.Keeper, gk gammKeeper.Keeper, msk managersplitterKeeper.Keeper) []wasmKeeper.Option {
	badgeQueryOptions := wasmKeeper.WithQueryPlugins(badgeQueryPlugins(keeper, anchorKeeper, mapsKeeper, gk, msk))
	return []wasmKeeper.Option{badgeQueryOptions}
}

func badgeMsgEncoders() *wasmKeeper.MessageEncoders {
	return &wasmKeeper.MessageEncoders{
		Custom: customBindings.EncodeBitBadgesModuleMessage(),
	}
}

// badgeQueryPlugins needs to be registered in test setup to handle custom query callbacks
func badgeQueryPlugins(bk badgeKeeper.Keeper, anchorKeeper anchorKeeper.Keeper, mapsKeeper mapsKeeper.Keeper, gk gammKeeper.Keeper, msk managersplitterKeeper.Keeper) *wasmKeeper.QueryPlugins {
	return &wasmKeeper.QueryPlugins{
		Custom: customBindings.PerformCustomBitBadgesModuleQuery(bk, anchorKeeper, mapsKeeper, gk, msk),
	}
}
