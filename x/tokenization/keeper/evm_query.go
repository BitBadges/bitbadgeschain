package keeper

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/bitbadges/bitbadgeschain/pkg/evmcompat"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmconfig "github.com/cosmos/evm/server/config"
	evmkeeper "github.com/cosmos/evm/x/vm/keeper"
	"github.com/cosmos/evm/x/vm/statedb"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
)

type evmQueryContextKey struct{}

// ValidateEVMQueryChallengesAreContracts ensures every EVM query challenge's contract address
// has code on the EVM (is a contract). Call this when storing collections or approvals that
// include EVM query challenges. When evmKeeper is nil, validation is skipped (e.g. in tests
// or chains without EVM); execution will still fail at runtime if the address has no code.
func (k Keeper) ValidateEVMQueryChallengesAreContracts(ctx sdk.Context, challenges []*types.EVMQueryChallenge) error {
	if len(challenges) == 0 {
		return nil
	}
	if k.evmKeeper == nil {
		return nil // Skip when EVM not available; execution-time check will fail if needed
	}
	for i, challenge := range challenges {
		if challenge == nil || challenge.ContractAddress == "" {
			continue
		}
		var contractAddr common.Address
		addr := strings.TrimSpace(challenge.ContractAddress)
		if len(addr) >= 2 && strings.ToLower(addr[:2]) == "0x" {
			contractAddr = common.HexToAddress(addr)
		} else {
			accAddr, err := sdk.AccAddressFromBech32(addr)
			if err != nil {
				return fmt.Errorf("EVM query challenge %d: invalid contract address: %w", i, err)
			}
			contractAddr = common.BytesToAddress(accAddr.Bytes())
		}
		if !k.evmKeeper.IsContract(ctx, contractAddr) {
			return fmt.Errorf("EVM query challenge %d: address is not a contract (no code): %s", i, challenge.ContractAddress)
		}
	}
	return nil
}

// ExecuteEVMQuery performs a read-only EVM contract call (uses zero address as caller).
// Use ExecuteEVMQueryWithCaller when the caller must be an existing account (e.g. for approval/invariant flows).
func (k Keeper) ExecuteEVMQuery(ctx sdk.Context, contractAddress string, calldata []byte, gasLimit uint64) ([]byte, error) {
	return k.ExecuteEVMQueryWithCaller(ctx, "", contractAddress, calldata, gasLimit)
}

