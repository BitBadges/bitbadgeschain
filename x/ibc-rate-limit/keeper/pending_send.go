package keeper

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/types"
)

// SetPendingSendWindow records the window start an outbound debit landed in for one limit of a packet
func (k Keeper) SetPendingSendWindow(ctx sdk.Context, port, channel string, sequence uint64, scope string, timeframeType types.TimeframeType, timeframeDuration int64, windowStart int64) {
	store := ctx.KVStore(k.storeKey)
	key := types.PendingSendKey(port, channel, sequence, scope, int32(timeframeType), timeframeDuration)
	store.Set(key, sdk.Uint64ToBigEndian(uint64(windowStart)))
}

// GetPendingSendWindow returns the recorded window start for one limit of a packet
func (k Keeper) GetPendingSendWindow(ctx sdk.Context, port, channel string, sequence uint64, scope string, timeframeType types.TimeframeType, timeframeDuration int64) (int64, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.PendingSendKey(port, channel, sequence, scope, int32(timeframeType), timeframeDuration)
	bz := store.Get(key)
	if bz == nil {
		return 0, false
	}
	return int64(sdk.BigEndianToUint64(bz)), true
}

// DeletePendingSend removes every pending-send record of a packet
func (k Keeper) DeletePendingSend(ctx sdk.Context, port, channel string, sequence uint64) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.PendingSendPrefix(port, channel, sequence))
	var keys [][]byte
	for ; iterator.Valid(); iterator.Next() {
		keys = append(keys, iterator.Key())
	}
	iterator.Close()
	for _, key := range keys {
		store.Delete(key)
	}
}
