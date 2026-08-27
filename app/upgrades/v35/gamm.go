package v35

import (
	"fmt"
	"math"
	"sort"

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
	// TotalLiquidityRekeyed reports whether the denom-keyed total-liquidity
	// index entry was moved onto the new denom.
	TotalLiquidityRekeyed bool
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

	// The denom-keyed total-liquidity index sits beside the pools and holds the
	// summed reserves per denom. It is written by the liquidity hooks and read
	// by queries and genesis export, so the rename orphans it exactly the way it
	// orphans a taker-fee key — with the added twist that its value is an amount
	// and needs the same 10^9 as the pool assets below.
	res.TotalLiquidityRekeyed = gk.RedenominateTotalLiquidity(ctx, legacyDenom(), newDenom(), ConversionFactor)

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
			liquidity, factors, err := rescaleStableswapLegs(p.PoolLiquidity, p.ScalingFactors)
			if err != nil {
				return res, fmt.Errorf("stableswap pool %d: %w", p.Id, err)
			}
			p.PoolLiquidity, p.ScalingFactors = liquidity, factors
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
		"total_liquidity_rekeyed", res.TotalLiquidityRekeyed,
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

// rescaleStableswapLegs converts a stableswap pool's native-denom leg.
//
// Stableswap prices off liquidity[i] / scalingFactors[i], and the two arrays are
// index-aligned with each other. So the migration has to do two things the
// balancer path does not:
//
//   - scale the BADGE scaling factor along with the BADGE liquidity. Scaling
//     only the liquidity makes the BADGE leg look 10^9 times deeper than it is,
//     the pool quotes BADGE at effectively zero, and the first swapper drains
//     the other side of the pool. This finding carries no denom string in the
//     pool record, which is why the exported-genesis sweep cannot see it.
//   - re-sort BOTH arrays together. PoolLiquidity is sorted by denom and abadge
//     sorts where ubadge did not, so re-sorting the coins alone would leave
//     every scaling factor attached to the wrong asset.
func rescaleStableswapLegs(liquidity sdk.Coins, factors []uint64) (sdk.Coins, []uint64, error) {
	if len(liquidity) != len(factors) {
		return nil, nil, fmt.Errorf(
			"pool has %d liquidity legs but %d scaling factors; refusing to guess the pairing",
			len(liquidity), len(factors),
		)
	}

	type leg struct {
		coin   sdk.Coin
		factor uint64
	}

	legs := make([]leg, 0, len(liquidity))
	changed := false
	for i, coin := range liquidity {
		factor := factors[i]
		if coin.Denom == legacyDenom() {
			scaled, err := scaleScalingFactor(factor)
			if err != nil {
				return nil, nil, err
			}
			coin = sdk.NewCoin(newDenom(), coin.Amount.Mul(ConversionFactor))
			factor = scaled
			changed = true
		}
		legs = append(legs, leg{coin: coin, factor: factor})
	}

	if !changed {
		return liquidity, factors, nil
	}

	sort.SliceStable(legs, func(i, j int) bool { return legs[i].coin.Denom < legs[j].coin.Denom })

	outCoins := make(sdk.Coins, len(legs))
	outFactors := make([]uint64, len(legs))
	for i, l := range legs {
		outCoins[i] = l.coin
		outFactors[i] = l.factor
	}
	return outCoins, outFactors, nil
}

// maxScalingFactor is the largest scaling factor x/gamm can actually use.
//
// The field is uint64, but validateScalingFactors rejects anything where
// int64(scalingFactor) <= 0 and getDescaledPoolAmt multiplies by
// int64(scalingFactor), so the usable ceiling is MaxInt64 rather than MaxUint64.
// Scaling by 10^9 crosses it for any factor above ~9.2e9 — reachable, unlike the
// sdkmath.Int overflow bounds elsewhere in this migration.
const maxScalingFactor = uint64(math.MaxInt64)

func scaleScalingFactor(factor uint64) (uint64, error) {
	multiplier := ConversionFactor.Uint64()
	if factor > maxScalingFactor/multiplier {
		return 0, fmt.Errorf(
			"scaling factor %d overflows when multiplied by %d (max usable factor is %d): "+
				"drain or re-parameterise this pool before the upgrade height",
			factor, multiplier, maxScalingFactor/multiplier,
		)
	}
	return factor * multiplier, nil
}
