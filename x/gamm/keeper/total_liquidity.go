package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/v2/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	"github.com/bitbadges/bitbadgeschain/x/gamm/types"
)

func (k Keeper) GetTotalLiquidity(ctx sdk.Context) (sdk.Coins, error) {
	coins := sdk.Coins{}
	k.IterateDenomLiquidity(ctx, func(coin sdk.Coin) bool {
		coins = coins.Add(coin)
		return false
	})
	return coins, nil
}

func (k Keeper) setTotalLiquidity(ctx sdk.Context, coins sdk.Coins) {
	for _, coin := range coins {
		k.setDenomLiquidity(ctx, coin.Denom, coin.Amount)
	}
}

func (k Keeper) setDenomLiquidity(ctx sdk.Context, denom string, amount osmomath.Int) {
	store := ctx.KVStore(k.storeKey)
	bz, err := amount.Marshal()
	if err != nil {
		panic(err)
	}
	store.Set(types.GetDenomPrefix(denom), bz)
}

func (k Keeper) GetDenomLiquidity(ctx sdk.Context, denom string) osmomath.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetDenomPrefix(denom))
	if bz == nil {
		return osmomath.NewInt(0)
	}

	var amount osmomath.Int
	if err := amount.Unmarshal(bz); err != nil {
		panic(err)
	}
	return amount
}

func (k Keeper) IterateDenomLiquidity(ctx sdk.Context, cb func(sdk.Coin) bool) {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.KeyTotalLiquidity)

	iterator := prefixStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var amount osmomath.Int
		if err := amount.Unmarshal(iterator.Value()); err != nil {
			panic(err)
		}

		if cb(sdk.NewCoin(string(iterator.Key()), amount)) {
			break
		}
	}
}

func (k Keeper) RecordTotalLiquidityIncrease(ctx sdk.Context, coins sdk.Coins) {
	for _, coin := range coins {
		amount := k.GetDenomLiquidity(ctx, coin.Denom)
		amount = amount.Add(coin.Amount)
		k.setDenomLiquidity(ctx, coin.Denom, amount)
	}
}

func (k Keeper) RecordTotalLiquidityDecrease(ctx sdk.Context, coins sdk.Coins) {
	for _, coin := range coins {
		amount := k.GetDenomLiquidity(ctx, coin.Denom)
		amount = amount.Sub(coin.Amount)
		k.setDenomLiquidity(ctx, coin.Denom, amount)
	}
}

// RedenominateTotalLiquidity moves the total-liquidity index entry for one denom
// onto another, scaling the amount it holds.
//
// The index is keyed by denom (types.GetDenomPrefix) with the summed liquidity
// across all pools as the value, so a denom rename orphans the entry: the old
// key holds a figure nothing reads and the new denom reports zero liquidity.
//
// Low stakes on its own — the index is written at InitGenesis and by the
// liquidity record hooks, and read only by queries and genesis export — but
// leaving it wrong means the exported genesis of a redenominated chain disagrees
// with its own pools, and a stale entry survives every future round-trip through
// export/import.
//
// Exported for the v35 upgrade handler, which cannot reach setDenomLiquidity.
func (k Keeper) RedenominateTotalLiquidity(ctx sdk.Context, legacyDenom, newDenom string, factor osmomath.Int) bool {
	store := ctx.KVStore(k.storeKey)

	oldKey := types.GetDenomPrefix(legacyDenom)
	bz := store.Get(oldKey)
	if bz == nil {
		return false
	}

	var amount osmomath.Int
	if err := amount.Unmarshal(bz); err != nil {
		panic(err)
	}

	// An entry may already exist under the new denom on a re-run; add rather
	// than overwrite so nothing is lost either way.
	scaled := amount.Mul(factor).Add(k.GetDenomLiquidity(ctx, newDenom))

	store.Delete(oldKey)
	k.setDenomLiquidity(ctx, newDenom, scaled)
	return true
}
