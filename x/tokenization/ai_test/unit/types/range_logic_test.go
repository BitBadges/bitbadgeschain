package types_test

import (
	"math"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// RangeLogicTestSuite tests the range merging, splitting, and validation logic
type RangeLogicTestSuite struct {
	suite.Suite
}

func TestRangeLogicTestSuite(t *testing.T) {
	suite.Run(t, new(RangeLogicTestSuite))
}

// =============================================================================
// Range Merging Tests
// =============================================================================

// TestRangeMerging_AdjacentRanges tests that adjacent ranges merge correctly
func (suite *RangeLogicTestSuite) TestRangeMerging_AdjacentRanges() {
	// [1-10] and [11-20] should merge to [1-20]
	ranges := []*types.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
		{Start: sdkmath.NewUint(11), End: sdkmath.NewUint(20)},
	}

	merged := types.SortUintRangesAndMergeAdjacentAndIntersecting(ranges)

	suite.Require().Len(merged, 1, "adjacent ranges should merge into one")
	suite.Require().Equal(sdkmath.NewUint(1), merged[0].Start)
	suite.Require().Equal(sdkmath.NewUint(20), merged[0].End)
}

// TestRangeMerging_OverlappingRanges tests that overlapping ranges merge correctly
func (suite *RangeLogicTestSuite) TestRangeMerging_OverlappingRanges() {
	// [1-15] and [10-20] should merge to [1-20]
	ranges := []*types.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(15)},
		{Start: sdkmath.NewUint(10), End: sdkmath.NewUint(20)},
	}

	merged := types.SortUintRangesAndMergeAdjacentAndIntersecting(ranges)

	suite.Require().Len(merged, 1, "overlapping ranges should merge into one")
	suite.Require().Equal(sdkmath.NewUint(1), merged[0].Start)
	suite.Require().Equal(sdkmath.NewUint(20), merged[0].End)
}

// TestRangeMerging_ContainedRanges tests that contained ranges are handled correctly
func (suite *RangeLogicTestSuite) TestRangeMerging_ContainedRanges() {
	// [1-100] and [10-20] - second is fully contained in first
	ranges := []*types.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
		{Start: sdkmath.NewUint(10), End: sdkmath.NewUint(20)},
	}

	merged := types.SortUintRangesAndMergeAdjacentAndIntersecting(ranges)

	suite.Require().Len(merged, 1, "contained range should be absorbed")
	suite.Require().Equal(sdkmath.NewUint(1), merged[0].Start)
	suite.Require().Equal(sdkmath.NewUint(100), merged[0].End)
}

// TestRangeMerging_DisjointRanges tests that disjoint ranges remain separate
func (suite *RangeLogicTestSuite) TestRangeMerging_DisjointRanges() {
	// [1-10] and [20-30] should remain separate
	ranges := []*types.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
		{Start: sdkmath.NewUint(20), End: sdkmath.NewUint(30)},
	}

	merged := types.SortUintRangesAndMergeAdjacentAndIntersecting(ranges)

	suite.Require().Len(merged, 2, "disjoint ranges should remain separate")
	suite.Require().Equal(sdkmath.NewUint(1), merged[0].Start)
	suite.Require().Equal(sdkmath.NewUint(10), merged[0].End)
	suite.Require().Equal(sdkmath.NewUint(20), merged[1].Start)
	suite.Require().Equal(sdkmath.NewUint(30), merged[1].End)
}

// TestRangeMerging_MultipleAdjacentRanges tests merging multiple adjacent ranges
func (suite *RangeLogicTestSuite) TestRangeMerging_MultipleAdjacentRanges() {
	// [1-10], [11-20], [21-30] should merge to [1-30]
	ranges := []*types.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
		{Start: sdkmath.NewUint(11), End: sdkmath.NewUint(20)},
		{Start: sdkmath.NewUint(21), End: sdkmath.NewUint(30)},
	}

	merged := types.SortUintRangesAndMergeAdjacentAndIntersecting(ranges)

	suite.Require().Len(merged, 1, "multiple adjacent ranges should merge into one")
	suite.Require().Equal(sdkmath.NewUint(1), merged[0].Start)
	suite.Require().Equal(sdkmath.NewUint(30), merged[0].End)
}

