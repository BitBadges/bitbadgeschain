package keeper

import (
	sdkmath "cosmossdk.io/math"

	sdkerrors "cosmossdk.io/errors"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CheckMustOwnTokens checks if the initiatedBy address owns the required tokens
func (k Keeper) CheckMustOwnTokens(
	ctx sdk.Context,
	mustOwnTokens []*types.MustOwnTokens,
	initiatedBy string,
	fromAddress string,
	toAddress string,
) error {
	failedMustOwnTokens := false
	for _, mustOwnToken := range mustOwnTokens {
		collection, found := k.GetCollectionFromStore(ctx, mustOwnToken.CollectionId)
		if !found {
			failedMustOwnTokens = true
			break
		}

		// Determine which party to check ownership for
		partyToCheck := initiatedBy // default to initiator

		switch mustOwnToken.OwnershipCheckParty {
		case "initiator":
			partyToCheck = initiatedBy
		case "sender":
			partyToCheck = fromAddress
		case "recipient":
			partyToCheck = toAddress
		case types.MintAddress:
			partyToCheck = collection.MintEscrowAddress
		case "":
			partyToCheck = initiatedBy
		default:
			// TODO: We could support any bb1 address here as well
			partyToCheck = initiatedBy
		}

		partyBalances, _ := k.GetBalanceOrApplyDefault(ctx, collection, partyToCheck)
		balances := partyBalances.Balances

		if mustOwnToken.OverrideWithCurrentTime {
			currTime := sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli()))
			mustOwnToken.OwnershipTimes = []*types.UintRange{{Start: currTime, End: currTime}}
		}

		fetchedBalances, err := types.GetBalancesForIds(ctx, mustOwnToken.TokenIds, mustOwnToken.OwnershipTimes, balances)
		if err != nil {
			failedMustOwnTokens = true
			break
		}

		satisfiesRequirementsForOne := false
		for _, fetchedBalance := range fetchedBalances {
			//check if amount is within range
			minAmount := mustOwnToken.AmountRange.Start
			maxAmount := mustOwnToken.AmountRange.End

			if fetchedBalance.Amount.LT(minAmount) || fetchedBalance.Amount.GT(maxAmount) {
				failedMustOwnTokens = true
			} else {
				satisfiesRequirementsForOne = true
			}
		}

		if mustOwnToken.MustSatisfyForAllAssets && failedMustOwnTokens {
			break
		} else if !mustOwnToken.MustSatisfyForAllAssets && satisfiesRequirementsForOne {
			failedMustOwnTokens = false
			break
		}
	}

	if failedMustOwnTokens {
		return sdkerrors.Wrapf(ErrInadequateApprovals, "failed token ownership requirements")
	}

	return nil
}
