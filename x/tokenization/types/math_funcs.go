package types

import (
	sdkmath "cosmossdk.io/math"
)

// Add adds two sdkmath.Uints. Note: This function does not check for overflow.
// Use SafeAddWithOverflowCheck if overflow protection is needed.
func Add(left sdkmath.Uint, right sdkmath.Uint) sdkmath.Uint {
	return left.Add(right)
}

// SafeAddWithOverflowCheck adds two sdkmath.Uints and returns an error if the result overflows.
// Note: sdkmath.Uint uses big.Int internally which handles arbitrary precision, but we keep this
// check for defense-in-depth in case the underlying implementation changes.
func SafeAddWithOverflowCheck(left sdkmath.Uint, right sdkmath.Uint) (sdkmath.Uint, error) {
	result := left.Add(right)
	// If result is less than either operand, overflow occurred (security fix: changed && to ||)
	if result.LT(left) || result.LT(right) {
		return sdkmath.NewUint(0), ErrOverflow
	}
	return result, nil
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
