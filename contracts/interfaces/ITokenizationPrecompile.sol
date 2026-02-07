// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../types/TokenizationTypes.sol";

/// @title ITokenizationPrecompile
/// @notice Interface for the BitBadges tokenization precompile
/// @dev Precompile address: 0x0000000000000000000000000000000000000800
///      All types are imported from TokenizationTypes for full proto compatibility
interface ITokenizationPrecompile {
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

    /// @notice Transfer tokens from the caller to specified addresses
    /// @param collectionId The collection ID
    /// @param toAddresses The recipient addresses
    /// @param amount The amount to transfer to each recipient
    /// @param tokenIds The token ID ranges to transfer
    /// @param ownershipTimes The ownership time ranges
    /// @return success True if transfer succeeded
    function transferTokens(
        uint256 collectionId,
        address[] calldata toAddresses,
        uint256 amount,
        UintRange[] calldata tokenIds,
        UintRange[] calldata ownershipTimes
    ) external returns (bool success);

    /// @notice Set an incoming approval for the caller
    /// @param collectionId The collection ID
    /// @param approval The incoming approval with full criteria support
    /// @return success True if approval was set
    function setIncomingApproval(
        uint256 collectionId,
        UserIncomingApproval calldata approval
    ) external returns (bool success);

    /// @notice Set an outgoing approval for the caller
    /// @param collectionId The collection ID
    /// @param approval The outgoing approval with full criteria support
    /// @return success True if approval was set
    function setOutgoingApproval(
        uint256 collectionId,
        UserOutgoingApproval calldata approval
    ) external returns (bool success);

    /// @notice Delete a collection (creator only)
    /// @param collectionId The collection ID to delete
    /// @return success True if deletion succeeded
    function deleteCollection(
        uint256 collectionId
    ) external returns (bool success);

    /// @notice Delete an incoming approval
    /// @param collectionId The collection ID
    /// @param approvalId The approval ID to delete
    /// @return success True if deletion succeeded
    function deleteIncomingApproval(
        uint256 collectionId,
        string calldata approvalId
    ) external returns (bool success);

    /// @notice Delete an outgoing approval
    /// @param collectionId The collection ID
    /// @param approvalId The approval ID to delete
    /// @return success True if deletion succeeded
    function deleteOutgoingApproval(
        uint256 collectionId,
        string calldata approvalId
    ) external returns (bool success);

    /// @notice Create a new dynamic store
    /// @param defaultValue The default boolean value
    /// @param uri Optional URI for metadata
    /// @param customData Optional custom data
    /// @return storeId The newly created store ID
    function createDynamicStore(
        bool defaultValue,
        string calldata uri,
        string calldata customData
    ) external returns (uint256 storeId);

    /// @notice Update an existing dynamic store
    /// @param storeId The store ID to update
    /// @param defaultValue The new default value
    /// @param globalEnabled Whether the store is globally enabled
    /// @param uri New URI for metadata
    /// @param customData New custom data
    /// @return success True if update succeeded
    function updateDynamicStore(
        uint256 storeId,
        bool defaultValue,
        bool globalEnabled,
        string calldata uri,
        string calldata customData
    ) external returns (bool success);

    /// @notice Delete a dynamic store (creator only)
    /// @param storeId The store ID to delete
    /// @return success True if deletion succeeded
    function deleteDynamicStore(
        uint256 storeId
    ) external returns (bool success);

    /// @notice Set a value in a dynamic store for an address
    /// @param storeId The store ID
    /// @param address_ The address to set the value for
    /// @param value The boolean value to set
    /// @return success True if value was set
    function setDynamicStoreValue(
        uint256 storeId,
        address address_,
        bool value
    ) external returns (bool success);

    /// @notice Set custom data for a collection
    /// @param collectionId The collection ID
    /// @param customData The new custom data
    /// @return resultCollectionId The collection ID (unchanged)
    function setCustomData(
        uint256 collectionId,
        string calldata customData
    ) external returns (uint256 resultCollectionId);

