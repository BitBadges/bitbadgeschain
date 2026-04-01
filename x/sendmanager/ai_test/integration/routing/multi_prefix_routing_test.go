package routing

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/sendmanager/keeper"
)

type MultiPrefixRoutingTestSuite struct {
	testutil.AITestSuite
}

func TestMultiPrefixRoutingTestSuite(t *testing.T) {
	suite.Run(t, new(MultiPrefixRoutingTestSuite))
}

func (suite *MultiPrefixRoutingTestSuite) TestMultiPrefixRouting_OnlySinglePrefixSupported() {
	// Only "badgeslp:" prefix is supported; registering other prefixes must fail
	router := testutil.GenerateMockRouter(keeper.AliasDenomPrefix)

	err := suite.Keeper.RegisterRouter(keeper.AliasDenomPrefix, router)
	suite.Require().NoError(err)

	// Verify only one prefix is registered
	prefixes := suite.Keeper.GetRegisteredPrefixes()
	suite.Require().Len(prefixes, 1)
	suite.Require().Contains(prefixes, keeper.AliasDenomPrefix)

	// Attempting to register a different prefix should fail
	router2 := testutil.GenerateMockRouter("tokens:")
	err = suite.Keeper.RegisterRouter("tokens:", router2)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "only prefix")
}

func (suite *MultiPrefixRoutingTestSuite) TestMultiPrefixRouting_BadgeslpRoutedCorrectly() {
	router := testutil.GenerateMockRouter(keeper.AliasDenomPrefix)

	err := suite.Keeper.RegisterRouter(keeper.AliasDenomPrefix, router)
	suite.Require().NoError(err)

	aliceAddr, err := sdk.AccAddressFromBech32(suite.Alice)
	suite.Require().NoError(err)
	bobAddr, err := sdk.AccAddressFromBech32(suite.Bob)
	suite.Require().NoError(err)

	// Send badgeslp denom - should route through the alias router
	coin := sdk.NewCoin("badgeslp:123:456", sdkmath.NewInt(1000))
	err = suite.Keeper.SendCoinWithAliasRouting(suite.Ctx, aliceAddr, bobAddr, &coin)
	suite.Require().NoError(err)

	// Verify router was called
	calls := router.GetSendCalls()
	suite.Require().Len(calls, 1)
	suite.Require().Equal("badgeslp:123:456", calls[0].Denom)
}

func (suite *MultiPrefixRoutingTestSuite) TestMultiPrefixRouting_NonBadgeslpDenomNotRouted() {
	router := testutil.GenerateMockRouter(keeper.AliasDenomPrefix)

	err := suite.Keeper.RegisterRouter(keeper.AliasDenomPrefix, router)
	suite.Require().NoError(err)

	aliceAddr, err := sdk.AccAddressFromBech32(suite.Alice)
	suite.Require().NoError(err)
	bobAddr, err := sdk.AccAddressFromBech32(suite.Bob)
	suite.Require().NoError(err)

	// A denom with an unsupported prefix should NOT be routed through the alias router
	coin := sdk.NewCoin("tokens:789:012", sdkmath.NewInt(500))
	err = suite.Keeper.SendCoinWithAliasRouting(suite.Ctx, aliceAddr, bobAddr, &coin)
	// Falls through to bank, which has no balance
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "insufficient funds")

	// Router should not have been called
	calls := router.GetSendCalls()
	suite.Require().Len(calls, 0)
}
