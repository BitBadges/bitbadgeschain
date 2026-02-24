package types

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func TestValidateEVMQueryChallenges_Valid(t *testing.T) {
	tests := []struct {
		name       string
		challenges []*EVMQueryChallenge
	}{
		{
			name:       "empty challenges",
			challenges: []*EVMQueryChallenge{},
		},
		{
			name:       "nil challenges",
			challenges: nil,
		},
		{
			name: "single valid challenge",
			challenges: []*EVMQueryChallenge{
				{
					ContractAddress:    "0x1234567890123456789012345678901234567890",
					Calldata:           "70a08231000000000000000000000000abcdef",
					ExpectedResult:     "0000000000000000000000000000000000000000000000000000000000000001",
					ComparisonOperator: "eq",
					GasLimit:           sdkmath.NewUint(100000),
				},
			},
		},
		{
			name: "multiple valid challenges",
			challenges: []*EVMQueryChallenge{
				{
					ContractAddress:    "0x1234567890123456789012345678901234567890",
					Calldata:           "70a08231",
					ComparisonOperator: "gte",
					GasLimit:           sdkmath.NewUint(100000),
				},
				{
					ContractAddress:    "0xabcdef1234567890abcdef1234567890abcdef12",
					Calldata:           "18160ddd",
					ComparisonOperator: "lt",
					GasLimit:           sdkmath.NewUint(50000),
				},
			},
		},
		{
			name: "all comparison operators",
			challenges: []*EVMQueryChallenge{
				{ContractAddress: "0x1234567890123456789012345678901234567890", Calldata: "70a0823100000000", ComparisonOperator: "eq", GasLimit: sdkmath.NewUint(100000)},
				{ContractAddress: "0x1234567890123456789012345678901234567890", Calldata: "18160ddd00000000", ComparisonOperator: "ne", GasLimit: sdkmath.NewUint(100000)},
				{ContractAddress: "0x1234567890123456789012345678901234567890", Calldata: "70a0823100000001", ComparisonOperator: "gt", GasLimit: sdkmath.NewUint(100000)},
				{ContractAddress: "0x1234567890123456789012345678901234567890", Calldata: "70a0823100000002", ComparisonOperator: "gte", GasLimit: sdkmath.NewUint(100000)},
				{ContractAddress: "0x1234567890123456789012345678901234567890", Calldata: "70a0823100000003", ComparisonOperator: "lt", GasLimit: sdkmath.NewUint(100000)},
				{ContractAddress: "0x1234567890123456789012345678901234567890", Calldata: "70a0823100000004", ComparisonOperator: "lte", GasLimit: sdkmath.NewUint(100000)},
			},
		},
		{
			name: "empty comparison operator (defaults to eq)",
			challenges: []*EVMQueryChallenge{
				{
					ContractAddress:    "0x1234567890123456789012345678901234567890",
					Calldata:           "70a08231",
					ComparisonOperator: "",
					GasLimit:           sdkmath.NewUint(100000),
				},
			},
		},
		{
			name: "gas limit at max",
			challenges: []*EVMQueryChallenge{
				{
					ContractAddress: "0x1234567890123456789012345678901234567890",
					Calldata:        "70a08231",
					GasLimit:        sdkmath.NewUint(500000),
				},
			},
		},
		{
			name: "zero gas limit (will use default)",
			challenges: []*EVMQueryChallenge{
				{
					ContractAddress: "0x1234567890123456789012345678901234567890",
					Calldata:        "70a08231",
					GasLimit:        sdkmath.NewUint(0),
				},
			},
		},
		{
			name: "nil challenge in list should be skipped",
			challenges: []*EVMQueryChallenge{
				nil,
				{
					ContractAddress: "0x1234567890123456789012345678901234567890",
					Calldata:        "70a08231",
					GasLimit:        sdkmath.NewUint(100000),
				},
			},
		},
		{
			name: "challenge with placeholders",
			challenges: []*EVMQueryChallenge{
				{
					ContractAddress: "0x1234567890123456789012345678901234567890",
					Calldata:        "70a08231$initiator",
					GasLimit:        sdkmath.NewUint(100000),
				},
			},
		},
		{
			name: "challenge with $collectionId placeholder",
			challenges: []*EVMQueryChallenge{
				{
					ContractAddress: "0x1234567890123456789012345678901234567890",
					Calldata:        "a1b2c3d4$collectionId0000000000000000000000000000000000000000000000000000000000000003",
					GasLimit:        sdkmath.NewUint(100000),
				},
			},
		},
		{
			name: "challenge with URI and custom data",
			challenges: []*EVMQueryChallenge{
				{
					ContractAddress: "0x1234567890123456789012345678901234567890",
					Calldata:        "70a08231",
					Uri:             "https://example.com/challenge",
					CustomData:      `{"description": "Balance check"}`,
					GasLimit:        sdkmath.NewUint(100000),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEVMQueryChallenges(tt.challenges)
			require.NoError(t, err, "expected no error for valid challenges: %s", tt.name)
		})
	}
}

func TestValidateEVMQueryChallenges_Invalid(t *testing.T) {
	tests := []struct {
		name          string
		challenges    []*EVMQueryChallenge
		expectedError string
	}{
		{
			name: "missing contract address",
			challenges: []*EVMQueryChallenge{
				{
					ContractAddress: "",
					Calldata:        "70a08231",
					GasLimit:        sdkmath.NewUint(100000),
				},
			},
			expectedError: "contract address required",
		},
		{
			name: "missing calldata",
			challenges: []*EVMQueryChallenge{
				{
					ContractAddress: "0x1234567890123456789012345678901234567890",
					Calldata:        "",
					GasLimit:        sdkmath.NewUint(100000),
				},
			},
			expectedError: "calldata required",
		},
		{
			name: "invalid comparison operator",
			challenges: []*EVMQueryChallenge{
				{
					ContractAddress:    "0x1234567890123456789012345678901234567890",
					Calldata:           "70a08231",
					ComparisonOperator: "invalid",
					GasLimit:           sdkmath.NewUint(100000),
				},
			},
			expectedError: "invalid comparison operator",
		},
		{
			name: "gas limit exceeds maximum",
			challenges: []*EVMQueryChallenge{
				{
					ContractAddress: "0x1234567890123456789012345678901234567890",
					Calldata:        "70a08231",
					GasLimit:        sdkmath.NewUint(500001),
				},
			},
			expectedError: "gas limit exceeds maximum",
		},
		{
			name:          "too many challenges",
			challenges:    makeManyEVMQueryChallenges(11),
			expectedError: "too many EVM query challenges",
		},
		{
			name: "second challenge invalid",
			challenges: []*EVMQueryChallenge{
				{
					ContractAddress: "0x1234567890123456789012345678901234567890",
					Calldata:        "70a0823100000000",
					GasLimit:        sdkmath.NewUint(100000),
				},
				{
					ContractAddress: "", // Invalid
					Calldata:        "70a0823100000000",
					GasLimit:        sdkmath.NewUint(100000),
				},
			},
			expectedError: "EVM query challenge 1: contract address required",
		},
		{
			name: "various invalid operators",
			challenges: []*EVMQueryChallenge{
				{
					ContractAddress:    "0x1234567890123456789012345678901234567890",
					Calldata:           "70a0823100000000",
					ComparisonOperator: "equals", // Should be "eq"
					GasLimit:           sdkmath.NewUint(100000),
				},
			},
			expectedError: "invalid comparison operator equals",
		},
		{
			name: "invalid contract address format - 0x too short",
			challenges: []*EVMQueryChallenge{
				{
					ContractAddress: "0x123456789012345678901234567890123456789",
					Calldata:        "70a0823100000000",
					GasLimit:        sdkmath.NewUint(100000),
				},
			},
			expectedError: "invalid contract address format",
		},
		{
			name: "invalid contract address format - non-hex",
			challenges: []*EVMQueryChallenge{
				{
					ContractAddress: "0x123456789012345678901234567890123456789g",
					Calldata:        "70a0823100000000",
					GasLimit:        sdkmath.NewUint(100000),
				},
			},
			expectedError: "invalid contract address format",
		},
		{
			name: "calldata too short for selector",
			challenges: []*EVMQueryChallenge{
				{
					ContractAddress: "0x1234567890123456789012345678901234567890",
					Calldata:        "70a082",
					GasLimit:        sdkmath.NewUint(100000),
				},
			},
			expectedError: "at least 4 bytes",
		},
		{
			name: "calldata odd hex length",
			challenges: []*EVMQueryChallenge{
				{
					ContractAddress: "0x1234567890123456789012345678901234567890",
					Calldata:        "70a0823",
					GasLimit:        sdkmath.NewUint(100000),
				},
			},
			expectedError: "even hex length",
		},
		{
			name: "calldata non-hex character",
			challenges: []*EVMQueryChallenge{
				{
					ContractAddress: "0x1234567890123456789012345678901234567890",
					Calldata:        "70a08231zz000000",
					GasLimit:        sdkmath.NewUint(100000),
				},
			},
			expectedError: "non-hex",
		},
		{
			name: "expected result odd hex length",
			challenges: []*EVMQueryChallenge{
				{
					ContractAddress: "0x1234567890123456789012345678901234567890",
					Calldata:        "70a0823100000000",
					ExpectedResult:  "000000000000000000000000000000000000000000000000000000000000001",
					GasLimit:        sdkmath.NewUint(100000),
				},
			},
			expectedError: "expected result must have even hex length",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEVMQueryChallenges(tt.challenges)
			require.Error(t, err, "expected error for invalid challenges: %s", tt.name)
			require.Contains(t, err.Error(), tt.expectedError, "error should contain expected message for: %s", tt.name)
		})
	}
}

