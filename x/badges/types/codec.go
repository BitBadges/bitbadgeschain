package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"

	encodingcodec "github.com/bitbadges/bitbadgeschain/encoding/codec"
)

//HACK: Even though the miscellaneous encoding/codec stuff is not used in the module, we register it here w/ the tokens stuff (just needs to be registered once)

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
	cdc.RegisterConcrete(&MsgSetIncomingApproval{}, "badges/SetIncomingApproval", nil)
	cdc.RegisterConcrete(&MsgDeleteIncomingApproval{}, "badges/DeleteIncomingApproval", nil)
	cdc.RegisterConcrete(&MsgSetOutgoingApproval{}, "badges/SetOutgoingApproval", nil)
	cdc.RegisterConcrete(&MsgDeleteOutgoingApproval{}, "badges/DeleteOutgoingApproval", nil)
	cdc.RegisterConcrete(&MsgPurgeApprovals{}, "badges/PurgeApprovals", nil)
	cdc.RegisterConcrete(&MsgSetValidBadgeIds{}, "badges/SetValidBadgeIds", nil)
	cdc.RegisterConcrete(&MsgSetManager{}, "badges/SetManager", nil)
	cdc.RegisterConcrete(&MsgSetCollectionMetadata{}, "badges/SetCollectionMetadata", nil)
	cdc.RegisterConcrete(&MsgSetBadgeMetadata{}, "badges/SetBadgeMetadata", nil)
	cdc.RegisterConcrete(&MsgSetCustomData{}, "badges/SetCustomData", nil)
	cdc.RegisterConcrete(&MsgSetStandards{}, "badges/SetStandards", nil)
	cdc.RegisterConcrete(&MsgSetCollectionApprovals{}, "badges/SetCollectionApprovals", nil)
	cdc.RegisterConcrete(&MsgSetIsArchived{}, "badges/SetIsArchived", nil)
	cdc.RegisterConcrete(&MsgUnwrapIBCDenom{}, "badges/UnwrapIBCDenom", nil)

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
		&MsgSetIncomingApproval{},
		&MsgDeleteIncomingApproval{},
		&MsgSetOutgoingApproval{},
		&MsgDeleteOutgoingApproval{},
		&MsgPurgeApprovals{},
		&MsgSetValidBadgeIds{},
		&MsgSetManager{},
		&MsgSetCollectionMetadata{},
		&MsgSetBadgeMetadata{},
		&MsgSetCustomData{},
		&MsgSetStandards{},
		&MsgSetCollectionApprovals{},
		&MsgSetIsArchived{},
		&MsgUnwrapIBCDenom{},
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
