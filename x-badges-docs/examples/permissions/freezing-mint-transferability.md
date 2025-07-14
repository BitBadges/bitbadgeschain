# Freezing Mint Transferability

This example demonstrates how to permanently freeze minting capabilities by making mint-related collection approvals immutable.

## Overview

By setting `permanentlyForbiddenTimes` for mint approval updates, you can ensure that no new minting approvals can be added and existing ones cannot be modified.

## Permission Configuration

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
            fromListId: 'Mint',
            toListId: 'All',
            initiatedByListId: 'All',
            transferTimes: FullTimeRanges,
            badgeIds: FullTimeRanges,
            ownershipTimes: FullTimeRanges,
            approvalId: 'All',

            // What is status of this approval at any given time? (Unhandled = soft-enabled)
            permanentlyPermittedTimes: [],
            permanentlyForbiddenTimes: FullTimeRanges,
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
        // Include any initial mint approvals here
        // These will be the ONLY mint approvals ever possible
        {
            fromListId: 'Mint',
            toListId: 'creator-address',
            // ... initial mint approval configuration
        },
    ],
};
```

## Important Notes

### ⚠️ Irreversible Action

Once set to permanently forbidden, mint permissions cannot be restored. Carefully configure initial mint approvals before freezing. Ensure all mint approvals you will ever need are set.

## Related Examples

-   [Building Collection Permissions](../building-collection-permissions.md) - General permission patterns
-   [Building Collection Approvals](../building-collection-approvals.md) - Approval configuration

## Related Concepts

-   [Permissions System](../../concepts/permissions/README.md) - Permission mechanics
-   [Collection Permissions](../../concepts/permissions/permission-system.md) - Collection-level controls
-   [Action Permissions](../../concepts/permissions/action-permission.md) - Specific action controls
