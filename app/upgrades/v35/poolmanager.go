package v35

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	poolmanager "github.com/bitbadges/bitbadgeschain/x/poolmanager"
)

// PoolManagerMigrationResult reports what RekeyTakerFees touched.
type PoolManagerMigrationResult struct {
	Rekeyed int
	Total   int
}

// RekeyTakerFees moves per-pair taker-fee overrides onto the new denom.
//
// This one is easy to miss because nothing about it looks like an amount.
// x/poolmanager stores taker-fee overrides under a key built by string
// concatenation:
//
//	FormatDenomTradePairKey -> <prefix>|<tokenInDenom>|<tokenOutDenom>
//
// The denom is *in the key*, not just the value. Renaming ubadge to abadge
// therefore orphans every override involving the native denom: the record still
// exists but nothing will ever look it up again, and the pair silently falls
// back to the default taker fee. A fee change nobody asked for, on every BADGE
// pair, with no error anywhere.
//
// Deleting the stale key uses the keeper's own quirk rather than reaching into
// the store: SetDenomPairTakerFee deletes the entry when the fee it is given
// equals the default. That keeps this migration on exported API.
//
// Note that a pair whose override already equals the default is not stored at
// all, so it never appears here and needs no migration.
func RekeyTakerFees(ctx sdk.Context, pmk poolmanager.Keeper) (PoolManagerMigrationResult, error) {
	var res PoolManagerMigrationResult

	pairs, err := pmk.GetAllTradingPairTakerFees(ctx)
	if err != nil {
		return res, fmt.Errorf("reading taker fees: %w", err)
	}
	res.Total = len(pairs)

	defaultFee := pmk.GetDefaultTakerFee(ctx)

	for _, pair := range pairs {
		inDenom, outDenom := pair.TokenInDenom, pair.TokenOutDenom
		if inDenom != legacyDenom() && outDenom != legacyDenom() {
			continue
		}

		newIn, newOut := inDenom, outDenom
		if newIn == legacyDenom() {
			newIn = newDenom()
		}
		if newOut == legacyDenom() {
			newOut = newDenom()
		}

		// Guard against clobbering: if the override happens to equal the
		// default, writing it would delete rather than store, quietly losing
		// the record instead of moving it.
		if pair.TakerFee.Equal(defaultFee) {
			return res, fmt.Errorf(
				"taker fee for %s/%s equals the default; SetDenomPairTakerFee would delete rather than move it",
				inDenom, outDenom,
			)
		}

		pmk.SetDenomPairTakerFee(ctx, newIn, newOut, pair.TakerFee)

		// Retire the old key. Passing the default fee is how the keeper spells
		// "delete this entry".
		pmk.SetDenomPairTakerFee(ctx, inDenom, outDenom, defaultFee)

		res.Rekeyed++
	}

	ctx.Logger().Info(
		"v35: re-keyed poolmanager taker fees",
		"from", legacyDenom(),
		"to", newDenom(),
		"rekeyed", res.Rekeyed,
		"total_pairs", res.Total,
	)

	return res, nil
}
