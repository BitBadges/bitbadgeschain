package sendmanager_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	sendmanager "github.com/bitbadges/bitbadgeschain/x/sendmanager/precompile"
	"github.com/bitbadges/bitbadgeschain/x/sendmanager/precompile/test/helpers"
)

type ErrorTestSuite struct {
	suite.Suite
	Precompile *sendmanager.Precompile
	TestSuite  *helpers.TestSuite
}

func TestErrorTestSuite(t *testing.T) {
	suite.Run(t, new(ErrorTestSuite))
}

func (suite *ErrorTestSuite) SetupTest() {
	suite.TestSuite = helpers.NewTestSuite(suite.T())
	suite.Precompile = suite.TestSuite.Precompile
}

func (suite *ErrorTestSuite) TestErrorCodeInvalidInput() {
	err := sendmanager.ErrInvalidInput("test error")
	suite.NotNil(err)
	suite.Equal(sendmanager.ErrorCodeInvalidInput, err.Code)
	suite.Contains(err.Error(), "invalid input")
}

func (suite *ErrorTestSuite) TestErrorCodeSendFailed() {
	err := sendmanager.ErrSendFailed("test error")
	suite.NotNil(err)
	suite.Equal(sendmanager.ErrorCodeSendFailed, err.Code)
	suite.Contains(err.Error(), "send failed")
}

func (suite *ErrorTestSuite) TestErrorCodeInsufficientBalance() {
	err := sendmanager.ErrInsufficientBalance("test error")
	suite.NotNil(err)
	suite.Equal(sendmanager.ErrorCodeInsufficientBalance, err.Code)
	suite.Contains(err.Error(), "insufficient balance")
}

func (suite *ErrorTestSuite) TestWrapError() {
	// WrapError only maps Cosmos SDK errors, not PrecompileErrors
	// So a PrecompileError will use the default code
	originalErr := sendmanager.ErrInvalidInput("original")
	wrapped := sendmanager.WrapError(originalErr, sendmanager.ErrorCodeSendFailed, "wrapped")
	suite.NotNil(wrapped)
	// Since PrecompileError is not a Cosmos SDK error, it will use the default code
	suite.Equal(sendmanager.ErrorCodeSendFailed, wrapped.Code)
}

func (suite *ErrorTestSuite) TestWrapError_NoMapping() {
	// Error that doesn't map to a specific code
	originalErr := sendmanager.ErrSendFailed("test")
	wrapped := sendmanager.WrapError(originalErr, sendmanager.ErrorCodeInternalError, "wrapped")
	suite.NotNil(wrapped)
	suite.Equal(sendmanager.ErrorCodeInternalError, wrapped.Code) // Should use default code
}

func (suite *ErrorTestSuite) TestReadOnlyMode() {
	caller := suite.TestSuite.AliceEVM

	method := suite.Precompile.ABI.Methods["send"]
	require.NotNil(suite.T(), method)

	jsonMsg, err := helpers.BuildSendJSON(
		suite.TestSuite.Alice.String(),
		suite.TestSuite.Bob.String(),
		"1000",
		"stake",
	)
	suite.NoError(err)

	input, err := helpers.PackMethodCall(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	
	// Try to execute in read-only mode
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, true)
	suite.Error(err)
	suite.Nil(result)
	// Error might say "read-only", "write protection", or "does not support read-only" depending on the layer
	errorMsg := err.Error()
	hasReadOnly := strings.Contains(errorMsg, "read-only")
	hasWriteProtection := strings.Contains(errorMsg, "write protection")
	hasDoesNotSupport := strings.Contains(errorMsg, "does not support read-only")
	suite.True(
		hasReadOnly || hasWriteProtection || hasDoesNotSupport,
		"Error should mention read-only or write protection: %s", errorMsg,
	)
}

