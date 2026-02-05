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
import "./interfaces/IBadgesPrecompile.sol";

IBadgesPrecompile precompile = IBadgesPrecompile(0x0000000000000000000000000000000000001001);

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

