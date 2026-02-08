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

	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	gammkeeper "github.com/bitbadges/bitbadgeschain/x/gamm/keeper"
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
// For methods with dynamic inputs (arrays), it attempts to parse the input
// to calculate accurate gas costs. If parsing fails, it falls back to base gas.
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

	// Try to unpack arguments to calculate dynamic gas
	// If unpacking fails, fall back to base gas
	argsData := input[4:]
	args, err := method.Inputs.Unpack(argsData)
	if err != nil {
		// Parsing failed, return base gas as fallback
		return p.getBaseGas(method.Name)
	}

	// Calculate dynamic gas based on method and arguments
	return p.calculateDynamicGas(method.Name, args)
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

// calculateDynamicGas calculates gas based on method name and unpacked arguments
func (p Precompile) calculateDynamicGas(methodName string, args []interface{}) uint64 {
	baseGas := p.getBaseGas(methodName)
	if baseGas == 0 {
		return 0
	}

	switch methodName {
	case JoinPoolMethod:
		// args: poolId, shareOutAmount, tokenInMaxs[]
		if len(args) >= 3 {
			if tokenInMaxs, ok := args[2].([]interface{}); ok {
				return CalculateDynamicGas(baseGas, 0, len(tokenInMaxs), 0)
			}
		}
		return baseGas

	case ExitPoolMethod:
		// args: poolId, shareInAmount, tokenOutMins[]
		if len(args) >= 3 {
			if tokenOutMins, ok := args[2].([]interface{}); ok {
				return CalculateDynamicGas(baseGas, 0, len(tokenOutMins), 0)
			}
		}
		return baseGas

	case SwapExactAmountInMethod:
		// args: routes[], tokenIn, tokenOutMinAmount, affiliates[]
		numRoutes := 0
		numAffiliates := 0
		if len(args) >= 1 {
			if routes, ok := args[0].([]interface{}); ok {
				numRoutes = len(routes)
			}
		}
		if len(args) >= 4 {
			if affiliates, ok := args[3].([]interface{}); ok {
				numAffiliates = len(affiliates)
			}
		}
		return CalculateDynamicGas(baseGas, numRoutes, 0, numAffiliates)

	case SwapExactAmountInWithIBCTransferMethod:
		// args: routes[], tokenIn, tokenOutMinAmount, ibcTransferInfo, affiliates[]
		numRoutes := 0
		numAffiliates := 0
		memoLength := 0
		if len(args) >= 1 {
			if routes, ok := args[0].([]interface{}); ok {
				numRoutes = len(routes)
			}
		}
		if len(args) >= 3 {
			if ibcInfo, ok := args[3].(map[string]interface{}); ok {
				if memo, ok := ibcInfo["memo"].(string); ok {
					memoLength = len(memo)
				}
			}
		}
		if len(args) >= 5 {
			if affiliates, ok := args[4].([]interface{}); ok {
				numAffiliates = len(affiliates)
			}
		}
		gas := CalculateDynamicGas(baseGas, numRoutes, 0, numAffiliates)
		// Add gas for memo bytes
		gas += uint64(memoLength) * GasPerMemoByte
		return gas

	case CalcJoinPoolNoSwapSharesMethod, CalcJoinPoolSharesMethod:
		// args: poolId, tokensIn[]
		if len(args) >= 2 {
			if tokensIn, ok := args[1].([]interface{}); ok {
				return CalculateDynamicGas(baseGas, 0, len(tokensIn), 0)
			}
		}
		return baseGas

	default:
		// For other methods, return base gas
		return baseGas
	}
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

	var bz []byte
	switch method.Name {
	// Transactions
	case JoinPoolMethod:
		bz, err = p.JoinPool(ctx, method, args, contract)
	case ExitPoolMethod:
		bz, err = p.ExitPool(ctx, method, args, contract)
	case SwapExactAmountInMethod:
		bz, err = p.SwapExactAmountIn(ctx, method, args, contract)
	case SwapExactAmountInWithIBCTransferMethod:
		bz, err = p.SwapExactAmountInWithIBCTransfer(ctx, method, args, contract)
	// Queries
	case GetPoolMethod:
		bz, err = p.GetPool(ctx, method, args)
	case GetPoolsMethod:
		bz, err = p.GetPools(ctx, method, args)
	case GetPoolTypeMethod:
		bz, err = p.GetPoolType(ctx, method, args)
	case CalcJoinPoolNoSwapSharesMethod:
		bz, err = p.CalcJoinPoolNoSwapShares(ctx, method, args)
	case CalcExitPoolCoinsFromSharesMethod:
		bz, err = p.CalcExitPoolCoinsFromShares(ctx, method, args)
	case CalcJoinPoolSharesMethod:
		bz, err = p.CalcJoinPoolShares(ctx, method, args)
	case GetPoolParamsMethod:
		bz, err = p.GetPoolParams(ctx, method, args)
	case GetTotalSharesMethod:
		bz, err = p.GetTotalShares(ctx, method, args)
	case GetTotalLiquidityMethod:
		bz, err = p.GetTotalLiquidity(ctx, method, args)
	default:
		return nil, method.Name, fmt.Errorf(cmn.ErrUnknownMethod, method.Name)
	}

	return bz, method.Name, err
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
