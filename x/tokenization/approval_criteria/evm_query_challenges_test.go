package approval_criteria

import (
	"encoding/hex"
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// Valid 0x+40 hex addresses for use in Check(to, from, initiator, ...) so placeholder replacement succeeds
const (
	testAddrTo       = "0x0000000000000000000000000000000000000001"
	testAddrFrom     = "0x0000000000000000000000000000000000000002"
	testAddrInitiator = "0x0000000000000000000000000000000000000003"
)

// mockEVMQueryService is a mock implementation of EVMQueryService for testing
type mockEVMQueryService struct {
	returnValue []byte
	returnError error
	// Track calls for verification
	lastContractAddress string
	lastCalldata        []byte
	lastGasLimit        uint64
}

func (m *mockEVMQueryService) ExecuteEVMQuery(ctx sdk.Context, callerAddress string, contractAddress string, calldata []byte, gasLimit uint64) ([]byte, error) {
	m.lastContractAddress = contractAddress
	m.lastCalldata = calldata
	m.lastGasLimit = gasLimit
	return m.returnValue, m.returnError
}

// Helper to create a mock context
func mockContext() sdk.Context {
	return sdk.Context{}
}

// Helper to create a basic approval with EVM query challenges
func createApprovalWithEVMQueryChallenges(challenges []*types.EVMQueryChallenge) *types.CollectionApproval {
	return &types.CollectionApproval{
		ApprovalId: "test-approval",
		ApprovalCriteria: &types.ApprovalCriteria{
			EvmQueryChallenges: challenges,
		},
	}
}

// Helper to create a basic collection
func createCollection() *types.TokenCollection {
	return &types.TokenCollection{
		CollectionId: sdkmath.NewUint(1),
	}
}

func TestEVMQueryChallengesChecker_Name(t *testing.T) {
	checker := NewEVMQueryChallengesChecker(&mockEVMQueryService{})
	require.Equal(t, "EVMQueryChallenges", checker.Name())
}

func TestEVMQueryChallengesChecker_NoChallenges(t *testing.T) {
	mockService := &mockEVMQueryService{}
	checker := NewEVMQueryChallengesChecker(mockService)

	// Test with nil approval criteria
	approval := &types.CollectionApproval{}
	errMsg, err := checker.Check(mockContext(), approval, createCollection(), testAddrTo, testAddrFrom, testAddrInitiator, "collection", "approver", nil, nil, "", false)
	require.NoError(t, err)
	require.Empty(t, errMsg)

	// Test with empty challenges
	approval = createApprovalWithEVMQueryChallenges([]*types.EVMQueryChallenge{})
	errMsg, err = checker.Check(mockContext(), approval, createCollection(), testAddrTo, testAddrFrom, testAddrInitiator, "collection", "approver", nil, nil, "", false)
	require.NoError(t, err)
	require.Empty(t, errMsg)
}

func TestEVMQueryChallengesChecker_BasicSuccess(t *testing.T) {
	// Return value that matches expected
	expectedResult, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000001")
	mockService := &mockEVMQueryService{
		returnValue: expectedResult,
	}
	checker := NewEVMQueryChallengesChecker(mockService)

	challenges := []*types.EVMQueryChallenge{
		{
			ContractAddress:    "0x1234567890123456789012345678901234567890",
			Calldata:           "70a08231000000000000000000000000abcdef1234567890abcdef1234567890abcdef12",
			ExpectedResult:     "0000000000000000000000000000000000000000000000000000000000000001",
			ComparisonOperator: "eq",
			GasLimit:           sdkmath.NewUint(100000),
		},
	}
	approval := createApprovalWithEVMQueryChallenges(challenges)

	errMsg, err := checker.Check(mockContext(), approval, createCollection(), testAddrTo, testAddrFrom, testAddrInitiator, "collection", "approver", nil, nil, "", false)
	require.NoError(t, err)
	require.Empty(t, errMsg)
}

func TestEVMQueryChallengesChecker_ResultMismatch(t *testing.T) {
	// Return value that doesn't match expected
	returnedResult, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000002")
	mockService := &mockEVMQueryService{
		returnValue: returnedResult,
	}
	checker := NewEVMQueryChallengesChecker(mockService)

	challenges := []*types.EVMQueryChallenge{
		{
			ContractAddress:    "0x1234567890123456789012345678901234567890",
			Calldata:           "70a08231",
			ExpectedResult:     "0000000000000000000000000000000000000000000000000000000000000001",
			ComparisonOperator: "eq",
			GasLimit:           sdkmath.NewUint(100000),
		},
	}
	approval := createApprovalWithEVMQueryChallenges(challenges)

	errMsg, err := checker.Check(mockContext(), approval, createCollection(), testAddrTo, testAddrFrom, testAddrInitiator, "collection", "approver", nil, nil, "", false)
	require.Error(t, err)
	require.Contains(t, errMsg, "result mismatch")
}

func TestEVMQueryChallengesChecker_ComparisonOperators(t *testing.T) {
	tests := []struct {
		name           string
		returnValue    string // hex string
		expectedResult string // hex string
		operator       string
		shouldPass     bool
	}{
		// Equal operator
		{"eq_equal", "0000000000000000000000000000000000000000000000000000000000000005", "0000000000000000000000000000000000000000000000000000000000000005", "eq", true},
		{"eq_not_equal", "0000000000000000000000000000000000000000000000000000000000000005", "0000000000000000000000000000000000000000000000000000000000000006", "eq", false},
		{"empty_op_equal", "0000000000000000000000000000000000000000000000000000000000000005", "0000000000000000000000000000000000000000000000000000000000000005", "", true},

		// Not equal operator
		{"ne_not_equal", "0000000000000000000000000000000000000000000000000000000000000005", "0000000000000000000000000000000000000000000000000000000000000006", "ne", true},
		{"ne_equal", "0000000000000000000000000000000000000000000000000000000000000005", "0000000000000000000000000000000000000000000000000000000000000005", "ne", false},

		// Greater than
		{"gt_greater", "0000000000000000000000000000000000000000000000000000000000000006", "0000000000000000000000000000000000000000000000000000000000000005", "gt", true},
		{"gt_equal", "0000000000000000000000000000000000000000000000000000000000000005", "0000000000000000000000000000000000000000000000000000000000000005", "gt", false},
		{"gt_less", "0000000000000000000000000000000000000000000000000000000000000004", "0000000000000000000000000000000000000000000000000000000000000005", "gt", false},

		// Greater than or equal
		{"gte_greater", "0000000000000000000000000000000000000000000000000000000000000006", "0000000000000000000000000000000000000000000000000000000000000005", "gte", true},
		{"gte_equal", "0000000000000000000000000000000000000000000000000000000000000005", "0000000000000000000000000000000000000000000000000000000000000005", "gte", true},
		{"gte_less", "0000000000000000000000000000000000000000000000000000000000000004", "0000000000000000000000000000000000000000000000000000000000000005", "gte", false},

		// Less than
		{"lt_less", "0000000000000000000000000000000000000000000000000000000000000004", "0000000000000000000000000000000000000000000000000000000000000005", "lt", true},
		{"lt_equal", "0000000000000000000000000000000000000000000000000000000000000005", "0000000000000000000000000000000000000000000000000000000000000005", "lt", false},
		{"lt_greater", "0000000000000000000000000000000000000000000000000000000000000006", "0000000000000000000000000000000000000000000000000000000000000005", "lt", false},

		// Less than or equal
		{"lte_less", "0000000000000000000000000000000000000000000000000000000000000004", "0000000000000000000000000000000000000000000000000000000000000005", "lte", true},
		{"lte_equal", "0000000000000000000000000000000000000000000000000000000000000005", "0000000000000000000000000000000000000000000000000000000000000005", "lte", true},
		{"lte_greater", "0000000000000000000000000000000000000000000000000000000000000006", "0000000000000000000000000000000000000000000000000000000000000005", "lte", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			returnValue, _ := hex.DecodeString(tc.returnValue)
			mockService := &mockEVMQueryService{
				returnValue: returnValue,
			}
			checker := NewEVMQueryChallengesChecker(mockService)

			challenges := []*types.EVMQueryChallenge{
				{
					ContractAddress:    "0x1234567890123456789012345678901234567890",
					Calldata:           "70a08231",
					ExpectedResult:     tc.expectedResult,
					ComparisonOperator: tc.operator,
					GasLimit:           sdkmath.NewUint(100000),
				},
			}
			approval := createApprovalWithEVMQueryChallenges(challenges)

			errMsg, err := checker.Check(mockContext(), approval, createCollection(), testAddrTo, testAddrFrom, testAddrInitiator, "collection", "approver", nil, nil, "", false)
			if tc.shouldPass {
				require.NoError(t, err, "expected check to pass for %s", tc.name)
				require.Empty(t, errMsg)
			} else {
				require.Error(t, err, "expected check to fail for %s", tc.name)
				require.Contains(t, errMsg, "result mismatch")
			}
		})
	}
}

