package types

import (
	sdkmath "cosmossdk.io/math"
)

// MaxVelocityBuckets is the fixed number of time buckets used for sliding-window
// velocity tracking. Using a fixed bucket count bounds storage per tracker
// regardless of window size while providing ±(windowDuration/24) precision.
const MaxVelocityBuckets = 24

// VelocityLimit defines a sliding-window rate limit for approvals.
// Unlike ResetTimeIntervals which resets counters at fixed boundaries (and can
// be gamed by timing transfers at boundary edges), VelocityLimit enforces a
// rolling window using time-bucketed tracking.
//
// WindowDuration is the length of the rolling window in milliseconds.
// The window is divided into MaxVelocityBuckets (24) buckets, so each bucket
// covers windowDuration/24 milliseconds. Precision is ±(windowDuration/24).
type VelocityLimit struct {
	// The rolling window duration in milliseconds.
	WindowDuration sdkmath.Uint `json:"windowDuration"`
}

// TimeBucket stores the transfer data for a single time bucket within the
// sliding window. Each bucket covers a duration of windowDuration/24.
type TimeBucket struct {
	// The absolute start time of this bucket in milliseconds.
	// Used to determine if the bucket has expired (older than currentTime - windowDuration).
	StartTime sdkmath.Uint `json:"startTime"`
	// Cumulative transfer amounts recorded in this bucket.
	Amounts []*Balance `json:"amounts,omitempty"`
	// Number of transfers recorded in this bucket.
	NumTransfers sdkmath.Uint `json:"numTransfers"`
}

// IsVelocityLimitBasicallyNil returns true if the velocity limit is nil or has
// a zero/nil window duration, meaning it should be treated as not set.
func IsVelocityLimitBasicallyNil(vl *VelocityLimit) bool {
	return vl == nil || vl.WindowDuration.IsNil() || vl.WindowDuration.IsZero()
}

// GetBucketSize returns the duration each bucket covers: windowDuration / MaxVelocityBuckets.
// Panics if windowDuration is zero (caller must validate first).
func (vl *VelocityLimit) GetBucketSize() sdkmath.Uint {
	return vl.WindowDuration.Quo(sdkmath.NewUint(MaxVelocityBuckets))
}

// GetBucketIndex returns the deterministic bucket index for the given timestamp.
// bucketIndex = (currentTimeMs / bucketSizeMs) % MaxVelocityBuckets
func (vl *VelocityLimit) GetBucketIndex(currentTimeMs sdkmath.Uint) uint64 {
	bucketSize := vl.GetBucketSize()
	if bucketSize.IsZero() {
		return 0
	}
	// globalBucketNum = currentTimeMs / bucketSizeMs
	globalBucketNum := currentTimeMs.Quo(bucketSize)
	// index = globalBucketNum % 24
	return globalBucketNum.Mod(sdkmath.NewUint(MaxVelocityBuckets)).Uint64()
}

// GetBucketStartTime returns the start timestamp of the bucket that contains currentTimeMs.
// bucketStart = (currentTimeMs / bucketSizeMs) * bucketSizeMs
func (vl *VelocityLimit) GetBucketStartTime(currentTimeMs sdkmath.Uint) sdkmath.Uint {
	bucketSize := vl.GetBucketSize()
	if bucketSize.IsZero() {
		return sdkmath.NewUint(0)
	}
	return currentTimeMs.Quo(bucketSize).Mul(bucketSize)
}

// IsBucketExpired returns true if the bucket's start time is older than
// (currentTimeMs - windowDuration), meaning it is outside the rolling window.
func (vl *VelocityLimit) IsBucketExpired(bucketStartTime, currentTimeMs sdkmath.Uint) bool {
	if currentTimeMs.LT(vl.WindowDuration) {
		// If current time is less than the window, only bucket start time 0 would be valid
		return false
	}
	windowStart := currentTimeMs.Sub(vl.WindowDuration)
	return bucketStartTime.LT(windowStart)
}

// PruneAndSumBuckets iterates over all time buckets, prunes expired ones (by zeroing them),
// and returns the rolling sum of amounts and transfer counts across all non-expired buckets.
// This implements lazy pruning: expired buckets are cleared only when accessed.
//
// Parameters:
//   - buckets: the array of MaxVelocityBuckets time buckets from the tracker
//   - currentTimeMs: the current block time in milliseconds
//   - isNumTransfers: if true, sum NumTransfers; if false, sum Amounts (balances)
//
// Returns:
//   - rollingAmounts: the summed balances across all active buckets (nil if isNumTransfers)
//   - rollingNumTransfers: the summed transfer count across all active buckets (zero if !isNumTransfers)
//   - prunedBuckets: the buckets array with expired entries zeroed out
func (vl *VelocityLimit) PruneAndSumBuckets(buckets []TimeBucket, currentTimeMs sdkmath.Uint, isNumTransfers bool) (
	rollingAmounts []*Balance,
	rollingNumTransfers sdkmath.Uint,
	prunedBuckets []TimeBucket,
) {
	rollingNumTransfers = sdkmath.NewUint(0)
	rollingAmounts = []*Balance{}
	prunedBuckets = make([]TimeBucket, len(buckets))
	copy(prunedBuckets, buckets)

	for i := range prunedBuckets {
		bucket := &prunedBuckets[i]

		// If the bucket's start time is zero/nil, it has never been written to — skip it
		if bucket.StartTime.IsNil() || bucket.StartTime.IsZero() {
			continue
		}

		// Prune expired buckets by zeroing them out
		if vl.IsBucketExpired(bucket.StartTime, currentTimeMs) {
			bucket.StartTime = sdkmath.NewUint(0)
			bucket.Amounts = []*Balance{}
			bucket.NumTransfers = sdkmath.NewUint(0)
			continue
		}

		// Accumulate non-expired bucket data
		if isNumTransfers {
			if !bucket.NumTransfers.IsNil() {
				rollingNumTransfers = rollingNumTransfers.Add(bucket.NumTransfers)
			}
		} else {
			if len(bucket.Amounts) > 0 {
				// We accumulate amounts; caller will handle the actual balance addition
				rollingAmounts = append(rollingAmounts, bucket.Amounts...)
			}
		}
	}

	return rollingAmounts, rollingNumTransfers, prunedBuckets
}

// InitBuckets creates a fresh array of MaxVelocityBuckets empty time buckets.
func InitBuckets() []TimeBucket {
	buckets := make([]TimeBucket, MaxVelocityBuckets)
	for i := range buckets {
		buckets[i] = TimeBucket{
			StartTime:    sdkmath.NewUint(0),
			Amounts:      []*Balance{},
			NumTransfers: sdkmath.NewUint(0),
		}
	}
	return buckets
}
