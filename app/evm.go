package app

import (
	"github.com/ethereum/go-ethereum/common"

	feemarket "github.com/cosmos/evm/x/feemarket"
	feemarketkeeper "github.com/cosmos/evm/x/feemarket/keeper"
	feemarkettypes "github.com/cosmos/evm/x/feemarket/types"
	evmmodule "github.com/cosmos/evm/x/vm"
	evmkeeper "github.com/cosmos/evm/x/vm/keeper"
	evmtypes "github.com/cosmos/evm/x/vm/types"

	storetypes "cosmossdk.io/store/types"

	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	tokenizationprecompile "github.com/bitbadges/bitbadgeschain/x/evm/precompiles/tokenization"
)

// registerEVMModules registers the FeeMarket and EVM modules with tokenization precompile
func (app *App) registerEVMModules(appOpts servertypes.AppOptions) error {
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)

	// Create FeeMarket store keys
	feemarketKey := storetypes.NewKVStoreKey(feemarkettypes.StoreKey)
	feemarketTransientKey := storetypes.NewTransientStoreKey(feemarkettypes.TransientKey)

	// Register feemarket store keys
	if err := app.RegisterStores(feemarketKey, feemarketTransientKey); err != nil {
		return err
	}

	// Create FeeMarket keeper
	app.FeeMarketKeeper = feemarketkeeper.NewKeeper(
		app.appCodec,
		authority,
		feemarketKey,
		feemarketTransientKey,
	)

	// Register FeeMarket module
	feemarketModule := feemarket.NewAppModule(app.FeeMarketKeeper)
	if err := app.RegisterModules(feemarketModule); err != nil {
		return err
	}

	// Create EVM store keys
	evmKey := storetypes.NewKVStoreKey(evmtypes.StoreKey)
	evmTransientKey := storetypes.NewTransientStoreKey(evmtypes.TransientKey)

	// Register EVM store keys
	if err := app.RegisterStores(evmKey, evmTransientKey); err != nil {
		return err
	}

	// Create EVM keeper
	// Note: Use DefaultEVMChainID to allow SetChainConfig to work correctly in parallel tests
	// SetChainConfig only allows overwriting if the existing chain ID equals DefaultEVMChainID
	// If you use a custom chain ID, SetChainConfig will panic on the second call in parallel tests
	evmChainID := evmtypes.DefaultEVMChainID

	storeKeys := make(map[string]*storetypes.KVStoreKey)
	storeKeys[evmtypes.StoreKey] = evmKey

	app.EVMKeeper = configureEVMKeeper(evmkeeper.NewKeeper(
		app.appCodec,
		evmKey,
		evmTransientKey,
		storeKeys,
		authority,
		app.AccountKeeper,
		app.PreciseBankKeeper, // Use PreciseBankKeeper for fractional balance support
		app.StakingKeeper,
		app.FeeMarketKeeper, // Use FeeMarket keeper
		app.ConsensusParamsKeeper,
		nil, // ERC20Keeper - can be nil for basic setup
		evmChainID,
		"", // tracer - empty for now
	))

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
