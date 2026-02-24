// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../types/TokenizationTypes.sol";

/// @title ITokenizationPrecompile
/// @notice Interface for the BitBadges tokenization precompile
/// @dev Precompile address: 0x0000000000000000000000000000000000001001
///      All methods use JSON string parameters matching protobuf JSON format.
///      The caller address (creator/sender) is automatically set from msg.sender.
///      Use helper libraries to construct JSON strings from Solidity types.
interface ITokenizationPrecompile {
    // ============ Types ============
    
    /// @notice Input structure for multi-message execution
    struct MessageInput {
        string messageType;  // e.g., "createCollection", "transferTokens"
        string msgJson;      // JSON matching the protobuf format
    }
    // ============ Events ============

    /// @notice Emitted when tokens are transferred
    /// @param collectionId The collection ID
    /// @param from The sender address
    /// @param to The recipient addresses
    /// @param amount The amount transferred
    event TransferTokens(
        uint256 indexed collectionId,
        address indexed from,
        address[] to,
        uint256 amount
    );

    /// @notice Emitted when an incoming approval is set
    /// @param collectionId The collection ID
    /// @param from The address setting the approval
    /// @param approvalId The approval ID
    event SetIncomingApproval(
        uint256 indexed collectionId,
        address indexed from,
        string approvalId
    );

    /// @notice Emitted when an outgoing approval is set
    /// @param collectionId The collection ID
    /// @param from The address setting the approval
    /// @param approvalId The approval ID
    event SetOutgoingApproval(
        uint256 indexed collectionId,
        address indexed from,
        string approvalId
    );

    /// @notice Emitted when a collection is created
    /// @param collectionId The new collection ID
    /// @param creator The creator address
    event CollectionCreated(
        uint256 indexed collectionId,
        address indexed creator
    );

    /// @notice Emitted when a collection is updated
    /// @param collectionId The collection ID
    /// @param updater The updater address
    event CollectionUpdated(
        uint256 indexed collectionId,
        address indexed updater
    );

    /// @notice Emitted when a collection is deleted
    /// @param collectionId The collection ID
    /// @param deleter The deleter address
    event CollectionDeleted(
        uint256 indexed collectionId,
        address indexed deleter
    );

    /// @notice Emitted when address lists are created
    /// @param creator The creator address
    /// @param listCount The number of lists created
    event AddressListsCreated(
        address indexed creator,
        uint256 listCount
    );

    /// @notice Emitted when a dynamic store is created
    /// @param storeId The new store ID
    /// @param creator The creator address
    event DynamicStoreCreated(
        uint256 indexed storeId,
        address indexed creator
    );

    /// @notice Emitted when a vote is cast
    /// @param collectionId The collection ID
    /// @param voter The voter address
    /// @param proposalId The proposal ID
    /// @param yesWeight The weight of the yes vote
    event VoteCast(
        uint256 indexed collectionId,
        address indexed voter,
        string proposalId,
        uint256 yesWeight
    );

    // ============ Transactions ============
    // NOTE: All methods now use a JSON string parameter (msgJson) instead of individual parameters.
    // The JSON must match the protobuf JSON format for the corresponding Msg type.
    // Use helper libraries or construct JSON manually. The caller address is automatically set from msg.sender.

    /// @notice Transfer tokens from the caller to specified addresses
    /// @param msgJson JSON string matching MsgTransferTokens protobuf format
    /// @return success True if transfer succeeded
    function transferTokens(
        string calldata msgJson
    ) external returns (bool success);

    /// @notice Set an incoming approval for the caller
    /// @param msgJson JSON string matching MsgSetIncomingApproval protobuf format
    /// @return success True if approval was set
    function setIncomingApproval(
        string calldata msgJson
    ) external returns (bool success);

    /// @notice Set an outgoing approval for the caller
    /// @param msgJson JSON string matching MsgSetOutgoingApproval protobuf format
    /// @return success True if approval was set
    function setOutgoingApproval(
        string calldata msgJson
    ) external returns (bool success);

    /// @notice Delete a collection (creator only)
    /// @param msgJson JSON string matching MsgDeleteCollection protobuf format
    /// @return success True if deletion succeeded
    function deleteCollection(
        string calldata msgJson
    ) external returns (bool success);

