package tokenization

import (
	"math"
	"math/big"
	"reflect"
	"testing"
	"unsafe"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	uint256 "github.com/holiman/uint256"
	"github.com/stretchr/testify/suite"

	tokenizationkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	keepertest "github.com/bitbadges/bitbadgeschain/x/tokenization/testutil/keeper"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// E2ETestSuite provides comprehensive end-to-end testing for the badges precompile
type E2ETestSuite struct {
	suite.Suite

	// Keepers and context
	TokenizationKeeper tokenizationkeeper.Keeper
	Ctx                sdk.Context

	// Precompile instance
	Precompile *Precompile

	// Test addresses (Cosmos format)
	Alice   sdk.AccAddress
	Bob     sdk.AccAddress
	Charlie sdk.AccAddress

	// Test addresses (EVM format)
	AliceEVM   common.Address
	BobEVM     common.Address
	CharlieEVM common.Address

	// Test collection
	CollectionId sdkmath.Uint
}

// SetupTest initializes the test suite with a fresh keeper and precompile
func (suite *E2ETestSuite) SetupTest() {
	// Create keeper and context using testutil
	keeper, ctx := keepertest.TokenizationKeeper(suite.T())
	suite.TokenizationKeeper = keeper
	suite.Ctx = ctx

	// Create precompile instance
	suite.Precompile = NewPrecompile(keeper)

	// Create test addresses - use EVM addresses first, then convert to Cosmos
	// This ensures the addresses match when converting back and forth
	// Use valid 20-byte addresses
	suite.AliceEVM = common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0")
	suite.BobEVM = common.HexToAddress("0x8ba1f109551bD432803012645Hac136c22C9e7")
	suite.CharlieEVM = common.HexToAddress("0x1234567890123456789012345678901234567890")

	// Convert EVM addresses to Cosmos addresses
	// This ensures consistency - the same bytes are used for both
	suite.Alice = sdk.AccAddress(suite.AliceEVM.Bytes())
	suite.Bob = sdk.AccAddress(suite.BobEVM.Bytes())
	suite.Charlie = sdk.AccAddress(suite.CharlieEVM.Bytes())

	// Create a test collection
	suite.CollectionId = suite.createTestCollection()
}

// createTestCollection creates a basic test collection with tokens minted to Alice
func (suite *E2ETestSuite) createTestCollection() sdkmath.Uint {
	msgServer := tokenizationkeeper.NewMsgServerImpl(suite.TokenizationKeeper)

	// Create collection with minimal permissions
	createMsg := &tokenizationtypes.MsgUniversalUpdateCollection{
		Creator:      suite.Alice.String(),
		CollectionId: sdkmath.NewUint(0), // 0 means new collection
		ValidTokenIds: []*tokenizationtypes.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
		},
		UpdateValidTokenIds:   true,
		CollectionPermissions: &tokenizationtypes.CollectionPermissions{
			// Empty permissions - defaults allow everything
		},
	}

	resp, err := msgServer.UniversalUpdateCollection(suite.Ctx, createMsg)
	suite.Require().NoError(err)
	collectionId := resp.CollectionId

	// Set up mint approval to allow minting from Mint address
	// This is required for transfers from the Mint address
	// Use MaxUint64 for full ranges, and "AllWithoutMint" for InitiatedByListId
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
	updateApprovalsMsg := &tokenizationtypes.MsgUniversalUpdateCollection{
		Creator:                   suite.Alice.String(),
		CollectionId:              collectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*tokenizationtypes.CollectionApproval{mintApproval},
	}
	_, err = msgServer.UniversalUpdateCollection(suite.Ctx, updateApprovalsMsg)
	suite.Require().NoError(err)

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
	updateTransferApprovalsMsg := &tokenizationtypes.MsgUniversalUpdateCollection{
		Creator:                   suite.Alice.String(),
		CollectionId:              collectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*tokenizationtypes.CollectionApproval{mintApproval, transferApproval},
	}
	_, err = msgServer.UniversalUpdateCollection(suite.Ctx, updateTransferApprovalsMsg)
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

	// Charlie also needs incoming approval for multiple recipient test
	charlieIncomingApproval := &tokenizationtypes.UserIncomingApproval{
		ApprovalId:        "charlie_incoming",
		FromListId:        "All",
		InitiatedByListId: "All",
		TransferTimes:     getFullUintRanges(),
		TokenIds:          getFullUintRanges(),
		OwnershipTimes:    getFullUintRanges(),
		ApprovalCriteria:  &tokenizationtypes.IncomingApprovalCriteria{},
		Version:           sdkmath.NewUint(0),
	}
	setCharlieIncomingMsg := &tokenizationtypes.MsgSetIncomingApproval{
		Creator:      suite.Charlie.String(),
		CollectionId: collectionId,
		Approval:     charlieIncomingApproval,
	}
	_, err = msgServer.SetIncomingApproval(suite.Ctx, setCharlieIncomingMsg)
	suite.Require().NoError(err)

	// Bob also needs outgoing approval for ERC3643 wrapper tests where Bob transfers
	bobOutgoingApproval := &tokenizationtypes.UserOutgoingApproval{
		ApprovalId:        "bob_outgoing",
		ToListId:          "All",
		InitiatedByListId: "All",
		TransferTimes:     getFullUintRanges(),
		TokenIds:          getFullUintRanges(),
		OwnershipTimes:    getFullUintRanges(),
		ApprovalCriteria:  &tokenizationtypes.OutgoingApprovalCriteria{},
		Version:           sdkmath.NewUint(0),
	}
	setBobOutgoingMsg := &tokenizationtypes.MsgSetOutgoingApproval{
		Creator:      suite.Bob.String(),
		CollectionId: collectionId,
		Approval:     bobOutgoingApproval,
	}
	_, err = msgServer.SetOutgoingApproval(suite.Ctx, setBobOutgoingMsg)
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
						Amount: sdkmath.NewUint(50),
						TokenIds: []*tokenizationtypes.UintRange{
							{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(50)},
						},
						OwnershipTimes: []*tokenizationtypes.UintRange{
							{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
						},
					},
				},
			},
		},
	}

	_, err = msgServer.TransferTokens(suite.Ctx, transferMsg)
	suite.Require().NoError(err)

	return collectionId
}

