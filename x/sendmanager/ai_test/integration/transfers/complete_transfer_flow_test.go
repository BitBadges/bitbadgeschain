package transfers

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/sendmanager/keeper"
	"github.com/bitbadges/bitbadgeschain/x/sendmanager/types"
)

type CompleteTransferFlowTestSuite struct {
	testutil.AITestSuite
}

func TestCompleteTransferFlowTestSuite(t *testing.T) {
	suite.Run(t, new(CompleteTransferFlowTestSuite))
}

func (suite *CompleteTransferFlowTestSuite) TestCompleteTransferFlow_AliasDenom() {
	router := testutil.GenerateMockRouter("tokenization:")
	err := suite.Keeper.RegisterRouter("tokenization:", router)
	suite.Require().NoError(err)

	// Recreate msg server after registering router (msgServer embeds Keeper by value)
	suite.MsgServer = keeper.NewMsgServerImpl(suite.Keeper)

	msg := &types.MsgSendWithAliasRouting{
		FromAddress: suite.Alice,
		ToAddress:   suite.Bob,
		Amount: sdk.Coins{
			sdk.NewCoin("tokenization:123:456", sdkmath.NewInt(1000)),
		},
	}

	resp, err := suite.MsgServer.SendWithAliasRouting(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err)
	suite.Require().NotNil(resp)

	// Verify router was called
	calls := router.GetSendCalls()
	suite.Require().Len(calls, 1)
}

func (suite *CompleteTransferFlowTestSuite) TestCompleteTransferFlow_BankDenom() {
	// Pre-fund Alice's account
	aliceAddr, err := sdk.AccAddressFromBech32(suite.Alice)
	suite.Require().NoError(err)
	suite.MockBank.SetBalance(aliceAddr, sdk.NewCoin("uatom", sdkmath.NewInt(1000)))

	msg := &types.MsgSendWithAliasRouting{
		FromAddress: suite.Alice,
		ToAddress:   suite.Bob,
		Amount: sdk.Coins{
			sdk.NewCoin("uatom", sdkmath.NewInt(1000)),
		},
	}

	resp, err := suite.MsgServer.SendWithAliasRouting(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err)
	suite.Require().NotNil(resp)

	// Verify balance was transferred
	bobAddr, err := sdk.AccAddressFromBech32(suite.Bob)
	suite.Require().NoError(err)
	bobBalance := suite.MockBank.GetBalance(suite.Ctx, bobAddr, "uatom")
	suite.Require().Equal(sdkmath.NewInt(1000), bobBalance.Amount)
}

func (suite *CompleteTransferFlowTestSuite) TestCompleteTransferFlow_MultipleCoins() {
	router := testutil.GenerateMockRouter("tokenization:")
	err := suite.Keeper.RegisterRouter("tokenization:", router)
	suite.Require().NoError(err)

	// Recreate msg server after registering router (msgServer embeds Keeper by value)
	suite.MsgServer = keeper.NewMsgServerImpl(suite.Keeper)

	msg := &types.MsgSendWithAliasRouting{
		FromAddress: suite.Alice,
		ToAddress:   suite.Bob,
		Amount: sdk.Coins{
			sdk.NewCoin("tokenization:123:456", sdkmath.NewInt(1000)),
			sdk.NewCoin("uatom", sdkmath.NewInt(500)),
		},
	}

	// Mixed denoms - alias denom succeeds, but bank denom fails due to insufficient funds
	_, err = suite.MsgServer.SendWithAliasRouting(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "insufficient funds")
}
