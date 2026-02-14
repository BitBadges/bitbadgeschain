// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/IGammPrecompile.sol";
import "../types/GammTypes.sol";
import "../libraries/GammWrappers.sol";
import "../libraries/GammBuilders.sol";
import "../libraries/GammHelpers.sol";
import "../libraries/GammJSONHelpers.sol";
import "../libraries/GammErrors.sol";

/**
 * @title GammHelperLibrariesTestContract
 * @notice Comprehensive test contract for all gamm helper libraries
 * @dev Tests all helper libraries E2E to verify:
 *      - Proper data population
 *      - Correct JSON construction
 *      - Correct underlying message calls
 *      - Type safety and error handling
 */
contract GammHelperLibrariesTestContract {
    IGammPrecompile constant GAMM = 
        IGammPrecompile(0x0000000000000000000000000000000000001002);
    
    // ============ Events for Test Verification ============
    
    event TestResult(
        string testName,
        bool success,
        bytes returnData
    );
    
    event JSONConstructed(
        string methodName,
        string json
    );
    
    // ============ GammWrappers Tests ============
    
    /**
     * @notice Test joinPool wrapper
     */
    function testJoinPoolWrapper(
        uint64 poolId,
        uint256 shareOutAmount,
        GammTypes.Coin[] memory tokenInMaxs
    ) external returns (uint256, GammTypes.Coin[] memory) {
        (uint256 shares, GammTypes.Coin[] memory tokens) = GammWrappers.joinPool(
            GAMM,
            poolId,
            shareOutAmount,
            tokenInMaxs
        );
        emit TestResult("joinPoolWrapper", true, abi.encode(shares, tokens));
        return (shares, tokens);
    }
    
    /**
     * @notice Test exitPool wrapper
     */
    function testExitPoolWrapper(
        uint64 poolId,
        uint256 shareInAmount,
        GammTypes.Coin[] memory tokenOutMins
    ) external returns (GammTypes.Coin[] memory) {
        GammTypes.Coin[] memory tokens = GammWrappers.exitPool(
            GAMM,
            poolId,
            shareInAmount,
            tokenOutMins
        );
        emit TestResult("exitPoolWrapper", true, abi.encode(tokens));
        return tokens;
    }
    
    /**
     * @notice Test swapExactAmountIn wrapper
     */
    function testSwapExactAmountInWrapper(
        GammTypes.SwapAmountInRoute[] memory routes,
        GammTypes.Coin memory tokenIn,
        uint256 tokenOutMinAmount,
        GammTypes.Affiliate[] memory affiliates
    ) external returns (uint256) {
        uint256 tokenOutAmount = GammWrappers.swapExactAmountIn(
            GAMM,
            routes,
            tokenIn,
            tokenOutMinAmount,
            affiliates
        );
        emit TestResult("swapExactAmountInWrapper", true, abi.encode(tokenOutAmount));
        return tokenOutAmount;
    }
    
    /**
     * @notice Test getPool wrapper
     */
    function testGetPoolWrapper(
        uint64 poolId
    ) external returns (bytes memory) {
        bytes memory pool = GammWrappers.getPool(GAMM, poolId);
        emit TestResult("getPoolWrapper", true, pool);
        return pool;
    }
    
    /**
     * @notice Test getTotalShares wrapper
     */
    function testGetTotalSharesWrapper(
        uint64 poolId
    ) external returns (GammTypes.Coin memory) {
        GammTypes.Coin memory totalShares = GammWrappers.getTotalShares(GAMM, poolId);
        emit TestResult("getTotalSharesWrapper", true, abi.encode(totalShares));
        return totalShares;
    }
    
    /**
     * @notice Test getTotalLiquidity wrapper
     */
    function testGetTotalLiquidityWrapper(
        uint64 poolId
    ) external returns (GammTypes.Coin[] memory) {
        GammTypes.Coin[] memory liquidity = GammWrappers.getTotalLiquidity(GAMM, poolId);
        emit TestResult("getTotalLiquidityWrapper", true, abi.encode(liquidity));
        return liquidity;
    }
    
    // ============ GammBuilders Tests ============
    
    /**
     * @notice Test JoinPoolBuilder
     */
    function testJoinPoolBuilder(
        uint64 poolId,
        uint256 shareOutAmount,
        string memory denom1,
        uint256 amount1,
        string memory denom2,
        uint256 amount2
    ) external returns (string memory) {
        GammBuilders.JoinPoolBuilder memory builder = GammBuilders.newJoinPool(poolId, shareOutAmount);
        builder = GammBuilders.addTokenInMax(builder, denom1, amount1);
        builder = GammBuilders.addTokenInMax(builder, denom2, amount2);
        string memory json = GammBuilders.build(builder);
        emit JSONConstructed("joinPoolBuilder", json);
        return json;
    }
    
    /**
     * @notice Test SwapBuilder
     */
    function testSwapBuilder(
        uint64 poolId,
        string memory tokenOutDenom,
        string memory tokenInDenom,
        uint256 tokenInAmount,
        uint256 tokenOutMinAmount
    ) external returns (string memory) {
        GammBuilders.SwapBuilder memory builder = GammBuilders.newSwap();
        builder = GammBuilders.addRoute(builder, poolId, tokenOutDenom);
        builder = GammBuilders.withTokenIn(builder, tokenInDenom, tokenInAmount);
        builder = GammBuilders.withTokenOutMinAmount(builder, tokenOutMinAmount);
        string memory json = GammBuilders.build(builder);
        emit JSONConstructed("swapBuilder", json);
        return json;
    }
    
    // ============ GammJSONHelpers Tests ============
    
    /**
     * @notice Test joinPoolJSON
     */
    function testJoinPoolJSON(
        uint64 poolId,
        uint256 shareOutAmount,
        string memory tokenInMaxsJson
    ) external returns (string memory) {
        string memory json = GammJSONHelpers.joinPoolJSON(poolId, shareOutAmount, tokenInMaxsJson);
        emit JSONConstructed("joinPoolJSON", json);
        return json;
    }
    
    /**
     * @notice Test exitPoolJSON
     */
    function testExitPoolJSON(
        uint64 poolId,
        uint256 shareInAmount,
        string memory tokenOutMinsJson
    ) external returns (string memory) {
        string memory json = GammJSONHelpers.exitPoolJSON(poolId, shareInAmount, tokenOutMinsJson);
        emit JSONConstructed("exitPoolJSON", json);
        return json;
    }
    
    /**
     * @notice Test swapExactAmountInJSON
     */
    function testSwapExactAmountInJSON(
        string memory routesJson,
        string memory tokenInJson,
        uint256 tokenOutMinAmount,
        string memory affiliatesJson
    ) external returns (string memory) {
        string memory json = GammJSONHelpers.swapExactAmountInJSON(
            routesJson,
            tokenInJson,
            tokenOutMinAmount,
            affiliatesJson
        );
        emit JSONConstructed("swapExactAmountInJSON", json);
        return json;
    }
    
    /**
     * @notice Test getPoolJSON
     */
    function testGetPoolJSON(
        uint64 poolId
    ) external returns (string memory) {
        string memory json = GammJSONHelpers.getPoolJSON(poolId);
        emit JSONConstructed("getPoolJSON", json);
        return json;
    }
    
    /**
     * @notice Test coinsToJson
     */
    function testCoinsToJson(
        GammTypes.Coin[] memory coins
    ) external returns (string memory) {
        string memory json = GammJSONHelpers.coinsToJson(coins);
        emit JSONConstructed("coinsToJson", json);
        return json;
    }
    
    /**
     * @notice Test swapRoutesToJson
     */
    function testSwapRoutesToJson(
        GammTypes.SwapAmountInRoute[] memory routes
    ) external returns (string memory) {
        string memory json = GammJSONHelpers.swapRoutesToJson(routes);
        emit JSONConstructed("swapRoutesToJson", json);
        return json;
    }
    
    /**
     * @notice Test affiliatesToJson
     */
    function testAffiliatesToJson(
        GammTypes.Affiliate[] memory affiliates
    ) external returns (string memory) {
        string memory json = GammJSONHelpers.affiliatesToJson(affiliates);
        emit JSONConstructed("affiliatesToJson", json);
        return json;
    }
}

