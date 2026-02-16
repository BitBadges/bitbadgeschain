package v24_patch

import (
	"context"
	"fmt"
	"sort"

	"github.com/ethereum/go-ethereum/common"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	gammprecompile "github.com/bitbadges/bitbadgeschain/x/gamm/precompile"
	sendmanagerprecompile "github.com/bitbadges/bitbadgeschain/x/sendmanager/precompile"
	tokenizationprecompile "github.com/bitbadges/bitbadgeschain/x/tokenization/precompile"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	evmkeeper "github.com/cosmos/evm/x/vm/keeper"
	evmtypes "github.com/cosmos/evm/x/vm/types"
)

const UpgradeName = "v24-patch"

var precompileAddressStrings = []string{
	"0x0000000000000000000000000000000000000800", // Staking
	"0x0000000000000000000000000000000000000801", // Distribution
	"0x0000000000000000000000000000000000000802", // ICS20 (IBC)
	"0x0000000000000000000000000000000000000803", // Vesting
	"0x0000000000000000000000000000000000000804", // Bank
	"0x0000000000000000000000000000000000000805", // Governance
	"0x0000000000000000000000000000000000000806", // Slashing
	tokenizationprecompile.TokenizationPrecompileAddress,
	gammprecompile.GammPrecompileAddress,
	sendmanagerprecompile.SendManagerPrecompileAddress,
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

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	evmKeeper *evmkeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
) func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)

		// Ensure ubadge denom metadata exists
		if _, found := bankKeeper.GetDenomMetaData(ctx, "ubadge"); !found {
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
		}

		// Fix EVM params
		evmParams := evmKeeper.GetParams(sdkCtx)
		evmParams.EvmDenom = "ubadge"
		evmParams.ExtendedDenomOptions = &evmtypes.ExtendedDenomOptions{ExtendedDenom: "abadge"}
		evmParams.ActiveStaticPrecompiles = getAllPrecompileAddressStrings()
		if err := evmKeeper.SetParams(sdkCtx, evmParams); err != nil {
			return nil, fmt.Errorf("failed to set EVM params: %w", err)
		}

		// Fix EVM coin info
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

		sdkCtx.Logger().Info("v24-patch upgrade complete",
			"evmDenom", finalParams.EvmDenom,
			"precompiles", len(finalParams.ActiveStaticPrecompiles))

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
