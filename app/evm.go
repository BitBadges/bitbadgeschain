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
	tokenizationprecompile "github.com/bitbadges/bitbadgeschain/x/evm/precompiles/tokenization"
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

	storeKeys := make(map[string]*storetypes.KVStoreKey)
	storeKeys[evmtypes.StoreKey] = evmKey

	// Create EVM keeper first with pointer to ERC20 keeper (ERC20 keeper not yet initialized)
	// This follows the EVMD pattern: we pass &app.ERC20Keeper even though it's not initialized yet
	// The EVM keeper will hold a reference to the ERC20 keeper, which will be initialized below
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
