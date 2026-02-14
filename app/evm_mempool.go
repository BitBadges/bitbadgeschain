package app

import (
	"fmt"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdkmempool "github.com/cosmos/cosmos-sdk/types/mempool"

	evmmempool "github.com/cosmos/evm/mempool"
	evmtypes "github.com/cosmos/evm/x/vm/types"
)

// configureEVMMempool sets up the EVM mempool and related handlers using app options.
// This is an advanced configuration - only needed if JSON-RPC is enabled.
// The mempool is NOT created automatically - the app must create it.
func (app *App) configureEVMMempool(appOpts servertypes.AppOptions, logger log.Logger) error {
	if evmtypes.GetChainConfig() == nil {
		logger.Debug("evm chain config is not set, skipping mempool configuration")
		return nil
	}

	// Check if JSON-RPC is enabled - mempool is only needed for JSON-RPC
	// If JSON-RPC is disabled, we can skip mempool configuration (it's an advanced feature)
	jsonRPCEnabled := false
	if val := appOpts.Get("json-rpc.enable"); val != nil {
		if b, ok := val.(bool); ok {
			jsonRPCEnabled = b
		}
	}
	if !jsonRPCEnabled {
		logger.Debug("JSON-RPC is disabled, skipping EVM mempool configuration (advanced feature)")
		return nil
	}

	// Get cosmos pool max tx from app options (defaults to 1000 if not set)
	// Note: server.GetCosmosPoolMaxTx is not exported in v0.5.1, so we read directly
	cosmosPoolMaxTx := 1000 // Default value
	if val := appOpts.Get("mempool.max-txs"); val != nil {
		if i, ok := val.(int); ok && i > 0 {
			cosmosPoolMaxTx = i
		}
	}
	if cosmosPoolMaxTx < 0 {
		logger.Debug("app-side mempool is disabled, skipping evm mempool configuration")
		return nil
	}

	mempoolConfig, err := app.createMempoolConfig(appOpts, logger)
	if err != nil {
		return fmt.Errorf("failed to get mempool config: %w", err)
	}

	// Create a minimal client context - the server will update it later with full context
	clientCtx := client.Context{}.
		WithCodec(app.appCodec).
		WithTxConfig(app.txConfig).
		WithInterfaceRegistry(app.interfaceRegistry)

	// Create EVM mempool - clientCtx will be updated later by the server
	evmMempool := evmmempool.NewExperimentalEVMMempool(
		app.CreateQueryContext,
		logger,
		app.EVMKeeper,
		app.FeeMarketKeeper,
		app.txConfig,
		clientCtx,
		mempoolConfig,
		cosmosPoolMaxTx,
	)
	app.EVMMempool = evmMempool
	app.SetMempool(evmMempool)
	checkTxHandler := evmmempool.NewCheckTxHandler(evmMempool)
	app.SetCheckTxHandler(checkTxHandler)

	abciProposalHandler := baseapp.NewDefaultProposalHandler(evmMempool, app)
	abciProposalHandler.SetSignerExtractionAdapter(
		evmmempool.NewEthSignerExtractionAdapter(
			sdkmempool.NewDefaultSignerExtractionAdapter(),
		),
	)
	app.SetPrepareProposal(abciProposalHandler.PrepareProposalHandler())

	return nil
}

// createMempoolConfig creates a new EVMMempoolConfig with default configuration.
// Note: server helper functions (GetLegacyPoolConfig, etc.) are not exported in v0.5.1,
// so we use simple defaults. These can be configured via app.toml in the future.
func (app *App) createMempoolConfig(appOpts servertypes.AppOptions, logger log.Logger) (*evmmempool.EVMMempoolConfig, error) {
	// Get block gas limit from app options (defaults to 100M if not set, matching config.yml max_gas)
	// This allows configuration via app.toml: [evm.mempool] block-gas-limit = 100000000
	blockGasLimit := uint64(100_000_000) // Default: 100M gas (matches config.yml max_gas)
	if val := appOpts.Get("evm.mempool.block-gas-limit"); val != nil {
		if i, ok := val.(uint64); ok && i > 0 {
			blockGasLimit = i
		} else if i, ok := val.(int); ok && i > 0 {
			blockGasLimit = uint64(i)
		} else if i, ok := val.(int64); ok && i > 0 {
			blockGasLimit = uint64(i)
		}
	}
	
	logger.Info("EVM mempool block gas limit configured", "block_gas_limit", blockGasLimit)

	return &evmmempool.EVMMempoolConfig{
		AnteHandler:      app.BaseApp.AnteHandler(),
		LegacyPoolConfig: nil, // Default to nil - can be configured via app.toml
		BlockGasLimit:    blockGasLimit,
		MinTip:           nil, // Default to nil - can be configured via app.toml
	}, nil
}
