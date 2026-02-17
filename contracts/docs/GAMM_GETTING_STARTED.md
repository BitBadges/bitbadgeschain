# Getting Started with BitBadges Gamm Precompile

This guide will help you get started building Solidity contracts that interact with the BitBadges gamm precompile for liquidity pool operations.

## Overview

The BitBadges gamm precompile provides a Solidity interface to the gamm module, enabling you to:
- Join and exit liquidity pools
- Perform token swaps (single-hop and multi-hop)
- Query pool information and calculate swap amounts
- Manage affiliate fees for swaps

## Precompile Address

```
0x0000000000000000000000000000000000001002
```

## Installation

### 1. Import the Required Files

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./interfaces/IGammPrecompile.sol";
import "./types/GammTypes.sol";
import "./libraries/GammWrappers.sol";
import "./libraries/GammHelpers.sol";
```

### 2. Initialize the Precompile Interface

```solidity
contract MyPoolContract {
    IGammPrecompile constant GAMM = 
        IGammPrecompile(0x0000000000000000000000000000000000001002);
    
    // Your contract code here
}
```

## Quick Start Examples

### Example 1: Join a Pool

```solidity
import "./libraries/GammWrappers.sol";
import "./libraries/GammHelpers.sol";

function joinPool(
    uint64 poolId,
    uint256 shareOutAmount
) external {
    // Create token inputs
    GammTypes.Coin[] memory tokenInMaxs = new GammTypes.Coin[](2);
    tokenInMaxs[0] = GammHelpers.createCoin("uatom", 1000000);
    tokenInMaxs[1] = GammHelpers.createCoin("uosmo", 2000000);
    
    // Use typed wrapper for type safety
    (uint256 shares, GammTypes.Coin[] memory tokens) = GammWrappers.joinPool(
        GAMM,
        poolId,
        shareOutAmount,
        tokenInMaxs
    );
    
    // Use shares and tokens as needed
}
```

### Example 2: Perform a Swap

```solidity
function swapTokens(
    uint64 poolId,
    string memory tokenInDenom,
    uint256 tokenInAmount,
    string memory tokenOutDenom,
    uint256 tokenOutMinAmount
) external {
    // Create swap route
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

### Example 3: Query Pool Information

```solidity
function getPoolInfo(uint64 poolId) external view returns (GammTypes.Coin memory) {
    // Get total shares
    GammTypes.Coin memory totalShares = GammWrappers.getTotalShares(GAMM, poolId);
    return totalShares;
}
```

## Using Builders

For complex operations, use the builder pattern:

```solidity
import "./libraries/GammBuilders.sol";

function joinPoolWithBuilder(uint64 poolId, uint256 shareOutAmount) external {
    GammBuilders.JoinPoolBuilder memory builder = GammBuilders.newJoinPool(poolId, shareOutAmount);
    builder = builder.addTokenInMax("uatom", 1000000);
    builder = builder.addTokenInMax("uosmo", 2000000);
    
    string memory json = builder.build();
    
    (uint256 shares, GammTypes.Coin[] memory tokens) = GAMM.joinPool(json);
}
```

## Next Steps

- Read the [API Reference](GAMM_API_REFERENCE.md) for complete method documentation
- Check out [Best Practices](GAMM_BEST_PRACTICES.md) for security and optimization tips
- See [Examples](GAMM_EXAMPLES.md) for more complex use cases















