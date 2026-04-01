package keeper

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/ai_test/testutil"
	sendmanagerkeeper "github.com/bitbadges/bitbadgeschain/x/sendmanager/keeper"
)

type RouterRegistrationValidationTestSuite struct {
	testutil.AITestSuite
}

func TestRouterRegistrationValidationTestSuite(t *testing.T) {
	suite.Run(t, new(RouterRegistrationValidationTestSuite))
}

func (suite *RouterRegistrationValidationTestSuite) TestRegisterRouter_EmptyPrefix() {
	router := testutil.GenerateMockRouter("")
	err := suite.Keeper.RegisterRouter("", router)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "only prefix")
}

func (suite *RouterRegistrationValidationTestSuite) TestRegisterRouter_NonBadgeslpPrefix_Rejected() {
	router := testutil.GenerateMockRouter("badges:")

	err := suite.Keeper.RegisterRouter("badges:", router)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "only prefix")
}

func (suite *RouterRegistrationValidationTestSuite) TestRegisterRouter_AnotherInvalidPrefix_Rejected() {
	router := testutil.GenerateMockRouter("a:")

	err := suite.Keeper.RegisterRouter("a:", router)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "only prefix")
}

func (suite *RouterRegistrationValidationTestSuite) TestRegisterRouter_YetAnotherInvalidPrefix_Rejected() {
	router := testutil.GenerateMockRouter("a:b:")

	err := suite.Keeper.RegisterRouter("a:b:", router)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "only prefix")
}

func (suite *RouterRegistrationValidationTestSuite) TestRegisterRouter_BadgeslpPrefix_Succeeds() {
	router := testutil.GenerateMockRouter(sendmanagerkeeper.AliasDenomPrefix)

	err := suite.Keeper.RegisterRouter(sendmanagerkeeper.AliasDenomPrefix, router)
	suite.Require().NoError(err)

	prefixes := suite.Keeper.GetRegisteredPrefixes()
	suite.Require().Len(prefixes, 1)
	suite.Require().Contains(prefixes, sendmanagerkeeper.AliasDenomPrefix)
}

func (suite *RouterRegistrationValidationTestSuite) TestRegisterRouter_BadgeslpDuplicate_Overwrites() {
	router1 := testutil.GenerateMockRouter(sendmanagerkeeper.AliasDenomPrefix)
	router2 := testutil.GenerateMockRouter(sendmanagerkeeper.AliasDenomPrefix)

	err := suite.Keeper.RegisterRouter(sendmanagerkeeper.AliasDenomPrefix, router1)
	suite.Require().NoError(err)

	// Re-registering badgeslp: should succeed (overwrites)
	err = suite.Keeper.RegisterRouter(sendmanagerkeeper.AliasDenomPrefix, router2)
	suite.Require().NoError(err)

	prefixes := suite.Keeper.GetRegisteredPrefixes()
	suite.Require().Len(prefixes, 1)
}
