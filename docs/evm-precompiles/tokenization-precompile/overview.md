# Tokenization Precompile Overview

## Introduction

The Tokenization Precompile is a precompiled contract that provides Solidity smart contracts with native access to the BitBadges tokenization module. It enables developers to build sophisticated token management systems directly in Solidity while leveraging the full power of the Cosmos SDK tokenization module.

## Key Features

### 1. Token Transfers

Transfer tokens with support for:
- Multiple recipients in a single transaction
- Token ID ranges (transfer multiple token IDs efficiently)
- Ownership time ranges (time-based ownership)
- Complex approval systems

### 2. Collection Management

Create and manage token collections:
- Create collections with custom metadata
- Update collection properties
- Set valid token ID ranges
- Manage collection permissions
- Archive/unarchive collections

### 3. Approval System

Sophisticated approval system with:
- Incoming approvals (who can send to you)
- Outgoing approvals (who you can send to)
- Complex approval criteria:
  - Merkle challenges
  - Predetermined balances
  - Voting challenges
  - Time-based checks
  - Address checks

### 4. Balance Queries

Query balances with:
- Token ID filtering
- Ownership time filtering
- Amount calculations
- Total supply queries

### 5. Dynamic Stores

Flexible key-value storage:
- Create custom stores
- Set values for addresses
- Query store values
- Global enable/disable

### 6. Governance

Participate in collection governance:
- Cast votes on proposals
- Query vote status
- Track voting challenges

## Architecture

The precompile consists of:

1. **Precompile Contract** (Go): Implements the EVM precompile interface
2. **ABI Definition** (JSON): Defines the Solidity interface
3. **Type Registry** (Solidity): Type definitions matching proto types
4. **Conversion Layer** (Go): Converts between Solidity and Cosmos SDK types
5. **Handler Methods** (Go): Implements business logic

## Address

The precompile is available at a fixed address:

```
0x0000000000000000000000000000000000001001
```

This address cannot be changed and is part of the chain's configuration.

## Type System

The precompile uses comprehensive type definitions that match the Cosmos SDK proto types:

- **UintRange**: Represents ranges of token IDs or ownership times
- **Balance**: User balance with token IDs and ownership times
- **CollectionMetadata**: Collection metadata (URI, custom data)
- **ApprovalCriteria**: Complex approval conditions
- **And many more...**

See [Types & Data Structures](types.md) for complete type reference.

## Security

The precompile implements comprehensive security measures:

- **Caller Verification**: Prevents impersonation attacks
- **Input Validation**: Validates all inputs before processing
- **DoS Protection**: Array size limits prevent DoS attacks
- **Error Sanitization**: Prevents information leakage

See [Security](security.md) for security best practices.

## Gas Costs

Gas costs are calculated dynamically based on operation complexity:

- **Base Gas**: Fixed cost per operation
- **Per-Element Costs**: Additional gas for arrays, ranges, etc.
- **Query Operations**: Lower gas costs (read-only)

See [Gas & Costs](gas.md) for detailed gas information.

## Error Handling

All methods use structured error handling:

- **Error Codes**: Categorized error types
- **Error Messages**: Clear, user-friendly messages
- **Error Sanitization**: Sensitive information removed

See [Error Handling](errors.md) for error codes and handling.

## Integration with Cosmos SDK

The precompile integrates with the Cosmos SDK EVM module, which provides:

- EVM compatibility layer
- Address conversion (EVM ↔ Cosmos)
- Gas accounting
- Transaction execution

For more information, see the [Cosmos SDK EVM documentation](https://docs.cosmos.network/evm/v0.5.0/documentation/overview).

## Limitations

### Known Limitations

1. **Invariants and Paths**: Cannot be set through the precompile. Use native Cosmos SDK interface.
2. **Complex Nested Types**: Some deeply nested structures may have simplified representations in Solidity.
3. **ABI Load Failure**: If `abi.json` is corrupted, the precompile will be disabled but the chain will still start.

### Design Decisions

1. **Creator Field**: Always set from `msg.sender` (cannot be specified in input)
2. **Overlapping Ranges**: Allowed for token IDs and ownership times (by design)
3. **Return Types**: Complex types returned as protobuf-encoded bytes for full compatibility

## Next Steps

- Read the [API Reference](api-reference.md)
- Explore [Transaction Methods](transactions.md)
- Review [Query Methods](queries.md)
- Check out [Examples](examples.md)
