package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	// Proto messages will be registered here once proto files are added.
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateCouncil{},
		&MsgPropose{},
		&MsgVote{},
		&MsgExecuteProposal{},
	)
}

func init() {
	RegisterCodec(Amino)
	Amino.Seal()
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
	AminoCdc  = codec.NewAminoCodec(Amino)
)
