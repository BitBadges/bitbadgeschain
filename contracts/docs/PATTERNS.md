# Common Patterns

Common patterns and use cases for building contracts with the BitBadges tokenization precompile.

## Table of Contents

- [Basic Transfer Patterns](#basic-transfer-patterns)
- [Compliance Patterns](#compliance-patterns)
- [Time-Bound Ownership](#time-bound-ownership)
- [Multi-Token Collections](#multi-token-collections)
- [Approval Patterns](#approval-patterns)
- [Dynamic Store Patterns](#dynamic-store-patterns)
- [Collection Management](#collection-management)

## Basic Transfer Patterns

### Simple Transfer

```solidity
function transfer(
    uint256 collectionId,
    address to,
    uint256 amount,
    uint256 tokenId
) external {
    TokenizationWrappers.transferSingleToken(
        TOKENIZATION,
        collectionId,
        to,
        amount,
        tokenId
    );
}
```

### Batch Transfer

```solidity
function batchTransfer(
    uint256 collectionId,
    address[] memory recipients,
    uint256[] memory amounts,
    uint256 tokenId
) external {
    require(recipients.length == amounts.length, "Length mismatch");
    
    TokenizationTypes.UintRange[] memory tokenIds = new TokenizationTypes.UintRange[](1);
    tokenIds[0] = TokenizationHelpers.createSingleTokenIdRange(tokenId);
    
    TokenizationTypes.UintRange[] memory ownershipTimes = new TokenizationTypes.UintRange[](1);
    ownershipTimes[0] = TokenizationHelpers.createFullOwnershipTimeRange();
    
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

## Compliance Patterns

### KYC Registry

```solidity
contract CompliantToken {
    ITokenizationPrecompile constant TOKENIZATION = 
        ITokenizationPrecompile(0x0000000000000000000000000000000000001001);
    
    uint256 public kycRegistryId;
    address public complianceOfficer;
    
    constructor() {
        complianceOfficer = msg.sender;
        // Create KYC registry (default: not KYC'd)
        kycRegistryId = TokenizationWrappers.createDynamicStore(
            TOKENIZATION,
            false,
            "ipfs://kyc-registry-metadata",
            "KYC Registry"
        );
    }
    
    modifier onlyComplianceOfficer() {
        require(msg.sender == complianceOfficer, "Not authorized");
        _;
    }
    
    function setKYCStatus(address user, bool isKYCd) external onlyComplianceOfficer {
        TokenizationWrappers.setDynamicStoreValue(
            TOKENIZATION,
            kycRegistryId,
            user,
            isKYCd
        );
    }
    
    function transferWithKYC(
        uint256 collectionId,
        address to,
        uint256 amount,
        uint256 tokenId
    ) external {
        // Check KYC status (simplified - full check may require off-chain)
        // In production, you might emit an event and check off-chain
        // or use a more sophisticated on-chain check
        
        TokenizationWrappers.transferSingleToken(
            TOKENIZATION,
            collectionId,
            to,
            amount,
            tokenId
        );
    }
}
```

### Whitelist Pattern

```solidity
contract WhitelistedToken {
    uint256 public whitelistId;
    
    function createWhitelist() external {
        // Create address list for whitelist
        // Note: Address lists are created via createAddressLists
        // This is a simplified example
    }
    
    function transferToWhitelisted(
        uint256 collectionId,
        address to,
        uint256 amount,
        uint256 tokenId
    ) external {
        // Check whitelist (implementation depends on your setup)
        // Perform transfer
        TokenizationWrappers.transferSingleToken(
            TOKENIZATION,
            collectionId,
            to,
            amount,
            tokenId
        );
    }
}
```

## Time-Bound Ownership

### Lock-Up Period

```solidity
function transferWithLockup(
    uint256 collectionId,
    address to,
    uint256 amount,
    uint256 tokenId,
    uint256 lockupDuration
) external {
    TokenizationTypes.UintRange[] memory tokenIds = new TokenizationTypes.UintRange[](1);
    tokenIds[0] = TokenizationHelpers.createSingleTokenIdRange(tokenId);
    
    // Ownership expires after lockup period
    TokenizationTypes.UintRange[] memory ownershipTimes = new TokenizationTypes.UintRange[](1);
    ownershipTimes[0] = TokenizationHelpers.createOwnershipTimeRange(
        block.timestamp,
        lockupDuration
    );
    
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

### Expiring Tokens

```solidity
function issueExpiringToken(
    uint256 collectionId,
    address to,
    uint256 amount,
    uint256 tokenId,
    uint256 expirationTime
) external {
    TokenizationTypes.UintRange[] memory tokenIds = new TokenizationTypes.UintRange[](1);
    tokenIds[0] = TokenizationHelpers.createSingleTokenIdRange(tokenId);
    
    // Token expires at specific time
    TokenizationTypes.UintRange[] memory ownershipTimes = new TokenizationTypes.UintRange[](1);
    ownershipTimes[0] = TokenizationHelpers.createUintRange(
        block.timestamp,
        expirationTime
    );
    
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

## Multi-Token Collections

### Token ID Ranges

```solidity
function transferMultipleTokenRanges(
    uint256 collectionId,
    address to,
    uint256 amount
) external {
    // Transfer tokens across multiple ID ranges
    TokenizationTypes.UintRange[] memory tokenIds = new TokenizationTypes.UintRange[](3);
    tokenIds[0] = TokenizationHelpers.createUintRange(1, 100);
    tokenIds[1] = TokenizationHelpers.createUintRange(200, 300);
    tokenIds[2] = TokenizationHelpers.createUintRange(500, 600);
    
    TokenizationTypes.UintRange[] memory ownershipTimes = new TokenizationTypes.UintRange[](1);
    ownershipTimes[0] = TokenizationHelpers.createFullOwnershipTimeRange();
    
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

## Approval Patterns

### Setting Approvals

```solidity
function setApproval(
    uint256 collectionId,
    string memory approvalId,
    // ... approval parameters
) external {
    // Construct approval JSON using TokenizationJSONHelpers
    // or use builder pattern for complex approvals
    string memory json = ""; // Construct approval JSON
    TOKENIZATION.setIncomingApproval(json);
}
```

### Checking Approvals Before Transfer

```solidity
function transferWithApprovalCheck(
    uint256 collectionId,
    address to,
    uint256 amount,
    uint256 tokenId
) external {
    // In production, you might:
    // 1. Query approval status off-chain
    // 2. Use collection-level approvals
    // 3. Check user balance and approvals
    
    // Perform transfer (approvals are checked by the precompile)
    TokenizationWrappers.transferSingleToken(
        TOKENIZATION,
        collectionId,
        to,
        amount,
        tokenId
    );
}
```

## Dynamic Store Patterns

### Boolean Registry

```solidity
contract RegistryManager {
    mapping(string => uint256) public registries;
    
    function createRegistry(string memory name) external returns (uint256) {
        uint256 storeId = TokenizationWrappers.createDynamicStore(
            TOKENIZATION,
            false,  // Default value
            string(abi.encodePacked("ipfs://", name)),
            name
        );
        registries[name] = storeId;
        return storeId;
    }
    
    function setRegistryValue(
        string memory registryName,
        address user,
        bool value
    ) external {
        uint256 storeId = registries[registryName];
        require(storeId != 0, "Registry not found");
        
        TokenizationWrappers.setDynamicStoreValue(
            TOKENIZATION,
            storeId,
            user,
            value
        );
    }
    
    function getRegistryValue(
        string memory registryName,
        address user
    ) external view returns (bytes memory) {
        uint256 storeId = registries[registryName];
        require(storeId != 0, "Registry not found");
        
        return TokenizationWrappers.getDynamicStoreValue(
            TOKENIZATION,
            storeId,
            user
        );
    }
}
```

### Multi-Registry Pattern

```solidity
contract MultiComplianceToken {
    uint256 public kycRegistryId;
    uint256 public accreditedInvestorRegistryId;
    uint256 public blacklistRegistryId;
    
    function initializeRegistries() external {
        kycRegistryId = TokenizationWrappers.createDynamicStore(
            TOKENIZATION, false, "ipfs://kyc", "KYC"
        );
        accreditedInvestorRegistryId = TokenizationWrappers.createDynamicStore(
            TOKENIZATION, false, "ipfs://accredited", "Accredited"
        );
        blacklistRegistryId = TokenizationWrappers.createDynamicStore(
            TOKENIZATION, false, "ipfs://blacklist", "Blacklist"
        );
    }
    
    function checkCompliance(address user) external view returns (bool) {
        // Check all registries
        // Note: Full decoding may require off-chain tools
        bytes memory kyc = TokenizationWrappers.getDynamicStoreValue(
            TOKENIZATION, kycRegistryId, user
        );
        // ... check other registries
        return true; // Simplified
    }
}
```

## Collection Management

### Creating Collections with Builder

```solidity
function createMyCollection(
    string memory name,
    string memory symbol
) external returns (uint256) {
    TokenizationBuilders.CollectionBuilder memory builder = 
        TokenizationBuilders.newCollection();
    
    // Set required fields
    builder = builder.withValidTokenIdRange(1, 10000);
    builder = builder.withManager(TokenizationHelpers.addressToCosmosString(msg.sender));
    
    // Set metadata
    TokenizationTypes.CollectionMetadata memory metadata = 
        TokenizationHelpers.createCollectionMetadata(
            string(abi.encodePacked("ipfs://", name)),
            string(abi.encodePacked('{"name":"', name, '","symbol":"', symbol, '"}'))
        );
    builder = builder.withMetadata(metadata);
    
    // Set standards
    string[] memory standards = new string[](1);
    standards[0] = "ERC-3643";
    builder = builder.withStandards(standards);
    
    // Set default balances (auto-approve self-initiated)
    builder = builder.withDefaultBalancesFromFlags(true, true, false);
    
    // Build and create
    string memory json = builder.build();
    return TOKENIZATION.createCollection(json);
}
```

### Querying Collections

```solidity
function getCollectionInfo(uint256 collectionId) external view {
    // Get collection bytes (requires off-chain decoding for full info)
    bytes memory collectionBytes = TokenizationWrappers.getCollection(
        TOKENIZATION,
        collectionId
    );
    
    // For simple queries, use direct methods
    uint256 totalSupply = TokenizationWrappers.getTotalSupply(
        TOKENIZATION,
        collectionId,
        tokenIds,
        ownershipTimes
    );
    
    uint256 userBalance = TokenizationWrappers.getBalanceAmount(
        TOKENIZATION,
        collectionId,
        msg.sender,
        tokenIds,
        ownershipTimes
    );
}
```

## Best Practices

1. **Use typed wrappers** (`TokenizationWrappers`) for better type safety
2. **Validate inputs** before calling precompile methods
3. **Handle errors** appropriately (use custom errors when possible)
4. **Emit events** for important state changes
5. **Use builders** for complex operations like collection creation
6. **Cache registry IDs** to avoid repeated lookups
7. **Consider gas costs** when designing batch operations

See [Best Practices](./BEST_PRACTICES.md) for more detailed guidance.













