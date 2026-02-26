package tokenization

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	sdkmath "cosmossdk.io/math"

	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// ValidateAddress validates that an address is not zero
func ValidateAddress(addr common.Address, fieldName string) error {
	if addr == (common.Address{}) {
		return ErrInvalidInput(fmt.Sprintf("%s cannot be zero address", fieldName))
	}
	return nil
}

// ValidateAddresses validates that an address array is not empty and contains no zero addresses
func ValidateAddresses(addrs []common.Address, fieldName string) error {
	if len(addrs) == 0 {
		return ErrInvalidInput(fmt.Sprintf("%s cannot be empty", fieldName))
	}
	if len(addrs) > MaxRecipients {
		return ErrInvalidInput(fmt.Sprintf("%s size (%d) exceeds maximum allowed size (%d)", fieldName, len(addrs), MaxRecipients))
	}
	for i, addr := range addrs {
		if addr == (common.Address{}) {
			return ErrInvalidInput(fmt.Sprintf("%s[%d] cannot be zero address", fieldName, i))
		}
	}
	return nil
}

// ValidateUintRange validates that a UintRange is valid (start <= end)
func ValidateUintRange(r *tokenizationtypes.UintRange, fieldName string) error {
	if r == nil {
		return ErrInvalidInput(fmt.Sprintf("%s cannot be nil", fieldName))
	}
	if r.Start.GT(r.End) {
		return ErrInvalidInput(fmt.Sprintf("%s: start (%s) cannot be greater than end (%s)", fieldName, r.Start.String(), r.End.String()))
	}
	return nil
}

// ValidateUintRanges validates that all UintRanges in a slice are valid
func ValidateUintRanges(ranges []*tokenizationtypes.UintRange, fieldName string) error {
	if len(ranges) == 0 {
		return ErrInvalidInput(fmt.Sprintf("%s cannot be empty", fieldName))
	}
	for i, r := range ranges {
		if err := ValidateUintRange(r, fmt.Sprintf("%s[%d]", fieldName, i)); err != nil {
			return err
		}
	}
	return nil
}

// ValidateBigIntRange validates that a big.Int range is valid:
// - Both start and end are non-nil
// - Both are non-negative
// - Both are within uint256 bounds (for EVM compatibility)
// - start <= end
func ValidateBigIntRange(start, end *big.Int, fieldName string) error {
	// Check for nil
	if start == nil {
		return ErrInvalidInput(fmt.Sprintf("%s.start cannot be nil", fieldName))
	}
	if end == nil {
		return ErrInvalidInput(fmt.Sprintf("%s.end cannot be nil", fieldName))
	}
	// Check overflow (includes negative check)
	if err := CheckOverflow(start, fmt.Sprintf("%s.start", fieldName)); err != nil {
		return err
	}
	if err := CheckOverflow(end, fmt.Sprintf("%s.end", fieldName)); err != nil {
		return err
	}
	// Check start <= end
	if start.Cmp(end) > 0 {
		return ErrInvalidInput(fmt.Sprintf("%s: start (%s) cannot be greater than end (%s)", fieldName, start.String(), end.String()))
	}
	return nil
}

// ValidateBigIntRanges validates that all big.Int ranges in a slice are valid
func ValidateBigIntRanges(ranges []struct {
	Start *big.Int `json:"start"`
	End   *big.Int `json:"end"`
}, fieldName string,
) error {
	if len(ranges) == 0 {
		return ErrInvalidInput(fmt.Sprintf("%s cannot be empty", fieldName))
	}
	for i, r := range ranges {
		if err := ValidateBigIntRange(r.Start, r.End, fmt.Sprintf("%s[%d]", fieldName, i)); err != nil {
			return err
		}
	}
	return nil
}

// ValidateAmount validates that an amount is valid:
// - Non-nil
// - Greater than zero
// - Within uint256 bounds (for EVM compatibility)
func ValidateAmount(amount *big.Int, fieldName string) error {
	if amount == nil {
		return ErrInvalidInput(fmt.Sprintf("%s cannot be nil", fieldName))
	}
	if amount.Sign() <= 0 {
		return ErrInvalidInput(fmt.Sprintf("%s must be greater than zero, got %s", fieldName, amount.String()))
	}
	// Check uint256 overflow
	if err := CheckOverflow(amount, fieldName); err != nil {
		return err
	}
	return nil
}

