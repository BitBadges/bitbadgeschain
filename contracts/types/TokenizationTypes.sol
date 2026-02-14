// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * @title TokenizationTypes
 * @notice Comprehensive type registry for BitBadges tokenization module
 * @dev This file contains all Solidity struct definitions that mirror proto message types
 *      from the tokenization module. All types are 1:1 mappings from proto definitions.
 *      NOTE: Creator fields are NOT included in input structs - they are always msg.sender
 */

// ============================================================================
// Core Types (balances.proto)
// ============================================================================

/**
 * @notice Range of IDs from start to end (inclusive)
 * @dev Used for token IDs, ownership times, transfer times, etc.
 */
struct UintRange {
    uint256 start;
    uint256 end;
}

/**
 * @notice Balance of a token for a specific user
 * @dev User owns amount of tokens for the specified token IDs and ownership times
 */
struct Balance {
    uint256 amount;
    UintRange[] ownershipTimes;
    UintRange[] tokenIds;
}

/**
 * @notice Options for precalculating balances
 */
struct PrecalculationOptions {
    uint256 overrideTimestamp;
    UintRange[] tokenIdsOverride;
}

// ============================================================================
// Metadata Types (metadata.proto)
// ============================================================================

/**
 * @notice Metadata for specific token IDs
 */
struct TokenMetadata {
    string uri;
    string customData;
    UintRange[] tokenIds;
}

/**
 * @notice Metadata for a collection
 */
struct CollectionMetadata {
    string uri;
    string customData;
}

/**
 * @notice Metadata for paths (alias paths and cosmos coin wrapper paths)
 */
struct PathMetadata {
    string uri;
    string customData;
}

// ============================================================================
// Permission Types (permissions.proto)
// ============================================================================

/**
 * @notice Permission for performing an action
 * @dev Simple permission that only checks permitted/forbidden times
 */
struct ActionPermission {
    UintRange[] permanentlyPermittedTimes;
    UintRange[] permanentlyForbiddenTimes;
}

/**
 * @notice Permission for performing an action for specific tokens
 */
struct TokenIdsActionPermission {
    UintRange[] tokenIds;
    UintRange[] permanentlyPermittedTimes;
    UintRange[] permanentlyForbiddenTimes;
}

/**
 * @notice Permission for updating collection approvals
 */
struct CollectionApprovalPermission {
    string fromListId;
    string toListId;
    string initiatedByListId;
    UintRange[] transferTimes;
    UintRange[] tokenIds;
    UintRange[] ownershipTimes;
    string approvalId;
    UintRange[] permanentlyPermittedTimes;
    UintRange[] permanentlyForbiddenTimes;
}

/**
 * @notice Permission for updating user outgoing approvals
 */
struct UserOutgoingApprovalPermission {
    string toListId;
    string initiatedByListId;
    UintRange[] transferTimes;
    UintRange[] tokenIds;
    UintRange[] ownershipTimes;
    string approvalId;
    UintRange[] permanentlyPermittedTimes;
    UintRange[] permanentlyForbiddenTimes;
}

/**
 * @notice Permission for updating user incoming approvals
 */
struct UserIncomingApprovalPermission {
    string fromListId;
    string initiatedByListId;
    UintRange[] transferTimes;
    UintRange[] tokenIds;
    UintRange[] ownershipTimes;
    string approvalId;
    UintRange[] permanentlyPermittedTimes;
    UintRange[] permanentlyForbiddenTimes;
}

/**
 * @notice Permissions for collection operations
 */
struct CollectionPermissions {
    ActionPermission[] canDeleteCollection;
    ActionPermission[] canArchiveCollection;
    ActionPermission[] canUpdateStandards;
    ActionPermission[] canUpdateCustomData;
    ActionPermission[] canUpdateManager;
    ActionPermission[] canUpdateCollectionMetadata;
    TokenIdsActionPermission[] canUpdateValidTokenIds;
    TokenIdsActionPermission[] canUpdateTokenMetadata;
    CollectionApprovalPermission[] canUpdateCollectionApprovals;
    ActionPermission[] canAddMoreAliasPaths;
    ActionPermission[] canAddMoreCosmosCoinWrapperPaths;
}

/**
 * @notice Permissions for user operations
 */
struct UserPermissions {
    UserOutgoingApprovalPermission[] canUpdateOutgoingApprovals;
    UserIncomingApprovalPermission[] canUpdateIncomingApprovals;
    ActionPermission[] canUpdateAutoApproveSelfInitiatedOutgoingTransfers;
    ActionPermission[] canUpdateAutoApproveSelfInitiatedIncomingTransfers;
    ActionPermission[] canUpdateAutoApproveAllIncomingTransfers;
}

// ============================================================================
// Approval Tracking Types (approval_tracking.proto)
// ============================================================================

/**
 * @notice Options for auto-deletion of approvals
 */
