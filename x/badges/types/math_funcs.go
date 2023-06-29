package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Safe adds two sdkmath.Uints and returns an error if the result overflows sdkmath.Uint.
func SafeAdd(left sdkmath.Uint, right sdkmath.Uint) (sdkmath.Uint, error) {
	//try to add the two numbers and catch panic
	defer func() (sdkmath.Uint, error) {
		if r := recover(); r != nil {
			return sdk.NewUint(0), ErrOverflow
		}
		return left.Add(right), nil
	}()

	return left.Add(right), nil
}

// Safe subtracts two sdkmath.Uints and returns an error if the result underflows sdkmath.Uint.
func SafeSubtract(left sdkmath.Uint, right sdkmath.Uint) (sdkmath.Uint, error) {
	// Cosmos SDK's Uint type does not overflow, so this check is not needed.
	if right.GT(left) {
		return sdk.NewUint(0), ErrUnderflow
	}
	return left.Sub(right), nil
}