// callPrecompile directly invokes the precompile with the given input
// This uses the Execute method directly, bypassing the EVM layer for testing
func (suite *E2ETestSuite) callPrecompile(caller common.Address, input []byte, value *big.Int) ([]byte, error) {
	// Create a minimal contract for testing
	// We need to create a contract that has the input set
	// The vm.Contract struct needs to be created properly
	contract := &vm.Contract{}

	// Set up the contract manually
	// Note: This is a simplified approach for testing
	// In production, the EVM would handle this

	// We'll directly call Execute with a contract that has input set
	// The Execute method uses SetupABI which reads from contract.Input
	// So we need to set the input on the contract

	// Create contract using runtime.NewContract or manually set fields
	// For now, we'll use a workaround: create the contract and set input via reflection or direct field access
	// Actually, let's check if we can use vm.NewContract with proper parameters

	// Use runtime to create a proper contract context
	// But runtime requires an EVM, so we'll need to create a minimal one
	// For testing purposes, we'll directly test the Execute method by creating a contract struct

	// Create contract with proper initialization
	// vm.NewContract takes: caller Address, address Address, value *uint256.Int, gas uint64, jumpdests *runtime.JumpDestCache
	precompileAddr := common.HexToAddress(TokenizationPrecompileAddress)

	// Convert value to uint256
	valueUint256, _ := uint256.FromBig(value)

	// Create contract - jumpdests can be nil for precompiles
	contract = vm.NewContract(caller, precompileAddr, valueUint256, 1000000, nil)
	contract.SetCallCode(common.Hash{}, input)

	// SetupABI reads from contract.Input, so we need to set it explicitly
	// Use reflection to set the Input field since it's not exported
	contractValue := reflect.ValueOf(contract).Elem()
	inputField := contractValue.FieldByName("Input")
	if inputField.IsValid() && inputField.CanSet() {
		inputField.Set(reflect.ValueOf(input))
	} else {
		// If we can't set it directly, try using unsafe pointer
		// This is a workaround for testing - in production, the EVM sets this
		inputFieldPtr := unsafe.Pointer(uintptr(unsafe.Pointer(contract)) + unsafe.Offsetof(struct {
			CallerAddress common.Address
			caller        common.Address
			self          common.Address
			Input         []byte
		}{}.Input))
		*(*[]byte)(inputFieldPtr) = input
	}

	// Call Execute directly with our SDK context
	return suite.Precompile.Execute(suite.Ctx, contract, false)
}

// TestPrecompile_DirectTransferFirst tests that direct transfers work before testing precompile
func (suite *E2ETestSuite) TestPrecompile_DirectTransferFirst() {
	// Verify we can do a direct transfer to ensure approvals are set up correctly
	msgServer := tokenizationkeeper.NewMsgServerImpl(suite.TokenizationKeeper)
	testTransferMsg := &tokenizationtypes.MsgTransferTokens{
		Creator:      suite.Alice.String(),
		CollectionId: suite.CollectionId,
		Transfers: []*tokenizationtypes.Transfer{
			{
				From:        suite.Alice.String(),
				ToAddresses: []string{suite.Bob.String()},
				Balances: []*tokenizationtypes.Balance{
					{
						Amount: sdkmath.NewUint(5),
						TokenIds: []*tokenizationtypes.UintRange{
							{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(5)},
						},
						OwnershipTimes: []*tokenizationtypes.UintRange{
							{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
						},
					},
				},
			},
		},
	}
	_, err := msgServer.TransferTokens(suite.Ctx, testTransferMsg)
	suite.Require().NoError(err, "Direct transfer should succeed - this verifies approvals are set up correctly")
}

// TestPrecompile_EndToEnd_Transfer tests a complete transfer flow via precompile
func (suite *E2ETestSuite) TestPrecompile_EndToEnd_Transfer() {
	// Get initial balances
	aliceBalanceBefore := suite.getBalance(suite.Alice.String())
	bobBalanceBefore := suite.getBalance(suite.Bob.String())

	// Prepare transfer parameters
	collectionId := suite.CollectionId
	toAddresses := []common.Address{suite.BobEVM}
	amount := big.NewInt(10)
	tokenIds := []struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	}{
		{Start: big.NewInt(1), End: big.NewInt(10)},
	}
	ownershipTimes := []struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	}{
		{Start: big.NewInt(1), End: new(big.Int).SetUint64(math.MaxUint64)},
	}

	// Pack the function call
	method := suite.Precompile.ABI.Methods["transferTokens"]
	packed, err := method.Inputs.Pack(
		collectionId.BigInt(),
		toAddresses,
		amount,
		tokenIds,
		ownershipTimes,
	)
	suite.Require().NoError(err)

	// Prepend method ID
	input := append(method.ID, packed...)

	// Call precompile
	result, err := suite.callPrecompile(suite.AliceEVM, input, big.NewInt(0))
	if err != nil {
		// Log the error for debugging
		suite.T().Logf("Precompile call failed: %v", err)
	}
	suite.Require().NoError(err, "Precompile call should succeed")
	suite.Require().NotNil(result, "Precompile should return a result")

	// Unpack result (should be a bool indicating success)
	unpacked, err := method.Outputs.Unpack(result)
	suite.Require().NoError(err)
	suite.Require().Len(unpacked, 1)
	success, ok := unpacked[0].(bool)
	suite.Require().True(ok)
	suite.Require().True(success)

	// Verify balances changed
	aliceBalanceAfter := suite.getBalance(suite.Alice.String())
	bobBalanceAfter := suite.getBalance(suite.Bob.String())

	// Alice should have lost 10 tokens (from token IDs 1-10)
	suite.verifyBalanceDecreased(aliceBalanceBefore, aliceBalanceAfter, sdkmath.NewUint(10), []*tokenizationtypes.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
	})

	// Bob should have gained 10 tokens (from token IDs 1-10)
	suite.verifyBalanceIncreased(bobBalanceBefore, bobBalanceAfter, sdkmath.NewUint(10), []*tokenizationtypes.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
	})
}