func TestEVMQueryChallengesChecker_PlaceholderReplacement(t *testing.T) {
	mockService := &mockEVMQueryService{
		returnValue: []byte{0x01}, // Any non-empty result
	}
	checker := NewEVMQueryChallengesChecker(mockService)

	// Use placeholders in calldata - the placeholder gets replaced with a 64-char hex string (32 bytes padded)
	// Format: function selector (4 bytes = 8 hex chars) + $placeholder (gets replaced with 64 hex chars)
	challenges := []*types.EVMQueryChallenge{
		{
			ContractAddress: "0x1234567890123456789012345678901234567890",
			Calldata:        "70a08231$sender", // balanceOf with sender placeholder - after replacement becomes valid hex
			GasLimit:        sdkmath.NewUint(100000),
		},
	}
	approval := createApprovalWithEVMQueryChallenges(challenges)

	// Use a valid 0x hex address for the sender (will be converted to 64-char padded hex)
	sender := "0xabcdef1234567890abcdef1234567890abcdef12"

	_, err := checker.Check(mockContext(), approval, createCollection(), testAddrTo, sender, testAddrInitiator, "collection", "approver", nil, nil, "", false)
	require.NoError(t, err)

	// Verify the calldata was transformed - check that the service received proper calldata
	// 4 bytes function selector + 32 bytes address = 36 bytes
	require.NotNil(t, mockService.lastCalldata)
	require.Equal(t, 36, len(mockService.lastCalldata), "calldata should be 4 bytes selector + 32 bytes address")
}

