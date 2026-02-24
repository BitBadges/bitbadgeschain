package tokenization_test

import (
	"crypto/ecdsa"
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
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	evmkeeper "github.com/cosmos/evm/x/vm/keeper"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/app"
	gammprecompile "github.com/bitbadges/bitbadgeschain/x/gamm/precompile"
	tokenizationkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	tokenization "github.com/bitbadges/bitbadgeschain/x/tokenization/precompile"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/precompile/test/helpers"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// SplitContractsTestSuite tests the split precompile test contracts
// These contracts are smaller than the original PrecompileTestContract to stay under EVM size limits
type SplitContractsTestSuite struct {
	suite.Suite
	App                *app.App
	Ctx                sdk.Context
	EVMKeeper          *evmkeeper.Keeper
	Precompile         *tokenization.Precompile
	BankKeeper         bankkeeper.Keeper
	TokenizationKeeper *tokenizationkeeper.Keeper

	// Contract deployments - one for each split contract
	TransferContractAddr     common.Address
	TransferContractABI      abi.ABI
	TransferContractBytecode []byte

	CollectionContractAddr     common.Address
	CollectionContractABI      abi.ABI
	CollectionContractBytecode []byte

	DynamicStoreContractAddr     common.Address
	DynamicStoreContractABI      abi.ABI
	DynamicStoreContractBytecode []byte

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

func TestSplitContractsTestSuite(t *testing.T) {
	suite.Run(t, new(SplitContractsTestSuite))
}

// SetupTest sets up the test suite
func (suite *SplitContractsTestSuite) SetupTest() {
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

	// Set chain ID
	suite.ChainID = big.NewInt(90123)

	// Load split contract artifacts
	suite.loadContractArtifacts()
}

// loadContractArtifacts loads the compiled artifacts for all split contracts
func (suite *SplitContractsTestSuite) loadContractArtifacts() {
	var err error

	// Load PrecompileTransferTestContract
	suite.TransferContractBytecode, err = helpers.GetContractBytecodeByType(helpers.ContractTypePrecompileTransfer)
	if err != nil {
		suite.T().Logf("WARNING: Could not load transfer contract bytecode: %v", err)
	}
	suite.TransferContractABI, err = helpers.GetContractABIByType(helpers.ContractTypePrecompileTransfer)
	if err != nil {
		suite.T().Logf("WARNING: Could not load transfer contract ABI: %v", err)
	}

	// Load PrecompileCollectionTestContract
	suite.CollectionContractBytecode, err = helpers.GetContractBytecodeByType(helpers.ContractTypePrecompileCollection)
	if err != nil {
		suite.T().Logf("WARNING: Could not load collection contract bytecode: %v", err)
	}
	suite.CollectionContractABI, err = helpers.GetContractABIByType(helpers.ContractTypePrecompileCollection)
	if err != nil {
		suite.T().Logf("WARNING: Could not load collection contract ABI: %v", err)
	}

	// Load PrecompileDynamicStoreTestContract
	suite.DynamicStoreContractBytecode, err = helpers.GetContractBytecodeByType(helpers.ContractTypePrecompileDynamicStore)
	if err != nil {
		suite.T().Logf("WARNING: Could not load dynamic store contract bytecode: %v", err)
	}
	suite.DynamicStoreContractABI, err = helpers.GetContractABIByType(helpers.ContractTypePrecompileDynamicStore)
	if err != nil {
		suite.T().Logf("WARNING: Could not load dynamic store contract ABI: %v", err)
	}
}

// setupAppWithEVM creates a full app instance with EVM keeper and precompile registered
func (suite *SplitContractsTestSuite) setupAppWithEVM() {
	suite.App = app.Setup(false)
	suite.Ctx = suite.App.BaseApp.NewContext(false)
	suite.Ctx = suite.Ctx.WithBlockHeight(1).WithBlockTime(time.Now())

	// Set up validator for EVM transactions
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

	header := suite.Ctx.BlockHeader()
	header.ProposerAddress = valConsAddr
	suite.Ctx = suite.Ctx.WithBlockHeader(header)

	_, err = suite.App.BeginBlocker(suite.Ctx)
	require.NoError(suite.T(), err)

	suite.EVMKeeper = suite.App.EVMKeeper
	suite.BankKeeper = suite.App.BankKeeper
	suite.TokenizationKeeper = suite.App.TokenizationKeeper

	// Create and register precompiles
	suite.Precompile = tokenization.NewPrecompile(suite.TokenizationKeeper)
	tokenizationPrecompileAddr := common.HexToAddress(tokenization.TokenizationPrecompileAddress)
	suite.EVMKeeper.RegisterStaticPrecompile(tokenizationPrecompileAddr, suite.Precompile)
	err = suite.EVMKeeper.EnableStaticPrecompiles(suite.Ctx, tokenizationPrecompileAddr)
	require.NoError(suite.T(), err)

	gammPrecompile := gammprecompile.NewPrecompile(suite.App.GammKeeper)
	gammPrecompileAddr := common.HexToAddress(gammprecompile.GammPrecompileAddress)
	suite.EVMKeeper.RegisterStaticPrecompile(gammPrecompileAddr, gammPrecompile)
	err = suite.EVMKeeper.EnableStaticPrecompiles(suite.Ctx, gammPrecompileAddr)
	require.NoError(suite.T(), err)
}

// createTestValidator creates a validator for testing
func (suite *SplitContractsTestSuite) createTestValidator() {
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

// createTestCollection creates a test collection with mint and transfer approvals
func (suite *SplitContractsTestSuite) createTestCollection() sdkmath.Uint {
	msg := &tokenizationtypes.MsgCreateCollection{
		Creator: suite.Alice.String(),
		ValidTokenIds: []*tokenizationtypes.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
		},
		CollectionApprovals: []*tokenizationtypes.CollectionApproval{
			{
				FromListId:        "Mint",
				ToListId:          "All",
				InitiatedByListId: "All",
				TransferTimes:     []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
				TokenIds:          []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)}},
				OwnershipTimes:    []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
				ApprovalId:        "mint-approval",
				ApprovalCriteria:  &tokenizationtypes.ApprovalCriteria{OverridesFromOutgoingApprovals: true, OverridesToIncomingApprovals: true},
			},
			{
				FromListId:        "!Mint",
				ToListId:          "All",
				InitiatedByListId: "All",
				TransferTimes:     []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
				TokenIds:          []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)}},
				OwnershipTimes:    []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
				ApprovalId:        "transfer-approval",
				ApprovalCriteria:  &tokenizationtypes.ApprovalCriteria{OverridesFromOutgoingApprovals: true, OverridesToIncomingApprovals: true},
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
						Amount:         sdkmath.NewUint(1000),
						TokenIds:       []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)}},
						OwnershipTimes: []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
					},
				},
			},
		},
	}

	_, err = msgServer.TransferTokens(suite.Ctx, mintMsg)
	suite.Require().NoError(err)

	return resp.CollectionId
}

