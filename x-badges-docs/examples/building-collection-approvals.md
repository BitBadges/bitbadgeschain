# Building Your Collection Approvals

The collection-level transferability is determined by the collection-level approvals. The important thing to consider here is that any approval that allows transfers from the "Mint" address will mint balances out of thin air.

## Approval Categories

It is typically recommended to split into two categories:

-   **Mint Approvals** (`fromListId: 'Mint'`)
-   **Post-Mint Approvals** (`fromListId: '!Mint'`)

## Important Notes

1. The reserved "All" list ID includes Mint. Do not use "All" for the fromListId for post-mint approvals.
2. To function, the "Mint" approval must forcefully override the user-level outgoing approval because it cannot be managed.

## Code Example

Mix and match the approvals as you see fit. See the examples in the Approvals folder for a bunch of examples.

-   [Transferable Approval](./approvals/transferable-approval.md)
-   [Burnable Approval](./approvals/burnable-approval.md)

```typescript
const mintApprovals = [
    // Mint approvals with fromListId: 'Mint'
];

const postMintApprovals = [
    // Post-mint approvals with fromListId: '!Mint'
    transferableApproval,
    burnableApproval,
];

const collectionApprovals = [...mintApprovals, ...postMintApprovals];

const collectionApprovalTimeline = [
    {
        timelineTimes: FullTimeRanges,
        collectionApprovals,
    },
];
```

## Related Concepts

-   [Transferability / Approvals](../concepts/transferability-approvals.md)
-   [Address Lists](../concepts/address-lists.md)
-   [Timeline System](../concepts/timeline-system.md)
-   [Approval Examples](./approvals/)
