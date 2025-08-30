package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"

	encodingcodec "github.com/bitbadges/bitbadgeschain/encoding/codec"
	oldtypes "github.com/bitbadges/bitbadgeschain/x/badges/types/v13"
)

const (
	OldVersion     = "V13"
	CurrentVersion = "V14"
)

//HACK: Even though the miscellaneous encoding/codec stuff is not used in the module, we register it here w/ the tokens stuff (just needs to be registered once)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&oldtypes.MsgTransferBadges{}, "badges/TransferBadges"+OldVersion, nil)
	cdc.RegisterConcrete(&oldtypes.MsgDeleteCollection{}, "badges/DeleteCollection"+OldVersion, nil)
	cdc.RegisterConcrete(&oldtypes.MsgUpdateUserApprovals{}, "badges/UpdateUserApprovals"+OldVersion, nil)
	cdc.RegisterConcrete(&oldtypes.MsgUniversalUpdateCollection{}, "badges/UniversalUpdateCollection"+OldVersion, nil)
	cdc.RegisterConcrete(&oldtypes.MsgCreateAddressLists{}, "badges/CreateAddressLists"+OldVersion, nil)
	cdc.RegisterConcrete(&oldtypes.MsgCreateCollection{}, "badges/CreateCollection"+OldVersion, nil)
	cdc.RegisterConcrete(&oldtypes.MsgUpdateCollection{}, "badges/UpdateCollection"+OldVersion, nil)
	cdc.RegisterConcrete(&oldtypes.MsgCreateDynamicStore{}, "badges/CreateDynamicStore"+OldVersion, nil)
	cdc.RegisterConcrete(&oldtypes.MsgUpdateDynamicStore{}, "badges/UpdateDynamicStore"+OldVersion, nil)
	cdc.RegisterConcrete(&oldtypes.MsgDeleteDynamicStore{}, "badges/DeleteDynamicStore"+OldVersion, nil)
	cdc.RegisterConcrete(&oldtypes.MsgSetDynamicStoreValue{}, "badges/SetDynamicStoreValue"+OldVersion, nil)
	cdc.RegisterConcrete(&oldtypes.MsgSetIncomingApproval{}, "badges/SetIncomingApproval"+OldVersion, nil)
	cdc.RegisterConcrete(&oldtypes.MsgDeleteIncomingApproval{}, "badges/DeleteIncomingApproval"+OldVersion, nil)
	cdc.RegisterConcrete(&oldtypes.MsgSetOutgoingApproval{}, "badges/SetOutgoingApproval"+OldVersion, nil)
	cdc.RegisterConcrete(&oldtypes.MsgDeleteOutgoingApproval{}, "badges/DeleteOutgoingApproval"+OldVersion, nil)
	cdc.RegisterConcrete(&oldtypes.MsgPurgeApprovals{}, "badges/PurgeApprovals"+OldVersion, nil)
	cdc.RegisterConcrete(&oldtypes.MsgSetValidBadgeIds{}, "badges/SetValidBadgeIds"+OldVersion, nil)
	cdc.RegisterConcrete(&oldtypes.MsgSetManager{}, "badges/SetManager"+OldVersion, nil)
	cdc.RegisterConcrete(&oldtypes.MsgSetCollectionMetadata{}, "badges/SetCollectionMetadata"+OldVersion, nil)
	cdc.RegisterConcrete(&oldtypes.MsgSetBadgeMetadata{}, "badges/SetBadgeMetadata"+OldVersion, nil)
	cdc.RegisterConcrete(&oldtypes.MsgSetCustomData{}, "badges/SetCustomData"+OldVersion, nil)
	cdc.RegisterConcrete(&oldtypes.MsgSetStandards{}, "badges/SetStandards"+OldVersion, nil)
	cdc.RegisterConcrete(&oldtypes.MsgSetCollectionApprovals{}, "badges/SetCollectionApprovals"+OldVersion, nil)
	cdc.RegisterConcrete(&oldtypes.MsgSetIsArchived{}, "badges/SetIsArchived"+OldVersion, nil)

	//Register new types but versioned
	cdc.RegisterConcrete(&MsgTransferBadges{}, "badges/TransferBadges"+CurrentVersion, nil)
	cdc.RegisterConcrete(&MsgDeleteCollection{}, "badges/DeleteCollection"+CurrentVersion, nil)
	cdc.RegisterConcrete(&MsgUpdateUserApprovals{}, "badges/UpdateUserApprovals"+CurrentVersion, nil)
	cdc.RegisterConcrete(&MsgUniversalUpdateCollection{}, "badges/UniversalUpdateCollection"+CurrentVersion, nil)
	cdc.RegisterConcrete(&MsgCreateAddressLists{}, "badges/CreateAddressLists"+CurrentVersion, nil)
	cdc.RegisterConcrete(&MsgCreateCollection{}, "badges/CreateCollection"+CurrentVersion, nil)
	cdc.RegisterConcrete(&MsgUpdateCollection{}, "badges/UpdateCollection"+CurrentVersion, nil)
	cdc.RegisterConcrete(&MsgCreateDynamicStore{}, "badges/CreateDynamicStore"+CurrentVersion, nil)
	cdc.RegisterConcrete(&MsgUpdateDynamicStore{}, "badges/UpdateDynamicStore"+CurrentVersion, nil)
	cdc.RegisterConcrete(&MsgDeleteDynamicStore{}, "badges/DeleteDynamicStore"+CurrentVersion, nil)
	cdc.RegisterConcrete(&MsgSetDynamicStoreValue{}, "badges/SetDynamicStoreValue"+CurrentVersion, nil)
	cdc.RegisterConcrete(&MsgSetIncomingApproval{}, "badges/SetIncomingApproval"+CurrentVersion, nil)
	cdc.RegisterConcrete(&MsgDeleteIncomingApproval{}, "badges/DeleteIncomingApproval"+CurrentVersion, nil)
	cdc.RegisterConcrete(&MsgSetOutgoingApproval{}, "badges/SetOutgoingApproval"+CurrentVersion, nil)
	cdc.RegisterConcrete(&MsgDeleteOutgoingApproval{}, "badges/DeleteOutgoingApproval"+CurrentVersion, nil)
	cdc.RegisterConcrete(&MsgPurgeApprovals{}, "badges/PurgeApprovals"+CurrentVersion, nil)
	cdc.RegisterConcrete(&MsgSetValidBadgeIds{}, "badges/SetValidBadgeIds"+CurrentVersion, nil)
	cdc.RegisterConcrete(&MsgSetManager{}, "badges/SetManager"+CurrentVersion, nil)
	cdc.RegisterConcrete(&MsgSetCollectionMetadata{}, "badges/SetCollectionMetadata"+CurrentVersion, nil)
	cdc.RegisterConcrete(&MsgSetBadgeMetadata{}, "badges/SetBadgeMetadata"+CurrentVersion, nil)
	cdc.RegisterConcrete(&MsgSetCustomData{}, "badges/SetCustomData"+CurrentVersion, nil)
	cdc.RegisterConcrete(&MsgSetStandards{}, "badges/SetStandards"+CurrentVersion, nil)
	cdc.RegisterConcrete(&MsgSetCollectionApprovals{}, "badges/SetCollectionApprovals"+CurrentVersion, nil)
	cdc.RegisterConcrete(&MsgSetIsArchived{}, "badges/SetIsArchived"+CurrentVersion, nil)

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

		&oldtypes.MsgTransferBadges{},
		&oldtypes.MsgDeleteCollection{},
		&oldtypes.MsgUpdateUserApprovals{},
		&oldtypes.MsgUniversalUpdateCollection{},
		&oldtypes.MsgCreateAddressLists{},
		&oldtypes.MsgCreateCollection{},
		&oldtypes.MsgUpdateCollection{},
		&oldtypes.MsgCreateDynamicStore{},
		&oldtypes.MsgUpdateDynamicStore{},
		&oldtypes.MsgDeleteDynamicStore{},
		&oldtypes.MsgSetDynamicStoreValue{},
		&oldtypes.MsgSetIncomingApproval{},
		&oldtypes.MsgDeleteIncomingApproval{},
		&oldtypes.MsgSetOutgoingApproval{},
		&oldtypes.MsgDeleteOutgoingApproval{},
		&oldtypes.MsgPurgeApprovals{},
		&oldtypes.MsgSetValidBadgeIds{},
		&oldtypes.MsgSetManager{},
		&oldtypes.MsgSetCollectionMetadata{},
		&oldtypes.MsgSetBadgeMetadata{},
		&oldtypes.MsgSetCustomData{},
		&oldtypes.MsgSetStandards{},
		&oldtypes.MsgSetCollectionApprovals{},
		&oldtypes.MsgSetIsArchived{},
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
