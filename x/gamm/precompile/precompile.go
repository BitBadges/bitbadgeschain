// Package gamm implements a precompiled contract for the BitBadges gamm module.
// This precompile enables Solidity smart contracts to interact with liquidity pools
// through a standardized EVM interface.
//
// The precompile is available at address 0x0000000000000000000000000000000000001002 and provides
// both transaction methods (state-changing operations) and query methods (read-only operations).
//
// Transaction Methods:
//   - joinPool: Join a liquidity pool by providing tokens
//   - exitPool: Exit a liquidity pool by burning shares
//   - swapExactAmountIn: Swap tokens with exact input amount
//   - swapExactAmountInWithIBCTransfer: Swap tokens and transfer via IBC
//
// Query Methods:
//   - getPool: Query pool data by ID
//   - getPools: Query all pools with pagination
//   - calcJoinPoolNoSwapShares: Calculate shares for joining pool without swap
//   - calcExitPoolCoinsFromShares: Calculate tokens received for exiting pool
//   - And more query methods for pool information
//
// All methods use structured error handling with error codes for consistent error reporting.
// Input validation is performed on all parameters to ensure security and correctness.
package gamm

import (
	"embed"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	cmn "github.com/cosmos/evm/precompiles/common"

	"github.com/cosmos/gogoproto/proto"

	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	gammkeeper "github.com/bitbadges/bitbadgeschain/x/gamm/keeper"
	gammtypes "github.com/bitbadges/bitbadgeschain/x/gamm/types"
)

var _ vm.PrecompiledContract = &Precompile{}

var (
	// Embed abi json file to the executable binary. Needed when importing as dependency.
	//
	//go:embed abi.json
	f   embed.FS
	ABI abi.ABI
	// abiLoadError stores any error from ABI loading for lazy error reporting
	abiLoadError error
)

func init() {
	ABI, abiLoadError = cmn.LoadABI(f, "abi.json")
	if abiLoadError != nil {
		// Log the error but don't panic - the error will be returned when the precompile is used
		// This allows the chain to start even if the ABI is malformed, but the precompile will be disabled
		fmt.Printf("WARNING: Failed to load gamm precompile ABI: %v\n", abiLoadError)
	}
}

// GetABILoadError returns any error that occurred during ABI loading
// This can be checked by callers to verify the precompile is properly initialized
func GetABILoadError() error {
	return abiLoadError
}

// ValidatePrecompileEnabled validates that the precompile is properly configured and enabled.
// This checks:
//   - ABI loaded successfully
//   - Precompile address is valid
//   - Basic structure is correct
//
// Note: This cannot verify that the precompile is enabled in the EVM keeper without
// access to the keeper. For production, ensure the precompile address is in the
// genesis state's active_static_precompiles array.
func ValidatePrecompileEnabled() error {
	// Check ABI loading
	if abiLoadError != nil {
		return fmt.Errorf("gamm precompile ABI failed to load: %w", abiLoadError)
	}

	// Check ABI is valid
	if len(ABI.Methods) == 0 {
		return fmt.Errorf("gamm precompile ABI has no methods")
	}

	// Check precompile address is valid
	precompileAddr := common.HexToAddress(GammPrecompileAddress)
	if precompileAddr == (common.Address{}) {
		return fmt.Errorf("gamm precompile address is zero")
	}

	// Verify required methods exist
	requiredMethods := []string{
		JoinPoolMethod,
		ExitPoolMethod,
		SwapExactAmountInMethod,
		GetPoolMethod,
	}
	for _, methodName := range requiredMethods {
		if _, found := ABI.Methods[methodName]; !found {
			return fmt.Errorf("required method %s not found in ABI", methodName)
		}
	}

	return nil
}

// Precompile defines the gamm precompile
type Precompile struct {
	cmn.Precompile

	abi.ABI
	gammKeeper gammkeeper.Keeper
}

