# BitBadges Blockchain

**The most feature-rich tokenization standard ever built — exclusively available as a Cosmos SDK module.**

Enterprise-grade tokenization with 20+ compliance primitives, protocol-level enforcement, and full EVM compatibility. No smart contracts required. No code required. No per-project audits.

## Why BitBadges?

BitBadges provides a **no-code, plug-and-play tokenization module** designed for RWAs, compliant assets, custom stablecoins, and any use case requiring advanced transferability controls. Unlike existing standards (ERC-20, ERC-3643, ICS20), compliance is checked on **every transfer, everywhere** — in DEXs, liquidity pools, IBC transfers — automatically and at the protocol level.

**Key Capabilities:**
- ✅ **20+ Compliance Primitives** — Time-gating, KYC requirements, spend limits, freeze/revoke, multi-sig, royalties
- ✅ **Four Levels of Control** — Chain-level, issuer-level, sender-level, recipient-level rules
- ✅ **Full EVM Compatibility** — Solidity contracts call module via precompile; works with ERC-3643, DeFi protocols
- ✅ **Time-Dependent Balances** — Native vesting, auto-expiring subscriptions, time-locked ownership
- ✅ **IBC Interoperability** — Wrap to ICS20, tap into any IBC liquidity, cross-chain settlement
- ✅ **No-Code Deployment** — Full features without smart contracts or audits
- ✅ **7000+ Off-Chain Integrations** — Gate transfers by Discord, email, KYC providers, custom APIs

