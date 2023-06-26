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

	detailsToCheck := []*types.UniversalPermissionDetails{}
	for _, balanceObj := range supplysAndAmounts {
		for _, badgeIdRange := range balanceObj.BadgeIds {
			detailsToCheck = append(detailsToCheck, &types.UniversalPermissionDetails{
				BadgeId: badgeIdRange,
			})
		}
	}

	err = CheckActionWithBadgeIdsPermission(ctx, detailsToCheck, collection.Permissions.CanCreateMoreBadges)
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

		totalSupplys, err = AddBalancesForIdRanges(totalSupplys, balance.BadgeIds, balance.Times, balance.Amount)
		if err != nil {
			return types.BadgeCollection{}, err
		}

		unmintedSupplys, err = AddBalancesForIdRanges(unmintedSupplys, balance.BadgeIds, balance.Times, balance.Amount)
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

	if collection.BalancesType != sdk.NewUint(0) {
		err = k.HandleTransfers(ctx, collection, transfers, "Manager")
		if err != nil {
			return types.BadgeCollection{}, err
		}
	} else {
		if len(transfers) > 0 || len(collection.ApprovedTransfersTimeline) > 0 {
			return types.BadgeCollection{}, ErrOffChainBalances
		}
	}

	return collection, nil
}