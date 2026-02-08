package helpers

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

// CompileSolidityContract compiles a Solidity contract and returns bytecode and ABI
// contractPath: path to the .sol file
// contractName: name of the contract to compile (must match contract name in file)
// Returns: bytecode (hex string), ABI (abi.ABI), and error
func CompileSolidityContract(contractPath string, contractName string) ([]byte, abi.ABI, error) {
	// Check if contract file exists
	if _, err := os.Stat(contractPath); os.IsNotExist(err) {
		return nil, abi.ABI{}, fmt.Errorf("contract file not found: %s", contractPath)
	}

	// Try solcjs first (npm version), then fall back to solc
	var cmd *exec.Cmd
	var compilerName string
	contractDir := filepath.Dir(contractPath)
	
	// Check for solcjs
	if _, err := exec.LookPath("solcjs"); err == nil {
		compilerName = "solcjs"
		// solcjs uses --bin and --abi flags, output goes to current directory
		// We need to use --base-path and --include-path for imports
		contractFile := filepath.Base(contractPath)
		
		cmd = exec.Command("solcjs",
			"--bin",
			"--abi",
			"--base-path", contractDir,
			"--include-path", contractDir,
			"--output-dir", contractDir,
			contractFile,
		)
		cmd.Dir = contractDir
	} else if _, err := exec.LookPath("solc"); err == nil {
		compilerName = "solc"
		// solc uses --combined-json flag
		cmd = exec.Command("solc",
			"--combined-json", "bin,abi",
			"--allow-paths", contractDir,
			contractPath,
		)
	} else {
		return nil, abi.ABI{}, fmt.Errorf("neither solcjs nor solc found in PATH. Please install solc or solcjs")
	}

	// Execute compilation
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, abi.ABI{}, fmt.Errorf("compilation failed: %w\nOutput: %s", err, string(output))
	}

	// Parse output based on compiler
	if compilerName == "solcjs" {
		return parseSolcjsOutput(contractDir, contractName, string(output))
	} else {
		return parseSolcOutput(contractName, string(output))
	}
}

// parseSolcjsOutput parses output from solcjs compiler
// solcjs creates separate .bin and .abi files
func parseSolcjsOutput(contractDir, contractName, output string) ([]byte, abi.ABI, error) {
	// solcjs creates files like ContractName_sol_ContractName.bin and .abi
	// The exact naming depends on solcjs version, but typically: ContractName_sol_ContractName.bin
	binFile := filepath.Join(contractDir, fmt.Sprintf("%s_sol_%s.bin", contractName, contractName))
	abiFile := filepath.Join(contractDir, fmt.Sprintf("%s_sol_%s.abi", contractName, contractName))
	
	// Try alternative naming patterns
	if _, err := os.Stat(binFile); os.IsNotExist(err) {
		// Try with just contract name
		binFile = filepath.Join(contractDir, fmt.Sprintf("%s.bin", contractName))
		abiFile = filepath.Join(contractDir, fmt.Sprintf("%s.abi", contractName))
	}
	
	// Read bytecode
	bytecodeHex, err := os.ReadFile(binFile)
	if err != nil {
		return nil, abi.ABI{}, fmt.Errorf("failed to read bytecode file %s: %w\nCompiler output: %s", binFile, err, output)
	}
	
	// Remove newlines and trim
	bytecodeHexStr := strings.TrimSpace(string(bytecodeHex))
	bytecodeHexStr = strings.ReplaceAll(bytecodeHexStr, "\n", "")
	
	// Read ABI
	abiBytes, err := os.ReadFile(abiFile)
	if err != nil {
		return nil, abi.ABI{}, fmt.Errorf("failed to read ABI file %s: %w", abiFile, err)
	}
	
	// Parse ABI
	contractABI, err := abi.JSON(strings.NewReader(string(abiBytes)))
	if err != nil {
		return nil, abi.ABI{}, fmt.Errorf("failed to parse ABI: %w", err)
	}
	
	// Convert hex string to bytes
	bytecode, err := hexStringToBytes(bytecodeHexStr)
	if err != nil {
		return nil, abi.ABI{}, fmt.Errorf("failed to parse bytecode: %w", err)
	}
	
	return bytecode, contractABI, nil
}

// parseSolcOutput parses output from solc compiler (combined-json format)
func parseSolcOutput(contractName, output string) ([]byte, abi.ABI, error) {
	var result struct {
		Contracts map[string]map[string]struct {
			Bin string `json:"bin"`
			ABI string `json:"abi"`
		} `json:"contracts"`
	}
	
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		return nil, abi.ABI{}, fmt.Errorf("failed to parse solc output: %w\nOutput: %s", err, output)
	}
	
	// Find the contract in the output
	// solc output format: "contracts/File.sol:ContractName"
	var contractData struct {
		Bin string `json:"bin"`
		ABI string `json:"abi"`
	}
	
	found := false
	for key, contracts := range result.Contracts {
		if contract, ok := contracts[contractName]; ok {
			contractData = contract
			found = true
			break
		}
		// Also check if key contains contract name
		if strings.Contains(key, contractName) {
			// Get the first contract in this file
			for _, contract := range contracts {
				contractData = contract
				found = true
				break
			}
			if found {
				break
			}
		}
	}
	
	if !found {
		return nil, abi.ABI{}, fmt.Errorf("contract %s not found in compilation output\nAvailable contracts: %v", contractName, getContractKeysFromTyped(result.Contracts))
	}
	
	// Parse ABI
	contractABI, err := abi.JSON(strings.NewReader(contractData.ABI))
	if err != nil {
		return nil, abi.ABI{}, fmt.Errorf("failed to parse ABI: %w", err)
	}
	
	// Convert hex string to bytes
	bytecode, err := hexStringToBytes(contractData.Bin)
	if err != nil {
		return nil, abi.ABI{}, fmt.Errorf("failed to parse bytecode: %w", err)
	}
	
	return bytecode, contractABI, nil
}

// getContractKeys extracts all contract keys from the contracts map
func getContractKeys(contracts map[string]map[string]interface{}) []string {
	keys := make([]string, 0)
	for fileKey, fileContracts := range contracts {
		for contractKey := range fileContracts {
			keys = append(keys, fmt.Sprintf("%s:%s", fileKey, contractKey))
		}
	}
	return keys
}

// getContractKeysFromTyped extracts contract keys from typed contracts map
func getContractKeysFromTyped(contracts map[string]map[string]struct {
	Bin string `json:"bin"`
	ABI string `json:"abi"`
}) []string {
	keys := make([]string, 0)
	for fileKey, fileContracts := range contracts {
		for contractKey := range fileContracts {
			keys = append(keys, fmt.Sprintf("%s:%s", fileKey, contractKey))
		}
	}
	return keys
}

// hexStringToBytes converts a hex string to bytes
func hexStringToBytes(hexStr string) ([]byte, error) {
	// Remove 0x prefix if present
	hexStr = strings.TrimPrefix(hexStr, "0x")
	
	// Convert hex string to bytes
	bytecode := make([]byte, len(hexStr)/2)
	for i := 0; i < len(hexStr); i += 2 {
		var b byte
		_, err := fmt.Sscanf(hexStr[i:i+2], "%02x", &b)
		if err != nil {
			return nil, fmt.Errorf("invalid hex string: %w", err)
		}
		bytecode[i/2] = b
	}
	
	return bytecode, nil
}

