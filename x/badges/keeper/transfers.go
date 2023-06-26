package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"math"
)

func (k Keeper) HandleTransfers(ctx sdk.Context, collection types.BadgeCollection, transfers []*types.Transfer, initiatedBy string) (error) {
	//If "Mint" or "Manager" 
	unmintedBalances := types.UserBalanceStore{
		Balances: collection.UnmintedSupplys,
	}
	err := *new(error)

	if initiatedBy == "Manager" {
		initiatedBy = GetCurrentManager(ctx, collection)
	}

	for _, transfer := range transfers {
		fromBalanceKey := ConstructBalanceKey(transfer.From, collection.CollectionId)
		fromUserBalance := unmintedBalances
		found := true
		if transfer.From == "Manager" {
			transfer.From = GetCurrentManager(ctx, collection)
		} 
		
		if transfer.From != "Mint" {
			fromUserBalance, found = k.GetUserBalanceFromStore(ctx, fromBalanceKey)
			if !found {
				return ErrUserBalanceNotExists
			}
		}
		
		for _, to := range transfer.ToAddresses {
			if to == "Manager" {
				to = GetCurrentManager(ctx, collection)
			}

			toBalanceKey := ConstructBalanceKey(to, collection.CollectionId)
			toUserBalance, found := k.GetUserBalanceFromStore(ctx, toBalanceKey)
			if !found {
				toUserBalance = types.UserBalanceStore{
					Balances : []*types.Balance{},
					ApprovedTransfersTimeline: []*types.UserApprovedTransferTimeline{},
					NextTransferTrackerId: sdk.NewUint(1),
					Permissions: &types.UserPermissions{
						CanUpdateApprovedTransfers: []*types.UserApprovedTransferPermission{},
					},
				}
			}

			for _, balance := range transfer.Balances {
				amount := balance.Amount
				fromUserBalance, toUserBalance, err = HandleTransfer(ctx, collection, balance.BadgeIds, balance.Times, fromUserBalance, toUserBalance, amount, transfer.From, to, initiatedBy)
				if err != nil {
					return err
				}
			}
			
			//TODO: solutions

			if err := k.SetUserBalanceInStore(ctx, toBalanceKey, toUserBalance); err != nil {
				return err
			}
		}

		if err := k.SetUserBalanceInStore(ctx, fromBalanceKey, fromUserBalance); err != nil {
			return err
		}
	}

	return nil
}

// Forceful transfers will transfer the balances and deduct from approvals directly without adding it to pending.
func HandleTransfer(ctx sdk.Context, collection types.BadgeCollection, badgeIds []*types.IdRange, times []*types.IdRange, fromUserBalance types.UserBalanceStore, toUserBalance types.UserBalanceStore, amount sdk.Uint, from string, to string, initiatedBy string) (types.UserBalanceStore, types.UserBalanceStore, error) {
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
		fromUserBalance, err = DeductApprovalsIfNeeded(ctx, &fromUserBalance, badgeIds, times, from, to, initiatedBy, amount)
		if err != nil {
			return types.UserBalanceStore{}, types.UserBalanceStore{}, err
		}
	}

	fromUserBalance.Balances, err = SubtractBalancesForIdRanges(fromUserBalance.Balances, badgeIds, times, amount)
	if err != nil {
		return types.UserBalanceStore{}, types.UserBalanceStore{}, err
	}

	toUserBalance.Balances, err = AddBalancesForIdRanges(toUserBalance.Balances, badgeIds, times, amount)
	if err != nil {
		return types.UserBalanceStore{}, types.UserBalanceStore{}, err
	}

	return fromUserBalance, toUserBalance, nil
}

