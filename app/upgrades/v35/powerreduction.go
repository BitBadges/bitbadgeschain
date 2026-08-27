package v35

import (
	"context"
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	appparams "github.com/bitbadges/bitbadgeschain/app/params"
)

// LegacyPowerReduction is the value sdk.DefaultPowerReduction had before this
// upgrade: the SDK default of 10^6.
//
// appparams.PowerReduction is that scaled by the same ConversionFactor the
// balances move by, which is what keeps every validator's consensus power
// numerically identical across the upgrade. Consensus power is
// tokens / PowerReduction, so multiplying the numerator by 10^9 and leaving the
// denominator alone would multiply every validator's power by 10^9 — straight
// at CometBFT's MaxTotalVotingPower ceiling.
var LegacyPowerReduction = sdkmath.NewIntWithDecimal(1, 6)

// AssertPowerReductionMatchesBondDenom refuses to run the upgrade on a chain
// whose bond denom this migration does not redenominate.
//
// Why this guard exists, in full, because the divergence it protects has
// already misled one reviewer:
//
//   - **Mainnet (bitbadges-1) bonds ubadge.** Both
//     /cosmos/staking/v1beta1/params and /cosmos/mint/v1beta1/params return
//     "ubadge" on the live chain. ubadge is the denom this upgrade
//     redenominates, so PowerReduction moving 10^6 -> 10^15 moves in lockstep
//     with the token scale and consensus power is preserved exactly. Checked
//     against the largest live validator: 307,572,905,653,630 ubadge / 10^6 =
//     307,572,905 before; (307,572,905,653,630 * 10^9) / 10^15 = 307,572,905
//     after. Identical.
//
//   - **The local dev chain bonds ustake.** genesis.json, config.yml and
//     start-chain.sh all use "ustake" as the bond denom, and so does the
//     archived genesis-711316.json export — which is a *pre-v33* snapshot of
//     mainnet, from before the bond denom moved to ubadge, and is the likeliest
//     source of the belief that this chain still bonds ustake. ustake is NOT
//     redenominated by this upgrade, which is exactly why start-chain.sh's
//     ustake amounts were scaled by hand rather than migrated. On such a chain,
//     PowerReduction jumping to 10^15 while the bond denom stays at its original
//     scale divides every validator's power by 10^9 and collapses the validator
//     set to zero power: an instant, silent, unrecoverable consensus halt.
//
// The bond denom having demonstrably changed once already in this chain's
// history is the whole argument for checking it at runtime rather than
// hardcoding the assumption in a comment.
//
// So: PowerReduction is only allowed to have moved when the bond denom is the
// one being redenominated. Anywhere else, fail loudly at the top of the handler
// — before a single balance is touched — rather than let the chain discover it
// at the next EndBlocker.
//
// The new denom is accepted alongside the legacy one so a re-run of the handler
// against already-migrated state (which the idempotency test does, and which a
// replay can do) is not mistaken for a foreign chain.
func AssertPowerReductionMatchesBondDenom(goCtx context.Context, sk *stakingkeeper.Keeper) error {
	ctx := sdk.UnwrapSDKContext(goCtx)

	bondDenom, err := sk.BondDenom(ctx)
	if err != nil {
		return fmt.Errorf("reading bond denom: %w", err)
	}

	if bondDenom == legacyDenom() || bondDenom == newDenom() {
		return nil
	}

	// A chain bonding something else may still run this binary safely, but only
	// if PowerReduction was left where the SDK put it.
	if sdk.DefaultPowerReduction.Equal(LegacyPowerReduction) {
		ctx.Logger().Info(
			"v35: bond denom is not the redenominated denom, and PowerReduction was left at the SDK default",
			"bond_denom", bondDenom,
		)
		return nil
	}

	return fmt.Errorf(
		"refusing to run the v35 redenomination: the bond denom is %q, which this upgrade does not "+
			"redenominate, but PowerReduction has been moved from %s to %s. Consensus power is "+
			"tokens/PowerReduction, so every validator's power would be divided by %s and the validator "+
			"set would collapse to zero power at the next EndBlocker. Mainnet bonds %q and is unaffected; "+
			"the local dev chain bonds ustake and must not run this handler",
		bondDenom, LegacyPowerReduction, sdk.DefaultPowerReduction,
		sdk.DefaultPowerReduction.Quo(LegacyPowerReduction),
		appparams.LegacyBaseCoinUnit,
	)
}
