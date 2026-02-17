# Adding New Message Types to x/tokenization

This guide documents the process for adding new message types to the tokens module, including all necessary steps, gotchas, and commands.

## Overview

When adding new message types to the tokens module, you need to follow a specific sequence of steps to ensure proper integration with the Cosmos SDK, WASM bindings, CLI, and testing infrastructure.

## Step-by-Step Process

### 1. Add Proto Definitions

**Location**: `proto/tokenization/`

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

### 2. Update Genesis Files

**Location**: `proto/tokenization/` and `x/tokenization/`

**Files to modify**:

-   `proto/tokenization/genesis.proto` - Add new fields for state persistence
-   `x/tokenization/types/genesis.go` - Update default genesis state
-   `x/tokenization/module/genesis.go` - Add initialization and export logic

**Commands**:

```bash
# Generate Go code from proto definitions
ignite generate proto-go
```

**Gotchas**:

-   Add imports for new proto files in genesis.proto
-   Use appropriate field numbers (continue from existing sequence)
-   Add default values in types/genesis.go
-   Handle initialization and export in module/genesis.go
-   Include proper error handling for store operations

### 3. Update Badges Custom Message Types

**Location**: `x/tokenization/types/`

**Files to modify**:

-   Add new message types to the custom message wrapper types
-   Ensure proper JSON marshaling/unmarshaling

**Gotchas**:

-   Custom message types must be properly integrated with WASM bindings
-   All fields must be serializable

### 4. Update Custom Transaction Handler

**Location**: `custom-bindings/custom_tx.go`

**Purpose**: Enable WASM contracts to call the new message types

**Gotchas**:

-   Must set the `Creator` field to the sender address
-   Handle all message variants in the switch statement
-   Return proper error types

### 6. Add CLI Commands

**Location**: `x/tokenization/client/cli/`

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

### 7. Add Integration Test Helpers

**Location**: `x/tokenization/keeper/integration_msg_helpers_test.go`

**Purpose**: Provide helper functions for testing the new message types

**Gotchas**:

-   Follow existing patterns for error handling
-   Include proper validation calls
-   Use consistent naming conventions

### 8. Implement Core Functionality

**Location**: `x/tokenization/keeper/`

**Files to create/modify**:

-   `msg_server_<action>.go` - Main message handler
-   `msg_server_<action>_test.go` - Unit tests
-   `keeper.go` - Add any new keeper methods
-   `keys.go` - Add new store keys if needed

**Commands**:

```bash
# Run tests
go test ./x/tokenization/keeper/...
```

**Gotchas**:

-   Implement proper authorization checks
-   Handle all edge cases
-   Follow existing patterns for state management
-   Use proper error types from the types package

### 9. Add Simulation Support

**Location**: `x/tokenization/simulation/`

**Files to create/modify**:

-   `<action>.go` - Simulation function for the new message type
-   `simulation_test.go` - Add test case for the new simulation function

**Purpose**: Enable the message type to be used in Cosmos SDK simulations (stress testing, fuzzing, etc.)

**Gotchas**:

-   **IMPORTANT**: If you modify an existing message type (add/remove fields), you MUST update the corresponding simulation file
-   Follow the pattern: `SimulateMsg<Action><Entity>`
-   Use random values for all fields
-   Handle cases where required resources don't exist (return `NoOpMsg`)
-   For update operations, verify the resource exists and use the correct creator account
-   Add the simulation function to `simulation_test.go` to ensure it doesn't panic

**Example**:

```go
func SimulateMsgUpdateDynamicStore(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// Implementation...
	}
}
```

**Commands**:

```bash
# Run simulation tests
go test ./x/tokenization/simulation/...
```

### 11. Register in Codec

**Location**: `x/tokenization/types/codec.go`

**Purpose**: Register new message types with the Cosmos SDK codec

**Gotchas**:

-   Must register in both `RegisterCodec` and `RegisterInterfaces`
-   Use proper concrete type names
-   Follow existing naming patterns