struct AutoDeletionOptions {
    bool afterOneUse;
    bool afterOverallMaxNumTransfers;
    bool allowCounterpartyPurge;
    bool allowPurgeIfExpired;
}

/**
 * @notice Time intervals to reset trackers at
 */
struct ResetTimeIntervals {
    uint256 startTime;
    uint256 intervalLength;
}

/**
 * @notice Approval amounts per unique address
 */
struct ApprovalAmounts {
    uint256 overallApprovalAmount;
    uint256 perToAddressApprovalAmount;
    uint256 perFromAddressApprovalAmount;
    uint256 perInitiatedByAddressApprovalAmount;
    string amountTrackerId;
    ResetTimeIntervals resetTimeIntervals;
}

/**
 * @notice Maximum number of transfers per unique address
 */
struct MaxNumTransfers {
    uint256 overallMaxNumTransfers;
    uint256 perToAddressMaxNumTransfers;
    uint256 perFromAddressMaxNumTransfers;
    uint256 perInitiatedByAddressMaxNumTransfers;
    string amountTrackerId;
    ResetTimeIntervals resetTimeIntervals;
}

/**
 * @notice Tracker for approvals
 */
struct ApprovalTracker {
    uint256 numTransfers;
    Balance[] amounts;
    uint256 lastUpdatedAt;
}

// ============================================================================
// Challenge Types (challenges.proto)
// ============================================================================

/**
 * @notice Merkle challenge for approval
 */
struct MerkleChallenge {
    string root;
    uint256 expectedProofLength;
    bool useCreatorAddressAsLeaf;
    uint256 maxUsesPerLeaf;
    string uri;
    string customData;
    string challengeTrackerId;
    string leafSigner;
}

/**
 * @notice ETH signature challenge for approval
 */
struct ETHSignatureChallenge {
    string signer;
    string challengeTrackerId;
    string uri;
    string customData;
}

/**
 * @notice Voting challenge for approval
 */
struct VotingChallenge {
    string proposalId;
    uint256 quorumThreshold;
    Voter[] voters;
    string uri;
    string customData;
}

/**
 * @notice Voter in a voting challenge
 */
struct Voter {
    string address_;
    uint256 weight;
}

/**
 * @notice Item in a Merkle path
 */
struct MerklePathItem {
    string aunt;
    bool onRight;
}

/**
 * @notice Merkle proof
 */
struct MerkleProof {
    string leaf;
    MerklePathItem[] aunts;
    string leafSignature;
}

/**
 * @notice ETH signature proof
 */
struct ETHSignatureProof {
    string nonce;
    string signature;
}

/**
 * @notice Vote proof for voting challenge
 */
struct VoteProof {
    string proposalId;
    string voter;
    uint256 yesWeight;
}

// ============================================================================
// Approval Conditions Types (approval_conditions.proto)
// ============================================================================

/**
 * @notice Coin transfer requirement
 */
struct CoinTransfer {
    string to;
    // Note: Cosmos Coin type would need special handling - using string representation
    string[] coinDenoms;
    uint256[] coinAmounts;
    bool overrideFromWithApproverAddress;
    bool overrideToWithInitiator;
}

/**
 * @notice Must own tokens requirement
 */
struct MustOwnTokens {
    uint256 collectionId;
    UintRange amountRange;
    UintRange[] ownershipTimes;
    UintRange[] tokenIds;
    bool overrideWithCurrentTime;
    bool mustSatisfyForAllAssets;
    string ownershipCheckParty;
}

/**
 * @notice Dynamic store challenge
 */
struct DynamicStoreChallenge {
    uint256 storeId;
    string ownershipCheckParty;
}

/**
 * @notice Address checks
 */
struct AddressChecks {
    bool mustBeEvmContract;
    bool mustNotBeEvmContract;
    bool mustBeLiquidityPool;
    bool mustNotBeLiquidityPool;
}

/**
 * @notice Alternative time-based checks
 */
struct AltTimeChecks {
    UintRange[] offlineHours;
    UintRange[] offlineDays;
}

/**
 * @notice User royalties
 */
struct UserRoyalties {
    uint256 percentage;
    string payoutAddress;
}

// ============================================================================
// Predetermined Balances Types (predetermined_balances.proto)
// ============================================================================

/**
 * @notice Manual balances list
 */
struct ManualBalances {
    Balance[] balances;
}

/**
 * @notice Recurring ownership times
 */
struct RecurringOwnershipTimes {
    uint256 startTime;
    uint256 intervalLength;
    uint256 chargePeriodLength;
}

/**
 * @notice Incremented balances
 */
struct IncrementedBalances {
    Balance[] startBalances;
    uint256 incrementTokenIdsBy;
    uint256 incrementOwnershipTimesBy;
    uint256 durationFromTimestamp;
    bool allowOverrideTimestamp;
    RecurringOwnershipTimes recurringOwnershipTimes;
    bool allowOverrideWithAnyValidToken;
}

