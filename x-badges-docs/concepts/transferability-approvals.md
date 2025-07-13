# Transferability / Approvals

First, read [Transferability ](broken-reference/)for an overview of approved transfers.

Note: The [Approved Transfers](transferability-approvals.md) and [Permissions ](broken-reference/)are the most powerful features of the interface, but they can also be the most confusing. Please ask for help if needed.

Collections with "Off-Chain" balances and address lists do not utilize on-chain transferability, so this page is not applicable to them.

## Approvals Overview

### Approval Levels - Collection, Incoming, Outgoing

Approved transfers encompass three hierarchical levels: collection, incoming, and outgoing, as previously elaborated. The interfaces for these three levels share common elements, with slight variations in functionality:

-   Incoming approvals exclude the "to" fields as they are automatically populated with the recipient's address.
-   Outgoing approvals omit the "from" fields, as they are automatically filled with the sender's address.
-   The collection level holds the capacity to override user-level (incoming / outgoing) approvals, but not vice versa.

**For a transfer to be approved, it has to satisfy the collection-level approvals, and if not overriden forcefully by the collection-level approvals, the user incoming / outgoing also have to be satisfied.**

<figure><img src="../../../.gitbook/assets/image (33).png" alt=""><figcaption></figcaption></figure>

### Approvals != Escrows

When it comes to approved transfers, it's important to note that they are just approvals and may not necessarily correspond directly to underlying balances. Approvers must ensure that sufficient balances are available to uphold the approvals' integrity, accounting for potential revokes or freezes. This responsibility extends to the collection-level transferability as well.

On a similar note, if approvals become no longer valid (such as approving a badge but it was revoked via a different approval), the former approval doesn't automatically get cancelled. It is the approver's responsibility to handle them accordingly.

### Mint - Unlimited Balances

For on-chain balances, the Mint address has unlimited balances. So, it is extra critical to set these approvals correctly.

The Mint address does not have incoming / outgoing approvals that can be controlled. For all Mint approvals, they must forcefully override the user-level outgoing approval because it cannot be managed.

### Representation

In the collection interface, they are represented as the following:

```json
{
    ...
    "collectionApprovals": [{ ... }, ...]
    ...
}
```

```typescript
export interface CollectionApproval<T extends NumberType> {
    toListId: string;
    fromListId: string;
    initiatedByListId: string;
    transferTimes: UintRange<T>[];
    badgeIds: UintRange<T>[];
    ownershipTimes: UintRange<T>[];
    approvalId: string;

    uri?: string;
    customData?: string;

    approvalCriteria?: ApprovalCriteria<T>;

    version: T;
}
```

Note; User incoming / outgoing approvals follow the same interface except **toList** is auto-populated with the user's address for incoming approvals and similarly the **fromList** for outgoingApprovals.

**Approved vs Unapproved**

Approvals are simply a set of criteria, so it is entirely possible the same transfer could satisfy multiple approvals on the same level. We handle approvals per level in the following manner:

