package approval_criteria

import (
	"fmt"

	sdkerrors "cosmossdk.io/errors"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
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

// Check validates dynamic store challenges for the initiator
// It checks if the initiator has a true value for each challenge (read-only check)
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

		// If globalEnabled = true, proceed with existing per-address logic
		// Get the current value for the initiator
		dynamicStoreValue, found := c.dynamicStoreService.GetDynamicStoreValue(ctx, storeId, initiator)

		var val bool
		if found {
			val = dynamicStoreValue.Value
		} else {
			// If no specific value found, use the default value from the store
			val = dynamicStore.DefaultValue
		}

		// Check if the initiator has a true value (read-only check, no updates)
		if !val {
			detErrMsg := fmt.Sprintf("initiator does not have permission for dynamic store challenge storeId %s", storeId.String())
			return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
		}
	}

	return "", nil
}
