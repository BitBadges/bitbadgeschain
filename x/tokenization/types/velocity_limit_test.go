package types

import (
	"math"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

// --- IsVelocityLimitBasicallyNil ---

func TestIsVelocityLimitBasicallyNil_Nil(t *testing.T) {
	require.True(t, IsVelocityLimitBasicallyNil(nil))
}

func TestIsVelocityLimitBasicallyNil_ZeroWindow(t *testing.T) {
	vl := &VelocityLimit{WindowDuration: sdkmath.NewUint(0)}
	require.True(t, IsVelocityLimitBasicallyNil(vl))
}

func TestIsVelocityLimitBasicallyNil_NilWindow(t *testing.T) {
	vl := &VelocityLimit{}
	require.True(t, IsVelocityLimitBasicallyNil(vl))
}

func TestIsVelocityLimitBasicallyNil_ValidWindow(t *testing.T) {
	vl := &VelocityLimit{WindowDuration: sdkmath.NewUint(86400000)}
	require.False(t, IsVelocityLimitBasicallyNil(vl))
}

// --- GetBucketSize ---

func TestGetBucketSize_24HourWindow(t *testing.T) {
	// 24 hours = 86,400,000 ms. Bucket size = 86400000 / 24 = 3,600,000 ms (1 hour).
	vl := VelocityLimit{WindowDuration: sdkmath.NewUint(86400000)}
	expected := sdkmath.NewUint(3600000)
	require.True(t, vl.GetBucketSize().Equal(expected),
		"expected bucket size %s, got %s", expected, vl.GetBucketSize())
}

func TestGetBucketSize_MinimumWindow24ms(t *testing.T) {
	// Minimum valid: 24ms => bucket size = 1ms per bucket.
	vl := VelocityLimit{WindowDuration: sdkmath.NewUint(24)}
	expected := sdkmath.NewUint(1)
	require.True(t, vl.GetBucketSize().Equal(expected),
		"expected bucket size %s, got %s", expected, vl.GetBucketSize())
}

func TestGetBucketSize_30DayWindow(t *testing.T) {
	// 30 days = 2,592,000,000 ms. Bucket size = 108,000,000 ms (30 hours).
	windowMs := uint64(30 * 24 * 3600 * 1000) // 2592000000
	vl := VelocityLimit{WindowDuration: sdkmath.NewUint(windowMs)}
	expected := sdkmath.NewUint(windowMs / MaxVelocityBuckets)
	require.True(t, vl.GetBucketSize().Equal(expected),
		"expected bucket size %s, got %s", expected, vl.GetBucketSize())
}

func TestGetBucketSize_NotEvenlyDivisible(t *testing.T) {
	// 25ms / 24 = 1 (integer division truncation)
	vl := VelocityLimit{WindowDuration: sdkmath.NewUint(25)}
	expected := sdkmath.NewUint(1)
	require.True(t, vl.GetBucketSize().Equal(expected))
}

// --- GetBucketIndex ---

func TestGetBucketIndex_VariousTimestamps(t *testing.T) {
	// 24-hour window, bucket size = 3,600,000 ms (1 hour)
	vl := VelocityLimit{WindowDuration: sdkmath.NewUint(86400000)}
	bucketSize := uint64(3600000)

	tests := []struct {
		name      string
		timeMs    uint64
		wantIndex uint64
	}{
		{"time 0", 0, 0},
		{"start of bucket 1", bucketSize, 1},
		{"middle of bucket 0", bucketSize / 2, 0},
		{"start of bucket 23", 23 * bucketSize, 23},
		// Bucket 24 wraps to index 0
		{"wraps to 0 at bucket 24", 24 * bucketSize, 0},
		{"wraps to 1 at bucket 25", 25 * bucketSize, 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			idx := vl.GetBucketIndex(sdkmath.NewUint(tc.timeMs))
			require.Equal(t, tc.wantIndex, idx, "timestamp %d", tc.timeMs)
		})
	}
}

