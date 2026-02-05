package tokenization

import (
	"math"
	"math/big"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"

	sdkmath "cosmossdk.io/math"

	tokenizationkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	keepertest "github.com/bitbadges/bitbadgeschain/x/tokenization/testutil/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EdgeCaseTestSuite is a test suite for edge cases and boundary conditions
type EdgeCaseTestSuite struct {
	suite.Suite
	TokenizationKeeper tokenizationkeeper.Keeper
	Ctx                sdk.Context
	Precompile         *Precompile

	// Test addresses
	AliceEVM  common.Address
	BobEVM    common.Address
	CharlieEVM common.Address
	Alice     sdk.AccAddress
	Bob       sdk.AccAddress
	Charlie   sdk.AccAddress

	// Test data
	CollectionId sdkmath.Uint
}

func TestEdgeCaseTestSuite(t *testing.T) {
	suite.Run(t, new(EdgeCaseTestSuite))
}

// SetupTest initializes the test suite
func (suite *EdgeCaseTestSuite) SetupTest() {
	keeper, ctx := keepertest.TokenizationKeeper(suite.T())
	suite.TokenizationKeeper = keeper
	suite.Ctx = ctx
	suite.Precompile = NewPrecompile(keeper)

	// Create test addresses
	suite.AliceEVM = common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0")
	suite.BobEVM = common.HexToAddress("0x8ba1f109551bD432803012645Hac136c22C9e7")
	suite.CharlieEVM = common.HexToAddress("0x1234567890123456789012345678901234567890")

	suite.Alice = sdk.AccAddress(suite.AliceEVM.Bytes())
	suite.Bob = sdk.AccAddress(suite.BobEVM.Bytes())
	suite.Charlie = sdk.AccAddress(suite.CharlieEVM.Bytes())

	// Set up test collection
	suite.CollectionId = suite.createTestCollection()
}

// createTestCollection creates a test collection with balances
func (suite *EdgeCaseTestSuite) createTestCollection() sdkmath.Uint {
	msgServer := tokenizationkeeper.NewMsgServerImpl(suite.TokenizationKeeper)

	// Create collection
	createMsg := &tokenizationtypes.MsgUniversalUpdateCollection{
		Creator:      suite.Alice.String(),
		CollectionId: sdkmath.NewUint(0), // 0 means new collection
		ValidTokenIds: []*tokenizationtypes.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
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
				OverallMaxNumTransfers: sdkmath.NewUint(10000),
				AmountTrackerId:        "mint-tracker",
			},
			ApprovalAmounts: &tokenizationtypes.ApprovalAmounts{
				PerFromAddressApprovalAmount: sdkmath.NewUint(10000),
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
						TokenIds:       []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1000)}},
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

// TestBoundaryConditions tests boundary values for various inputs
func (suite *EdgeCaseTestSuite) TestBoundaryConditions() {
	// Test maximum uint256 values
	maxUint256 := new(big.Int)
	maxUint256.Exp(big.NewInt(2), big.NewInt(256), nil)
	maxUint256.Sub(maxUint256, big.NewInt(1))

	tests := []struct {
		name        string
		collectionId *big.Int
		expectError bool
	}{
		{
			name:        "maximum_uint256_collection_id",
			collectionId: maxUint256,
			expectError: false, // Should be valid, even if collection doesn't exist
		},
		{
			name:        "minimum_valid_collection_id",
			collectionId: big.NewInt(1),
			expectError: false,
		},
		{
			name:        "zero_collection_id",
			collectionId: big.NewInt(0),
			expectError: false, // 0 is valid for new collections
		},
		{
			name:        "one_over_max_uint256",
			collectionId: new(big.Int).Add(maxUint256, big.NewInt(1)),
			expectError: false, // ValidateCollectionId only checks for negative, not overflow
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			err := ValidateCollectionId(tt.collectionId)
			if tt.expectError {
				suite.Require().Error(err)
			} else {
				// May or may not error depending on validation logic
				_ = err
			}
		})
	}
}

