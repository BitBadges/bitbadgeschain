// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/ITokenizationPrecompile.sol";
import "../types/TokenizationTypes.sol";
import "./TokenizationJSONHelpers.sol";
import "./TokenizationHelpers.sol";

/**
 * @title TokenizationWrappers
 * @notice Typed wrapper functions for the tokenization precompile
 * @dev Provides type-safe wrappers that accept structs instead of JSON strings.
 *      These wrappers internally convert structs to JSON using TokenizationJSONHelpers.
 *      Use these for better compile-time type checking and reduced boilerplate.
 * 
 * Usage example:
 * ```solidity
 * import "./libraries/TokenizationWrappers.sol";
 * 
 * ITokenizationPrecompile precompile = ITokenizationPrecompile(0x0000000000000000000000000000000000001001);
 * 
 * // Type-safe transfer
 * UintRange[] memory tokenIds = new UintRange[](1);
 * tokenIds[0] = TokenizationHelpers.createUintRange(1, 1);
 * UintRange[] memory ownershipTimes = new UintRange[](1);
 * ownershipTimes[0] = TokenizationHelpers.createFullOwnershipTimeRange();
 * 
 * address[] memory recipients = new address[](1);
 * recipients[0] = recipient;
 * 
 * bool success = TokenizationWrappers.transferTokens(
 *     precompile,
 *     collectionId,
 *     recipients,
 *     amount,
 *     tokenIds,
 *     ownershipTimes
 * );
 * ```
 */
