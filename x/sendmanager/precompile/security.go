package precompile

import (
	"fmt"
	"math/big"
)

// Security considerations for the sendmanager precompile:
//
// This package implements security measures to protect against common attack vectors
// in EVM precompiles. The security model relies on both EVM-level protections and
// application-level validation.
//
// 1. Reentrancy Protection:
//    - The EVM itself prevents reentrancy through its call stack mechanism
//    - Cosmos SDK's state machine ensures atomicity of operations
//    - Precompile operations are executed atomically within a single transaction
//
// 2. Caller Verification:
//    - contract.Caller() returns the address that initiated the call
//    - This is verified by the EVM and cannot be spoofed
//    - The caller is used as the "from_address" for send operations
//    - Zero address callers are explicitly rejected via VerifyCaller()
//
// 3. Overflow Protection:
//    - All big.Int values are validated before conversion to sdkmath.Int
//    - Amount validation ensures values don't exceed uint256 max
//    - CheckOverflow() validates that values are non-negative and within bounds
//
// 4. Input Validation:
//    - All inputs are validated via ValidateBasic() on the Cosmos message
//    - Zero addresses are rejected
//    - Empty denoms are rejected
//
// 5. State Consistency:
//    - All state changes go through the Cosmos SDK keeper, ensuring consistency
//    - Transactions are atomic - either all state changes succeed or all fail

// MaxUint256 is the maximum value for Solidity uint256 (2^256 - 1)
// Pre-computed once for efficiency and reused across all overflow checks.
// This is critical for EVM compatibility - values exceeding this cannot be
// represented in Solidity and must be rejected to prevent silent truncation.
var MaxUint256 = func() *big.Int {
	max := new(big.Int)
	max.Lsh(big.NewInt(1), 256)        // 2^256
	max.Sub(max, big.NewInt(1))        // 2^256 - 1
	return max
}()

// MaxInt256 is the maximum value for Solidity int256 (2^255 - 1)
// Used for signed integer values in Cosmos SDK (sdkmath.Int).
var MaxInt256 = func() *big.Int {
	max := new(big.Int)
	max.Lsh(big.NewInt(1), 255)        // 2^255
	max.Sub(max, big.NewInt(1))        // 2^255 - 1
	return max
}()

// CheckUint256Overflow checks if a big.Int value would overflow Solidity uint256.
// This MUST be called before packing any big.Int value for EVM return to ensure
// compatibility and prevent silent data truncation.
//
// Returns error if:
// - value is nil
// - value is negative
// - value exceeds 2^256-1 (MaxUint256)
func CheckUint256Overflow(value *big.Int, fieldName string) error {
	if value == nil {
		return ErrInvalidInput(fmt.Sprintf("%s cannot be nil", fieldName))
	}
	if value.Sign() < 0 {
		return ErrInvalidInput(fmt.Sprintf("%s cannot be negative", fieldName))
	}
	if value.Cmp(MaxUint256) > 0 {
		return ErrInvalidInput(fmt.Sprintf("%s overflow: value exceeds maximum uint256 (2^256-1)", fieldName))
	}
	return nil
}

// CheckInt256Overflow checks if a big.Int value would overflow Solidity int256.
// Used for signed values from sdkmath.Int.
//
// Returns error if:
// - value is nil
// - value exceeds 2^255-1 (MaxInt256) or is less than -(2^255)
func CheckInt256Overflow(value *big.Int, fieldName string) error {
	if value == nil {
		return ErrInvalidInput(fmt.Sprintf("%s cannot be nil", fieldName))
	}
	// For signed int256, check both positive and negative bounds
	minInt256 := new(big.Int).Neg(new(big.Int).Add(MaxInt256, big.NewInt(1)))
	if value.Cmp(MaxInt256) > 0 {
		return ErrInvalidInput(fmt.Sprintf("%s overflow: value exceeds maximum int256 (2^255-1)", fieldName))
	}
	if value.Cmp(minInt256) < 0 {
		return ErrInvalidInput(fmt.Sprintf("%s underflow: value is less than minimum int256 (-(2^255))", fieldName))
	}
	return nil
}

// Maximum allowed sizes for arrays (DoS protection)
const (
	MaxCoins        = 20    // Maximum coins in a send
	MaxStringLength = 10000 // Maximum length for strings
)

// ValidateArraySize validates that an array size is within reasonable bounds
// This helps prevent DoS attacks through extremely large arrays
func ValidateArraySize(size int, maxSize int, fieldName string) error {
	if size == 0 {
		return ErrInvalidInput(fmt.Sprintf("%s cannot be empty", fieldName))
	}
	if size > maxSize {
		return ErrInvalidInput(fmt.Sprintf("%s size (%d) exceeds maximum allowed size (%d)", fieldName, size, maxSize))
	}
	return nil
}

// ValidateStringLength validates that strings are within size limits
func ValidateStringLength(s string, fieldName string) error {
	if len(s) > MaxStringLength {
		return ErrInvalidInput(fmt.Sprintf("%s length (%d) exceeds maximum allowed length (%d)", fieldName, len(s), MaxStringLength))
	}
	return nil
}