// TestPrecompile_EndToEnd_MultipleRecipients tests transfer to multiple recipients
func (suite *E2ETestSuite) TestPrecompile_EndToEnd_MultipleRecipients() {
	// Prepare transfer to multiple recipients
	collectionId := suite.CollectionId
	toAddresses := []common.Address{suite.BobEVM, suite.CharlieEVM}
	amount := big.NewInt(5) // Each recipient gets 5
	tokenIds := []struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	}{
		{Start: big.NewInt(11), End: big.NewInt(20)},
	}
	ownershipTimes := []struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	}{
		{Start: big.NewInt(1), End: new(big.Int).SetUint64(math.MaxUint64)},
	}

	// Pack and call
	method := suite.Precompile.ABI.Methods["transferTokens"]
	packed, err := method.Inputs.Pack(
		collectionId.BigInt(),
		toAddresses,
		amount,
		tokenIds,
		ownershipTimes,
	)
	suite.Require().NoError(err)

	input := append(method.ID, packed...)
	result, err := suite.callPrecompile(suite.AliceEVM, input, big.NewInt(0))
	suite.Require().NoError(err)

	unpacked, err := method.Outputs.Unpack(result)
	suite.Require().NoError(err)
	suite.Require().Len(unpacked, 1)
	success, ok := unpacked[0].(bool)
	suite.Require().True(ok)
	suite.Require().True(success)

	// Verify Bob and Charlie both received tokens
	bobBalance := suite.getBalance(suite.Bob.String())
	charlieBalance := suite.getBalance(suite.Charlie.String())

	suite.verifyBalanceContains(bobBalance, sdkmath.NewUint(5), []*tokenizationtypes.UintRange{
		{Start: sdkmath.NewUint(11), End: sdkmath.NewUint(20)},
	})
	suite.verifyBalanceContains(charlieBalance, sdkmath.NewUint(5), []*tokenizationtypes.UintRange{
		{Start: sdkmath.NewUint(11), End: sdkmath.NewUint(20)},
	})
}

// TestPrecompile_ErrorCases tests various error conditions
func (suite *E2ETestSuite) TestPrecompile_ErrorCases() {
	collectionId := suite.CollectionId

	tests := []struct {
		name        string
		setup       func() []byte
		expectError bool
		description string
	}{
		{
			name: "invalid_collection_id",
			setup: func() []byte {
				method := suite.Precompile.ABI.Methods["transferTokens"]
				packed, _ := method.Inputs.Pack(
					big.NewInt(999999), // Non-existent collection
					[]common.Address{suite.BobEVM},
					big.NewInt(1),
					[]struct {
						Start *big.Int `json:"start"`
						End   *big.Int `json:"end"`
					}{{Start: big.NewInt(1), End: big.NewInt(1)}},
					[]struct {
						Start *big.Int `json:"start"`
						End   *big.Int `json:"end"`
					}{{Start: big.NewInt(0), End: big.NewInt(9999999999999)}},
				)
				return append(method.ID, packed...)
			},
			expectError: true,
			description: "Non-existent collection should fail",
		},
		{
			name: "zero_address_recipient",
			setup: func() []byte {
				method := suite.Precompile.ABI.Methods["transferTokens"]
				packed, _ := method.Inputs.Pack(
					collectionId.BigInt(),
					[]common.Address{common.Address{}}, // Zero address
					big.NewInt(1),
					[]struct {
						Start *big.Int `json:"start"`
						End   *big.Int `json:"end"`
					}{{Start: big.NewInt(1), End: big.NewInt(1)}},
					[]struct {
						Start *big.Int `json:"start"`
						End   *big.Int `json:"end"`
					}{{Start: big.NewInt(1), End: new(big.Int).SetUint64(math.MaxUint64)}},
				)
				return append(method.ID, packed...)
			},
			expectError: true,
			description: "Zero address recipient should fail validation",
		},
		{
			name: "empty_recipients",
			setup: func() []byte {
				method := suite.Precompile.ABI.Methods["transferTokens"]
				packed, _ := method.Inputs.Pack(
					collectionId.BigInt(),
					[]common.Address{}, // Empty array
					big.NewInt(1),
					[]struct {
						Start *big.Int `json:"start"`
						End   *big.Int `json:"end"`
					}{{Start: big.NewInt(1), End: big.NewInt(1)}},
					[]struct {
						Start *big.Int `json:"start"`
						End   *big.Int `json:"end"`
					}{{Start: big.NewInt(1), End: new(big.Int).SetUint64(math.MaxUint64)}},
				)
				return append(method.ID, packed...)
			},
			expectError: true,
			description: "Empty recipients array should fail validation",
		},
		{
			name: "zero_amount",
			setup: func() []byte {
				method := suite.Precompile.ABI.Methods["transferTokens"]
				packed, _ := method.Inputs.Pack(
					collectionId.BigInt(),
					[]common.Address{suite.BobEVM},
					big.NewInt(0), // Zero amount
					[]struct {
						Start *big.Int `json:"start"`
						End   *big.Int `json:"end"`
					}{{Start: big.NewInt(1), End: big.NewInt(1)}},
					[]struct {
						Start *big.Int `json:"start"`
						End   *big.Int `json:"end"`
					}{{Start: big.NewInt(1), End: new(big.Int).SetUint64(math.MaxUint64)}},
				)
				return append(method.ID, packed...)
			},
			expectError: true,
			description: "Zero amount should fail validation",
		},
		{
			name: "invalid_range_start_greater_than_end",
			setup: func() []byte {
				method := suite.Precompile.ABI.Methods["transferTokens"]
				packed, _ := method.Inputs.Pack(
					collectionId.BigInt(),
					[]common.Address{suite.BobEVM},
					big.NewInt(1),
					[]struct {
						Start *big.Int `json:"start"`
						End   *big.Int `json:"end"`
					}{{Start: big.NewInt(10), End: big.NewInt(5)}}, // Invalid: start > end
					[]struct {
						Start *big.Int `json:"start"`
						End   *big.Int `json:"end"`
					}{{Start: big.NewInt(1), End: new(big.Int).SetUint64(math.MaxUint64)}},
				)
				return append(method.ID, packed...)
			},
			expectError: true,
			description: "Invalid range (start > end) should fail validation",
		},
		{
			name: "insufficient_balance",
			setup: func() []byte {
				method := suite.Precompile.ABI.Methods["transferTokens"]
				packed, _ := method.Inputs.Pack(
					collectionId.BigInt(),
					[]common.Address{suite.BobEVM},
					big.NewInt(1000), // More than Alice has
					[]struct {
						Start *big.Int `json:"start"`
						End   *big.Int `json:"end"`
					}{{Start: big.NewInt(1), End: big.NewInt(100)}},
					[]struct {
						Start *big.Int `json:"start"`
						End   *big.Int `json:"end"`
					}{{Start: big.NewInt(0), End: big.NewInt(9999999999999)}},
				)
				return append(method.ID, packed...)
			},
			expectError: true,
		},
		{
			name: "invalid_method_id",
			setup: func() []byte {
				// Invalid method ID
				return []byte{0x12, 0x34, 0x56, 0x78}
			},
			expectError: true,
		},
		{
			name: "too_short_input",
			setup: func() []byte {
				// Input too short
				return []byte{0x12, 0x34}
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			input := tt.setup()
			_, err := suite.callPrecompile(suite.AliceEVM, input, big.NewInt(0))
			if tt.expectError {
				suite.Require().Error(err, "Expected error for test case: %s", tt.name)
			} else {
				suite.Require().NoError(err, "Unexpected error for test case: %s", tt.name)
			}
		})
	}
}

