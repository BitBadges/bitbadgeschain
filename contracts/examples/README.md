# ERC-3643 Example Contracts

These example contracts demonstrate how to build compliant security tokens using the BitBadges tokenization precompile. Each example follows the ERC-3643 standard while leveraging BitBadges' native features for enhanced functionality.

## Documentation

For comprehensive documentation, see the [docs](../docs/) directory:

- **[Getting Started](../docs/GETTING_STARTED.md)** - Quick start guide and basic examples
- **[API Reference](../docs/API_REFERENCE.md)** - Complete API documentation
- **[Common Patterns](../docs/PATTERNS.md)** - Common patterns and use cases
- **[Extended Examples](../docs/EXAMPLES.md)** - More detailed examples
- **[Troubleshooting](../docs/TROUBLESHOOTING.md)** - Common issues and solutions
- **[Best Practices](../docs/BEST_PRACTICES.md)** - Security and optimization tips
- **[Gas Optimization](../docs/GAS_OPTIMIZATION.md)** - Gas optimization strategies

## Overview

The Tokenization precompile is available at address `0x0000000000000000000000000000000000001001` and provides:

- **Collections**: Token issuance and management
- **Dynamic Stores**: On-chain boolean registries (perfect for KYC/compliance)
- **Ownership Times**: Built-in support for time-bound ownership (lock-ups, expirations)
- **Token ID Ranges**: Multi-asset support within a single collection

## Example Contracts

### 1. TwoFactorSecurityToken.sol

**Use Case**: Security token requiring 2FA token ownership for transfers

**Key Features**:
- Separate 2FA token collection for authentication
- Time-limited 2FA tokens (auto-expire via `ownershipTimes`)
- Current-time ownership check as transfer requirement
- Cooldown between 2FA token requests
- High-value transfers require both sender AND recipient 2FA

**Tokenization Features Used**:
```solidity
import "./libraries/TokenizationJSONHelpers.sol";

// Issue time-limited 2FA token
string memory tokenIdsJson = TokenizationJSONHelpers.uintRangeToJson(1, 1);
string memory ownershipTimesJson = TokenizationJSONHelpers.uintRangeToJson(
    block.timestamp, 
    block.timestamp + 24 hours
);
address[] memory recipients = new address[](1);
recipients[0] = user;
string memory transferJson = TokenizationJSONHelpers.transferTokensJSON(
    twoFactorCollectionId, recipients, 1, tokenIdsJson, ownershipTimesJson
);
TOKENIZATION.transferTokens(transferJson);

// Check 2FA ownership at CURRENT TIME (single tokenId/ownershipTime)
string memory balanceJson = TokenizationJSONHelpers.getBalanceAmountJSON(
    twoFactorCollectionId,
    user,
    currentNonce,             // Nonce as token ID
    block.timestamp * 1000    // Current time in milliseconds
);
uint256 balance = TOKENIZATION.getBalanceAmount(balanceJson);
bool hasValid2FA = balance > 0;
```

**Flow**:
1. User requests 2FA token from authority
2. Authority issues time-limited token (e.g., 24 hours)
3. User can perform transfers while 2FA token is valid
4. Token auto-expires - no revocation needed

---

### 2. RealEstateSecurityToken.sol

**Use Case**: Tokenized real estate (Reg D offering)

**Key Features**:
- Identity Registry via Dynamic Stores (`kycRegistryId`, `accreditedRegistryId`)
- Address freezing for regulatory compliance
- Recovery transfer mechanism for lost wallets
- Token pausability

**Tokenization Features Used**:
```solidity
import "./libraries/TokenizationJSONHelpers.sol";

// KYC verification using Dynamic Stores
string memory createStoreJson = TokenizationJSONHelpers.createDynamicStoreJSON(
    false, "ipfs://...", "..."
);
kycRegistryId = TOKENIZATION.createDynamicStore(createStoreJson);

string memory setValueJson = TokenizationJSONHelpers.setDynamicStoreValueJSON(
    kycRegistryId, investor, true
);
TOKENIZATION.setDynamicStoreValue(setValueJson);

// Token transfers with compliance checks
string memory transferJson = TokenizationJSONHelpers.transferTokensJSON(
    collectionId, recipients, amount, tokenIdsJson, ownershipTimesJson
);
TOKENIZATION.transferTokens(transferJson);

// Balance queries (single tokenId/ownershipTime)
string memory balanceJson = TokenizationJSONHelpers.getBalanceAmountJSON(
    collectionId,
    account,
    1,                        // Single token ID
    block.timestamp * 1000    // Current time in milliseconds
);
TOKENIZATION.getBalanceAmount(balanceJson);
```

### 3. CarbonCreditToken.sol

**Use Case**: Verified carbon credits with vintage tracking