// ExecuteEVMQueryWithCaller performs a read-only EVM contract call with the given caller.
// If callerAddress is empty, the zero address is used.
// Note: cosmos/evm requires the caller account to exist in auth keeper for GetSequence().
// For read-only queries, we ensure the zero address account exists temporarily if needed.
func (k Keeper) ExecuteEVMQueryWithCaller(ctx sdk.Context, callerAddress string, contractAddress string, calldata []byte, gasLimit uint64) ([]byte, error) {
	if ctx.Value(evmQueryContextKey{}) != nil {
		return nil, fmt.Errorf("nested EVM query is not supported")
	}
	if k.evmKeeper == nil {
		return nil, fmt.Errorf("EVM keeper not available")
	}

	// Convert contract address to common.Address
	var contractAddr common.Address
	if len(contractAddress) >= 2 && strings.ToLower(contractAddress[:2]) == "0x" {
		contractAddr = common.HexToAddress(contractAddress)
	} else {
		accAddr, err := sdk.AccAddressFromBech32(contractAddress)
		if err != nil {
			return nil, fmt.Errorf("invalid contract address: %w", err)
		}
		contractAddr = common.BytesToAddress(accAddr.Bytes())
	}

	// Verify it's a contract
	if !k.evmKeeper.IsContract(ctx, contractAddr) {
		return nil, fmt.Errorf("address is not a contract: %s", contractAddress)
	}

	// Resolve caller: use provided address, otherwise zero address
	var callerAddr common.Address
	if callerAddress != "" {
		if len(callerAddress) >= 2 && strings.ToLower(callerAddress[:2]) == "0x" {
			callerAddr = common.HexToAddress(callerAddress)
		} else {
			accAddr, err := sdk.AccAddressFromBech32(callerAddress)
			if err != nil {
				return nil, fmt.Errorf("invalid caller address: %w", err)
			}
			callerAddr = common.BytesToAddress(accAddr.Bytes())
		}
	}
	// If callerAddress is empty, callerAddr will be the zero address (default value)

	if gasLimit == 0 {
		gasLimit = evmconfig.DefaultGasCap
	}
	if remaining := ctx.GasMeter().GasRemaining(); gasLimit > remaining {
		gasLimit = remaining
	}
	if gasLimit == 0 {
		return nil, fmt.Errorf("no gas remaining for EVM query")
	}
	gasCap := new(big.Int).SetUint64(gasLimit)

	// Isolate the read-only query in an atomic context so any snapshot frames or
	// transient writes are rolled back on return while preserving gas billing on ctx.
	atomic := evmcompat.NewAtomicContext(ctx)
	defer atomic.Rollback()

	// Run the EVM call against an isolated gas meter.
	//
	// cosmos/evm v0.7's stateful-precompile entrypoint (RunNativeAction) starts
	// by charging the *already consumed* gas of the incoming context against a
	// fresh meter capped at the gas the EVM forwarded to the precompile:
	//
	//	initialGas := ctx.GasMeter().GasConsumed()
	//	ctx = ctx.WithGasMeter(storetypes.NewGasMeter(contract.Gas))
	//	ctx.GasMeter().ConsumeGas(initialGas, "creating a new gas meter")
	//
	// A normal EVM transaction survives that because its context starts with a
	// fresh ante-handler meter. EVM queries do not: they run in the middle of a
	// Cosmos message (approval checks and post-transfer invariants), so the
	// incoming meter has already consumed far more than the forwarded gas and
	// the precompile aborts with "out of gas" before its handler ever runs.
	//
	// Metering the call separately is also what this function already assumes:
	// the EVM gas actually used is charged back onto the caller's meter below.
	evmCtx := atomic.Ctx().WithValue(evmQueryContextKey{}, true).WithGasMeter(storetypes.NewGasMeter(gasLimit))

	var (
		response *evmtypes.MsgEthereumTxResponse
		err      error
	)
	if realKeeper, ok := k.evmKeeper.(*evmkeeper.Keeper); ok {
		response, err = executeStaticEVMCall(evmCtx, realKeeper, callerAddr, contractAddr, calldata, gasLimit)
	} else {
		// Mock keeper path in unit tests: keep any account setup inside the rolled-back context.
		callerAccAddr := sdk.AccAddress(callerAddr.Bytes())
		if k.accountKeeper != nil && k.accountKeeper.GetAccount(evmCtx, callerAccAddr) == nil {
			acc := k.accountKeeper.NewAccountWithAddress(evmCtx, callerAccAddr)
			k.accountKeeper.SetAccount(evmCtx, acc)
		}
		response, err = k.evmKeeper.CallEVMWithData(evmCtx, nil, callerAddr, &contractAddr, calldata, false, false, gasCap)
	}

	// Charge for the EVM work BEFORE deciding whether the query succeeded.
	// The EVM ran on evmCtx's isolated meter, so we bill the measured (or limit)
	// gas directly onto the caller's outer gas meter.
	chargeEVMQueryGas(ctx, response, err, gasLimit)

	if err != nil {
		return nil, fmt.Errorf("EVM call failed: %w", err)
	}

	if response.VmError != "" {
		return nil, fmt.Errorf("EVM execution error: %s", response.VmError)
	}

	return response.Ret, nil
}

