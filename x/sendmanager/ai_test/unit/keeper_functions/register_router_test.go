package keeper_functions

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/ai_test/testutil"
)

type RegisterRouterTestSuite struct {
	testutil.AITestSuite
}

func TestRegisterRouterTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterRouterTestSuite))
}

func (suite *RegisterRouterTestSuite) TestRegisterRouter_Success() {
	router := testutil.GenerateMockRouter("badges:")
	
	err := suite.Keeper.RegisterRouter("badges:", router)
	suite.Require().NoError(err)
	
	prefixes := suite.Keeper.GetRegisteredPrefixes()
	suite.Require().Contains(prefixes, "badges:")
}

func (suite *RegisterRouterTestSuite) TestRegisterRouter_EmptyPrefix() {
	router := testutil.GenerateMockRouter("")
	
	err := suite.Keeper.RegisterRouter("", router)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "prefix cannot be empty")
}

func (suite *RegisterRouterTestSuite) TestRegisterRouter_DuplicatePrefix() {
	router1 := testutil.GenerateMockRouter("badges:")
	router2 := testutil.GenerateMockRouter("badges:")
	
	err := suite.Keeper.RegisterRouter("badges:", router1)
	suite.Require().NoError(err)
	
	err = suite.Keeper.RegisterRouter("badges:", router2)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "already registered")
}

func (suite *RegisterRouterTestSuite) TestRegisterRouter_OverlappingPrefixes_Prevented() {
	router1 := testutil.GenerateMockRouter("a:")
	router2 := testutil.GenerateMockRouter("a:b:")
	
	// Register longer prefix first
	err := suite.Keeper.RegisterRouter("a:b:", router2)
	suite.Require().NoError(err)
	
	// Try to register shorter prefix that overlaps - should fail
	// "a:b:" starts with "a:", so they overlap
	err = suite.Keeper.RegisterRouter("a:", router1)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "overlaps")
}

func (suite *RegisterRouterTestSuite) TestRegisterRouter_OverlappingPrefixes_ReverseOrder() {
	router1 := testutil.GenerateMockRouter("a:")
	router2 := testutil.GenerateMockRouter("a:b:")
	
	// Register shorter prefix first
	err := suite.Keeper.RegisterRouter("a:", router1)
	suite.Require().NoError(err)
	
	// Try to register longer prefix that overlaps - should fail
	// "a:b:" starts with "a:", so they overlap
	err = suite.Keeper.RegisterRouter("a:b:", router2)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "overlaps")
}

func (suite *RegisterRouterTestSuite) TestRegisterRouter_LongestPrefixMatching() {
	router1 := testutil.GenerateMockRouter("a:")
	router2 := testutil.GenerateMockRouter("a:b:")
	
	// Register shorter prefix first
	err := suite.Keeper.RegisterRouter("a:", router1)
	suite.Require().NoError(err)
	
	// This should fail due to overlap
	// "a:b:" starts with "a:", so they overlap
	err = suite.Keeper.RegisterRouter("a:b:", router2)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "overlaps")
}

func (suite *RegisterRouterTestSuite) TestRegisterRouter_MultipleNonOverlappingPrefixes() {
	router1 := testutil.GenerateMockRouter("badges:")
	router2 := testutil.GenerateMockRouter("tokens:")
	
	err := suite.Keeper.RegisterRouter("badges:", router1)
	suite.Require().NoError(err)
	
	err = suite.Keeper.RegisterRouter("tokens:", router2)
	suite.Require().NoError(err)
	
	prefixes := suite.Keeper.GetRegisteredPrefixes()
	suite.Require().Len(prefixes, 2)
	suite.Require().Contains(prefixes, "badges:")
	suite.Require().Contains(prefixes, "tokens:")
}

