package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkmath "cosmossdk.io/math"
)

// ApplyVelocityLimit handles the sliding-window velocity tracking for a single transfer.
// It is called instead of ResetApprovalTrackerIfNeeded when a VelocityLimit is configured.
//
// The approach:
//  1. Initialize time buckets if they don't exist yet
//  2. Determine which bucket the current time falls into
//  3. Prune expired buckets (lazy pruning — only on write)
//  4. Add the transfer data to the current bucket
//  5. Sum all non-expired buckets for the rolling total
//  6. Return the rolling totals so the caller can check thresholds
//
// Parameters:
//   - ctx: the SDK context (provides block time)
//   - tracker: the current ApprovalTracker state
//   - velocityLimit: the configured velocity limit with window duration
//   - transferAmounts: the balances being transferred (nil if isNumTransfers)
//   - isNumTransfers: if true, we're tracking transfer count; if false, tracking amounts
//
// Returns:
//   - updatedTracker: the tracker with updated time buckets
//   - rollingAmounts: summed balances across the rolling window (only if !isNumTransfers)
//   - rollingNumTransfers: summed transfer count across the rolling window (only if isNumTransfers)
func (k Keeper) ApplyVelocityLimit(
	ctx sdk.Context,
	tracker *types.ApprovalTracker,
	velocityLimit *types.VelocityLimit,
	transferAmounts []*types.Balance,
	isNumTransfers bool,
) (
	updatedTracker types.ApprovalTracker,
	rollingAmounts []*types.Balance,
	rollingNumTransfers sdkmath.Uint,
) {
	now := sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli()))

	// Initialize time buckets if this tracker hasn't been used with velocity limits yet
	if len(tracker.TimeBuckets) == 0 {
		tracker.TimeBuckets = types.InitBuckets()
	}

	// Determine which bucket the current time falls into.
	// bucketIndex = (currentTimeMs / bucketSizeMs) % 24
	bucketIdx := velocityLimit.GetBucketIndex(now)
	bucketStart := velocityLimit.GetBucketStartTime(now)

	// Check if the current bucket needs to be reset (it belongs to an older cycle).
	// A bucket at index i might hold data from a previous window cycle. If its startTime
	// doesn't match the expected start time for this cycle, it's stale and must be reset.
	currentBucket := &tracker.TimeBuckets[bucketIdx]
	if !currentBucket.StartTime.Equal(bucketStart) {
		// This bucket is from a different time period — reset it
		currentBucket.StartTime = bucketStart
		currentBucket.Amounts = []*types.Balance{}
		currentBucket.NumTransfers = sdkmath.NewUint(0)
	}

	// Add the current transfer to the bucket
	if isNumTransfers {
		currentBucket.NumTransfers = currentBucket.NumTransfers.Add(sdkmath.NewUint(1))
	} else if transferAmounts != nil {
		currentBucket.Amounts = append(currentBucket.Amounts, transferAmounts...)
	}

	// Prune expired buckets and compute rolling totals.
	// This is lazy pruning: we only clear old buckets when we access the tracker.
	rollingAmounts, rollingNumTransfers, prunedBuckets := velocityLimit.PruneAndSumBuckets(
		tracker.TimeBuckets, now, isNumTransfers,
	)
	tracker.TimeBuckets = prunedBuckets

	// Update the last-updated timestamp
	tracker.LastUpdatedAt = now

	return *tracker, rollingAmounts, rollingNumTransfers
}

// GetApprovalTrackerWithVelocity fetches the tracker from the store and applies
// velocity-limit logic if a VelocityLimit is configured. This is the velocity-aware
// counterpart to GetApprovalTrackerFromStoreAndResetIfNeeded.
//
// When velocityLimit is set, the tracker is NOT reset like with ResetTimeIntervals.
// Instead, the rolling window is computed by summing non-expired time buckets.
// The returned tracker has its Amounts/NumTransfers set to the rolling totals
// so existing threshold-checking code works without modification.
func (k Keeper) GetApprovalTrackerWithVelocity(
	ctx sdk.Context,
	collectionId sdkmath.Uint,
	addressForApproval string,
	approvalId string,
	amountTrackerId string,
	level string,
	trackerType string,
	address string,
	velocityLimit *types.VelocityLimit,
	isNumTransfers bool,
) (types.ApprovalTracker, error) {
	tracker, found := k.GetApprovalTrackerFromStore(
		ctx, collectionId, addressForApproval, approvalId,
		amountTrackerId, level, trackerType, address,
	)
	if !found {
		tracker = types.ApprovalTracker{
			Amounts:       []*types.Balance{},
			NumTransfers:  sdkmath.NewUint(0),
			LastUpdatedAt: sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli())),
			TimeBuckets:   types.InitBuckets(),
		}
	}

	now := sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli()))

	// Initialize time buckets if needed
	if len(tracker.TimeBuckets) == 0 {
		tracker.TimeBuckets = types.InitBuckets()
	}

	// Prune expired buckets and compute rolling totals WITHOUT adding any new transfer.
	// This gives us the current rolling window state for threshold checking.
	rollingAmounts, rollingNumTransfers, prunedBuckets := velocityLimit.PruneAndSumBuckets(
		tracker.TimeBuckets, now, isNumTransfers,
	)
	tracker.TimeBuckets = prunedBuckets

	// Set the tracker's Amounts/NumTransfers to the rolling totals so existing
	// threshold-checking code works transparently.
	if isNumTransfers {
		tracker.NumTransfers = rollingNumTransfers
	} else {
		tracker.Amounts = rollingAmounts
	}

	tracker.LastUpdatedAt = now
	return tracker, nil
}

// IncrementVelocityTracker adds a transfer to the appropriate time bucket and
// returns the updated tracker. Called after threshold checks pass, to record
// the transfer in the sliding window.
//
// This is the velocity-aware counterpart to the simple increment done in
// IncrementApprovalsAndAssertWithinThreshold.
func (k Keeper) IncrementVelocityTracker(
	ctx sdk.Context,
	tracker *types.ApprovalTracker,
	velocityLimit *types.VelocityLimit,
	transferAmounts []*types.Balance,
	isNumTransfers bool,
) types.ApprovalTracker {
	now := sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli()))

	// Initialize time buckets if needed
	if len(tracker.TimeBuckets) == 0 {
		tracker.TimeBuckets = types.InitBuckets()
	}

	// Determine the bucket for the current time
	bucketIdx := velocityLimit.GetBucketIndex(now)
	bucketStart := velocityLimit.GetBucketStartTime(now)

	currentBucket := &tracker.TimeBuckets[bucketIdx]
	if !currentBucket.StartTime.Equal(bucketStart) {
		// Stale bucket from a different time period — reset it
		currentBucket.StartTime = bucketStart
		currentBucket.Amounts = []*types.Balance{}
		currentBucket.NumTransfers = sdkmath.NewUint(0)
	}

	// Record the transfer in the current bucket
	if isNumTransfers {
		currentBucket.NumTransfers = currentBucket.NumTransfers.Add(sdkmath.NewUint(1))
	} else if transferAmounts != nil {
		currentBucket.Amounts = append(currentBucket.Amounts, transferAmounts...)
	}

	tracker.LastUpdatedAt = now
	return *tracker
}
