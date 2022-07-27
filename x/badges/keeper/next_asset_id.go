package keeper

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetNextAssetId(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	next_id, err := strconv.ParseUint(string((store.Get(NextAssetIDKey))), 10, 64);
	if err != nil {
		panic("Failed to get next asset ID");
	}
	return next_id
}

// Used in InitGenesis
func (k Keeper) SetNextAssetId(ctx sdk.Context, next_id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(NextAssetIDKey, []byte(strconv.FormatInt(int64(next_id), 10)))
}



