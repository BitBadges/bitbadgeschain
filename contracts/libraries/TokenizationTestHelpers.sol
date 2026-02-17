// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../types/sol";
import "./TokenizationHelpers.sol";

/**
 * @title TokenizationTestHelpers
 * @notice Testing utilities for tokenization contracts
 * @dev Provides mock data generators, test fixtures, and assertion helpers for writing tests.
 *      These helpers make it easier to write comprehensive tests for contracts using the tokenization precompile.
 * 
 * Usage example:
 * ```solidity
 * import "./libraries/TokenizationTestHelpers.sol";
 * 
 * function testTransfer() public {
 *     uint256 collectionId = 1;
 *     address recipient = address(0x123);
 *     uint256 amount = 100;
 *     
 *     // Generate test data
 *     UintRange[] memory tokenIds = 
 *         TokenizationTestHelpers.generateTokenIdRanges(1, 10);
 *     UintRange[] memory ownershipTimes = 
 *         TokenizationTestHelpers.generateFullOwnershipTimes();
 *     
 *     // Perform transfer
 *     // ...
 * }
 * ```
 */
library TokenizationTestHelpers {
    // ============ Mock Data Generators ============

    /**
     * @notice Generate an array of consecutive token ID ranges
     * @param startTokenId The first token ID
     * @param endTokenId The last token ID (inclusive)
     * @param rangeSize The size of each range (e.g., 10 means ranges of 10 tokens each)
     * @return ranges Array of token ID ranges
     */
    function generateTokenIdRanges(
        uint256 startTokenId,
        uint256 endTokenId,
        uint256 rangeSize
    ) internal pure returns (UintRange[] memory ranges) {
        require(startTokenId <= endTokenId, "TokenizationTestHelpers: startTokenId must be <= endTokenId");
        require(rangeSize > 0, "TokenizationTestHelpers: rangeSize must be > 0");
        
        uint256 numRanges = (endTokenId - startTokenId + rangeSize) / rangeSize;
        ranges = new UintRange[](numRanges);
        
        uint256 currentStart = startTokenId;
        for (uint256 i = 0; i < numRanges; i++) {
            uint256 currentEnd = currentStart + rangeSize - 1;
            if (currentEnd > endTokenId) {
                currentEnd = endTokenId;
            }
            ranges[i] = TokenizationHelpers.createUintRange(currentStart, currentEnd);
            currentStart = currentEnd + 1;
            if (currentStart > endTokenId) {
                break;
            }
        }
        return ranges;
    }

    /**
     * @notice Generate a single token ID range covering all tokens
     * @param startTokenId The first token ID
     * @param endTokenId The last token ID (inclusive)
     * @return range A single range covering all token IDs
     */
    function generateSingleTokenIdRange(
        uint256 startTokenId,
        uint256 endTokenId
    ) internal pure returns (UintRange[] memory range) {
        range = new UintRange[](1);
        range[0] = TokenizationHelpers.createUintRange(startTokenId, endTokenId);
        return range;
    }

    /**
     * @notice Generate ownership time ranges with full ownership (1 to max)
     * @return ranges Array with a single full ownership time range
     */
    function generateFullOwnershipTimes() internal pure returns (UintRange[] memory ranges) {
        ranges = new UintRange[](1);
        ranges[0] = TokenizationHelpers.createFullOwnershipTimeRange();
        return ranges;
    }

    /**
     * @notice Generate ownership time ranges from current time to a future time
     * @param duration The duration in seconds
     * @return ranges Array with a single ownership time range
     */
    function generateOwnershipTimesFromNow(
        uint256 duration
    ) internal view returns (UintRange[] memory ranges) {
        ranges = new UintRange[](1);
        ranges[0] = TokenizationHelpers.createOwnershipTimeRange(block.timestamp, duration);
        return ranges;
    }

    /**
     * @notice Generate ownership time ranges from a start time to an end time
     * @param startTime The start timestamp
     * @param endTime The end timestamp
     * @return ranges Array with a single ownership time range
     */
    function generateOwnershipTimes(
        uint256 startTime,
        uint256 endTime
    ) internal pure returns (UintRange[] memory ranges) {
        ranges = new UintRange[](1);
        ranges[0] = TokenizationHelpers.createUintRange(startTime, endTime);
        return ranges;
    }

    /**
     * @notice Generate a test collection metadata
     * @param name The collection name (used in URI)
     * @return metadata A CollectionMetadata struct with test data
     */
    function generateCollectionMetadata(
        string memory name
    ) internal pure returns (CollectionMetadata memory metadata) {
        string memory uri = string(abi.encodePacked("https://example.com/collections/", name));
        string memory customData = string(abi.encodePacked('{"name":"', name, '","test":true}'));
        return TokenizationHelpers.createCollectionMetadata(uri, customData);
    }

    /**
     * @notice Generate test token metadata
     * @param tokenId The token ID
     * @return metadata A TokenMetadata struct with test data
     */
    function generateTokenMetadata(
        uint256 tokenId
    ) internal pure returns (TokenMetadata memory metadata) {
        string memory uri = string(abi.encodePacked("https://example.com/tokens/", uintToString(tokenId)));
        string memory customData = string(abi.encodePacked('{"tokenId":', uintToString(tokenId), '}'));
        UintRange[] memory tokenIds = new UintRange[](1);
        tokenIds[0] = TokenizationHelpers.createSingleTokenIdRange(tokenId);
        return TokenizationHelpers.createTokenMetadata(uri, customData, tokenIds);
    }

    /**
     * @notice Generate an array of test addresses
     * @param count The number of addresses to generate
     * @return addresses Array of test addresses
     */
    function generateTestAddresses(uint256 count) internal pure returns (address[] memory addresses) {
        addresses = new address[](count);
        for (uint256 i = 0; i < count; i++) {
            // Generate deterministic test addresses
            addresses[i] = address(uint160(0x1000 + i));
        }
        return addresses;
    }

    /**
     * @notice Generate a test user balance store with default values
     * @return store A UserBalanceStore with default test values
     */
    function generateTestUserBalanceStore() internal pure returns (UserBalanceStore memory store) {
        return TokenizationHelpers.createEmptyUserBalanceStore();
    }

    /**
     * @notice Generate a test user balance store with auto-approve flags
     * @param autoApproveOutgoing Whether to auto-approve outgoing transfers
     * @param autoApproveIncoming Whether to auto-approve incoming transfers
     * @return store A UserBalanceStore with specified auto-approve flags
     */
    function generateTestUserBalanceStoreWithFlags(
        bool autoApproveOutgoing,
        bool autoApproveIncoming
    ) internal pure returns (UserBalanceStore memory store) {
        store = TokenizationHelpers.createEmptyUserBalanceStore();
        store.autoApproveSelfInitiatedOutgoingTransfers = autoApproveOutgoing;
        store.autoApproveSelfInitiatedIncomingTransfers = autoApproveIncoming;
        return store;
    }

    // ============ Test Fixtures ============

    /**
     * @notice Create a minimal collection configuration for testing
     * @param tokenIdStart The start token ID
     * @param tokenIdEnd The end token ID
     * @return validTokenIds Array of valid token ID ranges
     * @return metadata Collection metadata
     * @return defaultBalances Default user balance store
     */
    function createMinimalCollectionConfig(
        uint256 tokenIdStart,
        uint256 tokenIdEnd
    ) internal pure returns (
        UintRange[] memory validTokenIds,
        CollectionMetadata memory metadata,
        UserBalanceStore memory defaultBalances
    ) {
        validTokenIds = generateSingleTokenIdRange(tokenIdStart, tokenIdEnd);
        metadata = generateCollectionMetadata("TestCollection");
        defaultBalances = generateTestUserBalanceStore();
    }

    // ============ Assertion Helpers ============

    /**
     * @notice Assert that a UintRange is valid
     * @param range The range to validate
     * @param message Custom error message
     */
    function assertValidRange(
        UintRange memory range,
        string memory message
    ) internal pure {
        require(
            TokenizationHelpers.validateUintRange(range),
            string(abi.encodePacked("Invalid range: ", message))
        );
    }

    /**
     * @notice Assert that two addresses are equal
     * @param expected The expected address
     * @param actual The actual address
     * @param message Custom error message
     */
    function assertAddressEqual(
        address expected,
        address actual,
        string memory message
    ) internal pure {
        require(expected == actual, string(abi.encodePacked("Address mismatch: ", message)));
    }

    /**
     * @notice Assert that two uint256 values are equal
     * @param expected The expected value
     * @param actual The actual value
     * @param message Custom error message
     */
    function assertUintEqual(
        uint256 expected,
        uint256 actual,
        string memory message
    ) internal pure {
        require(expected == actual, string(abi.encodePacked("Uint mismatch: ", message)));
    }

    // ============ Helper Functions ============

    /**
     * @notice Convert uint256 to string (for test data generation)
     * @param value The uint256 value
     * @return str The string representation
     */
    function uintToString(uint256 value) internal pure returns (string memory str) {
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
}