/**
 * @notice Order calculation method for predetermined balances
 */
struct PredeterminedOrderCalculationMethod {
    bool useOverallNumTransfers;
    bool usePerToAddressNumTransfers;
    bool usePerFromAddressNumTransfers;
    bool usePerInitiatedByAddressNumTransfers;
    bool useMerkleChallengeLeafIndex;
    string challengeTrackerId;
}

/**
 * @notice Predetermined balances
 */
struct PredeterminedBalances {
    ManualBalances[] manualBalances;
    IncrementedBalances incrementedBalances;
    PredeterminedOrderCalculationMethod orderCalculationMethod;
}

// ============================================================================
// Approval Criteria Types (approval_criteria.proto)
// ============================================================================

/**
 * @notice Criteria for collection-level approvals
 */
struct ApprovalCriteria {
    MerkleChallenge[] merkleChallenges;
    PredeterminedBalances predeterminedBalances;
    ApprovalAmounts approvalAmounts;
    MaxNumTransfers maxNumTransfers;
    CoinTransfer[] coinTransfers;
    bool requireToEqualsInitiatedBy;
    bool requireFromEqualsInitiatedBy;
    bool requireToDoesNotEqualInitiatedBy;
    bool requireFromDoesNotEqualInitiatedBy;
    bool overridesFromOutgoingApprovals;
    bool overridesToIncomingApprovals;
    AutoDeletionOptions autoDeletionOptions;
    UserRoyalties userRoyalties;
    MustOwnTokens[] mustOwnTokens;
    DynamicStoreChallenge[] dynamicStoreChallenges;
    ETHSignatureChallenge[] ethSignatureChallenges;
    AddressChecks senderChecks;
    AddressChecks recipientChecks;
    AddressChecks initiatorChecks;
    AltTimeChecks altTimeChecks;
    bool mustPrioritize;
    VotingChallenge[] votingChallenges;
    bool allowBackedMinting;
    bool allowSpecialWrapping;
}

/**
 * @notice Criteria for outgoing approvals
 */
struct OutgoingApprovalCriteria {
    MerkleChallenge[] merkleChallenges;
    PredeterminedBalances predeterminedBalances;
    ApprovalAmounts approvalAmounts;
    MaxNumTransfers maxNumTransfers;
    CoinTransfer[] coinTransfers;
    bool requireToEqualsInitiatedBy;
    bool requireToDoesNotEqualInitiatedBy;
    AutoDeletionOptions autoDeletionOptions;
    MustOwnTokens[] mustOwnTokens;
    DynamicStoreChallenge[] dynamicStoreChallenges;
    ETHSignatureChallenge[] ethSignatureChallenges;
    AddressChecks recipientChecks;
    AddressChecks initiatorChecks;
    AltTimeChecks altTimeChecks;
    bool mustPrioritize;
    VotingChallenge[] votingChallenges;
}

/**
 * @notice Criteria for incoming approvals
 */
struct IncomingApprovalCriteria {
    MerkleChallenge[] merkleChallenges;
    PredeterminedBalances predeterminedBalances;
    ApprovalAmounts approvalAmounts;
    MaxNumTransfers maxNumTransfers;
    CoinTransfer[] coinTransfers;
    bool requireFromEqualsInitiatedBy;
    bool requireFromDoesNotEqualInitiatedBy;
    AutoDeletionOptions autoDeletionOptions;
    MustOwnTokens[] mustOwnTokens;
    DynamicStoreChallenge[] dynamicStoreChallenges;
    ETHSignatureChallenge[] ethSignatureChallenges;
    AddressChecks senderChecks;
    AddressChecks initiatorChecks;
    AltTimeChecks altTimeChecks;
    bool mustPrioritize;
    VotingChallenge[] votingChallenges;
}

// ============================================================================
// Approval Types (approvals.proto)
// ============================================================================

/**
 * @notice Collection-level approval
 */
struct CollectionApproval {
    string fromListId;
    string toListId;
    string initiatedByListId;
    UintRange[] transferTimes;
    UintRange[] tokenIds;
    UintRange[] ownershipTimes;
    string uri;
    string customData;
    string approvalId;
    ApprovalCriteria approvalCriteria;
    uint256 version;
}

/**
 * @notice User incoming approval
 */
struct UserIncomingApproval {
    string fromListId;
    string initiatedByListId;
    UintRange[] transferTimes;
    UintRange[] tokenIds;
    UintRange[] ownershipTimes;
    string uri;
    string customData;
    string approvalId;
    IncomingApprovalCriteria approvalCriteria;
    uint256 version;
}

/**
 * @notice User outgoing approval
 */
