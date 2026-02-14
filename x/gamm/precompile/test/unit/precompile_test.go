package gamm_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/third_party/apptesting"
	gamm "github.com/bitbadges/bitbadgeschain/x/gamm/precompile"
)

type PrecompileTestSuite struct {
	apptesting.KeeperTestHelper
}

func TestPrecompileTestSuite(t *testing.T) {
	suite.Run(t, new(PrecompileTestSuite))
}

func (suite *PrecompileTestSuite) SetupTest() {
	// T is automatically set by suite.Run, so we can use Reset
	suite.Reset()
}

func (suite *PrecompileTestSuite) TestPrecompile_RequiredGas() {
	precompile := gamm.NewPrecompile(suite.App.GammKeeper)

	// Test with valid method ID - get the method selector manually
	methodID := precompile.ABI.Methods["joinPool"].ID
	gas := precompile.RequiredGas(methodID[:])
	suite.Equal(uint64(gamm.GasJoinPoolBase), gas)

	// Test with invalid input (too short)
	gas = precompile.RequiredGas([]byte{0x12, 0x34})
	suite.Equal(uint64(0), gas)
}

func (suite *PrecompileTestSuite) TestPrecompile_Structure() {
	precompile := gamm.NewPrecompile(suite.App.GammKeeper)

	// Verify precompile is created correctly
	suite.NotNil(precompile)
	suite.NotNil(precompile.ABI)
	suite.Equal(gamm.GammPrecompileAddress, precompile.ContractAddress.Hex())

	// Verify key methods exist
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
		method, found := precompile.ABI.Methods[methodName]
		suite.True(found, "method %s should exist", methodName)
		suite.NotNil(method, "method %s should not be nil", methodName)
	}
}

func (suite *PrecompileTestSuite) TestPrecompile_IsTransaction() {
	precompile := gamm.NewPrecompile(suite.App.GammKeeper)

	// Test transaction methods
	transactionMethods := []string{
		"joinPool",
		"exitPool",
		"swapExactAmountIn",
		"swapExactAmountInWithIBCTransfer",
	}

	for _, methodName := range transactionMethods {
		method, found := precompile.ABI.Methods[methodName]
		suite.True(found, "method %s should exist", methodName)
		suite.True(precompile.IsTransaction(&method), "method %s should be a transaction", methodName)
	}

	// Test query methods (should not be transactions)
	queryMethods := []string{
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

	for _, methodName := range queryMethods {
		method, found := precompile.ABI.Methods[methodName]
		suite.True(found, "method %s should exist", methodName)
		suite.False(precompile.IsTransaction(&method), "method %s should not be a transaction", methodName)
	}
}

func TestGetABILoadError(t *testing.T) {
	// Test that ABI loading error can be retrieved
	err := gamm.GetABILoadError()
	// ABI should load successfully, so error should be nil
	require.NoError(t, err)
}

func TestPrecompileAddress(t *testing.T) {
	// Verify the precompile address is correct
	require.Equal(t, "0x0000000000000000000000000000000000001002", gamm.GammPrecompileAddress)
}