library TokenizationWrappers {
    // ============ Transaction Wrappers ============

    /**
     * @notice Transfer tokens with typed parameters
     * @param precompile The precompile interface instance
     * @param collectionId The collection ID
     * @param toAddresses Array of recipient addresses
     * @param amount The amount to transfer
     * @param tokenIds Array of token ID ranges
     * @param ownershipTimes Array of ownership time ranges
     * @return success True if transfer succeeded
     */
    function transferTokens(
        ITokenizationPrecompile precompile,
        uint256 collectionId,
        address[] memory toAddresses,
        uint256 amount,
        UintRange[] memory tokenIds,
        UintRange[] memory ownershipTimes
    ) internal returns (bool success) {
        string memory tokenIdsJson = _uintRangeArrayToJson(tokenIds);
        string memory ownershipTimesJson = _uintRangeArrayToJson(ownershipTimes);
        string memory json = TokenizationJSONHelpers.transferTokensJSON(
            collectionId,
            toAddresses,
            amount,
            tokenIdsJson,
            ownershipTimesJson
        );
        return precompile.transferTokens(json);
    }

    /**
     * @notice Transfer a single token (convenience wrapper)
     * @param precompile The precompile interface instance
     * @param collectionId The collection ID
     * @param to Recipient address
     * @param amount The amount to transfer
     * @param tokenId The token ID (single value, not range)
     * @return success True if transfer succeeded
     */
    function transferSingleToken(
        ITokenizationPrecompile precompile,
        uint256 collectionId,
        address to,
        uint256 amount,
        uint256 tokenId
    ) internal returns (bool success) {
        address[] memory recipients = new address[](1);
        recipients[0] = to;
        UintRange[] memory tokenIds = new UintRange[](1);
        tokenIds[0] = TokenizationHelpers.createSingleTokenIdRange(tokenId);
        UintRange[] memory ownershipTimes = new UintRange[](1);
        ownershipTimes[0] = TokenizationHelpers.createFullOwnershipTimeRange();
        return transferTokens(precompile, collectionId, recipients, amount, tokenIds, ownershipTimes);
    }

    /**
     * @notice Transfer tokens with full ownership time range (convenience wrapper)
     * @param precompile The precompile interface instance
     * @param collectionId The collection ID
     * @param toAddresses Array of recipient addresses
     * @param amount The amount to transfer
     * @param tokenIds Array of token ID ranges
     * @return success True if transfer succeeded
     */
    function transferTokensWithFullOwnership(
        ITokenizationPrecompile precompile,
        uint256 collectionId,
        address[] memory toAddresses,
        uint256 amount,
        UintRange[] memory tokenIds
    ) internal returns (bool success) {
        UintRange[] memory ownershipTimes = new UintRange[](1);
        ownershipTimes[0] = TokenizationHelpers.createFullOwnershipTimeRange();
        return transferTokens(precompile, collectionId, toAddresses, amount, tokenIds, ownershipTimes);
    }

    /**
     * @notice Create a dynamic store with typed parameters
     * @param precompile The precompile interface instance
     * @param defaultValue Default value for addresses not explicitly set
     * @param uri URI for metadata
     * @param customData Custom data string
     * @return storeId The newly created store ID
     */
    function createDynamicStore(
        ITokenizationPrecompile precompile,
        bool defaultValue,
        string memory uri,
        string memory customData
    ) internal returns (uint256 storeId) {
        string memory json = TokenizationJSONHelpers.createDynamicStoreJSON(
            defaultValue,
            uri,
            customData
        );
        return precompile.createDynamicStore(json);
    }

    /**
     * @notice Set a dynamic store value with typed parameters
     * @param precompile The precompile interface instance
     * @param storeId The store ID
     * @param address_ The address to set the value for
     * @param value The boolean value to set
     * @return success True if value was set
     */
    function setDynamicStoreValue(
        ITokenizationPrecompile precompile,
        uint256 storeId,
        address address_,
        bool value
    ) internal returns (bool success) {
        string memory json = TokenizationJSONHelpers.setDynamicStoreValueJSON(
            storeId,
            address_,
            value
        );
        return precompile.setDynamicStoreValue(json);
    }

    /**
     * @notice Update a dynamic store with typed parameters
     * @param precompile The precompile interface instance
     * @param storeId The store ID
     * @param defaultValue Default value for addresses not explicitly set
     * @param globalEnabled Whether the store is globally enabled
     * @param uri URI for metadata
     * @param customData Custom data string
     * @return success True if update succeeded
     */
    function updateDynamicStore(
        ITokenizationPrecompile precompile,
        uint256 storeId,
        bool defaultValue,
        bool globalEnabled,
        string memory uri,
        string memory customData
    ) internal returns (bool success) {
        string memory json = TokenizationJSONHelpers.updateDynamicStoreJSON(
            storeId,
            defaultValue,
            globalEnabled,
            uri,
            customData
        );
        return precompile.updateDynamicStore(json);
    }

    /**
     * @notice Delete a dynamic store
     * @param precompile The precompile interface instance
     * @param storeId The store ID
     * @return success True if deletion succeeded
     */
    function deleteDynamicStore(
        ITokenizationPrecompile precompile,
        uint256 storeId
    ) internal returns (bool success) {
        string memory json = TokenizationJSONHelpers.deleteDynamicStoreJSON(storeId);
        return precompile.deleteDynamicStore(json);
    }

    /**
     * @notice Delete a collection
     * @param precompile The precompile interface instance
     * @param collectionId The collection ID
     * @return success True if deletion succeeded
     */
    function deleteCollection(
        ITokenizationPrecompile precompile,
        uint256 collectionId
    ) internal returns (bool success) {
        string memory json = TokenizationJSONHelpers.deleteCollectionJSON(collectionId);
        return precompile.deleteCollection(json);
    }

    /**
     * @notice Delete an incoming approval
     * @param precompile The precompile interface instance
     * @param collectionId The collection ID
     * @param approvalId The approval ID
     * @return success True if deletion succeeded
     */
    function deleteIncomingApproval(
        ITokenizationPrecompile precompile,
        uint256 collectionId,
        string memory approvalId
    ) internal returns (bool success) {
        string memory json = TokenizationJSONHelpers.deleteIncomingApprovalJSON(
            collectionId,
            approvalId
        );
        return precompile.deleteIncomingApproval(json);
    }

    /**
     * @notice Delete an outgoing approval
     * @param precompile The precompile interface instance
     * @param collectionId The collection ID
     * @param approvalId The approval ID
     * @return success True if deletion succeeded
     */
    function deleteOutgoingApproval(
        ITokenizationPrecompile precompile,
        uint256 collectionId,
        string memory approvalId
    ) internal returns (bool success) {
        string memory json = TokenizationJSONHelpers.deleteOutgoingApprovalJSON(
            collectionId,
            approvalId
        );
        return precompile.deleteOutgoingApproval(json);
    }

    /**
     * @notice Set custom data for a collection
     * @param precompile The precompile interface instance
     * @param collectionId The collection ID
     * @param customData The custom data string
     * @return resultCollectionId The collection ID (unchanged)
     */
    function setCustomData(
        ITokenizationPrecompile precompile,
        uint256 collectionId,
        string memory customData
    ) internal returns (uint256 resultCollectionId) {
        // Note: This requires constructing JSON for MsgSetCustomData
        // For now, we'll need to use JSON helpers directly or extend them
        // This is a placeholder - actual implementation would need JSON builder
        revert("setCustomData: Use JSON helpers directly for now");
    }

    /**
     * @notice Set whether a collection is archived
     * @param precompile The precompile interface instance
     * @param collectionId The collection ID
     * @param isArchived Whether the collection is archived
     * @return resultCollectionId The collection ID (unchanged)
     */
    function setIsArchived(
        ITokenizationPrecompile precompile,
        uint256 collectionId,
        bool isArchived
    ) internal returns (uint256 resultCollectionId) {
        // Note: This requires constructing JSON for MsgSetIsArchived
        // For now, we'll need to use JSON helpers directly or extend them
        revert("setIsArchived: Use JSON helpers directly for now");
    }

    // ============ Query Wrappers ============

    /**
     * @notice Get collection details by ID
     * @param precompile The precompile interface instance
     * @param collectionId The collection ID
     * @return collection Protobuf-encoded TokenCollection bytes
     * @dev Note: Returns raw protobuf bytes. Use TokenizationDecoders to decode.
     */
    function getCollection(
        ITokenizationPrecompile precompile,
        uint256 collectionId
    ) internal view returns (bytes memory collection) {
        string memory json = TokenizationJSONHelpers.getCollectionJSON(collectionId);
        return precompile.getCollection(json);
    }

    /**
     * @notice Get user balance for a collection
     * @param precompile The precompile interface instance
     * @param collectionId The collection ID
     * @param userAddress The user address
     * @return balance Protobuf-encoded UserBalanceStore bytes
     * @dev Note: Returns raw protobuf bytes. Use TokenizationDecoders to decode.
     */
    function getBalance(
        ITokenizationPrecompile precompile,
        uint256 collectionId,
        address userAddress
    ) internal view returns (bytes memory balance) {
        string memory json = TokenizationJSONHelpers.getBalanceJSON(collectionId, userAddress);
        return precompile.getBalance(json);
    }

    /**
     * @notice Get balance amount for specific token/ownership ranges
     * @param precompile The precompile interface instance
     * @param collectionId The collection ID
     * @param userAddress The user address
     * @param tokenId Single token ID to query
     * @param ownershipTime Single ownership time to query (typically ms timestamp)
     * @return amount The exact balance amount for the specified (tokenId, ownershipTime)
     */
    function getBalanceAmount(
        ITokenizationPrecompile precompile,
        uint256 collectionId,
        address userAddress,
        uint256 tokenId,
        uint256 ownershipTime
    ) internal view returns (uint256 amount) {
        string memory json = TokenizationJSONHelpers.getBalanceAmountJSON(
            collectionId,
            userAddress,
            tokenId,
            ownershipTime
        );
        return precompile.getBalanceAmount(json);
    }

    /**
     * @notice Get total supply for a specific (tokenId, ownershipTime) combination
     * @param precompile The precompile interface instance
     * @param collectionId The collection ID
     * @param tokenId Single token ID to query
     * @param ownershipTime Single ownership time to query (typically ms timestamp)
     * @return amount The exact total supply for the specified (tokenId, ownershipTime)
     */
    function getTotalSupply(
        ITokenizationPrecompile precompile,
        uint256 collectionId,
        uint256 tokenId,
        uint256 ownershipTime
    ) internal view returns (uint256 amount) {
        string memory json = TokenizationJSONHelpers.getTotalSupplyJSON(
            collectionId,
            tokenId,
            ownershipTime
        );
        return precompile.getTotalSupply(json);
    }

    /**
     * @notice Get a dynamic store by ID
     * @param precompile The precompile interface instance
     * @param storeId The store ID
     * @return store Protobuf-encoded DynamicStore bytes
     * @dev Note: Returns raw protobuf bytes. Use TokenizationDecoders to decode.
     */
    function getDynamicStore(
        ITokenizationPrecompile precompile,
        uint256 storeId
    ) internal view returns (bytes memory store) {
        string memory json = TokenizationJSONHelpers.getDynamicStoreJSON(storeId);
        return precompile.getDynamicStore(json);
    }

    /**
     * @notice Get a dynamic store value for a specific address
     * @param precompile The precompile interface instance
     * @param storeId The store ID
     * @param userAddress The user address
     * @return value Protobuf-encoded DynamicStoreValue bytes
     * @dev Note: Returns raw protobuf bytes. Use TokenizationDecoders to decode.
     */
    function getDynamicStoreValue(
        ITokenizationPrecompile precompile,
        uint256 storeId,
        address userAddress
    ) internal view returns (bytes memory value) {
        string memory json = TokenizationJSONHelpers.getDynamicStoreValueJSON(
            storeId,
            userAddress
        );
        return precompile.getDynamicStoreValue(json);
    }

    /**
     * @notice Get an address list by ID
     * @param precompile The precompile interface instance
     * @param listId The list ID
     * @return list Protobuf-encoded AddressList bytes
     * @dev Note: Returns raw protobuf bytes. Use TokenizationDecoders to decode.
     */
    function getAddressList(
        ITokenizationPrecompile precompile,
        string memory listId
    ) internal view returns (bytes memory list) {
        string memory json = TokenizationJSONHelpers.getAddressListJSON(listId);
        return precompile.getAddressList(json);
    }

    /**
     * @notice Check if an address is a reserved protocol address
     * @param precompile The precompile interface instance
     * @param addr The address to check
     * @return isReserved True if the address is reserved for protocol use
     */
    function isAddressReservedProtocol(
        ITokenizationPrecompile precompile,
        address addr
    ) internal view returns (bool isReserved) {
        string memory json = TokenizationJSONHelpers.isAddressReservedProtocolJSON(addr);
        return precompile.isAddressReservedProtocol(json);
    }

    /**
     * @notice Get all reserved protocol addresses
     * @param precompile The precompile interface instance
     * @return addresses Array of reserved protocol addresses
     */
    function getAllReservedProtocolAddresses(
        ITokenizationPrecompile precompile
    ) internal view returns (address[] memory addresses) {
        string memory json = TokenizationJSONHelpers.getAllReservedProtocolAddressesJSON();
        return precompile.getAllReservedProtocolAddresses(json);
    }

    /**
     * @notice Get module parameters
     * @param precompile The precompile interface instance
     * @return params Protobuf-encoded Params bytes
     * @dev Note: Returns raw protobuf bytes. Use TokenizationDecoders to decode.
     */
    function params(
        ITokenizationPrecompile precompile
    ) internal view returns (bytes memory) {
        string memory json = TokenizationJSONHelpers.paramsJSON();
        return precompile.params(json);
    }

    // ============ Internal Helpers ============

    /**
     * @notice Convert UintRange array to JSON string
     * @param ranges Array of UintRange structs
     * @return json JSON string representation
     */
    function _uintRangeArrayToJson(
        UintRange[] memory ranges
    ) private pure returns (string memory json) {
        if (ranges.length == 0) {
            return "[]";
        }
        uint256[] memory starts = new uint256[](ranges.length);
        uint256[] memory ends = new uint256[](ranges.length);
        for (uint256 i = 0; i < ranges.length; i++) {
            starts[i] = ranges[i].start;
            ends[i] = ranges[i].end;
        }
        return TokenizationJSONHelpers.uintRangeArrayToJson(starts, ends);
    }
}

