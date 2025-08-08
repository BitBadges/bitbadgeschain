package keeper

import (
	sdkmath "cosmossdk.io/math"

	sdkerrors "cosmossdk.io/errors"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CheckMustOwnBadges checks if the initiatedBy address owns the required tokens
func (k Keeper) CheckMustOwnBadges(
	ctx sdk.Context,
	mustOwnBadges []*types.MustOwnBadges,
	initiatedBy string,
) error {
	failedMustOwnBadges := false
	for _, mustOwnBadge := range mustOwnBadges {
		collection, found := k.GetCollectionFromStore(ctx, mustOwnBadge.CollectionId)
		if !found {
			failedMustOwnBadges = true
			break
		}

		initiatorBalances, _ := k.GetBalanceOrApplyDefault(ctx, collection, initiatedBy)
		balances := initiatorBalances.Balances

		if mustOwnBadge.OverrideWithCurrentTime {
			currTime := sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli()))
			mustOwnBadge.OwnershipTimes = []*types.UintRange{{Start: currTime, End: currTime}}
		}

		fetchedBalances, err := types.GetBalancesForIds(ctx, mustOwnBadge.BadgeIds, mustOwnBadge.OwnershipTimes, balances)
		if err != nil {
			failedMustOwnBadges = true
			break
		}

		satisfiesRequirementsForOne := false
		for _, fetchedBalance := range fetchedBalances {
			//check if amount is within range
			minAmount := mustOwnBadge.AmountRange.Start
			maxAmount := mustOwnBadge.AmountRange.End

			if fetchedBalance.Amount.LT(minAmount) || fetchedBalance.Amount.GT(maxAmount) {
				failedMustOwnBadges = true
			} else {
				satisfiesRequirementsForOne = true
			}
		}

		if mustOwnBadge.MustSatisfyForAllAssets && failedMustOwnBadges {
			break
		} else if !mustOwnBadge.MustSatisfyForAllAssets && satisfiesRequirementsForOne {
			failedMustOwnBadges = false
			break
		}
	}

	if failedMustOwnBadges {
		return sdkerrors.Wrapf(ErrInadequateApprovals, "failed token ownership requirements")
	}

	return nil
}