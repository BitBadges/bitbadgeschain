package tokenization

import (
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/suite"

	sdkmath "cosmossdk.io/math"

	tokenizationkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	keepertest "github.com/bitbadges/bitbadgeschain/x/tokenization/testutil/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ERC3643ComplianceTestSuite tests ERC-3643 standard compliance
// ERC-3643 is a minimal token standard requiring:
// - transfer(address to, uint256 amount) returns (bool)
// - balanceOf(address account) returns (uint256)
// - totalSupply() returns (uint256)
// - Transfer event emission
type ERC3643ComplianceTestSuite struct {
	suite.Suite
	TokenizationKeeper tokenizationkeeper.Keeper
	Ctx                sdk.Context
	Precompile         *Precompile

	// Test addresses
	AliceEVM common.Address
	BobEVM   common.Address
	Alice    sdk.AccAddress
	Bob      sdk.AccAddress

	// Test data
	CollectionId sdkmath.Uint
}

func TestERC3643ComplianceTestSuite(t *testing.T) {
	suite.Run(t, new(ERC3643ComplianceTestSuite))
}

// SetupTest initializes the test suite
func (suite *ERC3643ComplianceTestSuite) SetupTest() {
	keeper, ctx := keepertest.TokenizationKeeper(suite.T())
	suite.TokenizationKeeper = keeper
	suite.Ctx = ctx
	suite.Precompile = NewPrecompile(keeper)

	// Create test addresses
	suite.AliceEVM = common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0")
	suite.BobEVM = common.HexToAddress("0x8ba1f109551bD432803012645Hac136c22C9e7")

	suite.Alice = sdk.AccAddress(suite.AliceEVM.Bytes())
	suite.Bob = sdk.AccAddress(suite.BobEVM.Bytes())

	// Set up test collection
	suite.CollectionId = suite.createTestCollection()
}

// createTestCollection creates a test collection with balances
func (suite *ERC3643ComplianceTestSuite) createTestCollection() sdkmath.Uint {
	msgServer := tokenizationkeeper.NewMsgServerImpl(suite.TokenizationKeeper)

	// Create collection
	createMsg := &tokenizationtypes.MsgUniversalUpdateCollection{
		Creator:      suite.Alice.String(),
		CollectionId: sdkmath.NewUint(0),
		ValidTokenIds: []*tokenizationtypes.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
		},
		UpdateValidTokenIds:   true,
		CollectionPermissions: &tokenizationtypes.CollectionPermissions{},
	}

	resp, err := msgServer.UniversalUpdateCollection(suite.Ctx, createMsg)
	suite.Require().NoError(err)
	collectionId := resp.CollectionId

	// Set up mint approval
	getFullUintRanges := func() []*tokenizationtypes.UintRange {
		return []*tokenizationtypes.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
		}
	}

	mintApproval := &tokenizationtypes.CollectionApproval{
		ApprovalId:        "mint_approval",
		FromListId:        tokenizationtypes.MintAddress,
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		TransferTimes:     getFullUintRanges(),
		TokenIds:          getFullUintRanges(),
		OwnershipTimes:    getFullUintRanges(),
		ApprovalCriteria: &tokenizationtypes.ApprovalCriteria{
			MaxNumTransfers: &tokenizationtypes.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(1000),
				AmountTrackerId:        "mint-tracker",
			},
			ApprovalAmounts: &tokenizationtypes.ApprovalAmounts{
				PerFromAddressApprovalAmount: sdkmath.NewUint(1000),
				AmountTrackerId:              "mint-tracker",
			},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
		Version: sdkmath.NewUint(0),
	}

	// Also set up a collection approval for regular transfers (AllWithoutMint to All)
	transferApproval := &tokenizationtypes.CollectionApproval{
		ApprovalId:        "transfer_approval",
		FromListId:        "AllWithoutMint",
		ToListId:          "All",
		InitiatedByListId: "AllWithoutMint",
		TransferTimes:     getFullUintRanges(),
		TokenIds:          getFullUintRanges(),
		OwnershipTimes:    getFullUintRanges(),
		ApprovalCriteria: &tokenizationtypes.ApprovalCriteria{
			MaxNumTransfers: &tokenizationtypes.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(1000),
				AmountTrackerId:        "transfer-tracker",
			},
			ApprovalAmounts: &tokenizationtypes.ApprovalAmounts{
				PerFromAddressApprovalAmount: sdkmath.NewUint(1000),
				AmountTrackerId:              "transfer-tracker",
			},
		},
		Version: sdkmath.NewUint(0),
	}

	updateApprovalsMsg := &tokenizationtypes.MsgUniversalUpdateCollection{
		Creator:                   suite.Alice.String(),
		CollectionId:              collectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*tokenizationtypes.CollectionApproval{mintApproval, transferApproval},
	}
	_, err = msgServer.UniversalUpdateCollection(suite.Ctx, updateApprovalsMsg)
	suite.Require().NoError(err)

	// Set up user-level approvals for transfers
	// Alice needs outgoing approval
	outgoingApproval := &tokenizationtypes.UserOutgoingApproval{
		ApprovalId:        "alice_outgoing",
		ToListId:          "All",
		InitiatedByListId: "All",
		TransferTimes:     getFullUintRanges(),
		TokenIds:          getFullUintRanges(),
		OwnershipTimes:    getFullUintRanges(),
		ApprovalCriteria:  &tokenizationtypes.OutgoingApprovalCriteria{},
		Version:           sdkmath.NewUint(0),
	}
	setOutgoingMsg := &tokenizationtypes.MsgSetOutgoingApproval{
		Creator:      suite.Alice.String(),
		CollectionId: collectionId,
		Approval:     outgoingApproval,
	}
	_, err = msgServer.SetOutgoingApproval(suite.Ctx, setOutgoingMsg)
	suite.Require().NoError(err)

	// Bob needs incoming approval
	incomingApproval := &tokenizationtypes.UserIncomingApproval{
		ApprovalId:        "bob_incoming",
		FromListId:        "All",
		InitiatedByListId: "All",
		TransferTimes:     getFullUintRanges(),
		TokenIds:          getFullUintRanges(),
		OwnershipTimes:    getFullUintRanges(),
		ApprovalCriteria:  &tokenizationtypes.IncomingApprovalCriteria{},
		Version:           sdkmath.NewUint(0),
	}
	setIncomingMsg := &tokenizationtypes.MsgSetIncomingApproval{
		Creator:      suite.Bob.String(),
		CollectionId: collectionId,
		Approval:     incomingApproval,
	}
	_, err = msgServer.SetIncomingApproval(suite.Ctx, setIncomingMsg)
	suite.Require().NoError(err)

	// Mint tokens to Alice
	transferMsg := &tokenizationtypes.MsgTransferTokens{
		Creator:      suite.Alice.String(),
		CollectionId: collectionId,
		Transfers: []*tokenizationtypes.Transfer{
			{
				From:        tokenizationtypes.MintAddress,
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
	_, err = msgServer.TransferTokens(suite.Ctx, transferMsg)
	suite.Require().NoError(err)

	return collectionId
}

// TestERC3643_Transfer tests the transfer function compliance
// ERC-3643 requires: transfer(address to, uint256 amount) returns (bool)
func (suite *ERC3643ComplianceTestSuite) TestERC3643_Transfer() {
	method := suite.Precompile.ABI.Methods["transferTokens"]
	suite.Require().NotNil(method)

	// Test: Transfer should succeed with valid inputs
	suite.Run("transfer_succeeds", func() {
		args := []interface{}{
			suite.CollectionId.BigInt(),
			[]common.Address{suite.BobEVM},
			big.NewInt(100),
			[]struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			}{{Start: big.NewInt(1), End: big.NewInt(100)}},
			[]struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			}{{Start: big.NewInt(1), End: new(big.Int).SetUint64(math.MaxUint64)}},
		}

		// Create proper vm.Contract for testing
		precompileAddr := common.HexToAddress(TokenizationPrecompileAddress)
		valueUint256, _ := uint256.FromBig(big.NewInt(0))
		contract := vm.NewContract(suite.AliceEVM, precompileAddr, valueUint256, 1000000, nil)

		result, err := suite.Precompile.TransferTokens(suite.Ctx, &method, args, contract)
		suite.Require().NoError(err)
		suite.Require().NotNil(result)

		// Unpack result (should be bool)
		unpacked, err := method.Outputs.Unpack(result)
		suite.Require().NoError(err)
		suite.Require().Len(unpacked, 1)
		success, ok := unpacked[0].(bool)
		suite.Require().True(ok)
		suite.True(success, "Transfer should return true on success")
	})

	// Test: Transfer should fail with zero address
	suite.Run("transfer_fails_zero_address", func() {
		args := []interface{}{
			suite.CollectionId.BigInt(),
			[]common.Address{{}}, // Zero address
			big.NewInt(100),
			[]struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			}{{Start: big.NewInt(1), End: big.NewInt(100)}},
			[]struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			}{{Start: big.NewInt(1), End: new(big.Int).SetUint64(math.MaxUint64)}},
		}

		precompileAddr := common.HexToAddress(TokenizationPrecompileAddress)
		valueUint256, _ := uint256.FromBig(big.NewInt(0))
		contract := vm.NewContract(suite.AliceEVM, precompileAddr, valueUint256, 1000000, nil)
		_, err := suite.Precompile.TransferTokens(suite.Ctx, &method, args, contract)
		suite.Require().Error(err)
		suite.Contains(err.Error(), "cannot be zero address")
	})

	// Test: Transfer should fail with zero amount
	suite.Run("transfer_fails_zero_amount", func() {
		args := []interface{}{
			suite.CollectionId.BigInt(),
			[]common.Address{suite.BobEVM},
			big.NewInt(0), // Zero amount
			[]struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			}{{Start: big.NewInt(1), End: big.NewInt(100)}},
			[]struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			}{{Start: big.NewInt(1), End: new(big.Int).SetUint64(math.MaxUint64)}},
		}

		precompileAddr := common.HexToAddress(TokenizationPrecompileAddress)
		valueUint256, _ := uint256.FromBig(big.NewInt(0))
		contract := vm.NewContract(suite.AliceEVM, precompileAddr, valueUint256, 1000000, nil)
		_, err := suite.Precompile.TransferTokens(suite.Ctx, &method, args, contract)
		suite.Require().Error(err)
		suite.Contains(err.Error(), "must be greater than zero")
	})
}

