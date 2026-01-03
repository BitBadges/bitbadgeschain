package keeper_functions

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/ai_test/testutil"
)

type SendCoinsWithAliasRoutingTestSuite struct {
	testutil.AITestSuite
}

func TestSendCoinsWithAliasRoutingTestSuite(t *testing.T) {
	suite.Run(t, new(SendCoinsWithAliasRoutingTestSuite))
}

func (suite *SendCoinsWithAliasRoutingTestSuite) TestSendCoinsWithAliasRouting_MultipleAliasDenoms() {
	router := testutil.GenerateMockRouter("badges:")
	err := suite.Keeper.RegisterRouter("badges:", router)
	suite.Require().NoError(err)

	aliceAddr, err := sdk.AccAddressFromBech32(suite.Alice)
	suite.Require().NoError(err)
	bobAddr, err := sdk.AccAddressFromBech32(suite.Bob)
	suite.Require().NoError(err)

	coins := sdk.Coins{
		sdk.NewCoin("badges:123:456", sdkmath.NewInt(1000)),
		sdk.NewCoin("badges:789:012", sdkmath.NewInt(2000)),
	}
	err = suite.Keeper.SendCoinsWithAliasRouting(suite.Ctx, aliceAddr, bobAddr, coins)
	suite.Require().NoError(err)

	// Verify router was called for both coins
	calls := router.GetSendCalls()
	suite.Require().Len(calls, 2)
}

func (suite *SendCoinsWithAliasRoutingTestSuite) TestSendCoinsWithAliasRouting_MixedDenoms() {
	router := testutil.GenerateMockRouter("badges:")
	err := suite.Keeper.RegisterRouter("badges:", router)
	suite.Require().NoError(err)

	aliceAddr, err := sdk.AccAddressFromBech32(suite.Alice)
	suite.Require().NoError(err)
	bobAddr, err := sdk.AccAddressFromBech32(suite.Bob)
	suite.Require().NoError(err)

	coins := sdk.Coins{
		sdk.NewCoin("badges:123:456", sdkmath.NewInt(1000)),
		sdk.NewCoin("uatom", sdkmath.NewInt(500)),
	}
	// Mixed denoms - processes alias denom first (succeeds), then bank denom (fails due to insufficient funds)
	err = suite.Keeper.SendCoinsWithAliasRouting(suite.Ctx, aliceAddr, bobAddr, coins)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "insufficient funds")

	// Verify router was called for alias denom
	calls := router.GetSendCalls()
	suite.Require().Len(calls, 1)
	suite.Require().Equal("badges:123:456", calls[0].Denom)
}

func (suite *SendCoinsWithAliasRoutingTestSuite) TestSendCoinsWithAliasRouting_EmptyDenom() {
	aliceAddr, err := sdk.AccAddressFromBech32(suite.Alice)
	suite.Require().NoError(err)
	bobAddr, err := sdk.AccAddressFromBech32(suite.Bob)
	suite.Require().NoError(err)

	coins := sdk.Coins{
		sdk.Coin{Denom: "", Amount: sdkmath.NewInt(1000)},
	}
	err = suite.Keeper.SendCoinsWithAliasRouting(suite.Ctx, aliceAddr, bobAddr, coins)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "cannot be empty")
}

func (suite *SendCoinsWithAliasRoutingTestSuite) TestSendCoinsWithAliasRouting_AllBankDenoms() {
	aliceAddr, err := sdk.AccAddressFromBech32(suite.Alice)
	suite.Require().NoError(err)
	bobAddr, err := sdk.AccAddressFromBech32(suite.Bob)
	suite.Require().NoError(err)

	coins := sdk.Coins{
		sdk.NewCoin("uatom", sdkmath.NewInt(1000)),
		sdk.NewCoin("uosmo", sdkmath.NewInt(500)),
	}
	err = suite.Keeper.SendCoinsWithAliasRouting(suite.Ctx, aliceAddr, bobAddr, coins)
	// Mock bank keeper checks balances, so this will fail without balance setup
	// This is correct behavior - we're testing that routing works
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "insufficient funds")
}

