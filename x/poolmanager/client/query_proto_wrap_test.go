package client_test

import (
	"github.com/bitbadges/bitbadgeschain/third_party/apptesting"
	"github.com/bitbadges/bitbadgeschain/x/poolmanager/client"
	"github.com/bitbadges/bitbadgeschain/x/poolmanager/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestPrimitiveRouteLengths(t *testing.T) {
	q := client.NewQuerier(nil)
	for _, denoms := range [][]string{nil, {"uatom", "uosmo"}} {
		require.NotPanics(t, func() {
			_, err := q.EstimateSwapExactAmountInWithPrimitiveTypes(sdk.Context{}, types.EstimateSwapExactAmountInWithPrimitiveTypesRequest{TokenIn: "10uatom", RoutesPoolId: []uint64{1}, RoutesTokenOutDenom: denoms})
			require.Equal(t, codes.InvalidArgument, status.Code(err))
		})
		require.NotPanics(t, func() {
			_, err := q.EstimateSwapExactAmountOutWithPrimitiveTypes(sdk.Context{}, types.EstimateSwapExactAmountOutWithPrimitiveTypesRequest{TokenOut: "10uatom", RoutesPoolId: []uint64{1}, RoutesTokenInDenom: denoms})
			require.Equal(t, codes.InvalidArgument, status.Code(err))
		})
	}
}

func TestPrimitiveExactOutMatchesStructuredRoutes(t *testing.T) {
	s := new(apptesting.KeeperTestHelper)
	s.SetT(t)
	s.Setup()
	pool := s.PrepareBalancerPoolWithCoins(sdk.NewInt64Coin("uatom", 1000000), sdk.NewInt64Coin("uosmo", 1000000))
	q := client.NewQuerier(&s.App.PoolManagerKeeper)
	expected, err := q.EstimateSwapExactAmountOut(s.Ctx, types.EstimateSwapExactAmountOutRequest{TokenOut: "10uosmo", Routes: []types.SwapAmountOutRoute{{PoolId: pool, TokenInDenom: "uatom"}}})
	require.NoError(t, err)
	actual, err := q.EstimateSwapExactAmountOutWithPrimitiveTypes(s.Ctx, types.EstimateSwapExactAmountOutWithPrimitiveTypesRequest{TokenOut: "10uosmo", RoutesPoolId: []uint64{pool}, RoutesTokenInDenom: []string{"uatom"}})
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}
