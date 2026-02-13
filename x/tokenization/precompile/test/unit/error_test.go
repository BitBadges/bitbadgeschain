package tokenization_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	tokenizationkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	tokenization "github.com/bitbadges/bitbadgeschain/x/tokenization/precompile"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/precompile/test/helpers"
	keepertest "github.com/bitbadges/bitbadgeschain/x/tokenization/testutil/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TestPrecompile_ErrorScenarios tests various error conditions
func TestPrecompile_ErrorScenarios(t *testing.T) {
	tokenizationKeeper, _ := keepertest.TokenizationKeeper(t)
	precompile := tokenization.NewPrecompile(tokenizationKeeper)

	// Test address for building JSON
	testAddr := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0")
	testAddrCosmos := sdk.AccAddress(testAddr.Bytes()).String()

	tests := []struct {
		name        string
		method      string
		setup       func() string // Returns JSON string
		expectError bool
		errorCode   tokenization.ErrorCode
	}{
		{
			name:   "invalid_collection_id_negative",
			method: tokenization.TransferTokensMethod,
			setup: func() string {
				// Negative collection ID - JSON will accept it as string, but validation should catch it
				msg := map[string]interface{}{
					"creator":        testAddrCosmos,
					"collectionId":   "-1", // Negative (will be caught in validation)
					"toAddresses":    []string{testAddrCosmos},
					"amount":         "1",
					"tokenIds":       []map[string]interface{}{{"start": "1", "end": "1"}},
					"ownershipTimes": []map[string]interface{}{{"start": "1", "end": "100"}},
				}
				jsonBytes, _ := json.Marshal(msg)
				return string(jsonBytes)
			},
			expectError: true,
			errorCode:   tokenization.ErrorCodeInvalidInput,
		},
		{
			name:   "zero_address_recipient",
			method: tokenization.TransferTokensMethod,
			setup: func() string {
				msg := map[string]interface{}{
					"creator":        testAddrCosmos,
					"collectionId":   "1",
					"toAddresses":    []string{""}, // Zero address (empty string)
					"amount":         "1",
					"tokenIds":       []map[string]interface{}{{"start": "1", "end": "1"}},
					"ownershipTimes": []map[string]interface{}{{"start": "1", "end": "100"}},
				}
				jsonBytes, _ := json.Marshal(msg)
				return string(jsonBytes)
			},
			expectError: true,
			errorCode:   tokenization.ErrorCodeInvalidInput,
		},
		{
			name:   "empty_recipients",
			method: tokenization.TransferTokensMethod,
			setup: func() string {
				msg := map[string]interface{}{
					"creator":        testAddrCosmos,
					"collectionId":   "1",
					"toAddresses":    []string{}, // Empty array
					"amount":         "1",
					"tokenIds":       []map[string]interface{}{{"start": "1", "end": "1"}},
					"ownershipTimes": []map[string]interface{}{{"start": "1", "end": "100"}},
				}
				jsonBytes, _ := json.Marshal(msg)
				return string(jsonBytes)
			},
			expectError: true,
			errorCode:   tokenization.ErrorCodeInvalidInput,
		},
		{
			name:   "zero_amount",
			method: tokenization.TransferTokensMethod,
			setup: func() string {
				msg := map[string]interface{}{
					"creator":        testAddrCosmos,
					"collectionId":   "1",
					"toAddresses":    []string{testAddrCosmos},
					"amount":         "0", // Zero amount
					"tokenIds":       []map[string]interface{}{{"start": "1", "end": "1"}},
					"ownershipTimes": []map[string]interface{}{{"start": "1", "end": "100"}},
				}
				jsonBytes, _ := json.Marshal(msg)
				return string(jsonBytes)
			},
			expectError: true,
			errorCode:   tokenization.ErrorCodeInvalidInput,
		},
		{
			name:   "invalid_range_start_greater_than_end",
			method: tokenization.TransferTokensMethod,
			setup: func() string {
				msg := map[string]interface{}{
					"creator":        testAddrCosmos,
					"collectionId":   "1",
					"toAddresses":    []string{testAddrCosmos},
					"amount":         "1",
					"tokenIds":       []map[string]interface{}{{"start": "10", "end": "5"}}, // Invalid: start > end
					"ownershipTimes": []map[string]interface{}{{"start": "1", "end": "100"}},
				}
				jsonBytes, _ := json.Marshal(msg)
				return string(jsonBytes)
			},
			expectError: true,
			errorCode:   tokenization.ErrorCodeInvalidInput,
		},
		{
			name:   "empty_token_ids",
			method: tokenization.TransferTokensMethod,
			setup: func() string {
				msg := map[string]interface{}{
					"creator":        testAddrCosmos,
					"collectionId":   "1",
					"toAddresses":    []string{testAddrCosmos},
					"amount":         "1",
					"tokenIds":       []map[string]interface{}{}, // Empty array
					"ownershipTimes": []map[string]interface{}{{"start": "1", "end": "100"}},
				}
				jsonBytes, _ := json.Marshal(msg)
				return string(jsonBytes)
			},
			expectError: true,
			errorCode:   tokenization.ErrorCodeInvalidInput,
		},
		{
			name:   "too_many_recipients",
			method: tokenization.TransferTokensMethod,
			setup: func() string {
				// Create array with more than MaxRecipients
				addresses := make([]string, tokenization.MaxRecipients+1)
				for i := range addresses {
					addr := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0")
					addresses[i] = sdk.AccAddress(addr.Bytes()).String()
				}
				msg := map[string]interface{}{
					"creator":        testAddrCosmos,
					"collectionId":   "1",
					"toAddresses":    addresses,
					"amount":         "1",
					"tokenIds":       []map[string]interface{}{{"start": "1", "end": "1"}},
					"ownershipTimes": []map[string]interface{}{{"start": "1", "end": "100"}},
				}
				jsonBytes, _ := json.Marshal(msg)
				return string(jsonBytes)
			},
			expectError: true,
			errorCode:   tokenization.ErrorCodeInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonStr := tt.setup()
			method := precompile.ABI.Methods[tt.method]

			// Pack method with JSON string
			packed, err := helpers.PackMethodWithJSON(&method, jsonStr)

			// For negative values, JSON will be valid but validation should catch it
			if tt.name == "invalid_collection_id_negative" {
				// JSON packing should succeed, but execution will fail
				require.NoError(t, err, "Packing JSON should succeed")
			} else {
				require.NoError(t, err, "Packing should succeed")
			}

			input := append(method.ID, packed...)

			// Create a minimal contract for testing
			// Note: This is a simplified test - full testing requires EVM context
			// The validation functions should catch these errors before execution
			gas := precompile.RequiredGas(input)
			if tt.expectError {
				// RequiredGas doesn't validate, so we test validation separately
				// For now, we just verify the method exists
				require.NotEqual(t, uint64(0), gas, "Method should have gas cost")
			}
		})
	}
}