// TestRangeMerging_UnsortedRanges tests that ranges are sorted before merging
func (suite *RangeLogicTestSuite) TestRangeMerging_UnsortedRanges() {
	// [20-30], [1-10], [11-20] (unsorted) should merge to [1-30]
	ranges := []*types.UintRange{
		{Start: sdkmath.NewUint(20), End: sdkmath.NewUint(30)},
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
		{Start: sdkmath.NewUint(11), End: sdkmath.NewUint(19)},
	}

	merged := types.SortUintRangesAndMergeAdjacentAndIntersecting(ranges)

	suite.Require().Len(merged, 1, "unsorted adjacent ranges should be sorted and merged")
	suite.Require().Equal(sdkmath.NewUint(1), merged[0].Start)
	suite.Require().Equal(sdkmath.NewUint(30), merged[0].End)
}

// TestRangeMerging_EmptyInput tests merging with empty input
func (suite *RangeLogicTestSuite) TestRangeMerging_EmptyInput() {
	ranges := []*types.UintRange{}

	merged := types.SortUintRangesAndMergeAdjacentAndIntersecting(ranges)

	suite.Require().Len(merged, 0, "empty input should return empty output")
}

// TestRangeMerging_SingleRange tests merging with single range
func (suite *RangeLogicTestSuite) TestRangeMerging_SingleRange() {
	ranges := []*types.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
	}

	merged := types.SortUintRangesAndMergeAdjacentAndIntersecting(ranges)

	suite.Require().Len(merged, 1, "single range should remain unchanged")
	suite.Require().Equal(sdkmath.NewUint(1), merged[0].Start)
	suite.Require().Equal(sdkmath.NewUint(100), merged[0].End)
}

// =============================================================================
// Range Splitting Tests
// =============================================================================

// TestRangeSplitting_RemoveFromMiddle tests removal from middle creates two disjoint ranges
func (suite *RangeLogicTestSuite) TestRangeSplitting_RemoveFromMiddle() {
	// Remove [50-60] from [1-100] should give [1-49] and [61-100]
	rangeToRemoveFrom := &types.UintRange{
		Start: sdkmath.NewUint(1),
		End:   sdkmath.NewUint(100),
	}
	rangeToRemove := &types.UintRange{
		Start: sdkmath.NewUint(50),
		End:   sdkmath.NewUint(60),
	}

	remaining, removed := types.RemoveUintRangeFromUintRange(rangeToRemove, rangeToRemoveFrom)

	suite.Require().Len(remaining, 2, "removal from middle should create two ranges")
	suite.Require().Equal(sdkmath.NewUint(1), remaining[0].Start)
	suite.Require().Equal(sdkmath.NewUint(49), remaining[0].End)
	suite.Require().Equal(sdkmath.NewUint(61), remaining[1].Start)
	suite.Require().Equal(sdkmath.NewUint(100), remaining[1].End)

	// Verify removed ranges
	suite.Require().True(len(removed) > 0, "should have removed ranges")
}

// TestRangeSplitting_RemoveFromStart tests removal from start
func (suite *RangeLogicTestSuite) TestRangeSplitting_RemoveFromStart() {
	// Remove [1-50] from [1-100] should give [51-100]
	rangeToRemoveFrom := &types.UintRange{
		Start: sdkmath.NewUint(1),
		End:   sdkmath.NewUint(100),
	}
	rangeToRemove := &types.UintRange{
		Start: sdkmath.NewUint(1),
		End:   sdkmath.NewUint(50),
	}

	remaining, removed := types.RemoveUintRangeFromUintRange(rangeToRemove, rangeToRemoveFrom)

	suite.Require().Len(remaining, 1, "removal from start should create one range")
	suite.Require().Equal(sdkmath.NewUint(51), remaining[0].Start)
	suite.Require().Equal(sdkmath.NewUint(100), remaining[0].End)
	suite.Require().True(len(removed) > 0, "should have removed ranges")
}

