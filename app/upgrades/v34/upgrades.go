package v34

import (
	"context"

	ibcratelimitkeeper "github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/keeper"
	poolmanagerkeeper "github.com/bitbadges/bitbadgeschain/x/poolmanager"
	tokenizationkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

const (
	UpgradeName = "v34"
)

// CustomUpgradeHandlerLogic runs chain-specific migration work for v34.
//
// v34 is primarily a dependency/infrastructure upgrade — cosmos-sdk v0.54,
// ibc-go v11 (which absorbed packet-forward-middleware), and cosmos/evm v0.7.
// It carries no tokenization state-schema change, so unlike v33 there is no
// tokenization keeper migration. Store deletions (x/crisis, x/group) are
// declared as StoreUpgrades in app/upgrades.go rather than performed here.
//
// It does carry one state migration: accounts still holding this chain's
// pre-cosmos/evm `/ethereum.PubKey` are rewritten to cosmos/evm's key type.
// Those accounts cannot transact today — see pubkeys.go and BB-12.
//
// This is in a separate function so we can test it locally with a snapshot.
func CustomUpgradeHandlerLogic(ctx context.Context, accountKeeper authkeeper.AccountKeeper, tokenizationKeeper tokenizationkeeper.Keeper, poolManagerKeeper poolmanagerkeeper.Keeper, rateLimitKeeper ibcratelimitkeeper.Keeper) error {
	if _, err := MigrateLegacyEthereumPubKeys(sdk.UnwrapSDKContext(ctx), accountKeeper); err != nil {
		return err
	}
	return nil
}

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	accountKeeper authkeeper.AccountKeeper,
	tokenizationKeeper tokenizationkeeper.Keeper,
	poolManagerKeeper poolmanagerkeeper.Keeper,
	rateLimitKeeper ibcratelimitkeeper.Keeper,
) func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		err := CustomUpgradeHandlerLogic(ctx, accountKeeper, tokenizationKeeper, poolManagerKeeper, rateLimitKeeper)
		if err != nil {
			return nil, err
		}

		// Run module migrations
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
