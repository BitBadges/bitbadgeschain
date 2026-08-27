package v35

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	erc20keeper "github.com/cosmos/evm/x/erc20/keeper"

	ratelimitkeeper "github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/keeper"
)

// DenomRefMigrationResult reports what RepointDenomReferences touched.
type DenomRefMigrationResult struct {
	TokenPairs       int
	RateLimitConfigs int
	// RateLimitCaps counts the individual MaxAmount fields rescaled inside those
	// configs. A config can carry several.
	RateLimitCaps int
	// RateLimitFlows counts the in-flight accumulator records re-keyed onto the
	// new denom.
	RateLimitFlows int
}

// RepointDenomReferences updates state that names the native denom without
// holding an amount of it.
//
// These are the ones a value-conservation check cannot catch, because no number
// changes — only a string. They were found by sweeping the exported genesis for
// the retired denom rather than by reasoning about balances.
//
// x/erc20 TokenPair maps a Cosmos denom to an ERC20 contract. A pair still
// pointing at the retired denom silently detaches the token from its ERC20
// representation: conversions in both directions stop resolving.
//
// x/ibc-rate-limit configs are matched by denom, and the module's documented
// behaviour is that a transfer with no matching config is *allowed*. So a stale
// config does not fail closed — it fails open, leaving the new denom
// unthrottled across every channel while the old config sits there matching
// nothing. That is the most dangerous shape a missed rename can take here, and
// it is invisible to any balance check.
//
// The rename alone inverts that failure mode into a worse one. Once a config
// matches the new denom again, its caps start being enforced — and the caps are
// raw amounts in the old scale. Mainnet's ubadge configs carry supply-shift
// caps of 10^16 and 3*10^16 and per-address caps of 10^15 and 10^16; those are
// ubadge figures, so leaving them alone turns a 10,000,000-BADGE daily ceiling
// into a 0.01-BADGE one and every non-trivial BADGE IBC transfer is rejected
// from the upgrade height onward. Renaming without rescaling is strictly worse
// than doing neither, so the two happen together here.
//
// The in-flight accumulators are keyed by denom too — ChannelFlow.NetFlow and
// AddressTransferData.TotalAmount — and hold amounts in the old scale as well.
// They are re-keyed and rescaled by the module itself so the current window
// carries over instead of silently restarting at zero.
func RepointDenomReferences(
	ctx sdk.Context,
	ek erc20keeper.Keeper,
	rlk ratelimitkeeper.Keeper,
) (DenomRefMigrationResult, error) {
	var res DenomRefMigrationResult

	// --- x/erc20 token pairs ---
	//
	// The pair is keyed by denom, so the record has to be removed and re-added
	// rather than edited in place.
	for _, pair := range ek.GetTokenPairs(ctx) {
		if pair.Denom != legacyDenom() {
			continue
		}
		ek.DeleteTokenPair(ctx, pair)
		pair.Denom = newDenom()
		ek.SetTokenPair(ctx, pair)
		res.TokenPairs++
	}

	// --- x/ibc-rate-limit configs ---
	params := rlk.GetParams(ctx)
	changed := false
	for i := range params.RateLimits {
		if params.RateLimits[i].Denom != legacyDenom() {
			continue
		}
		params.RateLimits[i].Denom = newDenom()

		// Every cap on this config is an amount of the denom being
		// redenominated, so each one moves by the same factor. A zero cap means
		// "disabled" to the keeper and must stay zero; Mul keeps it zero.
		for j := range params.RateLimits[i].SupplyShiftLimits {
			params.RateLimits[i].SupplyShiftLimits[j].MaxAmount =
				params.RateLimits[i].SupplyShiftLimits[j].MaxAmount.Mul(ConversionFactor)
			res.RateLimitCaps++
		}
		for j := range params.RateLimits[i].AddressLimits {
			params.RateLimits[i].AddressLimits[j].MaxAmount =
				params.RateLimits[i].AddressLimits[j].MaxAmount.Mul(ConversionFactor)
			res.RateLimitCaps++
		}
		// UniqueSenderLimits cap a *count* of senders, not an amount. They must
		// not move.

		changed = true
		res.RateLimitConfigs++
	}
	if changed {
		if err := rlk.SetParams(ctx, params); err != nil {
			return res, fmt.Errorf("setting rate limit params: %w", err)
		}
	}

	flows, err := rlk.RedenominateFlows(ctx, legacyDenom(), newDenom(), ConversionFactor)
	if err != nil {
		return res, fmt.Errorf("re-keying rate limit flows: %w", err)
	}
	res.RateLimitFlows = flows

	ctx.Logger().Info(
		"v35: repointed denom references",
		"from", legacyDenom(),
		"to", newDenom(),
		"token_pairs", res.TokenPairs,
		"rate_limit_configs", res.RateLimitConfigs,
		"rate_limit_caps", res.RateLimitCaps,
		"rate_limit_flows", res.RateLimitFlows,
	)
	return res, nil
}
