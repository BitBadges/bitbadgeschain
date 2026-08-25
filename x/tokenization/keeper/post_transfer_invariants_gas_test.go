package keeper_test

import (
	"fmt"
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/evm/x/vm/statedb"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	testutilkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/testutil/keeper"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// The v33 post-transfer invariant gas budget, hardcoded because it is history.
// Deriving it from the current constants would make every assertion below
// vacuous — the whole point is to compare v34 against what v33 accepted.
const (
	v33DefaultPostTransferEVMQueryGasLimit = uint64(100000)
	v33MaxTotalPostTransferInvariantGas    = uint64(1000000)
)

// v33PostTransferDefaultGasCapacity is how many default-gas post-transfer
// invariants a collection could carry on v33 and still transfer.
const v33PostTransferDefaultGasCapacity = int(v33MaxTotalPostTransferInvariantGas / v33DefaultPostTransferEVMQueryGasLimit)

// postTransferMockEVMKeeper is a always-succeeds EVM keeper that records what
// gas each query was actually given.
type postTransferMockEVMKeeper struct {
	calls     int
	gasLimits []uint64
}

func (m *postTransferMockEVMKeeper) IsContract(sdk.Context, common.Address) bool { return true }

func (m *postTransferMockEVMKeeper) CallEVMWithData(
	_ sdk.Context, _ *statedb.StateDB, _ common.Address, _ *common.Address,
	_ []byte, _ bool, _ bool, gasCap *big.Int,
) (*evmtypes.MsgEthereumTxResponse, error) {
	m.calls++
	m.gasLimits = append(m.gasLimits, gasCap.Uint64())
	return &evmtypes.MsgEthereumTxResponse{Ret: []byte{0x01}}, nil
}

// postTransferInvariantCollection builds a collection carrying n EVM query
// invariants. A zero gasLimit means "use the module default".
func postTransferInvariantCollection(n int, gasLimit uint64) *types.TokenCollection {
	challenges := make([]*types.EVMQueryChallenge, n)
	for i := range challenges {
		challenges[i] = &types.EVMQueryChallenge{
			ContractAddress: "0x1234567890123456789012345678901234567890",
			Calldata:        "70a08231",
			GasLimit:        sdkmath.NewUint(gasLimit),
		}
	}
	return &types.TokenCollection{
		CollectionId: sdkmath.NewUint(1),
		Invariants:   &types.CollectionInvariants{EvmQueryChallenges: challenges},
	}
}

func postTransferKeeperWithMockEVM(t *testing.T) (*keeper.Keeper, sdk.Context, *postTransferMockEVMKeeper) {
	t.Helper()
	k, ctx := testutilkeeper.TokenizationKeeper(t)
	mock := &postTransferMockEVMKeeper{}
	k.SetEVMKeeper(mock)
	return k, ctx, mock
}

// TestPostTransferInvariants_DefaultGasCapacityNotReducedFromV33 is the
// regression guard for the v34 post-transfer gas budget. It is the same defect
// that was fixed for approval criteria in evm_query_challenges.go, in a second
// file that the first fix missed.
//
// v34 raised DefaultPostTransferEVMQueryGasLimit from 100000 to 250000 because
// cosmos/evm v0.7 stateful precompiles need more headroom. The total cap did not
// move, and the loop reserves each invariant's DECLARED limit rather than its
// measured usage, so capacity fell from 10 to 4.
//
// That is not a tightened policy, it is a break against live state. A collection
// already on chain with 5 default-gas invariants validated on v33
// (5 x 100000 = 500000 <= 1000000) and would fail EVERY transfer forever on v34
// (5 x 250000 = 1250000 > 1000000). The transferring user cannot work around it:
// the gas limits live in the collection's invariants, not in their transaction.
//
// ValidateEVMQueryChallenges permits up to 10 invariants per collection, at
// create and at update alike, with no total-gas rule at write time, so 10 is
// genuinely reachable state.
//
// This asserts on behaviour (a check that must not error), not on the constants,
// so it fails for the reason that actually matters.
func TestPostTransferInvariants_DefaultGasCapacityNotReducedFromV33(t *testing.T) {
	for n := 1; n <= v33PostTransferDefaultGasCapacity; n++ {
		t.Run(fmt.Sprintf("%d_default_gas_invariants", n), func(t *testing.T) {
			k, ctx, mock := postTransferKeeperWithMockEVM(t)

			collection := postTransferInvariantCollection(n, 0)
			err := k.CheckPostTransferInvariants(ctx, collection, bob, []string{alice}, charlie)

			require.NoError(t, err,
				"a collection with %d default-gas post-transfer invariants passed on v33 (%d x %d = %d <= %d) "+
					"and must still pass; otherwise the upgrade permanently bricks every transfer of it",
				n, n, v33DefaultPostTransferEVMQueryGasLimit,
				uint64(n)*v33DefaultPostTransferEVMQueryGasLimit, v33MaxTotalPostTransferInvariantGas)

			require.Equal(t, n, mock.calls, "every invariant must actually have been executed")
			// The gas each invariant ran with is the v34 default, not the v33
			// one: the point is the capacity, not the per-call gas.
			for i, got := range mock.gasLimits {
				require.Equal(t, keeper.DefaultPostTransferEVMQueryGasLimit, got, "invariant %d gas limit", i)
			}
		})
	}
}

// TestPostTransferInvariants_TotalGasCapStillRejectsExcess proves the DoS
// ceiling did not simply get switched off. The count is derived from the
// constants so it cannot rot.
func TestPostTransferInvariants_TotalGasCapStillRejectsExcess(t *testing.T) {
	// Smallest number of max-gas invariants whose declared total exceeds the cap.
	n := int(keeper.MaxTotalPostTransferInvariantGas/keeper.MaxPostTransferEVMQueryGasLimit) + 1

	k, ctx, mock := postTransferKeeperWithMockEVM(t)
	collection := postTransferInvariantCollection(n, keeper.MaxPostTransferEVMQueryGasLimit)

	err := k.CheckPostTransferInvariants(ctx, collection, bob, []string{alice}, charlie)

	require.Error(t, err, "%d invariants at %d gas each must exceed the %d cap",
		n, keeper.MaxPostTransferEVMQueryGasLimit, keeper.MaxTotalPostTransferInvariantGas)
	require.Contains(t, err.Error(), "exceed total gas limit")
	require.Less(t, mock.calls, n, "the over-budget invariant must be rejected before it runs")
}

// TestPostTransferInvariants_TotalGasCapAllowsExactlyTheCap pins the boundary:
// the check is `>` the cap, so a budget landing exactly on it passes.
func TestPostTransferInvariants_TotalGasCapAllowsExactlyTheCap(t *testing.T) {
	n := int(keeper.MaxTotalPostTransferInvariantGas / keeper.MaxPostTransferEVMQueryGasLimit)
	require.Greater(t, n, 0, "the cap must admit at least one max-gas invariant")

	k, ctx, mock := postTransferKeeperWithMockEVM(t)
	collection := postTransferInvariantCollection(n, keeper.MaxPostTransferEVMQueryGasLimit)

	err := k.CheckPostTransferInvariants(ctx, collection, bob, []string{alice}, charlie)
	require.NoError(t, err, "a total landing exactly on the cap must pass")
	require.Equal(t, n, mock.calls)
}

// TestPostTransferInvariantGasBudgetInvariants documents the arithmetic the two
// behavioural tests depend on, and names the intent for whoever next reaches for
// one of these constants.
func TestPostTransferInvariantGasBudgetInvariants(t *testing.T) {
	require.LessOrEqual(t, keeper.DefaultPostTransferEVMQueryGasLimit, keeper.MaxPostTransferEVMQueryGasLimit,
		"the default must be expressible as an explicit limit")
	require.GreaterOrEqual(t,
		keeper.MaxTotalPostTransferInvariantGas/keeper.DefaultPostTransferEVMQueryGasLimit,
		v33MaxTotalPostTransferInvariantGas/v33DefaultPostTransferEVMQueryGasLimit,
		"v34 must not fit fewer default-gas post-transfer invariants than v33 did")
}
