package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

// Appends the pending transfer to both parties balance informations and increments the nonce. Since we append, they will alwyas be sorted by the nonce.
func (k Keeper) AppendPendingTransferForBothParties(ctx sdk.Context, fromBadgeBalanceInfo types.BadgeBalanceInfo, toBadgeBalanceInfo types.BadgeBalanceInfo, subbadgeRange types.NumberRange, to uint64, from uint64, amount uint64, approvedBy uint64, sentByFrom bool, expirationTime uint64) (types.BadgeBalanceInfo, types.BadgeBalanceInfo, error) {
	if amount == 0 {
		return fromBadgeBalanceInfo, toBadgeBalanceInfo, ErrBalanceIsZero
	}

	fromBadgeBalanceInfo.Pending = append(fromBadgeBalanceInfo.Pending, &types.PendingTransfer{
		SubbadgeRange:     &subbadgeRange,
		Amount:            amount,
		ApprovedBy:        approvedBy,
		SendRequest:       sentByFrom, // different 
		To:                to,
		From:              from,
		ThisPendingNonce:  fromBadgeBalanceInfo.PendingNonce, // this / other nonces are swapped 
		OtherPendingNonce: toBadgeBalanceInfo.PendingNonce,
		ExpirationTime:    expirationTime,
	})

	toBadgeBalanceInfo.Pending = append(toBadgeBalanceInfo.Pending, &types.PendingTransfer{
		SubbadgeRange:     &subbadgeRange,
		Amount:            amount,
		ApprovedBy:        approvedBy,
		SendRequest:       !sentByFrom, // different 
		To:                to,
		From:              from,
		ThisPendingNonce:  toBadgeBalanceInfo.PendingNonce, // this / other nonces are swapped 
		OtherPendingNonce: fromBadgeBalanceInfo.PendingNonce,
		ExpirationTime:    expirationTime,
	})

	fromBadgeBalanceInfo.PendingNonce += 1
	toBadgeBalanceInfo.PendingNonce += 1

	return fromBadgeBalanceInfo, toBadgeBalanceInfo, nil
}

//Removes pending transfer from the badgeBalanceInfo. 
func (k Keeper) RemovePending(ctx sdk.Context, badgeBalanceInfo types.BadgeBalanceInfo, thisNonce uint64, other_nonce uint64) (types.BadgeBalanceInfo, error) {
	pending := badgeBalanceInfo.Pending
	low := 0
	high := len(pending) - 1

	foundIdx := -1
	for low <= high {
		median := int(uint(low + high) >> 1)
		currPending := pending[median]
		if currPending.ThisPendingNonce == thisNonce && currPending.OtherPendingNonce == other_nonce {
			foundIdx = median
			break
		} else if currPending.ThisPendingNonce > thisNonce {
			high = median - 1
		} else {
			low = median + 1
		}
	}

	if foundIdx == -1 {
		return badgeBalanceInfo, ErrPendingNotFound
	}
	
	newPending := []*types.PendingTransfer{}
	newPending = append(newPending, pending[:foundIdx]...)
	newPending = append(newPending, pending[foundIdx + 1:]...)
	badgeBalanceInfo.Pending = newPending

	return badgeBalanceInfo, nil
}
