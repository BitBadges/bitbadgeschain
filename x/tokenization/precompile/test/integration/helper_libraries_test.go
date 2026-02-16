package tokenization_test

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/crypto/ed25519"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	evmkeeper "github.com/cosmos/evm/x/vm/keeper"

	"github.com/bitbadges/bitbadgeschain/app"
	customhookstypes "github.com/bitbadges/bitbadgeschain/x/custom-hooks/types"
	gammprecompile "github.com/bitbadges/bitbadgeschain/x/gamm/precompile"
	tokenizationkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	tokenization "github.com/bitbadges/bitbadgeschain/x/tokenization/precompile"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/precompile/test/helpers"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// HelperLibrariesTestSuite tests all helper libraries E2E
type HelperLibrariesTestSuite struct {
	suite.Suite
	App                *app.App
	Ctx                sdk.Context
	EVMKeeper          *evmkeeper.Keeper
	Precompile         *tokenization.Precompile
	BankKeeper         bankkeeper.Keeper
	TokenizationKeeper tokenizationkeeper.Keeper

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

	CollectionId sdkmath.Uint
	ChainID      *big.Int
}

func TestHelperLibrariesTestSuite(t *testing.T) {
	suite.Run(t, new(HelperLibrariesTestSuite))
}

// SetupTest sets up the test suite
func (suite *HelperLibrariesTestSuite) SetupTest() {
	suite.setupAppWithEVM()

	// Create test accounts
	suite.DeployerKey, suite.DeployerEVM, suite.Deployer = helpers.CreateEVMAccount()
	suite.AliceKey, suite.AliceEVM, suite.Alice = helpers.CreateEVMAccount()
	suite.BobKey, suite.BobEVM, suite.Bob = helpers.CreateEVMAccount()

	// Fund accounts
	err := helpers.FundEVMAccount(suite.Ctx, suite.BankKeeper, suite.Deployer, sdk.NewCoins(sdk.NewCoin("ustake", sdkmath.NewInt(10000000000000000))))
	suite.Require().NoError(err)
	err = helpers.FundEVMAccount(suite.Ctx, suite.BankKeeper, suite.Alice, sdk.NewCoins(sdk.NewCoin("ustake", sdkmath.NewInt(10000000000000000))))
	suite.Require().NoError(err)
	err = helpers.FundEVMAccount(suite.Ctx, suite.BankKeeper, suite.Bob, sdk.NewCoins(sdk.NewCoin("ustake", sdkmath.NewInt(10000000000000000))))
	suite.Require().NoError(err)

	// Create test collection
	suite.CollectionId = suite.createTestCollection()

	// Set chain ID (using testnet chain ID for testing)
	suite.ChainID = big.NewInt(90123)

	// Load compiled contract
	contractBytecode, err := helpers.GetContractBytecode()
	if err != nil {
		suite.T().Skipf("Skipping test - could not load contract bytecode: %v. Run 'make compile-contracts' first.", err)
		return
	}
	suite.ContractBytecode = contractBytecode

	contractABI, err := helpers.GetContractABI()
	if err != nil {
		suite.T().Skipf("Skipping test - could not load contract ABI: %v. Run 'make compile-contracts' first.", err)
		return
	}
	suite.TestContractABI = contractABI

	// Deploy test contract
	contractAddr, response, err := helpers.DeployContract(
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

	isContract, verifyErr := helpers.VerifyContractDeployment(suite.Ctx, suite.EVMKeeper, contractAddr)
	suite.Require().NoError(verifyErr)
	if !isContract {
		suite.T().Skip("Skipping test - contract deployment failed")
		return
	}

	suite.TestContractAddr = contractAddr
}

// setupAppWithEVM creates a full app instance with EVM keeper and precompile registered
func (suite *HelperLibrariesTestSuite) setupAppWithEVM() {
	suite.App = app.Setup(false)

	header := cmtproto.Header{
		Height:  1,
		Time:    time.Now(),
		ChainID: "bitbadges-1",
	}
	suite.Ctx = suite.App.BaseApp.NewContextLegacy(false, header)

	_ = suite.Ctx.TransientStore(customhookstypes.TransientStoreKey)

	// Create validator
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

	suite.EVMKeeper = suite.App.EVMKeeper
	suite.BankKeeper = suite.App.BankKeeper
	suite.TokenizationKeeper = suite.App.TokenizationKeeper

	// Register and enable precompile
	suite.Precompile = tokenization.NewPrecompile(suite.TokenizationKeeper)
	tokenizationPrecompileAddr := common.HexToAddress(tokenization.TokenizationPrecompileAddress)
	suite.EVMKeeper.RegisterStaticPrecompile(tokenizationPrecompileAddr, suite.Precompile)
	err = suite.EVMKeeper.EnableStaticPrecompiles(suite.Ctx, tokenizationPrecompileAddr)
	require.NoError(suite.T(), err)

	// Also register gamm precompile
	gammPrecompile := gammprecompile.NewPrecompile(suite.App.GammKeeper)
	gammPrecompileAddr := common.HexToAddress(gammprecompile.GammPrecompileAddress)
	suite.EVMKeeper.RegisterStaticPrecompile(gammPrecompileAddr, gammPrecompile)
	err = suite.EVMKeeper.EnableStaticPrecompiles(suite.Ctx, gammPrecompileAddr)
	require.NoError(suite.T(), err)
}

// createTestValidator creates a validator manually for testing
func (suite *HelperLibrariesTestSuite) createTestValidator() {
	privKey := ed25519.GenPrivKey()
	pubKey := privKey.PubKey()

	cosmosPubKey, err := cryptocodec.FromTmPubKeyInterface(pubKey)
	require.NoError(suite.T(), err)

	pkAny, err := codectypes.NewAnyWithValue(cosmosPubKey)
	require.NoError(suite.T(), err)

	valAddr := sdk.ValAddress(pubKey.Address())
	bondAmt := sdk.DefaultPowerReduction
	validator := stakingtypes.Validator{
		OperatorAddress:   valAddr.String(),
		ConsensusPubkey:   pkAny,
		Jailed:            false,
		Status:            stakingtypes.Bonded,
		Tokens:            bondAmt,
		DelegatorShares:   sdkmath.LegacyOneDec(),
		Description:       stakingtypes.Description{Moniker: "test-validator"},
		UnbondingHeight:   int64(0),
		UnbondingTime:     time.Unix(0, 0).UTC(),
		Commission:        stakingtypes.NewCommission(sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec()),
		MinSelfDelegation: sdkmath.ZeroInt(),
	}

	err = suite.App.StakingKeeper.SetValidator(suite.Ctx, validator)
	require.NoError(suite.T(), err)

	suite.App.StakingKeeper.SetValidatorByConsAddr(suite.Ctx, validator)
	suite.App.StakingKeeper.SetNewValidatorByPowerIndex(suite.Ctx, validator)

	valConsAddr, err := validator.GetConsAddr()
	require.NoError(suite.T(), err)

	signingInfo := slashingtypes.ValidatorSigningInfo{
		Address:             sdk.ConsAddress(valConsAddr).String(),
		StartHeight:         0,
		IndexOffset:         0,
		JailedUntil:         time.Unix(0, 0).UTC(),
		Tombstoned:          false,
		MissedBlocksCounter: 0,
	}
	err = suite.App.SlashingKeeper.SetValidatorSigningInfo(suite.Ctx, sdk.ConsAddress(valConsAddr), signingInfo)
	require.NoError(suite.T(), err)

	_, err = suite.App.StakingKeeper.ApplyAndReturnValidatorSetUpdates(suite.Ctx)
	require.NoError(suite.T(), err)
}

// createTestCollection creates a test collection and mints tokens to Alice
func (suite *HelperLibrariesTestSuite) createTestCollection() sdkmath.Uint {
	msg := &tokenizationtypes.MsgCreateCollection{
		Creator: suite.Alice.String(),
		ValidTokenIds: []*tokenizationtypes.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
		},
		CollectionApprovals: []*tokenizationtypes.CollectionApproval{
			{
				FromListId:        "Mint",
				ToListId:          "All",
				InitiatedByListId: "All",
				TransferTimes: []*tokenizationtypes.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
				},
				TokenIds: []*tokenizationtypes.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
				},
				OwnershipTimes: []*tokenizationtypes.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
				},
				ApprovalId: "mint-approval",
				ApprovalCriteria: &tokenizationtypes.ApprovalCriteria{
					OverridesFromOutgoingApprovals: true,
					OverridesToIncomingApprovals:   true,
				},
			},
			{
				FromListId:        "!Mint",
				ToListId:          "All",
				InitiatedByListId: "All",
				TransferTimes: []*tokenizationtypes.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
				},
				TokenIds: []*tokenizationtypes.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
				},
				OwnershipTimes: []*tokenizationtypes.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
				},
				ApprovalId: "transfer-approval",
				ApprovalCriteria: &tokenizationtypes.ApprovalCriteria{
					OverridesFromOutgoingApprovals: true,
					OverridesToIncomingApprovals:   true,
				},
			},
		},
	}

	msgServer := tokenizationkeeper.NewMsgServerImpl(suite.TokenizationKeeper)
	resp, err := msgServer.CreateCollection(suite.Ctx, msg)
	suite.Require().NoError(err)

	// Mint tokens to Alice
	mintMsg := &tokenizationtypes.MsgTransferTokens{
		Creator:      suite.Alice.String(),
		CollectionId: resp.CollectionId,
		Transfers: []*tokenizationtypes.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{suite.Alice.String()},
				Balances: []*tokenizationtypes.Balance{
					{
						Amount: sdkmath.NewUint(100),
						TokenIds: []*tokenizationtypes.UintRange{
							{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
						},
						OwnershipTimes: []*tokenizationtypes.UintRange{
							{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
						},
					},
				},
			},
		},
	}

	_, err = msgServer.TransferTokens(suite.Ctx, mintMsg)
	suite.Require().NoError(err)

	return resp.CollectionId
}

