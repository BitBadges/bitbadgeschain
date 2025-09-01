package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers the necessary x/gamm interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgSwapExactAmountIn{}, "poolmanager/swap-exact-amount-in", nil)
	cdc.RegisterConcrete(&MsgSwapExactAmountOut{}, "poolmanager/swap-exact-amount-out", nil)
	cdc.RegisterConcrete(&MsgSplitRouteSwapExactAmountIn{}, "poolmanager/split-amount-in", nil)
	cdc.RegisterConcrete(&MsgSplitRouteSwapExactAmountOut{}, "poolmanager/split-amount-out", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgSwapExactAmountIn{},
		&MsgSwapExactAmountOut{},
		&MsgSplitRouteSwapExactAmountIn{},
		&MsgSplitRouteSwapExactAmountOut{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
