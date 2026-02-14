package v24

import (
	"context"
	"strconv"

	"github.com/ethereum/go-ethereum/common"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	appparams "github.com/bitbadges/bitbadgeschain/app/params"
	gammprecompile "github.com/bitbadges/bitbadgeschain/x/gamm/precompile"
	ibcratelimitkeeper "github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/keeper"
	poolmanagerkeeper "github.com/bitbadges/bitbadgeschain/x/poolmanager"
	sendmanagerprecompile "github.com/bitbadges/bitbadgeschain/x/sendmanager/precompile"
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
	sdkCtx.Logger().Debug("Running v24 upgrade handler")

	// CRITICAL: Set EVM params and coin info FIRST before any migrations run
	// The EVM module's migration may call InitEvmCoinInfo, which uses params.EvmDenom to look up metadata.
	// By setting both params and coin info first, InitEvmCoinInfo will find the coin info already set
	// and skip the metadata lookup entirely.

	// 1. Set EVM params first
	evmParams := evmKeeper.GetParams(sdkCtx)
	oldEvmDenom := evmParams.EvmDenom
	evmParams.EvmDenom = "ubadge"
	evmParams.ExtendedDenomOptions = &evmtypes.ExtendedDenomOptions{ExtendedDenom: "abadge"}
	if err := evmKeeper.SetParams(sdkCtx, evmParams); err != nil {
		return err
	}
	sdkCtx.Logger().Info("Updated EVM params before migrations", "oldEvmDenom", oldEvmDenom, "newEvmDenom", "ubadge")

	// 2. Set denom metadata for "ubadge" - must exist before EVM module migrations run
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

	// TEMPORARY: Create minimal metadata for "aatom" to prevent InitEvmCoinInfo from panicking
	// during EVM module InitGenesis. The EVM module's InitGenesis uses default genesis state
	// which has "aatom" as the denom. We'll fix params and coin info after migrations complete.
	// Must use 18 decimals (standard EVM) to prevent "unsupported decimals: 0" error.
	if _, found := bankKeeper.GetDenomMetaData(ctx, "aatom"); !found {
		sdkCtx.Logger().Info("Creating temporary metadata for 'aatom' to prevent InitEvmCoinInfo panic during migrations")
		aatomMetadata := banktypes.Metadata{
			Description: "Temporary metadata (will be replaced with ubadge after migrations)",
			DenomUnits: []*banktypes.DenomUnit{
				{Denom: "aatom", Exponent: 0},
				{Denom: "atom", Exponent: 18}, // 18 decimals for EVM compatibility
			},
			Base:    "aatom",
			Display: "atom",
			Name:    "atom",
			Symbol:  "ATOM",
		}
		bankKeeper.SetDenomMetaData(ctx, aatomMetadata)
	}

	// 3. Set EvmCoinInfo directly BEFORE migrations run
	// This ensures that if InitEvmCoinInfo is called during migrations, it will find
	// the coin info already set and skip the metadata lookup (which would fail with "aatom")
	if err := evmKeeper.SetEvmCoinInfo(sdkCtx, evmtypes.EvmCoinInfo{
		Denom:         "ubadge",
		ExtendedDenom: "abadge",
		DisplayDenom:  "BADGE",
		Decimals:      9,
	}); err != nil {
		sdkCtx.Logger().Info("EVM coin info already set or error setting", "error", err)
		// Don't fail the upgrade if it's already set
	}

	// 4. Ensure EVM chain ID is set correctly (from build-time flag) BEFORE migrations
	// This is critical for EVM transaction signing and replay protection
	// Must be done before migrations so EVM module InitGenesis uses correct chain ID
	if err := ensureEVMChainID(sdkCtx); err != nil {
		return err
	}

	// Run tokenization migrations
	if err := tokenizationKeeper.MigrateTokenizationKeeper(sdkCtx); err != nil {
		return err
	}

	// Enable all precompiles (both default Cosmos and custom BitBadges)
	// This ensures all registered precompiles are enabled for existing chains that upgrade to this version
	// Get all precompile addresses (default Cosmos + custom BitBadges)
	allPrecompileAddresses := []common.Address{}

	// Add default Cosmos precompiles (0x0800-0x0806)
	// Addresses from cosmos/evm/x/vm/types/precompiles.go
	allPrecompileAddresses = append(allPrecompileAddresses, []common.Address{
		common.HexToAddress("0x0000000000000000000000000000000000000800"), // Staking
		common.HexToAddress("0x0000000000000000000000000000000000000801"), // Distribution
		common.HexToAddress("0x0000000000000000000000000000000000000802"), // ICS20 (IBC)
		common.HexToAddress("0x0000000000000000000000000000000000000803"), // Vesting
		common.HexToAddress("0x0000000000000000000000000000000000000804"), // Bank
		common.HexToAddress("0x0000000000000000000000000000000000000805"), // Governance
		common.HexToAddress("0x0000000000000000000000000000000000000806"), // Slashing
	}...)

	// Add custom BitBadges precompiles (0x1001+)
	allPrecompileAddresses = append(allPrecompileAddresses, []common.Address{
		common.HexToAddress(tokenizationprecompile.TokenizationPrecompileAddress), // 0x1001 - Tokenization
		common.HexToAddress(gammprecompile.GammPrecompileAddress),                 // 0x1002 - Gamm
		common.HexToAddress(sendmanagerprecompile.SendManagerPrecompileAddress),   // 0x1003 - SendManager
		// Next available address: 0x1004
	}...)

	for _, addr := range allPrecompileAddresses {
		if err := evmKeeper.EnableStaticPrecompiles(sdkCtx, addr); err != nil {
			// Log error but don't fail the upgrade if precompile is already enabled
			// This allows the migration to be idempotent
			sdkCtx.Logger().Info("Precompile enable attempt", "error", err, "address", addr.Hex())
		}
	}

	return nil
}

