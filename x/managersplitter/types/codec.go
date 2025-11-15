package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateManagerSplitter{}, "managersplitter/CreateManagerSplitter", nil)
	cdc.RegisterConcrete(&MsgUpdateManagerSplitter{}, "managersplitter/UpdateManagerSplitter", nil)
	cdc.RegisterConcrete(&MsgDeleteManagerSplitter{}, "managersplitter/DeleteManagerSplitter", nil)
	cdc.RegisterConcrete(&MsgExecuteUniversalUpdateCollection{}, "managersplitter/ExecuteUniversalUpdateCollection", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateManagerSplitter{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateManagerSplitter{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgDeleteManagerSplitter{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgExecuteUniversalUpdateCollection{},
	)

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

