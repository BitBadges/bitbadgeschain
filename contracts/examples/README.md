# ERC-3643 Example Contracts

These example contracts demonstrate how to build compliant security tokens using the BitBadges tokenization precompile. Each example follows the ERC-3643 standard while leveraging BitBadges' native features for enhanced functionality.

## Overview

The Tokenization precompile is available at address `0x0000000000000000000000000000000000000808` and provides:

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
// Issue time-limited 2FA token
UintRange[] memory ownershipTimes = new UintRange[](1);
ownershipTimes[0] = UintRange(block.timestamp, block.timestamp + 24 hours);
TOKENIZATION.transferTokens(twoFactorCollectionId, [user], 1, tokenIds, ownershipTimes);

// Check 2FA ownership at CURRENT TIME
UintRange[] memory checkTime = new UintRange[](1);
checkTime[0] = UintRange(block.timestamp, block.timestamp);
uint256 balance = TOKENIZATION.getBalanceAmount(twoFactorCollectionId, user, tokenIds, checkTime);
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
// KYC verification using Dynamic Stores
kycRegistryId = TOKENIZATION.createDynamicStore(false, "ipfs://...", "...");
TOKENIZATION.setDynamicStoreValue(kycRegistryId, investor, true);

// Token transfers with compliance checks
TOKENIZATION.transferTokens(collectionId, recipients, amount, tokenIds, ownershipTimes);

// Balance queries
TOKENIZATION.getBalanceAmount(collectionId, account, tokenIds, ownershipTimes);
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
// Vintage years as token IDs
validTokenIds = UintRange(2020, 2100);

// Expiring ownership
ownershipTimes = UintRange(block.timestamp, vintages[vintage].expirationTime);

// Retirement via transfer to sink
TOKENIZATION.transferTokens(collectionId, [RETIREMENT_SINK], amount, tokenIds, ownershipTimes);
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
// Lock-up via Dynamic Store
lockUpRegistryId = TOKENIZATION.createDynamicStore(true, "...", "..."); // Default: locked

// Ownership time aligned with fund lifecycle
ownershipTimes = UintRange(block.timestamp, fundTermination);

// No auto-approval for controlled transfers
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
function transfer(address to, uint256 amount) external returns (bool) {
    require(canTransfer(msg.sender, to, amount), "Transfer not compliant");
    return TOKENIZATION.transferTokens(collectionId, ...);
}
```

### Dynamic Store for Boolean Registry
```solidity
// Create
uint256 registryId = TOKENIZATION.createDynamicStore(defaultValue, uri, customData);

// Set value
TOKENIZATION.setDynamicStoreValue(registryId, address, true/false);

// Read value
bytes memory result = TOKENIZATION.getDynamicStoreValue(registryId, address);
bool value = abi.decode(result, (bool));
```

### Time-Bound Ownership
```solidity
UintRange[] memory ownershipTimes = new UintRange[](1);
ownershipTimes[0] = UintRange(startTime, expirationTime);
```
