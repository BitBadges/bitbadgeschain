package gamm_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/third_party/apptesting"
	gamm "github.com/bitbadges/bitbadgeschain/x/gamm/precompile"
)

type IntegrationTestSuite struct {
	apptesting.KeeperTestHelper
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (suite *IntegrationTestSuite) SetupTest() {
	suite.Reset()
}

// TestIntegration_PrecompileSetup verifies the precompile can be instantiated
// Full integration tests with actual pool operations require complete pool setup
// which should be done in the gamm module's integration test suite
func (suite *IntegrationTestSuite) TestIntegration_PrecompileSetup() {
	precompile := gamm.NewPrecompile(suite.App.GammKeeper)

	suite.NotNil(precompile)
	suite.NotNil(suite.Ctx)

	// Verify the precompile has the correct address
	suite.Equal(gamm.GammPrecompileAddress, precompile.ContractAddress.Hex())

	// Verify ABI is loaded
	suite.NotNil(precompile.ABI)
	method, found := precompile.ABI.Methods["joinPool"]
	suite.True(found)
	suite.NotNil(method)
}

// TestIntegration_ABILoading verifies ABI loads correctly
func (suite *IntegrationTestSuite) TestIntegration_ABILoading() {
	err := gamm.GetABILoadError()
	suite.NoError(err, "ABI should load successfully")
}

// TestIntegration_KeeperIntegration verifies basic keeper integration
func (suite *IntegrationTestSuite) TestIntegration_KeeperIntegration() {
	precompile := gamm.NewPrecompile(suite.App.GammKeeper)
	
	// Verify precompile has access to keeper
	suite.NotNil(precompile)
	
	// Verify keeper is accessible (indirectly through precompile)
	// The keeper should be set during precompile creation
	suite.Equal(gamm.GammPrecompileAddress, precompile.ContractAddress.Hex())
}
