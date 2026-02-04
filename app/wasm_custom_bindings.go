package app

import (
	wasmKeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	anchorKeeper "github.com/bitbadges/bitbadgeschain/x/anchor/keeper"
	tokenizationKeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	gammKeeper "github.com/bitbadges/bitbadgeschain/x/gamm/keeper"
	managersplitterKeeper "github.com/bitbadges/bitbadgeschain/x/managersplitter/keeper"
	mapsKeeper "github.com/bitbadges/bitbadgeschain/x/maps/keeper"

	customBindings "github.com/bitbadges/bitbadgeschain/custom-bindings"
)

func GetCustomMsgEncodersOptions() []wasmKeeper.Option {
	tokenizationEncodingOptions := wasmKeeper.WithMessageEncoders(tokenizationMsgEncoders())
	return []wasmKeeper.Option{tokenizationEncodingOptions}
}

func GetCustomMsgQueryOptions(keeper tokenizationKeeper.Keeper, anchorKeeper anchorKeeper.Keeper, mapsKeeper mapsKeeper.Keeper, gk gammKeeper.Keeper, msk managersplitterKeeper.Keeper) []wasmKeeper.Option {
	tokenizationQueryOptions := wasmKeeper.WithQueryPlugins(tokenizationQueryPlugins(keeper, anchorKeeper, mapsKeeper, gk, msk))
	return []wasmKeeper.Option{tokenizationQueryOptions}
}

func tokenizationMsgEncoders() *wasmKeeper.MessageEncoders {
	return &wasmKeeper.MessageEncoders{
		Custom: customBindings.EncodeBitBadgesModuleMessage(),
	}
}

// tokenizationQueryPlugins needs to be registered in test setup to handle custom query callbacks
func tokenizationQueryPlugins(bk tokenizationKeeper.Keeper, anchorKeeper anchorKeeper.Keeper, mapsKeeper mapsKeeper.Keeper, gk gammKeeper.Keeper, msk managersplitterKeeper.Keeper) *wasmKeeper.QueryPlugins {
	return &wasmKeeper.QueryPlugins{
		Custom: customBindings.PerformCustomBitBadgesModuleQuery(bk, anchorKeeper, mapsKeeper, gk, msk),
	}
}
