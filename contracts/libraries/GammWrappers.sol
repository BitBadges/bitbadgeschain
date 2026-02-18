// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/IGammPrecompile.sol";
import "../types/GammTypes.sol";
import "./GammJSONHelpers.sol";

/**
 * @title GammWrappers
 * @notice Typed wrapper functions for the gamm precompile
 * @dev Provides type-safe wrappers that accept structs instead of JSON strings.
 *      These wrappers internally convert structs to JSON using GammJSONHelpers.
 *      Use these for better compile-time type checking and reduced boilerplate.
 * 
 * Usage example:
 * ```solidity
 * import "./libraries/GammWrappers.sol";
 * 
 * IGammPrecompile precompile = IGammPrecompile(0x0000000000000000000000000000000000001002);
 * 
 * // Type-safe join pool
 * GammTypes.Coin[] memory tokenInMaxs = new GammTypes.Coin[](2);
 * tokenInMaxs[0] = GammTypes.Coin({denom: "uatom", amount: 1000000});
 * tokenInMaxs[1] = GammTypes.Coin({denom: "uosmo", amount: 2000000});
 * 
 * (uint256 shares, GammTypes.Coin[] memory tokens) = GammWrappers.joinPool(
 *     precompile,
 *     1,
 *     1000000,
 *     tokenInMaxs
 * );
 * ```
 */
