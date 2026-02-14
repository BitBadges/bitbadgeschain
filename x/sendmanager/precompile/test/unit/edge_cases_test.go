package sendmanager_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	sendmanager "github.com/bitbadges/bitbadgeschain/x/sendmanager/precompile"
	"github.com/bitbadges/bitbadgeschain/x/sendmanager/precompile/test/helpers"
)

type EdgeCasesTestSuite struct {
	suite.Suite
	Precompile *sendmanager.Precompile
	TestSuite  *helpers.TestSuite
}

func TestEdgeCasesTestSuite(t *testing.T) {
	suite.Run(t, new(EdgeCasesTestSuite))
}

func (suite *EdgeCasesTestSuite) SetupTest() {
	suite.TestSuite = helpers.NewTestSuite(suite.T())
	suite.Precompile = suite.TestSuite.Precompile
}

func (suite *EdgeCasesTestSuite) TestSend_MaximumAmount() {
	caller := suite.TestSuite.AliceEVM
	toAddress := suite.TestSuite.Bob.String()

	method := suite.Precompile.ABI.Methods["send"]
	require.NotNil(suite.T(), method)

	// Test with maximum uint256 value (2^256 - 1)
	maxUint256 := new(big.Int)
	maxUint256.Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))
	maxAmountStr := maxUint256.String()

	jsonMsg, err := helpers.BuildSendJSON(
		suite.TestSuite.Alice.String(),
		toAddress,
		maxAmountStr,
		"stake",
	)
	suite.NoError(err)

	input, err := helpers.PackMethodCall(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)

	// Should fail due to insufficient balance (Alice doesn't have max uint256)
	suite.Error(err)
	suite.Nil(result)
	suite.Contains(err.Error(), "insufficient")
}

func (suite *EdgeCasesTestSuite) TestSend_MinimumAmount() {
	caller := suite.TestSuite.AliceEVM
	toAddress := suite.TestSuite.Bob.String()

	method := suite.Precompile.ABI.Methods["send"]
	require.NotNil(suite.T(), method)

	// Test with minimum amount (1)
	jsonMsg, err := helpers.BuildSendJSON(
		suite.TestSuite.Alice.String(),
		toAddress,
		"1",
		"stake",
	)
	suite.NoError(err)

	input, err := helpers.PackMethodCall(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)

	// Should succeed if Alice has balance, or fail with insufficient funds
	if err != nil {
		suite.Contains(err.Error(), "insufficient")
	} else {
		suite.NotNil(result)
	}
}

func (suite *EdgeCasesTestSuite) TestSend_MultipleCoins() {
	caller := suite.TestSuite.AliceEVM
	toAddress := suite.TestSuite.Bob.String()

	method := suite.Precompile.ABI.Methods["send"]
	require.NotNil(suite.T(), method)

	// Build JSON with multiple coins
	msg := map[string]interface{}{
		"from_address": suite.TestSuite.Alice.String(),
		"to_address":   toAddress,
		"amount": []map[string]interface{}{
			{"denom": "stake", "amount": "1000"},
			{"denom": "uosmo", "amount": "2000"},
		},
	}
	jsonMsg, err := helpers.BuildQueryJSON(msg)
	suite.NoError(err)

	input, err := helpers.PackMethodCall(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)

	// Should succeed if Alice has both denoms, or fail with insufficient funds
	if err != nil {
		suite.Contains(err.Error(), "insufficient")
	} else {
		suite.NotNil(result)
	}
}

func (suite *EdgeCasesTestSuite) TestSend_VeryLongDenom() {
	caller := suite.TestSuite.AliceEVM
	toAddress := suite.TestSuite.Bob.String()

	method := suite.Precompile.ABI.Methods["send"]
	require.NotNil(suite.T(), method)

	// Create a very long denom string (but still valid)
	longDenom := ""
	for i := 0; i < 200; i++ {
		longDenom += "a"
	}

	jsonMsg, err := helpers.BuildSendJSON(
		suite.TestSuite.Alice.String(),
		toAddress,
		"1000",
		longDenom,
	)
	suite.NoError(err)

	input, err := helpers.PackMethodCall(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)

	// Should fail due to invalid denom or insufficient balance
	suite.Error(err)
	suite.Nil(result)
}

func (suite *EdgeCasesTestSuite) TestSend_SameFromAndToAddress() {
	caller := suite.TestSuite.AliceEVM
	// Send to self
	toAddress := suite.TestSuite.Alice.String()

	method := suite.Precompile.ABI.Methods["send"]
	require.NotNil(suite.T(), method)

	jsonMsg, err := helpers.BuildSendJSON(
		suite.TestSuite.Alice.String(),
		toAddress,
		"1000",
		"stake",
	)
	suite.NoError(err)

	input, err := helpers.PackMethodCall(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)

	// Sending to self might succeed or fail depending on validation
	// The important thing is it doesn't panic
	if err != nil {
		suite.Nil(result)
	} else {
		suite.NotNil(result)
	}
}

