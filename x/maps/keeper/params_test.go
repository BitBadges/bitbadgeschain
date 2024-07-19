package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "bitbadgeschain/testutil/keeper"
	"bitbadgeschain/x/maps/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.MapsKeeper(t)
	params := types.DefaultParams()

	require.NoError(t, k.SetParams(ctx, params))
	require.EqualValues(t, params, k.GetParams(ctx))
}