    /// @notice Delete an incoming approval
    /// @param msgJson JSON string matching MsgDeleteIncomingApproval protobuf format
    /// @return success True if deletion succeeded
    function deleteIncomingApproval(
        string calldata msgJson
    ) external returns (bool success);

    /// @notice Delete an outgoing approval
    /// @param msgJson JSON string matching MsgDeleteOutgoingApproval protobuf format
    /// @return success True if deletion succeeded
    function deleteOutgoingApproval(
        string calldata msgJson
    ) external returns (bool success);

    /// @notice Create a new dynamic store
    /// @param msgJson JSON string matching MsgCreateDynamicStore protobuf format
    /// @return storeId The newly created store ID
    function createDynamicStore(
        string calldata msgJson
    ) external returns (uint256 storeId);

    /// @notice Update an existing dynamic store
    /// @param msgJson JSON string matching MsgUpdateDynamicStore protobuf format
    /// @return success True if update succeeded
    function updateDynamicStore(
        string calldata msgJson
    ) external returns (bool success);

    /// @notice Delete a dynamic store (creator only)
    /// @param msgJson JSON string matching MsgDeleteDynamicStore protobuf format
    /// @return success True if deletion succeeded
    function deleteDynamicStore(
        string calldata msgJson
    ) external returns (bool success);

    /// @notice Set a value in a dynamic store for an address
    /// @param msgJson JSON string matching MsgSetDynamicStoreValue protobuf format
    /// @return success True if value was set
    function setDynamicStoreValue(
        string calldata msgJson
    ) external returns (bool success);

    /// @notice Set custom data for a collection
    /// @param msgJson JSON string matching MsgSetCustomData protobuf format
    /// @return resultCollectionId The collection ID (unchanged)
    function setCustomData(
        string calldata msgJson
    ) external returns (uint256 resultCollectionId);

    /// @notice Set whether a collection is archived
    /// @param msgJson JSON string matching MsgSetIsArchived protobuf format
    /// @return resultCollectionId The collection ID (unchanged)
    function setIsArchived(
        string calldata msgJson
    ) external returns (uint256 resultCollectionId);

    /// @notice Set the manager of a collection
    /// @param msgJson JSON string matching MsgSetManager protobuf format
    /// @return resultCollectionId The collection ID (unchanged)
    function setManager(
        string calldata msgJson
    ) external returns (uint256 resultCollectionId);

    /// @notice Set collection metadata
    /// @param msgJson JSON string matching MsgSetCollectionMetadata protobuf format
    /// @return resultCollectionId The collection ID (unchanged)
    function setCollectionMetadata(
        string calldata msgJson
    ) external returns (uint256 resultCollectionId);

    /// @notice Set the standards for a collection
    /// @param msgJson JSON string matching MsgSetStandards protobuf format
    /// @return resultCollectionId The collection ID (unchanged)
    function setStandards(
        string calldata msgJson
    ) external returns (uint256 resultCollectionId);

    /// @notice Cast a vote for a proposal in a voting challenge
    /// @param msgJson JSON string matching MsgCastVote protobuf format
    /// @return success True if vote was cast
    function castVote(
        string calldata msgJson
    ) external returns (bool success);

    /// @notice Create a new token collection
    /// @param msgJson JSON string matching MsgCreateCollection protobuf format
    /// @return newCollectionId The newly created collection ID
    function createCollection(
        string calldata msgJson
    ) external returns (uint256 newCollectionId);

    /// @notice Update an existing collection
    /// @param msgJson JSON string matching MsgUpdateCollection protobuf format
    /// @return resultCollectionId The collection ID (unchanged)
    function updateCollection(
        string calldata msgJson
    ) external returns (uint256 resultCollectionId);

    /// @notice Update user approvals and permissions
    /// @param msgJson JSON string matching MsgUpdateUserApprovals protobuf format
    /// @return success True if update succeeded
    function updateUserApprovals(
        string calldata msgJson
    ) external returns (bool success);

    /// @notice Create one or more address lists
    /// @param msgJson JSON string matching MsgCreateAddressLists protobuf format
    /// @return success True if lists were created
    function createAddressLists(
        string calldata msgJson
    ) external returns (bool success);

