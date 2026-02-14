package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/bitbadges/bitbadgeschain/x/tokenization/testutil/keeper"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.TokenizationKeeper(t)
	params := types.DefaultParams()

	require.NoError(t, k.SetParams(ctx, params))
	retrievedParams := k.GetParams(ctx)
	require.EqualValues(t, params.AllowedDenoms, retrievedParams.AllowedDenoms)
	// AffiliatePercentage should be zero (default value)
	require.NotEqual(t, uint64(0), retrievedParams.AffiliatePercentage.Uint64())
}