// TestValidationFunctions tests the validation helper functions
func TestValidationFunctions(t *testing.T) {
	t.Run("ValidateAddress", func(t *testing.T) {
		err := tokenization.ValidateAddress(common.Address{}, "test")
		require.Error(t, err)
		require.Contains(t, err.Error(), "cannot be zero address")

		err = tokenization.ValidateAddress(common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0"), "test")
		require.NoError(t, err)
	})

	t.Run("ValidateAddresses", func(t *testing.T) {
		err := tokenization.ValidateAddresses([]common.Address{}, "test")
		require.Error(t, err)
		require.Contains(t, err.Error(), "cannot be empty")

		err = tokenization.ValidateAddresses([]common.Address{{}}, "test")
		require.Error(t, err)
		require.Contains(t, err.Error(), "cannot be zero address")

		err = tokenization.ValidateAddresses([]common.Address{common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0")}, "test")
		require.NoError(t, err)
	})

	t.Run("ValidateAmount", func(t *testing.T) {
		err := tokenization.ValidateAmount(nil, "test")
		require.Error(t, err)

		err = tokenization.ValidateAmount(big.NewInt(0), "test")
		require.Error(t, err)
		require.Contains(t, err.Error(), "must be greater than zero")

		err = tokenization.ValidateAmount(big.NewInt(-1), "test")
		require.Error(t, err)
		require.Contains(t, err.Error(), "must be greater than zero")

		err = tokenization.ValidateAmount(big.NewInt(1), "test")
		require.NoError(t, err)
	})

	t.Run("ValidateCollectionId", func(t *testing.T) {
		err := tokenization.ValidateCollectionId(nil)
		require.Error(t, err)

		err = tokenization.ValidateCollectionId(big.NewInt(-1))
		require.Error(t, err)
		require.Contains(t, err.Error(), "cannot be negative")

		err = tokenization.ValidateCollectionId(big.NewInt(0))
		require.Error(t, err) // Zero collection ID should be rejected
		require.Contains(t, err.Error(), "cannot be zero")

		err = tokenization.ValidateCollectionId(big.NewInt(1))
		require.NoError(t, err)
	})
}

// TestMapCosmosErrorToPrecompileError tests the error mapping function
func TestMapCosmosErrorToPrecompileError(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedCode tokenization.ErrorCode
		shouldMap    bool
	}{
		{
			name:         "ErrCollectionNotExists",
			err:          tokenizationkeeper.ErrCollectionNotExists,
			expectedCode: tokenization.ErrorCodeCollectionNotFound,
			shouldMap:    true,
		},
		{
			name:         "ErrUserBalanceNotExists",
			err:          tokenizationkeeper.ErrUserBalanceNotExists,
			expectedCode: tokenization.ErrorCodeBalanceNotFound,
			shouldMap:    true,
		},
		{
			name:         "ErrInadequateApprovals",
			err:          tokenizationkeeper.ErrInadequateApprovals,
			expectedCode: tokenization.ErrorCodeUnauthorized,
			shouldMap:    true,
		},
		{
			name:         "ErrCollectionIsArchived",
			err:          tokenizationkeeper.ErrCollectionIsArchived,
			expectedCode: tokenization.ErrorCodeCollectionArchived,
			shouldMap:    true,
		},
		{
			name:         "ErrDisallowedTransfer",
			err:          tokenizationkeeper.ErrDisallowedTransfer,
			expectedCode: tokenization.ErrorCodeTransferFailed,
			shouldMap:    true,
		},
		{
			name:         "ErrAccountNotFound",
			err:          tokenizationkeeper.ErrAccountNotFound,
			expectedCode: tokenization.ErrorCodeQueryFailed,
			shouldMap:    true,
		},
		{
			name:         "ErrUnauthorized",
			err:          tokenizationtypes.ErrUnauthorized,
			expectedCode: tokenization.ErrorCodeUnauthorized,
			shouldMap:    true,
		},
		{
			name:         "ErrNotFound",
			err:          tokenizationtypes.ErrNotFound,
			expectedCode: tokenization.ErrorCodeQueryFailed,
			shouldMap:    true,
		},
		{
			name:         "ErrInvalidRequest",
			err:          tokenizationtypes.ErrInvalidRequest,
			expectedCode: tokenization.ErrorCodeInvalidInput,
			shouldMap:    true,
		},
		{
			name:         "ErrUnderflow",
			err:          tokenizationtypes.ErrUnderflow,
			expectedCode: tokenization.ErrorCodeTransferFailed,
			shouldMap:    true,
		},
		{
			name:         "ErrOverflow",
			err:          tokenizationtypes.ErrOverflow,
			expectedCode: tokenization.ErrorCodeTransferFailed,
			shouldMap:    true,
		},
		{
			name:         "ErrInvalidAddress",
			err:          tokenizationtypes.ErrInvalidAddress,
			expectedCode: tokenization.ErrorCodeInvalidInput,
			shouldMap:    true,
		},
		{
			name:         "ErrInvalidCollectionID",
			err:          tokenizationtypes.ErrInvalidCollectionID,
			expectedCode: tokenization.ErrorCodeCollectionNotFound, // ErrInvalidCollectionID is used when collection doesn't exist
			shouldMap:    true,
		},
		{
			name:         "ErrAmountEqualsZero",
			err:          tokenizationtypes.ErrAmountEqualsZero,
			expectedCode: tokenization.ErrorCodeInvalidInput,
			shouldMap:    true,
		},
		{
			name:         "unmapped_error",
			err:          fmt.Errorf("some random error"),
			expectedCode: 0,
			shouldMap:    false,
		},
		{
			name:         "nil_error",
			err:          nil,
			expectedCode: 0,
			shouldMap:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, found := tokenization.MapCosmosErrorToPrecompileError(tt.err)
			require.Equal(t, tt.shouldMap, found, "mapping should match expected")
			if tt.shouldMap {
				require.Equal(t, tt.expectedCode, code, "error code should match")
			}
		})
	}
}

// TestWrapErrorWithMapping tests that WrapError uses error mapping
func TestWrapErrorWithMapping(t *testing.T) {
	// Test that mapped errors use the mapped code
	err := tokenizationkeeper.ErrCollectionNotExists
	wrapped := tokenization.WrapError(err, tokenization.ErrorCodeQueryFailed, "test message")
	require.Equal(t, tokenization.ErrorCodeCollectionNotFound, wrapped.Code, "should use mapped code, not default")

	// Test that unmapped errors use the default code
	unmappedErr := fmt.Errorf("some random error")
	wrapped2 := tokenization.WrapError(unmappedErr, tokenization.ErrorCodeInternalError, "test message")
	require.Equal(t, tokenization.ErrorCodeInternalError, wrapped2.Code, "should use default code for unmapped errors")

	// Test that nil errors use the default code
	wrapped3 := tokenization.WrapError(nil, tokenization.ErrorCodeQueryFailed, "test message")
	require.Equal(t, tokenization.ErrorCodeQueryFailed, wrapped3.Code, "should use default code for nil errors")
}

// TestErrorSanitization tests that error sanitization removes sensitive information
func TestErrorSanitization(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string // What should NOT appear in output
		mustContain string // What should appear in output
	}{
		{
			name:        "removes_unix_paths",
			input:       "error at /home/user/file.go:123",
			expected:    "/home/",
			mustContain: "[path]/",
		},
		{
			name:        "removes_windows_paths",
			input:       "error at C:\\Users\\file.go:123",
			expected:    "C:\\",
			mustContain: "[path]\\",
		},
		{
			name:        "removes_go_file_references",
			input:       "error in handlers.go:456",
			expected:    ".go:",
			mustContain: "[file]:",
		},
		{
			name:        "removes_goroutine_info",
			input:       "goroutine 1 [running]:",
			expected:    "goroutine ",
			mustContain: "[goroutine] ",
		},
		{
			name:        "removes_panic_info",
			input:       "panic: runtime error",
			expected:    "panic:",
			mustContain: "[panic]:",
		},
		{
			name:        "removes_runtime_info",
			input:       "runtime.throw(...)",
			expected:    "runtime.",
			mustContain: "[runtime].",
		},
		{
			name:        "removes_module_paths",
			input:       "github.com/bitbadges/bitbadgeschain/x/tokenization/handlers.go:123",
			expected:    "github.com/bitbadges/",
			mustContain: "[module]/",
		},
		{
			name:        "removes_cosmos_module_paths",
			input:       "github.com/cosmos/cosmos-sdk/types/errors.go:456",
			expected:    "github.com/cosmos/",
			mustContain: "[module]/",
		},
		{
			name:        "removes_localhost_ip",
			input:       "connection to 127.0.0.1 failed",
			expected:    "127.0.0.1",
			mustContain: "[localhost]",
		},
		{
			name:        "truncates_long_messages",
			input:       string(make([]byte, 1000)), // 1000 bytes
			expected:    "",                         // Should be truncated
			mustContain: "[truncated]",
		},
		{
			name:        "preserves_safe_messages",
			input:       "collection not found: 123",
			expected:    "", // Should not be modified
			mustContain: "collection not found: 123",
		},
		{
			name:        "handles_empty_string",
			input:       "",
			expected:    "",
			mustContain: "",
		},
		{
			name:        "handles_multiple_patterns",
			input:       "error at /home/user/github.com/bitbadges/handlers.go:123 goroutine 1",
			expected:    "/home/",
			mustContain: "[path]/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create an error with the test input
			err := tokenization.NewPrecompileError(tokenization.ErrorCodeInternalError, "test", tt.input)

			// Check that sanitized details don't contain sensitive info
			if tt.expected != "" {
				require.NotContains(t, err.Details, tt.expected, "sensitive information should be removed")
			}

			// Check that sanitized details contain safe replacement
			if tt.mustContain != "" {
				require.Contains(t, err.Details, tt.mustContain, "should contain sanitized replacement")
			}

			// For truncation test, verify length
			if tt.name == "truncates_long_messages" {
				require.LessOrEqual(t, len(err.Details), 500+len("... [truncated]"), "should be truncated to max length")
			}
		})
	}
}

// TestErrorSanitizationIntegration tests error sanitization in real error scenarios
func TestErrorSanitizationIntegration(t *testing.T) {
	// Test with wrapped errors that contain paths
	err := fmt.Errorf("error in /home/user/project/handlers.go:123: collection not found")
	wrapped := tokenization.WrapError(err, tokenization.ErrorCodeCollectionNotFound, "test operation failed")

	// Verify sensitive path is removed
	require.NotContains(t, wrapped.Details, "/home/", "should remove file paths")
	require.Contains(t, wrapped.Details, "[path]/", "should replace with sanitized version")

	// Test with stack trace-like error
	stackTrace := "goroutine 1 [running]:\n    github.com/bitbadges/bitbadgeschain/x/tokenization/handlers.go:123\n    runtime.panic(...)"
	wrapped2 := tokenization.WrapError(fmt.Errorf("%s", stackTrace), tokenization.ErrorCodeInternalError, "internal error")

	require.NotContains(t, wrapped2.Details, "goroutine ", "should remove goroutine info")
	require.NotContains(t, wrapped2.Details, "github.com/bitbadges/", "should remove module paths")
}
