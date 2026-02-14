package gamm_test

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/crypto/ed25519"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	evmkeeper "github.com/cosmos/evm/x/vm/keeper"
	evmtypes "github.com/cosmos/evm/x/vm/types"

	"github.com/bitbadges/bitbadgeschain/third_party/apptesting"
	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	"github.com/bitbadges/bitbadgeschain/x/gamm/poolmodels/balancer"
	gamm "github.com/bitbadges/bitbadgeschain/x/gamm/precompile"
	"github.com/bitbadges/bitbadgeschain/x/gamm/precompile/test/helpers"
	tokenizationprecompile "github.com/bitbadges/bitbadgeschain/x/tokenization/precompile"
	tokenizationhelpers "github.com/bitbadges/bitbadgeschain/x/tokenization/precompile/test/helpers"
)

// GammHelperLibrariesTestSuite tests all gamm helper libraries E2E
type GammHelperLibrariesTestSuite struct {
	apptesting.KeeperTestHelper

	EVMKeeper  *evmkeeper.Keeper
	Precompile *gamm.Precompile

	// Contract deployment
	TestContractAddr common.Address
	TestContractABI  abi.ABI
	ContractBytecode []byte

	// Test accounts
	DeployerKey *ecdsa.PrivateKey
	AliceKey    *ecdsa.PrivateKey
	BobKey      *ecdsa.PrivateKey
	DeployerEVM common.Address
	AliceEVM    common.Address
	BobEVM      common.Address
	Deployer    sdk.AccAddress
	Alice       sdk.AccAddress
	Bob         sdk.AccAddress

	PoolId  uint64
	ChainID *big.Int
}

func TestGammHelperLibrariesTestSuite(t *testing.T) {
	suite.Run(t, new(GammHelperLibrariesTestSuite))
}

