package types

import (
	sdkmath "cosmossdk.io/math"
)

// Safe adds two sdkmath.Uints
func SafeAdd(left sdkmath.Uint, right sdkmath.Uint) (sdkmath.Uint, error) {
	defer func() (sdkmath.Uint, error) {
		if r := recover(); r != nil {
			return sdkmath.NewUint(0), ErrOverflow
		}
		return left.Add(right), nil
	}()

	return left.Add(right), nil
}

// Safe subtracts two sdkmath.Uints and returns an error if the result underflows sdkmath.Uint.
func SafeSubtract(left sdkmath.Uint, right sdkmath.Uint) (sdkmath.Uint, error) {
	if right.GT(left) {
		return sdkmath.NewUint(0), ErrUnderflow
	}
	return left.Sub(right), nil
}
