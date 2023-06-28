package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// Safe adds two sdk.Uints and returns an error if the result overflows sdk.Uint.
func SafeAdd(left sdk.Uint, right sdk.Uint) (sdk.Uint, error) {
	//try to add the two numbers and catch panic
	defer func() (sdk.Uint, error) {
		if r := recover(); r != nil {
			return sdk.NewUint(0), ErrOverflow
		}
		return left.Add(right), nil
	}()

	return left.Add(right), nil
}

// Safe subtracts two sdk.Uints and returns an error if the result underflows sdk.Uint.
func SafeSubtract(left sdk.Uint, right sdk.Uint) (sdk.Uint, error) {
	// Cosmos SDK's Uint type does not overflow, so this check is not needed.
	if right.GT(left) {
		return sdk.NewUint(0), ErrUnderflow
	}
	return left.Sub(right), nil
}
