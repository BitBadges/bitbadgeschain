# EVM Integration Documentation

## Overview

The BitBadges chain integrates the Cosmos `x/evm` module to enable Ethereum Virtual Machine (EVM) compatibility. This allows Solidity smart contracts to interact with the BitBadges tokenization module through a precompile interface.

## Architecture

```
Solidity Contract → EVM Precompile (0x0000000000000000000000000000000000001001) → Tokenization Module → Cosmos SDK State
```

### Components

1. **EVM Module**: Cosmos `x/evm` module providing EVM runtime
2. **Tokenization Precompile**: Precompiled contract at address `0x0000000000000000000000000000000000001001`
3. **ERC-3643 Wrapper**: Example Solidity contract demonstrating integration

## Precompile Address

The tokenization precompile is available at:
```
0x0000000000000000000000000000000000001001
```

## API Reference

### Transaction Methods

#### `transferTokens`

Transfers tokens from the caller to one or more recipients.

**Signature:**
```solidity
function transferTokens(
    uint256 collectionId,
    address[] calldata toAddresses,
    uint256 amount,
    UintRange[] calldata tokenIds,
    UintRange[] calldata ownershipTimes
) external returns (bool);
```

**Parameters:**
- `collectionId`: The collection ID to transfer from
- `toAddresses`: Array of recipient addresses
- `amount`: Amount to transfer to each recipient
- `tokenIds`: Array of token ID ranges to transfer
- `ownershipTimes`: Array of ownership time ranges to transfer

**Returns:**
- `bool`: `true` if transfer succeeded

**Gas Cost:** Base 30,000 + 5,000 per recipient + 1,000 per token ID range + 1,000 per ownership time range

#### `setIncomingApproval`

Sets an incoming approval for the caller.

**Signature:**
```solidity
function setIncomingApproval(
    uint256 collectionId,
    UserIncomingApproval calldata approval
) external returns (bool);
```

**Gas Cost:** Base 20,000 + dynamic based on ranges

#### `setOutgoingApproval`

Sets an outgoing approval for the caller.

**Signature:**
```solidity
function setOutgoingApproval(
    uint256 collectionId,
    UserOutgoingApproval calldata approval
) external returns (bool);
```

**Gas Cost:** Base 20,000 + dynamic based on ranges

### Query Methods

#### `getBalanceAmount`

Gets the balance amount for a user with specific token IDs and ownership times.

**Signature:**
```solidity
function getBalanceAmount(
    uint256 collectionId,
    address userAddress,
    UintRange[] calldata tokenIds,
    UintRange[] calldata ownershipTimes
) external view returns (uint256);
```

**Returns:**
- `uint256`: Total balance amount matching the criteria

#### `getTotalSupply`

Gets the total supply for a collection with specific token IDs and ownership times.

**Signature:**
```solidity
function getTotalSupply(
    uint256 collectionId,
    UintRange[] calldata tokenIds,
    UintRange[] calldata ownershipTimes
) external view returns (uint256);
```

**Returns:**
- `uint256`: Total supply matching the criteria

#### Other Query Methods

- `getCollection(uint256 collectionId) returns (bytes)`
- `getBalance(uint256 collectionId, address userAddress) returns (bytes)`
- `getAddressList(string listId) returns (bytes)`
- `getApprovalTracker(...) returns (bytes)`
- `getChallengeTracker(...) returns (uint256)`
- `getETHSignatureTracker(...) returns (uint256)`
- `getDynamicStore(uint256 storeId) returns (bytes)`
- `getDynamicStoreValue(uint256 storeId, address userAddress) returns (bytes)`
- `getWrappableBalances(string denom, address userAddress) returns (uint256)`
- `isAddressReservedProtocol(address addr) returns (bool)`
- `getAllReservedProtocolAddresses() returns (address[])`
- `getVote(...) returns (bytes)`
- `getVotes(...) returns (bytes)`
- `params() returns (bytes)`

## ERC-3643 Integration

The `ERC3643Badges.sol` contract provides an ERC-3643 compliant interface that wraps the precompile.

### Usage Example

