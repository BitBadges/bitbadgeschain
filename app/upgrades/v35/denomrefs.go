package v35

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	erc20keeper "github.com/cosmos/evm/x/erc20/keeper"

	ratelimitkeeper "github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/keeper"
)

// DenomRefMigrationResult reports what RepointDenomReferences touched.
type DenomRefMigrationResult struct {
	TokenPairs      int
	RateLimitConfigs int
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
		changed = true
		res.RateLimitConfigs++
	}
	if changed {
		if err := rlk.SetParams(ctx, params); err != nil {
			return res, fmt.Errorf("setting rate limit params: %w", err)
		}
	}

	ctx.Logger().Info(
		"v35: repointed denom references",
		"from", legacyDenom(),
		"to", newDenom(),
		"token_pairs", res.TokenPairs,
		"rate_limit_configs", res.RateLimitConfigs,
	)
	return res, nil
}