// NewPrecompile creates a new gamm Precompile instance implementing the
// PrecompiledContract interface.
func NewPrecompile(
	gammKeeper gammkeeper.Keeper,
) *Precompile {
	return &Precompile{
		Precompile: cmn.Precompile{
			KvGasConfig:          storetypes.GasConfig{},
			TransientKVGasConfig: storetypes.GasConfig{},
			ContractAddress:      common.HexToAddress(GammPrecompileAddress),
		},
		ABI:        ABI,
		gammKeeper: gammKeeper,
	}
}

// GammPrecompileAddress is the address of the gamm precompile
// Using standard precompile address range: 0x0000000000000000000000000000000000001002
const GammPrecompileAddress = "0x0000000000000000000000000000000000001002"

// GetCallerAddress gets the caller address and converts it to Cosmos format
// This should be used for ALL transaction methods to set the Sender field
// SECURITY: This ensures the sender is always the actual caller, preventing impersonation
// The caller is obtained from contract.Caller() which returns the EVM msg.sender
// and cannot be spoofed by malicious contracts
func (p Precompile) GetCallerAddress(contract *vm.Contract) (string, error) {
	caller := contract.Caller()
	if err := VerifyCaller(caller); err != nil {
		return "", err
	}
	return sdk.AccAddress(caller.Bytes()).String(), nil
}

// RequiredGas calculates the precompiled contract's base gas rate.
// For methods that use JSON string arguments, we return base gas since
// parsing JSON for gas calculation is complex and not worth the overhead.
func (p Precompile) RequiredGas(input []byte) uint64 {
	// NOTE: This check avoid panicking when trying to decode the method ID
	if len(input) < 4 {
		return 0
	}

	methodID := input[:4]

	method, err := p.MethodById(methodID)
	if err != nil {
		// This should never happen since this method is going to fail during Run
		return 0
	}

	// All methods now use JSON string arguments, so return base gas
	// Dynamic gas calculation would require parsing JSON which is not worth the overhead
	return p.getBaseGas(method.Name)
}

// getBaseGas returns the base gas cost for a method
func (p Precompile) getBaseGas(methodName string) uint64 {
	switch methodName {
	// Transaction methods
	case JoinPoolMethod:
		return GasJoinPoolBase
	case ExitPoolMethod:
		return GasExitPoolBase
	case SwapExactAmountInMethod:
		return GasSwapExactAmountInBase
	case SwapExactAmountInWithIBCTransferMethod:
		return GasSwapExactAmountInWithIBCTransferBase
	// Query methods
	case GetPoolMethod:
		return GasGetPoolBase
	case GetPoolsMethod:
		return GasGetPoolsBase
	case GetPoolTypeMethod:
		return GasGetPoolTypeBase
	case CalcJoinPoolNoSwapSharesMethod:
		return GasCalcJoinPoolNoSwapSharesBase
	case CalcExitPoolCoinsFromSharesMethod:
		return GasCalcExitPoolCoinsFromSharesBase
	case CalcJoinPoolSharesMethod:
		return GasCalcJoinPoolSharesBase
	case GetPoolParamsMethod:
		return GasGetPoolParamsBase
	case GetTotalSharesMethod:
		return GasGetTotalSharesBase
	case GetTotalLiquidityMethod:
		return GasGetTotalLiquidityBase
	}
	return 0
}

func (p Precompile) Run(evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	// Check if ABI loaded successfully during init
	if abiLoadError != nil {
		return nil, fmt.Errorf("gamm precompile unavailable: ABI failed to load: %w", abiLoadError)
	}

	return p.RunNativeAction(evm, contract, func(ctx sdk.Context) ([]byte, error) {
		result, methodName, err := p.ExecuteWithMethodName(ctx, contract, readonly)

		// Gas is tracked by the EVM, we log the method for monitoring
		LogPrecompileUsage(ctx, methodName, err == nil, 0, err)

		return result, err
	})
}

// Execute executes the precompiled contract gamm methods defined in the ABI.
// Deprecated: Use ExecuteWithMethodName instead for better performance (avoids double method lookup).
func (p Precompile) Execute(ctx sdk.Context, contract *vm.Contract, readOnly bool) ([]byte, error) {
	bz, _, err := p.ExecuteWithMethodName(ctx, contract, readOnly)
	return bz, err
}

