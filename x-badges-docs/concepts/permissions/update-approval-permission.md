# Update Approval Permission

Pre-Readings: [Transferability](../balances-transfers/transferability-approvals.md) and [Approval Criteria](../balances-transfers/approval-criteria/)

The ApprovalPermissions refer to the UPDATABILITY of the currently set approvals / transferability. These can be leveraged to freeze specific transferability / approvals in order to give users more confidence that they cannot be changed in the future. Note this refers to the updatability of them and has no bearing on what they are currently set to.

For what transfer combinations (see [Representing Transfers](../balances-transfers/transferability-approvals.md)), can I create / delete / update approvals?&#x20;

```json
"userPermissions": {
    "canUpdateIncomingApprovals": [...],
    "canUpdateOutgoingApprovals": [...],
    ...
}
```

```json
"collectionPermissions": {
    "canUpdateCollectionApprovals": [...]
    ...
}
```

The **canUpdateIncomingApprovals** and **canUpdateOutgoingApprovals** follow the same interface as **canUpdateCollectionApprovals** minus automatically populating the user's address for to / from for incoming / outgoing, respectively. We only explain the collection approval permission to avoid repetition.

```typescript
export interface CollectionApprovalPermission<T extends NumberType> {
  fromListId: string;
  toListId: string;
  initiatedByListId: string;
  transferTimes: UintRange<T>[];
  badgeIds: UintRange<T>[];
  ownershipTimes: UintRange<T>[];
  approvalId: string
  
  permanentlyPermittedTimes: UintRange<T>[];
  permanentlyForbiddenTimes: UintRange<T>[];
}
```

Ex: I can/cannot update the approvals for the transfer combinations ("All", "All", "All", 1-100, 1-10, 1-10, "All",  "All",  "All") tuple.

## ID Shorthands

IDs are used for locking specific approvals. This is because sometimes it may not be sufficient to just lock a specific (from, to, initiator, time, badge IDs, ownershipTimes) combination because multiple approvals could match to it.

To specify IDs, you can use the "All" reserved ID to represent all IDs, or you can use other shorthand methods such as "!xyz" to denote all IDs but xyz. These shorthands are the same as reserved lists, so we refer you[ there for more info](../../core-concepts/address-lists-lists.md). Just replace the addresses with the IDs.

## Break Down Logic

Because of the way breakdown logic is performed, the following is allowed.

Permission: Badge IDs 2-10 are locked as forbidden

Before: 1-10 -> Criteria ABC

After: 1 -> Criteria XYZ, 2-10 -> Criteria ABC

## Expected Behavior vs Non-Updatability

Approval permission updates are slightly tricky because even though an approval may be non-updatable according to the permissions set, its expected behavior may change due to how approvals are designed (i.e. using trackers).

* For example, lets say we want to freeze IDs 501-1000 and have an incrementing mint of x1 of ID 1, x1 of ID 2, up to ID 1000. If we simply freeze the IDs 501-1000, the approval could still be deleted for IDs 1-500, and the increment number (tracker) will then never reach 501 because it will be out of bounds every time. Thus, expected behavior for 501-1000 changes even though it is frozen.

### Definitions

To explain things easier, let's start with some definitions:

#### **Approval Tuple**

We define an approval tuple as a set of values (**from, to, initiated by, badgeIds, transferTimes, ownershipTimes, approvalId**). The tuple for a specific approval that is currently set will consist of all values for that approval.

Note that to match, all N criteria must match. If one doesn't, it isn't a match. For example, in the example below, badge ID 1 will never match to the tuple.

<pre class="language-json"><code class="lang-json"><strong>{
</strong>  "fromListId": "AllWithMint",
  "toListId": "AllWithMint",
  "initiatedByListId": "AllWithMint",
  "badgeIds":  [{ "start": "2", "end": "10" }],
  "transferTimes":  [{ "start": "1", "end": "18446744073709551615" }],
  "ownershipTimes":  [{ "start": "1", "end": "18446744073709551615" }],
  "approvalId": "All"
}
</code></pre>

#### **Non-Updatable**

For a given approval tuple, it is considered non-updatable according to the permissions if all possible combinations of the entire tuple are **permanently forbidden** from being updated in the permissions.&#x20;

