package app

import (
	storetypes "cosmossdk.io/store/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	v30 "github.com/bitbadges/bitbadgeschain/app/upgrades/v30"
)

// RegisterUpgradeHandlers registers all upgrade handlers
func (app *App) RegisterUpgradeHandlers() {
	app.UpgradeKeeper.SetUpgradeHandler(
		v30.UpgradeName,
		v30.CreateUpgradeHandler(
			app.ModuleManager,
			app.Configurator(),
			*app.TokenizationKeeper,
			app.PoolManagerKeeper,
			app.IBCRateLimitKeeper,
		),
	)

	// When a planned upgrade height is reached, the old binary will panic
	// writing on disk the height and name of the upgrade that triggered it
	// This will read that value, and execute the preparations for the upgrade.
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}

	var storeUpgrades *storetypes.StoreUpgrades

	switch upgradeInfo.Name {
	case v30.UpgradeName:
		storeUpgrades = &storetypes.StoreUpgrades{
			Renamed: []storetypes.StoreRename{},
			// v30: remove x/anchor, x/maps, x/group, x/crisis
			// (group/crisis removed by SDK v0.54 migration; anchor/maps deprecated)
			Deleted: []string{"anchor", "maps", "group", "crisis"},
			Added:   []string{},
		}
	}

	if storeUpgrades != nil {
		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, storeUpgrades))
	}
}