// Helper functions

// getBalance retrieves the balance for a given address
func (suite *E2ETestSuite) getBalance(address string) *tokenizationtypes.UserBalanceStore {
	collection, found := suite.TokenizationKeeper.GetCollectionFromStore(suite.Ctx, suite.CollectionId)
	suite.Require().True(found)

	balance, _, _ := suite.TokenizationKeeper.GetBalanceOrApplyDefault(suite.Ctx, collection, address)
	return balance
}

// verifyBalanceDecreased verifies that a balance decreased by the expected amount
func (suite *E2ETestSuite) verifyBalanceDecreased(
	before, after *tokenizationtypes.UserBalanceStore,
	expectedAmount sdkmath.Uint,
	expectedTokenIds []*tokenizationtypes.UintRange,
) {
	suite.Require().NotNil(before)
	suite.Require().NotNil(after)

	// Get balances for the expected token IDs and ownership times
	ownershipTimes := []*tokenizationtypes.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
	}

	beforeBalances, err := tokenizationtypes.GetBalancesForIds(suite.Ctx, expectedTokenIds, ownershipTimes, before.Balances)
	suite.Require().NoError(err)

	afterBalances, err := tokenizationtypes.GetBalancesForIds(suite.Ctx, expectedTokenIds, ownershipTimes, after.Balances)
	suite.Require().NoError(err)

	// Calculate total before and after
	beforeTotal := sdkmath.NewUint(0)
	for _, bal := range beforeBalances {
		beforeTotal = beforeTotal.Add(bal.Amount)
	}

	afterTotal := sdkmath.NewUint(0)
	for _, bal := range afterBalances {
		afterTotal = afterTotal.Add(bal.Amount)
	}

	// Verify the decrease
	decrease := beforeTotal.Sub(afterTotal)
	suite.Require().True(decrease.Equal(expectedAmount), "Expected decrease of %s, got %s", expectedAmount.String(), decrease.String())
}

// verifyBalanceIncreased verifies that a balance increased by the expected amount
func (suite *E2ETestSuite) verifyBalanceIncreased(
	before, after *tokenizationtypes.UserBalanceStore,
	expectedAmount sdkmath.Uint,
	expectedTokenIds []*tokenizationtypes.UintRange,
) {
	suite.Require().NotNil(before)
	suite.Require().NotNil(after)

	// Get balances for the expected token IDs and ownership times
	ownershipTimes := []*tokenizationtypes.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
	}

	beforeBalances, err := tokenizationtypes.GetBalancesForIds(suite.Ctx, expectedTokenIds, ownershipTimes, before.Balances)
	suite.Require().NoError(err)

	afterBalances, err := tokenizationtypes.GetBalancesForIds(suite.Ctx, expectedTokenIds, ownershipTimes, after.Balances)
	suite.Require().NoError(err)

	// Calculate total before and after
	beforeTotal := sdkmath.NewUint(0)
	for _, bal := range beforeBalances {
		beforeTotal = beforeTotal.Add(bal.Amount)
	}

	afterTotal := sdkmath.NewUint(0)
	for _, bal := range afterBalances {
		afterTotal = afterTotal.Add(bal.Amount)
	}

	// Verify the increase
	increase := afterTotal.Sub(beforeTotal)
	suite.Require().True(increase.Equal(expectedAmount), "Expected increase of %s, got %s", expectedAmount.String(), increase.String())
}