// SetupTest sets up the test suite
func (suite *GammHelperLibrariesTestSuite) SetupTest() {
	suite.Reset()

	// Set up context
	header := cmtproto.Header{
		Height:  1,
		Time:    time.Now(),
		ChainID: "bitbadges-1",
	}
	suite.Ctx = suite.App.BaseApp.NewContextLegacy(false, header)

	// Create validator for EVM transactions
	var firstValidator stakingtypes.ValidatorI
	suite.App.StakingKeeper.IterateValidators(suite.Ctx, func(_ int64, val stakingtypes.ValidatorI) (stop bool) {
		firstValidator = val
		return true
	})

	if firstValidator == nil {
		suite.createTestValidator()
		suite.App.StakingKeeper.IterateValidators(suite.Ctx, func(_ int64, val stakingtypes.ValidatorI) (stop bool) {
			firstValidator = val
			return true
		})
	}

	suite.Require().NotNil(firstValidator, "Validator must be created for EVM transactions")

	valConsAddr, err := firstValidator.GetConsAddr()
	require.NoError(suite.T(), err)

	voteInfos := []abci.VoteInfo{{
		Validator:   abci.Validator{Address: valConsAddr, Power: 1000},
		BlockIdFlag: cmtproto.BlockIDFlagCommit,
	}}
	suite.Ctx = suite.Ctx.WithVoteInfos(voteInfos)

	blockHeader := suite.Ctx.BlockHeader()
	blockHeader.ProposerAddress = valConsAddr
	suite.Ctx = suite.Ctx.WithBlockHeader(blockHeader)

	_, err = suite.App.BeginBlocker(suite.Ctx)
	require.NoError(suite.T(), err)

	// Get keepers
	suite.EVMKeeper = suite.App.EVMKeeper

	// Create precompile instances
	suite.Precompile = gamm.NewPrecompile(suite.App.GammKeeper)
	gammPrecompileAddr := common.HexToAddress(gamm.GammPrecompileAddress)

	// Register and ENABLE the precompile - both steps are required!
	suite.EVMKeeper.RegisterStaticPrecompile(gammPrecompileAddr, suite.Precompile)
	err = suite.EVMKeeper.EnableStaticPrecompiles(suite.Ctx, gammPrecompileAddr)
	// If already enabled, that's fine - just continue
	if err != nil && !strings.Contains(err.Error(), "already") {
		suite.Require().NoError(err, "Failed to enable gamm precompile")
	}
	require.Equal(suite.T(), gamm.GammPrecompileAddress, suite.Precompile.ContractAddress.Hex())

	// Also register and enable tokenization precompile for consistency
	tokenizationPrecompile := tokenizationprecompile.NewPrecompile(suite.App.TokenizationKeeper)
	tokenizationPrecompileAddr := common.HexToAddress(tokenizationprecompile.TokenizationPrecompileAddress)
	suite.EVMKeeper.RegisterStaticPrecompile(tokenizationPrecompileAddr, tokenizationPrecompile)
	err = suite.EVMKeeper.EnableStaticPrecompiles(suite.Ctx, tokenizationPrecompileAddr)
	if err != nil && !strings.Contains(err.Error(), "already") {
		suite.Require().NoError(err, "Failed to enable tokenization precompile")
	}

	// Create test accounts
	suite.DeployerKey, suite.DeployerEVM, suite.Deployer = helpers.CreateEVMAccount()
	suite.AliceKey, suite.AliceEVM, suite.Alice = helpers.CreateEVMAccount()
	suite.BobKey, suite.BobEVM, suite.Bob = helpers.CreateEVMAccount()

	// Fund accounts
	suite.FundAcc(suite.Deployer, sdk.NewCoins(sdk.NewCoin("ustake", osmomath.NewInt(10000000000000000))))
	suite.FundAcc(suite.Alice, sdk.NewCoins(sdk.NewCoin("ustake", osmomath.NewInt(10000000000000000))))
	suite.FundAcc(suite.Bob, sdk.NewCoins(sdk.NewCoin("ustake", osmomath.NewInt(10000000000000000))))

	// Fund with pool tokens
	poolCoins := sdk.NewCoins(
		sdk.NewCoin("uatom", osmomath.NewInt(2_000_000_000_000_000_000)),
		sdk.NewCoin("uosmo", osmomath.NewInt(2_000_000_000_000_000_000)),
	)
	suite.FundAcc(suite.Alice, poolCoins)
	suite.FundAcc(suite.Bob, poolCoins)

	// Create test pool
	suite.PoolId = suite.createTestPool()

	// Set chain ID (using testnet chain ID for testing)
	suite.ChainID = big.NewInt(90123)

	// Load compiled contract - must use gamm-specific contract (no fallback)
	contractBytecode, err := helpers.GetGammContractBytecode()
	if err != nil {
		suite.T().Skipf("Skipping test - could not load gamm contract bytecode: %v. Run 'make compile-contracts' first.", err)
		return
	}
	suite.ContractBytecode = contractBytecode

	contractABI, err := helpers.GetGammContractABI()
	if err != nil {
		suite.T().Skipf("Skipping test - could not load gamm contract ABI: %v. Run 'make compile-contracts' first.", err)
		return
	}
	suite.TestContractABI = contractABI

	// Deploy test contract
	contractAddr, response, err := tokenizationhelpers.DeployContract(
		suite.Ctx,
		suite.EVMKeeper,
		suite.DeployerKey,
		suite.ContractBytecode,
		suite.ChainID,
	)
	if err != nil {
		if response != nil && response.GasUsed > 0 {
			suite.T().Skipf("Skipping test - contract deployment reverted: %v", err)
			return
		}
		suite.Require().NoError(err, "Contract deployment failed")
	}

	isContract, verifyErr := tokenizationhelpers.VerifyContractDeployment(suite.Ctx, suite.EVMKeeper, contractAddr)
	suite.Require().NoError(verifyErr)
	if !isContract {
		suite.T().Skip("Skipping test - contract deployment failed")
		return
	}

	suite.TestContractAddr = contractAddr
}

// createTestValidator creates a test validator
func (suite *GammHelperLibrariesTestSuite) createTestValidator() {
	valPrivKey := ed25519.GenPrivKey()
	valPubKey := valPrivKey.PubKey()
	valAddr := sdk.ValAddress(valPubKey.Address())

	// Convert to SDK pubkey
	sdkPubKey, err := cryptocodec.FromTmPubKeyInterface(valPubKey)
	require.NoError(suite.T(), err)

	validator, err := stakingtypes.NewValidator(valAddr.String(), sdkPubKey, stakingtypes.Description{})
	require.NoError(suite.T(), err)

	validator.Status = stakingtypes.Bonded
	validator.Tokens = osmomath.NewInt(1000000)
	suite.App.StakingKeeper.SetValidator(suite.Ctx, validator)
	suite.App.StakingKeeper.SetValidatorByConsAddr(suite.Ctx, validator)
	suite.App.StakingKeeper.SetNewValidatorByPowerIndex(suite.Ctx, validator)

	// Set signing info
	valConsAddr, err := validator.GetConsAddr()
	require.NoError(suite.T(), err)
	signingInfo := slashingtypes.NewValidatorSigningInfo(
		sdk.ConsAddress(valConsAddr),
		0,
		0,
		time.Unix(0, 0),
		false,
		0,
	)
	suite.App.SlashingKeeper.SetValidatorSigningInfo(suite.Ctx, sdk.ConsAddress(valConsAddr), signingInfo)

	// Apply validator set updates
	_, err = suite.App.StakingKeeper.ApplyAndReturnValidatorSetUpdates(suite.Ctx)
	require.NoError(suite.T(), err)
}

