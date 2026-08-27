package app

import (
	"fmt"
	"runtime"

	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	"github.com/cosmos/cosmos-sdk/baseapp/txnrunner"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	vmrunner "github.com/cosmos/evm/x/vm/runner"
	evmtypes "github.com/cosmos/evm/x/vm/types"
)

// EVMParallelExecutionFlag opts a node into block-STM parallel block execution.
//
// Off by default, and deliberately so — see configureTxRunner.
const EVMParallelExecutionFlag = "evm.parallel-execution"

// configureTxRunner installs the EVM block tx runner.
//
// Installing *a* runner is required for correctness, not performance. In
// cosmos/evm v0.6 the StateDB numbered logs block-globally during execution:
//
//	log.Index = s.txConfig.LogIndex + uint(len(s.logs))
//
// v0.7 changed that to per-transaction numbering restarting at 0 for every tx,
// and moved block-global renumbering into evmtypes.PatchTxResponses, which only
// runs if a runner is installed. baseapp silently falls back to a bare
// DefaultRunner when none is, so blocks execute and everything looks healthy
// while every tx in the block reports logIndex from 0 — colliding on
// (blockHash, logIndex), the uniqueness key every indexer and explorer uses.
// PatchTxResponses rewrites res.Data and res.Events, so it feeds
// LastResultsHash and could not be corrected without another coordinated
// upgrade. app/evm_log_index_test.go pins this.
//
// The runner *choice* is the open question (BB-8). The sequential DefaultRunner
// is the default here, and switching to block-STM is opt-in per node via
// --evm.parallel-execution. That asymmetry is deliberate:
//
//   - Parallel execution across this chain's custom modules (x/tokenization,
//     x/gamm, x/poolmanager, x/sendmanager, x/custom-hooks) and their
//     precompiles, which share state, is a real determinism risk. block-STM's
//     conflict detection is only as good as the store-key surface it sees, and
//     a precompile reaching state outside that surface produces a conflict it
//     cannot detect — which shows up as a consensus fault, not a test failure.
//   - Execution order is consensus-critical. A flag that changes it must not be
//     set per-node in production: two validators disagreeing on this flag is a
//     chain halt. It exists so the option can be measured on an isolated
//     devnet, not so operators can tune it.
//   - Adopting block-STM permanently forecloses re-enabling the block gas
//     meter: baseapp's SetBlockSTMTxRunner panics when handed an *STMRunner
//     while the meter is on. cosmos-sdk v0.54 defaults disableBlockGasMeter to
//     true, so this chain would not trip it today — but it relies on that
//     default rather than setting it explicitly the way evmd does.
//
// Before this is considered for a default, BB-8 wants determinism verified
// across repeated runs of a block containing conflicting transactions. That
// test does not exist yet, and nothing here should be read as evidence the
// parallel path is safe.
func (app *App) configureTxRunner(appOpts servertypes.AppOptions) error {
	txDecoder := app.txConfig.TxDecoder()

	if !parallelExecutionEnabled(appOpts) {
		vmrunner.SetRunner(app.BaseApp, txnrunner.NewDefaultRunner(txDecoder))
		return nil
	}

	if len(app.evmNonTransientKeys) == 0 {
		return fmt.Errorf(
			"%s is set but no EVM store keys were collected; parallel execution cannot detect conflicts without them",
			EVMParallelExecutionFlag,
		)
	}

	workers := min(runtime.GOMAXPROCS(0), runtime.NumCPU())
	if workers < 1 {
		workers = 1
	}

	app.Logger().Info(
		"EVM parallel block execution (block-STM) is ENABLED — experimental; "+
			"execution order is consensus-critical and every validator must agree on this setting",
		"workers", workers,
		"stores", len(app.evmNonTransientKeys),
	)

	vmrunner.SetRunner(app.BaseApp, txnrunner.NewSTMRunner(
		txDecoder,
		app.evmNonTransientKeys,
		workers,
		true, // estimate: pre-scan each tx's likely read/write set
		func(storetypes.MultiStore) string { return evmtypes.GetEVMCoinDenom() },
	))
	return nil
}

func parallelExecutionEnabled(appOpts servertypes.AppOptions) bool {
	if appOpts == nil {
		return false
	}
	v := appOpts.Get(EVMParallelExecutionFlag)
	if v == nil {
		return false
	}
	enabled, ok := v.(bool)
	return ok && enabled
}
