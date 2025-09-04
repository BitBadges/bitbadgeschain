package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	"github.com/bitbadges/bitbadgeschain/app/apptesting"
	"github.com/bitbadges/bitbadgeschain/x/gamm/poolmodels/balancer"
	"github.com/bitbadges/bitbadgeschain/x/gamm/poolmodels/stableswap"
	"github.com/bitbadges/bitbadgeschain/x/gamm/types"
)

type KeeperTestSuite struct {
	apptesting.KeeperTestHelper

	queryClient types.QueryClient
}

var (
	defaultSpreadFactor = osmomath.MustNewDecFromStr("0.025")
	defaultZeroExitFee  = osmomath.ZeroDec()
)

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	s.Reset()
	s.queryClient = types.NewQueryClient(s.QueryHelper)
}

func (s *KeeperTestSuite) prepareCustomBalancerPool(
	balances sdk.Coins,
	poolAssets []balancer.PoolAsset,
	poolParams balancer.PoolParams,
) uint64 {
	s.fundAllAccountsWith(balances)

	poolID, err := s.App.PoolManagerKeeper.CreatePool(
		s.Ctx,
		balancer.NewMsgCreateBalancerPool(s.TestAccs[0], poolParams, poolAssets),
	)
	s.Require().NoError(err)

	return poolID
}

func (s *KeeperTestSuite) prepareCustomStableswapPool(
	balances sdk.Coins,
	poolParams stableswap.PoolParams,
	initialLiquidity sdk.Coins,
	scalingFactors []uint64,
) uint64 {
	s.fundAllAccountsWith(balances)

	poolID, err := s.App.PoolManagerKeeper.CreatePool(
		s.Ctx,
		stableswap.NewMsgCreateStableswapPool(s.TestAccs[0], poolParams, initialLiquidity, scalingFactors, ""),
	)
	s.Require().NoError(err)

	return poolID
}

func (s *KeeperTestSuite) fundAllAccountsWith(balances sdk.Coins) {
	for _, acc := range s.TestAccs {
		s.FundAcc(acc, balances)
	}
}
