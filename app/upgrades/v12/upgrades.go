package v12

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

const (
	UpgradeName = "v12"
)

// This is in a separate function so we can test it locally with a snapshot
func CustomUpgradeHandlerLogic(ctx context.Context, badgesKeeper keeper.Keeper, wasmKeeper wasmkeeper.Keeper) error {
	// Run migrations
	if err := badgesKeeper.MigrateBadgesKeeper(sdk.UnwrapSDKContext(ctx)); err != nil {
		return err
	}

	return nil
}

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	badgesKeeper keeper.Keeper,
) func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// Run module migrations
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
