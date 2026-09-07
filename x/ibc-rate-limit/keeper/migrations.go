package keeper

import (
	"bytes"
	"math"
	"strconv"

	"github.com/bitbadges/bitbadgeschain/pkg/storewalk"
	"github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	now := ctx.BlockTime().Unix()
	height := ctx.BlockHeight()
	legacyParts := 2
	if bytes.Equal(prefix, types.KeyPrefixUniqueSendersWindow) {
		legacyParts = 1
	}
	if bytes.Equal(prefix, types.KeyPrefixAddressTransferWindow) {
		legacyParts = 3
	}
	return storewalk.Prefix(ctx, store, prefix, func(key, value []byte) error {
		// Unsuffixed windows are not consulted by the timeframe-aware limiter.
		if len(bytes.Split(key, []byte("|"))) == legacyParts {
			return nil
		}
		timeframeType, timeframeDuration, err := parseWindowKeyTimeframe(key)
		if err != nil {
			return err
		}
		if timeframeType == types.TimeframeType_TIMEFRAME_TYPE_BLOCK {
			return nil
		}
		var window types.ChannelFlowWindow
		if err := k.cdc.Unmarshal(value, &window); err != nil {
			return err
		}
		durationSeconds := types.TimeframeDurationInSeconds(timeframeType, timeframeDuration)
		if window.WindowDuration == durationSeconds {
			return nil
		}
		elapsedBlocks := height - window.WindowStart
		window.WindowStart = now - elapsedBlocks*legacyBlockTimeSeconds
		window.WindowDuration = durationSeconds
		store.Set(key, k.cdc.MustMarshal(&window))
		return nil
	})
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
	if timeframeDuration <= 0 || timeframeDuration > math.MaxInt64/86400 {
		return 0, 0, types.ErrInvalidWindowKey.Wrap("invalid timeframe duration")
	}
	switch types.TimeframeType(timeframeType) {
	case types.TimeframeType_TIMEFRAME_TYPE_BLOCK, types.TimeframeType_TIMEFRAME_TYPE_HOUR, types.TimeframeType_TIMEFRAME_TYPE_DAY:
	default:
		return 0, 0, types.ErrInvalidWindowKey.Wrap("unknown timeframe type")
	}
	return types.TimeframeType(timeframeType), timeframeDuration, nil
}