// ExecuteWithMethodName executes the precompiled contract and returns the method name for logging.
// This avoids the double MethodById() lookup that occurs when logging separately.
func (p Precompile) ExecuteWithMethodName(ctx sdk.Context, contract *vm.Contract, readOnly bool) ([]byte, string, error) {
	method, args, err := cmn.SetupABI(p.ABI, contract, readOnly, p.IsTransaction)
	if err != nil {
		return nil, "unknown", err
	}

	// Extract JSON string from args
	if len(args) != 1 {
		return nil, method.Name, ErrInvalidInput(fmt.Sprintf("expected 1 argument (JSON string), got %d", len(args)))
	}

	jsonStr, ok := args[0].(string)
	if !ok {
		return nil, method.Name, ErrInvalidInput("expected JSON string as first argument")
	}

	// Route to transaction or query handler
	var bz []byte
	if p.IsTransaction(method) {
		bz, err = p.HandleTransaction(ctx, method, jsonStr, contract)
	} else {
		bz, err = p.HandleQuery(ctx, method, jsonStr)
	}

	return bz, method.Name, err
}

// HandleTransaction handles a transaction by unmarshaling JSON and executing via keeper
func (p Precompile) HandleTransaction(ctx sdk.Context, method *abi.Method, jsonStr string, contract *vm.Contract) ([]byte, error) {
	// Unmarshal JSON to Msg
	msg, err := p.unmarshalMsgFromJSON(method.Name, jsonStr, contract)
	if err != nil {
		return nil, err
	}

	// Execute message via keeper
	msgServer := gammkeeper.NewMsgServerImpl(&p.gammKeeper)

	// Route to appropriate handler based on message type
	var resp interface{}
	switch m := msg.(type) {
	case *gammtypes.MsgJoinPool:
		resp, err = msgServer.JoinPool(ctx, m)
	case *gammtypes.MsgExitPool:
		resp, err = msgServer.ExitPool(ctx, m)
	case *gammtypes.MsgSwapExactAmountIn:
		resp, err = msgServer.SwapExactAmountIn(ctx, m)
	case *gammtypes.MsgSwapExactAmountInWithIBCTransfer:
		resp, err = msgServer.SwapExactAmountInWithIBCTransfer(ctx, m)
	default:
		return nil, ErrInvalidInput(fmt.Sprintf("unsupported message type for method: %s", method.Name))
	}

	if err != nil {
		return nil, WrapError(err, ErrorCodeSwapFailed, fmt.Sprintf("transaction failed: %v", err))
	}

	// Pack response based on method output type
	switch method.Name {
	case JoinPoolMethod:
		if joinResp, ok := resp.(*gammtypes.MsgJoinPoolResponse); ok {
			return method.Outputs.Pack(joinResp.ShareOutAmount.BigInt(), ConvertCoinsToEVM(joinResp.TokenIn))
		}
		return nil, WrapError(fmt.Errorf("invalid response type for joinPool"), ErrorCodeInternalError, "expected MsgJoinPoolResponse")
	case ExitPoolMethod:
		if exitResp, ok := resp.(*gammtypes.MsgExitPoolResponse); ok {
			return method.Outputs.Pack(ConvertCoinsToEVM(exitResp.TokenOut))
		}
		return nil, WrapError(fmt.Errorf("invalid response type for exitPool"), ErrorCodeInternalError, "expected MsgExitPoolResponse")
	case SwapExactAmountInMethod, SwapExactAmountInWithIBCTransferMethod:
		if swapResp, ok := resp.(*gammtypes.MsgSwapExactAmountInResponse); ok {
			return method.Outputs.Pack(swapResp.TokenOutAmount.BigInt())
		}
		return nil, WrapError(fmt.Errorf("invalid response type for swapExactAmountIn"), ErrorCodeInternalError, "expected MsgSwapExactAmountInResponse")
	default:
		return nil, WrapError(fmt.Errorf("unsupported transaction method: %s", method.Name), ErrorCodeInternalError, "method not handled in response packing")
	}
}

