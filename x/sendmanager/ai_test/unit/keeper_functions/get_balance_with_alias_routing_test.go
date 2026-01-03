package keeper_functions

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/ai_test/testutil"
)

type GetBalanceWithAliasRoutingTestSuite struct {
	testutil.AITestSuite
}

func TestGetBalanceWithAliasRoutingTestSuite(t *testing.T) {
	suite.Run(t, new(GetBalanceWithAliasRoutingTestSuite))
}

func (suite *GetBalanceWithAliasRoutingTestSuite) TestGetBalanceWithAliasRouting_AliasDenom() {
	router := testutil.GenerateMockRouter("badges:")
	err := suite.Keeper.RegisterRouter("badges:", router)
	suite.Require().NoError(err)

	aliceAddr, err := sdk.AccAddressFromBech32(suite.Alice)
	suite.Require().NoError(err)

	balance, err := suite.Keeper.GetBalanceWithAliasRouting(suite.Ctx, aliceAddr, "badges:123:456")
	suite.Require().NoError(err)
	suite.Require().Equal("badges:123:456", balance.Denom)
}

func (suite *GetBalanceWithAliasRoutingTestSuite) TestGetBalanceWithAliasRouting_BankDenom() {
	aliceAddr, err := sdk.AccAddressFromBech32(suite.Alice)
	suite.Require().NoError(err)

	balance, err := suite.Keeper.GetBalanceWithAliasRouting(suite.Ctx, aliceAddr, "uatom")
	suite.Require().NoError(err)
	suite.Require().Equal("uatom", balance.Denom)
}

func (suite *GetBalanceWithAliasRoutingTestSuite) TestGetBalanceWithAliasRouting_EmptyDenom() {
	aliceAddr, err := sdk.AccAddressFromBech32(suite.Alice)
	suite.Require().NoError(err)

	_, err = suite.Keeper.GetBalanceWithAliasRouting(suite.Ctx, aliceAddr, "")
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "cannot be empty")
}

