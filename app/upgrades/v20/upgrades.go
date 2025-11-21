package v20

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	badgeskeeper "github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	poolmanagerkeeper "github.com/bitbadges/bitbadgeschain/x/poolmanager"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

const (
	UpgradeName = "v20"
)

// This is in a separate function so we can test it locally with a snapshot
func CustomUpgradeHandlerLogic(ctx context.Context, badgesKeeper badgeskeeper.Keeper, poolManagerKeeper poolmanagerkeeper.Keeper) error {
	// Run badges migrations
	if err := badgesKeeper.MigrateBadgesKeeper(sdk.UnwrapSDKContext(ctx)); err != nil {
		return err
	}

	// Update poolmanager default taker fee to 0.1% (0.001)
	if err := badgeskeeper.MigratePoolManagerTakerFee(sdk.UnwrapSDKContext(ctx), poolManagerKeeper); err != nil {
		return err
	}

	return nil
}

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	badgesKeeper badgeskeeper.Keeper,
	poolManagerKeeper poolmanagerkeeper.Keeper,
) func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		err := CustomUpgradeHandlerLogic(ctx, badgesKeeper, poolManagerKeeper)
		if err != nil {
			return nil, err
		}

		// Run module migrations
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