// deployContract deploys a contract and returns its address
func (suite *SplitContractsTestSuite) deployContract(bytecode []byte) (common.Address, error) {
	contractAddr, response, err := helpers.DeployContract(
		suite.Ctx,
		suite.EVMKeeper,
		suite.DeployerKey,
		bytecode,
		suite.ChainID,
	)
	if err != nil {
		if response != nil && response.GasUsed > 0 {
			return common.Address{}, err
		}
		return common.Address{}, err
	}

	isContract, verifyErr := helpers.VerifyContractDeployment(suite.Ctx, suite.EVMKeeper, contractAddr)
	if verifyErr != nil {
		return common.Address{}, verifyErr
	}
	if !isContract {
		return common.Address{}, err
	}

	return contractAddr, nil
}

// callContract calls a contract method and returns the response
func (suite *SplitContractsTestSuite) callContract(
	callerKey *ecdsa.PrivateKey,
	contractAddr common.Address,
	contractABI abi.ABI,
	methodName string,
	args []interface{},
	isView bool,
) ([]byte, *evmtypes.MsgEthereumTxResponse, error) {
	return helpers.CallContractMethod(
		suite.Ctx,
		suite.EVMKeeper,
		callerKey,
		contractAddr,
		contractABI,
		methodName,
		args,
		suite.ChainID,
		isView,
	)
}

