package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
)

// ValidateUintId checks a collection / dynamic store id: it must be initialised and fit in a
// uint64 (store keys are fixed-width 8 bytes). Zero is allowed only where the caller resolves
// it (new collection or auto-prev).
func ValidateUintId(id sdkmath.Uint, allowZero bool) error {
	if id.IsNil() {
		return sdkerrors.Wrapf(ErrUintUnititialized, "id cannot be nil")
	}
	if !allowZero && id.IsZero() {
		return sdkerrors.Wrapf(ErrInvalidRequest, "id cannot be zero")
	}
	if !id.BigInt().IsUint64() {
		return sdkerrors.Wrapf(ErrUintGreaterThanMax, "id %s exceeds the maximum uint64 value", id.String())
	}
	return nil
}
