package keeper

import (
	"fmt"

	sdkmath "cosmossdk.io/math"

	sdkerrors "cosmossdk.io/errors"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CheckMustOwnTokens checks if the initiatedBy address owns the required tokens
// Returns (deterministicErrorMsg, error) where deterministicErrorMsg is a deterministic error string
// All requirements in the array must pass for the overall check to succeed
func (k Keeper) CheckMustOwnTokens(
	ctx sdk.Context,
	mustOwnTokens []*types.MustOwnTokens,
	initiatedBy string,
	fromAddress string,
	toAddress string,
) (string, error) {
	for idx, mustOwnToken := range mustOwnTokens {
		// Check if this requirement passes
		requirementPassed, errMsg := k.checkSingleRequirement(ctx, mustOwnToken, idx, initiatedBy, fromAddress, toAddress)
		if !requirementPassed {
			return errMsg, sdkerrors.Wrap(ErrInadequateApprovals, errMsg)
		}
	}

	return "", nil
}

// checkSingleRequirement checks if a single MustOwnTokens requirement is satisfied
// Returns (passed bool, errorMsg string)
func (k Keeper) checkSingleRequirement(
	ctx sdk.Context,
	mustOwnToken *types.MustOwnTokens,
	requirementIdx int,
	initiatedBy string,
	fromAddress string,
	toAddress string,
) (bool, string) {
	// Check if collection exists
	collection, found := k.GetCollectionFromStore(ctx, mustOwnToken.CollectionId)
	if !found {
		errMsg := fmt.Sprintf("token ownership requirement idx %d failed: collection %s not found",
			requirementIdx, mustOwnToken.CollectionId.String())
		return false, errMsg
	}

	// Determine which party to check ownership for
	partyToCheck := k.determinePartyToCheck(mustOwnToken.OwnershipCheckParty, initiatedBy, fromAddress, toAddress, collection)

	// Get balances for the party
	partyBalances, _ := k.GetBalanceOrApplyDefault(ctx, collection, partyToCheck)
	balances := partyBalances.Balances

	// Determine ownership times to use (override with current time if needed)
	ownershipTimesToUse := mustOwnToken.OwnershipTimes
	if mustOwnToken.OverrideWithCurrentTime {
		currTime := sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli()))
		ownershipTimesToUse = []*types.UintRange{{Start: currTime, End: currTime}}
	}

	// Fetch balances matching the token IDs and ownership times
	fetchedBalances, err := types.GetBalancesForIds(ctx, mustOwnToken.TokenIds, ownershipTimesToUse, balances)
	if err != nil {
		errMsg := fmt.Sprintf("token ownership requirement idx %d failed: party %s does not meet requirements for collection %s",
			requirementIdx, partyToCheck, mustOwnToken.CollectionId.String())
		return false, errMsg
	}

	// Check if amounts are within the required range
	if mustOwnToken.AmountRange == nil {
		errMsg := fmt.Sprintf("token ownership requirement idx %d failed: amount range is nil",
			requirementIdx)
		return false, errMsg
	}
	requirementPassed := k.checkAmountRange(fetchedBalances, mustOwnToken.AmountRange, mustOwnToken.MustSatisfyForAllAssets)

	if !requirementPassed {
		errMsg := fmt.Sprintf("token ownership requirement idx %d failed: party %s does not meet requirements for collection %s",
			requirementIdx, partyToCheck, mustOwnToken.CollectionId.String())
		return false, errMsg
	}

	return true, ""
}

// determinePartyToCheck determines which party's ownership should be checked
func (k Keeper) determinePartyToCheck(
	ownershipCheckParty string,
	initiatedBy string,
	fromAddress string,
	toAddress string,
	collection *types.TokenCollection,
) string {
	switch ownershipCheckParty {
	case "initiator":
		return initiatedBy
	case "sender":
		return fromAddress
	case "recipient":
		return toAddress
	case types.MintAddress:
		return collection.MintEscrowAddress
	case "":
		return initiatedBy
	default:
		// TODO: We could support any bb1 address here as well
		return initiatedBy
	}
}

// checkAmountRange checks if the fetched balances satisfy the amount range requirement
// Returns true if the requirement is satisfied, false otherwise
func (k Keeper) checkAmountRange(
	fetchedBalances []*types.Balance,
	amountRange *types.UintRange,
	mustSatisfyForAllAssets bool,
) bool {
	if len(fetchedBalances) == 0 {
		return false
	}

	minAmount := amountRange.Start
	maxAmount := amountRange.End

	hasAtLeastOnePass := false
	hasAnyFail := false

	for _, balance := range fetchedBalances {
		amountInRange := !balance.Amount.LT(minAmount) && !balance.Amount.GT(maxAmount)

		if amountInRange {
			hasAtLeastOnePass = true
		} else {
			hasAnyFail = true
		}
	}

	// If MustSatisfyForAllAssets is true, all balances must pass (no failures)
	// If MustSatisfyForAllAssets is false, at least one balance must pass
	if mustSatisfyForAllAssets {
		return !hasAnyFail
	}
	return hasAtLeastOnePass
}
