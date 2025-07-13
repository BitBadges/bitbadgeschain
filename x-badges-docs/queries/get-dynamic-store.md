# GetDynamicStore

Retrieves information about a dynamic store.

## Proto Definition

```protobuf
message QueryGetDynamicStoreRequest {
  string storeId = 1;
}

message QueryGetDynamicStoreResponse {
  DynamicStore store = 1;
}

message DynamicStore {
  // The unique identifier for this dynamic store. This is assigned by the blockchain.
  string storeId = 1 [(gogoproto.customtype) = "Uint", (gogoproto.nullable) = false];
  // The address of the creator of this dynamic store.
  string createdBy = 2;
  // The default value for uninitialized addresses.
  bool defaultValue = 3;
}
```

## Usage Example

```bash
# CLI query
bitbadgeschaind query badges get-dynamic-store [store-id]

# REST API
curl "https://lcd.bitbadges.io/bitbadges/bitbadgeschain/badges/get_dynamic_store/1"
```

### Response Example

```json
{
    "store": {
        "storeId": "1",
        "createdBy": "bb1...",
        "defaultValue": false
    }
}
```