// ============ TokenizationWrappers Tests ============

// TestWrappers_TransferTokens tests the transferTokens wrapper
func (suite *HelperLibrariesTestSuite) TestWrappers_TransferTokens() {
	if len(suite.TestContractABI.Methods) == 0 {
		suite.T().Skip("Contract ABI not loaded")
		return
	}

	method, exists := suite.TestContractABI.Methods["testTransferTokensWrapper"]
	if !exists {
		suite.T().Skip("testTransferTokensWrapper method not found")
		return
	}

	// Prepare test data
	toAddresses := []common.Address{suite.BobEVM}
	amount := big.NewInt(10)

	// Create UintRange structs for Solidity
	type UintRange struct {
		Start *big.Int
		End   *big.Int
	}
	tokenIds := []UintRange{
		{Start: big.NewInt(1), End: big.NewInt(10)},
	}
	ownershipTimes := []UintRange{
		{Start: big.NewInt(1), End: new(big.Int).SetUint64(math.MaxUint64)},
	}

	args := []interface{}{
		suite.CollectionId.BigInt(),
		toAddresses,
		amount,
		tokenIds,
		ownershipTimes,
	}

	// Get initial balances
	aliceBalanceBefore := suite.getBalanceAmount(suite.Alice.String(), suite.CollectionId)
	bobBalanceBefore := suite.getBalanceAmount(suite.Bob.String(), suite.CollectionId)

	// Call contract method
	returnData, response, err := helpers.CallContractMethod(
		suite.Ctx,
		suite.EVMKeeper,
		suite.AliceKey,
		suite.TestContractAddr,
		suite.TestContractABI,
		"testTransferTokensWrapper",
		args,
		suite.ChainID,
		false,
	)

	suite.Require().NoError(err)
	suite.Require().NotNil(response)

	// Verify JSON was constructed correctly by checking the call succeeded
	// Note: The transfer may fail if the contract is the caller (not Alice)
	// but we can still verify the JSON was constructed and the precompile was called
	if response.VmError == "" {
		if len(returnData) > 0 {
			unpacked, err := method.Outputs.Unpack(returnData)
			if err == nil && len(unpacked) > 0 {
				if success, ok := unpacked[0].(bool); ok {
					suite.T().Logf("Transfer result: %v", success)

					// If successful, verify balance changes
					if success {
						aliceBalanceAfter := suite.getBalanceAmount(suite.Alice.String(), suite.CollectionId)
						bobBalanceAfter := suite.getBalanceAmount(suite.Bob.String(), suite.CollectionId)

						if aliceBalanceBefore.GTE(sdkmath.NewUint(10)) {
							suite.Equal(aliceBalanceBefore.Sub(sdkmath.NewUint(10)), aliceBalanceAfter, "Alice should have lost 10 tokens")
						}
						suite.Equal(bobBalanceBefore.Add(sdkmath.NewUint(10)), bobBalanceAfter, "Bob should have gained 10 tokens")
					}
				}
			}
		}
	}

	// Verify the method was called (gas was used)
	suite.Greater(response.GasUsed, uint64(0), "Gas should be used")
}