**Key Features**:
- Multiple vintages as different token IDs (2020-2100)
- Time-bound ownership (credits expire)
- Retirement tracking (permanent removal from circulation)
- Verified buyer/seller registries

**Tokenization Features Used**:
```solidity
import "./libraries/TokenizationJSONHelpers.sol";

// Vintage years as token IDs
string memory validTokenIdsJson = TokenizationJSONHelpers.uintRangeToJson(2020, 2100);

// Expiring ownership
string memory ownershipTimesJson = TokenizationJSONHelpers.uintRangeToJson(
    block.timestamp, 
    vintages[vintage].expirationTime
);

// Retirement via transfer to sink
address[] memory sinkRecipients = new address[](1);
sinkRecipients[0] = RETIREMENT_SINK;
string memory tokenIdsJson = TokenizationJSONHelpers.uintRangeToJson(vintage, vintage);
string memory retireJson = TokenizationJSONHelpers.transferTokensJSON(
    collectionId, sinkRecipients, amount, tokenIdsJson, ownershipTimesJson
);
TOKENIZATION.transferTokens(retireJson);
```

### 4. PrivateEquityToken.sol

**Use Case**: Private equity fund LP interests (Reg D 506(c))

**Key Features**:
- Lock-up period enforcement
- Maximum ownership limits (25% default)
- Qualified Purchaser and Accredited Investor verification
- Capital call processing
- Blacklist registry

**Tokenization Features Used**:
```solidity
import "./libraries/TokenizationJSONHelpers.sol";

// Lock-up via Dynamic Store
string memory createStoreJson = TokenizationJSONHelpers.createDynamicStoreJSON(
    true, "...", "..."  // Default: locked
);
lockUpRegistryId = TOKENIZATION.createDynamicStore(createStoreJson);

// Ownership time aligned with fund lifecycle
string memory ownershipTimesJson = TokenizationJSONHelpers.uintRangeToJson(
    block.timestamp, 
    fundTermination
);

// No auto-approval for controlled transfers (set in createCollection JSON)
defaultBalances.autoApproveSelfInitiatedOutgoingTransfers = false;
```

## ERC-3643 Compliance Mapping

| ERC-3643 Component | Tokenization Implementation |
|-------------------|-------------------------|
| Identity Registry | Dynamic Stores |
| Compliance Module | Solidity contract logic + Dynamic Store checks |
| Token | Tokenization Collection |
| Transfer Restrictions | `canTransfer()` checks before `transferTokens()` |
| Forced Transfers | Collection manager permissions |
| Token Recovery | Compliance agent with collection approvals |

## Deployment Flow

1. Deploy the contract (creates Dynamic Stores for compliance)
2. Call `initializeCollection()` to create the Tokenization collection
3. Register investors in compliance registries
4. Issue/transfer tokens with built-in compliance checks

## Key Patterns

### Compliance Check Before Transfer
```solidity
import "./libraries/TokenizationJSONHelpers.sol";

function transfer(address to, uint256 amount) external returns (bool) {
    require(canTransfer(msg.sender, to, amount), "Transfer not compliant");
    
    address[] memory recipients = new address[](1);
    recipients[0] = to;
    string memory tokenIdsJson = TokenizationJSONHelpers.uintRangeToJson(1, 1);
    // Note: Use uint64.max for ownership times (BitBadges internal limit)
    string memory ownershipTimesJson = TokenizationJSONHelpers.uintRangeToJson(1, type(uint64).max);
    string memory transferJson = TokenizationJSONHelpers.transferTokensJSON(
        collectionId, recipients, amount, tokenIdsJson, ownershipTimesJson
    );
    return TOKENIZATION.transferTokens(transferJson);
}
```

### Dynamic Store for Boolean Registry
```solidity
import "./libraries/TokenizationJSONHelpers.sol";

// Create
string memory createJson = TokenizationJSONHelpers.createDynamicStoreJSON(
    defaultValue, uri, customData
);
uint256 registryId = TOKENIZATION.createDynamicStore(createJson);

// Set value
string memory setValueJson = TokenizationJSONHelpers.setDynamicStoreValueJSON(
    registryId, address, true
);
TOKENIZATION.setDynamicStoreValue(setValueJson);

// Read value
string memory getValueJson = TokenizationJSONHelpers.getDynamicStoreValueJSON(
    registryId, address
);
bytes memory result = TOKENIZATION.getDynamicStoreValue(getValueJson);
bool value = abi.decode(result, (bool));
```

### Time-Bound Ownership
```solidity
import "./libraries/TokenizationJSONHelpers.sol";

string memory ownershipTimesJson = TokenizationJSONHelpers.uintRangeToJson(
    startTime, 
    expirationTime
);
```
