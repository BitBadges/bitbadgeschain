package keeper_functions

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/ai_test/testutil"
	sendmanagerkeeper "github.com/bitbadges/bitbadgeschain/x/sendmanager/keeper"
)

type SendFromModuleToAccountTestSuite struct {
	testutil.AITestSuite
}

func TestSendFromModuleToAccountTestSuite(t *testing.T) {
	suite.Run(t, new(SendFromModuleToAccountTestSuite))
}

func (suite *SendFromModuleToAccountTestSuite) TestSendCoinsFromModuleToAccountWithAliasRouting_AliasDenom() {
	router := testutil.GenerateMockRouter(sendmanagerkeeper.AliasDenomPrefix)
	err := suite.Keeper.RegisterRouter(sendmanagerkeeper.AliasDenomPrefix, router)
	suite.Require().NoError(err)

	bobAddr, err := sdk.AccAddressFromBech32(suite.Bob)
	suite.Require().NoError(err)

	coins := sdk.Coins{
		sdk.NewCoin("badgeslp:123:456", sdkmath.NewInt(1000)),
	}
	err = suite.Keeper.SendCoinsFromModuleToAccountWithAliasRouting(suite.Ctx, "mymodule", bobAddr, coins)
	suite.Require().NoError(err)
}

func (suite *SendFromModuleToAccountTestSuite) TestSendCoinsFromModuleToAccountWithAliasRouting_BankDenom() {
	bobAddr, err := sdk.AccAddressFromBech32(suite.Bob)
	suite.Require().NoError(err)

	coins := sdk.Coins{
		sdk.NewCoin("uatom", sdkmath.NewInt(1000)),
	}
	err = suite.Keeper.SendCoinsFromModuleToAccountWithAliasRouting(suite.Ctx, "mymodule", bobAddr, coins)
	// Mock bank keeper checks balances for module-to-account transfers
	// This will fail because the module doesn't have balance
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "insufficient funds")
}

func (suite *SendFromModuleToAccountTestSuite) TestSendCoinsFromModuleToAccountWithAliasRouting_EmptyDenom() {
	bobAddr, err := sdk.AccAddressFromBech32(suite.Bob)
	suite.Require().NoError(err)

	coins := sdk.Coins{
		sdk.Coin{Denom: "", Amount: sdkmath.NewInt(1000)},
	}
	err = suite.Keeper.SendCoinsFromModuleToAccountWithAliasRouting(suite.Ctx, "mymodule", bobAddr, coins)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "cannot be empty")
}
