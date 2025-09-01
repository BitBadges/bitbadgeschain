package app

import (
	"cosmossdk.io/core/appmodule"
	storetypes "cosmossdk.io/store/types"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	// this line is used by starport scaffolding # ibc/app/import

	"github.com/bitbadges/bitbadgeschain/x/gamm"
	gammkeeper "github.com/bitbadges/bitbadgeschain/x/gamm/keeper"
	gammtypes "github.com/bitbadges/bitbadgeschain/x/gamm/types"
	"github.com/bitbadges/bitbadgeschain/x/poolmanager"
	poolmanagermodule "github.com/bitbadges/bitbadgeschain/x/poolmanager/module"
	poolmanagertypes "github.com/bitbadges/bitbadgeschain/x/poolmanager/types"
)

// registerIBCModules register IBC keepers and non dependency inject modules.
func (app *App) registerGammModules(appOpts servertypes.AppOptions) error {
	// set up non depinject support modules store keys
	if err := app.RegisterStores(
		storetypes.NewKVStoreKey(gammtypes.StoreKey),
		storetypes.NewKVStoreKey(poolmanagertypes.StoreKey),
	); err != nil {
		return err
	}

	// register thee gamm params
	app.ParamsKeeper.Subspace(poolmanagertypes.ModuleName).WithKeyTable(poolmanagertypes.ParamKeyTable())
	app.ParamsKeeper.Subspace(gammtypes.ModuleName).WithKeyTable(gammtypes.ParamKeyTable())

	app.GammKeeper = gammkeeper.NewKeeper(
		app.appCodec,
		app.GetKey(gammtypes.StoreKey),
		app.GetSubspace(gammtypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.DistrKeeper,
		app.BadgesKeeper,
	)

	app.PoolManagerKeeper = *poolmanager.NewKeeper(
		app.GetKey(poolmanagertypes.StoreKey),
		app.GetSubspace(poolmanagertypes.ModuleName),
		app.GammKeeper,
		app.BankKeeper,
		app.AccountKeeper,
		app.DistrKeeper,
		app.StakingKeeper,
	)

	app.GammKeeper.SetPoolManager(&app.PoolManagerKeeper)

	// register IBC modules
	if err := app.RegisterModules(
		gamm.NewAppModule(app.appCodec, app.GammKeeper, app.AccountKeeper, app.BankKeeper),
		poolmanagermodule.NewAppModule(app.PoolManagerKeeper, app.GammKeeper),
	); err != nil {
		return err
	}

	return nil
}

// RegisterIBC Since the IBC modules don't support dependency injection,
// we need to manually register the modules on the client side.
// This needs to be removed after IBC supports App Wiring.
func RegisterGamm(registry cdctypes.InterfaceRegistry) map[string]appmodule.AppModule {
	modules := map[string]appmodule.AppModule{
		gammtypes.ModuleName:        gamm.AppModule{},
		poolmanagertypes.ModuleName: poolmanagermodule.AppModule{},
	}

	for name, m := range modules {
		module.CoreAppModuleBasicAdaptor(name, m).RegisterInterfaces(registry)
	}

	return modules
}
