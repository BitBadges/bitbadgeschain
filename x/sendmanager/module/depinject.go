package sendmanager

import (
	"fmt"

	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"
	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/keeper"
	"github.com/bitbadges/bitbadgeschain/x/sendmanager/types"
)

var _ depinject.OnePerModuleType = AppModule{}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (AppModule) IsOnePerModuleType() {}

func init() {
	appconfig.Register(
		&types.Module{},
		appconfig.Provide(ProvideModule),
	)
}

type ModuleInputs struct {
	depinject.In

	Config       *types.Module
	StoreService store.KVStoreService
	Cdc          codec.Codec
	AddressCodec address.Codec

	AuthKeeper         types.AuthKeeper
	BankKeeper         types.BankKeeper
	DistributionKeeper types.DistributionKeeper
	// TokenizationKeeper is not needed at initialization - router registration happens later
	// TokenizationKeeper       tokenizationkeeper.Keeper `optional:"true"`
}

type ModuleOutputs struct {
	depinject.Out

	SendmanagerKeeper keeper.Keeper
	Module            appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	// default to governance authority if not provided
	authority := authtypes.NewModuleAddress("gov")
	if in.Config != nil && in.Config.Authority != "" {
		authority = authtypes.NewModuleAddressOrBech32Address(in.Config.Authority)
	}
	authorityBytes, err := in.AddressCodec.StringToBytes(authority.String())
	if err != nil {
		panic(fmt.Sprintf("failed to convert authority to bytes: %v", err))
	}

	k := keeper.NewKeeper(
		in.Cdc,
		in.StoreService,
		in.AddressCodec,
		authorityBytes,
		in.BankKeeper,
		in.DistributionKeeper,
	)

	// Note: Tokenization router registration is deferred to app.go after both keepers are created
	// This avoids a circular dependency (sendmanager needs tokenization for router, tokenization needs sendmanager)

	m := NewAppModule(in.Cdc, k, in.AuthKeeper, in.BankKeeper)

	return ModuleOutputs{SendmanagerKeeper: k, Module: m}
}
