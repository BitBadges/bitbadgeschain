package v35

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	gammkeeper "github.com/bitbadges/bitbadgeschain/x/gamm/keeper"
	balancer "github.com/bitbadges/bitbadgeschain/x/gamm/poolmodels/balancer"
	stableswap "github.com/bitbadges/bitbadgeschain/x/gamm/poolmodels/stableswap"
)

// GammMigrationResult reports what RescaleGammPools touched.
type GammMigrationResult struct {
	BalancerPools   int
	StableswapPools int
	Skipped         int
}

// RescaleGammPools converts the legacy-denom amounts recorded inside pool
// objects.
//
// A pool's actual reserves are coins held at the pool's address, so
// RedenominateBank already moved them. What it cannot reach is the pool record's
// own copy of those balances — x/gamm keeps PoolAssets (and, for stableswap,
// PoolLiquidity) in its store and prices swaps off that copy rather than
// re-reading the bank. Leave it stale and every swap prices the pool as if it
// held 10^-9 of its real reserves.
//
// TotalShares deliberately does NOT scale, which is worth stating because the
// intuition runs the other way. LP shares are their own denom (gamm/pool/N),
// untouched by a change to the native denom's precision. Scaling the assets
// while leaving shares alone is exactly what keeps each share worth the same:
//
//	before: 1000 ubadge / 100 shares  = 10 ubadge per share
//	after:  10^12 abadge / 100 shares = 10^10 abadge per share = the same value
//
// Scaling shares too would multiply every LP position's redemption value by
// 10^9. Pool weights are dimensionless ratios and are likewise untouched.
func RescaleGammPools(ctx sdk.Context, gk gammkeeper.Keeper) (GammMigrationResult, error) {
	var res GammMigrationResult

	pools, err := gk.GetPools(ctx)
	if err != nil {
		return res, fmt.Errorf("reading gamm pools: %w", err)
	}

	for _, pool := range pools {
		switch p := pool.(type) {
		case *balancer.Pool:
			for i := range p.PoolAssets {
				if p.PoolAssets[i].Token.Denom == legacyDenom() {
					p.PoolAssets[i].Token = sdk.NewCoin(newDenom(), p.PoolAssets[i].Token.Amount.Mul(ConversionFactor))
				}
			}
			// PoolAssets are documented as sorted by denomination, and abadge
			// sorts where ubadge did not. Re-sorting keeps that invariant.
			sortBalancerAssets(p)
			if err := gk.OverwritePoolV15MigrationUnsafe(ctx, p); err != nil {
				return res, fmt.Errorf("writing balancer pool %d: %w", p.Id, err)
			}
			res.BalancerPools++

		case *stableswap.Pool:
			p.PoolLiquidity = convertCoins(p.PoolLiquidity, legacyDenom(), newDenom())
			if err := gk.OverwritePoolV15MigrationUnsafe(ctx, p); err != nil {
				return res, fmt.Errorf("writing stableswap pool %d: %w", p.Id, err)
			}
			res.StableswapPools++

		default:
			// An unrecognised pool type is not something to shrug at: if it
			// holds legacy-denom reserves it would be left mispriced by a
			// factor of 10^9. Fail rather than silently skip.
			return res, fmt.Errorf("unhandled gamm pool type %T (pool %d): refusing to leave it unconverted", pool, pool.GetId())
		}
	}

	ctx.Logger().Info(
		"v35: rescaled gamm pools",
		"factor", ConversionFactor.String(),
		"balancer", res.BalancerPools,
		"stableswap", res.StableswapPools,
	)

	return res, nil
}

// sortBalancerAssets restores the sorted-by-denom invariant PoolAssets carries.
func sortBalancerAssets(p *balancer.Pool) {
	assets := p.PoolAssets
	for i := 1; i < len(assets); i++ {
		for j := i; j > 0 && assets[j].Token.Denom < assets[j-1].Token.Denom; j-- {
			assets[j], assets[j-1] = assets[j-1], assets[j]
		}
	}
}
