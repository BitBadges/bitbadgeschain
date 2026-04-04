package keeper_test

import (
	"math"
	"time"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkmath "cosmossdk.io/math"
)

// Helper to create a simple balance for testing.
func newTestBalance(amount uint64) []*types.Balance {
	return []*types.Balance{
		{
			Amount: sdkmath.NewUint(amount),
			OwnershipTimes: []*types.UintRange{
				{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
			},
			TokenIds: []*types.UintRange{
				{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
			},
		},
	}
}

// setBlockTime sets the suite context block time from a millisecond timestamp.
func (suite *TestSuite) setBlockTimeMs(ms uint64) {
	suite.ctx = suite.ctx.WithBlockTime(time.UnixMilli(int64(ms)))
}

// --- IncrementVelocityTracker ---

func (suite *TestSuite) TestIncrementVelocityTracker_AddsToCorrectBucket() {
	// 24-hour window, bucket size = 3,600,000 ms (1 hour)
	vl := &types.VelocityLimit{WindowDuration: sdkmath.NewUint(86400000)}
	tracker := &types.ApprovalTracker{
		Amounts:       []*types.Balance{},
		NumTransfers:  sdkmath.NewUint(0),
		LastUpdatedAt: sdkmath.NewUint(0),
		TimeBuckets:   types.InitBuckets(),
	}

	// Time = 5,000,000 ms => bucket index = (5000000/3600000) % 24 = 1
	suite.setBlockTimeMs(5000000)

	updated := suite.app.TokenizationKeeper.IncrementVelocityTracker(
		suite.ctx, tracker, vl, newTestBalance(100), false,
	)

	expectedIdx := vl.GetBucketIndex(sdkmath.NewUint(5000000))
	suite.Require().Equal(uint64(1), expectedIdx)
	suite.Require().Equal(1, len(updated.TimeBuckets[expectedIdx].Amounts))
	suite.Require().True(updated.TimeBuckets[expectedIdx].Amounts[0].Amount.Equal(sdkmath.NewUint(100)))
}

func (suite *TestSuite) TestIncrementVelocityTracker_AccumulatesInSameBucket() {
	vl := &types.VelocityLimit{WindowDuration: sdkmath.NewUint(86400000)}
	tracker := &types.ApprovalTracker{
		Amounts:       []*types.Balance{},
		NumTransfers:  sdkmath.NewUint(0),
		LastUpdatedAt: sdkmath.NewUint(0),
		TimeBuckets:   types.InitBuckets(),
	}

	// Two transfers at the same time => same bucket
	suite.setBlockTimeMs(5000000)

	updated := suite.app.TokenizationKeeper.IncrementVelocityTracker(
		suite.ctx, tracker, vl, newTestBalance(100), false,
	)

	updated2 := suite.app.TokenizationKeeper.IncrementVelocityTracker(
		suite.ctx, &updated, vl, newTestBalance(50), false,
	)

	expectedIdx := vl.GetBucketIndex(sdkmath.NewUint(5000000))
	suite.Require().Equal(2, len(updated2.TimeBuckets[expectedIdx].Amounts),
		"two transfers in same bucket should yield 2 amount entries")
}

func (suite *TestSuite) TestIncrementVelocityTracker_DifferentBucketsForDifferentTimes() {
	vl := &types.VelocityLimit{WindowDuration: sdkmath.NewUint(86400000)}
	bucketSize := uint64(3600000)
	tracker := &types.ApprovalTracker{
		Amounts:       []*types.Balance{},
		NumTransfers:  sdkmath.NewUint(0),
		LastUpdatedAt: sdkmath.NewUint(0),
		TimeBuckets:   types.InitBuckets(),
	}

	// Transfer at time 0 => bucket 0
	suite.setBlockTimeMs(0)
	updated := suite.app.TokenizationKeeper.IncrementVelocityTracker(
		suite.ctx, tracker, vl, newTestBalance(100), false,
	)

	// Transfer at time = bucketSize => bucket 1
	suite.setBlockTimeMs(bucketSize)
	updated = suite.app.TokenizationKeeper.IncrementVelocityTracker(
		suite.ctx, &updated, vl, newTestBalance(200), false,
	)

	idx0 := vl.GetBucketIndex(sdkmath.NewUint(0))
	idx1 := vl.GetBucketIndex(sdkmath.NewUint(bucketSize))
	suite.Require().NotEqual(idx0, idx1, "different times should yield different buckets")
	suite.Require().Equal(1, len(updated.TimeBuckets[idx0].Amounts))
	suite.Require().Equal(1, len(updated.TimeBuckets[idx1].Amounts))
}

func (suite *TestSuite) TestIncrementVelocityTracker_NumTransfersMode() {
	vl := &types.VelocityLimit{WindowDuration: sdkmath.NewUint(86400000)}
	tracker := &types.ApprovalTracker{
		Amounts:       []*types.Balance{},
		NumTransfers:  sdkmath.NewUint(0),
		LastUpdatedAt: sdkmath.NewUint(0),
		TimeBuckets:   types.InitBuckets(),
	}

	suite.setBlockTimeMs(5000000)

	updated := suite.app.TokenizationKeeper.IncrementVelocityTracker(
		suite.ctx, tracker, vl, nil, true,
	)
	updated = suite.app.TokenizationKeeper.IncrementVelocityTracker(
		suite.ctx, &updated, vl, nil, true,
	)

	expectedIdx := vl.GetBucketIndex(sdkmath.NewUint(5000000))
	suite.Require().True(updated.TimeBuckets[expectedIdx].NumTransfers.Equal(sdkmath.NewUint(2)),
		"two increments should yield numTransfers=2")
}

func (suite *TestSuite) TestIncrementVelocityTracker_ResetsStaleOldBucket() {
	// Verify that when a bucket index wraps around and the old bucket has stale data,
	// it gets reset before writing new data.
	vl := &types.VelocityLimit{WindowDuration: sdkmath.NewUint(86400000)}
	bucketSize := uint64(3600000)
	tracker := &types.ApprovalTracker{
		Amounts:       []*types.Balance{},
		NumTransfers:  sdkmath.NewUint(0),
		LastUpdatedAt: sdkmath.NewUint(0),
		TimeBuckets:   types.InitBuckets(),
	}

	// Write to bucket at time 0 (bucket index 0)
	suite.setBlockTimeMs(0)
	updated := suite.app.TokenizationKeeper.IncrementVelocityTracker(
		suite.ctx, tracker, vl, newTestBalance(999), false,
	)

	// Jump forward 24 full buckets (one full cycle) so bucket index 0 wraps around
	newTime := 24 * bucketSize
	suite.setBlockTimeMs(newTime)
	updated = suite.app.TokenizationKeeper.IncrementVelocityTracker(
		suite.ctx, &updated, vl, newTestBalance(50), false,
	)

	idx := vl.GetBucketIndex(sdkmath.NewUint(newTime))
	suite.Require().Equal(uint64(0), idx, "should wrap to index 0")
	// The old 999 should be gone; only 50 remains
	suite.Require().Equal(1, len(updated.TimeBuckets[idx].Amounts))
	suite.Require().True(updated.TimeBuckets[idx].Amounts[0].Amount.Equal(sdkmath.NewUint(50)))
}

// --- GetApprovalTrackerWithVelocity ---

func (suite *TestSuite) TestGetApprovalTrackerWithVelocity_RollingTotals() {
	// Test that GetApprovalTrackerWithVelocity correctly computes rolling totals
	// by building up tracker state through IncrementVelocityTracker (in-memory),
	// then verifying PruneAndSumBuckets produces the right totals.
	//
	// Note: TimeBucket is not a protobuf message and does not survive store
	// round-trips. The velocity system works by maintaining TimeBuckets in memory
	// within a single transaction's execution path.
	vl := &types.VelocityLimit{WindowDuration: sdkmath.NewUint(86400000)}
	bucketSize := uint64(3600000)

	tracker := types.ApprovalTracker{
		Amounts:       []*types.Balance{},
		NumTransfers:  sdkmath.NewUint(0),
		LastUpdatedAt: sdkmath.NewUint(0),
		TimeBuckets:   types.InitBuckets(),
	}

	baseTime := uint64(200000000)

	// 5 transfers in hour 0
	suite.setBlockTimeMs(baseTime + 100)
	for i := 0; i < 5; i++ {
		tracker = suite.app.TokenizationKeeper.IncrementVelocityTracker(
			suite.ctx, &tracker, vl, nil, true,
		)
	}

	// 3 transfers in hour 1
	suite.setBlockTimeMs(baseTime + bucketSize + 100)
	for i := 0; i < 3; i++ {
		tracker = suite.app.TokenizationKeeper.IncrementVelocityTracker(
			suite.ctx, &tracker, vl, nil, true,
		)
	}

	// Sum the rolling window at the current time
	now := sdkmath.NewUint(baseTime + bucketSize + 100)
	_, rollingNum, _ := vl.PruneAndSumBuckets(tracker.TimeBuckets, now, true)

	suite.Require().True(rollingNum.Equal(sdkmath.NewUint(8)),
		"rolling total should be 5+3=8, got %s", rollingNum)
}

func (suite *TestSuite) TestGetApprovalTrackerWithVelocity_ExcludesExpiredBuckets() {
	// Verify that expired buckets are excluded from rolling totals.
	vl := &types.VelocityLimit{WindowDuration: sdkmath.NewUint(86400000)}

	tracker := types.ApprovalTracker{
		Amounts:       []*types.Balance{},
		NumTransfers:  sdkmath.NewUint(0),
		LastUpdatedAt: sdkmath.NewUint(0),
		TimeBuckets:   types.InitBuckets(),
	}

	// Transfer at baseTime
	baseTime := uint64(200000000)
	suite.setBlockTimeMs(baseTime)
	tracker = suite.app.TokenizationKeeper.IncrementVelocityTracker(
		suite.ctx, &tracker, vl, nil, true,
	)

	// Transfer at baseTime + 2 hours (different bucket, still within window)
	suite.setBlockTimeMs(baseTime + 7200000)
	tracker = suite.app.TokenizationKeeper.IncrementVelocityTracker(
		suite.ctx, &tracker, vl, nil, true,
	)

	// Advance time well past the window to expire both buckets
	expiredTime := baseTime + uint64(86400000) + uint64(7200001)
	_, rollingNum, _ := vl.PruneAndSumBuckets(tracker.TimeBuckets, sdkmath.NewUint(expiredTime), true)

	suite.Require().True(rollingNum.IsZero(),
		"all buckets should be expired, got %s", rollingNum)
}

func (suite *TestSuite) TestGetApprovalTrackerWithVelocity_NotFoundInitializesEmpty() {
	vl := &types.VelocityLimit{WindowDuration: sdkmath.NewUint(86400000)}

	suite.setBlockTimeMs(100000000)

	// Tracker doesn't exist in store yet
	result, err := suite.app.TokenizationKeeper.GetApprovalTrackerWithVelocity(
		suite.ctx, sdkmath.NewUint(999), "addr", "approval-nonexistent", "tracker-nonexistent",
		"collection", "overall", "", vl, true,
	)
	suite.Require().Nil(err)
	suite.Require().True(result.NumTransfers.IsZero(), "new tracker should have 0 transfers")
	suite.Require().Equal(types.MaxVelocityBuckets, len(result.TimeBuckets))
}

// --- ApplyVelocityLimit ---

func (suite *TestSuite) TestApplyVelocityLimit_PrunesAndIncrements() {
	vl := &types.VelocityLimit{WindowDuration: sdkmath.NewUint(86400000)}
	tracker := &types.ApprovalTracker{
		Amounts:       []*types.Balance{},
		NumTransfers:  sdkmath.NewUint(0),
		LastUpdatedAt: sdkmath.NewUint(0),
		TimeBuckets:   types.InitBuckets(),
	}

	// Add an expired bucket
	tracker.TimeBuckets[2].StartTime = sdkmath.NewUint(1000)
	tracker.TimeBuckets[2].NumTransfers = sdkmath.NewUint(99)

	suite.setBlockTimeMs(200000000)

	updatedTracker, _, rollingNumTransfers := suite.app.TokenizationKeeper.ApplyVelocityLimit(
		suite.ctx, tracker, vl, nil, true,
	)

	// The expired bucket should not count. Only the new transfer (1) should.
	suite.Require().True(rollingNumTransfers.Equal(sdkmath.NewUint(1)),
		"only the new transfer should be in the rolling sum, got %s", rollingNumTransfers)
	// Expired bucket should be pruned
	suite.Require().True(updatedTracker.TimeBuckets[2].StartTime.IsZero())
}

func (suite *TestSuite) TestApplyVelocityLimit_InitializesBucketsIfEmpty() {
	vl := &types.VelocityLimit{WindowDuration: sdkmath.NewUint(86400000)}
	tracker := &types.ApprovalTracker{
		Amounts:       []*types.Balance{},
		NumTransfers:  sdkmath.NewUint(0),
		LastUpdatedAt: sdkmath.NewUint(0),
		TimeBuckets:   []types.TimeBucket{}, // empty!
	}

	suite.setBlockTimeMs(100000000)

	updatedTracker, _, _ := suite.app.TokenizationKeeper.ApplyVelocityLimit(
		suite.ctx, tracker, vl, nil, true,
	)

	suite.Require().Equal(types.MaxVelocityBuckets, len(updatedTracker.TimeBuckets))
}

// --- Boundary gaming prevention ---
// This is the primary value proposition of sliding windows over fixed resets.

func (suite *TestSuite) TestBoundaryGamingPrevention_SlidingWindowBlocksReset() {
	// Scenario: "100 per 24 hours" limit.
	//
	// With fixed resets (resets at midnight), a user can:
	//   - Transfer 100 at 23:59 (just before reset)
	//   - Transfer 100 at 00:01 (just after reset)
	//   => 200 in 2 minutes, "legally" under the per-period limit.
	//
	// With sliding window (24hr), the second transfer should fail because
	// the window still sees the first 100 within the last 24 hours.

	windowDuration := uint64(86400000) // 24 hours in ms
	vl := &types.VelocityLimit{WindowDuration: sdkmath.NewUint(windowDuration)}

	tracker := &types.ApprovalTracker{
		Amounts:       []*types.Balance{},
		NumTransfers:  sdkmath.NewUint(0),
		LastUpdatedAt: sdkmath.NewUint(0),
		TimeBuckets:   types.InitBuckets(),
	}

	// Transfer 1: at time T = 86,399,000 (23h 59m 59s into the day)
	t1 := uint64(86399000)
	suite.setBlockTimeMs(t1)
	updated, rollingAmounts1, _ := suite.app.TokenizationKeeper.ApplyVelocityLimit(
		suite.ctx, tracker, vl, newTestBalance(100), false,
	)

	// Rolling sum after first transfer should include the 100.
	suite.Require().Equal(1, len(rollingAmounts1), "first transfer: should have 1 balance entry")

	// Transfer 2: at time T = 86,401,000 (2 seconds later, "past midnight" in a fixed-reset world)
	t2 := uint64(86401000)
	suite.setBlockTimeMs(t2)
	updated2, rollingAmounts2, _ := suite.app.TokenizationKeeper.ApplyVelocityLimit(
		suite.ctx, &updated, vl, newTestBalance(100), false,
	)

	// The sliding window from t2 looks back to t2 - 86400000 = 1000 ms.
	// The first transfer at 86399000 is well within that window (86399000 > 1000).
	// So the rolling sum should include BOTH transfers.
	_ = updated2
	totalAmountEntries := len(rollingAmounts2)
	suite.Require().Equal(2, totalAmountEntries,
		"sliding window should see BOTH transfers (boundary gaming prevention): got %d entries", totalAmountEntries)

	// If this were a fixed-reset system, the first transfer would have been "forgotten" at the
	// reset boundary, allowing the second 100 to pass. The sliding window prevents this.
}

func (suite *TestSuite) TestBoundaryGamingPrevention_NumTransfers() {
	// Same scenario but with transfer count tracking.
	// Limit: 1 transfer per 24 hours.
	// Transfer at 23:59, another at 00:01 => sliding window should count both.

	windowDuration := uint64(86400000)
	vl := &types.VelocityLimit{WindowDuration: sdkmath.NewUint(windowDuration)}

	tracker := &types.ApprovalTracker{
		Amounts:       []*types.Balance{},
		NumTransfers:  sdkmath.NewUint(0),
		LastUpdatedAt: sdkmath.NewUint(0),
		TimeBuckets:   types.InitBuckets(),
	}

	// Transfer 1 at 23:59
	suite.setBlockTimeMs(86399000)
	updated, _, rollingNum1 := suite.app.TokenizationKeeper.ApplyVelocityLimit(
		suite.ctx, tracker, vl, nil, true,
	)
	suite.Require().True(rollingNum1.Equal(sdkmath.NewUint(1)),
		"first transfer: rolling count should be 1, got %s", rollingNum1)

	// Transfer 2 at 00:01 (2s later)
	suite.setBlockTimeMs(86401000)
	_, _, rollingNum2 := suite.app.TokenizationKeeper.ApplyVelocityLimit(
		suite.ctx, &updated, vl, nil, true,
	)
	suite.Require().True(rollingNum2.Equal(sdkmath.NewUint(2)),
		"second transfer: rolling count should be 2 (boundary gaming blocked), got %s", rollingNum2)
}

func (suite *TestSuite) TestBoundaryGaming_TransferExpiresAfterFullWindow() {
	// Verify that a transfer DOES expire once the full window has passed.
	// Use epoch-like timestamps so bucket start times are never zero (zero is
	// treated as "never written" by PruneAndSumBuckets).
	windowDuration := uint64(86400000)
	vl := &types.VelocityLimit{WindowDuration: sdkmath.NewUint(windowDuration)}

	tracker := &types.ApprovalTracker{
		Amounts:       []*types.Balance{},
		NumTransfers:  sdkmath.NewUint(0),
		LastUpdatedAt: sdkmath.NewUint(0),
		TimeBuckets:   types.InitBuckets(),
	}

	// Transfer at a realistic epoch time (well past 0)
	baseTime := uint64(200000000)
	suite.setBlockTimeMs(baseTime)
	updated, _, rollingNum1 := suite.app.TokenizationKeeper.ApplyVelocityLimit(
		suite.ctx, tracker, vl, nil, true,
	)
	suite.Require().True(rollingNum1.Equal(sdkmath.NewUint(1)),
		"first transfer: expected rolling count 1, got %s", rollingNum1)

	// Jump forward past the window + 1 full bucket so the original bucket is fully expired.
	expiredTime := baseTime + windowDuration + uint64(3600001)
	suite.setBlockTimeMs(expiredTime)

	_, _, rollingNum2 := suite.app.TokenizationKeeper.ApplyVelocityLimit(
		suite.ctx, &updated, vl, nil, true,
	)
	// Only the new transfer at expiredTime should count; the old one is expired.
	suite.Require().True(rollingNum2.Equal(sdkmath.NewUint(1)),
		"old transfer should have expired, expected 1 (new transfer only), got %s", rollingNum2)
}

func (suite *TestSuite) TestMultipleTransfersAcrossMultipleBuckets() {
	// Simulate 5 transfers spread across different hours within the window.
	windowDuration := uint64(86400000)
	vl := &types.VelocityLimit{WindowDuration: sdkmath.NewUint(windowDuration)}
	bucketSize := windowDuration / types.MaxVelocityBuckets // 3,600,000

	tracker := types.ApprovalTracker{
		Amounts:       []*types.Balance{},
		NumTransfers:  sdkmath.NewUint(0),
		LastUpdatedAt: sdkmath.NewUint(0),
		TimeBuckets:   types.InitBuckets(),
	}

	// Base time: start of a window cycle
	base := uint64(100000000)

	// 5 transfers each in a different hour-bucket
	for i := uint64(0); i < 5; i++ {
		t := base + (i * bucketSize) + 100 // +100 to be inside the bucket, not at boundary
		suite.setBlockTimeMs(t)
		tracker = suite.app.TokenizationKeeper.IncrementVelocityTracker(
			suite.ctx, &tracker, vl, nil, true,
		)
	}

	// Now query the rolling total at the time of the last transfer
	lastTime := base + (4*bucketSize) + 100
	suite.setBlockTimeMs(lastTime)

	// Prune and sum directly
	_, numTransfers, _ := vl.PruneAndSumBuckets(tracker.TimeBuckets, sdkmath.NewUint(lastTime), true)
	suite.Require().True(numTransfers.Equal(sdkmath.NewUint(5)),
		"all 5 transfers should be within the window, got %s", numTransfers)
}
