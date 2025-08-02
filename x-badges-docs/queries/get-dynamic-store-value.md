# GetDynamicStoreValue

Retrieves the numeric value for a specific address in a dynamic store. This is the number of uses an address has left for a dynamic store.

## Proto Definition

```protobuf
message QueryGetDynamicStoreValueRequest {
  string storeId = 1; // ID of dynamic store to query
  string address = 2; // Address to get value for
}

message QueryGetDynamicStoreValueResponse {
  DynamicStoreValue value = 1;
}

message DynamicStoreValue {
  string storeId = 1; // The dynamic store ID
  string address = 2; // The address this value applies to
  string value = 3; // The numeric value
}
```

## Usage Example

```bash
# CLI query
bitbadgeschaind query badges get-dynamic-store-value [store-id] [address]

# REST API
curl "https://lcd.bitbadges.io/bitbadges/bitbadgeschain/badges/get_dynamic_store_value/1/bb1..."
```

### Response Example

```json
{
    "value": {
        "storeId": "1",
        "address": "bb1...",
        "value": "100"
    }
}
```