struct UserOutgoingApproval {
    string toListId;
    string initiatedByListId;
    UintRange[] transferTimes;
    UintRange[] tokenIds;
    UintRange[] ownershipTimes;
    string uri;
    string customData;
    string approvalId;
    OutgoingApprovalCriteria approvalCriteria;
    uint256 version;
}

/**
 * @notice Approval identifier details
 */
struct ApprovalIdentifierDetails {
    string approvalId;
    string approvalLevel;
    string approverAddress;
    uint256 version;
}

// ============================================================================
// Collection Types (collections.proto)
// ============================================================================

/**
 * @notice Conversion side A (amount only)
 */
struct ConversionSideA {
    uint256 amount;
}

/**
 * @notice Conversion side A with denom
 */
struct ConversionSideAWithDenom {
    uint256 amount;
    string denom;
}

/**
 * @notice Conversion without denom
 */
struct ConversionWithoutDenom {
    ConversionSideA sideA;
    Balance[] sideB;
}

/**
 * @notice Conversion with denom
 */
struct Conversion {
    ConversionSideAWithDenom sideA;
    Balance[] sideB;
}

/**
 * @notice Denomination unit
 */
struct DenomUnit {
    uint256 decimals;
    string symbol;
    bool isDefaultDisplay;
    PathMetadata metadata;
}

/**
 * @notice Cosmos coin wrapper path
 */
struct CosmosCoinWrapperPath {
    string wrapperAddress; // Renamed from 'address' to avoid Solidity reserved keyword
    string denom;
    ConversionWithoutDenom conversion;
    string symbol;
    DenomUnit[] denomUnits;
    bool allowOverrideWithAnyValidToken;
    PathMetadata metadata;
}

/**
 * @notice Alias path
 */
struct AliasPath {
    string denom;
    ConversionWithoutDenom conversion;
    string symbol;
    DenomUnit[] denomUnits;
    PathMetadata metadata;
}

/**
 * @notice Cosmos coin backed path
 */
struct CosmosCoinBackedPath {
    string addr; // Renamed from 'address' to avoid Solidity reserved keyword
    Conversion conversion;
}

/**
 * @notice Collection invariants
 */
struct CollectionInvariants {
    bool noCustomOwnershipTimes;
    uint256 maxSupplyPerId;
    CosmosCoinBackedPath cosmosCoinBackedPath;
    bool noForcefulPostMintTransfers;
    bool disablePoolCreation;
}

/**
 * @notice Token collection
 * @dev NOTE: createdBy field is NOT included - it's always the caller
 */
struct TokenCollection {
    uint256 collectionId;
    CollectionMetadata collectionMetadata;
    TokenMetadata[] tokenMetadata;
    string customData;
    string manager;
    CollectionPermissions collectionPermissions;
    CollectionApproval[] collectionApprovals;
    string[] standards;
    bool isArchived;
    UserBalanceStore defaultBalances;
    // createdBy field excluded - always msg.sender
    UintRange[] validTokenIds;
    string mintEscrowAddress;
    CosmosCoinWrapperPath[] cosmosCoinWrapperPaths;
    CollectionInvariants invariants;
    AliasPath[] aliasPaths;
}

// ============================================================================
// Transfer Types (transfers.proto)
// ============================================================================

/**
 * @notice Precalculate balances from approval details
 */
struct PrecalculateBalancesFromApprovalDetails {
    string approvalId;
    string approvalLevel;
    string approverAddress;
    uint256 version;
    PrecalculationOptions precalculationOptions;
}

/**
 * @notice Transfer details
 * @dev NOTE: from field is NOT included - it's always the caller
 */
struct Transfer {
    // from field excluded - always msg.sender
    string[] toAddresses;
    Balance[] balances;
    PrecalculateBalancesFromApprovalDetails precalculateBalancesFromApproval;
    MerkleProof[] merkleProofs;
    ETHSignatureProof[] ethSignatureProofs;
    string memo;
    ApprovalIdentifierDetails[] prioritizedApprovals;
    bool onlyCheckPrioritizedCollectionApprovals;
    bool onlyCheckPrioritizedIncomingApprovals;
    bool onlyCheckPrioritizedOutgoingApprovals;
}

// ============================================================================
// Address List Types (address_lists.proto)
// ============================================================================

/**
 * @notice Address list
 */
struct AddressList {
    string listId;
    string[] addresses;
    bool whitelist;
    string uri;
    string customData;
    string createdBy;
}

/**
 * @notice Address list input (for creation)
 * @dev NOTE: createdBy field is NOT included - it's always the caller
 */
struct AddressListInput {
    string listId;
    string[] addresses;
    bool whitelist;
    string uri;
    string customData;
    // createdBy field excluded - always msg.sender
}

// ============================================================================
// Dynamic Store Types (dynamic_stores.proto)
// ============================================================================

/**
 * @notice Dynamic store
 */
