package keeper

import (
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//Create badges and update the unminted / total supplys for the collection
func (k Keeper) CreateBadges(ctx sdk.Context, collection *types.BadgeCollection, badgesToCreate []*types.Balance, transfers []*types.Transfer) (*types.BadgeCollection, error) {
	if IsInheritedBalances(collection) {
		return &types.BadgeCollection{}, ErrInheritedBalances
	}
	
	totalSupplys := collection.TotalSupplys 
	unmintedSupplys := collection.UnmintedSupplys

	err := *new(error)
	detailsToCheck := []*types.UniversalPermissionDetails{}
	for _, balanceObj := range badgesToCreate {
		for _, badgeIdRange := range balanceObj.BadgeIds {
			for _, time := range balanceObj.Times {
				detailsToCheck = append(detailsToCheck, &types.UniversalPermissionDetails{
					BadgeId: badgeIdRange,
					TransferTime: time,
				})
			}
		}
	}

	err = k.CheckActionWithBadgeIdsAndTimesPermission(ctx, detailsToCheck, collection.Permissions.CanCreateMoreBadges)
	if err != nil {
		return &types.BadgeCollection{}, err
	}
	nextBadgeId := collection.NextBadgeId

	//Update supplys and mint total supply for each to manager.
	//Subasset supplys are stored as []*types.Balance, so we can use the balance update functions
	for _, balance := range badgesToCreate {
		if balance.Amount.IsZero() {
			return &types.BadgeCollection{}, ErrSupplyEqualsZero
		}

		for _, badgeIdRange := range balance.BadgeIds {
			if badgeIdRange.End.GTE(nextBadgeId) {
				nextBadgeId = badgeIdRange.End.Add(sdkmath.NewUint(1))
			}
		}

		totalSupplys, err = types.AddBalance(totalSupplys, balance)
		if err != nil {
			return &types.BadgeCollection{}, err
		}

		unmintedSupplys, err = types.AddBalance(unmintedSupplys, balance)
		if err != nil {
			return &types.BadgeCollection{}, err
		}
	}

	collection.UnmintedSupplys = unmintedSupplys
	collection.TotalSupplys = totalSupplys
	collection.NextBadgeId = nextBadgeId

	if IsOnChainBalances(collection) {
		err = k.HandleTransfers(ctx, collection, transfers, "Manager")
		if err != nil {
			return &types.BadgeCollection{}, err
		}
	} else {
		if len(transfers) > 0 || len(collection.ApprovedTransfersTimeline) > 0 {
			return &types.BadgeCollection{}, ErrWrongBalancesType
		}
	}

	return collection, nil
}