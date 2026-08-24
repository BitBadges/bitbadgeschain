package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// parallelOpts turns on the block-STM tx runner for a test app.
func parallelOpts() map[string]interface{} {
	return map[string]interface{}{EVMParallelExecutionFlag: true}
}

// TestEVMParallelRunnerPreservesLogIndexing checks that the block-STM runner
// produces the same block-scoped log numbering as the sequential one.
//
// This is evidence for BB-8, not a decision. It shows the parallel path is
// wired correctly and does not break the invariant the sequential runner exists
// to provide — it does NOT show that parallel execution is deterministic across
// this chain's custom modules, which is the actual open risk and needs a
// conflicting-transaction determinism harness that does not exist yet.
//
// Two transactions calling the same contract in one block is also, deliberately,
// a conflicting workload: both touch the same account nonce and the same
// contract. If block-STM's conflict detection did not see those, this is where
// it would show up first.
func TestEVMParallelRunnerPreservesLogIndexing(t *testing.T) {
	app := SetupWithAppOptions(false, parallelOpts())
	ctx := app.NewContext(false)

	key, from := newFundedEVMAccount(t, app, ctx)
	contract := deployLogEmitter(t, app, key, from, 1)

	nonce := nonceOf(t, app, from)
	results := deliverBlock(t, app, 2, from,
		callLogEmitter(t, key, contract, nonce),
		callLogEmitter(t, key, contract, nonce+1),
	)

	first := ethLogs(t, results[0])
	second := ethLogs(t, results[1])
	require.Len(t, first, logsPerCall)
	require.Len(t, second, logsPerCall)

	require.Equal(t,
		[]uint64{0, 1, 2, 3},
		[]uint64{first[0].Index, first[1].Index, second[0].Index, second[1].Index},
		"block-STM must number logs block-globally, exactly as the sequential runner does")

	require.Equal(t,
		[]uint64{0, 0, 1, 1},
		[]uint64{first[0].TxIndex, first[1].TxIndex, second[0].TxIndex, second[1].TxIndex},
		"block-STM must preserve transaction ordering in TxIndex")
}
