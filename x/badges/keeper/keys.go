package keeper

import (
	"strconv"
	"strings"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

var (
	CollectionKey         = []byte{0x01}
	UserBalanceKey        = []byte{0x02}
	NextCollectionIdKey   = []byte{0x03}
	UsedClaimChallengeKey = []byte{0x04}
	AddressListKey        = []byte{0x06}
	ApprovalTrackerKey    = []byte{0x07}

	AccountGenerationPrefix = []byte{0x08}
	AddressGenerationPrefix = []byte{0x09}

	NextAddressListIdKey = []byte{0x0A}

	GlobalArchiveKey = []byte{0x0B}

	Delimiter   = []byte{0xDD}
	Placeholder = []byte{0xFF}

	IDLength = 8

	BalanceKeyDelimiter = "-"
)

// StoreKey is the store key string for nft
const StoreKey = types.ModuleName

type BalanceKeyDetails struct {
	collectionId sdkmath.Uint
	address      string
}

// Helper functions to manipulate the balance keys. These aren't prefixed. They will be after they are passed into the functions further down in this file.

// Creates the balance key from an address and collectionId. Note this is not prefixed yet. It is just performing a delimited string concatenation.
func ConstructBalanceKey(address string, id sdkmath.Uint) string {
	collection_id_str := id.String()
	address_str := address
	return collection_id_str + BalanceKeyDelimiter + address_str
}

func ConstructAddressListKey(addressListId string) string {
	return addressListId
}

func ConstructApprovalTrackerKey(collectionId sdkmath.Uint, addressForApproval string, amountTrackerId string, level string, trackerType string, address string) string {
	collection_id_str := collectionId.String()
	tracker_id_str := amountTrackerId
	return collection_id_str + BalanceKeyDelimiter + addressForApproval + BalanceKeyDelimiter + tracker_id_str + BalanceKeyDelimiter + level + BalanceKeyDelimiter + trackerType + BalanceKeyDelimiter + address
}

// Creates the used claim data key from an id and data. Note this is not prefixed yet. It is just performing a delimited string concatenation.
func ConstructUsedClaimDataKey(collectionId sdkmath.Uint, claimId sdkmath.Uint) string {
	collection_id_str := collectionId.String()
	claim_id_str := claimId.String()
	return collection_id_str + BalanceKeyDelimiter + claim_id_str
}

func ConstructUsedClaimChallengeKey(collectionId sdkmath.Uint, addressForChallenge string, challengeLevel string, challengeId string, codeLeafIndex sdkmath.Uint) string {
	collection_id_str := collectionId.String()

	code_leaf_index_str := codeLeafIndex.String()
	challenge_id_str := challengeId
	address_for_challenge_str := addressForChallenge
	challenge_level_str := challengeLevel
	return collection_id_str + BalanceKeyDelimiter + address_for_challenge_str + BalanceKeyDelimiter + challenge_level_str + BalanceKeyDelimiter + challenge_id_str + BalanceKeyDelimiter + code_leaf_index_str
}

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
