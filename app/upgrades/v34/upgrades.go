package v34

import (
	"context"

	ibcratelimitkeeper "github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/keeper"
	poolmanagerkeeper "github.com/bitbadges/bitbadgeschain/x/poolmanager"
	tokenizationkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

const (
	UpgradeName = "v34"
)

// CustomUpgradeHandlerLogic runs chain-specific migration work for v34.
//
// v34 is a dependency/infrastructure upgrade — cosmos-sdk v0.54, ibc-go v11
// (which absorbed packet-forward-middleware), and cosmos/evm v0.7. It carries
// no tokenization state-schema change, so unlike v33 there is no tokenization
// keeper migration to run here. The store deletions this upgrade needs
// (x/crisis and x/precisebank, both removed upstream) are declared as
// StoreUpgrades in app/upgrades.go rather than performed here.
//
// This is in a separate function so we can test it locally with a snapshot.
func CustomUpgradeHandlerLogic(ctx context.Context, tokenizationKeeper tokenizationkeeper.Keeper, poolManagerKeeper poolmanagerkeeper.Keeper, rateLimitKeeper ibcratelimitkeeper.Keeper) error {
	return nil
}

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	tokenizationKeeper tokenizationkeeper.Keeper,
	poolManagerKeeper poolmanagerkeeper.Keeper,
	rateLimitKeeper ibcratelimitkeeper.Keeper,
) func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		err := CustomUpgradeHandlerLogic(ctx, tokenizationKeeper, poolManagerKeeper, rateLimitKeeper)
		if err != nil {
			return nil, err
		}

		// Run module migrations
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
