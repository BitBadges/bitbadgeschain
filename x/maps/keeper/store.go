package keeper

import (
	"strings"

	"bitbadgeschain/x/maps/types"

	sdkerrors "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
)

func (k Keeper) SetMapInStore(ctx sdk.Context, protocol *types.Map) error {
	marshaled_badge, err := k.cdc.Marshal(protocol)
	if err != nil {
		return sdkerrors.Wrap(err, "Marshal types.Map failed")
	}
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	store.Set(mapStoreKey(protocol.MapId), marshaled_badge)
	return nil
}

// Gets a badge from the store according to the mapId.
func (k Keeper) GetMapFromStore(ctx sdk.Context, mapId string) (*types.Map, bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	marshaled_protocol := store.Get(mapStoreKey(mapId))

	var protocol types.Map
	if len(marshaled_protocol) == 0 {
		return &protocol, false
	}
	k.cdc.MustUnmarshal(marshaled_protocol, &protocol)
	return &protocol, true
}

// GetMapsFromStore defines a method for returning all badges information by key.
func (k Keeper) GetMapsFromStore(ctx sdk.Context) (protocols []*types.Map) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	iterator := storetypes.KVStorePrefixIterator(store, MapKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var protocol types.Map
		k.cdc.MustUnmarshal(iterator.Value(), &protocol)
		protocols = append(protocols, &protocol)
	}
	return
}

// StoreHasMapID determines whether the specified mapId exists
func (k Keeper) StoreHasMapID(ctx sdk.Context, mapId string) bool {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	return store.Has(mapStoreKey(mapId))
}

// DeleteMapFromStore deletes a badge from the store.
func (k Keeper) DeleteMapFromStore(ctx sdk.Context, mapId string) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	store.Delete(mapStoreKey(mapId))
}

func (k Keeper) SetMapValueInStore(ctx sdk.Context, mapId string, key string, value string, address string) error {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})

	valueStore := &types.ValueStore{
		Key:       key,
		Value:     value,
		LastSetBy: address,
	}

	marshaled_protocol, err := k.cdc.Marshal(valueStore)
	if err != nil {
		return sdkerrors.Wrap(err, "Marshal types.ValueStore failed")
	}

	store.Set(mapValueStoreKey(ConstructMapValueKey(mapId, key)), marshaled_protocol)
	return nil
}

// Gets a badge from the store according to the mapId.
func (k Keeper) GetMapValueFromStore(ctx sdk.Context, mapId string, value string) types.ValueStore {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})

	var valueStore types.ValueStore
	marshaled_protocol := store.Get(mapValueStoreKey(ConstructMapValueKey(mapId, value)))
	if len(marshaled_protocol) == 0 {
		return valueStore
	}

	k.cdc.MustUnmarshal(marshaled_protocol, &valueStore)
	return valueStore
}

// DeleteMapFromStore deletes a badge from the store.
func (k Keeper) DeleteMapValueFromStore(ctx sdk.Context, mapId string, value string) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	store.Delete(mapValueStoreKey(ConstructMapValueKey(mapId, value)))
}

func (k Keeper) GetMapKeysAndValuesFromStore(ctx sdk.Context) (mapIds []string, keys []string, values []*types.ValueStore) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	iterator := storetypes.KVStorePrefixIterator(store, MapValueKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var valueStore types.ValueStore
		k.cdc.MustUnmarshal(iterator.Value(), &valueStore)
		values = append(values, &valueStore)

		mapId, key := GetDetailsFromKey(string(iterator.Key()[1:]))
		mapIds = append(mapIds, mapId)
		keys = append(keys, key)
	}
	return
}

func GetDetailsFromKey(id string) (string, string) {
	result := strings.Split(id, BalanceKeyDelimiter)
	key := result[1]
	mapId := result[0]

	return mapId, key
}

func (k Keeper) SetMapDuplicateValueInStore(ctx sdk.Context, mapId string, value string) error {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})

	store.Set(mapValueDuplicatesStoreKey(ConstructMapValueDuplicatesKey(mapId, value)), []byte("true"))
	return nil
}

func (k Keeper) GetMapDuplicateValueFromStore(ctx sdk.Context, mapId string, value string) bool {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})

	return store.Has(mapValueDuplicatesStoreKey(ConstructMapValueDuplicatesKey(mapId, value)))
}

func (k Keeper) DeleteMapDuplicateValueFromStore(ctx sdk.Context, mapId string, value string) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	store.Delete(mapValueDuplicatesStoreKey(ConstructMapValueDuplicatesKey(mapId, value)))
}

func (k Keeper) GetMapDuplicateKeysAndValuesFromStore(ctx sdk.Context) (mapIds []string, values []string) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	iterator := storetypes.KVStorePrefixIterator(store, MapValueDuplicatesKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		mapId, value := GetDetailsFromKey(string(iterator.Key()[1:]))
		mapIds = append(mapIds, mapId)
		values = append(values, value)
	}
	return
}