```solidity
// Deploy the contract with a collection ID
ERC3643Badges token = new ERC3643Badges(1);

// Transfer tokens
token.transfer(recipient, amount);

// Query balance
uint256 balance = token.balanceOf(account);

// Query total supply
uint256 supply = token.totalSupply();
```

### Contract Address

The ERC-3643 contract must be deployed separately for each collection. The contract address is determined by the deployment transaction.

## Error Handling

The precompile uses structured error codes:

- `ErrorCodeInvalidInput (1)`: Invalid input parameters
- `ErrorCodeCollectionNotFound (2)`: Collection not found
- `ErrorCodeBalanceNotFound (3)`: Balance not found
- `ErrorCodeTransferFailed (4)`: Transfer operation failed
- `ErrorCodeApprovalFailed (5)`: Approval operation failed
- `ErrorCodeQueryFailed (6)`: Query operation failed
- `ErrorCodeInternalError (7)`: Internal error
- `ErrorCodeUnauthorized (8)`: Unauthorized operation

## Security Considerations

### Input Validation

All inputs are validated:
- Zero addresses are rejected
- Empty arrays are rejected where appropriate
- Invalid ranges (start > end) are rejected
- Negative amounts are rejected
- Array sizes are limited to prevent DoS

### Limits

- Maximum recipients: 100
- Maximum token ID ranges: 100
- Maximum ownership time ranges: 100
- Maximum approval ranges: 100

### Reentrancy Protection

Reentrancy is prevented by:
- EVM call stack mechanism
- Cosmos SDK state machine atomicity
- Precompile operations execute atomically

### Caller Verification

The caller address is verified by the EVM and cannot be spoofed. The caller is used as the "from" address for transfers.

## Events

The precompile emits Cosmos SDK events for all operations:

- `precompile_transfer_tokens`: Emitted on token transfers
- `precompile_set_incoming_approval`: Emitted on incoming approval updates
- `precompile_set_outgoing_approval`: Emitted on outgoing approval updates

## Gas Costs

### Gas Cost Rationale

Gas costs are designed to:
1. **Prevent DoS attacks**: Limit array sizes and charge per element
2. **Reflect computational complexity**: More complex operations cost more
3. **Encourage efficient usage**: Batch operations when possible

### Base Costs

**Transaction Methods:**
- `transferTokens`: 30,000 base gas
- `setIncomingApproval`: 20,000 base gas
- `setOutgoingApproval`: 20,000 base gas

**Query Methods:**
- Simple queries (`getCollection`, `getBalance`, `params`): 2,000 - 3,000 gas
- Complex queries (trackers, votes, stores): 5,000 gas
- Range-based queries (`getBalanceAmount`, `getTotalSupply`): 3,000 base + 500 per range

### Dynamic Costs

- **Per recipient**: 5,000 gas (for `transferTokens`)
- **Per token ID range**: 1,000 gas (for transactions), 500 gas (for queries)
- **Per ownership time range**: 1,000 gas (for transactions), 500 gas (for queries)
- **Per approval field**: 500 gas (for approval methods)

### Gas Estimation Examples

**Example 1: Simple Transfer**
```solidity
// Transfer to 1 recipient, 1 token ID range, 1 ownership time range
// Gas = 30,000 (base) + 5,000 (1 recipient) + 1,000 (1 token range) + 1,000 (1 ownership range)
// Total = 37,000 gas
```

**Example 2: Multi-Recipient Transfer**
```solidity
// Transfer to 5 recipients, 2 token ID ranges, 1 ownership time range
// Gas = 30,000 (base) + 25,000 (5 recipients) + 2,000 (2 token ranges) + 1,000 (1 ownership range)
// Total = 58,000 gas
```

**Example 3: Balance Query**
```solidity
// Query balance with 3 token ID ranges, 2 ownership time ranges
// Gas = 3,000 (base) + 1,500 (3 token ranges) + 1,000 (2 ownership ranges)
// Total = 5,500 gas
```

### Gas Optimization Tips

