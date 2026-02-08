# Error Handling

The tokenization precompile uses structured error handling with error codes for consistent error reporting.

## Error Structure

All errors follow this structure:

```
precompile error [code=X]: Message: Details
```

Where:
- **code**: Numeric error code (see below)
- **Message**: Human-readable error message
- **Details**: Additional context (sanitized to prevent information leakage)

## Error Codes

| Code | Name | Description |
|------|------|-------------|
| 1 | `ErrorCodeInvalidInput` | Invalid input parameters |
| 2 | `ErrorCodeCollectionNotFound` | Collection does not exist |
| 3 | `ErrorCodeBalanceNotFound` | Balance not found |
| 4 | `ErrorCodeTransferFailed` | Transfer operation failed |
| 5 | `ErrorCodeApprovalFailed` | Approval operation failed |
| 6 | `ErrorCodeQueryFailed` | Query operation failed |
| 7 | `ErrorCodeInternalError` | Internal error |
| 8 | `ErrorCodeUnauthorized` | Unauthorized operation |
| 9 | `ErrorCodeCollectionArchived` | Collection is archived (read-only) |

## Error Handling in Solidity

### Basic Error Handling

```solidity
function safeTransfer(
    uint256 collectionId,
    address recipient,
    uint256 amount
) external {
    try precompile.transferTokens(...) returns (bool success) {
        require(success, "Transfer failed");
    } catch Error(string memory reason) {
        // Handle error
        revert(reason);
    } catch (bytes memory lowLevelData) {
        // Handle low-level error
        revert("Low-level error");
    }
}
```

### Error Code Extraction

Error codes are embedded in the error message. You can parse them:

```solidity
function parseErrorCode(string memory error) internal pure returns (uint256) {
    // Extract error code from "precompile error [code=X]: ..."
    // Implementation depends on your needs
}
```

## Common Errors

### Invalid Input

**Code:** 1  
**Causes:**
- Invalid collection ID (zero or negative)
- Invalid address (zero address)
- Invalid ranges (start > end)
- Array size exceeds limits
- Type mismatches

**Example:**
```
precompile error [code=1]: invalid input: collectionId cannot be zero
```

### Collection Not Found

**Code:** 2  
**Causes:**
- Collection ID doesn't exist
- Collection was deleted

**Example:**
```
precompile error [code=2]: collection not found: collection 123 does not exist
```

### Transfer Failed

**Code:** 4  
**Causes:**
- Insufficient balance
- Approval requirements not met
- Collection archived
- Invalid token IDs or ownership times

**Example:**
```
precompile error [code=4]: transfer failed: insufficient balance
```

### Unauthorized

**Code:** 8  
**Causes:**
- Not the collection creator/manager
- Missing required permissions
- Zero address caller

**Example:**
```
precompile error [code=8]: unauthorized: caller is not the collection manager
```

### Collection Archived

**Code:** 9  
**Causes:**
- Collection is archived (read-only)
- Attempting to modify archived collection

**Example:**
```
precompile error [code=9]: collection archived: collection 123 is archived
```

## Error Sanitization

Error details are sanitized to prevent information leakage:

- File paths removed
- Stack traces truncated
- Internal module paths redacted
- IP addresses removed

This ensures that error messages are helpful for debugging without exposing sensitive system information.

## Best Practices

### 1. Always Check Return Values

```solidity
bool success = precompile.transferTokens(...);
require(success, "Transfer failed");
```

### 2. Handle Errors Gracefully

```solidity
try precompile.transferTokens(...) returns (bool success) {
    if (!success) {
        // Handle failure
    }
} catch Error(string memory reason) {
    // Log error for debugging
    emit ErrorOccurred(reason);
    revert(reason);
}
```

### 3. Validate Inputs Before Calling

```solidity
function transferTokens(...) external {
    require(collectionId > 0, "Invalid collection ID");
    require(recipient != address(0), "Invalid recipient");
    // ... more validation
    
    precompile.transferTokens(...);
}
```

### 4. Use Events for Error Logging

```solidity
event TransferError(uint256 collectionId, string reason);

function safeTransfer(...) external {
    try precompile.transferTokens(...) {
        // Success
    } catch Error(string memory reason) {
        emit TransferError(collectionId, reason);
        revert(reason);
    }
}
```

## Error Recovery

Some errors are recoverable:

- **Invalid Input**: Fix input parameters and retry
- **Insufficient Balance**: Wait for balance or request transfer
- **Approval Required**: Set appropriate approvals first

Some errors are not recoverable:

- **Collection Not Found**: Collection doesn't exist
- **Unauthorized**: Caller doesn't have permission
- **Collection Archived**: Collection is read-only

## Debugging

When debugging errors:

1. **Check Error Code**: Identify the error category
2. **Review Error Message**: Understand what went wrong
3. **Validate Inputs**: Ensure all inputs are valid
4. **Check Permissions**: Verify caller has required permissions
5. **Check State**: Verify collection/balance exists

## Resources

- [API Reference](api-reference.md) - Method-specific error information
- [Security](security.md) - Security-related error handling
- [Troubleshooting](../troubleshooting.md) - Common issues and solutions