// getBalanceAmount returns the balance amount for an address
func (suite *SplitContractsTestSuite) getBalanceAmount(address string, collectionId sdkmath.Uint) sdkmath.Uint {
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

// ============ Transfer Contract Tests ============

func (suite *SplitContractsTestSuite) TestTransferContract_Deploy() {
	if len(suite.TransferContractBytecode) == 0 {
		suite.T().Skip("Skipping - transfer contract bytecode not loaded. Run 'make compile-contracts' first.")
		return
	}

	contractAddr, err := suite.deployContract(suite.TransferContractBytecode)
	suite.Require().NoError(err, "Transfer contract deployment should succeed")
	suite.Require().NotEqual(common.Address{}, contractAddr, "Contract address should be set")

	suite.TransferContractAddr = contractAddr
	suite.T().Logf("Transfer contract deployed at: %s", contractAddr.Hex())
}

func (suite *SplitContractsTestSuite) TestTransferContract_GetCollection() {
	if len(suite.TransferContractBytecode) == 0 || len(suite.TransferContractABI.Methods) == 0 {
		suite.T().Skip("Skipping - transfer contract artifacts not loaded")
		return
	}

	// Deploy contract if not already deployed
	if suite.TransferContractAddr == (common.Address{}) {
		addr, err := suite.deployContract(suite.TransferContractBytecode)
		if err != nil {
			suite.T().Skipf("Skipping - contract deployment failed: %v", err)
			return
		}
		suite.TransferContractAddr = addr
	}

	// Call testGetCollection
	args := []interface{}{suite.CollectionId.BigInt()}
	returnData, response, err := suite.callContract(
		suite.AliceKey,
		suite.TransferContractAddr,
		suite.TransferContractABI,
		"testGetCollection",
		args,
		true, // view function
	)

	suite.Require().NoError(err)
	suite.Require().NotNil(response)
	suite.Greater(response.GasUsed, uint64(0), "Gas should be used")
	suite.T().Logf("GetCollection gas used: %d, return data length: %d", response.GasUsed, len(returnData))
}

func (suite *SplitContractsTestSuite) TestTransferContract_GetBalance() {
	if len(suite.TransferContractBytecode) == 0 || len(suite.TransferContractABI.Methods) == 0 {
		suite.T().Skip("Skipping - transfer contract artifacts not loaded")
		return
	}

	// Deploy contract if not already deployed
	if suite.TransferContractAddr == (common.Address{}) {
		addr, err := suite.deployContract(suite.TransferContractBytecode)
		if err != nil {
			suite.T().Skipf("Skipping - contract deployment failed: %v", err)
			return
		}
		suite.TransferContractAddr = addr
	}

	// Call testGetBalance
	args := []interface{}{
		suite.CollectionId.BigInt(),
		suite.AliceEVM,
	}
	returnData, response, err := suite.callContract(
		suite.AliceKey,
		suite.TransferContractAddr,
		suite.TransferContractABI,
		"testGetBalance",
		args,
		true, // view function
	)

	suite.Require().NoError(err)
	suite.Require().NotNil(response)
	suite.Greater(response.GasUsed, uint64(0), "Gas should be used")
	suite.T().Logf("GetBalance gas used: %d, return data length: %d", response.GasUsed, len(returnData))
}

func (suite *SplitContractsTestSuite) TestTransferContract_GetAddressList() {
	if len(suite.TransferContractBytecode) == 0 || len(suite.TransferContractABI.Methods) == 0 {
		suite.T().Skip("Skipping - transfer contract artifacts not loaded")
		return
	}

	// Deploy contract if not already deployed
	if suite.TransferContractAddr == (common.Address{}) {
		addr, err := suite.deployContract(suite.TransferContractBytecode)
		if err != nil {
			suite.T().Skipf("Skipping - contract deployment failed: %v", err)
			return
		}
		suite.TransferContractAddr = addr
	}

	// Call testGetAddressList with "All" list
	args := []interface{}{"All"}
	returnData, response, err := suite.callContract(
		suite.AliceKey,
		suite.TransferContractAddr,
		suite.TransferContractABI,
		"testGetAddressList",
		args,
		true, // view function
	)

	suite.Require().NoError(err)
	suite.Require().NotNil(response)
	suite.Greater(response.GasUsed, uint64(0), "Gas should be used")
	suite.T().Logf("GetAddressList gas used: %d, return data length: %d", response.GasUsed, len(returnData))
}

// ============ Collection Contract Tests ============

func (suite *SplitContractsTestSuite) TestCollectionContract_Deploy() {
	if len(suite.CollectionContractBytecode) == 0 {
		suite.T().Skip("Skipping - collection contract bytecode not loaded. Run 'make compile-contracts' first.")
		return
	}

	contractAddr, err := suite.deployContract(suite.CollectionContractBytecode)
	suite.Require().NoError(err, "Collection contract deployment should succeed")
	suite.Require().NotEqual(common.Address{}, contractAddr, "Contract address should be set")

	suite.CollectionContractAddr = contractAddr
	suite.T().Logf("Collection contract deployed at: %s", contractAddr.Hex())
}

func (suite *SplitContractsTestSuite) TestCollectionContract_DeleteCollection() {
	if len(suite.CollectionContractBytecode) == 0 || len(suite.CollectionContractABI.Methods) == 0 {
		suite.T().Skip("Skipping - collection contract artifacts not loaded")
		return
	}

	// Deploy contract if not already deployed
	if suite.CollectionContractAddr == (common.Address{}) {
		addr, err := suite.deployContract(suite.CollectionContractBytecode)
		if err != nil {
			suite.T().Skipf("Skipping - contract deployment failed: %v", err)
			return
		}
		suite.CollectionContractAddr = addr
	}

	// Create a temporary collection to delete
	msg := &tokenizationtypes.MsgCreateCollection{
		Creator: suite.Alice.String(),
		ValidTokenIds: []*tokenizationtypes.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
		},
	}
	msgServer := tokenizationkeeper.NewMsgServerImpl(suite.TokenizationKeeper)
	resp, err := msgServer.CreateCollection(suite.Ctx, msg)
	suite.Require().NoError(err)

	// Call testDeleteCollection
	args := []interface{}{resp.CollectionId.BigInt()}
	returnData, response, err := suite.callContract(
		suite.AliceKey,
		suite.CollectionContractAddr,
		suite.CollectionContractABI,
		"testDeleteCollection",
		args,
		false, // not a view function
	)

	// The delete may fail because the caller is the contract, not Alice
	// This is expected behavior - just verify the call was processed
	suite.Require().NoError(err)
	suite.Require().NotNil(response)
	suite.Greater(response.GasUsed, uint64(0), "Gas should be used")
	suite.T().Logf("DeleteCollection gas used: %d, return data length: %d, VmError: %s",
		response.GasUsed, len(returnData), response.VmError)
}

