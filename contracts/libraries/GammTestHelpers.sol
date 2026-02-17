// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../types/GammTypes.sol";
import "./GammHelpers.sol";

/**
 * @title GammTestHelpers
 * @notice Testing utilities for gamm contracts
 * @dev Provides mock data generators, test fixtures, and assertion helpers for writing tests.
 *      These helpers make it easier to write comprehensive tests for contracts using the gamm precompile.
 * 
 * Usage example:
 * ```solidity
 * import "./libraries/GammTestHelpers.sol";
 * 
 * function testJoinPool() public {
 *     uint64 poolId = 1;
 *     uint256 shareOutAmount = 1000000;
 *     
 *     // Generate test data
 *     GammTypes.Coin[] memory tokenInMaxs = 
 *         GammTestHelpers.generateCoins(["uatom", "uosmo"], [1000000, 2000000]);
 *     
 *     // Perform join pool
 *     // ...
 * }
 * ```
 */
library GammTestHelpers {
    // ============ Mock Data Generators ============

    /**
     * @notice Generate an array of coins from parallel arrays
     * @param denoms Array of denominations
     * @param amounts Array of amounts
     * @return coins Array of Coin structs
     */
    function generateCoins(
        string[] memory denoms,
        uint256[] memory amounts
    ) internal pure returns (GammTypes.Coin[] memory coins) {
        require(denoms.length == amounts.length, "GammTestHelpers: arrays must have same length");
        coins = new GammTypes.Coin[](denoms.length);
        for (uint256 i = 0; i < denoms.length; i++) {
            coins[i] = GammHelpers.createCoin(denoms[i], amounts[i]);
        }
        return coins;
    }

    /**
     * @notice Generate a single coin
     * @param denom The denomination
     * @param amount The amount
     * @return coin A Coin struct
     */
    function generateCoin(string memory denom, uint256 amount) internal pure returns (GammTypes.Coin memory) {
        return GammHelpers.createCoin(denom, amount);
    }

    /**
     * @notice Generate swap routes for a multi-hop swap
     * @param poolIds Array of pool IDs
     * @param tokenOutDenoms Array of output token denominations
     * @return routes Array of SwapAmountInRoute structs
     */
    function generateSwapRoutes(
        uint64[] memory poolIds,
        string[] memory tokenOutDenoms
    ) internal pure returns (GammTypes.SwapAmountInRoute[] memory routes) {
        require(poolIds.length == tokenOutDenoms.length, "GammTestHelpers: arrays must have same length");
        routes = new GammTypes.SwapAmountInRoute[](poolIds.length);
        for (uint256 i = 0; i < poolIds.length; i++) {
            routes[i] = GammHelpers.createSwapRoute(poolIds[i], tokenOutDenoms[i]);
        }
        return routes;
    }

    /**
     * @notice Generate a single-hop swap route
     * @param poolId The pool ID
     * @param tokenOutDenom The output token denomination
     * @return routes Array with single route
     */
    function generateSingleHopRoute(
        uint64 poolId,
        string memory tokenOutDenom
    ) internal pure returns (GammTypes.SwapAmountInRoute[] memory) {
        return GammHelpers.createSingleHopRoute(poolId, tokenOutDenom);
    }

    /**
     * @notice Generate a two-hop swap route
     * @param poolId1 First pool ID
     * @param intermediateDenom Intermediate token denomination
     * @param poolId2 Second pool ID
     * @param tokenOutDenom Final output token denomination
     * @return routes Array with two routes
     */
    function generateTwoHopRoute(
        uint64 poolId1,
        string memory intermediateDenom,
        uint64 poolId2,
        string memory tokenOutDenom
    ) internal pure returns (GammTypes.SwapAmountInRoute[] memory) {
        return GammHelpers.createTwoHopRoute(poolId1, intermediateDenom, poolId2, tokenOutDenom);
    }

    /**
     * @notice Generate affiliates array
     * @param addresses Array of affiliate addresses
     * @param basisPointsFees Array of basis points fees
     * @return affiliates Array of Affiliate structs
     */
    function generateAffiliates(
        address[] memory addresses,
        uint256[] memory basisPointsFees
    ) internal pure returns (GammTypes.Affiliate[] memory affiliates) {
        require(addresses.length == basisPointsFees.length, "GammTestHelpers: arrays must have same length");
        affiliates = new GammTypes.Affiliate[](addresses.length);
        for (uint256 i = 0; i < addresses.length; i++) {
            affiliates[i] = GammHelpers.createAffiliate(addresses[i], basisPointsFees[i]);
        }
        return affiliates;
    }

    /**
     * @notice Generate empty affiliates array
     * @return affiliates Empty Affiliate array
     */
    function generateEmptyAffiliates() internal pure returns (GammTypes.Affiliate[] memory) {
        return GammHelpers.createEmptyAffiliates();
    }

    /**
     * @notice Generate IBC transfer info
     * @param sourceChannel The IBC source channel
     * @param receiver The receiver address
     * @param memo Optional memo
     * @param timeoutTimestamp The timeout timestamp (0 for default)
     * @return info IBCTransferInfo struct
     */
    function generateIBCTransferInfo(
        string memory sourceChannel,
        string memory receiver,
        string memory memo,
        uint64 timeoutTimestamp
    ) internal pure returns (GammTypes.IBCTransferInfo memory) {
        return GammHelpers.createIBCTransferInfo(sourceChannel, receiver, memo, timeoutTimestamp);
    }
}


















