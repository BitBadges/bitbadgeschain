# Gas & Costs

The tokenization precompile uses dynamic gas calculation based on operation complexity.

## Gas Calculation

Gas costs are calculated as:

```
Total Gas = Base Gas + (Per-Element Costs)
```

Where per-element costs depend on:
- Number of recipients
- Number of token ID ranges
- Number of ownership time ranges
- Number of approval fields
- Query range complexity

## Transaction Gas Costs

### Base Gas Costs

| Method | Base Gas |
|--------|----------|
| `transferTokens` | 30,000 |
| `setIncomingApproval` | 20,000 |
| `setOutgoingApproval` | 20,000 |
| `createCollection` | 50,000 |
| `updateCollection` | 40,000 |
| `deleteCollection` | 20,000 |
| `createAddressLists` | 30,000 |
| `updateUserApprovals` | 30,000 |
| `deleteIncomingApproval` | 15,000 |
| `deleteOutgoingApproval` | 15,000 |
| `purgeApprovals` | 25,000 |
| `createDynamicStore` | 20,000 |
| `updateDynamicStore` | 20,000 |
| `deleteDynamicStore` | 15,000 |
| `setDynamicStoreValue` | 15,000 |
| `setValidTokenIds` | 20,000 |
| `setManager` | 15,000 |
| `setCollectionMetadata` | 15,000 |
| `setTokenMetadata` | 20,000 |
| `setCustomData` | 15,000 |
| `setStandards` | 15,000 |
| `setCollectionApprovals` | 30,000 |
| `setIsArchived` | 15,000 |
| `castVote` | 15,000 |
| `universalUpdateCollection` | 50,000 |

### Per-Element Costs

| Element | Cost per Element |
|---------|------------------|
| Recipient | 5,000 |
| Token ID Range | 1,000 |
| Ownership Time Range | 1,000 |
| Approval Field | 500 |

### Example: Transfer Tokens

```solidity
// Transfer to 3 recipients with 2 token ID ranges
// Gas = 30,000 (base) + (3 * 5,000) + (2 * 1,000) = 47,000
address[] memory recipients = new address[](3);
TokenizationTypes.UintRange[] memory tokenIds = new TokenizationTypes.UintRange[](2);
// ... setup

precompile.transferTokens(collectionId, recipients, amount, tokenIds, ownershipTimes);
```

## Query Gas Costs

### Base Gas Costs

| Method | Base Gas |
|--------|----------|
| `getCollection` | 3,000 |
| `getBalance` | 3,000 |
| `getAddressList` | 5,000 |
| `getApprovalTracker` | 5,000 |
| `getChallengeTracker` | 5,000 |
| `getETHSignatureTracker` | 5,000 |
| `getDynamicStore` | 5,000 |
| `getDynamicStoreValue` | 5,000 |
| `getWrappableBalances` | 5,000 |
| `isAddressReservedProtocol` | 2,000 |
| `getAllReservedProtocolAddresses` | 5,000 |
| `getVote` | 5,000 |
| `getVotes` | 5,000 |
| `params` | 2,000 |
| `getBalanceAmount` | 3,000 |
| `getTotalSupply` | 3,000 |

### Per-Range Costs

| Element | Cost per Range |
|---------|----------------|
| Query Range | 500 |

### Example: Get Balance Amount

```solidity
// Query with 3 token ID ranges and 2 ownership time ranges
// Gas = 3,000 (base) + (3 * 500) + (2 * 500) = 5,500
TokenizationTypes.UintRange[] memory tokenIds = new TokenizationTypes.UintRange[](3);
TokenizationTypes.UintRange[] memory ownershipTimes = new TokenizationTypes.UintRange[](2);
// ... setup

bytes memory result = precompile.getBalanceAmount(collectionId, user, tokenIds, ownershipTimes);
```

## Gas Estimation

### Estimating Transaction Gas

```solidity
function estimateTransferGas(
    uint256 numRecipients,
    uint256 numTokenRanges,
    uint256 numOwnershipRanges
) public pure returns (uint256) {
    return 30_000 + // base
           (numRecipients * 5_000) +
           (numTokenRanges * 1_000) +
           (numOwnershipRanges * 1_000);
}
```

### Estimating Query Gas

```solidity
function estimateQueryGas(
    uint256 numTokenRanges,
    uint256 numOwnershipRanges
) public pure returns (uint256) {
    return 3_000 + // base (example: getBalanceAmount)
           (numTokenRanges * 500) +
           (numOwnershipRanges * 500);
}
```

## Gas Optimization Tips

### 1. Batch Operations

Instead of multiple transfers, use a single transfer with multiple recipients:

```solidity
// ❌ Inefficient: 3 separate transfers
precompile.transferTokens(collectionId, [recipient1], amount, tokenIds, ownershipTimes);
precompile.transferTokens(collectionId, [recipient2], amount, tokenIds, ownershipTimes);
precompile.transferTokens(collectionId, [recipient3], amount, tokenIds, ownershipTimes);

// ✅ Efficient: 1 transfer with 3 recipients
address[] memory recipients = new address[](3);
recipients[0] = recipient1;
recipients[1] = recipient2;
recipients[2] = recipient3;
precompile.transferTokens(collectionId, recipients, amount, tokenIds, ownershipTimes);
```

### 2. Minimize Ranges

Use fewer, larger ranges instead of many small ranges:

```solidity
// ❌ Inefficient: 10 small ranges
TokenizationTypes.UintRange[] memory tokenIds = new TokenizationTypes.UintRange[](10);
// ... 10 ranges

// ✅ Efficient: 1 large range
TokenizationTypes.UintRange[] memory tokenIds = new TokenizationTypes.UintRange[](1);
tokenIds[0] = TokenizationTypes.UintRange({start: 1, end: 1000});
```

### 3. Cache Query Results

Cache query results in your contract to avoid repeated queries:

```solidity
mapping(uint256 => bytes) private cachedBalances;

function getCachedBalance(uint256 collectionId, address user) external returns (bytes memory) {
    if (cachedBalances[collectionId] != bytes(0)) {
        return cachedBalances[collectionId];
    }
    bytes memory balance = precompile.getBalance(collectionId, user);
    cachedBalances[collectionId] = balance;
    return balance;
}
```

## Gas Limits

The EVM has a gas limit per transaction. For complex operations:

- **Maximum Recipients**: ~200 (with base gas)
- **Maximum Ranges**: ~1,000 (with base gas)
- **Complex Approvals**: Depends on approval criteria complexity

Always test with realistic data sizes to ensure operations complete within gas limits.

## Resources

- [Transaction Methods](transactions.md) - Method-specific gas information
- [Query Methods](queries.md) - Query-specific gas information
- [Cosmos SDK EVM Documentation](https://docs.cosmos.network/evm/v0.5.0/documentation/overview) - EVM gas model









