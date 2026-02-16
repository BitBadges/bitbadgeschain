// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../types/GammTypes.sol";

/**
 * @title GammHelpers
 * @notice Convenience functions for creating common gamm types
 * @dev Provides helper functions to create Coin, SwapAmountInRoute, Affiliate, and other types
 *      with sensible defaults and common patterns.
 */
library GammHelpers {
    /**
     * @notice Create a Coin struct
     * @param denom The denomination (e.g., "uatom", "uosmo")
     * @param amount The amount
     * @return coin The Coin struct
     */
    function createCoin(string memory denom, uint256 amount) internal pure returns (GammTypes.Coin memory) {
        return GammTypes.Coin({denom: denom, amount: amount});
    }

    /**
     * @notice Create a SwapAmountInRoute struct
     * @param poolId The pool ID
     * @param tokenOutDenom The output token denomination
     * @return route The SwapAmountInRoute struct
     */
    function createSwapRoute(uint64 poolId, string memory tokenOutDenom) internal pure returns (GammTypes.SwapAmountInRoute memory) {
        return GammTypes.SwapAmountInRoute({poolId: poolId, tokenOutDenom: tokenOutDenom});
    }

    /**
     * @notice Create an Affiliate struct
     * @param address_ The affiliate recipient address
     * @param basisPointsFee The fee in basis points (0-10000, where 10000 = 100%)
     * @return affiliate The Affiliate struct
     */
    function createAffiliate(address address_, uint256 basisPointsFee) internal pure returns (GammTypes.Affiliate memory) {
        return GammTypes.Affiliate({address_: address_, basisPointsFee: basisPointsFee});
    }

    /**
     * @notice Create an IBCTransferInfo struct
     * @param sourceChannel The IBC source channel
     * @param receiver The receiver address
     * @param memo Optional memo
     * @param timeoutTimestamp The timeout timestamp (0 for default)
     * @return info The IBCTransferInfo struct
     */
    function createIBCTransferInfo(
        string memory sourceChannel,
        string memory receiver,
        string memory memo,
        uint64 timeoutTimestamp
    ) internal pure returns (GammTypes.IBCTransferInfo memory) {
        return GammTypes.IBCTransferInfo({
            sourceChannel: sourceChannel,
            receiver: receiver,
            memo: memo,
            timeoutTimestamp: timeoutTimestamp
        });
    }

    /**
     * @notice Create an empty affiliates array
     * @return affiliates Empty Affiliate array
     */
    function createEmptyAffiliates() internal pure returns (GammTypes.Affiliate[] memory) {
        return new GammTypes.Affiliate[](0);
    }

    /**
     * @notice Create a single-hop swap route
     * @param poolId The pool ID
     * @param tokenOutDenom The output token denomination
     * @return routes Array with single route
     */
    function createSingleHopRoute(uint64 poolId, string memory tokenOutDenom) internal pure returns (GammTypes.SwapAmountInRoute[] memory) {
        GammTypes.SwapAmountInRoute[] memory routes = new GammTypes.SwapAmountInRoute[](1);
        routes[0] = createSwapRoute(poolId, tokenOutDenom);
        return routes;
    }

    /**
     * @notice Create a two-hop swap route
     * @param poolId1 First pool ID
     * @param intermediateDenom Intermediate token denomination
     * @param poolId2 Second pool ID
     * @param tokenOutDenom Final output token denomination
     * @return routes Array with two routes
     */
    function createTwoHopRoute(
        uint64 poolId1,
        string memory intermediateDenom,
        uint64 poolId2,
        string memory tokenOutDenom
    ) internal pure returns (GammTypes.SwapAmountInRoute[] memory) {
        GammTypes.SwapAmountInRoute[] memory routes = new GammTypes.SwapAmountInRoute[](2);
        routes[0] = createSwapRoute(poolId1, intermediateDenom);
        routes[1] = createSwapRoute(poolId2, tokenOutDenom);
        return routes;
    }
}














