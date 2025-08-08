# Burnable Approval

This example demonstrates how to create a burnable approval that allows tokens to be sent to the burn address (0x0000000000000000000000000000000000000000), effectively removing them from circulation.

## Overview

A burnable approval enables tokens to be permanently destroyed by sending them to the zero address.

## Code Example

```typescript
const burnableApproval = new CollectionApproval({
    fromListId: '!Mint', // Excludes the Mint address
    toListId: convertToBitBadgesAddress(
        '0x0000000000000000000000000000000000000000' //bb1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqs7gvmv
    ),
    initiatedByListId: 'All',
    transferTimes: UintRangeArray.FullRanges(),
    ownershipTimes: UintRangeArray.FullRanges(),
    badgeIds: UintRangeArray.FullRanges(),
    approvalId: 'burnable-approval',
    version: 0n,
    approvalCriteria: undefined, // No additional restrictions
});
```

## Related Concepts

-   [Transferability / Approvals](../../concepts/transferability-approvals.md)
-   [Address Lists](../../concepts/address-lists.md)
-   [Timeline System](../../concepts/timeline-system.md)
