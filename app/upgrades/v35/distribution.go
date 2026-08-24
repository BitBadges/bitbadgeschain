package v35

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// DistributionMigrationResult reports what RescaleDistribution touched.
type DistributionMigrationResult struct {
	OutstandingRewards int
	Commissions        int
	CurrentRewards     int
	HistoricalRewards  int
	StartingInfos      int
}

// RescaleDistribution multiplies every reward amount in x/distribution by
// ConversionFactor.
//
// Distribution is the easiest module to forget and the most expensive to get
// wrong: it holds accrued value that no bank balance reflects yet. The community
// pool and the reward pools are real coins sitting on the distribution module
// account, which RedenominateBank moves — but the *accounting* that says who is
// owed what lives in this module's own store as DecCoins. Leave it at the old
// scale and every pending reward is understated by a factor of 10^9 relative to
// the pool backing it, so the first withdrawals drain far less than they should
// and the module account never reconciles.
//
// DelegatorStartingInfo.Stake is a share amount, not coins, and scales for the
// same reason delegator shares do in RescaleStaking: the reward calculation is
// (current period cumulative ratio - starting period ratio) * stake, so stake
// must stay in the same units as the shares it was recorded against.
//
// Slash events are deliberately not touched: ValidatorSlashEvent stores a
// *fraction*, which is dimensionless and unaffected by a change of denomination.
func RescaleDistribution(ctx sdk.Context, dk distrkeeper.Keeper) (DistributionMigrationResult, error) {
	var res DistributionMigrationResult

	// Collected before writing throughout: these iterate open store iterators,
	// and writing back inside the callback mutates the store underneath them.

	// Community pool.
	feePool, err := dk.FeePool.Get(ctx)
	if err != nil {
		return res, fmt.Errorf("reading fee pool: %w", err)
	}
	feePool.CommunityPool = scaleDecCoins(feePool.CommunityPool)
	if err := dk.FeePool.Set(ctx, feePool); err != nil {
		return res, fmt.Errorf("setting fee pool: %w", err)
	}

	// Outstanding rewards per validator.
	type outstanding struct {
		val     sdk.ValAddress
		rewards distrtypes.ValidatorOutstandingRewards
	}
	var outstandings []outstanding
	dk.IterateValidatorOutstandingRewards(ctx, func(val sdk.ValAddress, rewards distrtypes.ValidatorOutstandingRewards) bool {
		outstandings = append(outstandings, outstanding{val: val, rewards: rewards})
		return false
	})
	for _, o := range outstandings {
		o.rewards.Rewards = scaleDecCoins(o.rewards.Rewards)
		if err := dk.SetValidatorOutstandingRewards(ctx, o.val, o.rewards); err != nil {
			return res, fmt.Errorf("setting outstanding rewards for %s: %w", o.val, err)
		}
		res.OutstandingRewards++
	}

	// Accumulated commission per validator.
	type commission struct {
		val sdk.ValAddress
		com distrtypes.ValidatorAccumulatedCommission
	}
	var commissions []commission
	dk.IterateValidatorAccumulatedCommissions(ctx, func(val sdk.ValAddress, com distrtypes.ValidatorAccumulatedCommission) bool {
		commissions = append(commissions, commission{val: val, com: com})
		return false
	})
	for _, c := range commissions {
		c.com.Commission = scaleDecCoins(c.com.Commission)
		if err := dk.SetValidatorAccumulatedCommission(ctx, c.val, c.com); err != nil {
			return res, fmt.Errorf("setting commission for %s: %w", c.val, err)
		}
		res.Commissions++
	}

	// Current rewards per validator.
	type current struct {
		val     sdk.ValAddress
		rewards distrtypes.ValidatorCurrentRewards
	}
	var currents []current
	dk.IterateValidatorCurrentRewards(ctx, func(val sdk.ValAddress, rewards distrtypes.ValidatorCurrentRewards) bool {
		currents = append(currents, current{val: val, rewards: rewards})
		return false
	})
	for _, c := range currents {
		c.rewards.Rewards = scaleDecCoins(c.rewards.Rewards)
		if err := dk.SetValidatorCurrentRewards(ctx, c.val, c.rewards); err != nil {
			return res, fmt.Errorf("setting current rewards for %s: %w", c.val, err)
		}
		res.CurrentRewards++
	}

	// Historical rewards. These are cumulative reward *ratios* per share, so
	// they scale with the coins, not with the shares — a delegator's payout is
	// ratio * stake, and stake scales separately in RescaleStaking. Scaling both
	// would square the factor, so only the ratio moves here and stake moves
	// there, exactly as tokens and shares do for staking.
	type historical struct {
		val     sdk.ValAddress
		period  uint64
		rewards distrtypes.ValidatorHistoricalRewards
	}
	var historicals []historical
	dk.IterateValidatorHistoricalRewards(ctx, func(val sdk.ValAddress, period uint64, rewards distrtypes.ValidatorHistoricalRewards) bool {
		historicals = append(historicals, historical{val: val, period: period, rewards: rewards})
		return false
	})
	for _, h := range historicals {
		h.rewards.CumulativeRewardRatio = scaleDecCoins(h.rewards.CumulativeRewardRatio)
		if err := dk.SetValidatorHistoricalRewards(ctx, h.val, h.period, h.rewards); err != nil {
			return res, fmt.Errorf("setting historical rewards for %s period %d: %w", h.val, h.period, err)
		}
		res.HistoricalRewards++
	}

	// Delegator starting info: Stake is in share units.
	type startingInfo struct {
		val  sdk.ValAddress
		del  sdk.AccAddress
		info distrtypes.DelegatorStartingInfo
	}
	var startingInfos []startingInfo
	dk.IterateDelegatorStartingInfos(ctx, func(val sdk.ValAddress, del sdk.AccAddress, info distrtypes.DelegatorStartingInfo) bool {
		startingInfos = append(startingInfos, startingInfo{val: val, del: del, info: info})
		return false
	})
	for _, s := range startingInfos {
		s.info.Stake = s.info.Stake.Mul(conversionDec)
		if err := dk.SetDelegatorStartingInfo(ctx, s.val, s.del, s.info); err != nil {
			return res, fmt.Errorf("setting starting info for %s/%s: %w", s.val, s.del, err)
		}
		res.StartingInfos++
	}

	ctx.Logger().Info(
		"v35: rescaled distribution",
		"factor", ConversionFactor.String(),
		"outstanding_rewards", res.OutstandingRewards,
		"commissions", res.Commissions,
		"current_rewards", res.CurrentRewards,
		"historical_rewards", res.HistoricalRewards,
		"starting_infos", res.StartingInfos,
	)

	return res, nil
}

// scaleDecCoins rewrites the legacy denom to the new one at ConversionFactor,
// leaving every other denom alone. Distribution can hold rewards in denoms other
// than the bond denom (anything sent to the community pool), and those must not
// move.
func scaleDecCoins(coins sdk.DecCoins) sdk.DecCoins {
	if coins.IsZero() {
		return coins
	}
	out := make(sdk.DecCoins, 0, len(coins))
	for _, c := range coins {
		if c.Denom == legacyDenom() {
			out = append(out, sdk.NewDecCoinFromDec(newDenom(), c.Amount.Mul(conversionDec)))
			continue
		}
		out = append(out, c)
	}
	return out.Sort()
}

// Guard: DecCoins arithmetic below assumes LegacyDec, which is what the
// distribution module uses.
var _ = sdkmath.LegacyDec{}
