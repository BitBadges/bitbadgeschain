package keeper

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/ai_test/testutil"
)

type PrefixMatchingValidationTestSuite struct {
	testutil.AITestSuite
}

func TestPrefixMatchingValidationTestSuite(t *testing.T) {
	suite.Run(t, new(PrefixMatchingValidationTestSuite))
}

func (suite *PrefixMatchingValidationTestSuite) TestGetBalanceWithAliasRouting_EmptyDenom() {
	router := testutil.GenerateMockRouter("badges:")
	err := suite.Keeper.RegisterRouter("badges:", router)
	suite.Require().NoError(err)

	aliceAddr, err := sdk.AccAddressFromBech32(suite.Alice)
	suite.Require().NoError(err)

	_, err = suite.Keeper.GetBalanceWithAliasRouting(suite.Ctx, aliceAddr, "")
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "cannot be empty")
}

func (suite *PrefixMatchingValidationTestSuite) TestSendCoinWithAliasRouting_EmptyDenom() {
	router := testutil.GenerateMockRouter("badges:")
	err := suite.Keeper.RegisterRouter("badges:", router)
	suite.Require().NoError(err)

	aliceAddr, err := sdk.AccAddressFromBech32(suite.Alice)
	suite.Require().NoError(err)
	bobAddr, err := sdk.AccAddressFromBech32(suite.Bob)
	suite.Require().NoError(err)

	coin := sdk.Coin{Denom: "", Amount: sdkmath.NewInt(1000)}
	err = suite.Keeper.SendCoinWithAliasRouting(suite.Ctx, aliceAddr, bobAddr, &coin)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "cannot be empty")
}

func (suite *PrefixMatchingValidationTestSuite) TestIsICS20Compatible_EmptyDenom() {
	// Empty denom should be considered ICS20 compatible
	compatible := suite.Keeper.IsICS20Compatible(suite.Ctx, "")
	suite.Require().True(compatible)
}

func (suite *PrefixMatchingValidationTestSuite) TestStandardName_EmptyDenom() {
	// Empty denom should default to "x/bank"
	name := suite.Keeper.StandardName(suite.Ctx, "")
	suite.Require().Equal("x/bank", name)
}

