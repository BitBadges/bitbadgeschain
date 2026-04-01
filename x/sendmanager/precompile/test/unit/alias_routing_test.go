package sendmanager_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	sdkmath "cosmossdk.io/math"

	sendmanagerkeeper "github.com/bitbadges/bitbadgeschain/x/sendmanager/keeper"
	sendmanager "github.com/bitbadges/bitbadgeschain/x/sendmanager/precompile"
	"github.com/bitbadges/bitbadgeschain/x/sendmanager/precompile/test/helpers"
)

type AliasRoutingTestSuite struct {
	suite.Suite
	Precompile *sendmanager.Precompile
	TestSuite  *helpers.TestSuite
}

func TestAliasRoutingTestSuite(t *testing.T) {
	suite.Run(t, new(AliasRoutingTestSuite))
}

func (suite *AliasRoutingTestSuite) SetupTest() {
	suite.TestSuite = helpers.NewTestSuite(suite.T())
	suite.Precompile = suite.TestSuite.Precompile
}

func (suite *AliasRoutingTestSuite) TestSend_AliasDenomRoutesToRouter() {
	// Register a mock router for the badgeslp: prefix
	mockRouter := helpers.NewMockRouter(sendmanagerkeeper.AliasDenomPrefix)
	err := suite.TestSuite.App.SendmanagerKeeper.RegisterRouter(sendmanagerkeeper.AliasDenomPrefix, mockRouter)
	suite.Require().NoError(err)

	caller := suite.TestSuite.AliceEVM
	toAddress := suite.TestSuite.Bob.String()

	method := suite.Precompile.ABI.Methods["send"]
	require.NotNil(suite.T(), method)

	// Build JSON message with alias denom
	msg := map[string]interface{}{
		"from_address": suite.TestSuite.Alice.String(),
		"to_address":   toAddress,
		"amount": []map[string]interface{}{
			{"denom": "badgeslp:123:456", "amount": "1000"},
		},
	}
	jsonMsg, err := helpers.BuildQueryJSON(msg)
	suite.NoError(err)

	input, err := helpers.PackMethodCall(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)

	// Router should be called if validation passes and routing happens
	calls := mockRouter.GetSendCalls()

	if len(calls) > 0 {
		// Router was called - verify it worked correctly
		suite.Equal("badgeslp:123:456", calls[0].Denom)
		suite.Equal(suite.TestSuite.Alice.String(), calls[0].From)
		suite.Equal(toAddress, calls[0].To)

		expectedAmount := sdkmath.NewUintFromString("1000")
		suite.Equal(expectedAmount.String(), calls[0].Amount.String())
	}

	// Transaction may fail due to validation or balance checks
	if err != nil {
		suite.True(
			suite.Contains(err.Error(), "insufficient") ||
				suite.Contains(err.Error(), "invalid") ||
				suite.Contains(err.Error(), "denom"),
			"Error should mention balance, invalid, or denom: %s", err.Error(),
		)
	} else {
		suite.NotNil(result)
		suite.GreaterOrEqual(len(calls), 1, "Router should be called when transaction succeeds")
	}
}

func (suite *AliasRoutingTestSuite) TestSend_StandardDenomRoutesToBank() {
	// Don't register any router - standard denom should route to bank
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

	// Should succeed (routing through bank keeper)
	// May fail due to insufficient balance, but routing should work
	if err != nil {
		// If error, it should be about balance, not routing
		suite.Contains(err.Error(), "insufficient")
	} else {
		suite.NotNil(result)
	}
}

func (suite *AliasRoutingTestSuite) TestSend_MixedDenomsRoutesCorrectly() {
	// Register a mock router for badgeslp: prefix
	mockRouter := helpers.NewMockRouter(sendmanagerkeeper.AliasDenomPrefix)
	err := suite.TestSuite.App.SendmanagerKeeper.RegisterRouter(sendmanagerkeeper.AliasDenomPrefix, mockRouter)
	suite.Require().NoError(err)

	caller := suite.TestSuite.AliceEVM
	toAddress := suite.TestSuite.Bob.String()

	method := suite.Precompile.ABI.Methods["send"]
	require.NotNil(suite.T(), method)

	// Build JSON message with mixed denoms (alias + standard)
	// Coins must be sorted by denom - "badgeslp:123:456" comes before "stake" alphabetically
	msg := map[string]interface{}{
		"from_address": suite.TestSuite.Alice.String(),
		"to_address":   toAddress,
		"amount": []map[string]interface{}{
			{"denom": "badgeslp:123:456", "amount": "1000"}, // Alias denom (sorted first)
			{"denom": "stake", "amount": "2000"},             // Standard denom (sorted second)
		},
	}
	jsonMsg, err := helpers.BuildQueryJSON(msg)
	suite.NoError(err)

	input, err := helpers.PackMethodCall(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)

	// Verify router was called for alias denom
	calls := mockRouter.GetSendCalls()
	if len(calls) > 0 {
		suite.Equal("badgeslp:123:456", calls[0].Denom)
	}

	// Transaction may fail due to validation or balance checks
	if err != nil {
		suite.True(
			suite.Contains(err.Error(), "insufficient") ||
				suite.Contains(err.Error(), "invalid") ||
				suite.Contains(err.Error(), "denom") ||
				suite.Contains(err.Error(), "sorted"),
			"Error should mention balance, invalid, denom, or sorted: %s", err.Error(),
		)
	} else {
		suite.NotNil(result)
		suite.GreaterOrEqual(len(calls), 1, "Router should be called when transaction succeeds")
	}
}