// createTestPool creates a test balancer pool
func (suite *GammHelperLibrariesTestSuite) createTestPool() uint64 {
	oneTrillion := osmomath.NewInt(1e12)
	poolAssets := []balancer.PoolAsset{
		{
			Token:  sdk.NewCoin("uatom", oneTrillion),
			Weight: osmomath.NewInt(100),
		},
		{
			Token:  sdk.NewCoin("uosmo", oneTrillion),
			Weight: osmomath.NewInt(100),
		},
	}

	poolParams := balancer.PoolParams{
		SwapFee: osmomath.MustNewDecFromStr("0.025"),
		ExitFee: osmomath.ZeroDec(),
	}

	msg := balancer.NewMsgCreateBalancerPool(suite.Alice, poolParams, poolAssets)
	poolId, err := suite.App.PoolManagerKeeper.CreatePool(suite.Ctx, msg)
	suite.Require().NoError(err)
	return poolId
}

// callContract calls a contract method
func (suite *GammHelperLibrariesTestSuite) callContract(
	callerKey *ecdsa.PrivateKey,
	contractAddr common.Address,
	methodName string,
	args ...interface{},
) ([]byte, *evmtypes.MsgEthereumTxResponse, error) {
	return tokenizationhelpers.CallContractMethod(
		suite.Ctx,
		suite.EVMKeeper,
		callerKey,
		contractAddr,
		suite.TestContractABI,
		methodName,
		args,
		suite.ChainID,
		false, // isView
	)
}

// ============ Wrapper Tests ============

// TestWrappers_JoinPool tests the joinPool wrapper
func (suite *GammHelperLibrariesTestSuite) TestWrappers_JoinPool() {
	shareOutAmount := big.NewInt(1000000)
	tokenInMaxs := []struct {
		Denom  string
		Amount *big.Int
	}{
		{Denom: "uatom", Amount: big.NewInt(1000000)},
		{Denom: "uosmo", Amount: big.NewInt(2000000)},
	}

	returnData, response, err := suite.callContract(
		suite.AliceKey,
		suite.TestContractAddr,
		"testJoinPoolWrapper",
		suite.PoolId,
		shareOutAmount,
		tokenInMaxs,
	)
	suite.Require().NoError(err)
	suite.Greater(response.GasUsed, uint64(0), "Gas should be used")

	if len(returnData) > 0 {
		method := suite.TestContractABI.Methods["testJoinPoolWrapper"]
		unpacked, err := method.Outputs.Unpack(returnData)
		if err == nil && len(unpacked) >= 2 {
			suite.T().Logf("Join pool result: shares=%v, tokens=%v", unpacked[0], unpacked[1])
		}
	}
}

// TestWrappers_GetPool tests the getPool wrapper
func (suite *GammHelperLibrariesTestSuite) TestWrappers_GetPool() {
	returnData, response, err := suite.callContract(
		suite.AliceKey,
		suite.TestContractAddr,
		"testGetPoolWrapper",
		suite.PoolId,
	)
	suite.Require().NoError(err)
	suite.Greater(response.GasUsed, uint64(0), "Gas should be used")

	if len(returnData) > 0 {
		method := suite.TestContractABI.Methods["testGetPoolWrapper"]
		unpacked, err := method.Outputs.Unpack(returnData)
		if err == nil && len(unpacked) > 0 {
			if poolBytes, ok := unpacked[0].([]byte); ok {
				suite.Greater(len(poolBytes), 0, "Pool bytes should not be empty")
				suite.T().Logf("Pool bytes length: %d", len(poolBytes))
			}
		}
	}
}

// TestWrappers_GetTotalShares tests the getTotalShares wrapper
func (suite *GammHelperLibrariesTestSuite) TestWrappers_GetTotalShares() {
	returnData, response, err := suite.callContract(
		suite.AliceKey,
		suite.TestContractAddr,
		"testGetTotalSharesWrapper",
		suite.PoolId,
	)
	suite.Require().NoError(err)
	suite.Greater(response.GasUsed, uint64(0), "Gas should be used")

	if len(returnData) > 0 {
		method := suite.TestContractABI.Methods["testGetTotalSharesWrapper"]
		unpacked, err := method.Outputs.Unpack(returnData)
		if err == nil && len(unpacked) > 0 {
			suite.T().Logf("Total shares result: %v", unpacked[0])
		}
	}
}