func TestGetBucketIndex_ZeroBucketSize(t *testing.T) {
	// If window < 24 (bucket size = 0), GetBucketIndex returns 0 to avoid div-by-zero.
	vl := VelocityLimit{WindowDuration: sdkmath.NewUint(1)}
	idx := vl.GetBucketIndex(sdkmath.NewUint(1000))
	require.Equal(t, uint64(0), idx)
}

// --- GetBucketStartTime ---

func TestGetBucketStartTime(t *testing.T) {
	vl := VelocityLimit{WindowDuration: sdkmath.NewUint(86400000)}
	bucketSize := uint64(3600000)

	tests := []struct {
		timeMs    uint64
		wantStart uint64
	}{
		{0, 0},
		{1000, 0},                                     // still in first bucket
		{bucketSize, bucketSize},                       // exactly at bucket boundary
		{bucketSize + 500, bucketSize},                 // mid-bucket
		{2*bucketSize - 1, bucketSize},                 // last ms of second bucket
		{23 * bucketSize, 23 * bucketSize},             // start of last bucket
		{23*bucketSize + 100, 23 * bucketSize},         // inside last bucket
		{24 * bucketSize, 24 * bucketSize},             // wraps around (new cycle)
	}

	for _, tc := range tests {
		start := vl.GetBucketStartTime(sdkmath.NewUint(tc.timeMs))
		require.True(t, start.Equal(sdkmath.NewUint(tc.wantStart)),
			"timestamp %d: expected start %d, got %s", tc.timeMs, tc.wantStart, start)
	}
}

func TestGetBucketStartTime_ZeroBucketSize(t *testing.T) {
	vl := VelocityLimit{WindowDuration: sdkmath.NewUint(1)}
	start := vl.GetBucketStartTime(sdkmath.NewUint(5000))
	require.True(t, start.Equal(sdkmath.NewUint(0)))
}

// --- IsBucketExpired ---

func TestIsBucketExpired_WithinWindow(t *testing.T) {
	vl := VelocityLimit{WindowDuration: sdkmath.NewUint(86400000)} // 24 hours
	// Current time = 100,000,000. Window starts at 100,000,000 - 86,400,000 = 13,600,000.
	// Bucket at 14,000,000 is within the window.
	require.False(t, vl.IsBucketExpired(sdkmath.NewUint(14000000), sdkmath.NewUint(100000000)))
}

func TestIsBucketExpired_OutsideWindow(t *testing.T) {
	vl := VelocityLimit{WindowDuration: sdkmath.NewUint(86400000)}
	// Bucket at 10,000,000 is before windowStart of 13,600,000.
	require.True(t, vl.IsBucketExpired(sdkmath.NewUint(10000000), sdkmath.NewUint(100000000)))
}

func TestIsBucketExpired_ExactlyAtWindowEdge(t *testing.T) {
	vl := VelocityLimit{WindowDuration: sdkmath.NewUint(86400000)}
	// Bucket start = currentTime - windowDuration = 13,600,000 is NOT expired (it is at the boundary).
	// windowStart = 100000000 - 86400000 = 13600000. bucketStart < windowStart is false (equal).
	require.False(t, vl.IsBucketExpired(sdkmath.NewUint(13600000), sdkmath.NewUint(100000000)))
}

func TestIsBucketExpired_CurrentTimeLessThanWindow(t *testing.T) {
	vl := VelocityLimit{WindowDuration: sdkmath.NewUint(86400000)}
	// Current time is 1000, which is less than the window duration.
	// Early in blockchain life, nothing should be expired.
	require.False(t, vl.IsBucketExpired(sdkmath.NewUint(0), sdkmath.NewUint(1000)))
}

// --- PruneAndSumBuckets ---

func newBalance(amount uint64) *Balance {
	return &Balance{
		Amount: sdkmath.NewUint(amount),
		OwnershipTimes: []*UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
		},
		TokenIds: []*UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
		},
	}
}

