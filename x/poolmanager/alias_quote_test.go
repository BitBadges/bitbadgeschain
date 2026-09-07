package poolmanager_test

import (
	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	"github.com/bitbadges/bitbadgeschain/x/gamm/poolmodels/balancer"
	"github.com/bitbadges/bitbadgeschain/x/poolmanager/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *KeeperTestSuite) TestAliasQuoteMatchesExemptFee() {
	s.SetupTest()
	alias := "badgeslp:1:foo"
	pool, err := balancer.NewBalancerPool(77, balancer.PoolParams{SwapFee: osmomath.ZeroDec(), ExitFee: osmomath.ZeroDec()}, []balancer.PoolAsset{
		{Token: sdk.NewInt64Coin(alias, 1_000_000), Weight: osmomath.OneInt()},
		{Token: sdk.NewInt64Coin(FOO, 1_000_000), Weight: osmomath.OneInt()},
	}, s.Ctx.BlockTime())
	s.Require().NoError(err)
	s.Require().NoError(s.App.GammKeeper.InitializePool(s.Ctx, &pool, s.TestAccs[0]))
	k := s.App.PoolManagerKeeper
	k.SetPoolRoute(s.Ctx, 77, types.Balancer)
	k.SetDenomPairTakerFee(s.Ctx, alias, FOO, osmomath.MustNewDecFromStr("0.01"))
	output := sdk.NewInt64Coin(FOO, 10000)
	raw, err := s.App.GammKeeper.CalcInAmtGivenOut(s.Ctx, &pool, output, alias, osmomath.ZeroDec())
	s.Require().NoError(err)
	chargedInput, fee, err := k.ChargeTakerFee(s.Ctx, raw, FOO, s.TestAccs[0], false)
	s.Require().NoError(err)
	s.Require().True(fee.IsZero())
	quote, err := k.MultihopEstimateInGivenExactAmountOut(s.Ctx, []types.SwapAmountOutRoute{{PoolId: 77, TokenInDenom: alias}}, output)
	s.Require().NoError(err)
	s.T().Logf("actual required alias input=%s, quote=%s", chargedInput.Amount, quote)
	s.Require().Equal(chargedInput.Amount.String(), quote.String(), "exact-out quote must match exempt fee execution")
}
