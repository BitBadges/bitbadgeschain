# Guide: Adding New Fields to Proto Definitions

This guide explains how to add new fields to protobuf definitions in the BitBadges blockchain and handle all related updates.

## Overview

When adding new fields to proto definitions, you need to follow a systematic approach to ensure all generated code, type definitions, and schemas are properly updated. This process involves four main steps:

1. **Update proto definitions** in `./proto` directory
2. **Generate Go code** using `ignite generate proto-go --yes`
3. **Update EIP712 schemas** in `chain-handlers/ethereum/ethereum/eip712/schemas.go`
4. **Handle business logic** in `x/badges` module (if required)

## Step 1: Update Proto Definitions

### Location

Proto definitions are located in the `./proto` directory, organized by module:

-   `./proto/badges/` - Badge-related definitions
-   `./proto/maps/` - Map-related definitions
-   `./proto/anchor/` - Anchor-related definitions
-   `./proto/ethereum/` - Ethereum-specific definitions
-   `./proto/bitcoin/` - Bitcoin-specific definitions
-   `./proto/solana/` - Solana-specific definitions

### Key Files to Modify

For badge-related changes, the main files are:

-   `./proto/badges/tx.proto` - Transaction messages
-   `./proto/badges/query.proto` - Query messages
-   `./proto/badges/transfers.proto` - Transfer-related types
-   `./proto/badges/balances.proto` - Balance-related types
-   `./proto/badges/permissions.proto` - Permission-related types
-   `./proto/badges/collections.proto` - Collection-related types
-   `./proto/badges/metadata.proto` - Metadata-related types
-   `./proto/badges/timelines.proto` - Timeline-related types
-   `./proto/badges/address_lists.proto` - Address list types

### Adding New Fields

When adding a new field to a message, follow these conventions:

```protobuf
message ExampleMessage {
  // Existing fields...
  string existing_field = 1;

  // New field - use the next available field number
  string newField = 2;  // Use camelCase for field names

  // For number types, use the custom Uint type with gogoproto annotations
  string collectionId = 3 [(gogoproto.customtype) = "Uint", (gogoproto.nullable) = false];

  // For optional fields, use the appropriate wrapper or add a boolean flag
  bool hasNewField = 4;
  string newFieldValue = 5;
}
```

### Field Naming Conventions

-   **Use camelCase** for field names (e.g., `newField`, `collectionId`, `denomUnits`)
-   **Use descriptive names** that clearly indicate the field's purpose
-   **Follow existing patterns** in the codebase for consistency

### Field Numbering

-   **Never reuse field numbers** - always use the next available number
-   **Don't skip numbers** unless you have a specific reason
-   **Use descriptive field names** in camelCase

### Number Type Conventions

For number types that need to be handled as Uint in Go:

```protobuf
// Use this pattern for collection IDs, badge IDs, and other numeric identifiers
string collectionId = 1 [(gogoproto.customtype) = "Uint", (gogoproto.nullable) = false];
string badgeId = 2 [(gogoproto.customtype) = "Uint", (gogoproto.nullable) = false];
```

This pattern:

-   Declares the field as `string` in proto
-   Uses `gogoproto.customtype = "Uint"` to generate Go code with the `Uint` type
-   Uses `gogoproto.nullable = false` to ensure the field is not a pointer
-   Automatically handles conversion between string and Uint types

### Import Considerations

If your new field uses types from other proto files, ensure proper imports:

```protobuf
import "badges/transfers.proto";
import "badges/balances.proto";
import "cosmos/base/v1beta1/coin.proto";
```

## Step 2: Generate Go Code

After updating proto definitions, regenerate the Go code:

```bash
ignite generate proto-go --yes
```

This command will:

-   Generate Go structs from proto messages
-   Create gRPC client and server code
-   Update type definitions in `api/` directory
-   Generate validation code

**Note**: The `--yes` flag automatically answers "yes" to any prompts, which is useful in automated environments or when you want to proceed without manual confirmation.

### What Gets Generated

The generated code appears in:

-   `api/badges/` - Generated Go types for badges module
-   `api/maps/` - Generated Go types for maps module
-   `api/anchor/` - Generated Go types for anchor module
-   `x/badges/types/` - Additional type definitions

### API Directory Cleanup

**Important**: After generation, remove all versioned folders from the API directory:

```bash
# First, list what's in the api/badges/ directory to see the v* folders
ls api/badges/

# Then explicitly remove the versioned directories (e.g., v6, v7, etc.)
rm -rf api/badges/v6
rm -rf api/badges/v7
# ... remove any other v* directories you see

# Or if you're sure about the pattern, you can use:
# rm -rf api/badges/v*
```

**Rationale**: We use v6 types for migration handling, but we shouldn't include them in the API. The versioned folders are only needed for internal migration logic, not for external API consumption.

**Safety Note**: Using `ls` first helps avoid accidentally deleting files that might start with 'v' but aren't version directories.

### Verification

After generation, verify that:

-   New fields appear in the generated Go structs
-   No compilation errors exist
-   Field types are correctly mapped

### Build Verification

After generating the code, test that the application builds successfully:

```bash
go build ./cmd/bitbadgeschaind
```

This command builds the main blockchain application and will catch any compilation errors introduced by the proto changes.

### Auto-Stage Generated Files

After successful generation and build verification, automatically stage the generated proto files:

```bash
# Stage all generated proto files
git add *.pb.go *.pulsar.go

# Or stage specific directories if needed
git add api/badges/*.pulsar.go
git add x/badges/types/*.pb.go
```

**Note**: This automatically stages all generated protobuf files (_.pb.go and _.pulsar.go) so you don't have to manually track which files were generated. This is especially useful since these files are auto-generated and should always be committed together with proto changes.

## Step 3: Update EIP712 Schemas

### Purpose

EIP712 schemas are used for Ethereum signature verification. They define the structure of messages that can be signed by Ethereum wallets.

### Location

Schemas are defined in `chain-handlers/ethereum/ethereum/eip712/schemas.go`

### Adding New Schema Entries

For each new message type or field, you need to add a corresponding schema entry. The schema should include **all possible fields** with empty/default values to ensure proper type generation.

#### Schema Structure

```go
schemas = append(schemas, `{
    "type": "badges/NewMessageType",
    "value": {
        "creator": "",
        "newField": "",
        "optionalField": false,
        "arrayField": [],
        "objectField": {
            "nestedField": ""
        }
    }
}`)
```

#### Important Notes

-   **Include all fields** - even optional ones
-   **Use empty strings** for string fields (`""`)
-   **Use false** for boolean fields
-   **Use empty arrays** for repeated fields (`[]`)
-   **Use empty objects** for nested message types (`{}`)
-   **Follow the exact field names** from the proto definition

#### Example: Adding a New Field to Existing Message

If you add a new field to `MsgCreateCollection`, you need to update the corresponding schema:

```go
// Before
schemas = append(schemas, `{
    "type": "badges/CreateCollection",
    "value": {
        "creator": "",
        "collectionId": "",
        // ... existing fields
    }
}`)

// After
schemas = append(schemas, `{
    "type": "badges/CreateCollection",
    "value": {
        "creator": "",
        "collectionId": "",
        "newField": "",  // Add the new field
        // ... existing fields
    }
}`)
```

#### Schema Locations to Update

Common message types that need schema updates:

-   `badges/CreateCollection`
-   `badges/UpdateCollection`
-   `badges/UniversalUpdateCollection`
-   `badges/TransferBadges`
-   `badges/UpdateUserApprovals`
-   `badges/CreateAddressLists`
-   `maps/CreateMap`
-   `maps/UpdateMap`
-   `maps/SetValue`

## Step 4: Handle Business Logic

### Location

Business logic is implemented in the `x/badges` module:

-   `x/badges/keeper/` - Core business logic
-   `x/badges/types/` - Type definitions and validation
-   `x/badges/module/` - Module initialization and routing

### Common Updates Needed

#### 1. Message Handler Updates

If you added a new message type, implement the handler in `x/badges/keeper/msg_server.go`:

```go
func (k msgServer) NewMessageType(goCtx context.Context, msg *types.MsgNewMessageType) (*types.MsgNewMessageTypeResponse, error) {
    ctx := sdk.UnwrapSDKContext(goCtx)

    // Implement your business logic here

    return &types.MsgNewMessageTypeResponse{}, nil
}
```

#### 2. Validation Logic

Add validation for new fields in `x/badges/types/` files:

```go
func (msg *MsgNewMessageType) ValidateBasic() error {
    if msg.Creator == "" {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "creator cannot be empty")
    }

    // Add validation for new fields
    if msg.NewField == "" {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "new field cannot be empty")
    }

    return nil
}
```