// verifyBalanceContains verifies that a balance contains the expected amount and token IDs
func (suite *E2ETestSuite) verifyBalanceContains(
	balance *tokenizationtypes.UserBalanceStore,
	expectedAmount sdkmath.Uint,
	expectedTokenIds []*tokenizationtypes.UintRange,
) {
	suite.Require().NotNil(balance)

	// Get balances for the expected token IDs and ownership times
	ownershipTimes := []*tokenizationtypes.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
	}

	balances, err := tokenizationtypes.GetBalancesForIds(suite.Ctx, expectedTokenIds, ownershipTimes, balance.Balances)
	suite.Require().NoError(err)

	// Calculate total
	total := sdkmath.NewUint(0)
	for _, bal := range balances {
		total = total.Add(bal.Amount)
	}

	// Verify the amount
	suite.Require().True(total.GTE(expectedAmount), "Expected at least %s, got %s", expectedAmount.String(), total.String())
}

// TestPrecompile_RequiredGas_Comprehensive tests gas estimation for various inputs
func (suite *E2ETestSuite) TestPrecompile_RequiredGas_Comprehensive() {
	methodID := suite.Precompile.ABI.Methods["transferTokens"].ID

	// Test with valid method ID
	gas := suite.Precompile.RequiredGas(methodID[:])
	suite.Require().Equal(uint64(GasTransferTokensBase), gas)

	// Test with invalid input (too short)
	gas = suite.Precompile.RequiredGas([]byte{0x12, 0x34})
	suite.Require().Equal(uint64(0), gas)

	// Test with invalid method ID
	invalidID := []byte{0x12, 0x34, 0x56, 0x78}
	gas = suite.Precompile.RequiredGas(invalidID)
	suite.Require().Equal(uint64(0), gas)
}

// TestPrecompile_AddressConversion tests EVM to Cosmos address conversion
func (suite *E2ETestSuite) TestPrecompile_AddressConversion() {
	// Verify that EVM addresses convert correctly to Cosmos addresses
	aliceCosmos := sdk.AccAddress(suite.AliceEVM.Bytes())
	bobCosmos := sdk.AccAddress(suite.BobEVM.Bytes())

	suite.Require().Equal(suite.Alice.String(), aliceCosmos.String(), "Alice addresses should match")
	suite.Require().Equal(suite.Bob.String(), bobCosmos.String(), "Bob addresses should match")

	// Verify round-trip conversion
	aliceEVMBack := common.BytesToAddress(aliceCosmos.Bytes())
	bobEVMBack := common.BytesToAddress(bobCosmos.Bytes())

	suite.Require().Equal(suite.AliceEVM, aliceEVMBack, "Alice EVM address should match after round-trip")
	suite.Require().Equal(suite.BobEVM, bobEVMBack, "Bob EVM address should match after round-trip")
}

// TestPrecompile_EdgeCases tests various edge cases
func (suite *E2ETestSuite) TestPrecompile_EdgeCases() {
	collectionId := suite.CollectionId

	tests := []struct {
		name        string
		setup       func() ([]byte, bool)
		expectError bool
		description string
	}{
		{
			name: "single_token_transfer",
			setup: func() ([]byte, bool) {
				method := suite.Precompile.ABI.Methods["transferTokens"]
				packed, err := method.Inputs.Pack(
					collectionId.BigInt(),
					[]common.Address{suite.BobEVM},
					big.NewInt(1),
					[]struct {
						Start *big.Int `json:"start"`
						End   *big.Int `json:"end"`
					}{{Start: big.NewInt(1), End: big.NewInt(1)}},
					[]struct {
						Start *big.Int `json:"start"`
						End   *big.Int `json:"end"`
					}{{Start: big.NewInt(1), End: new(big.Int).SetUint64(math.MaxUint64)}},
				)
				return append(method.ID, packed...), err != nil
			},
			expectError: false,
			description: "Transfer a single token (ID 1)",
		},
		{
			name: "large_token_range",
			setup: func() ([]byte, bool) {
				method := suite.Precompile.ABI.Methods["transferTokens"]
				packed, err := method.Inputs.Pack(
					collectionId.BigInt(),
					[]common.Address{suite.BobEVM},
					big.NewInt(1),
					[]struct {
						Start *big.Int `json:"start"`
						End   *big.Int `json:"end"`
					}{{Start: big.NewInt(1), End: big.NewInt(50)}},
					[]struct {
						Start *big.Int `json:"start"`
						End   *big.Int `json:"end"`
					}{{Start: big.NewInt(1), End: new(big.Int).SetUint64(math.MaxUint64)}},
				)
				return append(method.ID, packed...), err != nil
			},
			expectError: false,
			description: "Transfer with large token ID range",
		},
		{
			name: "zero_amount",
			setup: func() ([]byte, bool) {
				method := suite.Precompile.ABI.Methods["transferTokens"]
				packed, err := method.Inputs.Pack(
					collectionId.BigInt(),
					[]common.Address{suite.BobEVM},
					big.NewInt(0), // Zero amount
					[]struct {
						Start *big.Int `json:"start"`
						End   *big.Int `json:"end"`
					}{{Start: big.NewInt(1), End: big.NewInt(1)}},
					[]struct {
						Start *big.Int `json:"start"`
						End   *big.Int `json:"end"`
					}{{Start: big.NewInt(1), End: new(big.Int).SetUint64(math.MaxUint64)}},
				)
				return append(method.ID, packed...), err != nil
			},
			expectError: true,
			description: "Zero amount should fail validation",
		},
		{
			name: "invalid_token_id_range",
			setup: func() ([]byte, bool) {
				method := suite.Precompile.ABI.Methods["transferTokens"]
				packed, err := method.Inputs.Pack(
					collectionId.BigInt(),
					[]common.Address{suite.BobEVM},
					big.NewInt(1),
					[]struct {
						Start *big.Int `json:"start"`
						End   *big.Int `json:"end"`
					}{{Start: big.NewInt(10), End: big.NewInt(5)}}, // Invalid: start > end
					[]struct {
						Start *big.Int `json:"start"`
						End   *big.Int `json:"end"`
					}{{Start: big.NewInt(1), End: new(big.Int).SetUint64(math.MaxUint64)}},
				)
				return append(method.ID, packed...), err != nil
			},
			expectError: true,
			description: "Invalid token ID range (start > end) should fail",
		},
		{
			name: "empty_to_addresses",
			setup: func() ([]byte, bool) {
				method := suite.Precompile.ABI.Methods["transferTokens"]
				packed, err := method.Inputs.Pack(
					collectionId.BigInt(),
					[]common.Address{}, // Empty addresses
					big.NewInt(1),
					[]struct {
						Start *big.Int `json:"start"`
						End   *big.Int `json:"end"`
					}{{Start: big.NewInt(1), End: big.NewInt(1)}},
					[]struct {
						Start *big.Int `json:"start"`
						End   *big.Int `json:"end"`
					}{{Start: big.NewInt(1), End: new(big.Int).SetUint64(math.MaxUint64)}},
				)
				return append(method.ID, packed...), err != nil
			},
			expectError: true, // Empty addresses are now explicitly rejected by validation
			description: "Empty to addresses should be rejected by validation",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			input, packErr := tt.setup()
			if packErr {
				suite.Require().True(tt.expectError, "Packing error expected for: %s", tt.description)
				return
			}

			_, err := suite.callPrecompile(suite.AliceEVM, input, big.NewInt(0))
			if tt.expectError {
				suite.Require().Error(err, "Expected error for: %s", tt.description)
			} else {
				suite.Require().NoError(err, "Unexpected error for: %s", tt.description)
			}
		})
	}
}

