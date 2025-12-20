package keeper

import (
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
)

// CheckAltTimeChecks validates if the current UTC time falls within any offline hours or days specified in altTimeChecks.
// Returns an error if the time is within an offline period (i.e., the approval should be denied).
// CheckAltTimeChecks validates alternative time checks
// Returns (deterministicErrorMsg, error) where deterministicErrorMsg is a deterministic error string
func (k Keeper) CheckAltTimeChecks(ctx sdk.Context, altTimeChecks *types.AltTimeChecks) (string, error) {
	if altTimeChecks == nil {
		return "", nil
	}

	// Get UTC time from block time
	blockTime := ctx.BlockTime()
	utcTime := blockTime.UTC()

	// Get current hour (0-23) and day of week (0=Sunday, 1=Monday, ..., 6=Saturday)
	currentHour := sdkmath.NewUint(uint64(utcTime.Hour()))
	currentDay := sdkmath.NewUint(uint64(utcTime.Weekday()))

	// Check if current hour falls within any offline hours range
	if len(altTimeChecks.OfflineHours) > 0 {
		detErrMsg := fmt.Sprintf("transfer denied: current UTC hour %d falls within offline hours", currentHour.Uint64())
		found, err := types.SearchUintRangesForUint(currentHour, altTimeChecks.OfflineHours)
		if err != nil {
			return detErrMsg, sdkerrors.Wrapf(err, "error checking offline hours")
		}
		if found {
			return detErrMsg, sdkerrors.Wrap(ErrDisallowedTransfer, detErrMsg)
		}
	}

	// Check if current day falls within any offline days range
	if len(altTimeChecks.OfflineDays) > 0 {
		detErrMsg := fmt.Sprintf("transfer denied: current UTC day %d falls within offline days", currentDay.Uint64())
		found, err := types.SearchUintRangesForUint(currentDay, altTimeChecks.OfflineDays)
		if err != nil {
			return detErrMsg, sdkerrors.Wrapf(err, "error checking offline days")
		}
		if found {
			return detErrMsg, sdkerrors.Wrap(ErrDisallowedTransfer, detErrMsg)
		}
	}

	return "", nil
}