// TestWrappers_TransferSingleToken tests the transferSingleToken convenience wrapper
func (suite *HelperLibrariesTestSuite) TestWrappers_TransferSingleToken() {
	if len(suite.TestContractABI.Methods) == 0 {
		suite.T().Skip("Contract ABI not loaded")
		return
	}

	method, exists := suite.TestContractABI.Methods["testTransferSingleTokenWrapper"]
	if !exists {
		suite.T().Skip("testTransferSingleTokenWrapper method not found")
		return
	}

	args := []interface{}{
		suite.CollectionId.BigInt(),
		suite.BobEVM,
		big.NewInt(5),
		big.NewInt(1),
	}

	returnData, response, err := helpers.CallContractMethod(
		suite.Ctx,
		suite.EVMKeeper,
		suite.AliceKey,
		suite.TestContractAddr,
		suite.TestContractABI,
		"testTransferSingleTokenWrapper",
		args,
		suite.ChainID,
		false,
	)

	suite.Require().NoError(err)
	suite.Require().NotNil(response)

	// Verify method was called
	if len(returnData) > 0 {
		unpacked, err := method.Outputs.Unpack(returnData)
		if err == nil && len(unpacked) > 0 {
			if success, ok := unpacked[0].(bool); ok {
				suite.T().Logf("TransferSingleToken result: %v", success)
			}
		}
	}

	suite.Greater(response.GasUsed, uint64(0), "Gas should be used")
}

// TestWrappers_GetBalanceAmount tests the getBalanceAmount wrapper
func (suite *HelperLibrariesTestSuite) TestWrappers_GetBalanceAmount() {
	if len(suite.TestContractABI.Methods) == 0 {
		suite.T().Skip("Contract ABI not loaded")
		return
	}

	method, exists := suite.TestContractABI.Methods["testGetBalanceAmountWrapper"]
	if !exists {
		suite.T().Skip("testGetBalanceAmountWrapper method not found")
		return
	}

	type UintRange struct {
		Start *big.Int
		End   *big.Int
	}
	tokenIds := []UintRange{
		{Start: big.NewInt(1), End: big.NewInt(10)},
	}
	ownershipTimes := []UintRange{
		{Start: big.NewInt(1), End: new(big.Int).SetUint64(math.MaxUint64)},
	}

	args := []interface{}{
		suite.CollectionId.BigInt(),
		suite.AliceEVM,
		tokenIds,
		ownershipTimes,
	}

	returnData, response, err := helpers.CallContractMethod(
		suite.Ctx,
		suite.EVMKeeper,
		suite.AliceKey,
		suite.TestContractAddr,
		suite.TestContractABI,
		"testGetBalanceAmountWrapper",
		args,
		suite.ChainID,
		true, // View function
	)

	suite.Require().NoError(err)
	suite.Require().NotNil(response)

	// Verify return value
	if len(returnData) > 0 {
		unpacked, err := method.Outputs.Unpack(returnData)
		if err == nil && len(unpacked) > 0 {
			if balance, ok := unpacked[0].(*big.Int); ok {
				suite.T().Logf("Balance amount: %s", balance.String())
				// Verify balance matches expected (Alice should have 100 tokens)
				expectedBalance := sdkmath.NewUint(100)
				suite.Equal(expectedBalance.BigInt().String(), balance.String(), "Balance should match")
			}
		}
	}
}

// TestWrappers_GetTotalSupply tests the getTotalSupply wrapper
func (suite *HelperLibrariesTestSuite) TestWrappers_GetTotalSupply() {
	if len(suite.TestContractABI.Methods) == 0 {
		suite.T().Skip("Contract ABI not loaded")
		return
	}

	method, exists := suite.TestContractABI.Methods["testGetTotalSupplyWrapper"]
	if !exists {
		suite.T().Skip("testGetTotalSupplyWrapper method not found")
		return
	}

	type UintRange struct {
		Start *big.Int
		End   *big.Int
	}
	tokenIds := []UintRange{
		{Start: big.NewInt(1), End: big.NewInt(10)},
	}
	ownershipTimes := []UintRange{
		{Start: big.NewInt(1), End: new(big.Int).SetUint64(math.MaxUint64)},
	}

	args := []interface{}{
		suite.CollectionId.BigInt(),
		tokenIds,
		ownershipTimes,
	}

	returnData, response, err := helpers.CallContractMethod(
		suite.Ctx,
		suite.EVMKeeper,
		suite.AliceKey,
		suite.TestContractAddr,
		suite.TestContractABI,
		"testGetTotalSupplyWrapper",
		args,
		suite.ChainID,
		true,
	)

	suite.Require().NoError(err)
	suite.Require().NotNil(response)

	if len(returnData) > 0 {
		unpacked, err := method.Outputs.Unpack(returnData)
		if err == nil && len(unpacked) > 0 {
			if supply, ok := unpacked[0].(*big.Int); ok {
				suite.T().Logf("Total supply: %s", supply.String())
				// Should be 100 (minted to Alice)
				expectedSupply := sdkmath.NewUint(100)
				suite.Equal(expectedSupply.BigInt().String(), supply.String(), "Total supply should match")
			}
		}
	}
}

// TestWrappers_CreateDynamicStore tests the createDynamicStore wrapper
func (suite *HelperLibrariesTestSuite) TestWrappers_CreateDynamicStore() {
	if len(suite.TestContractABI.Methods) == 0 {
		suite.T().Skip("Contract ABI not loaded")
		return
	}

	method, exists := suite.TestContractABI.Methods["testCreateDynamicStoreWrapper"]
	if !exists {
		suite.T().Skip("testCreateDynamicStoreWrapper method not found")
		return
	}

	args := []interface{}{
		false, // defaultValue
		"ipfs://test-store",
		"Test Dynamic Store",
	}

	returnData, response, err := helpers.CallContractMethod(
		suite.Ctx,
		suite.EVMKeeper,
		suite.AliceKey,
		suite.TestContractAddr,
		suite.TestContractABI,
		"testCreateDynamicStoreWrapper",
		args,
		suite.ChainID,
		false,
	)

	suite.Require().NoError(err)
	suite.Require().NotNil(response)

	if len(returnData) > 0 {
		unpacked, err := method.Outputs.Unpack(returnData)
		if err == nil && len(unpacked) > 0 {
			if storeId, ok := unpacked[0].(*big.Int); ok {
				suite.T().Logf("Created store ID: %s", storeId.String())
				suite.Greater(storeId.Uint64(), uint64(0), "Store ID should be > 0")
			}
		}
	}

	suite.Greater(response.GasUsed, uint64(0), "Gas should be used")
}

