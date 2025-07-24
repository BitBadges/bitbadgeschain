package v11

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
	UpgradeName = "v11"
)

// This is in a separate function so we can test it locally with a snapshot
func CustomUpgradeHandlerLogic(ctx context.Context, badgesKeeper keeper.Keeper, wasmKeeper wasmkeeper.Keeper) error {
	// Run migrations
	if err := badgesKeeper.MigrateBadgesKeeper(sdk.UnwrapSDKContext(ctx)); err != nil {
		return err
	}

	// Migrate WASM permissions from Everybody to Nobody
	if err := migrateWasmPermissions(ctx, wasmKeeper); err != nil {
		return err
	}

	return nil
}

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	badgesKeeper keeper.Keeper,
	wasmKeeper wasmkeeper.Keeper,
) func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {

		// Run custom upgrade logic
		if err := CustomUpgradeHandlerLogic(ctx, badgesKeeper, wasmKeeper); err != nil {
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

	chainId := sdk.UnwrapSDKContext(ctx).ChainID()
	if chainId == "bitbadges-1" {
		params.InstantiateDefaultPermission = wasmtypes.AccessTypeNobody
	} else {
		params.InstantiateDefaultPermission = wasmtypes.AccessTypeEverybody
	}

	// Set the updated parameters
	if err := wasmKeeper.SetParams(ctx, params); err != nil {
		return err
	}

	return nil
}
