package app

import (
	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	v24 "github.com/bitbadges/bitbadgeschain/app/upgrades/v24"
	v24_patch "github.com/bitbadges/bitbadgeschain/app/upgrades/v24_patch"
	v24_patch2 "github.com/bitbadges/bitbadgeschain/app/upgrades/v24_patch2"
)

// RegisterUpgradeHandlers registers all upgrade handlers
func (app *App) RegisterUpgradeHandlers() {
	app.UpgradeKeeper.SetUpgradeHandler(
		v24.UpgradeName,
		v24.CreateUpgradeHandler(
			app.ModuleManager,
			app.Configurator(),
			app.TokenizationKeeper,
			app.PoolManagerKeeper,
			app.IBCRateLimitKeeper,
			app.EVMKeeper,
			app.BankKeeper,
			app.FeeMarketKeeper,
			app.ERC20Keeper,
		),
	)

	// v24-patch: Fix EVM configuration on testnets that already ran v24
	// This fixes the aatom -> ubadge denom issue and enables all precompiles
	app.UpgradeKeeper.SetUpgradeHandler(
		v24_patch.UpgradeName,
		v24_patch.CreateUpgradeHandler(
			app.ModuleManager,
			app.Configurator(),
			app.EVMKeeper,
			app.BankKeeper,
		),
	)

	// v24-patch2: Fix nextCollectionId on testnets that already ran v24
	// The badges->tokenization rename caused InitGenesis to be called, resetting nextCollectionId to 1
	app.UpgradeKeeper.SetUpgradeHandler(
		v24_patch2.UpgradeName,
		v24_patch2.CreateUpgradeHandler(
			app.ModuleManager,
			app.Configurator(),
			app.TokenizationKeeper,
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
	case v24.UpgradeName:
		storeUpgrades = &storetypes.StoreUpgrades{
			Renamed: []storetypes.StoreRename{
				{
					OldKey: "badges",
					NewKey: "tokenization",
				},
			},
			Deleted: []string{
				"wasm-store",  // Remove WASM store
				"wasmx-store", // Remove WASMX store
			},
			Added: []string{
				"erc20",       // Add ERC20 store (new keeper)
				"evm",         // Add EVM store (new keeper)
				"feemarket",   // Add FeeMarket store (new keeper)
				"precisebank", // Add PreciseBank store (new keeper)
			},
		}
	case v24_patch.UpgradeName:
		// No store changes needed - this is a patch to fix EVM configuration
		// All stores were already added in v24
		storeUpgrades = nil
	case v24_patch2.UpgradeName:
		// No store changes needed - this is a patch to fix nextCollectionId
		storeUpgrades = nil
	}

	if storeUpgrades != nil {
		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, storeUpgrades))
	}
}