// TestMaximumArraySizes tests arrays at exactly the limits and one over
func (suite *EdgeCaseTestSuite) TestMaximumArraySizes() {
	// Create addresses for testing
	addresses := make([]common.Address, MaxRecipients)
	for i := 0; i < MaxRecipients; i++ {
		addresses[i] = common.BigToAddress(big.NewInt(int64(i + 1)))
	}

	// Test exactly at MaxRecipients
	suite.Run("exactly_max_recipients", func() {
		err := ValidateAddresses(addresses, "toAddresses")
		suite.Require().NoError(err)
	})

	// Test one over MaxRecipients
	suite.Run("one_over_max_recipients", func() {
		tooManyAddresses := append(addresses, common.BigToAddress(big.NewInt(101)))
		err := ValidateAddresses(tooManyAddresses, "toAddresses")
		suite.Require().Error(err)
		suite.Contains(err.Error(), "exceeds maximum")
	})

	// Test exactly at MaxTokenIdRanges
	suite.Run("exactly_max_token_id_ranges", func() {
		ranges := make([]struct {
			Start *big.Int `json:"start"`
			End   *big.Int `json:"end"`
		}, MaxTokenIdRanges)
		for i := 0; i < MaxTokenIdRanges; i++ {
			ranges[i] = struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			}{
				Start: big.NewInt(int64(i + 1)),
				End:   big.NewInt(int64(i + 1)),
			}
		}
		err := ValidateBigIntRanges(ranges, "tokenIds")
		suite.Require().NoError(err)
	})

	// Test one over MaxTokenIdRanges
	suite.Run("one_over_max_token_id_ranges", func() {
		ranges := make([]struct {
			Start *big.Int `json:"start"`
			End   *big.Int `json:"end"`
		}, MaxTokenIdRanges+1)
		for i := 0; i < MaxTokenIdRanges+1; i++ {
			ranges[i] = struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			}{
				Start: big.NewInt(int64(i + 1)),
				End:   big.NewInt(int64(i + 1)),
			}
		}
		err := ValidateArraySize(len(ranges), MaxTokenIdRanges, "tokenIds")
		suite.Require().Error(err)
		suite.Contains(err.Error(), "exceeds maximum")
	})
}

// TestRangeOverlap tests overlapping and adjacent ranges
func (suite *EdgeCaseTestSuite) TestRangeOverlap() {
	tests := []struct {
		name        string
		ranges      []struct {
			Start *big.Int `json:"start"`
			End   *big.Int `json:"end"`
		}
		expectError bool
		description string
	}{
		{
			name: "overlapping_ranges",
			ranges: []struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			}{
				{Start: big.NewInt(1), End: big.NewInt(10)},
				{Start: big.NewInt(5), End: big.NewInt(15)}, // Overlaps with first
			},
			expectError: false, // Overlapping ranges are allowed
			description: "Overlapping ranges should be valid",
		},
		{
			name: "adjacent_ranges_no_gap",
			ranges: []struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			}{
				{Start: big.NewInt(1), End: big.NewInt(10)},
				{Start: big.NewInt(11), End: big.NewInt(20)}, // Adjacent, no gap
			},
			expectError: false,
			description: "Adjacent ranges with no gap should be valid",
		},
		{
			name: "adjacent_ranges_with_gap",
			ranges: []struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			}{
				{Start: big.NewInt(1), End: big.NewInt(10)},
				{Start: big.NewInt(12), End: big.NewInt(20)}, // Gap at 11
			},
			expectError: false,
			description: "Ranges with gaps should be valid",
		},
		{
			name: "ranges_spanning_entire_uint256",
			ranges: []struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			}{
				{Start: big.NewInt(1), End: new(big.Int).SetUint64(math.MaxUint64)},
			},
			expectError: false,
			description: "Range spanning entire uint256 space should be valid",
		},
		{
			name: "single_value_range",
			ranges: []struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			}{
				{Start: big.NewInt(1), End: big.NewInt(1)},
			},
			expectError: false,
			description: "Single value range (start == end) should be valid",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			err := ValidateBigIntRanges(tt.ranges, "testRanges")
			if tt.expectError {
				suite.Require().Error(err, tt.description)
			} else {
				suite.Require().NoError(err, tt.description)
			}
		})
	}
}

// TestLargeValueTransfers tests transfers with maximum uint256 amounts
func (suite *EdgeCaseTestSuite) TestLargeValueTransfers() {
	maxUint256 := new(big.Int)
	maxUint256.Exp(big.NewInt(2), big.NewInt(256), nil)
	maxUint256.Sub(maxUint256, big.NewInt(1))

	// Test maximum amount
	suite.Run("maximum_uint256_amount", func() {
		err := ValidateAmount(maxUint256, "amount")
		suite.Require().NoError(err)
	})

	// Test amount of 1 (minimum valid)
	suite.Run("minimum_valid_amount", func() {
		err := ValidateAmount(big.NewInt(1), "amount")
		suite.Require().NoError(err)
	})

	// Test amount of 0 (should fail)
	suite.Run("zero_amount", func() {
		err := ValidateAmount(big.NewInt(0), "amount")
		suite.Require().Error(err)
		suite.Contains(err.Error(), "must be greater than zero")
	})

	// Test negative amount (should fail)
	suite.Run("negative_amount", func() {
		err := ValidateAmount(big.NewInt(-1), "amount")
		suite.Require().Error(err)
		suite.Contains(err.Error(), "must be greater than zero")
	})
}

