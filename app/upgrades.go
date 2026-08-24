package app

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	v34 "github.com/bitbadges/bitbadgeschain/app/upgrades/v34"
	v35 "github.com/bitbadges/bitbadgeschain/app/upgrades/v35"
)

// V35Keepers is the keeper set the v35 decimals migration needs.
//
// Built here rather than inline so there is exactly one place that wires it.
// An earlier version constructed this inline while the tests built their own
// equivalent struct; when new keepers were added to only one of them the tests
// kept passing against a correctly-populated struct while the chain ran the
// migration with nil keepers and panicked at the upgrade height.
func (app *App) V35Keepers() v35.Keepers {
	return v35.Keepers{
		Account:      app.AccountKeeper,
		Authz:        app.AuthzKeeper,
		Bank:         app.BankKeeper.(bankkeeper.BaseKeeper),
		Distribution: app.DistrKeeper,
		ERC20:        app.ERC20Keeper,
		EVM:          app.EVMKeeper,
		FeeGrant:     app.FeeGrantKeeper,
		FeeMarket:    app.FeeMarketKeeper,
		Gamm:         app.GammKeeper,
		Gov:          app.GovKeeper,
		Mint:         app.MintKeeper,
		PoolManager:  app.PoolManagerKeeper,
		RateLimit:    app.IBCRateLimitKeeper,
		Staking:      app.StakingKeeper,
		Tokenization: *app.TokenizationKeeper,
		Transfer:     app.TransferKeeper,
	}
}

// RegisterUpgradeHandlers registers all upgrade handlers
func (app *App) RegisterUpgradeHandlers() {
	app.UpgradeKeeper.SetUpgradeHandler(
		v34.UpgradeName,
		v34.CreateUpgradeHandler(
			app.ModuleManager,
			app.Configurator(),
			app.AccountKeeper,
			*app.TokenizationKeeper,
			app.PoolManagerKeeper,
			app.IBCRateLimitKeeper,
		),
	)

	app.UpgradeKeeper.SetUpgradeHandler(
		v35.UpgradeName,
		v35.CreateUpgradeHandler(
			app.ModuleManager,
			app.Configurator(),
			app.V35Keepers(),
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
	case v35.UpgradeName:
		storeUpgrades = &storetypes.StoreUpgrades{
			Renamed: []storetypes.StoreRename{},
			// The 18-decimal migration makes base == extended denom, so
			// x/precisebank has nothing left to bridge and is unwired.
			Deleted: []string{"precisebank"},
			Added:   []string{},
		}
	case v34.UpgradeName:
		storeUpgrades = &storetypes.StoreUpgrades{
			Renamed: []storetypes.StoreRename{},
			// v34 drops three module stores:
			//   crisis - x/crisis left the SDK in v0.54
			//   group  - x/group moved out of the SDK in v0.54 (enterprise)
			//
			// precisebank is deliberately NOT deleted. cosmos/evm v0.7 moved the
			// module to contrib rather than deleting it, and this chain still
			// needs it (see app.go) because BADGE is 9-decimal. Its v33 store
			// carries the fractional balances and must survive the upgrade.
			//
			// All three are registered modules in v33 (see the v33 genesis, which
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
