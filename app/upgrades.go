package app

import (
	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	v7 "github.com/bitbadges/bitbadgeschain/app/upgrades/v7"
)

// RegisterUpgradeHandlers registers all upgrade handlers
func (app *App) RegisterUpgradeHandlers() {
	app.UpgradeKeeper.SetUpgradeHandler(
		v7.UpgradeName,
		v7.CreateUpgradeHandler(
			app.ModuleManager,
			app.Configurator(),
			app.BadgesKeeper,
			app.BankKeeper,
			*app.GovKeeper,
			app.MintKeeper,
			*app.StakingKeeper,
			app.DistrKeeper,
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
	case v7.UpgradeName:
		// Add any store upgrades here
		storeUpgrades = &storetypes.StoreUpgrades{}
	}

	if storeUpgrades != nil {
		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, storeUpgrades))
	}
}
