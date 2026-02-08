package app

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common"

	precompiletypes "github.com/cosmos/evm/precompiles/types"
	erc20 "github.com/cosmos/evm/x/erc20"
	erc20keeper "github.com/cosmos/evm/x/erc20/keeper"
	erc20types "github.com/cosmos/evm/x/erc20/types"
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

	appparams "github.com/bitbadges/bitbadgeschain/app/params"
	tokenizationprecompile "github.com/bitbadges/bitbadgeschain/x/tokenization/precompile"
)

// registerEVMModules registers the FeeMarket, ERC20, and EVM modules with tokenization precompile
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

	// Create ERC20 store keys
	erc20Key := storetypes.NewKVStoreKey(erc20types.StoreKey)

	// Register ERC20 store keys
	if err := app.RegisterStores(erc20Key); err != nil {
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
	evmChainID := appparams.EVMChainID
	evmChainIDUint64, err := strconv.ParseUint(evmChainID, 10, 64)
	if err != nil {
		return err
	}

	// Get all non-transient KV store keys for the EVM keeper
	// The EVM keeper needs access to all stores so precompiles can access any store they need
	// (e.g., account store, bank store, tokenization store, etc.)
	// Note: Our version of EVM keeper expects a map, not a slice
	// IMPORTANT: The EVM keeper's keys parameter only accepts KV store keys, not transient stores.
	// Transient stores from other modules (like customhooks_transient) cannot be included in this map.
	// The EVM keeper's snapshotmulti.Store is built from this map, so it won't have access to other modules' transient stores.
	//
	// NOTE: There is a known issue with the EVM keeper's snapshot management when precompiles return errors.
	// The snapshot stack can be empty when trying to revert, causing "snapshot index 0 out of bound [0..0)" panics.
	// This is a bug in the upstream cosmos/evm module. As a workaround, precompiles should handle errors gracefully
	// and avoid returning errors that would trigger EVM reverts.
	//
	// DEBUG: Collect all KV store keys for snapshotter initialization
	storeKeysMap := make(map[string]*storetypes.KVStoreKey)
	allStoreKeys := app.GetStoreKeys()
	
	// Log store key collection for debugging
	kvStoreCount := 0
	transientStoreCount := 0
	otherStoreCount := 0
	
	for _, key := range allStoreKeys {
		switch k := key.(type) {
		case *storetypes.KVStoreKey:
			storeKeysMap[k.Name()] = k
			kvStoreCount++
		case *storetypes.TransientStoreKey:
			transientStoreCount++
			// Transient stores are intentionally excluded from EVM keeper's storeKeysMap
			// because the EVM keeper's snapshotmulti.Store only supports KV stores
		default:
			otherStoreCount++
			// Note: ObjectStoreKeys are not currently used in this codebase
			// If evmd pattern shows we need them, we would need to check if EVM keeper supports them
		}
	}
	
	// Log store registration summary (only in debug/test builds)
	// This helps diagnose snapshot issues by verifying all stores are registered
	if len(storeKeysMap) == 0 {
		panic("EVM keeper requires at least one KV store key for snapshotter initialization")
	}
	
	// Verify critical stores are included (stores that precompiles might access)
	criticalStores := []string{
		"acc",           // Account store (for address lookups)
		"bank",          // Bank store (for balance operations)
		"tokenization", // Tokenization store (for precompile operations)
		"evm",           // EVM store (for EVM state)
	}
	
	missingStores := []string{}
	for _, storeName := range criticalStores {
		if _, found := storeKeysMap[storeName]; !found {
			missingStores = append(missingStores, storeName)
		}
	}
	
	if len(missingStores) > 0 {
		// Log warning but don't panic - some stores might be optional
		// This is a diagnostic aid, not a hard requirement
	}

	// Create EVM keeper first with pointer to ERC20 keeper (ERC20 keeper not yet initialized)
	// This follows the EVMD pattern: we pass &app.ERC20Keeper even though it's not initialized yet
	// The EVM keeper will hold a reference to the ERC20 keeper, which will be initialized below
	// Note: For object key, we use evmKey as a fallback since object store keys may not be registered
	// The EVM keeper will work with just the KV store key if object keys aren't available
	app.EVMKeeper = configureEVMKeeper(evmkeeper.NewKeeper(
		app.appCodec,
		evmKey,
		evmTransientKey, // EVM keeper's own transient key (for EVM module use only)
		storeKeysMap,    // Only KV store keys - cannot include other modules' transient stores
		authority,
		app.AccountKeeper,
		app.PreciseBankKeeper, // Use PreciseBankKeeper for fractional balance support
		app.StakingKeeper,
		app.FeeMarketKeeper, // Use FeeMarket keeper
		app.ConsensusParamsKeeper,
		&app.ERC20Keeper, // Pass pointer to ERC20 keeper (not yet initialized, but that's ok)
		evmChainIDUint64,
		"", // tracer - empty for now
	).WithStaticPrecompiles(
		precompiletypes.DefaultStaticPrecompiles(
			*app.StakingKeeper,
			app.DistrKeeper,
			app.PreciseBankKeeper,
			&app.ERC20Keeper,
			&app.TransferKeeper, // Cosmos/evm transfer keeper (wraps ibc-go, adds ERC20 support for IBC v2/eureka)
			app.IBCKeeper.ChannelKeeper,
			*app.GovKeeper,
			app.SlashingKeeper,
			app.appCodec,
		),
	))

	// Create ERC20 keeper with pointer to transfer keeper
	// The transfer keeper is created in registerIBCModules (which is called before registerEVMModules)
	// but it holds a pointer to &app.ERC20Keeper, so we can pass &app.TransferKeeper here
	app.ERC20Keeper = erc20keeper.NewKeeper(
		erc20Key,
		app.appCodec,
		authority,
		app.AccountKeeper,
		app.PreciseBankKeeper, // Use PreciseBankKeeper for bank operations
		app.EVMKeeper,
		app.StakingKeeper,
		&app.TransferKeeper, // Cosmos/evm transfer keeper (created in registerIBCModules with pointer to this ERC20 keeper)
	)

	// Register ERC20 module
	erc20Module := erc20.NewAppModule(app.ERC20Keeper, app.AccountKeeper)
	if err := app.RegisterModules(erc20Module); err != nil {
		return err
	}

	// Register tokenization precompile
	tokenizationPrecompile := tokenizationprecompile.NewPrecompile(app.TokenizationKeeper)
	tokenizationPrecompileAddr := common.HexToAddress(tokenizationprecompile.TokenizationPrecompileAddress)
	app.EVMKeeper.RegisterStaticPrecompile(tokenizationPrecompileAddr, tokenizationPrecompile)

	// Note: Precompiles must be both registered (RegisterStaticPrecompile) and enabled (EnableStaticPrecompiles)
	// The precompile is registered above, but it will be enabled during InitGenesis when the EVM module initializes
	// For production, ensure the precompile address is in the genesis state's active_static_precompiles array
	// For tests, we enable it programmatically in the test setup (see evm_keeper_integration_test.go)

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
