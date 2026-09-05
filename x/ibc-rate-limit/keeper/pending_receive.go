package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetPendingReceiveWindow(ctx sdk.Context, port, channel string, sequence uint64, scope string, timeframeType types.TimeframeType, timeframeDuration int64, windowStart int64) {
	ctx.KVStore(k.storeKey).Set(types.PendingReceiveKey(port, channel, sequence, scope, int32(timeframeType), timeframeDuration), sdk.Uint64ToBigEndian(uint64(windowStart)))
}

func (k Keeper) GetPendingReceiveWindow(ctx sdk.Context, port, channel string, sequence uint64, scope string, timeframeType types.TimeframeType, timeframeDuration int64) (int64, bool) {
	bz := ctx.KVStore(k.storeKey).Get(types.PendingReceiveKey(port, channel, sequence, scope, int32(timeframeType), timeframeDuration))
	if bz == nil {
		return 0, false
	}
	return int64(sdk.BigEndianToUint64(bz)), true
}

func (k Keeper) DeletePendingReceive(ctx sdk.Context, port, channel string, sequence uint64) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.PendingReceivePrefix(port, channel, sequence))
	var keys [][]byte
	for ; iterator.Valid(); iterator.Next() {
		keys = append(keys, append([]byte(nil), iterator.Key()...))
	}
	iterator.Close()
	for _, key := range keys {
		store.Delete(key)
	}
}
