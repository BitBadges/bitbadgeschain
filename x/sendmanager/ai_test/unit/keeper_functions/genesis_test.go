package keeper_functions

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/sendmanager/types"
)

type GenesisTestSuite struct {
	testutil.AITestSuite
}

func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}

func (suite *GenesisTestSuite) TestInitGenesis() {
	genState := types.GenesisState{
		Params: types.DefaultParams(),
	}
	err := suite.Keeper.InitGenesis(suite.Ctx, genState)
	suite.Require().NoError(err)
}

func (suite *GenesisTestSuite) TestExportGenesis() {
	genState := types.GenesisState{
		Params: types.DefaultParams(),
	}
	err := suite.Keeper.InitGenesis(suite.Ctx, genState)
	suite.Require().NoError(err)

	exported, err := suite.Keeper.ExportGenesis(suite.Ctx)
	suite.Require().NoError(err)
	suite.Require().NotNil(exported)
	suite.Require().Equal(genState.Params, exported.Params)
}

func (suite *GenesisTestSuite) TestInitGenesis_ExportGenesis_RoundTrip() {
	genState := types.GenesisState{
		Params: types.DefaultParams(),
	}
	err := suite.Keeper.InitGenesis(suite.Ctx, genState)
	suite.Require().NoError(err)

	exported, err := suite.Keeper.ExportGenesis(suite.Ctx)
	suite.Require().NoError(err)
	suite.Require().NotNil(exported)
	suite.Require().Equal(genState.Params, exported.Params)
}
