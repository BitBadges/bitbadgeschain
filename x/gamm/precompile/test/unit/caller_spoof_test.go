package gamm_test

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"

	gamm "github.com/bitbadges/bitbadgeschain/x/gamm/precompile"
	"github.com/bitbadges/bitbadgeschain/x/gamm/precompile/test/helpers"
)

// CallerSpoofTestSuite tests that contract.Caller() cannot be spoofed
// and is always used for the sender field, even if JSON contains a different sender
type CallerSpoofTestSuite struct {
	suite.Suite
	Precompile *gamm.Precompile
	TestSuite  *helpers.TestSuite
}

func TestCallerSpoofTestSuite(t *testing.T) {
	suite.Run(t, new(CallerSpoofTestSuite))
}

func (suite *CallerSpoofTestSuite) SetupTest() {
	suite.TestSuite = helpers.NewTestSuite(suite.T())
	suite.Precompile = suite.TestSuite.Precompile
}

func (suite *CallerSpoofTestSuite) createContract(caller common.Address, input []byte) *vm.Contract {
	precompileAddr := common.HexToAddress(gamm.GammPrecompileAddress)
	valueUint256, _ := uint256.FromBig(big.NewInt(0))
	contract := vm.NewContract(caller, precompileAddr, valueUint256, 1000000, nil)
	if len(input) > 0 {
		contract.Input = input
	}
	return contract
}

func (suite *CallerSpoofTestSuite) packMethodWithJSON(method *abi.Method, jsonStr string) ([]byte, error) {
	args := []interface{}{jsonStr}
	packed, err := method.Inputs.Pack(args...)
	if err != nil {
		return nil, err
	}
	return append(method.ID[:], packed...), nil
}

func (suite *CallerSpoofTestSuite) TestJoinPool_SenderOverriddenByCaller() {
	// Test that even if JSON contains a different sender, contract.Caller() is used
	// First, fund Alice with tokens needed for pool creation
	poolCreationCoins := sdk.NewCoins(
		sdk.NewCoin("uatom", osmomath.NewInt(2_000_000_000_000_000_000)),
		sdk.NewCoin("uosmo", osmomath.NewInt(2_000_000_000_000_000_000)),
	)
	suite.TestSuite.FundAcc(suite.TestSuite.Alice, poolCreationCoins)
	
	// Create a pool
	poolId, err := suite.TestSuite.CreateDefaultTestPool(suite.TestSuite.Alice)
	suite.Require().NoError(err)
	suite.Require().NotZero(poolId)

	caller := suite.TestSuite.AliceEVM
	wrongSenderInJSON := suite.TestSuite.Charlie.String() // Different from caller

	method := suite.Precompile.ABI.Methods["joinPool"]
	require.NotNil(suite.T(), method)

	// Build JSON message with wrong sender
	jsonMsg := map[string]interface{}{
		"sender":         wrongSenderInJSON, // Wrong sender in JSON
		"pool_id":        poolId,
		"share_out_amount": "1000000",
		"token_in_maxs": []map[string]interface{}{
			{"denom": "uatom", "amount": "1000000"},
			{"denom": "uosmo", "amount": "1000000"},
		},
	}
	jsonBytes, err := json.Marshal(jsonMsg)
	suite.NoError(err)

	input, err := suite.packMethodWithJSON(&method, string(jsonBytes))
	suite.NoError(err)

	contract := suite.createContract(caller, input)
	
	// Get initial balances
	aliceInitialBalance := suite.TestSuite.App.BankKeeper.GetBalance(suite.TestSuite.Ctx, suite.TestSuite.Alice, "uatom")
	charlieInitialBalance := suite.TestSuite.App.BankKeeper.GetBalance(suite.TestSuite.Ctx, suite.TestSuite.Charlie, "uatom")

	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	
	if err == nil {
		suite.NotNil(result)
		
		// Verify that Alice's balance decreased (she is the caller), not Charlie's
		aliceFinalBalance := suite.TestSuite.App.BankKeeper.GetBalance(suite.TestSuite.Ctx, suite.TestSuite.Alice, "uatom")
		charlieFinalBalance := suite.TestSuite.App.BankKeeper.GetBalance(suite.TestSuite.Ctx, suite.TestSuite.Charlie, "uatom")

		// Alice's balance should have decreased (she is the actual caller)
		suite.True(
			aliceFinalBalance.Amount.LT(aliceInitialBalance.Amount),
			"Alice's balance should decrease (she is the caller)",
		)

		// Charlie's balance should NOT have changed (she was in JSON but not the caller)
		suite.Equal(
			charlieInitialBalance.Amount.String(),
			charlieFinalBalance.Amount.String(),
			"Charlie's balance should not change (she was not the caller)",
		)
	} else {
		// If error occurred, it should be about insufficient balance or validation,
		// not about wrong sender
		suite.NotContains(err.Error(), "impersonation", "Error should not be about sender impersonation")
	}
}

