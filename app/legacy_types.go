package app

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
)

// Legacy types for backward compatibility during genesis export.
// These types are registered so that old governance proposals referencing
// removed/renamed module messages can be properly marshaled/unmarshaled.
//
// These stubs are needed because:
// - WASM module has been removed from the chain
// - badges module has been renamed to tokenization
// Old governance proposals may still reference these old message types.

// LegacyWasmMsgUpdateParams is the legacy cosmwasm.wasm.v1.MsgUpdateParams type.
// Type URL: /cosmwasm.wasm.v1.MsgUpdateParams
type LegacyWasmMsgUpdateParams struct {
	Authority string                `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority,omitempty"`
	Params    LegacyWasmParams      `protobuf:"bytes,2,opt,name=params,proto3" json:"params"`
}

func (m *LegacyWasmMsgUpdateParams) Reset()         { *m = LegacyWasmMsgUpdateParams{} }
func (m *LegacyWasmMsgUpdateParams) String() string { return proto.CompactTextString(m) }
func (m *LegacyWasmMsgUpdateParams) ProtoMessage()  {}

// LegacyWasmParams is the legacy cosmwasm.wasm.v1.Params type.
type LegacyWasmParams struct {
	CodeUploadAccess             LegacyWasmAccessConfig `protobuf:"bytes,1,opt,name=code_upload_access,json=codeUploadAccess,proto3" json:"code_upload_access"`
	InstantiateDefaultPermission int32                  `protobuf:"varint,2,opt,name=instantiate_default_permission,json=instantiateDefaultPermission,proto3" json:"instantiate_default_permission,omitempty"`
}

func (m *LegacyWasmParams) Reset()         { *m = LegacyWasmParams{} }
func (m *LegacyWasmParams) String() string { return proto.CompactTextString(m) }
func (m *LegacyWasmParams) ProtoMessage()  {}

// LegacyWasmAccessConfig is the legacy cosmwasm.wasm.v1.AccessConfig type.
type LegacyWasmAccessConfig struct {
	Permission int32    `protobuf:"varint,1,opt,name=permission,proto3" json:"permission,omitempty"`
	Addresses  []string `protobuf:"bytes,3,rep,name=addresses,proto3" json:"addresses,omitempty"`
}

func (m *LegacyWasmAccessConfig) Reset()         { *m = LegacyWasmAccessConfig{} }
func (m *LegacyWasmAccessConfig) String() string { return proto.CompactTextString(m) }
func (m *LegacyWasmAccessConfig) ProtoMessage()  {}

// Implement sdk.Msg interface for LegacyWasmMsgUpdateParams
var _ sdk.Msg = (*LegacyWasmMsgUpdateParams)(nil)

func (m *LegacyWasmMsgUpdateParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{addr}
}

func (m *LegacyWasmMsgUpdateParams) ValidateBasic() error {
	return nil
}

// RegisterLegacyWasmInterfaces registers the legacy WASM types for backward
// compatibility when exporting/importing genesis with old governance proposals.
func RegisterLegacyWasmInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&LegacyWasmMsgUpdateParams{},
	)
}

// =============================================================================
// Legacy Badges Module Types (renamed to tokenization)
// =============================================================================

// LegacyBadgesMsgUpdateParams is the legacy badges.MsgUpdateParams type.
// Type URL: /badges.MsgUpdateParams
type LegacyBadgesMsgUpdateParams struct {
	Authority string            `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority,omitempty"`
	Params    LegacyBadgesParams `protobuf:"bytes,2,opt,name=params,proto3" json:"params"`
}

func (m *LegacyBadgesMsgUpdateParams) Reset()         { *m = LegacyBadgesMsgUpdateParams{} }
func (m *LegacyBadgesMsgUpdateParams) String() string { return proto.CompactTextString(m) }
func (m *LegacyBadgesMsgUpdateParams) ProtoMessage()  {}

// LegacyBadgesParams is the legacy badges.Params type.
type LegacyBadgesParams struct{}

func (m *LegacyBadgesParams) Reset()         { *m = LegacyBadgesParams{} }
func (m *LegacyBadgesParams) String() string { return proto.CompactTextString(m) }
func (m *LegacyBadgesParams) ProtoMessage()  {}

// Implement sdk.Msg interface for LegacyBadgesMsgUpdateParams
var _ sdk.Msg = (*LegacyBadgesMsgUpdateParams)(nil)

func (m *LegacyBadgesMsgUpdateParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{addr}
}

func (m *LegacyBadgesMsgUpdateParams) ValidateBasic() error {
	return nil
}

// RegisterLegacyBadgesInterfaces registers the legacy badges types for backward
// compatibility when exporting/importing genesis with old governance proposals.
func RegisterLegacyBadgesInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&LegacyBadgesMsgUpdateParams{},
	)
}

func init() {
	// Register the type URLs for the legacy WASM types
	proto.RegisterType((*LegacyWasmMsgUpdateParams)(nil), "cosmwasm.wasm.v1.MsgUpdateParams")
	proto.RegisterType((*LegacyWasmParams)(nil), "cosmwasm.wasm.v1.Params")
	proto.RegisterType((*LegacyWasmAccessConfig)(nil), "cosmwasm.wasm.v1.AccessConfig")

	// Register the type URLs for the legacy badges types
	proto.RegisterType((*LegacyBadgesMsgUpdateParams)(nil), "badges.MsgUpdateParams")
	proto.RegisterType((*LegacyBadgesParams)(nil), "badges.Params")
}
