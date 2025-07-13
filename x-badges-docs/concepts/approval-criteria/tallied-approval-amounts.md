# Tallied Approval Amounts

Pre-Readings (IMPORTANT): [Approval Trackers](approval-trackers.md)

```typescript
export interface ApprovalCriteria<T extends NumberType> {
  ...
  approvalAmounts?: ApprovalAmounts<T>;
  ...
}
```

**Approval Amounts**

Approval amounts (**approvalAmounts**) allow you to specify the threshold amount that can be transferred for this approval. This is similar to other interfaces (such as approvals for ERC721), except we use an increment + threshold system as opposed to a decrement + greater than 0 system. This is facilitated via the concept of [Approval Trackers](approval-trackers.md).

The amounts approved are scoped to the **badgeIds** and **ownershipTimes** defined by the base approval (see transferability page). Also, note that the to addresses are bounded to the addresses in the **toList,** from addresses from the **fromList**, and initiated by addresses from the **initiatedByList**.

We define four levels (**trackerType** = "overall", "to", "from", "initiatedBy") that you can specify for approval amounts as seen below. You can define multiple if desired, and to be approved, the transfer must satisfy all.

* **Overall**: Overall will increment a universal, cumulative approval tracker for all transfers that match this approval, regardless of who sends, receives, or initiates them.
* **Per To Address**: If you specify an approval amount per to address, we will create unique cumulative trackers for every unique "to" address.
* **Per From Address**: Creates unique cumulative tallies for every unique "from" address.
* **Per Initiated By Address**: Creates unique cumulative tallies for every unique "initiatedBy" address.

If the amount set is nil value or "0", this means there is no limit (no amount restrictions).

**Example**

```json
"collectionApprovals": [
    {
      "fromListId": "Bob",
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
      "approvalId": "uniqueID",
      "version": "1",
      
      "approvalCriteria": {
        "approvalAmounts": {
           "overallApprovalAmount": "1000", //overall limit of x1000
           "perFromAddressApprovalAmount": "0", //no limit
           "perToAddressApprovalAmount": "0",
           "perInitiatedByAddressApprovalAmount": "10", //limit of x10 per initiator
           "amountTrackerId": "uniqueID",
           "resetTimeIntervals": {
              "startTime": "0",
              "intervalLength": "0"
            }
        },
        ...
      }
      ...
    }
  
```

Let's say we have the **approvalAmounts** defined above and Alice initiates a transfer of x10 from Bob. There are two separate trackers that get incremented here.

\#1) Tracker with the following ID `1-collection- -approvalId-uniqueID-overall-` gets incremented to x10 out of 1000. Any subsequent transfers (say from Charlie) will also increment this overall universal tracker as well.

\#2) Tracker with ID `1-collection- -approvalId-uniqueID-initiatedBy-alice`gets incremented to x10 out of 10 used. Alice has now fully used up her threshold for this tracker. This tracker is only incremented when Alice initiates the transfer. If Charlie initiates a transfer, his unique initiatedBy tracker will get incremented which is separate from Alice's.

Since there was an unlimited amount approved for the "to" and "from" trackers, we do not increment anything for those trackers (as-needed basis).

**Resets + ID Changes**

Let's say we update the **amountTrackerId** to "uniqueID2" from "uniqueID". This makes all tracker IDs different, and thus, all tallies will start from scratch.

`1-collection- -approvalId-uniqueID-initiatedBy-alice` ->

`1-collection- -approvalId-uniqueID2-initiatedBy-alice`

If in the future, you change back to "uniqueID", the starting point will be the previous tally. Using the examples above, x10/10 used for Alice's initiated by tracker.
