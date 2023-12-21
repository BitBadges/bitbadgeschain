package keeper

import (
	"strings"

	sdkerrors "cosmossdk.io/errors"
	"github.com/bitbadges/bitbadgeschain/x/protocols/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkmath "cosmossdk.io/math"
)


func (k Keeper) SetProtocolInStore(ctx sdk.Context, protocol *types.Protocol) error {
	marshaled_badge, err := k.cdc.Marshal(protocol)
	if err != nil {
		return sdkerrors.Wrap(err, "Marshal types.Protocol failed")
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(protocolStoreKey(protocol.Name), marshaled_badge)
	return nil
}

// Gets a badge from the store according to the protocolName.
func (k Keeper) GetProtocolFromStore(ctx sdk.Context, protocolName string) (*types.Protocol, bool) {
	store := ctx.KVStore(k.storeKey)
	marshaled_protocol := store.Get(protocolStoreKey(protocolName))

	var protocol types.Protocol
	if len(marshaled_protocol) == 0 {
		return &protocol, false
	}
	k.cdc.MustUnmarshal(marshaled_protocol, &protocol)
	return &protocol, true
}

// GetProtocolsFromStore defines a method for returning all badges information by key.
func (k Keeper) GetProtocolsFromStore(ctx sdk.Context) (protocols []*types.Protocol) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, ProtocolKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var protocol types.Protocol
		k.cdc.MustUnmarshal(iterator.Value(), &protocol)
		protocols = append(protocols, &protocol)
	}
	return
}

// StoreHasProtocolID determines whether the specified protocolName exists
func (k Keeper) StoreHasProtocolID(ctx sdk.Context, protocolName string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(protocolStoreKey(protocolName))
}

// DeleteProtocolFromStore deletes a badge from the store.
func (k Keeper) DeleteProtocolFromStore(ctx sdk.Context, protocolName string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(protocolStoreKey(protocolName))
}




func (k Keeper) SetProtocolCollectionInStore(ctx sdk.Context, protocolName string, address string, collectionId sdkmath.Uint) error {
	store := ctx.KVStore(k.storeKey)
	collection_id_str := collectionId.String()

	store.Set(collectionIdForProtocolStoreKey(ConstructCollectionIdForProtocolKey(protocolName, address)),  []byte(collection_id_str))
	return nil
}

// Gets a badge from the store according to the protocolName.
func (k Keeper) GetProtocolCollectionFromStore(ctx sdk.Context, protocolName string, address string) sdkmath.Uint {
	store := ctx.KVStore(k.storeKey)
	collection_id_str := string(store.Get(collectionIdForProtocolStoreKey(ConstructCollectionIdForProtocolKey(protocolName, address))))
	return sdkmath.NewUintFromString(collection_id_str)
}

// DeleteProtocolFromStore deletes a badge from the store.
func (k Keeper) DeleteProtocolCollectionFromStore(ctx sdk.Context, protocolName string, address string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(collectionIdForProtocolStoreKey(ConstructCollectionIdForProtocolKey(protocolName, address)))
}

// GetProtocolsFromStore defines a method for returning all badges information by key.
func (k Keeper) GetProtocolCollectionsFromStore(ctx sdk.Context) (names []string, addresses []string, collectionIds []sdkmath.Uint) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, ProtocolKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		collectionIds = append(collectionIds, sdkmath.NewUintFromString(string(iterator.Value())))

		name, address := GetDetailsFromKey(string(iterator.Key()))
		names = append(names, name)
		addresses = append(addresses, address)
	}
	return
}



func GetDetailsFromKey(id string) (string, string) {
	result := strings.Split(id, BalanceKeyDelimiter)
	address := result[1]
	name := result[0]

	return name, address
	
}