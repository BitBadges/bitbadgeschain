package keeper

import (
	"strconv"
	"strings"

	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

var (
	BadgeKey           = []byte{0x01}
	BadgeBalanceKey    = []byte{0x02}
	NextAssetIDKey     = []byte{0x03}
	TransferManagerKey = []byte{0x04}
	// ClassTotalSupply     = []byte{0x05}

	Delimiter   = []byte{0xDD}
	Placeholder = []byte{0xFF}

	IDLength = 8

	BalanceKeyDelimiter = "-"
)

// StoreKey is the store key string for nft
const StoreKey = types.ModuleName

// badgeStoreKey returns the byte representation of the badge key ([]byte{0x01} + badgeID)
func badgeStoreKey(badgeID uint64) []byte {
	key := make([]byte, len(BadgeKey)+IDLength)
	copy(key, BadgeKey)
	copy(key[len(BadgeKey):], []byte(strconv.FormatUint(badgeID, 10)))
	return key
}

// Creates the balance key from an accountNumber and badgeID. Note this is not prefixed yet. It is just performing a delimited string concatenation.
func ConstructBalanceKey(accountNumber uint64, id uint64) string {
	badge_id_str := strconv.FormatUint(id, 10)
	account_num_str := strconv.FormatUint(accountNumber, 10)
	return account_num_str + BalanceKeyDelimiter + badge_id_str
}

// badgeBalanceStoreKey returns the byte representation of the badge balance store key ([]byte{0x02} + balanceKey)
func badgeBalanceStoreKey(balanceKey string) []byte {
	key := make([]byte, len(BadgeBalanceKey)+len(balanceKey))
	copy(key, BadgeBalanceKey)
	copy(key[len(BadgeBalanceKey):], []byte(balanceKey))
	return key
}


// Creates the transfer manager request key from an accountNumber and badgeID. Note this is not prefixed yet. It is just performing a delimited string concatenation.
func ConstructTransferManagerRequestKey(badgeId uint64, accountNumber uint64) string {
	badge_id_str := strconv.FormatUint(badgeId, 10)
	account_num_str := strconv.FormatUint(accountNumber, 10)
	return badge_id_str + BalanceKeyDelimiter + account_num_str + BalanceKeyDelimiter
}

// managerTransferRequestKey returns the byte representation of the manager transfer store key ([]byte{0x04} + id)
func managerTransferRequestKey(id string) []byte {
	key := make([]byte, len(TransferManagerKey)+len(id))
	copy(key, TransferManagerKey)
	copy(key[len(TransferManagerKey):], []byte(id))
	return key
}

// nextAssetIdKey returns the byte representation of the next asset id key ([]byte{0x03})
func nextAssetIDKey() []byte {
	return NextAssetIDKey
}

type BalanceKeyDetails struct {
	badgeId    uint64
	accountNum uint64
}

// Helper function to unparse a balance key and get the information from it.
func GetDetailsFromBalanceKey(id string) BalanceKeyDetails {
	result := strings.Split(id, BalanceKeyDelimiter)
	account_num, _ := strconv.ParseUint(result[0], 10, 64)
	badge_id, _ := strconv.ParseUint(result[1], 10, 64)

	return BalanceKeyDetails{
		accountNum: account_num,
		badgeId:    badge_id,
	}
}