struct DynamicStore {
    uint256 storeId;
    string createdBy;
    bool defaultValue;
    bool globalEnabled;
    string uri;
    string customData;
}

/**
 * @notice Dynamic store value
 */
struct DynamicStoreValue {
    uint256 storeId;
    string address_;
    bool value;
}

// ============================================================================
// User Balance Store Types (user_balance_store.proto)
// ============================================================================

/**
 * @notice User balance store
 */
struct UserBalanceStore {
    Balance[] balances;
    UserOutgoingApproval[] outgoingApprovals;
    UserIncomingApproval[] incomingApprovals;
    bool autoApproveSelfInitiatedOutgoingTransfers;
    bool autoApproveSelfInitiatedIncomingTransfers;
    bool autoApproveAllIncomingTransfers;
    UserPermissions userPermissions;
}

// ============================================================================
// Params Types (params.proto)
// ============================================================================

/**
 * @notice Module parameters
 */
struct Params {
    string[] allowedDenoms;
    uint256 affiliatePercentage;
}

// ============================================================================
// Timeline Types (timelines.proto)
// ============================================================================

/**
 * @notice Collection metadata timeline
 */
struct CollectionMetadataTimeline {
    CollectionMetadata collectionMetadata;
    UintRange[] timelineTimes;
}

/**
 * @notice Token metadata timeline
 */
struct TokenMetadataTimeline {
    TokenMetadata[] tokenMetadata;
    UintRange[] timelineTimes;
}

/**
 * @notice Custom data timeline
 */
struct CustomDataTimeline {
    string customData;
    UintRange[] timelineTimes;
}

/**
 * @notice Manager timeline
 */
struct ManagerTimeline {
    string manager;
    UintRange[] timelineTimes;
}

/**
 * @notice Is archived timeline
 */
struct IsArchivedTimeline {
    bool isArchived;
    UintRange[] timelineTimes;
}

/**
 * @notice Standards timeline
 */
struct StandardsTimeline {
    string[] standards;
    UintRange[] timelineTimes;
}

// ============================================================================
// Query Request/Response Types (query.proto)
// ============================================================================

/**
 * @notice Query params request
 * @dev Empty structs are not allowed in Solidity, so we use a dummy field
 */
struct QueryParamsRequest {
    bool _dummy; // Dummy field - Solidity doesn't allow empty structs
}

/**
 * @notice Query params response
 */
struct QueryParamsResponse {
    Params params;
}

/**
 * @notice Query get collection request
 */
struct QueryGetCollectionRequest {
    uint256 collectionId;
}

/**
 * @notice Query get collection response
 */
struct QueryGetCollectionResponse {
    TokenCollection collection;
}

/**
 * @notice Query get balance request
 */
struct QueryGetBalanceRequest {
    uint256 collectionId;
    string addr; // Renamed from 'address' to avoid Solidity reserved keyword
}

/**
 * @notice Query get balance response
 */
struct QueryGetBalanceResponse {
    UserBalanceStore balance;
}

/**
 * @notice Query get address list request
 */
struct QueryGetAddressListRequest {
    string listId;
}

/**
 * @notice Query get address list response
 */
struct QueryGetAddressListResponse {
    AddressList list;
}

/**
 * @notice Query get approval tracker request
 */
struct QueryGetApprovalTrackerRequest {
    string amountTrackerId;
    string approvalLevel;
    string approverAddress;
    string trackerType;
    uint256 collectionId;
    string approvedAddress;
    string approvalId;
}

/**
 * @notice Query get approval tracker response
 */
struct QueryGetApprovalTrackerResponse {
    ApprovalTracker tracker;
}

/**
 * @notice Query get challenge tracker request
 */
struct QueryGetChallengeTrackerRequest {
    uint256 collectionId;
    string approvalLevel;
    string approverAddress;
    string challengeTrackerId;
    uint256 leafIndex;
    string approvalId;
}

/**
 * @notice Query get challenge tracker response
 */
struct QueryGetChallengeTrackerResponse {
    uint256 numUsed;
}

/**
 * @notice Query get ETH signature tracker request
 */
struct QueryGetETHSignatureTrackerRequest {
    uint256 collectionId;
    string approvalLevel;
    string approverAddress;
    string approvalId;
    string challengeTrackerId;
    string signature;
}

/**
 * @notice Query get ETH signature tracker response
 */
struct QueryGetETHSignatureTrackerResponse {
    uint256 numUsed;
}

/**
 * @notice Query get dynamic store request
 */
struct QueryGetDynamicStoreRequest {
    uint256 storeId;
}

/**
 * @notice Query get dynamic store response
 */
struct QueryGetDynamicStoreResponse {
    DynamicStore store;
}

/**
 * @notice Query get dynamic store value request
 */
struct QueryGetDynamicStoreValueRequest {
    uint256 storeId;
    string addr; // Renamed from 'address' to avoid Solidity reserved keyword
}

