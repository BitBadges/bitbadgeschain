package gamm

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Security considerations for the precompile:
//
// This package implements security measures to protect against common attack vectors
// in EVM precompiles. The security model relies on both EVM-level protections and
// application-level validation.
//
// 1. Reentrancy Protection:
//    - The EVM itself prevents reentrancy through its call stack mechanism
//    - Cosmos SDK's state machine ensures atomicity of operations
//    - Precompile operations are executed atomically within a single transaction
//    - No external calls are made during state transitions that could enable reentrancy
//
// 2. Caller Verification:
//    - contract.Caller() returns the address that initiated the call
//    - This is verified by the EVM and cannot be spoofed
//    - The caller is used as the "sender" address for pool operations
//    - Zero address callers are explicitly rejected via VerifyCaller()
//
// 3. Overflow Protection:
//    - All big.Int values are validated before conversion to sdkmath.Int
//    - Amount validation ensures amount > 0
//    - CheckOverflow() validates that values are non-negative before conversion
//    - sdkmath.Int provides overflow protection for arithmetic operations
//
// 4. Input Validation:
//    - All inputs are validated before processing
//    - Zero addresses are rejected via ValidateAddress()
//    - Empty arrays are rejected where appropriate
//    - Pool IDs are validated to exist
//    - String inputs are validated to be non-empty where required
//
// 5. Denial of Service (DoS) Protection:
//    - Array size limits prevent DoS through extremely large arrays:
//      * MaxRoutes = 10
//      * MaxCoins = 20
//      * MaxAffiliates = 10
//    - These limits are enforced via validation functions
//
// 6. Error Handling:
//    - Structured error types prevent information leakage
//    - Error sanitization removes sensitive details from error messages
//    - Error codes allow clients to handle errors appropriately without exposing internals
//
// 7. State Consistency:
//    - All state changes go through the Cosmos SDK keeper, ensuring consistency
//    - Transactions are atomic - either all state changes succeed or all fail
//    - No partial state updates are possible

// VerifyCaller verifies that the caller is not a zero address
// This is a security check to ensure valid callers
func VerifyCaller(caller common.Address) error {
	if caller == (common.Address{}) {
		return ErrUnauthorized("caller cannot be zero address")
	}
	return nil
}

// CheckOverflow checks if a big.Int value would overflow when converted to sdkmath.Int.
// sdkmath.Int uses int256 internally (signed), which can hold values up to 2^255-1.
// This function MUST be called before any conversion to ensure EVM compatibility
// and prevent silent data truncation.
//
// Returns error if:
// - value is nil
// - value is negative
// - value exceeds 2^255-1 (MaxInt256)
func CheckOverflow(value *big.Int, fieldName string) error {
	if value == nil {
		return ErrInvalidInput(fmt.Sprintf("%s cannot be nil", fieldName))
	}
	if value.Sign() < 0 {
		return ErrInvalidInput(fmt.Sprintf("%s cannot be negative", fieldName))
	}
	if value.Cmp(MaxInt256) > 0 {
		return ErrInvalidInput(fmt.Sprintf("%s overflow: value exceeds maximum int256 (2^255-1)", fieldName))
	}
	return nil
}

// CheckUint256Overflow checks if a big.Int value would overflow Solidity uint256.
// This is used for unsigned values that will be returned to EVM contracts.
// While sdkmath.Int uses int256 internally, some return values are packed as uint256.
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

// ValidateArraySizeAllowEmpty validates array size but allows empty arrays
// Use this for optional arrays (e.g., tokenOutMins where empty means no minimum)
func ValidateArraySizeAllowEmpty(size int, maxSize int, fieldName string) error {
	if size > maxSize {
		return ErrInvalidInput(fmt.Sprintf("%s size (%d) exceeds maximum allowed size (%d)", fieldName, size, maxSize))
	}
	return nil
}

// Maximum allowed sizes for arrays (DoS protection)
const (
	MaxRoutes          = 10  // Maximum swap routes
	MaxCoins           = 20  // Maximum coins in arrays
	MaxAffiliates      = 10  // Maximum affiliates per swap
	MaxStringLength    = 10000 // Maximum length for strings (channel, receiver, memo)
	MaxMemoLength      = 256  // Maximum length for IBC transfer memo
	MaxPaginationLimit = 100 // Maximum limit for pagination queries
)

// MaxInt256 is the maximum value for Solidity int256 (2^255 - 1)
// Pre-computed once for efficiency and reused across all overflow checks.
// GAMM uses sdkmath.Int (signed) so we check against int256 max, not uint256.
// This is critical for EVM compatibility - values exceeding this cannot be
// represented in Solidity and must be rejected to prevent silent truncation.
var MaxInt256 = func() *big.Int {
	max := new(big.Int)
	max.Lsh(big.NewInt(1), 255)        // 2^255
	max.Sub(max, big.NewInt(1))        // 2^255 - 1
	return max
}()

// MaxUint256 is the maximum value for Solidity uint256 (2^256 - 1)
// Used for unsigned values like pool shares and token amounts.
var MaxUint256 = func() *big.Int {
	max := new(big.Int)
	max.Lsh(big.NewInt(1), 256)        // 2^256
	max.Sub(max, big.NewInt(1))        // 2^256 - 1
	return max
}()

// ValidateStringLength validates that strings are within size limits
func ValidateStringLength(s string, fieldName string) error {
	if len(s) > MaxStringLength {
		return ErrInvalidInput(fmt.Sprintf("%s length (%d) exceeds maximum allowed length (%d)", fieldName, len(s), MaxStringLength))
	}
	return nil
}