// TestRangeSplitting_RemoveFromEnd tests removal from end
func (suite *RangeLogicTestSuite) TestRangeSplitting_RemoveFromEnd() {
	// Remove [50-100] from [1-100] should give [1-49]
	rangeToRemoveFrom := &types.UintRange{
		Start: sdkmath.NewUint(1),
		End:   sdkmath.NewUint(100),
	}
	rangeToRemove := &types.UintRange{
		Start: sdkmath.NewUint(50),
		End:   sdkmath.NewUint(100),
	}

	remaining, removed := types.RemoveUintRangeFromUintRange(rangeToRemove, rangeToRemoveFrom)

	suite.Require().Len(remaining, 1, "removal from end should create one range")
	suite.Require().Equal(sdkmath.NewUint(1), remaining[0].Start)
	suite.Require().Equal(sdkmath.NewUint(49), remaining[0].End)
	suite.Require().True(len(removed) > 0, "should have removed ranges")
}

// TestRangeSplitting_RemoveEntireRange tests removal of entire range
func (suite *RangeLogicTestSuite) TestRangeSplitting_RemoveEntireRange() {
	// Remove [1-100] from [1-100] should give nothing
	rangeToRemoveFrom := &types.UintRange{
		Start: sdkmath.NewUint(1),
		End:   sdkmath.NewUint(100),
	}
	rangeToRemove := &types.UintRange{
		Start: sdkmath.NewUint(1),
		End:   sdkmath.NewUint(100),
	}

	remaining, removed := types.RemoveUintRangeFromUintRange(rangeToRemove, rangeToRemoveFrom)

	suite.Require().Len(remaining, 0, "removal of entire range should leave nothing")
	suite.Require().Len(removed, 1, "should have one removed range")
}

// TestRangeSplitting_NoOverlap tests removal with no overlap
func (suite *RangeLogicTestSuite) TestRangeSplitting_NoOverlap() {
	// Remove [200-300] from [1-100] - no overlap
	rangeToRemoveFrom := &types.UintRange{
		Start: sdkmath.NewUint(1),
		End:   sdkmath.NewUint(100),
	}
	rangeToRemove := &types.UintRange{
		Start: sdkmath.NewUint(200),
		End:   sdkmath.NewUint(300),
	}

	remaining, removed := types.RemoveUintRangeFromUintRange(rangeToRemove, rangeToRemoveFrom)

	suite.Require().Len(remaining, 1, "no overlap should leave original range")
	suite.Require().Equal(sdkmath.NewUint(1), remaining[0].Start)
	suite.Require().Equal(sdkmath.NewUint(100), remaining[0].End)
	suite.Require().Len(removed, 0, "should have no removed ranges")
}

// TestRangeSplitting_RemoveMultipleRanges tests removing multiple ranges
func (suite *RangeLogicTestSuite) TestRangeSplitting_RemoveMultipleRanges() {
	// Remove [20-30] and [70-80] from [1-100]
	rangesToRemoveFrom := []*types.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
	}
	rangesToRemove := []*types.UintRange{
		{Start: sdkmath.NewUint(20), End: sdkmath.NewUint(30)},
		{Start: sdkmath.NewUint(70), End: sdkmath.NewUint(80)},
	}

	remaining, removed := types.RemoveUintRangesFromUintRanges(rangesToRemove, rangesToRemoveFrom)

	suite.Require().Len(remaining, 3, "removing two ranges should create three remaining ranges")
	suite.Require().True(len(removed) > 0, "should have removed ranges")
}

// =============================================================================
// Range Validation Tests
// =============================================================================

// TestRangeValidation_StartGreaterThanEndRejected tests that start > end is rejected
func (suite *RangeLogicTestSuite) TestRangeValidation_StartGreaterThanEndRejected() {
	ranges := []*types.UintRange{
		{Start: sdkmath.NewUint(100), End: sdkmath.NewUint(50)}, // Invalid: start > end
	}

	err := types.ValidateRangesAreValid(ranges, false, true)
	suite.Require().Error(err, "range with start > end should be rejected")
	suite.Require().Contains(err.Error(), "greater than end")
}

// TestRangeValidation_SingleValueRangeAllowed tests that single-value range (start == end) is allowed
func (suite *RangeLogicTestSuite) TestRangeValidation_SingleValueRangeAllowed() {
	ranges := []*types.UintRange{
		{Start: sdkmath.NewUint(50), End: sdkmath.NewUint(50)}, // Valid: start == end
	}

	err := types.ValidateRangesAreValid(ranges, false, true)
	suite.Require().NoError(err, "single-value range (start == end) should be allowed")
}

