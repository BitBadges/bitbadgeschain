package v35

import (
	"fmt"

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
// Three places hold them: collection-level approvals, and the incoming and
// outgoing approvals on each user balance store.
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
		changed := false

		for _, approval := range balance.IncomingApprovals {
			if approval == nil || approval.ApprovalCriteria == nil {
				continue
			}
			n := scaleCoinTransfers(approval.ApprovalCriteria.CoinTransfers)
			res.CoinTransfers += n
			changed = changed || n > 0
		}
		for _, approval := range balance.OutgoingApprovals {
			if approval == nil || approval.ApprovalCriteria == nil {
				continue
			}
			n := scaleCoinTransfers(approval.ApprovalCriteria.CoinTransfers)
			res.CoinTransfers += n
			changed = changed || n > 0
		}

		if !changed {
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
	)

	return res, nil
}

// scaleCoinTransfers rewrites legacy-denom coins in place and reports how many
// individual coin entries it changed. Other denoms are left alone — a paid mint
// priced in USDC must stay priced in USDC.
func scaleCoinTransfers(transfers []*tokenizationtypes.CoinTransfer) int {
	changed := 0
	for _, transfer := range transfers {
		if transfer == nil {
			continue
		}
		for _, coin := range transfer.Coins {
			if coin == nil || coin.Denom != legacyDenom() {
				continue
			}
			coin.Denom = newDenom()
			coin.Amount = coin.Amount.Mul(ConversionFactor)
			changed++
		}
	}
	return changed
}
