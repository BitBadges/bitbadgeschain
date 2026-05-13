package poolmanager_test

import (
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TestGetPoolModule_DeterministicGas verifies the fix for the chain-halt bug at block 10183154
// (mainnet, 2026-05-12). The original bug: x/poolmanager's GetPoolModule sync.Map cache
// charged gas via osmoutils.ChargeMockReadGas, which used a hardcoded default KVGasConfig
// instead of the ctx's actual config. EVM precompiles set ctx.KVGasConfig to zero, so cache
// misses charged 0 gas (via the wrapped KVStore) while cache hits charged ~1000+ gas (via
// the hardcoded default). Validators with warm caches diverged from validators with cold
// caches → feemarket end-block state diverged → AppHash split → halt.
//
// The fix stores byte lengths in poolModuleCacheValue and uses osmoutils.ChargeKVReadGas
// (which reads gas config from ctx at hit time). Now cache hit and cache miss charge the
// same gas regardless of ctx config.
func (s *KeeperTestSuite) TestGetPoolModule_DeterministicGas_PrecompileContext() {
	s.SetupTest()
	poolId := s.PrepareBalancerPool()
	s.App.PoolManagerKeeper.ResetCaches()

	mkCtx := func() sdk.Context {
		return s.Ctx.
			WithKVGasConfig(storetypes.GasConfig{}).
			WithTransientKVGasConfig(storetypes.GasConfig{}).
			WithGasMeter(storetypes.NewGasMeter(10_000_000)).
			WithExecMode(sdk.ExecModeFinalize)
	}

	// Cache MISS: store.Get auto-charges via ctx's wrapped KVStore (zero config → 0 gas)
	missCtx := mkCtx()
	gasBefore := missCtx.GasMeter().GasConsumed()
	_, err := s.App.PoolManagerKeeper.GetPoolModule(missCtx, poolId)
	s.Require().NoError(err)
	gasMiss := missCtx.GasMeter().GasConsumed() - gasBefore

	// Cache HIT: ChargeKVReadGas honors ctx.KVGasConfig (also zero → 0 gas)
	hitCtx := mkCtx()
	gasBefore = hitCtx.GasMeter().GasConsumed()
	_, err = s.App.PoolManagerKeeper.GetPoolModule(hitCtx, poolId)
	s.Require().NoError(err)
	gasHit := hitCtx.GasMeter().GasConsumed() - gasBefore

	s.T().Logf("Zero KVGasConfig (precompile context): miss=%d hit=%d", gasMiss, gasHit)
	s.Require().Equal(gasMiss, gasHit,
		"cache hit and miss must charge identical gas; mismatch causes consensus divergence")
}

func (s *KeeperTestSuite) TestGetPoolModule_DeterministicGas_DefaultContext() {
	s.SetupTest()
	poolId := s.PrepareBalancerPool()
	s.App.PoolManagerKeeper.ResetCaches()

	mkCtx := func() sdk.Context {
		return s.Ctx.
			WithKVGasConfig(storetypes.KVGasConfig()).
			WithTransientKVGasConfig(storetypes.TransientGasConfig()).
			WithGasMeter(storetypes.NewGasMeter(10_000_000)).
			WithExecMode(sdk.ExecModeFinalize)
	}

	missCtx := mkCtx()
	gasBefore := missCtx.GasMeter().GasConsumed()
	_, err := s.App.PoolManagerKeeper.GetPoolModule(missCtx, poolId)
	s.Require().NoError(err)
	gasMiss := missCtx.GasMeter().GasConsumed() - gasBefore

	hitCtx := mkCtx()
	gasBefore = hitCtx.GasMeter().GasConsumed()
	_, err = s.App.PoolManagerKeeper.GetPoolModule(hitCtx, poolId)
	s.Require().NoError(err)
	gasHit := hitCtx.GasMeter().GasConsumed() - gasBefore

	s.T().Logf("Default KVGasConfig: miss=%d hit=%d", gasMiss, gasHit)
	s.Require().Equal(gasMiss, gasHit,
		"cache hit and miss must charge identical gas in default context as well")
}
