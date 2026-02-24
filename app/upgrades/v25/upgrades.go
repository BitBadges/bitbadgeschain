package v25

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

const (
	UpgradeName = "v25"
)

// CreateUpgradeHandler creates the v25 upgrade handler.
// Module migrations (e.g. tokenization 24->25 MigrateCollectionStats) are run via RunMigrations.
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
) func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
