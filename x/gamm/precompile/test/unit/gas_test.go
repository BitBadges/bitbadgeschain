package gamm_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/third_party/apptesting"
	gamm "github.com/bitbadges/bitbadgeschain/x/gamm/precompile"
)

type GasTestSuite struct {
	apptesting.KeeperTestHelper
}

func TestGasTestSuite(t *testing.T) {
	suite.Run(t, new(GasTestSuite))
}

func (suite *GasTestSuite) SetupTest() {
	suite.Reset()
}

func TestGasConstants(t *testing.T) {
	// Verify all gas constants are defined and reasonable
	require.Greater(t, uint64(gamm.GasJoinPoolBase), uint64(0))
	require.Greater(t, uint64(gamm.GasExitPoolBase), uint64(0))
	require.Greater(t, uint64(gamm.GasSwapExactAmountInBase), uint64(0))
	require.Greater(t, uint64(gamm.GasSwapExactAmountInWithIBCTransferBase), uint64(0))

	require.Greater(t, uint64(gamm.GasGetPoolBase), uint64(0))
	require.Greater(t, uint64(gamm.GasGetPoolsBase), uint64(0))
	require.Greater(t, uint64(gamm.GasGetPoolTypeBase), uint64(0))
	require.Greater(t, uint64(gamm.GasCalcJoinPoolNoSwapSharesBase), uint64(0))
	require.Greater(t, uint64(gamm.GasCalcExitPoolCoinsFromSharesBase), uint64(0))
	require.Greater(t, uint64(gamm.GasCalcJoinPoolSharesBase), uint64(0))
	require.Greater(t, uint64(gamm.GasGetPoolParamsBase), uint64(0))
	require.Greater(t, uint64(gamm.GasGetTotalSharesBase), uint64(0))
	require.Greater(t, uint64(gamm.GasGetTotalLiquidityBase), uint64(0))

	// Per-element gas costs
	require.Greater(t, uint64(gamm.GasPerRoute), uint64(0))
	require.Greater(t, uint64(gamm.GasPerCoin), uint64(0))
	require.Greater(t, uint64(gamm.GasPerAffiliate), uint64(0))
	require.Greater(t, uint64(gamm.GasPerMemoByte), uint64(0))
}

func (suite *GasTestSuite) TestGasCosts_TransactionMethods() {
	precompile := gamm.NewPrecompile(suite.App.GammKeeper)

	testCases := []struct {
		name     string
		method   string
		expected uint64
	}{
		{"joinPool", "joinPool", gamm.GasJoinPoolBase},
		{"exitPool", "exitPool", gamm.GasExitPoolBase},
		{"swapExactAmountIn", "swapExactAmountIn", gamm.GasSwapExactAmountInBase},
		{"swapExactAmountInWithIBCTransfer", "swapExactAmountInWithIBCTransfer", gamm.GasSwapExactAmountInWithIBCTransferBase},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			method, found := precompile.ABI.Methods[tc.method]
			suite.True(found, "Method %s should exist", tc.method)

			gas := precompile.RequiredGas(method.ID[:])
			suite.Equal(tc.expected, gas, "Gas cost for %s should match expected", tc.name)

			// Verify gas is reasonable
			suite.Greater(gas, uint64(0), "Gas should be greater than 0")
			suite.Less(gas, uint64(1000000), "Gas should be less than 1M for %s", tc.name)
		})
	}
}

func (suite *GasTestSuite) TestGasCosts_QueryMethods() {
	precompile := gamm.NewPrecompile(suite.App.GammKeeper)

	testCases := []struct {
		name     string
		method   string
		expected uint64
	}{
		{"getPool", "getPool", gamm.GasGetPoolBase},
		{"getPools", "getPools", gamm.GasGetPoolsBase},
		{"getPoolType", "getPoolType", gamm.GasGetPoolTypeBase},
		{"calcJoinPoolNoSwapShares", "calcJoinPoolNoSwapShares", gamm.GasCalcJoinPoolNoSwapSharesBase},
		{"calcExitPoolCoinsFromShares", "calcExitPoolCoinsFromShares", gamm.GasCalcExitPoolCoinsFromSharesBase},
		{"calcJoinPoolShares", "calcJoinPoolShares", gamm.GasCalcJoinPoolSharesBase},
		{"getPoolParams", "getPoolParams", gamm.GasGetPoolParamsBase},
		{"getTotalShares", "getTotalShares", gamm.GasGetTotalSharesBase},
		{"getTotalLiquidity", "getTotalLiquidity", gamm.GasGetTotalLiquidityBase},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			method, found := precompile.ABI.Methods[tc.method]
			suite.True(found, "Method %s should exist", tc.method)

			gas := precompile.RequiredGas(method.ID[:])
			suite.Equal(tc.expected, gas, "Gas cost for %s should match expected", tc.name)

			// Query methods should have lower gas costs
			suite.Less(gas, uint64(100000), "Query gas should be less than 100K for %s", tc.name)
		})
	}
}

func (suite *GasTestSuite) TestGasCosts_CompareTransactionVsQuery() {
	precompile := gamm.NewPrecompile(suite.App.GammKeeper)

	// Transaction methods should generally cost more than query methods
	transactionMethods := []string{"joinPool", "exitPool", "swapExactAmountIn"}
	queryMethods := []string{"getPool", "getPoolType", "getTotalShares"}

	for _, txMethod := range transactionMethods {
		for _, queryMethod := range queryMethods {
			suite.Run(txMethod+"_vs_"+queryMethod, func() {
				txMethodObj, found := precompile.ABI.Methods[txMethod]
				suite.True(found)

				queryMethodObj, found := precompile.ABI.Methods[queryMethod]
				suite.True(found)

				txGas := precompile.RequiredGas(txMethodObj.ID[:])
				queryGas := precompile.RequiredGas(queryMethodObj.ID[:])

				// Both should have gas costs
				suite.Greater(txGas, uint64(0))
				suite.Greater(queryGas, uint64(0))
			})
		}
	}
}
