# MsgDeleteDynamicStore

Deletes a dynamic store.

## Proto Definition

```protobuf
message MsgDeleteDynamicStore {
  string creator = 1; // Address deleting the store (must be creator)
  string storeId = 2; // ID of dynamic store to delete
}

message MsgDeleteDynamicStoreResponse {}
```

## Usage Example

```bash
# CLI command
bitbadgeschaind tx badges delete-dynamic-store [store-id] --from creator-key
```

### JSON Example

```json
{
    "creator": "bb1...",
    "storeId": "1"
}
```
