package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Gets badge and throws error if it does not exist. Alternative to GetCollectionFromStore which returns a found bool, not an error.
func (k Keeper) GetCollectionE(ctx sdk.Context, badgeId sdk.Uint) (types.BadgeCollection, error) {
	badge, found := k.GetCollectionFromStore(ctx, badgeId)
	if !found {
		return types.BadgeCollection{}, ErrCollectionNotExists
	}

	return badge, nil
}

// Gets the badge details from the store if it exists. Throws error if badge ranges are invalid or a badge ID does not yet exist.
func (k Keeper) GetCollectionAndAssertBadgeIdsAreValid(ctx sdk.Context, collectionId sdk.Uint, badgeIdRanges []*types.IdRange) (types.BadgeCollection, error) {
	badge, err := k.GetCollectionE(ctx, collectionId)
	if err != nil {
		return badge, err
	}

	err = k.ValidateIdRanges(badge, badgeIdRanges, true)
	if err != nil {
		return types.BadgeCollection{}, err
	}

	return badge, nil
}

func (k Keeper) ValidateIdRanges(collection types.BadgeCollection, ranges []*types.IdRange, checkIfLessThanNextBadgeId bool) error {
	for _, badgeIdRange := range ranges {
		if badgeIdRange.Start.GT(badgeIdRange.End) {
			return ErrInvalidBadgeRange
		}

		if checkIfLessThanNextBadgeId && badgeIdRange.End.GTE(collection.NextBadgeId) {
			return ErrBadgeNotExists
		}
	}

	return nil
}

// For each (supply, amount) pair, we create (amount) badges with a supply of (supply). Error if IDs overflow.
// We assume that lengths of supplys and amountsToCreate are equal before entering this function. Also amountsToCreate[i] can never be zero.
func (k Keeper) CreateBadges(ctx sdk.Context, collection types.BadgeCollection, supplysAndAmounts []*types.BadgeSupply, transfers []*types.Transfer, claims []*types.Claim, creatorAddress string) (types.BadgeCollection, error) {
	maxSupplys := collection.MaxSupplys //get current
	unmintedSupplys := collection.UnmintedSupplys

	err := *new(error)

	//Update supplys and mint total supply for each to manager.
	//Subasset supplys are stored as []*types.Balance, so we can use the balance update functions
	for _, obj := range supplysAndAmounts {
		amount := obj.Amount
		supply := obj.Supply
		nextBadgeId := collection.NextBadgeId

		if supply.IsZero() || amount.IsZero() {
			return types.BadgeCollection{}, ErrSupplyEqualsZero
		}

		collection.NextBadgeId, err = SafeAdd(collection.NextBadgeId, amount) //error on ID overflow
		if err != nil {
			return types.BadgeCollection{}, err
		}

		maxSupplys, err = UpdateBalancesForIdRanges(
			[]*types.IdRange{
				{Start: nextBadgeId, End: nextBadgeId.Add(amount).SubUint64(1)},
			},
			supply,
			maxSupplys,
		)
		if err != nil {
			return types.BadgeCollection{}, err
		}

		unmintedSupplys, err = UpdateBalancesForIdRanges(
			[]*types.IdRange{
				{Start: nextBadgeId, End: nextBadgeId.Add(amount).SubUint64(1)},
			},
			supply,
			unmintedSupplys,
		)
		if err != nil {
			return types.BadgeCollection{}, err
		}
	}

	collection.UnmintedSupplys = unmintedSupplys
	collection.MaxSupplys = maxSupplys

	rangesToValidate := []*types.IdRange{}
	for _, transfer := range transfers {
		for _, balance := range transfer.Balances {
			rangesToValidate = append(rangesToValidate, balance.BadgeIds...)
		}
	}

	err = k.ValidateIdRanges(collection, rangesToValidate, false)
	if err != nil {
		return types.BadgeCollection{}, err
	}

	if !collection.IsOffChainBalances {
		collection, err = k.MintViaTransfer(ctx, collection, transfers)
		if err != nil {
			return types.BadgeCollection{}, err
		}

		collection, err = k.MintViaClaim(ctx, collection, claims)
		if err != nil {
			return types.BadgeCollection{}, err
		}
	} else {
		if len(transfers) > 0 || len(claims) > 0 {
			return types.BadgeCollection{}, ErrOffChainBalances
		}
	}

	return collection, nil
}

