# Building User-Level Permissions

User-level permissions allow individual users to control their ability to update their own approvals. Note that these are almost always never needed unless in advanced situations. Typically, you just leave these soft-enabled (empty arrays) for all. These are only really needed in advanced situations where you want to lock down a user's ability to update their own approvals, such as escrow accounts.

The canUpdateOutgoingApprovals and canUpdateIncomingApprovals work similarly to [canUpdateCollectionApprovals](./building-collection-permissions.md) with key restrictions. - `fromListId` is locked to the user's address for outgoing approvals - `toListId` is locked to the user's address for incoming approvals

## User Permission Structure

```typescript
const userPermissions = {
    canUpdateOutgoingApprovals: [
        {
            // fromListId: 'user-address', // Locked to user's address
            toListId: 'All', // Can specify recipients
            initiatedByListId: 'All',
            transferTimes: FullTimeRanges,
            badgeIds: FullTimeRanges,
            ownershipTimes: FullTimeRanges,
            approvalId: 'All',
            permanentlyPermittedTimes: [],
            permanentlyForbiddenTimes: FullTimeRanges, // Lock forever
        },
    ],
    canUpdateIncomingApprovals: [
        {
            fromListId: 'All', // Can specify senders
            //  toListId: 'user-address', // Locked to user's address
            initiatedByListId: 'All',
            transferTimes: FullTimeRanges,
            badgeIds: FullTimeRanges,
            ownershipTimes: FullTimeRanges,
            approvalId: 'All',
            permanentlyPermittedTimes: [],
            permanentlyForbiddenTimes: FullTimeRanges, // Lock forever
        },
    ],
    canUpdateAutoApproveSelfInitiatedOutgoingTransfers: [
        {
            permanentlyPermittedTimes: [],
            permanentlyForbiddenTimes: FullTimeRanges,
        },
    ],
    canUpdateAutoApproveSelfInitiatedIncomingTransfers: [
        {
            permanentlyPermittedTimes: [],
            permanentlyForbiddenTimes: FullTimeRanges,
    canUpdateAutoApproveAllIncomingTransfers: [
        {
            permanentlyPermittedTimes: [],
            permanentlyForbiddenTimes: FullTimeRanges,
        },
    ],
};
```

## Implementation

Users update their permissions via `MsgUpdateUserApprovals`:

```typescript
const updateUserApprovals = {
    creator: 'bb1...', // User's address
    collectionId: '1',
    updateUserPermissions: true,
    userPermissions,
    // ... other approval updates
};
```

## Related Examples

For permission patterns, see:

-   [Freezing Mint Transferability](./permissions/freezing-mint-transferability.md) - Collection permission example
-   [Locking Specific Approval ID](./permissions/locking-specific-approval-id.md) - Approval ID targeting
-   [Locking Specific Token IDs](./permissions/locking-specific-badge-ids.md) - Token ID targeting
-   [Building Collection Permissions](./building-collection-permissions.md) - Collection-level patterns

For user approval configuration, see:

-   [Building User Approvals](./building-user-approvals.md) - User approval setup

## Related Concepts

-   [Permissions System](../concepts/permissions/README.md) - Permission mechanics
-   [Update Approval Permission](../concepts/permissions/update-approval-permission.md) - Approval-specific controls
-   [Default Balances](../concepts/default-balances.md) - User permission structure
