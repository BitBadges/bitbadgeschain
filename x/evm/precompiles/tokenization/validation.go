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
		return fmt.Errorf("%s cannot be zero address", fieldName)
	}
	return nil
}

// ValidateAddresses validates that an address array is not empty and contains no zero addresses
func ValidateAddresses(addrs []common.Address, fieldName string) error {
	if len(addrs) == 0 {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}
	if len(addrs) > MaxRecipients {
		return fmt.Errorf("%s size (%d) exceeds maximum allowed size (%d)", fieldName, len(addrs), MaxRecipients)
	}
	for i, addr := range addrs {
		if addr == (common.Address{}) {
			return fmt.Errorf("%s[%d] cannot be zero address", fieldName, i)
		}
	}
	return nil
}

// ValidateUintRange validates that a UintRange is valid (start <= end)
func ValidateUintRange(r *tokenizationtypes.UintRange, fieldName string) error {
	if r == nil {
		return fmt.Errorf("%s cannot be nil", fieldName)
	}
	if r.Start.GT(r.End) {
		return fmt.Errorf("%s: start (%s) cannot be greater than end (%s)", fieldName, r.Start.String(), r.End.String())
	}
	return nil
}

// ValidateUintRanges validates that all UintRanges in a slice are valid
func ValidateUintRanges(ranges []*tokenizationtypes.UintRange, fieldName string) error {
	if len(ranges) == 0 {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}
	for i, r := range ranges {
		if err := ValidateUintRange(r, fmt.Sprintf("%s[%d]", fieldName, i)); err != nil {
			return err
		}
	}
	return nil
}

// ValidateBigIntRange validates that a big.Int range is valid (start <= end, both non-negative)
func ValidateBigIntRange(start, end *big.Int, fieldName string) error {
	if start == nil {
		return fmt.Errorf("%s.start cannot be nil", fieldName)
	}
	if end == nil {
		return fmt.Errorf("%s.end cannot be nil", fieldName)
	}
	if start.Sign() < 0 {
		return fmt.Errorf("%s.start cannot be negative, got %s", fieldName, start.String())
	}
	if end.Sign() < 0 {
		return fmt.Errorf("%s.end cannot be negative, got %s", fieldName, end.String())
	}
	if start.Cmp(end) > 0 {
		return fmt.Errorf("%s: start (%s) cannot be greater than end (%s)", fieldName, start.String(), end.String())
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
		return fmt.Errorf("%s cannot be empty", fieldName)
	}
	for i, r := range ranges {
		if err := ValidateBigIntRange(r.Start, r.End, fmt.Sprintf("%s[%d]", fieldName, i)); err != nil {
			return err
		}
	}
	return nil
}

// ValidateAmount validates that an amount is greater than zero
func ValidateAmount(amount *big.Int, fieldName string) error {
	if amount == nil {
		return fmt.Errorf("%s cannot be nil", fieldName)
	}
	if amount.Sign() <= 0 {
		return fmt.Errorf("%s must be greater than zero, got %s", fieldName, amount.String())
	}
	return nil
}

// ValidateCollectionId validates that a collection ID is valid
func ValidateCollectionId(collectionId *big.Int) error {
	if collectionId == nil {
		return fmt.Errorf("collectionId cannot be nil")
	}
	if collectionId.Sign() < 0 {
		return fmt.Errorf("collectionId cannot be negative, got %s", collectionId.String())
	}
	return nil
}

// ValidateString validates that a string is not empty
func ValidateString(s, fieldName string) error {
	if s == "" {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}
	return nil
}

// ValidateStringOptional validates that a string is either empty or non-empty (no validation if empty)
func ValidateStringOptional(s, fieldName string) error {
	// Optional strings are allowed to be empty
	return nil
}

// ConvertAndValidateBigIntRanges converts big.Int ranges to UintRange and validates them
func ConvertAndValidateBigIntRanges(ranges []struct {
	Start *big.Int `json:"start"`
	End   *big.Int `json:"end"`
}, fieldName string,
) ([]*tokenizationtypes.UintRange, error) {
	if err := ValidateBigIntRanges(ranges, fieldName); err != nil {
		return nil, err
	}

	uintRanges := make([]*tokenizationtypes.UintRange, len(ranges))
	for i, r := range ranges {
		// Check for overflow - Uint can handle up to 2^256-1, but we should validate reasonable bounds
		uintRanges[i] = &tokenizationtypes.UintRange{
			Start: sdkmath.NewUintFromBigInt(r.Start),
			End:   sdkmath.NewUintFromBigInt(r.End),
		}
		// Validate the converted range
		if err := ValidateUintRange(uintRanges[i], fmt.Sprintf("%s[%d]", fieldName, i)); err != nil {
			return nil, err
		}
	}

	return uintRanges, nil
}