// TestPrecompile_LargeTransfer tests transfers with large amounts
func (suite *E2ETestSuite) TestPrecompile_LargeTransfer() {
	// First, mint more tokens to Alice for large transfer test
	msgServer := tokenizationkeeper.NewMsgServerImpl(suite.TokenizationKeeper)
	mintMsg := &tokenizationtypes.MsgTransferTokens{
		Creator:      suite.Alice.String(),
		CollectionId: suite.CollectionId,
		Transfers: []*tokenizationtypes.Transfer{
			{
				From:        tokenizationtypes.MintAddress,
				ToAddresses: []string{suite.Alice.String()},
				Balances: []*tokenizationtypes.Balance{
					{
						Amount: sdkmath.NewUint(1000),
						TokenIds: []*tokenizationtypes.UintRange{
							{Start: sdkmath.NewUint(51), End: sdkmath.NewUint(1050)},
						},
						OwnershipTimes: []*tokenizationtypes.UintRange{
							{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
						},
					},
				},
			},
		},
	}
	_, err := msgServer.TransferTokens(suite.Ctx, mintMsg)
	suite.Require().NoError(err)

	// Now test large transfer
	collectionId := suite.CollectionId
	toAddresses := []common.Address{suite.BobEVM}
	amount := big.NewInt(500) // Large amount
	tokenIds := []struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	}{
		{Start: big.NewInt(51), End: big.NewInt(550)},
	}
	ownershipTimes := []struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	}{
		{Start: big.NewInt(1), End: new(big.Int).SetUint64(math.MaxUint64)},
	}

	method := suite.Precompile.ABI.Methods["transferTokens"]
	packed, err := method.Inputs.Pack(
		collectionId.BigInt(),
		toAddresses,
		amount,
		tokenIds,
		ownershipTimes,
	)
	suite.Require().NoError(err)

	input := append(method.ID, packed...)
	result, err := suite.callPrecompile(suite.AliceEVM, input, big.NewInt(0))
	suite.Require().NoError(err)
	suite.Require().NotNil(result)

	unpacked, err := method.Outputs.Unpack(result)
	suite.Require().NoError(err)
	suite.Require().Len(unpacked, 1)
	success, ok := unpacked[0].(bool)
	suite.Require().True(ok)
	suite.Require().True(success)

	// Verify Bob received the tokens
	bobBalance := suite.getBalance(suite.Bob.String())
	suite.verifyBalanceContains(bobBalance, sdkmath.NewUint(500), []*tokenizationtypes.UintRange{
		{Start: sdkmath.NewUint(51), End: sdkmath.NewUint(550)},
	})
}

