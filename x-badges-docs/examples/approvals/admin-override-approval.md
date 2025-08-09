# Admin Override Approval

This example demonstrates how to create an approval that allows a specific address to forcefully transfer tokens, overriding all user-level approvals. This provides complete administrative control for emergency situations or management purposes.

## Overview

An admin override approval grants a specific address the power to:

-   Transfer tokens from any address to any address
-   Override user-level incoming and outgoing approvals
-   Bypass normal approval restrictions
-   Maintain complete administrative control

⚠️ **Warning**: This approval type grants significant power and should be used carefully with trusted addresses only.

## Code Example

```typescript
const approveSelfForcefully = (address: string) => {
    const id = 'complete-admin-control';

    return {
        fromListId: 'Mint',
        toListId: 'All',
        initiatedByListId: address,
        transferTimes: UintRangeArray.FullRanges(),
        badgeIds: UintRangeArray.FullRanges(),
        ownershipTimes: UintRangeArray.FullRanges(),
        approvalId: id,
        version: 0n,
        approvalCriteria: {
            ...EmptyApprovalCriteria,
            overridesFromOutgoingApprovals: true,
            overridesToIncomingApprovals: true,
        },
    };
};
```

## Related Concepts

-   [Transferability / Approvals](../../concepts/transferability-approvals.md)
-   [Address Lists](../../concepts/address-lists.md)
-   [Approval Criteria](../../concepts/approval-criteria/)
