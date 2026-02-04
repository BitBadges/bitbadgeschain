package approval_criteria

import (
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ApprovalCriteriaChecker defines the interface for checking approval criteria.
// All implementations should return (deterministicErrorMsg, error) where:
// - deterministicErrorMsg is a user-friendly error message if the check fails
// - error is nil if the check passes, or an error if the check fails
type ApprovalCriteriaChecker interface {
	Name() string
	Check(ctx sdk.Context, approval *types.CollectionApproval, collection *types.TokenCollection, to string, from string, initiator string, approvalLevel string, approverAddress string, merkleProofs []*types.MerkleProof, ethSignatureProofs []*types.ETHSignatureProof, memo string, isPrioritized bool) (detErrMsg string, err error)
}