    /// @notice Purge expired or specified approvals
    /// @param msgJson JSON string matching MsgPurgeApprovals protobuf format
    /// @return numPurged Number of approvals purged
    function purgeApprovals(
        string calldata msgJson
    ) external returns (uint256 numPurged);

    /// @notice Set valid token IDs for a collection
    /// @param msgJson JSON string matching MsgSetValidTokenIds protobuf format
    /// @return resultCollectionId The collection ID (unchanged)
    function setValidTokenIds(
        string calldata msgJson
    ) external returns (uint256 resultCollectionId);

    /// @notice Set token metadata for specific token IDs
    /// @param msgJson JSON string matching MsgSetTokenMetadata protobuf format
    /// @return resultCollectionId The collection ID (unchanged)
    function setTokenMetadata(
        string calldata msgJson
    ) external returns (uint256 resultCollectionId);

    /// @notice Set collection-level approvals
    /// @param msgJson JSON string matching MsgSetCollectionApprovals protobuf format
    /// @return resultCollectionId The collection ID (unchanged)
    function setCollectionApprovals(
        string calldata msgJson
    ) external returns (uint256 resultCollectionId);

    /// @notice Universal update for all collection properties
    /// @dev Combines all update operations into a single call
    /// @param msgJson JSON string matching MsgUniversalUpdateCollection protobuf format
    /// @return resultCollectionId The collection ID (unchanged)
    function universalUpdateCollection(
        string calldata msgJson
    ) external returns (uint256 resultCollectionId);

    /// @notice Execute multiple messages sequentially in a single atomic transaction
    /// @dev All messages execute in order. If any message fails, the entire transaction is rolled back.
    ///      Each message's creator field is automatically set from msg.sender.
    ///      Results are returned as bytes array - decode based on message type.
    /// @param messages Array of MessageInput structs with messageType and msgJson
    /// @return success True if all messages executed successfully
    /// @return results Array of result bytes (decode based on message type)
    function executeMultiple(
        MessageInput[] calldata messages
    ) external returns (bool success, bytes[] memory results);
    
    // ============ Queries ============
    // NOTE: All query methods now use a JSON string parameter (msgJson) instead of individual parameters.
    // The JSON must match the protobuf JSON format for the corresponding QueryRequest type.

    /// @notice Get collection details by ID
    /// @dev Returns protobuf-encoded TokenCollection. Decode using appropriate codec.
    /// @param msgJson JSON string matching QueryCollectionRequest protobuf format
    /// @return collection Protobuf-encoded TokenCollection bytes
    function getCollection(
        string calldata msgJson
    ) external view returns (bytes memory collection);

    /// @notice Get collection stats (e.g. holder count) by ID
    /// @dev Returns protobuf-encoded QueryGetCollectionStatsResponse (stats field = CollectionStats with holderCount).
    /// @param msgJson JSON string matching QueryGetCollectionStatsRequest, e.g. {"collectionId":"1"}
    /// @return stats Protobuf-encoded CollectionStats bytes
    function getCollectionStats(
        string calldata msgJson
    ) external view returns (bytes memory stats);

    /// @notice Get user balance for a collection
    /// @dev Returns protobuf-encoded UserBalanceStore. Decode using appropriate codec.
    /// @param msgJson JSON string matching QueryBalanceRequest protobuf format
    /// @return balance Protobuf-encoded UserBalanceStore bytes
    function getBalance(
        string calldata msgJson
    ) external view returns (bytes memory balance);

    /// @notice Get an address list by ID
    /// @dev Returns protobuf-encoded AddressList. Decode using appropriate codec.
    /// @param msgJson JSON string matching QueryAddressListRequest protobuf format
    /// @return list Protobuf-encoded AddressList bytes
    function getAddressList(
        string calldata msgJson
    ) external view returns (bytes memory list);
    
    /// @notice Get approval tracker for amount/transfer tracking
    /// @dev Returns protobuf-encoded ApprovalTracker
    /// @param msgJson JSON string matching QueryApprovalTrackerRequest protobuf format
    /// @return Protobuf-encoded ApprovalTracker bytes
    function getApprovalTracker(
        string calldata msgJson
    ) external view returns (bytes memory);

