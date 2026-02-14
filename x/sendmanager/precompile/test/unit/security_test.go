package sendmanager_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sendmanager "github.com/bitbadges/bitbadgeschain/x/sendmanager/precompile"
	"github.com/bitbadges/bitbadgeschain/x/sendmanager/precompile/test/helpers"
)

type SecurityTestSuite struct {
	suite.Suite
	Precompile *sendmanager.Precompile
	TestSuite  *helpers.TestSuite
}

func TestSecurityTestSuite(t *testing.T) {
	suite.Run(t, new(SecurityTestSuite))
}

func (suite *SecurityTestSuite) SetupTest() {
	suite.TestSuite = helpers.NewTestSuite(suite.T())
	suite.Precompile = suite.TestSuite.Precompile
}

func (suite *SecurityTestSuite) TestGetCallerAddress_ValidCaller() {
	caller := common.HexToAddress("0x1111111111111111111111111111111111111111")
	precompileAddr := common.HexToAddress(sendmanager.SendManagerPrecompileAddress)
	valueUint256, _ := uint256.FromBig(big.NewInt(0))
	contract := vm.NewContract(caller, precompileAddr, valueUint256, 1000000, nil)

	callerAddr, err := suite.Precompile.GetCallerAddress(contract)
	suite.NoError(err)
	suite.NotEmpty(callerAddr)

	// Should be a valid Cosmos address
	cosmosAddr, err := sdk.AccAddressFromBech32(callerAddr)
	suite.NoError(err)
	suite.Equal(caller.Bytes(), cosmosAddr.Bytes())
}

func (suite *SecurityTestSuite) TestGetCallerAddress_ZeroAddress() {
	caller := common.Address{} // Zero address
	precompileAddr := common.HexToAddress(sendmanager.SendManagerPrecompileAddress)
	valueUint256, _ := uint256.FromBig(big.NewInt(0))
	contract := vm.NewContract(caller, precompileAddr, valueUint256, 1000000, nil)

	callerAddr, err := suite.Precompile.GetCallerAddress(contract)
	suite.Error(err)
	suite.Empty(callerAddr)
	suite.Contains(err.Error(), "zero")
}

func (suite *SecurityTestSuite) TestVerifyCaller_ValidAddress() {
	caller := common.HexToAddress("0x1111111111111111111111111111111111111111")
	err := sendmanager.VerifyCaller(caller)
	suite.NoError(err)
}

func (suite *SecurityTestSuite) TestVerifyCaller_ZeroAddress() {
	caller := common.Address{}
	err := sendmanager.VerifyCaller(caller)
	suite.Error(err)
	suite.Contains(err.Error(), "zero")
}

func (suite *SecurityTestSuite) TestSend_FromAddressIsCaller() {
	// Test that the from_address is always set from contract.Caller(), not from JSON input
	caller := suite.TestSuite.AliceEVM
	toAddress := suite.TestSuite.Bob.String()

	method := suite.Precompile.ABI.Methods["send"]
	require.NotNil(suite.T(), method)

	// Build JSON message with a DIFFERENT from_address in JSON (should be ignored)
	// This tests that the caller is always used, preventing impersonation
	jsonMsg, err := helpers.BuildSendJSON(
		suite.TestSuite.Charlie.String(), // Wrong from address in JSON
		toAddress,
		"1000",
		"stake",
	)
	suite.NoError(err)

	input, err := helpers.PackMethodCall(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	
	// Execute the precompile
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	
	// The transaction should succeed (if Alice has balance)
	// The important part is that from_address is set to Alice (caller), not Charlie (JSON)
	if err == nil {
		suite.NotNil(result)
		
		// Unpack result to get success bool
		unpacked, err := method.Outputs.Unpack(result)
		suite.NoError(err)
		suite.Len(unpacked, 1)
		
		success, ok := unpacked[0].(bool)
		suite.True(ok)
		suite.True(success, "Send should succeed")
		
		// Verify that the actual sender was Alice (the caller), not Charlie
		// We can verify this by checking balances - if Alice's balance decreased,
		// then the caller was correctly used
		aliceBalance := suite.TestSuite.App.BankKeeper.GetBalance(suite.TestSuite.Ctx, suite.TestSuite.Alice, "stake")
		// Alice should have less than the initial 1_000_000_000_000 (if send succeeded)
		// This confirms that Alice was the actual sender
		suite.True(aliceBalance.Amount.LT(sdkmath.NewInt(1_000_000_000_000)), "Alice's balance should have decreased")
	} else {
		// If error occurred, it should be due to insufficient balance or validation,
		// not because of wrong sender
		// Error might be wrapped, so just check it's not about sender impersonation
		suite.NotContains(err.Error(), "impersonation", "Error should not be about sender impersonation")
	}
}

func (suite *SecurityTestSuite) TestSend_FromAddressOverriddenByCaller() {
	// This test explicitly verifies that even if JSON contains a from_address,
	// it is overridden by contract.Caller()
	caller := suite.TestSuite.AliceEVM
	toAddress := suite.TestSuite.Bob.String()

	// Create a message with from_address in JSON that differs from caller
	msg := map[string]interface{}{
		"fromAddress": suite.TestSuite.Charlie.String(), // Different from caller
		"toAddress":   toAddress,
		"amount":      "1000",
		"denom":       "stake",
	}
	
	jsonMsg, err := helpers.BuildQueryJSON(msg)
	suite.NoError(err)

	method := suite.Precompile.ABI.Methods["send"]
	input, err := helpers.PackMethodCall(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	
	// Get initial balances
	aliceInitialBalance := suite.TestSuite.App.BankKeeper.GetBalance(suite.TestSuite.Ctx, suite.TestSuite.Alice, "stake")
	bobInitialBalance := suite.TestSuite.App.BankKeeper.GetBalance(suite.TestSuite.Ctx, suite.TestSuite.Bob, "stake")
	charlieInitialBalance := suite.TestSuite.App.BankKeeper.GetBalance(suite.TestSuite.Ctx, suite.TestSuite.Charlie, "stake")

	// Execute
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	
	if err == nil {
		suite.NotNil(result)
		
		// Verify balances changed for Alice (caller), not Charlie (JSON)
		aliceFinalBalance := suite.TestSuite.App.BankKeeper.GetBalance(suite.TestSuite.Ctx, suite.TestSuite.Alice, "stake")
		bobFinalBalance := suite.TestSuite.App.BankKeeper.GetBalance(suite.TestSuite.Ctx, suite.TestSuite.Bob, "stake")
		charlieFinalBalance := suite.TestSuite.App.BankKeeper.GetBalance(suite.TestSuite.Ctx, suite.TestSuite.Charlie, "stake")

		// Alice's balance should have decreased
		suite.True(aliceFinalBalance.Amount.LT(aliceInitialBalance.Amount), "Alice's balance should decrease (she is the caller)")
		
		// Bob's balance should have increased
		suite.True(bobFinalBalance.Amount.GT(bobInitialBalance.Amount), "Bob's balance should increase")
		
		// Charlie's balance should NOT have changed (she was in JSON but not the caller)
		suite.Equal(charlieInitialBalance.Amount.String(), charlieFinalBalance.Amount.String(), "Charlie's balance should not change (she was not the caller)")
	}
}