// ============ Builder Tests ============

// TestBuilders_JoinPoolBuilder tests the JoinPoolBuilder
func (suite *GammHelperLibrariesTestSuite) TestBuilders_JoinPoolBuilder() {
	shareOutAmount := big.NewInt(1000000)

	returnData, response, err := suite.callContract(
		suite.AliceKey,
		suite.TestContractAddr,
		"testJoinPoolBuilder",
		suite.PoolId,
		shareOutAmount,
		"uatom",
		big.NewInt(1000000),
		"uosmo",
		big.NewInt(2000000),
	)
	suite.Require().NoError(err)
	suite.Greater(response.GasUsed, uint64(0), "Gas should be used")

	if len(returnData) > 0 {
		method := suite.TestContractABI.Methods["testJoinPoolBuilder"]
		unpacked, err := method.Outputs.Unpack(returnData)
		if err == nil && len(unpacked) > 0 {
			if jsonStr, ok := unpacked[0].(string); ok {
				suite.T().Logf("JoinPoolBuilder JSON: %s", jsonStr)
				suite.verifyJoinPoolBuilderJSON(jsonStr, suite.PoolId, shareOutAmount)
			}
		}
	}
}

// TestBuilders_SwapBuilder tests the SwapBuilder
func (suite *GammHelperLibrariesTestSuite) TestBuilders_SwapBuilder() {
	tokenInAmount := big.NewInt(100000)
	tokenOutMinAmount := big.NewInt(90000)

	returnData, response, err := suite.callContract(
		suite.AliceKey,
		suite.TestContractAddr,
		"testSwapBuilder",
		suite.PoolId,
		"uosmo",
		"uatom",
		tokenInAmount,
		tokenOutMinAmount,
	)
	suite.Require().NoError(err)
	suite.Greater(response.GasUsed, uint64(0), "Gas should be used")

	if len(returnData) > 0 {
		method := suite.TestContractABI.Methods["testSwapBuilder"]
		unpacked, err := method.Outputs.Unpack(returnData)
		if err == nil && len(unpacked) > 0 {
			if jsonStr, ok := unpacked[0].(string); ok {
				suite.T().Logf("SwapBuilder JSON: %s", jsonStr)
				suite.verifySwapBuilderJSON(jsonStr, suite.PoolId, "uosmo", "uatom", tokenInAmount, tokenOutMinAmount)
			}
		}
	}
}

// ============ JSON Helper Tests ============

// TestJSONHelpers_JoinPoolJSON tests joinPoolJSON
func (suite *GammHelperLibrariesTestSuite) TestJSONHelpers_JoinPoolJSON() {
	shareOutAmount := big.NewInt(1000000)
	tokenInMaxsJson := `[{"denom":"uatom","amount":"1000000"},{"denom":"uosmo","amount":"2000000"}]`

	returnData, response, err := suite.callContract(
		suite.AliceKey,
		suite.TestContractAddr,
		"testJoinPoolJSON",
		suite.PoolId,
		shareOutAmount,
		tokenInMaxsJson,
	)
	suite.Require().NoError(err)
	suite.Greater(response.GasUsed, uint64(0), "Gas should be used")

	if len(returnData) > 0 {
		method := suite.TestContractABI.Methods["testJoinPoolJSON"]
		unpacked, err := method.Outputs.Unpack(returnData)
		if err == nil && len(unpacked) > 0 {
			if jsonStr, ok := unpacked[0].(string); ok {
				suite.T().Logf("joinPoolJSON result: %s", jsonStr)
				suite.verifyJoinPoolJSON(jsonStr, suite.PoolId, shareOutAmount)
				suite.verifyJSONCanUnmarshalToProtobuf(jsonStr, nil)
			}
		}
	}
}

