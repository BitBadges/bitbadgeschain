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
func SafeAddWithOverflowCheck(left sdkmath.Uint, right sdkmath.Uint) (sdkmath.Uint, error) {
	result := left.Add(right) // SDK uints already handle overflows
	if result.LT(left) && result.LT(right) {
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