    /// @notice Set whether a collection is archived
    /// @param collectionId The collection ID
    /// @param isArchived Whether the collection is archived
    /// @return resultCollectionId The collection ID (unchanged)
    function setIsArchived(
        uint256 collectionId,
        bool isArchived
    ) external returns (uint256 resultCollectionId);

    /// @notice Set the manager of a collection
    /// @param collectionId The collection ID
    /// @param manager The new manager address (as Cosmos address string)
    /// @return resultCollectionId The collection ID (unchanged)
    function setManager(
        uint256 collectionId,
        string calldata manager
    ) external returns (uint256 resultCollectionId);

    /// @notice Set collection metadata
    /// @param collectionId The collection ID
    /// @param uri The new metadata URI
    /// @param customData Additional custom data
    /// @return resultCollectionId The collection ID (unchanged)
    function setCollectionMetadata(
        uint256 collectionId,
        string calldata uri,
        string calldata customData
    ) external returns (uint256 resultCollectionId);

    /// @notice Set the standards for a collection
    /// @param collectionId The collection ID
    /// @param standards Array of standard identifiers (e.g., "ERC721", "ERC1155")
    /// @return resultCollectionId The collection ID (unchanged)
    function setStandards(
        uint256 collectionId,
        string[] calldata standards
    ) external returns (uint256 resultCollectionId);

    /// @notice Cast a vote for a proposal in a voting challenge
    /// @param collectionId The collection ID
    /// @param approvalLevel The approval level ("collection", "incoming", "outgoing")
    /// @param approverAddress The approver address
    /// @param approvalId The approval ID containing the voting challenge
    /// @param proposalId The proposal ID to vote on
    /// @param yesWeight The weight of the yes vote
    /// @return success True if vote was cast
    function castVote(
        uint256 collectionId,
        string calldata approvalLevel,
        string calldata approverAddress,
        string calldata approvalId,
        string calldata proposalId,
        uint256 yesWeight
    ) external returns (bool success);

    /// @notice Create a new token collection
    /// @param msg_ The collection creation parameters
    /// @return newCollectionId The newly created collection ID
    function createCollection(
        MsgCreateCollection calldata msg_
    ) external returns (uint256 newCollectionId);

    /// @notice Update an existing collection
    /// @param msg_ The collection update parameters
    /// @return resultCollectionId The collection ID (unchanged)
    function updateCollection(
        MsgUpdateCollection calldata msg_
    ) external returns (uint256 resultCollectionId);

    /// @notice Update user approvals and permissions
    /// @param msg_ The user approval update parameters
    /// @return success True if update succeeded
    function updateUserApprovals(
        MsgUpdateUserApprovals calldata msg_
    ) external returns (bool success);

    /// @notice Create one or more address lists
    /// @param addressLists Array of address list definitions
    /// @return success True if lists were created
    function createAddressLists(
        AddressListInput[] calldata addressLists
    ) external returns (bool success);

    /// @notice Purge expired or specified approvals
    /// @param collectionId The collection ID
    /// @param purgeExpired Whether to purge expired approvals
    /// @param approverAddress The approver address context
    /// @param purgeCounterpartyApprovals Whether to purge counterparty approvals
    /// @param approvalsToPurge Specific approvals to purge
    /// @return numPurged Number of approvals purged
    function purgeApprovals(
        uint256 collectionId,
        bool purgeExpired,
        string calldata approverAddress,
        bool purgeCounterpartyApprovals,
        ApprovalIdentifierDetails[] calldata approvalsToPurge
    ) external returns (uint256 numPurged);

    /// @notice Set valid token IDs for a collection
    /// @param collectionId The collection ID
    /// @param validTokenIds The valid token ID ranges
    /// @param canUpdateValidTokenIds Permissions for future updates
    /// @return resultCollectionId The collection ID (unchanged)
    function setValidTokenIds(
        uint256 collectionId,
        UintRange[] calldata validTokenIds,
        TokenIdsActionPermission[] calldata canUpdateValidTokenIds
    ) external returns (uint256 resultCollectionId);

