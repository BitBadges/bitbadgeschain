package v24_patch2

import (
	"context"

	sdkmath "cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	tokenizationkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

const UpgradeName = "v24-patch2"

// CorrectNextCollectionId is the correct value for nextCollectionId on testnet
// This was reset to 1 during the v24 upgrade due to the badges->tokenization rename
// not properly transferring the module version in fromVM
const CorrectNextCollectionId = 15

// CorrectNextDynamicStoreId is the correct value for nextDynamicStoreId on testnet
// Set to 0 if you don't need to fix this, otherwise set to the correct value
const CorrectNextDynamicStoreId = 0

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	tokenizationKeeper tokenizationkeeper.Keeper,
) func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)

		// Fix nextCollectionId
		currentNextCollectionId := tokenizationKeeper.GetNextCollectionId(sdkCtx)
		if currentNextCollectionId.LT(sdkmath.NewUint(CorrectNextCollectionId)) {
			tokenizationKeeper.SetNextCollectionId(sdkCtx, sdkmath.NewUint(CorrectNextCollectionId))
			sdkCtx.Logger().Info("Fixed nextCollectionId",
				"old", currentNextCollectionId.String(),
				"new", CorrectNextCollectionId)
		}

		// Fix nextDynamicStoreId if needed
		if CorrectNextDynamicStoreId > 0 {
			currentNextDynamicStoreId := tokenizationKeeper.GetNextDynamicStoreId(sdkCtx)
			if currentNextDynamicStoreId.LT(sdkmath.NewUint(CorrectNextDynamicStoreId)) {
				tokenizationKeeper.SetNextDynamicStoreId(sdkCtx, sdkmath.NewUint(CorrectNextDynamicStoreId))
				sdkCtx.Logger().Info("Fixed nextDynamicStoreId",
					"old", currentNextDynamicStoreId.String(),
					"new", CorrectNextDynamicStoreId)
			}
		}

		sdkCtx.Logger().Info("v24-patch2 upgrade complete",
			"nextCollectionId", tokenizationKeeper.GetNextCollectionId(sdkCtx).String(),
			"nextDynamicStoreId", tokenizationKeeper.GetNextDynamicStoreId(sdkCtx).String())

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