// TestRangeValidation_MaxUint64BoundaryHandled tests MaxUint64 boundary values
func (suite *RangeLogicTestSuite) TestRangeValidation_MaxUint64BoundaryHandled() {
	ranges := []*types.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
	}

	err := types.ValidateRangesAreValid(ranges, false, true)
	suite.Require().NoError(err, "range ending at MaxUint64 should be valid")
}

// TestRangeValidation_ZeroStartRejected tests that zero start is rejected when not allowing all uints
func (suite *RangeLogicTestSuite) TestRangeValidation_ZeroStartRejected() {
	ranges := []*types.UintRange{
		{Start: sdkmath.NewUint(0), End: sdkmath.NewUint(100)},
	}

	err := types.ValidateRangesAreValid(ranges, false, true)
	suite.Require().Error(err, "range with zero start should be rejected when allowAllUints is false")
}

// TestRangeValidation_ZeroStartAllowedWithFlag tests that zero start is allowed with allowAllUints flag
func (suite *RangeLogicTestSuite) TestRangeValidation_ZeroStartAllowedWithFlag() {
	ranges := []*types.UintRange{
		{Start: sdkmath.NewUint(0), End: sdkmath.NewUint(100)},
	}

	err := types.ValidateRangesAreValid(ranges, true, true) // allowAllUints = true
	suite.Require().NoError(err, "range with zero start should be allowed when allowAllUints is true")
}

// TestRangeValidation_OverlappingRangesRejected tests that overlapping ranges are rejected
func (suite *RangeLogicTestSuite) TestRangeValidation_OverlappingRangesRejected() {
	ranges := []*types.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(50)},
		{Start: sdkmath.NewUint(40), End: sdkmath.NewUint(100)}, // Overlaps with first
	}

	err := types.ValidateRangesAreValid(ranges, false, true)
	suite.Require().Error(err, "overlapping ranges should be rejected")
	suite.Require().Contains(err.Error(), "overlap")
}

// TestRangeValidation_NilRangeRejected tests that nil range is rejected
func (suite *RangeLogicTestSuite) TestRangeValidation_NilRangeRejected() {
	ranges := []*types.UintRange{
		nil,
	}

	err := types.ValidateRangesAreValid(ranges, false, true)
	suite.Require().Error(err, "nil range should be rejected")
}

// TestRangeValidation_EmptyRangesEmptyError tests that empty ranges can error when flag is set
func (suite *RangeLogicTestSuite) TestRangeValidation_EmptyRangesEmptyError() {
	ranges := []*types.UintRange{}

	// With errorOnEmpty = true
	err := types.ValidateRangesAreValid(ranges, false, true)
	suite.Require().Error(err, "empty ranges should error when errorOnEmpty is true")

	// With errorOnEmpty = false
	err = types.ValidateRangesAreValid(ranges, false, false)
	suite.Require().NoError(err, "empty ranges should be allowed when errorOnEmpty is false")
}

// TestRangeValidation_AdjacentRangesAllowed tests that adjacent (non-overlapping) ranges are allowed
func (suite *RangeLogicTestSuite) TestRangeValidation_AdjacentRangesAllowed() {
	ranges := []*types.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(50)},
		{Start: sdkmath.NewUint(51), End: sdkmath.NewUint(100)}, // Adjacent, not overlapping
	}

	err := types.ValidateRangesAreValid(ranges, false, true)
	suite.Require().NoError(err, "adjacent non-overlapping ranges should be valid")
}

// =============================================================================
// Range Search Tests
// =============================================================================

// TestRangeSearch_FindIdInRange tests searching for an ID in ranges
func (suite *RangeLogicTestSuite) TestRangeSearch_FindIdInRange() {
	ranges := []*types.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(50)},
		{Start: sdkmath.NewUint(100), End: sdkmath.NewUint(200)},
	}

	// Search for ID in first range
	found, err := types.SearchUintRangesForUint(sdkmath.NewUint(25), ranges)
	suite.Require().NoError(err)
	suite.Require().True(found, "ID 25 should be found in range [1-50]")

	// Search for ID in second range
	found, err = types.SearchUintRangesForUint(sdkmath.NewUint(150), ranges)
	suite.Require().NoError(err)
	suite.Require().True(found, "ID 150 should be found in range [100-200]")

	// Search for ID not in any range
	found, err = types.SearchUintRangesForUint(sdkmath.NewUint(75), ranges)
	suite.Require().NoError(err)
	suite.Require().False(found, "ID 75 should not be found (between ranges)")
}

