package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"

	encodingcodec "github.com/bitbadges/bitbadgeschain/encoding/codec"
)

//HACK: Even though the miscellaneous encoding/codec stuff is not used in the module, we register it here w/ the badges stuff (just needs to be registered once)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgTransferBadges{}, "badges/TransferBadges", nil)
	cdc.RegisterConcrete(&MsgDeleteCollection{}, "badges/DeleteCollection", nil)
	cdc.RegisterConcrete(&MsgUpdateUserApprovals{}, "badges/UpdateUserApprovals", nil)
	cdc.RegisterConcrete(&MsgUniversalUpdateCollection{}, "badges/UniversalUpdateCollection", nil)
	cdc.RegisterConcrete(&MsgCreateAddressLists{}, "badges/CreateAddressLists", nil)
	cdc.RegisterConcrete(&MsgCreateCollection{}, "badges/CreateCollection", nil)
	cdc.RegisterConcrete(&MsgUpdateCollection{}, "badges/UpdateCollection", nil)
	cdc.RegisterConcrete(&MsgCreateDynamicStore{}, "badges/CreateDynamicStore", nil)
	cdc.RegisterConcrete(&MsgUpdateDynamicStore{}, "badges/UpdateDynamicStore", nil)
	cdc.RegisterConcrete(&MsgDeleteDynamicStore{}, "badges/DeleteDynamicStore", nil)
	cdc.RegisterConcrete(&MsgSetDynamicStoreValue{}, "badges/SetDynamicStoreValue", nil)

	encodingcodec.RegisterLegacyAminoCodec(cdc)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgTransferBadges{},
		&MsgDeleteCollection{},
		&MsgUpdateUserApprovals{},
		&MsgUniversalUpdateCollection{},
		&MsgCreateAddressLists{},
		&MsgCreateCollection{},
		&MsgUpdateCollection{},
		&MsgCreateDynamicStore{},
		&MsgUpdateDynamicStore{},
		&MsgDeleteDynamicStore{},
		&MsgSetDynamicStoreValue{},
	)
	// this line is used by starport scaffolding # 3

	encodingcodec.RegisterInterfaces(registry)

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
	AminoCdc  = codec.NewAminoCodec(Amino)
)
