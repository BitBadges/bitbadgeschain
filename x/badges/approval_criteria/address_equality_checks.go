package approval_criteria

import (
	sdkerrors "cosmossdk.io/errors"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RequireFromDoesNotEqualInitiatedByChecker checks that from address does not equal initiated by
type RequireFromDoesNotEqualInitiatedByChecker struct{}

// NewRequireFromDoesNotEqualInitiatedByChecker creates a new RequireFromDoesNotEqualInitiatedByChecker
func NewRequireFromDoesNotEqualInitiatedByChecker() *RequireFromDoesNotEqualInitiatedByChecker {
	return &RequireFromDoesNotEqualInitiatedByChecker{}
}

// Name returns the name of this checker
func (c *RequireFromDoesNotEqualInitiatedByChecker) Name() string {
	return "RequireFromDoesNotEqualInitiatedBy"
}

// Check validates that from address does not equal initiated by
func (c *RequireFromDoesNotEqualInitiatedByChecker) Check(ctx sdk.Context, approval *types.CollectionApproval, collection *types.TokenCollection, to string, from string, initiator string, approvalLevel string, approverAddress string, merkleProofs []*types.MerkleProof, ethSignatureProofs []*types.ETHSignatureProof, memo string, isPrioritized bool) (string, error) {
	if approval.ApprovalCriteria == nil || !approval.ApprovalCriteria.RequireFromDoesNotEqualInitiatedBy {
		return "", nil
	}

	if from == initiator {
		detErrMsg := "from address equals initiated by"
		return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
	}

	return "", nil
}

// RequireFromEqualsInitiatedByChecker checks that from address equals initiated by
type RequireFromEqualsInitiatedByChecker struct{}

// NewRequireFromEqualsInitiatedByChecker creates a new RequireFromEqualsInitiatedByChecker
func NewRequireFromEqualsInitiatedByChecker() *RequireFromEqualsInitiatedByChecker {
	return &RequireFromEqualsInitiatedByChecker{}
}

// Name returns the name of this checker
func (c *RequireFromEqualsInitiatedByChecker) Name() string {
	return "RequireFromEqualsInitiatedBy"
}

// Check validates that from address equals initiated by
func (c *RequireFromEqualsInitiatedByChecker) Check(ctx sdk.Context, approval *types.CollectionApproval, collection *types.TokenCollection, to string, from string, initiator string, approvalLevel string, approverAddress string, merkleProofs []*types.MerkleProof, ethSignatureProofs []*types.ETHSignatureProof, memo string, isPrioritized bool) (string, error) {
	if approval.ApprovalCriteria == nil || !approval.ApprovalCriteria.RequireFromEqualsInitiatedBy {
		return "", nil
	}

	if from != initiator {
		detErrMsg := "from address does not equal initiated by"
		return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
	}

	return "", nil
}

// RequireToDoesNotEqualInitiatedByChecker checks that to address does not equal initiated by
type RequireToDoesNotEqualInitiatedByChecker struct{}

// NewRequireToDoesNotEqualInitiatedByChecker creates a new RequireToDoesNotEqualInitiatedByChecker
func NewRequireToDoesNotEqualInitiatedByChecker() *RequireToDoesNotEqualInitiatedByChecker {
	return &RequireToDoesNotEqualInitiatedByChecker{}
}

// Name returns the name of this checker
func (c *RequireToDoesNotEqualInitiatedByChecker) Name() string {
	return "RequireToDoesNotEqualInitiatedBy"
}

// Check validates that to address does not equal initiated by
func (c *RequireToDoesNotEqualInitiatedByChecker) Check(ctx sdk.Context, approval *types.CollectionApproval, collection *types.TokenCollection, to string, from string, initiator string, approvalLevel string, approverAddress string, merkleProofs []*types.MerkleProof, ethSignatureProofs []*types.ETHSignatureProof, memo string, isPrioritized bool) (string, error) {
	if approval.ApprovalCriteria == nil || !approval.ApprovalCriteria.RequireToDoesNotEqualInitiatedBy {
		return "", nil
	}

	if to == initiator {
		detErrMsg := "to address equals initiated by"
		return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
	}

	return "", nil
}

// RequireToEqualsInitiatedByChecker checks that to address equals initiated by
type RequireToEqualsInitiatedByChecker struct{}

// NewRequireToEqualsInitiatedByChecker creates a new RequireToEqualsInitiatedByChecker
func NewRequireToEqualsInitiatedByChecker() *RequireToEqualsInitiatedByChecker {
	return &RequireToEqualsInitiatedByChecker{}
}

// Name returns the name of this checker
func (c *RequireToEqualsInitiatedByChecker) Name() string {
	return "RequireToEqualsInitiatedBy"
}

// Check validates that to address equals initiated by
func (c *RequireToEqualsInitiatedByChecker) Check(ctx sdk.Context, approval *types.CollectionApproval, collection *types.TokenCollection, to string, from string, initiator string, approvalLevel string, approverAddress string, merkleProofs []*types.MerkleProof, ethSignatureProofs []*types.ETHSignatureProof, memo string, isPrioritized bool) (string, error) {
	if approval.ApprovalCriteria == nil || !approval.ApprovalCriteria.RequireToEqualsInitiatedBy {
		return "", nil
	}

	if to != initiator {
		detErrMsg := "to address does not equal initiated by"
		return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
	}

	return "", nil
}
