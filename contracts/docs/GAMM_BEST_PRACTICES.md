# Best Practices

Security, gas optimization, and design best practices for building contracts with the BitBadges gamm precompile.

## Table of Contents

- [Security Best Practices](#security-best-practices)
- [Gas Optimization](#gas-optimization)
- [Code Organization](#code-organization)
- [Error Handling](#error-handling)
- [Testing Strategies](#testing-strategies)

## Security Best Practices

### 1. Input Validation

Always validate inputs before calling precompile methods:

```solidity
function joinPool(
    uint64 poolId,
    uint256 shareOutAmount
) external {
    // Validate inputs
    GammErrors.requireValidPoolId(poolId);
    require(shareOutAmount > 0, "Share amount must be > 0");
    
    GammTypes.Coin[] memory tokenInMaxs = new GammTypes.Coin[](2);
    tokenInMaxs[0] = GammHelpers.createCoin("uatom", 1000000);
    tokenInMaxs[1] = GammHelpers.createCoin("uosmo", 2000000);
    
    // Perform join
    GammWrappers.joinPool(GAMM, poolId, shareOutAmount, tokenInMaxs);
}
```

### 2. Slippage Protection

Always use minimum amounts to protect against slippage:

```solidity
function swapWithSlippageProtection(
    uint64 poolId,
    uint256 tokenInAmount,
    uint256 minTokenOutAmount // Calculate based on expected price
) external {
    // Calculate minimum based on current pool state
    uint256 expectedOut = calculateExpectedOutput(poolId, tokenInAmount);
    uint256 minOut = (expectedOut * 95) / 100; // 5% slippage tolerance
    
    require(minTokenOutAmount <= minOut, "Slippage too high");
    
    // Perform swap
    GammTypes.SwapAmountInRoute[] memory routes = 
        GammHelpers.createSingleHopRoute(poolId, "uosmo");
    GammTypes.Coin memory tokenIn = GammHelpers.createCoin("uatom", tokenInAmount);
    
    uint256 tokenOut = GammWrappers.swapExactAmountIn(
        GAMM,
        routes,
        tokenIn,
        minTokenOutAmount,
        GammHelpers.createEmptyAffiliates()
    );
}
```

### 3. Reentrancy Protection

Use reentrancy guards for state-changing operations:

```solidity
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";

contract SecurePoolContract is ReentrancyGuard {
    function joinPool(...) external nonReentrant {
        // Join pool logic
    }
}
```

### 4. Check-Effects-Interactions Pattern

Follow the check-effects-interactions pattern:

```solidity
function exitPool(...) external {
    // 1. Checks
    require(userShares[poolId][msg.sender] >= shareInAmount, "Insufficient shares");
    
    // 2. Effects (state changes)
    userShares[poolId][msg.sender] -= shareInAmount;
    
    // 3. Interactions (external calls)
    GammWrappers.exitPool(GAMM, poolId, shareInAmount, tokenOutMins);
}
```

### 5. Address Validation

Always validate addresses:

```solidity
function addAffiliate(address affiliate) external {
    require(affiliate != address(0), "Invalid address");
    GammErrors.requireValidAffiliate(affiliate, 100); // 1% fee
}
```

### 6. Basis Points Validation

Validate affiliate fees are within bounds:

```solidity
function setAffiliateFee(uint256 basisPointsFee) external {
    require(basisPointsFee <= 10000, "Fee cannot exceed 100%");
    require(basisPointsFee > 0, "Fee must be > 0");
    // Use fee
}
```

## Gas Optimization

### 1. Use Wrappers Instead of Manual JSON

Wrappers are optimized and reduce gas costs:

```solidity
// ❌ Higher gas cost
string memory json = GammJSONHelpers.joinPoolJSON(...);
GAMM.joinPool(json);

// ✅ Lower gas cost (wrappers optimize JSON construction)
GammWrappers.joinPool(GAMM, poolId, shareOutAmount, tokenInMaxs);
```

### 2. Reuse Arrays When Possible

Avoid creating new arrays in loops:

```solidity
// ❌ Creates new array each iteration
for (uint256 i = 0; i < pools.length; i++) {
    GammTypes.Coin[] memory liquidity = GammWrappers.getTotalLiquidity(GAMM, pools[i]);
}

// ✅ Reuse array
GammTypes.Coin[] memory liquidity;
for (uint256 i = 0; i < pools.length; i++) {
    liquidity = GammWrappers.getTotalLiquidity(GAMM, pools[i]);
    // Process liquidity
}
```

### 3. Use View Functions for Calculations

Use query functions to calculate before executing:

```solidity
function calculateAndSwap(
    uint64 poolId,
    uint256 tokenInAmount
) external {
    // Calculate expected output first (view function, no gas)
    GammTypes.Coin[] memory tokenInMaxs = new GammTypes.Coin[](1);
    tokenInMaxs[0] = GammHelpers.createCoin("uatom", tokenInAmount);
    
    (uint256 expectedShares, ) = GammWrappers.calcJoinPoolNoSwapShares(
        GAMM,
        poolId,
        tokenInMaxs
    );
    
    // Only proceed if expected output meets requirements
    require(expectedShares >= minShares, "Output too low");
    
    // Execute swap
    // ...
}
```

### 4. Batch Operations

Batch multiple operations when possible:

```solidity
function batchSwap(
    uint64[] memory poolIds,
    uint256[] memory amounts
) external {
    require(poolIds.length == amounts.length, "Length mismatch");
    
    for (uint256 i = 0; i < poolIds.length; i++) {
        // Perform swaps in batch
        // ...
    }
}
```

## Code Organization

### 1. Use Libraries for Common Operations

Extract common patterns into libraries:

```solidity
library PoolOperations {
    function safeJoinPool(
        IGammPrecompile precompile,
        uint64 poolId,
        uint256 shareOutAmount,
        GammTypes.Coin[] memory tokenInMaxs
    ) internal returns (uint256) {
        GammErrors.requireValidPoolId(poolId);
        (uint256 shares, ) = GammWrappers.joinPool(
            precompile, poolId, shareOutAmount, tokenInMaxs
        );
        return shares;
    }
}
```

### 2. Separate Concerns

Separate swap logic from pool management:

```solidity
contract SwapRouter {
    // Swap-specific logic
}

contract PoolManager {
    // Pool management logic
}
```

### 3. Use Events for Off-Chain Tracking

Emit events for important operations:

```solidity
event PoolJoined(uint64 indexed poolId, address indexed user, uint256 shares);
event SwapExecuted(address indexed user, uint256 amountIn, uint256 amountOut);

function joinPool(...) external {
    (uint256 shares, ) = GammWrappers.joinPool(...);
    emit PoolJoined(poolId, msg.sender, shares);
}
```

## Error Handling

### 1. Use Custom Errors

Custom errors are more gas-efficient:

```solidity
function swap(...) external {
    if (poolId == 0) {
        revert GammErrors.InvalidPoolId(poolId);
    }
    // ...
}
```

### 2. Provide Clear Error Messages

Make error messages descriptive:

```solidity
require(tokenOutAmount >= minTokenOutAmount, 
    "Swap output below minimum: expected at least minTokenOutAmount");
```

### 3. Handle Precompile Errors

Check return values and handle errors:

```solidity
function safeSwap(...) external {
    try GammWrappers.swapExactAmountIn(...) returns (uint256 tokenOut) {
        // Success
    } catch Error(string memory reason) {
        revert(string(abi.encodePacked("Swap failed: ", reason)));
    }
}
```

## Testing Strategies

### 1. Test with Real Pool Data

Use actual pool IDs and realistic amounts:

```solidity
function testJoinPool() public {
    uint64 poolId = 1; // Use actual pool ID
    uint256 shareOutAmount = 1000000;
    // ...
}
```

### 2. Test Edge Cases

Test boundary conditions:

```solidity
function testSwapWithZeroAmount() public {
    // Should revert
    // ...
}

function testSwapWithMaxSlippage() public {
    // Should succeed
    // ...
}
```

### 3. Test Affiliate Fees

Verify affiliate fee calculations:

```solidity
function testAffiliateFee() public {
    uint256 basisPointsFee = 100; // 1%
    // Verify fee is calculated correctly
    // ...
}
```

### 4. Use Test Helpers

Use GammTestHelpers for consistent test data:

```solidity
import "../libraries/GammTestHelpers.sol";

function testSwap() public {
    GammTypes.Coin[] memory coins = GammTestHelpers.generateCoins(
        ["uatom", "uosmo"],
        [1000000, 2000000]
    );
    // ...
}
```














