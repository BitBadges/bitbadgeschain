package keeper

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"strings"

	errorsmod "cosmossdk.io/errors"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkmath "cosmossdk.io/math"
)

var (
	CollectionKey              = []byte{0x01}
	UserBalanceKey             = []byte{0x02}
	NextCollectionIdKey        = []byte{0x03}
	UsedClaimChallengeKey      = []byte{0x04}
	AddressListKey             = []byte{0x06}
	ApprovalTrackerKey         = []byte{0x07}
	AccountGenerationPrefix    = []byte{0x08}
	AddressGenerationPrefix    = []byte{0x09}
	NextAddressListIdKey       = []byte{0x0A}
	ApprovalVersionKey         = []byte{0x0B}
	DynamicStoreKey            = []byte{0x0D}
	NextDynamicStoreIdKey      = []byte{0x0E}
	DynamicStoreValueKey       = []byte{0x0F}
	ETHSignatureTrackerKey     = []byte{0x10}
	ReservedProtocolAddressKey = []byte{0x11}
	PoolAddressCacheKey        = []byte{0x13}
	VotingTrackerKey           = []byte{0x14}
	CollectionStatsKey         = []byte{0x15}
	VotingChallengeTrackerKey  = []byte{0x16}

	WrapperPathGenerationPrefix = []byte{0x0C}
	BackedPathGenerationPrefix  = []byte{0x12}

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
	return fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s%s",
		collectionId.String(), BalanceKeyDelimiter,
		addressForChallenge, BalanceKeyDelimiter,
		approvalLevel, BalanceKeyDelimiter,
		approvalId, BalanceKeyDelimiter,
		challengeId, BalanceKeyDelimiter,
		codeLeafIndex.String())
}

// ConstructETHSignatureTrackerKey constructs a unique key for tracking ETH signature usage.
// The key includes: collectionId, approverAddress (addressForChallenge), approvalLevel, approvalId, challengeId, and the signature itself.
// This key is used to track how many times a specific signature has been used for a given approval/challenge context.
// Note: The signature field in the tracker key is the actual signature bytes (from the proof), not part of what gets signed.
// The signed message includes: nonce + "-" + initiatorAddress + "-" + collectionId + "-" + approverAddress + "-" + approvalLevel + "-" + approvalId + "-" + challengeId
func ConstructETHSignatureTrackerKey(collectionId sdkmath.Uint, addressForChallenge string, approvalLevel string, approvalId string, challengeId string, signature string) string {
	collection_id_str := collectionId.String()
	challenge_id_str := challengeId
	address_for_challenge_str := addressForChallenge
	challenge_level_str := approvalLevel
	return collection_id_str + BalanceKeyDelimiter + address_for_challenge_str + BalanceKeyDelimiter + challenge_level_str + BalanceKeyDelimiter + approvalId + BalanceKeyDelimiter + challenge_id_str + BalanceKeyDelimiter + signature
}

// ConstructVotingTrackerKey constructs a unique key for tracking votes.
// The key includes: collectionId, approverAddress, approvalLevel, approvalId, proposalId, and voterAddress.
// This key is used to store and retrieve votes for a given voting challenge.
func ConstructVotingTrackerKey(collectionId sdkmath.Uint, approverAddress string, approvalLevel string, approvalId string, proposalId string, voterAddress string) string {
	collection_id_str := collectionId.String()
	proposal_id_str := proposalId
	approver_address_str := approverAddress
	approval_level_str := approvalLevel
	voter_address_str := voterAddress
	return collection_id_str + BalanceKeyDelimiter + approver_address_str + BalanceKeyDelimiter + approval_level_str + BalanceKeyDelimiter + approvalId + BalanceKeyDelimiter + proposal_id_str + BalanceKeyDelimiter + voter_address_str
}

// ConstructVotingChallengeTrackerKey constructs a unique key for the voting challenge tracker.
// Unlike ConstructVotingTrackerKey, this is per-proposal (not per-voter) and tracks quorum state.
func ConstructVotingChallengeTrackerKey(collectionId sdkmath.Uint, approverAddress string, approvalLevel string, approvalId string, proposalId string) string {
	return collectionId.String() + BalanceKeyDelimiter + approverAddress + BalanceKeyDelimiter + approvalLevel + BalanceKeyDelimiter + approvalId + BalanceKeyDelimiter + proposalId
}

// Note be careful when getting details from a key because there could be a "-" (BalanceKeyDelimiter) in other fields.

// Helper function to unparse a balance key and get the information from it.
func GetDetailsFromBalanceKey(id string) (BalanceKeyDetails, error) {
	result := strings.Split(id, BalanceKeyDelimiter)

	// Validate that we have the expected number of parts
	if len(result) != 2 {
		return BalanceKeyDetails{}, errorsmod.Wrapf(ErrInvalidBalanceKeyFormat, "expected 2 parts, got %d", len(result))
	}

	// Validate that the collection ID can be parsed
	if result[0] == "" {
		return BalanceKeyDetails{}, errorsmod.Wrap(ErrEmptyCollectionIDInBalanceKey, "")
	}

	collection_id, err := strconv.ParseUint(result[0], 10, 64)
	if err != nil {
		return BalanceKeyDetails{}, errorsmod.Wrapf(ErrInvalidBalanceKeyFormat, "invalid collection ID '%s' in balance key: %v", result[0], err)
	}

	// Validate collection ID is not zero (collection IDs start from 1)
	if collection_id == 0 {
		return BalanceKeyDetails{}, errorsmod.Wrap(ErrCollectionIDCannotBeZero, "")
	}

	// Validate that the address is not empty
	if result[1] == "" {
		return BalanceKeyDetails{}, errorsmod.Wrap(ErrEmptyAddressInBalanceKey, "")
	}

	// Additional validation: check for potential delimiter conflicts
	if strings.Contains(result[0], BalanceKeyDelimiter) || strings.Contains(result[1], BalanceKeyDelimiter) {
		return BalanceKeyDetails{}, errorsmod.Wrapf(ErrBalanceKeyDelimiterConflict, "delimiter '%s'", BalanceKeyDelimiter)
	}

	return BalanceKeyDetails{
		address:      result[1],
		collectionId: sdkmath.NewUint(collection_id),
	}, nil
}