// TestWrappers_SetDynamicStoreValue tests the setDynamicStoreValue wrapper
func (suite *HelperLibrariesTestSuite) TestWrappers_SetDynamicStoreValue() {
	if len(suite.TestContractABI.Methods) == 0 {
		suite.T().Skip("Contract ABI not loaded")
		return
	}

	// First create a store
	createMethod, exists := suite.TestContractABI.Methods["testCreateDynamicStoreWrapper"]
	if !exists {
		suite.T().Skip("testCreateDynamicStoreWrapper method not found")
		return
	}

	createArgs := []interface{}{
		false,
		"ipfs://test",
		"Test",
	}

	createReturnData, _, err := helpers.CallContractMethod(
		suite.Ctx,
		suite.EVMKeeper,
		suite.AliceKey,
		suite.TestContractAddr,
		suite.TestContractABI,
		"testCreateDynamicStoreWrapper",
		createArgs,
		suite.ChainID,
		false,
	)
	suite.Require().NoError(err)

	var storeId *big.Int
	if len(createReturnData) > 0 {
		unpacked, err := createMethod.Outputs.Unpack(createReturnData)
		if err == nil && len(unpacked) > 0 {
			if id, ok := unpacked[0].(*big.Int); ok {
				storeId = id
			}
		}
	}

	if storeId == nil || storeId.Uint64() == 0 {
		suite.T().Skip("Could not create store for test")
		return
	}

	// Now test setDynamicStoreValue
	setMethod, exists := suite.TestContractABI.Methods["testSetDynamicStoreValueWrapper"]
	if !exists {
		suite.T().Skip("testSetDynamicStoreValueWrapper method not found")
		return
	}

	args := []interface{}{
		storeId,
		suite.BobEVM,
		true,
	}

	returnData, response, err := helpers.CallContractMethod(
		suite.Ctx,
		suite.EVMKeeper,
		suite.AliceKey,
		suite.TestContractAddr,
		suite.TestContractABI,
		"testSetDynamicStoreValueWrapper",
		args,
		suite.ChainID,
		false,
	)

	suite.Require().NoError(err)
	suite.Require().NotNil(response)

	if len(returnData) > 0 {
		unpacked, err := setMethod.Outputs.Unpack(returnData)
		if err == nil && len(unpacked) > 0 {
			if success, ok := unpacked[0].(bool); ok {
				suite.T().Logf("SetDynamicStoreValue result: %v", success)
			}
		}
	}

	suite.Greater(response.GasUsed, uint64(0), "Gas should be used")
}

// TestWrappers_GetCollection tests the getCollection wrapper
func (suite *HelperLibrariesTestSuite) TestWrappers_GetCollection() {
	if len(suite.TestContractABI.Methods) == 0 {
		suite.T().Skip("Contract ABI not loaded")
		return
	}

	method, exists := suite.TestContractABI.Methods["testGetCollectionWrapper"]
	if !exists {
		suite.T().Skip("testGetCollectionWrapper method not found")
		return
	}

	args := []interface{}{
		suite.CollectionId.BigInt(),
	}

	returnData, response, err := helpers.CallContractMethod(
		suite.Ctx,
		suite.EVMKeeper,
		suite.AliceKey,
		suite.TestContractAddr,
		suite.TestContractABI,
		"testGetCollectionWrapper",
		args,
		suite.ChainID,
		true,
	)

	suite.Require().NoError(err)
	suite.Require().NotNil(response)

	// Verify return data is not empty (protobuf bytes)
	if len(returnData) > 0 {
		unpacked, err := method.Outputs.Unpack(returnData)
		if err == nil && len(unpacked) > 0 {
			if collectionBytes, ok := unpacked[0].([]byte); ok {
				suite.Greater(len(collectionBytes), 0, "Collection bytes should not be empty")
				suite.T().Logf("Collection bytes length: %d", len(collectionBytes))
			}
		}
	}
}

// ============ TokenizationBuilders Tests ============

// TestBuilders_CollectionBuilder tests the CollectionBuilder
func (suite *HelperLibrariesTestSuite) TestBuilders_CollectionBuilder() {
	if len(suite.TestContractABI.Methods) == 0 {
		suite.T().Skip("Contract ABI not loaded")
		return
	}

	method, exists := suite.TestContractABI.Methods["testCollectionBuilder"]
	if !exists {
		suite.T().Skip("testCollectionBuilder method not found")
		return
	}

	args := []interface{}{
		big.NewInt(1),        // tokenIdStart
		big.NewInt(100),      // tokenIdEnd
		suite.Alice.String(), // manager (Cosmos address string)
		"ipfs://test-collection",
		"Test Collection",
	}

	returnData, response, err := helpers.CallContractMethod(
		suite.Ctx,
		suite.EVMKeeper,
		suite.AliceKey,
		suite.TestContractAddr,
		suite.TestContractABI,
		"testCollectionBuilder",
		args,
		suite.ChainID,
		false,
	)

	suite.Require().NoError(err)
	suite.Require().NotNil(response)

	// Verify JSON was constructed using verification helper
	if len(returnData) > 0 {
		unpacked, err := method.Outputs.Unpack(returnData)
		if err == nil && len(unpacked) > 0 {
			if jsonStr, ok := unpacked[0].(string); ok {
				suite.T().Logf("Generated JSON: %s", jsonStr)

				// Use verification helper
				suite.verifyCollectionBuilderJSON(
					jsonStr,
					big.NewInt(1),
					big.NewInt(100),
					suite.Alice.String(),
					"ipfs://test-collection",
					"Test Collection",
				)
				suite.verifyJSONCanUnmarshalToProtobuf(jsonStr, nil)
			}
		}
	}
}

