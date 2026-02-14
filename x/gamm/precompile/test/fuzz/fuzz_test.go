package gamm

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	gamm "github.com/bitbadges/bitbadgeschain/x/gamm/precompile"
)

// FuzzValidatePoolId fuzzes the ValidatePoolId function
func FuzzValidatePoolId(f *testing.F) {
	// Seed with valid pool IDs
	f.Add(uint64(1))
	f.Add(uint64(0))
	f.Add(uint64(1000000))
	f.Add(uint64(18446744073709551615)) // max uint64

	f.Fuzz(func(t *testing.T, poolId uint64) {
		err := gamm.ValidatePoolId(poolId)

		// Zero pool ID should fail
		if poolId == 0 {
			if err == nil {
				t.Errorf("expected error for poolId 0, got nil")
			}
		} else {
			// Non-zero pool IDs should pass
			if err != nil {
				t.Errorf("unexpected error for valid poolId %d: %v", poolId, err)
			}
		}
	})
}

// FuzzValidateShareAmount fuzzes the ValidateShareAmount function
func FuzzValidateShareAmount(f *testing.F) {
	// Seed with valid amounts
	f.Add(int64(1))
	f.Add(int64(0))
	f.Add(int64(-1))
	f.Add(int64(1000000))
	// Max int256 (2^255 - 1)
	maxInt256 := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 255), big.NewInt(1))
	f.Add(maxInt256.Int64())

	f.Fuzz(func(t *testing.T, amount int64) {
		amountBig := big.NewInt(amount)
		err := gamm.ValidateShareAmount(amountBig, "testShareAmount")

		// Zero or negative amounts should fail
		if amount <= 0 {
			if err == nil {
				t.Errorf("expected error for shareAmount %d, got nil", amount)
			}
		} else {
			// Positive amounts should pass (unless they overflow)
			// CheckOverflow will catch overflow cases
			if err != nil && amount > 0 {
				// Only error if it's not an overflow error (overflow is expected for very large values)
				if amount < 0x7fffffffffffffff { // reasonable max for int64
					t.Errorf("unexpected error for valid shareAmount %d: %v", amount, err)
				}
			}
		}
	})
}

// FuzzValidateCoin fuzzes the ValidateCoin function
func FuzzValidateCoin(f *testing.F) {
	// Seed with valid coin
	f.Add("uatom", int64(1000))
	// Seed with empty denom
	f.Add("", int64(1000))
	// Seed with negative amount
	f.Add("uatom", int64(-1))
	// Seed with zero amount
	f.Add("uatom", int64(0))
	// Seed with very long denom
	longDenom := make([]byte, gamm.MaxStringLength+1)
	for i := range longDenom {
		longDenom[i] = 'a'
	}
	f.Add(string(longDenom), int64(1000))

	f.Fuzz(func(t *testing.T, denom string, amount int64) {
		coin := struct {
			Denom  string   `json:"denom"`
			Amount *big.Int `json:"amount"`
		}{
			Denom:  denom,
			Amount: big.NewInt(amount),
		}

		err := gamm.ValidateCoin(coin, "testCoin")

		// Empty denom or invalid amount should fail
		if denom == "" || amount <= 0 {
			if err == nil {
				t.Errorf("expected error for coin {denom: %q, amount: %d}, got nil", denom, amount)
			}
		} else if len(denom) <= gamm.MaxStringLength && amount > 0 {
			// Valid coins should pass (unless amount overflows)
			if err != nil && amount < 0x7fffffffffffffff {
				t.Errorf("unexpected error for valid coin {denom: %q, amount: %d}: %v", denom, amount, err)
			}
		}
	})
}

// FuzzValidateRoutes fuzzes the ValidateRoutes function
func FuzzValidateRoutes(f *testing.F) {
	// Seed with valid route
	f.Add(uint64(1), "uatom")
	// Seed with zero pool ID
	f.Add(uint64(0), "uatom")
	// Seed with empty denom
	f.Add(uint64(1), "")
	// Seed with very long denom
	longDenom := make([]byte, gamm.MaxStringLength+1)
	for i := range longDenom {
		longDenom[i] = 'a'
	}
	f.Add(uint64(1), string(longDenom))

	f.Fuzz(func(t *testing.T, poolId uint64, tokenOutDenom string) {
		routes := []struct {
			PoolId        uint64 `json:"poolId"`
			TokenOutDenom string `json:"tokenOutDenom"`
		}{
			{
				PoolId:        poolId,
				TokenOutDenom: tokenOutDenom,
			},
		}

		err := gamm.ValidateRoutes(routes, "testRoutes")

		// Zero pool ID or empty denom should fail
		if poolId == 0 || tokenOutDenom == "" {
			if err == nil {
				t.Errorf("expected error for route {poolId: %d, tokenOutDenom: %q}, got nil", poolId, tokenOutDenom)
			}
		} else if len(tokenOutDenom) <= gamm.MaxStringLength {
			// Valid routes should pass
			if err != nil {
				t.Errorf("unexpected error for valid route {poolId: %d, tokenOutDenom: %q}: %v", poolId, tokenOutDenom, err)
			}
		}
	})
}

