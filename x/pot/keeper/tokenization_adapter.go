package keeper

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	tokenizationkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// TokenizationKeeperAdapter adapts the concrete x/tokenization keeper to the
// x/pot TokenizationKeeper interface. It wraps all calls defensively so that
// no panic can escape — missing collections or balance errors return 0.
type TokenizationKeeperAdapter struct {
	keeper *tokenizationkeeper.Keeper
}

// NewTokenizationKeeperAdapter creates a new adapter wrapping the tokenization keeper.
func NewTokenizationKeeperAdapter(keeper *tokenizationkeeper.Keeper) *TokenizationKeeperAdapter {
	return &TokenizationKeeperAdapter{keeper: keeper}
}

// GetCredentialBalance returns the balance of a specific token for a given address.
// Returns 0 if the collection is not found, the balance query fails, or any other error.
// This method never panics.
func (a *TokenizationKeeperAdapter) GetCredentialBalance(ctx sdk.Context, collectionId uint64, tokenId uint64, address string) (balance uint64, err error) {
	// Top-level recover to guarantee no panic escapes.
	defer func() {
		if r := recover(); r != nil {
			balance = 0
			err = fmt.Errorf("panic in GetCredentialBalance: %v", r)
		}
	}()

	// Fetch collection from store.
	collection, found := a.keeper.GetCollectionFromStore(ctx, sdkmath.NewUint(collectionId))
	if !found {
		return 0, nil
	}

	// Fetch the balance store for this address.
	balanceStore, _, balErr := a.keeper.GetBalanceOrApplyDefault(ctx, collection, address)
	if balErr != nil {
		return 0, nil
	}

	// Build ranges for a single token ID at current block time.
	tokenIdRange := &tokenizationtypes.UintRange{
		Start: sdkmath.NewUint(tokenId),
		End:   sdkmath.NewUint(tokenId),
	}

	blockTimeMs := uint64(ctx.BlockTime().UnixMilli())
	timeRange := &tokenizationtypes.UintRange{
		Start: sdkmath.NewUint(blockTimeMs),
		End:   sdkmath.NewUint(blockTimeMs),
	}

	fetchedBalances, fetchErr := tokenizationtypes.GetBalancesForIds(
		ctx,
		[]*tokenizationtypes.UintRange{tokenIdRange},
		[]*tokenizationtypes.UintRange{timeRange},
		balanceStore.Balances,
	)
	if fetchErr != nil {
		return 0, nil
	}

	if len(fetchedBalances) == 0 {
		return 0, nil
	}

	// The amount is an sdkmath.Uint — convert to uint64.
	amount := fetchedBalances[0].Amount
	if amount.IsNil() {
		return 0, nil
	}

	return amount.Uint64(), nil
}
