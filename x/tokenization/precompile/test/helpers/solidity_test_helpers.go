package helpers

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	evmkeeper "github.com/cosmos/evm/x/vm/keeper"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	tokenization "github.com/bitbadges/bitbadgeschain/x/tokenization/precompile"
)

// ContractType represents the type of test contract to load
type ContractType string

const (
	// MinimalTestContract - minimal contract for basic testing
	ContractTypeMinimal ContractType = "MinimalTestContract"
	// HelperLibrariesTestContract - comprehensive test contract for helper libraries
	ContractTypeHelperLibraries ContractType = "HelperLibrariesTestContract"
	// GammHelperLibrariesTestContract - test contract for GAMM helper libraries
	ContractTypeGammHelperLibraries ContractType = "GammHelperLibrariesTestContract"
	// PrecompileTransferTestContract - transfer and query methods
	ContractTypePrecompileTransfer ContractType = "PrecompileTransferTestContract"
	// PrecompileCollectionTestContract - collection management methods
	ContractTypePrecompileCollection ContractType = "PrecompileCollectionTestContract"
	// PrecompileDynamicStoreTestContract - dynamic store methods
	ContractTypePrecompileDynamicStore ContractType = "PrecompileDynamicStoreTestContract"
)

// Contract compilation paths
// These are relative to the repository root
const (
	// ContractSourcePath is the path to the Solidity contract source
	// Using MinimalTestContract for testing - it's smaller and stays under EVM size limits
	ContractSourcePath = "contracts/test/MinimalTestContract.sol"
	// ContractName is the name of the contract to compile
	ContractName = "MinimalTestContract"
	// CompiledBytecodePath is where solcjs outputs the bytecode file
	// solcjs creates files with pattern: test_ContractName_sol_ContractName.bin
	CompiledBytecodePath = "contracts/test/test_MinimalTestContract_sol_MinimalTestContract.bin"
	// CompiledABIPath is where solcjs outputs the ABI file
	// solcjs creates files with pattern: test_ContractName_sol_ContractName.abi
	CompiledABIPath = "contracts/test/test_MinimalTestContract_sol_MinimalTestContract.abi"
)

// getContractPaths returns the bytecode and ABI paths for a given contract type
func getContractPaths(contractType ContractType) (bytecodePath, abiPath string) {
	switch contractType {
	case ContractTypeHelperLibraries:
		return "contracts/test/test_HelperLibrariesTestContract_sol_HelperLibrariesTestContract.bin",
			"contracts/test/test_HelperLibrariesTestContract_sol_HelperLibrariesTestContract.abi"
	case ContractTypeGammHelperLibraries:
		return "contracts/test/test_GammHelperLibrariesTestContract_sol_GammHelperLibrariesTestContract.bin",
			"contracts/test/test_GammHelperLibrariesTestContract_sol_GammHelperLibrariesTestContract.abi"
	case ContractTypePrecompileTransfer:
		return "contracts/test/test_PrecompileTransferTestContract_sol_PrecompileTransferTestContract.bin",
			"contracts/test/test_PrecompileTransferTestContract_sol_PrecompileTransferTestContract.abi"
	case ContractTypePrecompileCollection:
		return "contracts/test/test_PrecompileCollectionTestContract_sol_PrecompileCollectionTestContract.bin",
			"contracts/test/test_PrecompileCollectionTestContract_sol_PrecompileCollectionTestContract.abi"
	case ContractTypePrecompileDynamicStore:
		return "contracts/test/test_PrecompileDynamicStoreTestContract_sol_PrecompileDynamicStoreTestContract.bin",
			"contracts/test/test_PrecompileDynamicStoreTestContract_sol_PrecompileDynamicStoreTestContract.abi"
	case ContractTypeMinimal:
		fallthrough
	default:
		return CompiledBytecodePath, CompiledABIPath
	}
}

