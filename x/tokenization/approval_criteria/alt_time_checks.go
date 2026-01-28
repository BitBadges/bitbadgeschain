package approval_criteria

import (
	"fmt"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AltTimeChecksChecker implements ApprovalCriteriaChecker for AltTimeChecks
type AltTimeChecksChecker struct {
	altTimeChecks *types.AltTimeChecks
}

// NewAltTimeChecksChecker creates a new AltTimeChecksChecker
func NewAltTimeChecksChecker(altTimeChecks *types.AltTimeChecks) *AltTimeChecksChecker {
	return &AltTimeChecksChecker{
		altTimeChecks: altTimeChecks,
	}
}

// Name returns the name of this checker
func (c *AltTimeChecksChecker) Name() string {
	return "AltTimeChecks"
}

// Check validates if the current UTC time falls within any offline hours or days specified in altTimeChecks.
// Returns an error if the time is within an offline period (i.e., the approval should be denied).
func (c *AltTimeChecksChecker) Check(ctx sdk.Context, approval *types.CollectionApproval, collection *types.TokenCollection, to string, from string, initiator string, approvalLevel string, approverAddress string, merkleProofs []*types.MerkleProof, ethSignatureProofs []*types.ETHSignatureProof, memo string, isPrioritized bool) (string, error) {
	if c.altTimeChecks == nil {
		return "", nil
	}

	// Get UTC time from block time
	blockTime := ctx.BlockTime()
	utcTime := blockTime.UTC()

	// Get current hour (0-23) and day of week (0=Sunday, 1=Monday, ..., 6=Saturday)
	currentHour := sdkmath.NewUint(uint64(utcTime.Hour()))
	currentDay := sdkmath.NewUint(uint64(utcTime.Weekday()))

	// Check if current hour falls within any offline hours range
	if len(c.altTimeChecks.OfflineHours) > 0 {
		detErrMsg := fmt.Sprintf("transfer denied: current UTC hour %d falls within offline hours", currentHour.Uint64())
		found, err := types.SearchUintRangesForUint(currentHour, c.altTimeChecks.OfflineHours)
		if err != nil {
			return detErrMsg, sdkerrors.Wrapf(err, "error checking offline hours")
		}
		if found {
			return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
		}
	}

	// Check if current day falls within any offline days range
	if len(c.altTimeChecks.OfflineDays) > 0 {
		detErrMsg := fmt.Sprintf("transfer denied: current UTC day %d falls within offline days", currentDay.Uint64())
		found, err := types.SearchUintRangesForUint(currentDay, c.altTimeChecks.OfflineDays)
		if err != nil {
			return detErrMsg, sdkerrors.Wrapf(err, "error checking offline days")
		}
		if found {
			return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
		}
	}

	return "", nil
}
