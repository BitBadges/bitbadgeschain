package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

var (
	ProtocolKey     = []byte{0x01}
	CollectionIdKey = []byte{0x02}

	Delimiter   = []byte{0xDD}
	Placeholder = []byte{0xFF}

	BalanceKeyDelimiter = "-"
)

// StoreKey is the store key string for nft
const StoreKey = types.ModuleName

// protocolStoreKey returns the byte representation of the protocol key ([]byte{0x01} + protocolId)
func protocolStoreKey(protocolName string) []byte {
	key := make([]byte, len(ProtocolKey)+ len(protocolName))
	copy(key, ProtocolKey)
	copy(key[len(ProtocolKey):], []byte(protocolName))
	return key
}

func ConstructCollectionIdForProtocolKey(protocolName string, address string) string {
	protocol_name_str := protocolName
	address_str := address
	return protocol_name_str + BalanceKeyDelimiter + address_str
}

func collectionIdForProtocolStoreKey(constructedKey string) []byte {
	key := make([]byte, len(CollectionIdKey) + len(constructedKey))
	copy(key, CollectionIdKey)
	copy(key[len(CollectionIdKey):], []byte(constructedKey))
	return key
}
