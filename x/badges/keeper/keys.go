package keeper

import (
	"strconv"
	"strings"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

var (
	CollectionKey       = []byte{0x01}
	UserBalanceKey      = []byte{0x02}
	NextCollectionIdKey = []byte{0x03}
	TransferManagerKey  = []byte{0x04}
	ClaimKey            = []byte{0x05}
	NextClaimIdKey      = []byte{0x06}
	UsedClaimDataKey   	= []byte{0x07}
	UsedClaimCodeKey   	= []byte{0x08}
	UsedClaimAddressKey = []byte{0x09}
	WhitelistIndexKey   = []byte{0x0A}

	Delimiter   = []byte{0xDD}
	Placeholder = []byte{0xFF}

	IDLength = 8

	BalanceKeyDelimiter = "-"
)

// StoreKey is the store key string for nft
const StoreKey = types.ModuleName

type BalanceKeyDetails struct {
	collectionId uint64
	accountNum   uint64
}

// Helper functions to manipulate the balance keys. These aren't prefixed. They will be after they are passed into the functions further down in this file.

// Creates the balance key from an accountNumber and collectionId. Note this is not prefixed yet. It is just performing a delimited string concatenation.
func ConstructBalanceKey(accountNumber uint64, id uint64) string {
	collection_id_str := strconv.FormatUint(id, 10)
	account_num_str := strconv.FormatUint(accountNumber, 10)
	return account_num_str + BalanceKeyDelimiter + collection_id_str
}


// Creates the used claim data key from an id and data. Note this is not prefixed yet. It is just performing a delimited string concatenation.
func ConstructUsedClaimDataKey(collectionId uint64, claimId uint64) string {
	collection_id_str := strconv.FormatUint(collectionId, 10)
	claim_id_str := strconv.FormatUint(claimId, 10)
	return collection_id_str + BalanceKeyDelimiter + claim_id_str
}

func ConstructUsedClaimCodeKey(collectionId uint64, claimId uint64, codeLeafIndex uint64) string {
	collection_id_str := strconv.FormatUint(collectionId, 10)
	claim_id_str := strconv.FormatUint(claimId, 10)
	code_leaf_index_str := strconv.FormatUint(codeLeafIndex, 10)
	return collection_id_str + BalanceKeyDelimiter + claim_id_str  + BalanceKeyDelimiter + code_leaf_index_str
}

func ConstructUsedWhitelistIndexKey(collectionId uint64, claimId uint64, whitelistLeafIndex uint64) string {
	collection_id_str := strconv.FormatUint(collectionId, 10)
	claim_id_str := strconv.FormatUint(claimId, 10)
	whitelist_leaf_index_str := strconv.FormatUint(whitelistLeafIndex, 10)
	return collection_id_str + BalanceKeyDelimiter + claim_id_str  + BalanceKeyDelimiter + whitelist_leaf_index_str
}

func ConstructUsedClaimAddressKey(collectionId uint64, claimId uint64, address string) string {
	collection_id_str := strconv.FormatUint(collectionId, 10)
	claim_id_str := strconv.FormatUint(claimId, 10)
	return collection_id_str + BalanceKeyDelimiter + claim_id_str  + BalanceKeyDelimiter + address
}


// Creates the transfer manager request key from an accountNumber and collectionId. Note this is not prefixed yet. It is just performing a delimited string concatenation.
func ConstructTransferManagerRequestKey(collectionId uint64, accountNumber uint64) string {
	collection_id_str := strconv.FormatUint(collectionId, 10)
	account_num_str := strconv.FormatUint(accountNumber, 10)
	return collection_id_str + BalanceKeyDelimiter + account_num_str + BalanceKeyDelimiter
}

// Helper function to unparse a balance key and get the information from it.
func GetDetailsFromBalanceKey(id string) BalanceKeyDetails {
	result := strings.Split(id, BalanceKeyDelimiter)
	account_num, _ := strconv.ParseUint(result[0], 10, 64)
	collection_id, _ := strconv.ParseUint(result[1], 10, 64)

	return BalanceKeyDetails{
		accountNum:   account_num,
		collectionId: collection_id,
	}
}

// Prefixer functions

// collectionStoreKey returns the byte representation of the collection key ([]byte{0x01} + collectionId)
func collectionStoreKey(collectionId uint64) []byte {
	key := make([]byte, len(CollectionKey)+IDLength)
	copy(key, CollectionKey)
	copy(key[len(CollectionKey):], []byte(strconv.FormatUint(collectionId, 10)))
	return key
}

// userBalanceStoreKey returns the byte representation of the collection balance store key ([]byte{0x02} + balanceKey)
func userBalanceStoreKey(balanceKey string) []byte {
	key := make([]byte, len(UserBalanceKey)+len(balanceKey))
	copy(key, UserBalanceKey)
	copy(key[len(UserBalanceKey):], []byte(balanceKey))
	return key
}

// usedClaimDataStoreKey returns the byte representation of the used claim data store key ([]byte{0x07} + key)
func usedClaimDataStoreKey(usedClaimDataKey string) []byte {
	key := make([]byte, len(UsedClaimDataKey)+len(usedClaimDataKey))
	copy(key, UsedClaimDataKey)
	copy(key[len(UsedClaimDataKey):], []byte(usedClaimDataKey))
	return key
}

func usedClaimCodeStoreKey(usedClaimCodeKey string) []byte {
	key := make([]byte, len(UsedClaimCodeKey)+len(usedClaimCodeKey))
	copy(key, UsedClaimCodeKey)
	copy(key[len(UsedClaimCodeKey):], []byte(usedClaimCodeKey))
	return key
}

func usedWhitelistIndexStoreKey(whitelistIndexKey string) []byte {
	key := make([]byte, len(WhitelistIndexKey)+len(whitelistIndexKey))
	copy(key, WhitelistIndexKey)
	copy(key[len(WhitelistIndexKey):], []byte(whitelistIndexKey))
	return key
}

func usedClaimAddressStoreKey(usedClaimAddressKey string) []byte {
	key := make([]byte, len(UsedClaimAddressKey)+len(usedClaimAddressKey))
	copy(key, UsedClaimAddressKey)
	copy(key[len(UsedClaimAddressKey):], []byte(usedClaimAddressKey))
	return key
}

// claimStoreKey returns the byte representation of the claim store key ([]byte{0x05} + claimId)
func claimStoreKey(claimId uint64) []byte {
	key := make([]byte, len(ClaimKey)+IDLength)
	copy(key, ClaimKey)
	copy(key[len(ClaimKey):], []byte(strconv.FormatUint(claimId, 10)))
	return key
}

// managerTransferRequestKey returns the byte representation of the manager transfer store key ([]byte{0x04} + id)
func managerTransferRequestKey(id string) []byte {
	key := make([]byte, len(TransferManagerKey)+len(id))
	copy(key, TransferManagerKey)
	copy(key[len(TransferManagerKey):], []byte(id))
	return key
}

// nextCollectionIdKey returns the byte representation of the next asset id key ([]byte{0x03})
func nextCollectionIdKey() []byte {
	return NextCollectionIdKey
}

// nextClaimIdKey returns the byte representation of the next claim id key ([]byte{0x06})
func nextClaimIdKey() []byte {
	return NextClaimIdKey
}
