# MsgDecrementStoreValue

Decrements a numeric value for a specific address in a dynamic store.

## Proto Definition

```protobuf
message MsgDecrementStoreValue {
  string creator = 1; // Address decrementing the value (must be store creator)
  string storeId = 2; // ID of the dynamic store
  string address = 3; // Address to decrement the value for
  string amount = 4; // Amount to decrement by
  bool setToZeroOnUnderflow = 5; // Whether to set to 0 if result would be negative
}

message MsgDecrementStoreValueResponse {}
```

## Usage Example

```bash
# CLI command
bitbadgeschaind tx badges decrement-store-value [store-id] [address] [amount] [set-to-zero-on-underflow] --from creator-key
```

### JSON Example

```json
{
    "creator": "bb1...",
    "storeId": "1",
    "address": "bb1...",
    "amount": "5",
    "setToZeroOnUnderflow": false
}
```

## Underflow Behavior

The `setToZeroOnUnderflow` field controls what happens when decrementing would result in a negative value:

-   **`true`**: The value is set to 0 instead of going negative
-   **`false`**: The operation fails with an error if it would result in a negative value

## Related Messages

-   [MsgSetDynamicStoreValue](./msg-set-dynamic-store-value.md) - Set absolute values
-   [MsgIncrementStoreValue](./msg-increment-store-value.md) - Increase values
-   [MsgCreateDynamicStore](./msg-create-dynamic-store.md) - Create new stores
