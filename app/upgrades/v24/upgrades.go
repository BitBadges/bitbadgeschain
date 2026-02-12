package v24

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	ibcratelimitkeeper "github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/keeper"
	poolmanagerkeeper "github.com/bitbadges/bitbadgeschain/x/poolmanager"
	tokenizationkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	tokenizationprecompile "github.com/bitbadges/bitbadgeschain/x/tokenization/precompile"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	evmkeeper "github.com/cosmos/evm/x/vm/keeper"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	// Uncomment when configuring rate limits:
)

const (
	UpgradeName = "v24"
)

// This is in a separate function so we can test it locally with a snapshot
func CustomUpgradeHandlerLogic(
	ctx context.Context,
	tokenizationKeeper tokenizationkeeper.Keeper,
	poolManagerKeeper poolmanagerkeeper.Keeper,
	rateLimitKeeper ibcratelimitkeeper.Keeper,
	evmKeeper *evmkeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.Logger().Debug("Running v23 upgrade handler")

	// Run tokenization migrations
	if err := tokenizationKeeper.MigrateTokenizationKeeper(sdkCtx); err != nil {
		return err
	}

	// EVM upgrade logic: Set denom metadata if not already set
	// This ensures the EVM module can find the denom metadata
	denomMetadata := banktypes.Metadata{
		Description: "The native token of BitBadges Chain",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    "ubadge",
				Exponent: 0,
				Aliases:  nil,
			},
			{
				Denom:    "badge",
				Exponent: 9,
				Aliases:  nil,
			},
		},
		Base:    "ubadge",
		Display: "badge",
		Name:    "Badge",
		Symbol:  "BADGE",
	}
	bankKeeper.SetDenomMetaData(ctx, denomMetadata)

	// (Required for NON-18 denom chains *only)
	// Update EVM params to set EvmDenom and add Extended denom options
	// InitEvmCoinInfo uses params.EvmDenom to look up denom metadata, so we must set it to "ubadge"
	// We use "abadge" as the extended denom (18 decimals via precisebank)
	evmParams := evmKeeper.GetParams(sdkCtx)
	evmParams.EvmDenom = "ubadge" // Set EvmDenom to "ubadge" so InitEvmCoinInfo can find the metadata
	evmParams.ExtendedDenomOptions = &evmtypes.ExtendedDenomOptions{ExtendedDenom: "abadge"}
	if err := evmKeeper.SetParams(sdkCtx, evmParams); err != nil {
		return err
	}

	// Initialize EvmCoinInfo in the module store. Chains bootstrapped before v0.5.0
	// binaries never stored this information (it lived only in process globals),
	// so migrating nodes would otherwise see an empty EvmCoinInfo on upgrade.
	if err := evmKeeper.InitEvmCoinInfo(sdkCtx); err != nil {
		return err
	}

	// Enable tokenization precompile in active_static_precompiles
	// The precompile is registered during app initialization, but must be enabled
	// via EnableStaticPrecompiles to be callable. This ensures the precompile
	// is available for existing chains that upgrade to this version.
	tokenizationPrecompileAddr := common.HexToAddress(tokenizationprecompile.TokenizationPrecompileAddress)
	if err := evmKeeper.EnableStaticPrecompiles(sdkCtx, tokenizationPrecompileAddr); err != nil {
		// Log error but don't fail the upgrade if precompile is already enabled
		// This allows the migration to be idempotent
		sdkCtx.Logger().Info("Tokenization precompile enable attempt", "error", err, "address", tokenizationPrecompileAddr.Hex())
	}

	return nil
}

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	tokenizationKeeper tokenizationkeeper.Keeper,
	poolManagerKeeper poolmanagerkeeper.Keeper,
	rateLimitKeeper ibcratelimitkeeper.Keeper,
	evmKeeper *evmkeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
) func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		err := CustomUpgradeHandlerLogic(ctx, tokenizationKeeper, poolManagerKeeper, rateLimitKeeper, evmKeeper, bankKeeper)
		if err != nil {
			return nil, err
		}

		// Run module migrations
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