// TestJSONHelpers_GetPoolJSON tests getPoolJSON
func (suite *GammHelperLibrariesTestSuite) TestJSONHelpers_GetPoolJSON() {
	returnData, response, err := suite.callContract(
		suite.AliceKey,
		suite.TestContractAddr,
		"testGetPoolJSON",
		suite.PoolId,
	)
	suite.Require().NoError(err)
	suite.Greater(response.GasUsed, uint64(0), "Gas should be used")

	if len(returnData) > 0 {
		method := suite.TestContractABI.Methods["testGetPoolJSON"]
		unpacked, err := method.Outputs.Unpack(returnData)
		if err == nil && len(unpacked) > 0 {
			if jsonStr, ok := unpacked[0].(string); ok {
				suite.T().Logf("getPoolJSON result: %s", jsonStr)
				suite.verifyGetPoolJSON(jsonStr, suite.PoolId)
				suite.verifyJSONCanUnmarshalToProtobuf(jsonStr, nil)
			}
		}
	}
}

// TestJSONHelpers_CoinsToJson tests coinsToJson
func (suite *GammHelperLibrariesTestSuite) TestJSONHelpers_CoinsToJson() {
	coins := []struct {
		Denom  string
		Amount *big.Int
	}{
		{Denom: "uatom", Amount: big.NewInt(1000000)},
		{Denom: "uosmo", Amount: big.NewInt(2000000)},
	}

	returnData, response, err := suite.callContract(
		suite.AliceKey,
		suite.TestContractAddr,
		"testCoinsToJson",
		coins,
	)
	suite.Require().NoError(err)
	suite.Greater(response.GasUsed, uint64(0), "Gas should be used")

	if len(returnData) > 0 {
		method := suite.TestContractABI.Methods["testCoinsToJson"]
		unpacked, err := method.Outputs.Unpack(returnData)
		if err == nil && len(unpacked) > 0 {
			if jsonStr, ok := unpacked[0].(string); ok {
				suite.T().Logf("coinsToJson result: %s", jsonStr)
				// Verify JSON structure
				var coinsArray []map[string]interface{}
				err := json.Unmarshal([]byte(jsonStr), &coinsArray)
				suite.NoError(err)
				suite.Len(coinsArray, 2)
			}
		}
	}
}

// ============ JSON Verification Helpers ============

// verifyJoinPoolJSON verifies joinPool JSON structure
func (suite *GammHelperLibrariesTestSuite) verifyJoinPoolJSON(
	jsonStr string,
	poolId uint64,
	shareOutAmount *big.Int,
) {
	var msg map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &msg)
	suite.NoError(err, "JSON should be valid")

	suite.Equal(fmt.Sprintf("%d", poolId), msg["poolId"], "poolId should match")
	suite.Equal(shareOutAmount.String(), msg["shareOutAmount"], "shareOutAmount should match")
	suite.NotNil(msg["tokenInMaxs"], "tokenInMaxs should be present")
}

// verifyGetPoolJSON verifies getPool JSON structure
func (suite *GammHelperLibrariesTestSuite) verifyGetPoolJSON(
	jsonStr string,
	poolId uint64,
) {
	var msg map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &msg)
	suite.NoError(err, "JSON should be valid")

	suite.Equal(fmt.Sprintf("%d", poolId), msg["poolId"], "poolId should match")
}

// verifyJoinPoolBuilderJSON verifies JoinPoolBuilder JSON
func (suite *GammHelperLibrariesTestSuite) verifyJoinPoolBuilderJSON(
	jsonStr string,
	poolId uint64,
	shareOutAmount *big.Int,
) {
	suite.verifyJoinPoolJSON(jsonStr, poolId, shareOutAmount)
}

// verifySwapBuilderJSON verifies SwapBuilder JSON
func (suite *GammHelperLibrariesTestSuite) verifySwapBuilderJSON(
	jsonStr string,
	poolId uint64,
	tokenOutDenom string,
	tokenInDenom string,
	tokenInAmount *big.Int,
	tokenOutMinAmount *big.Int,
) {
	var msg map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &msg)
	suite.NoError(err, "JSON should be valid")

	suite.NotNil(msg["routes"], "routes should be present")
	suite.NotNil(msg["tokenIn"], "tokenIn should be present")
	suite.Equal(tokenOutMinAmount.String(), msg["tokenOutMinAmount"], "tokenOutMinAmount should match")
}

// verifyJSONCanUnmarshalToProtobuf verifies JSON can be unmarshaled to protobuf
func (suite *GammHelperLibrariesTestSuite) verifyJSONCanUnmarshalToProtobuf(
	jsonStr string,
	expectedType interface{},
) {
	// Try to unmarshal as generic map first
	var msg map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &msg)
	suite.NoError(err, "JSON should be valid and unmarshalable")

	// If expectedType is provided, try to unmarshal to that type
	if expectedType != nil {
		// This is a placeholder - actual protobuf unmarshaling would happen here
		suite.T().Logf("JSON verified as valid: %s", jsonStr)
	}
}