func TestEVMQueryChallengesChecker_GasLimitEnforcement(t *testing.T) {
	mockService := &mockEVMQueryService{
		returnValue: []byte{0x01},
	}
	checker := NewEVMQueryChallengesChecker(mockService)

	tests := []struct {
		name           string
		gasLimit       uint64
		expectedLimit  uint64
	}{
		{"default_gas_limit", 0, DefaultEVMQueryGasLimit},
		{"custom_gas_limit", 50000, 50000},
		{"max_gas_limit_exceeded", 1000000, MaxEVMQueryGasLimit}, // Should be capped
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			challenges := []*types.EVMQueryChallenge{
				{
					ContractAddress: "0x1234567890123456789012345678901234567890",
					Calldata:        "70a08231",
					GasLimit:        sdkmath.NewUint(tc.gasLimit),
				},
			}
			approval := createApprovalWithEVMQueryChallenges(challenges)

			_, _ = checker.Check(mockContext(), approval, createCollection(), testAddrTo, testAddrFrom, testAddrInitiator, "collection", "approver", nil, nil, "", false)

			require.Equal(t, tc.expectedLimit, mockService.lastGasLimit, "gas limit should be %d for %s", tc.expectedLimit, tc.name)
		})
	}
}

func TestEVMQueryChallengesChecker_TotalGasLimit(t *testing.T) {
	mockService := &mockEVMQueryService{
		returnValue: []byte{0x01},
	}
	checker := NewEVMQueryChallengesChecker(mockService)

	// Create 3 challenges with 500k gas each = 1.5M total, exceeds 1M limit
	challenges := []*types.EVMQueryChallenge{
		{
			ContractAddress: "0x1234567890123456789012345678901234567890",
			Calldata:        "70a08231",
			GasLimit:        sdkmath.NewUint(500000),
		},
		{
			ContractAddress: "0x1234567890123456789012345678901234567890",
			Calldata:        "70a08231",
			GasLimit:        sdkmath.NewUint(500000),
		},
		{
			ContractAddress: "0x1234567890123456789012345678901234567890",
			Calldata:        "70a08231",
			GasLimit:        sdkmath.NewUint(500000), // This would exceed 1M total
		},
	}
	approval := createApprovalWithEVMQueryChallenges(challenges)

	errMsg, err := checker.Check(mockContext(), approval, createCollection(), testAddrTo, testAddrFrom, testAddrInitiator, "collection", "approver", nil, nil, "", false)
	require.Error(t, err, "should fail when total gas exceeds limit")
	require.Contains(t, errMsg, "exceed total gas limit", "error should mention total gas limit")
}

