package validation

import (
	"math"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// ============================================================================
// UintRange Cross-Field Validation Tests
// ============================================================================

func TestUintRange_CrossField_StartEndConsistency(t *testing.T) {
	// Test that Start <= End is enforced across all contexts
	testCases := []struct {
		name  string
		start uint64
		end   uint64
		valid bool
	}{
		{"start equals end", 5, 5, true},
		{"start less than end", 1, 10, true},
		{"start greater than end", 10, 1, false},
		{"boundary: max values", math.MaxUint64, math.MaxUint64, true},
		{"boundary: min values", 1, 1, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			range1 := &types.UintRange{
				Start: sdkmath.NewUint(tc.start),
				End:   sdkmath.NewUint(tc.end),
			}

			err := types.ValidateRangesAreValid([]*types.UintRange{range1}, false, false)
			if tc.valid {
				require.NoError(t, err, "should be valid")
			} else {
				require.Error(t, err, "should be invalid")
			}
		})
	}
}

// ============================================================================
// Balance Cross-Field Validation Tests
// ============================================================================

func TestBalance_CrossField_TokenIdsOwnershipTimesConsistency(t *testing.T) {
	// Test that TokenIds and OwnershipTimes are validated together
	balance := &types.Balance{
		Amount: sdkmath.NewUint(1),
		TokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
		},
		OwnershipTimes: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
		},
	}

	ctx := CreateTestContext()
	_, err := types.ValidateBalances(ctx, []*types.Balance{balance}, false)
	require.NoError(t, err, "consistent token IDs and ownership times should pass")
}

func TestBalance_CrossField_OverlappingTokenIdsAndOwnershipTimes(t *testing.T) {
	// Test balance with overlapping token IDs (invalid)
	balance := &types.Balance{
		Amount: sdkmath.NewUint(1),
		TokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
			{Start: sdkmath.NewUint(5), End: sdkmath.NewUint(20)}, // Overlaps
		},
		OwnershipTimes: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
		},
	}

	ctx := CreateTestContext()
	_, err := types.ValidateBalances(ctx, []*types.Balance{balance}, false)
	require.Error(t, err, "overlapping token IDs should fail")
}

func TestBalance_CrossField_OverlappingOwnershipTimesAndTokenIds(t *testing.T) {
	// Test balance with overlapping ownership times (invalid)
	balance := &types.Balance{
		Amount: sdkmath.NewUint(1),
		TokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
		},
		OwnershipTimes: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
			{Start: sdkmath.NewUint(50), End: sdkmath.NewUint(200)}, // Overlaps
		},
	}

	ctx := CreateTestContext()
	_, err := types.ValidateBalances(ctx, []*types.Balance{balance}, false)
	require.Error(t, err, "overlapping ownership times should fail")
}

// ============================================================================
// UintRange + Balance Integration Tests
// ============================================================================

func TestUintRangeBalance_Integration_ValidBalanceWithMultipleRanges(t *testing.T) {
	// Test balance with multiple non-overlapping token ID and ownership time ranges
	balance := &types.Balance{
		Amount: sdkmath.NewUint(5),
		TokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
			{Start: sdkmath.NewUint(20), End: sdkmath.NewUint(30)},
		},
		OwnershipTimes: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(50)},
			{Start: sdkmath.NewUint(100), End: sdkmath.NewUint(200)},
		},
	}

	ctx := CreateTestContext()
	_, err := types.ValidateBalances(ctx, []*types.Balance{balance}, false)
	require.NoError(t, err, "valid balance with multiple ranges should pass")
}

func TestUintRangeBalance_Integration_InvalidBalanceWithOverlappingRanges(t *testing.T) {
	// Test balance with overlapping ranges in both TokenIds and OwnershipTimes
	balance := &types.Balance{
		Amount: sdkmath.NewUint(5),
		TokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
			{Start: sdkmath.NewUint(5), End: sdkmath.NewUint(15)}, // Overlaps
		},
		OwnershipTimes: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
			{Start: sdkmath.NewUint(50), End: sdkmath.NewUint(200)}, // Overlaps
		},
	}

	ctx := CreateTestContext()
	_, err := types.ValidateBalances(ctx, []*types.Balance{balance}, false)
	require.Error(t, err, "overlapping ranges should fail")
}

