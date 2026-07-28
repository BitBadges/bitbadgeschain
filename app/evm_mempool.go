package app

import (
	"cosmossdk.io/log/v2"
	"github.com/cosmos/cosmos-sdk/baseapp"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	evmmempool "github.com/cosmos/evm/mempool"
	"github.com/cosmos/evm/server"
	evmtypes "github.com/cosmos/evm/x/vm/types"
)

// configureEVMMempool sets up the EVM mempool and related handlers using app options.
// This is an advanced configuration - only needed if JSON-RPC is enabled.
// The mempool is NOT created automatically - the app must create it.
//
// cosmos/evm v0.7 removed NewExperimentalEVMMempool and the hand-rolled
// EVMMempoolConfig. The mempool is now built from server.ResolveMempoolConfig
// and driven through explicit Insert/Reap/CheckTx handlers on BaseApp, which
// requires `mempool.type = "app"` in config.toml to take effect.
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

	mpConfig := server.ResolveMempoolConfig(app.BaseApp.AnteHandler(), appOpts, logger)

	cosmosPoolMaxTx := server.GetCosmosPoolMaxTx(appOpts, logger)
	if cosmosPoolMaxTx < 0 {
		logger.Debug("app-side mempool is disabled, skipping evm mempool configuration")
		return nil
	}

	if err := server.ValidateReapBounds(appOpts, mpConfig.BlockGasLimit); err != nil {
		return err
	}

	txEncoder := evmmempool.NewTxEncoder(app.txConfig)
	evmRechecker := evmmempool.NewTxRechecker(mpConfig.AnteHandler, txEncoder)
	cosmosRechecker := evmmempool.NewTxRechecker(mpConfig.AnteHandler, txEncoder)
	checkTxTimeout := server.GetMempoolCheckTxTimeout(appOpts, logger)

	mempool := evmmempool.NewMempool(
		app.CreateQueryContext,
		logger,
		app.EVMKeeper,
		app.FeeMarketKeeper,
		app.txConfig,
		evmRechecker,
		cosmosRechecker,
		mpConfig,
		cosmosPoolMaxTx,
	)
	app.EVMMempool = mempool

	abciProposalHandler := baseapp.NewDefaultProposalHandler(mempool, app)
	app.SetPrepareProposal(abciProposalHandler.PrepareProposalHandler())
	app.SetInsertTxHandler(mempool.NewInsertTxHandler(app.TxDecode))
	app.SetReapTxsHandler(mempool.NewReapTxsHandler())
	app.SetCheckTxHandler(mempool.NewCheckTxHandler(app.TxDecode, checkTxTimeout))
	app.SetMempool(mempool)

	// v0.7 removed the VM keeper's end-blocker notification; the mempool is
	// told about new blocks here instead when no event bus is wired up.
	app.SetPrepareCheckStater(func(_ sdk.Context) {
		if !mempool.HasEventBus() {
			mempool.NotifyNewBlock()
		}
	})

	return nil
}
