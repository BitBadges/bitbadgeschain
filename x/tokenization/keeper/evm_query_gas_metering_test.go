package keeper_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/evm/x/vm/statedb"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	testutilkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/testutil/keeper"
)

const (
	gasMeteringContract = "0x1234567890123456789012345678901234567890"
	// The gas an attacker's contract burns before reverting.
	gasMeteringBurned = uint64(480000)
	gasMeteringLimit  = uint64(500000)
)

// gasMeteredMockEVMKeeper reproduces cosmos/evm v0.7's CallEVMWithData contract
// exactly, because the whole defect lives in that contract's shape: when the VM
// reverts, CallEVMWithData returns a NON-NIL response (carrying the real
// GasUsed) *alongside* a non-nil error. See x/vm/keeper/call_evm.go:97-100.
//
// A mock that returned (nil, err) on revert would let the buggy code pass.
type gasMeteredMockEVMKeeper struct {
	gasUsed uint64
	vmError string
	// hardErr makes the call fail before the VM produces any measurement, the
	// way an intrinsic-gas or sequence lookup failure does: (nil, err).
	hardErr error
}

func (m *gasMeteredMockEVMKeeper) IsContract(sdk.Context, common.Address) bool { return true }

func (m *gasMeteredMockEVMKeeper) CallEVMWithData(
	_ sdk.Context, _ *statedb.StateDB, _ common.Address, _ *common.Address,
	_ []byte, _ bool, _ bool, _ *big.Int,
) (*evmtypes.MsgEthereumTxResponse, error) {
	if m.hardErr != nil {
		return nil, m.hardErr
	}
	res := &evmtypes.MsgEthereumTxResponse{
		GasUsed: m.gasUsed,
		VmError: m.vmError,
		Ret:     []byte{0x01},
	}
	if m.vmError != "" {
		// Mirror call_evm.go: the response AND an error come back together.
		return res, errorsmod.Wrap(evmtypes.ErrVMExecution, m.vmError)
	}
	return res, nil
}

// measureEVMQueryGas runs one EVM query against a freshly built keeper and
// returns how much Cosmos gas the caller's meter was charged.
//
// Each call gets its own keeper so the fixed, non-EVM overhead (creating the
// zero-address caller account, store reads) is identical between runs and
// cancels when two measurements are compared.
func measureEVMQueryGas(t *testing.T, evmKeeper *gasMeteredMockEVMKeeper) uint64 {
	t.Helper()
	k, ctx := testutilkeeper.TokenizationKeeper(t)
	k.SetEVMKeeper(evmKeeper)

	calldata, err := hex.DecodeString("70a08231")
	require.NoError(t, err)

	before := ctx.GasMeter().GasConsumed()
	_, _ = k.ExecuteEVMQuery(ctx, gasMeteringContract, calldata, gasMeteringLimit)
	return ctx.GasMeter().GasConsumed() - before
}

// TestExecuteEVMQuery_RevertingCallIsChargedLikeASuccessfulOne is the guard for
// the unmetered-validator-CPU defect.
//
// ExecuteEVMQueryWithCaller returned on `err != nil` / `VmError != ""` BEFORE
// reaching its ConsumeGas call, so a contract that burned its entire gas limit
// and then reverted cost the caller zero Cosmos gas. The EVM work itself was
// still done by every validator.
//
// The gas limits live in the collection's approval criteria and invariants, not
// in the transaction, so an attacker can publish a collection whose approval
// carries several max-gas challenges pointed at a gas-burning reverting
// contract and then spam transfers: validators do the work, the transaction
// fails, nothing is charged. Raising the total budget 1M -> 2.5M multiplies the
// free work per approval by 2.5x, and one MsgTransferTokens evaluates approvals
// per transfer item across the collection, incoming and outgoing levels.
//
// The assertion is the property that actually matters: reverting must cost the
// same as succeeding, so failing is never the cheaper way to buy compute.
func TestExecuteEVMQuery_RevertingCallIsChargedLikeASuccessfulOne(t *testing.T) {
	succeeded := measureEVMQueryGas(t, &gasMeteredMockEVMKeeper{gasUsed: gasMeteringBurned})
	reverted := measureEVMQueryGas(t, &gasMeteredMockEVMKeeper{
		gasUsed: gasMeteringBurned,
		vmError: "execution reverted",
	})

	require.Equal(t, succeeded, reverted,
		"a query that burned %d gas and reverted must cost the caller exactly what the same query costs "+
			"when it succeeds; charging less makes reverting a way to buy free validator compute",
		gasMeteringBurned)
	require.GreaterOrEqual(t, reverted, gasMeteringBurned,
		"the caller must be charged at least the %d gas the EVM actually burned", gasMeteringBurned)
}

// TestExecuteEVMQuery_HardFailureChargesTheDeclaredLimit covers the other
// failure shape: cosmos/evm can fail before the VM produces any GasUsed at all
// (intrinsic gas, gas overflow, sequence lookup), returning (nil, err). There is
// nothing measured to charge, so the declared limit is charged instead — it is
// deterministic across nodes because it comes from the collection's own stored
// configuration, and it keeps early failure from being the cheap path.
func TestExecuteEVMQuery_HardFailureChargesTheDeclaredLimit(t *testing.T) {
	charged := measureEVMQueryGas(t, &gasMeteredMockEVMKeeper{
		hardErr: errorsmod.Wrap(evmtypes.ErrInvalidGasCap, "intrinsic gas too low"),
	})

	require.GreaterOrEqual(t, charged, gasMeteringLimit,
		"a call that fails before the VM reports usage must still cost the declared %d gas limit",
		gasMeteringLimit)
}

// TestExecuteEVMQuery_SuccessfulCallStillChargesMeasuredGas pins that the fix
// did not disturb the pre-existing success path.
func TestExecuteEVMQuery_SuccessfulCallStillChargesMeasuredGas(t *testing.T) {
	charged := measureEVMQueryGas(t, &gasMeteredMockEVMKeeper{gasUsed: gasMeteringBurned})
	require.GreaterOrEqual(t, charged, gasMeteringBurned)
}