func TestPruneAndSumBuckets_NoBuckets(t *testing.T) {
	vl := VelocityLimit{WindowDuration: sdkmath.NewUint(86400000)}
	buckets := []TimeBucket{} // empty

	amounts, numTransfers, pruned := vl.PruneAndSumBuckets(buckets, sdkmath.NewUint(100000000), false)
	require.Equal(t, 0, len(amounts))
	require.True(t, numTransfers.IsZero())
	require.Equal(t, 0, len(pruned))
}

func TestPruneAndSumBuckets_AllBucketsEmpty(t *testing.T) {
	vl := VelocityLimit{WindowDuration: sdkmath.NewUint(86400000)}
	buckets := InitBuckets() // all start at 0

	amounts, numTransfers, pruned := vl.PruneAndSumBuckets(buckets, sdkmath.NewUint(100000000), false)
	require.Equal(t, 0, len(amounts))
	require.True(t, numTransfers.IsZero())
	require.Equal(t, MaxVelocityBuckets, len(pruned))
}

func TestPruneAndSumBuckets_AllBucketsExpired(t *testing.T) {
	vl := VelocityLimit{WindowDuration: sdkmath.NewUint(86400000)}
	buckets := InitBuckets()

	// Set all buckets to old times (well outside the window)
	for i := range buckets {
		buckets[i].StartTime = sdkmath.NewUint(1000) // very old
		buckets[i].Amounts = []*Balance{newBalance(50)}
		buckets[i].NumTransfers = sdkmath.NewUint(3)
	}

	now := sdkmath.NewUint(200000000) // well past the window
	amounts, numTransfers, pruned := vl.PruneAndSumBuckets(buckets, now, false)

	require.Equal(t, 0, len(amounts), "all expired buckets should yield 0 amounts")
	require.True(t, numTransfers.IsZero())
	// All buckets should be zeroed out
	for _, b := range pruned {
		require.True(t, b.StartTime.IsZero())
		require.Equal(t, 0, len(b.Amounts))
		require.True(t, b.NumTransfers.IsZero())
	}
}

func TestPruneAndSumBuckets_MixedExpiredAndActive_Amounts(t *testing.T) {
	vl := VelocityLimit{WindowDuration: sdkmath.NewUint(86400000)}
	buckets := InitBuckets()

	now := sdkmath.NewUint(200000000)
	windowStart := now.Sub(vl.WindowDuration) // 113,600,000

	// Bucket 0: expired (start = 100,000,000 < 113,600,000)
	buckets[0].StartTime = sdkmath.NewUint(100000000)
	buckets[0].Amounts = []*Balance{newBalance(100)}

	// Bucket 1: active (start = 120,000,000 >= 113,600,000)
	buckets[1].StartTime = sdkmath.NewUint(120000000)
	buckets[1].Amounts = []*Balance{newBalance(50)}

	// Bucket 2: active (start = 150,000,000 >= 113,600,000)
	buckets[2].StartTime = sdkmath.NewUint(150000000)
	buckets[2].Amounts = []*Balance{newBalance(75)}

	amounts, _, pruned := vl.PruneAndSumBuckets(buckets, now, false)

	// Should have 2 active balances
	require.Equal(t, 2, len(amounts))
	// Bucket 0 should be pruned
	require.True(t, pruned[0].StartTime.IsZero())
	// Bucket 1 and 2 should remain
	require.True(t, pruned[1].StartTime.Equal(sdkmath.NewUint(120000000)))
	require.True(t, pruned[2].StartTime.Equal(sdkmath.NewUint(150000000)))

	// Verify the expired bucket is actually at 100000000 which is < windowStart
	require.True(t, sdkmath.NewUint(100000000).LT(windowStart))
}

