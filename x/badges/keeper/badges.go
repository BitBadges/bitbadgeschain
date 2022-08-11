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
		// Subbadge ranges can set end == 0 to save storage space. By convention, this means end == start
		if subbadgeRange.End == 0 {
			subbadgeRange.End = subbadgeRange.Start
		} 

		if subbadgeRange.Start > subbadgeRange.End {
			return types.BitBadge{}, ErrInvalidSubbadgeRange
		}

		if subbadgeRange.End >= badge.NextSubassetId {
			return types.BitBadge{}, ErrSubBadgeNotExists
		}
	}

	return badge, nil
}

func CreateSubassets(ctx sdk.Context, badge types.BitBadge, managerBalanceInfo types.UserBalanceInfo, supplys []uint64, amounts[]uint64) (types.BitBadge, types.UserBalanceInfo, error) {
	newSubassetSupplys := badge.SubassetSupplys
	defaultSupply := badge.DefaultSubassetSupply
	if badge.DefaultSubassetSupply == 0 {
		defaultSupply = 1
	}

	err := *new(error)
	// Update supplys and mint total supply for each to manager. Don't store if supply == default
	for i, supply := range supplys {
		for j := uint64(0); j < amounts[i]; j++ {
			nextSubassetId := badge.NextSubassetId

			// We conventionalize supply == 0 as default, so we don't store if it is the default
			if supply != 0 && supply != defaultSupply {
				newSubassetSupplys = UpdateBalanceForId(nextSubassetId, supply, newSubassetSupplys)
			}
			badge.NextSubassetId += 1

			managerBalanceInfo, err = AddBalanceForId(ctx, managerBalanceInfo, nextSubassetId, supply)
			if err != nil {
				return types.BitBadge{}, types.UserBalanceInfo{}, err
			}
		}
	}
	badge.SubassetSupplys = newSubassetSupplys

	return badge, managerBalanceInfo, nil
}