func TestEVMQueryChallengesChecker_TotalGasLimitPasses(t *testing.T) {
	mockService := &mockEVMQueryService{
		returnValue: []byte{0x01},
	}
	checker := NewEVMQueryChallengesChecker(mockService)

	// Create 2 challenges with 500k gas each = 1M total, exactly at limit
	challenges := []*types.EVMQueryChallenge{
		{
			ContractAddress: "0x1234567890123456789012345678901234567890",
			Calldata:        "70a08231",
			GasLimit:        sdkmath.NewUint(500000),
		},
		{
			ContractAddress: "0x1234567890123456789012345678901234567890",
			Calldata:        "70a08231",
			GasLimit:        sdkmath.NewUint(500000),
		},
	}
	approval := createApprovalWithEVMQueryChallenges(challenges)

	errMsg, err := checker.Check(mockContext(), approval, createCollection(), testAddrTo, testAddrFrom, testAddrInitiator, "collection", "approver", nil, nil, "", false)
	require.NoError(t, err, "should pass when total gas is at limit")
	require.Empty(t, errMsg)
}

func TestEVMQueryChallengesChecker_InvalidCalldata(t *testing.T) {
	mockService := &mockEVMQueryService{}
	checker := NewEVMQueryChallengesChecker(mockService)

	challenges := []*types.EVMQueryChallenge{
		{
			ContractAddress: "0x1234567890123456789012345678901234567890",
			Calldata:        "not-valid-hex-gg",
			GasLimit:        sdkmath.NewUint(100000),
		},
	}
	approval := createApprovalWithEVMQueryChallenges(challenges)

	errMsg, err := checker.Check(mockContext(), approval, createCollection(), testAddrTo, testAddrFrom, testAddrInitiator, "collection", "approver", nil, nil, "", false)
	require.Error(t, err)
	require.Contains(t, errMsg, "invalid calldata hex")
}

func TestEVMQueryChallengesChecker_InvalidExpectedResult(t *testing.T) {
	mockService := &mockEVMQueryService{
		returnValue: []byte{0x01},
	}
	checker := NewEVMQueryChallengesChecker(mockService)

	challenges := []*types.EVMQueryChallenge{
		{
			ContractAddress: "0x1234567890123456789012345678901234567890",
			Calldata:        "70a08231",
			ExpectedResult:  "not-valid-hex-gg",
			GasLimit:        sdkmath.NewUint(100000),
		},
	}
	approval := createApprovalWithEVMQueryChallenges(challenges)

	errMsg, err := checker.Check(mockContext(), approval, createCollection(), testAddrTo, testAddrFrom, testAddrInitiator, "collection", "approver", nil, nil, "", false)
	require.Error(t, err)
	require.Contains(t, errMsg, "invalid expected result hex")
}

func TestEVMQueryChallengesChecker_EVMCallError(t *testing.T) {
	mockService := &mockEVMQueryService{
		returnError: fmt.Errorf("contract reverted"),
	}
	checker := NewEVMQueryChallengesChecker(mockService)

	challenges := []*types.EVMQueryChallenge{
		{
			ContractAddress: "0x1234567890123456789012345678901234567890",
			Calldata:        "70a08231",
			GasLimit:        sdkmath.NewUint(100000),
		},
	}
	approval := createApprovalWithEVMQueryChallenges(challenges)

	errMsg, err := checker.Check(mockContext(), approval, createCollection(), testAddrTo, testAddrFrom, testAddrInitiator, "collection", "approver", nil, nil, "", false)
	require.Error(t, err)
	require.Contains(t, errMsg, "EVM query challenge 0 failed")
}