### 12. Ensure GetSignBytes Uses AminoCdc

**Location**: `x/tokenization/types/`

**Purpose**: Ensure proper signature generation for transactions

**Gotchas**:

-   Verify that `AminoCdc` is used for signing
-   Test signature verification

### 13. Update Precompile (if EVM-callable)

**Location**: `x/tokenization/precompile/`

**Purpose**: Enable Solidity smart contracts to call the new message type

**When Required**: If the new message type should be callable from EVM/Solidity contracts

**Files to modify**:

-   `x/tokenization/precompile/abi.json` - Add method signature
-   `x/tokenization/precompile/precompile.go` - Add method constant, gas cost, handler
-   `contracts/interfaces/ITokenizationPrecompile.sol` - Add Solidity interface
-   `contracts/libraries/TokenizationJSONHelpers.sol` - Add JSON builder helper (if needed)
-   `contracts/libraries/TokenizationWrappers.sol` - Add wrapper function (if needed)

**Steps**:

1. **Add to `abi.json`**:
   ```json
   {
     "name": "newMethodName",
     "type": "function",
     "inputs": [{"name": "jsonInput", "type": "string"}],
     "outputs": [{"name": "success", "type": "bool"}]
   }
   ```

2. **Add method constant to `precompile.go`**:
   ```go
   const (
       // ... existing methods ...
       MethodNewMethodName = "newMethodName"
   )
   ```

3. **Add gas cost constant**:
   ```go
   const (
       GasNewMethodName = 50000 // Adjust based on complexity
   )
   ```

4. **Add case to `RequiredGas` switch**

5. **Add case to `Execute` switch**

6. **Implement method handler** with proper validation

7. **Update Solidity interface** in `contracts/interfaces/ITokenizationPrecompile.sol`

8. **Add tests** in `x/tokenization/precompile/test/`

**Gotchas**:

-   Use JSON input format matching protobuf structure
-   Implement proper gas estimation
-   Add DoS protection limits for arrays/strings
-   Update both unit and integration tests

**Reference**: See `x/tokenization/precompile/README.md` and `app/PRECOMPILE_MANAGEMENT.md` for detailed guidance

### 14. Update SDK Precompile Support (if EVM-callable)

**Location**: `../bitbadgesjs/packages/bitbadgesjs-sdk/src/transactions/precompile/`

**Purpose**: Enable TypeScript/JavaScript applications to call the new precompile method

**Files to modify**:

-   `abi.json` - Copy updated ABI from chain (must match `x/tokenization/precompile/abi.json`)
-   `function-mapper.ts` - Add mapping from SDK message type to precompile function
-   `type-detector.ts` - Add type detection for new message (if needed)
-   `data-converter.ts` - Add data conversion logic (if needed)

**Steps**:

1. **Copy ABI from chain**:
   ```bash
   cp ../bitbadgeschain/x/tokenization/precompile/abi.json \
      ../bitbadgesjs/packages/bitbadgesjs-sdk/src/transactions/precompile/abi.json
   ```

2. **Update `function-mapper.ts`** to map new message type:
   ```typescript
   case 'MsgNewAction':
     return PrecompileFunction.NEW_ACTION;
   ```

3. **Update `type-detector.ts`** if adding new message type detection

4. **Rebuild SDK**:
   ```bash
   cd ../bitbadgesjs/packages/bitbadgesjs-sdk
   bun run build
   ```

**Gotchas**:

-   ABI must be identical between chain and SDK
-   Function names in SDK must match ABI method names exactly
-   Rebuild SDK after changes for symlink to pick up updates

## Testing Checklist

-   [ ] Unit tests for message handlers
-   [ ] Integration tests with other modules
-   [ ] CLI command tests
-   [ ] WASM binding tests
-   [ ] Authorization tests
-   [ ] Edge case handling
-   [ ] Genesis state tests
-   [ ] **Simulation tests** (if adding new message or modifying existing message fields)
-   [ ] **Precompile tests** (if message is EVM-callable)
-   [ ] **SDK precompile tests** (if message is EVM-callable and SDK support added)