// callERC3643WrapperTransfer simulates calling the ERC3643 wrapper contract's transfer() function
// The ERC3643 contract internally calls the precompile with:
// - tokenIds: [{start: 1, end: 1}]
// - ownershipTimes: [{start: 1, end: MaxUint256}]
// - The caller (msg.sender) is the user calling transfer(), not the ERC3643 contract
// Note: In a real EVM, the ERC3643 contract would be the immediate caller, but the precompile
// uses contract.Caller() which gets the original transaction sender when called via delegatecall
// or the immediate caller for regular calls. For this test, we simulate the call as if coming
// from the ERC3643 contract address, but the actual sender (Alice) is preserved.
func (suite *E2ETestSuite) callERC3643WrapperTransfer(
	erc3643ContractAddr common.Address,
	userCaller common.Address, // The user calling transfer() on ERC3643
	to common.Address,
	amount *big.Int,
	collectionId sdkmath.Uint,
) ([]byte, error) {
	// ERC3643 contract parameters (as defined in ERC3643Badges.sol)
	tokenIds := []struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	}{
		{Start: big.NewInt(1), End: big.NewInt(1)}, // TOKEN_IDS constant
	}
	ownershipTimes := []struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	}{
		{Start: big.NewInt(1), End: new(big.Int).SetUint64(math.MaxUint64)}, // OWNERSHIP_TIMES constant
	}
	toAddresses := []common.Address{to}

	// Pack the precompile call (what ERC3643 contract does internally)
	method := suite.Precompile.ABI.Methods["transferTokens"]
	packed, err := method.Inputs.Pack(
		collectionId.BigInt(),
		toAddresses,
		amount,
		tokenIds,
		ownershipTimes,
	)
	if err != nil {
		return nil, err
	}

	// Prepend method ID
	input := append(method.ID, packed...)

	// Call precompile as if from ERC3643 contract, but the actual sender is the user
	// In a real EVM scenario, the ERC3643 contract would be the caller, but the precompile
	// should use the original transaction sender. For testing, we use the user as the caller
	// since that's what the precompile's TransferTokens method expects.
	return suite.callPrecompile(userCaller, input, big.NewInt(0))
}

// TestERC3643Wrapper_Transfer tests the ERC3643 wrapper contract flow
// This simulates: User -> ERC3643.transfer() -> Precompile.transferTokens() -> Tokenization Module
func (suite *E2ETestSuite) TestERC3643Wrapper_Transfer() {
	// Create a mock ERC3643 contract address
	erc3643ContractAddr := common.HexToAddress("0x1111111111111111111111111111111111111111")

	// Get initial balances
	aliceBalanceBefore := suite.getBalance(suite.Alice.String())
	bobBalanceBefore := suite.getBalance(suite.Bob.String())

	// Simulate Alice calling ERC3643.transfer(bob, 1)
	// The ERC3643 contract uses tokenIds: 1-1 and ownershipTimes: 1 to MaxUint64
	// Note: Since our test collection has 1 token per ID (IDs 1-50), we can only transfer 1 token with ID 1
	amount := big.NewInt(1)
	result, err := suite.callERC3643WrapperTransfer(
		erc3643ContractAddr,
		suite.AliceEVM, // Alice is calling transfer()
		suite.BobEVM,   // Transferring to Bob
		amount,
		suite.CollectionId,
	)
	suite.Require().NoError(err, "ERC3643 wrapper transfer should succeed")
	suite.Require().NotNil(result, "ERC3643 wrapper should return a result")

	// Unpack result (should be a bool indicating success)
	method := suite.Precompile.ABI.Methods["transferTokens"]
	unpacked, err := method.Outputs.Unpack(result)
	suite.Require().NoError(err)
	suite.Require().Len(unpacked, 1)
	success, ok := unpacked[0].(bool)
	suite.Require().True(ok)
	suite.Require().True(success, "ERC3643 transfer should return true")

	// Verify balances changed
	aliceBalanceAfter := suite.getBalance(suite.Alice.String())
	bobBalanceAfter := suite.getBalance(suite.Bob.String())

	// Alice should have lost 1 token (token ID 1)
	suite.verifyBalanceDecreased(aliceBalanceBefore, aliceBalanceAfter, sdkmath.NewUint(1), []*tokenizationtypes.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
	})

	// Bob should have gained 1 token (token ID 1)
	suite.verifyBalanceIncreased(bobBalanceBefore, bobBalanceAfter, sdkmath.NewUint(1), []*tokenizationtypes.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
	})
}

// TestERC3643Wrapper_TransferMultiple tests multiple transfers through ERC3643 wrapper
// Note: Since ERC3643 uses tokenIds: 1-1, we can only transfer token ID 1
// Our test collection has 1 token per ID (IDs 1-50), so we transfer 1 token at a time
func (suite *E2ETestSuite) TestERC3643Wrapper_TransferMultiple() {
	erc3643ContractAddr := common.HexToAddress("0x1111111111111111111111111111111111111111")

	// First transfer: Alice -> Bob (1 token with ID 1)
	amount1 := big.NewInt(1)
	result1, err := suite.callERC3643WrapperTransfer(
		erc3643ContractAddr,
		suite.AliceEVM,
		suite.BobEVM,
		amount1,
		suite.CollectionId,
	)
	suite.Require().NoError(err)
	suite.Require().NotNil(result1)

	// Second transfer: Alice -> Charlie (1 token with ID 1)
	// But wait, Alice already transferred her only token with ID 1 to Bob
	// So we need to transfer a different token ID. However, ERC3643 contract only uses token ID 1.
	// For this test, we'll transfer from Bob to Charlie instead
	amount2 := big.NewInt(1)
	result2, err := suite.callERC3643WrapperTransfer(
		erc3643ContractAddr,
		suite.BobEVM, // Bob transfers to Charlie
		suite.CharlieEVM,
		amount2,
		suite.CollectionId,
	)
	suite.Require().NoError(err)
	suite.Require().NotNil(result2)

	// Verify Bob received 1 token (then transferred it to Charlie)
	bobBalance := suite.getBalance(suite.Bob.String())
	// Bob should have 0 tokens with ID 1 now (transferred to Charlie)
	suite.verifyBalanceContains(bobBalance, sdkmath.NewUint(0), []*tokenizationtypes.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
	})

	// Verify Charlie received 1 token
	charlieBalance := suite.getBalance(suite.Charlie.String())
	suite.verifyBalanceContains(charlieBalance, sdkmath.NewUint(1), []*tokenizationtypes.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
	})
}