/**
 * @notice Query get dynamic store value response
 */
struct QueryGetDynamicStoreValueResponse {
    DynamicStoreValue value;
}

/**
 * @notice Query get wrappable balances request
 */
struct QueryGetWrappableBalancesRequest {
    string denom;
    string addr; // Renamed from 'address' to avoid Solidity reserved keyword
}

/**
 * @notice Query get wrappable balances response
 */
struct QueryGetWrappableBalancesResponse {
    uint256 amount;
}

/**
 * @notice Query is address reserved protocol request
 */
struct QueryIsAddressReservedProtocolRequest {
    string addr; // Renamed from 'address' to avoid Solidity reserved keyword
}

/**
 * @notice Query is address reserved protocol response
 */
struct QueryIsAddressReservedProtocolResponse {
    bool isReservedProtocol;
}

/**
 * @notice Query get all reserved protocol addresses request
 * @dev Empty struct - using dummy field to satisfy Solidity compiler
 */
struct QueryGetAllReservedProtocolAddressesRequest {
    bool _dummy; // Dummy field - struct cannot be empty in Solidity
}

/**
 * @notice Query get all reserved protocol addresses response
 */
struct QueryGetAllReservedProtocolAddressesResponse {
    string[] addresses;
}

/**
 * @notice Query get vote request
 */
struct QueryGetVoteRequest {
    uint256 collectionId;
    string approvalLevel;
    string approverAddress;
    string approvalId;
    string proposalId;
    string voterAddress;
}

/**
 * @notice Query get vote response
 */
struct QueryGetVoteResponse {
    VoteProof vote;
}

/**
 * @notice Query get votes request
 */
struct QueryGetVotesRequest {
    uint256 collectionId;
    string approvalLevel;
    string approverAddress;
    string approvalId;
    string proposalId;
}

/**
 * @notice Query get votes response
 */
struct QueryGetVotesResponse {
    VoteProof[] votes;
}

// ============================================================================
// Msg Types (tx.proto) - Input types only (creator field excluded)
// ============================================================================

/**
 * @notice Msg create collection
 * @dev NOTE: creator field is NOT included - it's always msg.sender
 */
struct MsgCreateCollection {
    // creator field excluded - always msg.sender
    UserBalanceStore defaultBalances;
    UintRange[] validTokenIds;
    CollectionPermissions collectionPermissions;
    string manager;
    CollectionMetadata collectionMetadata;
    TokenMetadata[] tokenMetadata;
    string customData;
    CollectionApproval[] collectionApprovals;
    string[] standards;
    bool isArchived;
    // Note: mintEscrowCoinsToTransfer would need Cosmos Coin type handling
    CosmosCoinWrapperPath[] cosmosCoinWrapperPathsToAdd;
    CollectionInvariants invariants;
    AliasPath[] aliasPathsToAdd;
}

/**
 * @notice Msg update collection
 * @dev NOTE: creator field is NOT included - it's always msg.sender
 */
struct MsgUpdateCollection {
    // creator field excluded - always msg.sender
    uint256 collectionId;
    bool updateValidTokenIds;
    UintRange[] validTokenIds;
    bool updateCollectionPermissions;
    CollectionPermissions collectionPermissions;
    bool updateManager;
    string manager;
    bool updateCollectionMetadata;
    CollectionMetadata collectionMetadata;
    bool updateTokenMetadata;
    TokenMetadata[] tokenMetadata;
    bool updateCustomData;
    string customData;
    bool updateCollectionApprovals;
    CollectionApproval[] collectionApprovals;
    bool updateStandards;
    string[] standards;
    bool updateIsArchived;
    bool isArchived;
    // Note: mintEscrowCoinsToTransfer would need Cosmos Coin type handling
    CosmosCoinWrapperPath[] cosmosCoinWrapperPathsToAdd;
    CollectionInvariants invariants;
    AliasPath[] aliasPathsToAdd;
}

/**
 * @notice Msg universal update collection
 * @dev NOTE: creator field is NOT included - it's always msg.sender
 */
struct MsgUniversalUpdateCollection {
    // creator field excluded - always msg.sender
    uint256 collectionId;
    UserBalanceStore defaultBalances;
    bool updateValidTokenIds;
    UintRange[] validTokenIds;
    bool updateCollectionPermissions;
    CollectionPermissions collectionPermissions;
    bool updateManager;
    string manager;
    bool updateCollectionMetadata;
    CollectionMetadata collectionMetadata;
    bool updateTokenMetadata;
    TokenMetadata[] tokenMetadata;
    bool updateCustomData;
    string customData;
    bool updateCollectionApprovals;
    CollectionApproval[] collectionApprovals;
    bool updateStandards;
    string[] standards;
    bool updateIsArchived;
    bool isArchived;
    // Note: mintEscrowCoinsToTransfer would need Cosmos Coin type handling
    CosmosCoinWrapperPath[] cosmosCoinWrapperPathsToAdd;
    CollectionInvariants invariants;
    AliasPath[] aliasPathsToAdd;
}