    /// @notice Set token metadata for specific token IDs
    /// @param collectionId The collection ID
    /// @param tokenMetadata Array of token metadata entries
    /// @param canUpdateTokenMetadata Permissions for future updates
    /// @return resultCollectionId The collection ID (unchanged)
    function setTokenMetadata(
        uint256 collectionId,
        TokenMetadata[] calldata tokenMetadata,
        TokenIdsActionPermission[] calldata canUpdateTokenMetadata
    ) external returns (uint256 resultCollectionId);

    /// @notice Set collection-level approvals
    /// @param collectionId The collection ID
    /// @param collectionApprovals Array of collection approvals
    /// @param canUpdateCollectionApprovals Permissions for future updates
    /// @return resultCollectionId The collection ID (unchanged)
    function setCollectionApprovals(
        uint256 collectionId,
        CollectionApproval[] calldata collectionApprovals,
        CollectionApprovalPermission[] calldata canUpdateCollectionApprovals
    ) external returns (uint256 resultCollectionId);

    /// @notice Universal update for all collection properties
    /// @dev Combines all update operations into a single call
    /// @param msg_ The universal update parameters
    /// @return resultCollectionId The collection ID (unchanged)
    function universalUpdateCollection(
        MsgUniversalUpdateCollection calldata msg_
    ) external returns (uint256 resultCollectionId);
    
    // ============ Queries ============

    /// @notice Get collection details by ID
    /// @dev Returns protobuf-encoded TokenCollection. Decode using appropriate codec.
    /// @param collectionId The collection ID to query
    /// @return collection Protobuf-encoded TokenCollection bytes
    function getCollection(
        uint256 collectionId
    ) external view returns (bytes memory collection);

    /// @notice Get user balance for a collection
    /// @dev Returns protobuf-encoded UserBalanceStore. Decode using appropriate codec.
    /// @param collectionId The collection ID
    /// @param userAddress The user address to query
    /// @return balance Protobuf-encoded UserBalanceStore bytes
    function getBalance(
        uint256 collectionId,
        address userAddress
    ) external view returns (bytes memory balance);

    /// @notice Get an address list by ID
    /// @dev Returns protobuf-encoded AddressList. Decode using appropriate codec.
    /// @param listId The list ID to query
    /// @return list Protobuf-encoded AddressList bytes
    function getAddressList(
        string calldata listId
    ) external view returns (bytes memory list);
    
    /// @notice Get approval tracker for amount/transfer tracking
    /// @dev Returns protobuf-encoded ApprovalTracker
    /// @param collectionId The collection ID
    /// @param approvalLevel The approval level ("collection", "incoming", "outgoing")
    /// @param approverAddress The approver address
    /// @param amountTrackerId The amount tracker ID
    /// @param trackerType The tracker type ("amounts" or "numTransfers")
    /// @param approvedAddress The approved address being tracked
    /// @param approvalId The approval ID
    /// @return Protobuf-encoded ApprovalTracker bytes
    function getApprovalTracker(
        uint256 collectionId,
        string calldata approvalLevel,
        address approverAddress,
        string calldata amountTrackerId,
        string calldata trackerType,
        address approvedAddress,
        string calldata approvalId
    ) external view returns (bytes memory);

    /// @notice Get challenge tracker for Merkle challenges
    /// @dev Returns protobuf-encoded challenge tracker data
    /// @param collectionId The collection ID
    /// @param approvalLevel The approval level
    /// @param approverAddress The approver address
    /// @param challengeTrackerId The challenge tracker ID
    /// @param leafIndex The leaf index in the Merkle tree
    /// @param approvalId The approval ID
    /// @return Protobuf-encoded challenge tracker bytes
    function getChallengeTracker(
        uint256 collectionId,
        string calldata approvalLevel,
        address approverAddress,
        string calldata challengeTrackerId,
        uint256 leafIndex,
        string calldata approvalId
    ) external view returns (bytes memory);

