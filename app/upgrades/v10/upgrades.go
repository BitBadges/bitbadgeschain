package v10

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

const (
	UpgradeName = "v10"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	badgesKeeper keeper.Keeper,
	wasmKeeper wasmkeeper.Keeper,
) func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// Run migrations
		if err := badgesKeeper.MigrateBadgesKeeper(sdk.UnwrapSDKContext(ctx)); err != nil {
			return nil, err
		}

		// Migrate WASM permissions from Everybody to Nobody
		if err := migrateWasmPermissions(ctx, wasmKeeper); err != nil {
			return nil, err
		}

		// Run module migrations
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}

// migrateWasmPermissions migrates WASM permissions from Everybody to Nobody
func migrateWasmPermissions(ctx context.Context, wasmKeeper wasmkeeper.Keeper) error {
	// Try to get current WASM parameters, but handle the case where they might not exist
	var params wasmtypes.Params

	// Use recover to handle potential panic when params don't exist
	func() {
		defer func() {
			if r := recover(); r != nil {
				// If panic occurs, use default params
				params = wasmtypes.DefaultParams()
			}
		}()
		params = wasmKeeper.GetParams(ctx)
	}()

	// Update the permission
	params.InstantiateDefaultPermission = wasmtypes.AccessTypeNobody

	// Set the updated parameters
	if err := wasmKeeper.SetParams(ctx, params); err != nil {
		return err
	}

	return nil
}