## Common Issues and Solutions

### Proto Generation Fails

-   Check field numbers are unique
-   Ensure all imports are correct
-   Verify syntax is valid protobuf

### Genesis State Issues

-   Ensure all new fields have proper default values
-   Check that initialization and export logic is complete
-   Verify store operations handle errors properly

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

See the implementation of existing message types like `MsgTransferTokens`, `MsgCreateCollection`, etc. for reference patterns.

## Commands Reference

```bash
# Generate proto code
ignite generate proto-go --yes

# Clean up versioned API folders (IMPORTANT!)
ls api/tokenization/
rm -rf api/tokenization/v*

# Build binary
make install

# Build verification
go build ./cmd/bitbadgeschaind

# Auto-stage generated files
git add *.pb.go *.pulsar.go

# Run tests
go test ./x/tokenization/keeper/...

# Run specific test
go test ./x/tokenization/keeper/ -run TestMsgCreateDynamicStore

# Check for linting issues
golangci-lint run ./x/tokenization/...

# Format code
go fmt ./x/tokenization/...
```

## Implementation Status for Dynamic Store Messages

### ✅ Completed Steps:

1. **Proto Definitions**: Added `dynamic_stores.proto`, updated `tx.proto` and `query.proto`
2. **Genesis Files**: Updated `genesis.proto`, `types/genesis.go`, and `module/genesis.go`
3. **Custom Transaction Handler**: Updated `custom_tx.go` with dynamic store handlers
4. **CLI Commands**: Created CLI commands for create, update, delete, and query operations
5. **Integration Test Helpers**: Added helper functions in `integration_msg_helpers_test.go`
7. **Core Functionality**: Implemented message handlers and store methods
8. **Codec Registration**: Registered new message types in `codec.go`

### 🔄 Remaining Steps:

9. **Write Tests**: Need to create comprehensive unit tests
10. **Documentation**: Update API documentation

### 📝 Notes:

-   Dynamic stores use incrementing IDs like collections
-   Only creators can update/delete their dynamic stores
-   Proper authorization checks are implemented
-   Events are emitted for all operations
-   **DynamicStore no longer has data or metadata fields.**
-   **ValidateBasic errs on the side of caution (checks for empty/invalid fields).**
-   **Test files for Msg types are present in both `/keeper` and `/types` as `msg_*_test.go`.**
-   **Genesis state now includes dynamic stores and next dynamic store ID.**
-   **Global Kill Switch**: All dynamic stores now have a `globalEnabled` field that acts as a global kill switch. When `globalEnabled = false`, all approvals using that store via `DynamicStoreChallenge` will fail immediately. This enables quick halting of approvals (e.g., when a 2FA protocol is compromised). See `_docs/PROTO_ADDITION_GUIDE.md` for details on the implementation.

## Next Steps

After completing all steps:

1. Create a pull request
2. Include comprehensive tests
3. Update documentation
4. Consider backward compatibility
5. Review security implications
6. **Update simulation files** if message fields changed (see Step 9)
7. **Update CLI commands** if message fields changed (see Step 6)

## Important Reminders

### When Modifying Existing Messages

If you modify an existing message type (add, remove, or change fields), you MUST update:

1. ✅ **CLI commands** (`x/tokenization/client/cli/tx_*.go`) - Update argument parsing
2. ✅ **Simulation files** (`x/tokenization/simulation/*.go`) - Update simulation functions
3. ✅ **Tests** - Update existing tests and add new ones for new fields
4. ✅ **Documentation** - Update relevant docs in `_docs/` folder
5. ✅ **Precompile** (if EVM-callable) - Update `abi.json`, handlers, and Solidity interfaces in `x/tokenization/precompile/` and `contracts/`
6. ✅ **SDK Precompile** (if EVM-callable) - Copy updated `abi.json` and update function mappers in `bitbadgesjs-sdk/src/transactions/precompile/`

These are often forgotten but critical for maintaining a complete and tested codebase.