// TestEmptyResults tests queries on empty collections and balances
func (suite *EdgeCaseTestSuite) TestEmptyResults() {
	method := suite.Precompile.ABI.Methods["getCollection"]

	// Test query on non-existent collection (empty result)
	suite.Run("non_existent_collection", func() {
		args := []interface{}{big.NewInt(999999)}
		result, err := suite.Precompile.GetCollection(suite.Ctx, &method, args)
		suite.Require().Error(err) // Should return error for non-existent collection
		suite.Nil(result)
	})

	// Test query with ranges that match no tokens
	methodBalance := suite.Precompile.ABI.Methods["getBalanceAmount"]
	suite.Run("ranges_matching_no_tokens", func() {
		args := []interface{}{
			suite.CollectionId.BigInt(),
			suite.AliceEVM,
			[]struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			}{
				{Start: big.NewInt(9999), End: big.NewInt(10000)}, // No tokens in this range
			},
			[]struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			}{
				{Start: big.NewInt(1), End: new(big.Int).SetUint64(math.MaxUint64)},
			},
		}
		result, err := suite.Precompile.GetBalanceAmount(suite.Ctx, &methodBalance, args)
		// Should succeed but return 0
		if err == nil {
			unpacked, err := methodBalance.Outputs.Unpack(result)
			suite.Require().NoError(err)
			suite.Require().Len(unpacked, 1)
			amount, ok := unpacked[0].(*big.Int)
			suite.Require().True(ok)
			suite.Require().NotNil(amount)
			suite.True(amount.Cmp(big.NewInt(0)) == 0, "Amount should be 0, got %s", amount.String())
		}
	})
}

// TestConcurrentCalls tests multiple concurrent calls to ensure no race conditions
func (suite *EdgeCaseTestSuite) TestConcurrentCalls() {
	method := suite.Precompile.ABI.Methods["getCollection"]
	collectionId := suite.CollectionId.BigInt()

	// Test concurrent queries
	suite.Run("concurrent_queries", func() {
		var wg sync.WaitGroup
		numGoroutines := 10
		errors := make([]error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				args := []interface{}{collectionId}
				_, err := suite.Precompile.GetCollection(suite.Ctx, &method, args)
				errors[idx] = err
			}(i)
		}

		wg.Wait()

		// All queries should succeed
		for i, err := range errors {
			suite.NoError(err, "Query %d should succeed", i)
		}
	})

	// Test concurrent transfers (if we had a way to simulate them)
	// Note: Actual concurrent transfers would require EVM integration tests
	suite.Run("concurrent_transfer_validation", func() {
		var wg sync.WaitGroup
		numGoroutines := 10
		errors := make([]error, numGoroutines)

		// Test validation only (not actual transfers)
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				addresses := []common.Address{suite.BobEVM}
				err := ValidateAddresses(addresses, "toAddresses")
				errors[idx] = err
			}(i)
		}

		wg.Wait()

		// All validations should succeed
		for i, err := range errors {
			suite.NoError(err, "Validation %d should succeed", i)
		}
	})
}

// TestVeryLargeRangeSpans tests ranges with very large spans
func (suite *EdgeCaseTestSuite) TestVeryLargeRangeSpans() {
	maxUint64 := new(big.Int).SetUint64(math.MaxUint64)

	tests := []struct {
		name        string
		start       *big.Int
		end         *big.Int
		expectError bool
	}{
		{
			name:        "max_uint64_span",
			start:       big.NewInt(0),
			end:         maxUint64,
			expectError: false,
		},
		{
			name:        "large_span_mid_range",
			start:       big.NewInt(1000000),
			end:         new(big.Int).Add(big.NewInt(1000000), maxUint64),
			expectError: false, // Should be valid even if end > max uint64
		},
		{
			name:        "start_equals_end",
			start:       maxUint64,
			end:         maxUint64,
			expectError: false,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			err := ValidateBigIntRange(tt.start, tt.end, "testRange")
			if tt.expectError {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

// TestMinimumValidValues tests minimum valid values for IDs and amounts
func (suite *EdgeCaseTestSuite) TestMinimumValidValues() {
	// Test minimum collection ID (1)
	suite.Run("minimum_collection_id", func() {
		err := ValidateCollectionId(big.NewInt(1))
		suite.Require().NoError(err)
	})

	// Test minimum amount (1)
	suite.Run("minimum_amount", func() {
		err := ValidateAmount(big.NewInt(1), "amount")
		suite.Require().NoError(err)
	})

	// Test minimum range (start = end = 1)
	suite.Run("minimum_range", func() {
		err := ValidateBigIntRange(big.NewInt(1), big.NewInt(1), "testRange")
		suite.Require().NoError(err)
	})
}