1. **Minimize Recipients**: Each recipient adds 5,000 gas. Batch transfers when possible.
2. **Consolidate Ranges**: Fewer ranges = lower gas costs
3. **Use Direct Return Methods**: `getBalanceAmount` and `getTotalSupply` avoid protobuf decoding overhead
4. **Cache Query Results**: Store frequently accessed data in contract storage
5. **Batch Operations**: Group multiple operations in a single transaction when possible

## Deployment Guide

### Prerequisites Checklist

Before deploying, ensure:

- [ ] Chain has EVM module enabled
- [ ] FeeMarket module is initialized
- [ ] Tokenization module is enabled
- [ ] Precompile address is available: `0x0000000000000000000000000000000000001001`
- [ ] EVM chain ID is configured correctly
- [ ] Genesis state includes EVM configuration

### Step-by-Step Deployment

#### 1. Genesis Configuration

The EVM module must be configured in genesis. Add the following to your genesis file:

```json
{
  "app_state": {
    "evm": {
      "params": {
        "evm_denom": "ustake",
        "enable_create": true,
        "enable_call": true,
        "extra_eips": [],
        "chain_config": {
          "chain_id": "9000",
          "homestead_block": "0",
          "dao_fork_block": "0",
          "dao_fork_support": true,
          "eip150_block": "0",
          "eip150_hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
          "eip155_block": "0",
          "eip158_block": "0",
          "byzantium_block": "0",
          "constantinople_block": "0",
          "petersburg_block": "0",
          "istanbul_block": "0",
          "muir_glacier_block": "0",
          "berlin_block": "0",
          "london_block": "0"
        }
      }
    },
    "feemarket": {
      "params": {
        "no_base_fee": false,
        "base_fee_change_denominator": 8,
        "elasticity_multiplier": 2,
        "enable_height": "0",
        "base_fee": "1000000000",
        "min_gas_price": "0",
        "min_gas_multiplier": "0.5"
      }
    }
  }
}
```

#### 2. Precompile Registration

The precompile is automatically registered when the EVM module is initialized in `app/evm.go`:

```go
tokenizationPrecompile := tokenizationprecompile.NewPrecompile(app.TokenizationKeeper)
tokenizationPrecompileAddr := common.HexToAddress(tokenizationprecompile.TokenizationPrecompileAddress)
app.EVMKeeper.RegisterStaticPrecompile(tokenizationPrecompileAddr, tokenizationPrecompile)
```

#### 3. Verify Precompile Registration

After deployment, verify the precompile is registered:

```bash
# Query precompile address (should return code)
cast code 0x0000000000000000000000000000000000001001 --rpc-url <your-rpc-url>
```

Or in Solidity:
```solidity
address precompileAddr = 0x0000000000000000000000000000000000001001;
uint256 codeSize;
assembly {
    codeSize := extcodesize(precompileAddr)
}
require(codeSize > 0, "Precompile not registered");
```

#### 4. Test Deployment

Create a test contract to verify the precompile works:

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IBadgesPrecompile {
    function params() external view returns (bytes);
}

contract PrecompileTest {
    IBadgesPrecompile constant PRECOMPILE = IBadgesPrecompile(0x0000000000000000000000000000000000001001);
    
    function testPrecompile() external view returns (bool) {
        try PRECOMPILE.params() returns (bytes memory) {
            return true;
        } catch {
            return false;
        }
    }
}
```

### Common Deployment Issues

1. **Precompile Not Found:**
   - Verify EVM module is enabled in genesis
   - Check that precompile registration code is executed
   - Ensure the precompile address matches exactly

2. **"Module Not Found" Errors:**
   - Verify tokenization module is enabled
   - Check that all required modules are initialized in correct order
   - Ensure FeeMarket is initialized before EVM

3. **Gas Estimation Failures:**
   - Verify EVM chain configuration is correct
   - Check that FeeMarket parameters are set
   - Ensure base fee is configured appropriately

4. **Transaction Reverts:**
   - Check that collections exist before transferring
   - Verify approvals are set correctly
   - Ensure sufficient balance for transfers

### Post-Deployment Verification

1. **Test Basic Query:**
   ```solidity
   bytes memory params = precompile.params();
   require(params.length > 0, "Params query failed");
   ```

2. **Test Collection Query:**
   ```solidity
   bytes memory collection = precompile.getCollection(collectionId);
   require(collection.length > 0, "Collection query failed");
   ```

3. **Test Transfer (if collection exists):**
   ```solidity
   bool success = precompile.transferTokens(...);
   require(success, "Transfer failed");
   ```

### Production Considerations

1. **Gas Price Configuration:**
   - Set appropriate base fee in FeeMarket
   - Configure min gas price to prevent spam
   - Monitor gas usage patterns

2. **Rate Limiting:**
   - Consider implementing rate limiting at chain level
   - Monitor for DoS attempts (large arrays, many ranges)
   - Set appropriate transaction size limits

3. **Monitoring:**
   - Monitor precompile usage via events
   - Track gas consumption patterns
   - Alert on error rate spikes

4. **Security:**
   - Review all input validations
   - Monitor for unusual transaction patterns
   - Keep EVM module and dependencies updated

## Testing

### Running Tests

```bash
# Run all precompile tests
go test ./x/evm/precompiles/tokenization/...

