package tokenization

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	tokenization "github.com/bitbadges/bitbadgeschain/x/tokenization/precompile"
)

// FuzzValidateAddress fuzzes the ValidateAddress function
func FuzzValidateAddress(f *testing.F) {
	// Seed with valid address
	f.Add([]byte{0x74, 0x2d, 0x35, 0xCc, 0x66, 0x34, 0xC0, 0x53, 0x29, 0x25, 0xa3, 0xb8, 0x44, 0xBc, 0x9e, 0x75, 0x95, 0xf0, 0xbE, 0xb0})
	// Seed with zero address
	f.Add([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})

	f.Fuzz(func(t *testing.T, addrBytes []byte) {
		// Ensure we have exactly 20 bytes for an address
		if len(addrBytes) != 20 {
			return
		}

		addr := common.BytesToAddress(addrBytes)
		err := tokenization.ValidateAddress(addr, "testAddress")

		// If it's a zero address, we expect an error
		if addr == (common.Address{}) {
			if err == nil {
				t.Errorf("expected error for zero address, got nil")
			}
		} else {
			// For non-zero addresses, validation should pass
			if err != nil {
				t.Errorf("unexpected error for valid address: %v", err)
			}
		}
	})
}

// FuzzValidateCollectionId fuzzes the ValidateCollectionId function
func FuzzValidateCollectionId(f *testing.F) {
	// Seed with valid IDs
	f.Add(int64(1))
	f.Add(int64(0))
	f.Add(int64(-1))
	f.Add(int64(1000000))

	f.Fuzz(func(t *testing.T, id int64) {
		collectionId := big.NewInt(id)
		err := tokenization.ValidateCollectionId(collectionId)

		// Negative IDs should fail
		if id < 0 {
			if err == nil {
				t.Errorf("expected error for negative collectionId %d, got nil", id)
			}
			return
		}
		// Zero should fail (only valid when creating new collections, not for queries)
		if id == 0 {
			if err == nil {
				t.Errorf("expected error for collectionId 0, got nil")
			}
			return
		}
		// Positive IDs should pass
		if err != nil {
			t.Errorf("unexpected error for valid collectionId %d: %v", id, err)
		}
	})
}

// FuzzValidateAmount fuzzes the ValidateAmount function
func FuzzValidateAmount(f *testing.F) {
	// Seed with valid amounts
	f.Add(int64(1))
	f.Add(int64(0))
	f.Add(int64(-1))
	f.Add(int64(1000000))

	f.Fuzz(func(t *testing.T, amount int64) {
		amountBig := big.NewInt(amount)
		err := tokenization.ValidateAmount(amountBig, "testAmount")

		// Zero or negative amounts should fail
		if amount <= 0 {
			if err == nil {
				t.Errorf("expected error for amount %d, got nil", amount)
			}
		} else {
			// Positive amounts should pass
			if err != nil {
				t.Errorf("unexpected error for valid amount %d: %v", amount, err)
			}
		}
	})
}

// FuzzValidateBigIntRanges fuzzes the ValidateBigIntRanges function
func FuzzValidateBigIntRanges(f *testing.F) {
	// Seed with valid range
	f.Add(int64(1), int64(10))
	// Seed with invalid range (start > end)
	f.Add(int64(10), int64(1))
	// Seed with equal start and end
	f.Add(int64(5), int64(5))

	f.Fuzz(func(t *testing.T, start, end int64) {
		ranges := []struct {
			Start *big.Int `json:"start"`
			End   *big.Int `json:"end"`
		}{
			{Start: big.NewInt(start), End: big.NewInt(end)},
		}

		err := tokenization.ValidateBigIntRanges(ranges, "testRanges")

		// Negative values or start > end should fail
		if start < 0 || end < 0 || start > end {
			if err == nil {
				t.Errorf("expected error for range [%d, %d], got nil", start, end)
			}
		} else if start > 0 && end > 0 && start <= end {
			// Valid ranges should pass
			if err != nil {
				t.Errorf("unexpected error for valid range [%d, %d]: %v", start, end, err)
			}
		}
	})
}

// FuzzValidateString fuzzes the ValidateString function with security edge cases
func FuzzValidateString(f *testing.F) {
	// Seed with valid strings
	f.Add("valid-collection-name")
	// Seed with empty string
	f.Add("")
	// Seed with unicode
	f.Add("\u0000\uFFFF")
	// Seed with null bytes
	f.Add("test\x00value")
	// Seed with very long string
	longStr := make([]byte, 10000)
	for i := range longStr {
		longStr[i] = 'a'
	}
	f.Add(string(longStr))
	// Seed with special characters
	f.Add("<script>alert('xss')</script>")
	f.Add("../../../etc/passwd")

	f.Fuzz(func(t *testing.T, s string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("panic on input string (len=%d): %v", len(s), r)
			}
		}()

		err := tokenization.ValidateString(s, "testString")

		// Empty strings should fail
		if s == "" {
			if err == nil {
				t.Errorf("expected error for empty string, got nil")
			}
		}
		// For non-empty strings, should either succeed or fail gracefully (no panic)
		_ = err
	})
}

