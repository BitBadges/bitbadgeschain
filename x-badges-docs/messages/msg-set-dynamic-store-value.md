# MsgSetDynamicStoreValue

Sets a boolean value for a specific address in a dynamic store.

## Proto Definition

```protobuf
message MsgSetDynamicStoreValue {
  string creator = 1; // Address setting the value (must be store creator)
  string storeId = 2; // ID of the dynamic store
  string address = 3; // Address to set the value for
  bool value = 4; // Boolean value to set (true/false)
}

message MsgSetDynamicStoreValueResponse {}
```

## Usage Example

```bash
# CLI command
bitbadgeschaind tx badges set-dynamic-store-value [store-id] [address] [value] --from creator-key
```

### JSON Example
```json
{
  "creator": "bb1...",
  "storeId": "1",
  "address": "bb1...",
  "value": true
}
```