// HandleQuery handles a query by unmarshaling JSON and executing via keeper
func (p Precompile) HandleQuery(ctx sdk.Context, method *abi.Method, jsonStr string) ([]byte, error) {
	// Unmarshal JSON to QueryRequest
	queryReq, err := p.unmarshalQueryFromJSON(method.Name, jsonStr)
	if err != nil {
		return nil, err
	}

	// Execute query via keeper querier
	querier := gammkeeper.NewQuerier(p.gammKeeper)
	var resp interface{}

	switch req := queryReq.(type) {
	case *gammtypes.QueryPoolRequest:
		resp, err = querier.Pool(ctx, req)
	case *gammtypes.QueryPoolsRequest:
		resp, err = querier.Pools(ctx, req)
	case *gammtypes.QueryPoolTypeRequest:
		resp, err = querier.PoolType(ctx, req)
	case *gammtypes.QueryCalcJoinPoolNoSwapSharesRequest:
		resp, err = querier.CalcJoinPoolNoSwapShares(ctx, req)
	case *gammtypes.QueryCalcExitPoolCoinsFromSharesRequest:
		resp, err = querier.CalcExitPoolCoinsFromShares(ctx, req)
	case *gammtypes.QueryCalcJoinPoolSharesRequest:
		resp, err = querier.CalcJoinPoolShares(ctx, req)
	case *gammtypes.QueryPoolParamsRequest:
		resp, err = querier.PoolParams(ctx, req)
	case *gammtypes.QueryTotalSharesRequest:
		resp, err = querier.TotalShares(ctx, req)
	case *gammtypes.QueryTotalLiquidityRequest:
		resp, err = querier.TotalLiquidity(ctx, req)
	default:
		return nil, ErrInvalidInput(fmt.Sprintf("unsupported query type for method: %s", method.Name))
	}

	if err != nil {
		return nil, WrapError(err, ErrorCodeQueryFailed, fmt.Sprintf("query failed: %v", err))
	}

	// Handle different return types
	switch method.Name {
	case GetPoolTypeMethod:
		if poolTypeResp, ok := resp.(*gammtypes.QueryPoolTypeResponse); ok {
			return method.Outputs.Pack(poolTypeResp.PoolType)
		}
	case GetTotalSharesMethod:
		// ABI expects struct Coin (tuple), not bytes
		if totalSharesResp, ok := resp.(*gammtypes.QueryTotalSharesResponse); ok {
			coinStruct := ConvertCoinToEVM(totalSharesResp.TotalShares)
			return method.Outputs.Pack(coinStruct)
		}
		return nil, WrapError(fmt.Errorf("response is not QueryTotalSharesResponse"), ErrorCodeInternalError, "invalid response type")
	case GetTotalLiquidityMethod:
		// ABI expects struct Coin[] (tuple[]), not bytes
		if totalLiquidityResp, ok := resp.(*gammtypes.QueryTotalLiquidityResponse); ok {
			coinsStruct := ConvertCoinsToEVM(totalLiquidityResp.Liquidity)
			return method.Outputs.Pack(coinsStruct)
		}
		return nil, WrapError(fmt.Errorf("response is not QueryTotalLiquidityResponse"), ErrorCodeInternalError, "invalid response type")
	case CalcJoinPoolNoSwapSharesMethod:
		// ABI expects (tokensOut tuple[], sharesOut uint256), not bytes
		if calcResp, ok := resp.(*gammtypes.QueryCalcJoinPoolNoSwapSharesResponse); ok {
			tokensOutStruct := ConvertCoinsToEVM(calcResp.TokensOut)
			sharesOutBigInt := calcResp.SharesOut.BigInt()
			return method.Outputs.Pack(tokensOutStruct, sharesOutBigInt)
		}
		return nil, WrapError(fmt.Errorf("response is not QueryCalcJoinPoolNoSwapSharesResponse"), ErrorCodeInternalError, "invalid response type")
	case CalcExitPoolCoinsFromSharesMethod:
		// ABI expects tokensOut tuple[], not bytes
		if calcResp, ok := resp.(*gammtypes.QueryCalcExitPoolCoinsFromSharesResponse); ok {
			tokensOutStruct := ConvertCoinsToEVM(calcResp.TokensOut)
			return method.Outputs.Pack(tokensOutStruct)
		}
		return nil, WrapError(fmt.Errorf("response is not QueryCalcExitPoolCoinsFromSharesResponse"), ErrorCodeInternalError, "invalid response type")
	case CalcJoinPoolSharesMethod:
		// ABI expects (shareOutAmount uint256, tokensOut tuple[]), not bytes
		if calcResp, ok := resp.(*gammtypes.QueryCalcJoinPoolSharesResponse); ok {
			shareOutBigInt := calcResp.ShareOutAmount.BigInt()
			tokensOutStruct := ConvertCoinsToEVM(calcResp.TokensOut)
			return method.Outputs.Pack(shareOutBigInt, tokensOutStruct)
		}
		return nil, WrapError(fmt.Errorf("response is not QueryCalcJoinPoolSharesResponse"), ErrorCodeInternalError, "invalid response type")
	case GetPoolMethod, GetPoolsMethod, GetPoolParamsMethod:
		// Marshal response to bytes (protobuf)
		if protoMsg, ok := resp.(proto.Message); ok {
			bz, err := proto.Marshal(protoMsg)
			if err != nil {
				return nil, WrapError(err, ErrorCodeInternalError, "failed to marshal query response")
			}
			return method.Outputs.Pack(bz)
		}
		return nil, WrapError(fmt.Errorf("response is not a proto.Message"), ErrorCodeInternalError, "invalid response type")
	}

	// Default: marshal to bytes
	if protoMsg, ok := resp.(proto.Message); ok {
		bz, err := proto.Marshal(protoMsg)
		if err != nil {
			return nil, WrapError(err, ErrorCodeInternalError, "failed to marshal query response")
		}
		return method.Outputs.Pack(bz)
	}
	return nil, WrapError(fmt.Errorf("response is not a proto.Message"), ErrorCodeInternalError, "invalid response type")
}

