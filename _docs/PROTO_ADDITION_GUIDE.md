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

-   `./proto/badges/` - Token-related definitions
-   `./proto/maps/` - Map-related definitions
-   `./proto/anchor/` - Anchor-related definitions
-   `./proto/ethereum/` - Ethereum-specific definitions
-   `./proto/bitcoin/` - Bitcoin-specific definitions
-   `./proto/solana/` - Solana-specific definitions

### Key Files to Modify

For token-related changes, the main files are:

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
// Use this pattern for collection IDs, token IDs, and other numeric identifiers
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
-   `badges/TransferTokens`
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

### Step 4.5: Handle New Stored Value Types

If your new fields introduce new types that need to be stored persistently (like new collections, balances, or custom data structures), you need to handle genesis state integration.

#### When This Applies

This step is necessary when you add:

-   New collection types
-   New balance types
-   New address list types
-   New approval tracker types
-   New dynamic store types
-   Any new data structure that needs to persist across chain restarts

#### Required Updates

##### 1. Update Genesis Proto

Add new fields to `proto/badges/genesis.proto`:

```protobuf
message GenesisState {
  // ... existing fields ...

  // New stored value types
  repeated NewDataType newDataTypes = 14;
  string nextNewDataTypeId = 15 [(gogoproto.customtype) = "Uint", (gogoproto.nullable) = false];
  repeated NewDataTypeValue newDataTypeValues = 16;

  // this line is used by starport scaffolding # genesis/proto/state
}
```

**Important**:

-   Use the next available field number (continue from existing sequence)
-   Add imports for new proto files if needed
-   Use `Uint` custom type for ID fields
-   Use `repeated` for collections of data

##### 2. Update Genesis Types

Add default values in `x/badges/types/genesis.go`:

```go
func DefaultGenesis() *GenesisState {
    return &GenesisState{
        PortId: PortID,
        // this line is used by starport scaffolding # genesis/types/default
        Params:               DefaultParams(),
        NextCollectionId:     types.NewUint(1),
        NextNewDataTypeId:    types.NewUint(1),  // Add default for new type
    }
}
```

##### 3. Update Genesis Module

Add initialization and export logic in `x/badges/module/genesis.go`:

```go
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
    // ... existing initialization ...

    // Set next new data type ID if defined; default 0
    if genState.NextNewDataTypeId.Equal(sdkmath.NewUint(0)) {
        genState.NextNewDataTypeId = sdkmath.NewUint(1)
    }
    k.SetNextNewDataTypeId(ctx, genState.NextNewDataTypeId)

    // Initialize new data types
    for _, newDataType := range genState.NewDataTypes {
        if err := k.SetNewDataTypeInStore(ctx, *newDataType); err != nil {
            panic(err)
        }
    }

    // Initialize new data type values
    for _, newDataTypeValue := range genState.NewDataTypeValues {
        if err := k.SetNewDataTypeValueInStore(ctx, newDataTypeValue.Id, newDataTypeValue.Address, newDataTypeValue.Value); err != nil {
            panic(err)
        }
    }
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
    genesis := types.DefaultGenesis()
    // ... existing export logic ...

    genesis.NextNewDataTypeId = k.GetNextNewDataTypeId(ctx)
    genesis.NewDataTypes = k.GetNewDataTypesFromStore(ctx)
    genesis.NewDataTypeValues = k.GetAllNewDataTypeValuesFromStore(ctx)

    return genesis
}
```

##### 4. Add Store Methods

Implement storage methods in `x/badges/keeper/store.go`:

