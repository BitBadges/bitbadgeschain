package routing

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/ai_test/testutil"
)

type MultiPrefixRoutingTestSuite struct {
	testutil.AITestSuite
}

func TestMultiPrefixRoutingTestSuite(t *testing.T) {
	suite.Run(t, new(MultiPrefixRoutingTestSuite))
}

func (suite *MultiPrefixRoutingTestSuite) TestMultiPrefixRouting_RegisterMultiplePrefixes() {
	router1 := testutil.GenerateMockRouter("badges:")
	router2 := testutil.GenerateMockRouter("tokens:")

	err := suite.Keeper.RegisterRouter("badges:", router1)
	suite.Require().NoError(err)

	err = suite.Keeper.RegisterRouter("tokens:", router2)
	suite.Require().NoError(err)

	// Verify both prefixes are registered
	prefixes := suite.Keeper.GetRegisteredPrefixes()
	suite.Require().Contains(prefixes, "badges:")
	suite.Require().Contains(prefixes, "tokens:")
}

func (suite *MultiPrefixRoutingTestSuite) TestMultiPrefixRouting_RouteToCorrectRouter() {
	router1 := testutil.GenerateMockRouter("badges:")
	router2 := testutil.GenerateMockRouter("tokens:")

	err := suite.Keeper.RegisterRouter("badges:", router1)
	suite.Require().NoError(err)

	err = suite.Keeper.RegisterRouter("tokens:", router2)
	suite.Require().NoError(err)

	aliceAddr, err := sdk.AccAddressFromBech32(suite.Alice)
	suite.Require().NoError(err)
	bobAddr, err := sdk.AccAddressFromBech32(suite.Bob)
	suite.Require().NoError(err)

	// Send badges denom - should route to router1
	coin1 := sdk.NewCoin("badges:123:456", sdkmath.NewInt(1000))
	err = suite.Keeper.SendCoinWithAliasRouting(suite.Ctx, aliceAddr, bobAddr, &coin1)
	suite.Require().NoError(err)

	// Send tokens denom - should route to router2
	coin2 := sdk.NewCoin("tokens:789:012", sdkmath.NewInt(500))
	err = suite.Keeper.SendCoinWithAliasRouting(suite.Ctx, aliceAddr, bobAddr, &coin2)
	suite.Require().NoError(err)

	// Verify routers were called
	calls1 := router1.GetSendCalls()
	suite.Require().Len(calls1, 1)
	suite.Require().Equal("badges:123:456", calls1[0].Denom)

	calls2 := router2.GetSendCalls()
	suite.Require().Len(calls2, 1)
	suite.Require().Equal("tokens:789:012", calls2[0].Denom)
}

