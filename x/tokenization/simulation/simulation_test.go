package simulation

import (
	"math/rand"
	"testing"

	keepertestutil "github.com/bitbadges/bitbadgeschain/x/tokenization/testutil/keeper"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

// TestSimulationOperations tests that all simulation operations can be called without panicking
func TestSimulationOperations(t *testing.T) {
	// Setup keeper and context using the testutil
	k, ctx := keepertestutil.TokenizationKeeper(t)
	
	// For simulation testing, we can pass nil keepers - the simulation functions
	// will handle this gracefully by returning NoOpMsg when needed
	// In a real simulation, these would be the actual app keepers
	var ak types.AccountKeeper = nil
	var bk types.BankKeeper = nil
	
	// Create test accounts with valid bech32 addresses
	acc1, _ := sdk.AccAddressFromBech32("bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430")
	acc2, _ := sdk.AccAddressFromBech32("bb15cftznlenkhfl0ykzwl525mczzrvt7y87thrwx")
	acc3, _ := sdk.AccAddressFromBech32("bb1qypqxpq9qcrsszg2pvxq6rs0zqg3yyc5l5jys0")
	
	accs := []simtypes.Account{
		{Address: acc1, PrivKey: nil},
		{Address: acc2, PrivKey: nil},
		{Address: acc3, PrivKey: nil},
	}
	
	r := rand.New(rand.NewSource(1))
	
	// Test each simulation function - they should return valid OperationMsg (possibly NoOpMsg)
	// but should not panic
	t.Run("SimulateMsgCreateCollection", func(t *testing.T) {
		op := SimulateMsgCreateCollection(ak, bk, k)
		opMsg, futureOps, err := op(r, nil, ctx, accs, "test-chain")
		require.NoError(t, err)
		// OperationMsg is a struct, always valid - check that Route field is set
		require.NotEmpty(t, opMsg.Route)
		// futureOps is a slice, can be nil or empty - both are valid
		_ = futureOps
	})
	
	t.Run("SimulateMsgUniversalUpdateCollection", func(t *testing.T) {
		op := SimulateMsgUniversalUpdateCollection(ak, bk, k)
		opMsg, futureOps, err := op(r, nil, ctx, accs, "test-chain")
		require.NoError(t, err)
		require.NotEmpty(t, opMsg.Route)
		_ = futureOps
	})
	
	t.Run("SimulateMsgDeleteCollection", func(t *testing.T) {
		op := SimulateMsgDeleteCollection(ak, bk, k)
		opMsg, futureOps, err := op(r, nil, ctx, accs, "test-chain")
		require.NoError(t, err)
		require.NotEmpty(t, opMsg.Route)
		_ = futureOps
	})
	
	t.Run("SimulateMsgTransferTokens", func(t *testing.T) {
		op := SimulateMsgTransferTokens(ak, bk, k)
		opMsg, futureOps, err := op(r, nil, ctx, accs, "test-chain")
		require.NoError(t, err)
		require.NotEmpty(t, opMsg.Route)
		_ = futureOps
	})
	
	t.Run("SimulateMsgUpdateUserApprovals", func(t *testing.T) {
		op := SimulateMsgUpdateUserApprovals(ak, bk, k)
		opMsg, futureOps, err := op(r, nil, ctx, accs, "test-chain")
		require.NoError(t, err)
		require.NotEmpty(t, opMsg.Route)
		_ = futureOps
	})
	
	t.Run("SimulateMsgSetIncomingApproval", func(t *testing.T) {
		op := SimulateMsgSetIncomingApproval(ak, bk, k)
		opMsg, futureOps, err := op(r, nil, ctx, accs, "test-chain")
		require.NoError(t, err)
		require.NotEmpty(t, opMsg.Route)
		_ = futureOps
	})
	
	t.Run("SimulateMsgSetOutgoingApproval", func(t *testing.T) {
		op := SimulateMsgSetOutgoingApproval(ak, bk, k)
		opMsg, futureOps, err := op(r, nil, ctx, accs, "test-chain")
		require.NoError(t, err)
		require.NotEmpty(t, opMsg.Route)
		_ = futureOps
	})
	
	t.Run("SimulateMsgPurgeApprovals", func(t *testing.T) {
		op := SimulateMsgPurgeApprovals(ak, bk, k)
		opMsg, futureOps, err := op(r, nil, ctx, accs, "test-chain")
		require.NoError(t, err)
		require.NotEmpty(t, opMsg.Route)
		_ = futureOps
	})
	
	t.Run("SimulateMsgCreateAddressLists", func(t *testing.T) {
		op := SimulateMsgCreateAddressLists(ak, bk, k)
		opMsg, futureOps, err := op(r, nil, ctx, accs, "test-chain")
		require.NoError(t, err)
		require.NotEmpty(t, opMsg.Route)
		_ = futureOps
	})
	
	t.Run("SimulateMsgSetDynamicStoreValue", func(t *testing.T) {
		op := SimulateMsgSetDynamicStoreValue(ak, bk, k)
		opMsg, futureOps, err := op(r, nil, ctx, accs, "test-chain")
		require.NoError(t, err)
		require.NotEmpty(t, opMsg.Route)
		_ = futureOps
	})

	t.Run("SimulateMsgCreateDynamicStore", func(t *testing.T) {
		op := SimulateMsgCreateDynamicStore(ak, bk, k)
		opMsg, futureOps, err := op(r, nil, ctx, accs, "test-chain")
		require.NoError(t, err)
		require.NotEmpty(t, opMsg.Route)
		_ = futureOps
	})

	t.Run("SimulateMsgUpdateDynamicStore", func(t *testing.T) {
		op := SimulateMsgUpdateDynamicStore(ak, bk, k)
		opMsg, futureOps, err := op(r, nil, ctx, accs, "test-chain")
		require.NoError(t, err)
		require.NotEmpty(t, opMsg.Route)
		_ = futureOps
	})
}

// TestSimulationOperationsWithExecution tests that simulation operations can actually execute
func TestSimulationOperationsWithExecution(t *testing.T) {
	// Setup keeper and context
	k, ctx := keepertestutil.TokenizationKeeper(t)
	
	var ak types.AccountKeeper = nil
	var bk types.BankKeeper = nil
	
	// Create test accounts
	acc1, _ := sdk.AccAddressFromBech32("bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430")
	acc2, _ := sdk.AccAddressFromBech32("bb15cftznlenkhfl0ykzwl525mczzrvt7y87thrwx")
	acc3, _ := sdk.AccAddressFromBech32("bb1qypqxpq9qcrsszg2pvxq6rs0zqg3yyc5l5jys0")
	
	accs := []simtypes.Account{
		{Address: acc1, PrivKey: nil},
		{Address: acc2, PrivKey: nil},
		{Address: acc3, PrivKey: nil},
	}
	
	// Setup initial state
	r := rand.New(rand.NewSource(1))
	err := SetupSimulationState(ctx, k, accs, r)
	require.NoError(t, err)
	
	// Track success rates
	successCount := 0
	totalCount := 0
	
	// Test CreateCollection with execution
	t.Run("SimulateMsgCreateCollection_WithExecution", func(t *testing.T) {
		op := SimulateMsgCreateCollection(ak, bk, k)
		for i := 0; i < 3; i++ {
			opMsg, _, err := op(r, nil, ctx, accs, "test-chain")
			require.NoError(t, err)
			totalCount++
			if opMsg.OK {
				// OperationMsg.Msg is []byte, we can't easily decode it here
				// But we can verify the operation generated a valid message
				require.NotEmpty(t, opMsg.Msg)
			}
		}
	})
	
	// Test CreateDynamicStore with execution
	t.Run("SimulateMsgCreateDynamicStore_WithExecution", func(t *testing.T) {
		op := SimulateMsgCreateDynamicStore(ak, bk, k)
		for i := 0; i < 3; i++ {
			opMsg, _, err := op(r, nil, ctx, accs, "test-chain")
			require.NoError(t, err)
			totalCount++
			if opMsg.OK {
				require.NotEmpty(t, opMsg.Msg)
				successCount++
			}
		}
	})
	
	// Test TransferTokens with execution (after creating collections)
	t.Run("SimulateMsgTransferTokens_WithExecution", func(t *testing.T) {
		op := SimulateMsgTransferTokens(ak, bk, k)
		for i := 0; i < 3; i++ {
			opMsg, _, err := op(r, nil, ctx, accs, "test-chain")
			require.NoError(t, err)
			totalCount++
			if opMsg.OK {
				require.NotEmpty(t, opMsg.Msg)
				successCount++
			}
		}
	})
	
	// Report success rate
	t.Logf("Simulation execution success rate: %d/%d (%.2f%%)", successCount, totalCount, float64(successCount)/float64(totalCount)*100)
}

// TestMultiRunOperation tests that MultiRunOperation works correctly
func TestMultiRunOperation(t *testing.T) {
	k, ctx := keepertestutil.TokenizationKeeper(t)
	
	var ak types.AccountKeeper = nil
	var bk types.BankKeeper = nil
	
	acc1, _ := sdk.AccAddressFromBech32("bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430")
	accs := []simtypes.Account{
		{Address: acc1, PrivKey: nil},
	}
	
	r := rand.New(rand.NewSource(1))
	
	// Setup state first
	err := SetupSimulationState(ctx, k, accs, r)
	require.NoError(t, err)
	
	// Test that MultiRunOperation wraps correctly
	op := SimulateMsgCreateCollection(ak, bk, k)
	wrappedOp := MultiRunOperation(op, 3)
	
	opMsg, _, err := wrappedOp(r, nil, ctx, accs, "test-chain")
	require.NoError(t, err)
	require.NotEmpty(t, opMsg.Route)
}

// TestBoundedTimelineTimes tests that bounded timeline generation works
func TestBoundedTimelineTimes(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	
	// Test with small bounds
	times := GetBoundedTimelineTimes(r, 5, 1, 10)
	require.Len(t, times, 5)
	for _, timeRange := range times {
		require.GreaterOrEqual(t, timeRange.Start.Uint64(), uint64(1))
		require.LessOrEqual(t, timeRange.End.Uint64(), uint64(10))
		require.LessOrEqual(t, timeRange.Start.Uint64(), timeRange.End.Uint64())
	}
	
	// Test default GetTimelineTimes uses bounds
	times2 := GetTimelineTimes(r, 3)
	require.Len(t, times2, 3)
	for _, timeRange := range times2 {
		require.GreaterOrEqual(t, timeRange.Start.Uint64(), uint64(MinTimelineRange))
		require.LessOrEqual(t, timeRange.End.Uint64(), uint64(MaxTimelineRange))
	}
}

// TestGetKnownGoodCollectionId tests the known-good collection ID helper
func TestGetKnownGoodCollectionId(t *testing.T) {
	k, ctx := keepertestutil.TokenizationKeeper(t)
	
	// Initially no collections
	collectionId, found := GetKnownGoodCollectionId(ctx, k)
	require.False(t, found)
	require.True(t, collectionId.IsZero())
	
	// Create a collection directly using GetOrCreateCollection
	acc1, _ := sdk.AccAddressFromBech32("bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430")
	accs := []simtypes.Account{{Address: acc1, PrivKey: nil}}
	r := rand.New(rand.NewSource(1))
	
	// Use GetOrCreateCollection which ensures a collection exists
	createdId, err := GetOrCreateCollection(ctx, k, acc1.String(), r, accs)
	require.NoError(t, err)
	require.False(t, createdId.IsZero())
	
	// Now should find a collection
	collectionId, found = GetKnownGoodCollectionId(ctx, k)
	require.True(t, found)
	require.Equal(t, sdkmath.NewUint(1), collectionId)
}

