package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgNewBadge{}, "badges/NewBadge", nil)
	cdc.RegisterConcrete(&MsgNewSubBadge{}, "badges/NewSubBadge", nil)
	cdc.RegisterConcrete(&MsgTransferBadge{}, "badges/TransferBadge", nil)
	cdc.RegisterConcrete(&MsgRequestTransferBadge{}, "badges/RequestTransferBadge", nil)
	cdc.RegisterConcrete(&MsgHandlePendingTransfer{}, "badges/HandlePendingTransfer", nil)
	cdc.RegisterConcrete(&MsgSetApproval{}, "badges/SetApproval", nil)
	cdc.RegisterConcrete(&MsgRevokeBadge{}, "badges/RevokeBadge", nil)
	cdc.RegisterConcrete(&MsgFreezeAddress{}, "badges/FreezeAddress", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgNewBadge{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgNewSubBadge{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgTransferBadge{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRequestTransferBadge{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgHandlePendingTransfer{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSetApproval{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRevokeBadge{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgFreezeAddress{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
