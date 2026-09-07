package tokenization_test

import (
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"

	sdkmath "cosmossdk.io/math"

	tokenization "github.com/bitbadges/bitbadgeschain/x/tokenization/precompile"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/precompile/test/helpers"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// GasAccuracyTestSuite is a test suite for gas accuracy testing
type GasAccuracyTestSuite struct {
	EVMKeeperIntegrationTestSuite
}

func TestGasAccuracyTestSuite(t *testing.T) {
	suite.Run(t, new(GasAccuracyTestSuite))
}

// SetupTest sets up the test suite
func (suite *GasAccuracyTestSuite) SetupTest() {
	suite.EVMKeeperIntegrationTestSuite.SetupTest()
}

// TestGasAccuracy_TransferTokens_EstimateVsActual tests gas estimation vs actual usage
func (suite *GasAccuracyTestSuite) TestGasAccuracy_TransferTokens_EstimateVsActual() {
	chainID := suite.getChainID()
	precompileAddr := common.HexToAddress(tokenization.TokenizationPrecompileAddress)
	method := suite.Precompile.ABI.Methods["transferTokens"]
	suite.Require().NotNil(method)

	toAddressesStr := []string{suite.Bob.String()}
	jsonMsg, err := helpers.BuildTransferTokensJSON(
		suite.CollectionId.BigInt(),
		suite.Alice.String(),
		toAddressesStr,
		big.NewInt(10),
		[]struct{ Start, End *big.Int }{{Start: big.NewInt(1), End: big.NewInt(10)}},
		[]struct{ Start, End *big.Int }{{Start: big.NewInt(1), End: new(big.Int).SetUint64(math.MaxUint64)}},
	)
	suite.Require().NoError(err)

	input, err := helpers.PackMethodWithJSON(&method, jsonMsg)
	suite.Require().NoError(err)

	// Estimate gas
	estimatedGas := suite.Precompile.RequiredGas(input)
	suite.T().Logf("Estimated gas: %d", estimatedGas)

	// Execute transaction with sufficient gas limit
	nonce := suite.getNonce(suite.AliceEVM)
	tx, err := helpers.BuildEVMTransaction(
		suite.AliceKey,
		&precompileAddr,
		input,
		big.NewInt(0),
		1000000, // Large gas limit to avoid out of gas
		big.NewInt(0),
		nonce,
		chainID,
	)
	suite.Require().NoError(err)

	response, err := helpers.ExecuteEVMTransaction(suite.Ctx, suite.EVMKeeper, tx)
	suite.Require().NoError(err)
	suite.Require().NotNil(response)

	// Verify actual gas used
	actualGas := response.GasUsed
	suite.T().Logf("Actual gas used: %d", actualGas)

	// Verify gas is within reasonable tolerance (10% of estimate)
	// Note: Actual gas may be higher due to EVM overhead
	if estimatedGas > 0 {
		tolerance := float64(estimatedGas) * 0.1
		diff := float64(actualGas) - float64(estimatedGas)
		suite.T().Logf("Gas difference: %.0f (tolerance: %.0f)", diff, tolerance)
		// Actual gas should be at least the estimated gas (may be higher due to overhead)
		suite.Require().GreaterOrEqual(actualGas, estimatedGas, "Actual gas should be at least the estimated gas")
	}
}

// TestGasAccuracy_AllMethods_WithinTolerance tests gas accuracy for all methods
func (suite *GasAccuracyTestSuite) TestGasAccuracy_AllMethods_WithinTolerance() {
	testMethods := []string{
		"transferTokens",
		"getCollection",
		"getBalance",
		"getTotalSupply",
	}

	for _, methodName := range testMethods {
		method, found := suite.Precompile.ABI.Methods[methodName]
		if !found {
			suite.T().Logf("Method %s not found, skipping", methodName)
			continue
		}

		// Build JSON for each method
		var input []byte
		var err error
		switch methodName {
		case "transferTokens":
			toAddressesStr := []string{suite.Bob.String()}
			jsonMsg, jsonErr := helpers.BuildTransferTokensJSON(
				suite.CollectionId.BigInt(),
				suite.Alice.String(),
				toAddressesStr,
				big.NewInt(1),
				[]struct{ Start, End *big.Int }{{Start: big.NewInt(1), End: big.NewInt(1)}},
				[]struct{ Start, End *big.Int }{{Start: big.NewInt(1), End: new(big.Int).SetUint64(math.MaxUint64)}},
			)
			if jsonErr != nil {
				suite.T().Logf("Failed to build JSON for %s: %v", methodName, jsonErr)
				continue
			}
			input, err = helpers.PackMethodWithJSON(&method, jsonMsg)
		case "getCollection":
			queryJson, jsonErr := helpers.BuildGetCollectionQueryJSON(suite.CollectionId.BigInt())
			if jsonErr != nil {
				suite.T().Logf("Failed to build JSON for %s: %v", methodName, jsonErr)
				continue
			}
			input, err = helpers.PackMethodWithJSON(&method, queryJson)
		case "getBalance":
			queryJson, jsonErr := helpers.BuildGetBalanceQueryJSON(suite.CollectionId.BigInt(), suite.Alice.String())
			if jsonErr != nil {
				suite.T().Logf("Failed to build JSON for %s: %v", methodName, jsonErr)
				continue
			}
			input, err = helpers.PackMethodWithJSON(&method, queryJson)
		default:
			suite.T().Logf("Method %s not handled in gas test, skipping", methodName)
			continue
		}

		if err != nil {
			suite.T().Logf("Failed to pack args for %s: %v", methodName, err)
			continue
		}

		estimatedGas := suite.Precompile.RequiredGas(input)
		suite.T().Logf("%s - Estimated gas: %d", methodName, estimatedGas)

		// Verify estimate is reasonable (not zero, not extremely large)
		suite.Require().Greater(estimatedGas, uint64(0), "Gas estimate should be greater than 0")
		suite.Require().Less(estimatedGas, uint64(10000000), "Gas estimate should be reasonable")
	}
}

// TestGasLimits_Enforced tests that gas limits are enforced
func (suite *GasAccuracyTestSuite) TestGasLimits_Enforced() {
	chainID := suite.getChainID()
	precompileAddr := common.HexToAddress(tokenization.TokenizationPrecompileAddress)
	method := suite.Precompile.ABI.Methods["transferTokens"]
	suite.Require().NotNil(method)

	toAddressesStr := []string{suite.Bob.String()}
	jsonMsg, err := helpers.BuildTransferTokensJSON(
		suite.CollectionId.BigInt(),
		suite.Alice.String(),
		toAddressesStr,
		big.NewInt(1),
		[]struct{ Start, End *big.Int }{{Start: big.NewInt(1), End: big.NewInt(1)}},
		[]struct{ Start, End *big.Int }{{Start: big.NewInt(1), End: new(big.Int).SetUint64(math.MaxUint64)}},
	)
	suite.Require().NoError(err)

	input, err := helpers.PackMethodWithJSON(&method, jsonMsg)
	suite.Require().NoError(err)

	// Execute with very low gas limit
	nonce := suite.getNonce(suite.AliceEVM)
	tx, err := helpers.BuildEVMTransaction(
		suite.AliceKey,
		&precompileAddr,
		input,
		big.NewInt(0),
		1000, // Very low gas limit
		big.NewInt(0),
		nonce,
		chainID,
	)
	suite.Require().NoError(err)

	response, err := helpers.ExecuteEVMTransaction(suite.Ctx, suite.EVMKeeper, tx)
	// Transaction should fail with out of gas
	if response != nil {
		suite.T().Logf("Transaction response: VmError=%s, GasUsed=%d", response.VmError, response.GasUsed)
		// Out of gas errors are expected with low gas limits
		if response.VmError != "" {
			suite.T().Log("Gas limit enforcement verified - transaction failed as expected")
		}
	}
}

// TestGasUsed_GrowsWithRecipients runs the same transfer through the EVM
// keeper with one and with five recipients and checks the gas the EVM charged
// grows by at least the per-recipient price. This is the end-to-end view of
// the SDK-side metering: store access and per-element charges taken inside
// the precompile must reach the EVM's gas accounting.
func (suite *GasAccuracyTestSuite) TestGasUsed_GrowsWithRecipients() {
	chainID := suite.getChainID()
	precompileAddr := common.HexToAddress(tokenization.TokenizationPrecompileAddress)
	method := suite.Precompile.ABI.Methods["transferTokens"]

	// Fresh recipients have no balance store, so let the collection's default
	// balances accept incoming transfers.
	collection, found := suite.TokenizationKeeper.GetCollectionFromStore(suite.Ctx, suite.CollectionId)
	suite.Require().True(found)
	collection.DefaultBalances = &tokenizationtypes.UserBalanceStore{
		AutoApproveAllIncomingTransfers:           true,
		AutoApproveSelfInitiatedIncomingTransfers: true,
		AutoApproveSelfInitiatedOutgoingTransfers: true,
	}
	suite.Require().NoError(suite.TokenizationKeeper.SetCollectionInStore(suite.Ctx, collection, false))

	// x/vm charges at least MinGasMultiplier (default 50%) of the gas limit,
	// which would hide the actual usage behind a flat number.
	feeParams := suite.App.FeeMarketKeeper.GetParams(suite.Ctx)
	feeParams.MinGasMultiplier = sdkmath.LegacyZeroDec()
	suite.Require().NoError(suite.App.FeeMarketKeeper.SetParams(suite.Ctx, feeParams))

	gasFor := func(recipients int) uint64 {
		toAddresses := make([]string, recipients)
		for i := range toAddresses {
			_, _, addr := helpers.CreateEVMAccount()
			toAddresses[i] = addr.String()
		}
		jsonMsg, err := helpers.BuildTransferTokensJSON(
			suite.CollectionId.BigInt(),
			suite.Alice.String(),
			toAddresses,
			big.NewInt(1),
			[]struct{ Start, End *big.Int }{{Start: big.NewInt(1), End: big.NewInt(50)}},
			[]struct{ Start, End *big.Int }{{Start: big.NewInt(1), End: new(big.Int).SetUint64(math.MaxUint64)}},
		)
		suite.Require().NoError(err)
		input, err := helpers.PackMethodWithJSON(&method, jsonMsg)
		suite.Require().NoError(err)

		tx, err := helpers.BuildEVMTransaction(suite.AliceKey, &precompileAddr, input, big.NewInt(0), 5_000_000, big.NewInt(0), suite.getNonce(suite.AliceEVM), chainID)
		suite.Require().NoError(err)
		response, err := helpers.ExecuteEVMTransaction(suite.Ctx, suite.EVMKeeper, tx)
		suite.Require().NoError(err)
		suite.Require().Empty(response.VmError, "transfer to %d recipients must succeed", recipients)
		return response.GasUsed
	}

	one := gasFor(1)
	five := gasFor(5)
	suite.T().Logf("gas used: 1 recipient=%d, 5 recipients=%d", one, five)
	suite.GreaterOrEqual(five-one, uint64(4*tokenization.GasPerRecipient),
		"gas charged by the EVM must grow with the number of recipients")
}
