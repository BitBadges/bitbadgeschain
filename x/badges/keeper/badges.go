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
func (k Keeper) CreateBadges(ctx sdk.Context, collection types.BadgeCollection, supplysAndAmounts []*types.Balance, transfers []*types.Transfer, creatorAddress string) (types.BadgeCollection, error) {
	totalSupplys := collection.TotalSupplys //get current
	unmintedSupplys := collection.UnmintedSupplys

	err := *new(error)

	allBadgeIds := []*types.IdRange{}
	for _, balanceObj := range supplysAndAmounts {
		allBadgeIds = append(allBadgeIds, balanceObj.BadgeIds...)
	}

	err = CheckActionWithBadgeIdsPermission(ctx, allBadgeIds, collection.Permissions.CanCreateMoreBadges)
	if err != nil {
		return types.BadgeCollection{}, err
	}
	nextBadgeId := collection.NextBadgeId

	//Update supplys and mint total supply for each to manager.
	//Subasset supplys are stored as []*types.Balance, so we can use the balance update functions
	for _, balance := range supplysAndAmounts {
		if balance.Amount.IsZero() {
			return types.BadgeCollection{}, ErrSupplyEqualsZero
		}

		for _, badgeIdRange := range balance.BadgeIds {
			if badgeIdRange.End.GTE(nextBadgeId) {
				nextBadgeId = badgeIdRange.End
			}
		}

		totalSupplys, err = AddBalancesForIdRanges(totalSupplys, balance.BadgeIds,	balance.Amount)
		if err != nil {
			return types.BadgeCollection{}, err
		}

		unmintedSupplys, err = AddBalancesForIdRanges(unmintedSupplys, balance.BadgeIds, balance.Amount)
		if err != nil {
			return types.BadgeCollection{}, err
		}
	}

	collection.UnmintedSupplys = unmintedSupplys
	collection.TotalSupplys = totalSupplys
	collection.NextBadgeId = nextBadgeId

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

	if collection.BalancesType.LTE(sdk.NewUint(0)) {
		collection, err = k.MintViaTransfer(ctx, collection, transfers)
		if err != nil {
			return types.BadgeCollection{}, err
		}
	} else {
		//TODO: add approvedTransfers check here?
		if len(transfers) > 0 {
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
				allowed, _ := IsTransferAllowed(ctx, balanceObj.BadgeIds, collection, "Mint", address, "Manager")
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