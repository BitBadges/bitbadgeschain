package keeper

import (
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//Create badges and update the unminted / total supplys for the collection
func (k Keeper) CreateBadges(ctx sdk.Context, collection *types.BadgeCollection, badgesToCreate []*types.Balance, transfers []*types.Transfer) (*types.BadgeCollection, error) {
	if IsInheritedBalances(collection) {
		return &types.BadgeCollection{}, ErrWrongBalancesType
	}

	//Check if we are allowed to create these badges
	err := *new(error)
	detailsToCheck := []*types.UniversalPermissionDetails{}
	for _, balanceObj := range badgesToCreate {
		for _, badgeIdRange := range balanceObj.BadgeIds {
			for _, time := range balanceObj.OwnershipTimes {
				detailsToCheck = append(detailsToCheck, &types.UniversalPermissionDetails{
					BadgeId: badgeIdRange,
					TransferTime: time,
				})
			}
		}
	}

	err = k.CheckBalancesActionPermission(ctx, detailsToCheck, collection.Permissions.CanCreateMoreBadges)
	if err != nil {
		return &types.BadgeCollection{}, err
	}


	//Create the badges and add newly created balances to unminted supplys
	for _, balance := range badgesToCreate {
		if balance.Amount.IsZero() {
			return &types.BadgeCollection{}, ErrSupplyEqualsZero
		}

		//Update nextBadgeId to be max badgeId + 1
		for _, badgeIdRange := range balance.BadgeIds {
			if badgeIdRange.End.GTE(collection.NextBadgeId) {
				collection.NextBadgeId = badgeIdRange.End.Add(sdkmath.NewUint(1))
			}
		}

		collection.TotalSupplys, err = types.AddBalance(collection.TotalSupplys, balance)
		if err != nil {
			return &types.BadgeCollection{}, err
		}

		collection.UnmintedSupplys, err = types.AddBalance(collection.UnmintedSupplys, balance)
		if err != nil {
			return &types.BadgeCollection{}, err
		}
	}

	
	if IsStandardBalances(collection) {
		//Handle any transfers defined in the Msg
		err = k.HandleTransfers(ctx, collection, transfers, "Manager")
		if err != nil {
			return &types.BadgeCollection{}, err
		}
	} else {
		//For readability, we do not allow transfers to happen on-chain, if not defined in the collection
		if len(transfers) > 0 || len(collection.CollectionApprovedTransfersTimeline) > 0 {
			return &types.BadgeCollection{}, ErrWrongBalancesType
		}
	}

	return collection, nil
}