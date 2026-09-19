package lifecycle

import (
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func bal(amount, start, end, ts, te uint64) *types.Balance {
	return &types.Balance{Amount: sdkmath.NewUint(amount), TokenIds: []*types.UintRange{{Start: sdkmath.NewUint(start), End: sdkmath.NewUint(end)}}, OwnershipTimes: []*types.UintRange{{Start: sdkmath.NewUint(ts), End: sdkmath.NewUint(te)}}}
}
func TestExactBalanceRanges(t *testing.T) {
	a := []*types.Balance{bal(2, 1, 2, 100, 200)}
	equal, err := equalBalances(a, []*types.Balance{bal(2, 1, 1, 100, 200), bal(2, 2, 2, 100, 200)})
	require.NoError(t, err)
	require.True(t, equal)
	for _, other := range [][]*types.Balance{{bal(2, 1, 2, 100, 199)}, {bal(2, 1, 2, 101, 200)}, {bal(2, 1, 3, 100, 200)}, {bal(1, 1, 2, 100, 200)}, {bal(2, 1, 2, 100, 200), bal(1, 2, 2, 150, 160)}} {
		equal, err = equalBalances(a, other)
		require.NoError(t, err)
		require.False(t, equal)
	}
	equal, err = equalBalances(a, []*types.Balance{bal(1, 1, 2, 100, 200), bal(1, 1, 2, 100, 200)})
	require.NoError(t, err)
	require.True(t, equal)
	invalid := bal(1, 1, 2, 100, 200)
	invalid.TokenIds = append(invalid.TokenIds, &types.UintRange{Start: sdkmath.NewUint(2), End: sdkmath.NewUint(3)})
	_, err = equalBalances(a, []*types.Balance{invalid})
	require.ErrorContains(t, err, "overlapping")
}
