# Adding New Message Types to x/badges

This guide documents the process for adding new message types to the badges module, including all necessary steps, gotchas, and commands.

## Overview

When adding new message types to the badges module, you need to follow a specific sequence of steps to ensure proper integration with the Cosmos SDK, WASM bindings, CLI, and testing infrastructure.

## Step-by-Step Process

### 1. Add Proto Definitions

**Location**: `proto/badges/`

**Files to modify**:

-   `tx.proto` - Add new message types
-   `query.proto` - Add new query types (if needed)
-   `dynamic_stores.proto` - Add new data structures (if needed)

**Commands**:

```bash
# Generate Go code from proto definitions
ignite generate proto-go
```

**Gotchas**:

-   Ensure all new types are properly defined with correct field numbers
-   Use appropriate Cosmos SDK types (e.g., `string` for addresses, `uint64` for IDs)
-   Follow naming conventions: `Msg<Action><Entity>` for messages
-   Add imports for new proto files in tx.proto and query.proto
-   Add service methods in the Msg service
-   Add message types to BadgeCustomMsgType for WASM bindings
-   Add query methods in the Query service if needed

### 2. Update Badges Custom Message Types

**Location**: `x/badges/types/`

**Files to modify**:

-   Add new message types to the custom message wrapper types
-   Ensure proper JSON marshaling/unmarshaling

**Gotchas**:

-   Custom message types must be properly integrated with WASM bindings
-   All fields must be serializable

### 3. Add to EIP712 Schemas

**Location**: `chain-handlers/ethereum/ethereum/eip712/schemas.go`

**Purpose**: Enable Ethereum EIP712 signature support for the new message types

**Gotchas**:

-   All fields must be represented as strings in the schema
-   Include all optional fields with empty string defaults
-   Follow the exact structure of existing schemas

### 4. Update Custom Transaction Handler

**Location**: `custom-bindings/custom_tx.go`

**Purpose**: Enable WASM contracts to call the new message types

**Gotchas**:

-   Must set the `Creator` field to the sender address
-   Handle all message variants in the switch statement
-   Return proper error types

### 5. Add CLI Commands

**Location**: `x/badges/client/cli/`

**Files to create/modify**:

-   `tx_<action>.go` - Transaction commands
-   `query_<action>.go` - Query commands (if needed)

**Commands**:

```bash
# Build the binary to test CLI commands
make install
```

**Gotchas**:

-   Use proper flag names and descriptions
-   Handle JSON input properly
-   Include proper validation

### 6. Add Integration Test Helpers

**Location**: `x/badges/keeper/integration_msg_helpers_test.go`

**Purpose**: Provide helper functions for testing the new message types

**Gotchas**:

-   Follow existing patterns for error handling
-   Include proper validation calls
-   Use consistent naming conventions

### 7. Implement Core Functionality

**Location**: `x/badges/keeper/`

**Files to create/modify**:

-   `msg_server_<action>.go` - Main message handler
-   `msg_server_<action>_test.go` - Unit tests
-   `keeper.go` - Add any new keeper methods
-   `keys.go` - Add new store keys if needed

**Commands**:

```bash
# Run tests
go test ./x/badges/keeper/...
```

**Gotchas**:

-   Implement proper authorization checks
-   Handle all edge cases
-   Follow existing patterns for state management
-   Use proper error types from the types package

### 8. Register in Codec

**Location**: `x/badges/types/codec.go`

**Purpose**: Register new message types with the Cosmos SDK codec

**Gotchas**:

-   Must register in both `RegisterCodec` and `RegisterInterfaces`
-   Use proper concrete type names
-   Follow existing naming patterns

### 9. Ensure GetSignBytes Uses AminoCdc

**Location**: `x/badges/types/`

**Purpose**: Ensure proper signature generation for transactions

**Gotchas**:

-   Verify that `AminoCdc` is used for signing
-   Test signature verification

## Testing Checklist

-   [ ] Unit tests for message handlers
-   [ ] Integration tests with other modules
-   [ ] CLI command tests
-   [ ] WASM binding tests
-   [ ] EIP712 signature tests
-   [ ] Authorization tests
-   [ ] Edge case handling

## Common Issues and Solutions

### Proto Generation Fails

-   Check field numbers are unique
-   Ensure all imports are correct
-   Verify syntax is valid protobuf

### WASM Binding Issues

-   Ensure custom message types are properly defined
-   Check that all fields are serializable
-   Verify creator field is set correctly

### CLI Command Not Found

-   Check that commands are properly registered
-   Verify flag names and descriptions
-   Ensure proper error handling

### Test Failures

-   Check authorization logic
-   Verify state management
-   Ensure proper error types are used

## Example Implementation

See the implementation of existing message types like `MsgTransferBadges`, `MsgCreateCollection`, etc. for reference patterns.

## Commands Reference

```bash
# Generate proto code
ignite generate proto-go

# Build binary
make install

# Run tests
go test ./x/badges/keeper/...

# Run specific test
go test ./x/badges/keeper/ -run TestMsgCreateDynamicStore

# Check for linting issues
golangci-lint run ./x/badges/...

# Format code
go fmt ./x/badges/...
```

## Implementation Status for Dynamic Store Messages

### ✅ Completed Steps:

1. **Proto Definitions**: Added `dynamic_stores.proto`, updated `tx.proto` and `query.proto`
2. **EIP712 Schemas**: Added schemas for all three dynamic store message types
3. **Custom Transaction Handler**: Updated `custom_tx.go` with dynamic store handlers
4. **CLI Commands**: Created CLI commands for create, update, delete, and query operations
5. **Integration Test Helpers**: Added helper functions in `integration_msg_helpers_test.go`
6. **Core Functionality**: Implemented message handlers and store methods
7. **Codec Registration**: Registered new message types in `codec.go`

### 🔄 Remaining Steps:

8. **Write Tests**: Need to create comprehensive unit tests
9. **Genesis Integration**: May need to add dynamic store support to genesis
10. **Documentation**: Update API documentation

### 📝 Notes:

-   Dynamic stores use incrementing IDs like collections
-   Only creators can update/delete their dynamic stores
-   Proper authorization checks are implemented
-   Events are emitted for all operations
-   **DynamicStore no longer has data or metadata fields.**
-   **ValidateBasic errs on the side of caution (checks for empty/invalid fields).**
-   **Test files for Msg types are present in both `/keeper` and `/types` as `msg_*_test.go`.**

## Next Steps

After completing all steps:

1. Create a pull request
2. Include comprehensive tests
3. Update documentation
4. Consider backward compatibility
5. Review security implications
 