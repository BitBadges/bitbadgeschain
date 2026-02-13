// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * @title GammDecoders
 * @notice Placeholder decoders for protobuf-encoded query responses
 * @dev Full protobuf decoding in Solidity is complex and not practical.
 *      These functions provide informative errors directing users to decode
 *      protobuf responses off-chain or use the precompile's typed return values.
 * 
 * Note: The precompile already returns typed values for most queries (Coin[], Coin, uint256, string).
 *       Use those return values directly instead of decoding bytes.
 */
library GammDecoders {
    /**
     * @notice Placeholder for decoding pool bytes
     * @dev Protobuf decoding in Solidity is not practical.
     *      Use the precompile's typed return values or decode off-chain.
     * @return This function always reverts with an informative error
     */
    function decodePool(bytes memory /* poolBytes */) internal pure returns (bytes memory) {
        // This is a placeholder - protobuf decoding in Solidity is complex
        // The precompile returns bytes for getPool, but you should decode this off-chain
        // or use a library that supports protobuf decoding
        revert("GammDecoders: Protobuf decoding not supported in Solidity. Decode off-chain or use precompile's typed return values.");
    }

    /**
     * @notice Placeholder for decoding pools bytes
     * @dev Protobuf decoding in Solidity is not practical.
     *      Use the precompile's typed return values or decode off-chain.
     * @return This function always reverts with an informative error
     */
    function decodePools(bytes memory /* poolsBytes */) internal pure returns (bytes memory) {
        revert("GammDecoders: Protobuf decoding not supported in Solidity. Decode off-chain or use precompile's typed return values.");
    }

    /**
     * @notice Placeholder for decoding pool params bytes
     * @dev Protobuf decoding in Solidity is not practical.
     *      Use the precompile's typed return values or decode off-chain.
     * @return This function always reverts with an informative error
     */
    function decodePoolParams(bytes memory /* paramsBytes */) internal pure returns (bytes memory) {
        revert("GammDecoders: Protobuf decoding not supported in Solidity. Decode off-chain or use precompile's typed return values.");
    }
}

