# Troubleshooting

Common issues and solutions when working with the BitBadges gamm precompile.

## Table of Contents

- [Common Errors](#common-errors)
- [JSON Construction Issues](#json-construction-issues)
- [Swap Issues](#swap-issues)
- [Pool Issues](#pool-issues)
- [Gas Issues](#gas-issues)
- [IBC Transfer Issues](#ibc-transfer-issues)

## Common Errors

### "Invalid pool ID"

**Error:** `InvalidPoolId` or pool not found.

**Solutions:**
1. Verify the pool ID is correct and exists
2. Check that the pool is active (not closed)
3. Ensure you're using the correct network/chain

```solidity
// Always validate pool ID
GammErrors.requireValidPoolId(poolId);

// Check pool exists by querying
GammTypes.Coin memory totalShares = GammWrappers.getTotalShares(GAMM, poolId);
require(totalShares.amount > 0, "Pool does not exist or has no liquidity");
```

### "Slippage tolerance exceeded"

**Error:** Swap fails because output is below minimum.

**Solutions:**
1. Calculate expected output first using `calcJoinPoolNoSwapShares` or `calcJoinPoolShares`
2. Set realistic minimum amounts (account for fees and price impact)
3. Use a slippage tolerance (e.g., 1-5%)

```solidity
// Calculate expected output first
GammTypes.Coin[] memory tokenInMaxs = new GammTypes.Coin[](1);
tokenInMaxs[0] = GammHelpers.createCoin("uatom", tokenInAmount);

(, GammTypes.Coin[] memory expectedTokens) = GammWrappers.calcJoinPoolNoSwapShares(
    GAMM, poolId, tokenInMaxs
);

// Set minimum with slippage tolerance (e.g., 5%)
uint256 minAmount = (expectedTokens[0].amount * 95) / 100;
```

### "Insufficient liquidity"

**Error:** Pool doesn't have enough liquidity for the operation.

**Solutions:**
1. Check pool liquidity using `getTotalLiquidity`
2. Reduce the amount you're trying to swap/join
3. Try a different pool if available

```solidity
// Check liquidity before operation
GammTypes.Coin[] memory liquidity = GammWrappers.getTotalLiquidity(GAMM, poolId);
require(liquidity[0].amount >= minRequired, "Insufficient pool liquidity");
```

## JSON Construction Issues

### Malformed JSON

**Error:** JSON string is invalid or malformed.

**Solutions:**
1. Use helper libraries instead of manual construction:
```solidity
// ❌ Don't do this:
string memory json = string(abi.encodePacked('{"poolId":"', poolId, '"}'));
// May have issues with escaping, formatting, etc.

// ✅ Do this:
string memory json = GammJSONHelpers.getPoolJSON(poolId);
```

2. Verify JSON structure matches protobuf format
3. Check string escaping for special characters

### Missing Required Fields

**Error:** JSON missing required fields.

**Solutions:**
1. Use wrappers which ensure all required fields are included
2. Check protobuf message structure for required fields
3. Use builders which validate required fields

```solidity
// ✅ Wrappers ensure all fields are included
GammWrappers.joinPool(GAMM, poolId, shareOutAmount, tokenInMaxs);

// ✅ Builders validate required fields
GammBuilders.SwapBuilder memory builder = GammBuilders.newSwap();
builder = builder.addRoute(poolId, tokenOutDenom);
builder = builder.withTokenIn(tokenInDenom, amount);
builder = builder.withTokenOutMinAmount(minAmount);
string memory json = builder.build(); // Validates all required fields
```

## Swap Issues

### "Route not found" or "Invalid route"

**Error:** Swap route is invalid.

**Solutions:**
1. Verify pool IDs in the route exist
2. Check that token denominations match pool assets
3. Ensure route is valid (intermediate tokens must match)

```solidity
// Validate route before swapping
function validateRoute(
    uint64[] memory poolIds,
    string[] memory denoms
) external view {
    require(poolIds.length == denoms.length, "Route length mismatch");
    
    for (uint256 i = 0; i < poolIds.length; i++) {
        GammErrors.requireValidPoolId(poolIds[i]);
        // Additional validation...
    }
}
```

### "Token out amount below minimum"

**Error:** Swap output is less than `tokenOutMinAmount`.

**Solutions:**
1. Calculate expected output first
2. Set realistic minimum amounts
3. Account for fees and price impact

```solidity
// Calculate expected output
// Then set minimum with slippage tolerance
uint256 expectedOut = calculateExpectedOutput(...);
uint256 minOut = (expectedOut * 95) / 100; // 5% slippage
```

## Pool Issues

### "Pool not found"

**Error:** Pool doesn't exist.

**Solutions:**
1. Verify pool ID is correct
2. Check pool exists using `getPool` or `getTotalShares`
3. Ensure you're on the correct network

### "Insufficient shares"

**Error:** Trying to exit more shares than owned.

**Solutions:**
1. Check actual shares before exiting:
```solidity
// Query total shares for the pool
GammTypes.Coin memory totalShares = GammWrappers.getTotalShares(GAMM, poolId);

// Check user's share balance (if tracking in contract)
require(userShares[poolId][msg.sender] >= shareInAmount, "Insufficient shares");
```

2. Track shares in your contract if needed
3. Use smaller amounts if shares are insufficient

## Gas Issues

### "Out of gas"

**Error:** Transaction runs out of gas.

**Solutions:**
1. Reduce operation complexity (fewer routes, smaller arrays)
2. Split large operations into smaller batches
3. Use view functions for calculations first
4. Optimize array usage (reuse, pre-allocate)

```solidity
// ❌ May run out of gas with many routes
GammTypes.SwapAmountInRoute[] memory routes = new GammTypes.SwapAmountInRoute[](10);
// ...

// ✅ Split into smaller operations
for (uint256 i = 0; i < 5; i++) {
    // Process in batches
}
```

### High Gas Costs

**Solutions:**
1. Use wrappers instead of manual JSON construction
2. Reuse arrays instead of creating new ones
3. Use view functions for calculations
4. Batch operations when possible

## IBC Transfer Issues

### "Invalid IBC channel"

**Error:** IBC channel doesn't exist or is invalid.

**Solutions:**
1. Verify channel ID is correct
2. Check channel exists and is active
3. Ensure channel supports the token denomination

```solidity
// Validate IBC transfer info
GammErrors.requireValidIBCTransferInfo(sourceChannel, receiver);

// Additional validation
require(bytes(sourceChannel).length > 0, "Channel cannot be empty");
require(bytes(receiver).length > 0, "Receiver cannot be empty");
```

### "IBC transfer timeout"

**Error:** IBC transfer times out.

**Solutions:**
1. Set appropriate timeout timestamp
2. Use default timeout (0) for automatic calculation
3. Ensure sufficient time for cross-chain transfer

```solidity
// Use default timeout (0 = automatic)
GammTypes.IBCTransferInfo memory ibcInfo = GammHelpers.createIBCTransferInfo(
    sourceChannel,
    receiver,
    memo,
    0 // Use default timeout
);
```

## Affiliate Issues

### "Invalid affiliate"

**Error:** Affiliate address or fee is invalid.

**Solutions:**
1. Validate affiliate address is not zero
2. Ensure basis points fee is between 0 and 10000
3. Use validation helpers:

```solidity
GammErrors.requireValidAffiliate(affiliateAddress, basisPointsFee);
```

### "Affiliate fee too high"

**Error:** Basis points fee exceeds maximum (10000 = 100%).

**Solutions:**
1. Ensure fee is between 0 and 10000
2. Validate before creating affiliate:

```solidity
require(basisPointsFee <= 10000, "Fee cannot exceed 100%");
GammTypes.Affiliate memory affiliate = GammHelpers.createAffiliate(
    address_,
    basisPointsFee
);
```

## Debugging Tips

### 1. Use Events

Emit events to track operations:

```solidity
event SwapAttempted(
    address indexed user,
    uint64 poolId,
    uint256 amountIn,
    uint256 minAmountOut
);

function swap(...) external {
    emit SwapAttempted(msg.sender, poolId, tokenInAmount, tokenOutMinAmount);
    // Perform swap
}
```

### 2. Check Return Values

Always check return values:

```solidity
(uint256 shares, GammTypes.Coin[] memory tokens) = GammWrappers.joinPool(...);
require(shares > 0, "No shares received");
require(tokens.length > 0, "No tokens provided");
```

### 3. Use View Functions First

Test with view functions before executing:

```solidity
// Test calculation first (free, view function)
(, GammTypes.Coin[] memory expectedTokens) = GammWrappers.calcJoinPoolNoSwapShares(
    GAMM, poolId, tokenInMaxs
);

// Only proceed if expected output is acceptable
require(expectedTokens[0].amount >= minAmount, "Expected output too low");

// Then execute
GammWrappers.joinPool(GAMM, poolId, shareOutAmount, tokenInMaxs);
```

### 4. Validate Pool State

Check pool state before operations:

```solidity
function safeSwap(...) external {
    // Check pool has liquidity
    GammTypes.Coin[] memory liquidity = GammWrappers.getTotalLiquidity(GAMM, poolId);
    require(liquidity.length > 0, "Pool has no liquidity");
    
    // Perform swap
    // ...
}
```


