// DeployContract deploys a Solidity contract via EVM transaction
// contractBytecode is the compiled bytecode of the contract
// Returns the deployed contract address and transaction response
func DeployContract(
	ctx sdk.Context,
	evmKeeper *evmkeeper.Keeper,
	deployerKey *ecdsa.PrivateKey,
	contractBytecode []byte,
	chainID *big.Int,
) (common.Address, *evmtypes.MsgEthereumTxResponse, error) {
	deployerAddr := crypto.PubkeyToAddress(deployerKey.PublicKey)
	nonce := evmKeeper.GetNonce(ctx, deployerAddr)

	// Build deployment transaction (to address is nil for contract creation)
	// Use zero gas price to avoid fee collector refund issues in tests
	tx, err := BuildEVMTransaction(
		deployerKey,
		nil, // nil for contract creation (not zero address!)
		contractBytecode,
		big.NewInt(0),
		10000000,      // Large gas limit for deployment
		big.NewInt(0), // Zero gas price for testing to avoid fee collector refund issues
		nonce,
		chainID,
	)
	if err != nil {
		return common.Address{}, nil, err
	}

	// Execute transaction
	response, err := ExecuteEVMTransaction(ctx, evmKeeper, tx)
	if err != nil {
		return common.Address{}, nil, err
	}

	// Check if transaction failed (but wasn't caught as an error)
	// If there's a VmError and it's not a snapshot error, the transaction likely failed
	if response.VmError != "" && !strings.Contains(response.VmError, "snapshot revert error") {
		// Transaction failed - contract wasn't deployed
		return common.Address{}, response, fmt.Errorf("contract deployment failed: %s", response.VmError)
	}

	// Check if nonce was incremented (indicates transaction was processed)
	newNonce := evmKeeper.GetNonce(ctx, deployerAddr)
	if newNonce <= nonce {
		// Nonce wasn't incremented, but if gas was used, the transaction was processed
		// In some test environments, nonces might not be incremented correctly
		// If gas was used, assume the transaction was processed
		if response.GasUsed == 0 {
			// No gas used - transaction definitely wasn't processed
			return common.Address{}, response, fmt.Errorf("contract deployment transaction was not processed (no gas used, nonce not incremented). VmError: %s", response.VmError)
		}
		// Gas was used but nonce wasn't incremented - this is unusual but might be a test environment quirk
		// Continue with contract address calculation
	}

	// Extract contract address from transaction receipt
	// The contract address is derived from deployer address and nonce
	// Note: We use the nonce BEFORE the transaction because CreateAddress expects the nonce
	// that was used when the transaction was created (the nonce that will be used for the transaction)
	contractAddr := crypto.CreateAddress(deployerAddr, nonce)

	// Verify the contract was actually deployed by checking if it has code
	// This helps catch cases where the transaction consumed gas but didn't deploy
	if !evmKeeper.IsContract(ctx, contractAddr) {
		// Contract doesn't have code - deployment likely failed
		// If gas was used, the transaction was processed but reverted
		// Return the address anyway so the caller can decide how to handle it
		// (Some test environments might have quirks where code isn't immediately available)
		if response.GasUsed > 0 {
			// Gas was used but no code - transaction reverted
			// Return address with error so caller can handle appropriately
			return contractAddr, response, fmt.Errorf("contract deployment reverted (gas used: %d but no code deployed). VmError: %s", response.GasUsed, response.VmError)
		}
		// No gas used - transaction wasn't processed
		return common.Address{}, response, fmt.Errorf("contract deployment transaction was not processed")
	}

	return contractAddr, response, nil
}

// CallContractMethod calls a method on a deployed contract
// contractAddr is the address of the deployed contract
// methodName is the name of the method to call
// contractABI is the ABI of the contract
// args are the arguments to pass to the method
// Returns the return data from the method call
func CallContractMethod(
	ctx sdk.Context,
	evmKeeper *evmkeeper.Keeper,
	callerKey *ecdsa.PrivateKey,
	contractAddr common.Address,
	contractABI abi.ABI,
	methodName string,
	args []interface{},
	chainID *big.Int,
	isView bool, // true for view functions (static calls), false for transactions
) ([]byte, *evmtypes.MsgEthereumTxResponse, error) {
	// Get method from ABI
	method, exists := contractABI.Methods[methodName]
	if !exists {
		return nil, nil, &tokenization.PrecompileError{
			Code:    tokenization.ErrorCodeInvalidInput,
			Message: "method not found",
			Details: fmt.Sprintf("method %s not found in contract ABI", methodName),
		}
	}

	// Pack method arguments
	packed, err := method.Inputs.Pack(args...)
	if err != nil {
		return nil, nil, err
	}

	// Create input with method ID
	input := append(method.ID, packed...)

	callerAddr := crypto.PubkeyToAddress(callerKey.PublicKey)

	if isView {
		// For view functions, use static call (no transaction needed)
		// This would require direct EVM call, which may not be available in keeper
		// For now, we'll use a transaction with zero value
		// TODO: Implement proper static call if EVM keeper supports it
		nonce := evmKeeper.GetNonce(ctx, callerAddr)
		tx, err := BuildEVMTransaction(
			callerKey,
			&contractAddr, // Pass pointer to address for calls
			input,
			big.NewInt(0),
			1000000,
			big.NewInt(0), // Zero gas price for testing to avoid fee collector refund issues
			nonce,
			chainID,
		)
		if err != nil {
			return nil, nil, err
		}

		response, err := ExecuteEVMTransaction(ctx, evmKeeper, tx)
		if err != nil {
			return nil, nil, err
		}

		// Unpack return data
		var returnData []byte
		if len(response.Ret) > 0 {
			returnData = response.Ret
		}

		// Unpack according to method outputs
		if len(method.Outputs) > 0 && len(returnData) > 0 {
			unpacked, err := method.Outputs.Unpack(returnData)
			if err != nil {
				return returnData, response, nil // Return raw data if unpacking fails
			}
			// Re-pack for consistency (caller can unpack as needed)
			repacked, _ := method.Outputs.Pack(unpacked...)
			return repacked, response, nil
		}

		return returnData, response, nil
	}

	// For non-view functions, execute as transaction
	// Use high gas limit to ensure precompile calls through contracts have enough gas
	nonce := evmKeeper.GetNonce(ctx, callerAddr)
	tx, err := BuildEVMTransaction(
		callerKey,
		&contractAddr, // Pass pointer to address for calls
		input,
		big.NewInt(0),
		5000000,       // High gas limit for contract->precompile calls
		big.NewInt(0), // Zero gas price for testing to avoid fee collector refund issues
		nonce,
		chainID,
	)
	if err != nil {
		return nil, nil, err
	}

	response, err := ExecuteEVMTransaction(ctx, evmKeeper, tx)
	if err != nil {
		return nil, nil, err
	}

	// Unpack return data if available
	var returnData []byte
	if len(response.Ret) > 0 {
		returnData = response.Ret
	}

	// Unpack according to method outputs
	if len(method.Outputs) > 0 && len(returnData) > 0 {
		unpacked, err := method.Outputs.Unpack(returnData)
		if err != nil {
			return returnData, response, nil // Return raw data if unpacking fails
		}
		// Re-pack for consistency
		repacked, _ := method.Outputs.Pack(unpacked...)
		return repacked, response, nil
	}

	return returnData, response, nil
}