/**
 * @notice Msg delete collection
 * @dev NOTE: creator field is NOT included - it's always msg.sender
 */
struct MsgDeleteCollection {
    // creator field excluded - always msg.sender
    uint256 collectionId;
}

/**
 * @notice Msg create address lists
 * @dev NOTE: creator field is NOT included - it's always msg.sender
 */
struct MsgCreateAddressLists {
    // creator field excluded - always msg.sender
    AddressListInput[] addressLists;
}

/**
 * @notice Msg transfer tokens
 * @dev NOTE: creator field is NOT included - it's always msg.sender
 */
struct MsgTransferTokens {
    // creator field excluded - always msg.sender
    uint256 collectionId;
    Transfer[] transfers;
}

/**
 * @notice Msg update user approvals
 * @dev NOTE: creator field is NOT included - it's always msg.sender
 */
struct MsgUpdateUserApprovals {
    // creator field excluded - always msg.sender
    uint256 collectionId;
    bool updateOutgoingApprovals;
    UserOutgoingApproval[] outgoingApprovals;
    bool updateIncomingApprovals;
    UserIncomingApproval[] incomingApprovals;
    bool updateAutoApproveSelfInitiatedOutgoingTransfers;
    bool autoApproveSelfInitiatedOutgoingTransfers;
    bool updateAutoApproveSelfInitiatedIncomingTransfers;
    bool autoApproveSelfInitiatedIncomingTransfers;
    bool updateAutoApproveAllIncomingTransfers;
    bool autoApproveAllIncomingTransfers;
    bool updateUserPermissions;
    UserPermissions userPermissions;
}

/**
 * @notice Msg set incoming approval
 * @dev NOTE: creator field is NOT included - it's always msg.sender
 */
struct MsgSetIncomingApproval {
    // creator field excluded - always msg.sender
    uint256 collectionId;
    UserIncomingApproval approval;
}

/**
 * @notice Msg delete incoming approval
 * @dev NOTE: creator field is NOT included - it's always msg.sender
 */
struct MsgDeleteIncomingApproval {
    // creator field excluded - always msg.sender
    uint256 collectionId;
    string approvalId;
}

/**
 * @notice Msg set outgoing approval
 * @dev NOTE: creator field is NOT included - it's always msg.sender
 */
struct MsgSetOutgoingApproval {
    // creator field excluded - always msg.sender
    uint256 collectionId;
    UserOutgoingApproval approval;
}

/**
 * @notice Msg delete outgoing approval
 * @dev NOTE: creator field is NOT included - it's always msg.sender
 */
struct MsgDeleteOutgoingApproval {
    // creator field excluded - always msg.sender
    uint256 collectionId;
    string approvalId;
}

/**
 * @notice Msg purge approvals
 * @dev NOTE: creator field is NOT included - it's always msg.sender
 */
struct MsgPurgeApprovals {
    // creator field excluded - always msg.sender
    uint256 collectionId;
    bool purgeExpired;
    string approverAddress;
    bool purgeCounterpartyApprovals;
    ApprovalIdentifierDetails[] approvalsToPurge;
}

/**
 * @notice Msg create dynamic store
 * @dev NOTE: creator field is NOT included - it's always msg.sender
 */
struct MsgCreateDynamicStore {
    // creator field excluded - always msg.sender
    bool defaultValue;
    string uri;
    string customData;
}

/**
 * @notice Msg update dynamic store
 * @dev NOTE: creator field is NOT included - it's always msg.sender
 */
struct MsgUpdateDynamicStore {
    // creator field excluded - always msg.sender
    uint256 storeId;
    bool defaultValue;
    bool globalEnabled;
    string uri;
    string customData;
}

/**
 * @notice Msg delete dynamic store
 * @dev NOTE: creator field is NOT included - it's always msg.sender
 */
struct MsgDeleteDynamicStore {
    // creator field excluded - always msg.sender
    uint256 storeId;
}

/**
 * @notice Msg set dynamic store value
 * @dev NOTE: creator field is NOT included - it's always msg.sender
 */
struct MsgSetDynamicStoreValue {
    // creator field excluded - always msg.sender
    uint256 storeId;
    string addr; // Renamed from 'address' to avoid Solidity reserved keyword
    bool value;
}

/**
 * @notice Msg set valid token IDs
 * @dev NOTE: creator field is NOT included - it's always msg.sender
 */
struct MsgSetValidTokenIds {
    // creator field excluded - always msg.sender
    uint256 collectionId;
    UintRange[] validTokenIds;
    TokenIdsActionPermission[] canUpdateValidTokenIds;
}