#### 3. State Management

If new fields affect state, update the keeper methods:

```go
func (k Keeper) SetNewField(ctx sdk.Context, collectionId string, newField string) {
    store := ctx.KVStore(k.storeKey)
    key := types.GetNewFieldKey(collectionId)
    store.Set(key, []byte(newField))
}

func (k Keeper) GetNewField(ctx sdk.Context, collectionId string) string {
    store := ctx.KVStore(k.storeKey)
    key := types.GetNewFieldKey(collectionId)
    return string(store.Get(key))
}
```

#### 4. Query Handlers

If new fields are queryable, add query handlers in `x/badges/keeper/query.go`:

```go
func (k Keeper) NewField(goCtx context.Context, req *types.QueryNewFieldRequest) (*types.QueryNewFieldResponse, error) {
    if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }

    ctx := sdk.UnwrapSDKContext(goCtx)
    newField := k.GetNewField(ctx, req.CollectionId)

    return &types.QueryNewFieldResponse{
        NewField: newField,
    }, nil
}
```

## Testing Your Changes

### 1. Compilation Test

```bash
go build ./...
```

### 2. Unit Tests

```bash
go test ./x/badges/...
```

### 3. Integration Tests

```bash
ignite chain test
```

### 4. EIP712 Signature Test

Test that new fields work with Ethereum signatures:

```bash
# Build and run the chain
ignite chain serve --skip-proto

# Test EIP712 signing with new fields
# (Use your preferred testing method)
```

## Common Pitfalls

### 1. Field Number Conflicts

-   **Problem**: Reusing field numbers causes runtime errors
-   **Solution**: Always use the next available field number

### 2. Missing Schema Entries

-   **Problem**: New fields don't work with Ethereum signatures
-   **Solution**: Add corresponding schema entries with empty values

### 3. Incomplete Validation

-   **Problem**: Invalid data can be stored in state
-   **Solution**: Implement proper validation in `ValidateBasic()` methods

### 4. Missing Query Handlers

-   **Problem**: New fields can't be queried
-   **Solution**: Add appropriate query handlers and update query.proto

### 5. State Key Conflicts

-   **Problem**: New state keys conflict with existing ones
-   **Solution**: Use unique key prefixes and update key generation functions

### 6. Forgetting API Cleanup

-   **Problem**: Versioned API folders are included in builds
-   **Solution**: Always remove `api/badges/v*` folders after generation

## Example: Complete Field Addition

Here's a complete example of adding a new `description` field to `MsgCreateCollection`:

### 1. Update Proto

```protobuf
// In proto/badges/tx.proto
message MsgCreateCollection {
  string creator = 1;
  string collectionId = 2;
  string description = 3;  // New field
  // ... other fields
}
```

### 2. Generate Code

```bash
ignite generate proto-go --yes
# First, list what's in the api/badges/ directory to see the v* folders
ls api/badges/
# Then explicitly remove the versioned directories (e.g., v6, v7, etc.)
rm -rf api/badges/v6
rm -rf api/badges/v7
# ... remove any other v* directories you see

# Build verification
go build ./cmd/bitbadgeschaind

# Auto-stage generated files
git add *.pb.go *.pulsar.go
```

### 3. Update Schema

```go
// In chain-handlers/ethereum/ethereum/eip712/schemas.go
schemas = append(schemas, `{
    "type": "badges/CreateCollection",
    "value": {
        "creator": "",
        "collectionId": "",
        "description": "",  // Add new field
        // ... other fields
    }
}`)
```

### 4. Update Business Logic

```go
// In x/badges/types/tx.pb.go (generated)
// The field will be automatically added to the struct

// In x/badges/types/message_create_collection.go
func (msg *MsgCreateCollection) ValidateBasic() error {
    if msg.Creator == "" {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "creator cannot be empty")
    }

    // Add validation for description
    if len(msg.Description) > 1000 {
        return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "description too long")
    }

    return nil
}
```

## Conclusion

Following this systematic approach ensures that new fields are properly integrated into all parts of the system:

-   Proto definitions define the structure
-   Generated code provides type safety
-   EIP712 schemas enable Ethereum compatibility
-   Business logic handles validation and state management
-   API cleanup prevents versioned folders from being included

Always test thoroughly after making changes, especially with Ethereum signature verification, as this is critical for cross-chain functionality.