// ParseContractEvents parses events from a contract transaction response
// contractABI is the ABI of the contract
// response is the transaction response containing logs
// Returns a map of event name to parsed event data
func ParseContractEvents(
	contractABI abi.ABI,
	response *evmtypes.MsgEthereumTxResponse,
) (map[string][]interface{}, error) {
	events := make(map[string][]interface{})

	// Parse logs from response
	// The response should contain logs that can be parsed using the ABI
	// Note: This is a simplified version - actual implementation may need
	// to handle log parsing differently based on EVM keeper's log format

	// For now, return empty map - actual event parsing will be implemented
	// when we have access to the log structure from evmtypes.MsgEthereumTxResponse
	// TODO: Implement proper event parsing based on EVM keeper's log format

	return events, nil
}

// VerifyContractDeployment verifies that a contract was deployed successfully
// by checking that the contract address has code
func VerifyContractDeployment(
	ctx sdk.Context,
	evmKeeper *evmkeeper.Keeper,
	contractAddr common.Address,
) (bool, error) {
	// Check if contract has code (is a contract)
	return evmKeeper.IsContract(ctx, contractAddr), nil
}

// findRepoRoot finds the repository root by looking for the contracts directory
func findRepoRoot() string {
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}

	// Start from current directory and walk up
	searchDir := cwd
	for {
		// Check if contracts directory exists
		if _, err := os.Stat(filepath.Join(searchDir, "contracts")); err == nil {
			return searchDir
		}

		// Move to parent directory
		parent := filepath.Dir(searchDir)
		if parent == searchDir {
			// Reached filesystem root
			break
		}
		searchDir = parent
	}

	// Fallback: return current directory
	return cwd
}

// GetContractABI returns the ABI for the MinimalTestContract (default)
// Loads from compiled ABI file (generated by `make compile-contracts`)
func GetContractABI() (abi.ABI, error) {
	return GetContractABIByType(ContractTypeMinimal)
}

// GetContractABIByType returns the ABI for a specific contract type
// Loads from compiled ABI file (generated by `make compile-contracts`)
func GetContractABIByType(contractType ContractType) (abi.ABI, error) {
	_, abiPath := getContractPaths(contractType)
	repoRoot := findRepoRoot()

	// Try multiple possible locations
	possiblePaths := []string{
		filepath.Join(repoRoot, abiPath), // From repo root (most reliable)
		abiPath,                          // Direct path (if running from repo root)
		filepath.Join("..", "..", "..", "..", "..", "..", abiPath), // From x/tokenization/precompile/test/integration
		filepath.Join("..", "..", "..", "..", abiPath),             // From x/tokenization/precompile
		filepath.Join(".", abiPath),
	}

	var abiBytes []byte
	var readErr error
	for _, path := range possiblePaths {
		abiBytes, readErr = os.ReadFile(path)
		if readErr == nil {
			break
		}
	}

	if readErr != nil {
		return abi.ABI{}, fmt.Errorf("failed to read contract ABI from %s (tried: %v). Run 'make compile-contracts' first: %w", abiPath, possiblePaths, readErr)
	}

	contractABI, parseErr := abi.JSON(strings.NewReader(string(abiBytes)))
	if parseErr != nil {
		return abi.ABI{}, fmt.Errorf("failed to parse ABI: %w", parseErr)
	}

	return contractABI, nil
}

