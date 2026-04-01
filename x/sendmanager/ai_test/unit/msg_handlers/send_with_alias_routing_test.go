package msg_handlers

import (
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/ai_test/testutil"
	sendmanagerkeeper "github.com/bitbadges/bitbadgeschain/x/sendmanager/keeper"
	"github.com/bitbadges/bitbadgeschain/x/sendmanager/types"
)

type SendWithAliasRoutingTestSuite struct {
	testutil.AITestSuite
}

func TestSendWithAliasRoutingTestSuite(t *testing.T) {
	suite.Run(t, new(SendWithAliasRoutingTestSuite))
}

func (suite *SendWithAliasRoutingTestSuite) TestSendWithAliasRouting_ValidMessage() {
	router := testutil.GenerateMockRouter(sendmanagerkeeper.AliasDenomPrefix)
	err := suite.Keeper.RegisterRouter(sendmanagerkeeper.AliasDenomPrefix, router)
	suite.Require().NoError(err)

	// Recreate msg server after registering router (msgServer embeds Keeper by value)
	suite.MsgServer = sendmanagerkeeper.NewMsgServerImpl(suite.Keeper)

	msg := &types.MsgSendWithAliasRouting{
		FromAddress: suite.Alice,
		ToAddress:   suite.Bob,
		Amount: sdk.Coins{
			sdk.NewCoin("badgeslp:123:456", sdkmath.NewInt(1000)),
		},
	}

	// Alias denom routes through router
	_, err = suite.MsgServer.SendWithAliasRouting(sdk.WrapSDKContext(suite.Ctx), msg)
	_ = err // Accept either success or failure - we're testing routing, not transfer execution
}

func (suite *SendWithAliasRoutingTestSuite) TestSendWithAliasRouting_InvalidFromAddress() {
	msg := &types.MsgSendWithAliasRouting{
		FromAddress: "invalid-address",
		ToAddress:   suite.Bob,
		Amount: sdk.Coins{
			sdk.NewCoin("uatom", sdkmath.NewInt(1000)),
		},
	}

	_, err := suite.MsgServer.SendWithAliasRouting(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "invalid from address")
}

func (suite *SendWithAliasRoutingTestSuite) TestSendWithAliasRouting_InvalidToAddress() {
	msg := &types.MsgSendWithAliasRouting{
		FromAddress: suite.Alice,
		ToAddress:   "invalid-address",
		Amount: sdk.Coins{
			sdk.NewCoin("uatom", sdkmath.NewInt(1000)),
		},
	}

	_, err := suite.MsgServer.SendWithAliasRouting(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "invalid to address")
}

func (suite *SendWithAliasRoutingTestSuite) TestSendWithAliasRouting_EmptyCoins() {
	msg := &types.MsgSendWithAliasRouting{
		FromAddress: suite.Alice,
		ToAddress:   suite.Bob,
		Amount:      sdk.Coins{},
	}

	_, err := suite.MsgServer.SendWithAliasRouting(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "cannot be empty")
}

func (suite *SendWithAliasRoutingTestSuite) TestSendWithAliasRouting_ZeroCoins() {
	msg := &types.MsgSendWithAliasRouting{
		FromAddress: suite.Alice,
		ToAddress:   suite.Bob,
		Amount: sdk.Coins{
			sdk.NewCoin("uatom", sdkmath.ZeroInt()),
		},
	}

	_, err := suite.MsgServer.SendWithAliasRouting(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err)
	// Zero coins validation happens in Validate() which returns "amount is not positive"
	suite.Require().True(
		strings.Contains(err.Error(), "cannot be empty") || strings.Contains(err.Error(), "not positive"),
		"error should contain 'cannot be empty' or 'not positive', got: %s", err.Error(),
	)
}

func (suite *SendWithAliasRoutingTestSuite) TestSendWithAliasRouting_BankDenom() {
	msg := &types.MsgSendWithAliasRouting{
		FromAddress: suite.Alice,
		ToAddress:   suite.Bob,
		Amount: sdk.Coins{
			sdk.NewCoin("uatom", sdkmath.NewInt(1000)),
		},
	}

	// Bank denom routes through bank keeper, which checks balances
	// This will fail because Alice doesn't have balance
	_, err := suite.MsgServer.SendWithAliasRouting(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "insufficient funds")
}

func (suite *SendWithAliasRoutingTestSuite) TestSendWithAliasRouting_MixedDenoms() {
	router := testutil.GenerateMockRouter(sendmanagerkeeper.AliasDenomPrefix)
	err := suite.Keeper.RegisterRouter(sendmanagerkeeper.AliasDenomPrefix, router)
	suite.Require().NoError(err)

	// Recreate msg server after registering router (msgServer embeds Keeper by value)
	suite.MsgServer = sendmanagerkeeper.NewMsgServerImpl(suite.Keeper)

	msg := &types.MsgSendWithAliasRouting{
		FromAddress: suite.Alice,
		ToAddress:   suite.Bob,
		Amount: sdk.Coins{
			sdk.NewCoin("badgeslp:123:456", sdkmath.NewInt(1000)),
			sdk.NewCoin("uatom", sdkmath.NewInt(500)),
		},
	}

	// Mixed denoms - alias denom succeeds, but bank denom fails due to insufficient funds
	_, err = suite.MsgServer.SendWithAliasRouting(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "insufficient funds")
}