// executeStaticEVMCall runs a read-only EVM StaticCall using the real EVM keeper.
// Using vm.EVM.StaticCall enforces EVM write-protection for state-changing opcodes
// and stateful precompile transaction methods.
func executeStaticEVMCall(
	ctx sdk.Context,
	realKeeper *evmkeeper.Keeper,
	callerAddr common.Address,
	contractAddr common.Address,
	calldata []byte,
	gasLimit uint64,
) (*evmtypes.MsgEthereumTxResponse, error) {
	cfg, err := realKeeper.EVMConfig(ctx, ctx.BlockHeader().ProposerAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to load evm config: %w", err)
	}

	msg := core.Message{
		From:       callerAddr,
		To:         &contractAddr,
		Nonce:      0,
		Value:      big.NewInt(0),
		GasLimit:   gasLimit,
		GasPrice:   big.NewInt(0),
		GasTipCap:  big.NewInt(0),
		GasFeeCap:  big.NewInt(0),
		Data:       calldata,
		AccessList: ethtypes.AccessList{},
	}

	sdb := statedb.New(ctx, realKeeper, statedb.NewEmptyTxConfig())
	ethCfg := evmtypes.GetEthChainConfig()
	evm := realKeeper.NewEVMWithOverridePrecompiles(ctx, msg, cfg, nil, sdb, true)
	rules := ethCfg.Rules(evm.Context.BlockNumber, true, evm.Context.Time)

	leftoverGas := msg.GasLimit
	intrinsicGas, err := realKeeper.GetEthIntrinsicGas(ctx, msg, ethCfg, false)
	if err != nil {
		return nil, fmt.Errorf("intrinsic gas failed: %w", err)
	}
	if leftoverGas < intrinsicGas {
		return nil, fmt.Errorf("%w: have %d, want %d", core.ErrIntrinsicGas, leftoverGas, intrinsicGas)
	}
	if rules.IsPrague {
		floorDataGas, err := core.FloorDataGas(msg.Data)
		if err != nil {
			return nil, err
		}
		if msg.GasLimit < floorDataGas {
			return nil, fmt.Errorf("%w: have %d, want %d", core.ErrFloorDataGas, msg.GasLimit, floorDataGas)
		}
	}
	leftoverGas -= intrinsicGas

	sdb.Prepare(rules, msg.From, common.Address{}, msg.To, evm.ActivePrecompiles(), msg.AccessList)
	if addr, ok := ethtypes.ParseDelegation(sdb.GetCode(contractAddr)); ok {
		sdb.AddAddressToAccessList(addr)
	}

	ret, leftoverGas, vmErr := evm.StaticCall(msg.From, contractAddr, msg.Data, leftoverGas)
	if msg.GasLimit < leftoverGas {
		return nil, fmt.Errorf("message gas limit < leftover gas (%d < %d)", msg.GasLimit, leftoverGas)
	}
	gasUsed := msg.GasLimit - leftoverGas

	var vmError string
	if vmErr != nil {
		vmError = vmErr.Error()
		if vmError == vm.ErrExecutionReverted.Error() {
			ret = evm.ReturnData()
		}
	}

	return &evmtypes.MsgEthereumTxResponse{
		GasUsed:    gasUsed,
		MaxUsedGas: gasUsed,
		VmError:    vmError,
		Ret:        ret,
	}, nil
}

// chargeEVMQueryGas bills the caller's Cosmos gas meter for an EVM query's
// computation whether or not the query succeeded, so that failing is never a
// cheaper way to buy validator CPU than succeeding.
//
// Both branches are deterministic across nodes, which matters because this is
// consensus-affecting:
//
//   - response.GasUsed is produced by ApplyMessageWithConfig from the message
//     gas limit, the EVM's leftover gas, the state refund counter and the
//     feemarket MinGasMultiplier — all consensus state. cosmos/evm populates it
//     on the revert path exactly as it does on success, and returns the response
//     alongside the error (x/vm/keeper/call_evm.go), so a revert that burned its
//     limit reports having burned its limit.
//   - When the call fails before the VM reports anything (intrinsic gas, gas
//     overflow, a sequence lookup failure) there is no measurement to charge, so
//     the declared limit is charged instead. That value comes from the
//     collection's own stored configuration, so every node computes the same
//     number.
func chargeEVMQueryGas(ctx sdk.Context, response *evmtypes.MsgEthereumTxResponse, callErr error, gasLimit uint64) {
	if response != nil && response.GasUsed > 0 {
		ctx.GasMeter().ConsumeGas(response.GasUsed, "evm_query_challenge")
		return
	}
	if callErr == nil && response != nil && response.VmError == "" {
		// A successful call that genuinely measured no gas: nothing to charge.
		return
	}
	ctx.GasMeter().ConsumeGas(gasLimit, "evm_query_challenge_failed")
}