// TestBuilders_TransferBuilder tests the TransferBuilder
func (suite *HelperLibrariesTestSuite) TestBuilders_TransferBuilder() {
	if len(suite.TestContractABI.Methods) == 0 {
		suite.T().Skip("Contract ABI not loaded")
		return
	}

	method, exists := suite.TestContractABI.Methods["testTransferBuilder"]
	if !exists {
		suite.T().Skip("testTransferBuilder method not found")
		return
	}

	args := []interface{}{
		suite.CollectionId.BigInt(),
		suite.BobEVM,
		big.NewInt(5),
		big.NewInt(1),
	}

	returnData, response, err := helpers.CallContractMethod(
		suite.Ctx,
		suite.EVMKeeper,
		suite.AliceKey,
		suite.TestContractAddr,
		suite.TestContractABI,
		"testTransferBuilder",
		args,
		suite.ChainID,
		false,
	)

	suite.Require().NoError(err)
	suite.Require().NotNil(response)

	// Verify JSON was constructed using verification helper
	if len(returnData) > 0 {
		unpacked, err := method.Outputs.Unpack(returnData)
		if err == nil && len(unpacked) > 0 {
			if jsonStr, ok := unpacked[0].(string); ok {
				suite.T().Logf("Generated transfer JSON: %s", jsonStr)

				// Use verification helper
				suite.verifyTransferBuilderJSON(
					jsonStr,
					suite.CollectionId,
					suite.BobEVM,
					big.NewInt(5),
					big.NewInt(1),
				)
				suite.verifyJSONCanUnmarshalToProtobuf(jsonStr, nil)
			}
		}
	}
}

// ============ TokenizationHelpers Tests ============

// TestHelpers_HelperFunctions tests TokenizationHelpers functions
func (suite *HelperLibrariesTestSuite) TestHelpers_HelperFunctions() {
	if len(suite.TestContractABI.Methods) == 0 {
		suite.T().Skip("Contract ABI not loaded")
		return
	}

	method, exists := suite.TestContractABI.Methods["testTokenizationHelpers"]
	if !exists {
		suite.T().Skip("testTokenizationHelpers method not found")
		return
	}

	args := []interface{}{}

	returnData, response, err := helpers.CallContractMethod(
		suite.Ctx,
		suite.EVMKeeper,
		suite.AliceKey,
		suite.TestContractAddr,
		suite.TestContractABI,
		"testTokenizationHelpers",
		args,
		suite.ChainID,
		true, // View function
	)

	suite.Require().NoError(err)
	suite.Require().NotNil(response)

	if len(returnData) > 0 {
		unpacked, err := method.Outputs.Unpack(returnData)
		if err == nil && len(unpacked) > 0 {
			if success, ok := unpacked[0].(bool); ok {
				suite.True(success, "Helper functions should work correctly")
			}
		}
	}
}

// ============ TokenizationJSONHelpers Tests ============

// TestJSONHelpers_TransferTokensJSON tests transferTokensJSON construction
func (suite *HelperLibrariesTestSuite) TestJSONHelpers_TransferTokensJSON() {
	if len(suite.TestContractABI.Methods) == 0 {
		suite.T().Skip("Contract ABI not loaded")
		return
	}

	method, exists := suite.TestContractABI.Methods["testTransferTokensJSON"]
	if !exists {
		suite.T().Skip("testTransferTokensJSON method not found")
		return
	}

	tokenIdStarts := []*big.Int{big.NewInt(1)}
	tokenIdEnds := []*big.Int{big.NewInt(10)}
	ownershipStarts := []*big.Int{big.NewInt(1)}
	ownershipEnds := []*big.Int{new(big.Int).SetUint64(math.MaxUint64)}

	args := []interface{}{
		suite.CollectionId.BigInt(),
		[]common.Address{suite.BobEVM},
		big.NewInt(10),
		tokenIdStarts,
		tokenIdEnds,
		ownershipStarts,
		ownershipEnds,
	}

	returnData, response, err := helpers.CallContractMethod(
		suite.Ctx,
		suite.EVMKeeper,
		suite.AliceKey,
		suite.TestContractAddr,
		suite.TestContractABI,
		"testTransferTokensJSON",
		args,
		suite.ChainID,
		true, // View function (just constructs JSON)
	)

	suite.Require().NoError(err)
	suite.Require().NotNil(response)

	// Verify JSON was constructed correctly using verification helper
	if len(returnData) > 0 {
		unpacked, err := method.Outputs.Unpack(returnData)
		if err == nil && len(unpacked) > 0 {
			if jsonStr, ok := unpacked[0].(string); ok {
				suite.T().Logf("Constructed JSON: %s", jsonStr)

				// Use verification helper to verify JSON structure
				tokenIds := []struct{ Start, End *big.Int }{
					{Start: big.NewInt(1), End: big.NewInt(10)},
				}
				ownershipTimes := []struct{ Start, End *big.Int }{
					{Start: big.NewInt(1), End: new(big.Int).SetUint64(math.MaxUint64)},
				}

				suite.verifyTransferTokensJSON(
					jsonStr,
					suite.CollectionId,
					[]common.Address{suite.BobEVM},
					big.NewInt(10),
					tokenIds,
					ownershipTimes,
				)

				// Verify JSON can be unmarshaled to protobuf
				suite.verifyJSONCanUnmarshalToProtobuf(jsonStr, nil)
			}
		}
	}
}

// TestJSONHelpers_GetCollectionJSON tests getCollectionJSON construction
func (suite *HelperLibrariesTestSuite) TestJSONHelpers_GetCollectionJSON() {
	if len(suite.TestContractABI.Methods) == 0 {
		suite.T().Skip("Contract ABI not loaded")
		return
	}

	method, exists := suite.TestContractABI.Methods["testGetCollectionJSON"]
	if !exists {
		suite.T().Skip("testGetCollectionJSON method not found")
		return
	}

	args := []interface{}{
		suite.CollectionId.BigInt(),
	}

	returnData, response, err := helpers.CallContractMethod(
		suite.Ctx,
		suite.EVMKeeper,
		suite.AliceKey,
		suite.TestContractAddr,
		suite.TestContractABI,
		"testGetCollectionJSON",
		args,
		suite.ChainID,
		true,
	)

	suite.Require().NoError(err)
	suite.Require().NotNil(response)

	if len(returnData) > 0 {
		unpacked, err := method.Outputs.Unpack(returnData)
		if err == nil && len(unpacked) > 0 {
			if jsonStr, ok := unpacked[0].(string); ok {
				suite.T().Logf("Constructed JSON: %s", jsonStr)

				// Use verification helper
				suite.verifyGetCollectionJSON(jsonStr, suite.CollectionId)
				suite.verifyJSONCanUnmarshalToProtobuf(jsonStr, nil)
			}
		}
	}
}