library GammWrappers {
    // ============ Transaction Wrappers ============

    /**
     * @notice Join a liquidity pool with typed parameters
     * @param precompile The precompile interface instance
     * @param poolId The pool ID to join
     * @param shareOutAmount The desired amount of pool shares
     * @param tokenInMaxs Maximum amounts of tokens to provide
     * @return shareOutAmount The actual amount of pool shares received
     * @return tokenIn The actual tokens provided
     */
    function joinPool(
        IGammPrecompile precompile,
        uint64 poolId,
        uint256 shareOutAmount,
        GammTypes.Coin[] memory tokenInMaxs
    ) internal returns (uint256, GammTypes.Coin[] memory) {
        string memory tokenInMaxsJson = GammJSONHelpers.coinsToJson(tokenInMaxs);
        string memory json = GammJSONHelpers.joinPoolJSON(poolId, shareOutAmount, tokenInMaxsJson);
        return precompile.joinPool(json);
    }

    /**
     * @notice Exit a liquidity pool with typed parameters
     * @param precompile The precompile interface instance
     * @param poolId The pool ID to exit
     * @param shareInAmount The amount of pool shares to burn
     * @param tokenOutMins Minimum amounts of tokens to receive
     * @return tokenOut The tokens received from exiting the pool
     */
    function exitPool(
        IGammPrecompile precompile,
        uint64 poolId,
        uint256 shareInAmount,
        GammTypes.Coin[] memory tokenOutMins
    ) internal returns (GammTypes.Coin[] memory) {
        string memory tokenOutMinsJson = GammJSONHelpers.coinsToJson(tokenOutMins);
        string memory json = GammJSONHelpers.exitPoolJSON(poolId, shareInAmount, tokenOutMinsJson);
        return precompile.exitPool(json);
    }

    /**
     * @notice Swap tokens with exact input amount
     * @param precompile The precompile interface instance
     * @param routes Array of swap routes
     * @param tokenIn Input token and amount
     * @param tokenOutMinAmount Minimum output amount
     * @param affiliates Array of affiliate fee recipients (can be empty)
     * @return tokenOutAmount The amount of output tokens received
     */
    function swapExactAmountIn(
        IGammPrecompile precompile,
        GammTypes.SwapAmountInRoute[] memory routes,
        GammTypes.Coin memory tokenIn,
        uint256 tokenOutMinAmount,
        GammTypes.Affiliate[] memory affiliates
    ) internal returns (uint256) {
        string memory routesJson = GammJSONHelpers.swapRoutesToJson(routes);
        string memory tokenInJson = GammJSONHelpers.coinToJson(tokenIn);
        string memory affiliatesJson = GammJSONHelpers.affiliatesToJson(affiliates);
        string memory json = GammJSONHelpers.swapExactAmountInJSON(
            routesJson,
            tokenInJson,
            tokenOutMinAmount,
            affiliatesJson
        );
        return precompile.swapExactAmountIn(json);
    }

    /**
     * @notice Swap tokens with exact input amount and IBC transfer
     * @param precompile The precompile interface instance
     * @param routes Array of swap routes
     * @param tokenIn Input token and amount
     * @param tokenOutMinAmount Minimum output amount
     * @param ibcTransferInfo IBC transfer information
     * @param affiliates Array of affiliate fee recipients (can be empty)
     * @return tokenOutAmount The amount of output tokens received
     */
    function swapExactAmountInWithIBCTransfer(
        IGammPrecompile precompile,
        GammTypes.SwapAmountInRoute[] memory routes,
        GammTypes.Coin memory tokenIn,
        uint256 tokenOutMinAmount,
        GammTypes.IBCTransferInfo memory ibcTransferInfo,
        GammTypes.Affiliate[] memory affiliates
    ) internal returns (uint256) {
        string memory routesJson = GammJSONHelpers.swapRoutesToJson(routes);
        string memory tokenInJson = GammJSONHelpers.coinToJson(tokenIn);
        string memory ibcTransferInfoJson = GammJSONHelpers.ibcTransferInfoToJson(ibcTransferInfo);
        string memory affiliatesJson = GammJSONHelpers.affiliatesToJson(affiliates);
        string memory json = GammJSONHelpers.swapExactAmountInWithIBCTransferJSON(
            routesJson,
            tokenInJson,
            tokenOutMinAmount,
            ibcTransferInfoJson,
            affiliatesJson
        );
        return precompile.swapExactAmountInWithIBCTransfer(json);
    }

    // ============ Query Wrappers ============

    /**
     * @notice Get pool data by ID
     * @param precompile The precompile interface instance
     * @param poolId The pool ID
     * @return pool The pool data as protobuf-encoded bytes
     */
    function getPool(
        IGammPrecompile precompile,
        uint64 poolId
    ) internal view returns (bytes memory) {
        string memory json = GammJSONHelpers.getPoolJSON(poolId);
        return precompile.getPool(json);
    }

    /**
     * @notice Get all pools with pagination
     * @param precompile The precompile interface instance
     * @param paginationJson JSON string for pagination (empty string for default)
     * @return pools The pools data as protobuf-encoded bytes
     */
    function getPools(
        IGammPrecompile precompile,
        string memory paginationJson
    ) internal view returns (bytes memory) {
        string memory json = GammJSONHelpers.getPoolsJSON(paginationJson);
        return precompile.getPools(json);
    }

    /**
     * @notice Get pool type by ID
     * @param precompile The precompile interface instance
     * @param poolId The pool ID
     * @return poolType The pool type string
     */
    function getPoolType(
        IGammPrecompile precompile,
        uint64 poolId
    ) internal view returns (string memory) {
        string memory json = GammJSONHelpers.getPoolTypeJSON(poolId);
        return precompile.getPoolType(json);
    }

    /**
     * @notice Calculate shares for joining pool without swap
     * @param precompile The precompile interface instance
     * @param poolId The pool ID
     * @param tokenInMaxs Maximum amounts of tokens to provide
     * @return tokensOut The tokens that would be provided
     * @return sharesOut The shares that would be received
     */
    function calcJoinPoolNoSwapShares(
        IGammPrecompile precompile,
        uint64 poolId,
        GammTypes.Coin[] memory tokenInMaxs
    ) internal view returns (GammTypes.Coin[] memory, uint256) {
        string memory tokenInMaxsJson = GammJSONHelpers.coinsToJson(tokenInMaxs);
        string memory json = GammJSONHelpers.calcJoinPoolNoSwapSharesJSON(poolId, tokenInMaxsJson);
        return precompile.calcJoinPoolNoSwapShares(json);
    }

    /**
     * @notice Calculate tokens received for exiting pool
     * @param precompile The precompile interface instance
     * @param poolId The pool ID
     * @param shareInAmount The amount of pool shares to burn
     * @return tokensOut The tokens that would be received
     */
    function calcExitPoolCoinsFromShares(
        IGammPrecompile precompile,
        uint64 poolId,
        uint256 shareInAmount
    ) internal view returns (GammTypes.Coin[] memory) {
        string memory json = GammJSONHelpers.calcExitPoolCoinsFromSharesJSON(poolId, shareInAmount);
        return precompile.calcExitPoolCoinsFromShares(json);
    }

    /**
     * @notice Calculate shares for joining pool (with swap)
     * @param precompile The precompile interface instance
     * @param poolId The pool ID
     * @param tokenInMaxs Maximum amounts of tokens to provide
     * @return shareOutAmount The shares that would be received
     * @return tokensOut The tokens that would be provided
     */
    function calcJoinPoolShares(
        IGammPrecompile precompile,
        uint64 poolId,
        GammTypes.Coin[] memory tokenInMaxs
    ) internal view returns (uint256, GammTypes.Coin[] memory) {
        string memory tokenInMaxsJson = GammJSONHelpers.coinsToJson(tokenInMaxs);
        string memory json = GammJSONHelpers.calcJoinPoolSharesJSON(poolId, tokenInMaxsJson);
        return precompile.calcJoinPoolShares(json);
    }

    /**
     * @notice Get pool parameters
     * @param precompile The precompile interface instance
     * @param poolId The pool ID
     * @return params The pool parameters as protobuf-encoded bytes
     */
    function getPoolParams(
        IGammPrecompile precompile,
        uint64 poolId
    ) internal view returns (bytes memory) {
        string memory json = GammJSONHelpers.getPoolParamsJSON(poolId);
        return precompile.getPoolParams(json);
    }

    /**
     * @notice Get total shares for a pool
     * @param precompile The precompile interface instance
     * @param poolId The pool ID
     * @return totalShares The total shares as a Coin struct
     */
    function getTotalShares(
        IGammPrecompile precompile,
        uint64 poolId
    ) internal view returns (GammTypes.Coin memory) {
        string memory json = GammJSONHelpers.getTotalSharesJSON(poolId);
        return precompile.getTotalShares(json);
    }

    /**
     * @notice Get total liquidity for a pool
     * @param precompile The precompile interface instance
     * @param poolId The pool ID
     * @return liquidity The total liquidity as an array of Coin structs
     */
    function getTotalLiquidity(
        IGammPrecompile precompile,
        uint64 poolId
    ) internal view returns (GammTypes.Coin[] memory) {
        string memory json = GammJSONHelpers.getTotalLiquidityJSON(poolId);
        return precompile.getTotalLiquidity(json);
    }
}

















