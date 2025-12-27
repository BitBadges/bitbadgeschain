package keeper_functions

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/sendmanager/types"
)

type ParamsTestSuite struct {
	testutil.AITestSuite
}

func TestParamsTestSuite(t *testing.T) {
	suite.Run(t, new(ParamsTestSuite))
}

func (suite *ParamsTestSuite) TestGetParams() {
	params := suite.Keeper.GetParams(suite.Ctx)
	suite.Require().NotNil(params)
}

func (suite *ParamsTestSuite) TestSetParams() {
	params := types.DefaultParams()
	err := suite.Keeper.SetParams(suite.Ctx, params)
	suite.Require().NoError(err)

	// Verify params were set
	retrievedParams := suite.Keeper.GetParams(suite.Ctx)
	suite.Require().Equal(params, retrievedParams)
}
