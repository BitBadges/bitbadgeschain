package v13

import (
	sdkmath "cosmossdk.io/math"
)

//Needed for custom gogoproto types

type (
	Int  = sdkmath.Int
	Uint = sdkmath.Uint
)

var (
	NewInt            = sdkmath.NewInt
	ZeroInt           = sdkmath.ZeroInt
	NewUint           = sdkmath.NewUint
	NewUintFromString = func(s string) Uint {
		if s == "" {
			return sdkmath.NewUint(0)
		}

		val := sdkmath.NewUintFromString(s)
		return val
	}
)

type (
	Dec = sdkmath.LegacyDec
)

var (
	NewDecWithPrec    = sdkmath.LegacyNewDecWithPrec
	NewDecFromInt     = sdkmath.LegacyNewDecFromInt
	NewDecFromStr     = sdkmath.LegacyNewDecFromStr
	MustNewDecFromStr = sdkmath.LegacyMustNewDecFromStr
)
