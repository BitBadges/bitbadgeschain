# Locking Specific Approval ID

This example demonstrates how to permanently lock a specific approval ID while keeping other approvals updatable.

## Overview

By targeting a specific `approvalId`, you can freeze that approval permanently while allowing updates to other approvals. The `!` operator can be used to target all approvals EXCEPT a specific ID.

## Lock Specific Approval ID

```typescript
const FullTimeRanges = [
    {
        start: '1',
        end: '18446744073709551615',
    },
];

const collectionPermissions = {
    canDeleteCollection: [],
    canArchiveCollection: [],
    canUpdateOffChainBalancesMetadata: [],
    canUpdateStandards: [],
    canUpdateCustomData: [],
    canUpdateManager: [],
    canUpdateCollectionMetadata: [],
    canUpdateValidBadgeIds: [],
    canUpdateBadgeMetadata: [],
    canUpdateCollectionApprovals: [
        {
            // Which approvals does this permission apply to? Approvals must match ALL criteria.
            fromListId: 'All',
            toListId: 'All',
            initiatedByListId: 'All',
            transferTimes: FullTimeRanges,
            badgeIds: FullTimeRanges,
            ownershipTimes: FullTimeRanges,
            approvalId: 'abc123', // Only targets this specific approval ID

            // What is status of this approval at any given time? (Unhandled = soft-enabled)
            permanentlyPermittedTimes: [],
            permanentlyForbiddenTimes: FullTimeRanges, // Permanently locked
        },
    ],
};
```

## Lock All EXCEPT Specific Approval ID

```typescript
const collectionPermissions = {
    canDeleteCollection: [],
    canArchiveCollection: [],
    canUpdateOffChainBalancesMetadata: [],
    canUpdateStandards: [],
    canUpdateCustomData: [],
    canUpdateManager: [],
    canUpdateCollectionMetadata: [],
    canUpdateValidBadgeIds: [],
    canUpdateBadgeMetadata: [],
    canUpdateCollectionApprovals: [
        {
            // Which approvals does this permission apply to? Approvals must match ALL criteria.
            fromListId: 'All',
            toListId: 'All',
            initiatedByListId: 'All',
            transferTimes: FullTimeRanges,
            badgeIds: FullTimeRanges,
            ownershipTimes: FullTimeRanges,
            approvalId: '!abc123', // All approvals EXCEPT abc123

            // What is status of this approval at any given time? (Unhandled = soft-enabled)
            permanentlyPermittedTimes: [],
            permanentlyForbiddenTimes: FullTimeRanges, // All others permanently locked
        },
    ],
};
```

## Implementation

```typescript
const createCollection = {
    // ... other collection fields
    collectionPermissions,
    collectionApprovals: [
        {
            approvalId: 'abc123',
            // ... this approval will be locked/unlocked based on configuration
        },
        {
            approvalId: 'other-approval',
            // ... this approval's updateability depends on configuration
        },
    ],
};
```

## Related Examples

-   [Freezing Mint Transferability](./freezing-mint-transferability.md) - Lock all mint approvals
-   [Building Collection Permissions](../building-collection-permissions.md) - General permission patterns

## Related Concepts

-   [Permissions System](../../concepts/permissions/README.md) - Permission mechanics
-   [Update Approval Permission](../../concepts/permissions/update-approval-permission.md) - Approval-specific controls
