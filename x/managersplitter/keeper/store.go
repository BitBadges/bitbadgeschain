package keeper

import (
	"encoding/binary"

	"github.com/bitbadges/bitbadgeschain/x/managersplitter/types"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetManagerSplitterInStore sets a manager splitter in the store
func (k Keeper) SetManagerSplitterInStore(ctx sdk.Context, ms *types.ManagerSplitter) error {
	marshaled, err := k.cdc.Marshal(ms)
	if err != nil {
		return sdkerrors.Wrap(err, "Marshal types.ManagerSplitter failed")
	}
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	store.Set(managerSplitterStoreKey(ms.Address), marshaled)
	return nil
}

// GetManagerSplitterFromStore gets a manager splitter from the store by address
func (k Keeper) GetManagerSplitterFromStore(ctx sdk.Context, addr string) (*types.ManagerSplitter, bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	marshaled := store.Get(managerSplitterStoreKey(addr))

	var ms types.ManagerSplitter
	if len(marshaled) == 0 {
		return &ms, false
	}
	k.cdc.MustUnmarshal(marshaled, &ms)
	return &ms, true
}

// GetAllManagerSplittersFromStore gets all manager splitters from the store
func (k Keeper) GetAllManagerSplittersFromStore(ctx sdk.Context) (managerSplitters []*types.ManagerSplitter) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	iterator := storetypes.KVStorePrefixIterator(store, ManagerSplitterKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var ms types.ManagerSplitter
		k.cdc.MustUnmarshal(iterator.Value(), &ms)
		managerSplitters = append(managerSplitters, &ms)
	}
	return
}

// DeleteManagerSplitterFromStore deletes a manager splitter from the store
func (k Keeper) DeleteManagerSplitterFromStore(ctx sdk.Context, addr string) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	store.Delete(managerSplitterStoreKey(addr))
}

// GetNextManagerSplitterId gets the next manager splitter ID
func (k Keeper) GetNextManagerSplitterId(ctx sdk.Context) sdkmath.Uint {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	bz := store.Get(ManagerSplitterCountKey)
	if len(bz) == 0 {
		return sdkmath.NewUint(1)
	}
	// Store as bytes - convert from big endian uint64
	var id uint64
	if len(bz) >= 8 {
		id = binary.BigEndian.Uint64(bz)
	} else {
		return sdkmath.NewUint(1)
	}
	return sdkmath.NewUint(id)
}

// SetNextManagerSplitterId sets the next manager splitter ID
func (k Keeper) SetNextManagerSplitterId(ctx sdk.Context, id sdkmath.Uint) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	// Store as bytes - convert to big endian uint64
	idUint64 := id.Uint64()
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, idUint64)
	store.Set(ManagerSplitterCountKey, bz)
}
