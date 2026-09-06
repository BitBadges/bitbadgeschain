package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	"github.com/stretchr/testify/require"
)

func TestUintRangeOrderingAndMerging(t *testing.T) {
	r := func(start, end uint64) *types.UintRange {
		return &types.UintRange{Start: sdkmath.NewUint(start), End: sdkmath.NewUint(end)}
	}
	first, equal := r(1, 3), r(1, 3)
	ranges := []*types.UintRange{r(8, 9), r(1, 5), first, equal, r(6, 7)}
	original := append([]*types.UintRange(nil), ranges...)
	require.True(t, types.DoRangesOverlap(ranges))
	for i := range ranges {
		require.Same(t, original[i], ranges[i], "overlap checks preserve input order")
	}
	merged, err := types.SortUintRangesAndMerge(ranges, true)
	require.NoError(t, err)
	require.Equal(t, []*types.UintRange{r(1, 9)}, merged)
	require.Same(t, first, ranges[0])
	require.Same(t, equal, ranges[1])
	_, err = types.SortUintRangesAndMerge(ranges, false)
	require.Error(t, err)
	disjoint := []*types.UintRange{r(8, 9), r(1, 2), r(3, 4)}
	require.False(t, types.DoRangesOverlap(disjoint))
	merged, err = types.SortUintRangesAndMerge(disjoint, false)
	require.NoError(t, err)
	require.Equal(t, []*types.UintRange{r(1, 4), r(8, 9)}, merged)
}

func BenchmarkUintRangeOrdering(b *testing.B) {
	ranges := make([]*types.UintRange, 1000)
	for i := range ranges {
		ranges[i] = types.CreateUintRange(sdkmath.NewUint(uint64(2*i+1)), sdkmath.NewUint(uint64(2*i+1)))
	}
	b.Run("overlap", func(b *testing.B) {
		for b.Loop() {
			types.DoRangesOverlap(ranges)
		}
	})
	b.Run("merge", func(b *testing.B) {
		for b.Loop() {
			_, _ = types.SortUintRangesAndMerge(ranges, false)
		}
	})
}
