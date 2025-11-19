package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

// RegisterLegacyAminoCodec registers the necessary x/gamm interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*CFMMPoolI)(nil), nil)
	cdc.RegisterConcrete(&MsgJoinPool{}, "gamm/JoinPool", nil)
	cdc.RegisterConcrete(&MsgExitPool{}, "gamm/ExitPool", nil)
	cdc.RegisterConcrete(&MsgSwapExactAmountIn{}, "gamm/SwapExactAmountIn", nil)
	cdc.RegisterConcrete(&MsgSwapExactAmountOut{}, "gamm/SwapExactAmountOut", nil)
	cdc.RegisterConcrete(&MsgJoinSwapExternAmountIn{}, "gamm/JoinSwapExternAmountIn", nil)
	cdc.RegisterConcrete(&MsgJoinSwapShareAmountOut{}, "gamm/JoinSwapShareAmountOut", nil)
	cdc.RegisterConcrete(&MsgExitSwapExternAmountOut{}, "gamm/ExitSwapExternAmountOut", nil)
	cdc.RegisterConcrete(&MsgExitSwapShareAmountIn{}, "gamm/ExitSwapShareAmountIn", nil)
	cdc.RegisterConcrete(&MsgSwapExactAmountInWithIBCTransfer{}, "gamm/SwapExactAmountInWithIBCTransfer", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgJoinPool{},
		&MsgExitPool{},
		&MsgSwapExactAmountIn{},
		&MsgSwapExactAmountOut{},
		&MsgJoinSwapExternAmountIn{},
		&MsgJoinSwapShareAmountOut{},
		&MsgExitSwapExternAmountOut{},
		&MsgExitSwapShareAmountIn{},
		&MsgSwapExactAmountInWithIBCTransfer{},
	)

	registry.RegisterImplementations(
		(*govtypesv1.Content)(nil),
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