// GetContractBytecode returns the bytecode for the MinimalTestContract (default)
// Loads from compiled bytecode file (generated by `make compile-contracts`)
func GetContractBytecode() ([]byte, error) {
	return GetContractBytecodeByType(ContractTypeMinimal)
}

// GetContractBytecodeByType returns the bytecode for a specific contract type
// Loads from compiled bytecode file (generated by `make compile-contracts`)
func GetContractBytecodeByType(contractType ContractType) ([]byte, error) {
	bytecodePath, _ := getContractPaths(contractType)
	repoRoot := findRepoRoot()

	// Try multiple possible locations
	possiblePaths := []string{
		filepath.Join(repoRoot, bytecodePath), // From repo root (most reliable)
		bytecodePath,                          // Direct path (if running from repo root)
		filepath.Join("..", "..", "..", "..", "..", "..", bytecodePath), // From x/tokenization/precompile/test/integration
		filepath.Join("..", "..", "..", "..", bytecodePath),             // From x/tokenization/precompile
		filepath.Join(".", bytecodePath),
	}

	var bytecodeHexBytes []byte
	var readErr error
	for _, path := range possiblePaths {
		bytecodeHexBytes, readErr = os.ReadFile(path)
		if readErr == nil {
			break
		}
	}

	if readErr != nil {
		return nil, fmt.Errorf("failed to read contract bytecode from %s (tried: %v). Run 'make compile-contracts' first: %w", bytecodePath, possiblePaths, readErr)
	}

	// Remove newlines and trim
	bytecodeHex := strings.TrimSpace(string(bytecodeHexBytes))
	bytecodeHex = strings.ReplaceAll(bytecodeHex, "\n", "")
	bytecodeHex = strings.TrimPrefix(bytecodeHex, "0x")

	// Convert hex string to bytes
	bytecode := make([]byte, len(bytecodeHex)/2)
	for i := 0; i < len(bytecodeHex); i += 2 {
		var b byte
		_, scanErr := fmt.Sscanf(bytecodeHex[i:i+2], "%02x", &b)
		if scanErr != nil {
			return nil, fmt.Errorf("invalid hex string in bytecode: %w", scanErr)
		}
		bytecode[i/2] = b
	}

	return bytecode, nil
}

// ContractEvent represents a parsed contract event
type ContractEvent struct {
	Name      string
	Args      []interface{}
	Address   common.Address
	TxHash    common.Hash
	BlockHash common.Hash
	Index     uint
}

// ParseEventLog parses a single event log from a transaction
func ParseEventLog(
	contractABI abi.ABI,
	eventName string,
	log types.Log,
) (*ContractEvent, error) {
	event, exists := contractABI.Events[eventName]
	if !exists {
		return nil, &tokenization.PrecompileError{
			Code:    tokenization.ErrorCodeInvalidInput,
			Message: "event not found",
			Details: fmt.Sprintf("event %s not found in contract ABI", eventName),
		}
	}

	// Separate indexed and non-indexed inputs
	var indexedInputs abi.Arguments
	var nonIndexedInputs abi.Arguments
	for _, input := range event.Inputs {
		if input.Indexed {
			indexedInputs = append(indexedInputs, input)
		} else {
			nonIndexedInputs = append(nonIndexedInputs, input)
		}
	}

	// Unpack non-indexed event data
	nonIndexedArgs, err := nonIndexedInputs.Unpack(log.Data)
	if err != nil {
		return nil, err
	}

	// Unpack indexed arguments from topics
	// Topics[0] is the event signature hash
	// Topics[1:] are the indexed arguments
	// Note: For indexed arguments, topics are 32-byte hashes
	// This is a simplified version - full implementation would properly decode
	// each topic based on its type (address, uint256, bytes32, etc.)
	// For now, we'll store the topic hash directly - callers can decode based on type
	indexedArgs := make([]interface{}, 0)
	for i := range indexedInputs {
		if i+1 < len(log.Topics) {
			// Store the topic hash - callers should decode based on the input type
			// Common types: address (20 bytes padded to 32), uint256 (32 bytes), bytes32 (32 bytes)
			indexedArgs = append(indexedArgs, log.Topics[i+1])
		}
	}

	// Combine indexed and non-indexed args in order
	// Note: The order should match the original event definition
	// For simplicity, we'll put indexed args first, then non-indexed
	allArgs := append(indexedArgs, nonIndexedArgs...)
	args := allArgs

	return &ContractEvent{
		Name:      eventName,
		Args:      args,
		Address:   log.Address,
		TxHash:    log.TxHash,
		BlockHash: log.BlockHash,
		Index:     log.Index,
	}, nil
}
