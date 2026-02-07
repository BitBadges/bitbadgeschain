package tokenization

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	tokenization "github.com/bitbadges/bitbadgeschain/x/evm/precompiles/tokenization"
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

