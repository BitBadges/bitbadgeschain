package approval_criteria

import (
	"fmt"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MustOwnTokensChecker implements ApprovalCriteriaChecker for MustOwnTokens checks
type MustOwnTokensChecker struct {
	collectionService CollectionService
}

// NewMustOwnTokensChecker creates a new MustOwnTokensChecker
func NewMustOwnTokensChecker(collectionService CollectionService) *MustOwnTokensChecker {
	return &MustOwnTokensChecker{
		collectionService: collectionService,
	}
}

// Name returns the name of this checker
func (c *MustOwnTokensChecker) Name() string {
	return "MustOwnTokens"
}

// Check validates that the required tokens are owned by the appropriate party
func (c *MustOwnTokensChecker) Check(ctx sdk.Context, approval *types.CollectionApproval, collection *types.TokenCollection, to string, from string, initiator string, approvalLevel string, approverAddress string, merkleProofs []*types.MerkleProof, ethSignatureProofs []*types.ETHSignatureProof, memo string, isPrioritized bool) (string, error) {
	if approval.ApprovalCriteria == nil || len(approval.ApprovalCriteria.MustOwnTokens) == 0 {
		return "", nil
	}

	mustOwnTokens := approval.ApprovalCriteria.MustOwnTokens
	for idx, mustOwnToken := range mustOwnTokens {
		// Check if this requirement passes
		requirementPassed, errMsg := c.checkSingleRequirement(ctx, mustOwnToken, idx, initiator, from, to)
		if !requirementPassed {
			return errMsg, sdkerrors.Wrap(types.ErrInvalidRequest, errMsg)
		}
	}

	return "", nil
}

// checkSingleRequirement checks if a single MustOwnTokens requirement is satisfied
// Returns (passed bool, errorMsg string)
func (c *MustOwnTokensChecker) checkSingleRequirement(
	ctx sdk.Context,
	mustOwnToken *types.MustOwnTokens,
	requirementIdx int,
	initiatedBy string,
	fromAddress string,
	toAddress string,
) (bool, string) {
	// Check if collection exists
	collection, found := c.collectionService.GetCollection(ctx, mustOwnToken.CollectionId)
	if !found {
		errMsg := fmt.Sprintf("token ownership requirement idx %d failed: collection %s not found",
			requirementIdx, mustOwnToken.CollectionId.String())
		return false, errMsg
	}

	// Determine which party to check ownership for
	partyToCheck := c.determinePartyToCheck(mustOwnToken.OwnershipCheckParty, initiatedBy, fromAddress, toAddress, collection)

	// Get balances for the party
	partyBalances, _ := c.collectionService.GetBalanceOrApplyDefault(ctx, collection, partyToCheck)
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
	requirementPassed := c.checkAmountRange(fetchedBalances, mustOwnToken.AmountRange, mustOwnToken.MustSatisfyForAllAssets)

	if !requirementPassed {
		errMsg := fmt.Sprintf("token ownership requirement idx %d failed: party %s does not meet requirements for collection %s",
			requirementIdx, partyToCheck, mustOwnToken.CollectionId.String())
		return false, errMsg
	}

	return true, ""
}

// determinePartyToCheck determines which party's ownership should be checked
func (c *MustOwnTokensChecker) determinePartyToCheck(
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
		if _, err := sdk.AccAddressFromBech32(ownershipCheckParty); err == nil {
			return ownershipCheckParty
		}

		// If not a valid address, fall back to default behavior
		return initiatedBy
	}
}

// checkAmountRange checks if the fetched balances satisfy the amount range requirement
// Returns true if the requirement is satisfied, false otherwise
func (c *MustOwnTokensChecker) checkAmountRange(
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
