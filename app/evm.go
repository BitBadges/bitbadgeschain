package app

import (
	"fmt"
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
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	appparams "github.com/bitbadges/bitbadgeschain/app/params"
	gammprecompile "github.com/bitbadges/bitbadgeschain/x/gamm/precompile"
	sendmanagerprecompile "github.com/bitbadges/bitbadgeschain/x/sendmanager/precompile"
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

	// Get EVM chain ID from build-time flag (set via ldflags)
	// Defaults to 90123 (local dev) if not set at build time
	evmChainIDStr := appparams.GetEVMChainID()
	evmChainIDUint64, err := strconv.ParseUint(evmChainIDStr, 10, 64)
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
		"acc",          // Account store (for address lookups)
		"bank",         // Bank store (for balance operations)
		"tokenization", // Tokenization store (for precompile operations)
		"evm",          // EVM store (for EVM state)
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

	// Register all custom precompiles
	// Note: DefaultStaticPrecompiles are already registered via WithStaticPrecompiles above
	app.registerCustomPrecompiles()

	// Validate no address collisions between custom precompiles
	// This helps catch bugs early if addresses are accidentally duplicated
	if err := ValidateNoAddressCollisions(); err != nil {
		panic(fmt.Sprintf("precompile address collision detected: %v", err))
	}

	// Note: Precompiles must be both registered (RegisterStaticPrecompile) and enabled (EnableStaticPrecompiles)
	// The precompiles are registered above, but they will be enabled during InitGenesis when the EVM module initializes
	// For production, ensure all precompile addresses are in the genesis state's active_static_precompiles array
	// For tests, we enable them programmatically in the test setup (see evm_keeper_integration_test.go)

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

// registerCustomPrecompiles registers all custom BitBadges precompiles
// This ensures all custom precompiles are registered and prevents address collisions
func (app *App) registerCustomPrecompiles() {
	// Register tokenization precompile
	tokenizationPrecompile := tokenizationprecompile.NewPrecompile(app.TokenizationKeeper)
	tokenizationPrecompileAddr := common.HexToAddress(tokenizationprecompile.TokenizationPrecompileAddress)
	app.EVMKeeper.RegisterStaticPrecompile(tokenizationPrecompileAddr, tokenizationPrecompile)

	// Register gamm precompile
	gammPrecompile := gammprecompile.NewPrecompile(app.GammKeeper)
	gammPrecompileAddr := common.HexToAddress(gammprecompile.GammPrecompileAddress)
	app.EVMKeeper.RegisterStaticPrecompile(gammPrecompileAddr, gammPrecompile)

	// Register sendmanager precompile
	sendManagerPrecompile := sendmanagerprecompile.NewPrecompile(app.SendmanagerKeeper)
	sendManagerPrecompileAddr := common.HexToAddress(sendmanagerprecompile.SendManagerPrecompileAddress)
	app.EVMKeeper.RegisterStaticPrecompile(sendManagerPrecompileAddr, sendManagerPrecompile)

	// Add additional custom precompiles here as needed
	// Next available address: 0x0000000000000000000000000000000000001004
}

// GetDefaultCosmosPrecompileAddresses returns all default Cosmos precompile addresses
// These are the standard precompiles provided by cosmos/evm (staking, distribution, bank, etc.)
// Addresses from cosmos/evm/x/vm/types/precompiles.go
func GetDefaultCosmosPrecompileAddresses() []common.Address {
	// Default Cosmos precompiles from cosmos/evm/x/vm/types/precompiles.go
	// These are the actual addresses used by the cosmos/evm module
	return []common.Address{
		common.HexToAddress("0x0000000000000000000000000000000000000800"), // Staking precompile
		common.HexToAddress("0x0000000000000000000000000000000000000801"), // Distribution precompile
		common.HexToAddress("0x0000000000000000000000000000000000000802"), // ICS20 (IBC) precompile
		common.HexToAddress("0x0000000000000000000000000000000000000803"), // Vesting precompile
		common.HexToAddress("0x0000000000000000000000000000000000000804"), // Bank precompile
		common.HexToAddress("0x0000000000000000000000000000000000000805"), // Governance precompile
		common.HexToAddress("0x0000000000000000000000000000000000000806"), // Slashing precompile
	}
}

// GetAllPrecompileAddresses returns all precompile addresses (both default Cosmos and custom BitBadges)
// This is used for enabling precompiles in genesis, upgrades, and tests
func (app *App) GetAllPrecompileAddresses() []common.Address {
	var addresses []common.Address

	// Add default Cosmos precompile addresses
	addresses = append(addresses, GetDefaultCosmosPrecompileAddresses()...)

	// Add custom BitBadges precompile addresses
	addresses = append(addresses, common.HexToAddress(tokenizationprecompile.TokenizationPrecompileAddress)) // 0x1001
	addresses = append(addresses, common.HexToAddress(gammprecompile.GammPrecompileAddress))                 // 0x1002
	// Note: Use Cosmos default bank precompile at 0x0804 (already included in default precompiles)

	return addresses
}

// EnableAllPrecompiles enables all registered precompiles (both default Cosmos and custom BitBadges)
// This should be called during genesis initialization or upgrade handlers
// It's safe to call multiple times (idempotent)
func (app *App) EnableAllPrecompiles(ctx sdk.Context) error {
	// Enable all precompiles (both default Cosmos and custom BitBadges)
	allAddresses := app.GetAllPrecompileAddresses()
	for _, addr := range allAddresses {
		if err := app.EVMKeeper.EnableStaticPrecompiles(ctx, addr); err != nil {
			// Log error but don't fail if precompile is already enabled (idempotent)
			// This allows the function to be called multiple times safely
			ctx.Logger().Info("Precompile enable attempt", "error", err, "address", addr.Hex())
		}
	}

	return nil
}
