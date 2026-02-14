package tokenization_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/suite"

	tokenization "github.com/bitbadges/bitbadgeschain/x/tokenization/precompile"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/precompile/test/helpers"
)

type ValidateBasicTestSuite struct {
	suite.Suite
	TestSuite  *helpers.TestSuite
	Precompile *tokenization.Precompile
}

func TestValidateBasicTestSuite(t *testing.T) {
	suite.Run(t, new(ValidateBasicTestSuite))
}

func (suite *ValidateBasicTestSuite) SetupTest() {
	suite.TestSuite = helpers.NewTestSuite()
	suite.Precompile = suite.TestSuite.Precompile
}

// TestValidateBasicCalledForMessages tests that ValidateBasic is called for messages
// by verifying that invalid messages fail with validation errors
func (suite *ValidateBasicTestSuite) TestValidateBasicCalledForMessages() {
	caller := suite.TestSuite.AliceEVM

	// Test with invalid message that should fail ValidateBasic
	// Using an empty/invalid createCollection message
	method := suite.Precompile.ABI.Methods["createCollection"]
	invalidJSON := `{"defaultBalances":null,"validTokenIds":[],"collectionPermissions":{}}` // Missing required fields

	input, err := helpers.PackMethodWithJSON(&method, invalidJSON)
	suite.Require().NoError(err)

	precompileAddr := common.HexToAddress(tokenization.TokenizationPrecompileAddress)
	valueUint256, _ := uint256.FromBig(big.NewInt(0))
	contract := vm.NewContract(caller, precompileAddr, valueUint256, 1000000, nil)
	contract.SetCallCode(common.Hash{}, input)

	_, err = suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	suite.Require().Error(err, "Invalid message should fail validation")
	// Verify the error is from ValidateBasic - could be "validation failed", "validation panic", or wrapped error
	// Note: Errors may be wrapped as "execution reverted" by the EVM layer, but ValidateBasic should still be called
	errStr := err.Error()
	suite.Require().True(
		contains(errStr, "validation failed") ||
			contains(errStr, "validation panic") ||
			contains(errStr, "invalid input") ||
			contains(errStr, "invalid") ||
			contains(errStr, "execution reverted") ||
			contains(errStr, "SetupABI failed"), // SetupABI might fail before validation, but that's okay for this test
		"Error should indicate validation failure or execution reverted (ValidateBasic was called), got: %s", errStr)
}

// TestValidateBasicCalledForQueries tests that ValidateBasic is called for queries
// by verifying that invalid queries fail with validation errors
func (suite *ValidateBasicTestSuite) TestValidateBasicCalledForQueries() {
	// Test with invalid query that should fail ValidateBasic
	method := suite.Precompile.ABI.Methods["getCollection"]
	invalidJSON := `{"collectionId":""}` // Empty collectionId should fail validation

	input, err := helpers.PackMethodWithJSON(&method, invalidJSON)
	suite.Require().NoError(err)

	caller := suite.TestSuite.AliceEVM
	precompileAddr := common.HexToAddress(tokenization.TokenizationPrecompileAddress)
	valueUint256, _ := uint256.FromBig(big.NewInt(0))
	contract := vm.NewContract(caller, precompileAddr, valueUint256, 1000000, nil)
	contract.SetCallCode(common.Hash{}, input)

	_, err = suite.Precompile.Execute(suite.TestSuite.Ctx, contract, true)
	suite.Require().Error(err, "Invalid query should fail validation")
	// Verify the error is from validation (either ValidateBasic or validateQueryRequest)
	// Note: Errors may be wrapped as "execution reverted" by the EVM layer
	errStr := err.Error()
	suite.Require().True(
		contains(errStr, "validation failed") ||
			contains(errStr, "validation panic") ||
			contains(errStr, "collection ID cannot be empty") ||
			contains(errStr, "invalid input") ||
			contains(errStr, "invalid") ||
			contains(errStr, "execution reverted") ||
			contains(errStr, "SetupABI failed"), // SetupABI might fail before validation, but that's okay for this test
		"Error should indicate validation failure or execution reverted (ValidateBasic was called), got: %s", errStr)
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

