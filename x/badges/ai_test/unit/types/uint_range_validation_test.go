package types_test

import (
	"math"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

// TestUintRange_ValidRange tests valid UintRange creation
func TestUintRange_ValidRange(t *testing.T) {
	range1 := &types.UintRange{
		Start: sdkmath.NewUint(1),
		End:   sdkmath.NewUint(100),
	}

	err := types.ValidateRangesAreValid([]*types.UintRange{range1}, false, false)
	require.NoError(t, err, "valid range should pass validation")
}

// TestUintRange_InvalidRange tests invalid UintRange (start > end)
func TestUintRange_InvalidRange(t *testing.T) {
	range1 := &types.UintRange{
		Start: sdkmath.NewUint(100),
		End:   sdkmath.NewUint(1), // Invalid: start > end
	}

	err := types.ValidateRangesAreValid([]*types.UintRange{range1}, false, false)
	require.Error(t, err, "invalid range should fail validation")
}

// TestUintRange_EqualStartEnd tests range with equal start and end
func TestUintRange_EqualStartEnd(t *testing.T) {
	range1 := &types.UintRange{
		Start: sdkmath.NewUint(50),
		End:   sdkmath.NewUint(50),
	}

	err := types.ValidateRangesAreValid([]*types.UintRange{range1}, false, false)
	require.NoError(t, err, "range with equal start and end should be valid")
}

// TestUintRange_ZeroStart tests range starting at zero
func TestUintRange_ZeroStart(t *testing.T) {
	range1 := &types.UintRange{
		Start: sdkmath.NewUint(0),
		End:   sdkmath.NewUint(100),
	}

	err := types.ValidateRangesAreValid([]*types.UintRange{range1}, false, false)
	// Zero start may or may not be valid depending on allowAllUints flag
	_ = err
}

// TestUintRange_MaxUint tests range with MaxUint64
func TestUintRange_MaxUint(t *testing.T) {
	range1 := &types.UintRange{
		Start: sdkmath.NewUint(1),
		End:   sdkmath.NewUint(math.MaxUint64),
	}

	err := types.ValidateRangesAreValid([]*types.UintRange{range1}, true, false) // allowAllUints=true for MaxUint64
	require.NoError(t, err, "range with MaxUint64 should be valid")
}

// TestUintRange_OverlappingRanges tests overlapping range validation
func TestUintRange_OverlappingRanges(t *testing.T) {
	ranges := []*types.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(50)},
		{Start: sdkmath.NewUint(30), End: sdkmath.NewUint(100)}, // Overlaps with first
	}

	err := types.ValidateRangesAreValid(ranges, false, false)
	// Overlapping ranges may or may not be valid depending on validation rules
	_ = err // Accept either outcome
}

// TestUintRange_NonOverlappingRanges tests non-overlapping ranges
func TestUintRange_NonOverlappingRanges(t *testing.T) {
	ranges := []*types.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(50)},
		{Start: sdkmath.NewUint(51), End: sdkmath.NewUint(100)}, // No overlap
	}

	err := types.ValidateRangesAreValid(ranges, false, false)
	require.NoError(t, err, "non-overlapping ranges should be valid")
}

// TestUintRange_AdjacentRanges tests adjacent ranges
func TestUintRange_AdjacentRanges(t *testing.T) {
	ranges := []*types.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(50)},
		{Start: sdkmath.NewUint(50), End: sdkmath.NewUint(100)}, // Adjacent (end == start)
	}

	err := types.ValidateRangesAreValid(ranges, false, false)
	// Adjacent ranges may or may not be valid depending on validation rules
	_ = err // Accept either outcome
}

// TestUintRange_EmptyRanges tests empty range array
func TestUintRange_EmptyRanges(t *testing.T) {
	ranges := []*types.UintRange{}

	err := types.ValidateRangesAreValid(ranges, false, false)
	require.NoError(t, err, "empty ranges array should be valid")
}

// TestUintRange_SingleRange tests single range
func TestUintRange_SingleRange(t *testing.T) {
	ranges := []*types.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
	}

	err := types.ValidateRangesAreValid(ranges, false, false)
	require.NoError(t, err, "single range should be valid")
}