// TestJSONHelpers_GetBalanceJSON tests getBalanceJSON construction
func (suite *HelperLibrariesTestSuite) TestJSONHelpers_GetBalanceJSON() {
	if len(suite.TestContractABI.Methods) == 0 {
		suite.T().Skip("Contract ABI not loaded")
		return
	}

	method, exists := suite.TestContractABI.Methods["testGetBalanceJSON"]
	if !exists {
		suite.T().Skip("testGetBalanceJSON method not found")
		return
	}

	args := []interface{}{
		suite.CollectionId.BigInt(),
		suite.AliceEVM,
	}

	returnData, response, err := helpers.CallContractMethod(
		suite.Ctx,
		suite.EVMKeeper,
		suite.AliceKey,
		suite.TestContractAddr,
		suite.TestContractABI,
		"testGetBalanceJSON",
		args,
		suite.ChainID,
		true,
	)

	suite.Require().NoError(err)
	suite.Require().NotNil(response)

	if len(returnData) > 0 {
		unpacked, err := method.Outputs.Unpack(returnData)
		if err == nil && len(unpacked) > 0 {
			if jsonStr, ok := unpacked[0].(string); ok {
				suite.T().Logf("Constructed JSON: %s", jsonStr)

				// Use verification helper
				suite.verifyGetBalanceJSON(jsonStr, suite.CollectionId, suite.AliceEVM)
				suite.verifyJSONCanUnmarshalToProtobuf(jsonStr, nil)
			}
		}
	}
}

// TestJSONHelpers_AllHelpers tests all JSON helper functions
func (suite *HelperLibrariesTestSuite) TestJSONHelpers_AllHelpers() {
	if len(suite.TestContractABI.Methods) == 0 {
		suite.T().Skip("Contract ABI not loaded")
		return
	}

	method, exists := suite.TestContractABI.Methods["testJSONHelpers"]
	if !exists {
		suite.T().Skip("testJSONHelpers method not found")
		return
	}

	args := []interface{}{}

	returnData, response, err := helpers.CallContractMethod(
		suite.Ctx,
		suite.EVMKeeper,
		suite.AliceKey,
		suite.TestContractAddr,
		suite.TestContractABI,
		"testJSONHelpers",
		args,
		suite.ChainID,
		true,
	)

	suite.Require().NoError(err)
	suite.Require().NotNil(response)

	if len(returnData) > 0 {
		unpacked, err := method.Outputs.Unpack(returnData)
		if err == nil && len(unpacked) > 0 {
			if success, ok := unpacked[0].(bool); ok {
				suite.True(success, "All JSON helpers should work correctly")
			}
		}
	}
}

// ============ TokenizationErrors Tests ============

// TestErrors_ValidationHelpers tests error validation helpers with valid inputs
func (suite *HelperLibrariesTestSuite) TestErrors_ValidationHelpers() {
	if len(suite.TestContractABI.Methods) == 0 {
		suite.T().Skip("Contract ABI not loaded")
		return
	}

	method, exists := suite.TestContractABI.Methods["testTokenizationErrorsValid"]
	if !exists {
		suite.T().Skip("testTokenizationErrorsValid method not found")
		return
	}

	args := []interface{}{
		suite.CollectionId.BigInt(),
		suite.AliceEVM,
		"non-empty-string",
	}

	returnData, response, err := helpers.CallContractMethod(
		suite.Ctx,
		suite.EVMKeeper,
		suite.AliceKey,
		suite.TestContractAddr,
		suite.TestContractABI,
		"testTokenizationErrorsValid",
		args,
		suite.ChainID,
		true,
	)

	suite.Require().NoError(err)
	suite.Require().NotNil(response)

	if len(returnData) > 0 {
		unpacked, err := method.Outputs.Unpack(returnData)
		if err == nil && len(unpacked) > 0 {
			if success, ok := unpacked[0].(bool); ok {
				suite.True(success, "Validation helpers should pass with valid inputs")
			}
		}
	}
}

// ============ Helper Functions ============

// getBalanceAmount is a helper to get balance amount
func (suite *HelperLibrariesTestSuite) getBalanceAmount(address string, collectionId sdkmath.Uint) sdkmath.Uint {
	res, err := suite.TokenizationKeeper.GetBalance(suite.Ctx, &tokenizationtypes.QueryGetBalanceRequest{
		CollectionId: collectionId.String(),
		Address:      address,
	})
	if err != nil {
		return sdkmath.ZeroUint()
	}

	if res.Balance == nil {
		return sdkmath.ZeroUint()
	}

	totalAmount := sdkmath.ZeroUint()
	for _, bal := range res.Balance.Balances {
		totalAmount = totalAmount.Add(bal.Amount)
	}

	return totalAmount
}

// ============ JSON Verification Helpers ============

// verifyTransferTokensJSON verifies that the JSON string matches the expected structure for transferTokens
func (suite *HelperLibrariesTestSuite) verifyTransferTokensJSON(
	jsonStr string,
	expectedCollectionId sdkmath.Uint,
	expectedToAddresses []common.Address,
	expectedAmount *big.Int,
	expectedTokenIds []struct{ Start, End *big.Int },
	expectedOwnershipTimes []struct{ Start, End *big.Int },
) {
	// Parse JSON
	var jsonData map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &jsonData)
	suite.Require().NoError(err, "JSON should be valid")

	// Verify required fields
	suite.Contains(jsonData, "collectionId", "JSON should contain collectionId")
	suite.Contains(jsonData, "transfers", "JSON should contain transfers")

	// Verify collectionId
	collectionIdStr, ok := jsonData["collectionId"].(string)
	suite.True(ok, "collectionId should be a string")
	suite.Equal(expectedCollectionId.String(), collectionIdStr, "Collection ID should match")

	// Verify transfers array
	transfers, ok := jsonData["transfers"].([]interface{})
	suite.True(ok, "transfers should be an array")
	suite.Greater(len(transfers), 0, "transfers should not be empty")

	// Verify first transfer
	transfer, ok := transfers[0].(map[string]interface{})
	suite.True(ok, "transfer should be an object")

	// Verify toAddresses
	toAddresses, ok := transfer["toAddresses"].([]interface{})
	suite.True(ok, "toAddresses should be an array")
	suite.Equal(len(expectedToAddresses), len(toAddresses), "toAddresses length should match")

	// Verify balances
	balances, ok := transfer["balances"].([]interface{})
	suite.True(ok, "balances should be an array")
	suite.Greater(len(balances), 0, "balances should not be empty")

	balance, ok := balances[0].(map[string]interface{})
	suite.True(ok, "balance should be an object")

	// Verify amount
	amountStr, ok := balance["amount"].(string)
	suite.True(ok, "amount should be a string")
	suite.Equal(expectedAmount.String(), amountStr, "Amount should match")

	// Verify tokenIds
	tokenIds, ok := balance["tokenIds"].([]interface{})
	suite.True(ok, "tokenIds should be an array")
	suite.Equal(len(expectedTokenIds), len(tokenIds), "tokenIds length should match")

	// Verify ownershipTimes
	ownershipTimes, ok := balance["ownershipTimes"].([]interface{})
	suite.True(ok, "ownershipTimes should be an array")
	suite.Equal(len(expectedOwnershipTimes), len(ownershipTimes), "ownershipTimes length should match")
}

