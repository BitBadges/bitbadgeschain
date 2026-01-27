package types

import (
	sdkerrors "cosmossdk.io/errors"
)

// ValidateMustOwnTokens validates a single MustOwnTokens requirement
// This is a standalone validation function that can be reused across modules
func ValidateMustOwnTokens(mustOwnToken *MustOwnTokens, idx int) error {
	if mustOwnToken == nil {
		return sdkerrors.Wrapf(ErrInvalidRequest, "MustOwnTokens requirement at index %d is nil", idx)
	}

	if mustOwnToken.CollectionId.IsNil() || mustOwnToken.CollectionId.IsZero() {
		return sdkerrors.Wrapf(ErrUintUnititialized, "CollectionId is uninitialized or zero for requirement at index %d", idx)
	}

	if mustOwnToken.AmountRange == nil {
		return sdkerrors.Wrapf(ErrInvalidRequest, "AmountRange is required for requirement at index %d", idx)
	}

	// Validate amount range
	if err := ValidateRangesAreValid([]*UintRange{mustOwnToken.AmountRange}, true, true); err != nil {
		return sdkerrors.Wrapf(err, "invalid amount range for requirement at index %d", idx)
	}

	// Validate token IDs if provided
	if len(mustOwnToken.TokenIds) > 0 {
		if err := ValidateRangesAreValid(mustOwnToken.TokenIds, false, false); err != nil {
			return sdkerrors.Wrapf(err, "invalid token IDs for requirement at index %d", idx)
		}
	}

	// Validate ownership times if not using override
	if !mustOwnToken.OverrideWithCurrentTime {
		if len(mustOwnToken.OwnershipTimes) == 0 {
			return sdkerrors.Wrapf(ErrInvalidRequest, "OwnershipTimes must be set or OverrideWithCurrentTime must be true for requirement at index %d", idx)
		}
		if err := ValidateRangesAreValid(mustOwnToken.OwnershipTimes, false, false); err != nil {
			return sdkerrors.Wrapf(err, "invalid ownership times for requirement at index %d", idx)
		}
	}

	return nil
}

// ValidateMustOwnTokensList validates a list of MustOwnTokens requirements
func ValidateMustOwnTokensList(mustOwnTokens []*MustOwnTokens) error {
	for idx, mustOwnToken := range mustOwnTokens {
		if err := ValidateMustOwnTokens(mustOwnToken, idx); err != nil {
			return err
		}
	}
	return nil
}

