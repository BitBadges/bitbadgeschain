package app

import (
	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	v21 "github.com/bitbadges/bitbadgeschain/app/upgrades/v21"
	sendmanagertypes "github.com/bitbadges/bitbadgeschain/x/sendmanager/types"
)

// RegisterUpgradeHandlers registers all upgrade handlers
func (app *App) RegisterUpgradeHandlers() {
	app.UpgradeKeeper.SetUpgradeHandler(
		v21.UpgradeName,
		v21.CreateUpgradeHandler(
			app.ModuleManager,
			app.Configurator(),
			app.BadgesKeeper,
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
	case v21.UpgradeName:
		storeUpgrades = &storetypes.StoreUpgrades{
			Added: []string{
				sendmanagertypes.StoreKey,
			},
		}
	}

	if storeUpgrades != nil {
		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, storeUpgrades))
	}
}