```go
// Set new data type in store
func (k Keeper) SetNewDataTypeInStore(ctx sdk.Context, newDataType types.NewDataType) error {
    marshaled_data, err := k.cdc.Marshal(&newDataType)
    if err != nil {
        return sdkerrors.Wrap(err, "Marshal types.NewDataType failed")
    }

    storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
    store := prefix.NewStore(storeAdapter, []byte{})
    store.Set(newDataTypeStoreKey(newDataType.Id), marshaled_data)
    return nil
}

// Get new data type from store
func (k Keeper) GetNewDataTypeFromStore(ctx sdk.Context, id sdkmath.Uint) (types.NewDataType, bool) {
    storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
    store := prefix.NewStore(storeAdapter, []byte{})
    marshaled_data := store.Get(newDataTypeStoreKey(id))

    var newDataType types.NewDataType
    if len(marshaled_data) == 0 {
        return newDataType, false
    }
    k.cdc.MustUnmarshal(marshaled_data, &newDataType)
    return newDataType, true
}

// Get all new data types from store
func (k Keeper) GetNewDataTypesFromStore(ctx sdk.Context) (newDataTypes []*types.NewDataType) {
    storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
    store := prefix.NewStore(storeAdapter, []byte{})
    iterator := storetypes.KVStorePrefixIterator(store, NewDataTypeKey)
    defer iterator.Close()
    for ; iterator.Valid(); iterator.Next() {
        var newDataType types.NewDataType
        k.cdc.MustUnmarshal(iterator.Value(), &newDataType)
        newDataTypes = append(newDataTypes, &newDataType)
    }
    return
}

// Get all new data type values from store
func (k Keeper) GetAllNewDataTypeValuesFromStore(ctx sdk.Context) (newDataTypeValues []*types.NewDataTypeValue) {
    storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
    store := prefix.NewStore(storeAdapter, []byte{})
    iterator := storetypes.KVStorePrefixIterator(store, NewDataTypeValueKey)
    defer iterator.Close()
    for ; iterator.Valid(); iterator.Next() {
        var newDataTypeValue types.NewDataTypeValue
        k.cdc.MustUnmarshal(iterator.Value(), &newDataTypeValue)
        newDataTypeValues = append(newDataTypeValues, &newDataTypeValue)
    }
    return
}
```

##### 5. Add Store Keys

Add key definitions in `x/badges/keeper/keys.go`:

```go
var (
    // ... existing keys ...
    NewDataTypeKey      = []byte{0x10}  // Use next available byte
    NewDataTypeValueKey = []byte{0x11}  // Use next available byte
)

// Store key functions
func newDataTypeStoreKey(id sdkmath.Uint) []byte {
    key := make([]byte, len(NewDataTypeKey)+IDLength)
    copy(key, NewDataTypeKey)
    copy(key[len(NewDataTypeKey):], []byte(id.String()))
    return key
}

func newDataTypeValueStoreKey(id sdkmath.Uint, address string) []byte {
    key := make([]byte, len(NewDataTypeValueKey)+IDLength+len(address))
    copy(key, NewDataTypeValueKey)
    copy(key[len(NewDataTypeValueKey):], []byte(id.String()))
    copy(key[len(NewDataTypeValueKey)+IDLength:], []byte(address))
    return key
}
```

##### 6. Regenerate Proto Code

After updating genesis.proto:

```bash
ignite generate proto-go --yes
# Remove versioned API folders
ls api/badges/
rm -rf api/badges/v*
# Build verification
go build ./cmd/bitbadgeschaind
# Auto-stage generated files
git add *.pb.go *.pulsar.go
```

##### 7. Test Genesis Integration

Test that genesis state works correctly:

```bash
# Test genesis module
go test ./x/badges/module/... -v

# Test types
go test ./x/badges/types/... -v

# Test keeper (if you have tests)
go test ./x/badges/keeper/... -v
```

#### Common Patterns

-   **Incrementing IDs**: Most new data types use incrementing IDs starting from 1
-   **Creator-based access**: Only creators can update/delete their data
-   **Value storage**: For data with per-address values, use separate value storage
-   **Error handling**: Use panic for genesis errors, return errors for runtime operations
-   **Key prefixes**: Use unique byte prefixes to avoid key collisions

#### Example: Dynamic Stores

See the implementation of dynamic stores for a complete example:

-   `proto/badges/dynamic_stores.proto` - Data structure definition
-   `proto/badges/genesis.proto` - Genesis state fields
-   `x/badges/keeper/store.go` - Storage methods
-   `x/badges/keeper/keys.go` - Key definitions
-   `x/badges/module/genesis.go` - Genesis integration

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