    /// @notice Get challenge tracker for Merkle challenges
    /// @dev Returns protobuf-encoded challenge tracker data
    /// @param msgJson JSON string matching QueryChallengeTrackerRequest protobuf format
    /// @return Protobuf-encoded challenge tracker bytes
    function getChallengeTracker(
        string calldata msgJson
    ) external view returns (bytes memory);

    /// @notice Get ETH signature tracker
    /// @dev Returns protobuf-encoded signature tracker data
    /// @param msgJson JSON string matching QueryETHSignatureTrackerRequest protobuf format
    /// @return Protobuf-encoded signature tracker bytes
    function getETHSignatureTracker(
        string calldata msgJson
    ) external view returns (bytes memory);

    /// @notice Get a dynamic store by ID
    /// @dev Returns protobuf-encoded DynamicStore
    /// @param msgJson JSON string matching QueryDynamicStoreRequest protobuf format
    /// @return Protobuf-encoded DynamicStore bytes
    function getDynamicStore(
        string calldata msgJson
    ) external view returns (bytes memory);

    /// @notice Get a dynamic store value for a specific address
    /// @dev Returns protobuf-encoded DynamicStoreValue
    /// @param msgJson JSON string matching QueryDynamicStoreValueRequest protobuf format
    /// @return Protobuf-encoded DynamicStoreValue bytes
    function getDynamicStoreValue(
        string calldata msgJson
    ) external view returns (bytes memory);

    /// @notice Get wrappable balances for a denomination
    /// @param msgJson JSON string matching QueryWrappableBalancesRequest protobuf format
    /// @return The wrappable balance amount
    function getWrappableBalances(
        string calldata msgJson
    ) external view returns (uint256);

    /// @notice Check if an address is a reserved protocol address
    /// @param msgJson JSON string matching QueryIsAddressReservedProtocolRequest protobuf format
    /// @return True if the address is reserved for protocol use
    function isAddressReservedProtocol(
        string calldata msgJson
    ) external view returns (bool);

    /// @notice Get all reserved protocol addresses
    /// @param msgJson JSON string (can be empty "{}" as no parameters needed)
    /// @return Array of reserved protocol addresses
    function getAllReservedProtocolAddresses(
        string calldata msgJson
    ) external view returns (address[] memory);

    /// @notice Get a specific vote in a voting challenge
    /// @dev Returns protobuf-encoded VoteProof
    /// @param msgJson JSON string matching QueryVoteRequest protobuf format
    /// @return Protobuf-encoded VoteProof bytes
    function getVote(
        string calldata msgJson
    ) external view returns (bytes memory);

    /// @notice Get all votes for a proposal
    /// @dev Returns protobuf-encoded array of VoteProof
    /// @param msgJson JSON string matching QueryVotesRequest protobuf format
    /// @return Protobuf-encoded VoteProof array bytes
    function getVotes(
        string calldata msgJson
    ) external view returns (bytes memory);

    /// @notice Get module parameters
    /// @dev Returns protobuf-encoded Params
    /// @param msgJson JSON string (can be empty "{}" as no parameters needed)
    /// @return Protobuf-encoded module parameters bytes
    function params(
        string calldata msgJson
    ) external view returns (bytes memory);

    /// @notice Get the balance amount for a specific (tokenId, ownershipTime) combination
    /// @dev JSON format: {"collectionId": "1", "address": "bb1...", "tokenId": "1", "ownershipTime": "1609459200000"}
    /// @param msgJson JSON string with collectionId, address, tokenId (single), ownershipTime (single)
    /// @return The exact balance amount for the specified token ID and ownership time
    function getBalanceAmount(
        string calldata msgJson
    ) external view returns (uint256);

    /// @notice Get the total supply for a specific (tokenId, ownershipTime) combination
    /// @dev JSON format: {"collectionId": "1", "tokenId": "1", "ownershipTime": "1609459200000"}
    /// @param msgJson JSON string with collectionId, tokenId (single), ownershipTime (single)
    /// @return The exact total supply for the specified token ID and ownership time
    function getTotalSupply(
        string calldata msgJson
    ) external view returns (uint256);
}

