package keeper

import (
	sdkmath "cosmossdk.io/math"

	sdkerrors "cosmossdk.io/errors"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CheckMustOwnBadges checks if the specified party owns the required badges
func (k Keeper) CheckMustOwnBadges(
	ctx sdk.Context,
	mustOwnBadges []*types.MustOwnBadges,
	initiatedBy string,
	fromAddress string,
	toAddress string,
) error {
	//Assert that the specified party owns the required badges
	failedMustOwnBadges := false
	for _, mustOwnBadge := range mustOwnBadges {
		collection, found := k.GetCollectionFromStore(ctx, mustOwnBadge.CollectionId)
		if !found {
			failedMustOwnBadges = true
			break
		}

		// Determine which party to check ownership for
		partyToCheck := initiatedBy // default to initiator
		if mustOwnBadge.OwnershipCheckParty != "" {
			switch mustOwnBadge.OwnershipCheckParty {
			case "initiator":
				partyToCheck = initiatedBy
			case "sender":
				partyToCheck = fromAddress
			case "recipient":
				partyToCheck = toAddress
			default:
				// If invalid value, default to initiator
				partyToCheck = initiatedBy
			}
		}

		partyBalances, _ := k.GetBalanceOrApplyDefault(ctx, collection, partyToCheck)
		balances := partyBalances.Balances

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
		return sdkerrors.Wrapf(ErrInadequateApprovals, "failed badge ownership requirements")
	}

	return nil
}