// ValidateNonOverlappingRanges checks if any ranges in the array overlap
func ValidateNonOverlappingRanges(ranges []struct {
	Start *big.Int `json:"start"`
	End   *big.Int `json:"end"`
}, fieldName string) error {
	if len(ranges) <= 1 {
		return nil
	}

	// Using a simple O(n^2) comparison since arrays are typically small
	for i := 0; i < len(ranges); i++ {
		for j := i + 1; j < len(ranges); j++ {
			// Check if ranges[i] and ranges[j] overlap
			// Ranges overlap if: start_i <= end_j AND start_j <= end_i
			if ranges[i].Start.Cmp(ranges[j].End) <= 0 && ranges[j].Start.Cmp(ranges[i].End) <= 0 {
				return ErrInvalidInput(fmt.Sprintf("%s contains overlapping ranges: [%s, %s] and [%s, %s]",
					fieldName,
					ranges[i].Start.String(), ranges[i].End.String(),
					ranges[j].Start.String(), ranges[j].End.String()))
			}
		}
	}

	return nil
}

// ValidateRangesWithOverlapCheck validates ranges and optionally checks for overlaps
func ValidateRangesWithOverlapCheck(ranges []struct {
	Start *big.Int `json:"start"`
	End   *big.Int `json:"end"`
}, fieldName string, allowOverlap bool) error {
	// First validate the basic range constraints
	if err := ValidateBigIntRanges(ranges, fieldName); err != nil {
		return err
	}

	// Check for overlaps if not allowed
	if !allowOverlap {
		if err := ValidateNonOverlappingRanges(ranges, fieldName); err != nil {
			return err
		}
	}

	return nil
}

// ValidateCollectionId validates that a collection ID is valid:
// - Non-nil
// - Non-negative
// - Non-zero (except for new collection creation where 0 means "assign new ID")
// - Within uint256 bounds (for EVM compatibility)
func ValidateCollectionId(collectionId *big.Int) error {
	if collectionId == nil {
		return ErrInvalidInput("collectionId cannot be nil")
	}
	// Check uint256 overflow (includes negative check)
	if err := CheckOverflow(collectionId, "collectionId"); err != nil {
		return err
	}
	if collectionId.Sign() == 0 {
		return ErrInvalidInput("collectionId cannot be zero")
	}
	return nil
}

// ValidateString validates that a string is not empty
func ValidateString(s, fieldName string) error {
	if s == "" {
		return ErrInvalidInput(fmt.Sprintf("%s cannot be empty", fieldName))
	}
	return nil
}

// ValidateStringOptional validates that a string is either empty or non-empty (no validation if empty)
func ValidateStringOptional(s, fieldName string) error {
	// Optional strings are allowed to be empty
	return nil
}

// ConvertAndValidateBigIntRanges converts big.Int ranges to UintRange and validates them.
// Validates that all values are within uint256 bounds before conversion to ensure EVM compatibility.
// Requires non-empty array - use ConvertAndValidateBigIntRangesAllowEmpty for optional arrays.
func ConvertAndValidateBigIntRanges(ranges []struct {
	Start *big.Int `json:"start"`
	End   *big.Int `json:"end"`
}, fieldName string,
) ([]*tokenizationtypes.UintRange, error) {
	// ValidateBigIntRanges calls ValidateBigIntRange which includes CheckOverflow
	if err := ValidateBigIntRanges(ranges, fieldName); err != nil {
		return nil, err
	}

	// Safe to convert - overflow has been checked
	uintRanges := make([]*tokenizationtypes.UintRange, len(ranges))
	for i, r := range ranges {
		uintRanges[i] = &tokenizationtypes.UintRange{
			Start: sdkmath.NewUintFromBigInt(r.Start),
			End:   sdkmath.NewUintFromBigInt(r.End),
		}
		// Validate the converted range (start <= end check on sdkmath.Uint)
		if err := ValidateUintRange(uintRanges[i], fmt.Sprintf("%s[%d]", fieldName, i)); err != nil {
			return nil, err
		}
	}

	return uintRanges, nil
}

// ConvertAndValidateBigIntRangesAllowEmpty converts big.Int ranges to UintRange and validates them.
// Validates that all values are within uint256 bounds before conversion to ensure EVM compatibility.
// Allows empty arrays (returns empty slice), validates individual entries if present.
func ConvertAndValidateBigIntRangesAllowEmpty(ranges []struct {
	Start *big.Int `json:"start"`
	End   *big.Int `json:"end"`
}, fieldName string,
) ([]*tokenizationtypes.UintRange, error) {
	if len(ranges) == 0 {
		return []*tokenizationtypes.UintRange{}, nil
	}

	uintRanges := make([]*tokenizationtypes.UintRange, len(ranges))
	for i, r := range ranges {
		// ValidateBigIntRange includes CheckOverflow for uint256 bounds
		if err := ValidateBigIntRange(r.Start, r.End, fmt.Sprintf("%s[%d]", fieldName, i)); err != nil {
			return nil, err
		}
		// Safe to convert - overflow has been checked
		uintRanges[i] = &tokenizationtypes.UintRange{
			Start: sdkmath.NewUintFromBigInt(r.Start),
			End:   sdkmath.NewUintFromBigInt(r.End),
		}
	}

	return uintRanges, nil
}
