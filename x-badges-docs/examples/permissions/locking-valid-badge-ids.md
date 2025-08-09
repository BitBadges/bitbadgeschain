# Locking Valid Token IDs

This example demonstrates how to control updates to the `validTokenIds` field, either locking it permanently or allowing controlled expansion. The `validTokenIds` field is used to control which token IDs are considered valid for the collection.

## Overview

The `canUpdateValidTokenIds` permission controls whether the valid token ID ranges can be modified.

## Lock Valid Token IDs Forever

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
    canUpdateValidTokenIds: [
        {
            // Which token IDs does this permission apply to?
            tokenIds: FullTimeRanges, // All token IDs

            // What is status of this permission at any given time?
            permanentlyPermittedTimes: [],
            permanentlyForbiddenTimes: FullTimeRanges, // Never allowed to update
        },
    ],
    canUpdateTokenMetadata: [],
    canUpdateCollectionApprovals: [],
};
```

## Lock Token IDs 1-100, Allow Future Expansion

```typescript
const collectionPermissions = {
    canDeleteCollection: [],
    canArchiveCollection: [],
    canUpdateOffChainBalancesMetadata: [],
    canUpdateStandards: [],
    canUpdateCustomData: [],
    canUpdateManager: [],
    canUpdateCollectionMetadata: [],
    canUpdateValidTokenIds: [
        {
            // Which token IDs does this permission apply to?
            tokenIds: [
                {
                    start: '1',
                    end: '100', // Only applies to tokens 1-100
                },
            ],

            // What is status of this permission at any given time?
            permanentlyPermittedTimes: [],
            permanentlyForbiddenTimes: FullTimeRanges, // Token IDs 1-100 locked forever
        },
        // Token IDs 101+ remain soft-enabled (can be updated by manager)
    ],
    canUpdateTokenMetadata: [],
    canUpdateCollectionApprovals: [],
};
```

## Implementation

```typescript
const createCollection = {
    // ... other collection fields
    collectionPermissions,
    validTokenIds: [
        {
            start: '1',
            end: '100', // Initial valid range
        },
    ],
};
```

## Important Notes

### ⚠️ Token ID Targeting

-   Permissions only apply to the specified token ID ranges
-   Unspecified ranges remain soft-enabled for manager updates
-   Cannot reduce valid token IDs once locked (only expansion possible for unlocked ranges)

## Related Examples

-   [Locking Specific Token IDs](./locking-specific-badge-ids.md) - Lock approval updates for token ranges
-   [Freezing Mint Transferability](./freezing-mint-transferability.md) - Lock mint approvals

## Related Concepts

-   [Valid Token IDs](../../concepts/valid-badge-ids.md) - Token ID range concept
-   [Token IDs Action Permission](../../concepts/permissions/balances-action-permission.md) - Token-specific permission controls
