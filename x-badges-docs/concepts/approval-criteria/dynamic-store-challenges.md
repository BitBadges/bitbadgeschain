# Dynamic Store Challenges

Require transfer initiators to pass boolean checks against dynamic stores. Typically, these are used with smart contracts.

## How It Works

Dynamic store challenges check if the transfer initiator has a `true` value in specified dynamic stores. The system:

1. **Checks Initiator**: Looks up the initiator's address in the specified dynamic store
2. **Evaluates Boolean**: Returns `true` or `false` for the initiator
3. **Requires All True**: All challenges must return `true` for approval
4. **Fails if Any False**: If any challenge returns `false`, transfer is denied

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
        { "storeId": "1" }, // Member status store
        { "storeId": "2" } // Subscription status store
    ]
}
```

## Managing Dynamic Stores

### Creating Stores

Use [MsgCreateDynamicStore](../../messages/msg-create-dynamic-store.md) to create new dynamic stores with default boolean values.

### Setting Values

Use [MsgSetDynamicStoreValue](../../messages/msg-set-dynamic-store-value.md) to set boolean values for specific addresses.

### Querying Values

Use [GetDynamicStoreValue](../../queries/get-dynamic-store-value.md) to check current values for addresses.

## Alternative

For fully off-chain solutions, consider [Merkle Challenges](merkle-challenges.md) to save gas costs.
