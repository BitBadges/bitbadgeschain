// Package tokenization_test provides E2E tests for EVM query challenges in approval and invariant flows.
package tokenization_test

import (
	"crypto/ecdsa"
	"encoding/hex"
	"math"
	"math/big"
	"strings"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/cometbft/cometbft/crypto/ed25519"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/app"
	gammprecompile "github.com/bitbadges/bitbadgeschain/x/gamm/precompile"
	tokenizationkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	tokenization "github.com/bitbadges/bitbadgeschain/x/tokenization/precompile"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/precompile/test/helpers"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// EVMQueryChallengesE2ESuite runs E2E tests for EVM query challenges (approval + invariant).
type EVMQueryChallengesE2ESuite struct {
	suite.Suite
	App                *app.App
	Ctx                sdk.Context
	TokenizationKeeper *tokenizationkeeper.Keeper

	DeployerKey *ecdsa.PrivateKey
	AliceKey    *ecdsa.PrivateKey
	BobKey      *ecdsa.PrivateKey
	Deployer    sdk.AccAddress
	Alice       sdk.AccAddress
	Bob         sdk.AccAddress
	Charlie     sdk.AccAddress
	Dave        sdk.AccAddress
	ChainID     *big.Int
}

func TestEVMQueryChallengesE2ESuite(t *testing.T) {
	suite.Run(t, new(EVMQueryChallengesE2ESuite))
}

func (s *EVMQueryChallengesE2ESuite) SetupTest() {
	s.App = app.Setup(false)
	s.Ctx = s.App.BaseApp.NewContext(false).WithBlockHeight(1).WithBlockTime(time.Now())

	// Set up validator for EVM
	var firstValidator stakingtypes.ValidatorI
	s.App.StakingKeeper.IterateValidators(s.Ctx, func(_ int64, val stakingtypes.ValidatorI) (stop bool) {
		firstValidator = val
		return true
	})
	if firstValidator == nil {
		s.createTestValidator()
		s.App.StakingKeeper.IterateValidators(s.Ctx, func(_ int64, val stakingtypes.ValidatorI) (stop bool) {
			firstValidator = val
			return true
		})
	}
	s.Require().NotNil(firstValidator)
	valConsAddr, err := firstValidator.GetConsAddr()
	require.NoError(s.T(), err)
	voteInfos := []abci.VoteInfo{{Validator: abci.Validator{Address: valConsAddr, Power: 1000}, BlockIdFlag: cmtproto.BlockIDFlagCommit}}
	s.Ctx = s.Ctx.WithVoteInfos(voteInfos)
	header := s.Ctx.BlockHeader()
	header.ProposerAddress = valConsAddr
	s.Ctx = s.Ctx.WithBlockHeader(header)
	_, err = s.App.BeginBlocker(s.Ctx)
	require.NoError(s.T(), err)

	s.TokenizationKeeper = s.App.TokenizationKeeper

	// Create test accounts
	s.DeployerKey, _, s.Deployer = helpers.CreateEVMAccount()
	s.AliceKey, _, s.Alice = helpers.CreateEVMAccount()
	_, _, s.Bob = helpers.CreateEVMAccount()
	_, _, s.Charlie = helpers.CreateEVMAccount()
	_, _, s.Dave = helpers.CreateEVMAccount()

	// Fund with ustake for gas
	for _, acc := range []sdk.AccAddress{s.Deployer, s.Alice, s.Bob, s.Charlie, s.Dave} {
		err := helpers.FundEVMAccount(s.Ctx, s.App.BankKeeper, acc, sdk.NewCoins(sdk.NewCoin("ustake", sdkmath.NewInt(10000000000000000))))
		s.Require().NoError(err)
	}
	s.ChainID = big.NewInt(90123)

	// Register and enable tokenization + gamm precompiles (required for contracts and EVM queries)
	tokenizationPrecompile := tokenization.NewPrecompile(s.TokenizationKeeper)
	tokenizationPrecompileAddr := common.HexToAddress(tokenization.TokenizationPrecompileAddress)
	s.App.EVMKeeper.RegisterStaticPrecompile(tokenizationPrecompileAddr, tokenizationPrecompile)
	err = s.App.EVMKeeper.EnableStaticPrecompiles(s.Ctx, tokenizationPrecompileAddr)
	require.NoError(s.T(), err)
	gammPrecompile := gammprecompile.NewPrecompile(s.App.GammKeeper)
	gammPrecompileAddr := common.HexToAddress(gammprecompile.GammPrecompileAddress)
	s.App.EVMKeeper.RegisterStaticPrecompile(gammPrecompileAddr, gammPrecompile)
	_ = s.App.EVMKeeper.EnableStaticPrecompiles(s.Ctx, gammPrecompileAddr)
	// Enable all precompiles (including bank 0x0804) so MinBankBalanceChecker can call the bank precompile
	err = s.App.EnableAllPrecompiles(s.Ctx)
	require.NoError(s.T(), err)
}

func (s *EVMQueryChallengesE2ESuite) createTestValidator() {
	privKey := ed25519.GenPrivKey()
	pubKey := privKey.PubKey()
	cosmosPubKey, _ := cryptocodec.FromTmPubKeyInterface(pubKey)
	pkAny, _ := codectypes.NewAnyWithValue(cosmosPubKey)
	valAddr := sdk.ValAddress(pubKey.Address())
	bondAmt := sdk.DefaultPowerReduction
	validator := stakingtypes.Validator{
		OperatorAddress:   valAddr.String(),
		ConsensusPubkey:   pkAny,
		Jailed:            false,
		Status:            stakingtypes.Bonded,
		Tokens:            bondAmt,
		DelegatorShares:   sdkmath.LegacyOneDec(),
		Description:       stakingtypes.Description{Moniker: "test"},
		UnbondingHeight:   int64(0),
		UnbondingTime:     time.Unix(0, 0).UTC(),
		Commission:        stakingtypes.NewCommission(sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec()),
		MinSelfDelegation: sdkmath.ZeroInt(),
	}
	_ = s.App.StakingKeeper.SetValidator(s.Ctx, validator)
	_ = s.App.StakingKeeper.SetValidatorByConsAddr(s.Ctx, validator)
	_ = s.App.StakingKeeper.SetNewValidatorByPowerIndex(s.Ctx, validator)
	valConsAddr, _ := validator.GetConsAddr()
	signingInfo := slashingtypes.ValidatorSigningInfo{
		Address: sdk.ConsAddress(valConsAddr).String(),
		StartHeight: 0, IndexOffset: 0, JailedUntil: time.Unix(0, 0).UTC(),
		Tombstoned: false, MissedBlocksCounter: 0,
	}
	_ = s.App.SlashingKeeper.SetValidatorSigningInfo(s.Ctx, sdk.ConsAddress(valConsAddr), signingInfo)
	_, _ = s.App.StakingKeeper.ApplyAndReturnValidatorSetUpdates(s.Ctx)
}

// TestE2E_ApprovalFlow_MinBankBalance checks min balance using simulated balances in the approval flow.
// This test uses a contract with simulated balances instead of the bank precompile to avoid
// ERC20 token pair registration requirements in the test environment.
func (s *EVMQueryChallengesE2ESuite) TestE2E_ApprovalFlow_MinBankBalance() {
	bytecode, err := helpers.GetContractBytecodeByType(helpers.ContractTypeMinBankBalanceChecker)
	if err != nil {
		s.T().Skipf("MinBankBalanceChecker bytecode not found (run make compile-contracts): %v", err)
		return
	}

	// Deploy MinBankBalanceChecker
	contractAddr, _, err := helpers.DeployContract(s.Ctx, s.App.EVMKeeper, s.DeployerKey, bytecode, s.ChainID)
	s.Require().NoError(err)
	isContract, _ := helpers.VerifyContractDeployment(s.Ctx, s.App.EVMKeeper, contractAddr)
	s.Require().True(isContract, "MinBankBalanceChecker should be deployed")

	// Set simulated balance for Alice (1000 units) so checkMinBalance(account, 1000) passes
	// This replaces the bank precompile approach that requires ERC20 token pair registration
	setBalanceABI := `[{"inputs":[{"internalType":"address","name":"account","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"setSimulatedBalance","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
	setBalanceContractABI, err := abi.JSON(strings.NewReader(setBalanceABI))
	s.Require().NoError(err)
	aliceEVMAddr := common.BytesToAddress(s.Alice.Bytes())
	_, _, err = helpers.CallContractMethod(s.Ctx, s.App.EVMKeeper, s.DeployerKey, contractAddr, setBalanceContractABI, "setSimulatedBalance", []interface{}{aliceEVMAddr, big.NewInt(1000)}, s.ChainID, false)
	s.Require().NoError(err, "setSimulatedBalance should succeed")

	// Build calldata for checkMinBalance(address,uint256): selector + $sender (placeholder) + minAmount 1000
	abiDef := `[{"inputs":[{"internalType":"address","name":"account","type":"address"},{"internalType":"uint256","name":"minAmount","type":"uint256"}],"name":"checkMinBalance","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"}]`
	contractABI, err := abi.JSON(strings.NewReader(abiDef))
	s.Require().NoError(err)
	packed, err := contractABI.Pack("checkMinBalance", common.Address{}, big.NewInt(1000))
	s.Require().NoError(err)
	// Calldata template: selector (4 bytes = 8 hex) + "$sender" (7 chars) + minAmount 1000 (32 bytes = 64 hex)
	calldataHex := hex.EncodeToString(packed[:4]) + "$sender" + hex.EncodeToString(packed[36:])
	s.Require().Len(calldataHex, 8+7+64, "calldata template length")

	// Create collection with transfer approval that includes EVM query challenge (min bank balance)
	getFullRanges := func() []*tokenizationtypes.UintRange {
		return []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}}
	}
	evmChallenge := &tokenizationtypes.EVMQueryChallenge{
		ContractAddress:    contractAddr.Hex(),
		Calldata:           calldataHex,
		ExpectedResult:     "0x0000000000000000000000000000000000000000000000000000000000000001",
		ComparisonOperator: "eq",
		GasLimit:           sdkmath.NewUint(0), // 0 = use default in checker
		Uri:                "https://example.com/min-bank-balance",
		CustomData:         `{"description":"min bank balance check"}`,
	}

	createMsg := &tokenizationtypes.MsgUniversalUpdateCollection{
		Creator:               s.Alice.String(),
		CollectionId:          sdkmath.NewUint(0),
		ValidTokenIds:         []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)}},
		UpdateValidTokenIds:   true,
		CollectionPermissions: &tokenizationtypes.CollectionPermissions{},
	}
	msgServer := tokenizationkeeper.NewMsgServerImpl(s.TokenizationKeeper)
	resp, err := msgServer.UniversalUpdateCollection(s.Ctx, createMsg)
	s.Require().NoError(err)
	collectionId := resp.CollectionId

	// Mint approval
	mintApproval := &tokenizationtypes.CollectionApproval{
		ApprovalId: "mint", FromListId: tokenizationtypes.MintAddress, ToListId: "AllWithoutMint", InitiatedByListId: "AllWithoutMint",
		TransferTimes: getFullRanges(), TokenIds: getFullRanges(), OwnershipTimes: getFullRanges(),
		ApprovalCriteria: &tokenizationtypes.ApprovalCriteria{
			MaxNumTransfers:                &tokenizationtypes.MaxNumTransfers{OverallMaxNumTransfers: sdkmath.NewUint(1000), AmountTrackerId: "mint"},
			ApprovalAmounts:                &tokenizationtypes.ApprovalAmounts{PerFromAddressApprovalAmount: sdkmath.NewUint(1000), AmountTrackerId: "mint"},
			OverridesFromOutgoingApprovals: true, OverridesToIncomingApprovals: true,
		},
		Version: sdkmath.NewUint(0),
	}
	// Transfer approval with EVM query challenge (min bank balance)
	transferApproval := &tokenizationtypes.CollectionApproval{
		ApprovalId: "transfer", FromListId: "AllWithoutMint", ToListId: "All", InitiatedByListId: "AllWithoutMint",
		TransferTimes: getFullRanges(), TokenIds: getFullRanges(), OwnershipTimes: getFullRanges(),
		ApprovalCriteria: &tokenizationtypes.ApprovalCriteria{
			MaxNumTransfers:    &tokenizationtypes.MaxNumTransfers{OverallMaxNumTransfers: sdkmath.NewUint(1000), AmountTrackerId: "tr"},
			ApprovalAmounts:    &tokenizationtypes.ApprovalAmounts{PerFromAddressApprovalAmount: sdkmath.NewUint(1000), AmountTrackerId: "tr"},
			EvmQueryChallenges: []*tokenizationtypes.EVMQueryChallenge{evmChallenge},
		},
		Version: sdkmath.NewUint(0),
	}
	_, err = msgServer.UniversalUpdateCollection(s.Ctx, &tokenizationtypes.MsgUniversalUpdateCollection{
		Creator: s.Alice.String(), CollectionId: collectionId, UpdateCollectionApprovals: true,
		CollectionApprovals: []*tokenizationtypes.CollectionApproval{mintApproval, transferApproval},
	})
	s.Require().NoError(err)

	// User approvals
	for _, user := range []sdk.AccAddress{s.Alice, s.Bob} {
		_, _ = msgServer.SetOutgoingApproval(s.Ctx, &tokenizationtypes.MsgSetOutgoingApproval{
			Creator: user.String(), CollectionId: collectionId,
			Approval: &tokenizationtypes.UserOutgoingApproval{ApprovalId: "o", ToListId: "All", InitiatedByListId: "All", TransferTimes: getFullRanges(), TokenIds: getFullRanges(), OwnershipTimes: getFullRanges(), ApprovalCriteria: &tokenizationtypes.OutgoingApprovalCriteria{}, Version: sdkmath.NewUint(0)},
		})
		_, _ = msgServer.SetIncomingApproval(s.Ctx, &tokenizationtypes.MsgSetIncomingApproval{
			Creator: user.String(), CollectionId: collectionId,
			Approval: &tokenizationtypes.UserIncomingApproval{ApprovalId: "i", FromListId: "All", InitiatedByListId: "All", TransferTimes: getFullRanges(), TokenIds: getFullRanges(), OwnershipTimes: getFullRanges(), ApprovalCriteria: &tokenizationtypes.IncomingApprovalCriteria{}, Version: sdkmath.NewUint(0)},
		})
	}

	// Mint to Alice
	_, err = msgServer.TransferTokens(s.Ctx, &tokenizationtypes.MsgTransferTokens{
		Creator: s.Alice.String(), CollectionId: collectionId,
		Transfers: []*tokenizationtypes.Transfer{{
			From: "Mint", ToAddresses: []string{s.Alice.String()},
			Balances: []*tokenizationtypes.Balance{{Amount: sdkmath.NewUint(10), TokenIds: getFullRanges(), OwnershipTimes: getFullRanges()}},
		}},
	})
	s.Require().NoError(err)

	// Transfer Alice -> Bob; approval runs EVM query (Alice has >= 1000 ustake), should pass
	_, err = msgServer.TransferTokens(s.Ctx, &tokenizationtypes.MsgTransferTokens{
		Creator: s.Alice.String(), CollectionId: collectionId,
		Transfers: []*tokenizationtypes.Transfer{{
			From: s.Alice.String(), ToAddresses: []string{s.Bob.String()},
			Balances: []*tokenizationtypes.Balance{{Amount: sdkmath.NewUint(5), TokenIds: getFullRanges(), OwnershipTimes: getFullRanges()}},
		}},
	})
	s.Require().NoError(err, "transfer with min bank balance approval should succeed")
}

// TestE2E_InvariantFlow_MaxUniqueHolders checks maxUniqueHolders using collection stats in the invariant flow.
func (s *EVMQueryChallengesE2ESuite) TestE2E_InvariantFlow_MaxUniqueHolders() {
	bytecode, err := helpers.GetContractBytecodeByType(helpers.ContractTypeMaxUniqueHoldersChecker)
	if err != nil {
		s.T().Skipf("MaxUniqueHoldersChecker bytecode not found (run make compile-contracts): %v", err)
		return
	}

	contractAddr, _, err := helpers.DeployContract(s.Ctx, s.App.EVMKeeper, s.DeployerKey, bytecode, s.ChainID)
	s.Require().NoError(err)
	isContract, _ := helpers.VerifyContractDeployment(s.Ctx, s.App.EVMKeeper, contractAddr)
	s.Require().True(isContract, "MaxUniqueHoldersChecker should be deployed")

	// Calldata: checkMaxHolders(uint256 collectionId, uint256 maxAllowed) with $collectionId and 3
	abiDef := `[{"inputs":[{"internalType":"uint256","name":"collectionId","type":"uint256"},{"internalType":"uint256","name":"maxAllowed","type":"uint256"}],"name":"checkMaxHolders","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"}]`
	contractABI, err := abi.JSON(strings.NewReader(abiDef))
	s.Require().NoError(err)
	packed, err := contractABI.Pack("checkMaxHolders", big.NewInt(0), big.NewInt(3))
	s.Require().NoError(err)
	calldataHex := hex.EncodeToString(packed[:4]) + "$collectionId" + hex.EncodeToString(packed[36:])

	evmChallenge := &tokenizationtypes.EVMQueryChallenge{
		ContractAddress:    contractAddr.Hex(),
		Calldata:           calldataHex,
		ExpectedResult:     "0x0000000000000000000000000000000000000000000000000000000000000001",
		ComparisonOperator: "eq",
		GasLimit:           sdkmath.NewUint(0), // 0 = use default in checker
		Uri:                "https://example.com/max-unique-holders",
		CustomData:         `{"description":"max unique holders invariant"}`,
	}

	getFullRanges := func() []*tokenizationtypes.UintRange {
		return []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}}
	}

	// Create collection with invariants containing EVM query (max 3 holders)
	createMsg := &tokenizationtypes.MsgUniversalUpdateCollection{
		Creator: s.Alice.String(), CollectionId: sdkmath.NewUint(0),
		ValidTokenIds:       []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)}},
		UpdateValidTokenIds: true, CollectionPermissions: &tokenizationtypes.CollectionPermissions{},
		Invariants: &tokenizationtypes.InvariantsAddObject{
			EvmQueryChallenges: []*tokenizationtypes.EVMQueryChallenge{evmChallenge},
		},
	}
	msgServer := tokenizationkeeper.NewMsgServerImpl(s.TokenizationKeeper)
	resp, err := msgServer.UniversalUpdateCollection(s.Ctx, createMsg)
	s.Require().NoError(err)
	collectionId := resp.CollectionId

	mintApproval := &tokenizationtypes.CollectionApproval{
		ApprovalId: "mint", FromListId: tokenizationtypes.MintAddress, ToListId: "AllWithoutMint", InitiatedByListId: "AllWithoutMint",
		TransferTimes: getFullRanges(), TokenIds: getFullRanges(), OwnershipTimes: getFullRanges(),
		ApprovalCriteria: &tokenizationtypes.ApprovalCriteria{
			MaxNumTransfers:                &tokenizationtypes.MaxNumTransfers{OverallMaxNumTransfers: sdkmath.NewUint(1000), AmountTrackerId: "mint"},
			ApprovalAmounts:                &tokenizationtypes.ApprovalAmounts{PerFromAddressApprovalAmount: sdkmath.NewUint(1000), AmountTrackerId: "mint"},
			OverridesFromOutgoingApprovals: true, OverridesToIncomingApprovals: true,
		},
		Version: sdkmath.NewUint(0),
	}
	transferApproval := &tokenizationtypes.CollectionApproval{
		ApprovalId: "transfer", FromListId: "AllWithoutMint", ToListId: "All", InitiatedByListId: "AllWithoutMint",
		TransferTimes: getFullRanges(), TokenIds: getFullRanges(), OwnershipTimes: getFullRanges(),
		ApprovalCriteria: &tokenizationtypes.ApprovalCriteria{
			MaxNumTransfers: &tokenizationtypes.MaxNumTransfers{OverallMaxNumTransfers: sdkmath.NewUint(1000), AmountTrackerId: "tr"},
			ApprovalAmounts: &tokenizationtypes.ApprovalAmounts{PerFromAddressApprovalAmount: sdkmath.NewUint(1000), AmountTrackerId: "tr"},
		},
		Version: sdkmath.NewUint(0),
	}
	_, err = msgServer.UniversalUpdateCollection(s.Ctx, &tokenizationtypes.MsgUniversalUpdateCollection{
		Creator: s.Alice.String(), CollectionId: collectionId, UpdateCollectionApprovals: true,
		CollectionApprovals: []*tokenizationtypes.CollectionApproval{mintApproval, transferApproval},
	})
	s.Require().NoError(err)

	for _, user := range []sdk.AccAddress{s.Alice, s.Bob, s.Charlie, s.Dave} {
		_, _ = msgServer.SetOutgoingApproval(s.Ctx, &tokenizationtypes.MsgSetOutgoingApproval{
			Creator: user.String(), CollectionId: collectionId,
			Approval: &tokenizationtypes.UserOutgoingApproval{ApprovalId: "o", ToListId: "All", InitiatedByListId: "All", TransferTimes: getFullRanges(), TokenIds: getFullRanges(), OwnershipTimes: getFullRanges(), ApprovalCriteria: &tokenizationtypes.OutgoingApprovalCriteria{}, Version: sdkmath.NewUint(0)},
		})
		_, _ = msgServer.SetIncomingApproval(s.Ctx, &tokenizationtypes.MsgSetIncomingApproval{
			Creator: user.String(), CollectionId: collectionId,
			Approval: &tokenizationtypes.UserIncomingApproval{ApprovalId: "i", FromListId: "All", InitiatedByListId: "All", TransferTimes: getFullRanges(), TokenIds: getFullRanges(), OwnershipTimes: getFullRanges(), ApprovalCriteria: &tokenizationtypes.IncomingApprovalCriteria{}, Version: sdkmath.NewUint(0)},
		})
	}

	// Mint to Alice
	_, err = msgServer.TransferTokens(s.Ctx, &tokenizationtypes.MsgTransferTokens{
		Creator: s.Alice.String(), CollectionId: collectionId,
		Transfers: []*tokenizationtypes.Transfer{{
			From: "Mint", ToAddresses: []string{s.Alice.String()},
			Balances: []*tokenizationtypes.Balance{{Amount: sdkmath.NewUint(30), TokenIds: getFullRanges(), OwnershipTimes: getFullRanges()}},
		}},
	})
	s.Require().NoError(err)

	// Alice -> Bob (2 holders): invariant should pass
	_, err = msgServer.TransferTokens(s.Ctx, &tokenizationtypes.MsgTransferTokens{
		Creator: s.Alice.String(), CollectionId: collectionId,
		Transfers: []*tokenizationtypes.Transfer{{
			From: s.Alice.String(), ToAddresses: []string{s.Bob.String()},
			Balances: []*tokenizationtypes.Balance{{Amount: sdkmath.NewUint(10), TokenIds: getFullRanges(), OwnershipTimes: getFullRanges()}},
		}},
	})
	s.Require().NoError(err, "transfer with 2 holders should pass invariant")

	// Alice -> Charlie (3 holders): invariant should pass
	_, err = msgServer.TransferTokens(s.Ctx, &tokenizationtypes.MsgTransferTokens{
		Creator: s.Alice.String(), CollectionId: collectionId,
		Transfers: []*tokenizationtypes.Transfer{{
			From: s.Alice.String(), ToAddresses: []string{s.Charlie.String()},
			Balances: []*tokenizationtypes.Balance{{Amount: sdkmath.NewUint(10), TokenIds: getFullRanges(), OwnershipTimes: getFullRanges()}},
		}},
	})
	s.Require().NoError(err, "transfer with 3 holders should pass invariant")

	// Alice -> Dave (4th holder): invariant should fail (max 3)
	_, err = msgServer.TransferTokens(s.Ctx, &tokenizationtypes.MsgTransferTokens{
		Creator: s.Alice.String(), CollectionId: collectionId,
		Transfers: []*tokenizationtypes.Transfer{{
			From: s.Alice.String(), ToAddresses: []string{s.Dave.String()},
			Balances: []*tokenizationtypes.Balance{{Amount: sdkmath.NewUint(5), TokenIds: getFullRanges(), OwnershipTimes: getFullRanges()}},
		}},
	})
	s.Require().Error(err, "transfer creating 4th holder should fail invariant")
	s.Require().Contains(err.Error(), "invariant", "error should mention invariant")
}
