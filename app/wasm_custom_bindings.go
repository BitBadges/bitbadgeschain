package app

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmKeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	badgeCustomBindings "github.com/bitbadges/bitbadgeschain/x/badges/custom-bindings"
	badgeKeeper "github.com/bitbadges/bitbadgeschain/x/badges/keeper"

	protocolCustomBindings "github.com/bitbadges/bitbadgeschain/x/protocols/custom-bindings"
	protocolKeeper "github.com/bitbadges/bitbadgeschain/x/protocols/keeper"
)

func GetCustomMsgEncodersOptions() []wasmKeeper.Option {
	badgeEncodingOptions := wasmKeeper.WithMessageEncoders(badgeMsgEncoders())
	protocolEncodingOptions := wasmKeeper.WithMessageEncoders(protocolMsgEncoders())
	return []wasm.Option{badgeEncodingOptions, protocolEncodingOptions}
}

func GetCustomMsgQueryOptions(keeper badgeKeeper.Keeper, keeper2 protocolKeeper.Keeper) []wasmKeeper.Option {
	badgeQueryOptions := wasmKeeper.WithQueryPlugins(badgeQueryPlugins(keeper))
	protocolQueryOptions := wasmKeeper.WithQueryPlugins(protocolQueryPlugins(keeper2))
	return []wasm.Option{badgeQueryOptions, protocolQueryOptions}
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


func protocolMsgEncoders() *wasmKeeper.MessageEncoders {
	return &wasmKeeper.MessageEncoders{
		Custom: protocolCustomBindings.EncodeProtocolMessage(),
	}
}

// protocolQueryPlugins needs to be registered in test setup to handle custom query callbacks
func protocolQueryPlugins(keeper protocolKeeper.Keeper) *wasmKeeper.QueryPlugins {
	return &wasmKeeper.QueryPlugins{
		Custom: protocolCustomBindings.PerformCustomProtocolQuery(keeper),
	}
}