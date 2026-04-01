package keeper_functions

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/ai_test/testutil"
	sendmanagerkeeper "github.com/bitbadges/bitbadgeschain/x/sendmanager/keeper"
)

type SpendFromCommunityPoolTestSuite struct {
	testutil.AITestSuite
}

func TestSpendFromCommunityPoolTestSuite(t *testing.T) {
	suite.Run(t, new(SpendFromCommunityPoolTestSuite))
}

func (suite *SpendFromCommunityPoolTestSuite) TestSpendFromCommunityPoolWithAliasRouting_AliasDenom() {
	router := testutil.GenerateMockRouter(sendmanagerkeeper.AliasDenomPrefix)
	err := suite.Keeper.RegisterRouter(sendmanagerkeeper.AliasDenomPrefix, router)
	suite.Require().NoError(err)

	bobAddr, err := sdk.AccAddressFromBech32(suite.Bob)
	suite.Require().NoError(err)

	coins := sdk.Coins{
		sdk.NewCoin("badgeslp:123:456", sdkmath.NewInt(1000)),
	}
	err = suite.Keeper.SpendFromCommunityPoolWithAliasRouting(suite.Ctx, bobAddr, coins)
	suite.Require().NoError(err)
}

func (suite *SpendFromCommunityPoolTestSuite) TestSpendFromCommunityPoolWithAliasRouting_BankDenom() {
	bobAddr, err := sdk.AccAddressFromBech32(suite.Bob)
	suite.Require().NoError(err)

	coins := sdk.Coins{
		sdk.NewCoin("uatom", sdkmath.NewInt(1000)),
	}
	err = suite.Keeper.SpendFromCommunityPoolWithAliasRouting(suite.Ctx, bobAddr, coins)
	// Mock bank keeper checks balances for module-to-account transfers
	// This will fail because the distribution module doesn't have balance
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "insufficient funds")
}

func (suite *SpendFromCommunityPoolTestSuite) TestSpendFromCommunityPoolWithAliasRouting_EmptyDenom() {
	bobAddr, err := sdk.AccAddressFromBech32(suite.Bob)
	suite.Require().NoError(err)

	coins := sdk.Coins{
		sdk.Coin{Denom: "", Amount: sdkmath.NewInt(1000)},
	}
	err = suite.Keeper.SpendFromCommunityPoolWithAliasRouting(suite.Ctx, bobAddr, coins)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "cannot be empty")
}
