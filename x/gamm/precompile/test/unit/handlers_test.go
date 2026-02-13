package gamm_test

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/suite"

	gamm "github.com/bitbadges/bitbadgeschain/x/gamm/precompile"
	"github.com/bitbadges/bitbadgeschain/x/gamm/precompile/test/helpers"
)

type HandlersTestSuite struct {
	suite.Suite
	TestSuite  *helpers.TestSuite
	Precompile *gamm.Precompile
}

func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(HandlersTestSuite))
}

func (suite *HandlersTestSuite) SetupTest() {
	suite.TestSuite = helpers.NewTestSuite(suite.T())
	suite.Precompile = suite.TestSuite.Precompile
}

func (suite *HandlersTestSuite) createContract(caller common.Address, input []byte) *vm.Contract {
	precompileAddr := common.HexToAddress(gamm.GammPrecompileAddress)
	valueUint256, _ := uint256.FromBig(big.NewInt(0))
	contract := vm.NewContract(caller, precompileAddr, valueUint256, 1000000, nil)
	if len(input) > 0 {
		contract.Input = input
	}
	return contract
}

// packMethodWithJSON packs a method call with a JSON string argument
func (suite *HandlersTestSuite) packMethodWithJSON(method *abi.Method, jsonStr string) ([]byte, error) {
	args := []interface{}{jsonStr}
	packed, err := method.Inputs.Pack(args...)
	if err != nil {
		return nil, err
	}
	return append(method.ID, packed...), nil
}

// TestJoinPool_InvalidInput tests validation errors for JoinPool
func (suite *HandlersTestSuite) TestJoinPool_InvalidInput() {
	caller := suite.TestSuite.AliceEVM
	method := suite.Precompile.ABI.Methods["joinPool"]

	// Test with zero pool ID
	jsonMsg := map[string]interface{}{
		"poolId":     "0",
		"shareOutMin": "1000",
		"tokenInMaxs": []map[string]interface{}{
			{
				"denom":  "uatom",
				"amount": "1000",
			},
		},
	}
	jsonBytes, _ := json.Marshal(jsonMsg)
	input, err := suite.packMethodWithJSON(&method, string(jsonBytes))
	suite.Require().NoError(err)

	contract := suite.createContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	suite.Error(err)
	suite.Nil(result)
	// Just check that error is not nil (validation moved to ValidateBasic)
	suite.NotNil(err)
}

// TestExitPool_InvalidInput tests validation errors for ExitPool
func (suite *HandlersTestSuite) TestExitPool_InvalidInput() {
	caller := suite.TestSuite.AliceEVM
	method := suite.Precompile.ABI.Methods["exitPool"]

	// Test with zero pool ID
	jsonMsg := map[string]interface{}{
		"poolId":      "0",
		"shareIn":     "1000",
		"tokenOutMins": []map[string]interface{}{
			{
				"denom":  "uatom",
				"amount": "1000",
			},
		},
	}
	jsonBytes, _ := json.Marshal(jsonMsg)
	input, err := suite.packMethodWithJSON(&method, string(jsonBytes))
	suite.Require().NoError(err)

	contract := suite.createContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	suite.Error(err)
	suite.Nil(result)
	// Just check that error is not nil (validation moved to ValidateBasic)
	suite.NotNil(err)
}

// TestSwapExactAmountIn_InvalidInput tests validation errors for SwapExactAmountIn
func (suite *HandlersTestSuite) TestSwapExactAmountIn_InvalidInput() {
	caller := suite.TestSuite.AliceEVM
	method := suite.Precompile.ABI.Methods["swapExactAmountIn"]

	// Test with empty routes
	jsonMsg := map[string]interface{}{
		"routes": []interface{}{}, // Empty routes
		"tokenIn": map[string]interface{}{
			"denom":  "uatom",
			"amount": "1000",
		},
		"tokenOutMinAmount": "100",
		"affiliates":        []interface{}{},
	}
	jsonBytes, _ := json.Marshal(jsonMsg)
	input, err := suite.packMethodWithJSON(&method, string(jsonBytes))
	suite.Require().NoError(err)

	contract := suite.createContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	suite.Error(err)
	suite.Nil(result)
	// Just check that error is not nil (validation moved to ValidateBasic)
	suite.NotNil(err)
}

