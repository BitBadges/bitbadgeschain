package sendmanager_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	sendmanager "github.com/bitbadges/bitbadgeschain/x/sendmanager/precompile"
	"github.com/bitbadges/bitbadgeschain/x/sendmanager/precompile/test/helpers"
)

type GasTestSuite struct {
	suite.Suite
	Precompile *sendmanager.Precompile
	TestSuite  *helpers.TestSuite
}

func TestGasTestSuite(t *testing.T) {
	suite.Run(t, new(GasTestSuite))
}

func (suite *GasTestSuite) SetupTest() {
	suite.TestSuite = helpers.NewTestSuite(suite.T())
	suite.Precompile = suite.TestSuite.Precompile
}

func (suite *GasTestSuite) TestRequiredGas_SendMethod() {
	method := suite.Precompile.ABI.Methods["send"]
	require.NotNil(suite.T(), method)

	jsonMsg, err := helpers.BuildSendJSON(
		suite.TestSuite.Alice.String(),
		suite.TestSuite.Bob.String(),
		"1000",
		"stake",
	)
	suite.NoError(err)

	input, err := helpers.PackMethodCall(&method, jsonMsg)
	suite.NoError(err)

	// RequiredGas adds a 150k buffer to base gas for Cosmos SDK operations (bank transfers)
	const txBuffer = 150_000
	gas := suite.Precompile.RequiredGas(input)
	suite.Equal(uint64(sendmanager.GasSendBase+txBuffer), gas, "Gas should equal GasSendBase + buffer for send method")
}

func (suite *GasTestSuite) TestRequiredGas_InvalidInput() {
	// Input too short
	shortInput := []byte{0x12, 0x34}
	gas := suite.Precompile.RequiredGas(shortInput)
	suite.Equal(uint64(0), gas, "Gas should be 0 for invalid input")
}

func (suite *GasTestSuite) TestRequiredGas_UnknownMethod() {
	// Create input with unknown method ID
	unknownInput := []byte{0xFF, 0xFF, 0xFF, 0xFF}
	gas := suite.Precompile.RequiredGas(unknownInput)
	suite.Equal(uint64(0), gas, "Gas should be 0 for unknown method")
}

func (suite *GasTestSuite) TestGetBaseGas_SendMethod() {
	// Test getBaseGas indirectly through RequiredGas
	method := suite.Precompile.ABI.Methods["send"]
	require.NotNil(suite.T(), method)

	jsonMsg, err := helpers.BuildSendJSON(
		suite.TestSuite.Alice.String(),
		suite.TestSuite.Bob.String(),
		"1000",
		"stake",
	)
	suite.NoError(err)

	input, err := helpers.PackMethodCall(&method, jsonMsg)
	suite.NoError(err)

	// RequiredGas adds a 150k buffer to base gas for Cosmos SDK operations
	const txBuffer = 150_000
	gas := suite.Precompile.RequiredGas(input)
	suite.Equal(uint64(sendmanager.GasSendBase+txBuffer), gas)
}

func (suite *GasTestSuite) TestGetBaseGas_UnknownMethod() {
	// Test with unknown method ID
	unknownInput := []byte{0xFF, 0xFF, 0xFF, 0xFF}
	gas := suite.Precompile.RequiredGas(unknownInput)
	suite.Equal(uint64(0), gas)
}