// TestERC3643_BalanceOf tests the balanceOf function compliance
// ERC-3643 requires: balanceOf(address account) returns (uint256)
func (suite *ERC3643ComplianceTestSuite) TestERC3643_BalanceOf() {
	method := suite.Precompile.ABI.Methods["getBalanceAmount"]
	suite.Require().NotNil(method)

	// Test: balanceOf should return correct balance
	suite.Run("balanceOf_returns_correct_balance", func() {
		args := []interface{}{
			suite.CollectionId.BigInt(),
			suite.AliceEVM,
			[]struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			}{{Start: big.NewInt(1), End: big.NewInt(100)}},
			[]struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			}{{Start: big.NewInt(1), End: new(big.Int).SetUint64(math.MaxUint64)}},
		}

		result, err := suite.Precompile.GetBalanceAmount(suite.Ctx, &method, args)
		suite.Require().NoError(err)
		suite.Require().NotNil(result)

		// Unpack result (should be uint256)
		unpacked, err := method.Outputs.Unpack(result)
		suite.Require().NoError(err)
		suite.Require().Len(unpacked, 1)
		balance, ok := unpacked[0].(*big.Int)
		suite.Require().True(ok)
		suite.Equal(big.NewInt(1000), balance, "Balance should be 1000")
	})

	// Test: balanceOf should return 0 for account with no balance
	suite.Run("balanceOf_returns_zero_for_empty_account", func() {
		args := []interface{}{
			suite.CollectionId.BigInt(),
			suite.BobEVM, // Bob has no balance initially
			[]struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			}{{Start: big.NewInt(1), End: big.NewInt(100)}},
			[]struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			}{{Start: big.NewInt(1), End: new(big.Int).SetUint64(math.MaxUint64)}},
		}

		result, err := suite.Precompile.GetBalanceAmount(suite.Ctx, &method, args)
		suite.Require().NoError(err)
		suite.Require().NotNil(result)

		unpacked, err := method.Outputs.Unpack(result)
		suite.Require().NoError(err)
		suite.Require().Len(unpacked, 1)
		balance, ok := unpacked[0].(*big.Int)
		suite.Require().True(ok)
		suite.Require().NotNil(balance)
		suite.True(balance.Cmp(big.NewInt(0)) == 0, "Balance should be 0 for empty account, got %s", balance.String())
	})
}

