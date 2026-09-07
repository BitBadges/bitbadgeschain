package storewalk

import (
	"fmt"

	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Prefix closes each iterator before writes and reads values at visit time so
// earlier key merges in the same batch cannot be overwritten by stale values.
func Prefix(ctx sdk.Context, store storetypes.KVStore, prefix []byte, visit func(key, value []byte) error) error {
	start := append([]byte(nil), prefix...)
	end := storetypes.PrefixEndBytes(prefix)
	visited := 0
	for {
		iterator := store.Iterator(start, end)
		keys := make([][]byte, 0, 1000)
		keyBytes := 0
		for ; iterator.Valid() && len(keys) < 1000 && keyBytes < 1024*1024; iterator.Next() {
			key := append([]byte(nil), iterator.Key()...)
			keys = append(keys, key)
			keyBytes += len(key)
		}
		iterator.Close()
		if len(keys) == 0 {
			break
		}
		start = append(append([]byte(nil), keys[len(keys)-1]...), 0)
		for _, key := range keys {
			if !store.Has(key) {
				continue
			}
			if err := visit(key, store.Get(key)); err != nil {
				return err
			}
			visited++
		}
		ctx.Logger().Info("migration prefix progress", "prefix", fmt.Sprintf("%x", prefix), "visited", visited)
	}
	return nil
}
