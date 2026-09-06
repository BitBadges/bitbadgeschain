package poolmanager_test

import (
	"github.com/bitbadges/bitbadgeschain/x/poolmanager/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *KeeperTestSuite) TestPoolRoutesFollowDiscardedContext() {
	s.SetupTest()
	const poolID uint64 = 99999
	child, _ := s.defaultCtx().CacheContext()
	s.App.PoolManagerKeeper.SetPoolRoute(child, poolID, types.Balancer)
	_, err := s.App.PoolManagerKeeper.GetPoolModule(child, poolID)
	s.Require().NoError(err)

	for _, lookup := range []func(sdk.Context) error{
		func(ctx sdk.Context) error { _, err := s.App.PoolManagerKeeper.GetPoolModule(ctx, poolID); return err },
		func(ctx sdk.Context) error { _, err := s.App.PoolManagerKeeper.GetPoolType(ctx, poolID); return err },
	} {
		ctx := s.defaultCtx()
		err = lookup(ctx)
		gas := ctx.GasMeter().GasConsumed()
		s.Require().Error(err, "discarded route must not be visible in the parent")
		s.App.PoolManagerKeeper.ResetCaches()
		freshCtx := s.defaultCtx()
		freshErr := lookup(freshCtx)
		s.Require().EqualError(freshErr, err.Error())
		s.Require().Equal(gas, freshCtx.GasMeter().GasConsumed())
	}
}
