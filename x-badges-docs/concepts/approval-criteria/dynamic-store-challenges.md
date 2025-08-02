# Dynamic Store Challenges

Require transfer initiators to pass checks against dynamic stores. Typically, these are used with smart contracts.

Dynamic stores are simply standalone (address -> number) stores where the number is the amount of uses an initiator has left. They are controlled by whoever creates them. These are powerful for creating dynamic approval criteria with smart contracts and other custom use cases.

## How It Works

Dynamic store challenges check if the transfer initiator has a value greater than 0 in specified dynamic stores. The system:

1. **Checks Initiator**: Looks up the initiator's address in the specified dynamic store
2. **Evaluates Number**: Returns the numeric value for the initiator
3. **Requires All > 0**: All challenges must return a value greater than 0 for approval
4. **Fails if Any ≤ 0**: If any challenge returns 0 or less, transfer is denied

## Interface

```typescript
interface DynamicStoreChallenge<T extends NumberType> {
    storeId: string; // Dynamic store ID to check
}
```

## Usage in Approval Criteria

```json
{
    "dynamicStoreChallenges": [
        { "storeId": "1" }, // Member points store (must have > 0 points)
        { "storeId": "2" } // Subscription level store (must have > 0 level)
    ]
}
```

## Managing Dynamic Stores

### Creating Stores

Use [MsgCreateDynamicStore](../../messages/msg-create-dynamic-store.md) to create new dynamic stores with default numeric values.

### Setting Values

Use [MsgSetDynamicStoreValue](../../messages/msg-set-dynamic-store-value.md) to set numeric values for specific addresses.

### Incrementing Values

Use [MsgIncrementStoreValue](../../messages/msg-increment-store-value.md) to increase values for specific addresses.

### Decrementing Values

Use [MsgDecrementStoreValue](../../messages/msg-decrement-store-value.md) to decrease values for specific addresses.

### Querying Values

Use [GetDynamicStoreValue](../../queries/get-dynamic-store-value.md) to check current values for addresses.

## Alternatives

For fully off-chain solutions, consider:

-   [Merkle Challenges](merkle-challenges.md) to save gas costs
-   [ETH Signature Challenges](eth-signature-challenges.md) for direct authorization
