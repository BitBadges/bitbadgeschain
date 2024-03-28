package keeper

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Create badges and update the unminted / total supplys for the collection
func (k Keeper) CreateBadges(ctx sdk.Context, collection *types.BadgeCollection, badgesToCreate []*types.Balance) (*types.BadgeCollection, error) {
	if !IsStandardBalances(collection) {
		//For readability, we do not allow transfers to happen on-chain, if not defined in the collection
		if len(collection.CollectionApprovals) > 0 {
			return &types.BadgeCollection{}, ErrWrongBalancesType
		}
	}


	//Check if we are allowed to create these badges
	err := *new(error)
	detailsToCheck := []*types.UniversalPermissionDetails{}
	for _, balanceObj := range badgesToCreate {
		for _, badgeIdRange := range balanceObj.BadgeIds {
			for _, time := range balanceObj.OwnershipTimes {
				detailsToCheck = append(detailsToCheck, &types.UniversalPermissionDetails{
					BadgeId:       badgeIdRange,
					OwnershipTime: time,
				})
			}
		}
	}

	err = k.CheckIfBalancesActionPermissionPermits(ctx, detailsToCheck, collection.CollectionPermissions.CanCreateMoreBadges, "can create more badges")
	if err != nil {
		return &types.BadgeCollection{}, err
	}

	unmintedSupplys, found := k.GetUserBalanceFromStore(ctx, ConstructBalanceKey("Mint", collection.CollectionId))
	if !found {
		unmintedSupplys = &types.UserBalanceStore{} //permissions and timelines do not matter because these are special addresses
	}

	totalSupplys, found := k.GetUserBalanceFromStore(ctx, ConstructBalanceKey("Total", collection.CollectionId))
	if !found {
		totalSupplys = &types.UserBalanceStore{} //permissions and timelines do not matter because these are special addresses
	}

	allBadgeIds := []*types.UintRange{}
	for _, balance := range totalSupplys.Balances {
		allBadgeIds = append(allBadgeIds, balance.BadgeIds...)
	}

	for _, balance := range badgesToCreate {
		allBadgeIds = append(allBadgeIds, balance.BadgeIds...)
	}

	allBadgeIds, err = types.SortUintRangesAndMerge(allBadgeIds, true)
	if err != nil {
		return &types.BadgeCollection{}, err
	}

	if len(allBadgeIds) > 1 || (len(allBadgeIds) == 1 && !allBadgeIds[0].Start.Equal(sdkmath.NewUint(1))) {
		return &types.BadgeCollection{}, sdkerrors.Wrapf(types.ErrNotSupported, "Badge Ids must be sequential starting from 1")
	}

	//Create the badges and add newly created balances to unminted supplys
	for _, balance := range badgesToCreate {
		if balance.Amount.IsZero() {
			return &types.BadgeCollection{}, ErrSupplyEqualsZero
		}

		totalSupplys.Balances, err = types.AddBalance(ctx, totalSupplys.Balances, balance)
		if err != nil {
			return &types.BadgeCollection{}, err
		}

		unmintedSupplys.Balances, err = types.AddBalance(ctx, unmintedSupplys.Balances, balance)
		if err != nil {
			return &types.BadgeCollection{}, err
		}
	}

	err = k.SetUserBalanceInStore(ctx, ConstructBalanceKey("Mint", collection.CollectionId), unmintedSupplys)
	if err != nil {
		return &types.BadgeCollection{}, err
	}

	err = k.SetUserBalanceInStore(ctx, ConstructBalanceKey("Total", collection.CollectionId), totalSupplys)
	if err != nil {
		return &types.BadgeCollection{}, err
	}

	return collection, nil
}
