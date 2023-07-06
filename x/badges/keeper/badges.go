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

	unmintedSupplys, found := k.GetUserBalanceFromStore(ctx, ConstructBalanceKey("Mint", collection.CollectionId))
	if !found {
		unmintedSupplys = &types.UserBalanceStore{} //permissions and timelines do not matter because these are special addresses
	}

	totalSupplys, found := k.GetUserBalanceFromStore(ctx, ConstructBalanceKey("Total", collection.CollectionId))
	if !found {
		totalSupplys = &types.UserBalanceStore{} //permissions and timelines do not matter because these are special addresses
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

		totalSupplys.Balances, err = types.AddBalance(totalSupplys.Balances, balance)
		if err != nil {
			return &types.BadgeCollection{}, err
		}

		unmintedSupplys.Balances, err = types.AddBalance(unmintedSupplys.Balances, balance)
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