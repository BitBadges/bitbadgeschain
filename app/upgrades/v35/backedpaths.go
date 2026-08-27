package v35

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	tokenizationkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
)

// BackedPathMigrationResult reports what RescaleTokenizationBackedPaths moved.
type BackedPathMigrationResult struct {
	// Collections counts the collections whose backed path was repointed.
	Collections int
	// EscrowsMoved counts the distinct escrow accounts whose balance was moved
	// to the re-derived address. Several collections can share one escrow, so
	// this is normally lower than Collections.
	EscrowsMoved int
	// Escrowed is the total moved.
	Escrowed sdkmath.Int
}

// RescaleTokenizationBackedPaths repoints collections backed by the native
// coin, and moves the coins escrowed behind them.
//
// A collection can declare Invariants.CosmosCoinBackedPath: a promise that its
// tokens are backed 1:N by a real bank coin, with the coins held in an escrow
// account and the exchange rate given by Conversion.SideA (an amount of the
// coin) against Conversion.SideB (an amount of tokens). Backing and unbacking
// move coins in and out of that escrow through the bank keeper.
//
// Two things have to move together here, and the second is the dangerous one.
//
// The rate. SideA.Amount is an amount of the coin named by SideA.Denom, so
// renaming ubadge to abadge without scaling the amount changes the exchange
// rate by 10^9. On mainnet, collection 73 declares 10^9 ubadge against 10^9
// tokens — one ubadge per token. Left unscaled that becomes one *abadge* per
// token, and the collection is suddenly backed by 10^-9 of what it claims.
//
// The escrow address. This is the part that strands money. The escrow address
// is derived by generatePathAddress as a module credential over
// (module name, BackedPathGenerationPrefix, the denom string) — the collection
// id is not an input. So the denom string is the *only* thing that determines
// the address, and changing it re-derives a completely different account. The
// coins already escrowed stay at the old address, where nothing will ever look
// for them again: unbacking would query the new address, find it empty, and the
// collection's backing would be permanently unredeemable.
//
// So the order here is: read the escrow balance at the *stored* address (which
// is where the coins actually are, whatever the derivation says), repoint the
// path, then move the balance to the newly derived address.
//
// Because the address is a function of the denom alone, every collection
// declaring the same denom shares one escrow. Mainnet shows this directly —
// collections 72, 76 and 78 declare the same IBC denom with three different
// SideA amounts and resolve to one address. The moves are therefore deduplicated
// by address, or the second collection would move an already-empty balance and
// the third would move it again.
//
// Runs after RedenominateBank, so the escrow balance is already in the new
// denom and the move is a straight transfer of the converted amount.
//
// Moved with UncheckedSetBalance rather than SendCoins for the same reason
// RedenominateBank uses it: this is not a transfer anyone initiated, it must not
// fire transfer hooks, and it must not be refusable by a blocked-address check.
func RescaleTokenizationBackedPaths(
	ctx sdk.Context,
	tk tokenizationkeeper.Keeper,
	bk bankkeeper.BaseKeeper,
	ak authkeeper.AccountKeeper,
) (BackedPathMigrationResult, error) {
	res := BackedPathMigrationResult{Escrowed: sdkmath.ZeroInt()}

	// Every backed path on the legacy denom resolves to the same new address,
	// but the *old* address is read from the path record rather than derived, so
	// a record whose address was never the derived one is still handled.
	type move struct{ from, to sdk.AccAddress }
	var moves []move
	seen := map[string]bool{}

	newAddr, err := tokenizationkeeper.DerivePathAddress(newDenom(), tokenizationkeeper.BackedPathGenerationPrefix)
	if err != nil {
		return res, fmt.Errorf("deriving the backed path address for %s: %w", newDenom(), err)
	}
	derivedLegacyAddr, err := tokenizationkeeper.DerivePathAddress(legacyDenom(), tokenizationkeeper.BackedPathGenerationPrefix)
	if err != nil {
		return res, fmt.Errorf("deriving the backed path address for %s: %w", legacyDenom(), err)
	}

	for _, collection := range tk.GetCollectionsFromStore(ctx) {
		if collection == nil || collection.Invariants == nil {
			continue
		}
		path := collection.Invariants.CosmosCoinBackedPath
		if path == nil || path.Conversion == nil || path.Conversion.SideA == nil {
			continue
		}
		if path.Conversion.SideA.Denom != legacyDenom() {
			continue
		}

		oldAddr, err := sdk.AccAddressFromBech32(path.Address)
		if err != nil {
			return res, fmt.Errorf(
				"collection %s declares a backed path on %s with an unparseable escrow address %q: %w",
				collection.CollectionId, legacyDenom(), path.Address, err)
		}
		if !oldAddr.Equals(derivedLegacyAddr) {
			// Not fatal: the stored address is where the coins are, and that is
			// what gets drained. Worth shouting about, because it means the
			// record predates or diverges from the current derivation.
			ctx.Logger().Error(
				"v35: backed path escrow address does not match the derivation for its denom",
				"collection", collection.CollectionId.String(),
				"stored", path.Address,
				"derived", derivedLegacyAddr.String(),
			)
		}

		// SideA.Amount is an amount of the denom being redenominated, so it
		// scales with it. SideB is a token balance and must not move.
		path.Conversion.SideA.Amount = path.Conversion.SideA.Amount.Mul(
			sdkmath.NewUintFromBigInt(ConversionFactor.BigInt()),
		)
		path.Conversion.SideA.Denom = newDenom()
		path.Address = newAddr.String()

		// skipInvariants for the same reason RescaleTokenizationCoinTransfers
		// does: the checks validate supply and approval shape against a chain
		// that is only half migrated at this point.
		if err := tk.SetCollectionInStore(ctx, collection, true); err != nil {
			return res, fmt.Errorf("writing collection %s: %w", collection.CollectionId, err)
		}
		res.Collections++

		if !seen[oldAddr.String()] {
			seen[oldAddr.String()] = true
			moves = append(moves, move{from: oldAddr, to: newAddr})
		}
	}

	for _, m := range moves {
		if m.from.Equals(m.to) {
			continue
		}

		amount := bk.GetBalance(ctx, m.from, newDenom()).Amount
		if !amount.IsPositive() {
			continue
		}

		// The re-derived escrow has never been transacted with, so it may have
		// no auth account. Create one, which is what a SendCoins into it would
		// have done, so later queries and any module that reads the account
		// record see a real account.
		if !ak.HasAccount(ctx, m.to) {
			ak.SetAccount(ctx, ak.NewAccountWithAddress(ctx, m.to))
		}

		credited := bk.GetBalance(ctx, m.to, newDenom()).Amount.Add(amount)
		if err := bk.UncheckedSetBalance(ctx, m.to, sdk.NewCoin(newDenom(), credited)); err != nil {
			return res, fmt.Errorf("crediting the re-derived backed path escrow %s: %w", m.to, err)
		}
		if err := bk.UncheckedSetBalance(ctx, m.from, sdk.NewCoin(newDenom(), sdkmath.ZeroInt())); err != nil {
			return res, fmt.Errorf("draining the retired backed path escrow %s: %w", m.from, err)
		}

		res.EscrowsMoved++
		res.Escrowed = res.Escrowed.Add(amount)

		ctx.Logger().Info(
			"v35: moved a backed path escrow to its re-derived address",
			"from", m.from.String(),
			"to", m.to.String(),
			"amount", amount.String(),
			"denom", newDenom(),
		)
	}

	ctx.Logger().Info(
		"v35: repointed tokenization backed paths",
		"from", legacyDenom(),
		"to", newDenom(),
		"collections", res.Collections,
		"escrows_moved", res.EscrowsMoved,
		"escrowed", res.Escrowed.String(),
	)

	return res, nil
}
