# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common Development Commands

### Building
```bash
# Build for current platform
make build-linux/amd64

# Build all platforms
make build-all

# Build with Ignite CLI (use --skip-proto flag due to manual proto corrections)
ignite chain build --skip-proto

# Build main binary for testing
go build ./cmd/bitbadgeschaind
```

### Testing
```bash
# Run all tokenization module tests
go test ./x/tokenization/...

# Run keeper tests specifically
go test ./x/tokenization/keeper/...

# Run specific test
go test ./x/tokenization/keeper/ -run TestMsgCreateDynamicStore

# Run integration tests
ignite chain test
```

### Linting and Formatting
```bash
# Check for linting issues
golangci-lint run ./x/tokenization/...

# Format code
go fmt ./x/tokenization/...
```

### Development Server
```bash
# Serve blockchain locally (must use --skip-proto)
ignite chain serve --skip-proto

# Initialize chain
ignite chain init --skip-proto
```

### Protocol Buffers
```bash
# Generate Go code from proto definitions
ignite generate proto-go --yes

# After generation, remove versioned API folders
ls api/tokenization/
rm -rf api/tokenization/v*

# Stage generated proto files
git add *.pb.go *.pulsar.go
```

## Architecture Overview

### Core Structure
This is a Cosmos SDK blockchain built with Ignite CLI that implements cross-chain digital token issuance and management.

### Key Modules
- **x/tokenization** - Core token functionality (collections, transfers, balances, permissions)
  - `precompile/` - EVM precompile for Solidity integration
- **x/gamm** - AMM liquidity pools
  - `precompile/` - EVM precompile for pool operations
- **x/sendmanager** - Native coin transfer management
  - `precompile/` - EVM precompile for sending native coins
- **x/maps** - Key-value mapping functionality  
- **x/anchor** - Anchoring and verification system
- **x/wasmx** - Extended WASM functionality

### EVM Integration
The blockchain includes full EVM compatibility via the `cosmos/evm` module:
- **EVM Module** (`github.com/cosmos/evm/x/vm`) - Core EVM functionality
- **ERC20 Module** (`github.com/cosmos/evm/x/erc20`) - ERC20 token wrapping
- **FeeMarket Module** (`github.com/cosmos/evm/x/feemarket`) - Dynamic fee pricing
- **Precompiles** - Custom precompiles for direct module access from Solidity
  - Registration in `app/evm.go:registerCustomPrecompiles()`
  - Enablement via genesis or upgrade handlers
  - Documentation in `app/PRECOMPILE_MANAGEMENT.md`

### Multi-Chain Support
The blockchain supports signatures from multiple chains:
- **Ethereum** - Uses EIP712 signatures (schemas in `chain-handlers/ethereum/ethereum/eip712/schemas.go`)
- **Bitcoin** - JSON schema with alphabetical sorting
- **Solana** - JSON schema with alphabetical sorting
- **Cosmos** - Standard Cosmos signatures

### Directory Structure

#### Core Implementation
- `x/tokenization/` - Tokenization module implementation
  - `keeper/` - Business logic and state management
  - `types/` - Type definitions and validation
  - `module/` - Module initialization and routing
- `x/maps/`, `x/anchor/`, `x/wasmx/` - Other custom modules

#### Protocol Definitions
- `proto/` - Protobuf definitions organized by module
- `api/` - Generated Go types (remove versioned folders after generation)

#### Chain Handlers
- `chain-handlers/` - Multi-chain signature support
  - `ethereum/` - EIP712 signature handling
  - `bitcoin/`, `solana/` - JSON schema signature handling

#### EVM Integration
- `app/evm.go` - EVM module registration and precompile setup
- `app/PRECOMPILE_MANAGEMENT.md` - Precompile registration and enablement guide
- `contracts/` - Solidity contracts, interfaces, and documentation
  - `interfaces/` - Precompile interface definitions
  - `libraries/` - Helper libraries for precompile interaction
  - `docs/` - EVM integration guides and API references

#### Generated Code
- `ts-client/` - TypeScript client generated from protos
- Generated `.pb.go` and `.pulsar.go` files throughout codebase

## Development Workflows

### Adding New Message Types
Follow the comprehensive guide in `_docs/ADDING_NEW_MSG_TYPES.md` which covers:
1. Proto definitions
2. Code generation  
3. EIP712 schema updates
4. WASM bindings
5. CLI commands
6. Tests and validation

### Adding New Proto Fields
Follow the guide in `_docs/PROTO_ADDITION_GUIDE.md` which covers:
1. Proto definition updates
2. Code generation with `ignite generate proto-go --yes`
3. EIP712 schema updates
4. Business logic integration
5. Genesis state handling for new stored types

### Proto Generation Requirements
- Always use `--skip-proto` flag with Ignite commands due to manual proto file corrections
- Remove versioned API folders after generation: `rm -rf api/tokenization/v*`
- Auto-stage generated files: `git add *.pb.go *.pulsar.go`

### Key Development Patterns

#### Store Keys and State Management
- Use unique byte prefixes in `x/tokenization/keeper/keys.go`
- Implement proper marshal/unmarshal in store methods
- Follow incrementing ID patterns for new data types

#### Message Validation
- Implement `ValidateBasic()` for all message types
- Use appropriate Cosmos SDK error types
- Validate all fields including new additions

#### Multi-Chain Compatibility
- Update EIP712 schemas when adding new fields (include all fields with empty defaults)
- Test cross-chain signature verification
- Maintain JSON schema compatibility for Bitcoin/Solana

#### Testing Strategy
- Unit tests for message handlers in `keeper/`
- Integration test helpers in `integration_msg_helpers_test.go`
- Test authorization and edge cases
- Verify genesis state import/export

## Important Notes

### Protocol Buffer Handling
- **Critical**: Always use `--skip-proto` flag with Ignite commands
- Manual corrections exist in generated query files
- Clean up API versioned folders after generation

### State Management
- Collections, maps, and dynamic stores use incrementing IDs
- Creator-based authorization for updates/deletions
- Proper error handling with panic for genesis, errors for runtime

### Cross-Chain Integration
- EIP712 schemas must include ALL possible fields
- Field names in schemas must match proto exactly
- Empty string defaults for string fields, false for booleans

### Custom Types
- Use `Uint` custom type for IDs with gogoproto annotations
- Non-nullable fields use `(gogoproto.nullable) = false`
- String representation in proto, Uint in Go code

### EVM Development
- **Precompile Addresses**: Use 0x1001+ range for custom precompiles (0x0800-0x0806 reserved for Cosmos defaults)
- **Registration**: Register in `app/evm.go:registerCustomPrecompiles()`
- **Enablement**: Add to genesis `active_static_precompiles` or upgrade handler
- **Testing**: Use helpers in `app/precompile_helpers.go` for test setup
- **Documentation**: Update `app/PRECOMPILE_MANAGEMENT.md` when adding new precompiles
- **ERC20 Wrapping**: Native coins can be wrapped via ERC20 keeper for standard ERC20 compatibility