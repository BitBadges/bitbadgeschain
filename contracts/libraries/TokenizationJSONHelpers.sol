// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * @title TokenizationJSONHelpers
 * @notice Helper library for constructing JSON strings for the tokenization precompile
 * @dev All methods return JSON strings that match the protobuf JSON format
 * 
 * Usage example:
 * ```solidity
 * string memory json = TokenizationJSONHelpers.transferTokensJSON(
 *     collectionId,
 *     recipients,
 *     amount,
 *     tokenIds,
 *     ownershipTimes
 * );
 * bool success = precompile.transferTokens(json);
 * ```
 */
library TokenizationJSONHelpers {
    /**
     * @notice Construct JSON for transferTokens
     */
    function transferTokensJSON(
        uint256 collectionId,
        address[] memory toAddresses,
        uint256 amount,
        string memory tokenIdsJson,
        string memory ownershipTimesJson
    ) internal pure returns (string memory) {
        string memory toAddressesJson = _addressArrayToJson(toAddresses);
        return string(abi.encodePacked(
            '{"collectionId":"', _uintToString(collectionId),
            '","toAddresses":', toAddressesJson,
            ',"amount":"', _uintToString(amount),
            '","tokenIds":', tokenIdsJson,
            ',"ownershipTimes":', ownershipTimesJson,
            '}'
        ));
    }

    /**
     * @notice Construct JSON for getCollection
     */
    function getCollectionJSON(
        uint256 collectionId
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"collectionId":"', _uintToString(collectionId), '"}'
        ));
    }

    /**
     * @notice Construct JSON for getBalance
     */
    function getBalanceJSON(
        uint256 collectionId,
        address userAddress
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"collectionId":"', _uintToString(collectionId),
            '","userAddress":"', _addressToString(userAddress), '"}'
        ));
    }

    /**
     * @notice Construct JSON for getBalanceAmount
     * @dev Queries exact balance for a single (tokenId, ownershipTime) combination
     * @param collectionId The collection ID
     * @param userAddress The user's address
     * @param tokenId Single token ID to query
     * @param ownershipTime Single ownership time to query (typically ms timestamp)
     */
    function getBalanceAmountJSON(
        uint256 collectionId,
        address userAddress,
        uint256 tokenId,
        uint256 ownershipTime
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"collectionId":"', _uintToString(collectionId),
            '","address":"', _addressToString(userAddress),
            '","tokenId":"', _uintToString(tokenId),
            '","ownershipTime":"', _uintToString(ownershipTime),
            '"}'
        ));
    }

    /**
     * @notice Construct JSON for getTotalSupply
     * @dev Queries total supply for a single (tokenId, ownershipTime) combination
     * @param collectionId The collection ID
     * @param tokenId Single token ID to query
     * @param ownershipTime Single ownership time to query (typically ms timestamp)
     */
    function getTotalSupplyJSON(
        uint256 collectionId,
        uint256 tokenId,
        uint256 ownershipTime
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"collectionId":"', _uintToString(collectionId),
            '","tokenId":"', _uintToString(tokenId),
            '","ownershipTime":"', _uintToString(ownershipTime),
            '"}'
        ));
    }

    /**
     * @notice Construct JSON for createDynamicStore
     */
    function createDynamicStoreJSON(
        bool defaultValue,
        string memory uri,
        string memory customData
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"defaultValue":', defaultValue ? 'true' : 'false',
            ',"uri":"', _escapeJsonString(uri),
            '","customData":"', _escapeJsonString(customData), '"}'
        ));
    }

    /**
     * @notice Construct JSON for setDynamicStoreValue
     */
    function setDynamicStoreValueJSON(
        uint256 storeId,
        address address_,
        bool value
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"storeId":"', _uintToString(storeId),
            '","address":"', _addressToString(address_),
            '","value":', value ? 'true' : 'false', '}'
        ));
    }

    /**
     * @notice Construct JSON for getDynamicStore
     */
    function getDynamicStoreJSON(
        uint256 storeId
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"storeId":"', _uintToString(storeId), '"}'
        ));
    }

    /**
     * @notice Construct JSON for getDynamicStoreValue
     */
    function getDynamicStoreValueJSON(
        uint256 storeId,
        address userAddress
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"storeId":"', _uintToString(storeId),
            '","userAddress":"', _addressToString(userAddress), '"}'
        ));
    }

    /**
     * @notice Construct JSON for updateDynamicStore
     */
    function updateDynamicStoreJSON(
        uint256 storeId,
        bool defaultValue,
        bool globalEnabled,
        string memory uri,
        string memory customData
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"storeId":"', _uintToString(storeId),
            '","defaultValue":', defaultValue ? 'true' : 'false',
            ',"globalEnabled":', globalEnabled ? 'true' : 'false',
            ',"uri":"', _escapeJsonString(uri),
            '","customData":"', _escapeJsonString(customData), '"}'
        ));
    }

    /**
     * @notice Construct JSON for deleteDynamicStore
     */
    function deleteDynamicStoreJSON(
        uint256 storeId
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"storeId":"', _uintToString(storeId), '"}'
        ));
    }

    /**
     * @notice Construct JSON for deleteCollection
     */
    function deleteCollectionJSON(
        uint256 collectionId
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"collectionId":"', _uintToString(collectionId), '"}'
        ));
    }

    /**
     * @notice Construct JSON for isAddressReservedProtocol
     */
    function isAddressReservedProtocolJSON(
        address addr
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"address":"', _addressToString(addr), '"}'
        ));
    }

    /**
     * @notice Construct JSON for getAllReservedProtocolAddresses (empty object)
     */
    function getAllReservedProtocolAddressesJSON() internal pure returns (string memory) {
        return "{}";
    }

    /**
     * @notice Construct JSON for params (empty object)
     */
    function paramsJSON() internal pure returns (string memory) {
        return "{}";
    }

    /**
     * @notice Construct JSON for deleteIncomingApproval
     */
    function deleteIncomingApprovalJSON(
        uint256 collectionId,
        string memory approvalId
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"collectionId":"', _uintToString(collectionId),
            '","approvalId":"', _escapeJsonString(approvalId), '"}'
        ));
    }

    /**
     * @notice Construct JSON for deleteOutgoingApproval
     */
    function deleteOutgoingApprovalJSON(
        uint256 collectionId,
        string memory approvalId
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"collectionId":"', _uintToString(collectionId),
            '","approvalId":"', _escapeJsonString(approvalId), '"}'
        ));
    }

    /**
     * @notice Construct JSON for getAddressList
     */
    function getAddressListJSON(
        string memory listId
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"listId":"', _escapeJsonString(listId), '"}'
        ));
    }

    /**
     * @notice Construct JSON for createCollection
     * @dev Accepts pre-constructed JSON strings for complex nested objects
     * @param validTokenIdsJson JSON array of UintRange (use uintRangeArrayToJson)
     * @param manager Manager address as Cosmos address string
     * @param collectionMetadataJson JSON object for CollectionMetadata (use collectionMetadataToJson)
     * @param defaultBalancesJson JSON object for UserBalanceStore (optional, can be "{}")
     * @param collectionPermissionsJson JSON object for CollectionPermissions (optional, can be "{}")
     * @param standardsJson JSON array of strings (use stringArrayToJson)
     * @param customData Custom data string
     * @param isArchived Whether collection is archived
     */
    function createCollectionJSON(
        string memory validTokenIdsJson,
        string memory manager,
        string memory collectionMetadataJson,
        string memory defaultBalancesJson,
        string memory collectionPermissionsJson,
        string memory standardsJson,
        string memory customData,
        bool isArchived
    ) internal pure returns (string memory) {
        string memory result = "{";
        
        // validTokenIds
        result = string(abi.encodePacked(result, '"validTokenIds":', validTokenIdsJson));
        
        // manager
        result = string(abi.encodePacked(result, ',"manager":"', _escapeJsonString(manager), '"'));
        
        // collectionMetadata
        result = string(abi.encodePacked(result, ',"collectionMetadata":', collectionMetadataJson));
        
        // defaultBalances (optional)
        if (bytes(defaultBalancesJson).length > 0 && keccak256(bytes(defaultBalancesJson)) != keccak256(bytes("{}"))) {
            result = string(abi.encodePacked(result, ',"defaultBalances":', defaultBalancesJson));
        }
        
        // collectionPermissions (optional)
        if (bytes(collectionPermissionsJson).length > 0 && keccak256(bytes(collectionPermissionsJson)) != keccak256(bytes("{}"))) {
            result = string(abi.encodePacked(result, ',"collectionPermissions":', collectionPermissionsJson));
        }
        
        // standards
        result = string(abi.encodePacked(result, ',"standards":', standardsJson));
        
        // customData
        if (bytes(customData).length > 0) {
            result = string(abi.encodePacked(result, ',"customData":"', _escapeJsonString(customData), '"'));
        }
        
        // isArchived
        result = string(abi.encodePacked(result, ',"isArchived":', isArchived ? 'true' : 'false'));
        
        result = string(abi.encodePacked(result, "}"));
        return result;
    }

    /**
     * @notice Construct JSON for CollectionMetadata
     */
    function collectionMetadataToJson(
        string memory uri,
        string memory customData
    ) internal pure returns (string memory) {
        string memory result = "{";
        if (bytes(uri).length > 0) {
            result = string(abi.encodePacked(result, '"uri":"', _escapeJsonString(uri), '"'));
        }
        if (bytes(customData).length > 0) {
            if (bytes(uri).length > 0) {
                result = string(abi.encodePacked(result, ","));
            }
            result = string(abi.encodePacked(result, '"customData":"', _escapeJsonString(customData), '"'));
        }
        result = string(abi.encodePacked(result, "}"));
        return result;
    }

    /**
     * @notice Convert string array to JSON array
     */
    function stringArrayToJson(
        string[] memory strings
    ) internal pure returns (string memory) {
        if (strings.length == 0) {
            return "[]";
        }
        string memory result = "[";
        for (uint256 i = 0; i < strings.length; i++) {
            if (i > 0) {
                result = string(abi.encodePacked(result, ","));
            }
            result = string(abi.encodePacked(result, '"', _escapeJsonString(strings[i]), '"'));
        }
        result = string(abi.encodePacked(result, "]"));
        return result;
    }

    /**
     * @notice Construct JSON for simple UserBalanceStore (with auto-approve flags)
     * @dev For more complex UserBalanceStore, construct JSON manually
     */
    function simpleUserBalanceStoreToJson(
        bool autoApproveSelfInitiatedOutgoingTransfers,
        bool autoApproveSelfInitiatedIncomingTransfers,
        bool autoApproveAllIncomingTransfers
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"autoApproveSelfInitiatedOutgoingTransfers":', 
            autoApproveSelfInitiatedOutgoingTransfers ? 'true' : 'false',
            ',"autoApproveSelfInitiatedIncomingTransfers":',
            autoApproveSelfInitiatedIncomingTransfers ? 'true' : 'false',
            ',"autoApproveAllIncomingTransfers":',
            autoApproveAllIncomingTransfers ? 'true' : 'false',
            ',"balances":[],"outgoingApprovals":[],"incomingApprovals":[],"userPermissions":{}}'
        ));
    }

    /**
     * @notice Convert UintRange array to JSON array
     * @dev Helper for constructing tokenIds and ownershipTimes JSON
     */
    function uintRangeArrayToJson(
        uint256[] memory starts,
        uint256[] memory ends
    ) internal pure returns (string memory) {
        require(starts.length == ends.length, "Arrays must have same length");
        if (starts.length == 0) {
            return "[]";
        }
        string memory result = "[";
        for (uint256 i = 0; i < starts.length; i++) {
            if (i > 0) {
                result = string(abi.encodePacked(result, ","));
            }
            result = string(abi.encodePacked(
                result,
                '{"start":"', _uintToString(starts[i]),
                '","end":"', _uintToString(ends[i]), '"}'
            ));
        }
        result = string(abi.encodePacked(result, "]"));
        return result;
    }

    /**
     * @notice Convert single UintRange to JSON
     */
    function uintRangeToJson(
        uint256 start,
        uint256 end
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '[{"start":"', _uintToString(start),
            '","end":"', _uintToString(end), '"}]'
        ));
    }

    // ============ Internal Helpers ============

    /**
     * @notice Convert uint256 to string (public helper)
     */
    function uintToString(uint256 value) internal pure returns (string memory) {
        return _uintToString(value);
    }

    /**
     * @notice Convert address to string (public helper)
     */
    function addressToString(address addr) internal pure returns (string memory) {
        return _addressToString(addr);
    }

    /**
     * @notice Convert uint256 to string
     */
    function _uintToString(uint256 value) private pure returns (string memory) {
        if (value == 0) {
            return "0";
        }
        uint256 temp = value;
        uint256 digits;
        while (temp != 0) {
            digits++;
            temp /= 10;
        }
        bytes memory buffer = new bytes(digits);
        while (value != 0) {
            digits -= 1;
            buffer[digits] = bytes1(uint8(48 + uint256(value % 10)));
            value /= 10;
        }
        return string(buffer);
    }

    /**
     * @notice Convert address to string (hex format)
     */
    function _addressToString(address addr) private pure returns (string memory) {
        bytes memory data = abi.encodePacked(addr);
        bytes memory alphabet = "0123456789abcdef";
        bytes memory str = new bytes(2 + data.length * 2);
        str[0] = "0";
        str[1] = "x";
        for (uint256 i = 0; i < data.length; i++) {
            str[2 + i * 2] = alphabet[uint8(data[i] >> 4)];
            str[3 + i * 2] = alphabet[uint8(data[i] & 0x0f)];
        }
        return string(str);
    }

    /**
     * @notice Convert address array to JSON array
     */
    function _addressArrayToJson(address[] memory addresses) private pure returns (string memory) {
        if (addresses.length == 0) {
            return "[]";
        }
        string memory result = "[";
        for (uint256 i = 0; i < addresses.length; i++) {
            if (i > 0) {
                result = string(abi.encodePacked(result, ","));
            }
            result = string(abi.encodePacked(result, '"', _addressToString(addresses[i]), '"'));
        }
        result = string(abi.encodePacked(result, "]"));
        return result;
    }

    /**
     * @notice Escape JSON string (escape quotes, backslashes, newlines, etc.)
     * @dev Basic escaping for JSON strings - escapes quotes, backslashes, and control characters
     */
    function _escapeJsonString(string memory str) private pure returns (string memory) {
        bytes memory strBytes = bytes(str);
        bytes memory result = new bytes(strBytes.length * 2); // Worst case: all chars need escaping
        uint256 resultIndex = 0;
        
        for (uint256 i = 0; i < strBytes.length; i++) {
            bytes1 char = strBytes[i];
            if (char == 0x22) { // "
                result[resultIndex++] = 0x5C; // \
                result[resultIndex++] = 0x22; // "
            } else if (char == 0x5C) { // \
                result[resultIndex++] = 0x5C; // \
                result[resultIndex++] = 0x5C; // \
            } else if (char == 0x0A) { // \n
                result[resultIndex++] = 0x5C; // \
                result[resultIndex++] = 0x6E; // n
            } else if (char == 0x0D) { // \r
                result[resultIndex++] = 0x5C; // \
                result[resultIndex++] = 0x72; // r
            } else if (char == 0x09) { // \t
                result[resultIndex++] = 0x5C; // \
                result[resultIndex++] = 0x74; // t
            } else if (char >= 0x20) { // Printable ASCII
                result[resultIndex++] = char;
            }
            // Control characters < 0x20 are skipped (except \n, \r, \t which are handled above)
        }
        
        // Resize result to actual length
        bytes memory finalResult = new bytes(resultIndex);
        for (uint256 i = 0; i < resultIndex; i++) {
            finalResult[i] = result[i];
        }
        return string(finalResult);
    }
}