    /// @notice Get ETH signature tracker
    /// @dev Returns protobuf-encoded signature tracker data
    /// @param collectionId The collection ID
    /// @param approvalLevel The approval level
    /// @param approverAddress The approver address
    /// @param approvalId The approval ID
    /// @param challengeTrackerId The challenge tracker ID
    /// @param signature The signature being tracked
    /// @return Protobuf-encoded signature tracker bytes
    function getETHSignatureTracker(
        uint256 collectionId,
        string calldata approvalLevel,
        address approverAddress,
        string calldata approvalId,
        string calldata challengeTrackerId,
        string calldata signature
    ) external view returns (bytes memory);

    /// @notice Get a dynamic store by ID
    /// @dev Returns protobuf-encoded DynamicStore
    /// @param storeId The store ID
    /// @return Protobuf-encoded DynamicStore bytes
    function getDynamicStore(
        uint256 storeId
    ) external view returns (bytes memory);

    /// @notice Get a dynamic store value for a specific address
    /// @dev Returns protobuf-encoded DynamicStoreValue
    /// @param storeId The store ID
    /// @param userAddress The address to query
    /// @return Protobuf-encoded DynamicStoreValue bytes
    function getDynamicStoreValue(
        uint256 storeId,
        address userAddress
    ) external view returns (bytes memory);

    /// @notice Get wrappable balances for a denomination
    /// @param denom The denomination to query
    /// @param userAddress The user address
    /// @return The wrappable balance amount
    function getWrappableBalances(
        string calldata denom,
        address userAddress
    ) external view returns (uint256);

    /// @notice Check if an address is a reserved protocol address
    /// @param addr The address to check
    /// @return True if the address is reserved for protocol use
    function isAddressReservedProtocol(
        address addr
    ) external view returns (bool);

    /// @notice Get all reserved protocol addresses
    /// @return Array of reserved protocol addresses
    function getAllReservedProtocolAddresses(
    ) external view returns (address[] memory);

    /// @notice Get a specific vote in a voting challenge
    /// @dev Returns protobuf-encoded VoteProof
    /// @param collectionId The collection ID
    /// @param approvalLevel The approval level
    /// @param approverAddress The approver address
    /// @param approvalId The approval ID
    /// @param proposalId The proposal ID
    /// @param voterAddress The voter address
    /// @return Protobuf-encoded VoteProof bytes
    function getVote(
        uint256 collectionId,
        string calldata approvalLevel,
        address approverAddress,
        string calldata approvalId,
        string calldata proposalId,
        address voterAddress
    ) external view returns (bytes memory);

    /// @notice Get all votes for a proposal
    /// @dev Returns protobuf-encoded array of VoteProof
    /// @param collectionId The collection ID
    /// @param approvalLevel The approval level
    /// @param approverAddress The approver address
    /// @param approvalId The approval ID
    /// @param proposalId The proposal ID
    /// @return Protobuf-encoded VoteProof array bytes
    function getVotes(
        uint256 collectionId,
        string calldata approvalLevel,
        address approverAddress,
        string calldata approvalId,
        string calldata proposalId
    ) external view returns (bytes memory);

    /// @notice Get module parameters
    /// @dev Returns protobuf-encoded Params
    /// @return Protobuf-encoded module parameters bytes
    function params(
    ) external view returns (bytes memory);

    /// @notice Get the balance amount for specific token/ownership ranges
    /// @param collectionId The collection ID
    /// @param userAddress The user address
    /// @param tokenIds The token ID ranges to query
    /// @param ownershipTimes The ownership time ranges
    /// @return The total balance amount
    function getBalanceAmount(
        uint256 collectionId,
        address userAddress,
        UintRange[] calldata tokenIds,
        UintRange[] calldata ownershipTimes
    ) external view returns (uint256);

    /// @notice Get the total supply for specific token/ownership ranges
    /// @param collectionId The collection ID
    /// @param tokenIds The token ID ranges to query
    /// @param ownershipTimes The ownership time ranges
    /// @return The total supply amount
    function getTotalSupply(
        uint256 collectionId,
        UintRange[] calldata tokenIds,
        UintRange[] calldata ownershipTimes
    ) external view returns (uint256);
}

