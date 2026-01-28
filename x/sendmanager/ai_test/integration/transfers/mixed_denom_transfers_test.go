package transfers

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/ai_test/testutil"
)

type MixedDenomTransfersTestSuite struct {
	testutil.AITestSuite
}

func TestMixedDenomTransfersTestSuite(t *testing.T) {
	suite.Run(t, new(MixedDenomTransfersTestSuite))
}

func (suite *MixedDenomTransfersTestSuite) TestMixedDenomTransfers_AliasAndBank() {
	router := testutil.GenerateMockRouter("tokenization:")
	err := suite.Keeper.RegisterRouter("tokenization:", router)
	suite.Require().NoError(err)

	aliceAddr, err := sdk.AccAddressFromBech32(suite.Alice)
	suite.Require().NoError(err)
	bobAddr, err := sdk.AccAddressFromBech32(suite.Bob)
	suite.Require().NoError(err)

	coins := sdk.Coins{
		sdk.NewCoin("tokenization:123:456", sdkmath.NewInt(1000)),
		sdk.NewCoin("uatom", sdkmath.NewInt(500)),
	}

	// Mixed denoms - alias denom succeeds, but bank denom fails due to insufficient funds
	err = suite.Keeper.SendCoinsWithAliasRouting(suite.Ctx, aliceAddr, bobAddr, coins)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "insufficient funds")

	// Verify router was called for alias denom (before bank denom fails)
	calls := router.GetSendCalls()
	suite.Require().Len(calls, 1)
	suite.Require().Equal("tokenization:123:456", calls[0].Denom)
}

func (suite *MixedDenomTransfersTestSuite) TestMixedDenomTransfers_MultipleAliasDenoms() {
	router1 := testutil.GenerateMockRouter("tokenization:")
	router2 := testutil.GenerateMockRouter("tokens:")

	err := suite.Keeper.RegisterRouter("tokenization:", router1)
	suite.Require().NoError(err)

	err = suite.Keeper.RegisterRouter("tokens:", router2)
	suite.Require().NoError(err)

	aliceAddr, err := sdk.AccAddressFromBech32(suite.Alice)
	suite.Require().NoError(err)
	bobAddr, err := sdk.AccAddressFromBech32(suite.Bob)
	suite.Require().NoError(err)

	coins := sdk.Coins{
		sdk.NewCoin("tokenization:123:456", sdkmath.NewInt(1000)),
		sdk.NewCoin("tokens:789:012", sdkmath.NewInt(500)),
	}

	err = suite.Keeper.SendCoinsWithAliasRouting(suite.Ctx, aliceAddr, bobAddr, coins)
	suite.Require().NoError(err)

	// Verify both routers were called
	calls1 := router1.GetSendCalls()
	suite.Require().Len(calls1, 1)

	calls2 := router2.GetSendCalls()
	suite.Require().Len(calls2, 1)
}

