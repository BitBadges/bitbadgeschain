package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Gets badge and throws error if it does not exist. Alternative to GetCollectionFromStore which returns a found bool, not an error.
func (k Keeper) GetCollectionE(ctx sdk.Context, badgeId uint64) (types.BadgeCollection, error) {
	badge, found := k.GetCollectionFromStore(ctx, badgeId)
	if !found {
		return types.BadgeCollection{}, ErrCollectionNotExists
	}

	return badge, nil
}

// Gets the badge details from the store if it exists. Throws error if badge ranges are invalid or a badge ID does not yet exist.
func (k Keeper) GetCollectionAndAssertBadgeIdsAreValid(ctx sdk.Context, collectionId uint64, badgeIdRanges []*types.IdRange) (types.BadgeCollection, error) {
	badge, err := k.GetCollectionE(ctx, collectionId)
	if err != nil {
		return badge, err
	}

	err = k.ValidateIdRanges(badge, badgeIdRanges)
	if err != nil {
		return types.BadgeCollection{}, err
	}

	return badge, nil
}

func (k Keeper) ValidateIdRanges(collection types.BadgeCollection, ranges []*types.IdRange) error {
	for _, badgeIdRange := range ranges {
		if badgeIdRange.Start > badgeIdRange.End {
			return ErrInvalidBadgeRange
		}

		if badgeIdRange.End >= collection.NextBadgeId {
			return ErrBadgeNotExists
		}
	}

	return nil
}

// For each (supply, amount) pair, we create (amount) badges with a supply of (supply). Error if IDs overflow.
// We assume that lengths of supplys and amountsToCreate are equal before entering this function. Also amountsToCreate[i] can never be zero.
func (k Keeper) CreateBadges(ctx sdk.Context, collection types.BadgeCollection, supplysAndAmounts []*types.BadgeSupplyAndAmount, transfers []*types.Transfers, claims []*types.Claim, creatorAddress string) (types.BadgeCollection, error) {
	maxSupplys := collection.MaxSupplys //get current
	unmintedSupplys := collection.UnmintedSupplys

	err := *new(error)

	//Update supplys and mint total supply for each to manager.
	//Subasset supplys are stored as []*types.Balance, so we can use the balance update functions
	for _, obj := range supplysAndAmounts {
		amount := obj.Amount
		supply := obj.Supply
		nextBadgeId := collection.NextBadgeId

		if supply == 0 || amount == 0 {
			return types.BadgeCollection{}, ErrSupplyEqualsZero
		}

		collection.NextBadgeId, err = SafeAdd(collection.NextBadgeId, amount) //error on ID overflow
		if err != nil {
			return types.BadgeCollection{}, err
		}

		maxSupplys, err = UpdateBalancesForIdRanges(
			[]*types.IdRange{
				{Start: nextBadgeId, End: nextBadgeId + amount - 1},
			},
			supply,
			maxSupplys,
		)
		if err != nil {
			return types.BadgeCollection{}, err
		}

		unmintedSupplys, err = UpdateBalancesForIdRanges(
			[]*types.IdRange{
				{Start: nextBadgeId, End: nextBadgeId + amount - 1},
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

	err = k.ValidateIdRanges(collection, rangesToValidate)
	if err != nil {
		return types.BadgeCollection{}, err
	}

	collection, err = k.MintViaTransfer(ctx, collection, transfers)
	if err != nil {
		return types.BadgeCollection{}, err
	}

	collection, err = k.MintViaClaim(ctx, collection, claims)
	if err != nil {
		return types.BadgeCollection{}, err
	}

	return collection, nil
}

func (k Keeper) MintViaTransfer(ctx sdk.Context, collection types.BadgeCollection, transfers []*types.Transfers) (types.BadgeCollection, error) {
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
				unmintedBalances, err = SubtractBalancesForIdRanges(unmintedBalances, balanceObj.BadgeIds, balanceObj.Amount)
				if err != nil {
					return types.BadgeCollection{}, err
				}

				recipientBalance, err = AddBalancesForIdRanges(recipientBalance, balanceObj.BadgeIds, balanceObj.Amount)
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

func (k Keeper) MintViaClaim(ctx sdk.Context, collection types.BadgeCollection, claims []*types.Claim) (types.BadgeCollection, error) {
	//Treat the unminted balances as another account and use our transfer function
	unmintedBalances := types.UserBalanceStore{
		Balances: collection.UnmintedSupplys,
	}

	currClaimId := collection.NextClaimId
	err := *new(error)
	for _, claim := range claims {
		for _, balance := range claim.Balances {
			unmintedBalances, err = SubtractBalancesForIdRanges(unmintedBalances, balance.BadgeIds, balance.Amount)
			if err != nil {
				return types.BadgeCollection{}, err
			}
		}

		err = k.SetClaimInStore(ctx, collection.CollectionId, currClaimId, *claim)
		if err != nil {
			return types.BadgeCollection{}, err
		}

		currClaimId, err = SafeAdd(currClaimId, 1)
		if err != nil {
			return types.BadgeCollection{}, err
		}
	}

	collection.NextClaimId = currClaimId
	collection.UnmintedSupplys = unmintedBalances.Balances

	return collection, nil
}
