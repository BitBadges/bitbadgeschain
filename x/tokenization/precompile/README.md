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

All precompile methods now use a JSON string parameter that matches the protobuf JSON format. Use the helper library to construct JSON strings:

```solidity
import "./interfaces/ITokenizationPrecompile.sol";
import "./libraries/TokenizationJSONHelpers.sol";

ITokenizationPrecompile precompile = ITokenizationPrecompile(0x0000000000000000000000000000000000001001);

// Transfer tokens using JSON helper
address[] memory recipients = new address[](1);
recipients[0] = recipient;

string memory tokenIdsJson = TokenizationJSONHelpers.uintRangeToJson(1, 1);
string memory ownershipTimesJson = TokenizationJSONHelpers.uintRangeToJson(1, type(uint256).max);

string memory transferJson = TokenizationJSONHelpers.transferTokensJSON(
    collectionId,
    recipients,
    amount,
    tokenIdsJson,
    ownershipTimesJson
);

bool success = precompile.transferTokens(transferJson);
```

### Direct JSON Usage

You can also construct JSON strings manually:

```solidity
// Simple query example
string memory queryJson = string(abi.encodePacked(
    '{"collectionId":"', uintToString(collectionId), '"}'
));
bytes memory collection = precompile.getCollection(queryJson);
```

### Helper Library

The `TokenizationJSONHelpers` library provides convenient functions for common operations:
- `transferTokensJSON()` - Construct transfer JSON
- `getCollectionJSON()` - Construct collection query JSON
- `getBalanceJSON()` - Construct balance query JSON
- `getBalanceAmountJSON()` - Construct balance amount query JSON
- `uintRangeArrayToJson()` - Convert UintRange arrays to JSON
- `uintRangeToJson()` - Convert single UintRange to JSON

## ABI Notes

### Return Types (Go to Solidity)

Struct output from the precompile is converted end-to-end to match the Solidity type definitions in `TokenizationTypes.sol`. All nested fields (approvalCriteria, cosmosCoinBackedPath, collectionPermissions, userPermissions, Conversion, etc.) are fully converted when the ABI returns a struct tuple. The `Pack*AsStruct` helpers only support struct/tuple outputs; there is no legacy bytes mode. The main query path (HandleQuery) currently returns marshalled proto bytes for all getters; to return struct tuples instead, the ABI would need to declare struct return types and the handler would need to use the Pack*AsStruct helpers.

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

