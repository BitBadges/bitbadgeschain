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
	router1 := testutil.GenerateMockRouter("tokenization:")
	router2 := testutil.GenerateMockRouter("tokenization:")

	err := suite.Keeper.RegisterRouter("tokenization:", router1)
	suite.Require().NoError(err)

	err = suite.Keeper.RegisterRouter("tokenization:", router2)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "already registered")
}

func (suite *RouterRegistrationValidationTestSuite) TestRegisterRouter_OverlappingPrefixes() {
	router1 := testutil.GenerateMockRouter("tokenization:")
	router2 := testutil.GenerateMockRouter("tokenization:lp:")

	err := suite.Keeper.RegisterRouter("tokenization:", router1)
	suite.Require().NoError(err)

	// This should fail because "tokenization:lp:" starts with "tokenization:"
	err = suite.Keeper.RegisterRouter("tokenization:lp:", router2)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "overlaps")
}

func (suite *RouterRegistrationValidationTestSuite) TestRegisterRouter_ReverseOverlappingPrefixes() {
	router1 := testutil.GenerateMockRouter("tokenization:lp:")
	router2 := testutil.GenerateMockRouter("tokenization:")

	err := suite.Keeper.RegisterRouter("tokenization:lp:", router1)
	suite.Require().NoError(err)

	// This should fail because "tokenization:" is a prefix of "tokenization:lp:"
	err = suite.Keeper.RegisterRouter("tokenization:", router2)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "overlaps")
}

