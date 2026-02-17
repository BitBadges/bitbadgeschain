// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../types/sol";

/**
 * @title TokenizationDecoders
 * @notice Utilities for decoding protobuf-encoded query responses
 * @dev IMPORTANT: Full protobuf decoding in Solidity is complex and gas-intensive.
 *      This library provides simplified decoders for common use cases.
 *      For full decoding, consider using off-chain tools or the TypeScript SDK.
 * 
 *      Protobuf encoding uses variable-length encoding and nested structures,
 *      making complete on-chain decoding impractical for complex types.
 * 
 * Usage example:
 * ```solidity
 * import "./libraries/TokenizationDecoders.sol";
 * 
 * bytes memory collectionBytes = precompile.getCollection(json);
 * // Note: Full decoding may require off-chain tools
 * // For simple fields, use extractors below
 * ```
 */
library TokenizationDecoders {
    // ============ Limitations ============

    /**
     * @notice Decode a collection from protobuf bytes
     * @param data The protobuf-encoded TokenCollection bytes
     * @return collection The decoded TokenCollection struct
     * @dev WARNING: Full protobuf decoding is not implemented in Solidity.
     *      This is a placeholder that will revert. Use off-chain decoding tools.
     *      Consider using the TypeScript SDK for full decoding.
     */
    function decodeCollection(bytes memory data) 
        internal pure returns (TokenCollection memory collection) 
    {
        // Protobuf decoding in Solidity is extremely complex and gas-intensive
        // For production use, decode off-chain using the TypeScript SDK
        revert("TokenizationDecoders: Full protobuf decoding not supported in Solidity. Use off-chain tools.");
    }

    /**
     * @notice Decode a balance from protobuf bytes
     * @param data The protobuf-encoded UserBalanceStore bytes
     * @return balance The decoded UserBalanceStore struct
     * @dev WARNING: Full protobuf decoding is not implemented in Solidity.
     *      This is a placeholder that will revert. Use off-chain decoding tools.
     */
    function decodeBalance(bytes memory data) 
        internal pure returns (UserBalanceStore memory balance) 
    {
        revert("TokenizationDecoders: Full protobuf decoding not supported in Solidity. Use off-chain tools.");
    }

    /**
     * @notice Decode an address list from protobuf bytes
     * @param data The protobuf-encoded AddressList bytes
     * @return list The decoded AddressList struct
     * @dev WARNING: Full protobuf decoding is not implemented in Solidity.
     *      This is a placeholder that will revert. Use off-chain decoding tools.
     */
    function decodeAddressList(bytes memory data) 
        internal pure returns (AddressList memory list) 
    {
        revert("TokenizationDecoders: Full protobuf decoding not supported in Solidity. Use off-chain tools.");
    }

    /**
     * @notice Decode a dynamic store from protobuf bytes
     * @param data The protobuf-encoded DynamicStore bytes
     * @return store The decoded DynamicStore struct
     * @dev WARNING: Full protobuf decoding is not implemented in Solidity.
     *      This is a placeholder that will revert. Use off-chain decoding tools.
     */
    function decodeDynamicStore(bytes memory data) 
        internal pure returns (DynamicStore memory store) 
    {
        revert("TokenizationDecoders: Full protobuf decoding not supported in Solidity. Use off-chain tools.");
    }

    // ============ Simple Field Extractors ============
    // These provide basic extraction for simple fields without full decoding

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
     * @notice Recommended approach for decoding query responses
     * @dev For production contracts:
     *      1. Use getBalanceAmount() and getTotalSupply() which return uint256 directly
     *      2. For complex types, decode off-chain using TypeScript SDK
     *      3. Pass decoded data to contracts via function parameters if needed
     *      4. Use events to emit decoded data for off-chain indexing
     * 
     * Example pattern:
     * ```solidity
     * // Instead of decoding full collection:
     * bytes memory collectionBytes = precompile.getCollection(json);
     * 
     * // Use direct queries for specific fields:
     * uint256 balance = precompile.getBalanceAmount(json);
     * 
     * // Or emit raw bytes for off-chain processing:
     * emit CollectionQueried(collectionId, collectionBytes);
     * ```
     */
    function recommendedApproach() internal pure {
        // This function exists only for documentation purposes
        // See NatSpec comments above
        // Suppress unused function warning
        bytes memory unused = "";
        unused;
    }
}