<pre class="language-json"><code class="lang-json"><strong>{
</strong>  "fromListId": "AllWithMint",
  "toListId": "AllWithMint",
  "initiatedByListId": "AllWithMint",
  "badgeIds":  [{ "start": "2", "end": "10" }],
  "transferTimes":  [{ "start": "1", "end": "18446744073709551615" }],
  "ownershipTimes":  [{ "start": "1", "end": "18446744073709551615" }],
  "approvalId": "All", //forbids approval "xyz" from being updated
  
  "permanentlyPermittedTimes": [],
  "permanentlyForbiddenTimes": [{ "start": "1", "end": "18446744073709551615" }]
}
</code></pre>

Note that non-updatability is scoped to the tuple itself.

**Brute Forced**

Commonly, you will make some values non-updatable by specifying some criteria (e.g. IDs 2-10) and setting everything else to all possible values. For example, the permission above does this for badge IDs 2-10. We refer to this as brute forcing (i.e. above badge IDs 2-10 are brute forced but badge IDs 2-11 are not). In other words, for some criteria, all possible combinations of that criteria are COMPLETELY forbidden and non-updatable.

#### **Expected Behavior**

As explained above, expected behavior not only encompasses non-updatability, but it also makes sure that nothing any other approval or update can do can affect the expected behavior of this approval. This designation is especially important.

### Freezing Specific Approval Tuples

With the way trackers work, it is important to handle approval permissions correctly to protect against break-down attacks to ensure expected behavior.

Below, we will walk through the process of making a specific approval tuple non-updatable AND keeping its expected behavior.&#x20;

**Specific Approval ID**

If you want to do this for a specific approval that is set, the approval tuple should consist of the specific values for that specific approval. Because the **approvalId** is unique and included in the tuple, you know there are no other approvals that overlap.

Thus, you can simply brute force this tuple in the permissions and call it a day.

**Tuples with Overlaps**

For tuples which may span multiple approvals, the algorithm is essentially the following:

1. Forbid updates for the exact tuple values you want to freeze in permissions
2. To keep expected behavior, you need to forbid updates for all specific approvals that are currently set and overlap (even partially) with the values you are trying to freeze.
   * Ex: If you are trying to freeze badge IDs 1-10, you should also entirely freeze the approval abc123 which is for badge IDs 1-100.
   * Technically, this only needs to be done for approvals that are set and can affect the behavior of the tuple values, but to be safe, we recommend freezing all overlapping approvals or restructuring so they do not overlap. It is difficult to manage when part of an approval is frozen.

**Example**

Let's say you want to forbid ever updating the transferability for badges 1-10, and you have the following approvals currently set (for badges 1-100).

```json
"collectionApprovals": [
    {
      "fromListId": "Mint",
      "toListId": "All",
      "initiatedByListId": "All",
      "transferTimes": [
        {
          "start": "1691978400000",
          "end": "1723514400000"
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
      "approvalId": "abc",
      "approvalCriteria": {
        .... //uses trackers
      }
    }
  ]
```

Step 1: Brute force IDs 1-10 in the permissions to be forbidden.

<pre class="language-json"><code class="lang-json"><strong>{
</strong>  "fromListId": "All",
  "toListId": "All",
  "initiatedByListId": "All",
  "badgeIds":  [{ "start": "1", "end": "10" }],
  "transferTimes":  [{ "start": "1", "end": "18446744073709551615" }],
  "ownershipTimes":  [{ "start": "1", "end": "18446744073709551615" }],
  "approvalId": "All",
  "permanentlyPermittedTimes": [],
  "permanentlyForbiddenTimes": [{ "start": "1", "end": "18446744073709551615" }]
}
</code></pre>

Step 2: Find all matching approvals. In this case, we only have one and it matches because it uses overlaps since it uses IDs 1-100.

```json
{
    "fromListId": "Mint",
    "toListId": "All",
    "initiatedByListId": "All",
    "transferTimes": [
      {
        "start": "1691978400000",
        "end": "1723514400000"
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
    "approvalId": "abc",
    "approvalCriteria": {
      ....
    }
  }
```

In this particular case, the approval is updatable (IDs 11-100 are) and is not already frozen. Thus, expected behavior of badges 1-10 may not be guaranteed.

We need to handle this, which we can do by adding another permission brute forcing approval "abc".

```json
{
    "fromListId": "All",
    "toListId": "All",
    "initiatedByListId": "All",
    "transferTimes": [
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
    "badgeIds": [
      {
        "start": "1",
        "end": "18446744073709551615"
      }
    ],
    "approvalId": "abc",    
     
    "permanentlyPermittedTimes": [],
    "permanentlyForbiddenTimes": [{ "start": "1", "end": "18446744073709551615" }]
  }
```
