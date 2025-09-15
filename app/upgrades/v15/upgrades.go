package v15

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
)

const (
	UpgradeName = "v15"
)

// This is in a separate function so we can test it locally with a snapshot
func CustomUpgradeHandlerLogic(ctx context.Context, badgesKeeper keeper.Keeper, mintKeeper mintkeeper.Keeper, slashingKeeper slashingkeeper.Keeper) error {
	// Run migrations
	if err := badgesKeeper.MigrateBadgesKeeper(sdk.UnwrapSDKContext(ctx)); err != nil {
		return err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Update mint module parameters to reflect new block time
	mintParams, err := mintKeeper.Params.Get(sdkCtx)
	if err != nil {
		return err
	}
	mintParams.BlocksPerYear = 15811200
	mintKeeper.Params.Set(sdkCtx, mintParams)

	// Update slashing module parameters
	slashingParams, err := slashingKeeper.GetParams(sdkCtx)
	if err != nil {
		return err
	}
	slashingParams.SignedBlocksWindow = 30000
	slashingKeeper.SetParams(sdkCtx, slashingParams)

	return nil
}

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	badgesKeeper keeper.Keeper,
	mintKeeper mintkeeper.Keeper,
	slashingKeeper slashingkeeper.Keeper,
) func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		err := CustomUpgradeHandlerLogic(ctx, badgesKeeper, mintKeeper, slashingKeeper)
		if err != nil {
			return nil, err
		}

		// Run module migrations
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
