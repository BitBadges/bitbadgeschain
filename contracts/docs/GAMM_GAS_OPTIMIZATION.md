# Gas Optimization

Strategies for optimizing gas costs when using the BitBadges gamm precompile.

## Table of Contents

- [JSON Construction Optimization](#json-construction-optimization)
- [Array Management](#array-management)
- [Function Selection](#function-selection)
- [Storage Optimization](#storage-optimization)
- [Batch Operations](#batch-operations)

## JSON Construction Optimization

### Use Wrappers Instead of Manual JSON

Wrappers optimize JSON construction internally:

```solidity
// ❌ Higher gas: Manual JSON construction
string memory tokenInMaxsJson = GammJSONHelpers.coinsToJson(tokenInMaxs);
string memory json = GammJSONHelpers.joinPoolJSON(poolId, shareOutAmount, tokenInMaxsJson);
GAMM.joinPool(json);

// ✅ Lower gas: Use wrappers (optimized JSON construction)
GammWrappers.joinPool(GAMM, poolId, shareOutAmount, tokenInMaxs);
```

**Gas Savings:** ~5,000-10,000 gas per call

### Use Builders for Complex Operations

Builders optimize array construction:

```solidity
// ❌ Higher gas: Manual array construction
GammTypes.Coin[] memory tokenInMaxs = new GammTypes.Coin[](2);
tokenInMaxs[0] = GammHelpers.createCoin("uatom", 1000000);
tokenInMaxs[1] = GammHelpers.createCoin("uosmo", 2000000);
string memory json = GammJSONHelpers.joinPoolJSON(poolId, shareOutAmount, GammJSONHelpers.coinsToJson(tokenInMaxs));

// ✅ Lower gas: Use builder
GammBuilders.JoinPoolBuilder memory builder = GammBuilders.newJoinPool(poolId, shareOutAmount);
builder = builder.addTokenInMax("uatom", 1000000);
builder = builder.addTokenInMax("uosmo", 2000000);
string memory json = builder.build();
```

## Array Management

### Reuse Arrays in Loops

Avoid creating new arrays in each iteration:

```solidity
// ❌ Higher gas: New array each iteration
for (uint256 i = 0; i < pools.length; i++) {
    GammTypes.Coin[] memory liquidity = GammWrappers.getTotalLiquidity(GAMM, pools[i]);
    // Process
}

// ✅ Lower gas: Reuse array
GammTypes.Coin[] memory liquidity;
for (uint256 i = 0; i < pools.length; i++) {
    liquidity = GammWrappers.getTotalLiquidity(GAMM, pools[i]);
    // Process
}
```

**Gas Savings:** ~20,000 gas per iteration

### Pre-allocate Arrays When Size is Known

```solidity
// ❌ Higher gas: Dynamic array growth
GammTypes.Coin[] memory coins;
for (uint256 i = 0; i < denoms.length; i++) {
    // Array grows dynamically
}

// ✅ Lower gas: Pre-allocate
GammTypes.Coin[] memory coins = new GammTypes.Coin[](denoms.length);
for (uint256 i = 0; i < denoms.length; i++) {
    coins[i] = GammHelpers.createCoin(denoms[i], amounts[i]);
}
```

## Function Selection

### Use View Functions for Calculations

Query functions are free (no gas for view calls):

```solidity
// ✅ Free: View function
(uint256 shares, ) = GammWrappers.calcJoinPoolNoSwapShares(
    GAMM, poolId, tokenInMaxs
);

// ❌ Costs gas: Transaction
(uint256 shares, ) = GammWrappers.joinPool(
    GAMM, poolId, shareOutAmount, tokenInMaxs
);
```

### Prefer Direct Queries Over Protobuf Decoding

Use typed return values when available:

```solidity
// ✅ Lower gas: Direct typed return
GammTypes.Coin memory totalShares = GammWrappers.getTotalShares(GAMM, poolId);

// ❌ Higher gas: Protobuf bytes (requires off-chain decoding anyway)
bytes memory poolBytes = GammWrappers.getPool(GAMM, poolId);
// Decoding not practical in Solidity
```

## Storage Optimization

### Use Memory for Temporary Data

Store temporary data in memory, not storage:

```solidity
// ❌ Higher gas: Storage
mapping(uint64 => GammTypes.Coin[]) public poolLiquidity;

function getLiquidity(uint64 poolId) external {
    poolLiquidity[poolId] = GammWrappers.getTotalLiquidity(GAMM, poolId);
}

// ✅ Lower gas: Memory
function getLiquidity(uint64 poolId) external view returns (GammTypes.Coin[] memory) {
    return GammWrappers.getTotalLiquidity(GAMM, poolId);
}
```

### Cache Frequently Accessed Values

```solidity
// ❌ Higher gas: Multiple storage reads
function swap(...) external {
    if (userShares[poolId][msg.sender] >= amount) {
        userShares[poolId][msg.sender] -= amount;
    }
}

// ✅ Lower gas: Cache storage value
function swap(...) external {
    uint256 shares = userShares[poolId][msg.sender];
    if (shares >= amount) {
        userShares[poolId][msg.sender] = shares - amount;
    }
}
```

## Batch Operations

### Batch Multiple Swaps

If possible, combine multiple operations:

```solidity
// ❌ Higher gas: Separate calls
function swapTwice(...) external {
    swap1(...);
    swap2(...);
}

// ✅ Lower gas: Single multi-hop swap
function swapMultiHop(...) external {
    GammTypes.SwapAmountInRoute[] memory routes = new GammTypes.SwapAmountInRoute[](2);
    routes[0] = GammHelpers.createSwapRoute(poolId1, "uosmo");
    routes[1] = GammHelpers.createSwapRoute(poolId2, "uion");
    // Single swap instead of two
}
```

### Use Affiliates Efficiently

If you need multiple affiliates, create the array once:

```solidity
// ❌ Higher gas: Create array multiple times
function swapWithAffiliates(...) external {
    // Array created each time
}

// ✅ Lower gas: Create once, reuse
GammTypes.Affiliate[] memory affiliates = new GammTypes.Affiliate[](2);
affiliates[0] = GammHelpers.createAffiliate(addr1, 50);
affiliates[1] = GammHelpers.createAffiliate(addr2, 50);
// Reuse for multiple swaps
```

## Gas Cost Estimates

Approximate gas costs for common operations:

| Operation | Gas Cost |
|-----------|----------|
| `joinPool` (2 tokens) | ~150,000 |
| `exitPool` (2 tokens) | ~120,000 |
| `swapExactAmountIn` (single-hop) | ~100,000 |
| `swapExactAmountIn` (two-hop) | ~150,000 |
| `getTotalShares` (view) | ~30,000 |
| `getTotalLiquidity` (view) | ~40,000 |
| `calcJoinPoolNoSwapShares` (view) | ~50,000 |

**Note:** Actual gas costs vary based on pool state, token amounts, and network conditions.

## Tips

1. **Use view functions** for calculations before executing transactions
2. **Reuse arrays** instead of creating new ones in loops
3. **Pre-allocate arrays** when size is known
4. **Use wrappers** instead of manual JSON construction
5. **Cache storage values** to reduce SLOAD operations
6. **Batch operations** when possible to reduce overhead
7. **Use memory** for temporary data instead of storage






