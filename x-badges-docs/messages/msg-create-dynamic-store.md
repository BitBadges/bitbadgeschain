# MsgCreateDynamicStore

Creates a new dynamic store for boolean key-value storage.

## Proto Definition

```protobuf
message MsgCreateDynamicStore {
  string creator = 1; // Address creating the dynamic store
  bool defaultValue = 2; // Default value for uninitialized addresses
}

message MsgCreateDynamicStoreResponse {
  string storeId = 1; // ID of the created dynamic store
}
```

## Usage Example

```bash
# CLI command
bitbadgeschaind tx badges create-dynamic-store [default-value] --from creator-key
```

### JSON Example
```json
{
  "creator": "bb1...",
  "defaultValue": false
}
```