package approval_criteria

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CheckAmountRange checks if the fetched balances satisfy the amount range requirement.
// Returns true if the requirement is satisfied, false otherwise.
// This is a standalone helper function that can be reused across modules.
func CheckAmountRange(
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

// DeterminePartyToCheckForMustOwnTokens determines which party's ownership should be checked for MustOwnTokens.
// This is a standalone helper function that can be reused across modules.
func DeterminePartyToCheckForMustOwnTokens(
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
		// Check if ownershipCheckParty is a valid bb1 address
		// If it is, return it directly (allows checking ownership for arbitrary addresses)
		// Use types.ValidateAddress to ensure bb1 prefix is handled correctly
		if err := types.ValidateAddress(ownershipCheckParty, false); err == nil {
			return ownershipCheckParty
		}

		// If not a valid address, fall back to default behavior
		return initiatedBy
	}
}

// DeterminePartyToCheckForDynamicStore determines which party's dynamic store value should be checked.
// This is a standalone helper function that can be reused across modules.
func DeterminePartyToCheckForDynamicStore(
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
		// Check if ownershipCheckParty is a valid bb1 address
		// If it is, return it directly (allows checking for arbitrary addresses)
		// Use types.ValidateAddress to ensure bb1 prefix is handled correctly
		if err := types.ValidateAddress(ownershipCheckParty, false); err == nil {
			return ownershipCheckParty
		}

		// If not a valid address, fall back to default behavior
		return initiatedBy
	}
}

// CheckMustOwnTokensRequirement checks if a single MustOwnTokens requirement is satisfied.
// Returns (passed bool, errorMsg string).
// This is a standalone helper function that can be reused across modules.
func CheckMustOwnTokensRequirement(
	ctx sdk.Context,
	mustOwnToken *types.MustOwnTokens,
	requirementIdx int,
	initiatedBy string,
	fromAddress string,
	toAddress string,
	collectionService CollectionService,
	errorPrefix string, // Prefix for error messages (e.g., "2FA requirement" or "token ownership requirement")
) (bool, string) {
	// Check if collection exists
	collection, found := collectionService.GetCollection(ctx, mustOwnToken.CollectionId)
	if !found {
		errMsg := fmt.Sprintf("%s idx %d failed: collection %s not found",
			errorPrefix, requirementIdx, mustOwnToken.CollectionId.String())
		return false, errMsg
	}

	// Determine which party to check ownership for
	partyToCheck := DeterminePartyToCheckForMustOwnTokens(mustOwnToken.OwnershipCheckParty, initiatedBy, fromAddress, toAddress, collection)

	// Get balances for the party
	partyBalances, _, err := collectionService.GetBalanceOrApplyDefault(ctx, collection, partyToCheck)
	if err != nil {
		errMsg := fmt.Sprintf("%s idx %d failed: %v", errorPrefix, requirementIdx, err)
		return false, errMsg
	}
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
		errMsg := fmt.Sprintf("%s idx %d failed: party %s does not meet requirements for collection %s",
			errorPrefix, requirementIdx, partyToCheck, mustOwnToken.CollectionId.String())
		return false, errMsg
	}

	// Check if amounts are within the required range
	if mustOwnToken.AmountRange == nil {
		errMsg := fmt.Sprintf("%s idx %d failed: amount range is nil",
			errorPrefix, requirementIdx)
		return false, errMsg
	}

	// If no balances were fetched, the requirement is not met
	if len(fetchedBalances) == 0 {
		errMsg := fmt.Sprintf("%s idx %d failed: party %s does not own any matching tokens for collection %s",
			errorPrefix, requirementIdx, partyToCheck, mustOwnToken.CollectionId.String())
		return false, errMsg
	}

	requirementPassed := CheckAmountRange(fetchedBalances, mustOwnToken.AmountRange, mustOwnToken.MustSatisfyForAllAssets)

	if !requirementPassed {
		errMsg := fmt.Sprintf("%s idx %d failed: party %s does not meet amount requirements for collection %s",
			errorPrefix, requirementIdx, partyToCheck, mustOwnToken.CollectionId.String())
		return false, errMsg
	}

	return true, ""
}

// CheckDynamicStoreChallenge checks if a single DynamicStoreChallenge requirement is satisfied.
// Returns (passed bool, errorMsg string).
// This is a standalone helper function that can be reused across modules.
func CheckDynamicStoreChallenge(
	ctx sdk.Context,
	challenge *types.DynamicStoreChallenge,
	challengeIdx int,
	initiatedBy string,
	fromAddress string,
	toAddress string,
	dynamicStoreService DynamicStoreService,
	collection *types.TokenCollection, // Can be nil if not needed
	errorPrefix string, // Prefix for error messages (e.g., "2FA dynamic store challenge" or "dynamic store challenge")
) (bool, string) {
	if challenge == nil {
		errMsg := fmt.Sprintf("%s idx %d is nil", errorPrefix, challengeIdx)
		return false, errMsg
	}

	storeId := challenge.StoreId

	// Get the dynamic store to check global kill switch first
	dynamicStore, foundStore := dynamicStoreService.GetDynamicStore(ctx, storeId)
	if !foundStore {
		errMsg := fmt.Sprintf("%s idx %d failed: dynamic store %s not found",
			errorPrefix, challengeIdx, storeId.String())
		return false, errMsg
	}

	// First check global kill switch - if disabled, fail immediately
	if !dynamicStore.GlobalEnabled {
		errMsg := fmt.Sprintf("%s idx %d failed: dynamic store %s is globally disabled",
			errorPrefix, challengeIdx, storeId.String())
		return false, errMsg
	}

	// Determine which party to check
	partyToCheck := DeterminePartyToCheckForDynamicStore(challenge.OwnershipCheckParty, initiatedBy, fromAddress, toAddress, collection)

	// Get the current value for the determined party
	dynamicStoreValue, found := dynamicStoreService.GetDynamicStoreValue(ctx, storeId, partyToCheck)

	var val bool
	if found {
		val = dynamicStoreValue.Value
	} else {
		// If no specific value found, use the default value from the store
		val = dynamicStore.DefaultValue
	}

	// Check if the party has a true value
	if !val {
		errMsg := fmt.Sprintf("%s idx %d failed: party %s does not have permission for dynamic store %s",
			errorPrefix, challengeIdx, partyToCheck, storeId.String())
		return false, errMsg
	}

	return true, ""
}

