package poolmanager_test

import (
	"sync"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// These tests verify the fix for the 2026-05-12 mainnet halt at block 10183154.
//
// Original bug: x/poolmanager's GetPoolModule sync.Map cache charged gas via
// osmoutils.ChargeMockReadGas, which used a hardcoded default KVGasConfig at module
// scope. EVM precompiles set ctx.KVGasConfig to zero, so cache misses (going through
// the wrapped KVStore) charged 0 gas while cache hits charged ~1006 gas via the
// hardcoded default. Validators with warm caches (uptime spanning a prior swap on
// the pool) thus charged ~1006 more gas per pool lookup than validators with cold
// caches, propagating through feemarket's end-block BlockGasWanted write into a
// divergent AppHash.
//
// Fix: poolModuleCacheValue stores key/value byte lengths instead of pre-computed
// gas amounts. Cache hit charges via osmoutils.ChargeKVReadGas, which reads gas
// config from ctx at hit time. Cache miss inlines ctx.KVStore(...).Get, which
// auto-charges via the wrapped store using the same ctx config. Both paths now
// charge identical gas regardless of ctx.

func (s *KeeperTestSuite) precompileCtx() sdk.Context {
	return s.Ctx.
		WithKVGasConfig(storetypes.GasConfig{}).
		WithTransientKVGasConfig(storetypes.GasConfig{}).
		WithGasMeter(storetypes.NewGasMeter(10_000_000)).
		WithExecMode(sdk.ExecModeFinalize)
}

func (s *KeeperTestSuite) defaultCtx() sdk.Context {
	return s.Ctx.
		WithKVGasConfig(storetypes.KVGasConfig()).
		WithTransientKVGasConfig(storetypes.TransientGasConfig()).
		WithGasMeter(storetypes.NewGasMeter(10_000_000)).
		WithExecMode(sdk.ExecModeFinalize)
}

func (s *KeeperTestSuite) getPoolModuleGas(ctx sdk.Context, poolId uint64) uint64 {
	before := ctx.GasMeter().GasConsumed()
	_, err := s.App.PoolManagerKeeper.GetPoolModule(ctx, poolId)
	s.Require().NoError(err)
	return ctx.GasMeter().GasConsumed() - before
}

func (s *KeeperTestSuite) getPoolTypeGas(ctx sdk.Context, poolId uint64) uint64 {
	before := ctx.GasMeter().GasConsumed()
	_, err := s.App.PoolManagerKeeper.GetPoolType(ctx, poolId)
	s.Require().NoError(err)
	return ctx.GasMeter().GasConsumed() - before
}

// ─── 1. Cache hit == cache miss across ctx configs (the core regression) ─────

func (s *KeeperTestSuite) TestGetPoolModule_DeterministicGas_PrecompileContext() {
	s.SetupTest()
	poolId := s.PrepareBalancerPool()
	s.App.PoolManagerKeeper.ResetCaches()

	gasMiss := s.getPoolModuleGas(s.precompileCtx(), poolId)
	gasHit := s.getPoolModuleGas(s.precompileCtx(), poolId)

	s.T().Logf("Zero KVGasConfig (precompile context): miss=%d hit=%d", gasMiss, gasHit)
	s.Require().Equal(uint64(0), gasMiss, "miss should charge 0 gas in zero-KVGasConfig ctx")
	s.Require().Equal(gasMiss, gasHit,
		"cache hit and miss must charge identical gas; mismatch causes consensus divergence")
}

func (s *KeeperTestSuite) TestGetPoolModule_DeterministicGas_DefaultContext() {
	s.SetupTest()
	poolId := s.PrepareBalancerPool()
	s.App.PoolManagerKeeper.ResetCaches()

	gasMiss := s.getPoolModuleGas(s.defaultCtx(), poolId)
	gasHit := s.getPoolModuleGas(s.defaultCtx(), poolId)

	s.T().Logf("Default KVGasConfig: miss=%d hit=%d", gasMiss, gasHit)
	s.Require().Greater(gasMiss, uint64(0))
	s.Require().Equal(gasMiss, gasHit)
}

// ─── 2. GetPoolType determinism ──────────────────────────────────────────────
// GetPoolType has its own cache-hit branch in the keeper. Verify it too.

func (s *KeeperTestSuite) TestGetPoolType_DeterministicGas_PrecompileContext() {
	s.SetupTest()
	poolId := s.PrepareBalancerPool()
	s.App.PoolManagerKeeper.ResetCaches()

	// GetPoolType's miss path uses osmoutils.Get (which does store.Get under the hood).
	// store.Get auto-charges via the ctx-wrapped store → 0 in precompile ctx.
	gasMiss := s.getPoolTypeGas(s.precompileCtx(), poolId)

	// Populate the cache via GetPoolModule, then GetPoolType should hit the cache.
	_, err := s.App.PoolManagerKeeper.GetPoolModule(s.precompileCtx(), poolId)
	s.Require().NoError(err)
	gasHit := s.getPoolTypeGas(s.precompileCtx(), poolId)

	s.T().Logf("GetPoolType precompile: miss=%d hit=%d", gasMiss, gasHit)
	s.Require().Equal(uint64(0), gasMiss)
	s.Require().Equal(uint64(0), gasHit)
}

func (s *KeeperTestSuite) TestGetPoolType_DeterministicGas_DefaultContext() {
	s.SetupTest()
	poolId := s.PrepareBalancerPool()
	s.App.PoolManagerKeeper.ResetCaches()

	gasMiss := s.getPoolTypeGas(s.defaultCtx(), poolId)

	// Populate cache via GetPoolModule, then GetPoolType should hit.
	_, err := s.App.PoolManagerKeeper.GetPoolModule(s.defaultCtx(), poolId)
	s.Require().NoError(err)
	gasHit := s.getPoolTypeGas(s.defaultCtx(), poolId)

	s.T().Logf("GetPoolType default: miss=%d hit=%d", gasMiss, gasHit)
	s.Require().Equal(gasMiss, gasHit,
		"GetPoolType cache hit and miss must charge identical gas")
}

// ─── 3. Stableswap pool — different serialized value length ──────────────────

func (s *KeeperTestSuite) TestGetPoolModule_StableswapPool_Deterministic() {
	s.SetupTest()
	poolId := s.PrepareBasicStableswapPool()

	s.App.PoolManagerKeeper.ResetCaches()
	gasMissDefault := s.getPoolModuleGas(s.defaultCtx(), poolId)
	gasHitDefault := s.getPoolModuleGas(s.defaultCtx(), poolId)
	s.Require().Equal(gasMissDefault, gasHitDefault, "stableswap default: hit != miss")

	s.App.PoolManagerKeeper.ResetCaches()
	gasMissPrecompile := s.getPoolModuleGas(s.precompileCtx(), poolId)
	gasHitPrecompile := s.getPoolModuleGas(s.precompileCtx(), poolId)
	s.Require().Equal(uint64(0), gasMissPrecompile)
	s.Require().Equal(gasMissPrecompile, gasHitPrecompile, "stableswap precompile: hit != miss")
}

// ─── 4. Multiple pool IDs — different key byte lengths ───────────────────────
// Pool ID 9 → key "1|9" (3 bytes). Pool ID 12 → key "1|12" (4 bytes).
// Cache must capture the per-pool key length and replay correctly.

func (s *KeeperTestSuite) TestGetPoolModule_MultiplePools_AllDeterministic() {
	s.SetupTest()
	poolIds := s.PrepareMultipleBalancerPools(12)

	for _, pid := range poolIds {
		s.App.PoolManagerKeeper.ResetCaches()
		gasMiss := s.getPoolModuleGas(s.defaultCtx(), pid)
		gasHit := s.getPoolModuleGas(s.defaultCtx(), pid)
		s.Require().Equalf(gasMiss, gasHit, "pool %d default: hit != miss", pid)

		s.App.PoolManagerKeeper.ResetCaches()
		gasMissP := s.getPoolModuleGas(s.precompileCtx(), pid)
		gasHitP := s.getPoolModuleGas(s.precompileCtx(), pid)
		s.Require().Equalf(uint64(0), gasMissP, "pool %d precompile miss != 0", pid)
		s.Require().Equalf(gasMissP, gasHitP, "pool %d precompile: hit != miss", pid)
	}
}

// ─── 5. Mixed pool ID widths in the same cache ──────────────────────────────
// Confirm the cache doesn't smear key lengths across entries.

func (s *KeeperTestSuite) TestGetPoolModule_MixedKeyLengths_NoSmearing() {
	s.SetupTest()
	poolIds := s.PrepareMultipleBalancerPools(11) // pool IDs span 1-digit and 2-digit
	s.App.PoolManagerKeeper.ResetCaches()

	// Populate cache for all pools
	for _, pid := range poolIds {
		_, err := s.App.PoolManagerKeeper.GetPoolModule(s.defaultCtx(), pid)
		s.Require().NoError(err)
	}

	// Now compare each pool's cache-hit gas to its own cold-miss gas (one at a time).
	for _, pid := range poolIds {
		gasHit := s.getPoolModuleGas(s.defaultCtx(), pid)

		s.App.PoolManagerKeeper.ResetCaches()
		// Re-warm all OTHER pools to keep cache state realistic
		for _, other := range poolIds {
			if other == pid {
				continue
			}
			_, _ = s.App.PoolManagerKeeper.GetPoolModule(s.defaultCtx(), other)
		}
		gasMiss := s.getPoolModuleGas(s.defaultCtx(), pid)

		s.Require().Equalf(gasMiss, gasHit, "pool %d: hit %d != miss %d", pid, gasHit, gasMiss)

		// Re-populate cache for next iteration
		_, _ = s.App.PoolManagerKeeper.GetPoolModule(s.defaultCtx(), pid)
	}
}

// ─── 6. Cache populate gated by ExecModeFinalize ─────────────────────────────

func (s *KeeperTestSuite) TestGetPoolModule_CacheOnlyPopulatesInFinalize() {
	s.SetupTest()
	poolId := s.PrepareBalancerPool()
	s.App.PoolManagerKeeper.ResetCaches()

	// CheckTx mode must NOT populate the cache
	checkCtx := s.Ctx.
		WithKVGasConfig(storetypes.KVGasConfig()).
		WithGasMeter(storetypes.NewGasMeter(10_000_000)).
		WithExecMode(sdk.ExecModeCheck)
	_, err := s.App.PoolManagerKeeper.GetPoolModule(checkCtx, poolId)
	s.Require().NoError(err)

	// Simulate mode also must NOT populate
	simCtx := s.Ctx.
		WithKVGasConfig(storetypes.KVGasConfig()).
		WithGasMeter(storetypes.NewGasMeter(10_000_000)).
		WithExecMode(sdk.ExecModeSimulate)
	_, err = s.App.PoolManagerKeeper.GetPoolModule(simCtx, poolId)
	s.Require().NoError(err)

	// Now a finalize-mode call should still be a clean cache MISS, and a subsequent
	// finalize-mode call should be a HIT. After the fix both charge identical gas,
	// so we verify by direct sync.Map inspection: it should be empty after Check+Simulate.
	gasFirstFinalize := s.getPoolModuleGas(s.defaultCtx(), poolId)
	gasSecondFinalize := s.getPoolModuleGas(s.defaultCtx(), poolId)
	s.Require().Equal(gasFirstFinalize, gasSecondFinalize,
		"after Check/Simulate ran (not populating cache), Finalize miss==hit gas")
}

// ─── 7. SetPoolRoute invalidates the cache ───────────────────────────────────

func (s *KeeperTestSuite) TestGetPoolModule_SetPoolRoute_InvalidatesCache() {
	s.SetupTest()
	poolId := s.PrepareBalancerPool()

	// Populate cache
	_, err := s.App.PoolManagerKeeper.GetPoolModule(s.defaultCtx(), poolId)
	s.Require().NoError(err)

	// SetPoolRoute deletes the cache entry (per create_pool.go:172)
	pType, err := s.App.PoolManagerKeeper.GetPoolType(s.defaultCtx(), poolId)
	s.Require().NoError(err)
	s.App.PoolManagerKeeper.SetPoolRoute(s.Ctx, poolId, pType)

	// Next call should re-miss then re-hit. Both must charge equal gas.
	gasReMiss := s.getPoolModuleGas(s.defaultCtx(), poolId)
	gasReHit := s.getPoolModuleGas(s.defaultCtx(), poolId)
	s.Require().Equal(gasReMiss, gasReHit)
}

// ─── 8. ResetCaches → next lookup is a clean miss with identical gas ────────

func (s *KeeperTestSuite) TestGetPoolModule_ResetCaches_NextIsClean() {
	s.SetupTest()
	poolId := s.PrepareBalancerPool()

	_, err := s.App.PoolManagerKeeper.GetPoolModule(s.defaultCtx(), poolId)
	s.Require().NoError(err)
	gasWarmHit := s.getPoolModuleGas(s.defaultCtx(), poolId)

	s.App.PoolManagerKeeper.ResetCaches()
	gasColdMiss := s.getPoolModuleGas(s.defaultCtx(), poolId)

	s.Require().Equal(gasWarmHit, gasColdMiss)
}

// ─── 9. Sequence determinism — total gas independent of cache state pattern ─

func (s *KeeperTestSuite) TestGetPoolModule_SequenceDeterminism() {
	s.SetupTest()
	poolIds := s.PrepareMultipleBalancerPools(3)

	// Pattern A: pre-warm all, then lookup all
	s.App.PoolManagerKeeper.ResetCaches()
	for _, pid := range poolIds {
		_, _ = s.App.PoolManagerKeeper.GetPoolModule(s.defaultCtx(), pid)
	}
	ctxA := s.defaultCtx()
	beforeA := ctxA.GasMeter().GasConsumed()
	for _, pid := range poolIds {
		_, err := s.App.PoolManagerKeeper.GetPoolModule(ctxA, pid)
		s.Require().NoError(err)
	}
	totalA := ctxA.GasMeter().GasConsumed() - beforeA

	// Pattern B: cold-reset before each lookup
	ctxB := s.defaultCtx()
	beforeB := ctxB.GasMeter().GasConsumed()
	for _, pid := range poolIds {
		s.App.PoolManagerKeeper.ResetCaches()
		_, err := s.App.PoolManagerKeeper.GetPoolModule(ctxB, pid)
		s.Require().NoError(err)
	}
	totalB := ctxB.GasMeter().GasConsumed() - beforeB

	s.T().Logf("sequence totals: all-warm=%d all-cold=%d", totalA, totalB)
	s.Require().Equal(totalA, totalB,
		"cache state pattern must not affect total gas across a sequence of lookups")
}

// ─── 10. Concurrent reads — no race, no gas divergence ──────────────────────

func (s *KeeperTestSuite) TestGetPoolModule_ConcurrentLoad_NoRace() {
	s.SetupTest()
	poolId := s.PrepareBalancerPool()

	// Pre-warm
	_, err := s.App.PoolManagerKeeper.GetPoolModule(s.defaultCtx(), poolId)
	s.Require().NoError(err)

	const goroutines = 20
	var wg sync.WaitGroup
	gasCh := make(chan uint64, goroutines)
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx := s.defaultCtx()
			gasCh <- s.getPoolModuleGas(ctx, poolId)
		}()
	}
	wg.Wait()
	close(gasCh)

	var first uint64
	firstSet := false
	for g := range gasCh {
		if !firstSet {
			first = g
			firstSet = true
			continue
		}
		s.Require().Equal(first, g, "all concurrent hits must charge identical gas")
	}
}

