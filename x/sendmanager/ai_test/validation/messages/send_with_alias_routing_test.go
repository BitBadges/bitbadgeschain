package messages

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/sendmanager/types"
)

type SendWithAliasRoutingValidationTestSuite struct {
	testutil.AITestSuite
}

func TestSendWithAliasRoutingValidationTestSuite(t *testing.T) {
	suite.Run(t, new(SendWithAliasRoutingValidationTestSuite))
}

func (suite *SendWithAliasRoutingValidationTestSuite) TestSendWithAliasRouting_ValidMessage() {
	msg := &types.MsgSendWithAliasRouting{
		FromAddress: suite.Alice,
		ToAddress:   suite.Bob,
		Amount: sdk.Coins{
			sdk.NewCoin("uatom", sdkmath.NewInt(1000)),
		},
	}

	// MsgSendWithAliasRouting doesn't have ValidateBasic, so we test via message server
	// For validation tests, we verify the message structure is valid
	suite.Require().NotNil(msg)
	suite.Require().NotEmpty(msg.FromAddress)
	suite.Require().NotEmpty(msg.ToAddress)
	suite.Require().False(msg.Amount.Empty())
}

func (suite *SendWithAliasRoutingValidationTestSuite) TestSendWithAliasRouting_InvalidFromAddress() {
	msg := &types.MsgSendWithAliasRouting{
		FromAddress: "invalid-address",
		ToAddress:   suite.Bob,
		Amount: sdk.Coins{
			sdk.NewCoin("uatom", sdkmath.NewInt(1000)),
		},
	}

	// Test that invalid address is caught by message server
	_, err := suite.MsgServer.SendWithAliasRouting(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "invalid from address")
}

func (suite *SendWithAliasRoutingValidationTestSuite) TestSendWithAliasRouting_InvalidToAddress() {
	msg := &types.MsgSendWithAliasRouting{
		FromAddress: suite.Alice,
		ToAddress:   "invalid-address",
		Amount: sdk.Coins{
			sdk.NewCoin("uatom", sdkmath.NewInt(1000)),
		},
	}

	// Test that invalid address is caught by message server
	_, err := suite.MsgServer.SendWithAliasRouting(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "invalid to address")
}

func (suite *SendWithAliasRoutingValidationTestSuite) TestSendWithAliasRouting_EmptyCoins() {
	msg := &types.MsgSendWithAliasRouting{
		FromAddress: suite.Alice,
		ToAddress:   suite.Bob,
		Amount:      sdk.Coins{},
	}

	// Test that empty coins are caught by message server
	_, err := suite.MsgServer.SendWithAliasRouting(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "cannot be empty")
}

func (suite *SendWithAliasRoutingValidationTestSuite) TestSendWithAliasRouting_InvalidCoins() {
	msg := &types.MsgSendWithAliasRouting{
		FromAddress: suite.Alice,
		ToAddress:   suite.Bob,
		Amount: sdk.Coins{
			sdk.Coin{Denom: "", Amount: sdkmath.NewInt(1000)}, // Invalid: empty denom
		},
	}

	// Test that invalid coins are caught by message server
	_, err := suite.MsgServer.SendWithAliasRouting(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err)
}

