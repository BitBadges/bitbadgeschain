package v35

import (
	"fmt"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"

	tokenizationkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// TokenizationMigrationResult reports what RescaleTokenizationCoinTransfers
// touched.
type TokenizationMigrationResult struct {
	Collections   int
	UserBalances  int
	CoinTransfers int
	// AllowedDenoms counts the entries repointed in per-approval
	// UserApprovalSettings.AllowedDenoms lists. Distinct from the params-level
	// AllowedDenoms list, which MigrateDenomParams handles.
	AllowedDenoms int
}

// RescaleTokenizationCoinTransfers converts legacy-denom amounts embedded in
// approval criteria.
//
// x/tokenization approvals can carry CoinTransfers — a payment leg attached to
// a token transfer, used for paid mints, royalties and payment requests. The
// amounts are authored by users and stored inside the approval, so they are
// frozen at whatever scale existed when the approval was written. Nothing else
// in this migration reaches them.
//
// The consequence of skipping this is not a display bug. A paid mint priced at
// "1 BADGE" would, after the redenomination, charge 10^-9 of that — every
// existing paid approval silently becomes free. Royalty legs likewise.
//
// Four places hold them: collection-level approvals, the incoming and outgoing
// approvals on each user balance store, and the collection's DefaultBalances —
// the balance store every user inherits the first time they touch the
// collection. Missing that last one is the same "paid mints become free" bug,
// deferred: it only bites users who join after the upgrade, so it survives any
// test that only looks at accounts which already exist.
func RescaleTokenizationCoinTransfers(ctx sdk.Context, tk tokenizationkeeper.Keeper) (TokenizationMigrationResult, error) {
	var res TokenizationMigrationResult

	// Collection approvals. GetCollectionsFromStore returns a materialised
	// slice, so writing back as we go is safe here.
	for _, collection := range tk.GetCollectionsFromStore(ctx) {
		if collection == nil {
			continue
		}
		changed := false
		for _, approval := range collection.CollectionApprovals {
			if approval == nil || approval.ApprovalCriteria == nil {
				continue
			}
			n := scaleCoinTransfers(approval.ApprovalCriteria.CoinTransfers)
			res.CoinTransfers += n
			changed = changed || n > 0

			d := repointAllowedDenoms(approval.ApprovalCriteria.UserApprovalSettings)
			res.AllowedDenoms += d
			changed = changed || d > 0
		}
		if n := scaleUserBalanceStoreApprovals(collection.DefaultBalances); n > 0 {
			res.CoinTransfers += n
			changed = true
		}
		if !changed {
			continue
		}
		// skipInvariants: the invariant checks validate token supply and
		// approval shape, neither of which this touches, and running them
		// mid-upgrade against a partially migrated chain is not meaningful.
		if err := tk.SetCollectionInStore(ctx, collection, true); err != nil {
			return res, fmt.Errorf("writing collection %s: %w", collection.CollectionId, err)
		}
		res.Collections++
	}

	// User-level approvals.
	balances, _, _, err := tk.GetUserBalancesFromStore(ctx)
	if err != nil {
		return res, fmt.Errorf("reading user balances: %w", err)
	}
	balanceKeys := tk.GetUserBalanceIdsFromStore(ctx)
	if len(balances) != len(balanceKeys) {
		return res, fmt.Errorf(
			"user balance store returned %d balances for %d keys: refusing to migrate on a mismatched pairing",
			len(balances), len(balanceKeys),
		)
	}

	for i, balance := range balances {
		if balance == nil {
			continue
		}
		n := scaleUserBalanceStoreApprovals(balance)
		res.CoinTransfers += n

		if n == 0 {
			continue
		}
		if err := tk.SetUserBalanceInStore(ctx, balanceKeys[i], balance, true); err != nil {
			return res, fmt.Errorf("writing user balance %s: %w", balanceKeys[i], err)
		}
		res.UserBalances++
	}

	ctx.Logger().Info(
		"v35: rescaled tokenization coin transfers",
		"factor", ConversionFactor.String(),
		"collections", res.Collections,
		"user_balances", res.UserBalances,
		"coin_transfers", res.CoinTransfers,
		"allowed_denoms", res.AllowedDenoms,
	)

	return res, nil
}

// scaleUserBalanceStoreApprovals converts the coin transfers on both approval
// lists of a balance store, and reports how many coin entries it changed.
//
// Shared between the per-user stores and a collection's DefaultBalances, which
// is the same type and carries the same approvals.
func scaleUserBalanceStoreApprovals(balance *tokenizationtypes.UserBalanceStore) int {
	if balance == nil {
		return 0
	}
	changed := 0
	for _, approval := range balance.IncomingApprovals {
		if approval == nil || approval.ApprovalCriteria == nil {
			continue
		}
		changed += scaleCoinTransfers(approval.ApprovalCriteria.CoinTransfers)
	}
	for _, approval := range balance.OutgoingApprovals {
		if approval == nil || approval.ApprovalCriteria == nil {
			continue
		}
		changed += scaleCoinTransfers(approval.ApprovalCriteria.CoinTransfers)
	}
	return changed
}

// repointAllowedDenoms renames the legacy denom in an approval's per-collection
// allowlist, reporting how many entries it changed.
//
// Only collection-level ApprovalCriteria carries UserApprovalSettings; the
// per-user incoming and outgoing criteria do not have the field at all, which
// matches where CheckAndHandleCoinTransfers reads it from.
//
// This is a *second*, independent allowlist. MigrateDenomParams already moves
// the module params one; this is the per-approval narrowing on top of it, read
// by CheckAndHandleCoinTransfers when it validates a user-level coin transfer
// (see x/tokenization/keeper/coin_transfers.go). An empty list means "no extra
// restriction", so only a populated one matters.
//
// It fails closed, which is what makes it worth catching: a collection that
// pinned its user-level transfers to ["ubadge"] would, after the rename, reject
// every coin transfer priced in the live denom — including the ones this
// migration just rescaled into it.
func repointAllowedDenoms(settings *tokenizationtypes.UserApprovalSettings) int {
	if settings == nil {
		return 0
	}
	changed := 0
	for i, denom := range settings.AllowedDenoms {
		if denom == legacyDenom() {
			settings.AllowedDenoms[i] = newDenom()
			changed++
		}
	}
	return changed
}

// scaleCoinTransfers rewrites legacy-denom coins and reports how many individual
// coin entries it changed. Other denoms are left alone — a paid mint priced in
// USDC must stay priced in USDC.
//
// The slice is rebuilt and re-sorted rather than mutated in place. Coins are
// held in denom order, and the rename moves BADGE's position: [ibc/..., ubadge]
// is sorted, [ibc/..., abadge] is not. An unsorted slice fails
// sdk.Coins.Validate and makes AmountOf binary-search past the entry it wants.
// Every other conversion in this migration re-sorts for the same reason.
func scaleCoinTransfers(transfers []*tokenizationtypes.CoinTransfer) int {
	changed := 0
	for _, transfer := range transfers {
		if transfer == nil {
			continue
		}
		converted := 0
		out := make([]*sdk.Coin, 0, len(transfer.Coins))
		for _, coin := range transfer.Coins {
			if coin == nil {
				out = append(out, coin)
				continue
			}
			if coin.Denom != legacyDenom() {
				out = append(out, coin)
				continue
			}
			scaled := sdk.NewCoin(newDenom(), coin.Amount.Mul(ConversionFactor))
			out = append(out, &scaled)
			converted++
		}
		if converted == 0 {
			continue
		}
		sort.SliceStable(out, func(i, j int) bool {
			if out[i] == nil || out[j] == nil {
				return out[j] == nil && out[i] != nil
			}
			return out[i].Denom < out[j].Denom
		})
		transfer.Coins = out
		changed += converted
	}
	return changed
}
