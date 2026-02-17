# Common Patterns

Common patterns and use cases for building contracts with the BitBadges gamm precompile.

## Table of Contents

- [Basic Pool Operations](#basic-pool-operations)
- [Swap Patterns](#swap-patterns)
- [Multi-Hop Swaps](#multi-hop-swaps)
- [Liquidity Management](#liquidity-management)
- [Affiliate Fees](#affiliate-fees)
- [IBC Transfers](#ibc-transfers)

## Basic Pool Operations

### Join a Pool

```solidity
function joinPool(
    uint64 poolId,
    uint256 shareOutAmount
) external {
    GammTypes.Coin[] memory tokenInMaxs = new GammTypes.Coin[](2);
    tokenInMaxs[0] = GammHelpers.createCoin("uatom", 1000000);
    tokenInMaxs[1] = GammHelpers.createCoin("uosmo", 2000000);
    
    (uint256 shares, GammTypes.Coin[] memory tokens) = GammWrappers.joinPool(
        GAMM,
        poolId,
        shareOutAmount,
        tokenInMaxs
    );
    
    // Use shares and tokens as needed
}
```

### Exit a Pool

```solidity
function exitPool(
    uint64 poolId,
    uint256 shareInAmount
) external {
    GammTypes.Coin[] memory tokenOutMins = new GammTypes.Coin[](2);
    tokenOutMins[0] = GammHelpers.createCoin("uatom", 500000);
    tokenOutMins[1] = GammHelpers.createCoin("uosmo", 1000000);
    
    GammTypes.Coin[] memory tokens = GammWrappers.exitPool(
        GAMM,
        poolId,
        shareInAmount,
        tokenOutMins
    );
    
    // Use tokens as needed
}
```

### Using Builders

```solidity
function joinPoolWithBuilder(uint64 poolId, uint256 shareOutAmount) external {
    GammBuilders.JoinPoolBuilder memory builder = 
        GammBuilders.newJoinPool(poolId, shareOutAmount);
    
    builder = builder.addTokenInMax("uatom", 1000000);
    builder = builder.addTokenInMax("uosmo", 2000000);
    
    string memory json = builder.build();
    (uint256 shares, GammTypes.Coin[] memory tokens) = GAMM.joinPool(json);
}
```

## Swap Patterns

### Simple Single-Hop Swap

```solidity
function swapTokens(
    uint64 poolId,
    string memory tokenInDenom,
    uint256 tokenInAmount,
    string memory tokenOutDenom,
    uint256 tokenOutMinAmount
) external {
    // Create single-hop route
    GammTypes.SwapAmountInRoute[] memory routes = 
        GammHelpers.createSingleHopRoute(poolId, tokenOutDenom);
    
    // Create input token
    GammTypes.Coin memory tokenIn = GammHelpers.createCoin(tokenInDenom, tokenInAmount);
    
    // No affiliates
    GammTypes.Affiliate[] memory affiliates = GammHelpers.createEmptyAffiliates();
    
    // Perform swap
    uint256 tokenOutAmount = GammWrappers.swapExactAmountIn(
        GAMM,
        routes,
        tokenIn,
        tokenOutMinAmount,
        affiliates
    );
    
    // Use tokenOutAmount as needed
}
```

### Swap with Builder

```solidity
function swapWithBuilder(
    uint64 poolId,
    string memory tokenInDenom,
    uint256 tokenInAmount,
    string memory tokenOutDenom,
    uint256 tokenOutMinAmount
) external {
    GammBuilders.SwapBuilder memory builder = GammBuilders.newSwap();
    
    builder = builder.addRoute(poolId, tokenOutDenom);
    builder = builder.withTokenIn(tokenInDenom, tokenInAmount);
    builder = builder.withTokenOutMinAmount(tokenOutMinAmount);
    
    string memory json = builder.build();
    uint256 tokenOut = GAMM.swapExactAmountIn(json);
}
```

## Multi-Hop Swaps

### Two-Hop Swap

```solidity
function twoHopSwap(
    uint64 poolId1,
    string memory intermediateDenom,
    uint64 poolId2,
    string memory tokenOutDenom,
    string memory tokenInDenom,
    uint256 tokenInAmount,
    uint256 tokenOutMinAmount
) external {
    // Create two-hop route
    GammTypes.SwapAmountInRoute[] memory routes = 
        GammHelpers.createTwoHopRoute(
            poolId1,
            intermediateDenom,
            poolId2,
            tokenOutDenom
        );
    
    GammTypes.Coin memory tokenIn = GammHelpers.createCoin(tokenInDenom, tokenInAmount);
    GammTypes.Affiliate[] memory affiliates = GammHelpers.createEmptyAffiliates();
    
    uint256 tokenOut = GammWrappers.swapExactAmountIn(
        GAMM,
        routes,
        tokenIn,
        tokenOutMinAmount,
        affiliates
    );
}
```

### Multi-Hop Swap with Builder

```solidity
function multiHopSwapWithBuilder(
    uint64[] memory poolIds,
    string[] memory intermediateDenoms,
    string memory finalDenom,
    string memory tokenInDenom,
    uint256 tokenInAmount,
    uint256 tokenOutMinAmount
) external {
    require(poolIds.length == intermediateDenoms.length + 1, "Invalid route");
    
    GammBuilders.SwapBuilder memory builder = GammBuilders.newSwap();
    
    // Add intermediate hops
    for (uint256 i = 0; i < intermediateDenoms.length; i++) {
        builder = builder.addRoute(poolIds[i], intermediateDenoms[i]);
    }
    
    // Add final hop
    builder = builder.addRoute(poolIds[poolIds.length - 1], finalDenom);
    builder = builder.withTokenIn(tokenInDenom, tokenInAmount);
    builder = builder.withTokenOutMinAmount(tokenOutMinAmount);
    
    string memory json = builder.build();
    uint256 tokenOut = GAMM.swapExactAmountIn(json);
}
```

## Liquidity Management

### Add Liquidity and Track Shares

```solidity
contract LiquidityManager {
    IGammPrecompile constant GAMM = 
        IGammPrecompile(0x0000000000000000000000000000000000001002);
    
    mapping(uint64 => mapping(address => uint256)) public userShares;
    
    function addLiquidity(
        uint64 poolId,
        uint256 shareOutAmount
    ) external {
        GammTypes.Coin[] memory tokenInMaxs = new GammTypes.Coin[](2);
        tokenInMaxs[0] = GammHelpers.createCoin("uatom", 1000000);
        tokenInMaxs[1] = GammHelpers.createCoin("uosmo", 2000000);
        
        (uint256 shares, ) = GammWrappers.joinPool(
            GAMM,
            poolId,
            shareOutAmount,
            tokenInMaxs
        );
        
        userShares[poolId][msg.sender] += shares;
    }
    
    function removeLiquidity(
        uint64 poolId,
        uint256 shareInAmount
    ) external {
        require(userShares[poolId][msg.sender] >= shareInAmount, "Insufficient shares");
        
        GammTypes.Coin[] memory tokenOutMins = new GammTypes.Coin[](2);
        tokenOutMins[0] = GammHelpers.createCoin("uatom", 500000);
        tokenOutMins[1] = GammHelpers.createCoin("uosmo", 1000000);
        
        GammTypes.Coin[] memory tokens = GammWrappers.exitPool(
            GAMM,
            poolId,
            shareInAmount,
            tokenOutMins
        );
        
        userShares[poolId][msg.sender] -= shareInAmount;
        
        // Use tokens as needed
    }
}
```

## Affiliate Fees

### Swap with Affiliate Fee

```solidity
function swapWithAffiliate(
    uint64 poolId,
    string memory tokenInDenom,
    uint256 tokenInAmount,
    string memory tokenOutDenom,
    uint256 tokenOutMinAmount,
    address affiliateAddress,
    uint256 basisPointsFee // e.g., 100 = 1%
) external {
    GammTypes.SwapAmountInRoute[] memory routes = 
        GammHelpers.createSingleHopRoute(poolId, tokenOutDenom);
    
    GammTypes.Coin memory tokenIn = GammHelpers.createCoin(tokenInDenom, tokenInAmount);
    
    // Create affiliate
    GammTypes.Affiliate[] memory affiliates = new GammTypes.Affiliate[](1);
    affiliates[0] = GammHelpers.createAffiliate(affiliateAddress, basisPointsFee);
    
    uint256 tokenOut = GammWrappers.swapExactAmountIn(
        GAMM,
        routes,
        tokenIn,
        tokenOutMinAmount,
        affiliates
    );
}
```

### Multiple Affiliates

```solidity
function swapWithMultipleAffiliates(
    uint64 poolId,
    address[] memory affiliateAddresses,
    uint256[] memory basisPointsFees
) external {
    require(affiliateAddresses.length == basisPointsFees.length, "Length mismatch");
    
    GammTypes.SwapAmountInRoute[] memory routes = 
        GammHelpers.createSingleHopRoute(poolId, "uosmo");
    
    GammTypes.Coin memory tokenIn = GammHelpers.createCoin("uatom", 1000000);
    
    // Create affiliates array
    GammTypes.Affiliate[] memory affiliates = new GammTypes.Affiliate[](affiliateAddresses.length);
    for (uint256 i = 0; i < affiliateAddresses.length; i++) {
        affiliates[i] = GammHelpers.createAffiliate(affiliateAddresses[i], basisPointsFees[i]);
    }
    
    uint256 tokenOut = GammWrappers.swapExactAmountIn(
        GAMM,
        routes,
        tokenIn,
        900000,
        affiliates
    );
}
```

## IBC Transfers

### Swap with IBC Transfer

```solidity
function swapAndTransferIBC(
    uint64 poolId,
    string memory sourceChannel,
    string memory receiver,
    string memory memo
) external {
    GammTypes.SwapAmountInRoute[] memory routes = 
        GammHelpers.createSingleHopRoute(poolId, "uosmo");
    
    GammTypes.Coin memory tokenIn = GammHelpers.createCoin("uatom", 1000000);
    
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
        900000,
        ibcInfo,
        affiliates
    );
}
```

## Query Patterns

### Check Pool Liquidity

```solidity
function checkPoolLiquidity(uint64 poolId) external view returns (uint256) {
    GammTypes.Coin[] memory liquidity = GammWrappers.getTotalLiquidity(GAMM, poolId);
    
    uint256 totalLiquidity = 0;
    for (uint256 i = 0; i < liquidity.length; i++) {
        totalLiquidity += liquidity[i].amount;
    }
    
    return totalLiquidity;
}
```

### Calculate Expected Output

```solidity
function calculateSwapOutput(
    uint64 poolId,
    string memory tokenInDenom,
    uint256 tokenInAmount
) external view returns (uint256) {
    GammTypes.Coin[] memory tokenInMaxs = new GammTypes.Coin[](1);
    tokenInMaxs[0] = GammHelpers.createCoin(tokenInDenom, tokenInAmount);
    
    (uint256 shares, ) = GammWrappers.calcJoinPoolNoSwapShares(
        GAMM,
        poolId,
        tokenInMaxs
    );
    
    return shares;
}
```















