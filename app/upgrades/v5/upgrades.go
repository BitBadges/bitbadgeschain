package v5

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

const (
	UpgradeName = "v5"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	badgesKeeper keeper.Keeper,
) func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// Run migrations
		if err := badgesKeeper.MigrateBadgesKeeper(sdk.UnwrapSDKContext(ctx)); err != nil {
			return nil, err
		}

		// Run module migrations
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