// FuzzCheckOverflow fuzzes the CheckOverflow function
func FuzzCheckOverflow(f *testing.F) {
	// Seed with valid values
	f.Add(int64(0))
	f.Add(int64(1))
	f.Add(int64(-1))
	f.Add(int64(1000000))
	// Max int256 (2^255 - 1)
	maxInt256 := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 255), big.NewInt(1))
	f.Add(maxInt256.Int64())

	f.Fuzz(func(t *testing.T, value int64) {
		valueBig := big.NewInt(value)
		err := gamm.CheckOverflow(valueBig, "testValue")

		// Negative values should fail
		if value < 0 {
			if err == nil {
				t.Errorf("expected error for negative value %d, got nil", value)
			}
		} else {
			// Non-negative values should pass (unless they overflow int256)
			// Max int256 is 2^255 - 1, which is much larger than int64 max
			// So int64 values should never overflow int256
			if err != nil && value >= 0 {
				// Only error if it's actually an overflow (which shouldn't happen for int64)
				maxInt256 := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 255), big.NewInt(1))
				if valueBig.Cmp(maxInt256) <= 0 {
					t.Errorf("unexpected error for valid value %d: %v", value, err)
				}
			}
		}
	})
}

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
		err := gamm.ValidateAddress(addr, "testAddress")

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

// FuzzValidateStringLength fuzzes the ValidateStringLength function
func FuzzValidateStringLength(f *testing.F) {
	// Seed with empty string
	f.Add("")
	// Seed with short string
	f.Add("uatom")
	// Seed with max length string
	maxString := make([]byte, gamm.MaxStringLength)
	for i := range maxString {
		maxString[i] = 'a'
	}
	f.Add(string(maxString))
	// Seed with over max length string
	overMaxString := make([]byte, gamm.MaxStringLength+1)
	for i := range overMaxString {
		overMaxString[i] = 'a'
	}
	f.Add(string(overMaxString))

	f.Fuzz(func(t *testing.T, s string) {
		err := gamm.ValidateStringLength(s, "testString")

		// Strings longer than max should fail
		if len(s) > gamm.MaxStringLength {
			if err == nil {
				t.Errorf("expected error for string length %d (max: %d), got nil", len(s), gamm.MaxStringLength)
			}
		} else {
			// Strings within limit should pass
			if err != nil {
				t.Errorf("unexpected error for valid string length %d: %v", len(s), err)
			}
		}
	})
}

// FuzzValidatePagination fuzzes the ValidatePagination function
func FuzzValidatePagination(f *testing.F) {
	// Seed with valid pagination
	f.Add(int64(0), int64(10))
	// Seed with negative offset
	f.Add(int64(-1), int64(10))
	// Seed with zero limit
	f.Add(int64(0), int64(0))
	// Seed with negative limit
	f.Add(int64(0), int64(-1))
	// Seed with over max limit (1000)
	f.Add(int64(0), int64(1001))
	// Seed with max limit
	f.Add(int64(0), int64(1000))

	f.Fuzz(func(t *testing.T, offset int64, limit int64) {
		offsetBig := big.NewInt(offset)
		limitBig := big.NewInt(limit)
		err := gamm.ValidatePagination(offsetBig, limitBig)

		// Negative offset, zero/negative limit, or limit > 1000 should fail
		if offset < 0 || limit <= 0 || limit > 1000 {
			if err == nil {
				t.Errorf("expected error for pagination {offset: %d, limit: %d}, got nil", offset, limit)
			}
		} else {
			// Valid pagination should pass
			if err != nil {
				t.Errorf("unexpected error for valid pagination {offset: %d, limit: %d}: %v", offset, limit, err)
			}
		}
	})
}

