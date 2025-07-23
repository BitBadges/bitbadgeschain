package app

import (
	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	v10 "github.com/bitbadges/bitbadgeschain/app/upgrades/v10"
)

// RegisterUpgradeHandlers registers all upgrade handlers
func (app *App) RegisterUpgradeHandlers() {
	app.UpgradeKeeper.SetUpgradeHandler(
		v10.UpgradeName,
		v10.CreateUpgradeHandler(
			app.ModuleManager,
			app.Configurator(),
			app.BadgesKeeper,
			// app.WasmKeeper,
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
	case v10.UpgradeName:
		// Add any store upgrades here
		storeUpgrades = &storetypes.StoreUpgrades{
			Added: []string{},
		}
	}

	if storeUpgrades != nil {
		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, storeUpgrades))
	}
}
