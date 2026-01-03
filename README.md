# BitBadges Blockchain

A Cosmos SDK blockchain for digital token issuance and management, enabling advanced permission systems and transferability controls.

## Overview

BitBadges blockchain is built on the Cosmos SDK and Tendermint consensus, providing a robust infrastructure for digital token issuance and management with advanced permission systems, multi-tier approval controls, and IBC compatibility.

For detailed implementation documentation, architecture details, and feature explanations, see [docs.bitbadges.io](https://docs.bitbadges.io/).

## Quick Start

### Prerequisites

**Ubuntu 23.10+:**

```bash
sudo apt-get install git curl make build-essential gcc
snap install go --classic  # Go 1.21+
```

**Cross-compilation (optional):**

```bash
# For ARM64
sudo apt-get install gcc-aarch64-linux-gnu

# For macOS (see osxcross project for setup)
```

### Build from Source

```bash
# Build for all platforms
make build-all

# Build for specific platform
make build-linux/amd64
make build-darwin/amd64
make build-linux/arm64
```

### Using Ignite CLI

```bash
# Install Ignite CLI from https://ignite.com/cli
ignite chain init --skip-proto
ignite chain build --skip-proto
ignite chain serve --skip-proto
```

**Note**: The `--skip-proto` flag is required due to manual corrections in generated query files.

## Development

### Project Structure

```
bitbadgeschain/
├── x/                          # Blockchain modules
│   ├── badges/                 # Core token functionality
│   ├── maps/                   # Key-value mappings
│   ├── anchor/                 # Data anchoring
│   └── wasmx/                  # WASM extensions
├── proto/                      # Protocol buffer definitions
├── api/                        # Generated Go types
├── ts-client/                  # TypeScript client
├── app/                        # Application configuration
├── cmd/                        # CLI commands
└── _docs/                      # Development documentation
```

### Common Commands

**Build & Test:**

```bash
# Build main binary
go build ./cmd/bitbadgeschaind

# Run module tests
go test ./x/badges/...
go test ./x/badges/keeper/...

# Run specific test
go test ./x/badges/keeper/ -run TestMsgCreateDynamicStore

# Integration tests
ignite chain test
```

**Linting & Formatting:**

```bash
golangci-lint run ./x/badges/...
go fmt ./x/badges/...
```

**Protocol Buffers:**

```bash
# Generate Go code from proto files
ignite generate proto-go --yes

# Clean up versioned API folders (required after generation)
rm -rf api/badges/v*

# Stage generated files
git add *.pb.go *.pulsar.go
```

### Development Workflows

For detailed development guides, see:

-   [Adding New Message Types](_docs/ADDING_NEW_MSG_TYPES.md)
-   [Proto Addition Guide](_docs/PROTO_ADDITION_GUIDE.md)
-   [Module Architecture](_docs/BADGES_MODULE_ARCHITECTURE.md)

## Configuration

### Local Development

Configure your development environment with `config.yml`:

```yaml
version: 1
accounts:
    - name: alice
      coins: ['1000000000000000ustake', '1ubadge']
    - name: bob
      coins: ['99999999999999996ubadge', '1000000000000000ustake']

validators:
    - name: alice
      bonded: 1000000000000000ustake
```

### Network Deployment

Production configurations available in:

-   `config.testnet.yml` - Testnet configuration
-   Genesis files for mainnet and testnet networks
-   Release configurations in `release-info/`

## API & Integration

### REST API

The blockchain exposes a REST API for querying collections, balances, and approvals. OpenAPI specification available at `docs/static/openapi.yml`.

### TypeScript Client

Generated TypeScript client available in `ts-client/` for easy integration with web applications.

### WASM Bindings

Smart contracts can interact with the tokens module through custom WASM bindings in `custom-bindings/`.

## Releases

### Creating Releases

```bash
# Tag and push new version
git tag v1.0.0
git push origin v1.0.0
```

This automatically creates a draft release with configured build targets.

### Upgrade Process

Blockchain upgrades are coordinated through governance proposals. See `release-info/` for historical upgrade information and `app/upgrades/` for upgrade handler implementations.

## Testing

### Unit Tests

```bash
# Run all token module tests
go test ./x/badges/...

# Run with coverage
go test -cover ./x/badges/...

# Run specific keeper tests
go test ./x/badges/keeper/ -run TestMsgCreateCollection
```

### Integration Tests

```bash
# Full integration test suite
ignite chain test

# Simulation tests
go test ./x/badges/simulation/...
```

### Test Helpers

The module includes comprehensive test helpers in `x/badges/keeper/integration_*_test.go` for setting up test scenarios.

## Community

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes following the coding conventions
4. Add tests for new functionality
5. Submit a pull request

### Development Guidelines

-   Follow existing code patterns and conventions
-   Add comprehensive tests for new features
-   Update documentation for API changes
-   Use the project's linting and formatting tools

### Support

-   [BitBadges Documentation](https://docs.bitbadges.io/)
-   [BitBadges Website](https://bitbadges.io)
-   [Cosmos SDK Documentation](https://docs.cosmos.network)
-   [Ignite CLI Documentation](https://docs.ignite.com)

## Acknowledgments

-   Built with [Cosmos SDK](https://cosmos.network/) and [Ignite CLI](https://ignite.com/cli)
-   Inter-blockchain communication powered by [IBC](https://ibcprotocol.org/)

## License

This repository is licensed under the Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International License and is registered under US Copyright.