// ensureEVMChainID ensures the EVM chain ID matches the build-time value.
// This is called during upgrade migrations to verify/set the correct EVM chain ID.
func ensureEVMChainID(ctx sdk.Context) error {
	// Get the expected EVM chain ID from build-time flag (defaults to 90123 for local dev)
	expectedEVMChainIDStr := appparams.GetEVMChainID()
	expectedEVMChainID, err := strconv.ParseUint(expectedEVMChainIDStr, 10, 64)
	if err != nil {
		ctx.Logger().Error("Failed to parse expected EVM chain ID", "error", err, "chainIDStr", expectedEVMChainIDStr)
		return err
	}

	// Get current chain config
	chainConfig := evmtypes.GetChainConfig()
	if chainConfig == nil {
		ctx.Logger().Info("EVM chain config not initialized")
		return nil
	}

	currentEVMChainID := chainConfig.GetChainId()

	// Log chain ID info
	ctx.Logger().Info("EVM chain ID validation during upgrade",
		"expectedEVMChainID", expectedEVMChainID,
		"currentEVMChainID", currentEVMChainID,
	)

	// Check if update is needed
	if currentEVMChainID != expectedEVMChainID {
		ctx.Logger().Info("EVM chain ID mismatch detected during upgrade",
			"current", currentEVMChainID,
			"expected", expectedEVMChainID,
		)

		// Create new chain config with correct chain ID
		newChainConfig := evmtypes.DefaultChainConfig(expectedEVMChainID)
		if err := evmtypes.SetChainConfig(newChainConfig); err != nil {
			// SetChainConfig may fail if already set - this is expected
			ctx.Logger().Info("Could not update EVM chain config (may already be set)",
				"error", err,
				"currentChainID", currentEVMChainID,
				"expectedChainID", expectedEVMChainID,
			)
			// Don't return error - the chain config was already set during keeper initialization
		} else {
			ctx.Logger().Info("Updated EVM chain ID during upgrade",
				"from", currentEVMChainID,
				"to", expectedEVMChainID,
			)
		}
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
		sdkCtx := sdk.UnwrapSDKContext(ctx)

		err := CustomUpgradeHandlerLogic(ctx, tokenizationKeeper, poolManagerKeeper, rateLimitKeeper, evmKeeper, bankKeeper)
		if err != nil {
			return nil, err
		}

		// Run module migrations
		// Note: For new modules (like EVM), RunMigrations calls InitGenesis with default genesis state
		// The EVM module's InitGenesis will set params from genesis (defaults to "aatom") and call InitEvmCoinInfo
		// We need to fix params and coin info AFTER migrations complete
		vm, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return nil, err
		}

		// CRITICAL: After migrations, the EVM module's InitGenesis may have overwritten our params with "aatom"
		// We must set them back to "ubadge" and ensure coin info is correct
		evmParams := evmKeeper.GetParams(sdkCtx)
		if evmParams.EvmDenom != "ubadge" {
			sdkCtx.Logger().Info("Fixing EVM params after migrations", "oldEvmDenom", evmParams.EvmDenom, "newEvmDenom", "ubadge")
			evmParams.EvmDenom = "ubadge"
			evmParams.ExtendedDenomOptions = &evmtypes.ExtendedDenomOptions{ExtendedDenom: "abadge"}
			if err := evmKeeper.SetParams(sdkCtx, evmParams); err != nil {
				return nil, err
			}
		}

		// Ensure coin info is set correctly (InitEvmCoinInfo may have failed during InitGenesis)
		if err := evmKeeper.SetEvmCoinInfo(sdkCtx, evmtypes.EvmCoinInfo{
			Denom:         "ubadge",
			ExtendedDenom: "abadge",
			DisplayDenom:  "BADGE",
			Decimals:      9,
		}); err != nil {
			sdkCtx.Logger().Info("EVM coin info already set or error setting after migrations", "error", err)
			// Don't fail the upgrade if it's already set
		}

		return vm, nil
	}
}
