package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgNewCollection{}, "badges/NewCollection", nil)
	cdc.RegisterConcrete(&MsgMintBadge{}, "badges/MintBadge", nil)
	cdc.RegisterConcrete(&MsgTransferBadge{}, "badges/TransferBadge", nil)
	cdc.RegisterConcrete(&MsgSetApproval{}, "badges/SetApproval", nil)
	cdc.RegisterConcrete(&MsgUpdateDisallowedTransfers{}, "badges/UpdateDisallowedTransfers", nil)
	cdc.RegisterConcrete(&MsgUpdateUris{}, "badges/UpdateUris", nil)
	cdc.RegisterConcrete(&MsgUpdatePermissions{}, "badges/UpdatePermissions", nil)
	cdc.RegisterConcrete(&MsgTransferManager{}, "badges/TransferManager", nil)
	cdc.RegisterConcrete(&MsgRequestTransferManager{}, "badges/RequestTransferManager", nil)
	cdc.RegisterConcrete(&MsgUpdateBytes{}, "badges/UpdateBytes", nil)
	cdc.RegisterConcrete(&MsgRegisterAddresses{}, "badges/RegisterAddresses", nil)
	cdc.RegisterConcrete(&MsgClaimBadge{}, "badges/ClaimBadge", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgNewCollection{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgMintBadge{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgTransferBadge{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSetApproval{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateDisallowedTransfers{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateUris{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdatePermissions{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgTransferManager{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRequestTransferManager{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateBytes{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRegisterAddresses{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgClaimBadge{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)

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