func (k Keeper) MintViaTransfer(ctx sdk.Context, collection types.BadgeCollection, transfers []*types.Transfer) (types.BadgeCollection, error) {
	//Treat the unminted balances as another account for compatibility with our transfer function
	unmintedBalances := types.UserBalanceStore{
		Balances: collection.UnmintedSupplys,
	}

	err := *new(error)
	for _, transfer := range transfers {
		for _, address := range transfer.ToAddresses {
			recipientBalance := types.UserBalanceStore{}
			hasBalance := k.StoreHasUserBalance(ctx, ConstructBalanceKey(address, collection.CollectionId))
			if hasBalance {
				recipientBalance, _ = k.GetUserBalanceFromStore(ctx, ConstructBalanceKey(address, collection.CollectionId))
			}

			for _, balanceObj := range transfer.Balances {
				//Check if minting via transfer is allowed to this address
				//ignore requiresApproval because we are minting
				allowed, _ := IsTransferAllowed(ctx, balanceObj.BadgeIds, collection, "Mint", address, "Mint")
				if !allowed {
					return types.BadgeCollection{}, ErrMintNotAllowed
				}

				unmintedBalances.Balances, err = SubtractBalancesForIdRanges(unmintedBalances.Balances, balanceObj.BadgeIds, balanceObj.Amount)
				if err != nil {
					return types.BadgeCollection{}, err
				}

				recipientBalance.Balances, err = AddBalancesForIdRanges(recipientBalance.Balances, balanceObj.BadgeIds, balanceObj.Amount)
				if err != nil {
					return types.BadgeCollection{}, err
				}
			}

			if err := k.SetUserBalanceInStore(ctx, ConstructBalanceKey(address, collection.CollectionId), recipientBalance); err != nil {
				return types.BadgeCollection{}, err
			}
		}
	}

	collection.UnmintedSupplys = unmintedBalances.Balances

	return collection, nil
}

// type BadgeIdsPlusTimeIntervals struct  {
// 	BadgeIds []*types.IdRange
// 	TimeIntervals []*types.IdRange
// }

