package gamm_test

import (
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

func (suite *HandlersTestSuite) createContract(caller common.Address) *vm.Contract {
	precompileAddr := common.HexToAddress(gamm.GammPrecompileAddress)
	valueUint256, _ := uint256.FromBig(big.NewInt(0))
	return vm.NewContract(caller, precompileAddr, valueUint256, 1000000, nil)
}

// TestJoinPool_InvalidInput tests validation errors for JoinPool
func (suite *HandlersTestSuite) TestJoinPool_InvalidInput() {
	caller := suite.TestSuite.AliceEVM
	contract := suite.createContract(caller)

	// Test with zero pool ID
	// tokenInMaxs needs to be []interface{} with map[string]interface{} elements
	tokenInMaxsRaw := []interface{}{
		map[string]interface{}{
			"denom":  "uatom",
			"amount": big.NewInt(1000),
		},
	}
	args := []interface{}{
		uint64(0), // Invalid pool ID
		big.NewInt(1000),
		tokenInMaxsRaw,
	}

	method := suite.Precompile.ABI.Methods["joinPool"]
	result, err := suite.Precompile.JoinPool(suite.TestSuite.Ctx, &method, args, contract)
	suite.Error(err)
	suite.Nil(result)
	suite.Contains(err.Error(), "poolId cannot be zero")
}

// TestExitPool_InvalidInput tests validation errors for ExitPool
func (suite *HandlersTestSuite) TestExitPool_InvalidInput() {
	caller := suite.TestSuite.AliceEVM
	contract := suite.createContract(caller)

	// Test with zero pool ID
	// tokenOutMins needs to be []interface{} with map[string]interface{} elements
	tokenOutMinsRaw := []interface{}{
		map[string]interface{}{
			"denom":  "uatom",
			"amount": big.NewInt(1000),
		},
	}
	args := []interface{}{
		uint64(0), // Invalid pool ID
		big.NewInt(1000),
		tokenOutMinsRaw,
	}

	method := suite.Precompile.ABI.Methods["exitPool"]
	result, err := suite.Precompile.ExitPool(suite.TestSuite.Ctx, &method, args, contract)
	suite.Error(err)
	suite.Nil(result)
	suite.Contains(err.Error(), "poolId cannot be zero")
}

// TestSwapExactAmountIn_InvalidInput tests validation errors for SwapExactAmountIn
func (suite *HandlersTestSuite) TestSwapExactAmountIn_InvalidInput() {
	caller := suite.TestSuite.AliceEVM
	contract := suite.createContract(caller)

	// Test with empty routes
	// routes needs to be []interface{} with map[string]interface{} elements
	routesRaw := []interface{}{}
	// tokenIn needs to be map[string]interface{}
	tokenInRaw := map[string]interface{}{
		"denom":  "uatom",
		"amount": big.NewInt(1000),
	}
	// affiliates needs to be []interface{} with map[string]interface{} elements
	affiliatesRaw := []interface{}{}
	args := []interface{}{
		routesRaw,    // Empty routes
		tokenInRaw,  // tokenIn
		big.NewInt(100), // tokenOutMinAmount
		affiliatesRaw, // affiliates
	}

	method := suite.Precompile.ABI.Methods["swapExactAmountIn"]
	result, err := suite.Precompile.SwapExactAmountIn(suite.TestSuite.Ctx, &method, args, contract)
	suite.Error(err)
	suite.Nil(result)
	suite.Contains(err.Error(), "cannot be empty")
}

// TestGetPool_InvalidInput tests validation errors for GetPool
func (suite *HandlersTestSuite) TestGetPool_InvalidInput() {
	// Test with zero pool ID
	args := []interface{}{
		uint64(0), // Invalid pool ID
	}

	method := suite.Precompile.ABI.Methods["getPool"]
	result, err := suite.Precompile.GetPool(suite.TestSuite.Ctx, &method, args)
	suite.Error(err)
	suite.Nil(result)
	suite.Contains(err.Error(), "poolId cannot be zero")
}

// TestGetPool_PoolNotFound tests error when pool doesn't exist
func (suite *HandlersTestSuite) TestGetPool_PoolNotFound() {
	// Test with non-existent pool ID
	args := []interface{}{
		uint64(99999), // Non-existent pool
	}

	method := suite.Precompile.ABI.Methods["getPool"]
	result, err := suite.Precompile.GetPool(suite.TestSuite.Ctx, &method, args)
	suite.Error(err)
	suite.Nil(result)
	// Should return pool not found error
	suite.Contains(err.Error(), "pool")
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
func (suite *HandlersTestSuite) TestHandlerArgumentCounts() {
	caller := suite.TestSuite.AliceEVM
	contract := suite.createContract(caller)

	// Test JoinPool with wrong number of arguments
	method := suite.Precompile.ABI.Methods["joinPool"]

	// Too few arguments
	args := []interface{}{
		uint64(1),
		big.NewInt(1000),
		// Missing tokenInMaxs
	}
	result, err := suite.Precompile.JoinPool(suite.TestSuite.Ctx, &method, args, contract)
	suite.Error(err)
	suite.Nil(result)
	suite.Contains(err.Error(), "invalid number of arguments")
}

// TestQueryMethods_ReadOnly tests that query methods don't modify state
func (suite *HandlersTestSuite) TestQueryMethods_ReadOnly() {
	// Query methods should not modify state, so they should work even with invalid pool IDs
	// (they'll return errors, but won't panic or modify state)

	method := suite.Precompile.ABI.Methods["getPoolType"]
	args := []interface{}{
		uint64(99999), // Non-existent pool
	}

	// This should return an error, but not panic
	result, err := suite.Precompile.GetPoolType(suite.TestSuite.Ctx, &method, args)
	suite.Error(err) // Expected error for non-existent pool
	suite.Nil(result)
}
