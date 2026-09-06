package types

import (
	"math/big"

	sdkmath "cosmossdk.io/math"
)

// Add adds two sdkmath.Uints. Note: This function does not check for overflow.
// Use SafeAddWithOverflowCheck if overflow protection is needed.
func Add(left sdkmath.Uint, right sdkmath.Uint) sdkmath.Uint {
	return left.Add(right)
}

// SafeAddWithOverflowCheck adds two sdkmath.Uints and returns an error if the result overflows.
// sdkmath.Uint panics once a value needs more than 256 bits, so the bound is checked up front.
func SafeAddWithOverflowCheck(left sdkmath.Uint, right sdkmath.Uint) (sdkmath.Uint, error) {
	sum := new(big.Int).Add(left.BigInt(), right.BigInt())
	if sum.BitLen() > sdkmath.MaxBitLen {
		return sdkmath.NewUint(0), ErrOverflow
	}
	return sdkmath.NewUintFromBigInt(sum), nil
}

// Safe subtracts two sdkmath.Uints and returns an error if the result underflows sdkmath.Uint.
func SafeSubtract(left sdkmath.Uint, right sdkmath.Uint) (sdkmath.Uint, error) {
	if right.GT(left) {
		return sdkmath.NewUint(0), ErrUnderflow
	}
	return left.Sub(right), nil
}

// SafeMulBoundedToUint64 multiplies two sdkmath.Uints and returns the product.
// Returns ErrOverflow if the product would exceed MaxUint64 (e.g. for Cosmos SDK coin amounts).
// Use this before converting amounts to bank coins or any uint64-bounded use.
func SafeMulBoundedToUint64(a, b sdkmath.Uint) (sdkmath.Uint, error) {
	if b.IsZero() {
		return sdkmath.NewUint(0), nil
	}
	maxSafeA := sdkmath.NewUint(MaxUint64Value).Quo(b)
	if a.GT(maxSafeA) {
		return sdkmath.NewUint(0), ErrOverflow
	}
	return a.Mul(b), nil
}
