// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../types/GammTypes.sol";

/// @title IGammPrecompile
/// @notice Interface for the BitBadges gamm precompile
/// @dev Precompile address: 0x0000000000000000000000000000000000001002
///      All methods use JSON string parameters matching protobuf JSON format.
///      The caller address (sender) is automatically set from msg.sender.
///      Use helper libraries to construct JSON strings from Solidity types.
interface IGammPrecompile {
    // ============ Transactions ============
    // NOTE: All methods now use a JSON string parameter (msgJson) instead of individual parameters.
    // The JSON format matches the protobuf JSON serialization format.

    /// @notice Join a liquidity pool by providing tokens
    /// @param msgJson JSON string matching MsgJoinPool protobuf JSON format
    /// @return shareOutAmount The amount of pool shares received
    /// @return tokenIn The actual tokens provided (may be less than tokenInMaxs)
    function joinPool(string memory msgJson)
        external
        returns (uint256 shareOutAmount, GammTypes.Coin[] memory tokenIn);

    /// @notice Exit a liquidity pool by burning shares
    /// @param msgJson JSON string matching MsgExitPool protobuf JSON format
    /// @return tokenOut The tokens received from exiting the pool
    function exitPool(string memory msgJson)
        external
        returns (GammTypes.Coin[] memory tokenOut);

    /// @notice Swap tokens with exact input amount
    /// @param msgJson JSON string matching MsgSwapExactAmountIn protobuf JSON format
    /// @return tokenOutAmount The amount of output tokens received
    function swapExactAmountIn(string memory msgJson)
        external
        returns (uint256 tokenOutAmount);

    /// @notice Swap tokens with exact input amount and transfer via IBC
    /// @param msgJson JSON string matching MsgSwapExactAmountInWithIBCTransfer protobuf JSON format
    /// @return tokenOutAmount The amount of output tokens received
    function swapExactAmountInWithIBCTransfer(string memory msgJson)
        external
        returns (uint256 tokenOutAmount);

    // ============ Queries ============

    /// @notice Get pool data by ID
    /// @param msgJson JSON string matching QueryPoolRequest protobuf JSON format
    /// @return pool The pool data as protobuf-encoded bytes
    function getPool(string memory msgJson)
        external
        view
        returns (bytes memory pool);

    /// @notice Get all pools with pagination
    /// @param msgJson JSON string matching QueryPoolsRequest protobuf JSON format
    /// @return pools The pools data as protobuf-encoded bytes
    function getPools(string memory msgJson)
        external
        view
        returns (bytes memory pools);

    /// @notice Get pool type by ID
    /// @param msgJson JSON string matching QueryPoolTypeRequest protobuf JSON format
    /// @return poolType The pool type string
    function getPoolType(string memory msgJson)
        external
        view
        returns (string memory poolType);

    /// @notice Calculate shares for joining pool without swap
    /// @param msgJson JSON string matching QueryCalcJoinPoolNoSwapSharesRequest protobuf JSON format
    /// @return tokensOut The tokens that would be provided
    /// @return sharesOut The shares that would be received
    function calcJoinPoolNoSwapShares(string memory msgJson)
        external
        view
        returns (GammTypes.Coin[] memory tokensOut, uint256 sharesOut);

    /// @notice Calculate tokens received for exiting pool
    /// @param msgJson JSON string matching QueryCalcExitPoolCoinsFromSharesRequest protobuf JSON format
    /// @return tokensOut The tokens that would be received
    function calcExitPoolCoinsFromShares(string memory msgJson)
        external
        view
        returns (GammTypes.Coin[] memory tokensOut);

    /// @notice Calculate shares for joining pool (with swap)
    /// @param msgJson JSON string matching QueryCalcJoinPoolSharesRequest protobuf JSON format
    /// @return shareOutAmount The shares that would be received
    /// @return tokensOut The tokens that would be provided
    function calcJoinPoolShares(string memory msgJson)
        external
        view
        returns (uint256 shareOutAmount, GammTypes.Coin[] memory tokensOut);

    /// @notice Get pool parameters
    /// @param msgJson JSON string matching QueryPoolParamsRequest protobuf JSON format
    /// @return params The pool parameters as protobuf-encoded bytes
    function getPoolParams(string memory msgJson)
        external
        view
        returns (bytes memory params);

    /// @notice Get total shares for a pool
    /// @param msgJson JSON string matching QueryTotalSharesRequest protobuf JSON format
    /// @return totalShares The total shares as a Coin struct
    function getTotalShares(string memory msgJson)
        external
        view
        returns (GammTypes.Coin memory totalShares);

    /// @notice Get total liquidity for a pool
    /// @param msgJson JSON string matching QueryTotalLiquidityRequest protobuf JSON format
    /// @return liquidity The total liquidity as an array of Coin structs
    function getTotalLiquidity(string memory msgJson)
        external
        view
        returns (GammTypes.Coin[] memory liquidity);
}

