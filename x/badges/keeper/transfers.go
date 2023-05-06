package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

// Handles a transfer from one address to another. If it can be a forceful transfer, it will forcefully transfer the balances and approvals. If it is a pending transfer, it will add it to the pending transfers.
func HandleTransfer(collection types.BadgeCollection, badgeIdRange *types.IdRange, fromUserBalance types.UserBalanceStore, toUserBalance types.UserBalanceStore, amount uint64, from string, to string, approvedBy string) (types.UserBalanceStore, types.UserBalanceStore, error) {
	err := *new(error)

	fromUserBalance, toUserBalance, err = ForcefulTransfer(collection, badgeIdRange, fromUserBalance, toUserBalance, amount, from, to, approvedBy)
	if err != nil {
		return types.UserBalanceStore{}, types.UserBalanceStore{}, err
	}

	return fromUserBalance, toUserBalance, nil
}

// Forceful transfers will transfer the balances and deduct from approvals directly without adding it to pending.
func ForcefulTransfer(collection types.BadgeCollection, badgeIdRange *types.IdRange, fromUserBalance types.UserBalanceStore, toUserBalance types.UserBalanceStore, amount uint64, from string, to string, approvedBy string) (types.UserBalanceStore, types.UserBalanceStore, error) {
	// 1. Check if the from address is frozen
	// 2. Remove approvals if approvedBy != from
	// 3. Deduct from "From" balance
	// 4. Add to "To" balance
	isManagerApprovedTransfer, err := AssertTransferAllowed(collection, from, to, approvedBy)
	if err != nil {
		return types.UserBalanceStore{}, types.UserBalanceStore{}, err
	}

	if !isManagerApprovedTransfer {
		fromUserBalance, err = DeductApprovals(fromUserBalance, collection, collection.CollectionId, badgeIdRange, from, to, approvedBy, amount)
		if err != nil {
			return types.UserBalanceStore{}, types.UserBalanceStore{}, err
		}
	}

	fromUserBalance, err = SubtractBalancesForIdRanges(fromUserBalance, []*types.IdRange{badgeIdRange}, amount)
	if err != nil {
		return types.UserBalanceStore{}, types.UserBalanceStore{}, err
	}

	toUserBalance, err = AddBalancesForIdRanges(toUserBalance, []*types.IdRange{badgeIdRange}, amount)
	if err != nil {
		return types.UserBalanceStore{}, types.UserBalanceStore{}, err
	}

	return fromUserBalance, toUserBalance, nil
}

// Deduct approvals from requester if requester != from
func DeductApprovals(UserBalance types.UserBalanceStore, collection types.BadgeCollection, collectionId uint64, rangeToDeduct *types.IdRange, from string, to string, requester string, amount uint64) (types.UserBalanceStore, error) {
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
func IsTransferAllowed(collection types.BadgeCollection, permissions types.Permissions, fromAddress string, toAddress string, approvedBy string) (bool, bool) {
	if approvedBy == collection.Manager {
		//Check if this is an approved transfer by the manager
		for _, managerApprovedTransfer := range collection.ManagerApprovedTransfers {
			fromFound := false
			toFound := false
			if managerApprovedTransfer.From.ManagerOptions == uint64(types.AddressOptions_IncludeManager) {
				fromFound = true
			} else if managerApprovedTransfer.From.ManagerOptions == uint64(types.AddressOptions_ExcludeManager) {
				fromFound = false
			} else {
				for _, addresss := range managerApprovedTransfer.From.Addresses {
					if addresss == fromAddress {
						fromFound = true
					}
				}
			}

			if managerApprovedTransfer.To.ManagerOptions == uint64(types.AddressOptions_IncludeManager) {
				toFound = true
			} else if managerApprovedTransfer.To.ManagerOptions == uint64(types.AddressOptions_ExcludeManager) {
				toFound = false
			} else {
				for _, addresss := range managerApprovedTransfer.To.Addresses {
					if addresss == toAddress {
						toFound = true
					}
				}
			}

			//As of now, we have checked whether the target addresses are in the list
			//If includeOnlySpecified is false, then we need to flip the values
			if !managerApprovedTransfer.To.IncludeOnlySpecified {
				toFound = !toFound
			}

			if !managerApprovedTransfer.From.IncludeOnlySpecified {
				fromFound = !fromFound
			}
			
			if fromFound && toFound {
				return true, true
			}
		}
	}	

	//Check if this is an allowed transfer
	for _, allowedTransfer := range collection.AllowedTransfers {
		fromFound := false
		toFound := false
		for _, addresss := range allowedTransfer.From.Addresses {
			if addresss == fromAddress {
				fromFound = true
			}
		}

		for _, addresss := range allowedTransfer.To.Addresses {
			if addresss == toAddress {
				toFound = true
			}
		}

		//As of now, we have checked whether the target addresses are in the list
		//If includeOnlySpecified is false, then we need to flip the values
		if !allowedTransfer.To.IncludeOnlySpecified {
			toFound = !toFound
		}

		if !allowedTransfer.From.IncludeOnlySpecified {
			fromFound = !fromFound
		}

		if fromFound && toFound {
			return true, false
		}
	}

	return false, false
}

// Returns an error if account is Frozen
func AssertTransferAllowed(collection types.BadgeCollection, from string, to string, approvedBy string) (bool, error) {
	permissions := types.GetPermissions(collection.Permissions)

	transferIsAllowed, isManagerApprovedTransfer := IsTransferAllowed(collection, permissions, from, to, approvedBy)
	if !transferIsAllowed {
		return false, ErrAddressFrozen
	}

	return isManagerApprovedTransfer, nil
}
