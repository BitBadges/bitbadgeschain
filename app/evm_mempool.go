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

	// NOTE: this deliberately does NOT gate on json-rpc.enable.
	//
	// Whether the EVM mempool is active is decided by app.toml's
	// mempool.max-txs, and cosmos/evm enforces that it agree with config.toml's
	// mempool.type (see server/config.ValidateCrossConfig). JSON-RPC is not part
	// of that contract.
	//
	// Gating here on json-rpc.enable broke the default configuration: the
	// cross-config check saw the EVM mempool as enabled and required
	// mempool.type = "app", but this function returned before installing any of
	// the handlers. CometBFT then routed to an app mempool with no ReapTxs
	// handler, so the node logged "ReapTxs handler not set" on every block and
	// could never include a transaction — while still producing empty blocks,
	// so it looked healthy.
	//
	// The cosmosPoolMaxTx < 0 check below is the correct gate and mirrors
	// upstream's own condition exactly.
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
