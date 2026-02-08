# Security

The tokenization precompile implements comprehensive security measures to protect against common attack vectors.

## Security Features

### 1. Caller Verification

All transaction methods verify the caller to prevent impersonation:

```go
caller := contract.Caller()
if err := VerifyCaller(caller); err != nil {
    return nil, err
}
```

**Protection:**
- Prevents zero address calls
- Ensures `msg.sender` is valid
- Cannot be spoofed by malicious contracts

### 2. Input Validation

All inputs are validated before processing:

- **Type Checking**: Ensures correct types
- **Range Validation**: Validates token ID and ownership time ranges
- **Size Limits**: Prevents DoS attacks
- **Business Rules**: Validates business logic constraints

**Example:**
```go
if collectionId.Sign() <= 0 {
    return nil, ErrInvalidInput("collectionId cannot be zero or negative")
}
```

### 3. DoS Protection

Array size limits prevent denial-of-service attacks:

| Field | Max Size |
|-------|----------|
| Recipients | 100 |
| Token ID Ranges | 100 |
| Ownership Time Ranges | 100 |
| Approval Ranges | 100 |

**Protection:**
- Prevents excessive gas consumption
- Limits memory usage
- Prevents transaction timeouts

### 4. Error Sanitization

Error messages are sanitized to prevent information leakage:

- File paths removed
- Stack traces truncated
- Internal module paths redacted
- IP addresses removed

**Example:**
```
// Before sanitization:
Error: /home/user/project/handlers.go:123: collection not found

// After sanitization:
Error: collection not found
```

### 5. Creator Field Protection

The `creator` field is always set from `msg.sender`:

```go
creatorCosmosAddr, err := p.GetCallerAddress(contract)
// creator field cannot be specified in input
```

**Protection:**
- Prevents impersonation attacks
- Creator field not exposed in ABI
- Always uses actual caller address

## Security Best Practices

### 1. Validate All Inputs

Always validate inputs in your Solidity contract:

```solidity
function transferTokens(
    uint256 collectionId,
    address recipient,
    uint256 amount
) external {
    require(collectionId > 0, "Invalid collection ID");
    require(recipient != address(0), "Invalid recipient");
    require(amount > 0, "Invalid amount");
    
    precompile.transferTokens(...);
}
```

### 2. Check Return Values

Always check return values:

```solidity
bool success = precompile.transferTokens(...);
require(success, "Transfer failed");
```

### 3. Handle Errors Gracefully

Use try-catch for error handling:

```solidity
try precompile.transferTokens(...) returns (bool success) {
    require(success, "Transfer failed");
} catch Error(string memory reason) {
    revert(reason);
}
```

### 4. Use Events for Logging

Log important operations:

```solidity
event TransferAttempted(uint256 collectionId, address recipient, bool success);

function transferTokens(...) external {
    bool success = precompile.transferTokens(...);
    emit TransferAttempted(collectionId, recipient, success);
    require(success, "Transfer failed");
}
```

### 5. Reentrancy Protection

The precompile is protected against reentrancy by the EVM's call stack, but you should still follow best practices:

```solidity
bool private locked;

modifier noReentrant() {
    require(!locked, "Reentrant call");
    locked = true;
    _;
    locked = false;
}

function transferTokens(...) external noReentrant {
    precompile.transferTokens(...);
}
```

## Threat Model

### Protected Against

✅ **Reentrancy Attacks**: Prevented by EVM call stack and atomic transactions  
✅ **Integer Overflow**: Prevented by validation and `sdkmath.Uint` type  
✅ **Invalid Input Attacks**: Prevented by comprehensive input validation  
✅ **DoS Attacks**: Prevented by array size limits  
✅ **Information Leakage**: Prevented by structured error handling  
✅ **State Corruption**: Prevented by atomic transactions and keeper validation  
✅ **Impersonation**: Prevented by caller verification  

### Known Limitations

⚠️ **Rate Limiting**: Not implemented at precompile level (can be added at chain level)  
⚠️ **Gas Price Manipulation**: Handled by EVM module, not precompile  
⚠️ **Access Control**: Handled by tokenization module's approval system  

## Security Audit

The precompile has been evaluated for production readiness:

- ✅ Comprehensive type safety
- ✅ Input validation
- ✅ Error handling
- ✅ DoS protection
- ✅ Security patterns documented

See the [Production Readiness Evaluation](../../x/evm/precompiles/tokenization/THIRD_EVALUATION.md) for details.

## Reporting Security Issues

If you discover a security vulnerability, please report it responsibly:

1. **Do not** open a public issue
2. Contact the security team directly
3. Provide detailed information about the vulnerability
4. Allow time for the issue to be addressed before disclosure

## Resources

- [Error Handling](errors.md) - Error codes and handling
- [API Reference](api-reference.md) - Method-specific security notes
- [Cosmos SDK Security](https://docs.cosmos.network/evm/v0.5.0/documentation/overview) - EVM security model