func TestPruneAndSumBuckets_NumTransfersMode(t *testing.T) {
	vl := VelocityLimit{WindowDuration: sdkmath.NewUint(86400000)}
	buckets := InitBuckets()

	now := sdkmath.NewUint(200000000)

	// Active bucket with 5 transfers
	buckets[0].StartTime = sdkmath.NewUint(180000000)
	buckets[0].NumTransfers = sdkmath.NewUint(5)

	// Active bucket with 3 transfers
	buckets[1].StartTime = sdkmath.NewUint(190000000)
	buckets[1].NumTransfers = sdkmath.NewUint(3)

	// Expired bucket with 10 transfers (should not count)
	buckets[2].StartTime = sdkmath.NewUint(50000000)
	buckets[2].NumTransfers = sdkmath.NewUint(10)

	_, numTransfers, _ := vl.PruneAndSumBuckets(buckets, now, true)
	require.True(t, numTransfers.Equal(sdkmath.NewUint(8)),
		"expected 8 transfers, got %s", numTransfers)
}

func TestPruneAndSumBuckets_RollingSumIsZeroWhenAllExpired(t *testing.T) {
	vl := VelocityLimit{WindowDuration: sdkmath.NewUint(86400000)}
	buckets := InitBuckets()

	// All buckets are very old
	for i := range buckets {
		buckets[i].StartTime = sdkmath.NewUint(1000)
		buckets[i].NumTransfers = sdkmath.NewUint(99)
	}

	now := sdkmath.NewUint(500000000)
	_, numTransfers, _ := vl.PruneAndSumBuckets(buckets, now, true)
	require.True(t, numTransfers.IsZero(), "all expired => rolling sum should be 0")
}

// --- InitBuckets ---

func TestInitBuckets_Length(t *testing.T) {
	buckets := InitBuckets()
	require.Equal(t, MaxVelocityBuckets, len(buckets))
	for _, b := range buckets {
		require.True(t, b.StartTime.IsZero())
		require.Equal(t, 0, len(b.Amounts))
		require.True(t, b.NumTransfers.IsZero())
	}
}

// --- Edge cases ---

func TestGetBucketIndex_LargeTimestamp(t *testing.T) {
	// 30-day window with very large timestamp (simulating real-world epoch ms).
	windowMs := uint64(30 * 24 * 3600 * 1000)
	vl := VelocityLimit{WindowDuration: sdkmath.NewUint(windowMs)}

	// Approximate current epoch in ms (Mar 2026): ~1,774,000,000,000
	epochMs := uint64(1774000000000)
	idx := vl.GetBucketIndex(sdkmath.NewUint(epochMs))
	require.Less(t, idx, uint64(MaxVelocityBuckets), "bucket index must be in [0, 24)")
}

func TestPruneAndSumBuckets_DoesNotMutateOriginal(t *testing.T) {
	vl := VelocityLimit{WindowDuration: sdkmath.NewUint(86400000)}
	buckets := InitBuckets()

	// Set an expired bucket
	buckets[0].StartTime = sdkmath.NewUint(1000)
	buckets[0].NumTransfers = sdkmath.NewUint(5)
	originalStart := buckets[0].StartTime

	_, _, _ = vl.PruneAndSumBuckets(buckets, sdkmath.NewUint(500000000), true)

	// The original slice should not be modified (copy semantics)
	require.True(t, buckets[0].StartTime.Equal(originalStart),
		"PruneAndSumBuckets should not mutate the original bucket slice")
}

func TestGetBucketIndex_DeterministicAcrossCycles(t *testing.T) {
	// The same bucket index repeats every 24 buckets (modular).
	vl := VelocityLimit{WindowDuration: sdkmath.NewUint(2400)} // bucket size = 100
	idx1 := vl.GetBucketIndex(sdkmath.NewUint(500))            // 500/100 = 5 mod 24 = 5
	idx2 := vl.GetBucketIndex(sdkmath.NewUint(2900))           // 2900/100 = 29 mod 24 = 5
	require.Equal(t, idx1, idx2, "same position in different cycles should map to same index")
	require.Equal(t, uint64(5), idx1)
}
