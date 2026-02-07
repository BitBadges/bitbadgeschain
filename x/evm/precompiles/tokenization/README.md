# Tokenization Precompile

This package implements a precompiled contract for the BitBadges tokenization module, enabling Solidity smart contracts to interact with the tokenization system.

## Precompile Address

```
0x0000000000000000000000000000000000001001
```

## Files

- `precompile.go`: Main precompile implementation
- `validation.go`: Input validation helpers
- `errors.go`: Structured error types
- `events.go`: Event emission helpers
- `gas.go`: Dynamic gas calculation
- `security.go`: Security utilities
- `abi.json`: ABI definition for the precompile
- `precompile_test.go`: Unit tests
- `integration_test.go`: Integration tests
- `e2e_test.go`: End-to-end test suite
- `error_test.go`: Error scenario tests

## Usage

### In Solidity

```solidity
import "./interfaces/ITokenizationPrecompile.sol";

ITokenizationPrecompile precompile = ITokenizationPrecompile(0x0000000000000000000000000000000000001001);

// Transfer tokens
address[] memory recipients = new address[](1);
recipients[0] = recipient;
UintRange[] memory tokenIds = new UintRange[](1);
tokenIds[0] = UintRange({start: 1, end: 1});
UintRange[] memory ownershipTimes = new UintRange[](1);
ownershipTimes[0] = UintRange({start: 1, end: type(uint256).max});

bool success = precompile.transferTokens(
    collectionId,
    recipients,
    amount,
    tokenIds,
    ownershipTimes
);
```

## ABI Notes

### Return Type Simplifications

Due to the complexity of nested structures in the tokenization types, some fields are returned as empty arrays in Solidity struct representations:

- **TokenCollection**: `approvalCriteria` (in collectionApprovals), `invariants.cosmosCoinBackedPath`, `userPermissions` (in defaultBalances)
- **UserBalanceStore**: `approvalCriteria` (in outgoing/incomingApprovals)
- **CollectionApproval**: `approvalCriteria`

For full access to these fields, use the raw bytes returned from query methods and decode with protobuf.

### Known Limitations

1. **Invariants and Paths**: Cannot be set through the precompile. Use native Cosmos SDK interface.
2. **Nested Conversion Errors**: Some deeply nested structures may silently skip invalid items to maintain backward compatibility.
3. **ABI Load Failure**: If `abi.json` is corrupted, the precompile will be disabled but the chain will still start.

### DoS Protection Limits

| Field | Max Size |
|-------|----------|
| Recipients | 100 |
| Token ID Ranges | 100 |
| Ownership Time Ranges | 100 |
| Approval Ranges | 100 |
| Denom Units | 50 |
| Merkle Challenges | 20 |
| Coin Transfers | 50 |
| Metadata Length | 10,000 chars |

## Development

### Adding New Methods

1. Add method to `abi.json`
2. Add method constant to `precompile.go`
3. Add gas cost constant
4. Add case to `RequiredGas` switch
5. Add case to `Execute` switch
6. Implement method handler
7. Add validation
8. Add tests

### Testing

```bash
go test ./x/evm/precompiles/tokenization/...
```

