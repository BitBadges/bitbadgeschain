# Best Practices

Security, gas optimization, and design best practices for building contracts with the BitBadges tokenization precompile.

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
function transfer(
    uint256 collectionId,
    address to,
    uint256 amount,
    uint256 tokenId
) external {
    // Validate inputs
    TokenizationErrors.requireValidCollectionId(collectionId);
    TokenizationErrors.requireValidAddress(to);
    require(amount > 0, "Amount must be > 0");
    require(tokenId > 0, "Token ID must be > 0");
    
    // Perform transfer
    TokenizationWrappers.transferSingleToken(
        TOKENIZATION, collectionId, to, amount, tokenId
    );
}
```

### 2. Access Control

Implement proper access control for sensitive operations:

```solidity
contract SecureToken {
    address public owner;
    mapping(address => bool) public authorized;
    
    modifier onlyOwner() {
        require(msg.sender == owner, "Not owner");
        _;
    }
    
    modifier onlyAuthorized() {
        require(authorized[msg.sender] || msg.sender == owner, "Not authorized");
        _;
    }
    
    function setKYCStatus(address user, bool status) external onlyAuthorized {
        // Only authorized addresses can set KYC
    }
}
```

### 3. Reentrancy Protection

Use reentrancy guards for state-changing operations:

```solidity
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";

contract SecureContract is ReentrancyGuard {
    function transfer(...) external nonReentrant {
        // Transfer logic
    }
}
```

### 4. Check-Effects-Interactions Pattern

Follow the check-effects-interactions pattern:

```solidity
function transfer(...) external {
    // 1. Checks
    require(balance >= amount, "Insufficient balance");
    
    // 2. Effects (state changes)
    balances[msg.sender] -= amount;
    balances[to] += amount;
    
    // 3. Interactions (external calls)
    TokenizationWrappers.transferTokens(...);
}
```

### 5. Address Validation

Always validate addresses:

```solidity
function setManager(address newManager) external {
    require(newManager != address(0), "Invalid address");
    require(newManager != address(this), "Cannot set self");
    // ...
}
```

### 6. Integer Overflow Protection

Solidity 0.8+ has built-in overflow protection, but be aware of edge cases:

```solidity
function safeTransfer(uint256 amount) external {
    require(amount > 0, "Amount must be > 0");
    require(amount <= type(uint256).max, "Amount too large");
    // ...
}
```

## Gas Optimization

### 1. Use Typed Wrappers

Typed wrappers are generally more gas-efficient than manual JSON construction:

```solidity
// ✅ Efficient
TokenizationWrappers.transferTokens(...);

// ❌ Less efficient (more string operations)
string memory json = TokenizationJSONHelpers.transferTokensJSON(...);
TOKENIZATION.transferTokens(json);
```

### 2. Cache Frequently Used Values

Cache registry IDs and other frequently accessed values:

```solidity
contract OptimizedContract {
    uint256 public kycRegistryId;  // Cached
    
    constructor() {
        kycRegistryId = TokenizationWrappers.createDynamicStore(...);
    }
    
    function checkKYC(address user) external view {
        // Use cached ID instead of looking it up
        return TokenizationWrappers.getDynamicStoreValue(
            TOKENIZATION, kycRegistryId, user
        );
    }
}
```

### 3. Use Direct Queries When Possible

Prefer direct queries that return simple types:

```solidity
// ✅ Gas efficient - returns uint256 directly
uint256 balance = TokenizationWrappers.getBalanceAmount(...);

// ❌ Less efficient - returns protobuf bytes
bytes memory balanceBytes = TokenizationWrappers.getBalance(...);
// Requires off-chain decoding
```

### 4. Batch Operations Efficiently

Batch operations to reduce transaction overhead:

```solidity
function batchTransfer(
    address[] memory recipients,
    uint256[] memory amounts
) external {
    // Single transaction for multiple transfers
    for (uint256 i = 0; i < recipients.length; i++) {
        TokenizationWrappers.transferSingleToken(...);
    }
}
```

### 5. Avoid Unnecessary Storage Reads

Cache storage values in memory:

```solidity
function transfer(...) external {
    // ✅ Cache in memory
    uint256 cachedCollectionId = collectionId;
    
    // Use cached value
    TokenizationWrappers.transferTokens(
        TOKENIZATION, cachedCollectionId, ...
    );
}
```

### 6. Use Events Instead of Storage

Emit events for data that can be indexed off-chain:

```solidity
event TransferExecuted(
    uint256 indexed collectionId,
    address indexed from,
    address indexed to,
    uint256 amount
);