func TestValidateEVMQueryChallenges_MaxChallenges(t *testing.T) {
	// Exactly 10 should be valid
	challenges := makeManyEVMQueryChallenges(10)
	err := ValidateEVMQueryChallenges(challenges)
	require.NoError(t, err, "10 challenges should be valid")

	// 11 should be invalid
	challenges = makeManyEVMQueryChallenges(11)
	err = ValidateEVMQueryChallenges(challenges)
	require.Error(t, err, "11 challenges should be invalid")
	require.Contains(t, err.Error(), "too many EVM query challenges")
}

func TestValidateEVMQueryChallenges_GasLimitBoundary(t *testing.T) {
	tests := []struct {
		name        string
		gasLimit    uint64
		expectError bool
	}{
		{"gas_limit_0", 0, false},
		{"gas_limit_1", 1, false},
		{"gas_limit_100000", 100000, false},
		{"gas_limit_499999", 499999, false},
		{"gas_limit_500000", 500000, false},
		{"gas_limit_500001", 500001, true},
		{"gas_limit_1000000", 1000000, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			challenges := []*EVMQueryChallenge{
				{
					ContractAddress: "0x1234567890123456789012345678901234567890",
					Calldata:        "70a08231",
					GasLimit:        sdkmath.NewUint(tt.gasLimit),
				},
			}
			err := ValidateEVMQueryChallenges(challenges)
			if tt.expectError {
				require.Error(t, err, "expected error for gas limit %d", tt.gasLimit)
				require.Contains(t, err.Error(), "gas limit exceeds maximum")
			} else {
				require.NoError(t, err, "expected no error for gas limit %d", tt.gasLimit)
			}
		})
	}
}

