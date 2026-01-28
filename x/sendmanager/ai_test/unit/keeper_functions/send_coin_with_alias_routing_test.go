package keeper_functions

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/ai_test/testutil"
)

type SendCoinWithAliasRoutingTestSuite struct {
	testutil.AITestSuite
}

func TestSendCoinWithAliasRoutingTestSuite(t *testing.T) {
	suite.Run(t, new(SendCoinWithAliasRoutingTestSuite))
}

func (suite *SendCoinWithAliasRoutingTestSuite) TestSendCoinWithAliasRouting_AliasDenom() {
	router := testutil.GenerateMockRouter("tokenization:")
	err := suite.Keeper.RegisterRouter("tokenization:", router)
	suite.Require().NoError(err)

	aliceAddr, err := sdk.AccAddressFromBech32(suite.Alice)
	suite.Require().NoError(err)
	bobAddr, err := sdk.AccAddressFromBech32(suite.Bob)
	suite.Require().NoError(err)

	coin := sdk.NewCoin("tokenization:123:456", sdkmath.NewInt(1000))
	err = suite.Keeper.SendCoinWithAliasRouting(suite.Ctx, aliceAddr, bobAddr, &coin)
	suite.Require().NoError(err)

	// Verify router was called
	calls := router.GetSendCalls()
	suite.Require().Len(calls, 1)
	suite.Require().Equal("tokenization:123:456", calls[0].Denom)
}

func (suite *SendCoinWithAliasRoutingTestSuite) TestSendCoinWithAliasRouting_BankDenom() {
	aliceAddr, err := sdk.AccAddressFromBech32(suite.Alice)
	suite.Require().NoError(err)
	bobAddr, err := sdk.AccAddressFromBech32(suite.Bob)
	suite.Require().NoError(err)

	// Bank denom should route through bank keeper (no router registered)
	// The mock bank keeper will return insufficient funds error, which is expected behavior
	coin := sdk.NewCoin("uatom", sdkmath.NewInt(500))
	err = suite.Keeper.SendCoinWithAliasRouting(suite.Ctx, aliceAddr, bobAddr, &coin)
	// Mock bank keeper checks balances, so this will fail without balance setup
	// This is actually correct behavior - we're testing that routing works
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "insufficient funds")
}

func (suite *SendCoinWithAliasRoutingTestSuite) TestSendCoinWithAliasRouting_EmptyDenom() {
	aliceAddr, err := sdk.AccAddressFromBech32(suite.Alice)
	suite.Require().NoError(err)
	bobAddr, err := sdk.AccAddressFromBech32(suite.Bob)
	suite.Require().NoError(err)

	// Create coin with empty denom - this will panic in sdk.NewCoin, so we test the keeper's validation
	// Instead, we create a coin struct directly to test the keeper's empty denom check
	coin := sdk.Coin{Denom: "", Amount: sdkmath.NewInt(1000)}
	err = suite.Keeper.SendCoinWithAliasRouting(suite.Ctx, aliceAddr, bobAddr, &coin)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "cannot be empty")
}

func (suite *SendCoinWithAliasRoutingTestSuite) TestSendCoinWithAliasRouting_MixedDenoms() {
	router := testutil.GenerateMockRouter("tokenization:")
	err := suite.Keeper.RegisterRouter("tokenization:", router)
	suite.Require().NoError(err)

	aliceAddr, err := sdk.AccAddressFromBech32(suite.Alice)
	suite.Require().NoError(err)
	bobAddr, err := sdk.AccAddressFromBech32(suite.Bob)
	suite.Require().NoError(err)

	// Test alias denom
	aliasCoin := sdk.NewCoin("tokenization:123:456", sdkmath.NewInt(1000))
	err = suite.Keeper.SendCoinWithAliasRouting(suite.Ctx, aliceAddr, bobAddr, &aliasCoin)
	suite.Require().NoError(err)

	// Test bank denom - will fail without balance, which is expected
	bankCoin := sdk.NewCoin("uatom", sdkmath.NewInt(500))
	err = suite.Keeper.SendCoinWithAliasRouting(suite.Ctx, aliceAddr, bobAddr, &bankCoin)
	// Mock bank keeper checks balances, so this will fail - this is correct behavior
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "insufficient funds")
}
