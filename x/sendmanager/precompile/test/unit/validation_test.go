package sendmanager_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	sendmanager "github.com/bitbadges/bitbadgeschain/x/sendmanager/precompile"
	"github.com/bitbadges/bitbadgeschain/x/sendmanager/precompile/test/helpers"
)

type ValidationTestSuite struct {
	suite.Suite
	Precompile *sendmanager.Precompile
	TestSuite  *helpers.TestSuite
}

func TestValidationTestSuite(t *testing.T) {
	suite.Run(t, new(ValidationTestSuite))
}

func (suite *ValidationTestSuite) SetupTest() {
	suite.TestSuite = helpers.NewTestSuite(suite.T())
	suite.Precompile = suite.TestSuite.Precompile
}

func (suite *ValidationTestSuite) TestSend_InvalidJSON() {
	caller := suite.TestSuite.AliceEVM

	method := suite.Precompile.ABI.Methods["send"]
	require.NotNil(suite.T(), method)

	// Invalid JSON
	invalidJSON := `{"fromAddress": "invalid", "toAddress": "invalid"` // Missing closing brace

	input, err := helpers.PackMethodCall(&method, invalidJSON)
	suite.NoError(err) // Packing should succeed

	contract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	suite.Error(err)
	suite.Nil(result)
	suite.Contains(err.Error(), "JSON")
}

func (suite *ValidationTestSuite) TestSend_MissingFields() {
	caller := suite.TestSuite.AliceEVM

	method := suite.Precompile.ABI.Methods["send"]
	require.NotNil(suite.T(), method)

	// Missing required fields
	msg := map[string]interface{}{
		"fromAddress": suite.TestSuite.Alice.String(),
		// Missing toAddress, amount, denom
	}
	jsonMsg, err := helpers.BuildQueryJSON(msg)
	suite.NoError(err)

	input, err := helpers.PackMethodCall(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	suite.Error(err)
	suite.Nil(result)
}

func (suite *ValidationTestSuite) TestSend_EmptyStringFields() {
	caller := suite.TestSuite.AliceEVM

	method := suite.Precompile.ABI.Methods["send"]
	require.NotNil(suite.T(), method)

	jsonMsg, err := helpers.BuildSendJSON(
		"", // Empty from address (will be overridden by caller anyway)
		"", // Empty to address
		"", // Empty amount
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

func (suite *ValidationTestSuite) TestSend_InvalidAmountFormat() {
	caller := suite.TestSuite.AliceEVM
	toAddress := suite.TestSuite.Bob.String()

	method := suite.Precompile.ABI.Methods["send"]
	require.NotNil(suite.T(), method)

	// Invalid amount (not a number)
	msg := map[string]interface{}{
		"fromAddress": suite.TestSuite.Alice.String(),
		"toAddress":   toAddress,
		"amount":      "not-a-number",
		"denom":       "stake",
	}
	jsonMsg, err := helpers.BuildQueryJSON(msg)
	suite.NoError(err)

	input, err := helpers.PackMethodCall(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	suite.Error(err)
	suite.Nil(result)
}

func (suite *ValidationTestSuite) TestSend_NegativeAmount() {
	caller := suite.TestSuite.AliceEVM
	toAddress := suite.TestSuite.Bob.String()

	method := suite.Precompile.ABI.Methods["send"]
	require.NotNil(suite.T(), method)

	jsonMsg, err := helpers.BuildSendJSON(
		suite.TestSuite.Alice.String(),
		toAddress,
		"-1000", // Negative amount
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

