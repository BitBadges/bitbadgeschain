# Building Your Collection Permissions

Collection permissions are executable by the manager. They are used to control who can perform various management actions on your collection and when those actions are allowed.

```typescript
const manager = collection.getCurrentManager();
```

## Setting Your Permissions

You have a few options for setting your permissions.

1. No Manager

If you simply don't want a manager, you can set the manager to an empty string. Then, the permission values never matter.

```typescript
const managerTimeline = [
    {
        manager: '',
        timelineTimes: FullTimeRanges,
    },
];
```

2. Complete Control - Soft Enabled

Each permission is enabled by default, unless you permanently disabled it. Thus, an empty array means that the permission is enabled for all times. However, it is soft enabled, meaning that the manager can disable it at any time. This configuration offers full control with ability to disable in the future.

```typescript
const collectionPermissions = {
    canDeleteCollection: [],
    canArchiveCollection: [],
    canUpdateOffChainBalancesMetadata: [],
    canUpdateStandards: [],
    canUpdateCustomData: [],
    canUpdateManager: [],
    canUpdateCollectionMetadata: [],
    canUpdateBadgeMetadata: [],
    canUpdateCollectionApprovals: [],
    canUpdateValidBadgeIds: [],
};
```

3. Custom Permissions

Oftentimes, you want a little more control over your permissions though.

Each permission follows the same pattern:

1. For the times `permanentlyPermittedTimes`, the permission is always permitted for the given values.
2. For the times `permanentlyForbiddenTimes`, the permission is always forbidden for the given values.
3. If the item is not explicity in either, then the permission is enabled for the given values, but the status can change.

```typescript
const CanArchiveCollection = {
    permanentlyPermittedTimes: [],
    permanentlyForbiddenTimes: FullTimeRanges,
    timelineTimes: FullTimeRanges,
};
```

Each permission type follows the same pattern of two categories:

```typescript
// Part 1. Enabled vs Disabled Times For The Execution Of The Permission
const permanentlyPermittedTimes = [];
const permanentlyForbiddenTimes = FullTimeRanges;

// Part 2. For what values (if any) does this apply? This is dependent on the permission type.
const {
    timelineTimes,
    badgeIds,
    fromListId,
    toListId,
    initiatedByListId,
    transferTimes,
    ownershipTimes,
    approvalId,
} = permission;
```

## Main Permissions To Consider

1. Should the number of token IDs in the collection be expandable? frozen upon genesis? -> Handle with `canUpdateValidBadgeIds`
2. What about the transferability? -> Handle with `canUpdateCollectionApprovals`
    - Should the transferability be frozen upon genesis?
    - Should we disallow updating transferability for only some token IDs? some approvals? Mint? Post-Mint?
    - This could be critical for enforcing total circulating supply. For example, if you can create more approvals from the Mint address, then you can theoretically mint however many tokens you want.

## Examples

We refer you to the [examples](../examples/permissions) or relevant concepts for more detailed examples.

## Related Concepts

-   [Permission System](../concepts/permissions/permission-system.md)
-   [Manager](../concepts/manager.md)
-   [Timeline System](../concepts/timeline-system.md)
-   [Timed Update Permission](../concepts/permissions/timed-update-permission.md)
