package types

import (
	sdkerrors "cosmossdk.io/errors"
)

// ValidateDynamicStoreChallenge validates a single DynamicStoreChallenge requirement
// This is a standalone validation function that can be reused across modules
func ValidateDynamicStoreChallenge(challenge *DynamicStoreChallenge, idx int) error {
	if challenge == nil {
		return sdkerrors.Wrapf(ErrInvalidRequest, "DynamicStoreChallenge at index %d is nil", idx)
	}

	if challenge.StoreId.IsNil() {
		return sdkerrors.Wrapf(ErrUintUnititialized, "StoreId is uninitialized for dynamic store challenge at index %d", idx)
	}

	if challenge.StoreId.IsZero() {
		return sdkerrors.Wrapf(ErrUintUnititialized, "StoreId cannot be zero for dynamic store challenge at index %d", idx)
	}

	return nil
}

// ValidateDynamicStoreChallengesList validates a list of DynamicStoreChallenge requirements
// It also checks for duplicate store IDs
func ValidateDynamicStoreChallengesList(challenges []*DynamicStoreChallenge) error {
	storeIds := make(map[string]bool)
	for idx, challenge := range challenges {
		if err := ValidateDynamicStoreChallenge(challenge, idx); err != nil {
			return err
		}

		// Check for duplicate store IDs
		storeIdStr := challenge.StoreId.String()
		if storeIds[storeIdStr] {
			return sdkerrors.Wrapf(ErrInvalidRequest, "duplicate dynamic store challenge storeId: %s at index %d", storeIdStr, idx)
		}
		storeIds[storeIdStr] = true
	}
	return nil
}

