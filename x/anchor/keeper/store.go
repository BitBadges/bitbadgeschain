package keeper

import (
	"strconv"

	"github.com/bitbadges/bitbadgeschain/x/anchor/types"

	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkmath "cosmossdk.io/math"
)

/** ****************************** NEXT ANCHOR ID ****************************** **/
// Gets the next anchor ID.
func (k Keeper) GetNextAnchorId(ctx sdk.Context) (sdkmath.Uint, error) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	nextId := store.Get(NextLocationIdKey)

	nextID := types.NewUintFromString(string(nextId))
	if nextID.IsZero() {
		return sdkmath.NewUint(1), nil
	}

	return nextID, nil
}

// * ****************************** ANCHORS ****************************** **/
// Set anchor location by ID
func (k Keeper) SetAnchorLocation(ctx sdk.Context, idx sdkmath.Uint, value string, creator string) error {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	blockTime := ctx.BlockTime().UnixMilli()
	anchor := types.AnchorData{
		Creator:   creator,
		Data:      value,
		Timestamp: sdkmath.NewUint(uint64(blockTime)),
	}

	marshaled_info, err := k.cdc.Marshal(&anchor)
	if err != nil {
		return err
	}

	locationKey := []byte{}
	locationKey = append(locationKey, NextLocationIdKey[0])
	locationKey = append(locationKey, strconv.FormatUint(uint64(idx.Uint64()), 10)...)

	store.Set(locationKey, marshaled_info)
	return nil
}

func (k Keeper) GetAnchorLocation(ctx sdk.Context, idx sdkmath.Uint) types.AnchorData {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	locationKey := []byte{}
	locationKey = append(locationKey, NextLocationIdKey[0])
	locationKey = append(locationKey, strconv.FormatUint(uint64(idx.Uint64()), 10)...)

	found := store.Has(locationKey)

	if !found {
		return types.AnchorData{}
	}

	val := store.Get(locationKey)

	var anchor types.AnchorData
	k.cdc.MustUnmarshal(val, &anchor)

	return anchor
}