// TestERC3643_TotalSupply tests the totalSupply function compliance
// ERC-3643 requires: totalSupply() returns (uint256)
func (suite *ERC3643ComplianceTestSuite) TestERC3643_TotalSupply() {
	method := suite.Precompile.ABI.Methods["getTotalSupply"]
	suite.Require().NotNil(method)

	// Test: totalSupply should return correct supply
	suite.Run("totalSupply_returns_correct_supply", func() {
		args := []interface{}{
			suite.CollectionId.BigInt(),
			[]struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			}{{Start: big.NewInt(1), End: big.NewInt(100)}},
			[]struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			}{{Start: big.NewInt(1), End: new(big.Int).SetUint64(math.MaxUint64)}},
		}

		result, err := suite.Precompile.GetTotalSupply(suite.Ctx, &method, args)
		suite.Require().NoError(err)
		suite.Require().NotNil(result)

		// Unpack result (should be uint256)
		unpacked, err := method.Outputs.Unpack(result)
		suite.Require().NoError(err)
		suite.Require().Len(unpacked, 1)
		supply, ok := unpacked[0].(*big.Int)
		suite.Require().True(ok)
		suite.Require().NotNil(supply)
		suite.True(supply.Cmp(big.NewInt(1000)) == 0, "Total supply should be 1000, got %s", supply.String())
	})
}