// transactionMethods is a map of method names that are transactions (state-changing).
// Using a map provides O(1) lookup instead of O(n) switch statement.
var transactionMethods = map[string]bool{
	JoinPoolMethod:                         true,
	ExitPoolMethod:                         true,
	SwapExactAmountInMethod:                true,
	SwapExactAmountInWithIBCTransferMethod: true,
}

// IsTransaction checks if the given method name corresponds to a transaction or query.
// Uses O(1) map lookup for better performance.
func (Precompile) IsTransaction(method *abi.Method) bool {
	return transactionMethods[method.Name]
}

// Method name constants
const (
	// Transaction methods
	JoinPoolMethod                         = "joinPool"
	ExitPoolMethod                         = "exitPool"
	SwapExactAmountInMethod                = "swapExactAmountIn"
	SwapExactAmountInWithIBCTransferMethod = "swapExactAmountInWithIBCTransfer"

	// Query methods
	GetPoolMethod                     = "getPool"
	GetPoolsMethod                    = "getPools"
	GetPoolTypeMethod                 = "getPoolType"
	CalcJoinPoolNoSwapSharesMethod    = "calcJoinPoolNoSwapShares"
	CalcExitPoolCoinsFromSharesMethod = "calcExitPoolCoinsFromShares"
	CalcJoinPoolSharesMethod          = "calcJoinPoolShares"
	GetPoolParamsMethod               = "getPoolParams"
	GetTotalSharesMethod              = "getTotalShares"
	GetTotalLiquidityMethod           = "getTotalLiquidity"
)

// LogPrecompileUsage logs precompile usage for monitoring
func LogPrecompileUsage(ctx sdk.Context, method string, success bool, gasUsed uint64, err error) {
	logger := ctx.Logger()

	if err != nil {
		// Extract error code if it's a PrecompileError
		if precompileErr, ok := err.(*PrecompileError); ok {
			logger.Error("precompile error",
				"method", method,
				"error_code", precompileErr.Code,
				"error_message", precompileErr.Message,
				"gas_used", gasUsed,
			)
		} else {
			logger.Error("precompile error",
				"method", method,
				"error", err.Error(),
				"gas_used", gasUsed,
			)
		}
	} else {
		logger.Info("precompile success",
			"method", method,
			"gas_used", gasUsed,
		)
	}
}