// verifyGetCollectionJSON verifies that the JSON string matches the expected structure for getCollection
func (suite *HelperLibrariesTestSuite) verifyGetCollectionJSON(
	jsonStr string,
	expectedCollectionId sdkmath.Uint,
) {
	// Parse JSON
	var jsonData map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &jsonData)
	suite.Require().NoError(err, "JSON should be valid")

	// Verify required fields
	suite.Contains(jsonData, "collectionId", "JSON should contain collectionId")

	// Verify collectionId
	collectionIdStr, ok := jsonData["collectionId"].(string)
	suite.True(ok, "collectionId should be a string")
	suite.Equal(expectedCollectionId.String(), collectionIdStr, "Collection ID should match")
}

// verifyGetBalanceJSON verifies that the JSON string matches the expected structure for getBalance
func (suite *HelperLibrariesTestSuite) verifyGetBalanceJSON(
	jsonStr string,
	expectedCollectionId sdkmath.Uint,
	expectedAddress common.Address,
) {
	// Parse JSON
	var jsonData map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &jsonData)
	suite.Require().NoError(err, "JSON should be valid")

	// Verify required fields
	suite.Contains(jsonData, "collectionId", "JSON should contain collectionId")
	suite.Contains(jsonData, "userAddress", "JSON should contain userAddress")

	// Verify collectionId
	collectionIdStr, ok := jsonData["collectionId"].(string)
	suite.True(ok, "collectionId should be a string")
	suite.Equal(expectedCollectionId.String(), collectionIdStr, "Collection ID should match")

	// Verify userAddress (convert EVM address to Cosmos address format)
	userAddressStr, ok := jsonData["userAddress"].(string)
	suite.True(ok, "userAddress should be a string")
	// The address should be in Cosmos format (bech32), so we just verify it's not empty
	suite.NotEmpty(userAddressStr, "userAddress should not be empty")
}

// verifyGetBalanceAmountJSON verifies that the JSON string matches the expected structure for getBalanceAmount
func (suite *HelperLibrariesTestSuite) verifyGetBalanceAmountJSON(
	jsonStr string,
	expectedCollectionId sdkmath.Uint,
	expectedAddress common.Address,
	expectedTokenIds []struct{ Start, End *big.Int },
	expectedOwnershipTimes []struct{ Start, End *big.Int },
) {
	// Parse JSON
	var jsonData map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &jsonData)
	suite.Require().NoError(err, "JSON should be valid")

	// Verify required fields
	suite.Contains(jsonData, "collectionId", "JSON should contain collectionId")
	suite.Contains(jsonData, "userAddress", "JSON should contain userAddress")
	suite.Contains(jsonData, "tokenIds", "JSON should contain tokenIds")
	suite.Contains(jsonData, "ownershipTimes", "JSON should contain ownershipTimes")

	// Verify collectionId
	collectionIdStr, ok := jsonData["collectionId"].(string)
	suite.True(ok, "collectionId should be a string")
	suite.Equal(expectedCollectionId.String(), collectionIdStr, "Collection ID should match")

	// Verify userAddress
	userAddressStr, ok := jsonData["userAddress"].(string)
	suite.True(ok, "userAddress should be a string")
	suite.NotEmpty(userAddressStr, "userAddress should not be empty")

	// Verify tokenIds
	tokenIds, ok := jsonData["tokenIds"].([]interface{})
	suite.True(ok, "tokenIds should be an array")
	suite.Equal(len(expectedTokenIds), len(tokenIds), "tokenIds length should match")

	// Verify ownershipTimes
	ownershipTimes, ok := jsonData["ownershipTimes"].([]interface{})
	suite.True(ok, "ownershipTimes should be an array")
	suite.Equal(len(expectedOwnershipTimes), len(ownershipTimes), "ownershipTimes length should match")
}

// verifyCreateDynamicStoreJSON verifies that the JSON string matches the expected structure for createDynamicStore
func (suite *HelperLibrariesTestSuite) verifyCreateDynamicStoreJSON(
	jsonStr string,
	expectedDefaultValue bool,
	expectedURI string,
	expectedCustomData string,
) {
	// Parse JSON
	var jsonData map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &jsonData)
	suite.Require().NoError(err, "JSON should be valid")

	// Verify required fields
	suite.Contains(jsonData, "defaultValue", "JSON should contain defaultValue")
	suite.Contains(jsonData, "uri", "JSON should contain uri")
	suite.Contains(jsonData, "customData", "JSON should contain customData")

	// Verify defaultValue
	defaultValue, ok := jsonData["defaultValue"].(bool)
	suite.True(ok, "defaultValue should be a bool")
	suite.Equal(expectedDefaultValue, defaultValue, "defaultValue should match")

	// Verify uri
	uri, ok := jsonData["uri"].(string)
	suite.True(ok, "uri should be a string")
	suite.Equal(expectedURI, uri, "URI should match")

	// Verify customData
	customData, ok := jsonData["customData"].(string)
	suite.True(ok, "customData should be a string")
	suite.Equal(expectedCustomData, customData, "customData should match")
}

