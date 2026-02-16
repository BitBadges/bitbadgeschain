// Package precompile implements a precompiled contract for the sendmanager module.
// This precompile enables Solidity smart contracts to execute MsgSendWithAliasRouting transactions
// without requiring ERC20 wrapping, keeping all accounting in x/bank (Cosmos side).
//
// The precompile is available at address 0x0000000000000000000000000000000000001003 and provides
// transaction methods for sending coins with alias denom routing support.
//
// Transaction Methods:
//   - send: Send native Cosmos coins from the caller to a recipient (supports alias denoms)
//
// All methods use structured error handling with error codes for consistent error reporting.
// Input validation is performed on all parameters to ensure security and correctness.
package precompile

import (
	"embed"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	cmn "github.com/cosmos/evm/precompiles/common"

	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	sendmanagerkeeper "github.com/bitbadges/bitbadgeschain/x/sendmanager/keeper"
	sendmanagertypes "github.com/bitbadges/bitbadgeschain/x/sendmanager/types"
)

const (
	// Base gas costs for transactions
	// IMPORTANT: These values must account for bank module operations and state changes.
	// If these values are too low, gas estimation will succeed but actual execution
	// may fail with "out of gas" errors that produce missing revert data.
	GasSendBase = 100_000 // Send involves balance checks, transfers, and state updates

	// Gas costs per element for dynamic calculations
	GasPerCoin = 10_000 // Each coin involves balance lookup and transfer
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
		fmt.Printf("WARNING: Failed to load sendmanager precompile ABI: %v\n", abiLoadError)
	}
}

// GetABILoadError returns any error that occurred during ABI loading
// This can be checked by callers to verify the precompile is properly initialized
func GetABILoadError() error {
	return abiLoadError
}

// Precompile defines the sendmanager precompile
type Precompile struct {
	cmn.Precompile

	abi.ABI
	sendManagerKeeper sendmanagerkeeper.Keeper
}

// NewPrecompile creates a new sendmanager Precompile instance implementing the
// PrecompiledContract interface.
func NewPrecompile(
	sendManagerKeeper sendmanagerkeeper.Keeper,
) *Precompile {
	return &Precompile{
		Precompile: cmn.Precompile{
			KvGasConfig:          storetypes.GasConfig{},
			TransientKVGasConfig: storetypes.GasConfig{},
			ContractAddress:      common.HexToAddress(SendManagerPrecompileAddress),
		},
		ABI:               ABI,
		sendManagerKeeper: sendManagerKeeper,
	}
}

// SendManagerPrecompileAddress is the address of the sendmanager precompile
// Using standard precompile address range: 0x0000000000000000000000000000000000001003
const SendManagerPrecompileAddress = "0x0000000000000000000000000000000000001003"

// SendMethod is the name of the send method in the ABI
const SendMethod = "send"

// GetCallerAddress gets the caller address and converts it to Cosmos format
// This should be used for ALL transaction methods to set the from_address field
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
	case SendMethod:
		return GasSendBase
	}
	return 0
}

func (p Precompile) Run(evm *vm.EVM, contract *vm.Contract, readonly bool) ([]byte, error) {
	// Check if ABI loaded successfully during init
	if abiLoadError != nil {
		return nil, fmt.Errorf("sendmanager precompile unavailable: ABI failed to load: %w", abiLoadError)
	}

	// Sendmanager precompile only supports transactions, not queries
	if readonly {
		return nil, fmt.Errorf("sendmanager precompile does not support read-only operations")
	}

	return p.RunNativeAction(evm, contract, func(ctx sdk.Context) ([]byte, error) {
		result, methodName, err := p.ExecuteWithMethodName(ctx, contract, readonly)

		// Gas is tracked by the EVM, we log the method for monitoring
		LogPrecompileUsage(ctx, methodName, err == nil, 0, err)

		return result, err
	})
}

// Execute executes the precompiled contract sendmanager methods defined in the ABI.
// This is a convenience method for testing that wraps ExecuteWithMethodName.
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

	// Route to transaction handler
	var bz []byte
	if p.IsTransaction(method) {
		bz, err = p.HandleTransaction(ctx, method, jsonStr, contract)
	} else {
		return nil, method.Name, ErrInvalidInput(fmt.Sprintf("method %s is not a transaction", method.Name))
	}

	return bz, method.Name, err
}

// IsTransaction returns true if the method is a transaction (state-changing operation)
func (p Precompile) IsTransaction(method *abi.Method) bool {
	switch method.Name {
	case SendMethod:
		return true
	}
	return false
}

// HandleTransaction handles a transaction by unmarshaling JSON and executing via keeper
func (p Precompile) HandleTransaction(ctx sdk.Context, method *abi.Method, jsonStr string, contract *vm.Contract) ([]byte, error) {
	// Unmarshal JSON to Msg
	msg, err := p.unmarshalMsgFromJSON(method.Name, jsonStr, contract)
	if err != nil {
		return nil, err
	}

	// Execute message via keeper
	msgServer := sendmanagerkeeper.NewMsgServerImpl(p.sendManagerKeeper)

	// Route to appropriate handler based on message type
	var resp interface{}
	switch m := msg.(type) {
	case *sendmanagertypes.MsgSendWithAliasRouting:
		_, err = msgServer.SendWithAliasRouting(ctx, m)
		resp = true // ABI: bool success
	default:
		return nil, ErrInvalidInput(fmt.Sprintf("unsupported message type for method: %s", method.Name))
	}

	if err != nil {
		return nil, WrapError(err, ErrorCodeSendFailed, "send failed")
	}

	// Pack response based on method output type
	switch method.Name {
	case SendMethod:
		// Pack bool success
		return method.Outputs.Pack(resp)
	default:
		return nil, ErrInvalidInput(fmt.Sprintf("unknown method: %s", method.Name))
	}
}