func (k Keeper) MintViaClaim(ctx sdk.Context, collection types.BadgeCollection, claims []*types.Claim) (types.BadgeCollection, error) {
	//It is all dependent on the to addresses as well.
	//Ex: a transfer may map to one allowed transfer for one user but a completely different allowed transfer for another user
	//Thus, we can not iterate through the allowed transfers and check if it is allowed for all users for the entire duration of the claim.

	// approvedTransfers := collection.ApprovedTransfers

	// for _, claim := range claims {
	// 	claimTimeIntervals := claim.TimeIntervals
	// 	nonHandledMap := []BadgeIdsPlusTimeIntervals{}
	// 	for _, balance := range claim.Balances {
	// 		nonHandledMap = append(nonHandledMap, BadgeIdsPlusTimeIntervals{
	// 			BadgeIds: balance.BadgeIds,
	// 			TimeIntervals: claimTimeIntervals,
	// 		})
	// 	}

	// 	//We need to ensure that all badge IDs will be transferable for the entire duration of the claim (i.e. explicitly allowed)
	// 	for _, allowedTransfer := range approvedTransfers {
	// 		// If it is not frozen, it is updatable and can potentially be changed.
	// 		// We need to ensure it is claimable for duration of claim.
	// 		// All claim transfers will be from = "Mint" and initiatedBy = "Manager". Thus, we do not want to allow from == initiatedBy or to == initiatedBy
	// 		if !allowedTransfer.IsFrozen || allowedTransfer.RequireFromEqualsInitiatedBy || allowedTransfer.RequireToEqualsInitiatedBy {
	// 			continue
	// 		}

	// 		//TODO: We should also probably handle a check of the "To" addresses

	// 		fromFound := CheckAddresses(allowedTransfer.From, "Mint", "Manager")
	// 		initiatedByFound := CheckAddresses(allowedTransfer.InitiatedBy, "Mint", "Manager")
	// 		if !fromFound || !initiatedByFound {
	// 			continue
	// 		}

	// 		newNonHandledMap := []BadgeIdsPlusTimeIntervals{}
	// 		for _, nonHandled := range nonHandledMap {
	// 			nonOverlappingBadgeIds := append([]*types.IdRange{}, nonHandled.BadgeIds...)
	// 			nonOverlappingTimeIntervals := append([]*types.IdRange{}, nonHandled.TimeIntervals...)

	// 			//Remove overlapping badge IDs from each entry
	// 			for _, badgeIdRange := range allowedTransfer.BadgeIds {
	// 				newNonOverlappingBadgeIds := []*types.IdRange{}
	// 				for _, nonOverlappingBadgeId := range nonOverlappingBadgeIds {
	// 					removedRanges, _ := RemoveIdsFromIdRange(badgeIdRange, nonOverlappingBadgeId)
	// 					newNonOverlappingBadgeIds = append(newNonOverlappingBadgeIds, removedRanges...)
	// 				}
	// 				nonOverlappingBadgeIds = newNonOverlappingBadgeIds
	// 			}

	// 			//Remove overlapping time intervals from each entry
	// 			for _, timeInterval := range allowedTransfer.TimeIntervals {
	// 				newNonOverlappingTimeIntervals := []*types.IdRange{}
	// 				for _, nonOverlappingTimeInterval := range nonOverlappingTimeIntervals {
	// 					removedRanges, _ := RemoveIdsFromIdRange(timeInterval, nonOverlappingTimeInterval)
	// 					newNonOverlappingTimeIntervals = append(newNonOverlappingTimeIntervals, removedRanges...)
	// 				}
	// 				nonOverlappingTimeIntervals = newNonOverlappingTimeIntervals
	// 			}

	// 			overlappingBadgeIds := append([]*types.IdRange{}, nonHandled.BadgeIds...)
	// 			for _, nonOverlappingBadgeId := range nonOverlappingBadgeIds {
	// 				newOverlappingBadgeIds := []*types.IdRange{}
	// 				for _, overlappingBadgeId := range overlappingBadgeIds {
	// 					removedRanges, _ := RemoveIdsFromIdRange(nonOverlappingBadgeId, overlappingBadgeId)
	// 					newOverlappingBadgeIds = append(newOverlappingBadgeIds, removedRanges...)
	// 				}
	// 				overlappingBadgeIds = newOverlappingBadgeIds
	// 			}

	// 			overlappingTimeIntervals := append([]*types.IdRange{}, nonHandled.TimeIntervals...)
	// 			for _, nonOverlappingTimeInterval := range nonOverlappingTimeIntervals {
	// 				newOverlappingTimeIntervals := []*types.IdRange{}
	// 				for _, overlappingTimeInterval := range overlappingTimeIntervals {
	// 					removedRanges, _ := RemoveIdsFromIdRange(nonOverlappingTimeInterval, overlappingTimeInterval)
	// 					newOverlappingTimeIntervals = append(newOverlappingTimeIntervals, removedRanges...)
	// 				}
	// 				overlappingTimeIntervals = newOverlappingTimeIntervals
	// 			}

	// 			//If there are any overlaps on a disallowed transfer, then we can't mint
	// 			if len(overlappingBadgeIds) > 0 && len(overlappingTimeIntervals) > 0 && !allowedTransfer.IsAllowed {
	// 				return types.BadgeCollection{}, ErrMintNotAllowed
	// 			}

	// 			if len(overlappingBadgeIds) > 0  && len(nonOverlappingTimeIntervals) > 0 {
	// 				newNonHandledMap = append(newNonHandledMap, BadgeIdsPlusTimeIntervals{
	// 					BadgeIds: overlappingBadgeIds,
	// 					TimeIntervals: nonOverlappingTimeIntervals,
	// 				})
	// 			}

	// 			if len(nonOverlappingBadgeIds) > 0 {
	// 				newNonHandledMap = append(newNonHandledMap, BadgeIdsPlusTimeIntervals{
	// 					BadgeIds: nonOverlappingBadgeIds,
	// 					TimeIntervals: nonHandled.TimeIntervals,
	// 				})
	// 			}
	// 		}

	// 		nonHandledMap = newNonHandledMap
	// 	}

	// 	//At the end, we should expect everything to be handled
	// 	if len(nonHandledMap) > 0 {
	// 		return types.BadgeCollection{}, ErrMintNotAllowed
	// 	}
	// }

	//Treat the unminted balances as another account and use our transfer function
	unmintedBalances := types.UserBalanceStore{
		Balances: collection.UnmintedSupplys,
	}

	currClaimId := collection.NextClaimId
	err := *new(error)
	for _, claim := range claims {
		for _, balance := range claim.Balances {
			unmintedBalances.Balances, err = SubtractBalancesForIdRanges(unmintedBalances.Balances, balance.BadgeIds, balance.Amount)
			if err != nil {
				return types.BadgeCollection{}, err
			}
		}

		err = k.SetClaimInStore(ctx, collection.CollectionId, currClaimId, *claim)
		if err != nil {
			return types.BadgeCollection{}, err
		}

		currClaimId, err = SafeAdd(currClaimId, sdk.NewUint(1))
		if err != nil {
			return types.BadgeCollection{}, err
		}
	}

	collection.NextClaimId = currClaimId
	collection.UnmintedSupplys = unmintedBalances.Balances

	return collection, nil
}
