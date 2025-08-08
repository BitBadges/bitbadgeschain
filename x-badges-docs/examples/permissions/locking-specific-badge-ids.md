# Locking Specific Token IDs

This example demonstrates how to permanently lock approvals for specific token IDs while keeping other approvals updatable.

## Overview

By targeting specific `badgeIds`, you can freeze approvals for those tokens permanently while allowing updates to approvals for other token IDs.

## Lock Token IDs 1-100

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
            badgeIds: [
                {
                    start: '1',
                    end: '100', // Only targets tokens 1-100
                },
            ],
            ownershipTimes: FullTimeRanges,
            approvalId: 'All',

            // What is status of this approval at any given time? (Unhandled = soft-enabled)
            permanentlyPermittedTimes: [],
            permanentlyForbiddenTimes: FullTimeRanges, // Permanently locked
        },
    ],
};
```

## Lock All Tokens EXCEPT 1-100

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
            badgeIds: [
                {
                    start: '101',
                    end: '18446744073709551615', // All tokens except 1-100
                },
            ],
            ownershipTimes: FullTimeRanges,
            approvalId: 'All',

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
            badgeIds: [{ start: '1', end: '50' }],
            // ... this approval will be locked if it overlaps with permission criteria
        },
        {
            badgeIds: [{ start: '150', end: '200' }],
            // ... this approval's updateability depends on configuration
        },
    ],
};
```

## Use Cases

- **Lock Founder Tokens**: Prevent modification of special token 1-100 transfer rules
- **Preserve Rare Items**: Keep limited edition tokens (1-100) immutable
- **Tier-Based Control**: Lock specific tiers while allowing others to evolve

## Important Notes

### ⚠️ ID Range Targeting

The permission only applies to approvals that overlap with the specified token ID ranges. Approvals targeting token IDs outside the range remain updatable.

## Related Examples

- [Locking Specific Approval ID](./locking-specific-approval-id.md) - Lock by approval ID
- [Freezing Mint Transferability](./freezing-mint-transferability.md) - Lock all mint approvals

## Related Concepts

- [Permissions System](../../concepts/permissions/README.md) - Permission mechanics
- [Timed Update With Token IDs Permission](../../concepts/permissions/timed-update-with-badge-ids-permission.md) - Token-specific controls