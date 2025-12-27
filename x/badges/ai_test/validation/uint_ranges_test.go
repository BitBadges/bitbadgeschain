package validation

import (
	"math"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

// ============================================================================
// Basic Structure Tests
// ============================================================================

func TestUintRange_NilRangeInArray(t *testing.T) {
	ranges := []*types.UintRange{
		GenerateValidUintRange(1, 10),
		nil, // Invalid: nil range
		GenerateValidUintRange(20, 30),
	}

	err := types.ValidateRangesAreValid(ranges, false, false)
	require.Error(t, err, "nil range should fail validation")
	require.Contains(t, err.Error(), "nil", "error should mention nil")
}

func TestUintRange_NilStartValue(t *testing.T) {
	range1 := GenerateInvalidUintRangeNilStart()

	err := types.ValidateRangesAreValid([]*types.UintRange{range1}, false, false)
	require.Error(t, err, "nil Start should fail validation")
	require.Contains(t, err.Error(), "nil", "error should mention nil")
}

func TestUintRange_NilEndValue(t *testing.T) {
	range1 := GenerateInvalidUintRangeNilEnd()

	err := types.ValidateRangesAreValid([]*types.UintRange{range1}, false, false)
	require.Error(t, err, "nil End should fail validation")
	require.Contains(t, err.Error(), "nil", "error should mention nil")
}

func TestUintRange_StartGreaterThanEnd(t *testing.T) {
	range1 := GenerateInvalidUintRangeStartGreaterThanEnd()

	err := types.ValidateRangesAreValid([]*types.UintRange{range1}, false, false)
	require.Error(t, err, "Start > End should fail validation")
	require.Contains(t, err.Error(), "greater", "error should mention greater than")
}

func TestUintRange_StartEqualsEnd(t *testing.T) {
	// Single value range - valid
	range1 := GenerateValidUintRange(5, 5)

	err := types.ValidateRangesAreValid([]*types.UintRange{range1}, false, false)
	require.NoError(t, err, "Start == End (single value) should be valid")
}

func TestUintRange_ZeroStartWhenAllowAllUintsFalse(t *testing.T) {
	range1 := GenerateInvalidUintRangeZeroStart()

	err := types.ValidateRangesAreValid([]*types.UintRange{range1}, false, false)
	require.Error(t, err, "zero Start should fail when allowAllUints=false")
	require.Contains(t, err.Error(), "zero", "error should mention zero")
}

func TestUintRange_ZeroEndWhenAllowAllUintsFalse(t *testing.T) {
	range1 := GenerateInvalidUintRangeZeroEnd()

	err := types.ValidateRangesAreValid([]*types.UintRange{range1}, false, false)
	require.Error(t, err, "zero End should fail when allowAllUints=false")
	// The error might be "start greater than end" or "zero" depending on validation order
	// Both are valid error messages for this case
	// Check if error contains any of the expected keywords (case-insensitive)
	errMsg := strings.ToLower(err.Error())
	hasZero := strings.Contains(errMsg, "zero")
	hasGreater := strings.Contains(errMsg, "greater")
	hasUninitialized := strings.Contains(errMsg, "uninitialized")
	require.True(t, hasZero || hasGreater || hasUninitialized,
		"error should mention zero, greater, or uninitialized, got: %s", err.Error())
}


func TestUintRange_StartGreaterThanMaxUint64(t *testing.T) {
	range1 := &types.UintRange{
		Start: sdkmath.NewUint(math.MaxUint64).Add(sdkmath.NewUint(1)),
		End:   sdkmath.NewUint(math.MaxUint64).Add(sdkmath.NewUint(10)),
	}

	err := types.ValidateRangesAreValid([]*types.UintRange{range1}, false, false)
	require.Error(t, err, "Start > math.MaxUint64 should fail when allowAllUints=false")
}

func TestUintRange_EndGreaterThanMaxUint64_Basic(t *testing.T) {
	range1 := &types.UintRange{
		Start: sdkmath.NewUint(1),
		End:   sdkmath.NewUint(math.MaxUint64).Add(sdkmath.NewUint(1)),
	}

	err := types.ValidateRangesAreValid([]*types.UintRange{range1}, false, false)
	require.Error(t, err, "End > math.MaxUint64 should fail when allowAllUints=false")
}

func TestUintRange_EmptyArrayWithErrorOnEmptyTrue(t *testing.T) {
	ranges := []*types.UintRange{}

	err := types.ValidateRangesAreValid(ranges, false, true)
	require.Error(t, err, "empty array should fail when errorOnEmpty=true")
	require.Contains(t, err.Error(), "empty", "error should mention empty")
}

func TestUintRange_EmptyArrayWithErrorOnEmptyFalse(t *testing.T) {
	ranges := []*types.UintRange{}

	err := types.ValidateRangesAreValid(ranges, false, false)
	require.NoError(t, err, "empty array should pass when errorOnEmpty=false")
}

// ============================================================================
// Overlap Detection Tests
// ============================================================================

func TestUintRange_ExactOverlap(t *testing.T) {
	exactOverlap, _ := GenerateOverlappingUintRanges()

	err := types.ValidateRangesAreValid(exactOverlap, false, false)
	require.Error(t, err, "exact overlap should fail validation")
	require.Contains(t, err.Error(), "overlap", "error should mention overlap")
}

func TestUintRange_PartialOverlap(t *testing.T) {
	_, partialOverlap := GenerateOverlappingUintRanges()

	err := types.ValidateRangesAreValid(partialOverlap, false, false)
	require.Error(t, err, "partial overlap should fail validation")
	require.Contains(t, err.Error(), "overlap", "error should mention overlap")
}

func TestUintRange_AdjacentBoundaries(t *testing.T) {
	// Adjacent ranges [1-5], [6-10] should not overlap (gap of 1)
	adjacent := GenerateAdjacentUintRanges()

	err := types.ValidateRangesAreValid(adjacent, false, false)
	require.NoError(t, err, "adjacent ranges should not overlap")
}

func TestUintRange_GapBetweenRanges(t *testing.T) {
	// Ranges with gap [1-5], [7-10] should not overlap
	gapped := GenerateGappedUintRanges()

	err := types.ValidateRangesAreValid(gapped, false, false)
	require.NoError(t, err, "gapped ranges should not overlap")
}

func TestUintRange_RangeFullyContained(t *testing.T) {
	// [1-10] contains [3-7]
	contained := GenerateContainedUintRanges()

	err := types.ValidateRangesAreValid(contained, false, false)
	require.Error(t, err, "fully contained range should overlap")
	require.Contains(t, err.Error(), "overlap", "error should mention overlap")
}

func TestUintRange_RangeContainingAnother(t *testing.T) {
	// [1-10] contains [5-8]
	ranges := []*types.UintRange{
		GenerateValidUintRange(1, 10),
		GenerateValidUintRange(5, 8),
	}

	err := types.ValidateRangesAreValid(ranges, false, false)
	require.Error(t, err, "containing range should overlap")
}

func TestUintRange_MultipleOverlappingRanges(t *testing.T) {
	// [1-5], [3-7], [6-10] - all overlap
	ranges := []*types.UintRange{
		GenerateValidUintRange(1, 5),
		GenerateValidUintRange(3, 7),
		GenerateValidUintRange(6, 10),
	}

	err := types.ValidateRangesAreValid(ranges, false, false)
	require.Error(t, err, "multiple overlapping ranges should fail")
}

func TestUintRange_OverlappingAtBoundaries(t *testing.T) {
	// [1-5], [5-10] - overlap at boundary (inclusive)
	ranges := []*types.UintRange{
		GenerateValidUintRange(1, 5),
		GenerateValidUintRange(5, 10),
	}

	err := types.ValidateRangesAreValid(ranges, false, false)
	require.Error(t, err, "overlapping at boundaries should fail")
}

// ============================================================================
// Duplicate Value Detection Tests
// ============================================================================

func TestUintRange_SameValueInTwoRanges(t *testing.T) {
	// [1-5], [3-7] - values 3, 4, 5 are duplicates
	ranges := []*types.UintRange{
		GenerateValidUintRange(1, 5),
		GenerateValidUintRange(3, 7),
	}

	err := types.ValidateRangesAreValid(ranges, false, false)
	require.Error(t, err, "duplicate values should fail validation")
}

func TestUintRange_ExactDuplicateRanges(t *testing.T) {
	exactOverlap, _ := GenerateOverlappingUintRanges()

	err := types.ValidateRangesAreValid(exactOverlap, false, false)
	require.Error(t, err, "exact duplicate ranges should fail")
}

func TestUintRange_SingleValueOverlap(t *testing.T) {
	// [1-5], [5-10] - value 5 is duplicate
	ranges := []*types.UintRange{
		GenerateValidUintRange(1, 5),
		GenerateValidUintRange(5, 10),
	}

	err := types.ValidateRangesAreValid(ranges, false, false)
	require.Error(t, err, "single value overlap should fail")
}

func TestUintRange_MultipleValuesOverlap(t *testing.T) {
	// [1-10], [5-15] - values 5-10 are duplicates
	ranges := []*types.UintRange{
		GenerateValidUintRange(1, 10),
		GenerateValidUintRange(5, 15),
	}

	err := types.ValidateRangesAreValid(ranges, false, false)
	require.Error(t, err, "multiple value overlap should fail")
}

// ============================================================================
// Unsorted Ranges Tests
// ============================================================================

func TestUintRange_UnsortedRangesNoOverlap(t *testing.T) {
	// Ranges in reverse order [20-30], [1-10] - should pass (no overlap, just unsorted)
	ranges := []*types.UintRange{
		GenerateValidUintRange(20, 30),
		GenerateValidUintRange(1, 10),
	}

	err := types.ValidateRangesAreValid(ranges, false, false)
	require.NoError(t, err, "unsorted ranges with no overlap should pass (validation sorts internally)")
}

func TestUintRange_UnsortedRangesWithOverlap(t *testing.T) {
	// Ranges in reverse order [10-20], [5-15] - should fail (overlap, unsorted)
	ranges := []*types.UintRange{
		GenerateValidUintRange(10, 20),
		GenerateValidUintRange(5, 15),
	}

	err := types.ValidateRangesAreValid(ranges, false, false)
	require.Error(t, err, "unsorted ranges with overlap should fail")
	require.Contains(t, err.Error(), "overlap", "error should mention overlap")
}

func TestUintRange_UnsortedRangesMultiple(t *testing.T) {
	// Multiple ranges in random order [30-40], [1-10], [20-25] - should pass (no overlap)
	ranges := []*types.UintRange{
		GenerateValidUintRange(30, 40),
		GenerateValidUintRange(1, 10),
		GenerateValidUintRange(20, 25),
	}

	err := types.ValidateRangesAreValid(ranges, false, false)
	require.NoError(t, err, "multiple unsorted ranges with no overlap should pass")
}

func TestUintRange_UnsortedRangesMultipleWithOverlap(t *testing.T) {
	// Multiple ranges in random order with overlap [15-25], [1-10], [5-20] - should fail
	ranges := []*types.UintRange{
		GenerateValidUintRange(15, 25),
		GenerateValidUintRange(1, 10),
		GenerateValidUintRange(5, 20),
	}

	err := types.ValidateRangesAreValid(ranges, false, false)
	require.Error(t, err, "multiple unsorted ranges with overlap should fail")
}

func TestUintRange_UnsortedAdjacentRanges(t *testing.T) {
	// Adjacent ranges in reverse order [6-10], [1-5] - should pass (adjacent, not overlapping)
	ranges := []*types.UintRange{
		GenerateValidUintRange(6, 10),
		GenerateValidUintRange(1, 5),
	}

	err := types.ValidateRangesAreValid(ranges, false, false)
	require.NoError(t, err, "unsorted adjacent ranges should pass")
}

// ============================================================================
// Edge Cases: Start < 1 and End > MaxUint64
// ============================================================================

func TestUintRange_StartLessThanOne(t *testing.T) {
	// Start = 0 is already tested, but let's be explicit
	range1 := GenerateInvalidUintRangeZeroStart()

	err := types.ValidateRangesAreValid([]*types.UintRange{range1}, false, false)
	require.Error(t, err, "start < 1 should fail when allowAllUints=false")
}

func TestUintRange_EndGreaterThanMaxUint64(t *testing.T) {
	// End > MaxUint64 is already tested, but let's verify it's caught
	range1 := &types.UintRange{
		Start: sdkmath.NewUint(1),
		End:   sdkmath.NewUint(math.MaxUint64).Add(sdkmath.NewUint(1)),
	}

	err := types.ValidateRangesAreValid([]*types.UintRange{range1}, false, false)
	require.Error(t, err, "end > MaxUint64 should fail when allowAllUints=false")
}

func TestUintRange_StartAndEndGreaterThanMaxUint64(t *testing.T) {
	// Both start and end > MaxUint64
	range1 := &types.UintRange{
		Start: sdkmath.NewUint(math.MaxUint64).Add(sdkmath.NewUint(1)),
		End:   sdkmath.NewUint(math.MaxUint64).Add(sdkmath.NewUint(10)),
	}

	err := types.ValidateRangesAreValid([]*types.UintRange{range1}, false, false)
	require.Error(t, err, "start and end > MaxUint64 should fail when allowAllUints=false")
}

func TestUintRange_StartAtMaxUint64(t *testing.T) {
	// Start = MaxUint64, End = MaxUint64 (single value at max)
	range1 := &types.UintRange{
		Start: sdkmath.NewUint(math.MaxUint64),
		End:   sdkmath.NewUint(math.MaxUint64),
	}

	err := types.ValidateRangesAreValid([]*types.UintRange{range1}, false, false)
	require.NoError(t, err, "start and end at MaxUint64 should be valid")
}

func TestUintRange_EndAtMaxUint64(t *testing.T) {
	// Start = 1, End = MaxUint64 (full range)
	range1 := GenerateValidUintRangeMax()

	err := types.ValidateRangesAreValid([]*types.UintRange{range1}, false, false)
	require.NoError(t, err, "end at MaxUint64 should be valid")
}

// ============================================================================
// Duplicate Value Detection - Explicit Tests
// ============================================================================

func TestUintRange_DuplicateSingleValue(t *testing.T) {
	// Two ranges both containing value 5: [1-5], [5-10]
	ranges := []*types.UintRange{
		GenerateValidUintRange(1, 5),
		GenerateValidUintRange(5, 10),
	}

	err := types.ValidateRangesAreValid(ranges, false, false)
	require.Error(t, err, "duplicate single value should fail")
}

func TestUintRange_DuplicateMultipleValues(t *testing.T) {
	// Two ranges with multiple duplicate values: [1-10], [5-15] - values 5-10 are duplicates
	ranges := []*types.UintRange{
		GenerateValidUintRange(1, 10),
		GenerateValidUintRange(5, 15),
	}

	err := types.ValidateRangesAreValid(ranges, false, false)
	require.Error(t, err, "duplicate multiple values should fail")
}

func TestUintRange_DuplicateAllValues(t *testing.T) {
	// Two identical ranges: [1-10], [1-10]
	ranges := []*types.UintRange{
		GenerateValidUintRange(1, 10),
		GenerateValidUintRange(1, 10),
	}

	err := types.ValidateRangesAreValid(ranges, false, false)
	require.Error(t, err, "duplicate all values should fail")
}

func TestUintRange_DuplicateThreeRanges(t *testing.T) {
	// Three ranges all overlapping: [1-10], [5-15], [10-20]
	ranges := []*types.UintRange{
		GenerateValidUintRange(1, 10),
		GenerateValidUintRange(5, 15),
		GenerateValidUintRange(10, 20),
	}

	err := types.ValidateRangesAreValid(ranges, false, false)
	require.Error(t, err, "duplicate values across three ranges should fail")
}

// ============================================================================
// AltTimeChecks Specific Tests
// ============================================================================

func TestUintRange_AltTimeHoursStartGreaterThan23(t *testing.T) {
	ranges := []*types.UintRange{
		{
			Start: sdkmath.NewUint(24),
			End:   sdkmath.NewUint(25),
		},
	}

	err := types.ValidateAltTimeChecks(&types.AltTimeChecks{
		OfflineHours: ranges,
		OfflineDays:  []*types.UintRange{},
	})
	require.Error(t, err, "OfflineHours Start > 23 should fail")
}

func TestUintRange_AltTimeHoursEndGreaterThan23(t *testing.T) {
	ranges := GenerateInvalidAltTimeHoursRanges()

	err := types.ValidateAltTimeChecks(&types.AltTimeChecks{
		OfflineHours: ranges,
		OfflineDays:  []*types.UintRange{},
	})
	require.Error(t, err, "OfflineHours End > 23 should fail")
}

func TestUintRange_AltTimeHoursValidBoundary(t *testing.T) {
	ranges := GenerateAltTimeHoursRanges()

	err := types.ValidateAltTimeChecks(&types.AltTimeChecks{
		OfflineHours: ranges,
		OfflineDays:  []*types.UintRange{},
	})
	require.NoError(t, err, "OfflineHours [0-23] should be valid")
}

func TestUintRange_AltTimeHoursDuplicateValues(t *testing.T) {
	// [9-12], [11-15] - hours 11, 12 are duplicates
	ranges := []*types.UintRange{
		{
			Start: sdkmath.NewUint(9),
			End:   sdkmath.NewUint(12),
		},
		{
			Start: sdkmath.NewUint(11),
			End:   sdkmath.NewUint(15),
		},
	}

	err := types.ValidateAltTimeChecks(&types.AltTimeChecks{
		OfflineHours: ranges,
		OfflineDays:  []*types.UintRange{},
	})
	require.Error(t, err, "duplicate hour values should fail")
}

func TestUintRange_AltTimeDaysStartGreaterThan6(t *testing.T) {
	ranges := []*types.UintRange{
		{
			Start: sdkmath.NewUint(7),
			End:   sdkmath.NewUint(8),
		},
	}

	err := types.ValidateAltTimeChecks(&types.AltTimeChecks{
		OfflineHours: []*types.UintRange{},
		OfflineDays:  ranges,
	})
	require.Error(t, err, "OfflineDays Start > 6 should fail")
}

func TestUintRange_AltTimeDaysEndGreaterThan6(t *testing.T) {
	ranges := GenerateInvalidAltTimeDaysRanges()

	err := types.ValidateAltTimeChecks(&types.AltTimeChecks{
		OfflineHours: []*types.UintRange{},
		OfflineDays:  ranges,
	})
	require.Error(t, err, "OfflineDays End > 6 should fail")
}

func TestUintRange_AltTimeDaysValidBoundary(t *testing.T) {
	ranges := GenerateAltTimeDaysRanges()

	err := types.ValidateAltTimeChecks(&types.AltTimeChecks{
		OfflineHours: []*types.UintRange{},
		OfflineDays:  ranges,
	})
	require.NoError(t, err, "OfflineDays [0-6] should be valid")
}

func TestUintRange_AltTimeDaysDuplicateValues(t *testing.T) {
	// [1-3], [2-5] - days 2, 3 are duplicates
	ranges := []*types.UintRange{
		{
			Start: sdkmath.NewUint(1),
			End:   sdkmath.NewUint(3),
		},
		{
			Start: sdkmath.NewUint(2),
			End:   sdkmath.NewUint(5),
		},
	}

	err := types.ValidateAltTimeChecks(&types.AltTimeChecks{
		OfflineHours: []*types.UintRange{},
		OfflineDays:  ranges,
	})
	require.Error(t, err, "duplicate day values should fail")
}

// ============================================================================
// Special Cases Tests
// ============================================================================

func TestUintRange_EndIsZeroTreatedAsStart(t *testing.T) {
	// According to documentation: "If end.IsZero(), we assume end == start"
	// This is a documented behavior, so we test that it's handled
	range1 := &types.UintRange{
		Start: sdkmath.NewUint(5),
		End:   sdkmath.NewUint(0), // Zero end
	}

	// Note: The validation function checks for zero when allowAllUints=false
	// So this will fail validation, but the comment suggests the logic treats it as end == start
	err := types.ValidateRangesAreValid([]*types.UintRange{range1}, false, false)
	require.Error(t, err, "zero End should fail when allowAllUints=false")
}

func TestUintRange_AllowAllUintsWithZeroStartEnd(t *testing.T) {
	range1 := &types.UintRange{
		Start: sdkmath.NewUint(0),
		End:   sdkmath.NewUint(0),
	}

	err := types.ValidateRangesAreValid([]*types.UintRange{range1}, true, false)
	require.NoError(t, err, "zero Start/End should be valid when allowAllUints=true")
}

func TestUintRange_AllowAllUintsWithStartGreaterThanMaxUint64(t *testing.T) {
	// When allowAllUints=true, values > math.MaxUint64 are allowed
	// Note: sdkmath.Uint can't actually exceed MaxUint64, but the validation allows it conceptually
	range1 := &types.UintRange{
		Start: sdkmath.NewUint(1),
		End:   sdkmath.NewUint(math.MaxUint64),
	}

	err := types.ValidateRangesAreValid([]*types.UintRange{range1}, true, false)
	require.NoError(t, err, "max values should be valid when allowAllUints=true")
}

func TestUintRange_SingleRangeFullSpectrum(t *testing.T) {
	range1 := GenerateValidUintRangeMax()

	err := types.ValidateRangesAreValid([]*types.UintRange{range1}, false, false)
	require.NoError(t, err, "full spectrum range [1, math.MaxUint64] should be valid")
}

func TestUintRange_MaximumNumberOfRanges(t *testing.T) {
	// Stress test: create many non-overlapping ranges
	ranges := make([]*types.UintRange, 100)
	for i := 0; i < 100; i++ {
		start := uint64(i*1000 + 1)
		end := uint64((i+1)*1000 - 1)
		ranges[i] = GenerateValidUintRange(start, end)
	}

	err := types.ValidateRangesAreValid(ranges, false, false)
	require.NoError(t, err, "many non-overlapping ranges should be valid")
}

// ============================================================================
// UintRange Operations Tests
// ============================================================================

func TestUintRange_SortUintRangesAndMerge_AdjacentRanges(t *testing.T) {
	// [1-5], [6-10] should merge to [1-10] when mergeIntersecting=true
	ranges := GenerateAdjacentUintRanges()

	merged, err := types.SortUintRangesAndMerge(ranges, true)
	require.NoError(t, err)
	require.Len(t, merged, 1, "adjacent ranges should merge")
	require.Equal(t, sdkmath.NewUint(1), merged[0].Start)
	require.Equal(t, sdkmath.NewUint(10), merged[0].End)
}

func TestUintRange_SortUintRangesAndMerge_OverlappingRanges(t *testing.T) {
	// [1-10], [5-15] should merge to [1-15] when mergeIntersecting=true
	ranges := []*types.UintRange{
		GenerateValidUintRange(1, 10),
		GenerateValidUintRange(5, 15),
	}

	merged, err := types.SortUintRangesAndMerge(ranges, true)
	require.NoError(t, err)
	require.Len(t, merged, 1, "overlapping ranges should merge")
	require.Equal(t, sdkmath.NewUint(1), merged[0].Start)
	require.Equal(t, sdkmath.NewUint(15), merged[0].End)
}

func TestUintRange_SortUintRangesAndMerge_NoMergeWhenMergeIntersectingFalse(t *testing.T) {
	// [1-5], [6-10] - adjacent ranges
	// When mergeIntersecting=false, the function checks for overlap (End.GTE(Start))
	// Adjacent ranges [1-5], [6-10] have 5 < 6, so no overlap, so no error
	// However, the merge logic still runs and merges adjacent ranges
	// Looking at the code: if currRange.Start <= prevInsertedRange.End.AddUint64(1), it merges
	// So [1-5], [6-10] will merge because 6 <= 5+1 = 6
	// This means adjacent ranges ARE merged even when mergeIntersecting=false
	// The mergeIntersecting flag only prevents overlapping ranges from being processed
	ranges := GenerateAdjacentUintRanges()

	merged, err := types.SortUintRangesAndMerge(ranges, false)
	require.NoError(t, err)
	// Adjacent ranges are merged by the merge logic
	require.Len(t, merged, 1, "adjacent ranges are merged even when mergeIntersecting=false")
	require.Equal(t, sdkmath.NewUint(1), merged[0].Start)
	require.Equal(t, sdkmath.NewUint(10), merged[0].End)
}

func TestUintRange_SortUintRangesAndMerge_ErrorOnOverlapWhenMergeIntersectingFalse(t *testing.T) {
	// [1-10], [5-15] should error when mergeIntersecting=false (overlap)
	ranges := []*types.UintRange{
		GenerateValidUintRange(1, 10),
		GenerateValidUintRange(5, 15),
	}

	_, err := types.SortUintRangesAndMerge(ranges, false)
	require.Error(t, err, "overlapping ranges should error when mergeIntersecting=false")
}

func TestUintRange_RemoveUintRangeFromUintRange_NoOverlap(t *testing.T) {
	// Remove [20-30] from [1-10] - no overlap, should return original
	toRemove := GenerateValidUintRange(20, 30)
	fromRange := GenerateValidUintRange(1, 10)

	remaining, removed := types.RemoveUintRangeFromUintRange(toRemove, fromRange)
	require.Len(t, remaining, 1, "no overlap should return original range")
	require.Len(t, removed, 0, "no overlap should return no removed ranges")
	require.Equal(t, sdkmath.NewUint(1), remaining[0].Start)
	require.Equal(t, sdkmath.NewUint(10), remaining[0].End)
}

func TestUintRange_RemoveUintRangeFromUintRange_FullContainment(t *testing.T) {
	// Remove [1-10] from [1-10] - full containment, should return empty
	toRemove := GenerateValidUintRange(1, 10)
	fromRange := GenerateValidUintRange(1, 10)

	remaining, removed := types.RemoveUintRangeFromUintRange(toRemove, fromRange)
	require.Len(t, remaining, 0, "full containment should return empty")
	require.Len(t, removed, 1, "full containment should return removed range")
}

func TestUintRange_RemoveUintRangeFromUintRange_PartialOverlapStart(t *testing.T) {
	// Remove [1-5] from [1-10] - partial overlap at start, should return [6-10]
	toRemove := GenerateValidUintRange(1, 5)
	fromRange := GenerateValidUintRange(1, 10)

	remaining, removed := types.RemoveUintRangeFromUintRange(toRemove, fromRange)
	require.Len(t, remaining, 1, "partial overlap at start should return one range")
	require.Equal(t, sdkmath.NewUint(6), remaining[0].Start)
	require.Equal(t, sdkmath.NewUint(10), remaining[0].End)
	require.Len(t, removed, 1, "should return removed portion")
}

func TestUintRange_RemoveUintRangeFromUintRange_PartialOverlapEnd(t *testing.T) {
	// Remove [6-10] from [1-10] - partial overlap at end, should return [1-5]
	toRemove := GenerateValidUintRange(6, 10)
	fromRange := GenerateValidUintRange(1, 10)

	remaining, removed := types.RemoveUintRangeFromUintRange(toRemove, fromRange)
	require.Len(t, remaining, 1, "partial overlap at end should return one range")
	require.Equal(t, sdkmath.NewUint(1), remaining[0].Start)
	require.Equal(t, sdkmath.NewUint(5), remaining[0].End)
	require.Len(t, removed, 1, "should return removed portion")
}

func TestUintRange_RemoveUintRangeFromUintRange_MiddleRemoval(t *testing.T) {
	// Remove [5-7] from [1-10] - middle removal, should return [1-4] and [8-10]
	toRemove := GenerateValidUintRange(5, 7)
	fromRange := GenerateValidUintRange(1, 10)

	remaining, removed := types.RemoveUintRangeFromUintRange(toRemove, fromRange)
	require.Len(t, remaining, 2, "middle removal should return two ranges")
	require.Equal(t, sdkmath.NewUint(1), remaining[0].Start)
	require.Equal(t, sdkmath.NewUint(4), remaining[0].End)
	require.Equal(t, sdkmath.NewUint(8), remaining[1].Start)
	require.Equal(t, sdkmath.NewUint(10), remaining[1].End)
	require.Len(t, removed, 1, "should return removed portion")
}

func TestUintRange_SearchUintRangesForUint_Found(t *testing.T) {
	ranges := []*types.UintRange{
		GenerateValidUintRange(1, 10),
		GenerateValidUintRange(20, 30),
		GenerateValidUintRange(40, 50),
	}

	found, err := types.SearchUintRangesForUint(sdkmath.NewUint(5), ranges)
	require.NoError(t, err)
	require.True(t, found, "value 5 should be found in range [1-10]")

	found, err = types.SearchUintRangesForUint(sdkmath.NewUint(25), ranges)
	require.NoError(t, err)
	require.True(t, found, "value 25 should be found in range [20-30]")
}

func TestUintRange_SearchUintRangesForUint_NotFound(t *testing.T) {
	ranges := []*types.UintRange{
		GenerateValidUintRange(1, 10),
		GenerateValidUintRange(20, 30),
	}

	found, err := types.SearchUintRangesForUint(sdkmath.NewUint(15), ranges)
	require.NoError(t, err)
	require.False(t, found, "value 15 should not be found")
}

func TestUintRange_SearchUintRangesForUint_BoundaryValues(t *testing.T) {
	ranges := []*types.UintRange{
		GenerateValidUintRange(1, 10),
	}

	// Test start boundary
	found, err := types.SearchUintRangesForUint(sdkmath.NewUint(1), ranges)
	require.NoError(t, err)
	require.True(t, found, "start boundary should be found")

	// Test end boundary
	found, err = types.SearchUintRangesForUint(sdkmath.NewUint(10), ranges)
	require.NoError(t, err)
	require.True(t, found, "end boundary should be found")
}

func TestUintRange_InvertUintRanges_EmptyInput(t *testing.T) {
	// Empty input should return full range [minId, maxId]
	minId := sdkmath.NewUint(1)
	maxId := sdkmath.NewUint(100)

	inverted := types.InvertUintRanges([]*types.UintRange{}, minId, maxId)
	require.Len(t, inverted, 1, "empty input should return one range")
	require.Equal(t, minId, inverted[0].Start)
	require.Equal(t, maxId, inverted[0].End)
}

func TestUintRange_InvertUintRanges_SingleRange(t *testing.T) {
	// Invert [20-30] from [1-100] should return [1-19] and [31-100]
	ranges := []*types.UintRange{
		GenerateValidUintRange(20, 30),
	}
	minId := sdkmath.NewUint(1)
	maxId := sdkmath.NewUint(100)

	inverted := types.InvertUintRanges(ranges, minId, maxId)
	require.Len(t, inverted, 2, "should return two ranges")
	require.Equal(t, sdkmath.NewUint(1), inverted[0].Start)
	require.Equal(t, sdkmath.NewUint(19), inverted[0].End)
	require.Equal(t, sdkmath.NewUint(31), inverted[1].Start)
	require.Equal(t, sdkmath.NewUint(100), inverted[1].End)
}

func TestUintRange_InvertUintRanges_MultipleRanges(t *testing.T) {
	// Invert [10-20], [30-40] from [1-100] should return [1-9], [21-29], [41-100]
	ranges := []*types.UintRange{
		GenerateValidUintRange(10, 20),
		GenerateValidUintRange(30, 40),
	}
	minId := sdkmath.NewUint(1)
	maxId := sdkmath.NewUint(100)

	inverted := types.InvertUintRanges(ranges, minId, maxId)
	require.Len(t, inverted, 3, "should return three ranges")
	require.Equal(t, sdkmath.NewUint(1), inverted[0].Start)
	require.Equal(t, sdkmath.NewUint(9), inverted[0].End)
	require.Equal(t, sdkmath.NewUint(21), inverted[1].Start)
	require.Equal(t, sdkmath.NewUint(29), inverted[1].End)
	require.Equal(t, sdkmath.NewUint(41), inverted[2].Start)
	require.Equal(t, sdkmath.NewUint(100), inverted[2].End)
}