func (suite *EdgeCasesTestSuite) TestSend_LargeCoinArray() {
	caller := suite.TestSuite.AliceEVM
	toAddress := suite.TestSuite.Bob.String()

	method := suite.Precompile.ABI.Methods["send"]
	require.NotNil(suite.T(), method)

	// Build JSON with many coins with different denoms (to avoid duplicate denom error)
	// Coins must be sorted by denom for validation
	denoms := []string{"dai", "stake", "uatom", "uion", "uosmo", "usdc", "uusdc", "uusdt", "wbtc", "weth"}
	coins := []map[string]interface{}{}
	for i := 0; i < 10; i++ {
		coins = append(coins, map[string]interface{}{
			"denom":  denoms[i],
			"amount": "100",
		})
	}

	msg := map[string]interface{}{
		"from_address": suite.TestSuite.Alice.String(),
		"to_address":   toAddress,
		"amount":       coins,
	}
	jsonMsg, err := helpers.BuildQueryJSON(msg)
	suite.NoError(err)

	input, err := helpers.PackMethodCall(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)

	// Should handle multiple coins (may fail due to insufficient balance or invalid denom)
	// The test verifies that the precompile can handle multiple coins without panicking
	if err != nil {
		// Error might be about insufficient balance, invalid denom, duplicate denoms, or sorting
		// All of these are acceptable - the test verifies the precompile handles the input
		suite.True(
			suite.Contains(err.Error(), "insufficient") ||
				suite.Contains(err.Error(), "invalid") ||
				suite.Contains(err.Error(), "denom") ||
				suite.Contains(err.Error(), "duplicate") ||
				suite.Contains(err.Error(), "sorted"),
			"Error should mention balance, denom, duplicate, or sorted: %s", err.Error(),
		)
	} else {
		suite.NotNil(result)
	}
}

func (suite *EdgeCasesTestSuite) TestSend_BoundaryAmounts() {
	caller := suite.TestSuite.AliceEVM
	toAddress := suite.TestSuite.Bob.String()

	method := suite.Precompile.ABI.Methods["send"]
	require.NotNil(suite.T(), method)

	testCases := []struct {
		name   string
		amount string
	}{
		{"one", "1"},
		{"max_int64", "9223372036854775807"},
		{"max_uint64", "18446744073709551615"},
		{"very_large", "999999999999999999999999999999"},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			jsonMsg, err := helpers.BuildSendJSON(
				suite.TestSuite.Alice.String(),
				toAddress,
				tc.amount,
				"stake",
			)
			suite.NoError(err)

			input, err := helpers.PackMethodCall(&method, jsonMsg)
			suite.NoError(err)

			contract := suite.TestSuite.CreateMockContract(caller, input)
			_, err = suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
			// Should not panic, may fail due to insufficient balance
			if err != nil {
				suite.Contains(err.Error(), "insufficient")
			}
		})
	}
}

func (suite *EdgeCasesTestSuite) TestSend_EmptyCoinsArray() {
	caller := suite.TestSuite.AliceEVM
	toAddress := suite.TestSuite.Bob.String()

	method := suite.Precompile.ABI.Methods["send"]
	require.NotNil(suite.T(), method)

	// Build JSON with empty coins array
	msg := map[string]interface{}{
		"from_address": suite.TestSuite.Alice.String(),
		"to_address":   toAddress,
		"amount":       []interface{}{},
	}
	jsonMsg, err := helpers.BuildQueryJSON(msg)
	suite.NoError(err)

	input, err := helpers.PackMethodCall(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)

	// Empty coins array should fail validation
	suite.Error(err)
	suite.Nil(result)
}

func (suite *EdgeCasesTestSuite) TestSend_InvalidAmountFormat() {
	caller := suite.TestSuite.AliceEVM
	toAddress := suite.TestSuite.Bob.String()

	method := suite.Precompile.ABI.Methods["send"]
	require.NotNil(suite.T(), method)

	testCases := []struct {
		name   string
		amount string
	}{
		{"decimal", "100.5"},
		{"scientific", "1e10"},
		{"negative_string", "-1000"},
		{"non_numeric", "abc"},
		// Note: "0xFF" might be parsed as 255, so we skip it
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			msg := map[string]interface{}{
				"from_address": suite.TestSuite.Alice.String(),
				"to_address":   toAddress,
				"amount": []map[string]interface{}{
					{"denom": "stake", "amount": tc.amount},
				},
			}
			jsonMsg, err := helpers.BuildQueryJSON(msg)
			suite.NoError(err)

			input, err := helpers.PackMethodCall(&method, jsonMsg)
			suite.NoError(err)

			contract := suite.TestSuite.CreateMockContract(caller, input)
			result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)

			// Should fail validation
			suite.Error(err)
			suite.Nil(result)
		})
	}
}
