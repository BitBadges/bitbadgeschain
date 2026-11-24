package app

import (
	"cosmossdk.io/core/appmodule"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
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
	// set up non depinject support modules store keys.
	// Keep the legacy "params" KV store mounted through v36 so the upgrade
	// handler can migrate poolmanager params before the store is retired.
	if err := app.RegisterStores(
		storetypes.NewKVStoreKey(gammtypes.StoreKey),
		storetypes.NewKVStoreKey(poolmanagertypes.StoreKey),
		storetypes.NewKVStoreKey("params"),
	); err != nil {
		return err
	}

	app.GammKeeper = gammkeeper.NewKeeper(
		app.appCodec,
		app.GetKey(gammtypes.StoreKey),
		app.AccountKeeper,
		app.BankKeeper,
		app.DistrKeeper,
		app.TokenizationKeeper,
		&app.SendmanagerKeeper,
		app.TransferKeeper,
		app.TransferICS4Wrapper,
		app.IBCKeeper.ChannelKeeper,
	)

	app.PoolManagerKeeper = *poolmanager.NewKeeper(
		app.GetKey(poolmanagertypes.StoreKey),
		&app.GammKeeper,
		app.BankKeeper,
		app.AccountKeeper,
		app.DistrKeeper,
		app.StakingKeeper,
	)

	app.GammKeeper.SetPoolManager(&app.PoolManagerKeeper)

	// register IBC modules
	if err := app.RegisterModules(
		gamm.NewAppModule(app.appCodec, app.GammKeeper, app.AccountKeeper, app.BankKeeper),
		poolmanagermodule.NewAppModule(&app.PoolManagerKeeper, app.GammKeeper),
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
