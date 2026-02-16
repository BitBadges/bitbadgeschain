# Gas Optimization Guide

Strategies for optimizing gas costs when using the BitBadges tokenization precompile.

## Table of Contents

- [Query Optimization](#query-optimization)
- [Transaction Optimization](#transaction-optimization)
- [Storage Patterns](#storage-patterns)
- [JSON Construction](#json-construction)
- [Batch Operations](#batch-operations)

## Query Optimization

### Use Direct Queries

Prefer queries that return simple types instead of protobuf bytes:

```solidity
// ✅ Gas efficient - ~3,000 gas
uint256 balance = TokenizationWrappers.getBalanceAmount(
    TOKENIZATION, collectionId, user, tokenIds, ownershipTimes
);

// ❌ Less efficient - ~5,000+ gas + decoding cost
bytes memory balanceBytes = TokenizationWrappers.getBalance(
    TOKENIZATION, collectionId, user
);
// Requires off-chain decoding
```

**Gas Savings:** ~2,000+ gas per query

### Cache Query Results

Cache frequently accessed query results:

```solidity
contract CachedQueries {
    mapping(address => uint256) private cachedBalances;
    mapping(address => uint256) private cacheTimestamp;
    uint256 private constant CACHE_DURATION = 1 hours;
    
    function getCachedBalance(address user) external view returns (uint256) {
        if (block.timestamp < cacheTimestamp[user] + CACHE_DURATION) {
            return cachedBalances[user];
        }
        // Refresh cache
        uint256 balance = TokenizationWrappers.getBalanceAmount(...);
        cachedBalances[user] = balance;
        cacheTimestamp[user] = block.timestamp;
        return balance;
    }
}
```

### Minimize Query Ranges

Use the smallest possible token ID and ownership time ranges:

```solidity
// ✅ Efficient - single range
TokenizationTypes.UintRange[] memory tokenIds = new TokenizationTypes.UintRange[](1);
tokenIds[0] = TokenizationHelpers.createSingleTokenIdRange(tokenId);

// ❌ Less efficient - multiple ranges
TokenizationTypes.UintRange[] memory tokenIds = new TokenizationTypes.UintRange[](10);
// ... many ranges
```

## Transaction Optimization

### Use Typed Wrappers

Typed wrappers are optimized for gas efficiency:

```solidity
// ✅ Efficient - optimized JSON construction
TokenizationWrappers.transferTokens(
    TOKENIZATION, collectionId, recipients, amount, tokenIds, ownershipTimes
);

// ❌ Less efficient - manual JSON construction
string memory json = TokenizationJSONHelpers.transferTokensJSON(...);
TOKENIZATION.transferTokens(json);
```

**Gas Savings:** ~500-1,000 gas per transaction

### Batch Operations

Combine multiple operations into single transactions:

```solidity
// ✅ Efficient - single transaction
function batchTransfer(
    address[] memory recipients,
    uint256[] memory amounts
) external {
    for (uint256 i = 0; i < recipients.length; i++) {
        TokenizationWrappers.transferSingleToken(...);
    }
}

// ❌ Less efficient - multiple transactions
// Each transfer requires separate transaction overhead
```

**Gas Savings:** ~21,000 gas per additional operation (transaction overhead)

### Optimize Range Arrays

Pre-allocate arrays when size is known:

```solidity
// ✅ Efficient - known size
TokenizationTypes.UintRange[] memory tokenIds = new TokenizationTypes.UintRange[](1);
tokenIds[0] = TokenizationHelpers.createSingleTokenIdRange(tokenId);

// ❌ Less efficient - dynamic sizing
TokenizationTypes.UintRange[] memory tokenIds;
// ... push operations
```

## Storage Patterns

### Use Memory for Temporary Data

Store temporary data in memory, not storage:

```solidity
function transfer(...) external {
    // ✅ Memory (cheaper)
    uint256 tempCollectionId = collectionId;
    
    // ❌ Storage (more expensive)
    // Accessing collectionId directly from storage multiple times
}
```

**Gas Savings:** ~100 gas per storage read avoided

### Pack Structs Efficiently

When creating custom structs, pack them efficiently:

```solidity
// ✅ Efficient - packed struct
struct TransferData {
    uint128 amount;      // Fits in 128 bits
    uint64 tokenId;      // Fits in 64 bits
    uint64 timestamp;    // Fits in 64 bits
    // Total: 256 bits = 1 storage slot
}

// ❌ Less efficient - not packed
struct TransferData {
    uint256 amount;      // 1 slot
    uint256 tokenId;     // 1 slot
    uint256 timestamp;   // 1 slot
    // Total: 3 storage slots
}
```

### Cache Storage Values

Cache frequently accessed storage values:

```solidity
function complexOperation() external {
    // ✅ Cache once
    uint256 cachedId = collectionId;
    uint256 cachedRegistry = kycRegistryId;
    
    // Use cached values multiple times
    // ...
}

// ❌ Access storage multiple times
function complexOperation() external {
    // Access collectionId from storage multiple times
    // ...
}
```

**Gas Savings:** ~100 gas per additional storage read avoided

## JSON Construction

### Use Builders for Complex Operations

Builders optimize JSON construction:

```solidity
// ✅ Efficient - optimized builder
TokenizationBuilders.CollectionBuilder memory builder = 
    TokenizationBuilders.newCollection();
builder = builder.withValidTokenIdRange(1, 1000);
string memory json = builder.build();

// ❌ Less efficient - manual concatenation
string memory json = string(abi.encodePacked(
    '{"validTokenIds":[{"start":"1","end":"1000"}]',
    // ... many more concatenations
));
```

### Minimize String Operations

Reduce string concatenations and operations:

```solidity
// ✅ Efficient - single operation
string memory json = TokenizationJSONHelpers.getCollectionJSON(collectionId);

// ❌ Less efficient - multiple concatenations
string memory json = string(abi.encodePacked(
    '{"collectionId":"',
    uintToString(collectionId),
    '"}'
));
```

## Batch Operations

### Efficient Batching

Batch operations to amortize transaction overhead:

```solidity
function efficientBatchTransfer(
    address[] memory recipients,
    uint256[] memory amounts,
    uint256 tokenId
) external {
    // Prepare ranges once
    TokenizationTypes.UintRange[] memory tokenIds = new TokenizationTypes.UintRange[](1);
    tokenIds[0] = TokenizationHelpers.createSingleTokenIdRange(tokenId);
    
    TokenizationTypes.UintRange[] memory ownershipTimes = new TokenizationTypes.UintRange[](1);
    ownershipTimes[0] = TokenizationHelpers.createFullOwnershipTimeRange();
    
    // Reuse ranges for all transfers
    for (uint256 i = 0; i < recipients.length; i++) {
        address[] memory singleRecipient = new address[](1);
        singleRecipient[0] = recipients[i];
        
        TokenizationWrappers.transferTokens(
            TOKENIZATION,
            collectionId,
            singleRecipient,
            amounts[i],
            tokenIds,
            ownershipTimes
        );
    }
}
```

**Gas Savings:** Reusing ranges saves ~500-1,000 gas per transfer

### Limit Batch Sizes

Balance batch size with gas limits:

```solidity
uint256 public constant MAX_BATCH_SIZE = 50;  // Adjust based on gas limits

function batchTransfer(...) external {
    require(recipients.length <= MAX_BATCH_SIZE, "Batch too large");
    // ...
}
```

## Additional Optimization Tips

### 1. Use Events for Off-Chain Data

Instead of storing data on-chain, emit events:

```solidity
// ✅ Efficient - event (cheap)
event TransferRecorded(
    uint256 indexed collectionId,
    address indexed from,
    address indexed to,
    uint256 amount
);

// ❌ Less efficient - storage (expensive)
mapping(address => Transfer[]) public transferHistory;
```

**Gas Savings:** ~20,000 gas per storage write avoided

### 2. Avoid Unnecessary Loops

Minimize loop iterations:

```solidity
// ✅ Efficient - single iteration
function transferToSingle(address to, uint256 amount) external {
    // Single transfer
}

// ❌ Less efficient - loop with checks
function transferToMany(address[] memory recipients) external {
    for (uint256 i = 0; i < recipients.length; i++) {
        // Multiple transfers with overhead
    }
}
```

### 3. Use View Functions

Use view functions for read-only operations:

```solidity
// ✅ Efficient - view function (no gas cost for callers)
function getBalance(address user) external view returns (uint256) {
    return TokenizationWrappers.getBalanceAmount(...);
}

// ❌ Less efficient - state-changing function
function getBalance(address user) external returns (uint256) {
    // Costs gas even for reads
}
```

### 4. Optimize Error Messages

Use custom errors instead of string errors:

```solidity
// ✅ Efficient - custom error
error InsufficientBalance(uint256 required, uint256 available);

// ❌ Less efficient - string error
require(balance >= amount, "Insufficient balance");
```

**Gas Savings:** ~200-500 gas per error

## Gas Cost Reference

Approximate gas costs for common operations:

| Operation | Gas Cost |
|-----------|----------|
| `getBalanceAmount()` | ~3,000 |
| `getBalance()` | ~5,000+ |
| `transferTokens()` | ~30,000+ |
| `createDynamicStore()` | ~20,000 |
| `setDynamicStoreValue()` | ~15,000 |
| `getCollection()` | ~5,000+ |

*Note: Actual gas costs vary based on input size and network conditions.*

## Measurement Tools

Use gas measurement tools to verify optimizations:

```solidity
// In tests
uint256 gasBefore = gasleft();
transfer(...);
uint256 gasUsed = gasBefore - gasleft();
console.log("Gas used:", gasUsed);
```

## See Also

- [Best Practices](./BEST_PRACTICES.md) for general optimization strategies
- [API Reference](./API_REFERENCE.md) for method details
- [Troubleshooting](./TROUBLESHOOTING.md) for gas-related issues













