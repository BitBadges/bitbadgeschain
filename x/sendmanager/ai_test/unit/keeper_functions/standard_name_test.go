package keeper_functions

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/ai_test/testutil"
	sendmanagerkeeper "github.com/bitbadges/bitbadgeschain/x/sendmanager/keeper"
)

type StandardNameTestSuite struct {
	testutil.AITestSuite
}

func TestStandardNameTestSuite(t *testing.T) {
	suite.Run(t, new(StandardNameTestSuite))
}

func (suite *StandardNameTestSuite) TestStandardName_AliasDenom() {
	router := testutil.GenerateMockRouter(sendmanagerkeeper.AliasDenomPrefix)
	err := suite.Keeper.RegisterRouter(sendmanagerkeeper.AliasDenomPrefix, router)
	suite.Require().NoError(err)

	// Alias denom should return "x/tokenization"
	name := suite.Keeper.StandardName(suite.Ctx, "badgeslp:123:456")
	suite.Require().Equal("x/tokenization", name)
}

func (suite *StandardNameTestSuite) TestStandardName_BankDenom() {
	// Bank denom should return "x/bank"
	name := suite.Keeper.StandardName(suite.Ctx, "uatom")
	suite.Require().Equal("x/bank", name)
}

func (suite *StandardNameTestSuite) TestStandardName_EmptyDenom() {
	// Empty denom defaults to "x/bank"
	name := suite.Keeper.StandardName(suite.Ctx, "")
	suite.Require().Equal("x/bank", name)
}

func (suite *StandardNameTestSuite) TestStandardName_UnregisteredPrefix() {
	// Unregistered prefix should return "x/bank"
	name := suite.Keeper.StandardName(suite.Ctx, "tokens:123")
	suite.Require().Equal("x/bank", name)
}
