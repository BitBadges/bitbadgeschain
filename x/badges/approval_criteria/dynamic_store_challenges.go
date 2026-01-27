package approval_criteria

import (
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

// Check validates dynamic store challenges for the appropriate party
// It checks if the specified party has a true value for each challenge (read-only check)
func (c *DynamicStoreChallengesChecker) Check(ctx sdk.Context, approval *types.CollectionApproval, collection *types.TokenCollection, to string, from string, initiator string, approvalLevel string, approverAddress string, merkleProofs []*types.MerkleProof, ethSignatureProofs []*types.ETHSignatureProof, memo string, isPrioritized bool) (string, error) {
	if approval.ApprovalCriteria == nil || len(approval.ApprovalCriteria.DynamicStoreChallenges) == 0 {
		return "", nil
	}

	challenges := approval.ApprovalCriteria.DynamicStoreChallenges
	for idx, challenge := range challenges {
		if challenge == nil {
			detErrMsg := "challenge is nil"
			return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
		}

		// Use shared helper function
		passed, errMsg := CheckDynamicStoreChallenge(ctx, challenge, idx, initiator, from, to, c.dynamicStoreService, collection, "dynamic store challenge")
		if !passed {
			return errMsg, sdkerrors.Wrap(types.ErrInvalidRequest, errMsg)
		}
	}

	return "", nil
}
