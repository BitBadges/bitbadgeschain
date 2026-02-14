# Extended Examples

Extended examples for the BitBadges gamm precompile. These examples demonstrate real-world use cases.

## Table of Contents

- [Simple Liquidity Provider](#simple-liquidity-provider)
- [Automated Market Maker Router](#automated-market-maker-router)
- [Affiliate Fee Distributor](#affiliate-fee-distributor)
- [Cross-Chain Swap Bridge](#cross-chain-swap-bridge)

## Simple Liquidity Provider

A contract that manages liquidity provision and removal.

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/IGammPrecompile.sol";
import "../libraries/GammWrappers.sol";
import "../libraries/GammHelpers.sol";
import "../libraries/GammBuilders.sol";

contract SimpleLiquidityProvider {
    IGammPrecompile constant GAMM = 
        IGammPrecompile(0x0000000000000000000000000000000000001002);
    
    mapping(uint64 => mapping(address => uint256)) public userShares;
    
    event LiquidityAdded(uint64 indexed poolId, address indexed user, uint256 shares);
    event LiquidityRemoved(uint64 indexed poolId, address indexed user, uint256 tokens);
    
    function addLiquidity(
        uint64 poolId,
        uint256 shareOutAmount,
        string[] memory denoms,
        uint256[] memory amounts
    ) external {
        require(denoms.length == amounts.length, "Length mismatch");
        
        // Build join pool request
        GammBuilders.JoinPoolBuilder memory builder = 
            GammBuilders.newJoinPool(poolId, shareOutAmount);
        
        for (uint256 i = 0; i < denoms.length; i++) {
            builder = builder.addTokenInMax(denoms[i], amounts[i]);
        }
        
        string memory json = builder.build();
        (uint256 shares, GammTypes.Coin[] memory tokens) = GAMM.joinPool(json);
        
        userShares[poolId][msg.sender] += shares;
        emit LiquidityAdded(poolId, msg.sender, shares);
    }
    
    function removeLiquidity(
        uint64 poolId,
        uint256 shareInAmount,
        string[] memory denoms,
        uint256[] memory minAmounts
    ) external {
        require(userShares[poolId][msg.sender] >= shareInAmount, "Insufficient shares");
        require(denoms.length == minAmounts.length, "Length mismatch");
        
        // Build exit pool request
        GammBuilders.ExitPoolBuilder memory builder = 
            GammBuilders.newExitPool(poolId, shareInAmount);
        
        for (uint256 i = 0; i < denoms.length; i++) {
            builder = builder.addTokenOutMin(denoms[i], minAmounts[i]);
        }
        
        string memory json = builder.build();
        GammTypes.Coin[] memory tokens = GAMM.exitPool(json);
        
        userShares[poolId][msg.sender] -= shareInAmount;
        emit LiquidityRemoved(poolId, msg.sender, tokens.length);
    }
    
    function getMyShares(uint64 poolId) external view returns (uint256) {
        return userShares[poolId][msg.sender];
    }
}
```

## Automated Market Maker Router

A contract that finds the best swap route and executes swaps.

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/IGammPrecompile.sol";
import "../libraries/GammWrappers.sol";
import "../libraries/GammHelpers.sol";
import "../libraries/GammBuilders.sol";

contract AMMRouter {
    IGammPrecompile constant GAMM = 
        IGammPrecompile(0x0000000000000000000000000000000000001002);
    
    struct SwapRoute {
        uint64[] poolIds;
        string[] intermediateDenoms;
        string finalDenom;
    }
    
    event SwapExecuted(
        address indexed user,
        string tokenIn,
        uint256 amountIn,
        string tokenOut,
        uint256 amountOut
    );
    
    function swap(
        SwapRoute memory route,
        string memory tokenInDenom,
        uint256 tokenInAmount,
        uint256 tokenOutMinAmount
    ) external returns (uint256) {
        require(route.poolIds.length > 0, "Empty route");
        require(route.poolIds.length == route.intermediateDenoms.length + 1, "Invalid route");
        
        // Build swap with multi-hop route
        GammBuilders.SwapBuilder memory builder = GammBuilders.newSwap();
        
        // Add intermediate hops
        for (uint256 i = 0; i < route.intermediateDenoms.length; i++) {
            builder = builder.addRoute(route.poolIds[i], route.intermediateDenoms[i]);
        }
        
        // Add final hop
        builder = builder.addRoute(
            route.poolIds[route.poolIds.length - 1],
            route.finalDenom
        );
        
        builder = builder.withTokenIn(tokenInDenom, tokenInAmount);
        builder = builder.withTokenOutMinAmount(tokenOutMinAmount);
        
        string memory json = builder.build();
        uint256 tokenOut = GAMM.swapExactAmountIn(json);
        
        emit SwapExecuted(
            msg.sender,
            tokenInDenom,
            tokenInAmount,
            route.finalDenom,
            tokenOut
        );
        
        return tokenOut;
    }
    
    function swapWithAffiliate(
        SwapRoute memory route,
        string memory tokenInDenom,
        uint256 tokenInAmount,
        uint256 tokenOutMinAmount,
        address affiliateAddress,
        uint256 basisPointsFee
    ) external returns (uint256) {
        require(route.poolIds.length > 0, "Empty route");
        
        GammBuilders.SwapBuilder memory builder = GammBuilders.newSwap();
        
        // Add routes
        for (uint256 i = 0; i < route.intermediateDenoms.length; i++) {
            builder = builder.addRoute(route.poolIds[i], route.intermediateDenoms[i]);
        }
        builder = builder.addRoute(
            route.poolIds[route.poolIds.length - 1],
            route.finalDenom
        );
        
        builder = builder.withTokenIn(tokenInDenom, tokenInAmount);
        builder = builder.withTokenOutMinAmount(tokenOutMinAmount);
        builder = builder.addAffiliate(affiliateAddress, basisPointsFee);
        
        string memory json = builder.build();
        return GAMM.swapExactAmountIn(json);
    }
}
```

## Affiliate Fee Distributor

A contract that collects and distributes affiliate fees from swaps.

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/IGammPrecompile.sol";
import "../libraries/GammWrappers.sol";
import "../libraries/GammHelpers.sol";

contract AffiliateFeeDistributor {
    IGammPrecompile constant GAMM = 
        IGammPrecompile(0x0000000000000000000000000000000000001002);
    
    address public owner;
    mapping(address => uint256) public affiliateBalances;
    
    event FeeCollected(address indexed affiliate, uint256 amount);
    event FeeWithdrawn(address indexed affiliate, uint256 amount);
    
    constructor() {
        owner = msg.sender;
    }
    
    modifier onlyOwner() {
        require(msg.sender == owner, "Not owner");
        _;
    }
    
    function executeSwapWithAffiliate(
        uint64 poolId,
        string memory tokenInDenom,
        uint256 tokenInAmount,
        string memory tokenOutDenom,
        uint256 tokenOutMinAmount,
        address affiliateAddress,
        uint256 basisPointsFee
    ) external returns (uint256) {
        GammTypes.SwapAmountInRoute[] memory routes = 
            GammHelpers.createSingleHopRoute(poolId, tokenOutDenom);
        
        GammTypes.Coin memory tokenIn = GammHelpers.createCoin(tokenInDenom, tokenInAmount);
        
        // Create affiliate pointing to this contract
        GammTypes.Affiliate[] memory affiliates = new GammTypes.Affiliate[](1);
        affiliates[0] = GammHelpers.createAffiliate(affiliateAddress, basisPointsFee);
        
        uint256 tokenOut = GammWrappers.swapExactAmountIn(
            GAMM,
            routes,
            tokenIn,
            tokenOutMinAmount,
            affiliates
        );
        
        // Note: Affiliate fees are automatically distributed by the precompile
        // This contract can track balances if needed
        
        return tokenOut;
    }
    
    function withdrawFees(address affiliate) external {
        require(affiliateBalances[affiliate] > 0, "No balance");
        
        uint256 amount = affiliateBalances[affiliate];
        affiliateBalances[affiliate] = 0;
        
        // Transfer logic here (implementation depends on token type)
        emit FeeWithdrawn(affiliate, amount);
    }
}
```

## Cross-Chain Swap Bridge

A contract that performs swaps and transfers via IBC.

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/IGammPrecompile.sol";
import "../libraries/GammWrappers.sol";
import "../libraries/GammHelpers.sol";

contract CrossChainSwapBridge {
    IGammPrecompile constant GAMM = 
        IGammPrecompile(0x0000000000000000000000000000000000001002);
    
    mapping(string => bool) public supportedChannels;
    address public bridgeOperator;
    
    event CrossChainSwap(
        address indexed user,
        string sourceChannel,
        string receiver,
        uint256 amountOut
    );
    
    constructor() {
        bridgeOperator = msg.sender;
    }
    
    modifier onlyOperator() {
        require(msg.sender == bridgeOperator, "Not operator");
        _;
    }
    
    function addSupportedChannel(string memory channel) external onlyOperator {
        supportedChannels[channel] = true;
    }
    
    function removeSupportedChannel(string memory channel) external onlyOperator {
        supportedChannels[channel] = false;
    }
    
    function swapAndBridge(
        uint64 poolId,
        string memory tokenInDenom,
        uint256 tokenInAmount,
        string memory tokenOutDenom,
        uint256 tokenOutMinAmount,
        string memory sourceChannel,
        string memory receiver,
        string memory memo
    ) external returns (uint256) {
        require(supportedChannels[sourceChannel], "Channel not supported");
        
        GammTypes.SwapAmountInRoute[] memory routes = 
            GammHelpers.createSingleHopRoute(poolId, tokenOutDenom);
        
        GammTypes.Coin memory tokenIn = GammHelpers.createCoin(tokenInDenom, tokenInAmount);
        
        // Create IBC transfer info
        GammTypes.IBCTransferInfo memory ibcInfo = GammHelpers.createIBCTransferInfo(
            sourceChannel,
            receiver,
            memo,
            0 // Use default timeout
        );
        
        GammTypes.Affiliate[] memory affiliates = GammHelpers.createEmptyAffiliates();
        
        uint256 tokenOut = GammWrappers.swapExactAmountInWithIBCTransfer(
            GAMM,
            routes,
            tokenIn,
            tokenOutMinAmount,
            ibcInfo,
            affiliates
        );
        
        emit CrossChainSwap(msg.sender, sourceChannel, receiver, tokenOut);
        
        return tokenOut;
    }
}
```











