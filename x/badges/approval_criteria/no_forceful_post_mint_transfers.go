package approval_criteria

import (
	sdkerrors "cosmossdk.io/errors"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NoForcefulPostMintTransfersChecker implements ApprovalCriteriaChecker for the NoForcefulPostMintTransfers invariant
type NoForcefulPostMintTransfersChecker struct{}

// NewNoForcefulPostMintTransfersChecker creates a new NoForcefulPostMintTransfersChecker
func NewNoForcefulPostMintTransfersChecker() *NoForcefulPostMintTransfersChecker {
	return &NoForcefulPostMintTransfersChecker{}
}

// Name returns the name of this checker
func (c *NoForcefulPostMintTransfersChecker) Name() string {
	return "NoForcefulPostMintTransfers"
}

// Check validates that forceful transfers (with overrides) are disallowed when the NoForcefulPostMintTransfers invariant is enabled
// This only applies when fromAddress does not equal "Mint"
func (c *NoForcefulPostMintTransfersChecker) Check(ctx sdk.Context, approval *types.CollectionApproval, collection *types.TokenCollection, to string, from string, initiator string, approvalLevel string, approverAddress string, merkleProofs []*types.MerkleProof, ethSignatureProofs []*types.ETHSignatureProof, memo string, isPrioritized bool) (string, error) {
	// Only check if collection is provided and the invariant is enabled
	if collection == nil || collection.Invariants == nil || !collection.Invariants.NoForcefulPostMintTransfers {
		return "", nil
	}

	// Only check if fromAddress is not "Mint"
	if from == types.MintAddress {
		return "", nil
	}

	approvalCriteria := approval.ApprovalCriteria
	if approvalCriteria == nil {
		return "", nil
	}

	// Check for overrides that bypass user approvals
	if approvalCriteria.OverridesFromOutgoingApprovals {
		detErrMsg := "forceful transfers that bypass user outgoing approvals are disallowed when noForcefulPostMintTransfers invariant is enabled (unless from address is Mint)"
		return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
	}

	if approvalCriteria.OverridesToIncomingApprovals {
		detErrMsg := "forceful transfers that bypass user incoming approvals are disallowed when noForcefulPostMintTransfers invariant is enabled (unless from address is Mint)"
		return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
	}

	return "", nil
}
