package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateMap{}, "maps/CreateMap", nil)
	cdc.RegisterConcrete(&MsgUpdateMap{}, "maps/UpdateMap", nil)
	cdc.RegisterConcrete(&MsgDeleteMap{}, "maps/DeleteMap", nil)
	cdc.RegisterConcrete(&MsgSetValue{}, "maps/SetValue", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateMap{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateMap{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgDeleteMap{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSetValue{},
	)

	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

// NOTE: This is required for the GetSignBytes function
func init() {
	RegisterCodec(Amino)
	Amino.Seal()
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
	// AminoCdc is a amino codec created to support amino JSON compatible msgs.
	AminoCdc = codec.NewAminoCodec(Amino)
)