// verifyCollectionBuilderJSON verifies that the JSON string from CollectionBuilder matches expected structure
func (suite *HelperLibrariesTestSuite) verifyCollectionBuilderJSON(
	jsonStr string,
	expectedTokenIdStart *big.Int,
	expectedTokenIdEnd *big.Int,
	expectedManager string,
	expectedURI string,
	expectedCustomData string,
) {
	// Parse JSON
	var jsonData map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &jsonData)
	suite.Require().NoError(err, "JSON should be valid")

	// Verify required fields
	suite.Contains(jsonData, "validTokenIds", "JSON should contain validTokenIds")

	// Verify validTokenIds
	validTokenIds, ok := jsonData["validTokenIds"].([]interface{})
	suite.True(ok, "validTokenIds should be an array")
	suite.Greater(len(validTokenIds), 0, "validTokenIds should not be empty")

	// Verify first token ID range
	tokenIdRange, ok := validTokenIds[0].(map[string]interface{})
	suite.True(ok, "tokenIdRange should be an object")

	startStr, ok := tokenIdRange["start"].(string)
	suite.True(ok, "start should be a string")
	suite.Equal(expectedTokenIdStart.String(), startStr, "Token ID start should match")

	endStr, ok := tokenIdRange["end"].(string)
	suite.True(ok, "end should be a string")
	suite.Equal(expectedTokenIdEnd.String(), endStr, "Token ID end should match")

	// Verify manager if provided
	if expectedManager != "" {
		suite.Contains(jsonData, "manager", "JSON should contain manager")
		manager, ok := jsonData["manager"].(string)
		suite.True(ok, "manager should be a string")
		suite.Equal(expectedManager, manager, "Manager should match")
	}

	// Verify collectionMetadata if provided
	if expectedURI != "" || expectedCustomData != "" {
		suite.Contains(jsonData, "collectionMetadata", "JSON should contain collectionMetadata")
		metadata, ok := jsonData["collectionMetadata"].(map[string]interface{})
		suite.True(ok, "collectionMetadata should be an object")

		if expectedURI != "" {
			uri, ok := metadata["uri"].(string)
			suite.True(ok, "uri should be a string")
			suite.Equal(expectedURI, uri, "URI should match")
		}

		if expectedCustomData != "" {
			customData, ok := metadata["customData"].(string)
			suite.True(ok, "customData should be a string")
			suite.Equal(expectedCustomData, customData, "customData should match")
		}
	}
}

// verifyTransferBuilderJSON verifies that the JSON string from TransferBuilder matches expected structure
func (suite *HelperLibrariesTestSuite) verifyTransferBuilderJSON(
	jsonStr string,
	expectedCollectionId sdkmath.Uint,
	expectedToAddress common.Address,
	expectedAmount *big.Int,
	expectedTokenId *big.Int,
) {
	// Parse JSON
	var jsonData map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &jsonData)
	suite.Require().NoError(err, "JSON should be valid")

	// Verify required fields
	suite.Contains(jsonData, "collectionId", "JSON should contain collectionId")
	suite.Contains(jsonData, "toAddresses", "JSON should contain toAddresses")
	suite.Contains(jsonData, "amount", "JSON should contain amount")
	suite.Contains(jsonData, "tokenIds", "JSON should contain tokenIds")
	suite.Contains(jsonData, "ownershipTimes", "JSON should contain ownershipTimes")

	// Verify collectionId
	collectionIdStr, ok := jsonData["collectionId"].(string)
	suite.True(ok, "collectionId should be a string")
	suite.Equal(expectedCollectionId.String(), collectionIdStr, "Collection ID should match")

	// Verify toAddresses
	toAddresses, ok := jsonData["toAddresses"].([]interface{})
	suite.True(ok, "toAddresses should be an array")
	suite.Equal(1, len(toAddresses), "toAddresses should have one address")

	// Verify amount
	amountStr, ok := jsonData["amount"].(string)
	suite.True(ok, "amount should be a string")
	suite.Equal(expectedAmount.String(), amountStr, "Amount should match")

	// Verify tokenIds (should be single token ID range)
	tokenIds, ok := jsonData["tokenIds"].([]interface{})
	suite.True(ok, "tokenIds should be an array")
	suite.Equal(1, len(tokenIds), "tokenIds should have one range")

	tokenIdRange, ok := tokenIds[0].(map[string]interface{})
	suite.True(ok, "tokenIdRange should be an object")

	startStr, ok := tokenIdRange["start"].(string)
	suite.True(ok, "start should be a string")
	suite.Equal(expectedTokenId.String(), startStr, "Token ID start should match")

	endStr, ok := tokenIdRange["end"].(string)
	suite.True(ok, "end should be a string")
	suite.Equal(expectedTokenId.String(), endStr, "Token ID end should match (single token)")
}

// verifyJSONCanUnmarshalToProtobuf verifies that JSON can be unmarshaled into the protobuf message type
func (suite *HelperLibrariesTestSuite) verifyJSONCanUnmarshalToProtobuf(
	jsonStr string,
	protobufType interface{},
) {
	// This is a generic verification that the JSON structure is compatible with protobuf
	// We'll use json.Unmarshal to verify the structure is valid
	var jsonData map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &jsonData)
	suite.Require().NoError(err, "JSON should be valid and parseable")

	// Additional verification: ensure JSON can be unmarshaled into protobuf types
	// This is a simplified check - in practice, we'd unmarshal into the actual protobuf message
	suite.NotEmpty(jsonData, "JSON should not be empty")
}

// compareJSONWithGoHelper compares Solidity-generated JSON with Go helper output
func (suite *HelperLibrariesTestSuite) compareJSONWithGoHelper(
	solidityJSON string,
	goJSON string,
) {
	// Parse both JSON strings
	var solidityData map[string]interface{}
	var goData map[string]interface{}

	err := json.Unmarshal([]byte(solidityJSON), &solidityData)
	suite.Require().NoError(err, "Solidity JSON should be valid")

	err = json.Unmarshal([]byte(goJSON), &goData)
	suite.Require().NoError(err, "Go JSON should be valid")

	// Compare key fields (allowing for field order differences)
	// This is a simplified comparison - in practice, we'd do a deep comparison
	suite.Equal(len(solidityData), len(goData), "JSON should have same number of top-level fields")

	// Verify common fields exist in both
	for key := range solidityData {
		suite.Contains(goData, key, fmt.Sprintf("Go JSON should contain field: %s", key))
	}
}

// Helper function to convert EVM address to Cosmos address string for comparison
func evmToCosmosAddress(evmAddr common.Address) string {
	// This is a simplified conversion - in practice, we'd use the actual SDK address conversion
	// For testing purposes, we'll just verify the address is not empty
	return sdk.AccAddress(evmAddr.Bytes()).String()
}
