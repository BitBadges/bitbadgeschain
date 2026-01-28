package routing

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/ai_test/testutil"
)

type PrefixPriorityTestSuite struct {
	testutil.AITestSuite
}

func TestPrefixPriorityTestSuite(t *testing.T) {
	suite.Run(t, new(PrefixPriorityTestSuite))
}

func (suite *PrefixPriorityTestSuite) TestPrefixPriority_LongestPrefixWins() {
	// Note: This test demonstrates that overlapping prefixes are prevented
	// If we had "tokenization:" and "tokenization:lp:", the longer one should win
	// But since overlapping is prevented, we test with non-overlapping prefixes
	
	router1 := testutil.GenerateMockRouter("tokenization:")
	router2 := testutil.GenerateMockRouter("badgeslp:")

	err := suite.Keeper.RegisterRouter("tokenization:", router1)
	suite.Require().NoError(err)

	err = suite.Keeper.RegisterRouter("badgeslp:", router2)
	suite.Require().NoError(err)

	aliceAddr, err := sdk.AccAddressFromBech32(suite.Alice)
	suite.Require().NoError(err)
	bobAddr, err := sdk.AccAddressFromBech32(suite.Bob)
	suite.Require().NoError(err)

	// Send tokenization denom - should route to router1
	coin1 := sdk.NewCoin("tokenization:123:456", sdkmath.NewInt(1000))
	err = suite.Keeper.SendCoinWithAliasRouting(suite.Ctx, aliceAddr, bobAddr, &coin1)
	suite.Require().NoError(err)

	// Send badgeslp denom - should route to router2
	coin2 := sdk.NewCoin("badgeslp:789:012", sdkmath.NewInt(500))
	err = suite.Keeper.SendCoinWithAliasRouting(suite.Ctx, aliceAddr, bobAddr, &coin2)
	suite.Require().NoError(err)

	// Verify routers were called correctly
	calls1 := router1.GetSendCalls()
	suite.Require().Len(calls1, 1)
	suite.Require().Equal("tokenization:123:456", calls1[0].Denom)

	calls2 := router2.GetSendCalls()
	suite.Require().Len(calls2, 1)
	suite.Require().Equal("badgeslp:789:012", calls2[0].Denom)
}

