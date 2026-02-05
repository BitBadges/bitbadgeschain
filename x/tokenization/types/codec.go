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
	cdc.RegisterConcrete(&MsgTransferTokens{}, "tokenization/TransferTokens", nil)
	cdc.RegisterConcrete(&MsgDeleteCollection{}, "tokenization/DeleteCollection", nil)
	cdc.RegisterConcrete(&MsgUpdateUserApprovals{}, "tokenization/UpdateUserApprovals", nil)
	cdc.RegisterConcrete(&MsgUniversalUpdateCollection{}, "tokenization/UniversalUpdateCollection", nil)
	cdc.RegisterConcrete(&MsgCreateAddressLists{}, "tokenization/CreateAddressLists", nil)
	cdc.RegisterConcrete(&MsgCreateCollection{}, "tokenization/CreateCollection", nil)
	cdc.RegisterConcrete(&MsgUpdateCollection{}, "tokenization/UpdateCollection", nil)
	cdc.RegisterConcrete(&MsgCreateDynamicStore{}, "tokenization/CreateDynamicStore", nil)
	cdc.RegisterConcrete(&MsgUpdateDynamicStore{}, "tokenization/UpdateDynamicStore", nil)
	cdc.RegisterConcrete(&MsgDeleteDynamicStore{}, "tokenization/DeleteDynamicStore", nil)
	cdc.RegisterConcrete(&MsgSetDynamicStoreValue{}, "tokenization/SetDynamicStoreValue", nil)
	cdc.RegisterConcrete(&MsgSetIncomingApproval{}, "tokenization/SetIncomingApproval", nil)
	cdc.RegisterConcrete(&MsgDeleteIncomingApproval{}, "tokenization/DeleteIncomingApproval", nil)
	cdc.RegisterConcrete(&MsgSetOutgoingApproval{}, "tokenization/SetOutgoingApproval", nil)
	cdc.RegisterConcrete(&MsgDeleteOutgoingApproval{}, "tokenization/DeleteOutgoingApproval", nil)
	cdc.RegisterConcrete(&MsgPurgeApprovals{}, "tokenization/PurgeApprovals", nil)
	cdc.RegisterConcrete(&MsgSetValidTokenIds{}, "tokenization/SetValidTokenIds", nil)
	cdc.RegisterConcrete(&MsgSetManager{}, "tokenization/SetManager", nil)
	cdc.RegisterConcrete(&MsgSetCollectionMetadata{}, "tokenization/SetCollectionMetadata", nil)
	cdc.RegisterConcrete(&MsgSetTokenMetadata{}, "tokenization/SetTokenMetadata", nil)
	cdc.RegisterConcrete(&MsgSetCustomData{}, "tokenization/SetCustomData", nil)
	cdc.RegisterConcrete(&MsgSetStandards{}, "tokenization/SetStandards", nil)
	cdc.RegisterConcrete(&MsgSetCollectionApprovals{}, "tokenization/SetCollectionApprovals", nil)
	cdc.RegisterConcrete(&MsgSetIsArchived{}, "tokenization/SetIsArchived", nil)
	cdc.RegisterConcrete(&MsgSetReservedProtocolAddress{}, "tokenization/SetReservedProtocolAddress", nil)
	cdc.RegisterConcrete(&MsgCastVote{}, "tokenization/CastVote", nil)

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
		&MsgCastVote{},
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
