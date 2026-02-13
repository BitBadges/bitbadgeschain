// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * @title GammTypes
 * @notice Comprehensive type registry for BitBadges gamm module
 * @dev This file contains all Solidity struct definitions that mirror proto message types
 *      from the gamm module. All types are 1:1 mappings from proto definitions.
 *      NOTE: Sender fields are NOT included in input structs - they are always msg.sender
 */

library GammTypes {
    // ============================================================================
    // Core Types
    // ============================================================================

    /**
     * @notice Coin represents a token with denomination and amount
     * @dev Used for token inputs, outputs, and liquidity amounts
     */
    struct Coin {
        string denom;
        uint256 amount;
    }

    /**
     * @notice SwapAmountInRoute represents a single hop in a swap route
     * @dev Used for multi-hop swaps through different pools
     */
    struct SwapAmountInRoute {
        uint64 poolId;
        string tokenOutDenom;
    }

    /**
     * @notice Affiliate represents a fee recipient for swaps
     * @dev Used to specify affiliate fees in basis points (1/10000, e.g., 100 = 1%)
     */
    struct Affiliate {
        address address_; // Using address_ to avoid Solidity keyword conflict
        uint256 basisPointsFee; // Basis points (0-10000, where 10000 = 100%)
    }

    /**
     * @notice IBCTransferInfo represents IBC transfer information
     * @dev Used for swaps that include IBC transfers
     */
    struct IBCTransferInfo {
        string sourceChannel;
        string receiver;
        string memo;
        uint64 timeoutTimestamp;
    }

    // ============================================================================
    // Message Types (for JSON construction)
    // ============================================================================

    /**
     * @notice MsgJoinPool represents a join pool message
     * @dev Sender is automatically set from msg.sender
     */
    struct MsgJoinPool {
        uint64 poolId;
        string shareOutAmount; // String representation of uint256
        Coin[] tokenInMaxs;
    }

    /**
     * @notice MsgExitPool represents an exit pool message
     * @dev Sender is automatically set from msg.sender
     */
    struct MsgExitPool {
        uint64 poolId;
        string shareInAmount; // String representation of uint256
        Coin[] tokenOutMins;
    }

    /**
     * @notice MsgSwapExactAmountIn represents a swap with exact input amount
     * @dev Sender is automatically set from msg.sender
     */
    struct MsgSwapExactAmountIn {
        SwapAmountInRoute[] routes;
        Coin tokenIn;
        string tokenOutMinAmount; // String representation of uint256
        Affiliate[] affiliates;
    }

    /**
     * @notice MsgSwapExactAmountInWithIBCTransfer represents a swap with IBC transfer
     * @dev Sender is automatically set from msg.sender
     */
    struct MsgSwapExactAmountInWithIBCTransfer {
        SwapAmountInRoute[] routes;
        Coin tokenIn;
        string tokenOutMinAmount; // String representation of uint256
        IBCTransferInfo ibcTransferInfo;
        Affiliate[] affiliates;
    }

    // ============================================================================
    // Query Types (for JSON construction)
    // ============================================================================

    /**
     * @notice QueryPoolRequest represents a request to get pool by ID
     */
    struct QueryPoolRequest {
        uint64 poolId;
    }

    /**
     * @notice QueryPoolsRequest represents a request to get all pools
     * @dev Pagination is optional - handled as JSON strings in helpers
     *      This struct is not used directly - use getPoolsJSON() with pagination JSON
     */
    // QueryPoolsRequest is not defined as a struct since pagination is optional
    // Use GammJSONHelpers.getPoolsJSON() with pagination JSON string

    /**
     * @notice QueryPoolTypeRequest represents a request to get pool type
     */
    struct QueryPoolTypeRequest {
        uint64 poolId;
    }

    /**
     * @notice QueryCalcJoinPoolNoSwapSharesRequest represents a request to calculate shares
     */
    struct QueryCalcJoinPoolNoSwapSharesRequest {
        uint64 poolId;
        Coin[] tokenInMaxs;
    }

    /**
     * @notice QueryCalcExitPoolCoinsFromSharesRequest represents a request to calculate exit tokens
     */
    struct QueryCalcExitPoolCoinsFromSharesRequest {
        uint64 poolId;
        string shareInAmount; // String representation of uint256
    }

    /**
     * @notice QueryCalcJoinPoolSharesRequest represents a request to calculate join shares (with swap)
     */
    struct QueryCalcJoinPoolSharesRequest {
        uint64 poolId;
        Coin[] tokenInMaxs;
    }

    /**
     * @notice QueryPoolParamsRequest represents a request to get pool parameters
     */
    struct QueryPoolParamsRequest {
        uint64 poolId;
    }

    /**
     * @notice QueryTotalSharesRequest represents a request to get total shares
     */
    struct QueryTotalSharesRequest {
        uint64 poolId;
    }

    /**
     * @notice QueryTotalLiquidityRequest represents a request to get total liquidity
     */
    struct QueryTotalLiquidityRequest {
        uint64 poolId;
    }
}