// TestERC3643_TransferEvent tests that Transfer events are emitted
// ERC-3643 requires: event Transfer(address indexed from, address indexed to, uint256 value)
func (suite *ERC3643ComplianceTestSuite) TestERC3643_TransferEvent() {
	// Note: Events are emitted as Cosmos SDK events, not Solidity events
	// The ERC-3643 wrapper contract should emit Solidity Transfer events
	// This test verifies the precompile emits Cosmos events that can be tracked

	method := suite.Precompile.ABI.Methods["transferTokens"]
	suite.Require().NotNil(method)

	args := []interface{}{
		suite.CollectionId.BigInt(),
		[]common.Address{suite.BobEVM},
		big.NewInt(50),
		[]struct {
			Start *big.Int `json:"start"`
			End   *big.Int `json:"end"`
		}{{Start: big.NewInt(1), End: big.NewInt(100)}},
		[]struct {
			Start *big.Int `json:"start"`
			End   *big.Int `json:"end"`
		}{{Start: big.NewInt(1), End: new(big.Int).SetUint64(math.MaxUint64)}},
	}

	precompileAddr := common.HexToAddress(TokenizationPrecompileAddress)
	valueUint256, _ := uint256.FromBig(big.NewInt(0))
	contract := vm.NewContract(suite.AliceEVM, precompileAddr, valueUint256, 1000000, nil)
	_, err := suite.Precompile.TransferTokens(suite.Ctx, &method, args, contract)
	suite.Require().NoError(err)

	// Verify event was emitted (check event manager)
	events := suite.Ctx.EventManager().Events()
	found := false
	for _, event := range events {
		if event.Type == "precompile_transfer_tokens" {
			found = true
			break
		}
	}
	suite.True(found, "Transfer event should be emitted")
}

// TestERC3643_Completeness tests that all required ERC-3643 functions are available
func (suite *ERC3643ComplianceTestSuite) TestERC3643_Completeness() {
	// Verify all required methods exist in ABI
	requiredMethods := []string{
		"transferTokens",   // Maps to transfer(address, uint256)
		"getBalanceAmount", // Maps to balanceOf(address)
		"getTotalSupply",   // Maps to totalSupply()
	}

	for _, methodName := range requiredMethods {
		method, found := suite.Precompile.ABI.Methods[methodName]
		suite.True(found, "Required method %s should exist", methodName)
		suite.NotNil(method, "Method %s should not be nil", methodName)
	}
}
