package sendmanager_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"

	gamm "github.com/bitbadges/bitbadgeschain/x/gamm/precompile"
	sendmanager "github.com/bitbadges/bitbadgeschain/x/sendmanager/precompile"
	tokenization "github.com/bitbadges/bitbadgeschain/x/tokenization/precompile"
	"github.com/bitbadges/bitbadgeschain/x/sendmanager/precompile/test/helpers"
)

// CallerVerificationTestSuite tests that contract.Caller() is always auto-set
// for all three manual precompiles: sendmanager, tokenization, and gamm
type CallerVerificationTestSuite struct {
	suite.Suite
	SendManagerPrecompile  *sendmanager.Precompile
	TokenizationPrecompile *tokenization.Precompile
	GammPrecompile         *gamm.Precompile
	TestSuite              *helpers.TestSuite
}

func TestCallerVerificationTestSuite(t *testing.T) {
	suite.Run(t, new(CallerVerificationTestSuite))
}

func (suite *CallerVerificationTestSuite) SetupTest() {
	suite.TestSuite = helpers.NewTestSuite(suite.T())
	suite.SendManagerPrecompile = suite.TestSuite.Precompile

	// Create tokenization precompile
	tokenizationKeeper := suite.TestSuite.App.TokenizationKeeper
	suite.TokenizationPrecompile = tokenization.NewPrecompile(tokenizationKeeper)

	// Create gamm precompile
	gammKeeper := suite.TestSuite.App.GammKeeper
	suite.GammPrecompile = gamm.NewPrecompile(gammKeeper)
}

func (suite *CallerVerificationTestSuite) TestSendManager_CallerIsAutoSet() {
	// Test that sendmanager precompile always uses contract.Caller() for from_address
	caller := suite.TestSuite.AliceEVM
	precompileAddr := common.HexToAddress(sendmanager.SendManagerPrecompileAddress)
	valueUint256, _ := uint256.FromBig(big.NewInt(0))
	contract := vm.NewContract(caller, precompileAddr, valueUint256, 1000000, nil)

	// Verify GetCallerAddress returns the caller
	callerAddr, err := suite.SendManagerPrecompile.GetCallerAddress(contract)
	suite.NoError(err)
	suite.NotEmpty(callerAddr)

	// Convert to Cosmos address and verify it matches
	cosmosAddr, err := sdk.AccAddressFromBech32(callerAddr)
	suite.NoError(err)
	suite.Equal(caller.Bytes(), cosmosAddr.Bytes(), "SendManager: Caller address should match contract.Caller()")
}

func (suite *CallerVerificationTestSuite) TestTokenization_CallerIsAutoSet() {
	// Test that tokenization precompile always uses contract.Caller() for creator
	caller := suite.TestSuite.AliceEVM
	precompileAddr := common.HexToAddress(tokenization.TokenizationPrecompileAddress)
	valueUint256, _ := uint256.FromBig(big.NewInt(0))
	contract := vm.NewContract(caller, precompileAddr, valueUint256, 1000000, nil)

	// Verify GetCallerAddress returns the caller
	callerAddr, err := suite.TokenizationPrecompile.GetCallerAddress(contract)
	suite.NoError(err)
	suite.NotEmpty(callerAddr)

	// Convert to Cosmos address and verify it matches
	cosmosAddr, err := sdk.AccAddressFromBech32(callerAddr)
	suite.NoError(err)
	suite.Equal(caller.Bytes(), cosmosAddr.Bytes(), "Tokenization: Caller address should match contract.Caller()")
}

func (suite *CallerVerificationTestSuite) TestGamm_CallerIsAutoSet() {
	// Test that gamm precompile always uses contract.Caller() for sender
	caller := suite.TestSuite.AliceEVM
	precompileAddr := common.HexToAddress(gamm.GammPrecompileAddress)
	valueUint256, _ := uint256.FromBig(big.NewInt(0))
	contract := vm.NewContract(caller, precompileAddr, valueUint256, 1000000, nil)

	// Verify GetCallerAddress returns the caller
	callerAddr, err := suite.GammPrecompile.GetCallerAddress(contract)
	suite.NoError(err)
	suite.NotEmpty(callerAddr)

	// Convert to Cosmos address and verify it matches
	cosmosAddr, err := sdk.AccAddressFromBech32(callerAddr)
	suite.NoError(err)
	suite.Equal(caller.Bytes(), cosmosAddr.Bytes(), "Gamm: Caller address should match contract.Caller()")
}

