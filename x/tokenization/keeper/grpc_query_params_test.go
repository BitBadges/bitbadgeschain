package keeper_test

import (
	"testing"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/bitbadges/bitbadgeschain/x/tokenization/testutil/keeper"
)

func TestParamsQuery(t *testing.T) {
	keeper, ctx := keepertest.TokenizationKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	params := types.DefaultParams()
	keeper.SetParams(ctx, params)

	response, err := keeper.Params(wctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, params.AllowedDenoms, response.Params.AllowedDenoms)
	// AffiliatePercentage should be 500 (default value)
	require.Equal(t, uint64(500), response.Params.AffiliatePercentage.Uint64())
}
