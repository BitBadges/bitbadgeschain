# MsgUpdateDynamicStore

Updates an existing dynamic store's default value.

## Proto Definition

```protobuf
message MsgUpdateDynamicStore {
  string creator = 1; // Address updating the store (must be creator)
  string storeId = 2; // ID of dynamic store to update
  bool defaultValue = 3; // New default value for uninitialized addresses
}

message MsgUpdateDynamicStoreResponse {}
```

## Usage Example

```bash
# CLI command
bitbadgeschaind tx badges update-dynamic-store [store-id] [default-value] --from creator-key
```

### JSON Example

```json
{
    "creator": "bb1...",
    "storeId": "1",
    "defaultValue": true
}
```
