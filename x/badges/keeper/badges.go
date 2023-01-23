package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//Gets badge and throws error if it does not exist. Alternative to GetCollectionFromStore which returns a found bool, not an error.
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

func (k Keeper) ValidateIdRanges(collection types.BadgeCollection, ranges []*types.IdRange) (error) {
	for _, badgeIdRange := range ranges {
		badgeIdRange = NormalizeIdRange(badgeIdRange)

		if badgeIdRange.Start > badgeIdRange.End {
			return ErrInvalidBadgeRange
		}

		if badgeIdRange.End >= collection.NextBadgeId {
			return ErrBadgeNotExists
		}
	}

	return nil
}

//For each (supply, amount) pair, we create (amount) badges with a supply of (supply). Error if IDs overflow.
//We assume that lengths of supplys and amountsToCreate are equal before entering this function. Also amountsToCreate[i] can never be zero.
func (k Keeper) CreateBadges(ctx sdk.Context, collection types.BadgeCollection, supplysAndAmounts []*types.BadgeSupplyAndAmount, transfers []*types.Transfers, claims []*types.Claim, creatorAddress string) (types.BadgeCollection, error) {
	maxSupplys := collection.MaxSupplys //get current
	unmintedSupplys := collection.UnmintedSupplys

	err := *new(error)
	//Update supplys and mint total supply for each to manager.
	//Subasset supplys are stored as []*types.Balance, so we can use the balance update functions
	for _, obj := range supplysAndAmounts {
		amountToCreate := obj.Amount
		supply := obj.Supply
		nextBadgeId := collection.NextBadgeId

		if supply == 0 {
			return types.BadgeCollection{}, ErrSupplyEqualsZero
		}

		maxSupplys = UpdateBalancesForIdRanges(
			[]*types.IdRange{
				{Start: nextBadgeId, End: nextBadgeId + amountToCreate - 1},
			},
			supply,
			maxSupplys,
		)

		unmintedSupplys = UpdateBalancesForIdRanges(
			[]*types.IdRange{
				{Start: nextBadgeId, End: nextBadgeId + amountToCreate - 1},
			},
			supply,
			unmintedSupplys,
		)

		collection.NextBadgeId, err = SafeAdd(collection.NextBadgeId, amountToCreate) //error on ID overflow
		if err != nil {
			return types.BadgeCollection{}, err
		}
	}

	collection.UnmintedSupplys = unmintedSupplys
	collection.MaxSupplys = maxSupplys

	accsToCheck := []uint64{}
	for _, transfer := range transfers {
		accsToCheck = append(accsToCheck, transfer.ToAddresses...)
	}

	rangesToValidate := []*types.IdRange{}
	for _, transfer := range transfers {
		for _, balance := range transfer.Balances {
			rangesToValidate = append(rangesToValidate, balance.BadgeIds...)
		}
	}

	_, _, err = k.UniversalValidate(ctx, UniversalValidationParams{
		Creator: creatorAddress,
		AccountsToCheckRegistration: accsToCheck,
		OnlyCheckAccounts: true,
	})
	if err != nil {
		return types.BadgeCollection{}, err
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
	//Treat the unminted balances as another account and use our transfer function
	unmintedBalances := types.UserBalance{
		Balances: collection.UnmintedSupplys,
	}
	

	err := *new(error)
	for _, transfer := range transfers {
		for _, address := range transfer.ToAddresses {
			recipientBalance := types.UserBalance{}
			hasBalance := k.StoreHasUserBalance(ctx, ConstructBalanceKey(address, collection.CollectionId))
			if hasBalance {
				recipientBalance, _ = k.GetUserBalanceFromStore(ctx, ConstructBalanceKey(address, collection.CollectionId))
			}

			for _, balanceObj := range transfer.Balances {
				unmintedBalances, err = SubtractBalancesForIdRanges(unmintedBalances, balanceObj.BadgeIds, balanceObj.Balance)
				if err != nil {
					return types.BadgeCollection{}, err
				}

				recipientBalance, err = AddBalancesForIdRanges(recipientBalance, balanceObj.BadgeIds, balanceObj.Balance)
				if err != nil {
					return types.BadgeCollection{}, err
				}
			}

			if err := k.SetUserBalanceInStore(ctx, ConstructBalanceKey(address, collection.CollectionId), GetBalanceToInsertToStorage(recipientBalance)); err != nil {
				return types.BadgeCollection{}, err
			}
		}
	}

	collection.UnmintedSupplys = unmintedBalances.Balances

	return collection, nil
}


func (k Keeper) MintViaClaim(ctx sdk.Context, collection types.BadgeCollection, claims []*types.Claim) (types.BadgeCollection, error) {
	//Treat the unminted balances as another account and use our transfer function
	unmintedBalances := types.UserBalance{
		Balances: collection.UnmintedSupplys,
	}
	
	err := *new(error)
	for _, claim := range claims {
		unmintedBalances, err = SubtractBalancesForIdRanges(unmintedBalances, claim.Balance.BadgeIds, claim.Balance.Balance)
		if err != nil {
			return types.BadgeCollection{}, err
		}

		nextClaimId := k.GetNextClaimId(ctx)
		err = k.SetClaimInStore(ctx, nextClaimId, *claim)
		if err != nil {
			return types.BadgeCollection{}, err
		}

		k.IncrementNextClaimId(ctx)
	}

	collection.UnmintedSupplys = unmintedBalances.Balances

	return collection, nil
}