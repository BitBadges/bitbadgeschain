# Balance System

The Balance system in BitBadges is designed to represent ownership of tokens across different IDs and time ranges. Ownership times are a new concept to BitBadges allowing you to set that someone owns a token during a specific time but not other times.

## Balance Interface

```typescript
export interface Balance<T extends NumberType> {
    amount: T;
    badgeIds: UintRange<T>[];
    ownershipTimes: UintRange<T>[];
}
```

-   `amount`: The quantity of tokens owned
-   `badgeIds`: An array of ID ranges representing the tokens owned
-   `ownershipTimes`: An array of time ranges during which the tokens are owned

## Interpreting Balances

When interpreting balances, it's crucial to understand that multiple ranges of token IDs and ownership times within a single Balance structure represent all possible combinations.

### Interpretation Algorithm

```javascript
for (balance of balances) {
    for (badgeIdRange of balance.badgeIds) {
        for (ownershipTimeRange of balance.ownershipTimes) {
            // User owns x(balance.amount) of (badgeIdRange) for the times (ownershipTimeRange)
        }
    }
}
```

### Example

Consider the following balance:

```json
{
    "amount": 1,
    "badgeIds": [
        { "start": 1, "end": 10 },
        { "start": 20, "end": 30 }
    ],
    "ownershipTimes": [
        { "start": 20, "end": 50 },
        { "start": 100, "end": 200 }
    ]
}
```

This balance expands to:

1. 1x of IDs 1-10 from times 20-50
2. 1x of IDs 1-10 from times 100-200
3. 1x of IDs 20-30 from times 20-50
4. 1x of IDs 20-30 from times 100-200

## Balance Subtraction

When subtracting balances, you may need to represent the result as multiple Balance objects. For example, if we subtract the first set of balances from the example above (1x of IDs 1-10 from times 20-50), the result would be:

```json
[
    {
        "amount": 1,
        "badgeIds": [
            { "start": 1, "end": 10 },
            { "start": 20, "end": 30 }
        ],
        "ownershipTimes": [{ "start": 100, "end": 200 }]
    },
    {
        "amount": 1,
        "badgeIds": [{ "start": 20, "end": 30 }],
        "ownershipTimes": [{ "start": 20, "end": 50 }]
    }
]
```

## Handling Duplicates

When duplicate token IDs are specified in balances, they are combined and their amounts are added. For example:

```json
{
    "amount": 1,
    "badgeIds": [
        { "start": 1, "end": 10 },
        { "start": 1, "end": 10 }
    ],
    "ownershipTimes": [{ "start": 100, "end": 200 }]
}
```

This is equivalent to and will be treated as:

```json
{
    "amount": 2,
    "badgeIds": [{ "start": 1, "end": 10 }],
    "ownershipTimes": [{ "start": 100, "end": 200 }]
}
```

## Best Practices

1. **Efficient Representation**: Try to represent balances in the most compact form possible by combining overlapping ranges
2. **Careful Subtraction**: When subtracting balances, ensure that you correctly split the remaining balances to accurately represent the result
3. **Avoid Duplicates**: While the system handles duplicates by combining them, it's more efficient to represent balances without duplicates in the first place
4. **Time-Aware Operations**: Always consider the time dimension when performing operations on balances, as ownership can vary over time
5. **Range Calculations**: Familiarize yourself with range operations, as they are crucial for correctly manipulating and interpreting balances
