package keeper_functions

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/ai_test/testutil"
	sendmanagerkeeper "github.com/bitbadges/bitbadgeschain/x/sendmanager/keeper"
)

type IsICS20CompatibleTestSuite struct {
	testutil.AITestSuite
}

func TestIsICS20CompatibleTestSuite(t *testing.T) {
	suite.Run(t, new(IsICS20CompatibleTestSuite))
}

func (suite *IsICS20CompatibleTestSuite) TestIsICS20Compatible_AliasDenom() {
	router := testutil.GenerateMockRouter(sendmanagerkeeper.AliasDenomPrefix)
	err := suite.Keeper.RegisterRouter(sendmanagerkeeper.AliasDenomPrefix, router)
	suite.Require().NoError(err)

	// Alias denom should not be ICS20 compatible
	compatible := suite.Keeper.IsICS20Compatible(suite.Ctx, "badgeslp:123:456")
	suite.Require().False(compatible)
}

func (suite *IsICS20CompatibleTestSuite) TestIsICS20Compatible_BankDenom() {
	// Bank denom should be ICS20 compatible
	compatible := suite.Keeper.IsICS20Compatible(suite.Ctx, "uatom")
	suite.Require().True(compatible)
}

func (suite *IsICS20CompatibleTestSuite) TestIsICS20Compatible_EmptyDenom() {
	// Empty denom is considered ICS20 compatible
	compatible := suite.Keeper.IsICS20Compatible(suite.Ctx, "")
	suite.Require().True(compatible)
}

func (suite *IsICS20CompatibleTestSuite) TestIsICS20Compatible_UnregisteredPrefix() {
	// Unregistered prefix should be ICS20 compatible
	compatible := suite.Keeper.IsICS20Compatible(suite.Ctx, "tokens:123")
	suite.Require().True(compatible)
}
