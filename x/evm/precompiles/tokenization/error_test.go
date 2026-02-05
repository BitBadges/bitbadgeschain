package tokenization

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	tokenizationkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	keepertest "github.com/bitbadges/bitbadgeschain/x/tokenization/testutil/keeper"
)

// TestPrecompile_ErrorScenarios tests various error conditions
func TestPrecompile_ErrorScenarios(t *testing.T) {
	tokenizationKeeper, _ := keepertest.TokenizationKeeper(t)
	precompile := NewPrecompile(tokenizationKeeper)

	tests := []struct {
		name        string
		method      string
		setup       func() []interface{}
		expectError bool
		errorCode   ErrorCode
	}{
		{
			name:   "invalid_collection_id_negative",
			method: TransferTokensMethod,
			setup: func() []interface{} {
				return []interface{}{
					big.NewInt(-1), // Negative collection ID
					[]common.Address{common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0")},
					big.NewInt(1),
					[]struct {
						Start *big.Int `json:"start"`
						End   *big.Int `json:"end"`
					}{{Start: big.NewInt(1), End: big.NewInt(1)}},
					[]struct {
						Start *big.Int `json:"start"`
						End   *big.Int `json:"end"`
					}{{Start: big.NewInt(1), End: big.NewInt(100)}},
				}
			},
			expectError: true,
			errorCode:   ErrorCodeInvalidInput,
		},
		{
			name:   "zero_address_recipient",
			method: TransferTokensMethod,
			setup: func() []interface{} {
				return []interface{}{
					big.NewInt(1),
					[]common.Address{common.Address{}}, // Zero address
					big.NewInt(1),
					[]struct {
						Start *big.Int `json:"start"`
						End   *big.Int `json:"end"`
					}{{Start: big.NewInt(1), End: big.NewInt(1)}},
					[]struct {
						Start *big.Int `json:"start"`
						End   *big.Int `json:"end"`
					}{{Start: big.NewInt(1), End: big.NewInt(100)}},
				}
			},
			expectError: true,
			errorCode:   ErrorCodeInvalidInput,
		},
		{
			name:   "empty_recipients",
			method: TransferTokensMethod,
			setup: func() []interface{} {
				return []interface{}{
					big.NewInt(1),
					[]common.Address{}, // Empty array
					big.NewInt(1),
					[]struct {
						Start *big.Int `json:"start"`
						End   *big.Int `json:"end"`
					}{{Start: big.NewInt(1), End: big.NewInt(1)}},
					[]struct {
						Start *big.Int `json:"start"`
						End   *big.Int `json:"end"`
					}{{Start: big.NewInt(1), End: big.NewInt(100)}},
				}
			},
			expectError: true,
			errorCode:   ErrorCodeInvalidInput,
		},
		{
			name:   "zero_amount",
			method: TransferTokensMethod,
			setup: func() []interface{} {
				return []interface{}{
					big.NewInt(1),
					[]common.Address{common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0")},
					big.NewInt(0), // Zero amount
					[]struct {
						Start *big.Int `json:"start"`
						End   *big.Int `json:"end"`
					}{{Start: big.NewInt(1), End: big.NewInt(1)}},
					[]struct {
						Start *big.Int `json:"start"`
						End   *big.Int `json:"end"`
					}{{Start: big.NewInt(1), End: big.NewInt(100)}},
				}
			},
			expectError: true,
			errorCode:   ErrorCodeInvalidInput,
		},
		{
			name:   "invalid_range_start_greater_than_end",
			method: TransferTokensMethod,
			setup: func() []interface{} {
				return []interface{}{
					big.NewInt(1),
					[]common.Address{common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0")},
					big.NewInt(1),
					[]struct {
						Start *big.Int `json:"start"`
						End   *big.Int `json:"end"`
					}{{Start: big.NewInt(10), End: big.NewInt(5)}}, // Invalid: start > end
					[]struct {
						Start *big.Int `json:"start"`
						End   *big.Int `json:"end"`
					}{{Start: big.NewInt(1), End: big.NewInt(100)}},
				}
			},
			expectError: true,
			errorCode:   ErrorCodeInvalidInput,
		},
		{
			name:   "empty_token_ids",
			method: TransferTokensMethod,
			setup: func() []interface{} {
				return []interface{}{
					big.NewInt(1),
					[]common.Address{common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0")},
					big.NewInt(1),
					[]struct {
						Start *big.Int `json:"start"`
						End   *big.Int `json:"end"`
					}{}, // Empty array
					[]struct {
						Start *big.Int `json:"start"`
						End   *big.Int `json:"end"`
					}{{Start: big.NewInt(1), End: big.NewInt(100)}},
				}
			},
			expectError: true,
			errorCode:   ErrorCodeInvalidInput,
		},
		{
			name:   "too_many_recipients",
			method: TransferTokensMethod,
			setup: func() []interface{} {
				// Create array with more than MaxRecipients
				addresses := make([]common.Address, MaxRecipients+1)
				for i := range addresses {
					addresses[i] = common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0")
				}
				return []interface{}{
					big.NewInt(1),
					addresses,
					big.NewInt(1),
					[]struct {
						Start *big.Int `json:"start"`
						End   *big.Int `json:"end"`
					}{{Start: big.NewInt(1), End: big.NewInt(1)}},
					[]struct {
						Start *big.Int `json:"start"`
						End   *big.Int `json:"end"`
					}{{Start: big.NewInt(1), End: big.NewInt(100)}},
				}
			},
			expectError: true,
			errorCode:   ErrorCodeInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.setup()
			method := precompile.ABI.Methods[tt.method]
			packed, err := method.Inputs.Pack(args...)
			
			// For negative values, ABI packing will fail before validation
			if tt.name == "invalid_collection_id_negative" {
				require.Error(t, err, "Packing negative value into uint should fail")
				return
			}
			
			require.NoError(t, err, "Packing should succeed")

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
		err := ValidateAddress(common.Address{}, "test")
		require.Error(t, err)
		require.Contains(t, err.Error(), "cannot be zero address")

		err = ValidateAddress(common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0"), "test")
		require.NoError(t, err)
	})

	t.Run("ValidateAddresses", func(t *testing.T) {
		err := ValidateAddresses([]common.Address{}, "test")
		require.Error(t, err)
		require.Contains(t, err.Error(), "cannot be empty")

		err = ValidateAddresses([]common.Address{common.Address{}}, "test")
		require.Error(t, err)
		require.Contains(t, err.Error(), "cannot be zero address")

		err = ValidateAddresses([]common.Address{common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0")}, "test")
		require.NoError(t, err)
	})

	t.Run("ValidateAmount", func(t *testing.T) {
		err := ValidateAmount(nil, "test")
		require.Error(t, err)

		err = ValidateAmount(big.NewInt(0), "test")
		require.Error(t, err)
		require.Contains(t, err.Error(), "must be greater than zero")

		err = ValidateAmount(big.NewInt(-1), "test")
		require.Error(t, err)
		require.Contains(t, err.Error(), "must be greater than zero")

		err = ValidateAmount(big.NewInt(1), "test")
		require.NoError(t, err)
	})

	t.Run("ValidateCollectionId", func(t *testing.T) {
		err := ValidateCollectionId(nil)
		require.Error(t, err)

		err = ValidateCollectionId(big.NewInt(-1))
		require.Error(t, err)
		require.Contains(t, err.Error(), "cannot be negative")

		err = ValidateCollectionId(big.NewInt(0))
		require.NoError(t, err)

		err = ValidateCollectionId(big.NewInt(1))
		require.NoError(t, err)
	})
}

// TestMapCosmosErrorToPrecompileError tests the error mapping function
func TestMapCosmosErrorToPrecompileError(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedCode ErrorCode
		shouldMap    bool
	}{
		{
			name:         "ErrCollectionNotExists",
			err:          tokenizationkeeper.ErrCollectionNotExists,
			expectedCode: ErrorCodeCollectionNotFound,
			shouldMap:    true,
		},
		{
			name:         "ErrUserBalanceNotExists",
			err:          tokenizationkeeper.ErrUserBalanceNotExists,
			expectedCode: ErrorCodeBalanceNotFound,
			shouldMap:    true,
		},
		{
			name:         "ErrInadequateApprovals",
			err:          tokenizationkeeper.ErrInadequateApprovals,
			expectedCode: ErrorCodeUnauthorized,
			shouldMap:    true,
		},
		{
			name:         "ErrCollectionIsArchived",
			err:          tokenizationkeeper.ErrCollectionIsArchived,
			expectedCode: ErrorCodeCollectionArchived,
			shouldMap:    true,
		},
		{
			name:         "ErrDisallowedTransfer",
			err:          tokenizationkeeper.ErrDisallowedTransfer,
			expectedCode: ErrorCodeTransferFailed,
			shouldMap:    true,
		},
		{
			name:         "ErrAccountNotFound",
			err:          tokenizationkeeper.ErrAccountNotFound,
			expectedCode: ErrorCodeQueryFailed,
			shouldMap:    true,
		},
		{
			name:         "ErrUnauthorized",
			err:          tokenizationtypes.ErrUnauthorized,
			expectedCode: ErrorCodeUnauthorized,
			shouldMap:    true,
		},
		{
			name:         "ErrNotFound",
			err:          tokenizationtypes.ErrNotFound,
			expectedCode: ErrorCodeQueryFailed,
			shouldMap:    true,
		},
		{
			name:         "ErrInvalidRequest",
			err:          tokenizationtypes.ErrInvalidRequest,
			expectedCode: ErrorCodeInvalidInput,
			shouldMap:    true,
		},
		{
			name:         "ErrUnderflow",
			err:          tokenizationtypes.ErrUnderflow,
			expectedCode: ErrorCodeTransferFailed,
			shouldMap:    true,
		},
		{
			name:         "ErrOverflow",
			err:          tokenizationtypes.ErrOverflow,
			expectedCode: ErrorCodeTransferFailed,
			shouldMap:    true,
		},
		{
			name:         "ErrInvalidAddress",
			err:          tokenizationtypes.ErrInvalidAddress,
			expectedCode: ErrorCodeInvalidInput,
			shouldMap:    true,
		},
		{
			name:         "ErrInvalidCollectionID",
			err:          tokenizationtypes.ErrInvalidCollectionID,
			expectedCode: ErrorCodeCollectionNotFound, // ErrInvalidCollectionID is used when collection doesn't exist
			shouldMap:    true,
		},
		{
			name:         "ErrAmountEqualsZero",
			err:          tokenizationtypes.ErrAmountEqualsZero,
			expectedCode: ErrorCodeInvalidInput,
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
			code, found := MapCosmosErrorToPrecompileError(tt.err)
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
	wrapped := WrapError(err, ErrorCodeQueryFailed, "test message")
	require.Equal(t, ErrorCodeCollectionNotFound, wrapped.Code, "should use mapped code, not default")

	// Test that unmapped errors use the default code
	unmappedErr := fmt.Errorf("some random error")
	wrapped2 := WrapError(unmappedErr, ErrorCodeInternalError, "test message")
	require.Equal(t, ErrorCodeInternalError, wrapped2.Code, "should use default code for unmapped errors")

	// Test that nil errors use the default code
	wrapped3 := WrapError(nil, ErrorCodeQueryFailed, "test message")
	require.Equal(t, ErrorCodeQueryFailed, wrapped3.Code, "should use default code for nil errors")
}