function transfer(...) external {
    // Emit event instead of storing in mapping
    emit TransferExecuted(collectionId, from, to, amount);
}
```

## Code Organization

### 1. Use Libraries for Reusable Logic

Organize reusable logic into libraries:

```solidity
library TransferHelpers {
    function safeTransfer(
        ITokenizationPrecompile precompile,
        uint256 collectionId,
        address to,
        uint256 amount,
        uint256 tokenId
    ) internal returns (bool) {
        TokenizationErrors.requireValidCollectionId(collectionId);
        TokenizationErrors.requireValidAddress(to);
        require(amount > 0, "Amount must be > 0");
        
        return TokenizationWrappers.transferSingleToken(
            precompile, collectionId, to, amount, tokenId
        );
    }
}
```

### 2. Separate Concerns

Separate business logic from precompile interactions:

```solidity
contract WellOrganizedContract {
    // Business logic
    function processTransfer(...) internal {
        // Validation
        // Business rules
        // State updates
    }
    
    // Precompile interaction
    function executeTransfer(...) internal {
        TokenizationWrappers.transferTokens(...);
    }
    
    // Public interface
    function transfer(...) external {
        processTransfer(...);
        executeTransfer(...);
    }
}
```

### 3. Use Constants

Define constants for magic numbers and addresses:

```solidity
contract ConstantsContract {
    ITokenizationPrecompile constant TOKENIZATION = 
        ITokenizationPrecompile(0x0000000000000000000000000000000000001001);
    
    uint256 public constant MIN_TRANSFER_AMOUNT = 1;
    uint256 public constant MAX_TRANSFER_AMOUNT = type(uint256).max;
}
```

## Error Handling

### 1. Use Custom Errors

Custom errors are more gas-efficient than string errors:

```solidity
// ✅ Gas efficient
error InsufficientBalance(uint256 required, uint256 available);

function transfer(uint256 amount) external {
    if (balance < amount) {
        revert InsufficientBalance(amount, balance);
    }
}

// ❌ Less efficient
require(balance >= amount, "Insufficient balance");
```

### 2. Provide Context in Errors

Include relevant context in error messages:

```solidity
error TransferFailed(
    uint256 collectionId,
    address from,
    address to,
    string reason
);

function transfer(...) external {
    bool success = TokenizationWrappers.transferTokens(...);
    if (!success) {
        revert TransferFailed(collectionId, from, to, "Precompile call failed");
    }
}
```

### 3. Validate Before External Calls

Always validate before making external calls:

```solidity
function transfer(...) external {
    // Validate first
    require(collectionId > 0, "Invalid collection");
    require(to != address(0), "Invalid recipient");
    
    // Then call
    TokenizationWrappers.transferTokens(...);
}
```

## Testing Strategies

### 1. Use Test Helpers

Leverage test helpers for consistent test data:

```solidity
import "./libraries/TokenizationTestHelpers.sol";

function testTransfer() public {
    uint256 collectionId = 1;
    address recipient = TokenizationTestHelpers.generateTestAddresses(1)[0];
    
    TokenizationTypes.UintRange[] memory tokenIds = 
        TokenizationTestHelpers.generateTokenIdRanges(1, 10, 1);
    
    // Test transfer
}
```

### 2. Test Edge Cases

Test boundary conditions and edge cases:

```solidity
function testTransferEdgeCases() public {
    // Test zero amount
    // Test zero address
    // Test invalid collection ID
    // Test expired ownership times
    // Test invalid token IDs
}
```

### 3. Test Error Conditions

Verify errors are thrown correctly:

```solidity
function testTransferFailsWithInsufficientBalance() public {
    uint256 largeAmount = type(uint256).max;
    
    vm.expectRevert("Insufficient balance");
    transfer(collectionId, recipient, largeAmount, tokenId);
}
```

### 4. Integration Tests

Test with actual precompile interactions:

```solidity
function testIntegrationTransfer() public {
    // Setup
    uint256 collectionId = createTestCollection();
    mintTokens(collectionId, user, 100);
    
    // Execute
    transfer(collectionId, recipient, 50, 1);
    
    // Verify
    uint256 balance = getBalance(collectionId, recipient);
    assertEq(balance, 50);
}
```

## Additional Recommendations

1. **Document Your Code**: Use NatSpec comments for all public functions
2. **Version Your Contracts**: Use semantic versioning
3. **Upgradeability**: Consider upgradeable patterns if needed
4. **Pausability**: Implement pause functionality for emergency stops
5. **Time Locks**: Use time locks for critical operations
6. **Multi-sig**: Consider multi-sig for administrative functions
7. **Audit**: Get professional security audits before mainnet deployment

## See Also

- [Gas Optimization Guide](./GAS_OPTIMIZATION.md) for detailed gas strategies
- [Troubleshooting](./TROUBLESHOOTING.md) for common issues
- [Patterns](./PATTERNS.md) for implementation patterns













