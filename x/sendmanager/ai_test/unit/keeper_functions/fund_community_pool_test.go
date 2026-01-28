package keeper_functions

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/ai_test/testutil"
)

type FundCommunityPoolTestSuite struct {
	testutil.AITestSuite
}

func TestFundCommunityPoolTestSuite(t *testing.T) {
	suite.Run(t, new(FundCommunityPoolTestSuite))
}

func (suite *FundCommunityPoolTestSuite) TestFundCommunityPoolWithAliasRouting_AliasDenom() {
	router := testutil.GenerateMockRouter("badges:")
	err := suite.Keeper.RegisterRouter("badges:", router)
	suite.Require().NoError(err)

	aliceAddr, err := sdk.AccAddressFromBech32(suite.Alice)
	suite.Require().NoError(err)

	coins := sdk.Coins{
		sdk.NewCoin("badges:123:456", sdkmath.NewInt(1000)),
	}
	err = suite.Keeper.FundCommunityPoolWithAliasRouting(suite.Ctx, aliceAddr, coins)
	suite.Require().NoError(err)
}

func (suite *FundCommunityPoolTestSuite) TestFundCommunityPoolWithAliasRouting_BankDenom() {
	aliceAddr, err := sdk.AccAddressFromBech32(suite.Alice)
	suite.Require().NoError(err)

	coins := sdk.Coins{
		sdk.NewCoin("uatom", sdkmath.NewInt(1000)),
	}
	err = suite.Keeper.FundCommunityPoolWithAliasRouting(suite.Ctx, aliceAddr, coins)
	// Mock distribution keeper doesn't check balances, so this should succeed
	// The distribution keeper's FundCommunityPool is a no-op in the mock
	suite.Require().NoError(err)
}

func (suite *FundCommunityPoolTestSuite) TestFundCommunityPoolWithAliasRouting_EmptyDenom() {
	aliceAddr, err := sdk.AccAddressFromBech32(suite.Alice)
	suite.Require().NoError(err)

	coins := sdk.Coins{
		sdk.Coin{Denom: "", Amount: sdkmath.NewInt(1000)},
	}
	err = suite.Keeper.FundCommunityPoolWithAliasRouting(suite.Ctx, aliceAddr, coins)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "cannot be empty")
}

func (suite *FundCommunityPoolTestSuite) TestFundCommunityPoolWithAliasRouting_MixedDenoms() {
	router := testutil.GenerateMockRouter("badges:")
	err := suite.Keeper.RegisterRouter("badges:", router)
	suite.Require().NoError(err)

	aliceAddr, err := sdk.AccAddressFromBech32(suite.Alice)
	suite.Require().NoError(err)

	coins := sdk.Coins{
		sdk.NewCoin("badges:123:456", sdkmath.NewInt(1000)),
		sdk.NewCoin("uatom", sdkmath.NewInt(500)),
	}
	err = suite.Keeper.FundCommunityPoolWithAliasRouting(suite.Ctx, aliceAddr, coins)
	suite.Require().NoError(err)
}

