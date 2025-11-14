package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"

	encodingcodec "github.com/bitbadges/bitbadgeschain/encoding/codec"
)

// NOTE: The miscellaneous encoding/codec registration is included here to ensure
// all necessary codec types are registered once for the tokens module.
// This is required for proper serialization/deserialization across the module.

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgTransferTokens{}, "badges/TransferTokens", nil)
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
	cdc.RegisterConcrete(&MsgSetValidTokenIds{}, "badges/SetValidTokenIds", nil)
	cdc.RegisterConcrete(&MsgSetManager{}, "badges/SetManager", nil)
	cdc.RegisterConcrete(&MsgSetCollectionMetadata{}, "badges/SetCollectionMetadata", nil)
	cdc.RegisterConcrete(&MsgSetTokenMetadata{}, "badges/SetTokenMetadata", nil)
	cdc.RegisterConcrete(&MsgSetCustomData{}, "badges/SetCustomData", nil)
	cdc.RegisterConcrete(&MsgSetStandards{}, "badges/SetStandards", nil)
	cdc.RegisterConcrete(&MsgSetCollectionApprovals{}, "badges/SetCollectionApprovals", nil)
	cdc.RegisterConcrete(&MsgSetIsArchived{}, "badges/SetIsArchived", nil)
	cdc.RegisterConcrete(&MsgSetReservedProtocolAddress{}, "badges/SetReservedProtocolAddress", nil)

	encodingcodec.RegisterLegacyAminoCodec(cdc)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgTransferTokens{},
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
		&MsgSetValidTokenIds{},
		&MsgSetManager{},
		&MsgSetCollectionMetadata{},
		&MsgSetTokenMetadata{},
		&MsgSetCustomData{},
		&MsgSetStandards{},
		&MsgSetCollectionApprovals{},
		&MsgSetIsArchived{},
		&MsgSetReservedProtocolAddress{},
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
