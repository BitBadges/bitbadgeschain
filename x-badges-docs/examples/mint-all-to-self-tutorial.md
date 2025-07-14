# Mint All Badges to Self - Tutorial

This tutorial walks through the process of creating a collection and minting all badges to yourself in a single transaction. This is useful for creating collections where you want to control the initial distribution.

## Overview

This is a two-step process that can be executed as a single multi-message transaction:

1. **Create Collection** with a mint approval that allows you to mint badges
2. **Execute Transfer** using that approval to mint badges to yourself

## Step 1: Create Mint Approval

First, create an approval that allows you to mint badges from the "Mint" address:

```typescript
// Step 1: Set up your mint approval
const mintApproval = {
    fromListId: 'Mint', // From the mint address
    toListId: 'All', // To any address
    initiatedByListId: myAddress, // Only you can initiate
    transferTimes: UintRangeArray.FullRanges(),
    badgeIds: UintRangeArray.FullRanges(), // All badge IDs
    ownershipTimes: UintRangeArray.FullRanges(),
    approvalId: 'mint-approval',
    version: 0n,
    approvalCriteria: {
        // No restrictions - you can mint unlimited amounts
        ...defaultNoRestrictionsApprovalCriteria,
        overridesFromOutgoingApprovals: true, // Required for mint address
    },
};

// Step 1: Create your collection with the mint approval
const collection = {
    ...BaseCollectionDetails,
    collectionApprovalTimeline: [
        {
            timelineTimes: FullTimeRanges,
            collectionApprovals: [mintApproval, ...otherApprovals],
        },
    ],
};

// Create the collection
```

## Step 2: Execute Mint Transfer

After creating the collection, use the mint approval to transfer badges to yourself:

```typescript
// Step 2: Mint badges to yourself using the approval
const transfers = [
    {
        from: 'Mint', // From mint address
        toAddresses: [myAddress], // To your address
        balances: [
            {
                badgeIds: [{ start: 1n, end: 100n }],
                ownershipTimes: UintRangeArray.FullRanges(),
                amount: 100n,
            },
        ],
        // ... other transfer details
    },
];
```

## Related Concepts

-   [Building Collection Approvals](./building-collection-approvals.md)
-   [Admin Override Approval](./approvals/admin-override-approval.md)
-   [Mint Escrow Address](../concepts/mint-escrow-address.md)
-   [MsgTransferBadges](../messages/msg-transfer-badges.md)
