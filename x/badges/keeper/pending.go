package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

// Appends the pending transfer to both parties balance informations and increments the nonce. Since we append, they will alwyas be sorted by the nonce.
func AppendPendingTransferForBothParties(ctx sdk.Context, fromUserBalanceInfo types.UserBalanceInfo, toUserBalanceInfo types.UserBalanceInfo, subbadgeRange types.IdRange, to uint64, from uint64, amount uint64, approvedBy uint64, sentByFrom bool, expirationTime uint64) (types.UserBalanceInfo, types.UserBalanceInfo, error) {
	if amount == 0 {
		return fromUserBalanceInfo, toUserBalanceInfo, ErrBalanceIsZero
	}

	fromUserBalanceInfo.Pending = append(fromUserBalanceInfo.Pending, &types.PendingTransfer{
		SubbadgeRange:     &subbadgeRange,
		Amount:            amount,
		ApprovedBy:        approvedBy,
		Sent:              sentByFrom, // different
		To:                to,
		From:              from,
		ThisPendingNonce:  fromUserBalanceInfo.PendingNonce, // this / other nonces are swapped
		OtherPendingNonce: toUserBalanceInfo.PendingNonce,
		ExpirationTime:    expirationTime,
	})

	toUserBalanceInfo.Pending = append(toUserBalanceInfo.Pending, &types.PendingTransfer{
		SubbadgeRange:     &subbadgeRange,
		Amount:            amount,
		ApprovedBy:        approvedBy,
		Sent:              !sentByFrom, // different
		To:                to,
		From:              from,
		ThisPendingNonce:  toUserBalanceInfo.PendingNonce, // this / other nonces are swapped
		OtherPendingNonce: fromUserBalanceInfo.PendingNonce,
		ExpirationTime:    expirationTime,
	})
	err := *new(error)
	fromUserBalanceInfo.PendingNonce, err = SafeAdd(fromUserBalanceInfo.PendingNonce, 1)
	if err != nil {
		return fromUserBalanceInfo, toUserBalanceInfo, err
	}
	
	toUserBalanceInfo.PendingNonce, err = SafeAdd(toUserBalanceInfo.PendingNonce, 1)
	if err != nil {
		return fromUserBalanceInfo, toUserBalanceInfo, err
	}

	return fromUserBalanceInfo, toUserBalanceInfo, nil
}

//Removes pending transfer from the userBalanceInfo.
func RemovePending(ctx sdk.Context, userBalanceInfo types.UserBalanceInfo, thisNonce uint64, other_nonce uint64) (types.UserBalanceInfo, error) {
	pending := userBalanceInfo.Pending
	low := 0
	high := len(pending) - 1

	foundIdx := -1
	for low <= high {
		median := int(uint(low+high) >> 1)
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
		return userBalanceInfo, ErrPendingNotFound
	}

	newPending := []*types.PendingTransfer{}
	newPending = append(newPending, pending[:foundIdx]...)
	newPending = append(newPending, pending[foundIdx+1:]...)
	userBalanceInfo.Pending = newPending

	return userBalanceInfo, nil
}

// Prunes pending transfers that have expired
func PruneExpiredPending(currTime uint64, accountNum uint64, pending []*types.PendingTransfer) []*types.PendingTransfer {
	prunedPending := make([]*types.PendingTransfer, 0)
	for _, pendingTransfer := range pending {
		//TODO: prune expired pending transfers
		// if pendingTransfer.ExpirationTime != 0 && pendingTransfer.ExpirationTime < currTime && !pendingTransfer.Sent {
		// 	continue
		// } else if pendingTransfer.ExpirationTime != 0 && pendingTransfer.ExpirationTime < currTime && pendingTransfer.Sent && pendingTransfer.From == accountNum {
		// 	continue
		// } else {
		prunedPending = append(prunedPending, pendingTransfer)
		// }
	}
	return prunedPending
}

// Binary search pending by nonce
func SearchPendingByNonce(pending []*types.PendingTransfer, nonce uint64) (int, bool) {
	low := 0
	high := len(pending) - 1

	for low <= high {
		median := int(uint(low+high) >> 1)
		currPending := pending[median]

		if currPending.ThisPendingNonce == nonce {
			return median, true
		} else if currPending.ThisPendingNonce > nonce {
			high = median - 1
		} else {
			low = median + 1
		}
	}

	return -1, false
}
