// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../types/TokenizationTypes.sol";

/**
 * @title TokenizationDecoders
 * @notice Utilities for decoding protobuf-encoded query responses
 * @dev IMPORTANT: Full protobuf decoding in Solidity is complex and gas-intensive.
 *      This library provides:
 *      1. Varint parsing helpers for reading protobuf fields
 *      2. Collection stats parsing (holderCount)
 *      3. Simple field extractors
 *
 *      For full struct decoding, use off-chain tools or the TypeScript SDK.
 *      See contracts/test/MaxUniqueHoldersChecker.sol for a complete example
 *      of using these helpers for EVM query challenges and invariants.
 *
 * v25 Changes:
 *      - Added protobuf parsing helpers (readVarint, varintSize, parseStringNumeric)
 *      - Added parseHolderCountFromStats for collection stats queries
 *      - These enable on-chain EVM query challenges and invariants
 *
 * Usage example:
 * ```solidity
 * import "./libraries/TokenizationDecoders.sol";
 *
 * // Parse holder count from collection stats
 * bytes memory statsBytes = precompile.getCollectionStats(json);
 * uint256 holderCount = TokenizationDecoders.parseHolderCountFromStats(statsBytes);
 *
 * // For full collection decoding, use direct queries or off-chain tools
 * uint256 balance = precompile.getBalanceAmount(json);
 * ```
 */
