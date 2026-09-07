package storewalk

import (
	"fmt"
	"testing"

	"cosmossdk.io/log/v2"
	"github.com/cosmos/cosmos-sdk/store/v2/mem"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestWalkClosesIteratorBeforeMutatingAndResumesAfterDeletedKey(t *testing.T) {
	store := &guardStore{KVStore: mem.NewStore(), t: t}
	for i := 0; i < 2501; i++ {
		store.Set([]byte(fmt.Sprintf("a%06d", i)), []byte("value"))
	}
	store.Set([]byte("b-other"), []byte("untouched"))
	visited := 0
	ctx := sdk.Context{}.WithLogger(log.NewNopLogger())
	require.NoError(t, Prefix(ctx, store, []byte("a"), func(key, value []byte) error {
		require.Equal(t, []byte("value"), value)
		store.Delete(key)
		visited++
		return nil
	}))
	require.Equal(t, 2501, visited)
	require.Equal(t, []byte("untouched"), store.Get([]byte("b-other")))
}

type guardStore struct {
	storetypes.KVStore
	t    *testing.T
	open int
}
type guardedIterator struct {
	storetypes.Iterator
	store *guardStore
}

func (s *guardStore) Iterator(start, end []byte) storetypes.Iterator {
	s.open++
	return &guardedIterator{s.KVStore.Iterator(start, end), s}
}
func (it *guardedIterator) Close() error { it.store.open--; return it.Iterator.Close() }
func (s *guardStore) Set(key, value []byte) {
	require.Zero(s.t, s.open, "writes with an open iterator are unsafe")
	s.KVStore.Set(key, value)
}
func (s *guardStore) Delete(key []byte) {
	require.Zero(s.t, s.open, "deletes with an open iterator are unsafe")
	s.KVStore.Delete(key)
}

func TestWalkReadsMergedValueAtVisitTime(t *testing.T) {
	store := mem.NewStore()
	store.Set([]byte("a1"), []byte("one"))
	store.Set([]byte("a2"), []byte("two"))
	require.NoError(t, Prefix(sdk.Context{}.WithLogger(log.NewNopLogger()), store, []byte("a"), func(key, value []byte) error {
		if string(key) == "a1" {
			store.Set([]byte("a2"), []byte("merged"))
			store.Delete(key)
		} else {
			require.Equal(t, "merged", string(value))
		}
		return nil
	}))
}
