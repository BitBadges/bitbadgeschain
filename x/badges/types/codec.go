package types

import (
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgNewCollection{}, "badges/NewCollection", nil)
	cdc.RegisterConcrete(&MsgMintAndDistributeBadges{}, "badges/MintAndDistributeBadges", nil)
	cdc.RegisterConcrete(&MsgTransferBadge{}, "badges/TransferBadge", nil)
	cdc.RegisterConcrete(&MsgUpdateCollectionApprovedTransfers{}, "badges/UpdateCollectionApprovedTransfers", nil)
	cdc.RegisterConcrete(&MsgUpdateMetadata{}, "badges/UpdateMetadata", nil)
	cdc.RegisterConcrete(&MsgUpdateCollectionPermissions{}, "badges/UpdateCollectionPermissions", nil)
	cdc.RegisterConcrete(&MsgUpdateManager{}, "badges/UpdateManager", nil)
	cdc.RegisterConcrete(&MsgDeleteCollection{}, "badges/DeleteCollection", nil)
	cdc.RegisterConcrete(&MsgArchiveCollection{}, "badges/ArchiveCollection", nil)
	cdc.RegisterConcrete(&MsgForkCollection{}, "badges/ForkCollection", nil)
	cdc.RegisterConcrete(&MsgUpdateUserApprovedTransfers{}, "badges/UpdateUserApprovedTransfers", nil)
	cdc.RegisterConcrete(&MsgUpdateUserPermissions{}, "badges/UpdateUserPermissions", nil)
// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgNewCollection{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgMintAndDistributeBadges{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgTransferBadge{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateCollectionApprovedTransfers{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateMetadata{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateCollectionPermissions{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateManager{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgDeleteCollection{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgArchiveCollection{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgForkCollection{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateUserApprovedTransfers{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
	&MsgUpdateUserPermissions{},
)
// this line is used by starport scaffolding # 3

msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)

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
