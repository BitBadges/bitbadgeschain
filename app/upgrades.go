package app

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	v34 "github.com/bitbadges/bitbadgeschain/app/upgrades/v34"
	v35 "github.com/bitbadges/bitbadgeschain/app/upgrades/v35"
)

// RegisterUpgradeHandlers registers all upgrade handlers
func (app *App) RegisterUpgradeHandlers() {
	app.UpgradeKeeper.SetUpgradeHandler(
		v34.UpgradeName,
		v34.CreateUpgradeHandler(
			app.ModuleManager,
			app.Configurator(),
			app.AccountKeeper,
		),
	)
	app.UpgradeKeeper.SetUpgradeHandler(
		v35.UpgradeName,
		v35.CreateUpgradeHandler(
			app.ModuleManager,
			app.Configurator(),
			v35.Keepers{
				ConsensusParams: app.ConsensusParamsKeeper,
				FeeMarket:       app.FeeMarketKeeper,
				IBCRateLimit:    app.IBCRateLimitKeeper,
				Tokenization:    app.TokenizationKeeper,
			},
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
	case v34.UpgradeName:
		storeUpgrades = &storetypes.StoreUpgrades{
			Renamed: []storetypes.StoreRename{},
			// v34 drops two module stores:
			//   crisis - x/crisis left the SDK in v0.54
			//   group  - x/group moved out of the SDK in v0.54 (enterprise)
			//
			// precisebank is deliberately NOT deleted. cosmos/evm v0.7 moved the
			// module to contrib rather than deleting it, and this chain still
			// needs it (see app.go) because BADGE is 9-decimal. Its v33 store
			// carries the fractional balances and must survive the upgrade.
			//
			// Both are registered modules in v33 (see the v33 genesis, which
			// carries "crisis" and "group" app_state) and are gone in v34, so
			// their stores would otherwise be orphaned on disk forever. The
			// upgrade completes either way — the SDK tolerates an unmounted
			// store — but leaving them behind keeps dead state around and is
			// inconsistent with how crisis is handled.
			Deleted: []string{"crisis", "group"},
			Added:   []string{},
		}
	}

	if storeUpgrades != nil {
		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, storeUpgrades))
	}
}
