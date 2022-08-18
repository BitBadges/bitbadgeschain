package keeper

import (
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

// Appends the pending transfer to both parties balance informations and increments the nonce. Since we append, they will alwyas be sorted by the nonce.
func AppendPendingTransferForBothParties(fromUserBalanceInfo types.UserBalanceInfo, toUserBalanceInfo types.UserBalanceInfo, subbadgeRange *types.IdRange, to uint64, from uint64, amount uint64, approvedBy uint64, sentByFrom bool, expirationTime uint64, cantCancelBeforeTime uint64) (types.UserBalanceInfo, types.UserBalanceInfo, error) {
	if amount == 0 {
		return fromUserBalanceInfo, toUserBalanceInfo, ErrBalanceIsZero
	}

	if expirationTime != 0 && cantCancelBeforeTime > expirationTime {
		return fromUserBalanceInfo, toUserBalanceInfo, ErrCancelTimeIsGreaterThanExpirationTime
	}

	fromUserBalanceInfo.Pending = append(fromUserBalanceInfo.Pending, &types.PendingTransfer{
		SubbadgeRange:     subbadgeRange,
		Amount:            amount,
		ApprovedBy:        approvedBy,
		Sent:              sentByFrom, // different
		To:                to,
		From:              from,
		ThisPendingNonce:  fromUserBalanceInfo.PendingNonce, // this / other nonces are swapped
		OtherPendingNonce: toUserBalanceInfo.PendingNonce,
		ExpirationTime:    expirationTime,
		CantCancelBeforeTime: cantCancelBeforeTime,
	})

	toUserBalanceInfo.Pending = append(toUserBalanceInfo.Pending, &types.PendingTransfer{
		SubbadgeRange:     subbadgeRange,
		Amount:            amount,
		ApprovedBy:        approvedBy,
		Sent:              !sentByFrom, // different
		To:                to,
		From:              from,
		ThisPendingNonce:  toUserBalanceInfo.PendingNonce, // this / other nonces are swapped
		OtherPendingNonce: fromUserBalanceInfo.PendingNonce,
		ExpirationTime:    expirationTime,
		CantCancelBeforeTime: cantCancelBeforeTime,
	})

	err := *new(error)
	fromUserBalanceInfo.PendingNonce, err = SafeAdd(fromUserBalanceInfo.PendingNonce, 1) //nonces shouldn't reach the case where they overflow but this is just for safety
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
func RemovePending(userBalanceInfo types.UserBalanceInfo, thisNonce uint64, other_nonce uint64) (types.UserBalanceInfo, error) {
	idx, found :=  SearchPendingByNonce(userBalanceInfo.Pending, thisNonce)
	if !found {
		return userBalanceInfo, ErrPendingNotFound
	}

	newPending := []*types.PendingTransfer{}
	newPending = append(newPending, userBalanceInfo.Pending[:idx]...)
	newPending = append(newPending, userBalanceInfo.Pending[idx+1:]...)
	userBalanceInfo.Pending = newPending
	return userBalanceInfo, nil
}

// Prunes pending transfers that have expired
func PruneExpiredPending(currTime uint64, accountNum uint64, pending []*types.PendingTransfer) []*types.PendingTransfer {
	prunedPending := []*types.PendingTransfer{}
	for _, pendingTransfer := range pending {
		expired := pendingTransfer.ExpirationTime != 0 && pendingTransfer.ExpirationTime < currTime
		
		//Only prune received pending transfers that you can't accept anymore
		if expired && !pendingTransfer.Sent {
			continue
		} else {
			prunedPending = append(prunedPending, pendingTransfer)
		}
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
