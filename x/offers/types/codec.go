package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	// this line is used by starport scaffolding # 1
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateProposal{}, "offers/CreateProposal", nil)
	cdc.RegisterConcrete(&MsgAcceptProposal{}, "offers/AcceptProposal", nil)
	cdc.RegisterConcrete(&MsgRejectAndDeleteProposal{}, "offers/RejectAndDeleteProposal", nil)
	cdc.RegisterConcrete(&MsgExecuteProposal{}, "offers/ExecuteProposal", nil)

	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateProposal{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAcceptProposal{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRejectAndDeleteProposal{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgExecuteProposal{},
	)
	// this line is used by starport scaffolding # 3

	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	// ModuleCdc references the global x/ibc-transfer module codec. Note, the codec
	// should ONLY be used in certain instances of tests and for JSON encoding.
	//
	// The actual codec used for serialization should be provided to x/ibc transfer and
	// defined at the application level.
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)

// NOTE: This is required for the GetSignBytes function
func init() {
	RegisterCodec(Amino)
	Amino.Seal()
}

var (
	Amino    = codec.NewLegacyAmino()
	AminoCdc = codec.NewAminoCodec(Amino)
)
