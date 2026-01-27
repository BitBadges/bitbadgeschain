package approval_criteria

import (
	sdkerrors "cosmossdk.io/errors"
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
	// Use shared helper function
	return CheckMustOwnTokensRequirement(ctx, mustOwnToken, requirementIdx, initiatedBy, fromAddress, toAddress, c.collectionService, "token ownership requirement")
}