1. If the transfer is unhandled (doesn't match to any approval), it is DISAPPROVED by default.
2. We allow users to specify **prioritizedApprovals** and **onlyCheckPrioritizedApprovals** (in [MsgTransferBadges](../../bitbadges-blockchain/cosmos-sdk-msgs/x-badges/msgtransferbadges.md)) when transferring, so they can only use up their desired approvals.
    1. Any approval with side effects (coin transfers, criteria, trackers, etc) MUST be specified explicity in prioritized approvals.
    2. Any approval without side effects (e.g. generic unlimited transferability collection approval) does not have to be specified.
    3. When scanning approvals, we check prioritized approvals first. Then, we attempt to scan any approvals without side effects.

We strongly recommend designing approvals in a way where no transfer can map to multiple. This improves the simplicity and readability of your collection.

**Who? When? What? - Main Fields**

To represent transfers, six main fields are used: **`toList`**, **`fromList`**, **`initiatedByList`**, **`transferTimes`**, **`badgeIds`**, and **`ownershipTimes`**. These fields collectively define the transfer details, such as the addresses involved, timing, and badge details. This representation leverages range logic, breaking down into individual tuples for enhanced comprehension.

-   **toList, fromList, initiatedByList**: [AddressLists](../../core-concepts/address-lists-lists.md) specifying which addresses can send, receive, and initiate the transfer. If we use **toListId, fromListId, initiatedByListId**, these refer to the lists IDs of the respective lists. IDs can either be reserved IDs (see [AddressLists](../../core-concepts/address-lists-lists.md)) or IDs of lists created on-chain through [MsgCreateAddressLists](../../bitbadges-blockchain/cosmos-sdk-msgs/). Note that on-chain approvals cannot access off-chain lists.
-   **transferTimes**: When can the transfer takes place? A [UintRange](../../core-concepts/uint-ranges.md)\[] of times (UNIX milliseconds).
-   **badgeIds**: What badge IDs can be transferred? A [UintRange](../../core-concepts/uint-ranges.md)\[] of badge IDs.
-   **ownershipTimes**: What ownership times for the badges are being transferred? (UNIX milliseconds)

For example, we might have something like the following:

-   ```json
    "fromListId": "Mint", //reserved list ID for the "Mint" addres
    "toListId": "All", //reserved list ID for all addresses (excluding "Mint")
    "initiatedByListId": "All",
    "transferTimes": [
      {
        "start": "1691931600000",
        "end": "1723554000000"
      }
    ],
    "ownershipTimes": [
      {
        "start": "1",
        "end": "18446744073709551615" //Max possible value
      }
    ],
    "badgeIds": [
      {
        "start": "1",
        "end": "100"
      }
    ],
    ```

Let's break down the definition above.

-   The "Mint" list ID is the reserved list corresponding to the Mint address. This approval only allows transfers from the "Mint" address. The transfer can be initiated by any user (because the AddressList "All" includes all addresses) and can have any address as the recipient.
-   The ownership rights for any time of badge IDs 1-100 can be transferred from UNIX time 1691931600000 (Aug 13, 2023) to time 1723554000000 (Aug 13, 2024).

Note the approval only applies to the details defined and must match ALL details. For example, badge ID #101 is not defined by this approval even if all other criteria matches.

**Transferring From Mint Address**

As mentioned before, we check the collection level approvals first, and if not overriden, we check the user-level incoming/outgoing approvals.

The Mint address is a special case. It technically has its own approvals, but since it is not a real address, the user approvals are always empty and never usable. Thus, it is important that when you attempt transfers from the Mint address, you **override the outgoing approvals** of the Mint address (see [Overrides](approval-criteria/#overrides) on the next page for how).

It is also recommended that when dealing with approvals from the "Mint" address, the approval's **fromList** is only the "Mint" address and no other address. This helps readability and simplicity and avoiding unintentionally approving users to mint, which could be very bad. See Example 2 below.

Again, remember the Mint address has unlimited balances.

#### Approval IDs

All approvals must have a unique **approvalId** for identification per level. This is simply used for identification.

```json
{
    ...
    "approvalId": "abc123",
}
```

**Metadata**

We provide an optional **uri** and **customData** to allow you to add a link to something about your approval. See [Compatibility](../../bitbadges-api/concepts/designing-for-compatibility.md) for the expected format for the BitBadges API / Indexer.

This can typically be used for providing names, descriptions about your approvals. Or, we also use it to host N - 1 layers of a Merkle tree for a Merkle challenge of codes (N - 1 to be able to construct the path but not give away the value of leaves which are to be secret). Or, for whitelist trees where no leaves are secret, we can host the full tree. Learn more in the approval criteria merkle challenges section.

#### Approval Criteria

The **`approvalCriteria`** section corresponds to additional restrictions or challenges necessary to be satisfied for approval. It defines aspects like the quantity approved, maximum transfers, and more. There is a lot here, so we have dedicated a page to just explaining the [approval details here](approval-criteria/).

For the rest of this page, you can simply think of it as the challenges or restrictions that need to be obeyed to be approved.

**Breaking Down Range Logic**

Even though our interface uses range logic (UintRanges, AddressLists), you can think of it as we break everything down into single-value tuples (e.g., `(bob, alice, bob, 1, 1, 1000)`) and check each singular value tuple separately. This simplifies the matching process and enhances clarity.

#### Matching Transfers to Approvals

The process of matching transfers to approvals involves several steps. This is done on a per-level basis.

1. We start with the collection-level approvals.
2. Expand all approval tuples with range logic (AllWithMint, ...., \[IDs 1-100]) to singular tuple values (e.g. (bob, alice, bob, badge ID #1, ....)
3. Expand the current transfer tuple to singular tuple values.
4. Find all matches (for approvals, linear first match by default but can be customized with **prioritizedApprovals** and **onlyCheckPrioritizedApprovals**).
    1. If anything is unhandled on any approval level (accounting for overrides), the overall transfer is disapproved.
    2. In the case of overflowing approvals (e.g. we are transferring x10 but have two approvals for x3 and x12), we deduct as much as possible from each one as we iterate. So using the previous example, we would end up with x3/3 of the first approval used and x7/12 of the second used.
    3. We check the **`approvalCriteria`** for each match and ensure everything is satisfied. If not, it is not a match.
5. For any amounts / balances that were approved but do not override incoming / outgoing approvals respectively, we go back to step 2 and check the recipient's incoming approvals and the senders' outgoing approvals for those balances.

### Defaults and Auto Approvals

**Auto Approvals**

If **autoApproveSelfInitiatedOutgoingTransfers** is set to true, we automatically apply an unlimited approval (with no amount restrictions) to the user's outgoing approvals when the sender is the same as the initiator.

If **autoApproveSelfInitiatedIncomingTransfers** is set to true, we automatically apply an unlimited approval (with no amount restrictions) to the user's incoming approvals when the recipient is the same as the initiator.

if **autoApproveAllIncomingTransfers** is set to true, we automatically apply an unlimited approval (with no amount restrictions) to the user's incoming approvals when the recipient is the same as the initiator.

In 99% of cases, the auto approvals should be true because the expected functionality is that if the user is initiating the transaction, they also approve it. However, this can be turned off and leveraged for specific use cases such as using an account for an escrow.

**Defaults**

We allow the collection to define default values for each user, and when the user first interacts with the collection, they will start with these values. The defaults include **balances**, **outgoingApprovals**, **incomingApprovals**, **autoApproveSelfInitiatedOutgoingTransfers,** **autoApproveAllIncomingTransfers**, and **autoApproveSelfInitiatedIncomingTransfers.**

For default balances, we refer you to the balance types and creating badges sections (i.e. these are the starting balances).

Note: Default approvals can NOT contain custom criteria checks. In other words, default approvals can not have side effects. They can only be a simple approval without any custom restrictions.

**Default Outgoing Approvals**

For outgoing approvals, the expected functionality is that everything is disapproved by default unless self initiated. Thus, the following is the typical default values.

```
"defaultOutgoingApprovals": [],
"autoApproveSelfInitiatedOutgoingTransfers": true
```

**Default Incoming Approvals - Forceful Transfers vs Opt-In Only**

However, with incoming approvals, the expected functionality is slightly different. There are a couple options. Do you want users to be able to transfer "forcefully" to an address without prior approval by default? Or, do you want users to have to self-initiate / opt-in first to receive badges?

In order to allow forceful transfers to an address without prior approval, the **incomingApprovals** must be set to something like below. Otherwise, if empty or \[], then all transfers must be initiated by or manually approved by the recipient by default (opt-in only).

```json
//"forceful" is allowed
"defaultIncomingApprovals": [
    {
      "fromListId": "AllWithMint",
      "initiatedByListId": "AllWithMint",
      "transferTimes": [
        {
          "start": "1",
          "end": "18446744073709551615"
        }
      ],
      "badgeIds": [
        {
          "start": "1",
          "end": "18446744073709551615"
        }
      ],
      "ownershipTimes": [
        {
          "start": "1",
          "end": "18446744073709551615"
        }
      ],
      "approvalId": "forceful-transfers-allowed"
    }
  ]
```

### **Approval Value vs Permission**

While the value may seem similar to the approval update permissions, the permission corresponds to the **updatability** of the approvals (i.e. **canUpdateCollectionApprovals**). The approvals themselves correspond to if a transfer is currently approved or not.

**Example**

Current Value - IDs 1-10 are approved

Permission - IDs 1-5 cannot be updated, 6-10 can

### **Example 1 - Putting It Together**

Bob is transferring x10 of **ownershipTimes** 1-Max of **badgeIds** 1-2 at time T to Alice.

#### Transfer Request:

Transfer: `(Bob, Alice, Bob, T, [1-2], [1-Max])`

#### Collection Approved Transfers:

1. `(Bob, Alice, Bob, T, [1-2], [1000-2000]) -> APPROVED`
2. `(Bob, Alice, All, T, [1-2], [1-2000]) -> APPROVED`
3. `(Bob, Alice, All, T, [1-2], [2001-Max]) -> APPROVED`

Let's say each approval has no amount restriction but does not override the user level incoming / outgoing approvals.

In this scenario, let's say the default "first match" approach is used:

1. The first approved transfer `(Bob, Alice, Bob, T, [1-2], [1000-2000])` matches the transfer request partially, but it only covers **ownershipTimes** from 1000 to 2000. We deduct this overlap but still have a lot remaining to be approved for (\[1-999], \[2001-Max]).
2. The second approved transfer `(Bob, Alice, All, T, [1-2], [1-2000]),` again partially matches, and we deduct. This partial match only handles 1-999 because we handled 1000-2000 already. Note that this approval says it can be initiated by "All" instead of Bob, but Bob's address is within the "All" list.
3. The third approved transfer `(Bob, Alice, All, T, [1-2], [1001-Max])` covers the rest. This transfer is approved on the collection level because the entire transfer was handled.

If Bob was requesting badge ID 3 to be transferred as well, it would fail because badge ID 3 is unhandled by all defined approvals (and disallowed by default if unhandled).

**Outgoing Approvals**

Because we did not override the user level approvals, we need to check that Bob approved this transfer in his approvals.

The process above would then be repeated for Bob's outgoing approvals. In this case, Bob is the initiator, so we automatically add an unlimited approval by default (see below).

**Incoming Approvals**

Likewise, we also need to check Alice's incoming approvals using the same process.

**Satisfied?**

If all levels are satisfied, the transfer is approved, and we deduct/increment the used approvals where necessary.

**Extending the Example: Prioritized Approvals**

Let's say that Bob only wants to use the second and third approvals from the collection level but not the first. By default, a first-match policy is applied, so it would by default use the first one as shown above.

When initiating the transfer (MsgTransferBadges), Bob can set **prioritizedApprovals** to be the second and third collection approvals. These would then be checked first, followed by the first approval.

If Bob additionally sets **onlyCheckPrioritizedApprovals** = true, we only check the ones specified in **prioritizedApprovals**.

### Example 2 - Collection

This would define a collection where badges 1-100 can be transferred from the Mint address (according to the first approval). Once transferred out of the Mint address, they can be transferred freely, thus making the collection transferable.

Note how the **fromList** of each approval are non-overlapping, so any transfer will only match to one of the two approvals (if either). The first approval is restricted to transfers from the Mint address whereas the second is all EXCEPT the Mint address.

```json
 "collectionApprovals": [
    {
      "fromListId": "Mint",
      "toListId": "AllWithMint",
      "initiatedByListId": "AllWithMint",
      "transferTimes": [
        {
          "start": "1691931600000",
          "end": "1723554000000"
        }
      ],
      "ownershipTimes": [
        {
          "start": "1",
          "end": "18446744073709551615"
        }
      ],
      "badgeIds": [
        {
          "start": "1",
          "end": "100"
        }
      ],
      "approvalId": "claim-from-mint-address",
      ... //other criteria (including the IMPORTANT overrideFromOutgoingApprovals = true since we are dealing with transfers from the Mint address)
    },
    {
      "fromListId": "AllWithoutMint",
      "toListId": "AllWithoutMint",
      "initiatedByListId": "AllWithoutMint",
      "badgeIds": [
        {
          "start": "1",
          "end": "100"
        }
      ],
      "ownershipTimes": [
        {
          "start": "1",
          "end": "18446744073709551615"
        }
      ],
      "transferTimes": [
        {
          "start": "1",
          "end": "18446744073709551615"
        }
      ],
      "approvalId": "transferable"
    }
  ],
```

### Example 3 - Outgoing Approvals

This would set approve Charlie to send badges to Bob on this user's behalf.

```json
"outgoingApprovals": [
  {
    "toListId": "Bob",
    "initiatedByListId": "Charlie",
    "transferTimes": [
      {
        "start": "1",
        "end": "18446744073709551615"
      }
    ],
    "badgeIds": [
      {
        "start": "1",
        "end": "18446744073709551615"
      }
    ],
    "ownershipTimes": [
      {
        "start": "1",
        "end": "18446744073709551615"
      }
    ],
    "approvalId": "test",
    //see next page
    "approvalCriteria": //define approval criteria (how much? challenges? etc here)
  }
]
```

### Example 4 - Incoming Approvals

This would set approve this user to receive any transfer from Bob.

```json
"incomingApprovals": [
  {
    "fromListId": "Bob",
    "initiatedByListId": "AllWithMint",
    "transferTimes": [
      {
        "start": "1",
        "end": "18446744073709551615"
      }
    ],
    "badgeIds": [
      {
        "start": "1",
        "end": "18446744073709551615"
      }
    ],
    "ownershipTimes": [
      {
        "start": "1",
        "end": "18446744073709551615"
      }
    ],
    "approvalId": "test",
  }
]
```
