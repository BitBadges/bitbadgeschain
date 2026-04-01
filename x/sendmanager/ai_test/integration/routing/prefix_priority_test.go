package routing

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/sendmanager/keeper"
)

type PrefixPriorityTestSuite struct {
	testutil.AITestSuite
}

func TestPrefixPriorityTestSuite(t *testing.T) {
	suite.Run(t, new(PrefixPriorityTestSuite))
}

func (suite *PrefixPriorityTestSuite) TestPrefixPriority_OnlyBadgeslpAccepted() {
	// Only "badgeslp:" prefix is supported; other prefixes are rejected.
	router := testutil.GenerateMockRouter(keeper.AliasDenomPrefix)

	err := suite.Keeper.RegisterRouter(keeper.AliasDenomPrefix, router)
	suite.Require().NoError(err)

	// Attempting a non-badgeslp prefix should fail
	router2 := testutil.GenerateMockRouter("badges:")
	err = suite.Keeper.RegisterRouter("badges:", router2)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "only prefix")
}

func (suite *PrefixPriorityTestSuite) TestPrefixPriority_BadgeslpDenomRoutedCorrectly() {
	router := testutil.GenerateMockRouter(keeper.AliasDenomPrefix)

	err := suite.Keeper.RegisterRouter(keeper.AliasDenomPrefix, router)
	suite.Require().NoError(err)

	aliceAddr, err := sdk.AccAddressFromBech32(suite.Alice)
	suite.Require().NoError(err)
	bobAddr, err := sdk.AccAddressFromBech32(suite.Bob)
	suite.Require().NoError(err)

	// Send badgeslp denom - should route through the alias router
	coin := sdk.NewCoin("badgeslp:789:012", sdkmath.NewInt(500))
	err = suite.Keeper.SendCoinWithAliasRouting(suite.Ctx, aliceAddr, bobAddr, &coin)
	suite.Require().NoError(err)

	calls := router.GetSendCalls()
	suite.Require().Len(calls, 1)
	suite.Require().Equal("badgeslp:789:012", calls[0].Denom)
}

func (suite *PrefixPriorityTestSuite) TestPrefixPriority_NonBadgeslpDenomFallsToBank() {
	router := testutil.GenerateMockRouter(keeper.AliasDenomPrefix)

	err := suite.Keeper.RegisterRouter(keeper.AliasDenomPrefix, router)
	suite.Require().NoError(err)

	aliceAddr, err := sdk.AccAddressFromBech32(suite.Alice)
	suite.Require().NoError(err)
	bobAddr, err := sdk.AccAddressFromBech32(suite.Bob)
	suite.Require().NoError(err)

	// A denom that does NOT start with "badgeslp:" should route through bank
	coin := sdk.NewCoin("uatom", sdkmath.NewInt(1000))
	err = suite.Keeper.SendCoinWithAliasRouting(suite.Ctx, aliceAddr, bobAddr, &coin)
	// Bank keeper fails with insufficient funds — that's expected, it means routing went to bank
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "insufficient funds")

	// Router should NOT have been called
	calls := router.GetSendCalls()
	suite.Require().Len(calls, 0)
}