// Deduct approvals from requester if requester != from
func DeductApprovalsIfNeeded(ctx sdk.Context, userBalance *types.UserBalanceStore, badgeIds []*types.IdRange, times []*types.IdRange, from string, to string, requester string, amount sdk.Uint) (types.UserBalanceStore, error) {
	//check all approvals, default disallow at the end unless from === initiator	
	newUserBalance := userBalance
	currApprovedTransfers := GetCurrentUserApprovedTransfers(ctx, userBalance)
	newCurrApprovedTransfers := []*types.UserApprovedTransfer{}
	for _, approvedTransfer := range currApprovedTransfers {
		for _, allowedCombination := range approvedTransfer.AllowedCombinations {
			badgeIds := approvedTransfer.BadgeIds
			if allowedCombination.InvertBadgeIds {
				badgeIds = types.InvertIdRanges(badgeIds, sdk.NewUint(math.MaxUint64))
			}

			times := approvedTransfer.TransferTimes
			if allowedCombination.InvertTransferTimes {
				times = types.InvertIdRanges(times, sdk.NewUint(uint64(ctx.BlockTime().UnixMilli())))
			}

			toMappingId := approvedTransfer.ToMappingId
			if allowedCombination.InvertTo {
				toMappingId = "!" + toMappingId
			}

			initiatedByMappingId := approvedTransfer.InitiatedByMappingId
			if allowedCombination.InvertInitiatedBy {
				initiatedByMappingId = "!" + initiatedByMappingId
			}

			newCurrApprovedTransfers = append(newCurrApprovedTransfers, &types.UserApprovedTransfer{
				ToMappingId: toMappingId,
				InitiatedByMappingId: initiatedByMappingId,
				TransferTimes: times,
				BadgeIds: badgeIds,
				AllowedCombinations: []*types.IsTransferAllowed{
					{
						IsAllowed: allowedCombination.IsAllowed,
					},
				},
				AmountRestrictions: approvedTransfer.AmountRestrictions,
				TransferTrackerId: approvedTransfer.TransferTrackerId,
				RequireToEqualsInitiatedBy: approvedTransfer.RequireToEqualsInitiatedBy,
				RequireToDoesNotEqualInitiatedBy: approvedTransfer.RequireToDoesNotEqualInitiatedBy,
			})
		}
	}


	newCurrApprovedTransfers = append(newCurrApprovedTransfers, &types.UserApprovedTransfer{
		//TODO:
		ToMappingId: "", //everyone
		InitiatedByMappingId: "", //only user address
		TransferTimes: []*types.IdRange{
			{
				Start: sdk.NewUint(0),
				End: sdk.NewUint(uint64(ctx.BlockTime().UnixMilli())),
			},
		},
		BadgeIds: []*types.IdRange{
			{
				Start: sdk.NewUint(1),
				End: sdk.NewUint(math.MaxUint64),
			},
		},
		AllowedCombinations: []*types.IsTransferAllowed{
			{
				IsAllowed: true,
			},
		},
		AmountRestrictions: []*types.AmountRestrictions{{}},
		TransferTrackerId: sdk.NewUint(0), //TODO: think about this
	})
	
	castedApprovedTransfers := CastUserApprovedTransferToUniversalPermission(newCurrApprovedTransfers)
	firstMatches := types.GetFirstMatchOnly(castedApprovedTransfers) //but could be duplicate mapping IDs so we need to be careful here

	unhandledBadgeIds := badgeIds
	for _, match := range firstMatches {
		//TODO: check if the match is valid
		//check to and initiatedBy mapping IDs
		_, timeFound := types.SearchIdRangesForId(sdk.NewUint(uint64(ctx.BlockTime().UnixMilli())), []*types.IdRange{match.TransferTime})
		removedBadges := []*types.IdRange{}
		
		if timeFound { //TODO: add mapping id checks
			unhandledBadgeIds, removedBadges = types.RemoveIdRangeFromIdRange([]*types.IdRange{match.BadgeId}, unhandledBadgeIds)
			if len(removedBadges) > 0 {
				//We have a valid match for at least some badges, procees to check restrictions
				approvedTransfer := match.ArbitraryValue.(*types.UserApprovedTransfer)
				isAllowed := approvedTransfer.AllowedCombinations[0].IsAllowed
				if !isAllowed {
					return types.UserBalanceStore{}, ErrInadequateApprovals
				}
				
				transferTrackerId := approvedTransfer.TransferTrackerId
				amountRestrictions := approvedTransfer.AmountRestrictions
				requireToEqualsInitiatedBy := approvedTransfer.RequireToEqualsInitiatedBy
				requireToDoesNotEqualInitiatedBy := approvedTransfer.RequireToDoesNotEqualInitiatedBy

				if requireToEqualsInitiatedBy && to != requester {
					return types.UserBalanceStore{}, ErrInadequateApprovals
				}

				if requireToDoesNotEqualInitiatedBy && to == requester {
					return types.UserBalanceStore{}, ErrInadequateApprovals
				}

				//TODO: Check amount restrictions from transfertrackerId

				
			}
		}
	}

	return newUserBalance, ErrInadequateApprovals
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
