// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/ITokenizationPrecompile.sol";

/**
 * @title MaxUniqueHoldersChecker
 * @notice Reference implementation for EVM query challenges and collection invariants (v25)
 * @dev This contract demonstrates how to use the v25 EVM query challenge feature to enforce
 *      post-transfer invariants. It checks that a collection's unique holder count does not
 *      exceed a maximum allowed value.
 *
 * ## Overview
 * EVM query challenges allow you to execute read-only contract calls as part of:
 * 1. Approval criteria (pre-transfer validation)
 * 2. Collection invariants (post-transfer validation)
 *
 * This example focuses on collection invariants - rules that must hold true after every transfer.
 *
 * ## How It Works
 * 1. This contract queries `getCollectionStats()` to get the holder count
 * 2. It parses the protobuf response to extract the numeric holder count
 * 3. It returns `bytes32(1)` if the count is within limits, `bytes32(0)` otherwise
 *
 * ## Usage in Collection Invariants
 * When creating a collection, add an EVM query challenge invariant:
 *
 * ```solidity
 * // Using TokenizationHelpers
 * EVMQueryChallenge memory challenge = TokenizationHelpers.createEVMQueryChallenge(
 *     address(maxHoldersChecker),           // This contract's address
 *     abi.encodeWithSelector(               // Calldata
 *         MaxUniqueHoldersChecker.checkMaxHolders.selector,
 *         "$collectionId",                  // Placeholder replaced at runtime
 *         100                               // Max 100 holders
 *     ),
 *     bytes32(uint256(1)),                  // Expected: pass (1)
 *     "eq",                                 // Comparison operator
 *     200000                                // Gas limit
 * );
 *
 * CollectionInvariants memory invariants = TokenizationHelpers.createInvariantsWithEVMChallenges(
 *     new EVMQueryChallenge[](1)
 * );
 * invariants.evmQueryChallenges[0] = challenge;
 * ```
 *
 * ## Gas Limits
 * - Default: 100,000 gas per challenge
 * - Maximum: 500,000 gas per challenge
 * - Total maximum: 1,000,000 gas across all challenges
 *
 * ## Protobuf Parsing
 * This contract demonstrates on-chain protobuf parsing for simple numeric fields.
 * For complex types, use the TypeScript SDK for off-chain decoding.
 * See TokenizationDecoders.sol for reusable parsing helpers.
 *
 * ## See Also
 * - TokenizationDecoders.sol: Reusable protobuf parsing helpers
 * - TokenizationHelpers.sol: Helper functions for creating EVM query challenges
 * - MinBankBalanceChecker.sol: Example of approval-level EVM query challenges
 */
contract MaxUniqueHoldersChecker {
    /// @notice Tokenization precompile address (always 0x1001)
    ITokenizationPrecompile constant TOKENIZATION = ITokenizationPrecompile(0x0000000000000000000000000000000000001001);

    /// @notice Check that collection's unique holder count is <= maxAllowed
    /// @param collectionId Collection ID (use $collectionId placeholder in invariant calldata)
    /// @param maxAllowed Maximum allowed number of unique holders
    /// @return Pass: bytes32(1) if holderCount <= maxAllowed, bytes32(0) otherwise
    /// @dev This function is called by the chain during post-transfer invariant validation.
    ///      The $collectionId placeholder in the calldata is replaced with the actual collection ID.
    ///
    /// Example calldata for invariant:
    /// ```
    /// abi.encodeWithSelector(
    ///     MaxUniqueHoldersChecker.checkMaxHolders.selector,
    ///     0,    // Placeholder for $collectionId
    ///     100   // Max 100 holders
    /// )
    /// ```
    function checkMaxHolders(uint256 collectionId, uint256 maxAllowed) external view returns (bytes32) {
        // Build JSON query for getCollectionStats
        string memory json = string(abi.encodePacked('{"collectionId":"', _uintToString(collectionId), '"}'));

        // Query collection stats from precompile
        bytes memory response = TOKENIZATION.getCollectionStats(json);
        if (response.length == 0) return bytes32(uint256(0));

        // Parse holder count from protobuf response
        uint256 holderCount = _parseHolderCountFromStats(response);

        // Return pass (1) or fail (0)
        return holderCount <= maxAllowed ? bytes32(uint256(1)) : bytes32(uint256(0));
    }

    /// @notice Parse holderCount from QueryGetCollectionStatsResponse protobuf
    /// @param data The protobuf-encoded response from getCollectionStats()
    /// @return The holder count as uint256
    /// @dev Response structure:
    ///      QueryGetCollectionStatsResponse {
    ///        field 1 (stats): CollectionStats {
    ///          field 1 (holderCount): string (numeric ASCII)
    ///        }
    ///      }
    ///      Wire format: 0x0a <outer_len_varint> 0x0a <inner_len_varint> <holderCount_ascii>
    ///
    ///      For a reusable implementation, see TokenizationDecoders.parseHolderCountFromStats()
    function _parseHolderCountFromStats(bytes memory data) internal pure returns (uint256) {
        uint256 i = 0;
        if (i >= data.length) return 0;
        if (data[i] != 0x0a) return 0;
        i++;
        uint256 outerLen = _readVarint(data, i);
        i += _varintSize(data, i);
        if (i + outerLen > data.length) return 0;
        uint256 innerEnd = i + outerLen;
        if (i >= data.length || data[i] != 0x0a) return 0;
        i++;
        if (i >= innerEnd) return 0;
        uint256 strLen = _readVarint(data, i);
        i += _varintSize(data, i);
        if (i + strLen > innerEnd) return 0;
        uint256 count = 0;
        for (uint256 j = 0; j < strLen && (i + j) < data.length; j++) {
            uint8 d = uint8(data[i + j]);
            if (d >= 48 && d <= 57) {
                count = count * 10 + (d - 48);
            }
        }
        return count;
    }

    function _readVarint(bytes memory data, uint256 start) internal pure returns (uint256) {
        uint256 v = 0;
        for (uint256 i = start; i < data.length && i < start + 10; i++) {
            v |= uint256(uint8(data[i]) & 0x7f) << (7 * (i - start));
            if (data[i] < 0x80) return v;
        }
        return v;
    }

    function _varintSize(bytes memory data, uint256 start) internal pure returns (uint256) {
        for (uint256 i = start; i < data.length && i < start + 10; i++) {
            if (data[i] < 0x80) return i - start + 1;
        }
        return 1;
    }

    function _uintToString(uint256 v) internal pure returns (string memory) {
        if (v == 0) return "0";
        uint256 j = v;
        uint256 len;
        while (j != 0) {
            len++;
            j /= 10;
        }
        bytes memory b = new bytes(len);
        uint256 k = len;
        while (v != 0) {
            k = k - 1;
            b[k] = bytes1(uint8(48 + v % 10));
            v /= 10;
        }
        return string(b);
    }
}
