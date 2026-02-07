# EVM Integration Developer Guide

This guide covers essential information for developers building applications on BitBadges Chain that interact with both EVM and Cosmos SDK functionality.

## Table of Contents

1. [Transaction Signing Capabilities](#transaction-signing-capabilities)
2. [Address Conversion](#address-conversion)
3. [Precompile Caller Limitations](#precompile-caller-limitations)
4. [Key Types and Compatibility](#key-types-and-compatibility)
5. [Best Practices](#best-practices)

## Transaction Signing Capabilities

### Overview

BitBadges Chain supports two types of transactions, each with specific signing requirements:

| Transaction Type | Signing Key Type | Hash Algorithm | Use Case |
|-----------------|------------------|----------------|----------|
| `MsgEthereumTx` | `ethsecp256k1` (Ethereum-style) | Keccak256 | EVM contract calls, precompile calls |
| Standard Cosmos Messages | `secp256k1` (Cosmos-style) | SHA256 | Native Cosmos SDK messages (e.g., `MsgDelegate`, `MsgTransferTokens`) |

### Who Can Sign What?

#### ✅ ETH Wallets (ethsecp256k1 keys)

**Can sign:**
- `MsgEthereumTx` transactions
  - Direct EVM contract calls
  - Precompile calls from Solidity contracts
  - Any Ethereum-compatible transaction

**Cannot sign:**
- Standard Cosmos SDK messages (different hash algorithm)
  - `MsgDelegate`
  - `MsgTransferTokens`
  - Other native Cosmos messages

**Example:**
```solidity
// ✅ This works - ETH wallet can sign MsgEthereumTx
ITokenizationPrecompile precompile = ITokenizationPrecompile(0x0000000000000000000000000000000000001001);
precompile.transferTokens(collectionId, recipients, amount, tokenIds, ownershipTimes);
```

#### ✅ Cosmos Wallets (standard secp256k1 keys)

**Can sign:**
- Standard Cosmos SDK messages
  - `MsgDelegate`
  - `MsgTransferTokens`
  - `MsgCreateCollection`
  - All native Cosmos SDK messages

**Cannot sign:**
- `MsgEthereumTx` transactions (different signature format)
  - Cannot directly call EVM contracts
  - Cannot call precompiles from native Cosmos transactions

**Example:**
```go
// ✅ This works - Cosmos wallet can sign standard message
msg := &tokenizationtypes.MsgTransferTokens{
    Creator:      cosmosAddress,
    CollectionId: collectionId,
    Transfers:    transfers,
}
```

#### ❌ Cross-Compatibility

**Not supported out of the box:**
- ETH wallets cannot sign standard Cosmos messages
- Cosmos wallets cannot sign `MsgEthereumTx`

**Why?**
- Different hash algorithms (Keccak256 vs SHA256)
- Different signature formats
- EVM ante handler routes by transaction type

**Workaround:**
- Use precompiles from Solidity contracts (ETH wallets) to access Cosmos SDK functionality
- Use native Cosmos messages (Cosmos wallets) for direct SDK access

## Address Conversion

### Overview

BitBadges Chain uses a unified address system where Ethereum addresses (20 bytes) and Cosmos addresses (Bech32 format) represent the same underlying account when using `ethsecp256k1` keys.

### Conversion Mechanics

#### EVM Address → Cosmos Address

When an EVM address is used in Cosmos SDK operations, it's automatically converted:

```go
// In precompile code
caller := contract.Caller()  // common.Address (20 bytes, e.g., 0x1234...)
cosmosAddr := sdk.AccAddress(caller.Bytes()).String()  // Bech32 (e.g., bb1abc...)
```

**Key Points:**
- Same 20-byte value, different encoding
- EVM: `0x` hex prefix (e.g., `0x1234...5678`)
- Cosmos: Bech32 with `bb` prefix (e.g., `bb1abc...xyz`)

#### Cosmos Address → EVM Address

The reverse conversion is also possible:

```go
cosmosAddr, _ := sdk.AccAddressFromBech32("bb1abc...xyz")
evmAddr := common.BytesToAddress(cosmosAddr.Bytes())  // 0x1234...5678
```

### Address Format Examples

| Format | Example | Use Case |
|--------|---------|----------|
| EVM (hex) | `0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb` | Solidity contracts, EVM transactions |
| Cosmos (Bech32) | `bb1abc123def456...` | Cosmos SDK messages, queries |

### Important Notes

1. **Same Account, Different Representations**
   - An account with an `ethsecp256k1` key has the same 20-byte address in both formats
   - The Bech32 encoding is just a different way to represent the same bytes

2. **Automatic Conversion in Precompiles**
   - Precompiles automatically convert EVM addresses to Cosmos addresses
   - No manual conversion needed when calling precompile methods

3. **Address Validation**
   - Cosmos SDK validates addresses using Bech32 format
   - EVM uses hex format
   - Both represent the same underlying account

## Precompile Caller Limitations

### Understanding `msg.sender` in Precompiles

When a Solidity contract calls a precompile, the precompile uses `contract.Caller()` to identify the caller. This has important implications for authorization and security.

### How It Works

```go
// In precompile code
func (p Precompile) GetCallerAddress(contract *vm.Contract) (string, error) {
    caller := contract.Caller()  // EVM msg.sender
    return sdk.AccAddress(caller.Bytes()).String(), nil  // Converted to Cosmos address
}
```

### Limitations

#### 1. **Direct Contract Calls Only**

The caller is always the immediate caller (the contract making the precompile call), not the original transaction signer.

**Example:**
```
User (0xAlice) → Contract A → Precompile
                ↑
                contract.Caller() = Contract A's address, NOT 0xAlice
```

**Implication:**
- Precompile sees the contract address as the caller
- Original user address is not directly available
- Authorization must be handled at the contract level

#### 2. **No Cross-Contract Delegation**

If Contract A calls Contract B, which then calls the precompile:
```
User → Contract A → Contract B → Precompile
                                    ↑
                                    contract.Caller() = Contract B
```

**Implication:**
- Precompile only sees Contract B as the caller
- Contract A's identity is not visible to the precompile

#### 3. **Authorization Patterns**

Because the precompile sees the contract address, not the user address, you need to design authorization carefully:

**Pattern 1: Contract-Level Authorization**
```solidity
contract MyContract {
    ITokenizationPrecompile precompile = ITokenizationPrecompile(0x0000000000000000000000000000000000001001);
    
    mapping(address => bool) public authorized;
    
    function transferTokens(...) external {
        require(authorized[msg.sender], "Not authorized");
        // Contract is authorized, so precompile call succeeds
        precompile.transferTokens(...);
    }
}
```

**Pattern 2: Pass User Address Explicitly**
```solidity
// If precompile method accepts a "from" parameter, pass msg.sender
precompile.transferTokensFrom(msg.sender, ...);
```

**Pattern 3: Use Approval System**
```solidity
// User approves contract, contract acts on user's behalf
// Precompile sees contract as caller, but contract has user's approval
```

### Security Considerations

1. **Impersonation Prevention**
   - Precompile always uses `contract.Caller()` to prevent spoofing
   - Cannot be manipulated by malicious contracts

2. **Authorization Design**
   - Design contracts to handle authorization before calling precompile
   - Don't rely on precompile seeing the original user address

3. **Multi-Contract Scenarios**
   - Be aware that intermediate contracts are invisible to precompile
   - Design authorization to account for this

## Key Types and Compatibility

### Key Type Comparison

| Key Type | Algorithm | Hash Function | Address Format | Can Sign MsgEthereumTx | Can Sign Cosmos Messages |
|----------|-----------|---------------|----------------|----------------------|-------------------------|
| `ethsecp256k1` | secp256k1 | Keccak256 | Both (EVM hex + Cosmos Bech32) | ✅ Yes | ❌ No (different hash) |
| `secp256k1` (Cosmos) | secp256k1 | SHA256 | Cosmos Bech32 only | ❌ No (different format) | ✅ Yes |

### Account Creation

When creating accounts:

**For EVM Compatibility:**
```bash
# Create account with ethsecp256k1 key
bitbadgeschaind keys add my-eth-account --keyring-backend test --algo eth_secp256k1
```

**For Cosmos-Only:**
```bash
# Create account with standard secp256k1 key (default)
bitbadgeschaind keys add my-cosmos-account --keyring-backend test
```

### Key Registration

The codebase registers both key types:

```go
// ethsecp256k1 keys are registered in the codec
registry.RegisterImplementations((*cryptotypes.PubKey)(nil), &ethsecp256k1.PubKey{})
registry.RegisterImplementations((*cryptotypes.PrivKey)(nil), &ethsecp256k1.PrivKey{})
```

This allows the ante handler to recognize and verify signatures from both key types.

## Best Practices

### 1. Choose the Right Key Type

- **Use `ethsecp256k1`** if you need:
  - EVM contract interaction
  - Precompile calls from Solidity
  - Compatibility with Ethereum tooling

- **Use standard `secp256k1`** if you need:
  - Only native Cosmos SDK functionality
  - No EVM interaction required

### 2. Address Handling

- **In Solidity contracts:** Use EVM addresses (hex format)
- **In Cosmos SDK code:** Use Bech32 addresses
- **Precompiles handle conversion automatically**

### 3. Authorization Patterns

- **For contracts calling precompiles:**
  - Implement authorization at the contract level
  - Don't rely on precompile seeing original user address
  - Use approval patterns when needed

### 4. Transaction Types

- **Use `MsgEthereumTx`** for:
  - EVM contract calls
  - Precompile calls from Solidity

- **Use standard Cosmos messages** for:
  - Direct SDK module interaction
  - Better gas efficiency for simple operations
  - Native Cosmos tooling compatibility

### 5. Testing Considerations

- Test with both key types if your app supports both
- Verify address conversions work correctly
- Test authorization patterns with contract intermediaries

## Summary

| Aspect | ETH Wallets (ethsecp256k1) | Cosmos Wallets (secp256k1) |
|--------|---------------------------|---------------------------|
| **Can sign MsgEthereumTx** | ✅ Yes | ❌ No |
| **Can sign Cosmos messages** | ❌ No | ✅ Yes |
| **Address format** | Both (hex + Bech32) | Bech32 only |
| **Precompile access** | ✅ Via Solidity contracts | ❌ No direct access |
| **Native SDK access** | ❌ Via precompiles only | ✅ Direct access |

## Additional Resources

- [EVM Precompiles Overview](./overview.md)
- [Tokenization Precompile Documentation](./tokenization-precompile/README.md)
- [Architecture Details](./architecture.md)
- [Security Best Practices](./tokenization-precompile/security.md)

