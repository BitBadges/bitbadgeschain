package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bitbadges/bitbadgeschain/x/wasmx/types"

	keepertest "github.com/bitbadges/bitbadgeschain/testutil/keeper"
)

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.WasmxKeeper(t)
	params := types.DefaultParams()

	require.NoError(t, k.SetParams(ctx, params))
	require.EqualValues(t, params, k.GetParams(ctx))
}
