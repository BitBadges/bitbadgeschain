package keeper_test

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/testutil"

	appparams "github.com/bitbadges/bitbadgeschain/app/params"
	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	"github.com/bitbadges/bitbadgeschain/x/gamm/poolmodels/balancer"
	"github.com/bitbadges/bitbadgeschain/x/gamm/types"
)

func (s *KeeperTestSuite) TestGammInitGenesis() {
	s.SetupTest()

	for i := 0; i < 3; i++ {
		s.PrepareBalancerPool()
	}

	pools, err := s.App.GammKeeper.GetPoolsAndPoke(s.Ctx)
	if err != nil {
		panic(err)
	}

	balancerPoolPreInit := pools[0]

	poolAnys := []*codectypes.Any{}
	for _, poolI := range pools {
		any, err := codectypes.NewAnyWithValue(poolI)
		if err != nil {
			panic(err)
		}
		poolAnys = append(poolAnys, any)
	}

	// Reset the testing env so that we can see if the pools get re-initialized via init genesis
	s.SetupTest()

	// Check if the pools were reset
	_, err = s.App.GammKeeper.GetPoolAndPoke(s.Ctx, 1)
	s.Require().Error(err)

	s.App.GammKeeper.InitGenesis(s.Ctx, types.GenesisState{
		Pools:          poolAnys,
		NextPoolNumber: 7,
		Params:         types.Params{},
	}, s.App.AppCodec())

	poolStored, err := s.App.GammKeeper.GetPoolAndPoke(s.Ctx, 1)
	s.Require().NoError(err)
	s.Require().Equal(balancerPoolPreInit.GetId(), poolStored.GetId())
	s.Require().Equal(balancerPoolPreInit.GetAddress(), poolStored.GetAddress())
	s.Require().Equal(balancerPoolPreInit.GetSpreadFactor(s.Ctx), poolStored.GetSpreadFactor(s.Ctx))
	s.Require().Equal(balancerPoolPreInit.GetExitFee(s.Ctx), poolStored.GetExitFee(s.Ctx))
	s.Require().Equal(balancerPoolPreInit.GetTotalShares(), poolStored.GetTotalShares())
	s.Require().Equal(balancerPoolPreInit.String(), poolStored.String())

	_, err = s.App.GammKeeper.GetPoolAndPoke(s.Ctx, 7)
	s.Require().Error(err)

	liquidity, err := s.App.GammKeeper.GetTotalLiquidity(s.Ctx)
	s.Require().NoError(err)
	expectedLiquidity := sdk.NewCoins(sdk.NewInt64Coin("bar", 15000000), sdk.NewInt64Coin("baz", 15000000), sdk.NewInt64Coin("foo", 15000000), sdk.NewInt64Coin(appparams.BaseCoinUnit, 15000000))
	s.Require().Equal(expectedLiquidity.String(), liquidity.String())
}

func (s *KeeperTestSuite) TestGammExportGenesis() {
	s.SetupTest()
	ctx := s.Ctx

	acc1 := s.TestAccs[0]
	err := testutil.FundAccount(ctx, s.App.BankKeeper, acc1, sdk.NewCoins(
		sdk.NewCoin(appparams.BaseCoinUnit, osmomath.NewInt(10000000000)),
		sdk.NewInt64Coin("foo", 100000),
		sdk.NewInt64Coin("bar", 100000),
	))
	s.Require().NoError(err)

	msg := balancer.NewMsgCreateBalancerPool(acc1, balancer.PoolParams{
		SwapFee: osmomath.NewDecWithPrec(1, 2),
		ExitFee: osmomath.ZeroDec(),
	}, []balancer.PoolAsset{{
		Weight: osmomath.NewInt(100),
		Token:  sdk.NewCoin("foo", osmomath.NewInt(10000)),
	}, {
		Weight: osmomath.NewInt(100),
		Token:  sdk.NewCoin("bar", osmomath.NewInt(10000)),
	}})
	_, err = s.App.PoolManagerKeeper.CreatePool(ctx, msg)
	s.Require().NoError(err)

	msg = balancer.NewMsgCreateBalancerPool(acc1, balancer.PoolParams{
		SwapFee: osmomath.NewDecWithPrec(1, 2),
		ExitFee: osmomath.ZeroDec(),
	}, []balancer.PoolAsset{{
		Weight: osmomath.NewInt(70),
		Token:  sdk.NewCoin("foo", osmomath.NewInt(10000)),
	}, {
		Weight: osmomath.NewInt(100),
		Token:  sdk.NewCoin("bar", osmomath.NewInt(10000)),
	}})
	_, err = s.App.PoolManagerKeeper.CreatePool(ctx, msg)
	s.Require().NoError(err)

	genesis := s.App.GammKeeper.ExportGenesis(ctx)
	s.Require().Len(genesis.Pools, 2)
}

func (s *KeeperTestSuite) TestMarshalUnmarshalGenesis() {
	s.SetupTest()
	ctx := s.Ctx

	acc1 := s.TestAccs[0]
	err := testutil.FundAccount(ctx, s.App.BankKeeper, acc1, sdk.NewCoins(
		sdk.NewCoin(appparams.BaseCoinUnit, osmomath.NewInt(10000000000)),
		sdk.NewInt64Coin("foo", 100000),
		sdk.NewInt64Coin("bar", 100000),
	))
	s.Require().NoError(err)

	msg := balancer.NewMsgCreateBalancerPool(acc1, balancer.PoolParams{
		SwapFee: osmomath.NewDecWithPrec(1, 2),
		ExitFee: osmomath.ZeroDec(),
	}, []balancer.PoolAsset{{
		Weight: osmomath.NewInt(100),
		Token:  sdk.NewCoin("foo", osmomath.NewInt(10000)),
	}, {
		Weight: osmomath.NewInt(100),
		Token:  sdk.NewCoin("bar", osmomath.NewInt(10000)),
	}})
	_, err = s.App.PoolManagerKeeper.CreatePool(ctx, msg)
	s.Require().NoError(err)

	genesis := s.App.GammKeeper.ExportGenesis(ctx)
	s.Assert().NotPanics(func() {
		s.App.GammKeeper.InitGenesis(s.Ctx, *genesis, s.App.AppCodec())
	})
}
