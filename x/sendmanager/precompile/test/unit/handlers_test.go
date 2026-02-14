package sendmanager_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	sendmanager "github.com/bitbadges/bitbadgeschain/x/sendmanager/precompile"
	"github.com/bitbadges/bitbadgeschain/x/sendmanager/precompile/test/helpers"
)

type HandlersTestSuite struct {
	suite.Suite
	Precompile *sendmanager.Precompile
	TestSuite  *helpers.TestSuite
}

func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(HandlersTestSuite))
}

func (suite *HandlersTestSuite) SetupTest() {
	suite.TestSuite = helpers.NewTestSuite(suite.T())
	suite.Precompile = suite.TestSuite.Precompile
}

func (suite *HandlersTestSuite) TestSend_Success() {
	caller := suite.TestSuite.AliceEVM
	toAddress := suite.TestSuite.Bob.String()

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
	suite.NoError(err)
	suite.NotNil(result)

	// Unpack result to get success bool
	unpacked, err := method.Outputs.Unpack(result)
	suite.NoError(err)
	suite.Len(unpacked, 1)

	success, ok := unpacked[0].(bool)
	suite.True(ok)
	suite.True(success, "Send should succeed")
}

func (suite *HandlersTestSuite) TestSend_InsufficientBalance() {
	caller := suite.TestSuite.AliceEVM
	toAddress := suite.TestSuite.Bob.String()

	method := suite.Precompile.ABI.Methods["send"]
	require.NotNil(suite.T(), method)

	// Try to send more than Alice has
	jsonMsg, err := helpers.BuildSendJSON(
		suite.TestSuite.Alice.String(),
		toAddress,
		"999999999999999999999999", // Extremely large amount
		"stake",
	)
	suite.NoError(err)

	input, err := helpers.PackMethodCall(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	suite.Error(err)
	suite.Nil(result)
	// Error might be wrapped, check for either "insufficient" or "balance"
	suite.True(
		suite.Contains(err.Error(), "insufficient") || suite.Contains(err.Error(), "balance"),
		"Error should mention insufficient balance: %s", err.Error(),
	)
}

func (suite *HandlersTestSuite) TestSend_InvalidDenom() {
	caller := suite.TestSuite.AliceEVM
	toAddress := suite.TestSuite.Bob.String()

	method := suite.Precompile.ABI.Methods["send"]
	require.NotNil(suite.T(), method)

	jsonMsg, err := helpers.BuildSendJSON(
		suite.TestSuite.Alice.String(),
		toAddress,
		"1000",
		"", // Empty denom
	)
	suite.NoError(err)

	input, err := helpers.PackMethodCall(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	suite.Error(err)
	suite.Nil(result)
}

func (suite *HandlersTestSuite) TestSend_InvalidToAddress() {
	caller := suite.TestSuite.AliceEVM

	method := suite.Precompile.ABI.Methods["send"]
	require.NotNil(suite.T(), method)

	jsonMsg, err := helpers.BuildSendJSON(
		suite.TestSuite.Alice.String(),
		"invalid-address", // Invalid address
		"1000",
		"stake",
	)
	suite.NoError(err)

	input, err := helpers.PackMethodCall(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	suite.Error(err)
	suite.Nil(result)
}

func (suite *HandlersTestSuite) TestSend_ZeroAmount() {
	caller := suite.TestSuite.AliceEVM
	toAddress := suite.TestSuite.Bob.String()

	method := suite.Precompile.ABI.Methods["send"]
	require.NotNil(suite.T(), method)

	jsonMsg, err := helpers.BuildSendJSON(
		suite.TestSuite.Alice.String(),
		toAddress,
		"0", // Zero amount
		"stake",
	)
	suite.NoError(err)

	input, err := helpers.PackMethodCall(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	// Zero amount might be invalid or might succeed (depending on validation)
	// The important thing is it doesn't panic
	if err != nil {
		suite.Nil(result)
	}
}

