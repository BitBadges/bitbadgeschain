package v23

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	badgeskeeper "github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	ibcratelimitkeeper "github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/keeper"
	poolmanagerkeeper "github.com/bitbadges/bitbadgeschain/x/poolmanager"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	// Uncomment when configuring rate limits:
)

const (
	UpgradeName = "v23"
)

// This is in a separate function so we can test it locally with a snapshot
func CustomUpgradeHandlerLogic(ctx context.Context, badgesKeeper badgeskeeper.Keeper, poolManagerKeeper poolmanagerkeeper.Keeper, rateLimitKeeper ibcratelimitkeeper.Keeper) error {
	// Run badges migrations
	if err := badgesKeeper.MigrateBadgesKeeper(sdk.UnwrapSDKContext(ctx)); err != nil {
		return err
	}

	return nil
}

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	badgesKeeper badgeskeeper.Keeper,
	poolManagerKeeper poolmanagerkeeper.Keeper,
	rateLimitKeeper ibcratelimitkeeper.Keeper,
) func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		err := CustomUpgradeHandlerLogic(ctx, badgesKeeper, poolManagerKeeper, rateLimitKeeper)
		if err != nil {
			return nil, err
		}

		// Run module migrations
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
