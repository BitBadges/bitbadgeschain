# Installation

This guide covers installing and setting up the tokenization precompile in your Solidity project.

## Prerequisites

- Solidity ^0.8.0
- Access to BitBadges Chain (mainnet or testnet)
- Understanding of EVM precompiles

## Installation Steps

### 1. Import Required Files

Copy the following files to your project:

```
contracts/
├── interfaces/
│   └── ITokenizationPrecompile.sol
├── types/
│   └── TokenizationTypes.sol
└── libraries/
    └── BadgesHelpers.sol (optional)
```

### 2. Import in Your Contract

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./interfaces/ITokenizationPrecompile.sol";
import "./types/TokenizationTypes.sol";
```

### 3. Initialize Precompile

```solidity
contract MyContract {
    ITokenizationPrecompile constant precompile = 
        ITokenizationPrecompile(0x0000000000000000000000000000000000001001);
    
    // Your contract code
}
```

## Precompile Address

The tokenization precompile is available at:

```
0x0000000000000000000000000000000000001001
```

This address is fixed and cannot be changed.

## Network Configuration

### Mainnet

- **Chain ID**: [Your mainnet chain ID]
- **Precompile Address**: `0x0000000000000000000000000000000000001001`
- **RPC Endpoint**: [Your mainnet RPC]

### Testnet

- **Chain ID**: [Your testnet chain ID]
- **Precompile Address**: `0x0000000000000000000000000000000000001001`
- **RPC Endpoint**: [Your testnet RPC]

## Verification

### Check Precompile Availability

```solidity
function checkPrecompile() external view returns (bool) {
    // Try to call a simple query method
    try precompile.params() returns (bytes memory) {
        return true;
    } catch {
        return false;
    }
}
```

### Verify Address

```solidity
function verifyPrecompileAddress() external pure returns (address) {
    return address(precompile);
    // Should return: 0x0000000000000000000000000000000000001001
}
```

## Dependencies

### Required

- **Solidity**: ^0.8.0
- **EVM Compatibility**: Full EVM compatibility required

### Optional

- **BadgesHelpers**: Helper library for common operations
- **OpenZeppelin**: For additional security utilities

## Troubleshooting

### Precompile Not Found

If the precompile is not found:

1. Verify you're on the correct network
2. Check the precompile address
3. Ensure the chain has the precompile enabled

### Type Errors

If you encounter type errors:

1. Ensure `TokenizationTypes.sol` is imported
2. Check Solidity version compatibility
3. Verify type definitions match your version

### Compilation Errors

If compilation fails:

1. Check Solidity version (^0.8.0 required)
2. Verify all imports are correct
3. Check for missing dependencies

## Next Steps

- Read the [Getting Started Guide](../getting-started.md)
- Explore [API Reference](api-reference.md)
- Check out [Examples](examples.md)

## Resources

- [Cosmos SDK EVM Documentation](https://docs.cosmos.network/evm/v0.5.0/documentation/overview)
- [Tokenization Module Documentation](../../../_docs/TOKENIZATION_MODULE_ARCHITECTURE.md)






