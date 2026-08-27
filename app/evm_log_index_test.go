package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestEVMLogIndexIsBlockGlobalAcrossTransactions pins the behaviour that
// app.go's vmrunner.SetRunner call exists to provide.
//
// cosmos/evm v0.6 numbered EVM logs block-globally inside the StateDB. v0.7
// changed the StateDB to number them per transaction, restarting at 0, and moved
// block-global renumbering into evmtypes.PatchTxResponses — which only runs if a
// tx runner is installed on baseapp. baseapp silently falls back to a bare
// DefaultRunner when none is, so blocks still execute and everything looks fine
// while every transaction in the block reports logIndex starting from 0.
//
// (blockHash, logIndex) is the uniqueness key every indexer, explorer and dapp
// uses for EVM logs, so that collision breaks eth_getLogs and receipts. And
// because PatchTxResponses rewrites res.Data and res.Events, it feeds
// LastResultsHash — so shipping it wrong could not be corrected without another
// coordinated upgrade.
//
// This test fails if vmrunner.SetRunner is removed from app.go. That is the
// point: the previous suite passed either way.
func TestEVMLogIndexIsBlockGlobalAcrossTransactions(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)

	key, from := newFundedEVMAccount(t, app, ctx)
	contract := deployLogEmitter(t, app, key, from, 1)

	// Two calls, same block. Each emits two logs.
	nonce := nonceOf(t, app, from)
	results := deliverBlock(t, app, 2, from,
		callLogEmitter(t, key, contract, nonce),
		callLogEmitter(t, key, contract, nonce+1),
	)

	first := ethLogs(t, results[0])
	second := ethLogs(t, results[1])

	require.Len(t, first, logsPerCall, "first tx should emit %d logs", logsPerCall)
	require.Len(t, second, logsPerCall, "second tx should emit %d logs", logsPerCall)

	// The assertion that matters. Without a runner the second transaction's
	// logs come back as 0,1 — colliding with the first transaction's.
	gotIndices := []uint64{
		first[0].Index, first[1].Index,
		second[0].Index, second[1].Index,
	}
	require.Equal(t, []uint64{0, 1, 2, 3}, gotIndices,
		"log Index must continue across transactions in a block, not restart per transaction; "+
			"a 0,1,0,1 result means vmrunner.SetRunner is missing and PatchTxResponses never ran")

	// TxIndex must identify which transaction each log came from.
	gotTxIndices := []uint64{
		first[0].TxIndex, first[1].TxIndex,
		second[0].TxIndex, second[1].TxIndex,
	}
	require.Equal(t, []uint64{0, 0, 1, 1}, gotTxIndices,
		"log TxIndex must be the transaction's position in the block")

	// (blockHash, logIndex) is the key indexers rely on; assert it is unique
	// across the block rather than only that the sequence looks right.
	seen := map[uint64]bool{}
	for _, idx := range gotIndices {
		require.False(t, seen[idx], "duplicate logIndex %d within one block", idx)
		seen[idx] = true
	}
}

// A single transaction in a block is the degenerate case, and it is the one a
// keeper-level test would also pass. Included so a regression that breaks only
// the multi-transaction path is clearly distinguishable from one that breaks
// log emission outright.
func TestEVMLogIndexSingleTransactionStartsAtZero(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)

	key, from := newFundedEVMAccount(t, app, ctx)
	contract := deployLogEmitter(t, app, key, from, 1)

	results := deliverBlock(t, app, 2, from, callLogEmitter(t, key, contract, nonceOf(t, app, from)))
	logs := ethLogs(t, results[0])

	require.Len(t, logs, logsPerCall)
	require.Equal(t, uint64(0), logs[0].Index)
	require.Equal(t, uint64(1), logs[1].Index)
	require.Equal(t, uint64(0), logs[0].TxIndex)
}

// Log indices restart at 0 in each new block. They are block-scoped, not
// chain-scoped — an indexer keying on (blockHash, logIndex) depends on that.
func TestEVMLogIndexRestartsEachBlock(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)

	key, from := newFundedEVMAccount(t, app, ctx)
	contract := deployLogEmitter(t, app, key, from, 1)

	firstBlock := deliverBlock(t, app, 2, from, callLogEmitter(t, key, contract, nonceOf(t, app, from)))
	secondBlock := deliverBlock(t, app, 3, from, callLogEmitter(t, key, contract, nonceOf(t, app, from)))

	firstLogs := ethLogs(t, firstBlock[0])
	secondLogs := ethLogs(t, secondBlock[0])

	require.Equal(t, uint64(0), firstLogs[0].Index)
	require.Equal(t, uint64(0), secondLogs[0].Index,
		"a new block must restart log numbering at 0")
}
