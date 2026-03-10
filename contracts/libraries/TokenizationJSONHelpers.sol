// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * @title TokenizationJSONHelpers
 * @notice Helper library for constructing JSON strings for the tokenization precompile
 * @dev All methods return JSON strings that match the protobuf JSON format
 *
 * IMPORTANT: BitBadges uses uint64 internally. Use the constants below for time/ID ranges.
 * Using type(uint256).max will cause "range overflow" errors!
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
    // ============ Critical Constants ============
    // IMPORTANT: BitBadges uses uint64 internally for timestamps and IDs.
    // Using type(uint256).max will cause "range overflow" errors!

    /// @notice Use for ownership times that never expire
    uint64 public constant FOREVER = type(uint64).max;  // 18446744073709551615

    /// @notice Maximum valid time value for BitBadges
    uint64 public constant MAX_TIME = type(uint64).max;

    /// @notice Maximum valid token ID value for BitBadges
    uint64 public constant MAX_ID = type(uint64).max;

    /// @notice Minimum valid time/ID value (typically 1, not 0)
    uint64 public constant MIN_ID = 1;

    /// @notice String version of FOREVER for direct JSON use
    string public constant FOREVER_STR = "18446744073709551615";
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
     * @dev Uses "address" field name to match proto definition
     */
    function getBalanceJSON(
        uint256 collectionId,
        address userAddress
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"collectionId":"', _uintToString(collectionId),
            '","address":"', _addressToString(userAddress), '"}'
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

    // ============ Dynamic Store (Boolean Address Flags) ============
    // NOTE: Despite the name "Dynamic Store", this is a BOOLEAN store per address.
    // It does NOT support arbitrary key-value pairs.
    // - Each store maps addresses to true/false
    // - defaultValue: returned for addresses not explicitly set
    // - globalEnabled: kill switch (false = all lookups return false)

    /**
     * @notice Construct JSON for createDynamicStore (boolean address flag store)
     * @param defaultValue The default boolean value for addresses not explicitly set
     * @param uri Metadata URI for the store
     * @param customData Custom data string
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
     * @dev Sets a boolean flag for the given address in the store
     *      The address is auto-converted from EVM to bech32 format by the precompile
     * @param storeId The dynamic store ID
     * @param address_ The address to set the value for (0x... format)
     * @param value The boolean value to set (true/false)
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
     * @dev Uses "address" field name to match proto definition
     */
    function getDynamicStoreValueJSON(
        uint256 storeId,
        address userAddress
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"storeId":"', _uintToString(storeId),
            '","address":"', _addressToString(userAddress), '"}'
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
     * @dev For more complex UserBalanceStore, use userBalanceStoreToJson()
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
     * @notice Construct JSON for full UserBalanceStore with all 7 fields
     * @dev Proto fields: balances[], outgoingApprovals[], incomingApprovals[],
     *      autoApproveSelfInitiatedOutgoingTransfers, autoApproveSelfInitiatedIncomingTransfers,
     *      autoApproveAllIncomingTransfers, userPermissions
     * @param balancesJson Pre-encoded Balance[] array JSON (use balanceArrayToJson)
     * @param outgoingApprovalsJson Pre-encoded UserOutgoingApproval[] array JSON
     * @param incomingApprovalsJson Pre-encoded UserIncomingApproval[] array JSON
     * @param autoApproveSelfInitiatedOutgoing Auto-approve self-initiated outgoing transfers
     * @param autoApproveSelfInitiatedIncoming Auto-approve self-initiated incoming transfers
     * @param autoApproveAllIncoming Auto-approve all incoming transfers
     * @param userPermissionsJson Pre-encoded UserPermissions JSON
     */
    function userBalanceStoreToJson(
        string memory balancesJson,
        string memory outgoingApprovalsJson,
        string memory incomingApprovalsJson,
        bool autoApproveSelfInitiatedOutgoing,
        bool autoApproveSelfInitiatedIncoming,
        bool autoApproveAllIncoming,
        string memory userPermissionsJson
    ) internal pure returns (string memory) {
        string memory result = "{";

        // balances
        if (bytes(balancesJson).length > 0 && keccak256(bytes(balancesJson)) != keccak256(bytes("[]"))) {
            result = string(abi.encodePacked(result, '"balances":', balancesJson));
        } else {
            result = string(abi.encodePacked(result, '"balances":[]'));
        }

        // outgoingApprovals
        if (bytes(outgoingApprovalsJson).length > 0 && keccak256(bytes(outgoingApprovalsJson)) != keccak256(bytes("[]"))) {
            result = string(abi.encodePacked(result, ',"outgoingApprovals":', outgoingApprovalsJson));
        } else {
            result = string(abi.encodePacked(result, ',"outgoingApprovals":[]'));
        }

        // incomingApprovals
        if (bytes(incomingApprovalsJson).length > 0 && keccak256(bytes(incomingApprovalsJson)) != keccak256(bytes("[]"))) {
            result = string(abi.encodePacked(result, ',"incomingApprovals":', incomingApprovalsJson));
        } else {
            result = string(abi.encodePacked(result, ',"incomingApprovals":[]'));
        }

        // Boolean flags
        result = string(abi.encodePacked(
            result,
            ',"autoApproveSelfInitiatedOutgoingTransfers":', autoApproveSelfInitiatedOutgoing ? 'true' : 'false',
            ',"autoApproveSelfInitiatedIncomingTransfers":', autoApproveSelfInitiatedIncoming ? 'true' : 'false',
            ',"autoApproveAllIncomingTransfers":', autoApproveAllIncoming ? 'true' : 'false'
        ));

        // userPermissions
        if (bytes(userPermissionsJson).length > 0 && keccak256(bytes(userPermissionsJson)) != keccak256(bytes("{}"))) {
            result = string(abi.encodePacked(result, ',"userPermissions":', userPermissionsJson));
        } else {
            result = string(abi.encodePacked(result, ',"userPermissions":{}'));
        }

        return string(abi.encodePacked(result, '}'));
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

    // ============ Invariants JSON Helpers (v25) ============

    /**
     * @notice Convert EVMQueryChallenge to JSON
     * @dev Used for building invariants with EVM query challenges
     * @param contractAddress The EVM contract address (hex string with 0x prefix)
     * @param callData The calldata for the static call (hex string)
     * @param expectedResult The expected result (hex string)
     * @param comparisonOperator The comparison operator ("eq", "ne", "gt", "gte", "lt", "lte")
     * @param gasLimit Gas limit for the call
     * @param uri Optional metadata URI
     * @param customData Optional custom data
     */
    function evmQueryChallengeToJson(
        string memory contractAddress,
        string memory callData,
        string memory expectedResult,
        string memory comparisonOperator,
        uint256 gasLimit,
        string memory uri,
        string memory customData
    ) internal pure returns (string memory) {
        string memory result = string(abi.encodePacked(
            '{"contractAddress":"', _escapeJsonString(contractAddress),
            '","calldata":"', _escapeJsonString(callData),
            '","expectedResult":"', _escapeJsonString(expectedResult),
            '","comparisonOperator":"', _escapeJsonString(comparisonOperator),
            '","gasLimit":"', _uintToString(gasLimit), '"'
        ));

        if (bytes(uri).length > 0) {
            result = string(abi.encodePacked(result, ',"uri":"', _escapeJsonString(uri), '"'));
        }
        if (bytes(customData).length > 0) {
            result = string(abi.encodePacked(result, ',"customData":"', _escapeJsonString(customData), '"'));
        }

        return string(abi.encodePacked(result, '}'));
    }

    /**
     * @notice Convert EVMQueryChallenge array to JSON array (simple version without uri/customData)
     * @dev Accepts parallel arrays matching EVMQueryChallenge fields. For full version use evmQueryChallengeArrayFullToJson.
     */
    function evmQueryChallengeArrayToJson(
        string[] memory contractAddresses,
        string[] memory callDatas,
        string[] memory expectedResults,
        string[] memory comparisonOperators,
        uint256[] memory gasLimits
    ) internal pure returns (string memory) {
        require(
            contractAddresses.length == callDatas.length &&
            callDatas.length == expectedResults.length &&
            expectedResults.length == comparisonOperators.length &&
            comparisonOperators.length == gasLimits.length,
            "Array lengths must match"
        );

        if (contractAddresses.length == 0) {
            return "[]";
        }

        string memory result = "[";
        for (uint256 i = 0; i < contractAddresses.length; i++) {
            if (i > 0) {
                result = string(abi.encodePacked(result, ","));
            }
            result = string(abi.encodePacked(
                result,
                evmQueryChallengeToJson(
                    contractAddresses[i],
                    callDatas[i],
                    expectedResults[i],
                    comparisonOperators[i],
                    gasLimits[i],
                    "",
                    ""
                )
            ));
        }
        return string(abi.encodePacked(result, "]"));
    }

    /**
     * @notice Convert EVMQueryChallenge array to JSON array with full field support
     * @dev Accepts parallel arrays matching all EVMQueryChallenge fields including uri and customData
     * @param contractAddresses Array of contract addresses
     * @param callDatas Array of calldatas
     * @param expectedResults Array of expected results
     * @param comparisonOperators Array of comparison operators
     * @param gasLimits Array of gas limits
     * @param uris Array of metadata URIs (can be empty strings)
     * @param customDatas Array of custom data strings (can be empty strings)
     */
    function evmQueryChallengeArrayFullToJson(
        string[] memory contractAddresses,
        string[] memory callDatas,
        string[] memory expectedResults,
        string[] memory comparisonOperators,
        uint256[] memory gasLimits,
        string[] memory uris,
        string[] memory customDatas
    ) internal pure returns (string memory) {
        require(
            contractAddresses.length == callDatas.length &&
            callDatas.length == expectedResults.length &&
            expectedResults.length == comparisonOperators.length &&
            comparisonOperators.length == gasLimits.length &&
            gasLimits.length == uris.length &&
            uris.length == customDatas.length,
            "Array lengths must match"
        );

        if (contractAddresses.length == 0) {
            return "[]";
        }

        string memory result = "[";
        for (uint256 i = 0; i < contractAddresses.length; i++) {
            if (i > 0) {
                result = string(abi.encodePacked(result, ","));
            }
            result = string(abi.encodePacked(
                result,
                evmQueryChallengeToJson(
                    contractAddresses[i],
                    callDatas[i],
                    expectedResults[i],
                    comparisonOperators[i],
                    gasLimits[i],
                    uris[i],
                    customDatas[i]
                )
            ));
        }
        return string(abi.encodePacked(result, "]"));
    }

    /**
     * @notice Convert CollectionInvariants to JSON (simple version)
     * @dev Used for collection creation/update with invariants. For full version use collectionInvariantsFullToJson.
     * @param noCustomOwnershipTimes Disallow custom ownership times
     * @param maxSupplyPerId Maximum supply per token ID (0 for unlimited)
     * @param noForcefulPostMintTransfers Prevent forceful transfers after minting
     * @param disablePoolCreation Disable liquidity pool creation
     * @param evmQueryChallengesJson Pre-constructed JSON array of EVM query challenges
     */
    function collectionInvariantsToJson(
        bool noCustomOwnershipTimes,
        uint256 maxSupplyPerId,
        bool noForcefulPostMintTransfers,
        bool disablePoolCreation,
        string memory evmQueryChallengesJson
    ) internal pure returns (string memory) {
        return collectionInvariantsFullToJson(
            noCustomOwnershipTimes,
            maxSupplyPerId,
            noForcefulPostMintTransfers,
            disablePoolCreation,
            evmQueryChallengesJson,
            ""
        );
    }

    /**
     * @notice Convert CollectionInvariants to JSON with all fields
     * @dev Includes cosmosCoinBackedPath field
     * @param noCustomOwnershipTimes Disallow custom ownership times
     * @param maxSupplyPerId Maximum supply per token ID (0 for unlimited)
     * @param noForcefulPostMintTransfers Prevent forceful transfers after minting
     * @param disablePoolCreation Disable liquidity pool creation
     * @param evmQueryChallengesJson Pre-constructed JSON array of EVM query challenges
     * @param cosmosCoinBackedPathJson Pre-constructed CosmosCoinBackedPath JSON
     */
    function collectionInvariantsFullToJson(
        bool noCustomOwnershipTimes,
        uint256 maxSupplyPerId,
        bool noForcefulPostMintTransfers,
        bool disablePoolCreation,
        string memory evmQueryChallengesJson,
        string memory cosmosCoinBackedPathJson
    ) internal pure returns (string memory) {
        string memory result = "{";

        result = string(abi.encodePacked(
            result,
            '"noCustomOwnershipTimes":', noCustomOwnershipTimes ? 'true' : 'false',
            ',"noForcefulPostMintTransfers":', noForcefulPostMintTransfers ? 'true' : 'false',
            ',"disablePoolCreation":', disablePoolCreation ? 'true' : 'false'
        ));

        if (maxSupplyPerId > 0) {
            result = string(abi.encodePacked(result, ',"maxSupplyPerId":"', _uintToString(maxSupplyPerId), '"'));
        }

        if (bytes(evmQueryChallengesJson).length > 0 && keccak256(bytes(evmQueryChallengesJson)) != keccak256(bytes("[]"))) {
            result = string(abi.encodePacked(result, ',"evmQueryChallenges":', evmQueryChallengesJson));
        }

        if (bytes(cosmosCoinBackedPathJson).length > 0 && keccak256(bytes(cosmosCoinBackedPathJson)) != keccak256(bytes("{}"))) {
            result = string(abi.encodePacked(result, ',"cosmosCoinBackedPath":', cosmosCoinBackedPathJson));
        }

        return string(abi.encodePacked(result, '}'));
    }

    /**
     * @notice Construct JSON for createCollection with invariants (v25 extended version)
     * @dev Extends createCollectionJSON with invariants support
     */
    function createCollectionWithInvariantsJSON(
        string memory validTokenIdsJson,
        string memory manager,
        string memory collectionMetadataJson,
        string memory defaultBalancesJson,
        string memory collectionPermissionsJson,
        string memory standardsJson,
        string memory customData,
        bool isArchived,
        string memory invariantsJson
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

        // invariants (v25)
        if (bytes(invariantsJson).length > 0 && keccak256(bytes(invariantsJson)) != keccak256(bytes("{}"))) {
            result = string(abi.encodePacked(result, ',"invariants":', invariantsJson));
        }

        result = string(abi.encodePacked(result, "}"));
        return result;
    }

    // ============ Nested Type JSON Helpers ============

    /**
     * @notice Convert Balance to JSON
     * @dev Proto fields: amount, tokenIds[], ownershipTimes[]
     * @param amount The balance amount
     * @param tokenIdsJson Pre-encoded UintRange[] JSON for token IDs
     * @param ownershipTimesJson Pre-encoded UintRange[] JSON for ownership times
     */
    function balanceToJson(
        uint256 amount,
        string memory tokenIdsJson,
        string memory ownershipTimesJson
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"amount":"', _uintToString(amount),
            '","badgeIds":', tokenIdsJson,
            ',"ownershipTimes":', ownershipTimesJson, '}'
        ));
    }

    /**
     * @notice Convert Balance array to JSON array
     * @param amounts Array of amounts
     * @param tokenIdsJsons Array of pre-encoded UintRange[] JSONs
     * @param ownershipTimesJsons Array of pre-encoded UintRange[] JSONs
     */
    function balanceArrayToJson(
        uint256[] memory amounts,
        string[] memory tokenIdsJsons,
        string[] memory ownershipTimesJsons
    ) internal pure returns (string memory) {
        require(
            amounts.length == tokenIdsJsons.length &&
            tokenIdsJsons.length == ownershipTimesJsons.length,
            "Array lengths must match"
        );

        if (amounts.length == 0) {
            return "[]";
        }

        string memory result = "[";
        for (uint256 i = 0; i < amounts.length; i++) {
            if (i > 0) {
                result = string(abi.encodePacked(result, ","));
            }
            result = string(abi.encodePacked(
                result,
                balanceToJson(amounts[i], tokenIdsJsons[i], ownershipTimesJsons[i])
            ));
        }
        return string(abi.encodePacked(result, "]"));
    }

    /**
     * @notice Convert PathMetadata to JSON
     * @dev Proto fields: uri, customData
     * @param uri Metadata URI
     * @param customData Custom data string
     */
    function pathMetadataToJson(
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
        return string(abi.encodePacked(result, "}"));
    }

    /**
     * @notice Convert TokenMetadata to JSON
     * @dev Proto fields: uri, customData (same as PathMetadata)
     * @param uri Metadata URI
     * @param customData Custom data string
     */
    function tokenMetadataToJson(
        string memory uri,
        string memory customData
    ) internal pure returns (string memory) {
        return pathMetadataToJson(uri, customData);
    }

    /**
     * @notice Convert DenomUnit to JSON
     * @dev Proto fields: decimals, symbol, isDefaultDisplay, metadata
     * @param decimals Number of decimals
     * @param symbol Token symbol
     * @param isDefaultDisplay Whether this is the default display unit
     * @param metadataJson Pre-encoded PathMetadata JSON
     */
    function denomUnitToJson(
        uint256 decimals,
        string memory symbol,
        bool isDefaultDisplay,
        string memory metadataJson
    ) internal pure returns (string memory) {
        string memory result = string(abi.encodePacked(
            '{"decimals":"', _uintToString(decimals),
            '","symbol":"', _escapeJsonString(symbol),
            '","isDefaultDisplay":', isDefaultDisplay ? 'true' : 'false'
        ));

        if (bytes(metadataJson).length > 0 && keccak256(bytes(metadataJson)) != keccak256(bytes("{}"))) {
            result = string(abi.encodePacked(result, ',"metadata":', metadataJson));
        }

        return string(abi.encodePacked(result, '}'));
    }

    /**
     * @notice Convert DenomUnit array to JSON array
     */
    function denomUnitArrayToJson(
        uint256[] memory decimalsArr,
        string[] memory symbols,
        bool[] memory isDefaultDisplays,
        string[] memory metadataJsons
    ) internal pure returns (string memory) {
        require(
            decimalsArr.length == symbols.length &&
            symbols.length == isDefaultDisplays.length &&
            isDefaultDisplays.length == metadataJsons.length,
            "Array lengths must match"
        );

        if (decimalsArr.length == 0) {
            return "[]";
        }

        string memory result = "[";
        for (uint256 i = 0; i < decimalsArr.length; i++) {
            if (i > 0) {
                result = string(abi.encodePacked(result, ","));
            }
            result = string(abi.encodePacked(
                result,
                denomUnitToJson(decimalsArr[i], symbols[i], isDefaultDisplays[i], metadataJsons[i])
            ));
        }
        return string(abi.encodePacked(result, "]"));
    }

    /**
     * @notice Convert CosmosCoinBackedPath to JSON
     * @dev Proto fields: address, conversion
     * @param addressStr The backing address
     * @param conversionJson Pre-encoded Conversion JSON (typically just a ratio)
     */
    function cosmosCoinBackedPathToJson(
        string memory addressStr,
        string memory conversionJson
    ) internal pure returns (string memory) {
        string memory result = "{";

        if (bytes(addressStr).length > 0) {
            result = string(abi.encodePacked(result, '"address":"', _escapeJsonString(addressStr), '"'));
        }

        if (bytes(conversionJson).length > 0 && keccak256(bytes(conversionJson)) != keccak256(bytes("{}"))) {
            if (bytes(addressStr).length > 0) {
                result = string(abi.encodePacked(result, ","));
            }
            result = string(abi.encodePacked(result, '"conversion":', conversionJson));
        }

        return string(abi.encodePacked(result, '}'));
    }

    /**
     * @notice Convert CosmosCoinWrapperPath to JSON
     * @dev Proto fields: address, denom, conversion, symbol, denomUnits[], allowOverrideWithAnyValidToken, metadata
     * @param addressStr The wrapper address
     * @param denom The cosmos coin denomination
     * @param conversionJson Pre-encoded Conversion JSON
     * @param symbol Token symbol
     * @param denomUnitsJson Pre-encoded DenomUnit[] JSON
     * @param allowOverrideWithAnyValidToken Whether to allow override
     * @param metadataJson Pre-encoded PathMetadata JSON
     */
    function cosmosCoinWrapperPathToJson(
        string memory addressStr,
        string memory denom,
        string memory conversionJson,
        string memory symbol,
        string memory denomUnitsJson,
        bool allowOverrideWithAnyValidToken,
        string memory metadataJson
    ) internal pure returns (string memory) {
        string memory result = string(abi.encodePacked(
            '{"address":"', _escapeJsonString(addressStr),
            '","denom":"', _escapeJsonString(denom), '"'
        ));

        if (bytes(conversionJson).length > 0 && keccak256(bytes(conversionJson)) != keccak256(bytes("{}"))) {
            result = string(abi.encodePacked(result, ',"conversion":', conversionJson));
        }

        result = string(abi.encodePacked(
            result,
            ',"symbol":"', _escapeJsonString(symbol), '"'
        ));

        if (bytes(denomUnitsJson).length > 0 && keccak256(bytes(denomUnitsJson)) != keccak256(bytes("[]"))) {
            result = string(abi.encodePacked(result, ',"denomUnits":', denomUnitsJson));
        }

        result = string(abi.encodePacked(
            result,
            ',"allowOverrideWithAnyValidToken":', allowOverrideWithAnyValidToken ? 'true' : 'false'
        ));

        if (bytes(metadataJson).length > 0 && keccak256(bytes(metadataJson)) != keccak256(bytes("{}"))) {
            result = string(abi.encodePacked(result, ',"metadata":', metadataJson));
        }

        return string(abi.encodePacked(result, '}'));
    }

    /**
     * @notice Convert AliasPath to JSON
     * @dev Proto fields: denom, conversion, symbol, denomUnits[], metadata
     * @param denom The alias denomination
     * @param conversionJson Pre-encoded Conversion JSON
     * @param symbol Token symbol
     * @param denomUnitsJson Pre-encoded DenomUnit[] JSON
     * @param metadataJson Pre-encoded PathMetadata JSON
     */
    function aliasPathToJson(
        string memory denom,
        string memory conversionJson,
        string memory symbol,
        string memory denomUnitsJson,
        string memory metadataJson
    ) internal pure returns (string memory) {
        string memory result = string(abi.encodePacked(
            '{"denom":"', _escapeJsonString(denom), '"'
        ));

        if (bytes(conversionJson).length > 0 && keccak256(bytes(conversionJson)) != keccak256(bytes("{}"))) {
            result = string(abi.encodePacked(result, ',"conversion":', conversionJson));
        }

        result = string(abi.encodePacked(
            result,
            ',"symbol":"', _escapeJsonString(symbol), '"'
        ));

        if (bytes(denomUnitsJson).length > 0 && keccak256(bytes(denomUnitsJson)) != keccak256(bytes("[]"))) {
            result = string(abi.encodePacked(result, ',"denomUnits":', denomUnitsJson));
        }

        if (bytes(metadataJson).length > 0 && keccak256(bytes(metadataJson)) != keccak256(bytes("{}"))) {
            result = string(abi.encodePacked(result, ',"metadata":', metadataJson));
        }

        return string(abi.encodePacked(result, '}'));
    }

    /**
     * @notice Convert CollectionApproval to JSON
     * @dev Proto fields: fromListId, toListId, initiatedByListId, transferTimes[], tokenIds[],
     *      ownershipTimes[], approvalId, uri, customData, approvalCriteria
     */
    function collectionApprovalToJson(
        string memory fromListId,
        string memory toListId,
        string memory initiatedByListId,
        string memory transferTimesJson,
        string memory tokenIdsJson,
        string memory ownershipTimesJson,
        string memory approvalId,
        string memory uri,
        string memory customData,
        string memory approvalCriteriaJson
    ) internal pure returns (string memory) {
        // Build in parts to avoid stack too deep
        string memory part1 = string(abi.encodePacked(
            '{"fromListId":"', _escapeJsonString(fromListId),
            '","toListId":"', _escapeJsonString(toListId),
            '","initiatedByListId":"', _escapeJsonString(initiatedByListId), '"'
        ));

        string memory part2 = string(abi.encodePacked(
            ',"transferTimes":', transferTimesJson,
            ',"badgeIds":', tokenIdsJson,
            ',"ownershipTimes":', ownershipTimesJson
        ));

        string memory part3 = string(abi.encodePacked(
            ',"approvalId":"', _escapeJsonString(approvalId), '"'
        ));

        if (bytes(uri).length > 0) {
            part3 = string(abi.encodePacked(part3, ',"uri":"', _escapeJsonString(uri), '"'));
        }

        if (bytes(customData).length > 0) {
            part3 = string(abi.encodePacked(part3, ',"customData":"', _escapeJsonString(customData), '"'));
        }

        if (bytes(approvalCriteriaJson).length > 0 && keccak256(bytes(approvalCriteriaJson)) != keccak256(bytes("{}"))) {
            part3 = string(abi.encodePacked(part3, ',"approvalCriteria":', approvalCriteriaJson));
        }

        return string(abi.encodePacked(part1, part2, part3, '}'));
    }

    /**
     * @notice Convert CollectionApproval array to JSON array
     */
    function collectionApprovalArrayToJson(
        string[] memory jsons
    ) internal pure returns (string memory) {
        if (jsons.length == 0) {
            return "[]";
        }

        string memory result = "[";
        for (uint256 i = 0; i < jsons.length; i++) {
            if (i > 0) {
                result = string(abi.encodePacked(result, ","));
            }
            result = string(abi.encodePacked(result, jsons[i]));
        }
        return string(abi.encodePacked(result, "]"));
    }

    /**
     * @notice Convert UserOutgoingApproval to JSON
     * @dev Similar to CollectionApproval but for user outgoing approvals
     */
    function userOutgoingApprovalToJson(
        string memory toListId,
        string memory initiatedByListId,
        string memory transferTimesJson,
        string memory tokenIdsJson,
        string memory ownershipTimesJson,
        string memory approvalId,
        string memory uri,
        string memory customData,
        string memory approvalCriteriaJson
    ) internal pure returns (string memory) {
        string memory part1 = string(abi.encodePacked(
            '{"toListId":"', _escapeJsonString(toListId),
            '","initiatedByListId":"', _escapeJsonString(initiatedByListId), '"'
        ));

        string memory part2 = string(abi.encodePacked(
            ',"transferTimes":', transferTimesJson,
            ',"badgeIds":', tokenIdsJson,
            ',"ownershipTimes":', ownershipTimesJson
        ));

        string memory part3 = string(abi.encodePacked(
            ',"approvalId":"', _escapeJsonString(approvalId), '"'
        ));

        if (bytes(uri).length > 0) {
            part3 = string(abi.encodePacked(part3, ',"uri":"', _escapeJsonString(uri), '"'));
        }

        if (bytes(customData).length > 0) {
            part3 = string(abi.encodePacked(part3, ',"customData":"', _escapeJsonString(customData), '"'));
        }

        if (bytes(approvalCriteriaJson).length > 0 && keccak256(bytes(approvalCriteriaJson)) != keccak256(bytes("{}"))) {
            part3 = string(abi.encodePacked(part3, ',"approvalCriteria":', approvalCriteriaJson));
        }

        return string(abi.encodePacked(part1, part2, part3, '}'));
    }

    /**
     * @notice Convert UserIncomingApproval to JSON
     * @dev Similar to CollectionApproval but for user incoming approvals
     */
    function userIncomingApprovalToJson(
        string memory fromListId,
        string memory initiatedByListId,
        string memory transferTimesJson,
        string memory tokenIdsJson,
        string memory ownershipTimesJson,
        string memory approvalId,
        string memory uri,
        string memory customData,
        string memory approvalCriteriaJson
    ) internal pure returns (string memory) {
        string memory part1 = string(abi.encodePacked(
            '{"fromListId":"', _escapeJsonString(fromListId),
            '","initiatedByListId":"', _escapeJsonString(initiatedByListId), '"'
        ));

        string memory part2 = string(abi.encodePacked(
            ',"transferTimes":', transferTimesJson,
            ',"badgeIds":', tokenIdsJson,
            ',"ownershipTimes":', ownershipTimesJson
        ));

        string memory part3 = string(abi.encodePacked(
            ',"approvalId":"', _escapeJsonString(approvalId), '"'
        ));

        if (bytes(uri).length > 0) {
            part3 = string(abi.encodePacked(part3, ',"uri":"', _escapeJsonString(uri), '"'));
        }

        if (bytes(customData).length > 0) {
            part3 = string(abi.encodePacked(part3, ',"customData":"', _escapeJsonString(customData), '"'));
        }

        if (bytes(approvalCriteriaJson).length > 0 && keccak256(bytes(approvalCriteriaJson)) != keccak256(bytes("{}"))) {
            part3 = string(abi.encodePacked(part3, ',"approvalCriteria":', approvalCriteriaJson));
        }

        return string(abi.encodePacked(part1, part2, part3, '}'));
    }

    /**
     * @notice Convert UserPermissions to JSON
     * @dev Proto fields: canUpdateIncomingApprovals[], canUpdateOutgoingApprovals[],
     *      canUpdateAutoApproveSelfInitiatedIncomingTransfers[],
     *      canUpdateAutoApproveSelfInitiatedOutgoingTransfers[],
     *      canUpdateAutoApproveAllIncomingTransfers[]
     */
    function userPermissionsToJson(
        string memory canUpdateIncomingApprovalsJson,
        string memory canUpdateOutgoingApprovalsJson,
        string memory canUpdateAutoApproveSelfInitiatedIncomingJson,
        string memory canUpdateAutoApproveSelfInitiatedOutgoingJson,
        string memory canUpdateAutoApproveAllIncomingJson
    ) internal pure returns (string memory) {
        string memory result = "{";
        bool hasField = false;

        if (bytes(canUpdateIncomingApprovalsJson).length > 0 && keccak256(bytes(canUpdateIncomingApprovalsJson)) != keccak256(bytes("[]"))) {
            result = string(abi.encodePacked(result, '"canUpdateIncomingApprovals":', canUpdateIncomingApprovalsJson));
            hasField = true;
        }

        if (bytes(canUpdateOutgoingApprovalsJson).length > 0 && keccak256(bytes(canUpdateOutgoingApprovalsJson)) != keccak256(bytes("[]"))) {
            if (hasField) result = string(abi.encodePacked(result, ","));
            result = string(abi.encodePacked(result, '"canUpdateOutgoingApprovals":', canUpdateOutgoingApprovalsJson));
            hasField = true;
        }

        if (bytes(canUpdateAutoApproveSelfInitiatedIncomingJson).length > 0 && keccak256(bytes(canUpdateAutoApproveSelfInitiatedIncomingJson)) != keccak256(bytes("[]"))) {
            if (hasField) result = string(abi.encodePacked(result, ","));
            result = string(abi.encodePacked(result, '"canUpdateAutoApproveSelfInitiatedIncomingTransfers":', canUpdateAutoApproveSelfInitiatedIncomingJson));
            hasField = true;
        }

        if (bytes(canUpdateAutoApproveSelfInitiatedOutgoingJson).length > 0 && keccak256(bytes(canUpdateAutoApproveSelfInitiatedOutgoingJson)) != keccak256(bytes("[]"))) {
            if (hasField) result = string(abi.encodePacked(result, ","));
            result = string(abi.encodePacked(result, '"canUpdateAutoApproveSelfInitiatedOutgoingTransfers":', canUpdateAutoApproveSelfInitiatedOutgoingJson));
            hasField = true;
        }

        if (bytes(canUpdateAutoApproveAllIncomingJson).length > 0 && keccak256(bytes(canUpdateAutoApproveAllIncomingJson)) != keccak256(bytes("[]"))) {
            if (hasField) result = string(abi.encodePacked(result, ","));
            result = string(abi.encodePacked(result, '"canUpdateAutoApproveAllIncomingTransfers":', canUpdateAutoApproveAllIncomingJson));
        }

        return string(abi.encodePacked(result, '}'));
    }

    /**
     * @notice Convert CollectionPermissions to JSON
     * @dev Proto fields: canDeleteCollection[], canArchiveCollection[], canUpdateStandards[],
     *      canUpdateCustomData[], canUpdateManager[], canUpdateCollectionMetadata[],
     *      canUpdateValidTokenIds[], canUpdateTokenMetadata[], canUpdateCollectionApprovals[],
     *      canAddMoreAliasPaths[], canAddMoreCosmosCoinWrapperPaths[]
     */
    function collectionPermissionsToJson(
        string memory canDeleteCollectionJson,
        string memory canArchiveCollectionJson,
        string memory canUpdateStandardsJson,
        string memory canUpdateCustomDataJson,
        string memory canUpdateManagerJson,
        string memory canUpdateCollectionMetadataJson,
        string memory canUpdateValidTokenIdsJson,
        string memory canUpdateTokenMetadataJson,
        string memory canUpdateCollectionApprovalsJson,
        string memory canAddMoreAliasPathsJson,
        string memory canAddMoreCosmosCoinWrapperPathsJson
    ) internal pure returns (string memory) {
        string memory result = "{";
        bool hasField = false;

        if (bytes(canDeleteCollectionJson).length > 0 && keccak256(bytes(canDeleteCollectionJson)) != keccak256(bytes("[]"))) {
            result = string(abi.encodePacked(result, '"canDeleteCollection":', canDeleteCollectionJson));
            hasField = true;
        }

        if (bytes(canArchiveCollectionJson).length > 0 && keccak256(bytes(canArchiveCollectionJson)) != keccak256(bytes("[]"))) {
            if (hasField) result = string(abi.encodePacked(result, ","));
            result = string(abi.encodePacked(result, '"canArchiveCollection":', canArchiveCollectionJson));
            hasField = true;
        }

        if (bytes(canUpdateStandardsJson).length > 0 && keccak256(bytes(canUpdateStandardsJson)) != keccak256(bytes("[]"))) {
            if (hasField) result = string(abi.encodePacked(result, ","));
            result = string(abi.encodePacked(result, '"canUpdateStandards":', canUpdateStandardsJson));
            hasField = true;
        }

        if (bytes(canUpdateCustomDataJson).length > 0 && keccak256(bytes(canUpdateCustomDataJson)) != keccak256(bytes("[]"))) {
            if (hasField) result = string(abi.encodePacked(result, ","));
            result = string(abi.encodePacked(result, '"canUpdateCustomData":', canUpdateCustomDataJson));
            hasField = true;
        }

        if (bytes(canUpdateManagerJson).length > 0 && keccak256(bytes(canUpdateManagerJson)) != keccak256(bytes("[]"))) {
            if (hasField) result = string(abi.encodePacked(result, ","));
            result = string(abi.encodePacked(result, '"canUpdateManager":', canUpdateManagerJson));
            hasField = true;
        }

        if (bytes(canUpdateCollectionMetadataJson).length > 0 && keccak256(bytes(canUpdateCollectionMetadataJson)) != keccak256(bytes("[]"))) {
            if (hasField) result = string(abi.encodePacked(result, ","));
            result = string(abi.encodePacked(result, '"canUpdateCollectionMetadata":', canUpdateCollectionMetadataJson));
            hasField = true;
        }

        if (bytes(canUpdateValidTokenIdsJson).length > 0 && keccak256(bytes(canUpdateValidTokenIdsJson)) != keccak256(bytes("[]"))) {
            if (hasField) result = string(abi.encodePacked(result, ","));
            result = string(abi.encodePacked(result, '"canUpdateValidTokenIds":', canUpdateValidTokenIdsJson));
            hasField = true;
        }

        if (bytes(canUpdateTokenMetadataJson).length > 0 && keccak256(bytes(canUpdateTokenMetadataJson)) != keccak256(bytes("[]"))) {
            if (hasField) result = string(abi.encodePacked(result, ","));
            result = string(abi.encodePacked(result, '"canUpdateTokenMetadata":', canUpdateTokenMetadataJson));
            hasField = true;
        }

        if (bytes(canUpdateCollectionApprovalsJson).length > 0 && keccak256(bytes(canUpdateCollectionApprovalsJson)) != keccak256(bytes("[]"))) {
            if (hasField) result = string(abi.encodePacked(result, ","));
            result = string(abi.encodePacked(result, '"canUpdateCollectionApprovals":', canUpdateCollectionApprovalsJson));
            hasField = true;
        }

        if (bytes(canAddMoreAliasPathsJson).length > 0 && keccak256(bytes(canAddMoreAliasPathsJson)) != keccak256(bytes("[]"))) {
            if (hasField) result = string(abi.encodePacked(result, ","));
            result = string(abi.encodePacked(result, '"canAddMoreAliasPaths":', canAddMoreAliasPathsJson));
            hasField = true;
        }

        if (bytes(canAddMoreCosmosCoinWrapperPathsJson).length > 0 && keccak256(bytes(canAddMoreCosmosCoinWrapperPathsJson)) != keccak256(bytes("[]"))) {
            if (hasField) result = string(abi.encodePacked(result, ","));
            result = string(abi.encodePacked(result, '"canAddMoreCosmosCoinWrapperPaths":', canAddMoreCosmosCoinWrapperPathsJson));
        }

        return string(abi.encodePacked(result, '}'));
    }

    /**
     * @notice Convert Conversion ratio to JSON
     * @dev Proto fields: numerator, denominator
     * @param numerator The numerator of the conversion ratio
     * @param denominator The denominator of the conversion ratio
     */
    function conversionToJson(
        uint256 numerator,
        uint256 denominator
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"numerator":"', _uintToString(numerator),
            '","denominator":"', _uintToString(denominator), '"}'
        ));
    }

    // ============ Query Helpers ============

    /**
     * @notice Construct JSON for getCollectionStats
     * @param collectionId The collection ID
     */
    function getCollectionStatsJSON(
        uint256 collectionId
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"collectionId":"', _uintToString(collectionId), '"}'
        ));
    }

    /**
     * @notice Construct JSON for getApprovalTracker
     * @param collectionId The collection ID
     * @param approvalLevel "collection", "incoming", or "outgoing"
     * @param approverAddress Approver address (empty for collection level)
     * @param approvalId The approval ID
     * @param trackerType "amountsTracker" or "numTransfersTracker"
     * @param trackedAddress Address to check tracker for
     */
    function getApprovalTrackerJSON(
        uint256 collectionId,
        string memory approvalLevel,
        string memory approverAddress,
        string memory approvalId,
        string memory trackerType,
        string memory trackedAddress
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"collectionId":"', _uintToString(collectionId),
            '","approvalLevel":"', approvalLevel,
            '","approverAddress":"', approverAddress,
            '","approvalId":"', approvalId,
            '","trackerType":"', trackerType,
            '","trackedAddress":"', trackedAddress, '"}'
        ));
    }

    /**
     * @notice Construct JSON for getChallengeTracker
     * @param collectionId The collection ID
     * @param approvalLevel "collection", "incoming", or "outgoing"
     * @param approverAddress Approver address (empty for collection level)
     * @param approvalId The approval ID
     * @param challengeId The challenge ID
     * @param leafIndex The leaf index in the Merkle tree
     */
    function getChallengeTrackerJSON(
        uint256 collectionId,
        string memory approvalLevel,
        string memory approverAddress,
        string memory approvalId,
        string memory challengeId,
        uint256 leafIndex
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"collectionId":"', _uintToString(collectionId),
            '","approvalLevel":"', approvalLevel,
            '","approverAddress":"', approverAddress,
            '","approvalId":"', approvalId,
            '","challengeId":"', challengeId,
            '","leafIndex":"', _uintToString(leafIndex), '"}'
        ));
    }

    /**
     * @notice Construct JSON for getWrappableBalances
     * @param denom The denomination to check
     * @param userAddress The user address
     */
    function getWrappableBalancesJSON(
        string memory denom,
        address userAddress
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"denom":"', denom,
            '","address":"', _addressToString(userAddress), '"}'
        ));
    }

    /**
     * @notice Construct JSON for getVote
     * @param collectionId The collection ID
     * @param approvalLevel "collection", "incoming", or "outgoing"
     * @param approverAddress Approver address (empty for collection level)
     * @param approvalId The approval ID
     * @param proposalId The proposal ID
     * @param voterAddress The voter address
     */
    function getVoteJSON(
        uint256 collectionId,
        string memory approvalLevel,
        string memory approverAddress,
        string memory approvalId,
        string memory proposalId,
        address voterAddress
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"collectionId":"', _uintToString(collectionId),
            '","approvalLevel":"', approvalLevel,
            '","approverAddress":"', approverAddress,
            '","approvalId":"', approvalId,
            '","proposalId":"', proposalId,
            '","voterAddress":"', _addressToString(voterAddress), '"}'
        ));
    }

    /**
     * @notice Construct JSON for getVotes (paginated)
     * @param collectionId The collection ID
     * @param approvalLevel "collection", "incoming", or "outgoing"
     * @param approverAddress Approver address (empty for collection level)
     * @param approvalId The approval ID
     * @param proposalId The proposal ID
     */
    function getVotesJSON(
        uint256 collectionId,
        string memory approvalLevel,
        string memory approverAddress,
        string memory approvalId,
        string memory proposalId
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"collectionId":"', _uintToString(collectionId),
            '","approvalLevel":"', approvalLevel,
            '","approverAddress":"', approverAddress,
            '","approvalId":"', approvalId,
            '","proposalId":"', proposalId, '"}'
        ));
    }

    // ============ Transaction Helpers - Approvals ============

    /**
     * @notice Construct JSON for setIncomingApproval
     * @param collectionId The collection ID
     * @param approvalJson Pre-encoded UserIncomingApproval JSON
     */
    function setIncomingApprovalJSON(
        uint256 collectionId,
        string memory approvalJson
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"collectionId":"', _uintToString(collectionId),
            '","approval":', approvalJson, '}'
        ));
    }

    /**
     * @notice Construct JSON for setOutgoingApproval
     * @param collectionId The collection ID
     * @param approvalJson Pre-encoded UserOutgoingApproval JSON
     */
    function setOutgoingApprovalJSON(
        uint256 collectionId,
        string memory approvalJson
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"collectionId":"', _uintToString(collectionId),
            '","approval":', approvalJson, '}'
        ));
    }

    /**
     * @notice Construct JSON for updateUserApprovals
     * @dev Set updateX flags to true to update the corresponding fields
     * @param collectionId The collection ID
     * @param updateOutgoingApprovals Whether to update outgoing approvals
     * @param outgoingApprovalsJson Pre-encoded outgoing approvals array JSON
     * @param updateIncomingApprovals Whether to update incoming approvals
     * @param incomingApprovalsJson Pre-encoded incoming approvals array JSON
     * @param updateAutoApproveSelfInitiatedOutgoingTransfers Whether to update this setting
     * @param autoApproveSelfInitiatedOutgoingTransfers New value
     * @param updateAutoApproveSelfInitiatedIncomingTransfers Whether to update this setting
     * @param autoApproveSelfInitiatedIncomingTransfers New value
     * @param updateAutoApproveAllIncomingTransfers Whether to update this setting
     * @param autoApproveAllIncomingTransfers New value
     * @param updateUserPermissions Whether to update user permissions
     * @param userPermissionsJson Pre-encoded user permissions JSON
     */
    function updateUserApprovalsJSON(
        uint256 collectionId,
        bool updateOutgoingApprovals,
        string memory outgoingApprovalsJson,
        bool updateIncomingApprovals,
        string memory incomingApprovalsJson,
        bool updateAutoApproveSelfInitiatedOutgoingTransfers,
        bool autoApproveSelfInitiatedOutgoingTransfers,
        bool updateAutoApproveSelfInitiatedIncomingTransfers,
        bool autoApproveSelfInitiatedIncomingTransfers,
        bool updateAutoApproveAllIncomingTransfers,
        bool autoApproveAllIncomingTransfers,
        bool updateUserPermissions,
        string memory userPermissionsJson
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"collectionId":"', _uintToString(collectionId),
            '","updateOutgoingApprovals":', updateOutgoingApprovals ? 'true' : 'false',
            ',"outgoingApprovals":', outgoingApprovalsJson,
            ',"updateIncomingApprovals":', updateIncomingApprovals ? 'true' : 'false',
            ',"incomingApprovals":', incomingApprovalsJson,
            ',"updateAutoApproveSelfInitiatedOutgoingTransfers":', updateAutoApproveSelfInitiatedOutgoingTransfers ? 'true' : 'false',
            ',"autoApproveSelfInitiatedOutgoingTransfers":', autoApproveSelfInitiatedOutgoingTransfers ? 'true' : 'false',
            ',"updateAutoApproveSelfInitiatedIncomingTransfers":', updateAutoApproveSelfInitiatedIncomingTransfers ? 'true' : 'false',
            ',"autoApproveSelfInitiatedIncomingTransfers":', autoApproveSelfInitiatedIncomingTransfers ? 'true' : 'false',
            ',"updateAutoApproveAllIncomingTransfers":', updateAutoApproveAllIncomingTransfers ? 'true' : 'false',
            ',"autoApproveAllIncomingTransfers":', autoApproveAllIncomingTransfers ? 'true' : 'false',
            ',"updateUserPermissions":', updateUserPermissions ? 'true' : 'false',
            ',"userPermissions":', userPermissionsJson, '}'
        ));
    }

    /**
     * @notice Construct JSON for purgeApprovals
     * @param collectionId The collection ID
     * @param purgeExpired Whether to purge expired approvals
     * @param approverAddress Address of user whose approvals to purge (empty for creator)
     * @param purgeCounterpartyApprovals Whether to purge counterparty approvals
     * @param approvalsToPurgeJson Pre-encoded ApprovalIdentifierDetails array JSON
     */
    function purgeApprovalsJSON(
        uint256 collectionId,
        bool purgeExpired,
        string memory approverAddress,
        bool purgeCounterpartyApprovals,
        string memory approvalsToPurgeJson
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"collectionId":"', _uintToString(collectionId),
            '","purgeExpired":', purgeExpired ? 'true' : 'false',
            ',"approverAddress":"', approverAddress,
            '","purgeCounterpartyApprovals":', purgeCounterpartyApprovals ? 'true' : 'false',
            ',"approvalsToPurge":', approvalsToPurgeJson, '}'
        ));
    }

    // ============ Transaction Helpers - Collection Updates ============

    /**
     * @notice Construct JSON for setCustomData
     * @param collectionId The collection ID
     * @param customData The custom data string
     * @param canUpdateCustomDataJson Pre-encoded ActionPermission array JSON
     */
    function setCustomDataJSON(
        uint256 collectionId,
        string memory customData,
        string memory canUpdateCustomDataJson
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"collectionId":"', _uintToString(collectionId),
            '","customData":"', customData,
            '","canUpdateCustomData":', canUpdateCustomDataJson, '}'
        ));
    }

    /**
     * @notice Construct JSON for setManager
     * @param collectionId The collection ID
     * @param manager The new manager address
     * @param canUpdateManagerJson Pre-encoded ActionPermission array JSON
     */
    function setManagerJSON(
        uint256 collectionId,
        string memory manager,
        string memory canUpdateManagerJson
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"collectionId":"', _uintToString(collectionId),
            '","manager":"', manager,
            '","canUpdateManager":', canUpdateManagerJson, '}'
        ));
    }

    /**
     * @notice Construct JSON for setStandards
     * @param collectionId The collection ID
     * @param standardsJson Pre-encoded string array JSON
     * @param canUpdateStandardsJson Pre-encoded ActionPermission array JSON
     */
    function setStandardsJSON(
        uint256 collectionId,
        string memory standardsJson,
        string memory canUpdateStandardsJson
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"collectionId":"', _uintToString(collectionId),
            '","standards":', standardsJson,
            ',"canUpdateStandards":', canUpdateStandardsJson, '}'
        ));
    }

    /**
     * @notice Construct JSON for setIsArchived
     * @param collectionId The collection ID
     * @param isArchived Whether the collection is archived
     * @param canArchiveCollectionJson Pre-encoded ActionPermission array JSON
     */
    function setIsArchivedJSON(
        uint256 collectionId,
        bool isArchived,
        string memory canArchiveCollectionJson
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"collectionId":"', _uintToString(collectionId),
            '","isArchived":', isArchived ? 'true' : 'false',
            ',"canArchiveCollection":', canArchiveCollectionJson, '}'
        ));
    }

    /**
     * @notice Construct JSON for setCollectionMetadata
     * @param collectionId The collection ID
     * @param collectionMetadataJson Pre-encoded CollectionMetadata JSON
     * @param canUpdateCollectionMetadataJson Pre-encoded ActionPermission array JSON
     */
    function setCollectionMetadataJSON(
        uint256 collectionId,
        string memory collectionMetadataJson,
        string memory canUpdateCollectionMetadataJson
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"collectionId":"', _uintToString(collectionId),
            '","collectionMetadata":', collectionMetadataJson,
            ',"canUpdateCollectionMetadata":', canUpdateCollectionMetadataJson, '}'
        ));
    }

    /**
     * @notice Construct JSON for setValidTokenIds
     * @param collectionId The collection ID
     * @param validTokenIdsJson Pre-encoded UintRange array JSON
     * @param canUpdateValidTokenIdsJson Pre-encoded TokenIdsActionPermission array JSON
     */
    function setValidTokenIdsJSON(
        uint256 collectionId,
        string memory validTokenIdsJson,
        string memory canUpdateValidTokenIdsJson
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"collectionId":"', _uintToString(collectionId),
            '","validTokenIds":', validTokenIdsJson,
            ',"canUpdateValidTokenIds":', canUpdateValidTokenIdsJson, '}'
        ));
    }

    /**
     * @notice Construct JSON for setTokenMetadata
     * @param collectionId The collection ID
     * @param tokenMetadataJson Pre-encoded TokenMetadata array JSON
     * @param canUpdateTokenMetadataJson Pre-encoded TokenIdsActionPermission array JSON
     */
    function setTokenMetadataJSON(
        uint256 collectionId,
        string memory tokenMetadataJson,
        string memory canUpdateTokenMetadataJson
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"collectionId":"', _uintToString(collectionId),
            '","tokenMetadata":', tokenMetadataJson,
            ',"canUpdateTokenMetadata":', canUpdateTokenMetadataJson, '}'
        ));
    }

    /**
     * @notice Construct JSON for setCollectionApprovals
     * @param collectionId The collection ID
     * @param collectionApprovalsJson Pre-encoded CollectionApproval array JSON
     * @param canUpdateCollectionApprovalsJson Pre-encoded CollectionApprovalPermission array JSON
     */
    function setCollectionApprovalsJSON(
        uint256 collectionId,
        string memory collectionApprovalsJson,
        string memory canUpdateCollectionApprovalsJson
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"collectionId":"', _uintToString(collectionId),
            '","collectionApprovals":', collectionApprovalsJson,
            ',"canUpdateCollectionApprovals":', canUpdateCollectionApprovalsJson, '}'
        ));
    }

    // ============ Transaction Helpers - Other ============

    /**
     * @notice Construct JSON for castVote
     * @param collectionId The collection ID
     * @param approvalLevel "collection", "incoming", or "outgoing"
     * @param approverAddress Approver address (empty for collection level)
     * @param approvalId The approval ID
     * @param proposalId The proposal ID
     * @param yesWeight The yes vote weight (0-100)
     */
    function castVoteJSON(
        uint256 collectionId,
        string memory approvalLevel,
        string memory approverAddress,
        string memory approvalId,
        string memory proposalId,
        uint256 yesWeight
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"collectionId":"', _uintToString(collectionId),
            '","approvalLevel":"', approvalLevel,
            '","approverAddress":"', approverAddress,
            '","approvalId":"', approvalId,
            '","proposalId":"', proposalId,
            '","yesWeight":"', _uintToString(yesWeight), '"}'
        ));
    }

    /**
     * @notice Construct JSON for createAddressLists
     * @param addressListsJson Pre-encoded AddressListInput array JSON
     */
    function createAddressListsJSON(
        string memory addressListsJson
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"addressLists":', addressListsJson, '}'
        ));
    }

    /**
     * @notice Construct JSON for a single AddressListInput
     * @param listId The list ID
     * @param addressesJson Pre-encoded addresses array JSON (strings)
     * @param whitelist True for whitelist, false for blacklist
     * @param uri Optional URI for the list
     * @param customData Optional custom data
     */
    function addressListInputToJson(
        string memory listId,
        string memory addressesJson,
        bool whitelist,
        string memory uri,
        string memory customData
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"listId":"', listId,
            '","addresses":', addressesJson,
            ',"whitelist":', whitelist ? 'true' : 'false',
            ',"uri":"', uri,
            '","customData":"', customData, '"}'
        ));
    }

    /**
     * @notice Construct JSON for ApprovalIdentifierDetails (used in purgeApprovals)
     * @param approvalLevel "collection", "incoming", or "outgoing"
     * @param approvalId The approval ID
     */
    function approvalIdentifierDetailsToJson(
        string memory approvalLevel,
        string memory approvalId
    ) internal pure returns (string memory) {
        return string(abi.encodePacked(
            '{"approvalLevel":"', approvalLevel,
            '","approvalId":"', approvalId, '"}'
        ));
    }

    /**
     * @notice Construct JSON array from ApprovalIdentifierDetails elements
     */
    function approvalIdentifierDetailsArrayToJson(
        string[] memory elements
    ) internal pure returns (string memory) {
        if (elements.length == 0) {
            return "[]";
        }
        string memory result = "[";
        for (uint256 i = 0; i < elements.length; i++) {
            if (i > 0) {
                result = string(abi.encodePacked(result, ","));
            }
            result = string(abi.encodePacked(result, elements[i]));
        }
        return string(abi.encodePacked(result, "]"));
    }

    // ============ Internal Helpers ============

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