For detailed documentation, see [docs.bitbadges.io](https://docs.bitbadges.io/).

---

## The Paradigm Shift: Smart Tokens

Traditional tokens are simple balances with permissionless transferability. **Smart tokens** are different — they have programmable rules, custom transferability, compliance checks, and ownership controls built directly into the token standard itself, enforced on every transfer, everywhere, automatically.

| Traditional Tokens | Smart Tokens (BitBadges) |
|--------------------|--------------------------|
| Simple mint/transfer/burn | Programmable rules on every transfer |
| Compliance in app layer (bypassable) | Compliance in protocol layer (enforced) |
| Custom contracts per use case | No-code, composable primitives |
| Per-project audits | Single audited module, reused |
| Limited control (1-2 levels) | Four levels (chain, issuer, sender, recipient) |
| Static balances | Time-dependent balances (vesting, expiry) |

This represents a fundamental shift: instead of fragmented, vulnerable contracts on the app layer, we now have an abstracted tokenization layer that is self-contained, programmable, composable, and reusable across thousands of deployments.

**The issuer has full control over the entire lifecycle** — transferability, compliance, ownership, revocability, freezability — all configured through composable building blocks, not custom code.

---

## Technical Overview

BitBadges blockchain is built on the Cosmos SDK and CometBFT consensus, providing robust infrastructure for enterprise-grade tokenization.

**Infrastructure Features:**
- ✅ **Cosmos EVM Integration** — Full Ethereum Virtual Machine (EVM) compatibility via `cosmos/evm` module
- ✅ **ERC20 Support** — Native Cosmos coins can be wrapped as ERC20 tokens
- ✅ **Custom Precompiles** — Direct access to tokenization, Gamm, and SendManager modules from Solidity
- ✅ **Dual Wallet Support** — Same address works for both Cosmos and EVM transactions
- ✅ **JSON-RPC API** — Standard Ethereum JSON-RPC endpoints for Web3 compatibility

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
│   ├── tokenization/           # Core token functionality
│   │   └── precompile/         # Tokenization EVM precompile
│   ├── gamm/                   # AMM liquidity pools
│   │   └── precompile/         # Gamm EVM precompile
│   ├── sendmanager/            # Native coin transfers
│   │   └── precompile/         # SendManager EVM precompile
│   ├── maps/                   # Key-value mappings
│   ├── anchor/                 # Data anchoring
│   └── wasmx/                  # WASM extensions
├── app/                        # Application configuration
│   ├── evm.go                  # EVM module registration
│   └── PRECOMPILE_MANAGEMENT.md # Precompile documentation
├── contracts/                  # Solidity contracts and interfaces
│   ├── docs/                   # EVM integration guides
│   ├── interfaces/             # Precompile interfaces
│   └── libraries/              # Helper libraries
├── proto/                      # Protocol buffer definitions
├── api/                        # Generated Go types
├── ts-client/                  # TypeScript client
├── cmd/                        # CLI commands
└── _docs/                      # Development documentation
```

### Common Commands

**Build & Test:**

```bash
# Build main binary
go build ./cmd/bitbadgeschaind

# Run module tests
go test ./x/tokenization/...
go test ./x/tokenization/keeper/...

# Run specific test
go test ./x/tokenization/keeper/ -run TestMsgCreateDynamicStore

# Integration tests
ignite chain test
```

**Linting & Formatting:**

```bash
golangci-lint run ./x/tokenization/...
go fmt ./x/tokenization/...
```

**Protocol Buffers:**

```bash
# Generate Go code from proto files
ignite generate proto-go --yes

# Clean up versioned API folders (required after generation)
rm -rf api/tokenization/v*

# Stage generated files
git add *.pb.go *.pulsar.go
```

### Development Workflows

For detailed development guides, see:

-   [Adding New Message Types](_docs/ADDING_NEW_MSG_TYPES.md)
-   [Proto Addition Guide](_docs/PROTO_ADDITION_GUIDE.md)
-   [Module Architecture](_docs/TOKENIZATION_MODULE_ARCHITECTURE.md)

## Configuration

### EVM Configuration

The EVM module is configured in `app/evm.go` and requires:
- **EVM Chain IDs**: Set in `app/params/constants.go`
  - **Mainnet**: `50024` (BitBadges Mainnet)
  - **Testnet**: `50025` (BitBadges Testnet)
- **Precompile Enablement**: Configured in genesis `active_static_precompiles` array
- **JSON-RPC**: Optional, can be enabled for Web3 compatibility

For precompile management, see `app/PRECOMPILE_MANAGEMENT.md`.

**Note**: These chain IDs are registered in the [ethereum-lists/chains](https://github.com/ethereum-lists/chains) registry. See `_docs/CHAIN_ID_REGISTRATION.md` for details on chain ID registration.

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

### EVM Integration

BitBadges chain includes full EVM compatibility, enabling Ethereum developers to deploy and interact with Solidity smart contracts.

#### EVM Chain IDs

**Mainnet:**
- **Chain ID**: `50024` (BitBadges Mainnet)
- **Network Name**: BitBadges
- **Native Currency**: BADGE (ubadge base unit)
- **Registry**: To be registered in [ethereum-lists/chains](https://github.com/ethereum-lists/chains)

**Testnet:**
- **Chain ID**: `50025` (BitBadges Testnet)
- **Network Name**: BitBadges Testnet
- **Native Currency**: BADGE (ubadge base unit)
- **Registry**: To be registered in [ethereum-lists/chains](https://github.com/ethereum-lists/chains)

**Local Development:**
- **Chain ID**: `90123` (defaults to local dev chain ID)
- Configured in `app/params/constants.go`

#### JSON-RPC Endpoints

The chain exposes standard Ethereum JSON-RPC endpoints:
- `http://localhost:8545` - EVM JSON-RPC (if enabled)
- `http://localhost:26657` - Tendermint RPC (also supports some EVM queries)

#### Precompiles

Precompiles provide direct access to Cosmos modules from Solidity:

**Default Cosmos Precompiles** (0x0800-0x0806):
- `0x0800` - Staking precompile
- `0x0801` - Distribution precompile
- `0x0802` - ICS20 (IBC) precompile
- `0x0803` - Vesting precompile
- `0x0804` - Bank precompile (read-only queries)
- `0x0805` - Governance precompile
- `0x0806` - Slashing precompile

**Custom BitBadges Precompiles** (0x1001+):
- `0x1001` - **Tokenization precompile** - Create collections, transfer tokens, manage approvals
- `0x1002` - **Gamm precompile** - AMM liquidity pool operations
- `0x1003` - **SendManager precompile** - Send native Cosmos coins from EVM

See `contracts/docs/` for detailed precompile documentation:
- [Tokenization Precompile Guide](contracts/docs/GETTING_STARTED.md)
- [Gamm Precompile Guide](contracts/docs/GAMM_PRECOMPILE.md)
- [SendManager Precompile Guide](contracts/docs/SENDMANAGER_PRECOMPILE.md)
- [EVM Send Options](contracts/docs/EVM_SEND_OPTIONS.md)

#### ERC20 Wrapper

Native Cosmos coins can be wrapped as ERC20 tokens for use in standard Ethereum tooling:
- Each native denom has a corresponding ERC20 contract address
- Wrap/unwrap operations via ERC20 keeper
- Supports IBC transfers with ERC20 compatibility

#### Development Tools

- **MetaMask**: Connect using Chain ID `50024` (mainnet) or `50025` (testnet)
- **Hardhat/Truffle**: Use standard Ethereum development tools
- **Web3.js/Ethers.js**: Full compatibility with standard libraries
- **Example dApp**: See `counter-dapp/` for a complete Next.js + MetaMask example

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
go test ./x/tokenization/...

# Run with coverage
go test -cover ./x/tokenization/...

# Run specific keeper tests
go test ./x/tokenization/keeper/ -run TestMsgCreateCollection
```

### Integration Tests

```bash
# Full integration test suite
ignite chain test

# Simulation tests
go test ./x/tokenization/simulation/...
```

### Test Helpers

The module includes comprehensive test helpers in `x/tokenization/keeper/integration_*_test.go` for setting up test scenarios.

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
-   **Email**: trevor@bitbadges.io

## Acknowledgments

-   Built with [Cosmos SDK](https://cosmos.network/) and [Ignite CLI](https://ignite.com/cli)
-   Inter-blockchain communication powered by [IBC](https://ibcprotocol.org/)

## License

This repository is licensed under the Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International License and is registered under US Copyright.
