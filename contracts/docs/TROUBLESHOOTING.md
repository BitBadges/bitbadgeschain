# Troubleshooting

Common issues and solutions when working with the BitBadges tokenization precompile.

## Table of Contents

- [Common Errors](#common-errors)
- [JSON Construction Issues](#json-construction-issues)
- [Type Conversion Problems](#type-conversion-problems)
- [Gas Issues](#gas-issues)
- [Query Response Decoding](#query-response-decoding)
- [Permission and Approval Issues](#permission-and-approval-issues)

## Common Errors

### "Collection not found"

**Error:** `CollectionNotFound` or similar error when querying/transferring.

**Solutions:**
1. Verify the collection ID is correct
2. Check that the collection exists using `getCollection()`
3. Ensure the collection is not archived
4. Verify you're using the correct network/chain

```solidity
// Always validate collection ID
TokenizationErrors.requireValidCollectionId(collectionId);

// Check collection exists
bytes memory collection = TokenizationWrappers.getCollection(
    TOKENIZATION, collectionId
);
require(collection.length > 0, "Collection does not exist");
```

### "Insufficient balance"

**Error:** Transfer fails due to insufficient balance.

**Solutions:**
1. Check actual balance before transferring:
```solidity
uint256 balance = TokenizationWrappers.getBalanceAmount(
    TOKENIZATION,
    collectionId,
    msg.sender,
    tokenIds,
    ownershipTimes
);
require(balance >= amount, "Insufficient balance");
```

2. Verify token ID ranges match your balance
3. Check ownership time ranges - tokens may have expired
4. Ensure you're querying with the correct ranges

### "Transfer failed"

**Error:** Generic transfer failure.

**Solutions:**
1. Check permissions - collection or user approvals may be required
2. Verify token IDs are valid for the collection
3. Check ownership time ranges are valid
4. Ensure recipient addresses are valid (not zero address)
5. Review collection permissions and approval requirements

### "Invalid token ID"

**Error:** Token ID is not valid for the collection.

**Solutions:**
1. Query collection to see valid token ID ranges:
```solidity
bytes memory collection = TokenizationWrappers.getCollection(
    TOKENIZATION, collectionId
);
// Decode to see validTokenIds (requires off-chain decoding)
```

2. Use token IDs within the collection's `validTokenIds` ranges
3. Verify token ID ranges in your transfer match collection configuration

## JSON Construction Issues

### Malformed JSON

**Error:** JSON string is invalid or malformed.

**Solutions:**
1. Use helper libraries instead of manual construction:
```solidity
// ❌ Don't do this:
string memory json = string(abi.encodePacked('{"collectionId":"', id, '"}'));

// ✅ Do this:
string memory json = TokenizationJSONHelpers.getCollectionJSON(collectionId);
```

2. Use typed wrappers when possible:
```solidity
// ✅ Best approach:
TokenizationWrappers.transferTokens(
    TOKENIZATION, collectionId, recipients, amount, tokenIds, ownershipTimes
);
```

3. Check JSON escaping for strings with special characters
4. Verify all required fields are present

### Missing Required Fields

**Error:** JSON missing required fields.

**Solutions:**
1. Review the protobuf message structure
2. Use builders for complex operations:
```solidity
TokenizationBuilders.CollectionBuilder memory builder = 
    TokenizationBuilders.newCollection();
builder = builder.withValidTokenIdRange(1, 1000);
// Builder ensures required fields are set
string memory json = builder.build();
```

3. Check documentation for required vs optional fields

## Type Conversion Problems

### Address Format Issues

**Error:** Address conversion between EVM and Cosmos formats.

**Solutions:**
1. Use helper function for Cosmos address strings:
```solidity
string memory cosmosAddr = TokenizationHelpers.addressToCosmosString(evmAddress);
```

2. Manager fields require Cosmos address format (hex string with 0x prefix)
3. User addresses in queries can use EVM address format

### UintRange Array Mismatch

**Error:** Array length mismatches when creating ranges.

**Solutions:**
1. Use helper functions:
```solidity
// ✅ Correct:
TokenizationTypes.UintRange[] memory ranges = 
    TokenizationHelpers.createUintRangeArray(starts, ends);

// ❌ Wrong:
// Manually creating arrays with mismatched lengths
```

2. Always validate array lengths match:
```solidity
require(starts.length == ends.length, "Array length mismatch");
```

## Gas Issues

### Out of Gas

**Error:** Transaction runs out of gas.

**Solutions:**
1. Break large operations into smaller batches
2. Use direct queries (`getBalanceAmount`) instead of full collection queries when possible
3. Optimize JSON construction - use typed wrappers which are more gas-efficient
4. Consider using events for off-chain indexing instead of on-chain queries

### High Gas Costs

**Issue:** Operations are expensive.

**Solutions:**
1. Use `getBalanceAmount()` instead of `getBalance()` when you only need the amount
2. Batch operations efficiently
3. Cache frequently accessed data (like registry IDs)
4. See [Gas Optimization](./GAS_OPTIMIZATION.md) for detailed strategies

## Query Response Decoding

### Protobuf Decoding Not Working

**Error:** Cannot decode query responses in Solidity.

**Solutions:**
1. **Use direct queries when possible:**
```solidity
// ✅ Returns uint256 directly:
uint256 balance = TokenizationWrappers.getBalanceAmount(...);

// ❌ Returns protobuf bytes (requires off-chain decoding):
bytes memory balance = TokenizationWrappers.getBalance(...);
```

2. **Decode off-chain:**
   - Use the TypeScript SDK for full decoding
   - Emit events with raw bytes for off-chain processing
   - Use view functions that return decoded values

3. **Extract simple fields:**
   - For boolean dynamic stores, you may be able to check bytes length
   - For complex types, always decode off-chain

### Empty Response Bytes

**Error:** Query returns empty bytes.

**Solutions:**
1. Verify the query parameters are correct
2. Check that the resource exists (collection, balance, etc.)
3. Ensure you're using the correct query method
4. Check network/chain configuration

## Permission and Approval Issues

### "Permission denied"

**Error:** Operation fails due to insufficient permissions.

**Solutions:**
1. Check collection permissions:
   - Verify you're the creator/manager if required
   - Review `CollectionPermissions` structure
   - Check time-based permission restrictions

2. Check user approvals:
   - Verify incoming/outgoing approvals are set correctly
   - Check approval criteria are met
   - Review approval expiration times

3. Use collection-level approvals for complex permission scenarios

### Approval Not Working

**Error:** Approval set but transfers still fail.

**Solutions:**
1. Verify approval ID matches exactly (case-sensitive)
2. Check approval criteria:
   - Token ID ranges
   - Ownership time ranges
   - Transfer time ranges
   - Address list requirements

3. Ensure approval is not expired
4. Check approval level (collection vs user)
5. Verify approval version matches

## Debugging Tips

### 1. Use Events

Emit events to track state changes:

```solidity
event DebugTransfer(
    uint256 collectionId,
    address from,
    address to,
    uint256 amount,
    bool success
);

function transfer(...) external {
    bool success = TokenizationWrappers.transferTokens(...);
    emit DebugTransfer(collectionId, from, to, amount, success);
}
```

### 2. Validate Before Calling

Always validate inputs before calling precompile:

```solidity
function safeTransfer(...) external {
    TokenizationErrors.requireValidCollectionId(collectionId);
    TokenizationErrors.requireValidAddress(to);
    require(amount > 0, "Amount must be > 0");
    
    // Then perform transfer
}
```

### 3. Check Return Values

Always check return values:

```solidity
bool success = TokenizationWrappers.transferTokens(...);
require(success, "Transfer failed");
```

### 4. Use Try-Catch for Complex Operations

For operations that might fail, use try-catch:

```solidity
try TOKENIZATION.createCollection(json) returns (uint256 collectionId) {
    // Success
} catch Error(string memory reason) {
    // Handle error
    revert(reason);
} catch {
    revert("Unknown error");
}
```

## Getting More Help

1. Check the [API Reference](./API_REFERENCE.md) for method signatures
2. Review [Examples](./EXAMPLES.md) for working code
3. See [Best Practices](./BEST_PRACTICES.md) for common patterns
4. Review example contracts in `contracts/examples/`
5. Check precompile implementation in `x/tokenization/precompile/`













