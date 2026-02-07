# EVM Precompiles Documentation

Welcome to the BitBadges Chain EVM Precompiles documentation. This documentation provides comprehensive guides for using precompiled contracts to interact with Cosmos SDK modules from Solidity smart contracts.

## Quick Links

- [Overview](overview.md) - Introduction to EVM precompiles
- [Developer Guide](developer-guide.md) - **Essential guide for app developers** - Transaction signing, address conversion, and limitations
- [Getting Started](getting-started.md) - Quick start guide
- [Architecture](architecture.md) - System architecture
- [Tokenization Precompile](tokenization-precompile/README.md) - Tokenization precompile documentation

## What are Precompiles?

Precompiles are special contract addresses that execute native Go code instead of EVM bytecode. They provide:

- **Native Performance**: Direct access to Cosmos SDK modules
- **Type Safety**: Full type conversion between Solidity and Cosmos SDK
- **Security**: Built-in validation and error handling
- **Gas Efficiency**: Optimized gas costs

## Available Precompiles

### Tokenization Precompile

**Address:** `0x0000000000000000000000000000000000001001`

Full access to the BitBadges tokenization module:
- Token transfers with approval systems
- Collection management
- Balance queries
- Approval management
- Dynamic stores
- Governance features

[Read the Tokenization Precompile Documentation →](tokenization-precompile/README.md)

## Documentation Structure

```
docs/evm-precompiles/
├── SUMMARY.md              # GitBook table of contents
├── overview.md             # Introduction to precompiles
├── developer-guide.md      # Essential guide: signing, addresses, limitations
├── getting-started.md      # Quick start guide
├── architecture.md         # System architecture
├── tokenization-precompile/
│   ├── README.md          # Tokenization precompile overview
│   ├── overview.md        # Detailed overview
│   ├── installation.md    # Setup guide
│   ├── api-reference.md   # Complete API reference
│   ├── transactions.md    # Transaction methods
│   ├── queries.md         # Query methods
│   ├── types.md           # Type definitions
│   ├── errors.md          # Error handling
│   ├── gas.md             # Gas costs
│   ├── security.md         # Security best practices
│   └── examples.md        # Code examples
├── examples/              # Tutorials and examples
├── integration-guide.md   # Integration guide
├── troubleshooting.md     # Common issues
└── faq.md                 # Frequently asked questions
```

## Getting Started

1. **Read the Overview**: Understand what precompiles are and how they work
2. **Read the Developer Guide**: **Essential reading** - Learn about transaction signing, address conversion, and limitations
3. **Follow the Getting Started Guide**: Set up your first precompile interaction
4. **Explore the Tokenization Precompile**: Learn about available methods
5. **Check Examples**: See real-world usage patterns

## Resources

- [Cosmos SDK EVM Documentation](https://docs.cosmos.network/evm/v0.5.0/documentation/overview) - Official Cosmos SDK EVM docs
- [Tokenization Module](../_docs/TOKENIZATION_MODULE_ARCHITECTURE.md) - Tokenization module architecture
- [Example Contracts](../contracts/examples/) - Solidity example contracts

## Contributing

Documentation improvements are welcome! Please:

1. Follow the existing documentation style
2. Include code examples where helpful
3. Link to relevant Cosmos SDK documentation
4. Keep examples up-to-date with the latest code

## Support

For questions or issues:

- Check the [FAQ](faq.md)
- Review [Troubleshooting](troubleshooting.md)
- Open an issue on GitHub

---

**Last Updated:** 2024-12-19  
**Version:** 1.0.0




