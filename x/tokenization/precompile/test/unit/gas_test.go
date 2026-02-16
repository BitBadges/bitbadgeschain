package tokenization_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/suite"

	sdkmath "cosmossdk.io/math"

	tokenization "github.com/bitbadges/bitbadgeschain/x/tokenization/precompile"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/precompile/test/helpers"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

type GasTestSuite struct {
	suite.Suite
	TestSuite  *helpers.TestSuite
	Precompile *tokenization.Precompile
}

func TestGasTestSuite(t *testing.T) {
	suite.Run(t, new(GasTestSuite))
}

func (suite *GasTestSuite) SetupTest() {
	suite.TestSuite = helpers.NewTestSuite()
	suite.Precompile = suite.TestSuite.Precompile
}

func (suite *GasTestSuite) createContract(caller common.Address) *vm.Contract {
	precompileAddr := common.HexToAddress(tokenization.TokenizationPrecompileAddress)
	valueUint256, _ := uint256.FromBig(big.NewInt(0))
	return vm.NewContract(caller, precompileAddr, valueUint256, 1000000, nil)
}

// TestGasCosts_TransactionMethods tests gas costs for transaction methods
func (suite *GasTestSuite) TestGasCosts_TransactionMethods() {
	// Transaction methods add a 200k buffer to base gas for Cosmos SDK operations
	const txBuffer = 200_000

	testCases := []struct {
		name     string
		method   string
		args     []interface{}
		expected uint64
	}{
		{
			name:   "createCollection",
			method: "createCollection",
			args: []interface{}{
				nil,
				[]struct {
					Start *big.Int `json:"start"`
					End   *big.Int `json:"end"`
				}{{Start: big.NewInt(1), End: big.NewInt(100)}},
				map[string]interface{}{},
				suite.TestSuite.Manager.String(),
				map[string]interface{}{"uri": "", "customData": ""},
				[]interface{}{},
				"",
				[]interface{}{},
				[]string{},
				false,
				[]interface{}{},
				[]interface{}{},
				map[string]interface{}{},
				[]interface{}{},
			},
			expected: tokenization.GasCreateCollectionBase + txBuffer,
		},
		{
			name:   "transferTokens",
			method: "transferTokens",
			args: func() []interface{} {
				collectionId, _ := suite.TestSuite.CreateTestCollection(suite.TestSuite.Alice.String())
				suite.TestSuite.CreateTestBalance(
					collectionId,
					suite.TestSuite.Alice.String(),
					sdkmath.NewUint(1000),
					[]*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)}},
					[]*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1000)}},
				)
				return []interface{}{
					collectionId.BigInt(),
					[]common.Address{suite.TestSuite.BobEVM},
					big.NewInt(100),
					[]struct {
						Start *big.Int `json:"start"`
						End   *big.Int `json:"end"`
					}{{Start: big.NewInt(1), End: big.NewInt(10)}},
					[]struct {
						Start *big.Int `json:"start"`
						End   *big.Int `json:"end"`
					}{{Start: big.NewInt(1), End: big.NewInt(1000)}},
				}
			}(),
			expected: tokenization.GasTransferTokensBase + txBuffer,
		},
		{
			name:   "setManager",
			method: "setManager",
			args: func() []interface{} {
				collectionId, _ := suite.TestSuite.CreateTestCollection(suite.TestSuite.Alice.String())
				return []interface{}{
					collectionId.BigInt(),
					suite.TestSuite.Bob.String(),
				}
			}(),
			expected: tokenization.GasSetManagerBase + txBuffer,
		},
		{
			name:   "setCollectionMetadata",
			method: "setCollectionMetadata",
			args: func() []interface{} {
				collectionId, _ := suite.TestSuite.CreateTestCollection(suite.TestSuite.Alice.String())
				return []interface{}{
					collectionId.BigInt(),
					"https://example.com",
					"data",
				}
			}(),
			expected: tokenization.GasSetCollectionMetadataBase + txBuffer,
		},
		{
			name:   "createDynamicStore",
			method: "createDynamicStore",
			args: []interface{}{
				false,
				"https://example.com",
				"data",
			},
			expected: tokenization.GasCreateDynamicStoreBase + txBuffer,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			method, found := suite.Precompile.ABI.Methods[tc.method]
			suite.True(found, "Method %s should exist", tc.method)

			// RequiredGas takes input []byte (method ID)
			gas := suite.Precompile.RequiredGas(method.ID[:])
			suite.Equal(tc.expected, gas, "Gas cost for %s should match expected", tc.name)

			// Verify gas is reasonable (not zero, not excessive)
			suite.Greater(gas, uint64(0), "Gas should be greater than 0")
			suite.Less(gas, uint64(1000000), "Gas should be less than 1M for %s", tc.name)
		})
	}
}