// FuzzValidateNonOverlappingRanges fuzzes the overlap detection with adversarial inputs
func FuzzValidateNonOverlappingRanges(f *testing.F) {
	// Seed with non-overlapping ranges
	f.Add(int64(1), int64(5), int64(6), int64(10))
	// Seed with overlapping ranges
	f.Add(int64(1), int64(10), int64(5), int64(15))
	// Seed with exact same ranges
	f.Add(int64(1), int64(10), int64(1), int64(10))
	// Seed with adjacent ranges
	f.Add(int64(1), int64(5), int64(5), int64(10))
	// Seed with boundary conditions
	f.Add(int64(0), int64(1), int64(1), int64(2))

	f.Fuzz(func(t *testing.T, s1, e1, s2, e2 int64) {
		// Skip invalid ranges
		if s1 < 0 || e1 < 0 || s2 < 0 || e2 < 0 || s1 > e1 || s2 > e2 {
			return
		}

		ranges := []struct {
			Start *big.Int `json:"start"`
			End   *big.Int `json:"end"`
		}{
			{Start: big.NewInt(s1), End: big.NewInt(e1)},
			{Start: big.NewInt(s2), End: big.NewInt(e2)},
		}

		err := tokenization.ValidateNonOverlappingRanges(ranges, "testRanges")

		// Check if ranges overlap
		overlaps := (s1 <= e2 && s2 <= e1)
		if overlaps && s1 > 0 && e1 > 0 && s2 > 0 && e2 > 0 {
			if err == nil {
				t.Errorf("expected error for overlapping ranges [%d,%d] and [%d,%d], got nil", s1, e1, s2, e2)
			}
		}
	})
}

// FuzzCollectionIdOverflow tests for integer overflow in collection ID handling
func FuzzCollectionIdOverflow(f *testing.F) {
	// Seed with normal values
	f.Add(int64(1))
	// Seed with max int64
	f.Add(int64(9223372036854775807))
	// Seed with zero
	f.Add(int64(0))

	f.Fuzz(func(t *testing.T, id int64) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("panic on collection ID %d: %v", id, r)
			}
		}()

		// Test with very large big.Int values
		collectionId := big.NewInt(id)
		// Also test with a value that exceeds int64 max
		if id > 0 {
			// Create a value that's 10x larger
			largeId := new(big.Int).Mul(collectionId, big.NewInt(10000000000))
			err := tokenization.ValidateCollectionId(largeId)
			// Should either succeed or fail gracefully
			_ = err
		}

		err := tokenization.ValidateCollectionId(collectionId)
		_ = err
	})
}

// FuzzMultipleAddresses tests address validation with multiple addresses
func FuzzMultipleAddresses(f *testing.F) {
	// Seed with valid address bytes
	validAddr := []byte{0x74, 0x2d, 0x35, 0xCc, 0x66, 0x34, 0xC0, 0x53, 0x29, 0x25, 0xa3, 0xb8, 0x44, 0xBc, 0x9e, 0x75, 0x95, 0xf0, 0xbE, 0xb0}
	zeroAddr := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	f.Add(validAddr, zeroAddr)
	f.Add(zeroAddr, validAddr)
	f.Add(validAddr, validAddr)
	f.Add(zeroAddr, zeroAddr)

	f.Fuzz(func(t *testing.T, addr1Bytes, addr2Bytes []byte) {
		// Ensure we have exactly 20 bytes for each address
		if len(addr1Bytes) != 20 || len(addr2Bytes) != 20 {
			return
		}

		addr1 := common.BytesToAddress(addr1Bytes)
		addr2 := common.BytesToAddress(addr2Bytes)

		addrs := []common.Address{addr1, addr2}
		err := tokenization.ValidateAddresses(addrs, "testAddresses")

		// If any address is zero, validation should fail
		hasZero := addr1 == (common.Address{}) || addr2 == (common.Address{})
		if hasZero {
			if err == nil {
				t.Errorf("expected error for addresses containing zero address, got nil")
			}
		} else {
			// All valid addresses should pass
			if err != nil {
				t.Errorf("unexpected error for valid addresses: %v", err)
			}
		}
	})
}

// FuzzRangeEdgeCases tests edge cases in range validation
func FuzzRangeEdgeCases(f *testing.F) {
	// Test edge cases around 0 and max values
	f.Add(int64(0), int64(0))
	f.Add(int64(1), int64(1))
	f.Add(int64(0), int64(1))
	f.Add(int64(9223372036854775807), int64(9223372036854775807))

	f.Fuzz(func(t *testing.T, start, end int64) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("panic on range [%d, %d]: %v", start, end, r)
			}
		}()

		err := tokenization.ValidateBigIntRange(big.NewInt(start), big.NewInt(end), "testRange")

		// start > end should always fail
		if start > end {
			if err == nil {
				t.Errorf("expected error for invalid range [%d, %d], got nil", start, end)
			}
		}
		// negative values should fail
		if start < 0 || end < 0 {
			if err == nil {
				t.Errorf("expected error for negative range values [%d, %d], got nil", start, end)
			}
		}
	})
}

