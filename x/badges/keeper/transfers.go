package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Forceful transfers will transfer the balances and deduct from approvals directly without adding it to pending.
func HandleTransfer(ctx sdk.Context, collection types.BadgeCollection, badgeIds []*types.IdRange, fromUserBalance types.UserBalanceStore, toUserBalance types.UserBalanceStore, amount sdk.Uint, from string, to string, initiatedBy string) (types.UserBalanceStore, types.UserBalanceStore, error) {
	// 1. Check if the from address is frozen
	// 2. Remove approvals if initiatedBy != from and not managerApprovedTransfer
	// 3. Deduct from "From" balance
	// 4. Add to "To" balance
	err := *new(error)

	isAllowedTransfer, requiresApproval := IsTransferAllowed(ctx, badgeIds, collection, from, to, initiatedBy)
	if !isAllowedTransfer {
		return types.UserBalanceStore{}, types.UserBalanceStore{}, ErrAddressFrozen
	}

	if requiresApproval {
		fromUserBalance, err = DeductApprovalsIfNeeded(ctx, fromUserBalance, collection, collection.CollectionId, badgeIds, from, to, initiatedBy, amount)
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
func DeductApprovalsIfNeeded(ctx sdk.Context, UserBalance types.UserBalanceStore, collection types.BadgeCollection, collectionId sdk.Uint, badgeIds []*types.IdRange, from string, to string, requester string, amount sdk.Uint) (types.UserBalanceStore, error) {
	newUserBalance := UserBalance

	if from != requester {
		err := *new(error)
		idx, found := SearchApprovals(requester, newUserBalance.Approvals)
		if !found {
			return UserBalance, ErrApprovalForAddressDoesntExist
		}

		approval := newUserBalance.Approvals[idx]
		validTime := false
		for _, timeInterval := range approval.TimeIntervals {
			currTime := sdk.NewUint(uint64(ctx.BlockTime().UnixMilli()))
			if currTime.GT(timeInterval.Start) && currTime.LT(timeInterval.End) {
				validTime = true
			}
		}

		if !validTime {
			return UserBalance, ErrInvalidTime
		}

		newUserBalance.Approvals, err = RemoveBalanceFromApproval(newUserBalance.Approvals, amount, requester, badgeIds)
		if err != nil {
			return UserBalance, err
		}
	}

	return newUserBalance, nil
}

func CheckMappingAddresses(addressMapping *types.AddressMapping, addressToCheck string, managerAddress string) bool {
	found := false

	for _, address := range addressMapping.Addresses {
		if address == addressToCheck {
			found = true
		}

		//Support the manager alias
		if address == "Manager" && (addressToCheck == managerAddress || addressToCheck == "Manager") {
			found = true
		}
	}

	if !addressMapping.IncludeOnlySpecified {
		found = !found
	}

	return found
}

// Checks if the from and to addresses are in the transfer approvedTransfer.
// Handles the manager options for the from and to addresses.
// If includeOnlySpecified is true, then we check if the address is in the Addresses field.
// If includeOnlySpecified is false, then we check if the address is NOT in the Addresses field.

// Note addresses matching does not mean the transfer is allowed. It just means the addresses match.
// All other criteria must also be met.
func CheckIfAddressesMatch(collectionApprovedTransfer *types.CollectionApprovedTransfer, from string, to string, initiatedBy string, managerAddress string) bool {
	if from == "Mint" && initiatedBy == "Mint" {
		return collectionApprovedTransfer.IncludeMints && CheckMappingAddresses(collectionApprovedTransfer.To, to, managerAddress)
	}

	fromFound := CheckMappingAddresses(collectionApprovedTransfer.From, from, managerAddress)
	toFound := CheckMappingAddresses(collectionApprovedTransfer.To, to, managerAddress)
	initiatedByFound := CheckMappingAddresses(collectionApprovedTransfer.InitiatedBy, initiatedBy, managerAddress)

	return fromFound && toFound && initiatedByFound
}

// Checks if account is frozen or not.
func IsTransferAllowed(ctx sdk.Context, badgeIds []*types.IdRange, collection types.BadgeCollection, fromAddress string, toAddress string, initiatedBy string) (bool, bool) {
	for _, allowedTransfer := range collection.ApprovedTransfers {
		//Check if addresses match. Handles as a "Mint" tx, if from and initiatedBy are both "Mint"
		addressesMatch := CheckIfAddressesMatch(allowedTransfer, fromAddress, toAddress, initiatedBy, collection.Manager)
		if !addressesMatch {
			continue
		}

		//Check if current time is valid
		validTime := false
		for _, timeInterval := range allowedTransfer.TimeIntervals {
			time := sdk.NewUint(uint64(ctx.BlockTime().UnixMilli()))
			if time.GT(timeInterval.Start) && time.LT(timeInterval.End) {
				validTime = true
			}
		}

		if !validTime {
			continue
		}

		//Check if badge ids match
		//TODO: Make this an id_ranges.go function
		matchingBadgeIds := false
		//Start with all badge IDs in the collection and remove as we handle them
		startRanges := SortAndMergeOverlapping(badgeIds)
		for _, badgeIdRange := range allowedTransfer.BadgeIds {
			newStartRanges := []*types.IdRange{}
			for _, idRange := range startRanges {
				removedRanges, _ := RemoveIdsFromIdRange(badgeIdRange, idRange)
				newStartRanges = append(newStartRanges, removedRanges...)
			}
			startRanges = newStartRanges
		}
		matchingBadgeIds = len(startRanges) == 0

		if !matchingBadgeIds {
			continue
		}

		//Lastly, return if the transfer is allowed and if it requires approval
		allowed := allowedTransfer.IsAllowed
		requiresApproval := !allowedTransfer.NoApprovalRequired

		return allowed, requiresApproval
	}

	//If not explicitly allowed
	return false, true
}
