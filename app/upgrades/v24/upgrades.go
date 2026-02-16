package v24

import (
	"context"
	"fmt"
	"sort"
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
	erc20keeper "github.com/cosmos/evm/x/erc20/keeper"
	erc20types "github.com/cosmos/evm/x/erc20/types"
	feemarketkeeper "github.com/cosmos/evm/x/feemarket/keeper"
	feemarkettypes "github.com/cosmos/evm/x/feemarket/types"
	precisebanktypes "github.com/cosmos/evm/x/precisebank/types"
	evmkeeper "github.com/cosmos/evm/x/vm/keeper"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	// Uncomment when configuring rate limits:
)

const (
	UpgradeName = "v24"
)

// precompileAddressStrings are all precompile addresses (sorted for determinism)
var precompileAddressStrings = []string{
	"0x0000000000000000000000000000000000000800",         // Staking
	"0x0000000000000000000000000000000000000801",         // Distribution
	"0x0000000000000000000000000000000000000802",         // ICS20 (IBC)
	"0x0000000000000000000000000000000000000803",         // Vesting
	"0x0000000000000000000000000000000000000804",         // Bank
	"0x0000000000000000000000000000000000000805",         // Governance
	"0x0000000000000000000000000000000000000806",         // Slashing
	tokenizationprecompile.TokenizationPrecompileAddress, // 0x1001
	gammprecompile.GammPrecompileAddress,                 // 0x1002
	sendmanagerprecompile.SendManagerPrecompileAddress,   // 0x1003
}

func getAllPrecompileAddressStrings() []string {
	result := make([]string, len(precompileAddressStrings))
	copy(result, precompileAddressStrings)
	sort.Strings(result)
	return result
}

func getAllPrecompileAddresses() []common.Address {
	addrs := make([]common.Address, len(precompileAddressStrings))
	for i, s := range precompileAddressStrings {
		addrs[i] = common.HexToAddress(s)
	}
	return addrs
}

// CustomUpgradeHandlerLogic runs pre-migration setup (can be tested with a snapshot)
func CustomUpgradeHandlerLogic(
	ctx context.Context,
	tokenizationKeeper tokenizationkeeper.Keeper,
	poolManagerKeeper poolmanagerkeeper.Keeper,
	rateLimitKeeper ibcratelimitkeeper.Keeper,
	evmKeeper *evmkeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Set denom metadata for "ubadge" - required for EVM module initialization
	bankKeeper.SetDenomMetaData(ctx, banktypes.Metadata{
		Description: "The native token of BitBadges Chain",
		DenomUnits: []*banktypes.DenomUnit{
			{Denom: "ubadge", Exponent: 0},
			{Denom: "badge", Exponent: 9},
		},
		Base:    "ubadge",
		Display: "badge",
		Name:    "Badge",
		Symbol:  "BADGE",
	})

	// Ensure EVM chain ID is set correctly (from build-time flag)
	if err := ensureEVMChainID(sdkCtx); err != nil {
		return err
	}

	// Run tokenization migrations
	return tokenizationKeeper.MigrateTokenizationKeeper(sdkCtx)
}

// ensureEVMChainID ensures the EVM chain ID matches the build-time value.
func ensureEVMChainID(ctx sdk.Context) error {
	expectedEVMChainIDStr := appparams.GetEVMChainID()
	expectedEVMChainID, err := strconv.ParseUint(expectedEVMChainIDStr, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse EVM chain ID %s: %w", expectedEVMChainIDStr, err)
	}

	chainConfig := evmtypes.GetChainConfig()
	if chainConfig == nil {
		return nil
	}

	if chainConfig.GetChainId() != expectedEVMChainID {
		newChainConfig := evmtypes.DefaultChainConfig(expectedEVMChainID)
		// SetChainConfig may fail if already set - this is expected
		_ = evmtypes.SetChainConfig(newChainConfig)
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
	feeMarketKeeper feemarketkeeper.Keeper,
	erc20Keeper erc20keeper.Keeper,
) func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)

		// Run pre-migration setup
		if err := CustomUpgradeHandlerLogic(ctx, tokenizationKeeper, poolManagerKeeper, rateLimitKeeper, evmKeeper, bankKeeper); err != nil {
			return nil, err
		}

		// Skip default InitGenesis for EVM modules to prevent "aatom" default denom
		// We manually initialize with "ubadge" configuration below
		fromVM[evmtypes.ModuleName] = 1
		fromVM[feemarkettypes.ModuleName] = 1
		fromVM[erc20types.ModuleName] = 1
		fromVM[precisebanktypes.ModuleName] = 1

		// Run module migrations (EVM modules skipped)
		vm, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return nil, err
		}

		// Initialize FeeMarket params
		if err := feeMarketKeeper.SetParams(sdkCtx, feemarkettypes.DefaultParams()); err != nil {
			return nil, fmt.Errorf("failed to set FeeMarket params: %w", err)
		}

		// Initialize ERC20 params
		if err := erc20Keeper.SetParams(sdkCtx, erc20types.DefaultParams()); err != nil {
			return nil, fmt.Errorf("failed to set ERC20 params: %w", err)
		}

		// Initialize EVM params with ubadge denom and precompiles
		evmParams := evmtypes.Params{
			EvmDenom:                "ubadge",
			ExtraEIPs:               evmtypes.DefaultExtraEIPs,
			ActiveStaticPrecompiles: getAllPrecompileAddressStrings(),
			EVMChannels:             evmtypes.DefaultEVMChannels,
			AccessControl:           evmtypes.DefaultAccessControl,
			HistoryServeWindow:      evmtypes.DefaultHistoryServeWindow,
			ExtendedDenomOptions:    &evmtypes.ExtendedDenomOptions{ExtendedDenom: "abadge"},
		}
		if err := evmKeeper.SetParams(sdkCtx, evmParams); err != nil {
			return nil, fmt.Errorf("failed to set EVM params: %w", err)
		}

		// Set EVM coin info (used by ante handler for fee validation)
		coinInfo := evmtypes.EvmCoinInfo{
			Denom:         "ubadge",
			ExtendedDenom: "abadge",
			DisplayDenom:  "BADGE",
			Decimals:      9,
		}
		if err := evmKeeper.SetEvmCoinInfo(sdkCtx, coinInfo); err != nil {
			return nil, fmt.Errorf("failed to set EVM coin info: %w", err)
		}

		// Enable all precompiles
		for _, addr := range getAllPrecompileAddresses() {
			_ = evmKeeper.EnableStaticPrecompiles(sdkCtx, addr)
		}

		// Verify configuration
		finalParams := evmKeeper.GetParams(sdkCtx)
		finalCoinInfo := evmKeeper.GetEvmCoinInfo(sdkCtx)
		if finalParams.EvmDenom != "ubadge" || finalCoinInfo.Denom != "ubadge" {
			return nil, fmt.Errorf("EVM configuration mismatch - params: %s, coinInfo: %s",
				finalParams.EvmDenom, finalCoinInfo.Denom)
		}

		sdkCtx.Logger().Info("v24 upgrade complete",
			"evmDenom", finalParams.EvmDenom,
			"precompiles", len(finalParams.ActiveStaticPrecompiles))

		return vm, nil
	}
}
