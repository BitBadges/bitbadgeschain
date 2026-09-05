package keeper

import (
	"bytes"
	"strconv"

	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/types"
)

// legacyBlockTimeSeconds is the block time HOUR/DAY windows were converted to
// blocks with before v35.
const legacyBlockTimeSeconds int64 = 3

// MigrateV35WindowsToBlockTime converts HOUR and DAY windows written before
// v35 (block height start, duration in blocks at 3 s/block) to block-time
// windows (unix-second start, duration in seconds) with the same remaining
// lifetime. BLOCK windows and already converted windows are left unchanged.
// Must be called from the v35 upgrade handler.
func (k Keeper) MigrateV35WindowsToBlockTime(ctx sdk.Context) error {
	for _, prefix := range [][]byte{
		types.KeyPrefixChannelFlowWindow,
		types.KeyPrefixUniqueSendersWindow,
		types.KeyPrefixAddressTransferWindow,
	} {
		if err := k.migrateWindowPrefix(ctx, prefix); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) migrateWindowPrefix(ctx sdk.Context, prefix []byte) error {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	now := ctx.BlockTime().Unix()
	height := ctx.BlockHeight()

	for ; iterator.Valid(); iterator.Next() {
		timeframeType, timeframeDuration, err := parseWindowKeyTimeframe(iterator.Key())
		if err != nil {
			return err
		}
		if timeframeType == types.TimeframeType_TIMEFRAME_TYPE_BLOCK {
			continue
		}

		var window types.ChannelFlowWindow
		k.cdc.MustUnmarshal(iterator.Value(), &window)

		durationSeconds := types.TimeframeDurationInSeconds(timeframeType, timeframeDuration)
		if window.WindowDuration == durationSeconds {
			continue // already converted
		}

		elapsedBlocks := height - window.WindowStart
		window.WindowStart = now - elapsedBlocks*legacyBlockTimeSeconds
		window.WindowDuration = durationSeconds
		store.Set(iterator.Key(), k.cdc.MustMarshal(&window))
	}
	return nil
}

// parseWindowKeyTimeframe reads the trailing "|timeframeType|timeframeDuration"
// segment shared by every window key layout.
func parseWindowKeyTimeframe(key []byte) (types.TimeframeType, int64, error) {
	parts := bytes.Split(key, []byte("|"))
	if len(parts) < 2 {
		return 0, 0, types.ErrInvalidWindowKey.Wrapf("key %x", key)
	}
	timeframeType, err := strconv.ParseInt(string(parts[len(parts)-2]), 10, 32)
	if err != nil {
		return 0, 0, types.ErrInvalidWindowKey.Wrapf("key %x: %v", key, err)
	}
	timeframeDuration, err := strconv.ParseInt(string(parts[len(parts)-1]), 10, 64)
	if err != nil {
		return 0, 0, types.ErrInvalidWindowKey.Wrapf("key %x: %v", key, err)
	}
	return types.TimeframeType(timeframeType), timeframeDuration, nil
}
