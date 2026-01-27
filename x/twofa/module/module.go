package twofa

import (
	"context"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	badgeskeeper "github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/twofa/keeper"
	twofatypes "github.com/bitbadges/bitbadgeschain/x/twofa/types"
)

var (
	_ module.AppModuleBasic = (*AppModule)(nil)
	_ module.AppModule      = (*AppModule)(nil)
	_ appmodule.AppModule   = (*AppModule)(nil)
)

// AppModule implements the AppModule interface
type AppModule struct {
	cdc    codec.Codec
	keeper keeper.Keeper
}

func NewAppModule(
	cdc codec.Codec,
	keeper keeper.Keeper,
) AppModule {
	return AppModule{
		cdc:    cdc,
		keeper: keeper,
	}
}

// IsAppModule implements the appmodule.AppModule interface.
func (AppModule) IsAppModule() {}

// Name returns the name of the module as a string.
func (AppModule) Name() string {
	return twofatypes.ModuleName
}

// RegisterLegacyAminoCodec registers the amino codec
func (AppModule) RegisterLegacyAminoCodec(*codec.LegacyAmino) {}

// RegisterInterfaces registers a module's interface types and their concrete implementations as proto.Message.
func (AppModule) RegisterInterfaces(registrar codectypes.InterfaceRegistry) {
	twofatypes.RegisterInterfaces(registrar)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the module.
func (AppModule) RegisterGRPCGatewayRoutes(client.Context, *runtime.ServeMux) {
	// No gRPC gateway routes for twofa module
}

// RegisterServices registers a gRPC query service to respond to the module-specific gRPC queries
func (am AppModule) RegisterServices(registrar grpc.ServiceRegistrar) error {
	twofatypes.RegisterMsgServer(registrar, keeper.NewMsgServerImpl(am.keeper))
	return nil
}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() {}

// BeginBlock contains the logic that is automatically triggered at the beginning of each block.
func (am AppModule) BeginBlock(_ context.Context) error {
	return nil
}

// EndBlock contains the logic that is automatically triggered at the end of each block.
func (am AppModule) EndBlock(_ context.Context) error {
	return nil
}

// GetKeeper returns the keeper
func (am AppModule) GetKeeper() keeper.Keeper {
	return am.keeper
}

// ----------------------------------------------------------------------------
// App Wiring Setup
// ----------------------------------------------------------------------------

func init() {
	appmodule.Register(
		&twofatypes.Module{},
		appmodule.Provide(ProvideModule),
	)
}

type ModuleInputs struct {
	depinject.In

	StoreService store.KVStoreService
	Cdc          codec.Codec
	Config       *twofatypes.Module
	Logger       log.Logger

	BadgesKeeper badgeskeeper.Keeper
}

type ModuleOutputs struct {
	depinject.Out

	TwoFAKeeper keeper.Keeper
	Module      appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	k := keeper.NewKeeper(
		in.Cdc,
		in.StoreService,
		in.Logger,
		&in.BadgesKeeper,
	)
	m := NewAppModule(
		in.Cdc,
		k,
	)

	return ModuleOutputs{TwoFAKeeper: k, Module: m}
}

