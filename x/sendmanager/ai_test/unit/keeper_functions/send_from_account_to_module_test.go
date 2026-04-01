package keeper_functions

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/ai_test/testutil"
	sendmanagerkeeper "github.com/bitbadges/bitbadgeschain/x/sendmanager/keeper"
)

type SendFromAccountToModuleTestSuite struct {
	testutil.AITestSuite
}

func TestSendFromAccountToModuleTestSuite(t *testing.T) {
	suite.Run(t, new(SendFromAccountToModuleTestSuite))
}

func (suite *SendFromAccountToModuleTestSuite) TestSendCoinsFromAccountToModuleWithAliasRouting_AliasDenom() {
	router := testutil.GenerateMockRouter(sendmanagerkeeper.AliasDenomPrefix)
	err := suite.Keeper.RegisterRouter(sendmanagerkeeper.AliasDenomPrefix, router)
	suite.Require().NoError(err)

	aliceAddr, err := sdk.AccAddressFromBech32(suite.Alice)
	suite.Require().NoError(err)

	coins := sdk.Coins{
		sdk.NewCoin("badgeslp:123:456", sdkmath.NewInt(1000)),
	}
	err = suite.Keeper.SendCoinsFromAccountToModuleWithAliasRouting(suite.Ctx, aliceAddr, "mymodule", coins)
	suite.Require().NoError(err)
}

func (suite *SendFromAccountToModuleTestSuite) TestSendCoinsFromAccountToModuleWithAliasRouting_BankDenom() {
	aliceAddr, err := sdk.AccAddressFromBech32(suite.Alice)
	suite.Require().NoError(err)

	coins := sdk.Coins{
		sdk.NewCoin("uatom", sdkmath.NewInt(1000)),
	}
	err = suite.Keeper.SendCoinsFromAccountToModuleWithAliasRouting(suite.Ctx, aliceAddr, "mymodule", coins)
	// Mock bank keeper checks balances
	// This will fail because Alice doesn't have balance
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "insufficient funds")
}

func (suite *SendFromAccountToModuleTestSuite) TestSendCoinsFromAccountToModuleWithAliasRouting_EmptyDenom() {
	aliceAddr, err := sdk.AccAddressFromBech32(suite.Alice)
	suite.Require().NoError(err)

	coins := sdk.Coins{
		sdk.Coin{Denom: "", Amount: sdkmath.NewInt(1000)},
	}
	err = suite.Keeper.SendCoinsFromAccountToModuleWithAliasRouting(suite.Ctx, aliceAddr, "mymodule", coins)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "cannot be empty")
}