// Prefixer functions

// storeKey safely builds a prefixed store key from a prefix and a string suffix.
// It guards against integer overflow when computing the allocation size.
func storeKey(prefix []byte, suffix string) []byte {
	if len(suffix) > math.MaxInt-len(prefix) {
		panic("store key allocation size overflow")
	}
	key := make([]byte, len(prefix)+len(suffix))
	copy(key, prefix)
	copy(key[len(prefix):], suffix)
	return key
}

// collectionStoreKey returns the byte representation of the collection key ([]byte{0x01} + collectionId as 8-byte big-endian)
// Uses fixed-width binary encoding to prevent key collisions from decimal string truncation
func collectionStoreKey(collectionId sdkmath.Uint) []byte {
	key := make([]byte, len(CollectionKey)+IDLength)
	copy(key, CollectionKey)
	// Use fixed-width binary encoding (8 bytes big-endian) instead of decimal string
	// This prevents truncation when collectionId >= 100,000,000 (9+ digits)
	binary.BigEndian.PutUint64(key[len(CollectionKey):], collectionId.Uint64())
	return key
}

// userBalanceStoreKey returns the byte representation of the collection balance store key ([]byte{0x02} + balanceKey)
func userBalanceStoreKey(balanceKey string) []byte {
	return storeKey(UserBalanceKey, balanceKey)
}

func usedClaimChallengeStoreKey(usedClaimChallengeKey string) []byte {
	return storeKey(UsedClaimChallengeKey, usedClaimChallengeKey)
}

// nextCollectionIdKey returns the byte representation of the next asset id key ([]byte{0x03})
func nextCollectionIdKey() []byte {
	return NextCollectionIdKey
}

func nextAddressListCounterKey() []byte {
	return NextAddressListIdKey
}

func addressListStoreKey(addressListKey string) []byte {
	return storeKey(AddressListKey, addressListKey)
}

func approvalTrackerStoreKey(approvalTrackerKey string) []byte {
	return storeKey(ApprovalTrackerKey, approvalTrackerKey)
}

func approvalVersionStoreKey(approvalVersionKey string) []byte {
	return storeKey(ApprovalVersionKey, approvalVersionKey)
}

func dynamicStoreStoreKey(storeId sdkmath.Uint) []byte {
	key := make([]byte, len(DynamicStoreKey)+IDLength)
	copy(key, DynamicStoreKey)
	// Use fixed-width binary encoding (8 bytes big-endian) instead of decimal string
	// This prevents truncation when storeId >= 100,000,000 (9+ digits)
	binary.BigEndian.PutUint64(key[len(DynamicStoreKey):], storeId.Uint64())
	return key
}

func nextDynamicStoreIdKey() []byte {
	return NextDynamicStoreIdKey
}

func dynamicStoreValueStoreKey(storeId sdkmath.Uint, address string) []byte {
	key := make([]byte, len(DynamicStoreValueKey)+IDLength+len(address))
	copy(key, DynamicStoreValueKey)
	// Use fixed-width binary encoding (8 bytes big-endian) instead of decimal string
	// This prevents truncation when storeId >= 100,000,000 (9+ digits)
	binary.BigEndian.PutUint64(key[len(DynamicStoreValueKey):], storeId.Uint64())
	copy(key[len(DynamicStoreValueKey)+IDLength:], []byte(address))
	return key
}

func ethSignatureTrackerStoreKey(ethSignatureTrackerKey string) []byte {
	return storeKey(ETHSignatureTrackerKey, ethSignatureTrackerKey)
}

func votingTrackerStoreKey(votingTrackerKey string) []byte {
	return storeKey(VotingTrackerKey, votingTrackerKey)
}

func votingChallengeTrackerStoreKey(key string) []byte {
	return storeKey(VotingChallengeTrackerKey, key)
}

func reservedProtocolAddressStoreKey(address string) []byte {
	return storeKey(ReservedProtocolAddressKey, address)
}

func poolAddressCacheStoreKey(address string) []byte {
	return storeKey(PoolAddressCacheKey, address)
}

// collectionStatsStoreKey returns the byte representation of the collection stats key ([]byte{0x15} + collectionId as 8-byte big-endian)
func collectionStatsStoreKey(collectionId sdkmath.Uint) []byte {
	key := make([]byte, len(CollectionStatsKey)+IDLength)
	copy(key, CollectionStatsKey)
	binary.BigEndian.PutUint64(key[len(CollectionStatsKey):], collectionId.Uint64())
	return key
}
