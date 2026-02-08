# EVM Precompiles Architecture

## Overview

The EVM precompile system bridges Solidity smart contracts with Cosmos SDK modules, enabling native access to blockchain functionality from EVM-compatible code.

## System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Solidity Smart Contract                  │
│  (EVM Address: 0x0000000000000000000000000000000000001001)  │
└────────────────────────┬────────────────────────────────────┘
                         │
                         │ ABI-encoded calls
                         │
┌────────────────────────▼────────────────────────────────────┐
│              Precompile Contract (Go)                       │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  ABI Decoder                                         │   │
│  │  - Method ID resolution                             │   │
│  │  - Parameter unpacking                              │   │
│  └──────────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Type Converter                                      │   │
│  │  - Solidity → Cosmos SDK types                      │   │
│  │  - Validation & sanitization                        │   │
│  └──────────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Security Layer                                      │   │
│  │  - Caller verification                              │   │
│  │  - Input validation                                 │   │
│  │  - DoS protection                                   │   │
│  └──────────────────────────────────────────────────────┘   │
└────────────────────────┬────────────────────────────────────┘
                         │
                         │ Cosmos SDK messages
                         │
┌────────────────────────▼────────────────────────────────────┐
│              Cosmos SDK Module (Keeper)                      │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Tokenization Keeper                                 │   │
│  │  - State management                                  │   │
│  │  - Business logic                                    │   │
│  │  - Validation                                        │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

## Components

### 1. Precompile Contract

The precompile contract (`x/evm/precompiles/tokenization/precompile.go`) implements the `vm.PrecompiledContract` interface:

- **RequiredGas()**: Calculates gas costs based on input complexity
- **Run()**: Executes the precompile method
- **Execute()**: Dispatches to appropriate handler based on method name

### 2. ABI Definition

The ABI (`abi.json`) defines the Solidity interface:

- Method signatures
- Parameter types
- Return types
- Event definitions

### 3. Type Registry

Solidity type definitions (`contracts/types/TokenizationTypes.sol`) mirror Cosmos SDK proto types:

- Struct definitions
- Type mappings
- Helper functions

### 4. Conversion Layer

Type conversion utilities (`conversions.go`) handle:

- Solidity structs → Cosmos SDK proto types
- Cosmos SDK proto types → Solidity structs
- Validation and error handling

### 5. Handler Methods

Transaction and query handlers (`handlers.go`) implement:

- Input validation
- Type conversion
- Keeper method calls
- Response formatting

## Data Flow

### Transaction Flow

1. **Solidity Call**: Contract calls precompile method
2. **ABI Decoding**: Precompile decodes method ID and parameters
3. **Type Conversion**: Solidity types converted to Cosmos SDK types
4. **Validation**: Input validation and security checks
5. **Keeper Call**: Cosmos SDK keeper method invoked
6. **State Update**: Blockchain state updated
7. **Event Emission**: Events emitted for indexing
8. **Response**: Success/failure returned to Solidity

### Query Flow

1. **Solidity Call**: Contract calls query method
2. **ABI Decoding**: Precompile decodes method ID and parameters
3. **Type Conversion**: Solidity types converted to Cosmos SDK types
4. **Keeper Query**: Cosmos SDK query method invoked
5. **Response Encoding**: Results encoded as protobuf bytes
6. **Return**: Bytes returned to Solidity for decoding

## Address Conversion

EVM addresses (20 bytes) are automatically converted to Cosmos addresses:

```go
// EVM address → Cosmos address
caller := contract.Caller()  // common.Address (20 bytes)
cosmosAddr := sdk.AccAddress(caller.Bytes()).String()  // Bech32 format
```

## Gas Calculation

Gas costs are calculated dynamically based on operation complexity:

```go
// Base gas + per-element costs
gas = baseGas + (numRecipients * gasPerRecipient) + 
      (numTokenRanges * gasPerRange) + ...
```

See [Gas & Costs](tokenization-precompile/gas.md) for detailed gas calculation.

## Security Architecture

### Caller Verification

All transaction methods verify the caller:

```go
caller := contract.Caller()
if err := VerifyCaller(caller); err != nil {
    return nil, err
}
```

### Input Validation

All inputs are validated before processing:

- Type checking
- Range validation
- Size limits (DoS protection)
- Business rule validation

### Error Handling

Structured errors prevent information leakage:

```go
return nil, ErrInvalidInput("invalid collectionId")
```

See [Security](tokenization-precompile/security.md) for security details.

## Integration with Cosmos SDK

The precompile integrates with Cosmos SDK through:

- **Keeper Interface**: Direct access to module keeper
- **Context**: SDK context for state access
- **Messages**: Standard Cosmos SDK message types
- **Queries**: Standard Cosmos SDK query types

For more information on Cosmos SDK EVM integration, see the [official documentation](https://docs.cosmos.network/evm/v0.5.0/documentation/overview).

## Performance Considerations

### Gas Optimization

- Dynamic gas calculation based on actual complexity
- Efficient type conversion
- Minimal overhead for simple operations

### Caching

- ABI loaded once at initialization
- Method lookup optimized (O(1) map lookup)
- Type conversion cached where possible

## Extension Points

### Adding New Methods

1. Define method in ABI (`abi.json`)
2. Add method constant (`precompile.go`)
3. Implement handler (`handlers.go`)
4. Add type conversions if needed (`conversions.go`)
5. Update gas calculation (`gas.go`)

### Adding New Precompiles

1. Create new precompile package
2. Implement `vm.PrecompiledContract` interface
3. Register in EVM module
4. Define ABI and Solidity interface

## Resources

- [Cosmos SDK EVM Documentation](https://docs.cosmos.network/evm/v0.5.0/documentation/overview)
- [Tokenization Module Architecture](../tokenization/TOKENIZATION_MODULE_ARCHITECTURE.md)
- [Precompile Implementation](../x/evm/precompiles/tokenization/)