# Run specific test suite
go test ./x/evm/precompiles/tokenization/... -run TestE2ESuite

# Run error scenario tests
go test ./x/evm/precompiles/tokenization/... -run TestPrecompile_ErrorScenarios
```

### Test Coverage

- Unit tests for validation functions
- Integration tests for precompile methods
- E2E tests for complete transfer flows
- Error scenario tests
- ERC-3643 wrapper tests

## Troubleshooting

### Common Issues

1. **"collection not found" (Error Code 2)**: 
   - Ensure the collection ID exists
   - Verify the collection was created successfully
   - Check that you're using the correct collection ID format (uint256)

2. **"caller cannot be zero address" (Error Code 1)**:
   - Ensure the transaction has a valid sender
   - Verify `msg.sender` is not the zero address
   - Check that the contract is being called correctly

3. **"invalid input parameters" (Error Code 1)**:
   - Check that all addresses are non-zero
   - Verify ranges are valid (start <= end, both non-negative)
   - Ensure amounts are greater than zero
   - Check array sizes don't exceed limits (max 100 recipients, 100 ranges)

4. **"transfer failed" (Error Code 4)**:
   - Check that the caller has sufficient balance
   - Verify approvals are set correctly (incoming/outgoing)
   - Ensure collection approvals allow the transfer
   - Check that the collection is not archived

5. **"unauthorized operation" (Error Code 8)**:
   - Verify the caller has the necessary permissions
   - Check that approvals are configured correctly
   - Ensure the collection manager permissions allow the operation

6. **"collection is archived" (Error Code 9)**:
   - The collection is in read-only mode
   - Transfers and modifications are not allowed
   - Only queries are permitted

7. **"query failed" (Error Code 6)**:
   - The requested resource does not exist
   - Check that collection IDs, addresses, and other identifiers are correct
   - Verify the resource was created before querying

### Debugging

Enable debug logging to see detailed error messages:
```go
// In precompile methods, errors are wrapped with context
// Check the error code and message for details
```

**Interpreting Error Codes in Solidity:**
```solidity
// Errors are returned as revert reasons
// The error message includes the error code
// Example: "precompile error [code=2]: collection not found: collectionId: 123"
```

**Common Error Patterns:**
- Error Code 1: Usually indicates validation failures - check all inputs
- Error Code 2: Resource doesn't exist - verify IDs and addresses
- Error Code 4: Operation failed - check balances, approvals, permissions
- Error Code 6: Query failed - resource may not exist or be inaccessible
- Error Code 8: Permission denied - check approvals and manager permissions

### Performance Troubleshooting

1. **High Gas Costs:**
   - Reduce number of recipients (each adds 5,000 gas)
   - Minimize number of ranges (each adds 1,000 gas)
   - Batch operations when possible

2. **Slow Query Performance:**
   - Use `getBalanceAmount` and `getTotalSupply` instead of protobuf methods when possible
   - Limit range queries to necessary ranges only
   - Cache results when appropriate

3. **Protobuf Decoding Issues:**
   - Use direct return methods (`getBalanceAmount`, `getTotalSupply`) when possible
   - Ensure protobuf library is correctly configured
   - Verify byte array length before decoding

## Error Code Reference

All precompile methods return structured errors with specific error codes:

| Code | Name | When Returned | How to Handle |
|------|------|---------------|---------------|
| 1 | ErrorCodeInvalidInput | Invalid parameters (zero addresses, negative values, invalid ranges, etc.) | Validate all inputs before calling |
| 2 | ErrorCodeCollectionNotFound | Collection does not exist | Verify collection ID is correct and collection exists |
| 3 | ErrorCodeBalanceNotFound | Balance does not exist | Check that user has balance in the collection |
| 4 | ErrorCodeTransferFailed | Transfer failed (insufficient balance, approval issues, etc.) | Check balances, approvals, and permissions |
| 5 | ErrorCodeApprovalFailed | Approval operation failed | Verify approval parameters and permissions |
| 6 | ErrorCodeQueryFailed | Query operation failed | Resource may not exist or be inaccessible |
| 7 | ErrorCodeInternalError | Internal error (marshaling, etc.) | Report as bug, check logs |
| 8 | ErrorCodeUnauthorized | Unauthorized operation | Check approvals and manager permissions |
| 9 | ErrorCodeCollectionArchived | Collection is archived (read-only) | Collection cannot be modified, only queried |

### Handling Errors in Solidity

```solidity
// Example: Handle errors gracefully
function safeTransfer(uint256 collectionId, address to, uint256 amount) external {
    try precompile.transferTokens(...) returns (bool success) {
        require(success, "Transfer returned false");
    } catch Error(string memory reason) {
        // Parse error code from reason string
        // Format: "precompile error [code=X]: message"
        revert(reason);
    } catch {
        revert("Unknown error");
    }
}
```

### Example Error Scenarios

**Scenario 1: Invalid Collection ID**
```solidity
// Error: ErrorCodeInvalidInput (1) or ErrorCodeCollectionNotFound (2)
bytes memory collection = precompile.getCollection(999999);
// Will revert with error code 2 if collection doesn't exist
```

**Scenario 2: Insufficient Balance**
```solidity
// Error: ErrorCodeTransferFailed (4)
bool success = precompile.transferTokens(collectionId, recipients, amount, ...);
// Will revert with error code 4 if balance is insufficient
```

**Scenario 3: Missing Approval**
```solidity
// Error: ErrorCodeUnauthorized (8) or ErrorCodeTransferFailed (4)
bool success = precompile.transferTokens(...);
// Will revert if approvals are not set correctly
```

## Limitations

1. **Gas Estimation**: Dynamic gas calculation is limited in `RequiredGas` since arguments aren't parsed yet. Actual gas may vary based on input size.

2. **Query Parsing**: Most query methods return protobuf-encoded bytes, requiring decoding in Solidity. Use `getBalanceAmount` and `getTotalSupply` for direct uint256 returns.

3. **Rate Limiting**: Currently no per-address rate limiting (can be added at chain level). Array size limits (max 100) provide DoS protection.

4. **Protobuf Decoding**: Solidity contracts need a protobuf library to decode query responses. Consider using direct return methods when possible.

5. **Event Emission**: Events are emitted as Cosmos SDK events, not Solidity events. Use event indexing services to monitor precompile usage.

## Additional Resources

- [Cosmos EVM Module Documentation](https://github.com/cosmos/evm)
- [BitBadges Tokenization Module](../tokenization/README.md)
- [API Reference](./EVM_PRECOMPILE_API.md)
- [Usage Examples](./EVM_USAGE_EXAMPLES.md)

## Future Enhancements

- Enhanced gas estimation with argument parsing
- Batch operations for multiple transfers
- More query helpers returning simple types
- Rate limiting per address
- Enhanced monitoring and metrics

## References

- [Cosmos EVM Module Documentation](https://github.com/cosmos/evm)
- [ERC-3643 Standard](https://eips.ethereum.org/EIPS/eip-3643)
- [BitBadges Tokenization Module](../x/tokenization/README.md)

