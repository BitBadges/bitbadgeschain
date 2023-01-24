package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

// Handles a transfer from one address to another. If it can be a forceful transfer, it will forcefully transfer the balances and approvals. If it is a pending transfer, it will add it to the pending transfers.
func HandleTransfer(collection types.BadgeCollection, badgeIdRange *types.IdRange, fromUserBalance types.UserBalance, toUserBalance types.UserBalance, amount uint64, from uint64, to uint64, approvedBy uint64) (types.UserBalance, types.UserBalance, error) {
	err := *new(error)

	fromUserBalance, toUserBalance, err = ForcefulTransfer(collection, badgeIdRange, fromUserBalance, toUserBalance, amount, from, to, approvedBy)
	if err != nil {
		return types.UserBalance{}, types.UserBalance{}, err
	}

	return fromUserBalance, toUserBalance, nil
}

//Forceful transfers will transfer the balances and deduct from approvals directly without adding it to pending.
func ForcefulTransfer(collection types.BadgeCollection, badgeIdRange *types.IdRange, fromUserBalance types.UserBalance, toUserBalance types.UserBalance, amount uint64, from uint64, to uint64, approvedBy uint64) (types.UserBalance, types.UserBalance, error) {
	if amount == 0 {
		return types.UserBalance{}, types.UserBalance{}, nil
	}

	// 1. Check if the from address is frozen
	// 2. Remove approvals if approvedBy != from
	// 3. Deduct from "From" balance
	// 4. Add to "To" balance
	err := AssertTransferAllowed(collection, from, to, approvedBy)
	if err != nil {
		return types.UserBalance{}, types.UserBalance{}, err
	}

	fromUserBalance, err = DeductApprovals(fromUserBalance, collection, collection.CollectionId, badgeIdRange, from, to, approvedBy, amount)
	if err != nil {
		return types.UserBalance{}, types.UserBalance{}, err
	}

	fromUserBalance, err = SubtractBalancesForIdRanges(fromUserBalance, []*types.IdRange{badgeIdRange}, amount)
	if err != nil {
		return types.UserBalance{}, types.UserBalance{}, err
	}

	toUserBalance, err = AddBalancesForIdRanges(toUserBalance, []*types.IdRange{badgeIdRange}, amount)
	if err != nil {
		return types.UserBalance{}, types.UserBalance{}, err
	}

	return fromUserBalance, toUserBalance, nil
}

// Deduct approvals from requester if requester != from
func DeductApprovals(UserBalance types.UserBalance, collection types.BadgeCollection, collectionId uint64, rangeToDeduct *types.IdRange, from uint64, to uint64, requester uint64, amount uint64) (types.UserBalance, error) {
	newUserBalance := UserBalance

	if from != requester {
		err := *new(error)
		newUserBalance, err = RemoveBalanceFromApproval(newUserBalance, amount, requester, []*types.IdRange{rangeToDeduct})
		if err != nil {
			return UserBalance, err
		}
	}

	return newUserBalance, nil
}

// Checks if account is frozen or not.
func IsTransferAllowed(collection types.BadgeCollection, permissions types.Permissions, fromAddress uint64, toAddress uint64, approvedBy uint64) bool {
	if approvedBy == collection.Manager {
		//Check if this is an approved transfer by the manager
		for _, managerApprovedTransfer := range collection.ManagerApprovedTransfers {
			fromFound := false
			toFound := false
			if managerApprovedTransfer.From.Options == types.AddressOptions_IncludeManager {
				fromFound = true
			} else if managerApprovedTransfer.From.Options == types.AddressOptions_ExcludeManager {
				fromFound = false
			} else {
				_, fromFound = SearchIdRangesForId(fromAddress, managerApprovedTransfer.From.AccountNums)
			}

			if managerApprovedTransfer.To.Options == types.AddressOptions_IncludeManager {
				toFound = true
			} else if managerApprovedTransfer.To.Options == types.AddressOptions_ExcludeManager {
				toFound = false
			} else {
				_, toFound = SearchIdRangesForId(toAddress, managerApprovedTransfer.To.AccountNums)
			}

			if fromFound && toFound {
				return true
			}
		}
	}

	for _, disallowedTransfer := range collection.DisallowedTransfers {
		_, fromFound := SearchIdRangesForId(fromAddress, disallowedTransfer.From.AccountNums)
		_, toFound := SearchIdRangesForId(toAddress, disallowedTransfer.To.AccountNums)

		if fromFound && toFound {
			return false
		}
	}

	return true
}

// Returns an error if account is Frozen
func AssertTransferAllowed(collection types.BadgeCollection, from uint64, to uint64, approvedBy uint64) error {
	permissions := types.GetPermissions(collection.Permissions)

	transferIsAllowed := IsTransferAllowed(collection, permissions, from, to, approvedBy)
	if !transferIsAllowed {
		return ErrAddressFrozen
	}

	return nil
}
