package keeper_functions

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/sendmanager/keeper"
)

type RegisterRouterTestSuite struct {
	testutil.AITestSuite
}

func TestRegisterRouterTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterRouterTestSuite))
}

func (suite *RegisterRouterTestSuite) TestRegisterRouter_Success() {
	router := testutil.GenerateMockRouter(keeper.AliasDenomPrefix)

	err := suite.Keeper.RegisterRouter(keeper.AliasDenomPrefix, router)
	suite.Require().NoError(err)

	prefixes := suite.Keeper.GetRegisteredPrefixes()
	suite.Require().Contains(prefixes, keeper.AliasDenomPrefix)
}

func (suite *RegisterRouterTestSuite) TestRegisterRouter_EmptyPrefix_Rejected() {
	router := testutil.GenerateMockRouter("")

	err := suite.Keeper.RegisterRouter("", router)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "only prefix")
}

func (suite *RegisterRouterTestSuite) TestRegisterRouter_NonBadgeslpPrefix_Rejected() {
	router := testutil.GenerateMockRouter("badges:")

	err := suite.Keeper.RegisterRouter("badges:", router)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "only prefix")
}

func (suite *RegisterRouterTestSuite) TestRegisterRouter_AnotherInvalidPrefix_Rejected() {
	router := testutil.GenerateMockRouter("tokens:")

	err := suite.Keeper.RegisterRouter("tokens:", router)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "only prefix")
}

func (suite *RegisterRouterTestSuite) TestRegisterRouter_OverlappingPrefix_Rejected() {
	router := testutil.GenerateMockRouter("a:b:")

	err := suite.Keeper.RegisterRouter("a:b:", router)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "only prefix")
}

func (suite *RegisterRouterTestSuite) TestRegisterRouter_DuplicatePrefix_OverwritesSuccessfully() {
	router1 := testutil.GenerateMockRouter(keeper.AliasDenomPrefix)
	router2 := testutil.GenerateMockRouter(keeper.AliasDenomPrefix)

	err := suite.Keeper.RegisterRouter(keeper.AliasDenomPrefix, router1)
	suite.Require().NoError(err)

	// Registering the same prefix again should succeed (SetAliasRouter just overwrites)
	err = suite.Keeper.RegisterRouter(keeper.AliasDenomPrefix, router2)
	suite.Require().NoError(err)

	prefixes := suite.Keeper.GetRegisteredPrefixes()
	suite.Require().Len(prefixes, 1)
	suite.Require().Contains(prefixes, keeper.AliasDenomPrefix)
}

func (suite *RegisterRouterTestSuite) TestSetAliasRouter_Success() {
	router := testutil.GenerateMockRouter(keeper.AliasDenomPrefix)

	suite.Keeper.SetAliasRouter(router)

	prefixes := suite.Keeper.GetRegisteredPrefixes()
	suite.Require().Len(prefixes, 1)
	suite.Require().Contains(prefixes, keeper.AliasDenomPrefix)
}

func (suite *RegisterRouterTestSuite) TestGetRegisteredPrefixes_NoRouter() {
	// Without setting a router, should return empty
	prefixes := suite.Keeper.GetRegisteredPrefixes()
	suite.Require().Empty(prefixes)
}
