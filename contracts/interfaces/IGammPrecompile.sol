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

    /// @notice Create a balancer liquidity pool
    /// @param msgJson JSON string matching MsgCreateBalancerPool protobuf JSON format
    /// @return poolId The ID of the newly created pool
    function createPool(string memory msgJson)
        external
        returns (uint256 poolId);

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

    /// @notice Get total liquidity across all pools
    /// @param msgJson JSON string matching QueryTotalLiquidityRequest protobuf JSON format (e.g. "{}")
    /// @return liquidity The total liquidity as an array of Coin structs
    function getTotalLiquidity(string memory msgJson)
        external
        view
        returns (GammTypes.Coin[] memory liquidity);

    /// @notice Get spot price for a pool asset pair (quote_asset_denom per 1 base_asset_denom, 18-dec string)
    /// @param msgJson JSON string matching QuerySpotPriceRequest ({"pool_id": 1, "base_asset_denom": "...", "quote_asset_denom": "..."})
    /// @return spotPrice The 18-decimal spot price string
    function spotPrice(string memory msgJson)
        external
        view
        returns (string memory spotPrice);

    /// @notice Get cumulative swap volume recorded for a specific pool
    /// @param msgJson JSON string matching poolmanager TotalVolumeForPoolRequest ({"pool_id": 1})
    /// @return volume The cumulative pool volume as an array of Coin structs
    function getTotalVolumeForPool(string memory msgJson)
        external
        view
        returns (GammTypes.Coin[] memory volume);

    /// @notice Get cumulative taker fees tracked by x/poolmanager
    /// @param msgJson JSON string ("{}")
    /// @return takerFeesToStakers Cumulative taker fees distributed to stakers
    /// @return takerFeesToCommunityPool Cumulative taker fees sent to the community pool
    /// @return heightAccountingStartsFrom Block height when taker fee accounting started
    function getTakerFeesTracker(string memory msgJson)
        external
        view
        returns (
            GammTypes.Coin[] memory takerFeesToStakers,
            GammTypes.Coin[] memory takerFeesToCommunityPool,
            int64 heightAccountingStartsFrom
        );

    /// @notice Estimate output amount for a swap given an exact input coin and routes
    /// @param msgJson JSON string matching QuerySwapExactAmountInRequest ({"token_in": "100ubadge", "routes": [{"pool_id": 1, "token_out_denom": "..."}]})
    /// @return tokenOutAmount The estimated output token amount
    function estimateSwapExactAmountIn(string memory msgJson)
        external
        view
        returns (uint256 tokenOutAmount);

    /// @notice Estimate input amount required for a swap given an exact output coin and routes
    /// @param msgJson JSON string matching QuerySwapExactAmountOutRequest ({"token_out": "100uusdc", "routes": [{"pool_id": 1, "token_in_denom": "..."}]})
    /// @return tokenInAmount The estimated input token amount
    function estimateSwapExactAmountOut(string memory msgJson)
        external
        view
        returns (uint256 tokenInAmount);

    /// @notice Get total liquidity for a specific pool
    /// @param msgJson JSON string matching QueryTotalPoolLiquidityRequest ({"pool_id": 1})
    /// @return liquidity The pool liquidity as an array of Coin structs
    function getTotalPoolLiquidity(string memory msgJson)
        external
        view
        returns (GammTypes.Coin[] memory liquidity);

    /// @notice Get the configured taker fee for an ordered trading pair (18-dec string)
    /// @param msgJson JSON string matching poolmanager TradingPairTakerFeeRequest ({"denom_0": "...", "denom_1": "..."})
    /// @return takerFee The taker fee decimal string
    function getTradingPairTakerFee(string memory msgJson)
        external
        view
        returns (string memory takerFee);

    /// @notice Get the total number of pools created
    /// @param msgJson JSON string ("{}")
    /// @return numPools Total number of pools
    function getNumPools(string memory msgJson)
        external
        view
        returns (uint64 numPools);
}

