package app

import (
	"github.com/ethereum/go-ethereum/common"

	evmmodule "github.com/cosmos/evm/x/vm"
	evmkeeper "github.com/cosmos/evm/x/vm/keeper"
	evmtypes "github.com/cosmos/evm/x/vm/types"

	storetypes "cosmossdk.io/store/types"

	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	tokenizationprecompile "github.com/bitbadges/bitbadgeschain/x/evm/precompiles/tokenization"
)

// registerEVMModules registers the EVM module and tokenization precompile
func (app *App) registerEVMModules(appOpts servertypes.AppOptions) error {
	// Use evmtypes.StoreKey to match what the EVM keeper expects
	evmKey := storetypes.NewKVStoreKey(evmtypes.StoreKey)

	// Register store key
	if err := app.RegisterStores(evmKey); err != nil {
		return err
	}

	// Create EVM keeper
	// Note: Use DefaultEVMChainID to allow SetChainConfig to work correctly in parallel tests
	// SetChainConfig only allows overwriting if the existing chain ID equals DefaultEVMChainID
	// If you use a custom chain ID, SetChainConfig will panic on the second call in parallel tests
	evmChainID := evmtypes.DefaultEVMChainID
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)

	storeKeys := make(map[string]*storetypes.KVStoreKey)
	storeKeys[evmtypes.StoreKey] = evmKey

	app.EVMKeeper = evmkeeper.NewKeeper(
		app.appCodec,
		evmKey,
		storetypes.NewTransientStoreKey(evmtypes.TransientKey),
		storeKeys,
		authority,
		app.AccountKeeper,
		app.PreciseBankKeeper, // Use PreciseBankKeeper for fractional balance support
		app.StakingKeeper,
		nil, // FeeMarketKeeper - can be nil for basic setup
		app.ConsensusParamsKeeper,
		nil, // ERC20Keeper - can be nil for basic setup
		evmChainID,
		"", // tracer - empty for now
	).WithDefaultEvmCoinInfo(evmtypes.EvmCoinInfo{
		Denom:         "ubadge",
		ExtendedDenom: "ubadge",
		DisplayDenom:  "BADGE",
		Decimals:      9,
	})

	// Register tokenization precompile
	tokenizationPrecompile := tokenizationprecompile.NewPrecompile(app.TokenizationKeeper)
	tokenizationPrecompileAddr := common.HexToAddress(tokenizationprecompile.TokenizationPrecompileAddress)
	app.EVMKeeper.RegisterStaticPrecompile(tokenizationPrecompileAddr, tokenizationPrecompile)

	// Note: Precompiles will be enabled during InitGenesis or via governance
	// We don't enable them here during app initialization as the store context isn't fully available yet

	// Register EVM module
	evmModule := evmmodule.NewAppModule(
		app.EVMKeeper,
		app.AccountKeeper,
		app.PreciseBankKeeper, // Use PreciseBankKeeper for EVM module
		app.AccountKeeper.AddressCodec(),
	)

	if err := app.RegisterModules(evmModule); err != nil {
		return err
	}

	return nil
}