func (suite *CallerVerificationTestSuite) TestAllPrecompiles_CallerCannotBeSpoofed() {
	// Test that all three precompiles reject zero address callers
	zeroCaller := common.Address{}
	
	// Test sendmanager
	precompileAddr1 := common.HexToAddress(sendmanager.SendManagerPrecompileAddress)
	valueUint256, _ := uint256.FromBig(big.NewInt(0))
	contract1 := vm.NewContract(zeroCaller, precompileAddr1, valueUint256, 1000000, nil)
	_, err1 := suite.SendManagerPrecompile.GetCallerAddress(contract1)
	suite.Error(err1, "SendManager should reject zero address caller")
	suite.Contains(err1.Error(), "zero")

	// Test tokenization
	precompileAddr2 := common.HexToAddress(tokenization.TokenizationPrecompileAddress)
	contract2 := vm.NewContract(zeroCaller, precompileAddr2, valueUint256, 1000000, nil)
	_, err2 := suite.TokenizationPrecompile.GetCallerAddress(contract2)
	suite.Error(err2, "Tokenization should reject zero address caller")
	suite.Contains(err2.Error(), "zero")

	// Test gamm
	precompileAddr3 := common.HexToAddress(gamm.GammPrecompileAddress)
	contract3 := vm.NewContract(zeroCaller, precompileAddr3, valueUint256, 1000000, nil)
	_, err3 := suite.GammPrecompile.GetCallerAddress(contract3)
	suite.Error(err3, "Gamm should reject zero address caller")
	suite.Contains(err3.Error(), "zero")
}

func (suite *CallerVerificationTestSuite) TestAllPrecompiles_CallerConsistency() {
	// Test that all three precompiles return consistent caller addresses
	// for the same contract.Caller() value
	caller := suite.TestSuite.AliceEVM
	valueUint256, _ := uint256.FromBig(big.NewInt(0))

	// Get caller addresses from all three precompiles
	precompileAddr1 := common.HexToAddress(sendmanager.SendManagerPrecompileAddress)
	contract1 := vm.NewContract(caller, precompileAddr1, valueUint256, 1000000, nil)
	callerAddr1, err1 := suite.SendManagerPrecompile.GetCallerAddress(contract1)
	suite.NoError(err1)

	precompileAddr2 := common.HexToAddress(tokenization.TokenizationPrecompileAddress)
	contract2 := vm.NewContract(caller, precompileAddr2, valueUint256, 1000000, nil)
	callerAddr2, err2 := suite.TokenizationPrecompile.GetCallerAddress(contract2)
	suite.NoError(err2)

	precompileAddr3 := common.HexToAddress(gamm.GammPrecompileAddress)
	contract3 := vm.NewContract(caller, precompileAddr3, valueUint256, 1000000, nil)
	callerAddr3, err3 := suite.GammPrecompile.GetCallerAddress(contract3)
	suite.NoError(err3)

	// All should return the same Cosmos address for the same EVM caller
	suite.Equal(callerAddr1, callerAddr2, "SendManager and Tokenization should return same caller address")
	suite.Equal(callerAddr2, callerAddr3, "Tokenization and Gamm should return same caller address")
	suite.Equal(callerAddr1, callerAddr3, "SendManager and Gamm should return same caller address")

	// Verify the address matches the expected Cosmos address
	expectedCosmosAddr := sdk.AccAddress(caller.Bytes()).String()
	suite.Equal(expectedCosmosAddr, callerAddr1, "Caller address should match expected Cosmos address")
	suite.Equal(expectedCosmosAddr, callerAddr2, "Caller address should match expected Cosmos address")
	suite.Equal(expectedCosmosAddr, callerAddr3, "Caller address should match expected Cosmos address")
}

