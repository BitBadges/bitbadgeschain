// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../types/GammTypes.sol";
import "./GammJSONHelpers.sol";
import "./GammHelpers.sol";

/**
 * @title GammBuilders
 * @notice Builder pattern utilities for constructing complex gamm operations
 * @dev Provides fluent builder APIs for complex operations like join pool, exit pool, and swaps.
 *      These builders make it easier to construct complex JSON structures step by step.
 * 
 * Usage example:
 * ```solidity
 * import "./libraries/GammBuilders.sol";
 * 
 * // Build a join pool request
 * JoinPoolBuilder memory builder = GammBuilders.newJoinPool(poolId, shareOutAmount);
 * builder = builder.addTokenInMax("uatom", 1000000);
 * builder = builder.addTokenInMax("uosmo", 2000000);
 * string memory json = builder.build();
 * (uint256 shares, Coin[] memory tokens) = precompile.joinPool(json);
 * ```
 */
library GammBuilders {
    // ============ Join Pool Builder ============

    /**
     * @notice Builder struct for joining pools
     */
    struct JoinPoolBuilder {
        uint64 poolId;
        uint256 shareOutAmount;
        GammTypes.Coin[] tokenInMaxs;
    }

    /**
     * @notice Create a new join pool builder
     * @param poolId The pool ID to join
     * @param shareOutAmount The desired amount of pool shares
     * @return builder A new JoinPoolBuilder
     */
    function newJoinPool(uint64 poolId, uint256 shareOutAmount) internal pure returns (JoinPoolBuilder memory builder) {
        builder.poolId = poolId;
        builder.shareOutAmount = shareOutAmount;
        builder.tokenInMaxs = new GammTypes.Coin[](0);
        return builder;
    }

    /**
     * @notice Add a token input maximum
     * @param builder The builder instance
     * @param denom The token denomination
     * @param amount The maximum amount
     * @return The builder with token added
     */
    function addTokenInMax(
        JoinPoolBuilder memory builder,
        string memory denom,
        uint256 amount
    ) internal pure returns (JoinPoolBuilder memory) {
        GammTypes.Coin[] memory newTokens = new GammTypes.Coin[](builder.tokenInMaxs.length + 1);
        for (uint256 i = 0; i < builder.tokenInMaxs.length; i++) {
            newTokens[i] = builder.tokenInMaxs[i];
        }
        newTokens[builder.tokenInMaxs.length] = GammHelpers.createCoin(denom, amount);
        builder.tokenInMaxs = newTokens;
        return builder;
    }

    /**
     * @notice Build the JSON string for join pool
     * @param builder The builder instance
     * @return json The JSON string ready for precompile call
     */
    function build(JoinPoolBuilder memory builder) internal pure returns (string memory) {
        string memory tokenInMaxsJson = GammJSONHelpers.coinsToJson(builder.tokenInMaxs);
        return GammJSONHelpers.joinPoolJSON(builder.poolId, builder.shareOutAmount, tokenInMaxsJson);
    }

    // ============ Exit Pool Builder ============

    /**
     * @notice Builder struct for exiting pools
     */
    struct ExitPoolBuilder {
        uint64 poolId;
        uint256 shareInAmount;
        GammTypes.Coin[] tokenOutMins;
    }

    /**
     * @notice Create a new exit pool builder
     * @param poolId The pool ID to exit
     * @param shareInAmount The amount of pool shares to burn
     * @return builder A new ExitPoolBuilder
     */
    function newExitPool(uint64 poolId, uint256 shareInAmount) internal pure returns (ExitPoolBuilder memory builder) {
        builder.poolId = poolId;
        builder.shareInAmount = shareInAmount;
        builder.tokenOutMins = new GammTypes.Coin[](0);
        return builder;
    }

    /**
     * @notice Add a token output minimum
     * @param builder The builder instance
     * @param denom The token denomination
     * @param amount The minimum amount
     * @return The builder with token added
     */
    function addTokenOutMin(
        ExitPoolBuilder memory builder,
        string memory denom,
        uint256 amount
    ) internal pure returns (ExitPoolBuilder memory) {
        GammTypes.Coin[] memory newTokens = new GammTypes.Coin[](builder.tokenOutMins.length + 1);
        for (uint256 i = 0; i < builder.tokenOutMins.length; i++) {
            newTokens[i] = builder.tokenOutMins[i];
        }
        newTokens[builder.tokenOutMins.length] = GammHelpers.createCoin(denom, amount);
        builder.tokenOutMins = newTokens;
        return builder;
    }

    /**
     * @notice Build the JSON string for exit pool
     * @param builder The builder instance
     * @return json The JSON string ready for precompile call
     */
    function build(ExitPoolBuilder memory builder) internal pure returns (string memory) {
        string memory tokenOutMinsJson = GammJSONHelpers.coinsToJson(builder.tokenOutMins);
        return GammJSONHelpers.exitPoolJSON(builder.poolId, builder.shareInAmount, tokenOutMinsJson);
    }

    // ============ Swap Builder ============

    /**
     * @notice Builder struct for swaps
     */
    struct SwapBuilder {
        GammTypes.SwapAmountInRoute[] routes;
        GammTypes.Coin tokenIn;
        uint256 tokenOutMinAmount;
        GammTypes.Affiliate[] affiliates;
        bool hasTokenIn;
        bool hasTokenOutMinAmount;
    }

    /**
     * @notice Create a new swap builder
     * @return builder A new SwapBuilder
     */
    function newSwap() internal pure returns (SwapBuilder memory builder) {
        builder.routes = new GammTypes.SwapAmountInRoute[](0);
        builder.affiliates = new GammTypes.Affiliate[](0);
        return builder;
    }

    /**
     * @notice Add a swap route
     * @param builder The builder instance
     * @param poolId The pool ID
     * @param tokenOutDenom The output token denomination
     * @return The builder with route added
     */
    function addRoute(
        SwapBuilder memory builder,
        uint64 poolId,
        string memory tokenOutDenom
    ) internal pure returns (SwapBuilder memory) {
        GammTypes.SwapAmountInRoute[] memory newRoutes = new GammTypes.SwapAmountInRoute[](builder.routes.length + 1);
        for (uint256 i = 0; i < builder.routes.length; i++) {
            newRoutes[i] = builder.routes[i];
        }
        newRoutes[builder.routes.length] = GammHelpers.createSwapRoute(poolId, tokenOutDenom);
        builder.routes = newRoutes;
        return builder;
    }

    /**
     * @notice Set the input token
     * @param builder The builder instance
     * @param denom The input token denomination
     * @param amount The input amount
     * @return The builder with tokenIn set
     */
    function withTokenIn(
        SwapBuilder memory builder,
        string memory denom,
        uint256 amount
    ) internal pure returns (SwapBuilder memory) {
        builder.tokenIn = GammHelpers.createCoin(denom, amount);
        builder.hasTokenIn = true;
        return builder;
    }

    /**
     * @notice Set the minimum output amount
     * @param builder The builder instance
     * @param amount The minimum output amount
     * @return The builder with tokenOutMinAmount set
     */
    function withTokenOutMinAmount(
        SwapBuilder memory builder,
        uint256 amount
    ) internal pure returns (SwapBuilder memory) {
        builder.tokenOutMinAmount = amount;
        builder.hasTokenOutMinAmount = true;
        return builder;
    }

    /**
     * @notice Add an affiliate
     * @param builder The builder instance
     * @param address_ The affiliate address
     * @param basisPointsFee The fee in basis points
     * @return The builder with affiliate added
     */
    function addAffiliate(
        SwapBuilder memory builder,
        address address_,
        uint256 basisPointsFee
    ) internal pure returns (SwapBuilder memory) {
        GammTypes.Affiliate[] memory newAffiliates = new GammTypes.Affiliate[](builder.affiliates.length + 1);
        for (uint256 i = 0; i < builder.affiliates.length; i++) {
            newAffiliates[i] = builder.affiliates[i];
        }
        newAffiliates[builder.affiliates.length] = GammHelpers.createAffiliate(address_, basisPointsFee);
        builder.affiliates = newAffiliates;
        return builder;
    }

    /**
     * @notice Build the JSON string for swap
     * @param builder The builder instance
     * @return json The JSON string ready for precompile call
     */
    function build(SwapBuilder memory builder) internal pure returns (string memory) {
        require(builder.hasTokenIn, "GammBuilders: tokenIn must be set");
        require(builder.hasTokenOutMinAmount, "GammBuilders: tokenOutMinAmount must be set");
        require(builder.routes.length > 0, "GammBuilders: at least one route required");

        string memory routesJson = GammJSONHelpers.swapRoutesToJson(builder.routes);
        string memory tokenInJson = GammJSONHelpers.coinToJson(builder.tokenIn);
        string memory affiliatesJson = GammJSONHelpers.affiliatesToJson(builder.affiliates);
        return GammJSONHelpers.swapExactAmountInJSON(
            routesJson,
            tokenInJson,
            builder.tokenOutMinAmount,
            affiliatesJson
        );
    }
}







