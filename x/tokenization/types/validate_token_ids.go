package types

import (
	sdkerrors "cosmossdk.io/errors"
)

// ValidateBalancesWithinValidTokenIds ensures every token id range in balances is contained
// in validTokenIds.
func ValidateBalancesWithinValidTokenIds(balances []*Balance, validTokenIds []*UintRange) error {
	for _, balance := range balances {
		if balance == nil {
			continue
		}
		outside, _ := RemoveUintRangesFromUintRanges(validTokenIds, DeepCopyRanges(balance.TokenIds))
		if len(outside) > 0 {
			return sdkerrors.Wrapf(ErrInvalidRequest, "token ids %s are outside the collection's valid token ids", outside)
		}
	}
	return nil
}
