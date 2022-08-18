package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

//Gets badge and throws error if it does not exist. Alternative to GetBadgeFromStore which returns a found bool, not an error.
func (k Keeper) GetBadgeE(ctx sdk.Context, badgeId uint64) (types.BitBadge, error) {
	badge, found := k.GetBadgeFromStore(ctx, badgeId)
	if !found {
		return types.BitBadge{}, ErrBadgeNotExists
	}

	return badge, nil
}

// Gets the badge details from the store if it exists. Throws error if subbadge ranges are invalid or the subbadge does not yet exist.
func (k Keeper) GetBadgeAndAssertSubbadgeRangesAreValid(ctx sdk.Context, badgeId uint64, subbadgeRanges []*types.IdRange) (types.BitBadge, error) {
	badge, err := k.GetBadgeE(ctx, badgeId)
	if err != nil {
		return badge, err
	}

	for _, subbadgeRange := range subbadgeRanges {
		subbadgeRange = NormalizeIdRange(subbadgeRange)

		if subbadgeRange.Start > subbadgeRange.End {
			return types.BitBadge{}, ErrInvalidSubbadgeRange
		}

		if subbadgeRange.End >= badge.NextSubassetId {
			return types.BitBadge{}, ErrSubBadgeNotExists
		}
	}

	return badge, nil
}

//For each (supply, amountToCreate) pair, we create amountToCreate subbadges with specified supply. We also mint total supply to manager. Error if IDs overflow.
//We assume that lengths of supplys and amountsToCreate are equal before entering this function. Also amountsToCreate[i] can never be zero.
func CreateSubassets(badge types.BitBadge, managerBalanceInfo types.UserBalanceInfo, supplys []uint64, amounts []uint64) (types.BitBadge, types.UserBalanceInfo, error) {
	newSubassetSupplys := badge.SubassetSupplys
	defaultSupply := badge.DefaultSubassetSupply
	if badge.DefaultSubassetSupply == 0 {
		defaultSupply = 1
	}

	err := *new(error)
	//Update supplys and mint total supply for each to manager. Don't store if supply == default
	//Subasset supplys are stored as []*types.BalanceObject, so we can use the balance update functions
	for i, supply := range supplys {
		amountToCreate := amounts[i]
		nextSubassetId := badge.NextSubassetId

		// We conventionalize supply == 0 as default, and we don't store if supply == default
		if supply != 0 && supply != defaultSupply {
			newSubassetSupplys = UpdateBalancesForIdRanges(
				[]*types.IdRange{
					{Start: nextSubassetId, End: nextSubassetId + amountToCreate - 1},
				}, 
				supply, 
				newSubassetSupplys,
			)
		}

		managerBalanceInfo, err = AddBalancesForIdRanges(managerBalanceInfo, []*types.IdRange{{Start: nextSubassetId, End: nextSubassetId + amountToCreate - 1}}, supply)
		if err != nil {
			return types.BitBadge{}, types.UserBalanceInfo{}, err
		}

		badge.NextSubassetId, err = SafeAdd(badge.NextSubassetId, amountToCreate) //error on ID overflow
		if err != nil {
			return types.BitBadge{}, types.UserBalanceInfo{}, err
		}
	}
	badge.SubassetSupplys = newSubassetSupplys

	return badge, managerBalanceInfo, nil
}
