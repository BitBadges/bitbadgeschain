package tokenization

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
//    - The caller is used as the "from" address for transfers
//    - Zero address callers are explicitly rejected via VerifyCaller()
//
// 3. Overflow Protection:
//    - All big.Int values are validated before conversion to sdkmath.Uint
//    - Range validation ensures start <= end
//    - Amount validation ensures amount > 0
//    - CheckOverflow() validates that values are non-negative before conversion
//    - sdkmath.Uint provides overflow protection for arithmetic operations
//
// 4. Input Validation:
//    - All inputs are validated before processing
//    - Zero addresses are rejected via ValidateAddress()
//    - Empty arrays are rejected where appropriate via ValidateArraySize()
//    - Invalid ranges (start > end) are rejected via ValidateBigIntRanges()
//    - Collection IDs are validated to be non-zero via ValidateCollectionId()
//    - String inputs are validated to be non-empty via ValidateString()
//
// 5. Denial of Service (DoS) Protection:
//    - Array size limits prevent DoS through extremely large arrays:
//      * MaxRecipients = 100
//      * MaxTokenIdRanges = 100
//      * MaxOwnershipTimeRanges = 100
//      * MaxApprovalRanges = 100
//    - These limits are enforced via ValidateArraySize() and ValidateTransferInputs()
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
//
// Threat Model:
//
// The precompile is designed to protect against:
// - Reentrancy attacks: Prevented by EVM call stack and atomic transactions
// - Integer overflow: Prevented by validation and sdkmath.Uint type
// - Invalid input attacks: Prevented by comprehensive input validation
// - DoS attacks: Prevented by array size limits
// - Information leakage: Prevented by structured error handling
// - State corruption: Prevented by atomic transactions and keeper validation
//
// Known Limitations:
// - Rate limiting is not implemented at the precompile level (can be added at chain level)
// - Gas price manipulation protection is handled by the EVM module, not the precompile
// - Access control is handled by the tokenization module's approval system

// VerifyCaller verifies that the caller is not a zero address
// This is a security check to ensure valid callers
func VerifyCaller(caller common.Address) error {
	if caller == (common.Address{}) {
		return ErrUnauthorized("caller cannot be zero address")
	}
	return nil
}

// CheckOverflow checks if a big.Int value would overflow when converted to sdkmath.Uint
// sdkmath.Uint can handle values up to 2^256-1, same as big.Int
// This function validates that the value is non-negative
func CheckOverflow(value *big.Int, fieldName string) error {
	if value == nil {
		return ErrInvalidInput(fmt.Sprintf("%s cannot be nil", fieldName))
	}
	if value.Sign() < 0 {
		return ErrInvalidInput(fmt.Sprintf("%s cannot be negative", fieldName))
	}
	// sdkmath.Uint uses uint256 internally, which can hold values up to 2^256-1
	// Check if value exceeds uint256 max (2^256 - 1)
	maxUint256 := new(big.Int)
	maxUint256.Lsh(big.NewInt(1), 256)        // 2^256
	maxUint256.Sub(maxUint256, big.NewInt(1)) // 2^256 - 1

	if value.Cmp(maxUint256) > 0 {
		return ErrInvalidInput(fmt.Sprintf("%s overflow: value %s exceeds maximum uint256 value (2^256-1)", fieldName, value.String()))
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

// Maximum allowed sizes for arrays (DoS protection)
const (
	MaxRecipients          = 100
	MaxTokenIdRanges       = 100
	MaxOwnershipTimeRanges = 100
	MaxApprovalRanges      = 100
	// Additional DoS limits for nested structures
	MaxDenomUnits          = 50  // Maximum denom units per path
	MaxMerkleChallenges    = 20  // Maximum merkle challenges per approval
	MaxCoinTransfers       = 50  // Maximum coin transfers per approval
	MaxDynamicStoreChallenges = 20 // Maximum dynamic store challenges
	MaxETHSignatureChallenges = 20 // Maximum ETH signature challenges
	MaxVotingChallenges    = 20  // Maximum voting challenges
	MaxMustOwnTokens       = 50  // Maximum must own token rules
	MaxAddressListEntries  = 1000 // Maximum addresses per address list
	MaxMetadataLength      = 10000 // Maximum length for metadata strings (URI, customData)
)

// ValidateDenomUnitsSize validates that denomUnits array size is within limits
func ValidateDenomUnitsSize(size int) error {
	if size > MaxDenomUnits {
		return ErrInvalidInput(fmt.Sprintf("denomUnits size (%d) exceeds maximum allowed size (%d)", size, MaxDenomUnits))
	}
	return nil
}

// ValidateMerkleChallengesSize validates that merkle challenges array size is within limits
func ValidateMerkleChallengesSize(size int) error {
	if size > MaxMerkleChallenges {
		return ErrInvalidInput(fmt.Sprintf("merkleChallenges size (%d) exceeds maximum allowed size (%d)", size, MaxMerkleChallenges))
	}
	return nil
}

// ValidateCoinTransfersSize validates that coin transfers array size is within limits
func ValidateCoinTransfersSize(size int) error {
	if size > MaxCoinTransfers {
		return ErrInvalidInput(fmt.Sprintf("coinTransfers size (%d) exceeds maximum allowed size (%d)", size, MaxCoinTransfers))
	}
	return nil
}

// ValidateMetadataLength validates that metadata strings are within size limits
func ValidateMetadataLength(s string, fieldName string) error {
	if len(s) > MaxMetadataLength {
		return ErrInvalidInput(fmt.Sprintf("%s length (%d) exceeds maximum allowed length (%d)", fieldName, len(s), MaxMetadataLength))
	}
	return nil
}

// ValidateTransferInputs performs comprehensive security validation for transfer inputs
func ValidateTransferInputs(
	toAddresses []common.Address,
	tokenIdsRanges []struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	},
	ownershipTimesRanges []struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	},
) error {
	if err := ValidateArraySize(len(toAddresses), MaxRecipients, "toAddresses"); err != nil {
		return err
	}
	if err := ValidateArraySize(len(tokenIdsRanges), MaxTokenIdRanges, "tokenIds"); err != nil {
		return err
	}
	if err := ValidateArraySize(len(ownershipTimesRanges), MaxOwnershipTimeRanges, "ownershipTimes"); err != nil {
		return err
	}
	return nil
}

// ValidateApprovalInputs performs comprehensive security validation for approval inputs
func ValidateApprovalInputs(
	transferTimes []struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	},
	tokenIds []struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	},
	ownershipTimes []struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	},
) error {
	if err := ValidateArraySize(len(transferTimes), MaxApprovalRanges, "transferTimes"); err != nil {
		return err
	}
	if err := ValidateArraySize(len(tokenIds), MaxTokenIdRanges, "tokenIds"); err != nil {
		return err
	}
	if err := ValidateArraySize(len(ownershipTimes), MaxOwnershipTimeRanges, "ownershipTimes"); err != nil {
		return err
	}
	return nil
}

