package keeper

import (
	"strconv"
	"strings"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkmath "cosmossdk.io/math"
)

var (
	CollectionKey           = []byte{0x01}
	UserBalanceKey          = []byte{0x02}
	NextCollectionIdKey     = []byte{0x03}
	UsedClaimChallengeKey   = []byte{0x04}
	AddressListKey          = []byte{0x06}
	ApprovalTrackerKey      = []byte{0x07}
	AccountGenerationPrefix = []byte{0x08}
	AddressGenerationPrefix = []byte{0x09}
	NextAddressListIdKey    = []byte{0x0A}
	ApprovalVersionKey      = []byte{0x0B}
	DynamicStoreKey         = []byte{0x0D}
	NextDynamicStoreIdKey   = []byte{0x0E}
	DynamicStoreValueKey    = []byte{0x0F}

	WrapperPathGenerationPrefix = []byte{0x0C}

	Delimiter   = []byte{0xDD}
	Placeholder = []byte{0xFF}

	IDLength = 8

	BalanceKeyDelimiter = "-"
)

const StoreKey = types.ModuleName

type BalanceKeyDetails struct {
	collectionId sdkmath.Uint
	address      string
}

// Helper functions to manipulate the balance keys. These aren't prefixed. They will be after they are passed into the functions further down in this file.

// Creates the balance key from an address and collectionId. Note this is not prefixed yet. It is just performing a delimited string concatenation.
func ConstructBalanceKey(address string, id sdkmath.Uint) string {
	keyParts := []string{
		id.String(),
		address,
	}
	return strings.Join(keyParts, BalanceKeyDelimiter)
}

func ConstructAddressListKey(addressListId string) string {
	return addressListId
}

func ConstructApprovalTrackerKey(collectionID sdkmath.Uint, addressForApproval, approvalID, amountTrackerID, level, trackerType, address string) string {
	keyParts := []string{
		collectionID.String(),
		addressForApproval,
		approvalID,
		amountTrackerID,
		level,
		trackerType,
		address,
	}
	return strings.Join(keyParts, BalanceKeyDelimiter)
}

func ConstructApprovalVersionKey(collectionId sdkmath.Uint, approvalLevel string, approverAddress string, approvalId string) string {
	keyParts := []string{
		collectionId.String(),
		approvalLevel,
		approverAddress,
		approvalId,
	}
	return strings.Join(keyParts, BalanceKeyDelimiter)
}

// Creates the used claim data key from an id and data. Note this is not prefixed yet. It is just performing a delimited string concatenation.
func ConstructUsedClaimDataKey(collectionId sdkmath.Uint, claimId sdkmath.Uint) string {
	collection_id_str := collectionId.String()
	claim_id_str := claimId.String()
	return collection_id_str + BalanceKeyDelimiter + claim_id_str
}

func ConstructUsedClaimChallengeKey(collectionId sdkmath.Uint, addressForChallenge string, approvalLevel string, approvalId string, challengeId string, codeLeafIndex sdkmath.Uint) string {
	collection_id_str := collectionId.String()

	code_leaf_index_str := codeLeafIndex.String()
	challenge_id_str := challengeId
	address_for_challenge_str := addressForChallenge
	challenge_level_str := approvalLevel
	return collection_id_str + BalanceKeyDelimiter + address_for_challenge_str + BalanceKeyDelimiter + challenge_level_str + BalanceKeyDelimiter + approvalId + BalanceKeyDelimiter + challenge_id_str + BalanceKeyDelimiter + code_leaf_index_str
}

// Note be careful when getting details from a key because there could be a "-" (BalanceKeyDelimiter) in other fields.

// Helper function to unparse a balance key and get the information from it.
func GetDetailsFromBalanceKey(id string) BalanceKeyDetails {
	result := strings.Split(id, BalanceKeyDelimiter)
	address := result[1]
	collection_id, _ := strconv.ParseUint(result[0], 10, 64)

	return BalanceKeyDetails{
		address:      address,
		collectionId: sdkmath.NewUint(collection_id),
	}
}

// Prefixer functions

// collectionStoreKey returns the byte representation of the collection key ([]byte{0x01} + collectionId)
func collectionStoreKey(collectionId sdkmath.Uint) []byte {
	key := make([]byte, len(CollectionKey)+IDLength)
	copy(key, CollectionKey)
	copy(key[len(CollectionKey):], []byte(collectionId.String()))
	return key
}

// userBalanceStoreKey returns the byte representation of the collection balance store key ([]byte{0x02} + balanceKey)
func userBalanceStoreKey(balanceKey string) []byte {
	key := make([]byte, len(UserBalanceKey)+len(balanceKey))
	copy(key, UserBalanceKey)
	copy(key[len(UserBalanceKey):], []byte(balanceKey))
	return key
}

func usedClaimChallengeStoreKey(usedClaimChallengeKey string) []byte {
	key := make([]byte, len(UsedClaimChallengeKey)+len(usedClaimChallengeKey))
	copy(key, UsedClaimChallengeKey)
	copy(key[len(UsedClaimChallengeKey):], []byte(usedClaimChallengeKey))
	return key
}

// nextCollectionIdKey returns the byte representation of the next asset id key ([]byte{0x03})
func nextCollectionIdKey() []byte {
	return NextCollectionIdKey
}

func nextAddressListCounterKey() []byte {
	return NextAddressListIdKey
}

func addressListStoreKey(addressListKey string) []byte {
	key := make([]byte, len(AddressListKey)+len(addressListKey))
	copy(key, AddressListKey)
	copy(key[len(AddressListKey):], []byte(addressListKey))
	return key
}

func approvalTrackerStoreKey(approvalTrackerKey string) []byte {
	key := make([]byte, len(ApprovalTrackerKey)+len(approvalTrackerKey))
	copy(key, ApprovalTrackerKey)
	copy(key[len(ApprovalTrackerKey):], []byte(approvalTrackerKey))
	return key
}

func approvalVersionStoreKey(approvalVersionKey string) []byte {
	key := make([]byte, len(ApprovalVersionKey)+len(approvalVersionKey))
	copy(key, ApprovalVersionKey)
	copy(key[len(ApprovalVersionKey):], []byte(approvalVersionKey))
	return key
}

func dynamicStoreStoreKey(storeId sdkmath.Uint) []byte {
	key := make([]byte, len(DynamicStoreKey)+IDLength)
	copy(key, DynamicStoreKey)
	copy(key[len(DynamicStoreKey):], []byte(storeId.String()))
	return key
}

func nextDynamicStoreIdKey() []byte {
	return NextDynamicStoreIdKey
}

func dynamicStoreValueStoreKey(storeId sdkmath.Uint, address string) []byte {
	key := make([]byte, len(DynamicStoreValueKey)+IDLength+len(address))
	copy(key, DynamicStoreValueKey)
	copy(key[len(DynamicStoreValueKey):], []byte(storeId.String()))
	copy(key[len(DynamicStoreValueKey)+IDLength:], []byte(address))
	return key
}
