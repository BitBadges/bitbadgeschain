# EVM Precompiles Overview

## Introduction

BitBadges Chain extends the Cosmos SDK with Ethereum Virtual Machine (EVM) compatibility, enabling developers to interact with Cosmos SDK modules directly from Solidity smart contracts through **precompiled contracts** (precompiles).

Precompiles are special contracts that execute native Go code instead of EVM bytecode, providing:
- **Native Performance**: Direct access to Cosmos SDK modules without overhead
- **Type Safety**: Full type conversion between Solidity and Cosmos SDK types
- **Security**: Built-in validation and error handling
- **Gas Efficiency**: Optimized gas costs for common operations

## What are Precompiles?

Precompiles are special contract addresses that execute native code when called. Unlike regular smart contracts, precompiles:

- Execute Go code directly (not EVM bytecode)
- Have fixed addresses (cannot be deployed)
- Provide native access to Cosmos SDK modules
- Support both transactions (state-changing) and queries (read-only)

For more information on the Cosmos SDK EVM module, see the [official documentation](https://docs.cosmos.network/evm/v0.5.0/documentation/overview).

## Available Precompiles

### Tokenization Precompile

**Address:** `0x0000000000000000000000000000000000001001`

The tokenization precompile provides full access to the BitBadges tokenization module, enabling:

- Token transfers with complex approval systems
- Collection management (create, update, delete)
- Balance queries with time-based ownership
- Approval management (incoming/outgoing)
- Dynamic store operations
- Voting and governance features

See the [Tokenization Precompile documentation](tokenization-precompile/README.md) for complete details.

## Key Concepts

### Address Conversion

EVM addresses (20 bytes) are automatically converted to Cosmos addresses (Bech32 format) when interacting with Cosmos SDK modules. The precompile handles this conversion transparently.

### Gas Costs

Precompiles use dynamic gas calculation based on operation complexity:
- Base gas for each operation
- Additional gas per element (recipients, ranges, etc.)
- Query operations have lower gas costs

See [Gas & Costs](tokenization-precompile/gas.md) for detailed information.

### Error Handling

All precompiles use structured error handling with error codes:
- `ErrorCodeInvalidInput`: Invalid parameters
- `ErrorCodeCollectionNotFound`: Collection doesn't exist
- `ErrorCodeTransferFailed`: Transfer operation failed
- And more...

See [Error Handling](tokenization-precompile/errors.md) for complete error reference.

### Security

Precompiles implement comprehensive security measures:
- Caller verification (prevents impersonation)
- Input validation (prevents invalid data)
- DoS protection (array size limits)
- Error sanitization (prevents information leakage)

See [Security](tokenization-precompile/security.md) for security best practices.

## Getting Started

1. **Read the Developer Guide**: **Essential first step** - Understand transaction signing, address conversion, and limitations
2. **Install Dependencies**: Import the precompile interface in your Solidity contract
3. **Initialize Precompile**: Get the precompile instance at its fixed address
4. **Call Methods**: Use transaction or query methods as needed

See the [Developer Guide](developer-guide.md) for essential information about transaction signing and address handling.  
See the [Getting Started Guide](getting-started.md) for a step-by-step tutorial.

## Architecture

The precompile system consists of:

- **Precompile Contract**: Go implementation that interfaces with Cosmos SDK modules
- **ABI Definition**: JSON ABI for Solidity integration
- **Type Registry**: Solidity types matching Cosmos SDK proto definitions
- **Conversion Layer**: Type conversion between Solidity and Cosmos SDK types

See [Architecture](architecture.md) for detailed architecture documentation.

## Resources

- [Cosmos SDK EVM Documentation](https://docs.cosmos.network/evm/v0.5.0/documentation/overview)
- [Tokenization Module Documentation](../tokenization/README.md)
- [Example Contracts](../contracts/examples/)




