# Building User-Level Approvals

User-level approvals allow individual users to control their badge transfers through incoming and outgoing approvals. These work similarly to [collection-level approvals](./building-collection-approvals.md) with key restrictions.

We refer you to the collection-level examples and just apply the same logic to the user-level types with these differences.

## Key Differences from Collection Approvals

-   **Fixed Address Lists**:
    -   Incoming approvals: `fromListId` is locked to the user's address
    -   Outgoing approvals: `toListId` is locked to the user's address
-   **No Override Functionality**: Cannot override other approval levels
-   **User-Controlled**: Only the user can update their own approvals

## Incoming Approvals

Control what badges the user can receive:

```typescript
const userIncomingApproval = {
    fromListId: 'user-address', // Locked to approver's address
    toListId: 'All', // Can specify recipients
    initiatedByListId: 'All',
    transferTimes: [{ start: '1', end: '18446744073709551615' }],
    badgeIds: [{ start: '1', end: '100' }],
    ownershipTimes: [{ start: '1', end: '18446744073709551615' }],
    approvalId: 'user-incoming-approval',

    // Use any approval criteria from collection examples
    approvalCriteria: {
        // See: transferable-approval.md, burnable-approval.md, etc.
        // OR use EmptyApprovalCriteria for no restrictions
        ...EmptyApprovalCriteria,
    },
};
```

## Outgoing Approvals

Control what badges the user can send:

```typescript
const userOutgoingApproval = {
    fromListId: 'All', // Can specify senders
    toListId: 'user-address', // Locked to approver's address
    initiatedByListId: 'All',
    transferTimes: [{ start: '1', end: '18446744073709551615' }],
    badgeIds: [{ start: '1', end: '100' }],
    ownershipTimes: [{ start: '1', end: '18446744073709551615' }],
    approvalId: 'user-outgoing-approval',

    // Use any approval criteria from collection examples
    approvalCriteria: {
        // See: transferable-approval.md, burnable-approval.md, etc.
        // OR use EmptyApprovalCriteria for no restrictions
        ...EmptyApprovalCriteria,
    },
};
```

## Implementation

Users update their approvals via `MsgUpdateUserApprovals`:

```typescript
const updateUserApprovals = {
    creator: 'bb1...', // Your address
    collectionId: '1',
    updateIncomingApprovals: true,
    incomingApprovals: [userIncomingApproval],
    updateOutgoingApprovals: true,
    outgoingApprovals: [userOutgoingApproval],
    // ...
};
```

## Reference

For approval criteria examples, see:

-   [Empty Approval Criteria](./empty-approval-criteria.md) - No restrictions template
-   [Transferable Approval](./approvals/transferable-approval.md) - Basic transfer restrictions
-   [Burnable Approval](./approvals/burnable-approval.md) - Burn functionality
-   [Building Collection Approvals](./building-collection-approvals.md) - Collection-level patterns

For concepts, see:

-   [Transferability / Approvals](../concepts/transferability-approvals.md)
-   [Approval Criteria](../concepts/approval-criteria/README.md)
