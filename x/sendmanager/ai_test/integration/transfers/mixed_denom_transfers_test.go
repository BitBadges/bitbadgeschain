package transfers

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/ai_test/testutil"
	sendmanagerkeeper "github.com/bitbadges/bitbadgeschain/x/sendmanager/keeper"
)

type MixedDenomTransfersTestSuite struct {
	testutil.AITestSuite
}

func TestMixedDenomTransfersTestSuite(t *testing.T) {
	suite.Run(t, new(MixedDenomTransfersTestSuite))
}

func (suite *MixedDenomTransfersTestSuite) TestMixedDenomTransfers_AliasAndBank() {
	router := testutil.GenerateMockRouter(sendmanagerkeeper.AliasDenomPrefix)
	err := suite.Keeper.RegisterRouter(sendmanagerkeeper.AliasDenomPrefix, router)
	suite.Require().NoError(err)

	aliceAddr, err := sdk.AccAddressFromBech32(suite.Alice)
	suite.Require().NoError(err)
	bobAddr, err := sdk.AccAddressFromBech32(suite.Bob)
	suite.Require().NoError(err)

	coins := sdk.Coins{
		sdk.NewCoin("badgeslp:123:456", sdkmath.NewInt(1000)),
		sdk.NewCoin("uatom", sdkmath.NewInt(500)),
	}

	// Mixed denoms - alias denom succeeds, but bank denom fails due to insufficient funds
	err = suite.Keeper.SendCoinsWithAliasRouting(suite.Ctx, aliceAddr, bobAddr, coins)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "insufficient funds")

	// Verify router was called for alias denom (before bank denom fails)
	calls := router.GetSendCalls()
	suite.Require().Len(calls, 1)
	suite.Require().Equal("badgeslp:123:456", calls[0].Denom)
}

func (suite *MixedDenomTransfersTestSuite) TestMixedDenomTransfers_MultipleAliasDenoms() {
	router := testutil.GenerateMockRouter(sendmanagerkeeper.AliasDenomPrefix)

	err := suite.Keeper.RegisterRouter(sendmanagerkeeper.AliasDenomPrefix, router)
	suite.Require().NoError(err)

	aliceAddr, err := sdk.AccAddressFromBech32(suite.Alice)
	suite.Require().NoError(err)
	bobAddr, err := sdk.AccAddressFromBech32(suite.Bob)
	suite.Require().NoError(err)

	// Both denoms use the single supported prefix
	coins := sdk.Coins{
		sdk.NewCoin("badgeslp:123:456", sdkmath.NewInt(1000)),
		sdk.NewCoin("badgeslp:789:012", sdkmath.NewInt(500)),
	}

	err = suite.Keeper.SendCoinsWithAliasRouting(suite.Ctx, aliceAddr, bobAddr, coins)
	suite.Require().NoError(err)

	// Verify router was called for both alias denoms
	calls := router.GetSendCalls()
	suite.Require().Len(calls, 2)
}