func TestValidateEVMQueryChallenges_AllOperators(t *testing.T) {
	validOperators := []string{"", "eq", "ne", "gt", "gte", "lt", "lte"}
	invalidOperators := []string{"equals", "not_equals", "greater", "less", "EQ", "NE", "GT", "GTE", "LT", "LTE"}

	for _, op := range validOperators {
		t.Run("valid_"+op, func(t *testing.T) {
			challenges := []*EVMQueryChallenge{
				{
					ContractAddress:    "0x1234567890123456789012345678901234567890",
					Calldata:           "70a08231",
					ComparisonOperator: op,
					GasLimit:           sdkmath.NewUint(100000),
				},
			}
			err := ValidateEVMQueryChallenges(challenges)
			require.NoError(t, err, "operator '%s' should be valid", op)
		})
	}

	for _, op := range invalidOperators {
		t.Run("invalid_"+op, func(t *testing.T) {
			challenges := []*EVMQueryChallenge{
				{
					ContractAddress:    "0x1234567890123456789012345678901234567890",
					Calldata:           "70a08231",
					ComparisonOperator: op,
					GasLimit:           sdkmath.NewUint(100000),
				},
			}
			err := ValidateEVMQueryChallenges(challenges)
			require.Error(t, err, "operator '%s' should be invalid", op)
			require.Contains(t, err.Error(), "invalid comparison operator")
		})
	}
}

// Helper function to create many EVM query challenges
func makeManyEVMQueryChallenges(count int) []*EVMQueryChallenge {
	challenges := make([]*EVMQueryChallenge, count)
	for i := 0; i < count; i++ {
		challenges[i] = &EVMQueryChallenge{
			ContractAddress: "0x1234567890123456789012345678901234567890",
			Calldata:        "70a08231",
			GasLimit:        sdkmath.NewUint(100000),
		}
	}
	return challenges
}