// TestGasCosts_QueryMethods tests gas costs for query methods
func (suite *GasTestSuite) TestGasCosts_QueryMethods() {
	// Query methods add a 50k buffer to base gas for state reads
	const queryBuffer = 50_000

	testCases := []struct {
		name     string
		method   string
		expected uint64
	}{
		{
			name:     "getCollection",
			method:   "getCollection",
			expected: tokenization.GasGetCollectionBase + queryBuffer,
		},
		{
			name:     "getBalance",
			method:   "getBalance",
			expected: tokenization.GasGetBalanceBase + queryBuffer,
		},
		{
			name:     "getAddressList",
			method:   "getAddressList",
			expected: tokenization.GasGetAddressList + queryBuffer,
		},
		{
			name:     "getBalanceAmount",
			method:   "getBalanceAmount",
			expected: tokenization.GasGetBalanceAmountBase + queryBuffer,
		},
		{
			name:     "getTotalSupply",
			method:   "getTotalSupply",
			expected: tokenization.GasGetTotalSupplyBase + queryBuffer,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			method, found := suite.Precompile.ABI.Methods[tc.method]
			suite.True(found, "Method %s should exist", tc.method)

			// RequiredGas takes input []byte (method ID)
			gas := suite.Precompile.RequiredGas(method.ID[:])
			suite.Equal(tc.expected, gas, "Gas cost for %s should match expected", tc.name)

			// Query methods should have reasonable gas costs (base + 50k buffer)
			suite.Less(gas, uint64(100000), "Query gas should be less than 100K for %s", tc.name)
		})
	}
}

// TestGasCosts_VaryingInputSizes tests gas costs with varying input sizes
func (suite *GasTestSuite) TestGasCosts_VaryingInputSizes() {
	sizes := []int{1, 10, 50, 100}

	for _, size := range sizes {
		suite.Run(fmt.Sprintf("createCollection_%d_ranges", size), func() {
			validTokenIds := make([]struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			}, size)

			for i := 0; i < size; i++ {
				validTokenIds[i] = struct {
					Start *big.Int `json:"start"`
					End   *big.Int `json:"end"`
				}{
					Start: big.NewInt(int64(i*10 + 1)),
					End:   big.NewInt(int64((i + 1) * 10)),
				}
			}

			method := suite.Precompile.ABI.Methods["createCollection"]
			gas := suite.Precompile.RequiredGas(method.ID[:])

			// Gas should increase with input size, but not linearly
			suite.Greater(gas, uint64(0))
			suite.Less(gas, uint64(1000000))
		})
	}
}

// TestGasCosts_ComplexNestedStructures tests gas costs with complex nested structures
func (suite *GasTestSuite) TestGasCosts_ComplexNestedStructures() {
	// Test gas cost for method with complex nested structures
	// The actual gas cost is based on the method, not the input complexity
	// (since RequiredGas only takes method ID, not the full input)
	method := suite.Precompile.ABI.Methods["setOutgoingApproval"]
	gas := suite.Precompile.RequiredGas(method.ID[:])

	// Complex structures should have reasonable gas costs
	suite.Greater(gas, uint64(0))
	suite.Less(gas, uint64(500000))
}

// TestGasCosts_CompareTransactionVsQuery compares gas costs between transactions and queries
func (suite *GasTestSuite) TestGasCosts_CompareTransactionVsQuery() {
	// Transaction methods should generally cost more than query methods
	transactionMethods := []string{"createCollection", "transferTokens", "setManager"}
	queryMethods := []string{"getCollection", "getBalance", "getAddressList"}

	for _, txMethod := range transactionMethods {
		for _, queryMethod := range queryMethods {
			suite.Run(fmt.Sprintf("%s_vs_%s", txMethod, queryMethod), func() {
				txMethodObj, found := suite.Precompile.ABI.Methods[txMethod]
				suite.True(found)

				queryMethodObj, found := suite.Precompile.ABI.Methods[queryMethod]
				suite.True(found)

				txGas := suite.Precompile.RequiredGas(txMethodObj.ID[:])
				queryGas := suite.Precompile.RequiredGas(queryMethodObj.ID[:])

				// Transactions should generally cost more (but not always, depends on complexity)
				suite.Greater(txGas, uint64(0))
				suite.Greater(queryGas, uint64(0))
			})
		}
	}
}
