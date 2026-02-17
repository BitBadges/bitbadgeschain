# Getting Started with BitBadges Tokenization Precompile

This guide will help you get started building Solidity contracts that interact with the BitBadges tokenization precompile.

## Overview

The BitBadges tokenization precompile provides a Solidity interface to the tokenization module, enabling you to:
- Create and manage token collections
- Transfer tokens with time-bound ownership
- Set up approvals and permissions
- Use dynamic stores for boolean registries (e.g., KYC/compliance)
- Query balances, collections, and more

## Precompile Address

```
0x0000000000000000000000000000000000001001
```

## Installation

### 1. Import the Required Files

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./interfaces/ITokenizationPrecompile.sol";
import "./types/TokenizationTypes.sol";
import "./libraries/TokenizationWrappers.sol";
import "./libraries/TokenizationHelpers.sol";
```

### 2. Initialize the Precompile Interface

```solidity
contract MyTokenContract {
    ITokenizationPrecompile constant TOKENIZATION = 
        ITokenizationPrecompile(0x0000000000000000000000000000000000001001);
    
    // Your contract code here
}
```

## Quick Start Examples

### Example 1: Simple Token Transfer

```solidity
import "./libraries/TokenizationWrappers.sol";
import "./libraries/TokenizationHelpers.sol";

function transferToken(
    uint256 collectionId,
    address to,
    uint256 amount,
    uint256 tokenId
) external {
    // Use typed wrapper for type safety
    bool success = TokenizationWrappers.transferSingleToken(
        TOKENIZATION,
        collectionId,
        to,
        amount,
        tokenId
    );
    require(success, "Transfer failed");
}
```

### Example 2: Query Balance

```solidity
function getBalance(
    uint256 collectionId,
    address user
) external view returns (uint256) {
    // Create token ID and ownership time ranges
    TokenizationTypes.UintRange[] memory tokenIds = new TokenizationTypes.UintRange[](1);
    tokenIds[0] = TokenizationHelpers.createSingleTokenIdRange(1);
    
    TokenizationTypes.UintRange[] memory ownershipTimes = new TokenizationTypes.UintRange[](1);
    ownershipTimes[0] = TokenizationHelpers.createFullOwnershipTimeRange();
    
    // Query balance amount
    return TokenizationWrappers.getBalanceAmount(
        TOKENIZATION,
        collectionId,
        user,
        tokenIds,
        ownershipTimes
    );
}
```

### Example 3: Create a Dynamic Store (KYC Registry)

```solidity
function createKYCRegistry() external returns (uint256) {
    // Create a dynamic store with default value false (not KYC'd)
    uint256 storeId = TokenizationWrappers.createDynamicStore(
        TOKENIZATION,
        false,  // defaultValue: addresses are not KYC'd by default
        "ipfs://...",  // URI for metadata
        "KYC Registry"  // customData
    );
    return storeId;
}

function setKYCStatus(address user, bool isKYCd) external {
    require(kycRegistryId != 0, "Registry not initialized");
    TokenizationWrappers.setDynamicStoreValue(
        TOKENIZATION,
        kycRegistryId,
        user,
        isKYCd
    );
}
```

## Library Overview

### TokenizationWrappers
Type-safe wrappers that accept structs instead of JSON strings. Use these for better compile-time checking.

```solidity
// Instead of constructing JSON manually:
TokenizationWrappers.transferTokens(
    TOKENIZATION,
    collectionId,
    recipients,
    amount,
    tokenIds,
    ownershipTimes
);
```

### TokenizationHelpers
Utilities for creating and validating tokenization types.

```solidity
// Create common ranges
TokenizationTypes.UintRange memory fullOwnership = 
    TokenizationHelpers.createFullOwnershipTimeRange();

TokenizationTypes.UintRange memory singleToken = 
    TokenizationHelpers.createSingleTokenIdRange(1);
```

### TokenizationBuilders
Fluent builder APIs for complex operations.

```solidity
// Build a collection creation request
TokenizationBuilders.CollectionBuilder memory builder = 
    TokenizationBuilders.newCollection();
builder = builder.withValidTokenIdRange(1, 1000);
builder = builder.withManager(managerAddress);
string memory json = builder.build();
uint256 collectionId = TOKENIZATION.createCollection(json);
```

### TokenizationJSONHelpers
Low-level JSON construction helpers. Use when you need fine-grained control.

```solidity
string memory json = TokenizationJSONHelpers.transferTokensJSON(
    collectionId,
    recipients,
    amount,
    tokenIdsJson,
    ownershipTimesJson
);
```

## Common Patterns

### Pattern 1: Transfer with Full Ownership

```solidity
function transferWithFullOwnership(
    uint256 collectionId,
    address to,
    uint256 amount,
    uint256 tokenId
) external {
    TokenizationTypes.UintRange[] memory tokenIds = new TokenizationTypes.UintRange[](1);
    tokenIds[0] = TokenizationHelpers.createSingleTokenIdRange(tokenId);
    
    address[] memory recipients = new address[](1);
    recipients[0] = to;
    
    TokenizationWrappers.transferTokensWithFullOwnership(
        TOKENIZATION,
        collectionId,
        recipients,
        amount,
        tokenIds
    );
}
```

### Pattern 2: Time-Bound Ownership

```solidity
function transferWithExpiration(
    uint256 collectionId,
    address to,
    uint256 amount,
    uint256 tokenId,
    uint256 expirationTime
) external {
    TokenizationTypes.UintRange[] memory tokenIds = new TokenizationTypes.UintRange[](1);
    tokenIds[0] = TokenizationHelpers.createSingleTokenIdRange(tokenId);
    
    TokenizationTypes.UintRange[] memory ownershipTimes = new TokenizationTypes.UintRange[](1);
    ownershipTimes[0] = TokenizationHelpers.createOwnershipTimeRangeToExpiration(expirationTime);
    
    address[] memory recipients = new address[](1);
    recipients[0] = to;
    
    TokenizationWrappers.transferTokens(
        TOKENIZATION,
        collectionId,
        recipients,
        amount,
        tokenIds,
        ownershipTimes
    );
}
```

### Pattern 3: Compliance Check Before Transfer

```solidity
function transferWithCompliance(
    uint256 collectionId,
    address to,
    uint256 amount,
    uint256 tokenId
) external {
    // Check KYC status
    bytes memory storeValue = TokenizationWrappers.getDynamicStoreValue(
        TOKENIZATION,
        kycRegistryId,
        to
    );
    // Note: Full decoding requires off-chain tools
    // For boolean stores, you may need to check the raw bytes or use events
    
    // Perform transfer
    TokenizationWrappers.transferSingleToken(
        TOKENIZATION,
        collectionId,
        to,
        amount,
        tokenId
    );
}
```

## Next Steps

- Read the [API Reference](./API_REFERENCE.md) for complete method documentation
- Explore [Common Patterns](./PATTERNS.md) for advanced use cases
- Check out [Examples](./EXAMPLES.md) for real-world contract implementations
- Review [Best Practices](./BEST_PRACTICES.md) for security and gas optimization

## Getting Help

- Check [Troubleshooting](./TROUBLESHOOTING.md) for common issues
- Review the example contracts in `contracts/examples/`
- See the main README in `contracts/examples/README.md`