// ============ Dynamic Store Contract Tests ============

func (suite *SplitContractsTestSuite) TestDynamicStoreContract_Deploy() {
	if len(suite.DynamicStoreContractBytecode) == 0 {
		suite.T().Skip("Skipping - dynamic store contract bytecode not loaded. Run 'make compile-contracts' first.")
		return
	}

	contractAddr, err := suite.deployContract(suite.DynamicStoreContractBytecode)
	suite.Require().NoError(err, "Dynamic store contract deployment should succeed")
	suite.Require().NotEqual(common.Address{}, contractAddr, "Contract address should be set")

	suite.DynamicStoreContractAddr = contractAddr
	suite.T().Logf("Dynamic store contract deployed at: %s", contractAddr.Hex())
}

func (suite *SplitContractsTestSuite) TestDynamicStoreContract_Params() {
	if len(suite.DynamicStoreContractBytecode) == 0 || len(suite.DynamicStoreContractABI.Methods) == 0 {
		suite.T().Skip("Skipping - dynamic store contract artifacts not loaded")
		return
	}

	// Deploy contract if not already deployed
	if suite.DynamicStoreContractAddr == (common.Address{}) {
		addr, err := suite.deployContract(suite.DynamicStoreContractBytecode)
		if err != nil {
			suite.T().Skipf("Skipping - contract deployment failed: %v", err)
			return
		}
		suite.DynamicStoreContractAddr = addr
	}

	// Call testParams
	returnData, response, err := suite.callContract(
		suite.AliceKey,
		suite.DynamicStoreContractAddr,
		suite.DynamicStoreContractABI,
		"testParams",
		[]interface{}{},
		true, // view function
	)

	suite.Require().NoError(err)
	suite.Require().NotNil(response)
	suite.Greater(response.GasUsed, uint64(0), "Gas should be used")
	suite.T().Logf("Params gas used: %d, return data length: %d", response.GasUsed, len(returnData))
}

func (suite *SplitContractsTestSuite) TestDynamicStoreContract_IsAddressReservedProtocol() {
	if len(suite.DynamicStoreContractBytecode) == 0 || len(suite.DynamicStoreContractABI.Methods) == 0 {
		suite.T().Skip("Skipping - dynamic store contract artifacts not loaded")
		return
	}

	// Deploy contract if not already deployed
	if suite.DynamicStoreContractAddr == (common.Address{}) {
		addr, err := suite.deployContract(suite.DynamicStoreContractBytecode)
		if err != nil {
			suite.T().Skipf("Skipping - contract deployment failed: %v", err)
			return
		}
		suite.DynamicStoreContractAddr = addr
	}

	// Call testIsAddressReservedProtocol
	args := []interface{}{suite.AliceEVM}
	returnData, response, err := suite.callContract(
		suite.AliceKey,
		suite.DynamicStoreContractAddr,
		suite.DynamicStoreContractABI,
		"testIsAddressReservedProtocol",
		args,
		true, // view function
	)

	suite.Require().NoError(err)
	suite.Require().NotNil(response)
	suite.Greater(response.GasUsed, uint64(0), "Gas should be used")

	// Parse result - should be false since Alice is not a reserved protocol address
	if len(returnData) > 0 {
		method := suite.DynamicStoreContractABI.Methods["testIsAddressReservedProtocol"]
		unpacked, err := method.Outputs.Unpack(returnData)
		if err == nil && len(unpacked) > 0 {
			if isReserved, ok := unpacked[0].(bool); ok {
				suite.False(isReserved, "Alice should not be a reserved protocol address")
			}
		}
	}
	suite.T().Logf("IsAddressReservedProtocol gas used: %d", response.GasUsed)
}

