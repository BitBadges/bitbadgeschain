# Transferable Approval

This example demonstrates how to create a basic transferable approval that allows tokens to be freely transferred between any users after minting.

## Overview

A transferable approval enables tokens to be moved between addresses without restrictions.

## Code Example

```typescript
const transferableApproval = new CollectionApproval({
    fromListId: '!Mint', // Excludes the Mint address
    toListId: 'All',
    initiatedByListId: 'All',
    transferTimes: UintRangeArray.FullRanges(),
    ownershipTimes: UintRangeArray.FullRanges(),
    tokenIds: UintRangeArray.FullRanges(),
    approvalId: 'transferable-approval',
    version: 0n,
    approvalCriteria: undefined, // No additional restrictions
});
```

## Related Concepts

-   [Transferability / Approvals](../../concepts/transferability-approvals.md)
-   [Address Lists](../../concepts/address-lists.md)
-   [Timeline System](../../concepts/timeline-system.md)
