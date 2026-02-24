package keeper

import (
	sdkmath "cosmossdk.io/math"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// isExcludedFromHolderCount checks if address should be excluded from holder counting.
func (k Keeper) isExcludedFromHolderCount(ctx sdk.Context, collection *types.TokenCollection, address string) bool {
	if types.IsMintOrTotalAddress(address) {
		return true
	}
	if k.IsBackingPathAddress(ctx, collection, address) {
		return true
	}
	if k.IsWrappingPathAddress(ctx, collection, address) {
		return true
	}
	return false
}

// hasNonZeroBalance checks if UserBalanceStore has any non-zero balance.
func hasNonZeroBalance(balance *types.UserBalanceStore) bool {
	if balance == nil || len(balance.Balances) == 0 {
		return false
	}
	for _, bal := range balance.Balances {
		if bal.Amount.GT(sdkmath.ZeroUint()) {
			return true
		}
	}
	return false
}

// UpdateHolderCount updates holder count when balance changes.
func (k Keeper) UpdateHolderCount(
	ctx sdk.Context,
	collection *types.TokenCollection,
	address string,
	oldBalance *types.UserBalanceStore,
	newBalance *types.UserBalanceStore,
) error {
	if k.isExcludedFromHolderCount(ctx, collection, address) {
		return nil
	}

	hadBalance := hasNonZeroBalance(oldBalance)
	hasBalance := hasNonZeroBalance(newBalance)

	if hadBalance == hasBalance {
		return nil
	}

	stats, _ := k.GetCollectionStatsFromStore(ctx, collection.CollectionId)

	if !hadBalance && hasBalance {
		stats.HolderCount = stats.HolderCount.Add(sdkmath.OneUint())
	} else if hadBalance && !hasBalance {
		if stats.HolderCount.GT(sdkmath.ZeroUint()) {
			stats.HolderCount = stats.HolderCount.Sub(sdkmath.OneUint())
		}
	}

	return k.SetCollectionStatsInStore(ctx, collection.CollectionId, stats)
}

// UpdateCirculatingSupplyOnBacking adds/subtracts from circulating supply using Balance[]
func (k Keeper) UpdateCirculatingSupplyOnBacking(
	ctx sdk.Context,
	collectionId sdkmath.Uint,
	balances []*types.Balance,
	isBacking bool,
) error {
	stats, _ := k.GetCollectionStatsFromStore(ctx, collectionId)

	var err error
	if isBacking {
		// Backing removes from circulation
		stats.Balances, err = types.SubtractBalancesWithZeroForUnderflows(ctx, balances, stats.Balances)
		if err != nil {
			return err
		}
	} else {
		// Unbacking adds to circulation
		stats.Balances, err = types.AddBalances(ctx, balances, stats.Balances)
		if err != nil {
			return err
		}
	}

	return k.SetCollectionStatsInStore(ctx, collectionId, stats)
}

// IncrementCirculatingSupplyOnMint adds minted balances to circulating supply
func (k Keeper) IncrementCirculatingSupplyOnMint(
	ctx sdk.Context,
	collectionId sdkmath.Uint,
	mintedBalances []*types.Balance,
) error {
	if len(mintedBalances) == 0 {
		return nil
	}

	stats, _ := k.GetCollectionStatsFromStore(ctx, collectionId)

	var err error
	stats.Balances, err = types.AddBalances(ctx, mintedBalances, stats.Balances)
	if err != nil {
		return err
	}

	return k.SetCollectionStatsInStore(ctx, collectionId, stats)
}
