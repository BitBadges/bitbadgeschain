package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Handles a transfer from one address to another. If it can be a forceful transfer, it will forcefully transfer the balances and approvals. If it is a pending transfer, it will add it to the pending transfers.
func HandleTransfer(collection types.BadgeCollection, badgeIds []*types.IdRange, fromUserBalance types.UserBalanceStore, toUserBalance types.UserBalanceStore, amount sdk.Uint, from string, to string, approvedBy string) (types.UserBalanceStore, types.UserBalanceStore, error) {
	err := *new(error)

	fromUserBalance, toUserBalance, err = ForcefulTransfer(collection, badgeIds, fromUserBalance, toUserBalance, amount, from, to, approvedBy)
	if err != nil {
		return types.UserBalanceStore{}, types.UserBalanceStore{}, err
	}

	return fromUserBalance, toUserBalance, nil
}

// Forceful transfers will transfer the balances and deduct from approvals directly without adding it to pending.
func ForcefulTransfer(collection types.BadgeCollection, badgeIds []*types.IdRange, fromUserBalance types.UserBalanceStore, toUserBalance types.UserBalanceStore, amount sdk.Uint, from string, to string, approvedBy string) (types.UserBalanceStore, types.UserBalanceStore, error) {
	// 1. Check if the from address is frozen
	// 2. Remove approvals if approvedBy != from and not managerApprovedTransfer
	// 3. Deduct from "From" balance
	// 4. Add to "To" balance
	err := *new(error)

	isAllowedTransfer, isManagerApprovedTransfer := AssertTransferAllowed(collection, from, to, approvedBy)
	if !isAllowedTransfer {
		return types.UserBalanceStore{}, types.UserBalanceStore{}, ErrAddressFrozen
	}

	if !isManagerApprovedTransfer {
		fromUserBalance, err = DeductApprovalsIfNeeded(fromUserBalance, collection, collection.CollectionId, badgeIds, from, to, approvedBy, amount)
		if err != nil {
			return types.UserBalanceStore{}, types.UserBalanceStore{}, err
		}
	}

	fromUserBalance.Balances, err = SubtractBalancesForIdRanges(fromUserBalance.Balances, badgeIds, amount)
	if err != nil {
		return types.UserBalanceStore{}, types.UserBalanceStore{}, err
	}

	toUserBalance.Balances, err = AddBalancesForIdRanges(toUserBalance.Balances, badgeIds, amount)
	if err != nil {
		return types.UserBalanceStore{}, types.UserBalanceStore{}, err
	}

	return fromUserBalance, toUserBalance, nil
}

// Deduct approvals from requester if requester != from
func DeductApprovalsIfNeeded(UserBalance types.UserBalanceStore, collection types.BadgeCollection, collectionId sdk.Uint, badgeIds []*types.IdRange, from string, to string, requester string, amount sdk.Uint) (types.UserBalanceStore, error) {
	newUserBalance := UserBalance

	if from != requester {
		err := *new(error)
		newUserBalance.Approvals, err = RemoveBalanceFromApproval(newUserBalance.Approvals, amount, requester, badgeIds)
		if err != nil {
			return UserBalance, err
		}
	}

	return newUserBalance, nil
}

//Within each address mapping, we can specify the manager options to include/exclude the manager address
//from the Addresses field or do nothing.
//This is useful because the manager address is not always a fixed address and can transfer.
func HandleManagerOptions(addressMapping *types.AddressesMapping, managerAddress string) {
	
	if addressMapping.ManagerOptions.Equal(sdk.NewUint(uint64(types.AddressOptions_IncludeManager))) {
		addressMapping.Addresses = append(addressMapping.Addresses, managerAddress)
	} else if addressMapping.ManagerOptions.Equal(sdk.NewUint(uint64(types.AddressOptions_ExcludeManager))) {
		//Remove from Addresses
		newAddresses := []string{}
		for _, address := range addressMapping.Addresses {
			if address != managerAddress {
				newAddresses = append(newAddresses, address)
			}
		}
		addressMapping.Addresses = newAddresses
	}
}

// Checks if the from and to addresses are in the transfer mapping.
// Handles the manager options for the from and to addresses.
// If includeOnlySpecified is true, then we check if the address is in the Addresses field.
// If includeOnlySpecified is false, then we check if the address is NOT in the Addresses field.
func CheckIfInTransferMapping(transferMapping *types.TransferMapping, from string, to string, managerAddress string) (bool,bool) {
	fromFound := false
	toFound := false

	HandleManagerOptions(transferMapping.From, managerAddress)
	HandleManagerOptions(transferMapping.To, managerAddress)

	for _, address := range transferMapping.From.Addresses {
		if address == from {
			fromFound = true
		}
	}

	for _, address := range transferMapping.To.Addresses {
		if address == to {
			toFound = true
		}
	}

	if !transferMapping.From.IncludeOnlySpecified {
		fromFound = !fromFound
	}

	if !transferMapping.To.IncludeOnlySpecified {
		toFound = !toFound
	}

	return fromFound, toFound
}


// Checks if account is frozen or not.
func IsTransferAllowed(collection types.BadgeCollection, fromAddress string, toAddress string, approvedBy string) (bool, bool) {
	if approvedBy == collection.Manager {
		for _, managerApprovedTransfer := range collection.ManagerApprovedTransfers {
		  fromFound, toFound := CheckIfInTransferMapping(managerApprovedTransfer, fromAddress, toAddress, collection.Manager)
			
			//If both are true, then this is a manager approved transfer
			if fromFound && toFound {
				return true, true
			}
		}

		//If not, we handle it as a normal transfer
	}

	for _, allowedTransfer := range collection.AllowedTransfers {
		fromFound, toFound := CheckIfInTransferMapping(allowedTransfer, fromAddress, toAddress, collection.Manager)

		if fromFound && toFound {
			return true, false
		}
	}

	return false, false
}

// Returns an error if account is Frozen
func AssertTransferAllowed(collection types.BadgeCollection, from string, to string, approvedBy string) (bool, bool) {
	transferIsAllowed, isManagerApprovedTransfer := IsTransferAllowed(collection, from, to, approvedBy)
	return transferIsAllowed, isManagerApprovedTransfer
}
