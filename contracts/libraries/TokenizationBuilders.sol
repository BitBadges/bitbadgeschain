// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../types/sol";
import "./TokenizationJSONHelpers.sol";
import "./TokenizationHelpers.sol";

/**
 * @title TokenizationBuilders
 * @notice Builder pattern utilities for constructing complex tokenization operations
 * @dev Provides fluent builder APIs for complex operations like collection creation.
 *      These builders make it easier to construct complex JSON structures step by step.
 * 
 * @example
 * ```solidity
 * import "./libraries/TokenizationBuilders.sol";
 * 
 * // Build a collection creation request
 * CollectionBuilder memory builder = TokenizationBuilders.newCollection();
 * builder = builder.withValidTokenIds(tokenIds);
 * builder = builder.withManager(managerAddress);
 * builder = builder.withMetadata(metadata);
 * string memory json = builder.build();
 * uint256 collectionId = precompile.createCollection(json);
 * ```
 */
library TokenizationBuilders {
    // ============ Collection Builder ============

    /**
     * @notice Builder struct for creating collections
     * @dev All fields are optional except validTokenIds
     */
    struct CollectionBuilder {
        UintRange[] validTokenIds;
        string manager;
        CollectionMetadata metadata;
        UserBalanceStore defaultBalances;
        CollectionPermissions collectionPermissions;
        TokenMetadata[] tokenMetadata;
        string customData;
        CollectionApproval[] collectionApprovals;
        string[] standards;
        bool isArchived;
        bool hasValidTokenIds;
        bool hasManager;
        bool hasMetadata;
        bool hasDefaultBalances;
        bool hasCollectionPermissions;
        bool hasTokenMetadata;
        bool hasCustomData;
        bool hasCollectionApprovals;
        bool hasStandards;
        bool hasIsArchived;
    }

    /**
     * @notice Create a new collection builder
     * @return builder A new CollectionBuilder with default values
     */
    function newCollection() internal pure returns (CollectionBuilder memory builder) {
        // All fields start as unset
        return builder;
    }

    /**
     * @notice Set valid token IDs for the collection
     * @param builder The builder instance
     * @param tokenIds Array of token ID ranges
     * @return The builder with validTokenIds set
     */
    function withValidTokenIds(
        CollectionBuilder memory builder,
        UintRange[] memory tokenIds
    ) internal pure returns (CollectionBuilder memory) {
        builder.validTokenIds = tokenIds;
        builder.hasValidTokenIds = true;
        return builder;
    }

    /**
     * @notice Set valid token IDs from start and end arrays
     * @param builder The builder instance
     * @param starts Array of start values
     * @param ends Array of end values
     * @return The builder with validTokenIds set
     */
    function withValidTokenIdsFromArrays(
        CollectionBuilder memory builder,
        uint256[] memory starts,
        uint256[] memory ends
    ) internal pure returns (CollectionBuilder memory) {
        builder.validTokenIds = TokenizationHelpers.createUintRangeArray(starts, ends);
        builder.hasValidTokenIds = true;
        return builder;
    }

    /**
     * @notice Set a single token ID range
     * @param builder The builder instance
     * @param start The start token ID
     * @param end The end token ID
     * @return The builder with validTokenIds set
     */
    function withValidTokenIdRange(
        CollectionBuilder memory builder,
        uint256 start,
        uint256 end
    ) internal pure returns (CollectionBuilder memory) {
        UintRange[] memory tokenIds = new UintRange[](1);
        tokenIds[0] = TokenizationHelpers.createUintRange(start, end);
        builder.validTokenIds = tokenIds;
        builder.hasValidTokenIds = true;
        return builder;
    }

    /**
     * @notice Set the manager address
     * @param builder The builder instance
     * @param manager The manager address (as Cosmos address string)
     * @return The builder with manager set
     */
    function withManager(
        CollectionBuilder memory builder,
        string memory manager
    ) internal pure returns (CollectionBuilder memory) {
        builder.manager = manager;
        builder.hasManager = true;
        return builder;
    }

    /**
     * @notice Set collection metadata
     * @param builder The builder instance
     * @param metadata The collection metadata
     * @return The builder with metadata set
     */
    function withMetadata(
        CollectionBuilder memory builder,
        CollectionMetadata memory metadata
    ) internal pure returns (CollectionBuilder memory) {
        builder.metadata = metadata;
        builder.hasMetadata = true;
        return builder;
    }

    /**
     * @notice Set collection metadata from URI and custom data
     * @param builder The builder instance
     * @param uri The metadata URI
     * @param customData The custom data
     * @return The builder with metadata set
     */
    function withMetadataFromStrings(
        CollectionBuilder memory builder,
        string memory uri,
        string memory customData
    ) internal pure returns (CollectionBuilder memory) {
        builder.metadata = TokenizationHelpers.createCollectionMetadata(uri, customData);
        builder.hasMetadata = true;
        return builder;
    }

    /**
     * @notice Set default balances
     * @param builder The builder instance
     * @param defaultBalances The default user balance store
     * @return The builder with defaultBalances set
     */
    function withDefaultBalances(
        CollectionBuilder memory builder,
        UserBalanceStore memory defaultBalances
    ) internal pure returns (CollectionBuilder memory) {
        builder.defaultBalances = defaultBalances;
        builder.hasDefaultBalances = true;
        return builder;
    }

    /**
     * @notice Set default balances with auto-approve flags
     * @param builder The builder instance
     * @param autoApproveSelfInitiatedOutgoing Transfers Whether to auto-approve self-initiated outgoing transfers
     * @param autoApproveSelfInitiatedIncoming Transfers Whether to auto-approve self-initiated incoming transfers
     * @param autoApproveAllIncomingTransfers Whether to auto-approve all incoming transfers
     * @return The builder with defaultBalances set
     */
    function withDefaultBalancesFromFlags(
        CollectionBuilder memory builder,
        bool autoApproveSelfInitiatedOutgoingTransfers,
        bool autoApproveSelfInitiatedIncomingTransfers,
        bool autoApproveAllIncomingTransfers
    ) internal pure returns (CollectionBuilder memory) {
        builder.defaultBalances = TokenizationHelpers.createEmptyUserBalanceStore();
        builder.defaultBalances.autoApproveSelfInitiatedOutgoingTransfers = autoApproveSelfInitiatedOutgoingTransfers;
        builder.defaultBalances.autoApproveSelfInitiatedIncomingTransfers = autoApproveSelfInitiatedIncomingTransfers;
        builder.defaultBalances.autoApproveAllIncomingTransfers = autoApproveAllIncomingTransfers;
        builder.hasDefaultBalances = true;
        return builder;
    }

    /**
     * @notice Set collection permissions
     * @param builder The builder instance
     * @param permissions The collection permissions
     * @return The builder with collectionPermissions set
     */
    function withCollectionPermissions(
        CollectionBuilder memory builder,
        CollectionPermissions memory permissions
    ) internal pure returns (CollectionBuilder memory) {
        builder.collectionPermissions = permissions;
        builder.hasCollectionPermissions = true;
        return builder;
    }

    /**
     * @notice Set token metadata
     * @param builder The builder instance
     * @param tokenMetadata Array of token metadata
     * @return The builder with tokenMetadata set
     */
    function withTokenMetadata(
        CollectionBuilder memory builder,
        TokenMetadata[] memory tokenMetadata
    ) internal pure returns (CollectionBuilder memory) {
        builder.tokenMetadata = tokenMetadata;
        builder.hasTokenMetadata = true;
        return builder;
    }

    /**
     * @notice Set custom data
     * @param builder The builder instance
     * @param customData The custom data string
     * @return The builder with customData set
     */
    function withCustomData(
        CollectionBuilder memory builder,
        string memory customData
    ) internal pure returns (CollectionBuilder memory) {
        builder.customData = customData;
        builder.hasCustomData = true;
        return builder;
    }

    /**
     * @notice Set collection approvals
     * @param builder The builder instance
     * @param approvals Array of collection approvals
     * @return The builder with collectionApprovals set
     */
    function withCollectionApprovals(
        CollectionBuilder memory builder,
        CollectionApproval[] memory approvals
    ) internal pure returns (CollectionBuilder memory) {
        builder.collectionApprovals = approvals;
        builder.hasCollectionApprovals = true;
        return builder;
    }

    /**
     * @notice Set standards
     * @param builder The builder instance
     * @param standards Array of standard strings
     * @return The builder with standards set
     */
    function withStandards(
        CollectionBuilder memory builder,
        string[] memory standards
    ) internal pure returns (CollectionBuilder memory) {
        builder.standards = standards;
        builder.hasStandards = true;
        return builder;
    }

    /**
     * @notice Set whether the collection is archived
     * @param builder The builder instance
     * @param isArchived Whether the collection is archived
     * @return The builder with isArchived set
     */
    function withIsArchived(
        CollectionBuilder memory builder,
        bool isArchived
    ) internal pure returns (CollectionBuilder memory) {
        builder.isArchived = isArchived;
        builder.hasIsArchived = true;
        return builder;
    }

    /**
     * @notice Build the JSON string for collection creation
     * @param builder The builder instance
     * @return json The JSON string ready for createCollection
     * @dev Requires validTokenIds to be set
     */
    function build(CollectionBuilder memory builder) internal pure returns (string memory json) {
        require(builder.hasValidTokenIds, "TokenizationBuilders: validTokenIds is required");

        // Convert validTokenIds to JSON
        uint256[] memory starts = new uint256[](builder.validTokenIds.length);
        uint256[] memory ends = new uint256[](builder.validTokenIds.length);
        for (uint256 i = 0; i < builder.validTokenIds.length; i++) {
            starts[i] = builder.validTokenIds[i].start;
            ends[i] = builder.validTokenIds[i].end;
        }
        string memory validTokenIdsJson = TokenizationJSONHelpers.uintRangeArrayToJson(starts, ends);

        // Convert metadata to JSON
        string memory collectionMetadataJson = TokenizationJSONHelpers.collectionMetadataToJson(
            builder.hasMetadata ? builder.metadata.uri : "",
            builder.hasMetadata ? builder.metadata.customData : ""
        );

        // Convert defaultBalances to JSON (simplified version)
        string memory defaultBalancesJson = "{}";
        if (builder.hasDefaultBalances) {
            defaultBalancesJson = TokenizationJSONHelpers.simpleUserBalanceStoreToJson(
                builder.defaultBalances.autoApproveSelfInitiatedOutgoingTransfers,
                builder.defaultBalances.autoApproveSelfInitiatedIncomingTransfers,
                builder.defaultBalances.autoApproveAllIncomingTransfers
            );
        }

        // Convert collectionPermissions to JSON (empty for now - complex structure)
        string memory collectionPermissionsJson = "{}";

        // Convert standards to JSON
        string memory standardsJson = builder.hasStandards
            ? TokenizationJSONHelpers.stringArrayToJson(builder.standards)
            : "[]";

        // Build the final JSON
        return TokenizationJSONHelpers.createCollectionJSON(
            validTokenIdsJson,
            builder.hasManager ? builder.manager : "",
            collectionMetadataJson,
            defaultBalancesJson,
            collectionPermissionsJson,
            standardsJson,
            builder.hasCustomData ? builder.customData : "",
            builder.hasIsArchived ? builder.isArchived : false
        );
    }

    // ============ Transfer Builder ============

    /**
     * @notice Builder struct for creating transfers
     */
    struct TransferBuilder {
        address[] toAddresses;
        uint256 amount;
        UintRange[] tokenIds;
        UintRange[] ownershipTimes;
    }

    /**
     * @notice Create a new transfer builder
     * @return builder A new TransferBuilder
     */
    function newTransfer() internal pure returns (TransferBuilder memory builder) {
        return builder;
    }

    /**
     * @notice Set recipient addresses
     * @param builder The builder instance
     * @param toAddresses Array of recipient addresses
     * @return The builder with toAddresses set
     */
    function withRecipients(
        TransferBuilder memory builder,
        address[] memory toAddresses
    ) internal pure returns (TransferBuilder memory) {
        builder.toAddresses = toAddresses;
        return builder;
    }

    /**
     * @notice Set a single recipient
     * @param builder The builder instance
     * @param to The recipient address
     * @return The builder with toAddresses set
     */
    function withRecipient(
        TransferBuilder memory builder,
        address to
    ) internal pure returns (TransferBuilder memory) {
        address[] memory recipients = new address[](1);
        recipients[0] = to;
        builder.toAddresses = recipients;
        return builder;
    }

    /**
     * @notice Set the transfer amount
     * @param builder The builder instance
     * @param amount The amount to transfer
     * @return The builder with amount set
     */
    function withAmount(
        TransferBuilder memory builder,
        uint256 amount
    ) internal pure returns (TransferBuilder memory) {
        builder.amount = amount;
        return builder;
    }

    /**
     * @notice Set token ID ranges
     * @param builder The builder instance
     * @param tokenIds Array of token ID ranges
     * @return The builder with tokenIds set
     */
    function withTokenIds(
        TransferBuilder memory builder,
        UintRange[] memory tokenIds
    ) internal pure returns (TransferBuilder memory) {
        builder.tokenIds = tokenIds;
        return builder;
    }

    /**
     * @notice Set a single token ID
     * @param builder The builder instance
     * @param tokenId The token ID
     * @return The builder with tokenIds set
     */
    function withTokenId(
        TransferBuilder memory builder,
        uint256 tokenId
    ) internal pure returns (TransferBuilder memory) {
        UintRange[] memory tokenIds = new UintRange[](1);
        tokenIds[0] = TokenizationHelpers.createSingleTokenIdRange(tokenId);
        builder.tokenIds = tokenIds;
        return builder;
    }

    /**
     * @notice Set ownership time ranges
     * @param builder The builder instance
     * @param ownershipTimes Array of ownership time ranges
     * @return The builder with ownershipTimes set
     */
    function withOwnershipTimes(
        TransferBuilder memory builder,
        UintRange[] memory ownershipTimes
    ) internal pure returns (TransferBuilder memory) {
        builder.ownershipTimes = ownershipTimes;
        return builder;
    }

    /**
     * @notice Set full ownership time range (1 to max)
     * @param builder The builder instance
     * @return The builder with ownershipTimes set to full range
     */
    function withFullOwnershipTime(
        TransferBuilder memory builder
    ) internal pure returns (TransferBuilder memory) {
        UintRange[] memory ownershipTimes = new UintRange[](1);
        ownershipTimes[0] = TokenizationHelpers.createFullOwnershipTimeRange();
        builder.ownershipTimes = ownershipTimes;
        return builder;
    }

    /**
     * @notice Build the JSON string for transfer
     * @param builder The builder instance
     * @param collectionId The collection ID
     * @return json The JSON string ready for transferTokens
     */
    function buildTransfer(
        TransferBuilder memory builder,
        uint256 collectionId
    ) internal pure returns (string memory json) {
        require(builder.toAddresses.length > 0, "TokenizationBuilders: recipients required");
        require(builder.amount > 0, "TokenizationBuilders: amount must be > 0");
        require(builder.tokenIds.length > 0, "TokenizationBuilders: tokenIds required");
        require(builder.ownershipTimes.length > 0, "TokenizationBuilders: ownershipTimes required");

        uint256[] memory tokenStarts = new uint256[](builder.tokenIds.length);
        uint256[] memory tokenEnds = new uint256[](builder.tokenIds.length);
        for (uint256 i = 0; i < builder.tokenIds.length; i++) {
            tokenStarts[i] = builder.tokenIds[i].start;
            tokenEnds[i] = builder.tokenIds[i].end;
        }
        string memory tokenIdsJson = TokenizationJSONHelpers.uintRangeArrayToJson(tokenStarts, tokenEnds);

        uint256[] memory ownershipStarts = new uint256[](builder.ownershipTimes.length);
        uint256[] memory ownershipEnds = new uint256[](builder.ownershipTimes.length);
        for (uint256 i = 0; i < builder.ownershipTimes.length; i++) {
            ownershipStarts[i] = builder.ownershipTimes[i].start;
            ownershipEnds[i] = builder.ownershipTimes[i].end;
        }
        string memory ownershipTimesJson = TokenizationJSONHelpers.uintRangeArrayToJson(ownershipStarts, ownershipEnds);

        return TokenizationJSONHelpers.transferTokensJSON(
            collectionId,
            builder.toAddresses,
            builder.amount,
            tokenIdsJson,
            ownershipTimesJson
        );
    }
}

