library TokenizationDecoders {
    // ============ Protobuf Parsing Helpers ============
    // These helpers enable on-chain parsing of simple protobuf fields

    /**
     * @notice Read a varint from protobuf data
     * @param data The protobuf-encoded bytes
     * @param start The starting position in the data
     * @return value The decoded varint value
     * @dev Varints use 7 bits per byte with MSB as continuation flag.
     *      Supports up to 10 bytes (70 bits) for uint64 compatibility.
     */
    function readVarint(bytes memory data, uint256 start) internal pure returns (uint256 value) {
        uint256 v = 0;
        for (uint256 i = start; i < data.length && i < start + 10; i++) {
            v |= uint256(uint8(data[i]) & 0x7f) << (7 * (i - start));
            if (data[i] < 0x80) return v;
        }
        return v;
    }

    /**
     * @notice Calculate the byte size of a varint at a position
     * @param data The protobuf-encoded bytes
     * @param start The starting position in the data
     * @return size The number of bytes used by the varint
     * @dev Scans for the first byte without MSB set (< 0x80)
     */
    function varintSize(bytes memory data, uint256 start) internal pure returns (uint256 size) {
        for (uint256 i = start; i < data.length && i < start + 10; i++) {
            if (data[i] < 0x80) return i - start + 1;
        }
        return 1;
    }

    /**
     * @notice Parse a numeric value from an ASCII string in protobuf data
     * @param data The protobuf-encoded bytes
     * @param start The starting position of the string
     * @param length The length of the string
     * @return value The parsed numeric value
     * @dev Parses ASCII digits (0x30-0x39) into a uint256.
     *      Used for protobuf string-encoded numbers like holderCount.
     */
    function parseStringNumeric(
        bytes memory data,
        uint256 start,
        uint256 length
    ) internal pure returns (uint256 value) {
        uint256 count = 0;
        for (uint256 j = 0; j < length && (start + j) < data.length; j++) {
            uint8 d = uint8(data[start + j]);
            if (d >= 48 && d <= 57) {
                count = count * 10 + (d - 48);
            }
        }
        return count;
    }

    // ============ Collection Stats Parsing ============

    /**
     * @notice Parse holderCount from QueryGetCollectionStatsResponse
     * @param data The protobuf-encoded response from getCollectionStats()
     * @return holderCount The number of unique holders for the collection
     * @dev Response structure:
     *      QueryGetCollectionStatsResponse {
     *        field 1 (stats): CollectionStats {
     *          field 1 (holderCount): string (numeric)
     *          field 2 (balances): repeated Balance (circulating supply)
     *        }
     *      }
     *      Wire format: 0x0a <outer_len> 0x0a <inner_len> <ascii_digits>
     *
     * Example usage in EVM query challenge / invariant:
     * ```solidity
     * bytes memory response = TOKENIZATION.getCollectionStats(json);
     * uint256 holders = TokenizationDecoders.parseHolderCountFromStats(response);
     * require(holders <= maxAllowed, "Too many holders");
     * ```
     */
    function parseHolderCountFromStats(bytes memory data) internal pure returns (uint256 holderCount) {
        uint256 i = 0;

        // Empty or invalid response
        if (i >= data.length) return 0;

        // Expect field 1 (stats) with wire type 2 (length-delimited) = 0x0a
        if (data[i] != 0x0a) return 0;
        i++;

        // Read outer length (stats message length)
        uint256 outerLen = readVarint(data, i);
        i += varintSize(data, i);
        if (i + outerLen > data.length) return 0;
        uint256 innerEnd = i + outerLen;

        // Inside stats message, expect field 1 (holderCount) with wire type 2 = 0x0a
        if (i >= data.length || data[i] != 0x0a) return 0;
        i++;
        if (i >= innerEnd) return 0;

        // Read string length
        uint256 strLen = readVarint(data, i);
        i += varintSize(data, i);
        if (i + strLen > innerEnd) return 0;

        // Parse ASCII digits to number
        return parseStringNumeric(data, i, strLen);
    }

    // ============ Full Struct Decoders (Not Supported) ============
    // Full protobuf decoding is not practical on-chain. Use off-chain tools.

    /**
     * @notice Decode a collection from protobuf bytes
     * @param data The protobuf-encoded TokenCollection bytes
     * @return collection The decoded TokenCollection struct
     * @dev WARNING: Full protobuf decoding is not implemented in Solidity.
     *      Use the TypeScript SDK for full decoding, or use specific field
     *      extractors like parseHolderCountFromStats for needed fields.
     */
    function decodeCollection(bytes memory data)
        internal pure returns (TokenCollection memory collection)
    {
        revert("TokenizationDecoders: Full protobuf decoding not supported. Use TypeScript SDK or field extractors.");
    }

    /**
     * @notice Decode a balance from protobuf bytes
     * @param data The protobuf-encoded UserBalanceStore bytes
     * @return balance The decoded UserBalanceStore struct
     * @dev WARNING: Full protobuf decoding is not implemented.
     *      Use getBalanceAmount() for direct balance queries instead.
     */
    function decodeBalance(bytes memory data)
        internal pure returns (UserBalanceStore memory balance)
    {
        revert("TokenizationDecoders: Full protobuf decoding not supported. Use getBalanceAmount() instead.");
    }

    /**
     * @notice Decode an address list from protobuf bytes
     * @param data The protobuf-encoded AddressList bytes
     * @return list The decoded AddressList struct
     * @dev WARNING: Full protobuf decoding is not implemented.
     */
    function decodeAddressList(bytes memory data)
        internal pure returns (AddressList memory list)
    {
        revert("TokenizationDecoders: Full protobuf decoding not supported. Use off-chain tools.");
    }

    /**
     * @notice Decode a dynamic store from protobuf bytes
     * @param data The protobuf-encoded DynamicStore bytes
     * @return store The decoded DynamicStore struct
     * @dev WARNING: Full protobuf decoding is not implemented.
     */
    function decodeDynamicStore(bytes memory data)
        internal pure returns (DynamicStore memory store)
    {
        revert("TokenizationDecoders: Full protobuf decoding not supported. Use off-chain tools.");
    }

    // ============ Simple Field Extractors ============

    /**
     * @notice Check if protobuf data is empty
     * @param data The protobuf-encoded bytes
     * @return isEmpty_ True if data is empty or zero-length
     */
    function isEmpty(bytes memory data) internal pure returns (bool isEmpty_) {
        return data.length == 0;
    }

    /**
     * @notice Get the length of protobuf data
     * @param data The protobuf-encoded bytes
     * @return length The length of the data
     */
    function getLength(bytes memory data) internal pure returns (uint256 length) {
        return data.length;
    }

    // ============ Recommendations ============

    /**
     * @notice Recommended approaches for working with query responses
     * @dev For production contracts:
     *
     *      1. Use direct queries when available:
     *         - getBalanceAmount() returns uint256 directly
     *         - getTotalSupply() returns uint256 directly
     *
     *      2. Use field extractors for specific data:
     *         - parseHolderCountFromStats() for collection holder count
     *
     *      3. For EVM query challenges / invariants:
     *         See contracts/test/MaxUniqueHoldersChecker.sol for a complete example
     *         that uses parseHolderCountFromStats() to validate holder limits.
     *
     *      4. For complex data, decode off-chain:
     *         - Use TypeScript SDK for full struct decoding
     *         - Pass decoded data to contracts via function parameters
     *         - Use events to emit raw bytes for off-chain indexing
     *
     * Example patterns:
     * ```solidity
     * // Pattern 1: Direct query
     * uint256 balance = precompile.getBalanceAmount(json);
     *
     * // Pattern 2: Field extraction
     * bytes memory stats = precompile.getCollectionStats(json);
     * uint256 holders = TokenizationDecoders.parseHolderCountFromStats(stats);
     *
     * // Pattern 3: Invariant checker contract
     * // See MaxUniqueHoldersChecker.sol
     * ```
     */
    function recommendedApproach() internal pure {
        // Documentation-only function
    }
}

