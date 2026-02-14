package sendmanager_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	sdkmath "cosmossdk.io/math"

	sendmanager "github.com/bitbadges/bitbadgeschain/x/sendmanager/precompile"
	"github.com/bitbadges/bitbadgeschain/x/sendmanager/precompile/test/helpers"
)

type IntegrationTestSuite struct {
	suite.Suite
	Precompile *sendmanager.Precompile
	TestSuite  *helpers.TestSuite
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (suite *IntegrationTestSuite) SetupTest() {
	suite.TestSuite = helpers.NewTestSuite(suite.T())
	suite.Precompile = suite.TestSuite.Precompile
}

func (suite *IntegrationTestSuite) TestSend_CompleteFlow() {
	caller := suite.TestSuite.AliceEVM
	toAddress := suite.TestSuite.Bob.String()

	// Get initial balances
	aliceInitialBalance := suite.TestSuite.App.BankKeeper.GetBalance(suite.TestSuite.Ctx, suite.TestSuite.Alice, "stake")
	bobInitialBalance := suite.TestSuite.App.BankKeeper.GetBalance(suite.TestSuite.Ctx, suite.TestSuite.Bob, "stake")

	require.True(suite.T(), aliceInitialBalance.Amount.GT(sdkmath.ZeroInt()), "Alice should have initial balance")

	method := suite.Precompile.ABI.Methods["send"]
	require.NotNil(suite.T(), method)

	sendAmount := "1000"
	jsonMsg, err := helpers.BuildSendJSON(
		suite.TestSuite.Alice.String(),
		toAddress,
		sendAmount,
		"stake",
	)
	suite.NoError(err)

	input, err := helpers.PackMethodCall(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	suite.NoError(err)
	suite.NotNil(result)

	// Unpack result
	unpacked, err := method.Outputs.Unpack(result)
	suite.NoError(err)
	suite.Len(unpacked, 1)

	success, ok := unpacked[0].(bool)
	suite.True(ok)
	suite.True(success, "Send should succeed")

	// Verify balances changed
	aliceFinalBalance := suite.TestSuite.App.BankKeeper.GetBalance(suite.TestSuite.Ctx, suite.TestSuite.Alice, "stake")
	bobFinalBalance := suite.TestSuite.App.BankKeeper.GetBalance(suite.TestSuite.Ctx, suite.TestSuite.Bob, "stake")

	// Alice's balance should have decreased
	suite.True(aliceFinalBalance.Amount.LT(aliceInitialBalance.Amount), "Alice's balance should decrease")

	// Bob's balance should have increased
	suite.True(bobFinalBalance.Amount.GT(bobInitialBalance.Amount), "Bob's balance should increase")

	// Verify the amount transferred
	expectedAmount, ok := sdkmath.NewIntFromString(sendAmount)
	suite.True(ok, "Should parse amount string")
	actualTransfer := bobFinalBalance.Amount.Sub(bobInitialBalance.Amount)
	suite.Equal(expectedAmount.String(), actualTransfer.String(), "Transfer amount should match")
}

func (suite *IntegrationTestSuite) TestSend_MultipleTransfers() {
	caller := suite.TestSuite.AliceEVM

	method := suite.Precompile.ABI.Methods["send"]
	require.NotNil(suite.T(), method)

	// Send to Bob
	jsonMsg1, err := helpers.BuildSendJSON(
		suite.TestSuite.Alice.String(),
		suite.TestSuite.Bob.String(),
		"500",
		"stake",
	)
	suite.NoError(err)

	input1, err := helpers.PackMethodCall(&method, jsonMsg1)
	suite.NoError(err)

	contract1 := suite.TestSuite.CreateMockContract(caller, input1)
	result1, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract1, false)
	suite.NoError(err)
	suite.NotNil(result1)

	// Send to Charlie
	jsonMsg2, err := helpers.BuildSendJSON(
		suite.TestSuite.Alice.String(),
		suite.TestSuite.Charlie.String(),
		"300",
		"stake",
	)
	suite.NoError(err)

	input2, err := helpers.PackMethodCall(&method, jsonMsg2)
	suite.NoError(err)

	contract2 := suite.TestSuite.CreateMockContract(caller, input2)
	result2, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract2, false)
	suite.NoError(err)
	suite.NotNil(result2)

	// Verify all balances changed correctly
	aliceBalance := suite.TestSuite.App.BankKeeper.GetBalance(suite.TestSuite.Ctx, suite.TestSuite.Alice, "stake")
	bobBalance := suite.TestSuite.App.BankKeeper.GetBalance(suite.TestSuite.Ctx, suite.TestSuite.Bob, "stake")
	charlieBalance := suite.TestSuite.App.BankKeeper.GetBalance(suite.TestSuite.Ctx, suite.TestSuite.Charlie, "stake")

	// All balances should be non-zero
	suite.True(aliceBalance.Amount.GT(sdkmath.ZeroInt()), "Alice should have remaining balance")
	suite.True(bobBalance.Amount.GT(sdkmath.ZeroInt()), "Bob should have received tokens")
	suite.True(charlieBalance.Amount.GT(sdkmath.ZeroInt()), "Charlie should have received tokens")
}

func (suite *IntegrationTestSuite) TestSend_DifferentDenoms() {
	caller := suite.TestSuite.AliceEVM
	toAddress := suite.TestSuite.Bob.String()

	method := suite.Precompile.ABI.Methods["send"]
	require.NotNil(suite.T(), method)

	// Send stake
	jsonMsg1, err := helpers.BuildSendJSON(
		suite.TestSuite.Alice.String(),
		toAddress,
		"1000",
		"stake",
	)
	suite.NoError(err)

	input1, err := helpers.PackMethodCall(&method, jsonMsg1)
	suite.NoError(err)

	contract1 := suite.TestSuite.CreateMockContract(caller, input1)
	result1, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract1, false)
	suite.NoError(err)
	suite.NotNil(result1)

	// Send uosmo
	jsonMsg2, err := helpers.BuildSendJSON(
		suite.TestSuite.Alice.String(),
		toAddress,
		"2000",
		"uosmo",
	)
	suite.NoError(err)

	input2, err := helpers.PackMethodCall(&method, jsonMsg2)
	suite.NoError(err)

	contract2 := suite.TestSuite.CreateMockContract(caller, input2)
	result2, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract2, false)
	suite.NoError(err)
	suite.NotNil(result2)

	// Verify both denoms were transferred
	bobStakeBalance := suite.TestSuite.App.BankKeeper.GetBalance(suite.TestSuite.Ctx, suite.TestSuite.Bob, "stake")
	bobUosmoBalance := suite.TestSuite.App.BankKeeper.GetBalance(suite.TestSuite.Ctx, suite.TestSuite.Bob, "uosmo")

	suite.True(bobStakeBalance.Amount.GT(sdkmath.ZeroInt()), "Bob should have stake")
	suite.True(bobUosmoBalance.Amount.GT(sdkmath.ZeroInt()), "Bob should have uosmo")
}