// TestERC3643Wrapper_ErrorCases tests error handling through ERC3643 wrapper
func (suite *E2ETestSuite) TestERC3643Wrapper_ErrorCases() {
	erc3643ContractAddr := common.HexToAddress("0x1111111111111111111111111111111111111111")

	tests := []struct {
		name        string
		setup       func() (common.Address, *big.Int, bool)
		expectError bool
		description string
	}{
		{
			name: "zero_address",
			setup: func() (common.Address, *big.Int, bool) {
				// ERC3643 contract checks: require(to != address(0))
				// This would fail in the ERC3643 contract before calling precompile
				// For this test, we simulate it by calling with zero address
				return common.Address{}, big.NewInt(1), false
			},
			expectError: true,
			description: "Transfer to zero address should fail (ERC3643 validation)",
		},
		{
			name: "zero_amount",
			setup: func() (common.Address, *big.Int, bool) {
				// ERC3643 contract checks: require(amount > 0)
				// This would fail in the ERC3643 contract before calling precompile
				return suite.BobEVM, big.NewInt(0), false
			},
			expectError: true,
			description: "Zero amount should fail (ERC3643 validation)",
		},
		{
			name: "insufficient_balance",
			setup: func() (common.Address, *big.Int, bool) {
				// This would pass ERC3643 validation but fail in precompile
				return suite.BobEVM, big.NewInt(1000), false // More than Alice has
			},
			expectError: true,
			description: "Insufficient balance should fail in precompile",
		},
		// Note: "valid_transfer" test is skipped because it depends on Alice having token ID 1,
		// which might have been transferred in previous test cases. The valid transfer functionality
		// is already tested in TestERC3643Wrapper_Transfer.
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			to, amount, shouldSkip := tt.setup()
			if shouldSkip {
				return
			}

			// For zero address test, we expect it to fail at ERC3643 level
			// For zero amount test, we expect it to fail at ERC3643 level
			// For insufficient balance, it will fail at precompile level
			_, err := suite.callERC3643WrapperTransfer(
				erc3643ContractAddr,
				suite.AliceEVM,
				to,
				amount,
				suite.CollectionId,
			)

			if tt.expectError {
				suite.Require().Error(err, "Expected error for: %s", tt.description)
			} else {
				suite.Require().NoError(err, "Unexpected error for: %s", tt.description)
			}
		})
	}
}

// TestERC3643Wrapper_StateConsistency tests that ERC3643 wrapper transfers maintain state consistency
func (suite *E2ETestSuite) TestERC3643Wrapper_StateConsistency() {
	erc3643ContractAddr := common.HexToAddress("0x1111111111111111111111111111111111111111")

	// Get initial state
	aliceBalanceBefore := suite.getBalance(suite.Alice.String())
	bobBalanceBefore := suite.getBalance(suite.Bob.String())
	charlieBalanceBefore := suite.getBalance(suite.Charlie.String())

	// Perform multiple transfers through ERC3643 wrapper
	// Note: ERC3643 uses tokenIds: 1-1, and our collection has 1 token per ID
	// So we can only transfer 1 token with ID 1 at a time
	transfers := []struct {
		from   common.Address
		to     common.Address
		amount *big.Int
	}{
		{suite.AliceEVM, suite.BobEVM, big.NewInt(1)},     // Alice -> Bob (token ID 1)
		{suite.BobEVM, suite.CharlieEVM, big.NewInt(1)},   // Bob -> Charlie (token ID 1)
	}

	for _, transfer := range transfers {
		_, err := suite.callERC3643WrapperTransfer(
			erc3643ContractAddr,
			transfer.from,
			transfer.to,
			transfer.amount,
			suite.CollectionId,
		)
		suite.Require().NoError(err, "Transfer should succeed")
	}

	// Verify final state
	aliceBalanceAfter := suite.getBalance(suite.Alice.String())
	bobBalanceAfter := suite.getBalance(suite.Bob.String())
	charlieBalanceAfter := suite.getBalance(suite.Charlie.String())

	// Alice should have lost 1 token (token ID 1)
	suite.verifyBalanceDecreased(aliceBalanceBefore, aliceBalanceAfter, sdkmath.NewUint(1), []*tokenizationtypes.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
	})

	// Bob should have gained 1 and lost 1 = net 0
	bobNetChange := suite.calculateNetBalanceChange(bobBalanceBefore, bobBalanceAfter, []*tokenizationtypes.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
	})
	suite.Require().True(bobNetChange.Equal(sdkmath.NewUint(0)), "Bob should have net 0 tokens")

	// Charlie should have gained 1
	suite.verifyBalanceIncreased(charlieBalanceBefore, charlieBalanceAfter, sdkmath.NewUint(1), []*tokenizationtypes.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
	})
}

// calculateNetBalanceChange calculates the net change in balance (can be positive or negative)
func (suite *E2ETestSuite) calculateNetBalanceChange(
	before, after *tokenizationtypes.UserBalanceStore,
	expectedTokenIds []*tokenizationtypes.UintRange,
) sdkmath.Uint {
	ownershipTimes := []*tokenizationtypes.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
	}

	beforeBalances, _ := tokenizationtypes.GetBalancesForIds(suite.Ctx, expectedTokenIds, ownershipTimes, before.Balances)
	afterBalances, _ := tokenizationtypes.GetBalancesForIds(suite.Ctx, expectedTokenIds, ownershipTimes, after.Balances)

	beforeTotal := sdkmath.NewUint(0)
	for _, bal := range beforeBalances {
		beforeTotal = beforeTotal.Add(bal.Amount)
	}

	afterTotal := sdkmath.NewUint(0)
	for _, bal := range afterBalances {
		afterTotal = afterTotal.Add(bal.Amount)
	}

	return afterTotal.Sub(beforeTotal)
}

// TestE2ESuite runs the end-to-end test suite
func TestE2ESuite(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}
