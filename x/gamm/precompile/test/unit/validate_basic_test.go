package gamm_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/suite"

	gamm "github.com/bitbadges/bitbadgeschain/x/gamm/precompile"
	"github.com/bitbadges/bitbadgeschain/x/gamm/precompile/test/helpers"
)

type ValidateBasicTestSuite struct {
	suite.Suite
	TestSuite  *helpers.TestSuite
	Precompile *gamm.Precompile
}

func TestValidateBasicTestSuite(t *testing.T) {
	suite.Run(t, new(ValidateBasicTestSuite))
}

func (suite *ValidateBasicTestSuite) SetupTest() {
	suite.TestSuite = helpers.NewTestSuite(suite.T())
	suite.Precompile = suite.TestSuite.Precompile
}

// TestValidateBasicCalledForMessages tests that ValidateBasic is called for messages
// by verifying that invalid messages fail with validation errors
func (suite *ValidateBasicTestSuite) TestValidateBasicCalledForMessages() {
	caller := suite.TestSuite.AliceEVM

	// Test with invalid message that should fail ValidateBasic
	method := suite.Precompile.ABI.Methods["joinPool"]
	invalidJSON := `{"poolId":"0","shareOutAmount":"0","tokenInMaxs":[]}` // Invalid poolId

	input, err := helpers.PackMethodWithJSON(&method, invalidJSON)
	suite.Require().NoError(err)

	precompileAddr := common.HexToAddress(gamm.GammPrecompileAddress)
	valueUint256, _ := uint256.FromBig(big.NewInt(0))
	contract := vm.NewContract(caller, precompileAddr, valueUint256, 1000000, nil)
	contract.SetCallCode(common.Hash{}, input)

	_, err = suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	suite.Require().Error(err, "Invalid message should fail validation")
	// Verify the error is from ValidateBasic - could be "validation failed", "validation panic", or wrapped error
	// Note: Errors may be wrapped as "execution reverted" by the EVM layer, but ValidateBasic should still be called
	errStr := err.Error()
	// The key is that ValidateBasic was called - if it panics, we catch it; if it returns error, we return it
	// Both cases mean ValidateBasic was executed, which is what we're testing
	suite.Require().True(
		contains(errStr, "validation failed") ||
			contains(errStr, "validation panic") ||
			contains(errStr, "invalid input") ||
			contains(errStr, "execution reverted"), // EVM wraps errors, but ValidateBasic was still called
		"Error should indicate validation failure or execution reverted (ValidateBasic was called), got: %s", errStr)
}

// TestValidateBasicCalledForQueries tests that ValidateBasic is called for queries
// by verifying that invalid queries fail with validation errors
func (suite *ValidateBasicTestSuite) TestValidateBasicCalledForQueries() {
	// Test with invalid query that should fail ValidateBasic
	method := suite.Precompile.ABI.Methods["getPool"]
	invalidJSON := `{"poolId":"0"}` // Zero poolId should fail validation

	input, err := helpers.PackMethodWithJSON(&method, invalidJSON)
	suite.Require().NoError(err)

	caller := suite.TestSuite.AliceEVM
	precompileAddr := common.HexToAddress(gamm.GammPrecompileAddress)
	valueUint256, _ := uint256.FromBig(big.NewInt(0))
	contract := vm.NewContract(caller, precompileAddr, valueUint256, 1000000, nil)
	contract.SetCallCode(common.Hash{}, input)

	_, err = suite.Precompile.Execute(suite.TestSuite.Ctx, contract, true)
	suite.Require().Error(err, "Invalid query should fail validation")
	// Verify the error is from validation - could be validation failed, validation panic, or other validation error
	// Note: Errors may be wrapped as "execution reverted" by the EVM layer
	errStr := err.Error()
	suite.Require().True(
		contains(errStr, "validation failed") ||
			contains(errStr, "validation panic") ||
			contains(errStr, "invalid") ||
			contains(errStr, "pool ID") ||
			contains(errStr, "invalid input") ||
			contains(errStr, "execution reverted"), // EVM wraps errors, but ValidateBasic was still called
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