func TestUintRangeBalance_Integration_BalanceOperationsWithRanges(t *testing.T) {
	// Test that balance operations handle ranges correctly
	existing := []*types.Balance{
		{
			Amount: sdkmath.NewUint(5),
			TokenIds: []*types.UintRange{
				{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
			},
			OwnershipTimes: []*types.UintRange{
				{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
			},
		},
	}

	toAdd := []*types.Balance{
		{
			Amount: sdkmath.NewUint(3),
			TokenIds: []*types.UintRange{
				{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
			},
			OwnershipTimes: []*types.UintRange{
				{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
			},
		},
	}

	ctx := CreateTestContext()
	result, err := types.AddBalances(ctx, toAdd, existing)
	require.NoError(t, err, "adding balances should succeed")
	require.Len(t, result, 1, "should have one balance after merge")
	require.Equal(t, sdkmath.NewUint(8), result[0].Amount, "amounts should be added")
}

func TestUintRangeBalance_Integration_ComplexBalanceOperations(t *testing.T) {
	// Test complex balance operations with multiple ranges
	existing := []*types.Balance{
		{
			Amount: sdkmath.NewUint(5),
			TokenIds: []*types.UintRange{
				{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
				{Start: sdkmath.NewUint(20), End: sdkmath.NewUint(30)},
			},
			OwnershipTimes: []*types.UintRange{
				{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
			},
		},
	}

	toAdd := []*types.Balance{
		{
			Amount: sdkmath.NewUint(3),
			TokenIds: []*types.UintRange{
				{Start: sdkmath.NewUint(5), End: sdkmath.NewUint(15)}, // Overlaps with first range
			},
			OwnershipTimes: []*types.UintRange{
				{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
			},
		},
	}

	ctx := CreateTestContext()
	result, err := types.AddBalances(ctx, toAdd, existing)
	require.NoError(t, err, "complex balance operations should succeed")
	require.NotNil(t, result, "should return result")
}

// ============================================================================
// Multi-Level Validation Tests
// ============================================================================

func TestMultiLevel_TransferWithInvalidBalances(t *testing.T) {
	// Test transfer with invalid balances (zero amount)
	transfer := &types.Transfer{
		From:        "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		ToAddresses: []string{"bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q"},
		Balances: []*types.Balance{
			{
				Amount: sdkmath.NewUint(0), // Invalid
				TokenIds: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
				},
				OwnershipTimes: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
				},
			},
		},
	}

	ctx := CreateTestContext()
	err := types.ValidateTransfer(ctx, transfer, false)
	require.Error(t, err, "transfer with invalid balances should fail")
}

func TestMultiLevel_ApprovalWithInvalidRanges(t *testing.T) {
	// Test approval with invalid ranges
	approval := &types.CollectionApproval{
		ApprovalId:        "test_approval",
		FromListId:        "All",
		ToListId:          "All",
		InitiatedByListId: "All",
		TokenIds: []*types.UintRange{
			{
				Start: sdkmath.NewUint(100),
				End:   sdkmath.NewUint(1), // Invalid: start > end
			},
		},
		TransferTimes: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
		},
		OwnershipTimes: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
		},
		ApprovalCriteria: &types.ApprovalCriteria{},
		Version:          sdkmath.NewUint(0),
	}

	ctx := CreateTestContext()
	err := types.ValidateCollectionApprovals(ctx, []*types.CollectionApproval{approval}, false)
	require.Error(t, err, "approval with invalid ranges should fail")
}

func TestMultiLevel_AddressListWithInvalidAddresses(t *testing.T) {
	// Test address list with invalid addresses
	addressList := &types.AddressList{
		ListId: "test_list",
		Addresses: []string{
			"bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
			"invalid_address", // Invalid
		},
	}

	err := types.ValidateAddressList(addressList)
	require.Error(t, err, "address list with invalid addresses should fail")
}

func TestMultiLevel_MessageWithMultipleValidationFailures(t *testing.T) {
	// Test message with multiple validation failures
	msg := &types.MsgTransferTokens{
		Creator:      "invalid_address", // Invalid creator
		CollectionId: sdkmath.NewUint(0), // Zero collection ID
		Transfers: []*types.Transfer{
			{
				From:        "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
				ToAddresses: []string{}, // Empty to addresses
				Balances: []*types.Balance{
					{
						Amount: sdkmath.NewUint(0), // Zero amount
						TokenIds: []*types.UintRange{
							{Start: sdkmath.NewUint(100), End: sdkmath.NewUint(1)}, // Invalid range
						},
						OwnershipTimes: []*types.UintRange{
							{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
						},
					},
				},
			},
		},
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "message with multiple failures should fail")
	// Should fail on first validation error (creator address)
}