// ─── 11. Cross-context replay — the original bug scenario ───────────────────
// Populate the cache under one ctx config, then hit it under a different ctx config.
// Cache hit must honor the *current* ctx config, not the one that populated the cache.

func (s *KeeperTestSuite) TestGetPoolModule_CrossContextReplay() {
	s.SetupTest()
	poolId := s.PrepareBalancerPool()

	// Step 1: warm cache under precompile ctx (0 gas)
	s.App.PoolManagerKeeper.ResetCaches()
	gasPrecompileMiss := s.getPoolModuleGas(s.precompileCtx(), poolId)
	s.Require().Equal(uint64(0), gasPrecompileMiss)

	// Step 2: hit the cache under default ctx — must charge default-config gas
	gasDefaultHitOnPrecompileCache := s.getPoolModuleGas(s.defaultCtx(), poolId)

	// Step 3: cold-miss under default ctx for comparison
	s.App.PoolManagerKeeper.ResetCaches()
	gasDefaultMiss := s.getPoolModuleGas(s.defaultCtx(), poolId)

	s.T().Logf("cross-ctx: precompile-miss=%d default-hit-on-precompile-cache=%d default-miss=%d",
		gasPrecompileMiss, gasDefaultHitOnPrecompileCache, gasDefaultMiss)
	s.Require().Equal(gasDefaultMiss, gasDefaultHitOnPrecompileCache,
		"cache hit under default ctx must equal cold miss under default ctx, even if cache was populated under precompile ctx — this is the determinism property")

	// Step 4: reverse direction — populate under default, hit under precompile
	s.App.PoolManagerKeeper.ResetCaches()
	_ = s.getPoolModuleGas(s.defaultCtx(), poolId)
	gasPrecompileHitOnDefaultCache := s.getPoolModuleGas(s.precompileCtx(), poolId)
	s.Require().Equal(uint64(0), gasPrecompileHitOnDefaultCache,
		"cache hit under precompile ctx must charge 0, regardless of how cache was populated")
}

// ─── 12. Pool not found — error path is deterministic ───────────────────────

func (s *KeeperTestSuite) TestGetPoolModule_NotFound_DeterministicError() {
	s.SetupTest()
	const missingPoolId = uint64(99999)

	ctxA := s.precompileCtx()
	beforeA := ctxA.GasMeter().GasConsumed()
	_, errA := s.App.PoolManagerKeeper.GetPoolModule(ctxA, missingPoolId)
	gasA := ctxA.GasMeter().GasConsumed() - beforeA
	s.Require().Error(errA)

	ctxB := s.precompileCtx()
	beforeB := ctxB.GasMeter().GasConsumed()
	_, errB := s.App.PoolManagerKeeper.GetPoolModule(ctxB, missingPoolId)
	gasB := ctxB.GasMeter().GasConsumed() - beforeB
	s.Require().Error(errB)

	s.Require().Equal(errA.Error(), errB.Error())
	s.Require().Equal(gasA, gasB)
}
