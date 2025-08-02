package keeper_test

import (
	"testing"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/bitbadges/bitbadgeschain/x/badges/testutil/keeper"
)

func TestParamsQuery(t *testing.T) {
	keeper, ctx := keepertest.BadgesKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	params := types.DefaultParams()
	keeper.SetParams(ctx, params)

	response, err := keeper.Params(wctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, params.AllowedDenoms, response.Params.AllowedDenoms)
	// AffiliatePercentage should be zero (default value)
	require.Equal(t, uint64(0), response.Params.AffiliatePercentage.Uint64())
}
