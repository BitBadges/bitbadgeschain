package balancer

import (
	poolmanagertypes "github.com/bitbadges/bitbadgeschain/x/poolmanager/types"

	types "github.com/bitbadges/bitbadgeschain/x/gamm/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	proto "github.com/cosmos/gogoproto/proto"
)

// RegisterLegacyAminoCodec registers the necessary x/gamm interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&Pool{}, "gamm/BalancerPool", nil)
	cdc.RegisterConcrete(&MsgCreateBalancerPool{}, "gamm/create-balancer-pool", nil)
	cdc.RegisterConcrete(&PoolParams{}, "gamm/BalancerPoolParams", nil)
}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterInterface(
		"poolmanager.v1beta1.PoolI",
		(*poolmanagertypes.PoolI)(nil),
		&Pool{},
	)
	registry.RegisterInterface(
		"gamm.v1beta1.PoolI", // N.B.: the old proto-path is preserved for backwards-compatibility.
		(*types.CFMMPoolI)(nil),
		&Pool{},
	)
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreateBalancerPool{},
	)
	registry.RegisterImplementations(
		(*proto.Message)(nil),
		&PoolParams{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino = codec.NewLegacyAmino()
)

func init() {
	RegisterLegacyAminoCodec(amino)
	amino.Seal()
}
