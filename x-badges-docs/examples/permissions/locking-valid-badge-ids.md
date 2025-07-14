# Locking Valid Badge IDs

This example demonstrates how to control updates to the `validBadgeIds` field, either locking it permanently or allowing controlled expansion. The `validBadgeIds` field is used to control which badge IDs are considered valid for the collection.

## Overview

The `canUpdateValidBadgeIds` permission controls whether the valid badge ID ranges can be modified.

## Lock Valid Badge IDs Forever

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
    canUpdateValidBadgeIds: [
        {
            // Which badge IDs does this permission apply to?
            badgeIds: FullTimeRanges, // All badge IDs

            // What is status of this permission at any given time?
            permanentlyPermittedTimes: [],
            permanentlyForbiddenTimes: FullTimeRanges, // Never allowed to update
        },
    ],
    canUpdateBadgeMetadata: [],
    canUpdateCollectionApprovals: [],
};
```

## Lock Badge IDs 1-100, Allow Future Expansion

```typescript
const collectionPermissions = {
    canDeleteCollection: [],
    canArchiveCollection: [],
    canUpdateOffChainBalancesMetadata: [],
    canUpdateStandards: [],
    canUpdateCustomData: [],
    canUpdateManager: [],
    canUpdateCollectionMetadata: [],
    canUpdateValidBadgeIds: [
        {
            // Which badge IDs does this permission apply to?
            badgeIds: [
                {
                    start: '1',
                    end: '100', // Only applies to badges 1-100
                },
            ],

            // What is status of this permission at any given time?
            permanentlyPermittedTimes: [],
            permanentlyForbiddenTimes: FullTimeRanges, // Badge IDs 1-100 locked forever
        },
        // Badge IDs 101+ remain soft-enabled (can be updated by manager)
    ],
    canUpdateBadgeMetadata: [],
    canUpdateCollectionApprovals: [],
};
```

## Implementation

```typescript
const createCollection = {
    // ... other collection fields
    collectionPermissions,
    validBadgeIds: [
        {
            start: '1',
            end: '100', // Initial valid range
        },
    ],
};
```

## Important Notes

### ⚠️ Badge ID Targeting

-   Permissions only apply to the specified badge ID ranges
-   Unspecified ranges remain soft-enabled for manager updates
-   Cannot reduce valid badge IDs once locked (only expansion possible for unlocked ranges)

## Related Examples

-   [Locking Specific Badge IDs](./locking-specific-badge-ids.md) - Lock approval updates for badge ranges
-   [Freezing Mint Transferability](./freezing-mint-transferability.md) - Lock mint approvals

## Related Concepts

-   [Valid Badge IDs](../../concepts/valid-badge-ids.md) - Badge ID range concept
-   [Badge IDs Action Permission](../../concepts/permissions/balances-action-permission.md) - Badge-specific permission controls