func TestEVMQueryChallengesChecker_MultipleChallenges(t *testing.T) {
	// All challenges must pass
	result1, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000001")
	result2, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000002")

	callCount := 0
	mockService := &mockEVMQueryService{}

	checker := &EVMQueryChallengesChecker{
		evmQueryService: &sequentialMockService{
			results: [][]byte{result1, result2},
			callCount: &callCount,
		},
	}

	challenges := []*types.EVMQueryChallenge{
		{
			ContractAddress:    "0x1111111111111111111111111111111111111111",
			Calldata:           "70a08231",
			ExpectedResult:     "0000000000000000000000000000000000000000000000000000000000000001",
			ComparisonOperator: "eq",
			GasLimit:           sdkmath.NewUint(100000),
		},
		{
			ContractAddress:    "0x2222222222222222222222222222222222222222",
			Calldata:           "18160ddd",
			ExpectedResult:     "0000000000000000000000000000000000000000000000000000000000000002",
			ComparisonOperator: "eq",
			GasLimit:           sdkmath.NewUint(100000),
		},
	}
	approval := createApprovalWithEVMQueryChallenges(challenges)

	errMsg, err := checker.Check(mockContext(), approval, createCollection(), testAddrTo, testAddrFrom, testAddrInitiator, "collection", "approver", nil, nil, "", false)
	require.NoError(t, err)
	require.Empty(t, errMsg)
	require.Equal(t, 2, callCount)

	// Reset and test failure on second challenge
	callCount = 0
	checker.evmQueryService = &sequentialMockService{
		results: [][]byte{result1, result1}, // Second result doesn't match expected
		callCount: &callCount,
	}

	errMsg, err = checker.Check(mockContext(), approval, createCollection(), testAddrTo, testAddrFrom, testAddrInitiator, "collection", "approver", nil, nil, "", false)
	require.Error(t, err)
	require.Contains(t, errMsg, "EVM query challenge 1")
	_ = mockService // Avoid unused variable
}

func TestEVMQueryChallengesChecker_NoExpectedResult(t *testing.T) {
	// When no expected result is provided, any non-error result should pass
	mockService := &mockEVMQueryService{
		returnValue: []byte{0x01, 0x02, 0x03},
	}
	checker := NewEVMQueryChallengesChecker(mockService)

	challenges := []*types.EVMQueryChallenge{
		{
			ContractAddress: "0x1234567890123456789012345678901234567890",
			Calldata:        "70a08231",
			ExpectedResult:  "", // No expected result
			GasLimit:        sdkmath.NewUint(100000),
		},
	}
	approval := createApprovalWithEVMQueryChallenges(challenges)

	errMsg, err := checker.Check(mockContext(), approval, createCollection(), testAddrTo, testAddrFrom, testAddrInitiator, "collection", "approver", nil, nil, "", false)
	require.NoError(t, err)
	require.Empty(t, errMsg)
}

func TestEVMQueryChallengesChecker_NilChallenge(t *testing.T) {
	mockService := &mockEVMQueryService{
		returnValue: []byte{0x01},
	}
	checker := NewEVMQueryChallengesChecker(mockService)

	// Mix of nil and valid challenges
	challenges := []*types.EVMQueryChallenge{
		nil, // Should be skipped
		{
			ContractAddress: "0x1234567890123456789012345678901234567890",
			Calldata:        "70a08231",
			GasLimit:        sdkmath.NewUint(100000),
		},
	}
	approval := createApprovalWithEVMQueryChallenges(challenges)

	errMsg, err := checker.Check(mockContext(), approval, createCollection(), testAddrTo, testAddrFrom, testAddrInitiator, "collection", "approver", nil, nil, "", false)
	require.NoError(t, err)
	require.Empty(t, errMsg)
}

func TestAddressToHexPadded(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantLen     int  // Expected length of output when err == nil
		wantErr     bool // Expect an error (invalid address)
	}{
		{"hex_address", "0x1234567890123456789012345678901234567890", 64, false},
		{"short_hex", "0x1234", 0, true},
		{"invalid_address", "invalid", 0, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := addressToHexPadded(tc.input)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.wantLen, len(result), "expected length %d for %s", tc.wantLen, tc.name)
		})
	}
}

// sequentialMockService returns different results for each call
type sequentialMockService struct {
	results   [][]byte
	callCount *int
}

func (m *sequentialMockService) ExecuteEVMQuery(ctx sdk.Context, callerAddress string, contractAddress string, calldata []byte, gasLimit uint64) ([]byte, error) {
	idx := *m.callCount
	*m.callCount++
	if idx < len(m.results) {
		return m.results[idx], nil
	}
	return nil, fmt.Errorf("no more results")
}
