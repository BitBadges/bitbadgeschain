package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateProtocol{}, "protocols/CreateProtocol", nil)
	cdc.RegisterConcrete(&MsgUpdateProtocol{}, "protocols/UpdateProtocol", nil)
	cdc.RegisterConcrete(&MsgDeleteProtocol{}, "protocols/DeleteProtocol", nil)
	cdc.RegisterConcrete(&MsgSetCollectionForProtocol{}, "protocols/SetCollectionForProtocol", nil)
	cdc.RegisterConcrete(&MsgUnsetCollectionForProtocol{}, "protocols/UnsetCollectionForProtocol", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateProtocol{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateProtocol{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgDeleteProtocol{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSetCollectionForProtocol{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUnsetCollectionForProtocol{},
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