// TestRangeSearch_FindIdAtBoundary tests searching for an ID at range boundaries
func (suite *RangeLogicTestSuite) TestRangeSearch_FindIdAtBoundary() {
	ranges := []*types.UintRange{
		{Start: sdkmath.NewUint(10), End: sdkmath.NewUint(20)},
	}

	// Search for start boundary
	found, err := types.SearchUintRangesForUint(sdkmath.NewUint(10), ranges)
	suite.Require().NoError(err)
	suite.Require().True(found, "ID at start boundary should be found")

	// Search for end boundary
	found, err = types.SearchUintRangesForUint(sdkmath.NewUint(20), ranges)
	suite.Require().NoError(err)
	suite.Require().True(found, "ID at end boundary should be found")

	// Search just outside boundaries
	found, err = types.SearchUintRangesForUint(sdkmath.NewUint(9), ranges)
	suite.Require().NoError(err)
	suite.Require().False(found, "ID just before start should not be found")

	found, err = types.SearchUintRangesForUint(sdkmath.NewUint(21), ranges)
	suite.Require().NoError(err)
	suite.Require().False(found, "ID just after end should not be found")
}

// =============================================================================
// Range Inversion Tests
// =============================================================================

// TestRangeInversion_InvertMiddleRange tests inverting ranges
func (suite *RangeLogicTestSuite) TestRangeInversion_InvertMiddleRange() {
	// Invert [50-100] from [1-200] should give [1-49] and [101-200]
	rangesToInvert := []*types.UintRange{
		{Start: sdkmath.NewUint(50), End: sdkmath.NewUint(100)},
	}

	inverted := types.InvertUintRanges(rangesToInvert, sdkmath.NewUint(1), sdkmath.NewUint(200))

	suite.Require().Len(inverted, 2, "inverting middle range should create two ranges")
	suite.Require().Equal(sdkmath.NewUint(1), inverted[0].Start)
	suite.Require().Equal(sdkmath.NewUint(49), inverted[0].End)
	suite.Require().Equal(sdkmath.NewUint(101), inverted[1].Start)
	suite.Require().Equal(sdkmath.NewUint(200), inverted[1].End)
}

// TestRangeInversion_InvertFullRange tests inverting the entire range
func (suite *RangeLogicTestSuite) TestRangeInversion_InvertFullRange() {
	// Invert [1-200] from [1-200] should give nothing
	rangesToInvert := []*types.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(200)},
	}

	inverted := types.InvertUintRanges(rangesToInvert, sdkmath.NewUint(1), sdkmath.NewUint(200))

	suite.Require().Len(inverted, 0, "inverting entire range should give nothing")
}

// TestRangeInversion_InvertEmpty tests inverting empty ranges
func (suite *RangeLogicTestSuite) TestRangeInversion_InvertEmpty() {
	// Invert nothing from [1-200] should give [1-200]
	rangesToInvert := []*types.UintRange{}

	inverted := types.InvertUintRanges(rangesToInvert, sdkmath.NewUint(1), sdkmath.NewUint(200))

	suite.Require().Len(inverted, 1, "inverting nothing should give full range")
	suite.Require().Equal(sdkmath.NewUint(1), inverted[0].Start)
	suite.Require().Equal(sdkmath.NewUint(200), inverted[0].End)
}

// =============================================================================
// Deep Copy Tests
// =============================================================================

// TestDeepCopy_RangesAreIndependent tests that deep copied ranges are independent
func (suite *RangeLogicTestSuite) TestDeepCopy_RangesAreIndependent() {
	original := []*types.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
	}

	copied := types.DeepCopyRanges(original)

	// Modify the original
	original[0].Start = sdkmath.NewUint(50)

	// Verify copy is unchanged
	suite.Require().Equal(sdkmath.NewUint(1), copied[0].Start, "deep copy should be independent of original")
}
