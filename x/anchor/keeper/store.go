package keeper

import (
	"strconv"

	"github.com/bitbadges/bitbadgeschain/x/anchor/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkmath "cosmossdk.io/math"
)

/** ****************************** NEXT ANCHOR ID ****************************** **/
// Gets the next anchor ID.
func (k Keeper) GetNextAnchorId(ctx sdk.Context) sdkmath.Uint {
	store := ctx.KVStore(k.storeKey)
	nextID := types.NewUintFromString(string((store.Get(NextLocationIdKey))))
	if nextID.IsZero() {
		return sdkmath.NewUint(1)
	}

	return nextID
}

//* ****************************** BADGES ****************************** **/
// Set anchor location by ID
func (k Keeper) SetAnchorLocation(ctx sdk.Context, idx sdkmath.Uint, value string, creator string) error {
	store := ctx.KVStore(k.storeKey)

	blockTime := ctx.BlockTime().Unix()

	anchor := types.AnchorData{
		Creator: creator,
		Data: 	value,
		Timestamp: sdk.NewUint(uint64(blockTime)),
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
	store := ctx.KVStore(k.storeKey)

	locationKey := []byte{}
	locationKey = append(locationKey, NextLocationIdKey[0])
	locationKey = append(locationKey, strconv.FormatUint(uint64(idx.Uint64()), 10)...)

	if !store.Has(locationKey) {
		return types.AnchorData{}
	}

	var anchor types.AnchorData
	k.cdc.MustUnmarshal(store.Get(locationKey), &anchor)

	return anchor
}