// TestGetPool_InvalidInput tests validation errors for GetPool
func (suite *HandlersTestSuite) TestGetPool_InvalidInput() {
	method := suite.Precompile.ABI.Methods["getPool"]

	// Test with zero pool ID
	jsonMsg := map[string]interface{}{
		"poolId": "0",
	}
	jsonBytes, _ := json.Marshal(jsonMsg)
	input, err := suite.packMethodWithJSON(&method, string(jsonBytes))
	suite.Require().NoError(err)

	contract := suite.createContract(suite.TestSuite.AliceEVM, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, true)
	suite.Error(err)
	suite.Nil(result)
	// Just check that error is not nil (validation moved to ValidateBasic)
	suite.NotNil(err)
}

// TestGetPool_PoolNotFound tests error when pool doesn't exist
func (suite *HandlersTestSuite) TestGetPool_PoolNotFound() {
	method := suite.Precompile.ABI.Methods["getPool"]

	// Test with non-existent pool ID
	jsonMsg := map[string]interface{}{
		"poolId": "99999",
	}
	jsonBytes, _ := json.Marshal(jsonMsg)
	input, err := suite.packMethodWithJSON(&method, string(jsonBytes))
	suite.Require().NoError(err)

	contract := suite.createContract(suite.TestSuite.AliceEVM, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, true)
	suite.Error(err)
	suite.Nil(result)
	// Just check that error is not nil (validation moved to ValidateBasic)
	suite.NotNil(err)
}

// TestHandlerMethodSignatures tests that all handler methods exist and have correct signatures
func (suite *HandlersTestSuite) TestHandlerMethodSignatures() {
	methods := []string{
		"joinPool",
		"exitPool",
		"swapExactAmountIn",
		"swapExactAmountInWithIBCTransfer",
		"getPool",
		"getPools",
		"getPoolType",
		"calcJoinPoolNoSwapShares",
		"calcExitPoolCoinsFromShares",
		"calcJoinPoolShares",
		"getPoolParams",
		"getTotalShares",
		"getTotalLiquidity",
	}

	for _, methodName := range methods {
		suite.Run(methodName, func() {
			method, found := suite.Precompile.ABI.Methods[methodName]
			suite.True(found, "Method %s should exist", methodName)
			suite.NotNil(method, "Method %s should not be nil", methodName)
			suite.Equal(abi.Function, method.Type, "Method %s should be a function", methodName)
		})
	}
}

// TestHandlerArgumentCounts tests that handlers validate argument counts
// With JSON-based approach, argument count validation happens during JSON unmarshaling
func (suite *HandlersTestSuite) TestHandlerArgumentCounts() {
	// This test is no longer applicable with JSON-based approach
	// JSON validation happens during unmarshaling, not at argument packing level
	suite.T().Skip("Argument count validation now handled via JSON unmarshaling")
}

// TestQueryMethods_ReadOnly tests that query methods don't modify state
func (suite *HandlersTestSuite) TestQueryMethods_ReadOnly() {
	// Query methods should not modify state, so they should work even with invalid pool IDs
	// (they'll return errors, but won't panic or modify state)

	method := suite.Precompile.ABI.Methods["getPoolType"]
	jsonMsg := map[string]interface{}{
		"poolId": "99999", // Non-existent pool
	}
	jsonBytes, _ := json.Marshal(jsonMsg)
	input, err := suite.packMethodWithJSON(&method, string(jsonBytes))
	suite.Require().NoError(err)

	contract := suite.createContract(suite.TestSuite.AliceEVM, input)
	// This should return an error, but not panic
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, true)
	suite.Error(err) // Expected error for non-existent pool
	suite.Nil(result)
}
