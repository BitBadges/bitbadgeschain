package app

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/gogoproto/proto"

	ibcclienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"
)

// LegacyLocalhostClientState is the legacy 09-localhost ClientState type from IBC v8.
// This type is registered for backward compatibility during genesis import from
// chains that were running IBC v8 with the localhost light client.
// In IBC v10, the localhost light client was redesigned to be stateless and no
// longer uses a ClientState proto type.
//
// Type URL: /ibc.lightclients.localhost.v2.ClientState
type LegacyLocalhostClientState struct {
	LatestHeight ibcclienttypes.Height `protobuf:"bytes,1,opt,name=latest_height,json=latestHeight,proto3" json:"latest_height"`
}

var _ ibcexported.ClientState = (*LegacyLocalhostClientState)(nil)

func (cs *LegacyLocalhostClientState) Reset()         { *cs = LegacyLocalhostClientState{} }
func (cs *LegacyLocalhostClientState) String() string { return proto.CompactTextString(cs) }
func (cs *LegacyLocalhostClientState) ProtoMessage()  {}

// ClientType returns the localhost client type.
func (cs *LegacyLocalhostClientState) ClientType() string {
	return "09-localhost"
}

// Validate performs a basic validation of the client state fields.
func (cs *LegacyLocalhostClientState) Validate() error {
	return nil
}

// RegisterLegacyLocalhostInterfaces registers the legacy localhost v2 ClientState type
// for backward compatibility when importing genesis from IBC v8 chains.
func RegisterLegacyLocalhostInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*ibcexported.ClientState)(nil),
		&LegacyLocalhostClientState{},
	)
}

func init() {
	// Register the type URL for the legacy localhost ClientState
	proto.RegisterType((*LegacyLocalhostClientState)(nil), "ibc.lightclients.localhost.v2.ClientState")
}
