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
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	cmn "github.com/cosmos/evm/precompiles/common"

	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/pkg/evmcompat"
	sendmanagerkeeper "github.com/bitbadges/bitbadgeschain/x/sendmanager/keeper"
	sendmanagertypes "github.com/bitbadges/bitbadgeschain/x/sendmanager/types"
)

const (
	// Base gas costs for transactions
	// IMPORTANT: These values are DEDUCTED from the transaction gas before the precompile runs.
	// The actual execution gas comes from the remaining gas (contract.Gas after deduction).
	// Setting these too high causes "out of gas" errors. Keep these as minimal entry fees.
	GasSendBase = 30_000

	// Gas costs per element for dynamic calculations
	GasPerCoin       = 2_000
	GasPerInputChunk = 100 // Gas per 32-byte chunk of input, for JSON parsing
)

var _ vm.PrecompiledContract = &Precompile{}

// Name implements vm.PrecompiledContract (geth 1.17).
func (Precompile) Name() string {
	return "sendmanager"
}

var (
	// Embed abi json file to the executable binary. Needed when importing as dependency.
	//
	//go:embed abi.json
	f   []byte
	ABI abi.ABI
	// abiLoadError stores any error from ABI loading for lazy error reporting
	abiLoadError error
)

func init() {
	// abi.json is a Hardhat artifact ({_format, contractName, abi, ...}),
	// not a bare ABI array. cosmos/evm v0.6's cmn.LoadABI unwrapped it; v0.7
	// removed that helper and its replacement rejects contracts with no
	// bytecode, which precompile interfaces always are. So unwrap it here.
	var artifact struct {
		ABI json.RawMessage `json:"abi"`
	}
	if abiLoadError = json.Unmarshal(f, &artifact); abiLoadError == nil {
		ABI, abiLoadError = abi.JSON(bytes.NewReader(artifact.ABI))
	}
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
//
// bankKeeper feeds the balance handler that replays the bank events emitted
// by the keeper into the EVM StateDB, so that a balance the precompile moved
// is not overwritten by the StateDB at commit. It may be nil only for unit
// tests that never run inside the EVM.
func NewPrecompile(
	sendManagerKeeper sendmanagerkeeper.Keeper,
	bankKeeper cmn.BankKeeper,
) *Precompile {
	p := &Precompile{
		Precompile: cmn.Precompile{
			KvGasConfig:          storetypes.KVGasConfig(),
			TransientKVGasConfig: storetypes.TransientGasConfig(),
			ContractAddress:      common.HexToAddress(SendManagerPrecompileAddress),
		},
		ABI:               ABI,
		sendManagerKeeper: sendManagerKeeper,
	}
	if bankKeeper != nil {
		p.BalanceHandlerFactory = evmcompat.NewBalanceHandlerFactory(bankKeeper)
	}
	return p
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
// Returns a conservative estimate that accounts for Cosmos SDK operations to help
// estimateGas converge on a working value.
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

	// Get base gas and add buffer for Cosmos SDK operations
	// Send operations involve bank transfers which need significant gas
	baseGas := p.getBaseGas(method.Name)
	if baseGas == 0 {
		return 0
	}

	// Priced by input size: the JSON has to be parsed and validated before
	// the per-coin charge in HandleTransaction can apply.
	baseGas += uint64(len(input)/32) * GasPerInputChunk

	// Add buffer for Cosmos SDK state operations (bank transfers, etc.)
	return baseGas + 150_000
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

	return evmcompat.RunNativeAction(p.Precompile, evm, contract, func(ctx sdk.Context) ([]byte, error) {
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
		if err := ValidateArraySize(len(m.Amount), MaxCoins, "amount"); err != nil {
			return nil, err
		}
		ctx.GasMeter().ConsumeGas(uint64(len(m.Amount))*GasPerCoin, "precompile message elements")
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
