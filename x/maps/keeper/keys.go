package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

var (
	MapKey      = []byte{0x01}
	MapValueKey = []byte{0x02}

	MapValueDuplicatesKey = []byte{0x03}

	Delimiter   = []byte{0xDD}
	Placeholder = []byte{0xFF}

	BalanceKeyDelimiter = "-"
)

// StoreKey is the store key string for nft
const StoreKey = types.ModuleName

// mapStoreKey returns the byte representation of the protocol key ([]byte{0x01} + protocolId)
func mapStoreKey(mapId string) []byte {
	key := make([]byte, len(MapKey)+len(mapId))
	copy(key, MapKey)
	copy(key[len(MapKey):], []byte(mapId))
	return key
}

func ConstructMapValueKey(mapId string, key string) string {
	protocol_name_str := mapId
	address_str := key
	return protocol_name_str + BalanceKeyDelimiter + address_str
}

func mapValueStoreKey(constructedKey string) []byte {
	key := make([]byte, len(MapValueKey)+len(constructedKey))
	copy(key, MapValueKey)
	copy(key[len(MapValueKey):], []byte(constructedKey))
	return key
}

// Note be careful when getting details from a key because there could be a "-" (BalanceKeyDelimiter) in other fields.

func ConstructMapValueDuplicatesKey(mapId string, value string) string {
	protocol_name_str := mapId
	address_str := value
	return protocol_name_str + BalanceKeyDelimiter + address_str
}

func mapValueDuplicatesStoreKey(constructedKey string) []byte {
	key := make([]byte, len(MapValueDuplicatesKey)+len(constructedKey))
	copy(key, MapValueDuplicatesKey)
	copy(key[len(MapValueDuplicatesKey):], []byte(constructedKey))
	return key
}