/**
 * @notice Msg set manager
 * @dev NOTE: creator field is NOT included - it's always msg.sender
 */
struct MsgSetManager {
    // creator field excluded - always msg.sender
    uint256 collectionId;
    string manager;
    ActionPermission[] canUpdateManager;
}

/**
 * @notice Msg set collection metadata
 * @dev NOTE: creator field is NOT included - it's always msg.sender
 */
struct MsgSetCollectionMetadata {
    // creator field excluded - always msg.sender
    uint256 collectionId;
    CollectionMetadata collectionMetadata;
    ActionPermission[] canUpdateCollectionMetadata;
}

/**
 * @notice Msg set token metadata
 * @dev NOTE: creator field is NOT included - it's always msg.sender
 */
struct MsgSetTokenMetadata {
    // creator field excluded - always msg.sender
    uint256 collectionId;
    TokenMetadata[] tokenMetadata;
    TokenIdsActionPermission[] canUpdateTokenMetadata;
}

/**
 * @notice Msg set custom data
 * @dev NOTE: creator field is NOT included - it's always msg.sender
 */
struct MsgSetCustomData {
    // creator field excluded - always msg.sender
    uint256 collectionId;
    string customData;
    ActionPermission[] canUpdateCustomData;
}

/**
 * @notice Msg set standards
 * @dev NOTE: creator field is NOT included - it's always msg.sender
 */
struct MsgSetStandards {
    // creator field excluded - always msg.sender
    uint256 collectionId;
    string[] standards;
    ActionPermission[] canUpdateStandards;
}

/**
 * @notice Msg set collection approvals
 * @dev NOTE: creator field is NOT included - it's always msg.sender
 */
struct MsgSetCollectionApprovals {
    // creator field excluded - always msg.sender
    uint256 collectionId;
    CollectionApproval[] collectionApprovals;
    CollectionApprovalPermission[] canUpdateCollectionApprovals;
}

/**
 * @notice Msg set is archived
 * @dev NOTE: creator field is NOT included - it's always msg.sender
 */
struct MsgSetIsArchived {
    // creator field excluded - always msg.sender
    uint256 collectionId;
    bool isArchived;
    ActionPermission[] canArchiveCollection;
}

/**
 * @notice Msg set reserved protocol address
 * @dev NOTE: authority field is NOT included - handled separately for governance
 */
struct MsgSetReservedProtocolAddress {
    // authority field excluded - handled separately
    string addr; // Renamed from 'address' to avoid Solidity reserved keyword
    bool isReservedProtocol;
}

/**
 * @notice Msg cast vote
 * @dev NOTE: creator field is NOT included - it's always msg.sender
 */
struct MsgCastVote {
    // creator field excluded - always msg.sender
    uint256 collectionId;
    string approvalLevel;
    string approverAddress;
    string approvalId;
    string proposalId;
    uint256 yesWeight;
}

/**
 * @notice Msg update params
 * @dev NOTE: authority field is NOT included - handled separately for governance
 */
struct MsgUpdateParams {
    // authority field excluded - handled separately
    Params params;
}

// ============================================================================
// Msg Response Types
// ============================================================================

/**
 * @notice Response for create collection
 */
struct MsgCreateCollectionResponse {
    uint256 collectionId;
}

/**
 * @notice Response for update collection
 */
struct MsgUpdateCollectionResponse {
    uint256 collectionId;
}

/**
 * @notice Response for universal update collection
 */
struct MsgUniversalUpdateCollectionResponse {
    uint256 collectionId;
}

/**
 * @notice Response for purge approvals
 */
struct MsgPurgeApprovalsResponse {
    uint256 numPurged;
}

/**
 * @notice Response for create dynamic store
 */
struct MsgCreateDynamicStoreResponse {
    uint256 storeId;
}

// ============================================================================
// Helper Types for Msg Creation (from tx.proto)
// ============================================================================

/**
 * @notice Invariants add object (for collection creation/update)
 */
struct InvariantsAddObject {
    bool noCustomOwnershipTimes;
    uint256 maxSupplyPerId;
    CosmosCoinBackedPath cosmosCoinBackedPath;
    bool noForcefulPostMintTransfers;
    bool disablePoolCreation;
}

/**
 * @notice Cosmos coin wrapper path add object
 */
struct CosmosCoinWrapperPathAddObject {
    string denom;
    ConversionWithoutDenom conversion;
    string symbol;
    DenomUnit[] denomUnits;
    bool allowOverrideWithAnyValidToken;
    PathMetadata metadata;
}

/**
 * @notice Alias path add object
 */
struct AliasPathAddObject {
    string denom;
    ConversionWithoutDenom conversion;
    string symbol;
    DenomUnit[] denomUnits;
    PathMetadata metadata;
}

/**
 * @notice Cosmos coin backed path add object
 */
struct CosmosCoinBackedPathAddObject {
    Conversion conversion;
}

