package approval_criteria

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TransferTimesChecker implements ApprovalCriteriaChecker for checking if transfer time is valid
type TransferTimesChecker struct{}

// NewTransferTimesChecker creates a new TransferTimesChecker
func NewTransferTimesChecker() *TransferTimesChecker {
	return &TransferTimesChecker{}
}

// Name returns the name of this checker
func (c *TransferTimesChecker) Name() string {
	return "TransferTimes"
}

// Check validates that the current time falls within the approval's TransferTimes range
func (c *TransferTimesChecker) Check(ctx sdk.Context, approval *types.CollectionApproval, collection *types.TokenCollection, to string, from string, initiator string, approvalLevel string, approverAddress string, merkleProofs []*types.MerkleProof, ethSignatureProofs []*types.ETHSignatureProof, memo string, isPrioritized bool) (string, error) {
	if approval == nil {
		return "", sdkerrors.Wrap(types.ErrInvalidRequest, "approval cannot be nil")
	}

	currTime := sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli()))
	currTimeFound, err := types.SearchUintRangesForUint(currTime, approval.TransferTimes)
	if !currTimeFound || err != nil {
		detErrMsg := "transfer time not in range"
		return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
	}

	return "", nil
}