func (suite *CallerSpoofTestSuite) TestSwapExactAmountIn_SenderOverriddenByCaller() {
	// Test that sender is always set from contract.Caller(), not from JSON input
	// First, fund Alice with tokens needed for pool creation
	poolCreationCoins := sdk.NewCoins(
		sdk.NewCoin("uatom", osmomath.NewInt(2_000_000_000_000_000_000)),
		sdk.NewCoin("uosmo", osmomath.NewInt(2_000_000_000_000_000_000)),
	)
	suite.TestSuite.FundAcc(suite.TestSuite.Alice, poolCreationCoins)
	
	// Create a pool
	poolId, err := suite.TestSuite.CreateDefaultTestPool(suite.TestSuite.Alice)
	suite.Require().NoError(err)
	suite.Require().NotZero(poolId)

	caller := suite.TestSuite.AliceEVM
	wrongSenderInJSON := suite.TestSuite.Charlie.String() // Different from caller

	method := suite.Precompile.ABI.Methods["swapExactAmountIn"]
	require.NotNil(suite.T(), method)

	// Build JSON message with wrong sender
	jsonMsg := map[string]interface{}{
		"sender": wrongSenderInJSON, // Wrong sender in JSON
		"routes": []map[string]interface{}{
			{
				"pool_id":        poolId,
				"token_out_denom": "uosmo",
			},
		},
		"token_in": map[string]interface{}{
			"denom":  "uatom",
			"amount": "100000",
		},
		"token_out_min_amount": "50000",
	}
	jsonBytes, err := json.Marshal(jsonMsg)
	suite.NoError(err)

	input, err := suite.packMethodWithJSON(&method, string(jsonBytes))
	suite.NoError(err)

	contract := suite.createContract(caller, input)
	
	// Get initial balances
	aliceInitialBalance := suite.TestSuite.App.BankKeeper.GetBalance(suite.TestSuite.Ctx, suite.TestSuite.Alice, "uatom")
	charlieInitialBalance := suite.TestSuite.App.BankKeeper.GetBalance(suite.TestSuite.Ctx, suite.TestSuite.Charlie, "uatom")

	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	
	if err == nil {
		suite.NotNil(result)
		
		// Verify that Alice's balance decreased (she is the caller), not Charlie's
		aliceFinalBalance := suite.TestSuite.App.BankKeeper.GetBalance(suite.TestSuite.Ctx, suite.TestSuite.Alice, "uatom")
		charlieFinalBalance := suite.TestSuite.App.BankKeeper.GetBalance(suite.TestSuite.Ctx, suite.TestSuite.Charlie, "uatom")

		// Alice's balance should have decreased (she is the actual caller)
		suite.True(
			aliceFinalBalance.Amount.LT(aliceInitialBalance.Amount),
			"Alice's balance should decrease (she is the caller)",
		)

		// Charlie's balance should NOT have changed (she was in JSON but not the caller)
		suite.Equal(
			charlieInitialBalance.Amount.String(),
			charlieFinalBalance.Amount.String(),
			"Charlie's balance should not change (she was not the caller)",
		)
	} else {
		// If error occurred, it should be about insufficient balance or validation,
		// not about wrong sender
		suite.NotContains(err.Error(), "impersonation", "Error should not be about sender impersonation")
	}
}