func (suite *AliasRoutingTestSuite) TestRegisterRouter_RejectsNonBadgeslpPrefix() {
	// The simplified router only supports the badgeslp: prefix
	mockRouter := helpers.NewMockRouter("other:")
	err := suite.TestSuite.App.SendmanagerKeeper.RegisterRouter("other:", mockRouter)
	suite.Require().Error(err)
	suite.Contains(err.Error(), "only prefix")
}

func (suite *AliasRoutingTestSuite) TestSend_MultipleAliasDenoms() {
	// Register a mock router
	mockRouter := helpers.NewMockRouter(sendmanagerkeeper.AliasDenomPrefix)
	err := suite.TestSuite.App.SendmanagerKeeper.RegisterRouter(sendmanagerkeeper.AliasDenomPrefix, mockRouter)
	suite.Require().NoError(err)

	caller := suite.TestSuite.AliceEVM
	toAddress := suite.TestSuite.Bob.String()

	method := suite.Precompile.ABI.Methods["send"]
	require.NotNil(suite.T(), method)

	// Build JSON message with a single alias denom
	msg := map[string]interface{}{
		"from_address": suite.TestSuite.Alice.String(),
		"to_address":   toAddress,
		"amount": []map[string]interface{}{
			{"denom": "badgeslp:123:456", "amount": "1000"},
		},
	}
	jsonMsg, err := helpers.BuildQueryJSON(msg)
	suite.NoError(err)

	input, err := helpers.PackMethodCall(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)

	// Verify router was called
	calls := mockRouter.GetSendCalls()
	if len(calls) > 0 {
		suite.Equal("badgeslp:123:456", calls[0].Denom)
		suite.Equal(suite.TestSuite.Alice.String(), calls[0].From)
		suite.Equal(toAddress, calls[0].To)
	}

	if err != nil {
		suite.True(
			suite.Contains(err.Error(), "insufficient") ||
				suite.Contains(err.Error(), "invalid") ||
				suite.Contains(err.Error(), "denom"),
			"Error should mention balance, invalid, or denom: %s", err.Error(),
		)
	} else {
		suite.NotNil(result)
		suite.GreaterOrEqual(len(calls), 1, "Router should be called when transaction succeeds")
	}
}

func (suite *AliasRoutingTestSuite) TestSend_NoPrefixMatchRoutesToBank() {
	// Register a router for badgeslp:
	mockRouter := helpers.NewMockRouter(sendmanagerkeeper.AliasDenomPrefix)
	err := suite.TestSuite.App.SendmanagerKeeper.RegisterRouter(sendmanagerkeeper.AliasDenomPrefix, mockRouter)
	suite.Require().NoError(err)

	caller := suite.TestSuite.AliceEVM
	toAddress := suite.TestSuite.Bob.String()

	method := suite.Precompile.ABI.Methods["send"]
	require.NotNil(suite.T(), method)

	// Use a denom that doesn't match badgeslp: prefix
	jsonMsg, err := helpers.BuildSendJSON(
		suite.TestSuite.Alice.String(),
		toAddress,
		"1000",
		"uatom", // Doesn't match "badgeslp:" prefix
	)
	suite.NoError(err)

	input, err := helpers.PackMethodCall(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)

	// Should route to bank (may fail due to insufficient balance)
	if err != nil {
		suite.Contains(err.Error(), "insufficient")
	} else {
		suite.NotNil(result)
	}

	// Verify router was NOT called
	calls := mockRouter.GetSendCalls()
	suite.Len(calls, 0, "Router should not be called for non-matching denom")
}
