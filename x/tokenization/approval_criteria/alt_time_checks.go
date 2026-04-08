package approval_criteria

import (
	"fmt"
	"time"

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

// Check validates if the current UTC time falls within any offline hours, days, months, or weeks specified in altTimeChecks.
// Returns an error if the time is within an offline period (i.e., the approval should be denied).
func (c *AltTimeChecksChecker) Check(ctx sdk.Context, approval *types.CollectionApproval, collection *types.TokenCollection, to string, from string, initiator string, approvalLevel string, approverAddress string, merkleProofs []*types.MerkleProof, ethSignatureProofs []*types.ETHSignatureProof, memo string, isPrioritized bool) (string, error) {
	if c.altTimeChecks == nil {
		return "", nil
	}

	// Get block time and apply timezone offset
	blockTime := ctx.BlockTime()
	localTime := blockTime.UTC()
	if c.altTimeChecks.TimezoneOffsetMinutes.GT(sdkmath.NewUint(0)) {
		offsetDuration := time.Duration(c.altTimeChecks.TimezoneOffsetMinutes.Uint64()) * time.Minute
		if c.altTimeChecks.TimezoneOffsetNegative {
			localTime = localTime.Add(-offsetDuration)
		} else {
			localTime = localTime.Add(offsetDuration)
		}
	}

	// Get current hour (0-23) and day of week (0=Sunday, 1=Monday, ..., 6=Saturday)
	currentHour := sdkmath.NewUint(uint64(localTime.Hour()))
	currentDay := sdkmath.NewUint(uint64(localTime.Weekday()))

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

	// Check if current month falls within any offline months range (1-12)
	if len(c.altTimeChecks.OfflineMonths) > 0 {
		currentMonth := sdkmath.NewUint(uint64(localTime.Month())) // time.Month is 1-12
		detErrMsg := fmt.Sprintf("transfer denied: current UTC month %d falls within offline months", currentMonth.Uint64())
		found, err := types.SearchUintRangesForUint(currentMonth, c.altTimeChecks.OfflineMonths)
		if err != nil {
			return detErrMsg, sdkerrors.Wrapf(err, "error checking offline months")
		}
		if found {
			return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
		}
	}

	// Check if current day of month falls within any offline days of month range (1-31)
	if len(c.altTimeChecks.OfflineDaysOfMonth) > 0 {
		dayOfMonth := sdkmath.NewUint(uint64(localTime.Day())) // 1-31
		detErrMsg := fmt.Sprintf("transfer denied: current UTC day of month %d falls within offline days of month", dayOfMonth.Uint64())
		found, err := types.SearchUintRangesForUint(dayOfMonth, c.altTimeChecks.OfflineDaysOfMonth)
		if err != nil {
			return detErrMsg, sdkerrors.Wrapf(err, "error checking offline days of month")
		}
		if found {
			return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
		}
	}

	// Check if current week of year falls within any offline weeks of year range (1-52)
	if len(c.altTimeChecks.OfflineWeeksOfYear) > 0 {
		_, isoWeek := localTime.ISOWeek()                                      // ISO 8601 week number (1-52/53)
		weekOfYear := sdkmath.NewUint(uint64(isoWeek))
		detErrMsg := fmt.Sprintf("transfer denied: current UTC week of year %d falls within offline weeks of year", weekOfYear.Uint64())
		found, err := types.SearchUintRangesForUint(weekOfYear, c.altTimeChecks.OfflineWeeksOfYear)
		if err != nil {
			return detErrMsg, sdkerrors.Wrapf(err, "error checking offline weeks of year")
		}
		if found {
			return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
		}
	}

	return "", nil
}