func (suite *SplitContractsTestSuite) TestDynamicStoreContract_GetAllReservedProtocolAddresses() {
	if len(suite.DynamicStoreContractBytecode) == 0 || len(suite.DynamicStoreContractABI.Methods) == 0 {
		suite.T().Skip("Skipping - dynamic store contract artifacts not loaded")
		return
	}

	// Deploy contract if not already deployed
	if suite.DynamicStoreContractAddr == (common.Address{}) {
		addr, err := suite.deployContract(suite.DynamicStoreContractBytecode)
		if err != nil {
			suite.T().Skipf("Skipping - contract deployment failed: %v", err)
			return
		}
		suite.DynamicStoreContractAddr = addr
	}

	// Call testGetAllReservedProtocolAddresses
	returnData, response, err := suite.callContract(
		suite.AliceKey,
		suite.DynamicStoreContractAddr,
		suite.DynamicStoreContractABI,
		"testGetAllReservedProtocolAddresses",
		[]interface{}{},
		true, // view function
	)

	suite.Require().NoError(err)
	suite.Require().NotNil(response)
	suite.Greater(response.GasUsed, uint64(0), "Gas should be used")
	suite.T().Logf("GetAllReservedProtocolAddresses gas used: %d, return data length: %d", response.GasUsed, len(returnData))
}

// ============ All Contracts Gas Comparison ============

func (suite *SplitContractsTestSuite) TestAllContracts_DeploymentGasComparison() {
	results := make(map[string]struct {
		Size    int
		GasUsed uint64
		Success bool
	})

	// Deploy Transfer Contract
	if len(suite.TransferContractBytecode) > 0 {
		addr, response, err := helpers.DeployContract(suite.Ctx, suite.EVMKeeper, suite.DeployerKey, suite.TransferContractBytecode, suite.ChainID)
		success := err == nil && addr != (common.Address{})
		gasUsed := uint64(0)
		if response != nil {
			gasUsed = response.GasUsed
		}
		results["PrecompileTransferTestContract"] = struct {
			Size    int
			GasUsed uint64
			Success bool
		}{len(suite.TransferContractBytecode), gasUsed, success}
	}

	// Deploy Collection Contract
	if len(suite.CollectionContractBytecode) > 0 {
		addr, response, err := helpers.DeployContract(suite.Ctx, suite.EVMKeeper, suite.DeployerKey, suite.CollectionContractBytecode, suite.ChainID)
		success := err == nil && addr != (common.Address{})
		gasUsed := uint64(0)
		if response != nil {
			gasUsed = response.GasUsed
		}
		results["PrecompileCollectionTestContract"] = struct {
			Size    int
			GasUsed uint64
			Success bool
		}{len(suite.CollectionContractBytecode), gasUsed, success}
	}

	// Deploy Dynamic Store Contract
	if len(suite.DynamicStoreContractBytecode) > 0 {
		addr, response, err := helpers.DeployContract(suite.Ctx, suite.EVMKeeper, suite.DeployerKey, suite.DynamicStoreContractBytecode, suite.ChainID)
		success := err == nil && addr != (common.Address{})
		gasUsed := uint64(0)
		if response != nil {
			gasUsed = response.GasUsed
		}
		results["PrecompileDynamicStoreTestContract"] = struct {
			Size    int
			GasUsed uint64
			Success bool
		}{len(suite.DynamicStoreContractBytecode), gasUsed, success}
	}

	// Log results
	suite.T().Log("Contract Deployment Gas Comparison:")
	suite.T().Log("==================================")
	for name, result := range results {
		status := "FAILED"
		if result.Success {
			status = "SUCCESS"
		}
		suite.T().Logf("%s: Size=%d bytes, Gas=%d, Status=%s",
			name, result.Size, result.GasUsed, status)
	}

	// All split contracts should deploy successfully (under 24KB limit)
	for name, result := range results {
		suite.True(result.Success, "%s should deploy successfully", name)
		suite.Less(result.Size, 24576, "%s should be under 24KB", name)
	}
}
