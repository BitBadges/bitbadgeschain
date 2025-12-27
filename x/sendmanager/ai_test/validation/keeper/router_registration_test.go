package keeper

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/ai_test/testutil"
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
	suite.Require().Contains(err.Error(), "cannot be empty")
}

func (suite *RouterRegistrationValidationTestSuite) TestRegisterRouter_DuplicatePrefix() {
	router1 := testutil.GenerateMockRouter("badges:")
	router2 := testutil.GenerateMockRouter("badges:")

	err := suite.Keeper.RegisterRouter("badges:", router1)
	suite.Require().NoError(err)

	err = suite.Keeper.RegisterRouter("badges:", router2)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "already registered")
}

func (suite *RouterRegistrationValidationTestSuite) TestRegisterRouter_OverlappingPrefixes() {
	router1 := testutil.GenerateMockRouter("badges:")
	router2 := testutil.GenerateMockRouter("badges:lp:")

	err := suite.Keeper.RegisterRouter("badges:", router1)
	suite.Require().NoError(err)

	// This should fail because "badges:lp:" starts with "badges:"
	err = suite.Keeper.RegisterRouter("badges:lp:", router2)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "overlaps")
}

func (suite *RouterRegistrationValidationTestSuite) TestRegisterRouter_ReverseOverlappingPrefixes() {
	router1 := testutil.GenerateMockRouter("badges:lp:")
	router2 := testutil.GenerateMockRouter("badges:")

	err := suite.Keeper.RegisterRouter("badges:lp:", router1)
	suite.Require().NoError(err)

	// This should fail because "badges:" is a prefix of "badges:lp:"
	err = suite.Keeper.RegisterRouter("badges:", router2)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "overlaps")
}

