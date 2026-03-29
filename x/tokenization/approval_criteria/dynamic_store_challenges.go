package approval_criteria

import (
	"fmt"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DynamicStoreChallengesChecker implements ApprovalCriteriaChecker for DynamicStoreChallenges
type DynamicStoreChallengesChecker struct {
	dynamicStoreService DynamicStoreService
}

// NewDynamicStoreChallengesChecker creates a new DynamicStoreChallengesChecker
func NewDynamicStoreChallengesChecker(dynamicStoreService DynamicStoreService) *DynamicStoreChallengesChecker {
	return &DynamicStoreChallengesChecker{
		dynamicStoreService: dynamicStoreService,
	}
}

// Name returns the name of this checker
func (c *DynamicStoreChallengesChecker) Name() string {
	return "DynamicStoreChallenges"
}

// Check validates dynamic store challenges for the appropriate party
// It checks if the specified party's value satisfies the comparison condition (read-only check)
func (c *DynamicStoreChallengesChecker) Check(ctx sdk.Context, approval *types.CollectionApproval, collection *types.TokenCollection, to string, from string, initiator string, approvalLevel string, approverAddress string, merkleProofs []*types.MerkleProof, ethSignatureProofs []*types.ETHSignatureProof, memo string, isPrioritized bool) (string, error) {
	if approval.ApprovalCriteria == nil || len(approval.ApprovalCriteria.DynamicStoreChallenges) == 0 {
		return "", nil
	}

	challenges := approval.ApprovalCriteria.DynamicStoreChallenges
	for _, challenge := range challenges {
		if challenge == nil {
			detErrMsg := "challenge is nil"
			return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
		}

		storeId := challenge.StoreId

		// Get the dynamic store to check global kill switch first
		dynamicStore, foundStore := c.dynamicStoreService.GetDynamicStore(ctx, storeId)
		if !foundStore {
			detErrMsg := fmt.Sprintf("dynamic store not found for storeId %s", storeId.String())
			return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
		}

		// First check global kill switch - if disabled, fail immediately
		if !dynamicStore.GlobalEnabled {
			detErrMsg := fmt.Sprintf("dynamic store storeId %s is globally disabled", storeId.String())
			return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
		}

		// Determine which party to check ownership for
		partyToCheck := c.determinePartyToCheck(challenge.OwnershipCheckParty, initiator, from, to, collection)

		// If globalEnabled = true, proceed with per-address logic
		// Get the current value for the determined party
		dynamicStoreValue, found := c.dynamicStoreService.GetDynamicStoreValue(ctx, storeId, partyToCheck)

		var val sdkmath.Uint
		if found {
			val = dynamicStoreValue.Value
		} else {
			// If no specific value found, use the default value from the store
			val = dynamicStore.DefaultValue
		}

		// Evaluate the comparison
		if !evaluateComparison(val, challenge.ComparisonOperator, challenge.ComparisonValue) {
			detErrMsg := fmt.Sprintf("%s does not satisfy dynamic store challenge for storeId %s (value=%s, operator=%s, comparisonValue=%s)",
				partyToCheck, storeId.String(), val.String(), challenge.ComparisonOperator, challenge.ComparisonValue.String())
			return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
		}
	}

	return "", nil
}

// evaluateComparison evaluates a numeric comparison between val and comparisonValue.
// When operator is empty, defaults to legacy behavior: val != 0 (backward compatible with old boolean check).
func evaluateComparison(val sdkmath.Uint, operator string, comparisonValue sdkmath.Uint) bool {
	switch operator {
	case "eq":
		return val.Equal(comparisonValue)
	case "ne":
		return !val.Equal(comparisonValue)
	case "gt":
		return val.GT(comparisonValue)
	case "gte":
		return val.GTE(comparisonValue)
	case "lt":
		return val.LT(comparisonValue)
	case "lte":
		return val.LTE(comparisonValue)
	case "":
		// Legacy behavior: value != 0 (backward compatible with old boolean true/false)
		return !val.IsZero()
	default:
		// Unknown operator - fail safely
		return false
	}
}

// determinePartyToCheck determines which party's ownership should be checked
func (c *DynamicStoreChallengesChecker) determinePartyToCheck(
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
