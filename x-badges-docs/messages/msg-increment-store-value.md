# MsgIncrementStoreValue

Increments a numeric value for a specific address in a dynamic store.

## Proto Definition

```protobuf
message MsgIncrementStoreValue {
  string creator = 1; // Address incrementing the value (must be store creator)
  string storeId = 2; // ID of the dynamic store
  string address = 3; // Address to increment the value for
  string amount = 4; // Amount to increment by
}

message MsgIncrementStoreValueResponse {}
```

## Usage Example

```bash
# CLI command
bitbadgeschaind tx badges increment-store-value [store-id] [address] [amount] --from creator-key
```

### JSON Example

```json
{
    "creator": "bb1...",
    "storeId": "1",
    "address": "bb1...",
    "amount": "10"
}
```

## Related Messages

-   [MsgSetDynamicStoreValue](./msg-set-dynamic-store-value.md) - Set absolute values
-   [MsgDecrementStoreValue](./msg-decrement-store-value.md) - Decrease values
-   [MsgCreateDynamicStore](./msg-create-dynamic-store.md) - Create new stores
